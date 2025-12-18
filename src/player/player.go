package player

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	log "github.com/txdywmwddhxs/mahjong-engine-cli/src/clog"
	l "github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

var scanner = bufio.NewScanner(os.Stdin)

func (p *Player) initCardsPool() {
	p.CardsPool = c.InitCardsPool()
	p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.InitCardPoolDone), *p.CardsPool), p.counter)
}

func (p *Player) initPlayerCards() {
	*p.PlayerCards = (*p.CardsPool)[:13]
	p.PlayerCards.Sort()
	*p.CardsPool = (*p.CardsPool)[13:]
	p.logger.DebugS(fmt.Sprintf("%s: %v", log.Trans(l.InitPlayerCardDone), p.PlayerCards.Beauty()), p.counter)
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
	p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.InitRobotCardDone), p.RobotCards.Beauty()), p.counter)
}

func (p *Player) generateLimit() {
	goodCardLimit, _ := strconv.Atoi(utils.RateOfRobotGoodCard[:len(utils.RateOfRobotGoodCard)-1])
	babCardLimit, _ := strconv.Atoi(utils.RateOfRobotBadCard[:len(utils.RateOfRobotBadCard)-1])
	babCardLimit += goodCardLimit
	p.logger.DebugS(fmt.Sprintf("%s: %d%%, %d%%, %d%%",
		log.Trans(l.RateOfRobotsCard), goodCardLimit, babCardLimit-goodCardLimit, 100-goodCardLimit), p.counter)
	seed := utils.RandInt(0, 99)
	p.logger.DebugS(fmt.Sprintf("%s: %d", log.Trans(l.RandomNumberSeed), seed), p.counter)
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
	p.logger.DebugS(fmt.Sprintf("%s: %d", log.Trans(l.Limit), p.limit), p.counter)
	p.initCardsPool()
	p.initRobotCards()
	p.initPlayerCards()
}

func (p *Player) playerPlayCard(judge bool) {
	newCard := (*p.CardsPool)[0]
	p.PlayerCards.NewC(newCard)
	p.CardsPool.DelM(newCard)
	p.logger.InfoS(fmt.Sprintf("%s: %s", log.Trans(l.GetNewCard), newCard), p.counter)

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
					p.logger.InfoC(p.PlayerCards, p.counter)
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
		p.logger.InfoS(fmt.Sprintf("%s: %s", log.Trans(l.FixedCard), p.FixedCards.Beauty()), p.counter)
	}
	p.logger.InfoC(p.PlayerCards, p.counter)

	isWin := p.JudgeIsWin(true)
	if isWin {
		p.logger.DebugS(log.Trans(l.Win), p.counter)
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
	p.logger.DebugS(log.Trans(l.UserNotWin), p.counter)

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
		fmt.Printf("INFO: %s: ", log.Trans(l.PlayACard))
		scanner.Scan()
		input := scanner.Text()
		p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.PlayACardInputInfo), input), p.counter)
		input = strings.ToUpper(input)
		input = strings.TrimSpace(input)
		card = c.Card(input)
	} else {
		card = ""
	}

	if card == "" {
		if nc == "" {
			p.logger.InfoS(log.Trans(l.UnrecognizedInput), p.counter)
			return
		}
		p.PlayerCards.DelM(nc)
		if p.isWaiting {
			p.logger.InfoS(fmt.Sprintf("%s: %s", log.Trans(l.AutoPlay), nc), p.counter)
			time.Sleep(p.sleepTime)
		}
		p.logger.InfoC(p.PlayerCards, p.counter)
		return
	}

	if card == "END" || card == "ENDL" {
		p.logger.InfoS(log.Trans(l.End), p.counter)
		p.result = result.Cancel
		p.PlayerCards.DelL()
		p.ScoreItem = &score.Items{}
	} else if card == "WHOSYOURDADDY" {
		p.result = result.Win
		p.logger.InfoS(log.Trans(l.EnterScript), p.counter)
		p.PlayerCards.DelL()
		p.printResult()
		p.ScoreItem = &score.Items{}
	} else if card == "WARPTEN" {
		utils.ChangeQuickMode()
		if p.sleepTime != 0 {
			p.sleepTime = 0
			p.logger.InfoS(log.Trans(l.EnterScript), p.counter)
		} else {
			p.sleepTime = utils.SleepTime
			p.logger.InfoS(log.Trans(l.CancelScript), p.counter)
		}
	} else if strings.HasSuffix(string(card), "TING") || strings.HasSuffix(string(card), "TT") {
		inputArr := strings.Split(string(card), " ")
		if len(inputArr) < 2 {
			p.logger.InfoS(log.Trans(l.CannotWaitingPlayCard), p.counter)
			p.logger.InfoC(p.PlayerCards, p.counter)
			return
		}
		prepareToPlayCard := c.Card(strings.TrimSpace(inputArr[0]))
		err := p.PlayerCards.DelE(prepareToPlayCard)
		if err != nil {
			if p.FixedCards.ContainsC(prepareToPlayCard) {
				p.logger.InfoS(log.Trans(l.CannotPlayCard), p.counter)
			} else if c.CardsList.ContainsC(prepareToPlayCard) {
				p.logger.InfoS(log.Trans(l.WithoutThisCard), p.counter)
			} else {
				p.logger.InfoS(log.Trans(l.UnrecognizedInput), p.counter)
			}
			return
		}
		isWaiting, res := p.CheckIsWaiting()
		if isWaiting {
			p.isWaiting = true
			p.logger.InfoS(fmt.Sprintf("%s: %s", log.Trans(l.PlayerWaiting), res.Beauty()), p.counter)
			p.ScoreItem.AddI(score.Waiting)
			if len(*res) == 1 {
				p.ScoreItem.AddI(score.Single)
			}
		} else {
			p.PlayerCards.NewC(prepareToPlayCard)
			p.PlayerCards.Sort()
			p.logger.InfoS(log.Trans(l.CannotWaiting), p.counter)
		}
		p.logger.InfoC(p.PlayerCards, p.counter)
	} else {
		err := p.PlayerCards.DelE(card)
		if err != nil {
			if p.FixedCards.ContainsC(card) {
				p.logger.InfoS(log.Trans(l.CannotPlayCard), p.counter)
			} else if c.CardsList.ContainsC(card) {
				p.logger.InfoS(log.Trans(l.WithoutThisCard), p.counter)
			} else {
				p.logger.InfoS(log.Trans(l.UnrecognizedInput), p.counter)
			}
		}
		p.logger.InfoC(p.PlayerCards, p.counter)
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
		p.logger.InfoS(fmt.Sprintf("%s: %s", log.Trans(l.RobotPlayACard), card), p.counter)

		p.PlayerCards.NewC(card)
		if p.JudgeIsWin(true) {
			p.logger.DebugS(log.Trans(l.UserWin), p.counter)
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
		p.logger.InfoC(p.PlayerCards, p.counter)
	}
	var askPlayer, inputInfo string
	if pung {
		askPlayer, inputInfo = log.Trans(l.AskIfPung), log.Trans(l.AskPlayerPungInputInfo)
	} else {
		askPlayer, inputInfo = log.Trans(l.AskIfKong), log.Trans(l.AskPlayerKongInputInfo)
	}

	fmt.Printf("INFO: %s %s: ", askPlayer, card)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	info := scanner.Text()
	p.logger.DebugS(fmt.Sprintf("%s: %s", inputInfo, info), p.counter)
	info = strings.ToUpper(info)
	if info == "Y" || info == "YES" {
		return true, true
	} else if info == "N" || info == "NO" {
		return false, true
	} else {
		p.logger.InfoS(log.Trans(l.UnrecognizedInput), p.counter)
		return false, false
	}
}

