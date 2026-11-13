package player

import (
	"testing"
	"time"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func TestConsecutiveClaimsKeepRobotNext(t *testing.T) {
	for _, kong := range []bool{false, true} {
		name := "pung"
		hand := card.Cards{"B1", "B1", "T1", "T1", "W2", "W3", "W5", "W7", "B4", "B6", "T4", "T6", "X"}
		if kong {
			name = "kong"
			hand[4], hand[5] = "B1", "T1"
		}
		t.Run(name, func(t *testing.T) {
			initial := scenario(t, hand, card.Cards{}, card.Cards{"T1", "B1"})
			p := NewPlayerWithOptions(logging.NopLogger(), nil, language.New(utils.Chinese), utils.Chinese, Options{
				Initial: &initial, Input: replay.SourceFunc(autoAction), Sleep: func(time.Duration) {},
			})
			for _, turn := range []int{2, 3} {
				p.counter = turn
				p.playRobotTurn()
				if p.counter%2 == p.condition {
					t.Fatalf("after claim on turn %d: next turn would incorrectly draw for player (condition=%d)", turn, p.condition)
				}
			}
			if len(*p.FixedCards) != 6 && !kong || len(*p.FixedCards) != 8 && kong {
				t.Fatalf("expected two claims, fixed=%v", *p.FixedCards)
			}
		})
	}
}

// Exercise the real game loop: claims consume the player's discard, so the
// following action must be a robot discard, including after a second claim.
func TestReleaseClaimTurnOrder(t *testing.T) {
	for _, firstKong := range []bool{false, true} {
		for _, secondKong := range []bool{false, true} {
			name := "pung"
			if firstKong {
				name = "kong"
			}
			if secondKong {
				name += "_kong"
			} else {
				name += "_pung"
			}
			t.Run(name, func(t *testing.T) {
				hand := card.Cards{"B1", "B1", "T1", "T1", "W2", "W3", "W5", "W7", "B4", "B6", "T4", "T6", "X"}
				if firstKong {
					hand[4] = "B1"
				}
				if secondKong {
					hand[5] = "T1"
				}
				initial := scenario(t, hand, card.Cards{"W9", "T9", "B9"}, card.Cards{"W8", "T1", "B1"})
				initial.Limit = 6
				policy := replay.SourceFunc(func(req replay.Request, state replay.State) (replay.Action, error) {
					if req.Kind == "discard" && req.Tile == "" {
						tile := card.Card("W7")
						if req.Turn == 3 {
							tile = "W5"
						}
						return replay.Action{Kind: "discard", Tile: tile}, nil
					}
					return autoAction(req, state)
				})
				recording := recordGame(t, 1, &initial, policy)
				claims := 0
				robotTurns := []int{}
				for _, frame := range recording.Frames {
					if frame.Kind == "pung" || frame.Kind == "kong" {
						claims++
					}
					if frame.Kind == "robot_discard" {
						robotTurns = append(robotTurns, frame.State.Turn)
					}
					if frame.Kind == "draw" && frame.State.Turn == 4 {
						t.Fatal("player incorrectly drew immediately after second claim")
					}
				}
				if claims != 2 {
					t.Fatalf("expected two claims, got %d", claims)
				}
				if len(robotTurns) != 3 || robotTurns[0] != 2 || robotTurns[1] != 3 || robotTurns[2] != 4 {
					t.Fatalf("robot discard turns=%v, want [2 3 4]", robotTurns)
				}
				if err := VerifyReplay(recording); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
