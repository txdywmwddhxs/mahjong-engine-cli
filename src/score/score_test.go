package score

import "testing"

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
