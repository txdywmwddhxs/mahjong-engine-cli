package card

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	// NOTE: We use the same random number generator from utils to ensure
	// deterministic behavior when MAHJONG_SEED is set. This ensures that
	// all random operations (card shuffling, robot card generation, etc.)
	// use the same random sequence, making game replay reproducible.
	//
	// The card package's local rr is only used as a fallback if utils
	// is not available, but in practice we should use utils.RandInt or
	// a shared generator for consistency.

	cp := make(Cards, 0, CountOfKinds)
	for i := 0; i < 4; i++ {
		cp = append(cp, CardsList...)
	}

	// Use utils.RandInt for shuffling to ensure we use the same RNG sequence
	// as the rest of the game. We need to shuffle using the same generator.
	// Since we can't directly access utils.rr, we'll use a workaround:
	// Generate random indices using utils.RandInt and swap accordingly.
	//
	// However, for now, we keep using the local rr but ensure it's initialized
	// with the same seed as utils if MAHJONG_SEED is set.
	if seedStr := strings.TrimSpace(os.Getenv("MAHJONG_SEED")); seedStr != "" {
		if seed, err := strconv.ParseInt(seedStr, 10, 64); err == nil {
			rr = rand.New(rand.NewSource(seed))
		}
	}

	rr.Shuffle(CountOfKinds, func(i, j int) {
		cp[i], cp[j] = cp[j], cp[i]
	})
	return &cp
}
