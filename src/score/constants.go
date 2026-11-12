package score

const (
	WaitingMinScore      = 2
	WaitingMaxScore      = WaitingMinScore
	SingleMinScore       = 3
	SingleMaxScore       = SingleMinScore
	CounterMaxScore      = CounterMinScore
	CounterMinScore      = 2
	MissKindMinScore     = 3
	MissKindMaxScore     = MissKindMinScore
	OwnDrawTimes         = 1.0
	KongScore            = 2
	ThirteenOrphansScore = 20
	SevenPairsScore      = 10
	ContinuousLineScore  = 10
)

// Challenge-v1 score curves. Keep the operation order unchanged so existing
// replay settlements and floating-point rounding remain compatible.
const (
	winBaseScore            = 100
	winTurnDecay            = 0.3
	winScoreOffset          = 2
	lossBasePenalty         = 70
	lossTurnDecay           = 0.2
	lossReferenceTurn       = 100
	lossPenaltyOffset       = 9
	concealedKongMultiplier = 2
	maxMissingSuitBonus     = 2
	positiveRoundingOffset  = 0.5
)
