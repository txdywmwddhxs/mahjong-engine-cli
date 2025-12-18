package score

import (
	"fmt"
	"math"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

const (
	WaitingMinScore     = 2
	WaitingMaxScore     = 2
	SingleMinScore      = 3
	SingleMaxScore      = 3
	CounterMaxScore     = 2
	CounterMinScore     = 2
	MissKindMinScore    = 3
	MissKindMaxScore    = 3
	OwnDrawTimes        = 1.0
	KongScore           = 2
	ThirteenOneScore    = 20
	SevenPairsScore     = 10
	ContinuousLineScore = 10
)

var ItemGeneratorMap = map[Item]func() int{
	ExposedKong:    getExposedKongScore,
	ConcealedKong:  getConcealedScore,
	Waiting:        getWaitingScore,
	MissOneKind:    getMissOneScore,
	MissTwoKind:    getMissTwoScore,
	Counter:        getCounterScore,
	SevenPairs:     getSevenPairsScore,
	ThirteenOne:    getThirteenOneScore,
	ContinuousLine: getContinuousLineScore,
	Single:         getSingleScore,
}

func GetScore(items *Items, n int) (int, map[Item]int) {
	scoreDetail := make(map[Item]int)
	isOwnDraw := false
	if items.ContainsI(Lose) {
		scoreDetail[Lose] = getLoseScore(n)
		if items.ContainsI(Waiting) {
			scoreDetail[Waiting] = getWaitingScore()
		}
		return sum(scoreDetail), scoreDetail
	}

	if items.ContainsI(OwnDraw) {
		items.DelI(OwnDraw)
		isOwnDraw = true
	}

	if items.ContainsI(Win) {
		scoreDetail[Win] = getWinScore(n)
		items.DelI(Win)
	}

	for _, item := range *items {
		generator, ok := ItemGeneratorMap[item]
		if !ok {
			panic(fmt.Sprintf("no such generator: %s", item))
		}
		oldScore, ok := scoreDetail[item]
		if !ok {
			oldScore = 0
		}
		scoreDetail[item] = oldScore + generator()
	}

	if isOwnDraw {
		scoreDetail[OwnDraw] = int(float64(sum(scoreDetail))*OwnDrawTimes + 0.5)
	}

	return sum(scoreDetail), scoreDetail
}

func sum(scoreDetail map[Item]int) (result int) {
	for _, v := range scoreDetail {
		result += v
	}
	return result
}

func getWinScore(n int) int {
	return int(math.Ceil(100/(1.0+0.3*float64(n)) - 2))
}

func getExposedKongScore() int {
	return KongScore
}

func getConcealedScore() int {
	return 2 * KongScore
}

func getWaitingScore() int {
	return utils.RandInt(WaitingMinScore, WaitingMaxScore)
}

func getLoseScore(n int) int {
	return -int(math.Ceil(70/(1.0-0.2*float64(n-100)) + 9))
}

func getMissOneScore() int {
	return utils.RandInt(MissKindMinScore, MissKindMaxScore)
}

func getMissTwoScore() int {
	return getMissOneScore() * 2
}

func getCounterScore() int {
	return utils.RandInt(CounterMinScore, CounterMaxScore)
}

func getSevenPairsScore() int {
	return SevenPairsScore
}

func getThirteenOneScore() int {
	return ThirteenOneScore
}

func getSingleScore() int {
	return utils.RandInt(SingleMinScore, SingleMaxScore)
}

func getContinuousLineScore() int {
	return ContinuousLineScore
}
