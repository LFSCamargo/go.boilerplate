#!/usr/bin/env bash
# Unified quality gate — same entry point for Cursor, Claude, Codex, and CI.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODE="${1:-full}"

run() {
  echo ""
  bash "$SCRIPT_DIR/$1" || exit 1
}

case "$MODE" in
  light)
    run n-plus-one.sh
    run spaghetti.sh
    run security-light.sh
    run test-driven.sh
    ;;
  heavy)
    run security-heavy.sh
    ;;
  docs)
    run docs-sync.sh
    ;;
  full|*)
    run n-plus-one.sh
    run complexity.sh
    run spaghetti.sh
    run security-light.sh
    run test-driven.sh
    run deadcode.sh
    run docs-sync.sh
    run security-heavy.sh
    ;;
esac

echo ""
echo "All AI quality gates passed ($MODE)."
