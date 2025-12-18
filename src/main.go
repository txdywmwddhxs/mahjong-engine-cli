package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	l "github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/player"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/single"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

var (
	lang                   = utils.Config.Lang
	logger  logging.Logger = logging.NopLogger()
	console ui.UI          = ui.NewConsole(os.Stdin, os.Stdout)
	s                      = single.Single()
)

func main() {
	if s.IsRunning() {
		return
	}
	//s.CreatePidFile()
	//defer s.RemovePidFile()
	loopRun()
}

func loopRun() {
	utils.ChangeToCurrentVersion()
	for {
		fl, err := logging.NewFileLogger(utils.LogPath)
		if err != nil {
			// Keep the game playable even if logging fails.
			console.Info(fmt.Sprintf("%s: %v", l.MeetError(), err))
			logger = logging.NopLogger()
		} else {
			logger = fl
		}

		logger.Raw(fmt.Sprintf("\n%s\n", l.GameBegin()))
		logger.Raw(fmt.Sprintf("%s: %s\n", l.GameVersion(), utils.Config.Version))
		func() {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					logger.Debug(fmt.Sprintf("%s: %v, stack: %s", l.MeetError(), err, string(stack)), 0)
					if logger != nil {
						_ = logger.Close()
					}
				}
			}()

			p := player.NewPlayer(logger, console, lang)
			res, score, counter := p.Main()
			utils.UpdateConfig(res, score)
			msg := fmt.Sprintf("%s: %d", l.TotalScore(), utils.Config.Score)
			console.Info(msg)
			logger.Info(msg, counter)
			cs, _ := json.Marshal(utils.Config)
			logger.Debug(fmt.Sprintf("%s: %s", l.CurrentInfo(), string(cs)), counter+1)
			logger.Raw(fmt.Sprintf("%s\n", l.GameEnd()))
			_ = logger.Close()
		}()
		if !continue_() {
			return
		}
	}
}

func continue_() bool {
	input, _ := console.PromptPlain("\n" + l.Continue())
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
