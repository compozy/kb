# CLAUDE.md

## Project Overview

`kb` is a single-binary Go CLI for building and maintaining topic-based knowledge bases in the Karpathy KB pattern. It handles the non-LLM workflow: topic scaffolding, multi-source ingestion, structural linting, codebase analysis, QMD indexing/search, and KB-oriented inspection commands.

**Reference TypeScript source:** `~/dev/projects/kodebase`

## Source of Truth

- KB pivot tech spec: `.compozy/tasks/kb-pivot/_techspec.md`
- KB pivot task tracker: `.compozy/tasks/kb-pivot/_tasks.md`
- KB pivot workflow memory: `.compozy/tasks/kb-pivot/memory/`

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

- `make build` compiles `./...` and writes `bin/kb` from `cmd/kb`.

## Active Package Layout

| Path | Responsibility |
| --- | --- |
| `cmd/kb` | Program entrypoint for the `kb` binary |
| `internal/cli` | Cobra root, subcommands, flag resolution, and command I/O |
| `internal/topic` | Topic scaffolding, listing, and topic metadata lookup |
| `internal/ingest` | Ingest orchestration, frontmatter assembly, raw writes, and log entries |
| `internal/convert` | Converter registry and format-specific file converters |
| `internal/firecrawl` | Firecrawl REST client for `kb ingest url` |
| `internal/youtube` | YouTube transcript extraction and OpenRouter STT fallback |
| `internal/frontmatter` | Shared frontmatter parsing and generation helpers |
| `internal/lint` | KB structural lint engine and report rendering |
| `internal/generate` | Codebase-to-KB pipeline used by `kb ingest codebase` and the hidden legacy `generate` alias |
| `internal/scanner` | Source discovery and ignore handling for codebase ingest |
| `internal/adapter` | Tree-sitter parsing adapters |
| `internal/graph` | Graph normalization for codebase snapshots |
| `internal/metrics` | File, symbol, directory, and smell metrics |
| `internal/vault` | Topic path helpers, render/write logic, vault reads, and inspect snapshot loading |
| `internal/qmd` | QMD subprocess integration for index and search |
| `internal/output` | Table, JSON, and TSV formatting |
| `internal/models` | Shared domain models and interfaces |
| `internal/config` | TOML config plus env-backed runtime overrides |
| `internal/logger` | Slog logger setup |
| `internal/version` | Build metadata surfaced by `kb version` |

## Implementation Conventions

- Keep Cobra commands thin. Business logic belongs in packages like `internal/topic`, `internal/ingest`, `internal/lint`, `internal/generate`, `internal/vault`, and `internal/qmd`.
- `kb topic new` owns the topic skeleton under the selected vault root, including `raw/`, `wiki/`, `outputs/`, `bases/`, `CLAUDE.md`, `AGENTS.md`, and `log.md`.
- `kb ingest file`, `url`, `youtube`, and `bookmarks` should route through `internal/ingest` and the converter/client surfaces instead of duplicating write logic in CLI code.
- `kb ingest codebase` is the supported codebase entrypoint. The hidden `kb generate` command remains as a compatibility wrapper and should stay thin.
- Codebase inspection commands operate on `raw/codebase/` beneath the resolved topic. Keep inspect behavior topic-aware rather than vault-global.
- Raw KB documents must include frontmatter before being written; use `internal/frontmatter` helpers instead of hand-assembling YAML.
- CLI commands share the root `--vault` flag from `internal/cli/root.go`. Reuse the vault helpers instead of defining duplicate per-command vault flags.
- Prefer native Go integrations for converters and clients. Optional OCR or remote-provider fallbacks must degrade cleanly when prerequisites are missing.
- QMD integrations parse JSON from stdout and treat stderr as progress or diagnostics.

## CLI Surface

- `kb topic {new,list,info}`
- `kb ingest {url,file,youtube,codebase,bookmarks}`
- `kb lint [<slug>]`
- `kb inspect {smells|dead-code|complexity|blast-radius|coupling|symbol|file|backlinks|deps|circular-deps}`
- `kb index`
- `kb search <query>`
- `kb version`
- Hidden compatibility alias: `kb generate <path>`

## Runtime Config Notes

- `config.example.toml` documents every TOML section currently accepted by `internal/config`: `[app]`, `[log]`, `[firecrawl]`, and `[openrouter]`.
- `APP_CONFIG` selects the TOML file path.
- `.env` may supply `FIRECRAWL_API_KEY`, `FIRECRAWL_API_URL`, `OPENROUTER_API_KEY`, and `OPENROUTER_API_URL`.
- `openrouter.stt_model` is currently TOML-backed rather than env-overridden.

## Testing Conventions

- Default to table-driven tests with focused helpers and `t.TempDir()` for filesystem isolation.
- Keep integration tests co-located with the package under test behind `//go:build integration`.
- CLI integration tests should exercise real topic/ingest/lint/inspect flows instead of mocking Cobra wiring when the workflow itself is the behavior under test.
- QMD-related integration tests must isolate `HOME`, `XDG_CACHE_HOME`, and `XDG_CONFIG_HOME`.
- Treat test failures as behavior bugs first; do not weaken assertions to fit broken behavior.
