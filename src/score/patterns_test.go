package score

import (
	"reflect"
	"testing"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/hand"
)

func TestWinningPatterns(t *testing.T) {
	cases := []struct {
		name             string
		concealed, fixed c.Cards
		want             Items
	}{
		{"straight", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "T1", "T1", "T1", "W9", "W9"}, nil, Items{ContinuousLine}},
		{"seven pairs", c.Cards{"B1", "B1", "B2", "B2", "B4", "B4", "T1", "T1", "T3", "T3", "W5", "W5", "W9", "W9"}, nil, Items{SevenPairs}},
		{"orphans", c.Cards{"B1", "B9", "T1", "T9", "W1", "W9", "D", "N", "X", "B", "Z", "F", "W", "B1"}, nil, Items{ThirteenOrphans}},
		{"standard takes precedence", c.Cards{"B1", "B1", "B2", "B2", "B3", "B3", "T1", "T1", "T2", "T2", "T3", "T3", "W9", "W9"}, nil, nil},
		{"fixed sequence excluded by existing rule", c.Cards{"B4", "B5", "B6", "B7", "B8", "B9", "T1", "T1", "T1", "W9", "W9"}, c.Cards{"B1", "B2", "B3"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := hand.Analyze(tc.concealed, tc.fixed, 0)
			if !a.IsWin() {
				t.Fatal("fixture is not winning")
			}
			for i := 0; i < 2; i++ {
				if got := WinningPatterns(a); !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("got=%v want=%v", got, tc.want)
				}
			}
		})
	}
	if got := WinningPatterns(hand.Analysis{}); len(got) != 0 {
		t.Fatal("nonwinning hand received scoring patterns")
	}
}

func TestMissingSuitsIncludesFixedAndPreservesStorage(t *testing.T) {
	backing := c.Cards{"B1", "sentinel"}
	concealed := backing[:1]
	fixed := c.Cards{"T1", "T2", "T3"}
	if got := MissingSuits(concealed, fixed); !reflect.DeepEqual(got, Items{MissOneKind}) {
		t.Fatalf("got %v", got)
	}
	if backing[1] != "sentinel" {
		t.Fatal("overwrote concealed storage")
	}
	if got := MissingSuits(c.Cards{"D", "D"}, nil); !reflect.DeepEqual(got, Items{MissTwoKind}) {
		t.Fatalf("honors: got %v", got)
	}
	if got := MissingSuits(c.Cards{"B1", "T1"}, c.Cards{"W1"}); len(got) != 0 {
		t.Fatalf("three suits: got %v", got)
	}
}
