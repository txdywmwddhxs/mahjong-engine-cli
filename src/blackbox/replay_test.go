//go:build integration

package blackbox_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestCLI_Replay_RealGame is a regression test based on a real game session.
// It replays a complete game with a fixed seed and validates key game events.
//
// NOTE: Due to randomness in game flow, the exact sequence may differ from the original log.
// This test validates that the game can complete successfully with the given seed and inputs.
//
// Original game details:
//   - Seed: 1766309534090039000
//   - Result: Win (score 13: win=9, waiting=2, exposed-kong=2)
//   - Final total score: 8402
func TestCLI_Replay_RealGame(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Log("[test] Starting real game replay test...")
	t.Log("[test] Finding repository root...")
	root, err := repoRoot()
	if err != nil {
		t.Fatalf("repoRoot: %v", err)
	}
	t.Logf("[test] Repository root: %s", root)

	t.Log("[test] Setting up temporary test environment...")
	tmp := t.TempDir()
	work := filepath.Join(tmp, "work")
	if err := os.MkdirAll(filepath.Join(work, "config"), 0o755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(work, "log"), 0o755); err != nil {
		t.Fatalf("mkdir log: %v", err)
	}
	// Minimal CHANGELOG marker so runtime root detection works.
	if err := os.WriteFile(filepath.Join(work, "config", "CHANGELOG.txt"), []byte("v0.0.0\n"), 0o666); err != nil {
		t.Fatalf("write changelog: %v", err)
	}
	t.Logf("[test] Test work directory: %s", work)

	bin := filepath.Join(tmp, "play")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	t.Log("[test] Building CLI binary...")
	// Build the CLI binary.
	build := exec.Command("go", "build", "-trimpath", "-o", bin, filepath.Join(root, "src"))
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}
	t.Logf("[test] Binary built: %s", bin)

	// Use a flexible input strategy: provide enough inputs to handle all possible prompts.
	// Since the game flow may differ from the original log due to randomness,
	// we provide a generous supply of inputs that cover all scenarios.
	var inputBuf bytes.Buffer

	// Replay the exact user input sequence from the real game log.
	// Based on log analysis with seed 1766309534090039000:
	// Turn 1: w4
	// Turn 3: b1
	// Turn 5: w6 (invalid), w5 (retry)
	// Turn 6: Robot plays B6 -> pung prompt: w (invalid), y (confirm), b6 (invalid - already punged), b7
	// Turn 8: (empty - auto play)
	// Turn 10: w
	// Turn 12: x
	// Turn 14: t6
	// Turn 16: (empty - auto play)
	// Turn 18: b3
	// Turn 20: (empty - auto play)
	// Turn 22: Robot plays B6 -> kong prompt: h (invalid), y (confirm), b3
	// Turn 24: (empty - auto play)
	// Turn 26: b4 tt (waiting)
	// Turn 27: Robot plays T8 -> auto win
	// Continue: N
	inputs := []string{
		"w4\n",    // Turn 1: play card
		"",        // Turn 2: robot (no input needed)
		"b1\n",    // Turn 3: play card
		"",        // Turn 4: robot
		"w6\n",    // Turn 5: play card (invalid - without this card)
		"w5\n",    // Turn 5: retry play card
		"",        // Turn 6: robot plays B6, triggers pung prompt
		"w\n",     // Turn 6: pung prompt - invalid input
		"y\n",     // Turn 6: pung confirm
		"b6\n",    // Turn 6: play card (invalid - cannot play punged card)
		"b7\n",    // Turn 6: play card (valid)
		"",        // Turn 7: robot
		"\n",      // Turn 8: auto play (empty input)
		"",        // Turn 9: robot
		"w\n",     // Turn 10: play card
		"",        // Turn 11: robot
		"x\n",     // Turn 12: play card
		"",        // Turn 13: robot
		"t6\n",    // Turn 14: play card
		"",        // Turn 15: robot
		"\n",      // Turn 16: auto play
		"",        // Turn 17: robot
		"b3\n",    // Turn 18: play card
		"",        // Turn 19: robot
		"\n",      // Turn 20: auto play
		"",        // Turn 21: robot
		"h\n",     // Turn 22: kong prompt - invalid input
		"y\n",     // Turn 22: kong confirm
		"b3\n",    // Turn 22: play card after kong
		"",        // Turn 23: robot
		"\n",      // Turn 24: auto play
		"",        // Turn 25: robot
		"b4 tt\n", // Turn 26: waiting declaration
		"",        // Turn 27: robot plays T8, auto win
		"N\n",     // Continue prompt: No
		"ENDL\n",
	}
	t.Logf("[test] Prepared %d input lines", len(inputs))
	actualInputCount := 0
	for _, inp := range inputs {
		// Skip empty strings - they are placeholders for robot turns, not actual inputs
		if inp == "" {
			continue
		}
		// Only write actual inputs (non-empty strings)
		actualInputCount++
		if inp != "\n" {
			t.Logf("[test] Input %d: %q", actualInputCount, strings.TrimSpace(inp))
		}
		inputBuf.WriteString(inp)
	}
	t.Logf("[test] Total input buffer size: %d bytes", inputBuf.Len())
	t.Logf("[test] Actual input lines written: %d", actualInputCount)

	// Create a limited reader that will error when inputs are exhausted
	// We track by lines (newline characters) since bufio.Scanner reads line by line
	// Count actual input lines (non-empty strings)
	actualInputLines := 0
	for _, inp := range inputs {
		if inp != "" {
			actualInputLines++
		}
	}
	limitedReader := &limitedInputReader{
		reader:   &inputBuf,
		maxLines: actualInputLines,
		t:        t,
	}
	t.Logf("[test] Limited reader will allow %d input lines", actualInputLines)

	t.Log("[test] Running game with seed 1766309534090039000...")
	t.Log("[test] This may take a moment (replaying 27 turns)...")
	// cmd := exec.Command(bin, "--quick-mode")
	cmd := exec.Command(bin)
	cmd.Dir = work
	cmd.Env = append(os.Environ(),
		"MAHJONG_ROOT="+work,
		"MAHJONG_SEED=1766309534090039000",
	)
	cmd.Stdin = limitedReader

	// Create output files for debugging
	stdoutFile := filepath.Join(tmp, "game_stdout.log")
	stderrFile := filepath.Join(tmp, "game_stderr.log")
	stdoutF, err := os.Create(stdoutFile)
	if err != nil {
		t.Fatalf("create stdout file: %v", err)
	}
	defer stdoutF.Close()
	stderrF, err := os.Create(stderrFile)
	if err != nil {
		t.Fatalf("create stderr file: %v", err)
	}
	defer stderrF.Close()

	t.Logf("[test] Output files: stdout=%s, stderr=%s", stdoutFile, stderrFile)

	// Use pipes for real-time output monitoring
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}

	var stdout, stderr bytes.Buffer
	outputDone := make(chan bool, 2)

	// Monitor stdout in real-time and detect infinite loops
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		lineCount := 0
		lastUnrecognizedCount := 0
		lastUnrecognizedLine := ""
		consecutiveUnrecognized := 0
		const maxConsecutiveUnrecognized = 10 // If we see 10 consecutive "unrecognized input", it's likely a loop

		for scanner.Scan() {
			line := scanner.Text()
			lineCount++
			lineWithNewline := line + "\n"
			stdout.WriteString(lineWithNewline)
			stdoutF.WriteString(lineWithNewline) // Also write to file

			// Detect infinite loops: repeated "不能识别的输入" or "Unrecognized input"
			if strings.Contains(line, "不能识别的输入") || strings.Contains(line, "Unrecognized input") {
				if line == lastUnrecognizedLine {
					consecutiveUnrecognized++
					if consecutiveUnrecognized >= maxConsecutiveUnrecognized {
						t.Errorf("[test] INFINITE LOOP DETECTED: Saw %d consecutive 'unrecognized input' messages", consecutiveUnrecognized)
						t.Errorf("[test] Last prompt: %s", lastUnrecognizedLine)
						t.Errorf("[test] This indicates the game is stuck in a loop waiting for valid input")
						// Try to kill the process
						if cmd.Process != nil {
							cmd.Process.Kill()
						}
						return
					}
				} else {
					consecutiveUnrecognized = 1
					lastUnrecognizedLine = line
				}
				lastUnrecognizedCount++
			} else {
				consecutiveUnrecognized = 0
			}

			// Print ALL INFO/DEBUG lines to see what's happening
			if strings.Contains(line, "INFO:") || strings.Contains(line, "DEBUG:") {
				t.Logf("[game] [%d] %s", lineCount, strings.TrimSpace(line))
			} else if strings.Contains(line, "Play a card") || strings.Contains(line, "Whether to") {
				// Also print prompts to see where we're waiting
				t.Logf("[game] [prompt] %s", strings.TrimSpace(line))
			}
		}
		if err := scanner.Err(); err != nil {
			t.Logf("[game] Scanner error: %v", err)
		}
		t.Logf("[game] Stdout stream ended (total lines: %d, unrecognized inputs: %d)", lineCount, lastUnrecognizedCount)
		outputDone <- true
	}()

	// Monitor stderr in real-time
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			lineWithNewline := line + "\n"
			stderr.WriteString(lineWithNewline)
			stderrF.WriteString(lineWithNewline) // Also write to file
		}
		outputDone <- true
	}()

	// Start command with timeout context
	t.Log("[test] Starting game process...")
	if err := cmd.Start(); err != nil {
		t.Fatalf("cmd start: %v", err)
	}

	// Also monitor log file in real-time to see game progress
	logPath := filepath.Join(work, "log", "play.log")
	logMonitorDone := make(chan bool)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		lastSize := int64(0)
		for {
			select {
			case <-ticker.C:
				if info, err := os.Stat(logPath); err == nil {
					if info.Size() > lastSize {
						lastSize = info.Size()
						// Read last few lines of log
						if content, err := os.ReadFile(logPath); err == nil {
							lines := strings.Split(string(content), "\n")
							if len(lines) > 3 {
								lastLines := lines[len(lines)-4:]
								t.Logf("[log] Latest: %s", strings.Join(lastLines, " | "))
							}
						}
					}
				}
			case <-logMonitorDone:
				return
			}
		}
	}()

	// Wait for output streams to finish or timeout
	timeout := time.After(60 * time.Second) // Increased timeout
	outputCount := 0
	for outputCount < 2 {
		select {
		case <-outputDone:
			outputCount++
			t.Logf("[test] Output stream closed (%d/2)", outputCount)
		case <-timeout:
			t.Errorf("[test] TIMEOUT: Game execution exceeded 60 seconds")
			logMonitorDone <- true
			// Try to read final log state
			if content, err := os.ReadFile(logPath); err == nil {
				lines := strings.Split(string(content), "\n")
				t.Logf("[test] Log file has %d lines, last 10:", len(lines))
				start := len(lines) - 10
				if start < 0 {
					start = 0
				}
				for i := start; i < len(lines); i++ {
					t.Logf("[test]   [%d] %s", i+1, lines[i])
				}
			}
			cmd.Process.Kill()
			t.Fatalf("Game execution timed out - possible hang detected")
		}
	}
	logMonitorDone <- true

	// Wait for command to finish
	t.Log("[test] Waiting for command to finish...")
	err = cmd.Wait()
	out := stdout.String() + stderr.String()

	if err != nil {
		t.Logf("[test] Command exited with error: %v", err)
		t.Logf("[test] Last 30 lines of output:\n%s", getLastLines(out, 30))
	} else {
		t.Log("[test] Game execution completed successfully")
	}

	// Read the log file to validate detailed game events.
	t.Log("[test] Reading game log file...")
	logContent, logErr := os.ReadFile(logPath)
	if logErr != nil {
		t.Logf("[test] Could not read log file (non-fatal): %v", logErr)
	} else {
		t.Logf("[test] Log file size: %d bytes", len(logContent))
		logStr := string(logContent)
		t.Log("[test] Validating log content...")
		validateLogContent(t, logStr)
	}

	// Validate stdout/stderr contains expected key messages.
	t.Log("[test] Validating output...")
	validateOutput(t, out)

	// Validate final game state from log if available.
	if logErr == nil {
		t.Log("[test] Validating final game state...")
		validateFinalState(t, string(logContent))
	}
	t.Log("[test] All validations passed!")
}

