package player

import (
	"time"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	r "github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	s "github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

type Player struct {
	CardsPool      *c.Cards
	PlayerCards    *c.Cards
	RobotCards     *c.Cards
	FixedCards     *c.Cards
	GroupFourCount int
	result         r.Result
	counter        int
	ScoreItem      *s.Items
	lang           utils.Lang
	condition      int
	isWaiting      bool
	isClear        bool
	limit          int
	logger         logging.Logger
	ui             ui.UI
	t              language.Translator
	sleepTime      time.Duration
}

func NewPlayer(logger logging.Logger, ui ui.UI, t language.Translator, lang utils.Lang) (p *Player) {
	var st time.Duration
	if utils.Config.QuickMode {
		st = 0
	} else {
		st = utils.SleepTime
	}
	p = &Player{
		CardsPool:      &c.Cards{},
		PlayerCards:    &c.Cards{},
		RobotCards:     &c.Cards{},
		FixedCards:     &c.Cards{},
		GroupFourCount: 0,
		result:         r.Default,
		counter:        0,
		ScoreItem:      &s.Items{},
		lang:           lang,
		condition:      0,
		isWaiting:      false,
		isClear:        true,
		limit:          0,
		logger:         logger,
		ui:             ui,
		t:              t,
		sleepTime:      st,
	}
	p.init()
	return p
}
