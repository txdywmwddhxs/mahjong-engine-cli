package language

// Key is a stable message identifier. Business logic should depend on Keys,
// not on specific language strings.
type Key string

const (
	// Gameplay / prompts
	KeyLimit                  Key = "limit"
	KeyRateOfRobotsCard       Key = "rate_of_robots_card"
	KeyGetNewCard             Key = "get_new_card"
	KeyRandomNumberSeed       Key = "random_number_seed"
	KeyInitCardPoolDone       Key = "init_card_pool_done"
	KeyInitPlayerCardDone     Key = "init_player_card_done"
	KeyInitRobotCardDone      Key = "init_robot_card_done"
	KeyAskIfPung              Key = "ask_if_pung"
	KeyAskIfKong              Key = "ask_if_kong"
	KeyAskPlayerPungInputInfo Key = "ask_player_pung_input_info"
	KeyAskPlayerKongInputInfo Key = "ask_player_kong_input_info"
	KeyFixedCard              Key = "fixed_card"
	KeyEnd                    Key = "end"
	KeyUnrecognizedInput      Key = "unrecognized_input"
	KeyPlayerPung             Key = "player_pung"
	KeyPlayerKong             Key = "player_kong"
	KeyPlayerAddKong          Key = "player_add_kong"
	KeyUWin                   Key = "u_win"
	KeyULose                  Key = "u_lose"
	KeyDraw                   Key = "draw"
	KeyRobotsCard             Key = "robots_card"
	KeyCardAnalysisResult     Key = "card_analysis_result"
	KeyMatchSevenPairs        Key = "match_seven_pairs"
	KeyMatchThirteenOne       Key = "match_thirteen_one"
	KeyMatchWaiting           Key = "match_waiting"
	KeyNotMatchWaiting        Key = "not_match_waiting"
	KeyLastCard               Key = "last_card"
	KeyUserWin                Key = "user_win"
	KeyUserNotWin             Key = "user_not_win"
	KeyPlayACard              Key = "play_a_card"
	KeyPlayACardInputInfo     Key = "play_a_card_input_info"
	KeyAutoPlay               Key = "auto_play"
	KeyEnterScript            Key = "enter_script"
	KeyCancelScript           Key = "cancel_script"
	KeyPlayerWaiting          Key = "player_waiting"
	KeyCannotWaiting          Key = "cannot_waiting"
	KeyCannotWaitingPlayCard  Key = "cannot_waiting_play_card"
	KeyCannotPlayCard         Key = "cannot_play_card"
	KeyWithoutThisCard        Key = "without_this_card"
	KeyRobotPlayACard         Key = "robot_play_a_card"

	// Meta / logs
	KeyCurrentInfo Key = "current_info"
	KeyGameBegin   Key = "game_begin"
	KeyGameEnd     Key = "game_end"
	KeyGameVersion Key = "game_version"
	KeyGameSeed    Key = "game_seed"
	KeyThisScore   Key = "this_score"
	KeyTotalScore  Key = "total_score"
	KeyContinue    Key = "continue"
	KeyMeetError   Key = "meet_error"
	KeyCardKind    Key = "card_kind"
	KeyScoreItem   Key = "score_item"

	// Score item names
	KeyScoreWin            Key = "score_win"
	KeyScoreLose           Key = "score_lose"
	KeyScoreSingle         Key = "score_single"
	KeyScoreCounter        Key = "score_counter" // expects %d argument (CountOfWinAward)
	KeyScoreOwnDraw        Key = "score_own_draw"
	KeyScoreWaiting        Key = "score_waiting"
	KeyScoreExposedKong    Key = "score_exposed_kong"
	KeyScoreConcealedKong  Key = "score_concealed_kong"
	KeyScoreThirteenOne    Key = "score_thirteen_one"
	KeyScoreSevenPairs     Key = "score_seven_pairs"
	KeyScoreMissOneKind    Key = "score_miss_one_kind"
	KeyScoreMissTwoKind    Key = "score_miss_two_kind"
	KeyScoreContinuousLine Key = "score_continuous_line"
)

// AllKeys enumerates all translation keys used by the program. Keep this list
// updated when adding new Keys; tests will enforce catalog coverage.
var AllKeys = []Key{
	// Gameplay / prompts
	KeyLimit,
	KeyRateOfRobotsCard,
	KeyGetNewCard,
	KeyRandomNumberSeed,
	KeyInitCardPoolDone,
	KeyInitPlayerCardDone,
	KeyInitRobotCardDone,
	KeyAskIfPung,
	KeyAskIfKong,
	KeyAskPlayerPungInputInfo,
	KeyAskPlayerKongInputInfo,
	KeyFixedCard,
	KeyEnd,
	KeyUnrecognizedInput,
	KeyPlayerPung,
	KeyPlayerKong,
	KeyPlayerAddKong,
	KeyUWin,
	KeyULose,
	KeyDraw,
	KeyRobotsCard,
	KeyCardAnalysisResult,
	KeyMatchSevenPairs,
	KeyMatchThirteenOne,
	KeyMatchWaiting,
	KeyNotMatchWaiting,
	KeyLastCard,
	KeyUserWin,
	KeyUserNotWin,
	KeyPlayACard,
	KeyPlayACardInputInfo,
	KeyAutoPlay,
	KeyEnterScript,
	KeyCancelScript,
	KeyPlayerWaiting,
	KeyCannotWaiting,
	KeyCannotWaitingPlayCard,
	KeyCannotPlayCard,
	KeyWithoutThisCard,
	KeyRobotPlayACard,

	// Meta / logs
	KeyCurrentInfo,
	KeyGameBegin,
	KeyGameEnd,
	KeyGameVersion,
	KeyGameSeed,
	KeyThisScore,
	KeyTotalScore,
	KeyContinue,
	KeyMeetError,
	KeyCardKind,
	KeyScoreItem,

	// Score item names
	KeyScoreWin,
	KeyScoreLose,
	KeyScoreSingle,
	KeyScoreCounter,
	KeyScoreOwnDraw,
	KeyScoreWaiting,
	KeyScoreExposedKong,
	KeyScoreConcealedKong,
	KeyScoreThirteenOne,
	KeyScoreSevenPairs,
	KeyScoreMissOneKind,
	KeyScoreMissTwoKind,
	KeyScoreContinuousLine,
}
