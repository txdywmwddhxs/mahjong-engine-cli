# mahjong-engine-cli

A lightweight **CLI Mahjong** game written in Go.

项目故事：[一个麻将游戏的六年](docs/project-story.md) — 从 2020 年的离线命令行小游戏，到 2026 年的重新出发。

This project currently focuses on a single-player experience with scoring, logging, and multiple win patterns.  
**Fair PvP-style AI is a future goal**; the current “robot” behavior is part of the existing gameplay design (challenge-style).

## Features

- **Mahjong hand evaluation**
  - Standard 4 melds + 1 pair
  - Seven Pairs (七小对)
  - Thirteen Orphans (十三幺)
  - 1–9 Straight in one suit (一条龙 / “continuous line”)
- **Melds**
  - Pung (碰)
  - Kong (杠): exposed / concealed (明杠 / 暗杠), plus “pung to kong” (加杠)
- **Ting mode (听牌)**
  - `TING` / `TT` flow with validation and display of winning tiles
- **Scoring system**
  - Per-game scoring + cumulative stats persisted to config
- **Logging**
  - Game events and stats are written to log files
  - Historical logs are archived when the version changes
- **Bilingual prompts**
  - Chinese / English

## How to run

This repository contains source code only (recommended). Build with Go:

```bash
cd src
go run ./main.go
```

If you prefer a binary:

```bash
cd src
go build -o ../bin/play ./main.go
../bin/play
```

> Runtime paths are resolved from the detected project or package root. Set `MAHJONG_ROOT` to override it. Normal play creates the required directories; replay viewing and verification do not create runtime files.

## Gameplay notes (current design)

This project historically implements a **challenge-style “robot”**:

- The “robot hand” is generated at game start and is revealed on player loss.
- The game may use a **turn limit** mechanism (reach a certain number of turns and the player loses).

These mechanics are part of the current gameplay and will be revisited when implementing **fair PvP-style AI**.

## Config / Logs

- Config file (runtime stats): `config/config.json`  
  - It is intentionally **ignored by git** (see `.gitignore`).
- Version source: root `VERSION` file
- Changelog (update history only): `config/CHANGELOG.txt`
- Current log: `log/play.log`
- History logs: `log/history_log/`

Source runs and ordinary builds read the root `VERSION` file. Packaged builds
use a version injected at build time: `scripts/package.sh release` reads `VERSION`,
while `scripts/package.sh test` uses `dev`. Packages include a matching `VERSION`
file. Build-time injection takes priority; without injection or a readable,
nonempty `VERSION`, the application reports `dev`. Changelog contents never
determine the running version.

The console prints plain scrolling text without clearing or redrawing the terminal.

## Recording and replay

Run these commands from the repository root:

```bash
# Play normally; each completed round is automatically embedded in log/play.log.
go run ./src --quick-mode

# Optionally save separate JSON files as well.
go run ./src --record game.json

# Re-execute decisions and check every event, state and settlement.
go run ./src --replay game.json --verify

# Inspect recorded snapshots without executing game rules.
go run ./src --replay game.json

# Navigate snapshots interactively.
go run ./src --replay game.json --step

# Read the latest round directly from the log.
go run ./src --replay log/play.log --verify

# On macOS: copy a complete replay block from the log, then verify the clipboard.
pbpaste | go run ./src --replay - --verify
```

No recording flag is required for log replay. Readable log entries include the
seed, hands, decision IDs/types/actions, acceptance status and score. Each completed
round ends with a compact gzip/base64 replay block containing the full JSON record:

```text
-----BEGIN MAHJONG REPLAY V1 GAME_ID-----
... replay data ...
-----END MAHJONG REPLAY-----
```

Copy the entire block, including both markers, to share a reproducible game or save
it as a text file accepted by `--replay`. A copied whole round with its surrounding
readable log entries also works. The data is compressed so you do not need to copy
hundreds of repeated state snapshots. The normal terminal output stays readable.

A full log selects its latest recorded round by default. Add `--game GAME_ID` to
select another round, for example `--replay log/play.log --game GAME_ID --step`.
An unfinished latest round or a truncated/corrupt copied block produces an error;
it never silently selects an earlier completed round. EOF is logged as incomplete,
without updating stats. Older completed rounds remain selectable by ID. Native JSON
recordings continue to work. `--replay -` reads from stdin; use a file for interactive
`--step` navigation because stdin would otherwise also be needed for controls.

In step mode: Enter or `n` advances, `p` goes back, `g 12` jumps to frame 12,
`t 6` jumps to turn 6, and `q` exits. Frame 0 is the initial state.

Continuing a recorded session writes `game.2.json`, `game.3.json`, etc. Paths
are relative to the current working directory. The parent directory must exist;
existing files are never overwritten. Recording a normal game still updates normal
stats. `END` cancels a round and saves its cancellation without adding a result to
stats. EOF or an input failure aborts the round without saving an incomplete replay
or updating stats.

Each versioned JSON file contains the complete initial wall and hands, the actual
per-round seed, rule/code versions, numbered decision requests, typed actions,
accepted/rejected outcomes, semantic events, state snapshots with SHA-256 hashes,
and the final score breakdown. Rejected attempts are explicit decision/result
pairs; they never silently consume the next decision. The fixed seed controls new
round generation; replay restores the saved setup instead of reshuffling. A second
round uses a new per-round seed, so it can also be reproduced independently.

Verification has no prompts or delays and returns a nonzero exit status at the
first divergence, identifying the frame, turn and differing field. It also rejects
truncated, incomplete, malformed or incompatible records and unmatched decisions.
Format v1 is limited to 10,000 frames / 32 MiB. Hashes detect state corruption;
they are not signatures or proof that a record is authentic.

Viewing and verification do not write configuration, logs, stats or runtime
directories, including at process startup. Viewing uses snapshots and can display
records from other rule versions; verification requires the current rule version.
Old text logs / seed-only input scripts are not complete replay records and are
not automatically converted. Changing the UI language or terminal layout does not
affect replay verification.

Run regression checks with:

```bash
go test ./...
go test -tags integration ./src/blackbox -count=1 -timeout=30s
```

## Release acceptance check

Run one command before publishing:

```bash
./scripts/check.sh release
```

Requires Git, Bash, Go 1.24 or newer, and a C compiler for Go's race detector. The
command runs formatting and vet checks, uncached tests with the race detector,
CLI integration and golden replay tests, the 2,190,387-hand differential
experiment, 256 complete seeded games (including regeneration and replay), and
actual release/test package checks. It also verifies version/changelog
consistency and requires at least 50% statement coverage; `COVER_MIN` cannot
lower this threshold in release mode.

Package checks build and run binaries for the current host platform. They check
packaged files and version injection, startup outside the source directory,
multiple rounds, cancellation statistics, and read-only JSON/log replay
verification. All game data and package fixtures use temporary directories;
coverage is written to `dist/coverage.out` (or `COVER_OUT`). Normal saves and logs
are not used. Cross-platform runtime verification remains separate work.

Any failed stage exits nonzero. Only a fully successful run prints
`RELEASE CHECK PASSED`. Passing means the current release acceptance checks
succeeded; it does not prove the absence of all bugs. This is a correctness gate,
not a performance benchmark.

## License

MIT License. See `LICENSE`.


