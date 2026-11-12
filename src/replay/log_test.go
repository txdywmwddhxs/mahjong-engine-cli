package replay_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/logging"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/player"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/replay"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/ui"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

func recordedBlock(t *testing.T, seed int64) (*replay.File, string) {
	t.Helper()
	p := player.NewPlayerWithOptions(logging.NopLogger(), ui.NewConsole(strings.NewReader("WHOSYOURDADDY\n"), io.Discard), language.New(utils.English), utils.English, player.Options{Seed: &seed, Sleep: func(time.Duration) {}})
	recording := p.StartRecording()
	p.Main()
	if p.Err() != nil {
		t.Fatal(p.Err())
	}
	if err := recording.Finish(); err != nil {
		t.Fatal(err)
	}
	block, err := recording.File.LogBlock()
	if err != nil {
		t.Fatal(err)
	}
	return recording.File, block
}

func TestLogReplayCopyAndSelection(t *testing.T) {
	first, block1 := recordedBlock(t, 1)
	second, block2 := recordedBlock(t, 2)
	log := "v9.1.0\n\nDEBUG[0]: readable game details\n" + replay.LogStartPrefix + first.GameID + "\n" + block1 + "INFO[3]: total score\n" + replay.LogStartPrefix + second.GameID + "\n" + block2
	cases := []struct{ name, input, id, want string }{
		{"copied_block", block1, "", first.GameID},
		{"windows_line_endings", strings.ReplaceAll(block1, "\n", "\r\n"), "", first.GameID},
		{"complete_log_latest", log, "", second.GameID},
		{"select_first", log, first.GameID, first.GameID},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := replay.ReadSource(strings.NewReader(tc.input), tc.id)
			if err != nil {
				t.Fatal(err)
			}
			if f.GameID != tc.want {
				t.Fatalf("game %s, wanted %s", f.GameID, tc.want)
			}
			if err = player.VerifyReplay(f); err != nil {
				t.Fatal(err)
			}
		})
	}
	// Never replay an older game when the latest one crashed or was cut off.
	broken := log + replay.LogStartPrefix + "unfinished-game\nINFO[1]: new draw\n"
	if _, err := replay.ReadSource(strings.NewReader(broken), ""); err == nil || !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("silently fell back: %v", err)
	}
	f, err := replay.ReadSource(strings.NewReader(broken), first.GameID)
	if err != nil || f.GameID != first.GameID {
		t.Fatalf("could not select older completed game: %v", err)
	}
}

func TestLogReplayRejectsCorruptionAndPartialCopies(t *testing.T) {
	file, block := recordedBlock(t, 1)
	lines := strings.Split(block, "\n")
	badPayload := append([]string{}, lines...)
	badPayload[1] = "!" + badPayload[1][1:]
	cases := []struct{ name, input, want string }{
		{"no_replay", "INFO[1]: old text-only log\n", "no replay block"},
		{"truncated", strings.ReplaceAll(block, replay.LogBlockEnd, ""), "truncated"},
		{"missing_begin", strings.Join(lines[1:], "\n"), "without begin"},
		{"corrupt_base64", strings.Join(badPayload, "\n"), "corrupt"},
		{"mismatched_id", strings.Replace(block, file.GameID, "wrong-id", 1), "mismatch"},
		{"unsupported_format", strings.Replace(block, "REPLAY V1 ", "REPLAY V2 ", 1), "unsupported"},
		{"incomplete_status", block + replay.LogIncompletePrefix + file.GameID + " EOF\n", "incomplete"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := replay.ReadSource(strings.NewReader(tc.input), ""); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("wanted %s, got %v", tc.want, err)
			}
		})
	}
	if _, err := replay.ReadSource(strings.NewReader(block), "missing-game"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatal(err)
	}
}
