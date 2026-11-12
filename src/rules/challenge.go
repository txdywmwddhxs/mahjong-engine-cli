// Package rules holds shared challenge-mode policy without runtime initialization.
package rules

const (
	EarlyFinishTurnLimit = 25
	RobotMaxPungCount    = 2
	RobotGoodHandPercent = 32
	RobotBadHandPercent  = 10
	PercentScale         = 100
	GoodHandLimitMin     = 22
	GoodHandLimitMax     = 35
	BadHandLimitMin      = 60
	BadHandLimitMax      = 90

	// Keep the historical boundary semantics: a roll is in [0,99], and the
	// challenge generator uses inclusive <= thresholds. This is not a rebalance.
	MaxDifficultyRoll = PercentScale - 1

	// The last 14 wall tiles end a challenge round; this is a gameplay policy,
	// independent of the number of tiles in a winning hand.
	DrawWallThreshold       = 14
	FirstTurn               = 1
	TurnCycleLength         = 2
	InitialPlayerTurnParity = 0
	ClaimedPlayerTurnParity = 1

	// Safety limits are shared by the game loop and replay validation.
	MaxTurns            = 1000
	MaxRecordedTurn     = MaxTurns + 1 // Includes the final settlement counter.
	MaxDecisionRequests = 2000
)
