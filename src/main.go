package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/player"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/single"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

var (
	lang                   = utils.Config.Lang
	logger  logging.Logger = logging.NopLogger()
	console ui.UI          = ui.NewConsole(os.Stdin, os.Stdout)
	t                      = language.New(lang)
	s                      = single.Single()
)

func main() {
	if err := run(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	quick := flag.Bool("quick-mode", false, "Run without presentation delays")
	record := flag.String("record", "", "Record rounds to JSON (later rounds use .2, .3, ...)")
	replayPath := flag.String("replay", "", "View or verify JSON, a log, or stdin (-) without changing stats")
	gameID := flag.String("game", "", "Select a game ID from a replay log (default: latest round)")
	verify := flag.Bool("verify", false, "Strictly verify a replay with the current engine")
	step := flag.Bool("step", false, "Interactively view replay snapshots")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flag.Args())
	}
	replaySelected := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "replay" {
			replaySelected = true
		}
	})
	if replaySelected && *replayPath == "" {
		return fmt.Errorf("--replay requires a non-empty file path")
	}
	if *replayPath != "" {
		if *record != "" || (*verify && *step) {
			return fmt.Errorf("--replay cannot combine --record, or --verify with --step")
		}
		if *replayPath == "-" && *step {
			return fmt.Errorf("--step needs a replay file; stdin is used to read the pasted replay")
		}
		var file *replay.File
		var err error
		if *replayPath == "-" {
			file, err = replay.ReadSource(os.Stdin, *gameID)
		} else {
			file, err = replay.LoadGame(*replayPath, *gameID)
		}
		if err != nil {
			return err
		}
		if *verify {
			if err = player.VerifyReplay(file); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Replay verified: %s (%d frames)\n", file.GameID, len(file.Frames))
			return nil
		}
		return replay.View(file, os.Stdin, os.Stdout, *step)
	}
	if *verify || *step || *gameID != "" {
		return fmt.Errorf("--verify, --step and --game require --replay FILE")
	}
	if *quick {
		utils.Config.QuickMode = true
	}
	if s.IsRunning() {
		return fmt.Errorf("another game is running")
	}
	if err := utils.PrepareRuntime(); err != nil {
		return err
	}
	return loopRun(*record)
}

func roundRecordingPath(base string, round int) string {
	if round == 1 {
		return base
	}
	ext := filepath.Ext(base)
	return fmt.Sprintf("%s.%d%s", strings.TrimSuffix(base, ext), round, ext)
}

func loopRun(recordPath string) error {
	utils.ChangeToCurrentVersion()
	for round := 1; ; round++ {
		path := ""
		if recordPath != "" {
			path = roundRecordingPath(recordPath, round)
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("recording already exists: %s", path)
			} else if !os.IsNotExist(err) {
				return err
			}
			if info, err := os.Stat(filepath.Dir(path)); err != nil {
				return err
			} else if !info.IsDir() {
				return fmt.Errorf("recording parent is not a directory")
			}
		}
		fl, err := logging.NewFileLogger(utils.LogPath)
		if err != nil {
			console.Info(fmt.Sprintf("%s: %v", t.T(language.KeyMeetError), err))
			logger = logging.NopLogger()
		} else {
			logger = fl
		}
		err = runRound(path)
		_ = logger.Close()
		if err != nil {
			return err
		}
		if !continue_() {
			return nil
		}
	}
}

func runRound(recordPath string) (err error) {
	var recording *replay.Session
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("game failed: %v", recovered)
			logger.Debug(fmt.Sprintf("%v\n%s", err, debug.Stack()), 0)
		}
		if err != nil && recording != nil && !recording.File.Complete {
			logger.Raw(fmt.Sprintf("%s%s %q\n", replay.LogIncompletePrefix, recording.File.GameID, err.Error()))
		}
	}()
	seed := utils.NextRoundSeed()
	logger.Raw(fmt.Sprintf("\n%s\n", t.T(language.KeyGameBegin)))
	logger.Raw(fmt.Sprintf("%s: %s\n", t.T(language.KeyGameVersion), utils.Version))
	logger.Raw(fmt.Sprintf("%s: %d\n", t.T(language.KeyGameSeed), seed))
	p := player.NewPlayerWithOptions(logger, console, t, lang, player.Options{Seed: &seed})
	recording = p.StartRecording()
	logger.Raw(replay.LogStartPrefix + recording.File.GameID + "\n")
	res, amount, counter := p.Main()
	if err = p.Err(); err != nil {
		return err
	}
	if err = recording.Finish(); err != nil {
		return err
	}
	block, blockErr := recording.File.LogBlock()
	if blockErr != nil {
		return blockErr
	}
	if logErr := logger.WriteBlock(block); logErr != nil {
		console.Info("Could not save replay to log: " + logErr.Error())
	}
	if recordPath != "" {
		if err = recording.File.Save(recordPath); err != nil {
			return err
		}
		console.Info("Replay saved: " + recordPath)
	}

	utils.UpdateConfig(res, amount)
	msg := fmt.Sprintf("%s: %d", t.T(language.KeyTotalScore), utils.Config.Score)
	console.Info(msg)
	logger.Info(msg, counter)
	config, _ := json.Marshal(utils.Config)
	logger.Debug(fmt.Sprintf("%s: %s", t.T(language.KeyCurrentInfo), config), counter+1)
	logger.Raw(fmt.Sprintf("%s\n", t.T(language.KeyGameEnd)))
	return nil
}

func continue_() bool {
	input, err := console.PromptPlain("\n" + t.T(language.KeyContinue))
	if err != nil {
		return false
	}
	input = strings.ToUpper(strings.TrimSpace(input))
	switch input {
	case "":
		return true
	case "Y":
		return true
	case "YES":
		return true
	case "S":
		config, _ := json.Marshal(utils.Config)
		console.Plainln(string(config))
		return continue_()
	case "L":
		content := utils.LoadCurrentLog()
		console.Plainln(content)
		return continue_()
	default:
		return false
	}
}
