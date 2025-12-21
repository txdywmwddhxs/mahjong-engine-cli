package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileLogger_WritesAndCloses(t *testing.T) {
	tmp := t.TempDir()
	fp := filepath.Join(tmp, "play.log")

	l, err := NewFileLogger(fp)
	if err != nil {
		t.Fatalf("NewFileLogger: %v", err)
	}

	l.Raw("HEADER\n")
	l.Info("hello", 1)
	l.Debug("dbg", 2)
	l.Error("boom", 3)

	if err := l.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	b, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		"HEADER\n",
		"INFO[1]: hello\n",
		"DEBUG[2]: dbg\n",
		"ERROR[3]: boom\n",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in log:\n%s", want, s)
		}
	}
}
