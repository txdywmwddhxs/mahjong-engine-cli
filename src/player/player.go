package player

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
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
	view := ui.RenderCards(cards)
	p.logger.Info(view, p.counter)
	if p.ui != nil {
		p.ui.Info(view)
	}
}

func (p *Player) debug(msg string) {
	p.logger.Debug(msg, p.counter)
}

func (p *Player) initCardsPool() {
	p.CardsPool = c.InitCardsPool()
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyInitCardPoolDone), *p.CardsPool))
}

func (p *Player) initPlayerCards() {
	*p.PlayerCards = (*p.CardsPool)[:13]
	p.PlayerCards.Sort()
	*p.CardsPool = (*p.CardsPool)[13:]
	p.debug(fmt.Sprintf("%s: %v", p.t.T(language.KeyInitPlayerCardDone), p.PlayerCards.Beauty()))
}

func (p *Player) initRobotCards() {
	res := c.Cards{}
	//gen pair card group
	pairGroup := (*p.CardsPool)[0]
	res = append(res, pairGroup, pairGroup)
	p.CardsPool.DelM(pairGroup, pairGroup)

	nTripleCardGroup := utils.RandInt(0, 2)
	nSequenceGroup := 4 - nTripleCardGroup
	tripleCardGroup := make([]c.Cards, 0, nTripleCardGroup)
	sequenceGroup := make([]c.Cards, 0, nSequenceGroup)

	// gen triple card group
	if nTripleCardGroup != 0 {
		for {
			// seed: decide card kind
			seed := utils.RandInt(1, 4)
			var cn c.Card
			switch seed {
			case 1:
				k := utils.MajorCardList[seed-1]
				v := utils.RandInt(1, 9)
				cn = c.Card(string(k) + strconv.Itoa(v))
			case 2:
				k := utils.MajorCardList[seed-1]
				v := utils.RandInt(1, 9)
				cn = c.Card(string(k) + strconv.Itoa(v))
			case 3:
				k := utils.MajorCardList[seed-1]
				v := utils.RandInt(1, 9)
				cn = c.Card(string(k) + strconv.Itoa(v))
			default:
				seed := utils.RandInt(0, 6)
				cn = utils.MinorCardList[seed]
			}

			if p.CardsPool.CountC(cn) >= 3 {
				tripleCardGroup = append(tripleCardGroup, c.Cards{cn, cn, cn})
				p.CardsPool.DelT(cn, 3)
			}

			if len(tripleCardGroup) == nTripleCardGroup {
				break
			}
			if len(tripleCardGroup) > nTripleCardGroup {
				panic("generate tri error")
			}
		}
	}

	// gen sequence card group
	if nSequenceGroup != 0 {
		for {
			seed := utils.RandInt(0, 2)
			// card kind
			first := utils.MajorCardList[seed]
			// card num
			second := utils.RandInt(2, 8)
			fcn := first + c.Card(strconv.Itoa(second))
			scn := first + c.Card(strconv.Itoa(second-1))
			ncn := first + c.Card(strconv.Itoa(second+1))
			if p.CardsPool.ContainsC(fcn) && p.CardsPool.ContainsC(scn) && p.CardsPool.ContainsC(ncn) {
				sequenceGroup = append(sequenceGroup, c.Cards{fcn, scn, ncn})
				p.CardsPool.DelM(ncn, scn, fcn)
			}

			if len(sequenceGroup) == nSequenceGroup {
				break
			}
			if len(sequenceGroup) > nSequenceGroup {
				panic("generate seq error")
			}
		}
	}

	for _, i := range tripleCardGroup {
		res = append(res, i...)
	}
	for _, j := range sequenceGroup {
		res = append(res, j...)
	}

	res.Sort()
	*p.RobotCards = res
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyInitRobotCardDone), p.RobotCards.Beauty()))
}

func (p *Player) generateLimit() {
	goodCardLimit, _ := strconv.Atoi(utils.RateOfRobotGoodCard[:len(utils.RateOfRobotGoodCard)-1])
	babCardLimit, _ := strconv.Atoi(utils.RateOfRobotBadCard[:len(utils.RateOfRobotBadCard)-1])
	babCardLimit += goodCardLimit
	p.debug(fmt.Sprintf("%s: %d%%, %d%%, %d%%",
		p.t.T(language.KeyRateOfRobotsCard), goodCardLimit, babCardLimit-goodCardLimit, 100-goodCardLimit))
	seed := utils.RandInt(0, 99)
	p.debug(fmt.Sprintf("%s: %d", p.t.T(language.KeyRandomNumberSeed), seed))
	if seed <= goodCardLimit {
		p.limit = utils.RandInt(utils.GoodCardLimitRangeStart, utils.GoodCardLimitRangeEnd)
	} else if seed <= babCardLimit {
		p.limit = utils.RandInt(utils.GoodCardLimitRangeEnd, utils.BadCardLimitRangeStart)
	} else {
		p.limit = utils.RandInt(utils.BadCardLimitRangeStart, utils.BadCardLimitRangeEnd)
	}
}

