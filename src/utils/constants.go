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
	SleepTime             = 1250 * time.Millisecond
	DrawDisplayExtraDelay = 350 * time.Millisecond
)

var (
	MajorCardList       = card.Cards{"B", "T", "W"}
	MinorCardList       = card.Cards{"D", "N", "X", "B", "Z", "F", "W"}
	ThirteenOneCardList = append(card.Cards{"W1", "W9", "T1", "T9", "B1", "B9"}, MinorCardList...)
	ContinuousLineJudge = []card.Cards{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}}
)

// Statistics display keeps the historical significant digits and float width.
const (
	statisticsPercentScale   = 100
	winRateSignificantDigits = 4
	winRateFloatBitSize      = 32
)
