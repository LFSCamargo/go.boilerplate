#!/usr/bin/env bash
# Spaghetti detector — mixed contexts in one file (handler + db + routes).
source "$(dirname "$0")/lib/common.sh"

echo "== Spaghetti detector =="

collect_to_array files changed_go_files
if [[ ${#files[@]} -eq 0 ]]; then
  log_ok "no changed Go files"
  exit 0
fi

score_file() {
  local file=$1
  local score=0
  local contexts=0

  grep -q 'fiber\.Ctx\|fiber\.Router\|func Register' "$file" && contexts=$((contexts + 1))
  grep -qE 'gorm\.|\.DB\.|db\.|\.Preload\(|\.Find\(|\.Create\(' "$file" && contexts=$((contexts + 1))
  grep -qE 'router\.(Get|Post|Put|Delete|Patch)\(' "$file" && contexts=$((contexts + 1))
  grep -qE 'smtp\.|jwt\.|bcrypt\.|mail\.Send' "$file" && contexts=$((contexts + 1))

  # handlers/ and router.go should not contain gorm
  if [[ "$file" == *handlers/* ]] && grep -qE 'gorm|\.DB\.|db\.' "$file"; then
    log_fail "$file: handler contains DB logic — move to service/repository"
    score=$((score + 2))
  fi

  if [[ "$file" == *router.go ]]; then
    if grep -qE 'gorm|\.DB\.|return c\.JSON|return c\.Send' "$file"; then
      log_fail "$file: router contains business/DB logic — keep routing only"
    fi
    return 0
  fi

  if [[ "$contexts" -ge 3 ]]; then
    log_fail "$file: $contexts contexts mixed (HTTP + DB + routes + infra) — split files"
    score=$((score + 1))
  fi

  return "$score"
}

for file in "${files[@]}"; do
  [[ -f "$file" ]] || continue
  score_file "$file" || true
done

if [[ "$failures" -eq 0 ]]; then
  log_ok "no spaghetti patterns in changed files"
fi
exit_with_status
