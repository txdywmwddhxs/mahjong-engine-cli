package card

import (
	"math/rand"
	"time"
)

var ss = rand.NewSource(time.Now().UnixNano())
var rr = rand.New(ss)

const CountOfKinds = 136

var CardsList = Cards{
	"W1",
	"W2",
	"W3",
	"W4",
	"W5",
	"W6",
	"W7",
	"W8",
	"W9",
	"B1",
	"B2",
	"B3",
	"B4",
	"B5",
	"B6",
	"B7",
	"B8",
	"B9",
	"T1",
	"T2",
	"T3",
	"T4",
	"T5",
	"T6",
	"T7",
	"T8",
	"T9",
	"D",
	"N",
	"X",
	"B",
	"Z",
	"F",
	"W",
}

func InitCardsPool() *Cards {
	cp := make(Cards, 0, CountOfKinds)
	for i := 0; i < 4; i++ {
		cp = append(cp, CardsList...)
	}
	rr.Shuffle(CountOfKinds, func(i, j int) {
		cp[i], cp[j] = cp[j], cp[i]
	})
	return &cp
}