func (p *Player) init() {
	p.generateLimit()
	p.debug(fmt.Sprintf("%s: %d", p.t.T(language.KeyLimit), p.limit))
	p.initCardsPool()
	p.initRobotCards()
	p.initPlayerCards()
}

func (p *Player) playerPlayCard(judge bool) {
	newCard := (*p.CardsPool)[0]
	p.PlayerCards.NewC(newCard)
	p.CardsPool.DelM(newCard)
	p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyGetNewCard), newCard))

	if p.counter == 1 && judge {
		p.PlayerCards.Sort()
		handledCards := c.Cards{}

	search:
		for _, card := range *p.PlayerCards {
			if handledCards.ContainsC(card) {
				continue
			}
			handledCards.NewC(card)
			if p.PlayerCards.CountC(card) == 4 {
				userConfirmGroupFour := false
				isValidInput := false
				for {
					p.infoCards(p.PlayerCards)
					userConfirmGroupFour, isValidInput = p.userInput(card, false, false)
					if isValidInput && p.result == result.Default {
						break
					}
				}
				if userConfirmGroupFour {
					p.userGroupFour(card, 4, true)
					goto search
				}
			}
		}
	}
	if p.result != result.Default {
		return
	}

	if p.sleepTime != 0 && p.counter != 1 {
		time.Sleep(p.sleepTime + 350*time.Millisecond)
	}
	p.PlayerCards.Sort()
	if len(*p.FixedCards) != 0 {
		p.FixedCards.Sort()
		p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyFixedCard), p.FixedCards.Beauty()))
	}
	p.infoCards(p.PlayerCards)

	isWin := p.JudgeIsWin(true)
	if isWin {
		p.debug(p.t.T(language.KeyScoreWin))
		p.result = result.Win
		p.ScoreItem.AddI(score.Win)
		if p.counter <= utils.CountOfWinAward {
			p.ScoreItem.AddI(score.Counter)
		}
		if p.isClear {
			p.ScoreItem.AddI(score.OwnDraw)
		}
		p.printResult()
		return
	}
	p.debug(p.t.T(language.KeyUserNotWin))

	if p.PlayerCards.CountC(newCard) == 4 && !p.isWaiting {
		userConfirmGroupFour, isValidInput := false, false
		for {
			if isValidInput {
				break
			}
			userConfirmGroupFour, isValidInput = p.userInput(newCard, false, false)
		}
		if p.result != result.Default {
			return
		}
		if userConfirmGroupFour {
			p.userGroupFour(newCard, 4, true)
			return
		}
	}

	// pung to kong
	if p.FixedCards.CountC(newCard) == 3 && !p.isWaiting {
		userConfirm, isValidInput := false, false
		for {
			if isValidInput {
				break
			}
			userConfirm, isValidInput = p.userInput(newCard, false, false)
		}
		if p.result != result.Default {
			return
		}
		if userConfirm {
			p.changeGroupThreeToFour(newCard)
			return
		}
	}

	for {
		if len(*p.PlayerCards)+len(*p.FixedCards)-p.GroupFourCount == 13 {
			break
		} else if len(*p.PlayerCards)+len(*p.FixedCards)-p.GroupFourCount >= 15 {
			panic("Player card count error")
		}

		p.outCard(newCard)
	}
}

