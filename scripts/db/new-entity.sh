#!/usr/bin/env bash
# Scaffold a new GORM entity, SQL migration pair, and run migrations locally.
# Usage: ./scripts/db/new-entity.sh EntityName [field:type ...]
# Example: ./scripts/db/new-entity.sh Friendship requester_id:uuid addressee_id:uuid status:string
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

ENTITY="${1:-}"
shift || true

if [[ -z "$ENTITY" ]]; then
  echo "Usage: $0 EntityName [field:type ...]" >&2
  echo "Example: $0 Friendship requester_id:uuid addressee_id:uuid status:string" >&2
  exit 1
fi

if [[ ! "$ENTITY" =~ ^[A-Z][A-Za-z0-9]*$ ]]; then
  echo "Entity name must be PascalCase (e.g. Friendship)" >&2
  exit 1
fi

table_name=$(python3 -c "import re,sys; print(re.sub(r'(?<!^)(?=[A-Z])', '_', sys.argv[1]).lower())" "$ENTITY")

migration_dir="src/db/migrations"
last_num=$(ls "$migration_dir"/*_*.up.sql 2>/dev/null | sed -E 's|.*/([0-9]{6})_.*|\1|' | sort | tail -1)
if [[ -z "$last_num" ]]; then
  next_num="000001"
else
  next_num=$(printf "%06d" $((10#$last_num + 1)))
fi

migration_slug="create_${table_name}"
up_file="${migration_dir}/${next_num}_${migration_slug}.up.sql"
down_file="${migration_dir}/${next_num}_${migration_slug}.down.sql"
model_file="src/db/models/${table_name}.go"

if [[ -f "$model_file" ]]; then
  echo "Model already exists: $model_file" >&2
  exit 1
fi

# Build optional columns from field:type args
extra_up=""
extra_down=""
extra_model=""
for field_spec in "$@"; do
  field="${field_spec%%:*}"
  ftype="${field_spec#*:}"
  col="$field"
  case "$ftype" in
    uuid) sql_type="UUID NOT NULL" go_type="uuid.UUID" ;;
    string) sql_type="VARCHAR(255) NOT NULL" go_type="string" ;;
    text) sql_type="TEXT NOT NULL" go_type="string" ;;
    bool) sql_type="BOOLEAN NOT NULL DEFAULT FALSE" go_type="bool" ;;
    int) sql_type="INT NOT NULL DEFAULT 0" go_type="int" ;;
    timestamptz) sql_type="TIMESTAMPTZ NOT NULL DEFAULT NOW()" go_type="time.Time" ;;
    *)
      echo "Unknown field type: $ftype (use uuid|string|text|bool|int|timestamptz)" >&2
      exit 1
      ;;
  esac
  extra_up="${extra_up}    ${col} ${sql_type},\n"
  extra_model="${extra_model}\t${field} ${go_type} \`gorm:\"column:${col}\"\`\n"
done

cat >"$up_file" <<EOF
CREATE TABLE IF NOT EXISTS ${table_name} (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
$(printf "%b" "$extra_up")    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
EOF

cat >"$down_file" <<EOF
DROP TABLE IF EXISTS ${table_name};
EOF

cat >"$model_file" <<EOF
package models

import (
	"time"

	"github.com/google/uuid"
)

type ${ENTITY} struct {
	ID uuid.UUID \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"\`
$(printf "%b" "$extra_model")	CreatedAt time.Time \`gorm:"not null;autoCreateTime"\`
	UpdatedAt time.Time \`gorm:"not null;autoUpdateTime"\`
}

func (${ENTITY}) TableName() string { return "${table_name}" }
EOF

# Append to registry.go before closing brace of All()
registry="src/db/models/registry.go"
if grep -q "&${ENTITY}{}" "$registry"; then
  echo "Registry already contains ${ENTITY}" >&2
else
  python3 - <<PY "$registry" "$ENTITY"
import pathlib, sys
path = pathlib.Path(sys.argv[1])
entity = sys.argv[2]
text = path.read_text()
needle = "\t}\n"
insert = f"\t\t&{entity}{{}},\n"
if insert.strip() in text:
    raise SystemExit(0)
text = text.replace("\t}\n", insert + "\t}\n", 1)
path.write_text(text)
PY
fi

echo "Created:"
echo "  $model_file"
echo "  $up_file"
echo "  $down_file"
echo "Updated: $registry"
echo ""
echo "Running migrations against local environment..."

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker not found — run 'docker compose up postgres -d' then 'make db-migrate'" >&2
else
  if ! docker compose ps postgres 2>/dev/null | grep -q "running"; then
    echo "Starting postgres via docker compose..."
    docker compose up postgres -d
    sleep 3
  fi
fi

make db-migrate
echo "Done. Update docs/DESIGN_DOC.md with the new entity."
