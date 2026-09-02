# Codex — go.boilerplate

Follow [`AGENTS.md`](../AGENTS.md) for all instructions.

## Required before completing work

```bash
make ai-check
make test
make build
```

## Design doc

Source of truth: [`docs/DESIGN_DOC.md`](../docs/DESIGN_DOC.md)

Update it when changing structure, API, or data models.

## Project skills

| Skill | When |
|-------|------|
| `create-db-entity` | New GORM model + SQL migration |
| `create-email-template` | New React Email from Figma Protocol assets |

Skills live in `.cursor/skills/`, `.claude/skills/`, and `.codex/skills/`. For Figma email work, also follow **figma-design-to-code**.

## Quality scripts

All checks live in `scripts/ai-checks/`. Run via `make ai-check` — do not reimplement checks in prompts.
