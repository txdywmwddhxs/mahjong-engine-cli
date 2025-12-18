package player

import (
	"card"
	log "clog"
	"fmt"
	"testing"
	"utils"
)

func TestPlayer_CheckIsWaiting(t *testing.T) {
	lang := "English"
	logger := log.NewLogger(
		`C:\Users\Administrator\GolandProjects\card_player\log\test.log`,
		utils.Lang(lang))
	p := NewPlayer(logger, utils.Lang(lang))
	//p.PlayerCards = &card.Cards{"B3", "B3", "W2", "W3", "W5", "W6", "W7", "W7", "W8", "W9", "F", "F", "F"}
	p.PlayerCards = &card.Cards{"B1", "B1", "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "B9", "B9"}
	res, c := p.CheckIsWaiting()
	fmt.Println(res, *c)
}
