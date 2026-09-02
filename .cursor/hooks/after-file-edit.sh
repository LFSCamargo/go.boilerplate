#!/usr/bin/env bash
# Cursor hook: run light AI checks after file edits.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

# Consume stdin (hook payload)
cat >/dev/null

if [[ -x "$ROOT/scripts/ai-checks/run-all.sh" ]]; then
  bash "$ROOT/scripts/ai-checks/run-all.sh" light || {
    echo '{"additional_context": "Light AI quality gate failed. Run: make ai-check-light. Fix N+1, spaghetti, security-light, or missing tests before continuing."}' 
    exit 0
  }
fi

echo '{}'
exit 0
