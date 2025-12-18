package utils

import (
	"time"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
)

const (
	Chinese Lang = "Chinese"
	English Lang = "English"
)

const (
	RateOfRobotGoodCard = "32%"
	RateOfRobotBadCard  = "10%"
	SleepTime           = 1250 * time.Millisecond
	CountOfWinAward     = 25

	GoodCardLimitRangeStart = 22
	GoodCardLimitRangeEnd   = 35
	BadCardLimitRangeStart  = 60
	BadCardLimitRangeEnd    = 90
)

var (
	MajorCardList       = card.Cards{"B", "T", "W"}
	MinorCardList       = card.Cards{"D", "N", "X", "B", "Z", "F", "W"}
	ThirteenOneCardList = append(card.Cards{"W1", "W9", "T1", "T9", "B1", "B9"}, MinorCardList...)
	ContinuousLineJudge = []card.Cards{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}}
)
