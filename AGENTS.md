# AGENTS.md — AI Development Guide (Cursor, Claude, Codex)

This file is the **single source of truth** for AI agents working on **go.boilerplate**.

## Project purpose

Production-ready **Go HTTP API starter**. Ships auth (JWT, OTP, email verify / password reset), health, React Email + SMTP, Docker Compose, Huma OpenAPI, and AI quality gates. Add product features as new `src/modules/<name>` packages.

**Always read first:** [`docs/DESIGN_DOC.md`](docs/DESIGN_DOC.md)

## Stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.27 |
| HTTP | Fiber v3 + Huma OpenAPI (`/docs`, `/openapi.json`) |
| Config | Zog + `zenv` |
| DB | PostgreSQL + GORM + golang-migrate |
| Auth | JWT + revoked token table + OTP |
| Email | SMTP (Mailhog in dev) + React Email templates |
| Infra | Docker Compose (`make up-dev` Air live reload / `make up-prod`) |

## Architecture rules

```
src/
├── main.go              # bootstrap only
├── config/              # env validation
├── middleware/          # shared Fiber middleware (RequireAuth, ValidateBody)
├── db/                  # GORM connection + models + migrations
│   ├── models/          # entities
│   └── migrations/      # SQL up/down
├── openapi/             # Huma config, docs, playground
├── routes/              # wire modules (no business logic)
└── modules/<name>/
    ├── router.go        # routes only
    ├── handlers/        # HTTP handlers (thin)
    ├── services/        # business logic
    └── repositories/    # GORM queries (Preload to avoid N+1)
```

### Non-negotiables

1. **Handlers** — parse input, call service, return response. No GORM in handlers.
2. **Repositories** — all DB access; use `Preload`/`Joins` for related data.
3. **Tests** — every new handler/service/repository gets tests. Bugs → regression test first.
4. **Design doc** — update `docs/DESIGN_DOC.md` when changing API, models, modules, or infra.
5. **No spaghetti** — max ~300 lines/file; split when mixing HTTP + DB + routing.

## Before finishing any task

Run quality gates (same on all CLIs):

```bash
make ai-check-light   # after small edits
make ai-check         # before marking task complete
```

Individual checks:

```bash
make ai-n-plus-one
make ai-complexity
make ai-deadcode
make ai-spaghetti
make ai-security-light
make ai-security-heavy
make ai-test-driven
make ai-check-docs
make test
make build
```

If a gate fails, **fix before finishing**. Do not disable checks.

## Feature workflow (test-driven)

1. Read `docs/DESIGN_DOC.md` — confirm scope fits v1 constraints.
2. Write failing test(s) for the behavior.
3. Implement minimal code to pass.
4. Run `make ai-check` + `make test`.
5. Update `docs/DESIGN_DOC.md` (API, data model, rollout phase).
6. Summarize changes and test commands for the user.

## Security

- Never commit secrets; use `.env` (gitignored) and `.env.example`.
- Heavy scan: `make ai-security-heavy` (semgrep + gitleaks).
- Install optional tools: `semgrep`, `gitleaks`, `gosec`, `gocyclo`, `golangci-lint`.

## Module checklist (new feature)

- [ ] Router in `src/modules/<name>/router.go`
- [ ] Handlers in `handlers/`
- [ ] Service + repository if DB involved
- [ ] Tests (`*_test.go` or `router_test.go`)
- [ ] Registered in `src/routes/routes.go`
- [ ] `docs/DESIGN_DOC.md` updated
- [ ] `make ai-check` passes

## New database entity

Use skill **create-db-entity** (`.cursor/skills/create-db-entity/SKILL.md`) or:

```bash
make db-new-entity NAME=EntityName FIELDS="field:type ..."
make db-migrate
```

## Email (React Email + SMTP)

```bash
make emails-install
make emails-dev    # preview templates at localhost:3001
```

Templates live in `emails/emails/go_boilerplate/` (**Go Boilerplate** emails, Protocol style from [Figma SaaS Email Templates](https://www.figma.com/community/file/1626680546446620209/saas-email-templates)). Static assets: `emails/static/go_boilerplate/`. Brand name in copy: **Go Boilerplate**.

**New database entity:** use skill **create-db-entity** (`.cursor/skills/create-db-entity/SKILL.md`, also `.claude/skills/` and `.codex/skills/`).

**New or updated templates:** use skill **create-email-template** — always pull assets from Figma via `get_design_context` + `download_assets` (figma-design-to-code); never commit expiring MCP URLs or hand-drawn icons.

Go sends via `src/mail`:

```go
mailer, _ := mail.NewFromConfig(cfg)
mailer.SendVerifyEmail(ctx, mail.VerifyEmailParams{To: email, Name: name, Code: otp})
```

## Cross-tool parity

| Tool | Config file |
|------|-------------|
| **Cursor** | `.cursor/rules/project-purpose.mdc`, `.cursor/hooks.json` |
| **Claude Code** | `CLAUDE.md` → this file |
| **Codex** | `.codex/instructions.md` → this file |
| **CI / local** | `Makefile`, `scripts/ai-checks/` |

All tools must run the **same scripts** in `scripts/ai-checks/` — never duplicate check logic in prompts alone.