func (p *Player) outCard(nc c.Card) {
	var card c.Card
	if !p.isWaiting {
		input, _ := p.ui.PromptInfo(p.t.T(language.KeyPlayACard))
		p.logger.Debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayACardInputInfo), input), p.counter)
		input = strings.ToUpper(strings.TrimSpace(input))
		card = c.Card(input)
	} else {
		card = ""
	}

	if card == "" {
		if nc == "" {
			p.info(p.t.T(language.KeyUnrecognizedInput))
			return
		}
		p.PlayerCards.DelM(nc)
		if p.isWaiting {
			p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyAutoPlay), nc))
			time.Sleep(p.sleepTime)
		}
		p.infoCards(p.PlayerCards)
		return
	}

	if card == "END" || card == "ENDL" {
		p.info(p.t.T(language.KeyEnd))
		p.result = result.Cancel
		p.PlayerCards.DelL()
		p.ScoreItem = &score.Items{}
	} else if card == "WHOSYOURDADDY" {
		p.result = result.Win
		p.info(p.t.T(language.KeyEnterScript))
		p.PlayerCards.DelL()
		p.printResult()
		p.ScoreItem = &score.Items{}
	} else if card == "WARPTEN" {
		utils.ChangeQuickMode()
		if p.sleepTime != 0 {
			p.sleepTime = 0
			p.info(p.t.T(language.KeyEnterScript))
		} else {
			p.sleepTime = utils.SleepTime
			p.info(p.t.T(language.KeyCancelScript))
		}
	} else if strings.HasSuffix(string(card), "TING") || strings.HasSuffix(string(card), "TT") {
		inputArr := strings.Split(string(card), " ")
		if len(inputArr) < 2 {
			p.info(p.t.T(language.KeyCannotWaitingPlayCard))
			p.infoCards(p.PlayerCards)
			return
		}
		prepareToPlayCard := c.Card(strings.TrimSpace(inputArr[0]))
		err := p.PlayerCards.DelE(prepareToPlayCard)
		if err != nil {
			if p.FixedCards.ContainsC(prepareToPlayCard) {
				p.info(p.t.T(language.KeyCannotPlayCard))
			} else if c.CardsList.ContainsC(prepareToPlayCard) {
				p.info(p.t.T(language.KeyWithoutThisCard))
			} else {
				p.info(p.t.T(language.KeyUnrecognizedInput))
			}
			return
		}
		isWaiting, res := p.CheckIsWaiting()
		if isWaiting {
			p.isWaiting = true
			p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerWaiting), res.Beauty()))
			p.ScoreItem.AddI(score.Waiting)
			if len(*res) == 1 {
				p.ScoreItem.AddI(score.Single)
			}
		} else {
			p.PlayerCards.NewC(prepareToPlayCard)
			p.PlayerCards.Sort()
			p.info(p.t.T(language.KeyCannotWaiting))
		}
		p.infoCards(p.PlayerCards)
	} else {
		err := p.PlayerCards.DelE(card)
		if err != nil {
			if p.FixedCards.ContainsC(card) {
				p.info(p.t.T(language.KeyCannotPlayCard))
			} else if c.CardsList.ContainsC(card) {
				p.info(p.t.T(language.KeyWithoutThisCard))
			} else {
				p.info(p.t.T(language.KeyUnrecognizedInput))
			}
		}
		p.infoCards(p.PlayerCards)
	}
}

func (p *Player) robotPlayCard() {
	time.Sleep(p.sleepTime)
	if p.counter >= p.limit {
		p.result = result.Lose
		p.ScoreItem.AddI(score.Lose)
		if p.counter <= utils.CountOfWinAward {
			p.ScoreItem.AddI(score.Counter)
		}
		p.printResult()
	} else {
		card := p.CardsPool.DelL()
		p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyRobotPlayACard), card))

		p.PlayerCards.NewC(card)
		if p.JudgeIsWin(true) {
			p.debug(p.t.T(language.KeyUserWin))
			time.Sleep(p.sleepTime)
			p.result = result.Win
			p.ScoreItem.AddI(score.Win)
			if p.counter <= utils.CountOfWinAward {
				p.ScoreItem.AddI(score.Counter)
			}
			p.printResult()
			return
		}
		p.PlayerCards.DelM(card)

		if p.PlayerCards.CountC(card) == 3 && !p.isWaiting {
			uir, eff := false, false
			for {
				if eff {
					break
				}
				uir, eff = p.userInput(card, true, false)
			}
			if p.result != result.Default {
				return
			}
			if uir {
				p.userGroupFour(card, 3, false)
				p.changeCondition()
			}
		}

		if p.PlayerCards.CountC(card) >= 2 && !p.isWaiting {
			uir, eff := false, false
			for {
				if eff {
					break
				}

				uir, eff = p.userInput(card, true, true)
			}
			if p.result != result.Default {
				return
			}
			if uir {
				p.userGroup(card)
				p.changeCondition()
			}
		}
	}
}

