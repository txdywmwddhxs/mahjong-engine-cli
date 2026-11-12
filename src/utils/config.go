package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	r "github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
)

type Lang string

type ConfigType struct {
	Version     string  `json:"version"`
	Total       int     `json:"total"`
	Win         int     `json:"win"`
	Equal       int     `json:"equal"`
	Score       int     `json:"score"`
	WinScore    int     `json:"win_score"`
	LoseScore   int     `json:"lose_score"`
	WinScorePM  float64 `json:"win_score_pm"`
	LoseScorePM float64 `json:"lose_score_pm"`
	Rate        string  `json:"rate"`
	Lang        Lang    `json:"lang"`
	QuickMode   bool    `json:"quick_mode"`
}

var (
	Config  ConfigType
	Version string
)

func init() {
	initRuntimePaths()

	// Allow build-time injection (e.g. via -ldflags "-X .../utils.Version=...").
	// VERSION is the only file-based source; changelog is display content.
	Version = resolveVersion(Version)
	Config = loadConfig()

	// Defaults for first run / missing config.json.
	if Config.Lang == "" {
		Config.Lang = Chinese
	}
	if Config.Version == "" {
		Config.Version = Version
	}

}

func loadConfig() ConfigType {
	dc := ConfigType{}
	var c ConfigType
	f, err := os.Open(ConfigFilePath)
	if err != nil {
		return dc
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			// best effort
		}
	}(f)
	d, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return dc
	}
	err = json.Unmarshal(d, &c)
	if err != nil {
		return dc
	}
	return c
}

// resolveVersion prefers build-time injection so a packaged binary retains its
// version even if runtime metadata is missing. Unversioned builds use dev.
func resolveVersion(injected string) string {
	if version := strings.TrimSpace(injected); version != "" {
		return version
	}
	content, err := os.ReadFile(VersionFilePath)
	if err == nil {
		if version := strings.TrimSpace(string(content)); version != "" {
			return version
		}
	}
	return "dev"
}

func UpdateConfig(result r.Result, score int) {
	if result == r.Cancel {
		return
	}
	Config.Total++
	switch result {
	case r.Win:
		Config.Win++
		Config.WinScore += score
		winPMStr := fmt.Sprintf("%.2f", float64(Config.WinScore)/float64(Config.Win))
		Config.WinScorePM, _ = strconv.ParseFloat(winPMStr, 64)
	case r.Equal:
		Config.Equal++
	case r.Lose:
		Config.LoseScore += -score
		losePM := fmt.Sprintf("%.2f", float64(Config.LoseScore)/float64(Config.Total-Config.Win-Config.Equal))
		Config.LoseScorePM, _ = strconv.ParseFloat(losePM, 64)
	default:
		panic("Invalid result value")
	}
	Config.Rate = strconv.FormatFloat(float64(Config.Win)/float64(Config.Total)*statisticsPercentScale,
		'G', winRateSignificantDigits, winRateFloatBitSize) + "%"
	Config.Score += score
	UpdateConfigFile()
}
