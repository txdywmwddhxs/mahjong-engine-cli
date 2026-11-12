package hand

import c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"

var tileIndices = func() map[c.Card]int {
	indices := make(map[c.Card]int, len(c.CardsList))
	for i, tile := range c.CardsList {
		indices[tile] = i
	}
	return indices
}()

// validWinningTiles checks physical tile limits and the fixed melds before the
// concealed-hand solver runs. FixedCards is flat: an upgraded kong's fourth
// tile may have been appended after other melds, so order is not meaningful.
func validWinningTiles(concealed, fixedCards c.Cards, kongCount int) bool {
	if kongCount < 0 || kongCount > c.WinningMeldCount {
		return false
	}
	normalizedFixedCount := len(fixedCards) - kongCount*c.KongExtraTiles
	if normalizedFixedCount < 0 || normalizedFixedCount%c.MeldSize != 0 {
		return false
	}
	fixedMeldCount := normalizedFixedCount / c.MeldSize
	if fixedMeldCount > c.WinningMeldCount || kongCount > fixedMeldCount || len(concealed)+normalizedFixedCount != c.WinningTileCount {
		return false
	}
	var total, fixed [c.TileKindCount]uint8
	for group, cards := range []c.Cards{concealed, fixedCards} {
		for _, tile := range cards {
			i, ok := tileIndices[tile]
			if !ok || total[i] == c.CopiesPerTile {
				return false
			}
			total[i]++
			if group == 1 {
				fixed[i]++
			}
		}
	}
	return canFormFixedMelds(&fixed, kongCount)
}

// Require exactly the declared number of kongs, and otherwise only pungs or
// same-suit sequences. Backtracking handles overlapping sequence/pung choices.
func canFormFixedMelds(counts *[c.TileKindCount]uint8, kongCount int) bool {
	i := 0
	for i < len(counts) && counts[i] == 0 {
		i++
	}
	if i == len(counts) {
		return kongCount == 0
	}
	if kongCount > 0 && counts[i] == c.KongSize {
		counts[i] -= c.KongSize
		ok := canFormFixedMelds(counts, kongCount-1)
		counts[i] += c.KongSize
		if ok {
			return true
		}
	}
	if counts[i] >= c.MeldSize {
		counts[i] -= c.MeldSize
		ok := canFormFixedMelds(counts, kongCount)
		counts[i] += c.MeldSize
		if ok {
			return true
		}
	}
	if i < c.SuitedKindCount && i%c.RanksPerSuit <= c.MaxSequenceStartRank-c.MinRank && counts[i+1] > 0 && counts[i+2] > 0 {
		counts[i]--
		counts[i+1]--
		counts[i+2]--
		ok := canFormFixedMelds(counts, kongCount)
		counts[i]++
		counts[i+1]++
		counts[i+2]++
		if ok {
			return true
		}
	}
	return false
}
