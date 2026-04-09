# CLAUDE.md

## Project Overview

Kodebase CLI is a single-binary Go rewrite of the TypeScript `kodebase` CLI. The tool generates Obsidian knowledge vaults from source code repositories, exposes inspection commands over the generated snapshot, and integrates with QMD for indexing and search.

**Reference TypeScript source:** `~/dev/projects/kodebase`

## Source of Truth

- Rewrite implementation spec: `.compozy/tasks/rewrite/_techspec.md`
- Rewrite task tracker: `.compozy/tasks/rewrite/_tasks.md`
- Workflow memory: `.compozy/tasks/rewrite/memory/`

## Critical Rules

- `make verify` is the non-negotiable completion gate: `fmt -> lint -> test -> build -> boundaries`.
- `make lint` must report zero findings.
- Use `go get` for dependency changes.
- Never use destructive git restore/reset/checkout/clean/rm commands without explicit approval.
- Prefer local repository inspection over web search for codebase questions.

## Active Build Surface

```bash
make verify
make fmt
make lint
make test
make test-integration
make build
make deps
make help
```

## Active Package Layout

| Path | Responsibility |
| --- | --- |
| `cmd/kodebase` | Program entrypoint |
| `internal/cli` | Cobra command assembly and command adapters |
| `internal/generate` | Repository-to-vault orchestration |
| `internal/models` | Shared domain models and interfaces |
| `internal/scanner` | Source discovery and ignore handling |
| `internal/adapter` | Tree-sitter parsing adapters |
| `internal/graph` | Graph normalization |
| `internal/metrics` | Metrics computation |
| `internal/vault` | Rendering, writing, reading, and vault query helpers |
| `internal/qmd` | QMD shell client |
| `internal/output` | Structured output formatting |
| `internal/config` | TOML config plus env-backed secrets |
| `internal/logger` | Slog logger setup |
| `internal/version` | Build version metadata |

## Implementation Conventions

- The active import layout is `internal/...`, not `internal/kodebase/...`. Treat any remaining `internal/kodebase/...` references as stale documentation.
- Keep Cobra commands as thin adapters. Business logic belongs in packages like `internal/generate`, `internal/vault`, and `internal/qmd`.
- `generate` writes a JSON summary to stdout and uses structured stage logs through `slog`.
- `inspect` and `search` share the `internal/output` formatter and should continue to support `table`, `json`, and `tsv`.
- `search` defaults to hybrid QMD retrieval. `--lex` and `--vec` are mutually exclusive mode selectors.
- `index` is intentionally idempotent at the CLI layer by checking `qmd status` first and choosing add vs update.
- Parse QMD JSON from stdout only. Treat stderr as progress and diagnostics, not as standalone failure evidence.
- `vault.RenderDocuments` returns markdown content whose `Body` already includes YAML frontmatter. `.base` files are generated separately and stored as YAML definitions.
- Read-side circular dependency reporting reconstructs ordered cycles from file-level import relations instead of relying on serialized cycle lists.
- Integration tests that exercise QMD must isolate `HOME`, `XDG_CACHE_HOME`, and `XDG_CONFIG_HOME`.

## CLI Surface

- `kodebase generate <path>`
- `kodebase inspect {smells|dead-code|complexity|blast-radius|coupling|symbol|file|backlinks|deps|circular-deps}`
- `kodebase search <query>`
- `kodebase index`
- `kodebase index-vault`
- `kodebase version`

## Runtime Config Notes

- `config.example.toml` should only document keys accepted by `internal/config`.
- `APP_CONFIG` selects the TOML file path.
- `.env` may supply `DATABASE_URL` and `API_KEY`.
- Generate/search/index tuning currently lives on CLI flags rather than config file keys.

## Testing Conventions

- Default to table-driven tests with focused helpers and `t.TempDir()` for filesystem isolation.
- Keep integration tests co-located with the package under test behind `//go:build integration`.
- Always treat a test failure as a behavior bug until proven otherwise; do not weaken tests to fit broken behavior.
