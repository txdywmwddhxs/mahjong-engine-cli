package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var ss = rand.NewSource(time.Now().UnixNano())
var rr = rand.New(ss)

func ChangeToCurrentVersion() {
	if isCurrentVersion, version := isCurrentVersionLog(); !isCurrentVersion {
		newPath := fmt.Sprintf(HistoryLogPath, version)
		err := os.Rename(LogPath, newPath)
		if err != nil {
			panic(err)
		}
		f, err := os.OpenFile(LogPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println(f)
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

//func UpdateConfigFile() {
//	output, err := json.MarshalIndent(&Config, "", "  ")
//	if err != nil {
//		panic(err)
//	}
//	err = ioutil.WriteFile(configFilePath, output, 0666)
//	if err != nil {
//		panic(err)
//	}
//}

func UpdateConfigFile() {
	file, err := os.Create(configFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

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
