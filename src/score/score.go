package score

import (
	"fmt"
	"math"
)

var ItemScoreFunctions = map[Item]func() int{
	ExposedKong:     getExposedKongScore,
	ConcealedKong:   getConcealedKongScore,
	Waiting:         getWaitingScore,
	MissOneKind:     getMissingOneSuitScore,
	MissTwoKind:     getMissingTwoSuitsScore,
	Counter:         getCounterScore,
	SevenPairs:      getSevenPairsScore,
	ThirteenOrphans: getThirteenOrphansScore,
	ContinuousLine:  getContinuousLineScore,
	Single:          getSingleScore,
}

func GetScore(items *Items, turnCount int) (int, map[Item]int) {
	// Work on an owned copy: settlement must not consume recorded score items.
	remainingItems := append(Items(nil), (*items)...)
	items = &remainingItems
	scoreDetail := make(map[Item]int)
	hasOwnDraw := false
	if items.Contains(Lose) {
		scoreDetail[Lose] = getLoseScore(turnCount)
		if items.Contains(Waiting) {
			scoreDetail[Waiting] = getWaitingScore()
		}
		return sum(scoreDetail), scoreDetail
	}

	if items.Contains(OwnDraw) {
		items.Remove(OwnDraw)
		hasOwnDraw = true
	}

	if items.Contains(Win) {
		scoreDetail[Win] = getWinScore(turnCount)
		items.Remove(Win)
	}

	for _, item := range *items {
		scoreForItem, ok := ItemScoreFunctions[item]
		if !ok {
			panic(fmt.Sprintf("no such generator: %s", item))
		}
		accumulatedScore, ok := scoreDetail[item]
		if !ok {
			accumulatedScore = 0
		}
		scoreDetail[item] = accumulatedScore + scoreForItem()
	}

	if hasOwnDraw {
		scoreDetail[OwnDraw] = int(float64(sum(scoreDetail))*OwnDrawTimes + positiveRoundingOffset)
	}

	return sum(scoreDetail), scoreDetail
}

func sum(scoreDetail map[Item]int) (result int) {
	for _, v := range scoreDetail {
		result += v
	}
	return result
}

func getWinScore(turnCount int) int {
	return int(math.Ceil(winBaseScore/(1.0+winTurnDecay*float64(turnCount)) - winScoreOffset))
}

func getExposedKongScore() int {
	return KongScore
}

func getConcealedKongScore() int {
	return concealedKongMultiplier * KongScore
}

func getWaitingScore() int {
	return WaitingMinScore
}

func getLoseScore(turnCount int) int {
	return -int(math.Ceil(lossBasePenalty/(1.0-lossTurnDecay*float64(turnCount-lossReferenceTurn)) + lossPenaltyOffset))
}

func getMissingOneSuitScore() int {
	return MissKindMinScore
}

func getMissingTwoSuitsScore() int {
	return getMissingOneSuitScore() * maxMissingSuitBonus
}

func getCounterScore() int {
	return CounterMinScore
}

func getSevenPairsScore() int {
	return SevenPairsScore
}

func getThirteenOrphansScore() int {
	return ThirteenOrphansScore
}

func getSingleScore() int {
	return SingleMinScore
}

func getContinuousLineScore() int {
	return ContinuousLineScore
}
