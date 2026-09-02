# AI Development Setup

Unified quality gates for **Cursor**, **Claude Code**, and **Codex**.

## Quick start

```bash
make ai-check        # full gate (before merge / task done)
make ai-check-light  # fast gate (after edits)
```

## Configuration map

| Purpose | Location |
|---------|----------|
| Agent instructions | `AGENTS.md` |
| Claude | `CLAUDE.md` |
| Codex | `.codex/instructions.md` |
| Cursor rules | `.cursor/rules/*.mdc` |
| Cursor hooks | `.cursor/hooks.json` |
| Skills (all agents) | `.cursor/skills/`, `.claude/skills/`, `.codex/skills/` (`create-db-entity`, `create-email-template`) |
| Shared check scripts | `scripts/ai-checks/` |
| Design doc (must stay synced) | `docs/DESIGN_DOC.md` |

## Checks

| Check | Script | When |
|-------|--------|------|
| N+1 SQL | `n-plus-one.sh` | light + full |
| Complexity | `complexity.sh` | full |
| Spaghetti | `spaghetti.sh` | light + full |
| Security (light) | `security-light.sh` | light + full |
| Security (heavy) | `security-heavy.sh` | full |
| Test-driven | `test-driven.sh` | light + full |
| Docs sync | `docs-sync.sh` | full |

## Optional tools (install for full coverage)

```bash
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
brew install semgrep gitleaks golangci-lint   # or pip install semgrep
```

## Cursor hooks

Hooks call the same scripts as `make ai-check`:

- **afterFileEdit** → light checks
- **stop** → full checks (may inject follow-up if failed)

Restart Cursor after editing `.cursor/hooks.json`.
