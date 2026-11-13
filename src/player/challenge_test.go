package player

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

// Intn(100) consumes the upper 31 bits; subsequent draws select the range minimum.
type difficultySource struct{ roll int64 }

func (s *difficultySource) Seed(int64) {}
func (s *difficultySource) Int63() int64 {
	value := s.roll << 32
	s.roll = 0
	return value
}

type difficultyLogger struct {
	logging.Logger
	messages []string
}

func (l *difficultyLogger) Debug(message string, _ int) {
	l.messages = append(l.messages, message)
}

func TestGenerateLimitProbabilities(t *testing.T) {
	counts := map[int]int{}
	for roll := 0; roll < 100; roll++ {
		logger := &difficultyLogger{Logger: logging.NopLogger()}
		p := &Player{
			rng:    rand.New(&difficultySource{roll: int64(roll)}),
			logger: logger,
			t:      language.New(utils.Chinese),
		}
		p.generateLimit()
		want := 60
		if roll < 32 {
			want = 22
		} else if roll < 52 {
			want = 35
		}
		if p.limit != want {
			t.Fatalf("roll %d: limit=%d, want %d", roll, p.limit, want)
		}
		counts[p.limit]++
		if !strings.HasSuffix(logger.messages[0], ": 32%, 20%, 48%") {
			t.Fatalf("incorrect probability log: %s", logger.messages[0])
		}
	}
	if counts[22] != 32 || counts[35] != 20 || counts[60] != 48 {
		t.Fatalf("incorrect distribution: %v", counts)
	}
}
