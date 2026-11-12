package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestContinue_YesAndNo(tt *testing.T) {
	oldConsole := console
	oldTr := t
	console = ui.NewConsole(strings.NewReader("Y\nNO\n"), io.Discard)
	t = language.New(utils.Chinese)
	defer func() {
		console = oldConsole
		t = oldTr
	}()

	if !continue_() {
		tt.Fatalf("expected continue on Y")
	}
	if continue_() {
		tt.Fatalf("expected stop on NO")
	}
}

func TestContinue_ShowConfigAndLog(tt *testing.T) {
	tmp := tt.TempDir()
	oldLogPath := utils.LogPath
	utils.LogPath = filepath.Join(tmp, "play.log")
	defer func() { utils.LogPath = oldLogPath }()

	_ = os.WriteFile(utils.LogPath, []byte("v0.0.0\n\nhello-log"), 0o666)

	oldConsole := console
	oldTr := t
	console = ui.NewConsole(strings.NewReader("S\nL\nN\n"), io.Discard)
	t = language.New(utils.English)
	defer func() {
		console = oldConsole
		t = oldTr
	}()

	if continue_() {
		tt.Fatalf("expected stop after S/L then N")
	}
}
