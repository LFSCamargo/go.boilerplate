---
name: create-db-entity
description: >-
  Scaffolds a new GORM entity in src/db with SQL migration up/down files and
  runs migrations against the local Postgres environment. Use when the user
  asks to add a database entity, table, model, migration, or describes a new
  domain object to persist.
---

# Create DB Entity

Add a new persisted entity to **go.boilerplate** following the `src/db` domain pattern.

## Prerequisites

- Postgres running locally: `make up-dev-deps`
- Env configured: `cp .env.example .env`
- Read `docs/DESIGN_DOC.md` data model section before adding entities

## Quick path (preferred)

Use the scaffold script, then refine generated files:

```bash
# fields optional: name:type (uuid|string|text|bool|int|timestamptz)
make db-new-entity NAME=Friendship FIELDS="requester_id:uuid addressee_id:uuid status:string"
# or
bash scripts/db/new-entity.sh Friendship requester_id:uuid addressee_id:uuid status:string
```

The script will:

1. Create `src/db/models/<snake_case>.go`
2. Create `src/db/migrations/NNNNNN_create_<table>.up.sql` and `.down.sql`
3. Register the model in `src/db/models/registry.go`
4. Run `make db-migrate` against local Postgres

## Manual path (complex entities)

When the prompt describes relationships, enums, or indexes the script cannot infer:

### 1. GORM model

Location: `src/db/models/<snake_case>.go`

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

type Example struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    // fields...
    CreatedAt time.Time `gorm:"not null;autoCreateTime"`
    UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (Example) TableName() string { return "examples" }
```

### 2. SQL migration pair

Location: `src/db/migrations/`

- Next sequential number: `000005`, `000006`, …
- Files: `NNNNNN_create_<table>.up.sql` and `.down.sql`
- Match GORM model exactly; prefer explicit indexes and FK constraints
- Down migration must fully reverse up

### 3. Register model

Add `&Example{}` to `models.All()` in `src/db/models/registry.go`.

### 4. Run migration locally

```bash
make db-migrate
make db-migrate-version   # verify
```

### 5. Tests

Add `src/db/models/<snake_case>_test.go` or extend `models_test.go` with `TableName` test.

### 6. Docs

Update `docs/DESIGN_DOC.md`:

- ER diagram / data definition
- API section if entity is user-facing

### 7. Quality gates

```bash
make test
make ai-check-light
```

## Existing entities (reference)

| Model | Table | Purpose |
|-------|-------|---------|
| `User` | `users` | Registered accounts |
| `UserRegistrationConfig` | `user_registration_configs` | Registration policy singleton |
| `RevokedToken` | `revoked_tokens` | JWT invalidation on logout |
| `OTP` | `otps` | Email verify & password reset codes |

## Conventions

- UUID primary keys (`gen_random_uuid()`)
- `created_at` / `updated_at` on all tables
- Repository layer owns queries (handlers never import GORM directly)
- Use `Preload`/`Joins` — run `make ai-n-plus-one` after repository work
- OTP codes stored hashed only; purposes: `email_verify`, `password_reset`

## Prompt examples

**User:** "Add a Friendship entity with requester, addressee, and status"

**Agent:**

1. Run `make db-new-entity NAME=Friendship FIELDS="requester_id:uuid addressee_id:uuid status:string"`
2. Edit migration to add FK references to `users(id)` and unique constraint
3. Edit model with `User` relations if needed
4. `make db-migrate`
5. Update `docs/DESIGN_DOC.md`