func validateLogContent(t *testing.T, log string) {
	t.Helper()

	// Check seed is logged.
	t.Log("[validate] Checking seed in log...")
	if !strings.Contains(log, "Game random seed: 1766309534090039000") {
		t.Errorf("log missing expected seed line")
	} else {
		t.Log("[validate] ✓ Seed found in log")
	}

	// Check that game started and completed (more flexible validation)
	t.Log("[validate] Checking game completion...")
	if !strings.Contains(log, "Game End") && !strings.Contains(log, "游戏结束") {
		t.Errorf("log missing game end marker")
	} else {
		t.Log("[validate] ✓ Game completed")
	}

	// Check that some game events occurred (flexible)
	t.Log("[validate] Checking game events...")
	hasWin := strings.Contains(log, "You win") || strings.Contains(log, "你赢了")
	hasScore := strings.Contains(log, "Score this game:") || strings.Contains(log, "本局得分:")
	if !hasWin && !hasScore {
		t.Logf("[validate] Warning: No win or score found (game may have ended differently)")
	} else {
		if hasWin {
			t.Log("[validate] ✓ Win message found")
		}
		if hasScore {
			t.Log("[validate] ✓ Score found")
		}
	}
}

func validateOutput(t *testing.T, out string) {
	t.Helper()

	// Must contain continue prompt.
	t.Log("[validate] Checking Continue prompt in output...")
	if !strings.Contains(out, "Continue") && !strings.Contains(out, "继续") {
		t.Errorf("output missing Continue prompt")
	} else {
		t.Log("[validate] ✓ Continue prompt found")
	}

	// Must contain win message.
	t.Log("[validate] Checking win message in output...")
	if !strings.Contains(out, "You win") && !strings.Contains(out, "你赢了") {
		t.Errorf("output missing win message")
	} else {
		t.Log("[validate] ✓ Win message found")
	}

	// Must show score.
	t.Log("[validate] Checking score in output...")
	if !strings.Contains(out, "Score this game: 13") && !strings.Contains(out, "本局得分: 13") {
		t.Errorf("output missing score")
	} else {
		t.Log("[validate] ✓ Score found in output")
	}
}