func (p *Player) userGroup(card c.Card) {
	p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.PlayerWaiting), card), p.counter)
	p.FixedCards.NewC(card, card, card)
	p.PlayerCards.DelT(card, 2)
	p.isClear = false
	p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.PlayerPung), p.FixedCards.Beauty()), p.counter)
	p.logger.InfoS(fmt.Sprintf("%s: %s", log.Trans(l.FixedCard), p.FixedCards.Beauty()), p.counter)
	p.logger.InfoC(p.PlayerCards, p.counter)

	for {
		if len(*p.PlayerCards)+len(*p.FixedCards)-p.GroupFourCount == 13 {
			break
		}
		p.outCard("")
	}
}

func (p *Player) changeGroupThreeToFour(card c.Card) {
	p.logger.DebugS(fmt.Sprintf("%s: %s", l.PlayerAddKong(), card), p.counter)
	p.FixedCards.NewC(card)
	p.PlayerCards.DelM(card)
	p.GroupFourCount++
	p.ScoreItem.AddI(score.ExposedKong)
	p.playerPlayCard(true)
}

func (p *Player) userGroupFour(card c.Card, nDel int, self bool) {
	p.logger.DebugS(fmt.Sprintf("%s:, %s", log.Trans(l.PlayerKong), card), p.counter)
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
		p.logger.InfoS(log.Trans(l.UWin), p.counter)
	case result.Equal:
		p.logger.InfoS(log.Trans(l.Draw), p.counter)
	case result.Lose:
		p.logger.InfoS(log.Trans(l.ULose), p.counter)
		p.logger.InfoS(fmt.Sprintf("%s: %s", log.Trans(l.RobotsCard), p.RobotCards.Beauty()), p.counter)
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
				p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.CardAnalysisResult), complete), p.counter)
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
			p.logger.DebugS(log.Trans(l.MatchThirteenOne), p.counter)
			p.ScoreItem.AddI(score.ThirteenOne)
		}
		return true
	}

	if isSevenPairs && len(*p.FixedCards) == 0 {
		if final {
			p.logger.DebugS(log.Trans(l.MatchSevenPairs), p.counter)
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
		p.logger.DebugS(log.Trans(l.MatchWaiting), p.counter)
		return true, res
	} else {
		p.logger.DebugS(log.Trans(l.NotMatchWaiting), p.counter)
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
	p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.CardKind), cardKind), p.counter)
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
	p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.ScoreItem), *p.ScoreItem), p.counter)
	amount, scoreDetail := score.GetScore(p.ScoreItem, p.counter-1)
	p.logger.DebugS(fmt.Sprintf("%s: %s", log.Trans(l.LastCard), *p.CardsPool), p.counter)
	p.logger.InfoS(fmt.Sprintf("%s: %v", log.Trans(l.ThisScore), amount), p.counter)
	if amount != 0 {
		p.logger.PrintScore(scoreDetail, p.counter)
	}
	return p.result, amount, p.counter
}
