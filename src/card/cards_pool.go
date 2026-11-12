package card

import (
	"math/rand"
	"time"
)

var ss = rand.NewSource(time.Now().UnixNano())
var rr = rand.New(ss)

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

// InitCardsPool is retained for callers that do not own a game RNG.
func InitCardsPool() *Cards { return InitCardsPoolWithRNG(rr) }

// InitCardsPoolWithRNG uses the same per-round source as all other setup choices.
func InitCardsPoolWithRNG(rng *rand.Rand) *Cards {
	cp := make(Cards, 0, TotalTiles)
	for i := 0; i < CopiesPerTile; i++ {
		cp = append(cp, CardsList...)
	}
	rng.Shuffle(len(cp), func(i, j int) { cp[i], cp[j] = cp[j], cp[i] })
	return &cp
}
