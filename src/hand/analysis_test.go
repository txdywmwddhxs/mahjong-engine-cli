package hand

import (
	"reflect"
	"testing"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
)

func TestAnalysisOwnsItsGroups(t *testing.T) {
	concealed := c.Cards{"D", "D", "B3", "B2", "B1", "T3", "T2", "T1", "W3", "W2", "W1"}
	fixed := c.Cards{"N", "N", "N"}
	originalHand := append(c.Cards(nil), concealed...)
	originalFixed := append(c.Cards(nil), fixed...)
	a := Analyze(concealed, fixed, 0)
	if a.Kind != Standard || len(a.Groups) != 4 {
		t.Fatalf("unexpected analysis: %+v", a)
	}
	second := Analyze(concealed, fixed, 0)
	a.Groups[0][0] = "?"
	if !reflect.DeepEqual(second, Analyze(concealed, fixed, 0)) || !reflect.DeepEqual(concealed, originalHand) || !reflect.DeepEqual(fixed, originalFixed) {
		t.Fatal("analysis shares mutable state")
	}
}
