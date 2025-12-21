package player

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestPlayer_CheckIsWaiting(t *testing.T) {
	lang := "English"
	logPath := filepath.Join(os.TempDir(), "mahjong-engine-cli-test.log")
	logger, err := logging.NewFileLogger(logPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = logger.Close() }()

	console := ui.NewConsole(strings.NewReader(""), io.Discard)
	tr := language.New(utils.Lang(lang))
	p := NewPlayer(logger, console, tr, utils.Lang(lang))
	//p.PlayerCards = &card.Cards{"B3", "B3", "W2", "W3", "W5", "W6", "W7", "W7", "W8", "W9", "F", "F", "F"}
	p.PlayerCards = &card.Cards{"B1", "B1", "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "B9", "B9"}
	res, c := p.CheckIsWaiting()
	fmt.Println(res, *c)
}
