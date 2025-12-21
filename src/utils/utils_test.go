package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRandInt(t *testing.T) {
	for i := 0; i < 5; i++ {
		n := RandInt(0, 2)
		if n < 0 || n > 2 {
			t.Fatalf("out of range: %d", n)
		}
	}
}

func TestLoadCurrentLog(t *testing.T) {
	tmp := t.TempDir()
	old := LogPath
	LogPath = filepath.Join(tmp, "play.log")
	t.Cleanup(func() { LogPath = old })

	content := "v0.0.0\n\nfirst\n\nsecond"
	if err := os.WriteFile(LogPath, []byte(content), 0o666); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got := LoadCurrentLog()
	if strings.TrimSpace(got) != "second" {
		t.Fatalf("expected last block 'second', got %q", got)
	}
}

func TestChangeQuickMode_Toggles(t *testing.T) {
	tmp := t.TempDir()
	oldCfgPath := ConfigFilePath
	ConfigFilePath = filepath.Join(tmp, "config.json")
	t.Cleanup(func() { ConfigFilePath = oldCfgPath })

	// Ensure starting state and then toggle.
	Config.QuickMode = false
	ChangeQuickMode()
	if !Config.QuickMode {
		t.Fatalf("expected QuickMode=true after toggle")
	}
	ChangeQuickMode()
	if Config.QuickMode {
		t.Fatalf("expected QuickMode=false after second toggle")
	}
}
