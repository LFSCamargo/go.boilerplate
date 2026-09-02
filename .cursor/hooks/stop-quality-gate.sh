#!/usr/bin/env bash
# Cursor hook: full quality gate when agent stops — inject follow-up if checks fail.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

cat >/dev/null

OUTPUT=$(mktemp)
if bash "$ROOT/scripts/ai-checks/run-all.sh" full >"$OUTPUT" 2>&1; then
  echo '{}'
  rm -f "$OUTPUT"
  exit 0
fi

SUMMARY=$(tail -30 "$OUTPUT" | sed 's/"/\\"/g' | tr '\n' ' ')
rm -f "$OUTPUT"

cat <<EOF
{
  "followup_message": "Quality gates failed. Run \`make ai-check\` and fix all issues before finishing. Recent output: ${SUMMARY} Also ensure docs/DESIGN_DOC.md is updated if you changed API, models, or modules."
}
EOF
exit 0