func (p *Player) userInput(card c.Card, listCard bool, pung bool) (bool, bool) {
	if listCard {
		p.infoCards(p.PlayerCards)
	}
	var askPlayer, inputInfo string
	if pung {
		askPlayer, inputInfo = p.t.T(language.KeyAskIfPung), p.t.T(language.KeyAskPlayerPungInputInfo)
	} else {
		askPlayer, inputInfo = p.t.T(language.KeyAskIfKong), p.t.T(language.KeyAskPlayerKongInputInfo)
	}

	info, _ := p.ui.PromptInfo(fmt.Sprintf("%s %s", askPlayer, card))
	p.logger.Debug(fmt.Sprintf("%s: %s", inputInfo, info), p.counter)
	info = strings.ToUpper(strings.TrimSpace(info))
	if info == "Y" || info == "YES" {
		return true, true
	} else if info == "N" || info == "NO" {
		return false, true
	} else {
		p.info(p.t.T(language.KeyUnrecognizedInput))
		return false, false
	}
}

func (p *Player) userGroup(card c.Card) {
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerWaiting), card))
	p.FixedCards.NewC(card, card, card)
	p.PlayerCards.DelT(card, 2)
	p.isClear = false
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerPung), p.FixedCards.Beauty()))
	p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyFixedCard), p.FixedCards.Beauty()))
	p.infoCards(p.PlayerCards)

	for {
		if len(*p.PlayerCards)+len(*p.FixedCards)-p.GroupFourCount == 13 {
			break
		}
		p.outCard("")
	}
}

func (p *Player) changeGroupThreeToFour(card c.Card) {
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyPlayerAddKong), card))
	p.FixedCards.NewC(card)
	p.PlayerCards.DelM(card)
	p.GroupFourCount++
	p.ScoreItem.AddI(score.ExposedKong)
	p.playerPlayCard(true)
}

func (p *Player) userGroupFour(card c.Card, nDel int, self bool) {
	p.debug(fmt.Sprintf("%s:, %s", p.t.T(language.KeyPlayerKong), card))
	p.GroupFourCount++
	if self {
		p.ScoreItem.AddI(score.ConcealedKong)
	} else {
		p.ScoreItem.AddI(score.ExposedKong)
		p.isClear = false
	}
	p.FixedCards.NewC(card, card, card, card)
	p.PlayerCards.DelT(card, nDel)
	p.playerPlayCard(false)
}

func (p *Player) printResult() {
	switch p.result {
	case result.Win:
		p.info(p.t.T(language.KeyUWin))
	case result.Equal:
		p.info(p.t.T(language.KeyDraw))
	case result.Lose:
		p.info(p.t.T(language.KeyULose))
		p.info(fmt.Sprintf("%s: %s", p.t.T(language.KeyRobotsCard), p.RobotCards.Beauty()))
	default:
		panic("Invalid result value")
	}
}

