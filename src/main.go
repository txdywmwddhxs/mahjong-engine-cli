package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	log "github.com/txdywmwddhxs/mahjong-engine-cli/src/clog"
	l "github.com/txdywmwddhxs/mahjong-engine-cli/src/language"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/player"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/single"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

var (
	lang   = utils.Config.Lang
	logger *log.Log
	s      = single.Single()
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
		logger = log.NewLogger(utils.LogPath, lang)
		func() {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					logger.DebugS(fmt.Sprintf("%s: %v, stack: %s", log.Trans(l.MeetError), err, string(stack)), 0)
					if logger != nil {
						logger.Close()
					}
				}
			}()

			p := player.NewPlayer(logger, lang)
			res, score, counter := p.Main()
			utils.UpdateConfig(res, score)
			logger.InfoS(fmt.Sprintf("%s: %d", log.Trans(l.TotalScore), utils.Config.Score), counter)
			cs, _ := json.Marshal(utils.Config)
			logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.CurrentInfo), string(cs)), counter+1)
			logger.Close()
		}()
		if !continue_() {
			return
		}
	}
}

func continue_() bool {
	fmt.Printf("\n%s", log.Trans(l.Continue))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.ToUpper(scanner.Text())
	switch input {
	case "":
		return true
	case "Y":
		return true
	case "YES":
		return true
	case "S":
		config, _ := json.Marshal(utils.Config)
		fmt.Println(string(config))
		return continue_()
	case "L":
		content := utils.LoadCurrentLog()
		fmt.Println(content)
		return continue_()
	default:
		return false
	}
}
