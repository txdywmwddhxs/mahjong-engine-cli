package player

import (
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
)

func TestEarlyFinishAwardBoundary(t *testing.T) {
	for _, won := range []bool{false, true} {
		for _, turn := range []int{25, 26} {
			items := score.Items{}
			p := &Player{ScoreItems: &items, counter: turn}
			p.recordOutcome(won, false)
			if got, want := items.Contains(score.Counter), turn == 25; got != want {
				t.Errorf("turn=%d won=%v early bonus=%v want=%v", turn, won, got, want)
			}
		}
	}
}
