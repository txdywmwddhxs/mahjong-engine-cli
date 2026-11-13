package player

import (
	"fmt"
	"strconv"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/hand"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/rules"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func (p *Player) info(msg string) {
	p.logger.Info(msg, p.counter)
	if p.ui != nil {
		p.ui.Info(msg)
	}
}

func (p *Player) infoCards(cards *c.Cards) {
	p.info(ui.RenderCards(cards))
}

func (p *Player) debug(msg string) {
	p.logger.Debug(msg, p.counter)
}

func (p *Player) initCardsPool() {
	p.CardsPool = c.InitCardsPoolWithRNG(p.rng)
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyInitCardPoolDone), *p.CardsPool))
}

func (p *Player) initPlayerCards() {
	*p.PlayerCards = (*p.CardsPool)[:c.InitialHandSize]
	p.PlayerCards.Sort()
	*p.CardsPool = (*p.CardsPool)[c.InitialHandSize:]
	p.debug(fmt.Sprintf("%s: %v", p.t.T(language.KeyInitPlayerCardDone), p.PlayerCards.Display()))
}

func (p *Player) initRobotCards() {
	robotHand := c.Cards{}
	// Generate the pair.
	pairTile := (*p.CardsPool)[0]
	robotHand = append(robotHand, pairTile, pairTile)
	p.CardsPool.RemoveTiles(pairTile, pairTile)

	pungCount := p.randInt(0, rules.RobotMaxPungCount)
	sequenceCount := c.WinningMeldCount - pungCount
	pungs := make([]c.Cards, 0, pungCount)
	sequences := make([]c.Cards, 0, sequenceCount)

	// Generate pungs.
	if pungCount != 0 {
		for {
			// Select a numbered suit or honors.
			kindIndex := p.randInt(1, c.SuitCount+1)
			var pungTile c.Card
			switch kindIndex {
			case 1:
				suit := utils.MajorCardList[kindIndex-1]
				rank := p.randInt(c.MinRank, c.MaxRank)
				pungTile = c.Card(string(suit) + strconv.Itoa(rank))
			case 2:
				suit := utils.MajorCardList[kindIndex-1]
				rank := p.randInt(c.MinRank, c.MaxRank)
				pungTile = c.Card(string(suit) + strconv.Itoa(rank))
			case 3:
				suit := utils.MajorCardList[kindIndex-1]
				rank := p.randInt(c.MinRank, c.MaxRank)
				pungTile = c.Card(string(suit) + strconv.Itoa(rank))
			default:
				honorIndex := p.randInt(0, c.HonorKindCount-1)
				pungTile = utils.MinorCardList[honorIndex]
			}

			if p.CardsPool.Count(pungTile) >= c.MeldSize {
				pungs = append(pungs, c.Cards{pungTile, pungTile, pungTile})
				p.CardsPool.RemoveCopies(pungTile, c.MeldSize)
			}

			if len(pungs) == pungCount {
				break
			}
			if len(pungs) > pungCount {
				panic("generate tri error")
			}
		}
	}

	// Generate sequences.
	if sequenceCount != 0 {
		for {
			suitIndex := p.randInt(0, c.SuitCount-1)
			// Select the suit.
			suit := utils.MajorCardList[suitIndex]
			// Select the middle rank.
			middleRank := p.randInt(c.MinRank+1, c.MaxRank-1)
			middleTile := suit + c.Card(strconv.Itoa(middleRank))
			lowerTile := suit + c.Card(strconv.Itoa(middleRank-1))
			upperTile := suit + c.Card(strconv.Itoa(middleRank+1))
			if p.CardsPool.Contains(middleTile) && p.CardsPool.Contains(lowerTile) && p.CardsPool.Contains(upperTile) {
				sequences = append(sequences, c.Cards{middleTile, lowerTile, upperTile})
				p.CardsPool.RemoveTiles(upperTile, lowerTile, middleTile)
			}

			if len(sequences) == sequenceCount {
				break
			}
			if len(sequences) > sequenceCount {
				panic("generate seq error")
			}
		}
	}

	for _, pung := range pungs {
		robotHand = append(robotHand, pung...)
	}
	for _, sequence := range sequences {
		robotHand = append(robotHand, sequence...)
	}

	robotHand.Sort()
	*p.RobotCards = robotHand
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyInitRobotCardDone), p.RobotCards.Display()))
}

