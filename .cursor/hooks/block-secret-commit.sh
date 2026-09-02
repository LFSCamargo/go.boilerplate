#!/usr/bin/env bash
# Cursor hook: warn before git commit if .env or secrets might be staged.
set -euo pipefail

input=$(cat)
command=$(echo "$input" | python3 -c "import sys,json; print(json.load(sys.stdin).get('command',''))" 2>/dev/null || echo "")

if [[ "$command" == *"git commit"* ]]; then
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    if git diff --cached --name-only 2>/dev/null | grep -qE '^\.env$|credentials|\.pem$|id_rsa'; then
      echo '{
        "permission": "ask",
        "user_message": "Possible secret file staged for commit (.env, credentials, keys). Review before committing.",
        "agent_message": "Hook blocked auto-commit: sensitive files may be staged. Run git diff --cached."
      }'
      exit 0
    fi
  fi
fi

echo '{ "permission": "allow" }'
exit 0
