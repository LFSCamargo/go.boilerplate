# Claude Code — go.boilerplate

Follow [`AGENTS.md`](AGENTS.md) for all instructions.

## Session start

1. Read `docs/DESIGN_DOC.md`
2. Read `AGENTS.md` architecture rules

## Project skills

- **create-db-entity** — `.cursor/skills/create-db-entity/SKILL.md` (also `.claude/skills/create-db-entity/`)
- **create-email-template** — Figma Protocol → `emails/emails/go_boilerplate` (figma-design-to-code; assets in `emails/static/go_boilerplate/`)

## Session end (required)

```bash
make ai-check
make test
make build
```

If you changed API, models, routes, or Docker — update `docs/DESIGN_DOC.md`.

## Quality gates

Never skip. Same scripts as Cursor/Codex:

- `scripts/ai-checks/run-all.sh full`