func (p *Player) JudgeIsWin(final bool) bool {
	p.PlayerCards.Sort()
	length := len(*p.PlayerCards)
	var cards, originCards = make(c.Cards, length), make(c.Cards, length)
	copy(originCards, *p.PlayerCards)
	copy(cards, originCards)
	fixedCardsGroupLength := (len(*p.FixedCards) - p.GroupFourCount) / 3
	if len(cards)+len(*p.FixedCards)-p.GroupFourCount != 14 {
		return false
	}
	complete := make([]c.Cards, 0, 5)
	alreadyCards := c.Cards{}
	for _, card := range originCards {
		if alreadyCards.ContainsC(card) {
			continue
		}
		if cards.CountC(card) >= 2 {
			cards.DelT(card, 2)
			complete = append(complete, c.Cards{card, card})
		} else {
			continue
		}
		for i := 0; i < length; i++ {
			if len(cards) == 0 {
				break
			}
			headCard := cards[0]
			if cards.CountC(headCard) >= 3 {
				complete = append(complete, c.Cards{headCard, headCard, headCard})
				cards.DelT(headCard, 3)
			} else if len(headCard) == 2 { // head card is major card
				nj := c.Card(headCard[0]) + c.Card(rune(headCard[1]+1))
				nnj := c.Card(headCard[0]) + c.Card(rune(headCard[1]+2))
				if cards.ContainsC(nj) && cards.ContainsC(nnj) {
					complete = append(complete, c.Cards{headCard, nj, nnj})
					cards.DelM(headCard, nj, nnj)
				}
			}
		}
		if len(complete)+fixedCardsGroupLength == 5 {
			if final {
				p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyCardAnalysisResult), complete))
				if len(complete) >= 3 {
					// check if continuous line
					for _, k := range utils.MajorCardList {
						isContinuousLine := true
						for _, n := range utils.ContinuousLineJudge {
							in := false
							for _, group := range complete {
								if group.ContainsC(k+n[0]) && group.ContainsC(k+n[1]) && group.ContainsC(k+n[2]) {
									in = true
								}
							}
							if !in {
								isContinuousLine = false
								break
							}
						}
						if isContinuousLine {
							p.ScoreItem.AddI(score.ContinuousLine)
							break
						}
					}
				}
			}
			return true
		} else {
			cards = make(c.Cards, length)
			copy(cards, originCards)

			complete = make([]c.Cards, 0, 5)
			alreadyCards.NewC(card)
		}
	}

	isSevenPairs := true
	isThirteenOne := true

	hasPair := false
	alreadyCards = c.Cards{}
	for _, i := range cards {
		if alreadyCards.ContainsC(i) {
			continue
		}
		alreadyCards.NewC(i)
		num := cards.CountC(i)
		if isSevenPairs && num != 2 && num != 4 {
			isSevenPairs = false
		}

		if isThirteenOne && !utils.ThirteenOneCardList.ContainsC(i) {
			isThirteenOne = false
		}

		if isThirteenOne && num != 1 && num != 2 {
			isThirteenOne = false
		} else if isThirteenOne && num == 2 {
			if !hasPair {
				hasPair = true
			} else {
				isThirteenOne = false
			}
		}
		if !isThirteenOne && !isSevenPairs {
			break
		}
	}
	if isThirteenOne && hasPair && len(*p.FixedCards) == 0 {
		if final {
			p.debug(p.t.T(language.KeyMatchThirteenOne))
			p.ScoreItem.AddI(score.ThirteenOne)
		}
		return true
	}

	if isSevenPairs && len(*p.FixedCards) == 0 {
		if final {
			p.debug(p.t.T(language.KeyMatchSevenPairs))
			p.ScoreItem.AddI(score.SevenPairs)
		}
		return true
	}
	return false
}

func (p *Player) CheckIsWaiting() (bool, *c.Cards) {
	res := &c.Cards{}
	length := len(*p.PlayerCards)
	originCards := make(c.Cards, length)
	copy(originCards, *p.PlayerCards)
	for _, i := range c.CardsList {
		p.PlayerCards.NewC(i)
		if p.JudgeIsWin(false) {
			res.NewC(i)
		}
		*p.PlayerCards = make(c.Cards, length)
		copy(*p.PlayerCards, originCards)
	}
	if len(*res) != 0 {
		res.Sort()
		p.debug(p.t.T(language.KeyMatchWaiting))
		return true, res
	} else {
		p.debug(p.t.T(language.KeyNotMatchWaiting))
		return false, res
	}
}

func (p *Player) changeCondition() {
	if p.condition == 0 {
		p.condition = 1
	} else {
		p.condition = 1
	}
}

func (p *Player) getMissCard() {
	if p.result != result.Win {
		return
	}
	fullPlayCard := append(*p.PlayerCards, *p.FixedCards...)
	cardKind := make(c.Cards, 0, 3)
	for _, i := range fullPlayCard {
		if len(i) == 2 && !cardKind.ContainsC(c.Card(i[0])) {
			cardKind.NewC(c.Card(i[0]))
		}
	}
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyCardKind), cardKind))
	count := 3 - len(cardKind)
	if count == 1 {
		p.ScoreItem.AddI(score.MissOneKind)
	} else if count >= 2 {
		p.ScoreItem.AddI(score.MissTwoKind)
	}
}

func (p *Player) Main() (result.Result, int, int) {
	for {
		if p.result != result.Default {
			break
		}

		if len(*p.CardsPool) <= 14 {
			p.result = result.Equal
			*p.ScoreItem = score.Items{}
			break
		}
		if p.counter%2 == p.condition {
			p.counter++
			p.playerPlayCard(true)
		} else {
			p.counter++
			p.robotPlayCard()
		}
	}

	p.counter++
	p.getMissCard()
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyScoreItem), *p.ScoreItem))
	amount, scoreDetail := score.GetScore(p.ScoreItem, p.counter-1)
	p.debug(fmt.Sprintf("%s: %s", p.t.T(language.KeyLastCard), *p.CardsPool))
	p.info(fmt.Sprintf("%s: %v", p.t.T(language.KeyThisScore), amount))
	if amount != 0 {
		if view := ui.RenderScoreDetail(p.t, scoreDetail); view != "" {
			p.info(view)
		}
	}
	return p.result, amount, p.counter
}
