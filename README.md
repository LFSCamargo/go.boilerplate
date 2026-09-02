# go.boilerplate

Production-ready **Go HTTP API starter**: Fiber v3, Huma OpenAPI, JWT auth, React Email + SMTP, PostgreSQL, Docker Compose, and AI quality gates.

**Design doc (source of truth):** [`docs/DESIGN_DOC.md`](docs/DESIGN_DOC.md)  
**Agent guide:** [`AGENTS.md`](AGENTS.md)

## Quick start

```bash
cp .env.example .env
make install-deps
make up-dev-deps          # Postgres :5432 + Mailhog :8025
make db-migrate
make run-dev              # Air live reload → :5000
```

Full stack (nginx `:8080` → app `:5000` + Mailhog):

```bash
make up-dev
```

- API / Scalar docs: http://localhost:8080/docs
- Mailhog UI: http://localhost:8025

## Stack

| Layer | Choice |
|-------|--------|
| HTTP | Fiber v3 + Huma (`/docs`, `/openapi.json`) |
| Auth | JWT (`jti` revocation) + OTP email verify / password reset |
| DB | PostgreSQL 16 + GORM + golang-migrate |
| Email | React Email templates + SMTP (Mailhog in dev) |
| Infra | `docker-compose.dev.yml` / `docker-compose.prod.yml` |

## Commands

```bash
make test
make build
make ai-check-light
make ai-check
make emails-dev           # preview templates at :3001
```

## Adding a product module

1. Create `src/modules/<name>/` (router, handlers, services, repositories).
2. Register it in `src/routes/routes.go`.
3. Add entities with `make db-new-entity NAME=…` (skill `create-db-entity`).
4. Update `docs/DESIGN_DOC.md` and run `make ai-check`.
