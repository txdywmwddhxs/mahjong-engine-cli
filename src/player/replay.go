package player

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/rules"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

type inputAbort struct{ err error }

func (p *Player) randInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return min + p.rng.Intn(max-min+1)
}
func (p *Player) Seed() int64 { return p.seed }
func (p *Player) Err() error  { return p.err }
func (p *Player) Snapshot() replay.State {
	return (replay.State{Wall: *p.CardsPool, Hand: *p.PlayerCards, Robot: *p.RobotCards, Fixed: *p.FixedCards, Kongs: p.KongCount, Result: p.result, Turn: p.counter, Condition: p.condition, Waiting: p.isWaiting, Clear: p.isClear, Limit: p.limit, ScoreItems: *p.ScoreItems}).Clone()
}
func (p *Player) restore(state replay.State) {
	state = state.Clone()
	p.CardsPool = &state.Wall
	p.PlayerCards = &state.Hand
	p.RobotCards = &state.Robot
	p.FixedCards = &state.Fixed
	p.KongCount = state.Kongs
	p.result = state.Result
	p.counter = state.Turn
	p.condition = state.Condition
	p.isWaiting = state.Waiting
	p.isClear = state.Clear
	p.limit = state.Limit
	p.ScoreItems = &state.ScoreItems
}
func (p *Player) StartRecording() *replay.Session {
	p.session = replay.Record(p.Snapshot(), p.seed, utils.Version)
	return p.session
}
func (p *Player) emit(kind string, tile card.Card) {
	p.observe(replay.Frame{Kind: kind, Tile: tile, State: p.Snapshot()})
}
func (p *Player) observe(frame replay.Frame) {
	if err := p.session.Observe(frame); err != nil {
		panic(inputAbort{err})
	}
}
func (p *Player) decide(kind string, tile card.Card, prompt string) (replay.Action, int) {
	p.decisionID++
	if p.decisionID > rules.MaxDecisionRequests {
		panic(inputAbort{fmt.Errorf("decision limit exceeded")})
	}
	req := replay.Request{ID: p.decisionID, Turn: p.counter, Kind: kind, Tile: tile}
	action, err := p.session.Decision(req, p.Snapshot(), func() (replay.Action, error) {
		if p.input != nil {
			return p.input.Decide(req, p.Snapshot())
		}
		if p.ui == nil {
			return replay.Action{}, io.EOF
		}
		input, err := p.ui.PromptInfo(prompt)
		if err != nil {
			return replay.Action{}, fmt.Errorf("request %d turn %d (%s): %w", req.ID, req.Turn, req.Kind, err)
		}
		return replay.ParseAction(kind, input), nil
	})
	if err != nil {
		panic(inputAbort{err})
	}
	if err := action.Validate(kind); err != nil {
		panic(inputAbort{err})
	}
	entry, _ := json.Marshal(struct {
		Request replay.Request `json:"request"`
		Action  replay.Action  `json:"action"`
	}{req, action})
	p.logger.Debug("Decision: "+string(entry), p.counter)
	return action, req.ID
}
func (p *Player) decided(id int, accepted bool) {
	p.logger.Debug(fmt.Sprintf("Decision #%d accepted=%t", id, accepted), p.counter)
	p.observe(replay.Frame{Kind: "decision_result", DecisionID: id, Accepted: &accepted, State: p.Snapshot()})
}

// VerifyReplay restores the saved setup. It never reads stdin, changes config,
// records stats or sleeps, and stops at the first semantic divergence.
func VerifyReplay(file *replay.File) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("replay engine failed: %v", recovered)
		}
	}()
	session, err := replay.Verify(file)
	if err != nil {
		return err
	}
	p := NewPlayerWithOptions(logging.NopLogger(), ui.NewConsole(nil, io.Discard), language.New(utils.English), utils.English, Options{Seed: &file.Seed, Initial: &file.Initial, Session: session, Sleep: func(time.Duration) {}})
	p.Main()
	if p.Err() != nil {
		return p.Err()
	}
	return session.Finish()
}
