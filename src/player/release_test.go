package player

import (
	"os"
	"reflect"
	"testing"
)

// Run complete games under deterministic semantic input, then regenerate and
// replay each game. Golden scenario tests separately cover special win/meld paths.
func TestReleaseSeededGames(t *testing.T) {
	if os.Getenv("MAHJONG_RELEASE_CHECK") != "1" {
		t.Skip("run scripts/check.sh release")
	}
	for seed := int64(1); seed <= 256; seed++ {
		recording := recordGame(t, seed, nil, nil)
		if err := VerifyReplay(recording); err != nil {
			t.Fatalf("seed %d replay: %v", seed, err)
		}
		repeated := recordGame(t, seed, nil, nil)
		if !reflect.DeepEqual(recording.Initial, repeated.Initial) || !reflect.DeepEqual(recording.Frames, repeated.Frames) {
			t.Fatalf("seed %d produced nondeterministic state or settlement", seed)
		}
	}
	t.Log("256 complete games regenerated and replayed successfully")
}
