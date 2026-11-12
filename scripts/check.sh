#!/usr/bin/env bash
set -eu
set -o pipefail 2>/dev/null || true

usage() {
  cat <<'EOF'
Usage:
  ./scripts/check.sh [quick|full|release]

Modes:
  quick  Fast local checks (default): gofmt check + go vet + go test ./...
  full   Stricter checks:
         - gofmt check
         - go vet
         - go test -count=1 -short=false -race -timeout=5m ./...
         - go test -count=1 -tags=integration ./... (blackbox regression)
         - coverage gate >= 50% (override via COVER_MIN)

  release Release gate: full checks + 2,190,387-hand differential experiment,
          256 seeded game replays, native release/test package smoke tests,
          and version/changelog consistency. Coverage must be >= 50%.

Notes:
  - Requires Git, Go and a C compiler for the race detector; package checks use Bash.
  - Does not rewrite source or normal runtime data. full/release write coverage
    under dist/; game and package fixtures are isolated in temporary directories.
  - Coverage is computed with: go test -coverpkg=./... -coverprofile=... ./...
EOF
}

MODE="${1:-quick}"
if [[ "$MODE" != "quick" && "$MODE" != "full" && "$MODE" != "release" ]]; then
  usage
  exit 2
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "[check] mode=$MODE root=$ROOT"

echo "[check] gofmt (no changes allowed)"
UNFORMATTED="$(gofmt -l ./src || true)"
if [[ -n "$UNFORMATTED" ]]; then
  echo "ERROR: gofmt would change these files:" >&2
  echo "$UNFORMATTED" >&2
  echo "Run: gofmt -w src" >&2
  exit 1
fi

echo "[check] go vet"
go vet ./...

if [[ "$MODE" == "full" || "$MODE" == "release" ]]; then
  echo "[check] go test (no cache, race)"
  go test -count=1 -short=false -race -timeout=5m ./...

  echo "[check] go test (blackbox regression, integration tag)"
  if [[ "$MODE" == "release" ]]; then
    MAHJONG_WIN_EXPERIMENT=1 MAHJONG_RELEASE_CHECK=1 \
      go test -count=1 -short=false -tags=integration -timeout=10m ./...
  else
    go test -count=1 -short=false -tags=integration -timeout=5m ./...
  fi

  COVER_MIN="${COVER_MIN:-50}"
  if [[ "$MODE" == "release" ]]; then COVER_MIN=50; fi
  COVER_OUT="${COVER_OUT:-$ROOT/dist/coverage.out}"
  mkdir -p "$(dirname "$COVER_OUT")"

  echo "[check] coverage (min=${COVER_MIN}%)"
  go test -count=1 -coverpkg=./... -coverprofile="$COVER_OUT" ./... >/dev/null
  COVER_TOTAL="$(go tool cover -func="$COVER_OUT" | tail -n 1 | awk '{print $3}' | tr -d '%')"
  echo "[check] coverage total=${COVER_TOTAL}% (profile: $COVER_OUT)"
  awk -v got="$COVER_TOTAL" -v min="$COVER_MIN" 'BEGIN { exit !(got+0 >= min+0) }' || {
    echo "ERROR: coverage ${COVER_TOTAL}% is below minimum ${COVER_MIN}%" >&2
    exit 1
  }
else
  echo "[check] go test"
  go test ./...
fi

if [[ "$MODE" == "release" ]]; then
  echo "[check] RELEASE CHECK PASSED ($(go env GOHOSTOS)/$(go env GOHOSTARCH))"
else
  echo "[check] OK"
fi


