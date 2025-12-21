#!/usr/bin/env bash
set -eu
set -o pipefail 2>/dev/null || true

find_root() {
  local root=""
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    echo "$root"
    return 0
  fi
  local dir="$PWD"
  while :; do
    if [[ -f "$dir/go.mod" ]]; then
      echo "$dir"
      return 0
    fi
    local parent
    parent="$(cd "$dir/.." && pwd)"
    if [[ "$parent" == "$dir" ]]; then
      break
    fi
    dir="$parent"
  done
  echo "$PWD"
}

usage() {
  cat <<'EOF'
Usage:
  ./scripts/check.sh [quick|full]

Modes:
  quick  Fast local checks (default): gofmt check + go vet + go test ./...
  full   Stricter checks:
         - gofmt check
         - go vet
         - go test -count=1 -race ./...
         - go test -count=1 -tags=integration ./...   (blackbox regression)
         - coverage gate >= 50% (override via COVER_MIN)

Notes:
  - This script does NOT modify files; it fails if gofmt would change anything.
  - Coverage is computed with: go test -coverpkg=./... -coverprofile=... ./...
EOF
}

MODE="${1:-quick}"
if [[ "$MODE" != "quick" && "$MODE" != "full" ]]; then
  usage
  exit 2
fi

ROOT="$(find_root)"
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

if [[ "$MODE" == "full" ]]; then
  echo "[check] go test (no cache, race)"
  go test -count=1 -race ./...

  echo "[check] go test (blackbox regression, integration tag)"
  go test -count=1 -tags=integration ./...

  COVER_MIN="${COVER_MIN:-50}"
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

echo "[check] OK"


