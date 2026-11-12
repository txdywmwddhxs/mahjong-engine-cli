package player

import (
	"math/rand"
	"os"
	"testing"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

// Independent count-vector oracle: try every pair and both meld branches.
// Index order is CardsList (three suits, then seven honors).
func experimentWin(counts [34]uint8, concealed int, fixed int) bool {
	if concealed != 14-3*fixed {
		return false
	}
	if fixed == 0 {
		pairs, orphans, duplicate := true, true, 0
		for i, n := range counts {
			if n%2 != 0 {
				pairs = false
			}
			terminal := i >= 27 || i%9 == 0 || i%9 == 8
			if terminal {
				if n == 2 {
					duplicate++
				} else if n != 1 {
					orphans = false
				}
			} else if n != 0 {
				orphans = false
			}
		}
		if pairs || (orphans && duplicate == 1) {
			return true
		}
	}
	var melds func() bool
	melds = func() bool {
		i := 0
		for i < 34 && counts[i] == 0 {
			i++
		}
		if i == 34 {
			return true
		}
		if counts[i] >= 3 {
			counts[i] -= 3
			ok := melds()
			counts[i] += 3
			if ok {
				return true
			}
		}
		if i < 27 && i%9 <= 6 && counts[i+1] > 0 && counts[i+2] > 0 {
			counts[i]--
			counts[i+1]--
			counts[i+2]--
			ok := melds()
			counts[i]++
			counts[i+1]++
			counts[i+2]++
			if ok {
				return true
			}
		}
		return false
	}
	for i, n := range counts {
		if n >= 2 {
			counts[i] -= 2
			ok := melds()
			counts[i] += 2
			if ok {
				return true
			}
		}
	}
	return false
}

// Opt-in diagnostic: reports every category and fails on disagreement, while
// keeping ordinary go test runs fast. Seed and generation order are fixed.
// MAHJONG_WIN_EXPERIMENT=1 go test ./src/player -run TestJudgeIsWinExperiment -v -count=1
func TestJudgeIsWinExperiment(t *testing.T) {
	if os.Getenv("MAHJONG_WIN_EXPERIMENT") != "1" {
		t.Skip("set MAHJONG_WIN_EXPERIMENT=1 for large differential experiment")
	}
	type stats struct{ total, wins, falsePositive, falseNegative int }
	results := map[string]*stats{}
	check := func(category string, counts [34]uint8, fixed c.Cards, kongs int) {
		hand := c.Cards{}
		for i, n := range counts {
			for j := uint8(0); j < n; j++ {
				hand = append(hand, c.CardsList[i])
			}
		}
		want := experimentWin(counts, len(hand), (len(fixed)-kongs)/3)
		p := &Player{PlayerCards: &hand, FixedCards: &fixed, KongCount: kongs}
		got := p.JudgeIsWin()
		s := results[category]
		if s == nil {
			s = &stats{}
			results[category] = s
		}
		s.total++
		if want {
			s.wins++
		}
		if got != want {
			if got {
				s.falsePositive++
			} else {
				s.falseNegative++
			}
			if s.falsePositive+s.falseNegative <= 5 {
				t.Logf("MISMATCH %s hand=%v fixed=%v kongs=%d got=%v want=%v", category, hand, fixed, kongs, got, want)
			}
		}
	}
	// All 0..4 counts of ranks 1..9, for every possible concealed
	// size 2,5,8,11,14. Fixed honor pungs supply remaining melds.
	var enumerate func(int, int, [34]uint8)
	enumerate = func(i, total int, counts [34]uint8) {
		if total > 14 {
			return
		}
		if i == 9 {
			if total >= 2 && total%3 == 2 {
				fixed := c.Cards{}
				for j := 0; j < (14-total)/3; j++ {
					for k := 0; k < 3; k++ {
						fixed = append(fixed, c.CardsList[27+j])
					}
				}
				check("exhaustive-single-suit", counts, fixed, 0)
			}
			return
		}
		for n := 0; n <= 4; n++ {
			counts[i] = uint8(n)
			enumerate(i+1, total+n, counts)
		}
	}
	enumerate(0, 0, [34]uint8{})
	rng := rand.New(rand.NewSource(20261112))
	// Uniform draws without replacement from the physical 136-tile wall.
	for sample := 0; sample < 1000000; sample++ {
		groups := rng.Intn(5)
		fixed := c.Cards{}
		kongs := 0
		used := [34]uint8{}
		for j := 0; j < groups; j++ {
			i := 27 + j
			n := 3 + rng.Intn(2)
			if n == 4 {
				kongs++
			}
			used[i] = uint8(n)
			for k := 0; k < n; k++ {
				fixed = append(fixed, c.CardsList[i])
			}
		}
		counts := [34]uint8{}
		for n := 0; n < 14-3*groups; {
			i := rng.Intn(136) / 4
			if used[i] < 4 {
				counts[i]++
				used[i]++
				n++
			}
		}
		check("random-wall", counts, fixed, kongs)
	}
	// Construct valid pair+meld decompositions, with arbitrary exposed pungs,
	// sequences and kongs, rejecting any fifth physical copy. Half are then
	// mutated by replacing a tile, so the oracle decides if they still win.
	for sample := 0; sample < 1000000; {
		groups := rng.Intn(5)
		counts := [34]uint8{}
		used := [34]uint8{}
		fixed := c.Cards{}
		kongs := 0
		pair := rng.Intn(34)
		counts[pair] = 2
		used[pair] = 2
		for g := 0; g < 4; g++ {
			ids := []int{}
			i := rng.Intn(34)
			if rng.Intn(2) == 0 {
				start := rng.Intn(3)*9 + rng.Intn(7)
				ids = []int{start, start + 1, start + 2}
			} else {
				ids = []int{i, i, i}
			}
			exposed := g < groups
			if exposed && ids[0] == ids[1] && rng.Intn(2) == 0 {
				ids = append(ids, ids[0])
				kongs++
			}
			for _, idx := range ids {
				used[idx]++
				if exposed {
					fixed = append(fixed, c.CardsList[idx])
				} else {
					counts[idx]++
				}
			}
		}
		valid := true
		for _, n := range used {
			if n > 4 {
				valid = false
			}
		}
		if !valid {
			continue
		}
		category := "constructed-win"
		if sample%2 == 1 {
			from, to := rng.Intn(34), rng.Intn(34)
			for counts[from] == 0 {
				from = rng.Intn(34)
			}
			for used[to] >= 4 || to == from {
				to = rng.Intn(34)
			}
			counts[from]--
			counts[to]++
			category = "mutated-win"
		}
		check(category, counts, fixed, kongs)
		sample++
	}
	// Special winning families are too rare in uniform random hands.
	for sample := 0; sample < 10000; sample++ {
		counts := [34]uint8{}
		if sample%2 == 0 {
			for n := 0; n < 7; {
				i := rng.Intn(34)
				if counts[i] <= 2 {
					counts[i] += 2
					n++
				}
			}
			check("seven-pairs", counts, nil, 0)
		} else {
			for i := 0; i < 34; i++ {
				if i >= 27 || i%9 == 0 || i%9 == 8 {
					counts[i] = 1
				}
			}
			i := rng.Intn(34)
			for counts[i] == 0 {
				i = rng.Intn(34)
			}
			counts[i]++
			check("thirteen-orphans", counts, nil, 0)
		}
	}
	mismatches := 0
	for _, category := range []string{"exhaustive-single-suit", "random-wall", "constructed-win", "mutated-win", "seven-pairs", "thirteen-orphans"} {
		s := results[category]
		t.Logf("%s: total=%d oracle_wins=%d false_positive=%d false_negative=%d", category, s.total, s.wins, s.falsePositive, s.falseNegative)
		mismatches += s.falsePositive + s.falseNegative
	}
	if mismatches > 0 {
		t.Fatalf("%d disagreements with independent exhaustive-backtracking oracle", mismatches)
	}
}

// Invalid inputs must be rejected even if their concealed tiles can be split.
func TestJudgeIsWinRejectsInvalidInput(t *testing.T) {
	tr := language.New(utils.Lang("English"))
	cases := []struct {
		name        string
		hand, fixed c.Cards
		kongs       int
	}{
		{"five-identical-tiles", c.Cards{"B1", "B1", "B1", "B1", "B1", "B2", "B2", "B2", "B3", "B3", "B3", "D", "D", "D"}, nil, 0},
		{"unknown-tile-pair", c.Cards{"?", "?", "B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D", "D"}, nil, 0},
		{"invalid-fixed-meld", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D"}, c.Cards{"W1", "W4", "W8"}, 0},
	}
	for _, tc := range cases {
		p := &Player{PlayerCards: &tc.hand, FixedCards: &tc.fixed, KongCount: tc.kongs}
		if p.JudgeIsWin() {
			t.Errorf("%s: accepted invalid hand", tc.name)
		}
	}
	hand := c.Cards{"B1", "B1", "B1", "B1", "B2", "B2", "B2", "B3", "B3", "B3", "D", "D", "D"}
	fixed := c.Cards{}
	p := &Player{PlayerCards: &hand, FixedCards: &fixed, logger: logging.NopLogger(), t: tr}
	original := append(c.Cards(nil), hand...)
	ok, waiting := p.CheckIsWaiting()
	if !ok || len(*waiting) != 3 || !waiting.Contains("B2") || !waiting.Contains("B3") || !waiting.Contains("B4") {
		t.Errorf("unexpected waits: %v", *waiting)
	}
	for i, tile := range original {
		if hand[i] != tile {
			t.Fatal("waiting check changed original hand")
		}
	}
}

func TestExperimentOracleKnownHands(t *testing.T) {
	cases := []struct {
		hand c.Cards
		want bool
	}{
		{c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D", "D", "N", "N"}, true},
		{c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D", "D", "N", "X"}, false},
		{c.Cards{"B1", "B1", "B2", "B2", "B4", "B4", "T1", "T1", "T3", "T3", "D", "D", "N", "N"}, true},
		{c.Cards{"B1", "B1", "B1", "B1", "B4", "B4", "T1", "T1", "T3", "T3", "D", "D", "N", "N"}, true},
		{c.Cards{"W1", "W9", "B1", "B9", "T1", "T9", "D", "N", "X", "B", "Z", "F", "W", "W"}, true},
		{c.Cards{"W1", "W9", "B1", "B9", "T1", "T9", "D", "N", "X", "B", "Z", "F", "F", "F"}, false},
	}
	for _, tc := range cases {
		counts := [34]uint8{}
		for _, tile := range tc.hand {
			for i, v := range c.CardsList {
				if tile == v {
					counts[i]++
					break
				}
			}
		}
		if got := experimentWin(counts, len(tc.hand), 0); got != tc.want {
			t.Errorf("oracle hand=%v got=%v want=%v", tc.hand, got, tc.want)
		}
		hand := append(c.Cards(nil), tc.hand...)
		fixed := c.Cards{}
		p := &Player{PlayerCards: &hand, FixedCards: &fixed}
		if got := p.JudgeIsWin(); got != tc.want {
			t.Errorf("JudgeIsWin hand=%v got=%v want=%v", tc.hand, got, tc.want)
		}
	}
}
