//go:build integration

package blackbox_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

// TestReleasePackages tests the actual packaging script and its native artifacts,
// not a separately built development executable. Enabled by check.sh release.
func TestReleasePackages(t *testing.T) {
	if os.Getenv("MAHJONG_RELEASE_CHECK") != "1" {
		t.Skip("run scripts/check.sh release")
	}
	root, err := repoRoot()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "VERSION"))
	if err != nil {
		t.Fatal(err)
	}
	version := strings.TrimSpace(string(data))
	versionPattern := regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
	if !versionPattern.MatchString(version) {
		t.Fatalf("invalid release VERSION %q", version)
	}
	changelog, err := os.ReadFile(filepath.Join(root, "config", "CHANGELOG.txt"))
	if err != nil {
		t.Fatal(err)
	}
	latest := ""
	for _, line := range strings.Split(string(changelog), "\n") {
		if line = strings.TrimSpace(line); versionPattern.MatchString(line) {
			latest = line
		}
	}
	if latest != version {
		t.Fatalf("VERSION %q differs from latest changelog %q", version, latest)
	}
	for _, mode := range []string{"release", "test"} {
		t.Run(mode, func(t *testing.T) {
			temporary := t.TempDir()
			packageRoot := filepath.Join(temporary, "package")
			env := releaseEnvironment(map[string]string{"OUTDIR": packageRoot, "GOOS": runtime.GOOS, "GOARCH": runtime.GOARCH, "MAHJONG_SEED": "1"})
			if output, err := releaseCommand(root, env, "", 2*time.Minute, "bash", filepath.Join(root, "scripts", "package.sh"), mode); err != nil {
				t.Fatalf("package: %v\n%s", err, output)
			}
			expectedVersion := version
			if mode == "test" {
				expectedVersion = "dev"
			}
			packagedVersion, err := os.ReadFile(filepath.Join(packageRoot, "VERSION"))
			if err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(string(packagedVersion)) != expectedVersion {
				t.Fatalf("package VERSION=%q want=%q", packagedVersion, expectedVersion)
			}
			for _, name := range []string{"README.md", "LICENSE", filepath.Join("config", "CHANGELOG.txt")} {
				if _, err := os.Stat(filepath.Join(packageRoot, name)); err != nil {
					t.Fatal(err)
				}
			}
			executable := "play"
			if runtime.GOOS == "windows" {
				executable += ".exe"
			}
			binary := filepath.Join(packageRoot, executable)
			work := filepath.Join(temporary, "outside-package")
			if err := os.Mkdir(work, 0755); err != nil {
				t.Fatal(err)
			}
			recordPath := filepath.Join(work, "game.json")
			// Cwd has no project metadata; runtime root must resolve beside the binary.
			if output, err := releaseCommand(work, env, "END\nY\nEND\nN\n", 10*time.Second, binary, "--quick-mode", "--record", recordPath); err != nil {
				t.Fatalf("play: %v\n%s", err, output)
			}
			records := []string{recordPath, filepath.Join(work, "game.2.json")}
			for _, path := range records {
				recording, err := replay.Load(path)
				if err != nil {
					t.Fatal(err)
				}
				last := recording.Frames[len(recording.Frames)-1]
				if recording.CodeVersion != expectedVersion || last.Settlement == nil || last.Settlement.Result != result.Cancel || last.Settlement.Score != 0 {
					t.Fatalf("incorrect packaged version or cancellation: %+v", last.Settlement)
				}
			}
			if _, err := os.Stat(filepath.Join(packageRoot, "config", "config.json")); err != nil {
				t.Fatalf("runtime root not detected: %v", err)
			}
			if _, err := os.Stat(filepath.Join(work, "config")); !os.IsNotExist(err) {
				t.Fatal("runtime config appeared in working directory")
			}
			// Canceled rounds must not add statistics.
			configData, err := os.ReadFile(filepath.Join(packageRoot, "config", "config.json"))
			if err != nil {
				t.Fatal(err)
			}
			var config utils.ConfigType
			if err := json.Unmarshal(configData, &config); err != nil {
				t.Fatal(err)
			}
			if config.Total != 0 || config.Win != 0 || config.Equal != 0 || config.Score != 0 {
				t.Fatalf("canceled rounds changed statistics: %+v", config)
			}
			before := snapshotTree(t, packageRoot)
			for _, path := range append(records, filepath.Join(packageRoot, "log", "play.log")) {
				if output, err := releaseCommand(work, env, "", 10*time.Second, binary, "--replay", path, "--verify"); err != nil || !strings.Contains(output, "Replay verified:") {
					t.Fatalf("verify: %v\n%s", err, output)
				}
			}
			if after := snapshotTree(t, packageRoot); !reflect.DeepEqual(before, after) {
				t.Fatal("package replay verification changed runtime files")
			}
			// Build-time identity is authoritative even if runtime VERSION disappears.
			if err := os.Remove(filepath.Join(packageRoot, "VERSION")); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(packageRoot, "config", "CHANGELOG.txt"), []byte("v999.0.0\n"), 0644); err != nil {
				t.Fatal(err)
			}
			fallbackEnv := releaseEnvironment(map[string]string{"MAHJONG_ROOT": packageRoot, "MAHJONG_SEED": "1"})
			fallbackPath := filepath.Join(work, "fallback.json")
			if output, err := releaseCommand(work, fallbackEnv, "END\nN\n", 10*time.Second, binary, "--quick-mode", "--record", fallbackPath); err != nil {
				t.Fatalf("injected version fallback: %v\n%s", err, output)
			}
			recording, err := replay.Load(fallbackPath)
			if err != nil {
				t.Fatal(err)
			}
			if recording.CodeVersion != expectedVersion {
				t.Fatalf("lost injected version: %q", recording.CodeVersion)
			}
		})
	}
}

func releaseEnvironment(overrides map[string]string) []string {
	environment := []string{}
	for _, entry := range os.Environ() {
		key := strings.SplitN(entry, "=", 2)[0]
		switch key {
		case "MAHJONG_ROOT", "MAHJONG_SEED", "GOOS", "GOARCH", "OUTDIR":
			continue
		}
		if _, replaced := overrides[key]; !replaced {
			environment = append(environment, entry)
		}
	}
	for key, value := range overrides {
		environment = append(environment, key+"="+value)
	}
	return environment
}

func releaseCommand(directory string, environment []string, input string, timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = directory
	command.Env = environment
	command.Stdin = strings.NewReader(input)
	output, err := command.CombinedOutput()
	if ctx.Err() != nil {
		return string(output), ctx.Err()
	}
	return string(output), err
}
