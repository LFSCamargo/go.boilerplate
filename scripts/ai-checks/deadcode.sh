#!/usr/bin/env bash
# Dead code gate — golang.org/x/tools/cmd/deadcode (Go official toolset).
source "$(dirname "$0")/lib/common.sh"

echo "== Dead code checker (deadcode) =="

if ! require_cmd deadcode; then
  log_warn "install deadcode: go install golang.org/x/tools/cmd/deadcode@latest"
  exit_with_status
fi

output=""
if ! output=$(deadcode -test ./... 2>&1); then
  echo "$output"
  log_fail "deadcode analysis failed"
  exit_with_status
fi

if [[ -n "$output" ]]; then
  echo "$output"
  log_fail "unreachable functions found — remove dead code or reach them from main/tests"
fi

exit_with_status
