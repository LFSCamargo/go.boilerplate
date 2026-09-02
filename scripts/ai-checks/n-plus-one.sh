#!/usr/bin/env bash
# Detect likely GORM N+1 query patterns in Go code.
source "$(dirname "$0")/lib/common.sh"

echo "== N+1 SQL detector =="

collect_to_array files changed_go_files
if [[ ${#files[@]} -eq 0 ]]; then
  log_ok "no changed Go files"
  exit 0
fi

for file in "${files[@]}"; do
  [[ -f "$file" ]] || continue

  has_loop=false
  has_query=false
  has_mitigation=false

  grep -q 'for .*range' "$file" && has_loop=true
  grep -qE '\.(Find|First|Take|Where)\(' "$file" && has_query=true
  grep -qE 'Preload|Joins|Raw\(' "$file" && has_mitigation=true

  if [[ "$has_loop" == true && "$has_query" == true && "$has_mitigation" == false ]]; then
    log_fail "$file: possible N+1 — loop + GORM query without Preload/Joins/batch"
  fi
done

if [[ "$failures" -eq 0 ]]; then
  log_ok "no N+1 patterns detected in changed files"
fi
exit_with_status
