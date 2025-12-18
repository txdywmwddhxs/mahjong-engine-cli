package clog

import (
	"card"
	"fmt"
	l "language"
	"os"
	"score"
	"strings"
	"utils"
)

type Log struct {
	file *os.File
	lang utils.Lang
}

func NewLogger(fp string, lang utils.Lang) *Log {
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	_, _ = file.WriteString(fmt.Sprintf("\n%s\n", Trans(l.GameBegin)))
	_, _ = file.WriteString(fmt.Sprintf("%s: %s\n", Trans(l.GameVersion), utils.Config.Version))
	return &Log{
		file: file,
		lang: lang,
	}
}

func (log *Log) Close() {
	_, _ = log.file.WriteString(fmt.Sprintf("%s\n", Trans(l.GameEnd)))
	err := log.file.Close()
	if err != nil {
		return
	}
}

func (log *Log) InfoS(s string, counter int) {
	logString := fmt.Sprintf(InfoLog, counter, s)
	_, err := log.file.WriteString(logString)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf(InfoPrint, s)
}

func (log *Log) InfoC(c *card.Cards, counter int) {
	str := ""
	for ix, i := range *c {
		if ix == 0 {
			str = "[" + string(i)
		} else {
			str = str + ", " + string(i)
		}
	}
	str += "]"
	logString := fmt.Sprintf(InfoLog, counter, str)
	_, err := log.file.WriteString(logString)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf(InfoPrint, str)
}

func (log *Log) DebugS(s string, counter int) {
	logString := fmt.Sprintf(DebugLog, counter, s)
	_, err := log.file.WriteString(logString)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (log *Log) PrintScore(scoreDetail map[score.Item]int, counter int) {
	var scoreList = make([]string, 0, 7)
	detail := sortMapByArray(scoreDetail)
	for _, its := range detail {
		scoreList = append(scoreList, fmt.Sprintf("%s: %d", scoreTrans[its.Item](), its.Score))
	}
	log.InfoS(strings.Join(scoreList, ", "), counter)
}
