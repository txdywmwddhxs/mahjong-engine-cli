package card

import (
	"math/rand"
	"testing"
)

func TestInitCardsPool_Counts(t *testing.T) {
	// Make shuffle deterministic for this test.
	rr = rand.New(rand.NewSource(1))

	cp := InitCardsPool()
	if cp == nil {
		t.Fatalf("expected non-nil pool")
	}
	if len(*cp) != CountOfKinds {
		t.Fatalf("expected %d cards, got %d", CountOfKinds, len(*cp))
	}

	counts := map[Card]int{}
	for _, c := range *cp {
		counts[c]++
	}
	for _, c := range CardsList {
		if counts[c] != 4 {
			t.Fatalf("expected 4 of %s, got %d", c, counts[c])
		}
	}
}
