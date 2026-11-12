// Package hand analyzes legal Mahjong hands without changing game state,
// calculating points, or emitting presentation output.
package hand

import c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"

type Kind uint8

const (
	NotWinning Kind = iota
	Standard
	SevenPairs
	ThirteenOrphans
)

// Analysis describes the selected winning decomposition. For compatibility,
// standard hands take precedence over special hands; Groups contains only
// concealed groups, with the pair first. All slices are owned by the result.
type Analysis struct {
	Kind   Kind
	Groups []c.Cards
}

func (a Analysis) IsWin() bool { return a.Kind != NotWinning }

// Analyze leaves both inputs untouched, including their order and spare capacity.
func Analyze(concealed, fixed c.Cards, kongCount int) Analysis {
	if !validWinningTiles(concealed, fixed, kongCount) {
		return Analysis{}
	}
	concealed = append(c.Cards(nil), concealed...)
	concealed.Sort()
	handSize := len(concealed)
	var remainingTiles, originalTiles = make(c.Cards, handSize), make(c.Cards, handSize)
	copy(originalTiles, concealed)
	copy(remainingTiles, originalTiles)
	fixedMeldCount := (len(fixed) - kongCount*c.KongExtraTiles) / c.MeldSize
	if len(remainingTiles)+len(fixed)-kongCount*c.KongExtraTiles != c.WinningTileCount {
		return Analysis{}
	}
	groups := make([]c.Cards, 0, c.WinningGroupCount)
	visitedTiles := c.Cards{}
	for _, tile := range originalTiles {
		if visitedTiles.Contains(tile) {
			continue
		}
		if remainingTiles.Count(tile) >= c.PairSize {
			remainingTiles.RemoveCopies(tile, c.PairSize)
			groups = append(groups, c.Cards{tile, tile})
		} else {
			continue
		}
		for i := 0; i < handSize; i++ {
			if len(remainingTiles) == 0 {
				break
			}
			lowestTile := remainingTiles[0]
			if remainingTiles.Count(lowestTile) >= c.MeldSize {
				groups = append(groups, c.Cards{lowestTile, lowestTile, lowestTile})
				remainingTiles.RemoveCopies(lowestTile, c.MeldSize)
			} else if len(lowestTile) == c.SuitedTileCodeLength { // The lowest tile belongs to a numbered suit.
				nextTile := c.Card(lowestTile[0]) + c.Card(rune(lowestTile[1]+1))
				thirdTile := c.Card(lowestTile[0]) + c.Card(rune(lowestTile[1]+2))
				if remainingTiles.Contains(nextTile) && remainingTiles.Contains(thirdTile) {
					groups = append(groups, c.Cards{lowestTile, nextTile, thirdTile})
					remainingTiles.RemoveTiles(lowestTile, nextTile, thirdTile)
				}
			}
		}
		if len(groups)+fixedMeldCount == c.WinningGroupCount {
			return Analysis{Kind: Standard, Groups: groups}
		} else {
			remainingTiles = make(c.Cards, handSize)
			copy(remainingTiles, originalTiles)

			groups = make([]c.Cards, 0, c.WinningGroupCount)
			visitedTiles.Add(tile)
		}
	}

	isSevenPairs := true
	isThirteenOrphans := true

	hasOrphanPair := false
	visitedTiles = c.Cards{}
	for _, tile := range remainingTiles {
		if visitedTiles.Contains(tile) {
			continue
		}
		visitedTiles.Add(tile)
		tileCount := remainingTiles.Count(tile)
		if isSevenPairs && tileCount != c.PairSize && tileCount != c.PairSize*2 {
			isSevenPairs = false
		}

		if isThirteenOrphans && !isOrphan(tile) {
			isThirteenOrphans = false
		}

		if isThirteenOrphans && tileCount != 1 && tileCount != c.PairSize {
			isThirteenOrphans = false
		} else if isThirteenOrphans && tileCount == c.PairSize {
			if !hasOrphanPair {
				hasOrphanPair = true
			} else {
				isThirteenOrphans = false
			}
		}
		if !isThirteenOrphans && !isSevenPairs {
			break
		}
	}
	if isThirteenOrphans && hasOrphanPair && len(fixed) == 0 {
		return Analysis{Kind: ThirteenOrphans}
	}
	if isSevenPairs && len(fixed) == 0 {
		return Analysis{Kind: SevenPairs}
	}
	return Analysis{}
}

func isOrphan(tile c.Card) bool {
	tileIndex := tileIndices[tile]
	return tileIndex >= c.SuitedKindCount || tileIndex%c.RanksPerSuit == 0 || tileIndex%c.RanksPerSuit == c.RanksPerSuit-1
}
