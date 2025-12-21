package player

import (
	"io"
	"testing"

	c "github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

type scriptUI struct {
	inputs []string
	i      int
}

func (s *scriptUI) Info(string)    {}
func (s *scriptUI) Plain(string)   {}
func (s *scriptUI) Plainln(string) {}
func (s *scriptUI) PromptPlain(string) (string, error) {
	return s.next()
}
func (s *scriptUI) PromptInfo(string) (string, error) {
	return s.next()
}

func (s *scriptUI) next() (string, error) {
	if s.i >= len(s.inputs) {
		return "", io.EOF
	}
	v := s.inputs[s.i]
	s.i++
	return v, nil
}

func TestUserInput_YNAndInvalid(t *testing.T) {
	tr := language.New(utils.English)
	u := &scriptUI{inputs: []string{"Y", "NO", "maybe"}}

	p := NewPlayer(logging.NopLogger(), u, tr, utils.English)
	card := c.Card("B1")

	yes, ok := p.userInput(card, false, true)
	if !ok || !yes {
		t.Fatalf("expected Y => (true,true)")
	}
	no, ok := p.userInput(card, false, false)
	if !ok || no {
		t.Fatalf("expected NO => (false,true)")
	}
	_, ok = p.userInput(card, false, false)
	if ok {
		t.Fatalf("expected invalid => ok=false")
	}
}

func TestOutCard_EnterAutoPlaysNewCard(t *testing.T) {
	tr := language.New(utils.English)
	u := &scriptUI{inputs: []string{""}}

	p := NewPlayer(logging.NopLogger(), u, tr, utils.English)
	nc := c.Card("B1")
	p.PlayerCards = &c.Cards{nc, "W1", "W2"}

	p.outCard(nc)
	if p.PlayerCards.ContainsC(nc) {
		t.Fatalf("expected new card %s to be removed on empty input", nc)
	}
}
