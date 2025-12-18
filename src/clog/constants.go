package clog

import (
	l "language"
	"score"
)

const (
	InfoLog   = "INFO[%d]: %s\n"
	DebugLog  = "DEBUG[%d]: %s\n"
	InfoPrint = "INFO: %s\n"
)

func Trans(o func() string) string {
	return o()
}

var scoreTrans = map[score.Item]func() string{
	score.Waiting:        l.Waiting,
	score.ExposedKong:    l.ExposedKong,
	score.ConcealedKong:  l.ConcealedKong,
	score.Win:            l.Win,
	score.OwnDraw:        l.OwnDraw,
	score.Lose:           l.Lose,
	score.Single:         l.Single,
	score.MissTwoKind:    l.MissTwoKind,
	score.MissOneKind:    l.MissOneKind,
	score.Counter:        l.Counter,
	score.SevenPairs:     l.SevenPairs,
	score.ThirteenOne:    l.ThirteenOne,
	score.ContinuousLine: l.ContinuousLine,
}
