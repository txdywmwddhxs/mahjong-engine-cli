package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	r "result"
	"strconv"
	"strings"
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
	Config  = loadConfig()
	Version = loadVersion()
)

func loadConfig() ConfigType {
	dc := ConfigType{}
	var c ConfigType
	f, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println(err)
		return dc
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	d, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return dc
	}
	err = json.Unmarshal(d, &c)
	if err != nil {
		fmt.Println(err)
		return dc
	}
	return c
}

func loadVersion() string {
	content, err := os.ReadFile(ChangeLogPath)
	if err != nil {
		return ""
	}
	contentList := strings.Split(string(content), "\n")
	length := len(contentList)
	for i := 0; i < length; i++ {
		line := strings.TrimSpace(contentList[length-i-1])
		if strings.HasPrefix(line, "v") {
			return line
		}
	}
	return ""
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
	}
	Config.Rate = strconv.FormatFloat(float64(Config.Win)/float64(Config.Total)*100,
		'G', 4, 32) + "%"
	Config.Score += score
	UpdateConfigFile()
}
