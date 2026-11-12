package utils

import (
	"os"
	"path/filepath"
)

var (
	// RootPath is the detected project root directory (best effort).
	RootPath string

	// VersionFilePath is the absolute path to the root VERSION file.
	VersionFilePath string

	// ConfigFilePath is the absolute path to config/config.json.
	ConfigFilePath string

	// ChangeLogPath is the absolute path to config/CHANGELOG.txt.
	ChangeLogPath string

	// LogPath is the absolute path to log/play.log.
	LogPath string

	// HistoryLogDir is the absolute path to log/history_log.
	HistoryLogDir string

	// PIDPath is the absolute path to pid/play.pid.
	PIDPath string
)

func initRuntimePaths() {
	RootPath = findProjectRoot()

	VersionFilePath = filepath.Join(RootPath, "VERSION")
	ConfigFilePath = filepath.Join(RootPath, "config", "config.json")
	ChangeLogPath = filepath.Join(RootPath, "config", "CHANGELOG.txt")
	LogPath = filepath.Join(RootPath, "log", "play.log")
	HistoryLogDir = filepath.Join(RootPath, "log", "history_log")
	PIDPath = filepath.Join(RootPath, "pid", "play.pid")

}

// PrepareRuntime is called only after selecting normal play. Package loading,
// replay, verification and command-line validation do not create runtime files.
func PrepareRuntime() error {
	for _, directory := range []string{filepath.Dir(ConfigFilePath), filepath.Dir(LogPath), HistoryLogDir, filepath.Dir(PIDPath)} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			return err
		}
	}
	return writeConfigFile()
}

func findProjectRoot() string {
	// Allow tests / packaging to override runtime root explicitly.
	if root := os.Getenv("MAHJONG_ROOT"); root != "" {
		return root
	}

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
			// VERSION marks both source trees and packaged runtime roots.
			if fileExists(filepath.Join(dir, "VERSION")) {
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
