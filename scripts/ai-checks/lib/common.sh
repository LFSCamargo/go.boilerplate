#!/usr/bin/env bash
# Shared helpers for AI quality gates (Cursor, Claude, Codex, CI).

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$ROOT"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

failures=0

log_ok()   { echo -e "${GREEN}[ok]${NC} $*"; }
log_warn() { echo -e "${YELLOW}[warn]${NC} $*"; }
log_fail() { echo -e "${RED}[fail]${NC} $*"; failures=$((failures + 1)); }

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    log_warn "$1 not installed — skipping dependent checks (install for full gate)"
    return 1
  fi
  return 0
}

# Portable replacement for mapfile (macOS bash 3.x)
read_lines() {
  local line
  while IFS= read -r line; do
    [[ -n "$line" ]] && echo "$line"
  done
}

# Emits changed file paths. Honors CI_DIFF_BASE (GitHub Actions PR/push base SHA).
git_changed_name_lines() {
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    return 0
  fi

  if [[ -n "${CI_DIFF_BASE:-}" && "${CI_DIFF_BASE}" != "0000000000000000000000000000000000000000" ]]; then
    if git rev-parse --verify "${CI_DIFF_BASE}^{commit}" >/dev/null 2>&1; then
      git diff --name-only --diff-filter=ACMRTUXB "${CI_DIFF_BASE}...HEAD" 2>/dev/null || true
      return 0
    fi
  fi

  {
    git diff --name-only --diff-filter=ACMRTUXB HEAD 2>/dev/null || true
    git diff --cached --name-only --diff-filter=ACMRTUXB 2>/dev/null || true
    git ls-files --others --exclude-standard 2>/dev/null || true
  }
}

changed_go_files() {
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    git_changed_name_lines | grep '\.go$' | grep -v '_test\.go$' | sort -u | read_lines
  else
    find src -name '*.go' ! -name '*_test.go' 2>/dev/null | read_lines
  fi
}

changed_files() {
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    git_changed_name_lines | sort -u | read_lines
  else
    find . -type f \
      ! -path './.git/*' \
      ! -path './node_modules/*' 2>/dev/null | read_lines
  fi
}

collect_to_array() {
  local var_name=$1
  shift
  local -a collected=()
  local line
  while IFS= read -r line; do
    [[ -n "$line" ]] && collected+=("$line")
  done < <("$@")
  # bash 3.2 + set -u: "${arr[@]}" is unbound when the array is empty
  if [[ ${#collected[@]} -eq 0 ]]; then
    eval "$var_name=()"
  else
    eval "$var_name=(\"\${collected[@]}\")"
  fi
}

exit_with_status() {
  if [[ "$failures" -gt 0 ]]; then
    echo ""
    log_fail "$failures check(s) failed"
    exit 1
  fi
  log_ok "all checks passed"
  exit 0
}
