package player

import (
	"math/rand"
	"time"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	r "github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/rules"
	s "github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

type Player struct {
	CardsPool   *c.Cards
	PlayerCards *c.Cards
	RobotCards  *c.Cards
	FixedCards  *c.Cards
	KongCount   int
	result      r.Result
	counter     int
	ScoreItems  *s.Items
	lang        utils.Lang
	condition   int
	isWaiting   bool
	isClear     bool
	limit       int
	logger      logging.Logger
	ui          ui.UI
	t           language.Translator
	sleepTime   time.Duration
	sleep       func(time.Duration)
	rng         *rand.Rand
	seed        int64
	session     *replay.Session
	decisionID  int
	err         error
	input       replay.DecisionSource
}

// Options injects setup, decisions and timing without coupling them to the UI.
type Options struct {
	Input   replay.DecisionSource
	Seed    *int64
	Initial *replay.State
	Session *replay.Session
	Sleep   func(time.Duration)
}

func NewPlayer(logger logging.Logger, surface ui.UI, t language.Translator, lang utils.Lang) *Player {
	return NewPlayerWithOptions(logger, surface, t, lang, Options{})
}

func NewPlayerWithOptions(logger logging.Logger, surface ui.UI, t language.Translator, lang utils.Lang, options Options) *Player {
	seed := int64(0)
	if options.Seed != nil {
		seed = *options.Seed
	} else {
		seed = utils.NextRoundSeed()
	}
	sleeper := options.Sleep
	if sleeper == nil {
		sleeper = time.Sleep
	}
	var st time.Duration
	if !utils.Config.QuickMode {
		st = utils.SleepTime
	}
	p := &Player{
		CardsPool: &c.Cards{}, PlayerCards: &c.Cards{}, RobotCards: &c.Cards{}, FixedCards: &c.Cards{},
		ScoreItems: &s.Items{}, lang: lang, isClear: true, condition: rules.InitialPlayerTurnParity, logger: logger, ui: surface, t: t,
		sleepTime: st, sleep: sleeper, rng: rand.New(rand.NewSource(seed)), seed: seed, session: options.Session, input: options.Input,
	}
	if options.Initial != nil {
		p.restore(*options.Initial)
	} else {
		p.init()
	}
	return p
}
