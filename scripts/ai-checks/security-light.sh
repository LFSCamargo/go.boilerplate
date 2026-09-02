#!/usr/bin/env bash
# Lightweight security checks — fast, runs on every edit/stop.
source "$(dirname "$0")/lib/common.sh"

echo "== Security checker (lightweight) =="

# Hardcoded secret patterns
while IFS= read -r file; do
  [[ -f "$file" ]] || continue
  [[ "$file" == *.example ]] && continue
  [[ "$file" == *go.sum ]] && continue
  if grep -nEi '(password|secret|api[_-]?key)\s*=\s*["\x27][^"\x27]{8,}' "$file" 2>/dev/null; then
    log_fail "$file: possible hardcoded secret"
  fi
done < <(changed_files)

# Go vet on changed packages
if require_cmd go; then
  collect_to_array go_files changed_go_files
  if [[ ${#go_files[@]} -gt 0 ]]; then
    pkgs=$(printf '%s\n' "${go_files[@]}" | xargs -n1 dirname | sort -u | sed 's|^|./|' )
    for pkg in $pkgs; do
      if ! go vet "$pkg/..." 2>/dev/null; then
        log_fail "go vet failed for $pkg"
      fi
    done
  fi
fi

# gosec quick pass if installed
if require_cmd gosec; then
  collect_to_array go_files changed_go_files
  if [[ ${#go_files[@]} -gt 0 ]]; then
    if ! gosec -quiet -exclude=G104 "${go_files[@]}" 2>/dev/null; then
      log_fail "gosec reported issues (light scan)"
    fi
  fi
else
  log_warn "install gosec: go install github.com/securego/gosec/v2/cmd/gosec@latest"
fi

if [[ "$failures" -eq 0 ]]; then
  log_ok "lightweight security checks passed"
fi
exit_with_status
