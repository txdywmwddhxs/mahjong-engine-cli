//go:build integration

package blackbox_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLI_OneRound_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	root, err := repoRoot()
	if err != nil {
		t.Fatalf("repoRoot: %v", err)
	}

	tmp := t.TempDir()
	work := filepath.Join(tmp, "work")
	if err := os.MkdirAll(filepath.Join(work, "config"), 0o755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	// Minimal CHANGELOG marker so runtime root detection works.
	if err := os.WriteFile(filepath.Join(work, "config", "CHANGELOG.txt"), []byte("v0.0.0\n"), 0o666); err != nil {
		t.Fatalf("write changelog: %v", err)
	}

	bin := filepath.Join(tmp, "play")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	// Build the CLI binary.
	build := exec.Command("go", "build", "-trimpath", "-o", bin, filepath.Join(root, "src"))
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}

	// Run one game, trigger the win-script and then exit at Continue?.
	// NOTE: this relies on deterministic randomness to avoid extra prompts.
	in := bytes.NewBufferString("WHOSYOURDADDY\nN\n")
	cmd := exec.Command(bin, "--quick-mode")
	cmd.Dir = work
	cmd.Env = append(os.Environ(),
		"MAHJONG_ROOT="+work,
		"MAHJONG_SEED=1",
	)
	cmd.Stdin = in
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	out := stdout.String() + stderr.String()
	if err != nil {
		t.Fatalf("cli run failed: %v\n%s", err, out)
	}

	// Minimal assertions: program started, printed a continue prompt, and printed a win message.
	if !strings.Contains(out, "Continue") && !strings.Contains(out, "继续") {
		t.Fatalf("expected Continue prompt in output, got:\n%s", out)
	}
	if !strings.Contains(out, "You win") && !strings.Contains(out, "你赢了") {
		t.Fatalf("expected win message in output, got:\n%s", out)
	}
}

func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
