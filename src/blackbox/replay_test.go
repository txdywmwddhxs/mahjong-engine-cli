//go:build integration

package blackbox_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

// Exercise actual files and process startup. Engine-level replay tests cover
// ordinary games, melds and waiting; this test focuses on CLI isolation.
func TestCLI_Replay_RecordVerifyReadOnly(t *testing.T) {
	root, err := repoRoot()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	work := filepath.Join(tmp, "work")
	if err = os.MkdirAll(filepath.Join(work, "config"), 0755); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(work, "VERSION"), []byte("v0.0.0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(work, "config", "config.json"), []byte(`{"version":"v0.0.0","lang":"English"}`), 0644); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(tmp, "play")
	build := exec.Command("go", "build", "-o", bin, filepath.Join(root, "src"))
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	run := func(runtimeRoot, input string, args ...string) (string, error) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, bin, args...)
		cmd.Dir = tmp
		cmd.Env = append(os.Environ(), "MAHJONG_ROOT="+runtimeRoot, "MAHJONG_SEED=1")
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("command hung: %v\n%s", args, out)
		}
		return string(out), err
	}
	path := filepath.Join(work, "game.json")
	out, err := run(work, "WHOSYOURDADDY\nY\nWHOSYOURDADDY\nN\n", "--quick-mode", "--record", path)
	if err != nil {
		t.Fatalf("record: %v\n%s", err, out)
	}
	first, err := replay.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	secondPath := filepath.Join(work, "game.2.json")
	second, err := replay.Load(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	if first.Seed == second.Seed || first.GameID == second.GameID || reflect.DeepEqual(first.Initial, second.Initial) {
		t.Fatal("rounds share identity or setup")
	}
	logPath := filepath.Join(work, "log", "play.log")
	block, err := first.LogBlock()
	if err != nil {
		t.Fatal(err)
	}
	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(logData), block) {
		t.Fatal("recording was not embedded in readable log")
	}
	before := snapshotTree(t, work)
	for _, recording := range []string{path, secondPath} {
		if out, err = run(work, "", "--replay", recording, "--verify", "--quick-mode"); err != nil || !strings.Contains(out, "Replay verified") {
			t.Fatalf("verify: %v\n%s", err, out)
		}
	}
	if out, err = run(work, "n\np\ng 2\nt 1\nq\n", "--replay", path, "--step"); err != nil || !strings.Contains(out, "Frame 2/") {
		t.Fatalf("step: %v\n%s", err, out)
	}
	if after := snapshotTree(t, work); !reflect.DeepEqual(before, after) {
		t.Fatal("replay modified runtime files or directories")
	}
	if out, err = run(work, block, "--replay", "-", "--verify"); err != nil || !strings.Contains(out, first.GameID) {
		t.Fatalf("clipboard replay: %v\n%s", err, out)
	}
	if out, err = run(work, "", "--replay", logPath, "--verify"); err != nil || !strings.Contains(out, second.GameID) {
		t.Fatalf("log replay: %v\n%s", err, out)
	}
	if out, err = run(work, "", "--replay", logPath, "--game", first.GameID, "--verify"); err != nil || !strings.Contains(out, first.GameID) {
		t.Fatalf("selected game replay: %v\n%s", err, out)
	}
	if after := snapshotTree(t, work); !reflect.DeepEqual(before, after) {
		t.Fatal("log/clipboard replay modified runtime files")
	}
	missingRoot := filepath.Join(tmp, "must-not-be-created")
	if out, err = run(missingRoot, "", "--replay", path, "--verify"); err != nil {
		t.Fatalf("isolated verify: %v\n%s", err, out)
	}
	if _, err = os.Stat(missingRoot); !os.IsNotExist(err) {
		t.Fatal("replay created runtime root")
	}
	// Changing the recorded result must produce a useful, immediate failure.
	first.Frames[len(first.Frames)-1].Settlement.Score++
	tampered := filepath.Join(tmp, "tampered.json")
	if err = first.Save(tampered); err != nil {
		t.Fatal(err)
	}
	if out, err = run(work, "", "--replay", tampered, "--verify"); err == nil || !strings.Contains(out, "settlement.score") || !strings.Contains(out, "frame") {
		t.Fatalf("tamper was not diagnosed: %v\n%s", err, out)
	}
	// EOF cancels the current round and must neither auto-play nor add a result.
	configPath := filepath.Join(work, "config", "config.json")
	data, _ := os.ReadFile(configPath)
	var oldConfig, newConfig utils.ConfigType
	if err = json.Unmarshal(data, &oldConfig); err != nil {
		t.Fatal(err)
	}
	if out, err = run(work, "", "--quick-mode"); err != nil {
		t.Fatalf("EOF: %v\n%s", err, out)
	}
	data, _ = os.ReadFile(configPath)
	if err = json.Unmarshal(data, &newConfig); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(oldConfig, newConfig) {
		t.Fatal("EOF changed stats or config")
	}
	// EOF is marked as incomplete, instead of accidentally selecting the prior game.
	if out, err = run(work, "", "--replay", logPath, "--verify"); err == nil || !strings.Contains(out, "incomplete") {
		t.Fatalf("unfinished latest game was silently skipped: %v\n%s", err, out)
	}
	// Normal play records to the log without needing an explicit --record option.
	if out, err = run(work, "WHOSYOURDADDY\nN\n", "--quick-mode"); err != nil {
		t.Fatalf("automatic recording: %v\n%s", err, out)
	}
	if out, err = run(work, "", "--replay", logPath, "--verify"); err != nil || !strings.Contains(out, "Replay verified") {
		t.Fatalf("automatic log replay: %v\n%s", err, out)
	}

}

func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, _ := filepath.Rel(root, path)
		if entry.IsDir() {
			files[relative] = "directory"
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[relative] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
