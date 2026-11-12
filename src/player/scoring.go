package player

import (
	"fmt"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/hand"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/rules"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
)

// Apply scoring and presentation only after the game accepts a win. Keeping
// this at the original transition preserves intermediate replay score items.
func (p *Player) recordWinningPatterns(analysis hand.Analysis) {
	switch analysis.Kind {
	case hand.Standard:
		p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyCardAnalysisResult), analysis.Groups))
	case hand.SevenPairs:
		p.debug(p.t.T(language.KeyMatchSevenPairs))
	case hand.ThirteenOrphans:
		p.debug(p.t.T(language.KeyMatchThirteenOne))
	}
	*p.ScoreItems = append(*p.ScoreItems, score.WinningPatterns(analysis)...)
}

func (p *Player) recordMissingSuits() {
	if p.result == result.Win {
		*p.ScoreItems = append(*p.ScoreItems, score.MissingSuits(*p.PlayerCards, *p.FixedCards)...)
	}
}

// These adapters retain the historical score-item order in replay snapshots.
// Game actions report what happened; score bookkeeping stays in this file.
func (p *Player) recordOutcome(won, ownDraw bool) {
	if won {
		p.ScoreItems.Add(score.Win)
	} else {
		p.ScoreItems.Add(score.Lose)
	}
	if p.counter <= rules.EarlyFinishTurnLimit {
		p.ScoreItems.Add(score.Counter)
	}
	if ownDraw {
		p.ScoreItems.Add(score.OwnDraw)
	}
}

func (p *Player) recordWaiting(waitCount int) {
	p.ScoreItems.Add(score.Waiting)
	if waitCount == 1 {
		p.ScoreItems.Add(score.Single)
	}
}

func (p *Player) recordKong(concealed bool) {
	if concealed {
		p.ScoreItems.Add(score.ConcealedKong)
	} else {
		p.ScoreItems.Add(score.ExposedKong)
	}
}

func (p *Player) settleScore() (int, map[score.Item]int) {
	p.recordMissingSuits()
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyScoreItem), *p.ScoreItems))
	return score.GetScore(p.ScoreItems, p.counter-1)
}
