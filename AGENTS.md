# AGENTS.md

## Project Overview

Kodebase CLI is a Go port of [kodebase](https://github.com/compozy/kodebase). It scans source repositories, builds a normalized code graph, computes metrics, and writes Obsidian-compatible knowledge vaults that can later be inspected and indexed with QMD.

**Reference TypeScript source:** `~/dev/projects/kodebase`

## Source of Truth

- Rewrite PRD and implementation plan: `.compozy/tasks/rewrite/_techspec.md`
- Rewrite task tracker: `.compozy/tasks/rewrite/_tasks.md`
- Task-local memory and handoff notes: `.compozy/tasks/rewrite/memory/`

## Critical Rules

- `make verify` is the blocking completion gate. It must pass with zero warnings and zero errors.
- `make lint` has zero-tolerance for golangci-lint findings.
- Use `go get` for dependency changes. Do not hand-edit `go.mod`.
- Do not run destructive git restore/reset/checkout/clean/rm commands without explicit approval.
- Prefer local code search (`rg`, `rg --files`) over web lookups for repository questions.

## Build Commands

```bash
make verify              # fmt -> lint -> test -> build -> boundaries
make fmt                 # gofmt over repository Go files
make lint                # golangci-lint v2
make test                # unit tests with -race via gotestsum
make test-integration    # unit + integration tests with -race and -tags integration
make build               # build ./... and bin/kodebase with ldflags
make deps                # go mod tidy
make help                # mage target list
```

## Package Layout

| Path | Responsibility |
| --- | --- |
| `cmd/kodebase` | CLI entrypoint |
| `internal/cli` | Cobra command tree and command-specific I/O |
| `internal/generate` | End-to-end generate pipeline orchestration |
| `internal/models` | Domain types, snapshots, metrics, and shared interfaces |
| `internal/scanner` | Workspace discovery and ignore filtering |
| `internal/adapter` | Tree-sitter adapters for Go and TS/JS |
| `internal/graph` | Graph normalization from parsed files |
| `internal/metrics` | File, symbol, and directory metric computation |
| `internal/vault` | Path/text helpers, document rendering, writing, reading, and query resolution |
| `internal/qmd` | Shell-backed QMD integration for search and indexing |
| `internal/output` | Table / JSON / TSV output rendering |
| `internal/config` | TOML loading plus env-backed runtime secrets |
| `internal/logger` | Structured slog logger construction |
| `internal/version` | Build metadata surfaced by `kodebase version` |

## CLI Commands

| Command | Purpose |
| --- | --- |
| `kodebase generate <path>` | Scan a repository and write a topic vault summary as JSON |
| `kodebase inspect smells` | List smell signals for symbols and files |
| `kodebase inspect dead-code` | List dead exports and orphan files |
| `kodebase inspect complexity` | Rank functions by cyclomatic complexity |
| `kodebase inspect blast-radius` | Rank symbols by blast radius |
| `kodebase inspect coupling` | Rank files by instability |
| `kodebase inspect symbol <name>` | Find symbols by case-insensitive substring |
| `kodebase inspect file <path>` | Resolve one source file by exact path |
| `kodebase inspect backlinks <name-or-path>` | Show inbound relations for a file or symbol |
| `kodebase inspect deps <name-or-path>` | Show outgoing relations for a file or symbol |
| `kodebase inspect circular-deps` | List detected circular dependency cycles |
| `kodebase search <query>` | Query a generated vault through QMD hybrid, lexical, or vector modes |
| `kodebase index` | Create or update a QMD collection for a generated topic |
| `kodebase index-vault` | Alias for `kodebase index` |
| `kodebase version` | Print build version metadata |

### Command Notes

- `generate` accepts `--output`, `--topic`, `--title`, `--domain`, `--include`, `--exclude`, and `--semantic`.
- `inspect` subcommands share `--vault`, `--topic`, and `--format` (`table`, `json`, `tsv`).
- `search` supports `--lex`, `--vec`, `--limit`, `--min-score`, `--full`, `--all`, `--collection`, `--vault`, `--topic`, and `--format`.
- `index` supports `--vault`, `--topic`, `--name`, `--embed`, `--context`, and `--force-embed`.

## Runtime Config

- `config.example.toml` documents the TOML keys currently supported by `internal/config`.
- `APP_CONFIG` overrides the config file path.
- `.env` is loaded automatically when present.
- `DATABASE_URL` and `API_KEY` remain environment-only runtime secrets.
- Generation, inspect, search, and index behavior is currently configured by CLI flags rather than TOML keys.

## Architecture Notes

- Active package layout is the shorter `internal/...` structure. Ignore stale `internal/kodebase/...` references when implementing new work.
- `internal/generate` is the orchestration layer. Keep Cobra commands thin and push behavior into internal packages.
- The pipeline is: scan -> adapter parse -> graph normalize -> metrics compute -> vault render -> vault write -> inspect/search/index read paths.
- `vault.RenderDocuments` returns markdown bodies that already include frontmatter. Base definitions are rendered separately and written as YAML `.base` files.

## Testing Notes

- Integration tests use the `integration` build tag and live next to the packages they exercise.
- QMD-related integration tests must isolate `HOME`, `XDG_CACHE_HOME`, and `XDG_CONFIG_HOME`.
- Treat failing tests as product bugs first. Fix production behavior instead of weakening assertions.
