package player

import (
	"reflect"
	"testing"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
)

func TestJudgeIsWinHasNoSideEffects(t *testing.T) {
	cases := []struct {
		name  string
		cards c.Cards
		want  bool
	}{
		{"standard", c.Cards{"W9", "T1", "B9", "B8", "B7", "B6", "B5", "B4", "B3", "B2", "B1", "T1", "T1", "W9"}, true},
		{"seven pairs", c.Cards{"W9", "W9", "T3", "T3", "T1", "T1", "B4", "B4", "B2", "B2", "B1", "B1", "W5", "W5"}, true},
		{"orphans", c.Cards{"W", "F", "Z", "B", "X", "N", "D", "W9", "W1", "T9", "T1", "B9", "B1", "B1"}, true},
		{"not winning", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D", "D", "N", "X"}, false},
		{"invalid", c.Cards{"?", "?", "B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D", "D"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Keep a spare-capacity sentinel to catch accidental append into the input.
			backing := append(append(c.Cards(nil), tc.cards...), c.Card("sentinel"))
			before := append(c.Cards(nil), backing...)
			cards := backing[:len(tc.cards)]
			fixed := c.Cards{}
			items := score.Items{score.Waiting, score.ExposedKong}
			p := &Player{PlayerCards: &cards, FixedCards: &fixed, ScoreItems: &items, counter: 7, isWaiting: true}
			beforePlayer := *p
			for i := 0; i < 3; i++ {
				if got := p.JudgeIsWin(); got != tc.want {
					t.Fatalf("got=%v want=%v", got, tc.want)
				}
			}
			if !reflect.DeepEqual(*p, beforePlayer) || !reflect.DeepEqual(backing, before) || !reflect.DeepEqual(items, score.Items{score.Waiting, score.ExposedKong}) {
				t.Fatal("JudgeIsWin changed player state or input storage")
			}
			// No translator, logger or score list is required to query a win.
			p.ScoreItems = nil
			if got := p.JudgeIsWin(); got != tc.want {
				t.Fatal("result depends on score state")
			}
		})
	}
}
