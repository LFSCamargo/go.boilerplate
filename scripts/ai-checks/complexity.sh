#!/usr/bin/env bash
# Complexity gate — blocks merge-unsafe code for production scale.
source "$(dirname "$0")/lib/common.sh"

MAX_FILE_LINES=300
MAX_FUNC_LINES=80
MAX_CYCLO=15

echo "== Complexity checker =="

collect_to_array files changed_go_files
if [[ ${#files[@]} -eq 0 ]]; then
  log_ok "no changed Go files"
  exit 0
fi

for file in "${files[@]}"; do
  [[ -f "$file" ]] || continue
  lines=$(wc -l < "$file" | tr -d ' ')
  if [[ "$lines" -gt "$MAX_FILE_LINES" ]]; then
    log_fail "$file: $lines lines (max $MAX_FILE_LINES) — split into smaller files"
  fi
done

# gocyclo if available
if require_cmd gocyclo; then
  for file in "${files[@]}"; do
    [[ -f "$file" ]] || continue
    while IFS= read -r line; do
      complexity=$(echo "$line" | awk '{print $1}')
      name=$(echo "$line" | awk '{print $2}')
      if [[ "$complexity" -gt "$MAX_CYCLO" ]]; then
        log_fail "$file: $name cyclomatic complexity $complexity (max $MAX_CYCLO)"
      fi
    done < <(gocyclo -over "$MAX_CYCLO" "$file" 2>/dev/null || true)
  done
else
  log_warn "install gocyclo: go install github.com/fzipp/gocyclo/cmd/gocyclo@latest"
fi

# golangci-lint complexity linter if configured
if require_cmd golangci-lint; then
  if ! golangci-lint run --disable-all -E gocyclo,funlen,gocognit "${files[@]}" 2>/dev/null; then
    log_fail "golangci-lint complexity thresholds exceeded"
  fi
fi

if [[ "$failures" -eq 0 ]]; then
  log_ok "complexity within limits"
fi
exit_with_status
