package card

import (
	"reflect"
	"testing"
)

func TestCardsOperations(t *testing.T) {
	cards := Cards{"B1", "B1", "W9", "D"}
	if !cards.Contains("B1") || cards.Contains("NOPE") {
		t.Fatal("unexpected membership result")
	}
	if cards.Count("B1") != 2 || cards.Count("NOPE") != 0 {
		t.Fatal("unexpected tile count")
	}
	original := append(Cards(nil), cards...)
	if err := cards.Remove("NOPE"); err == nil {
		t.Fatal("expected error removing absent tile")
	}
	if !reflect.DeepEqual(cards, original) {
		t.Fatal("failed removal changed tiles")
	}
	if tile := cards.PopLast(); tile != "D" {
		t.Fatalf("popped %q", tile)
	}
	if err := cards.Remove("B1"); err != nil {
		t.Fatal(err)
	}
	cards.Add("T3", "T1", "B1", "B1")
	cards.RemoveCopies("B1", 2)
	cards.RemoveTiles("W9", "T3")
	if want := (Cards{"T1", "B1"}); !reflect.DeepEqual(cards, want) {
		t.Fatalf("got %v want %v", cards, want)
	}
	cards.Sort()
	if want := (Cards{"B1", "T1"}); !reflect.DeepEqual(cards, want) {
		t.Fatalf("sorted %v want %v", cards, want)
	}
	if got := cards.Display(); got != "[B1, T1]" {
		t.Fatalf("display %q", got)
	}
}

func TestCardsRemoveCopiesAndTiles(t *testing.T) {
	cards := Cards{"B1", "B2", "B1", "T1", "B1"}
	cards.RemoveCopies("B1", 2)
	if want := (Cards{"B2", "T1", "B1"}); !reflect.DeepEqual(cards, want) {
		t.Fatalf("got %v want %v", cards, want)
	}
	cards.RemoveTiles("B1", "B2", "T1")
	if len(cards) != 0 || cards.Display() != "[]" {
		t.Fatalf("expected empty collection: %v", cards)
	}
	cards.Add("W1", "W1")
	cards.RemoveTiles("W1", "W1")
	if len(cards) != 0 {
		t.Fatal("duplicate arguments did not remove both copies")
	}
}
