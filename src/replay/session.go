package replay

import (
	"fmt"
	"time"
)

// Session either records a live game or strictly checks an existing timeline.
// The engine sends semantic events, independent of translated output.
type Session struct {
	File   *File
	verify bool
	cursor int
}

func Record(initial State, seed int64, code string) *Session {
	initial = initial.Clone()
	return &Session{File: &File{Format: FormatVersion, Rules: RulesVersion, CodeVersion: code, GameID: fmt.Sprintf("%d-%d", time.Now().UnixNano(), seed), Seed: seed, Initial: initial, InitialHash: initial.Hash(), Frames: []Frame{}}}
}
func Verify(f *File) (*Session, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	if f.Rules != RulesVersion {
		return nil, fmt.Errorf("incompatible rules version %q (engine %q); viewing is still available", f.Rules, RulesVersion)
	}
	return &Session{File: f, verify: true}, nil
}
func (s *Session) IsVerifying() bool { return s != nil && s.verify }

func (s *Session) Decision(req Request, state State, live func() (Action, error)) (Action, error) {
	if s == nil {
		return live()
	}
	if s.verify {
		if s.cursor >= len(s.File.Frames) {
			return Action{}, fmt.Errorf("replay exhausted at request %d, turn %d (%s)", req.ID, req.Turn, req.Kind)
		}
		expected := s.File.Frames[s.cursor]
		if expected.Kind != "decision" || expected.Request == nil || expected.Action == nil {
			return Action{}, s.failure(req.Turn, "expected a decision, recorded "+expected.Kind)
		}
		if d := Difference(*expected.Request, req); d != "" {
			return Action{}, s.failure(req.Turn, "request."+d)
		}
		if d := Difference(expected.State, state); d != "" {
			return Action{}, s.failure(req.Turn, "state."+d)
		}
		if err := expected.Action.Validate(req.Kind); err != nil {
			return Action{}, err
		}
		s.cursor++
		return *expected.Action, nil
	}
	action, err := live()
	if err != nil {
		return Action{}, err
	}
	f := Frame{Kind: "decision", Request: &req, Action: &action, State: state}
	return action, s.Observe(f)
}
func (s *Session) Observe(f Frame) error {
	if s == nil {
		return nil
	}
	f.State = f.State.Clone()
	f.Hash = f.State.Hash()
	if s.verify {
		if s.cursor >= len(s.File.Frames) {
			return s.failure(f.State.Turn, "unexpected event "+f.Kind+" after end of recording")
		}
		expected := s.File.Frames[s.cursor]
		if d := Difference(expected.State, f.State); d != "" {
			return s.failure(f.State.Turn, "state."+d)
		}
		if d := Difference(expected, f); d != "" {
			return s.failure(f.State.Turn, d)
		}
		s.cursor++
		return nil
	}
	if len(s.File.Frames) >= MaxFrames {
		return fmt.Errorf("recording exceeds %d frames", MaxFrames)
	}
	s.File.Frames = append(s.File.Frames, f)
	return nil
}
func (s *Session) Finish() error {
	if s == nil {
		return nil
	}
	if s.verify {
		if s.cursor != len(s.File.Frames) {
			return s.failure(0, fmt.Sprintf("game ended with %d unconsumed frames", len(s.File.Frames)-s.cursor))
		}
	} else {
		s.File.Complete = true
	}
	return nil
}
func (s *Session) failure(turn int, message string) error {
	return fmt.Errorf("replay diverged at frame %d, turn %d: %s", s.cursor+1, turn, message)
}
