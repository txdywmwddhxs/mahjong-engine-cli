package card

import "testing"

func TestCardsOps(t *testing.T) {
	c := &Cards{"B1", "B1", "W9", "D"}
	if !c.ContainsC("B1") {
		t.Fatalf("expected ContainsC")
	}
	if c.CountC("B1") != 2 {
		t.Fatalf("expected count 2")
	}
	if err := c.DelE("NOPE"); err == nil {
		t.Fatalf("expected error deleting non-existent card")
	}

	_ = c.DelL()
	c.DelM("B1")
	c.DelT("B1", 1)

	c.NewC("T3", "T1")
	c.Sort()
	_ = c.Beauty()
}
