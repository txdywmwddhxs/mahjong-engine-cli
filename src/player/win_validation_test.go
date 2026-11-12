package player

import (
	"reflect"
	"testing"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestJudgeIsWinFixedMeldValidation(t *testing.T) {
	cases := []struct {
		name        string
		hand, fixed c.Cards
		kongs       int
		want        bool
	}{
		{"valid pung", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D"}, c.Cards{"W1", "W1", "W1"}, 0, true},
		{"valid sequence", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D"}, c.Cards{"W1", "W2", "W3"}, 0, true},
		{"appended kong upgrade", c.Cards{"B1", "B2", "B3", "T1", "T2", "T3", "D", "D"}, c.Cards{"W1", "W1", "W1", "W4", "W4", "W4", "W1"}, 1, true},
		{"four kongs", c.Cards{"D", "D"}, c.Cards{"W1", "W1", "W1", "W1", "W2", "W2", "W2", "W2", "W3", "W3", "W3", "W3", "W4", "W4", "W4", "W4"}, 4, true},
		{"five copies across concealed and fixed", c.Cards{"W1", "W1", "B1", "B2", "B3", "T1", "T2", "T3", "D", "D", "D"}, c.Cards{"W1", "W1", "W1"}, 0, false},
		{"cross suit sequence", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D"}, c.Cards{"W9", "B1", "B2"}, 0, false},
		{"claimed kong without four identical tiles", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D"}, c.Cards{"W1", "W2", "W3", "W4"}, 1, false},
		{"negative kong count", c.Cards{"B1", "B2", "B3", "B4", "B5", "B6", "T1", "T2", "T3", "D", "D", "D", "N", "N"}, nil, -1, false},
		{"excess kong count", c.Cards{"D", "D"}, nil, 5, false},
		{"incomplete fixed group", c.Cards{"B1", "B2", "B3", "T1", "T2", "T3", "D", "D", "D", "N", "N", "N"}, c.Cards{"W1", "W1"}, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hand := append(c.Cards(nil), tc.hand...)
			fixed := append(c.Cards(nil), tc.fixed...)
			originalFixed := append(c.Cards(nil), fixed...)
			p := &Player{PlayerCards: &hand, FixedCards: &fixed, KongCount: tc.kongs, logger: logging.NopLogger(), t: language.New(utils.Lang("English"))}
			if got := p.JudgeIsWin(); got != tc.want {
				t.Errorf("got=%v want=%v", got, tc.want)
			}
			if !reflect.DeepEqual(fixed, originalFixed) {
				t.Fatal("fixed tiles changed")
			}
		})
	}
	for _, p := range []*Player{{}, {PlayerCards: &c.Cards{}}, {FixedCards: &c.Cards{}}} {
		if p.JudgeIsWin() {
			t.Error("accepted nil card slices")
		}
	}
}

func TestCheckIsWaitingCountsFixedTiles(t *testing.T) {
	hand := c.Cards{"B1", "B2", "B3", "B4", "T1", "T2", "T3", "D", "D", "D"}
	fixed := c.Cards{"B1", "B1", "B1"}
	original := append(c.Cards(nil), hand...)
	p := &Player{PlayerCards: &hand, FixedCards: &fixed, logger: logging.NopLogger(), t: language.New(utils.Lang("English"))}
	ok, waits := p.CheckIsWaiting()
	if !ok || waits.Contains("B1") || !waits.Contains("B4") {
		t.Errorf("unexpected waits: %v", *waits)
	}
	if !reflect.DeepEqual(hand, original) {
		t.Fatal("waiting check changed hand")
	}
}