func (p *Player) generateLimit() {
	goodHandThreshold := rules.RobotGoodHandPercent
	badHandThreshold := rules.RobotBadHandPercent
	badHandThreshold += goodHandThreshold
	p.debug(fmt.Sprintf("%s: %d%%, %d%%, %d%%",
		p.t.T(language.KeyRateOfRobotsCard), goodHandThreshold, badHandThreshold-goodHandThreshold, rules.PercentScale-badHandThreshold))
	difficultyRoll := p.randInt(0, rules.MaxDifficultyRoll)
	p.debug(fmt.Sprintf("%s: %d", p.t.T(language.KeyRandomNumberSeed), difficultyRoll))
	if difficultyRoll < goodHandThreshold {
		p.limit = p.randInt(rules.GoodHandLimitMin, rules.GoodHandLimitMax)
	} else if difficultyRoll < badHandThreshold {
		p.limit = p.randInt(rules.GoodHandLimitMax, rules.BadHandLimitMin)
	} else {
		p.limit = p.randInt(rules.BadHandLimitMin, rules.BadHandLimitMax)
	}
}

func (p *Player) init() {
	p.generateLimit()
	p.debug(fmt.Sprintf("%s: %d", p.t.T(language.KeyLimit), p.limit))
	p.initCardsPool()
	p.initRobotCards()
	p.initPlayerCards()
}

func (p *Player) drawAndPlay(checkOpeningKongs bool) {
	newCard := (*p.CardsPool)[0]
	p.PlayerCards.Add(newCard)
	p.CardsPool.RemoveTiles(newCard)
	p.emit("draw", newCard)
	p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyGetNewCard), newCard))

	if p.counter == rules.FirstTurn && checkOpeningKongs {
		p.PlayerCards.Sort()
		handledCards := c.Cards{}

	search:
		for _, card := range *p.PlayerCards {
			if handledCards.Contains(card) {
				continue
			}
			handledCards.Add(card)
			if p.PlayerCards.Count(card) == c.KongSize {
				kongAccepted := false
				isValidInput := false
				for {
					p.infoCards(p.PlayerCards)
					kongAccepted, isValidInput = p.confirmMeld(card, false, false)
					if isValidInput && p.result == result.Default {
						break
					}
				}
				if kongAccepted {
					p.declareKong(card, c.KongSize, true)
					goto search
				}
			}
		}
	}
	if p.result != result.Default {
		return
	}

	if p.sleepTime != 0 && p.counter != rules.FirstTurn {
		p.sleep(p.sleepTime + utils.DrawDisplayExtraDelay)
	}
	p.PlayerCards.Sort()
	if len(*p.FixedCards) != 0 {
		p.FixedCards.Sort()
		p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyFixedCard), p.FixedCards.Display()))
	}
	p.infoCards(p.PlayerCards)

	analysis := p.analyzeWin()
	if analysis.IsWin() {
		p.recordWinningPatterns(analysis)
		p.debug(p.t.T(language.KeyScoreWin))
		p.result = result.Win
		p.recordOutcome(true, p.isClear)
		p.printResult()
		return
	}
	p.debug(p.t.T(language.KeyUserNotWin))

	if p.PlayerCards.Count(newCard) == c.KongSize && !p.isWaiting {
		kongAccepted, isValidInput := false, false
		for {
			if isValidInput {
				break
			}
			kongAccepted, isValidInput = p.confirmMeld(newCard, false, false)
		}
		if p.result != result.Default {
			return
		}
		if kongAccepted {
			p.declareKong(newCard, c.KongSize, true)
			return
		}
	}

	// pung to kong
	if p.FixedCards.Count(newCard) == c.MeldSize && !p.isWaiting {
		kongAccepted, isValidInput := false, false
		for {
			if isValidInput {
				break
			}
			kongAccepted, isValidInput = p.confirmMeld(newCard, false, false)
		}
		if p.result != result.Default {
			return
		}
		if kongAccepted {
			p.upgradePungToKong(newCard)
			return
		}
	}

	for {
		if len(*p.PlayerCards)+len(*p.FixedCards)-p.KongCount*c.KongExtraTiles == c.InitialHandSize {
			break
		} else if len(*p.PlayerCards)+len(*p.FixedCards)-p.KongCount*c.KongExtraTiles > c.WinningTileCount {
			panic("Player card count error")
		}

		p.handleDiscard(newCard)
	}
}

