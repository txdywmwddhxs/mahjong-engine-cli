# mahjong-engine-cli

A lightweight **CLI Mahjong** game written in Go.

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

> Note: the project currently uses relative paths (see `src/utils/constants.go`) for config/log locations.

## Gameplay notes (current design)

This project historically implements a **challenge-style “robot”**:

- The “robot hand” is generated at game start and is revealed on player loss.
- The game may use a **turn limit** mechanism (reach a certain number of turns and the player loses).

These mechanics are part of the current gameplay and will be revisited when implementing **fair PvP-style AI**.

## Config / Logs

- Config file (runtime stats): `config/config.json`  
  - It is intentionally **ignored by git** (see `.gitignore`).
- Changelog: `config/ChangeLog`
- Current log: `log/play.log`
- History logs: `log/history_log/`

## License

MIT License. See `LICENSE`.


