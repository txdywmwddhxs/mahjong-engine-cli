package score

import (
	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/hand"
)

// WinningPatterns maps hand analysis to the existing scoring rules. A straight
// is awarded only from the selected concealed decomposition, as in challenge-v1.
func WinningPatterns(analysis hand.Analysis) Items {
	switch analysis.Kind {
	case hand.SevenPairs:
		return Items{SevenPairs}
	case hand.ThirteenOrphans:
		return Items{ThirteenOrphans}
	case hand.Standard:
		for _, suit := range []byte{'B', 'T', 'W'} {
			sequenceStarts := [c.RanksPerSuit / c.MeldSize]bool{}
			for _, group := range analysis.Groups {
				if len(group) != c.MeldSize {
					continue
				}
				firstTile := group[0]
				if len(firstTile) != c.SuitedTileCodeLength || firstTile[0] != suit {
					continue
				}
				rank := firstTile[1]
				for sequenceIndex := range sequenceStarts {
					startRank := byte('0' + c.MinRank + sequenceIndex*c.MeldSize)
					if rank == startRank && group[1] == c.Card([]byte{suit, rank + 1}) && group[2] == c.Card([]byte{suit, rank + 2}) {
						sequenceStarts[sequenceIndex] = true
					}
				}
			}
			hasFullStraight := true
			for _, present := range sequenceStarts {
				if !present {
					hasFullStraight = false
					break
				}
			}
			if hasFullStraight {
				return Items{ContinuousLine}
			}
		}
	}
	return nil
}

// MissingSuits evaluates the final tiles without appending into either input.
func MissingSuits(concealed, fixed c.Cards) Items {
	suits := map[byte]bool{}
	for _, cards := range []c.Cards{concealed, fixed} {
		for _, tile := range cards {
			if len(tile) == c.SuitedTileCodeLength {
				suits[tile[0]] = true
			}
		}
	}
	missingSuitCount := c.SuitCount - len(suits)
	if missingSuitCount >= maxMissingSuitBonus {
		return Items{MissTwoKind}
	}
	if missingSuitCount == 1 {
		return Items{MissOneKind}
	}
	return nil
}
