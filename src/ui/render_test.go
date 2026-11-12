package ui

import (
	"strings"
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestRenderCards(t *testing.T) {
	c := &card.Cards{"B1", "W9"}
	out := RenderCards(c)
	if out != "[B1, W9]" {
		t.Fatalf("unexpected RenderCards output: %q", out)
	}
}

func TestRenderScoreDetail_UsesTranslatorAndOrder(t *testing.T) {
	tr := language.New(utils.Chinese)
	detail := map[score.Item]int{
		score.Win:         5,
		score.Lose:        -1,
		score.Counter:     2,
		score.Item("zzz"): 9, // extra item should still render
	}
	out := RenderScoreDetail(tr, detail)
	if out == "" {
		t.Fatalf("expected non-empty output")
	}
	if !strings.Contains(out, "胡牌") { // zh win
		t.Fatalf("expected zh win name in output, got: %q", out)
	}
	if !strings.Contains(out, "zzz: 9") {
		t.Fatalf("expected extra item to render, got: %q", out)
	}
	// sortBase begins with lose then win
	if !strings.HasPrefix(out, "输牌:") {
		t.Fatalf("expected output to start with lose entry, got: %q", out)
	}
}
