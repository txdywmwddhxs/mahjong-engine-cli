package player

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	log "github.com/txdywmwddhxs/mahjong-engine-cli/src/clog"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestPlayer_CheckIsWaiting(t *testing.T) {
	lang := "English"
	logPath := filepath.Join(os.TempDir(), "mahjong-engine-cli-test.log")
	logger := log.NewLogger(
		logPath,
		utils.Lang(lang))
	p := NewPlayer(logger, utils.Lang(lang))
	//p.PlayerCards = &card.Cards{"B3", "B3", "W2", "W3", "W5", "W6", "W7", "W7", "W8", "W9", "F", "F", "F"}
	p.PlayerCards = &card.Cards{"B1", "B1", "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "B9", "B9"}
	res, c := p.CheckIsWaiting()
	fmt.Println(res, *c)
}
