package card

import (
	"errors"
	"slices"
	"sort"
)

type Card string
type Cards []Card

func (cards *Cards) Init() Cards {
	return *new(Cards)
}

func (cards *Cards) Contains(card Card) bool {
	i := cards.index(card)
	return i >= 0
}

func (cards *Cards) index(card Card) int {
	return slices.Index(*cards, card)
}

func (cards *Cards) Count(card Card) int {
	count := 0
	for _, tile := range *cards {
		if card == tile {
			count++
		}
	}
	return count
}

// RemoveCopies removes exactly count copies. The caller must ensure they exist.
func (cards *Cards) RemoveCopies(card Card, count int) {
	for i := 0; i < count; i++ {
		cards.removeOne(card)
	}
}

// Remove removes one copy, returning an error if the tile is absent.
func (cards *Cards) Remove(card Card) error {
	if !cards.Contains(card) {
		return errors.New("error")
	}
	cards.removeOne(card)
	return nil
}

// RemoveTiles removes one copy of each argument. All requested tiles must exist.
func (cards *Cards) RemoveTiles(tiles ...Card) {
	for _, tile := range tiles {
		cards.removeOne(tile)
	}
}

// PopLast removes and returns the last tile. The collection must be nonempty.
func (cards *Cards) PopLast() Card {
	lastIndex := len(*cards) - 1
	tile := (*cards)[lastIndex]
	*cards = (*cards)[:lastIndex]
	return tile
}

func (cards *Cards) removeOne(card Card) {
	i := cards.index(card)
	*cards = append((*cards)[:i], (*cards)[i+1:]...)
}

func (cards *Cards) Add(tiles ...Card) {
	*cards = append(*cards, tiles...)
}

func (cards *Cards) Sort() {
	suitedTiles := Cards{}
	honorTiles := Cards{}
	for _, tile := range *cards {
		if len(tile) == SuitedTileCodeLength {
			suitedTiles = append(suitedTiles, tile)
		} else {
			honorTiles = append(honorTiles, tile)
		}
	}
	sort.Sort(&suitedTiles)
	sort.Sort(&honorTiles)
	*cards = append(suitedTiles, honorTiles...)
}

// Display returns the bracketed tile list used by the CLI and readable logs.
func (cards *Cards) Display() string {
	display := "["
	for index, tile := range *cards {
		if index == 0 {
			display += string(tile)
		} else {
			display = display + ", " + string(tile)
		}
	}
	display += "]"
	return display
}

func (cards *Cards) Less(i, j int) bool {
	return (*cards)[i] < (*cards)[j]
}

func (cards *Cards) Swap(i, j int) {
	(*cards)[i], (*cards)[j] = (*cards)[j], (*cards)[i]
}

func (cards *Cards) Len() int {
	return len(*cards)
}
