package utils

import (
	"os"
	"path/filepath"
)

var (
	// RootPath is the detected project root directory (best effort).
	RootPath string

	// ConfigFilePath is the absolute path to config/config.json.
	ConfigFilePath string

	// ChangeLogPath is the absolute path to config/ChangeLog.
	ChangeLogPath string

	// LogPath is the absolute path to log/play.log.
	LogPath string

	// HistoryLogDir is the absolute path to log/history_log.
	HistoryLogDir string

	// PIDPath is the absolute path to pid/play.pid.
	PIDPath string
)

func initRuntimePaths() {
	root := findProjectRoot()
	RootPath = root

	ConfigFilePath = filepath.Join(root, "config", "config.json")
	ChangeLogPath = filepath.Join(root, "config", "ChangeLog")
	LogPath = filepath.Join(root, "log", "play.log")
	HistoryLogDir = filepath.Join(root, "log", "history_log")
	PIDPath = filepath.Join(root, "pid", "play.pid")

	// Ensure required directories exist.
	_ = os.MkdirAll(filepath.Dir(ConfigFilePath), 0o755)
	_ = os.MkdirAll(filepath.Dir(LogPath), 0o755)
	_ = os.MkdirAll(HistoryLogDir, 0o755)
	_ = os.MkdirAll(filepath.Dir(PIDPath), 0o755)
}

func findProjectRoot() string {
	var starts []string

	if wd, err := os.Getwd(); err == nil && wd != "" {
		starts = append(starts, wd)
	}
	if exe, err := os.Executable(); err == nil && exe != "" {
		starts = append(starts, filepath.Dir(exe))
	}

	for _, start := range starts {
		dir := start
		for {
			// Prefer ChangeLog marker if present.
			if fileExists(filepath.Join(dir, "config", "ChangeLog")) {
				return dir
			}
			// Fallback marker for module root.
			if fileExists(filepath.Join(dir, "go.mod")) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// Best effort fallback.
	if len(starts) != 0 {
		return starts[0]
	}
	return "."
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
