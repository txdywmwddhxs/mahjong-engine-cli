package replay

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/result"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/rules"
)

func (f *File) Validate() error {
	if f.Format != FormatVersion {
		return fmt.Errorf("unsupported replay format %d", f.Format)
	}
	if f.Rules == "" || f.GameID == "" {
		return fmt.Errorf("missing replay metadata")
	}
	if !f.Complete {
		return fmt.Errorf("incomplete recording")
	}
	if f.InitialHash != f.Initial.Hash() {
		return fmt.Errorf("initial state hash mismatch")
	}
	if f.Initial.Turn != 0 || f.Initial.Result != result.Default || len(f.Initial.Hand) != card.InitialHandSize || len(f.Initial.Robot) != card.WinningTileCount || len(f.Initial.Fixed) != 0 || f.Initial.Limit < 1 || f.Initial.Condition != rules.InitialPlayerTurnParity || f.Initial.Kongs != 0 || f.Initial.Waiting || !f.Initial.Clear || len(f.Initial.ScoreItems) != 0 {
		return fmt.Errorf("invalid initial state")
	}
	counts := map[card.Card]int{}
	for _, group := range []card.Cards{f.Initial.Wall, f.Initial.Hand, f.Initial.Robot} {
		for _, tile := range group {
			counts[tile]++
		}
	}
	if len(counts) != len(card.CardsList) {
		return fmt.Errorf("initial state has invalid tile kinds")
	}
	for _, tile := range card.CardsList {
		if counts[tile] != card.CopiesPerTile {
			return fmt.Errorf("initial state must contain four %s tiles, got %d", tile, counts[tile])
		}
	}
	if len(f.Frames) == 0 || len(f.Frames) > MaxFrames {
		return fmt.Errorf("invalid frame count %d", len(f.Frames))
	}
	pending := map[int]bool{}
	nextID := 1
	for i, frame := range f.Frames {
		fail := func(msg string) error { return fmt.Errorf("invalid frame %d: %s", i+1, msg) }
		if frame.Hash != frame.State.Hash() {
			return fail("state hash mismatch")
		}
		if frame.State.Turn < 0 || frame.State.Turn > rules.MaxRecordedTurn || frame.State.Limit < 1 || frame.State.Kongs < 0 || frame.State.Kongs > card.WinningMeldCount || (frame.State.Condition != rules.InitialPlayerTurnParity && frame.State.Condition != rules.ClaimedPlayerTurnParity) || frame.State.Result < result.Default || frame.State.Result > result.Cancel {
			return fail("invalid state")
		}
		for _, group := range []card.Cards{frame.State.Wall, frame.State.Hand, frame.State.Robot, frame.State.Fixed} {
			if len(group) > card.TotalTiles {
				return fail("too many tiles")
			}
			for _, tile := range group {
				if !card.CardsList.Contains(tile) {
					return fail("unknown tile " + string(tile))
				}
			}
		}
		switch frame.Kind {
		case "decision":
			if frame.Request == nil || frame.Action == nil || frame.Request.ID != nextID || frame.Request.Turn != frame.State.Turn || frame.Accepted != nil || frame.Settlement != nil {
				return fail("malformed decision")
			}
			if err := frame.Action.Validate(frame.Request.Kind); err != nil {
				return fail(err.Error())
			}
			if frame.Request.Kind != "discard" && !card.CardsList.Contains(frame.Request.Tile) {
				return fail("confirmation missing tile")
			}
			pending[nextID] = true
			nextID++
		case "decision_result":
			if !pending[frame.DecisionID] || frame.Accepted == nil {
				return fail("unmatched decision result")
			}
			delete(pending, frame.DecisionID)
		case "draw", "robot_discard", "discard", "pung", "kong", "upgrade_kong":
			if !card.CardsList.Contains(frame.Tile) {
				return fail("event missing tile")
			}
		case "turn_end":
		case "settlement":
			if i != len(f.Frames)-1 || frame.Settlement == nil || frame.Settlement.Result == result.Default || frame.Settlement.Result != frame.State.Result || frame.Settlement.Counter != frame.State.Turn {
				return fail("invalid settlement")
			}
		default:
			return fail("unknown kind " + frame.Kind)
		}
	}
	if len(pending) != 0 || f.Frames[len(f.Frames)-1].Kind != "settlement" {
		return fmt.Errorf("recording missing decision results or settlement")
	}
	return nil
}

func Load(path string) (*File, error) { return LoadGame(path, "") }

func LoadGame(path, gameID string) (*File, error) {
	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	return ReadSource(in, gameID)
}

// Save publishes a complete file atomically and refuses to overwrite a recording.
func (f *File) Save(path string) error {
	if err := f.Validate(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	if len(data)+1 > MaxFileBytes {
		return fmt.Errorf("replay exceeds 32 MiB")
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".replay-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err = tmp.Write(append(data, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	if err = os.Link(tmp.Name(), path); err != nil {
		return fmt.Errorf("save replay (existing files are preserved): %w", err)
	}
	return nil
}
