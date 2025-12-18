#!/usr/bin/env bash
set -eu
# pipefail is supported by bash and zsh, but not by all shells.
set -o pipefail 2>/dev/null || true

find_root() {
  # Prefer git root if available.
  local root=""
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    echo "$root"
    return 0
  fi

  # Fallback: walk up from current directory to find go.mod.
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

  # Best-effort fallback.
  echo "$PWD"
  return 0
}

ROOT="$(find_root)"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/package.sh test
  ./scripts/package.sh release

Modes:
  test    Build a test package with Version=dev
  release Build a release package with Version read from ./VERSION

Environment overrides:
  GOOS / GOARCH  Cross-compile target (defaults to host go env)
  OUTDIR         Output directory (defaults to dist/play_<mode>_<ver>_<goos>_<goarch>)
EOF
}

MODE="${1:-}"
if [[ "$MODE" != "test" && "$MODE" != "release" ]]; then
  usage
  exit 2
fi

if [[ "$MODE" == "test" ]]; then
  VER="dev"
else
  if [[ ! -f "$ROOT/VERSION" ]]; then
    echo "ERROR: VERSION file not found at $ROOT/VERSION" >&2
    exit 1
  fi
  VER="$(head -n 1 "$ROOT/VERSION" | tr -d '\r' | xargs)"
  if [[ -z "$VER" ]]; then
    echo "ERROR: VERSION file is empty" >&2
    exit 1
  fi
fi

GOOS="${GOOS:-$(go env GOOS)}"
GOARCH="${GOARCH:-$(go env GOARCH)}"

PKG_NAME="play_${MODE}_${VER}_${GOOS}_${GOARCH}"
OUTDIR="${OUTDIR:-$ROOT/dist/$PKG_NAME}"

BIN_NAME="play"
if [[ "$GOOS" == "windows" ]]; then
  BIN_NAME="play.exe"
fi

rm -rf "$OUTDIR"
mkdir -p "$OUTDIR"

LDFLAGS="-s -w -X github.com/txdywmwddhxs/mahjong-engine-cli/src/utils.Version=$VER"

echo "Building $MODE package..."
echo "  Version: $VER"
echo "  Target : $GOOS/$GOARCH"
echo "  OutDir : $OUTDIR"

GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 \
  go build -trimpath -ldflags "$LDFLAGS" -o "$OUTDIR/$BIN_NAME" "$ROOT/src"

# Runtime files (minimal set required by current path detection + changelog display).
mkdir -p "$OUTDIR/config"
cp "$ROOT/config/CHANGELOG.txt" "$OUTDIR/config/CHANGELOG.txt"

# Convenience: include license/readme if present.
if [[ -f "$ROOT/LICENSE" ]]; then
  cp "$ROOT/LICENSE" "$OUTDIR/LICENSE"
fi
if [[ -f "$ROOT/README.md" ]]; then
  cp "$ROOT/README.md" "$OUTDIR/README.md"
fi

echo "Done."
echo "Binary: $OUTDIR/$BIN_NAME"

