// Package replay defines the language-independent, versioned replay format.
// It has no dependency on the game loop or terminal UI.
package replay

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
)

const FormatVersion = 1
const RulesVersion = "challenge-v1"
const MaxFrames = 10000
const MaxFileBytes = 32 << 20

// State includes all rule-affecting mutable game state; presentation and timing
// are deliberately excluded. Slice order is significant, especially the wall.
type State struct {
	Wall       card.Cards    `json:"wall"`
	Hand       card.Cards    `json:"hand"`
	Robot      card.Cards    `json:"robot"`
	Fixed      card.Cards    `json:"fixed"`
	Kongs      int           `json:"kongs"`
	Result     result.Result `json:"result"`
	Turn       int           `json:"turn"`
	Condition  int           `json:"condition"`
	Waiting    bool          `json:"waiting"`
	Clear      bool          `json:"clear"`
	Limit      int           `json:"limit"`
	ScoreItems score.Items   `json:"score_items"`
}

func (s State) Clone() State {
	s.Wall = append(card.Cards{}, s.Wall...)
	s.Hand = append(card.Cards{}, s.Hand...)
	s.Robot = append(card.Cards{}, s.Robot...)
	s.Fixed = append(card.Cards{}, s.Fixed...)
	s.ScoreItems = append(score.Items{}, s.ScoreItems...)
	return s
}
func (s State) Hash() string {
	b, _ := json.Marshal(s)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// Request IDs count attempts, including rejected input; turns count game turns.
type Request struct {
	ID   int       `json:"id"`
	Turn int       `json:"turn"`
	Kind string    `json:"kind"`
	Tile card.Card `json:"tile,omitempty"`
}

// DecisionSource can drive a game without any terminal prompts.
type DecisionSource interface {
	Decide(Request, State) (Action, error)
}
type SourceFunc func(Request, State) (Action, error)

func (f SourceFunc) Decide(r Request, s State) (Action, error) { return f(r, s) }

type Action struct {
	Kind string    `json:"kind"`
	Tile card.Card `json:"tile,omitempty"`
	Raw  string    `json:"raw,omitempty"` // only invalid input, for parser regression
}

func ParseAction(kind, input string) Action {
	input = strings.ToUpper(strings.TrimSpace(input))
	if kind != "discard" {
		switch input {
		case "Y", "YES":
			return Action{Kind: "accept"}
		case "N", "NO":
			return Action{Kind: "decline"}
		}
		return Action{Kind: "invalid", Raw: input}
	}
	switch input {
	case "":
		return Action{Kind: "discard_drawn"}
	case "END", "ENDL":
		return Action{Kind: "cancel"}
	case "WHOSYOURDADDY":
		return Action{Kind: "debug_win"}
	case "WARPTEN":
		return Action{Kind: "toggle_speed"}
	}
	fields := strings.Fields(input)
	if len(fields) == 2 && (fields[1] == "TING" || fields[1] == "TT") && card.CardsList.Contains(card.Card(fields[0])) {
		return Action{Kind: "declare_waiting", Tile: card.Card(fields[0])}
	}
	if card.CardsList.Contains(card.Card(input)) {
		return Action{Kind: "discard", Tile: card.Card(input)}
	}
	return Action{Kind: "invalid", Raw: input}
}

func (a Action) Input() string {
	switch a.Kind {
	case "accept":
		return "Y"
	case "decline":
		return "N"
	case "discard_drawn":
		return ""
	case "cancel":
		return "END"
	case "debug_win":
		return "WHOSYOURDADDY"
	case "toggle_speed":
		return "WARPTEN"
	case "declare_waiting":
		return string(a.Tile) + " TT"
	case "discard":
		return string(a.Tile)
	default:
		return a.Raw
	}
}

func (a Action) Validate(kind string) error {
	switch kind {
	case "discard", "confirm_pung", "confirm_kong":
	default:
		return fmt.Errorf("unknown request kind %q", kind)
	}
	if a.Kind == "invalid" {
		if a != ParseAction(kind, a.Raw) {
			return fmt.Errorf("invalid raw input action: %+v", a)
		}
		return nil
	}
	valid := a.Raw == ""
	if kind != "discard" {
		valid = valid && (a.Kind == "accept" || a.Kind == "decline") && a.Tile == ""
	} else {
		switch a.Kind {
		case "discard", "declare_waiting":
			valid = valid && card.CardsList.Contains(a.Tile)
		case "discard_drawn", "cancel", "debug_win", "toggle_speed":
			valid = valid && a.Tile == ""
		default:
			valid = false
		}
	}
	if !valid {
		return fmt.Errorf("invalid action for %s: %+v", kind, a)
	}
	return nil
}

type Settlement struct {
	Result  result.Result      `json:"result"`
	Score   int                `json:"score"`
	Counter int                `json:"counter"`
	Detail  map[score.Item]int `json:"detail"`
}

type Frame struct {
	Kind       string      `json:"kind"`
	Request    *Request    `json:"request,omitempty"`
	Action     *Action     `json:"action,omitempty"`
	Accepted   *bool       `json:"accepted,omitempty"`
	DecisionID int         `json:"decision_id,omitempty"`
	Tile       card.Card   `json:"tile,omitempty"`
	State      State       `json:"state"`
	Hash       string      `json:"hash"`
	Settlement *Settlement `json:"settlement,omitempty"`
}

type File struct {
	Format      int     `json:"format_version"`
	Rules       string  `json:"rules_version"`
	CodeVersion string  `json:"code_version"`
	GameID      string  `json:"game_id"`
	Seed        int64   `json:"seed"`
	Initial     State   `json:"initial"`
	InitialHash string  `json:"initial_hash"`
	Frames      []Frame `json:"frames"`
	Complete    bool    `json:"complete"`
}

// Difference identifies the first field that differs, not just its hash.
func Difference(expected, actual any) string {
	e, _ := json.Marshal(expected)
	a, _ := json.Marshal(actual)
	var ev, av any
	_ = json.Unmarshal(e, &ev)
	_ = json.Unmarshal(a, &av)
	return diff("", ev, av)
}
func diff(path string, e, a any) string {
	if reflect.DeepEqual(e, a) {
		return ""
	}
	if em, ok := e.(map[string]any); ok {
		if am, ok := a.(map[string]any); ok {
			keys := make(map[string]bool, len(em)+len(am))
			for k := range em {
				keys[k] = true
			}
			for k := range am {
				keys[k] = true
			}
			names := make([]string, 0, len(keys))
			for k := range keys {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, k := range names {
				if d := diff(strings.TrimPrefix(path+"."+k, "."), em[k], am[k]); d != "" {
					return d
				}
			}
		}
	}
	if es, ok := e.([]any); ok {
		if as, ok := a.([]any); ok {
			if len(es) != len(as) {
				return fmt.Sprintf("%s.length: expected %d, actual %d", path, len(es), len(as))
			}
			for i := range es {
				if d := diff(fmt.Sprintf("%s[%d]", path, i), es[i], as[i]); d != "" {
					return d
				}
			}
		}
	}
	return fmt.Sprintf("%s: expected %s, actual %s", path, compact(e), compact(a))
}
func compact(v any) string { b, _ := json.Marshal(v); return string(b) }
