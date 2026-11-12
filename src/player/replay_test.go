package player

import (
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func autoAction(req replay.Request, state replay.State) (replay.Action, error) {
	if req.Kind != "discard" {
		return replay.Action{Kind: "accept"}, nil
	}
	if req.Tile != "" {
		return replay.Action{Kind: "discard_drawn"}, nil
	}
	return replay.Action{Kind: "discard", Tile: state.Hand[0]}, nil
}
func recordGame(t *testing.T, seed int64, initial *replay.State, policy replay.SourceFunc) *replay.File {
	t.Helper()
	if policy == nil {
		policy = autoAction
	}
	p := NewPlayerWithOptions(logging.NopLogger(), nil, language.New(utils.Chinese), utils.Chinese, Options{Seed: &seed, Initial: initial, Input: policy, Sleep: func(time.Duration) {}})
	session := p.StartRecording()
	p.Main()
	if p.Err() != nil {
		t.Fatal(p.Err())
	}
	if err := session.Finish(); err != nil {
		t.Fatal(err)
	}
	if err := session.File.Validate(); err != nil {
		t.Fatal(err)
	}
	return session.File
}
func cloneFile(f *replay.File) *replay.File {
	b, _ := json.Marshal(f)
	var out replay.File
	_ = json.Unmarshal(b, &out)
	return &out
}

func TestReplayRoundTripAndIndependentGames(t *testing.T) {
	for _, seed := range []int64{1, 2, 42, 1766309534090039000} {
		f := recordGame(t, seed, nil, nil)
		for i := 0; i < 3; i++ {
			if err := VerifyReplay(f); err != nil {
				t.Fatalf("seed %d: %v", seed, err)
			}
		}
		again := recordGame(t, seed, nil, nil)
		if !reflect.DeepEqual(f.Initial, again.Initial) || !reflect.DeepEqual(f.Frames, again.Frames) {
			t.Fatalf("seed %d produced different games", seed)
		}
		// Replaying setup does not regenerate it, even if the diagnostic seed differs.
		again.Seed = 999
		if err := VerifyReplay(again); err != nil {
			t.Fatal(err)
		}
	}
}

// Build valid 136-tile initial states with selected future draws/discards.
func scenario(t *testing.T, hand, prefix, suffix card.Cards) replay.State {
	t.Helper()
	counts := map[card.Card]int{}
	for _, tile := range card.CardsList {
		counts[tile] = 4
	}
	robot := card.Cards{"D", "D", "F", "F", "F", "Z", "Z", "Z", "W", "W", "W", "N", "N", "N"}
	for _, group := range []card.Cards{hand, prefix, suffix, robot} {
		for _, tile := range group {
			counts[tile]--
			if counts[tile] < 0 {
				t.Fatalf("too many %s", tile)
			}
		}
	}
	wall := append(card.Cards{}, prefix...)
	for _, tile := range card.CardsList {
		for n := 0; n < counts[tile]; n++ {
			wall = append(wall, tile)
		}
	}
	wall = append(wall, suffix...)
	return replay.State{Wall: wall, Hand: hand, Robot: robot, Fixed: card.Cards{}, ScoreItems: score.Items{}, Clear: true, Limit: 12}
}

func TestReplayMeldsWaitingAndWin(t *testing.T) {
	base := card.Cards{"B6", "B6", "B2", "B3", "B4", "T1", "T3", "T5", "T7", "W2", "W4", "W6", "W8"}
	cases := []struct {
		name                 string
		hand, prefix, suffix card.Cards
		event                string
	}{
		{"pung_and_upgrade", base, card.Cards{"X", "B6"}, card.Cards{"X", "B6"}, "upgrade_kong"},
		{"exposed_kong", card.Cards{"B6", "B6", "B6", "B3", "B4", "T1", "T3", "T5", "T7", "W2", "W4", "W6", "W8"}, card.Cards{"X"}, card.Cards{"B6"}, "kong"},
		{"concealed_kong", card.Cards{"B3", "B3", "B3", "B6", "B4", "T1", "T3", "T5", "T7", "W2", "W4", "W6", "W8"}, card.Cards{"B3", "X"}, nil, "kong"},
		{"waiting_and_win", card.Cards{"B1", "B1", "B2", "B2", "B3", "B3", "T1", "T1", "T2", "T2", "T3", "T3", "W9"}, card.Cards{"X"}, card.Cards{"W9"}, "settlement"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initial := scenario(t, tc.hand, tc.prefix, tc.suffix)
			policy := replay.SourceFunc(autoAction)
			if tc.name == "waiting_and_win" {
				policy = func(req replay.Request, s replay.State) (replay.Action, error) {
					if req.Kind == "discard" {
						return replay.Action{Kind: "declare_waiting", Tile: req.Tile}, nil
					}
					return autoAction(req, s)
				}
			}
			f := recordGame(t, 1, &initial, policy)
			found := false
			for _, frame := range f.Frames {
				if frame.Kind == tc.event {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing %s event", tc.event)
			}
			if tc.name == "waiting_and_win" {
				end := f.Frames[len(f.Frames)-1]
				if !end.State.Waiting || end.Settlement.Result != result.Win || end.Settlement.Score != 68 || end.Settlement.Counter != 3 {
					t.Fatalf("bad waiting settlement: %+v", end.Settlement)
				}
			}
			if err := VerifyReplay(f); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestReplayRejectedInputAndTampering(t *testing.T) {
	rejected := false
	f := recordGame(t, 1, nil, func(req replay.Request, s replay.State) (replay.Action, error) {
		if !rejected {
			rejected = true
			return replay.Action{Kind: "invalid", Raw: "???"}, nil
		}
		return autoAction(req, s)
	})
	hasRejected := false
	for _, frame := range f.Frames {
		if frame.Accepted != nil && !*frame.Accepted {
			hasRejected = true
		}
	}
	if !hasRejected {
		t.Fatal("rejected input was not recorded")
	}
	if err := VerifyReplay(f); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, contains string
		mutate         func(*replay.File)
	}{
		{"state", "state.hand", func(f *replay.File) { f.Frames[0].State.Hand[0] = "W9"; f.Frames[0].Hash = f.Frames[0].State.Hash() }},
		{"request", "request.kind", func(f *replay.File) {
			for i := range f.Frames {
				x := &f.Frames[i]
				if x.Kind == "decision" {
					x.Request.Kind = "confirm_pung"
					x.Request.Tile = "B1"
					x.Action = &replay.Action{Kind: "decline"}
					break
				}
			}
		}},
		{"settlement", "settlement.score", func(f *replay.File) { f.Frames[len(f.Frames)-1].Settlement.Score++ }},
		{"truncated", "settlement", func(f *replay.File) { f.Frames = f.Frames[:len(f.Frames)-1] }},
		{"extra_frames", "settlement", func(f *replay.File) { f.Frames = append(f.Frames, f.Frames[len(f.Frames)-1]) }},
		{"format", "format", func(f *replay.File) { f.Format++ }},
		{"rules", "rules", func(f *replay.File) { f.Rules = "future-v2" }},
		{"incomplete", "incomplete", func(f *replay.File) { f.Complete = false }},
		{"hash", "hash", func(f *replay.File) { f.Frames[0].Hash = "bad" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			changed := cloneFile(f)
			tc.mutate(changed)
			err := VerifyReplay(changed)
			if err == nil || !strings.Contains(err.Error(), tc.contains) {
				t.Fatalf("wanted %q, got %v", tc.contains, err)
			}
		})
	}
}

func TestReplayFileAndViewer(t *testing.T) {
	f := recordGame(t, 2, nil, nil)
	path := filepath.Join(t.TempDir(), "game.json")
	if err := f.Save(path); err != nil {
		t.Fatal(err)
	}
	if err := f.Save(path); err == nil {
		t.Fatal("overwrote recording")
	}
	loaded, err := replay.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = VerifyReplay(loaded); err != nil {
		t.Fatal(err)
	}
	// Historical rules can be viewed, but not verified by this engine.
	loaded.Rules = "old-rules"
	var out strings.Builder
	if err = replay.View(loaded, strings.NewReader("n\np\ng 2\nt 2\nq\n"), &out, true); err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"Frame 0/", "Frame 1/", "Frame 2/", "turn 2"} {
		if !strings.Contains(out.String(), part) {
			t.Fatalf("viewer missing %s", part)
		}
	}
	if err = VerifyReplay(loaded); err == nil {
		t.Fatal("verified incompatible rules")
	}
}

func TestEOFUnwindsEveryDecision(t *testing.T) {
	cases := []struct {
		name    string
		initial *replay.State
		policy  replay.SourceFunc
	}{
		{"discard", nil, nil},
		{"opening_kong", nil, nil},
		{"after_pung", nil, nil},
	}
	kong := scenario(t, card.Cards{"B3", "B3", "B3", "B6", "B4", "T1", "T3", "T5", "T7", "W2", "W4", "W6", "W8"}, card.Cards{"B3"}, nil)
	cases[1].initial = &kong
	pung := scenario(t, card.Cards{"B6", "B6", "B2", "B3", "B4", "T1", "T3", "T5", "T7", "W2", "W4", "W6", "W8"}, card.Cards{"X"}, card.Cards{"B6"})
	cases[2].initial = &pung
	cases[2].policy = func(req replay.Request, s replay.State) (replay.Action, error) {
		if req.Kind == "discard" && req.Tile == "" {
			return replay.Action{}, io.EOF
		}
		return autoAction(req, s)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seed := int64(1)
			options := Options{Seed: &seed, Initial: tc.initial, Sleep: func(time.Duration) {}}
			if tc.policy != nil {
				options.Input = tc.policy
			}
			p := NewPlayerWithOptions(logging.NopLogger(), nil, language.New(utils.English), utils.English, options)
			res, _, _ := p.Main()
			if res != result.Cancel || !errors.Is(p.Err(), io.EOF) {
				t.Fatalf("res=%v err=%v", res, p.Err())
			}
		})
	}
}