func (p *Player) handleDiscard(drawnTile c.Card) {
	if p.isWaiting {
		p.applyAction(replay.Action{Kind: "discard_drawn"}, drawnTile)
		return
	}
	action, id := p.decide("discard", drawnTile, p.t.T(language.KeyPlayACard))
	p.logger.Debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayACardInputInfo), action.Input()), p.counter)
	accepted := false
	if action.Kind == "invalid" {
		p.info(p.t.T(language.KeyUnrecognizedInput))
	} else {
		accepted = p.applyAction(action, drawnTile)
	}
	p.decided(id, accepted)
}

func (p *Player) applyAction(action replay.Action, drawnTile c.Card) bool {
	if action.Kind == "discard_drawn" {
		if drawnTile == "" {
			p.info(p.t.T(language.KeyUnrecognizedInput))
			return false
		}
		p.PlayerCards.RemoveTiles(drawnTile)
		if p.isWaiting {
			p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyAutoPlay), drawnTile))
			p.sleep(p.sleepTime)
		}
		p.infoCards(p.PlayerCards)
		p.emit("discard", drawnTile)
		return true
	}
	switch action.Kind {
	case "cancel":
		p.info(p.t.T(language.KeyEnd))
		p.result = result.Cancel
		p.PlayerCards.PopLast()
		p.ScoreItems = &score.Items{}
		return true
	case "debug_win":
		p.result = result.Win
		p.info(p.t.T(language.KeyEnterScript))
		p.PlayerCards.PopLast()
		p.printResult()
		p.ScoreItems = &score.Items{}
		return true
	case "toggle_speed":
		// Speed is a presentation setting, never persisted by game actions.
		if p.sleepTime != 0 {
			p.sleepTime = 0
			p.info(p.t.T(language.KeyEnterScript))
		} else {
			p.sleepTime = utils.SleepTime
			p.info(p.t.T(language.KeyCancelScript))
		}
		return true
	}
	if action.Kind == "declare_waiting" {
		tile := action.Tile
		if p.PlayerCards.Remove(tile) != nil {
			p.badDiscard(tile)
			return false
		}
		waiting, tiles := p.CheckIsWaiting()
		if !waiting {
			p.PlayerCards.Add(tile)
			p.PlayerCards.Sort()
			p.info(p.t.T(language.KeyCannotWaiting))
			p.infoCards(p.PlayerCards)
			return false
		}
		p.isWaiting = true
		p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerWaiting), tiles.Display()))
		p.recordWaiting(len(*tiles))
		p.infoCards(p.PlayerCards)
		p.emit("discard", tile)
		return true
	}
	card := action.Tile
	if p.PlayerCards.Remove(card) != nil {
		p.badDiscard(card)
		return false
	}
	p.infoCards(p.PlayerCards)
	p.emit("discard", card)
	return true
}

func (p *Player) badDiscard(card c.Card) {
	if p.FixedCards.Contains(card) {
		p.info(p.t.T(language.KeyCannotPlayCard))
	} else if c.CardsList.Contains(card) {
		p.info(p.t.T(language.KeyWithoutThisCard))
	} else {
		p.info(p.t.T(language.KeyUnrecognizedInput))
	}
	p.infoCards(p.PlayerCards)
}

