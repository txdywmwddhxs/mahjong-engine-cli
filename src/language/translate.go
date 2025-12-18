package language

import (
	"fmt"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

var IsChinese = utils.Config.Lang == utils.Chinese

func Limit() string {
	if IsChinese {
		return "限制张数"
	}
	return "Count of limit"
}

func RateOfRobotsCard() string {
	if IsChinese {
		return "机器人手牌、坏牌和正常牌的概率为"
	}
	return "The probability that the robot's card is good, bad and normal is"
}

func GetNewCard() string {
	if IsChinese {
		return "摸到新牌"
	}
	return "Get new card"
}

func RandomNumberSeed() string {
	if IsChinese {
		return "随机数种子"
	}
	return "Random number seed"
}

func InitCardPoolDone() string {
	if IsChinese {
		return "初始化牌池完成"
	}
	return "Initial card pool is done"
}

func InitPlayerCardDone() string {
	if IsChinese {
		return "初始化用户牌型完成"
	}
	return "Initial player's card is done"
}

func InitRobotCardDone() string {
	if IsChinese {
		return "初始化机器人牌型完成"
	}
	return "Initial robot's card is done"
}

func AskIfPung() string {
	if IsChinese {
		return "是否碰牌"
	}
	return "Whether to pung"
}

func AskIfKong() string {
	if IsChinese {
		return "是否杠牌"
	}
	return "Whether to kong"
}

func AskPlayerPungInputInfo() string {
	if IsChinese {
		return "是否碰牌输入信息"
	}
	return "Whether to pung input message"
}

func AskPlayerKongInputInfo() string {
	if IsChinese {
		return "是否杠牌输入信息"
	}
	return "Whether to kong input message"
}

func FixedCard() string {
	if IsChinese {
		return "固定牌张"
	}
	return "Fixed card"
}

func End() string {
	if IsChinese {
		return "收到终止信号，牌局终止"
	}
	return "The card game is terminated when the termination signal is received"
}

func UnrecognizedInput() string {
	if IsChinese {
		return "不能识别的输入"
	}
	return "Unrecognized input"
}

func PlayerPung() string {
	if IsChinese {
		return "玩家碰牌"
	}
	return "Player pung"
}

func PlayerKong() string {
	if IsChinese {
		return "玩家杠牌"
	}
	return "Player kong"
}

func PlayerAddKong() string {
	if IsChinese {
		return "玩家加杠"
	}
	return "Player add kong"
}

func UWin() string {
	if IsChinese {
		return "你赢了"
	}
	return "You win"
}

func ULose() string {
	if IsChinese {
		return "你输了"
	}
	return "You lose"
}

func Draw() string {
	if IsChinese {
		return "流局"
	}
	return "Draw"
}

func RobotsCard() string {
	if IsChinese {
		return "机器人牌型为"
	}
	return "The robot's card is"
}

func CardAnalysisResult() string {
	if IsChinese {
		return "牌型分析结果"
	}
	return "Card analysis result"
}

func MatchSevenPairs() string {
	if IsChinese {
		return "符合七小对的牌型"
	}
	return "Match the card of seven pairs"
}

func MatchThirteenOne() string {
	if IsChinese {
		return "符合十三幺的牌型"
	}
	return "Match the card od thirteen-one"
}

func MatchWaiting() string {
	if IsChinese {
		return "符合听牌牌型"
	}
	return "Match the card of waiting"
}

func NotMatchWaiting() string {
	if IsChinese {
		return "不符合听牌牌型"
	}
	return "Not Match the card of waiting"
}

func LastCard() string {
	if IsChinese {
		return "牌池剩余牌张"
	}
	return "Last card of card pool"
}

func UserWin() string {
	if IsChinese {
		return "本轮用户牌型胜利"
	}
	return "This round of player card win"
}

func UserNotWin() string {
	if IsChinese {
		return "本轮用户牌型没有胜利"
	}
	return "This round of player card did not win"
}

func PlayACard() string {
	if IsChinese {
		return "打出一张牌"
	}
	return "Play a card"
}

func PlayACardInputInfo() string {
	if IsChinese {
		return "打出一张牌输入信息"
	}
	return "Play a card input message"
}

func AutoPlay() string {
	if IsChinese {
		return "玩家听牌，自动出牌"
	}
	return "Player waiting and automatic play cards"
}

func EnterScript() string {
	if IsChinese {
		return "输入秘籍"
	}
	return "Enter a secret script"
}

func CancelScript() string {
	if IsChinese {
		return "取消秘籍"
	}
	return "Cancel a secret script"
}

func PlayerWaiting() string {
	if IsChinese {
		return "用户听牌，听牌牌张"
	}
	return "Player waiting, waiting cards"
}

func CannotWaiting() string {
	if IsChinese {
		return "不能听牌"
	}
	return "Cannot waiting"
}

func CannotWaitingPlayCard() string {
	if IsChinese {
		return "不能听牌，需要打出一张牌"
	}
	return "Cannot waiting, play a card"
}

func CannotPlayCard() string {
	if IsChinese {
		return "不能打出已经碰或者杠的牌"
	}
	return "Cannot play card that have pung or kong"
}

func RobotPlayACard() string {
	if IsChinese {
		return "机器人打出一张牌"
	}
	return "Robot play a card"
}

func CurrentInfo() string {
	if IsChinese {
		return "当前版本牌局信息"
	}
	return "Game config of current version"
}

func GameBegin() string {
	if IsChinese {
		return "对局开始"
	}
	return "Game Start"
}

func GameEnd() string {
	if IsChinese {
		return "对局结束"
	}
	return "Game End"
}

func GameVersion() string {
	if IsChinese {
		return "牌局版本"
	}
	return "Game version"
}

func ThisScore() string {
	if IsChinese {
		return "本局得分"
	}
	return "Score this game"
}

func TotalScore() string {
	if IsChinese {
		return "当前总积分"
	}
	return "Current total points"
}

func Win() string {
	if IsChinese {
		return "胡牌"
	}
	return "win"
}

func Lose() string {
	if IsChinese {
		return "输牌"
	}
	return "lose"
}

func Single() string {
	if IsChinese {
		return "单吊"
	}
	return "single"
}

func Counter() string {
	if IsChinese {
		return fmt.Sprintf("%d 手内结束", utils.CountOfWinAward)
	}
	return fmt.Sprintf("end in %d", utils.CountOfWinAward)
}

func OwnDraw() string {
	if IsChinese {
		return "自摸"
	}
	return "own_draw"
}

func Waiting() string {
	if IsChinese {
		return "听牌"
	}
	return "waiting"
}

func ExposedKong() string {
	if IsChinese {
		return "明杠"
	}
	return "exposed-kong"
}

func ConcealedKong() string {
	if IsChinese {
		return "暗杠"
	}
	return "concealed-kong"
}

func WithoutThisCard() string {
	if IsChinese {
		return "没有这张牌"
	}
	return "Without this card"
}

func ThirteenOne() string {
	if IsChinese {
		return "十三幺"
	}
	return "thirteen-one"
}

func SevenPairs() string {
	if IsChinese {
		return "七小对"
	}
	return "seven-pairs"
}

func MissOneKind() string {
	if IsChinese {
		return "缺一门"
	}
	return "miss one kind"
}

func MissTwoKind() string {
	if IsChinese {
		return "缺两门"
	}
	return "miss two kind"
}

func CardKind() string {
	if IsChinese {
		return "用户牌型"
	}
	return "Card kind"
}

func ScoreItem() string {
	if IsChinese {
		return "得分项目"
	}
	return "Score items"
}

func ContinuousLine() string {
	if IsChinese {
		return "一条龙"
	}
	return "continuous-line"
}

func Continue() string {
	if IsChinese {
		return "继续？"
	}
	return "Continue? "
}

func MeetError() string {
	if IsChinese {
		return "遇到错误"
	}
	return "Meet error"
}
