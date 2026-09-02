<p align="center">
  <img src="emails/static/go_boilerplate/hero_image.png" alt="Go Boilerplate" width="640" />
</p>

<h1 align="center">go.boilerplate</h1>

<p align="center">
  Production-ready <strong>Go HTTP API starter</strong> — clone it, rename the module, and ship your product.
</p>

<p align="center">
  <a href="https://github.com/LFSCamargo/go.boilerplate/generate">Use this template</a>
  ·
  <a href="docs/DESIGN_DOC.md">Design doc</a>
  ·
  <a href="AGENTS.md">AI agent guide</a>
</p>

---

## Purpose

**go.boilerplate** is a GitHub template repository for starting new Go backend services without rebuilding the same foundations every time.

It gives you a working baseline for:

- User auth (register, login, logout, JWT revocation)
- Email verification and password reset (OTP + magic links)
- PostgreSQL migrations and a modular codebase
- Docker Compose for local dev and production-like runs
- OpenAPI docs, structured logging, and AI-friendly quality gates

This repo intentionally **does not** ship a product domain (no chat, billing, etc.). You add those as new modules under `src/modules/<name>` after cloning.

Use **[Use this template](https://github.com/LFSCamargo/go.boilerplate/generate)** on GitHub, or:

```bash
gh repo create my-api --template LFSCamargo/go.boilerplate --public --clone
```

Then rename the Go module path (`go.boilerplate` → `your.module`) and update branding in emails/docs.

---

## What's included

| Area | What you get |
|------|----------------|
| **HTTP** | [Fiber v3](https://gofiber.io/) + [Huma](https://huma.rocks/) OpenAPI at `/docs` |
| **Auth** | JWT with `jti` revocation, bcrypt passwords, OTP flows, avatar upload on register |
| **Database** | PostgreSQL 16, GORM models, SQL migrations via golang-migrate |
| **Email** | React Email templates (Protocol style) + SMTP; [Mailhog](https://github.com/mailhog/MailHog) in dev |
| **Logging** | Pino-style structured logs on Go `slog` with redaction |
| **Infra** | Docker Compose (Air live reload in dev, nginx reverse proxy) |
| **Quality** | Tests, CI, and shared AI checks for Cursor, Claude Code, and Codex |
| **Docs** | Design doc, agent rules, and skills for new entities and email templates |

### Auth endpoints (v1)

Base path: `/api/v1`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Create account (+ optional avatar) |
| POST | `/auth/login` | Issue JWT |
| POST | `/auth/logout` | Revoke current token |
| POST | `/auth/recover-password` | Send reset email |
| POST | `/auth/reset-password` | Reset with OTP or magic link |
| POST | `/auth/verify-email` | Resend verification (Bearer) |
| GET | `/auth/verify-email?token=` | Confirm via magic link |
| POST | `/auth/verify-email-code` | Confirm with OTP (Bearer) |
| GET | `/auth/me` | Current user |
| GET | `/health` | Liveness |

Full contract: [`docs/DESIGN_DOC.md`](docs/DESIGN_DOC.md)

---

## Repository layout

```
go.boilerplate/
├── src/
│   ├── main.go                 # bootstrap only
│   ├── config/                 # env validation (Zog)
│   ├── middleware/             # auth, rate limit, body validation
│   ├── db/                     # GORM + migrations + models
│   ├── mail/                   # React Email renderer + SMTP
│   ├── log/                    # structured logging
│   ├── openapi/                # Huma config + Scalar playground
│   ├── routes/                 # wire modules (no business logic)
│   └── modules/
│       ├── auth/               # handlers → services → repositories
│       └── health/
├── emails/
│   ├── emails/go_boilerplate/  # React Email templates
│   └── static/go_boilerplate/  # template assets (hero, logo, icons)
├── scripts/
│   ├── ai-checks/              # quality gates (N+1, security, docs sync, …)
│   └── db/new-entity.sh        # scaffold GORM entity + migration
├── docs/DESIGN_DOC.md          # architecture & API source of truth
├── AGENTS.md                   # AI development guide
├── docker-compose.dev.yml      # dev stack (Air + Postgres + Mailhog + nginx)
├── docker-compose.prod.yml     # production-like stack
└── .github/workflows/ci.yml    # test, build, ai-check
```

**Architecture rule:** handlers stay thin; repositories own all DB access; product features live in new `src/modules/*` packages.

---

## Quick start

**Requirements:** Go 1.27+, Docker, Node.js (for email templates)

```bash
cp .env.example .env
make install-deps
make up-dev-deps          # Postgres :5432 + Mailhog :8025
make db-migrate
make run-dev              # Air live reload on :5000
```

Or run the full stack (nginx `:8080` → app `:5000`):

```bash
make up-dev
```

| URL | Service |
|-----|---------|
| http://localhost:8080/docs | OpenAPI playground (Scalar) |
| http://localhost:8080/health | Health check |
| http://localhost:8025 | Mailhog (captured emails) |

---

## Common commands

```bash
make test              # all Go tests
make build             # compile
make ai-check          # full quality gate (before merge)
make ai-check-light    # fast gate after edits
make emails-dev        # preview templates at :3001
make db-new-entity NAME=Widget FIELDS="name:string active:bool"
make down-dev          # stop Docker dev stack
```

---

## Email templates

Transactional emails use the **Protocol** style ([Figma SaaS Email Templates](https://www.figma.com/community/file/1626680546446620209/saas-email-templates)).

| Template | File |
|----------|------|
| Verify email | `emails/emails/go_boilerplate/verify_email.tsx` |
| Password reset | `emails/emails/go_boilerplate/password_reset.tsx` |

Preview locally:

```bash
make emails-install
make emails-dev
```

Render from CLI:

```bash
make emails-render TEMPLATE=verify_email PROPS='{"code":"123456","name":"Ada","companyName":"Go Boilerplate"}'
```

After cloning, update brand copy and assets under `emails/static/<your-brand>/`.

---

## Adding your product

1. **Use this template** and rename the Go module in `go.mod` and imports.
2. Create `src/modules/<domain>/` with `router.go`, `handlers/`, `services/`, `repositories/`.
3. Register the module in `src/routes/routes.go`.
4. Add tables with `make db-new-entity` (see `.cursor/skills/create-db-entity/`).
5. Update `docs/DESIGN_DOC.md` and run `make ai-check`.

Do not grow chat, friends, or other product domains inside this starter — keep it a reusable base.

---

## AI-assisted development

This repo is built for human + AI workflows:

- **Cursor** — `.cursor/rules/`, hooks, skills
- **Claude Code** — `CLAUDE.md` → `AGENTS.md`
- **Codex** — `.codex/instructions.md`
- **CI** — same checks as local via `scripts/ai-checks/`

Read [`AGENTS.md`](AGENTS.md) before making changes. The design doc is the contract for API, models, and rollout.

---

## Stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.27 |
| HTTP | Fiber v3 + Huma OpenAPI |
| Config | Zog + `zenv` |
| Database | PostgreSQL 16 + GORM + golang-migrate |
| Auth | JWT + revoked tokens + OTP |
| Email | React Email + SMTP |
| Infra | Docker Compose + nginx + Air (dev) |

---

## License

No license file is included yet. Add one when you fork or publish a derivative project.

---

<p align="center">
  <sub>Extracted from the <a href="https://github.com/LFSCamargo/go.chat">go.chat</a> stack — same patterns, no product domain.</sub>
</p>