func (p *Player) playRobotTurn() {
	p.sleep(p.sleepTime)
	if p.counter >= p.limit {
		p.result = result.Lose
		p.recordOutcome(false, false)
		p.printResult()
	} else {
		card := p.CardsPool.PopLast()
		p.emit("robot_discard", card)
		p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyRobotPlayACard), card))

		p.PlayerCards.Add(card)
		// Keep gameplay/replay tile ordering explicit; analysis itself is pure.
		p.PlayerCards.Sort()
		analysis := p.analyzeWin()
		if analysis.IsWin() {
			p.recordWinningPatterns(analysis)
			p.debug(p.t.T(language.KeyUserWin))
			p.sleep(p.sleepTime)
			p.result = result.Win
			p.recordOutcome(true, false)
			p.printResult()
			return
		}
		p.PlayerCards.RemoveTiles(card)

		if p.PlayerCards.Count(card) == c.MeldSize && !p.isWaiting {
			meldAccepted, inputValid := false, false
			for {
				if inputValid {
					break
				}
				meldAccepted, inputValid = p.confirmMeld(card, true, false)
			}
			if p.result != result.Default {
				return
			}
			if meldAccepted {
				p.declareKong(card, c.MeldSize, false)
				p.changeCondition()
			}
		}

		if p.PlayerCards.Count(card) >= c.PairSize && !p.isWaiting {
			meldAccepted, inputValid := false, false
			for {
				if inputValid {
					break
				}

				meldAccepted, inputValid = p.confirmMeld(card, true, true)
			}
			if p.result != result.Default {
				return
			}
			if meldAccepted {
				p.claimPung(card)
				p.changeCondition()
			}
		}
	}
}

func (p *Player) confirmMeld(card c.Card, showHand bool, isPung bool) (bool, bool) {
	if showHand {
		p.infoCards(p.PlayerCards)
	}
	var prompt, inputLabel string
	if isPung {
		prompt, inputLabel = p.t.T(language.KeyAskIfPung), p.t.T(language.KeyAskPlayerPungInputInfo)
	} else {
		prompt, inputLabel = p.t.T(language.KeyAskIfKong), p.t.T(language.KeyAskPlayerKongInputInfo)
	}

	kind := "confirm_kong"
	if isPung {
		kind = "confirm_pung"
	}
	action, id := p.decide(kind, card, fmt.Sprintf("%s %s", prompt, card))
	p.logger.Debug(fmt.Sprintf("%s: %s", inputLabel, action.Input()), p.counter)
	valid := action.Kind == "accept" || action.Kind == "decline"
	p.decided(id, valid)
	if !valid {
		p.info(p.t.T(language.KeyUnrecognizedInput))
	}
	return action.Kind == "accept", valid
}

func (p *Player) claimPung(card c.Card) {
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerWaiting), card))
	p.FixedCards.Add(card, card, card)
	p.PlayerCards.RemoveCopies(card, c.MeldSize-1)
	p.isClear = false
	p.emit("pung", card)
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerPung), p.FixedCards.Display()))
	p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyFixedCard), p.FixedCards.Display()))
	p.infoCards(p.PlayerCards)

	for {
		if len(*p.PlayerCards)+len(*p.FixedCards)-p.KongCount*c.KongExtraTiles == c.InitialHandSize {
			break
		}
		p.handleDiscard("")
	}
}

func (p *Player) upgradePungToKong(card c.Card) {
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerAddKong), card))
	p.FixedCards.Add(card)
	p.PlayerCards.RemoveTiles(card)
	p.KongCount++
	p.recordKong(false)
	p.emit("upgrade_kong", card)
	p.drawAndPlay(true)
}