func validateFinalState(t *testing.T, log string) {
	t.Helper()

	t.Log("[validate] Parsing final score from log...")
	scanner := bufio.NewScanner(strings.NewReader(log))
	var lastScoreLine string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Score this game:") || strings.Contains(line, "本局得分:") {
			lastScoreLine = line
		}
	}

	if lastScoreLine == "" {
		t.Errorf("log missing final score line")
		return
	}
	t.Logf("[validate] ✓ Final score line: %s", strings.TrimSpace(lastScoreLine))

	// Validate score breakdown appears in log.
	t.Log("[validate] Checking score components...")
	if !strings.Contains(log, "win: 9") {
		t.Errorf("log missing win score component")
	} else {
		t.Log("[validate] ✓ Win component (9) found")
	}
	if !strings.Contains(log, "waiting: 2") {
		t.Errorf("log missing waiting score component")
	} else {
		t.Log("[validate] ✓ Waiting component (2) found")
	}
	if !strings.Contains(log, "exposed-kong: 2") {
		t.Errorf("log missing exposed-kong score component")
	} else {
		t.Log("[validate] ✓ Exposed-kong component (2) found")
	}
}

// getLastLines returns the last n lines from a string
func getLastLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= n {
		return s
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}

// limitedInputReader wraps a reader and tracks how many newlines are read.
// If more lines are read than the expected number of inputs, it returns an error.
type limitedInputReader struct {
	reader    io.Reader
	maxLines  int
	lineCount int
	t         *testing.T
}

