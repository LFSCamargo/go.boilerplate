#!/usr/bin/env bash
# Test-driven gate — changed production Go code must have tests.
source "$(dirname "$0")/lib/common.sh"

echo "== Test-driven checker =="

collect_to_array files changed_go_files
if [[ ${#files[@]} -eq 0 ]]; then
  log_ok "no changed Go files"
  exit 0
fi

for file in "${files[@]}"; do
  [[ -f "$file" ]] || continue

  # Skip main, generated, db migrations bootstrap
  [[ "$file" == */main.go ]] && continue
  [[ "$file" == */db/db.go ]] && continue

  dir=$(dirname "$file")
  base=$(basename "$file" .go)

  # Skip temporary stub handlers
  [[ "$base" == "stubs" ]] && continue

  has_test=false
  if [[ -f "${dir}/${base}_test.go" ]]; then
    has_test=true
  fi
  if [[ -f "${dir}/router_test.go" ]] && [[ "$file" == *router.go ]]; then
    has_test=true
  fi
  # handlers covered by module router_test.go
  if [[ "$file" == *handlers/* ]]; then
    module_dir=$(dirname "$dir")
    if [[ -f "${module_dir}/router_test.go" ]]; then
      has_test=true
    fi
  fi

  if [[ "$has_test" == false ]]; then
    log_fail "$file: no test file — add ${base}_test.go or router_test.go"
  fi
done

# Run tests for changed packages
if require_cmd go; then
  if ! go test ./... -count=1 2>&1; then
    log_fail "go test ./... failed"
  else
    log_ok "go test ./... passed"
  fi
fi

exit_with_status
