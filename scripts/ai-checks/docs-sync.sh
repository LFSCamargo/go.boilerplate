#!/usr/bin/env bash
# Docs synchronization — structural changes must update DESIGN_DOC.md.
source "$(dirname "$0")/lib/common.sh"

DESIGN_DOC="docs/DESIGN_DOC.md"
TRIGGER_PATTERNS='^src/(modules|routes|db)/|^src/main\.go$|docker-compose\.(dev|prod)\.yml$|Dockerfile$|^nginx/'

echo "== Docs synchronization checker =="

if [[ ! -f "$DESIGN_DOC" ]]; then
  log_fail "$DESIGN_DOC missing"
  exit_with_status
fi

collect_to_array changed changed_files
if [[ ${#changed[@]} -eq 0 ]]; then
  log_ok "no changed files"
  exit 0
fi

needs_doc=false
doc_changed=false

for f in "${changed[@]}"; do
  normalized="${f#./}"
  [[ "$normalized" == "$DESIGN_DOC" ]] && doc_changed=true
  if echo "$normalized" | grep -qE "$TRIGGER_PATTERNS"; then
    needs_doc=true
  fi
done

if [[ "$needs_doc" == true && "$doc_changed" == false ]]; then
  log_fail "structural changes detected but $DESIGN_DOC was not updated"
  echo "       Update sections: API Endpoints, Data definition, System Design, or Rollout Plan"
  echo "       AI agents: read docs/DESIGN_DOC.md and sync before finishing the task"
else
  log_ok "docs in sync (or no structural changes)"
fi

exit_with_status
