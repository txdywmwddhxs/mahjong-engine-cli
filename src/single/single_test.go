package single

import (
	"path/filepath"
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestScriptPidLifecycle(t *testing.T) {
	tmp := t.TempDir()
	old := utils.PIDPath
	utils.PIDPath = filepath.Join(tmp, "play.pid")
	t.Cleanup(func() { utils.PIDPath = old })

	s := Single()
	if s.IsRunning() {
		t.Fatalf("expected not running")
	}
	s.CreatePidFile()
	if !s.IsRunning() {
		t.Fatalf("expected running after pid file creation")
	}
	s.RemovePidFile()
	if s.IsRunning() {
		t.Fatalf("expected not running after removal")
	}
}
