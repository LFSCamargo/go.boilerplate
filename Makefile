.PHONY: test build ai-check ai-check-light ai-check-heavy ai-check-docs run \
	db-migrate db-migrate-down db-migrate-version db-new-entity \
	emails-install emails-dev emails-render \
	ai-n-plus-one ai-complexity ai-spaghetti ai-deadcode ai-security-light ai-security-heavy ai-test-driven \
	install-deps \
	up-dev up-dev-deps down-dev logs-dev run-dev \
	up-prod down-prod logs-prod

COMPOSE_DEV  = docker compose -f docker-compose.dev.yml
COMPOSE_PROD = docker compose -f docker-compose.prod.yml

install-deps:
	go mod tidy # install Go dependencies
	go install github.com/air-verse/air@v1.63.0
	go install golang.org/x/tools/cmd/deadcode@latest
	$(MAKE) emails-install # install email dependencies

test:
	go test ./... -count=1

build:
	go build ./...

run:
	go run ./src

# Host live reload (Air). Pair with `make up-dev-deps`.
run-dev:
	air -c .air.toml

# Full local stack with Air live reload (nginx :8080 → app :5000, Mailhog UI :8025).
up-dev:
	$(COMPOSE_DEV) up --build

# Postgres + Mailhog only — use with `make run` or `make run-dev` on the host.
up-dev-deps:
	$(COMPOSE_DEV) up -d postgres mailhog

down-dev:
	$(COMPOSE_DEV) down

logs-dev:
	$(COMPOSE_DEV) logs -f

# Production-like stack. Requires BEARER_SECRET, CORS_ORIGINS, SMTP_*, APP_PUBLIC_URL in .env.
up-prod:
	$(COMPOSE_PROD) up --build -d

down-prod:
	$(COMPOSE_PROD) down

logs-prod:
	$(COMPOSE_PROD) logs -f

db-migrate:
	go run ./cmd/migrate -direction=up

db-migrate-down:
	go run ./cmd/migrate -direction=down

db-migrate-version:
	go run ./cmd/migrate -direction=version

db-new-entity:
	@test -n "$(NAME)" || (echo "Usage: make db-new-entity NAME=EntityName FIELDS=\"field:type ...\"" && exit 1)
	bash scripts/db/new-entity.sh "$(NAME)" $(FIELDS)

emails-install:
	cd emails && npm install

emails-dev:
	# React Email serves `--dir`/static (emails/emails/static), not emails/static.
	ln -sfn ../static emails/emails/static
	cd emails && EMAIL_PREVIEW=1 npm run dev

emails-render:
	cd emails && npm run render -- $(TEMPLATE) '$(PROPS)'

ai-check:
	bash scripts/ai-checks/run-all.sh full

ai-check-light:
	bash scripts/ai-checks/run-all.sh light

ai-check-heavy:
	bash scripts/ai-checks/run-all.sh heavy

ai-check-docs:
	bash scripts/ai-checks/run-all.sh docs

ai-n-plus-one:
	bash scripts/ai-checks/n-plus-one.sh

ai-complexity:
	bash scripts/ai-checks/complexity.sh

ai-spaghetti:
	bash scripts/ai-checks/spaghetti.sh

ai-security-light:
	bash scripts/ai-checks/security-light.sh

ai-security-heavy:
	bash scripts/ai-checks/security-heavy.sh

ai-test-driven:
	bash scripts/ai-checks/test-driven.sh

ai-deadcode:
	bash scripts/ai-checks/deadcode.sh