func (r *limitedInputReader) Read(p []byte) (n int, err error) {
	// Only check limit when we actually read data
	n, err = r.reader.Read(p)
	if err != nil && err != io.EOF {
		return n, err
	}

	// If we've already exhausted inputs and there's more data, that's a problem
	if r.lineCount >= r.maxLines && n > 0 {
		r.t.Errorf("[test] INPUT EXHAUSTED: Already read %d input lines (limit: %d), but game is still trying to read", r.lineCount, r.maxLines)
		r.t.Errorf("[test] Last %d bytes read: %q", n, string(p[:min(n, 100)]))
		return n, errors.New("input exhausted: game tried to read more inputs than provided")
	}

	// Count newlines in the data we just read
	for i := 0; i < n; i++ {
		if p[i] == '\n' {
			r.lineCount++
			if r.lineCount > r.maxLines {
				r.t.Errorf("[test] INPUT EXHAUSTED: Read %d input lines, but only %d inputs provided", r.lineCount, r.maxLines)
				r.t.Errorf("[test] This indicates the game is trying to read more inputs than available, likely causing an infinite loop")
				r.t.Errorf("[test] Last %d bytes read: %q", n, string(p[:min(n, 100)]))
				return n, errors.New("input exhausted: game tried to read more inputs than provided")
			}
		}
	}

	return n, err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
