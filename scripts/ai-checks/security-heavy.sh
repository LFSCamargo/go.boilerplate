#!/usr/bin/env bash
# Heavy security — semgrep + gitleaks (pre-merge / stop hook).
source "$(dirname "$0")/lib/common.sh"

echo "== Security checker (heavy: semgrep + gitleaks) =="

if require_cmd semgrep; then
  if ! semgrep --config semgrep.yml --error --quiet . 2>/dev/null; then
    log_fail "semgrep found security issues — see output above"
  else
    log_ok "semgrep passed"
  fi
else
  log_warn "install semgrep: pip install semgrep  OR  brew install semgrep"
fi

if require_cmd gitleaks; then
  if ! gitleaks detect --source . --config .gitleaks.toml --no-banner --redact 2>/dev/null; then
    log_fail "gitleaks detected secrets in repository"
  else
    log_ok "gitleaks passed"
  fi
else
  log_warn "install gitleaks: brew install gitleaks"
fi

exit_with_status