func (p *Player) declareKong(card c.Card, concealedTilesToRemove int, concealed bool) {
	p.debug(fmt.Sprintf("%s:, %s", p.t.T(language.KeyPlayerKong), card))
	p.KongCount++
	p.recordKong(concealed)
	if !concealed {
		p.isClear = false
	}
	p.FixedCards.Add(card, card, card, card)
	p.PlayerCards.RemoveCopies(card, concealedTilesToRemove)
	p.emit("kong", card)
	p.drawAndPlay(false)
}

func (p *Player) printResult() {
	switch p.result {
	case result.Win:
		p.info(p.t.T(language.KeyUWin))
	case result.Equal:
		p.info(p.t.T(language.KeyDraw))
	case result.Lose:
		p.info(p.t.T(language.KeyULose))
		p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyRobotsCard), p.RobotCards.Display()))
	default:
		panic("Invalid result value")
	}
}

// JudgeIsWin only answers whether the current tiles form a legal winning hand.
// It does not change tile order, score items, logs, or any other player state.
func (p *Player) JudgeIsWin() bool {
	return p.analyzeWin().IsWin()
}

func (p *Player) analyzeWin() hand.Analysis {
	if p.PlayerCards == nil || p.FixedCards == nil {
		return hand.Analysis{}
	}
	return hand.Analyze(*p.PlayerCards, *p.FixedCards, p.KongCount)
}

func (p *Player) CheckIsWaiting() (bool, *c.Cards) {
	winningTiles := &c.Cards{}
	handSize := len(*p.PlayerCards)
	originalHand := make(c.Cards, handSize)
	copy(originalHand, *p.PlayerCards)
	for _, tile := range c.CardsList {
		if p.PlayerCards.Count(tile)+p.FixedCards.Count(tile) >= c.CopiesPerTile {
			continue
		}
		p.PlayerCards.Add(tile)
		if p.JudgeIsWin() {
			winningTiles.Add(tile)
		}
		*p.PlayerCards = make(c.Cards, handSize)
		copy(*p.PlayerCards, originalHand)
	}
	if len(*winningTiles) != 0 {
		winningTiles.Sort()
		p.debug(p.t.T(language.KeyMatchWaiting))
		return true, winningTiles
	} else {
		p.debug(p.t.T(language.KeyNotMatchWaiting))
		return false, winningTiles
	}
}

func (p *Player) changeCondition() {
	if p.condition == rules.InitialPlayerTurnParity {
		p.condition = rules.ClaimedPlayerTurnParity
	} else {
		p.condition = rules.InitialPlayerTurnParity
	}
}

func (p *Player) Main() (res result.Result, amount int, counter int) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if abort, ok := recovered.(inputAbort); ok {
				p.err = abort.err
				p.result = result.Cancel
				res, amount, counter = result.Cancel, 0, p.counter
			} else {
				panic(recovered)
			}
		}
	}()
	for {
		if p.result != result.Default {
			break
		}

		if p.counter >= rules.MaxTurns {
			panic(inputAbort{fmt.Errorf("turn limit exceeded")})
		}
		if len(*p.CardsPool) <= rules.DrawWallThreshold {
			p.result = result.Equal
			*p.ScoreItems = score.Items{}
			break
		}
		if p.counter%rules.TurnCycleLength == p.condition {
			p.counter++
			p.drawAndPlay(true)
		} else {
			p.counter++
			p.playRobotTurn()
		}
		p.emit("turn_end", "")
	}

	p.counter++
	amount, scoreDetail := p.settleScore()
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyLastCard), *p.CardsPool))
	p.info(fmt.Sprintf("%s: %v", p.t.T(language.KeyThisScore), amount))
	if amount != 0 {
		if view := ui.RenderScoreDetail(p.t, scoreDetail); view != "" {
			p.info(view)
		}
	}
	p.observe(replay.Frame{Kind: "settlement", State: p.Snapshot(), Settlement: &replay.Settlement{Result: p.result, Score: amount, Counter: p.counter, Detail: scoreDetail}})
	return p.result, amount, p.counter
}
