package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ss = rand.NewSource(time.Now().UnixNano())
var rr = rand.New(ss)

func ChangeToCurrentVersion() {
	if isCurrentVersion, version := isCurrentVersionLog(); !isCurrentVersion {
		// Archive current log (best effort) before writing a new header.
		if _, err := os.Stat(LogPath); err == nil {
			archiveName := strings.TrimSpace(version)
			if archiveName == "" {
				archiveName = "unknown"
			}
			newPath := filepath.Join(HistoryLogDir, fmt.Sprintf("%s.log", archiveName))
			_ = os.Rename(LogPath, newPath)
		}

		f, err := os.OpenFile(LogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				// best effort
			}
		}(f)

		_, err = f.WriteString(fmt.Sprintf("%s\n\n", Version))
		if err != nil {
			panic(err)
		}

		Config.Version = Version
		UpdateConfigFile()
	}
}

func isCurrentVersionLog() (bool, string) {
	// last version is written in the first line of LogPath
	f, err := os.Open(LogPath)
	if err != nil {
		return false, ""
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	r := bufio.NewScanner(f)
	var version string
	if r.Scan() {
		version = r.Text()
	}
	version = strings.TrimSpace(version)
	return version == Version, version
}

func RandInt(min, max int) int {
	if min == max {
		return max
	}
	if max < min {
		min, max = max, min
	}
	return rr.Intn(max-min+1) + min
}

func ChangeQuickMode() {
	Config.QuickMode = !Config.QuickMode
	UpdateConfigFile()
}

func UpdateConfigFile() {
	file, err := os.Create(ConfigFilePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&Config); err != nil {
		panic(err)
	}
}

func LoadCurrentLog() string {
	cb, err := os.ReadFile(LogPath)
	if err != nil {
		return ""
	}
	content := string(cb)
	cl := strings.Split(content, "\n\n")
	return cl[len(cl)-1]
}
