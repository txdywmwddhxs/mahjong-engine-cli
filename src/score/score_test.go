package score

import (
	"reflect"
	"testing"
)

func TestGetScore_LoseWithWaiting(t *testing.T) {
	items := &Items{Lose, Waiting}
	_, detail := GetScore(items, 10)
	if _, ok := detail[Lose]; !ok {
		t.Fatalf("expected Lose in detail")
	}
	if _, ok := detail[Waiting]; !ok {
		t.Fatalf("expected Waiting in detail when Losing with Waiting")
	}
}

func TestGetScore_Win_OwnDrawDoublesTotal(t *testing.T) {
	items := &Items{Win, OwnDraw, ExposedKong, MissOneKind}
	total, detail := GetScore(items, 10)
	if _, ok := detail[Win]; !ok {
		t.Fatalf("expected Win in detail")
	}
	if _, ok := detail[OwnDraw]; !ok {
		t.Fatalf("expected OwnDraw in detail")
	}

	own := detail[OwnDraw]
	withoutOwn := total - own
	if own != withoutOwn {
		t.Fatalf("expected OwnDraw to equal base sum (times=1.0); own=%d base=%d total=%d", own, withoutOwn, total)
	}
}

func TestGetScore_PanicsOnUnknownItem(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	items := &Items{Item("unknown_item")}
	_, _ = GetScore(items, 1)
}

func TestGetScoreDoesNotConsumeItems(t *testing.T) {
	items := Items{Win, OwnDraw, ExposedKong, ExposedKong, Waiting}
	original := append(Items(nil), items...)
	first, detail := GetScore(&items, 10)
	for n := 0; n < 3; n++ {
		total, again := GetScore(&items, 10)
		if total != first || !reflect.DeepEqual(again, detail) {
			t.Fatal("repeated settlement produced different results")
		}
	}
	if !reflect.DeepEqual(items, original) {
		t.Fatalf("mutated items: %v", items)
	}
	if detail[ExposedKong] != 2*KongScore {
		t.Fatal("lost repeated kong awards")
	}
}

// Concrete challenge-v1 settlements pin the original curves independently of
// the extracted constants, including the early-finish and difficulty boundaries.
func TestScoreCurveRegression(t *testing.T) {
	for _, tc := range []struct{ turn, win, loss int }{
		{0, 98, -13}, {1, 75, -13}, {25, 10, -14}, {60, 4, -17}, {90, 2, -33}, {100, 2, -79},
	} {
		winItems := Items{Win}
		lossItems := Items{Lose}
		if total, _ := GetScore(&winItems, tc.turn); total != tc.win {
			t.Errorf("turn %d win=%d want=%d", tc.turn, total, tc.win)
		}
		if total, _ := GetScore(&lossItems, tc.turn); total != tc.loss {
			t.Errorf("turn %d loss=%d want=%d", tc.turn, total, tc.loss)
		}
	}
}
