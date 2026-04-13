# AGENTS.md

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

## Build Commands

```bash
make verify              # fmt -> lint -> test -> build -> boundaries
make fmt                 # gofmt over repository Go files
make lint                # golangci-lint v2
make test                # unit tests with -race via gotestsum
make test-integration    # unit + integration tests with -race and -tags integration
make build               # build ./... and bin/kb with ldflags
make deps                # go mod tidy
make help                # mage target list
```

## Package Layout

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

## CLI Commands

| Command | Purpose |
| --- | --- |
| `kb topic new <slug> <title> <domain>` | Scaffold a new knowledge base topic |
| `kb topic list` | List scaffolded topics |
| `kb topic info <slug>` | Show metadata for one topic |
| `kb ingest url <url> --topic <slug>` | Scrape a web URL and ingest into a topic |
| `kb ingest file <path> --topic <slug>` | Convert a local file and ingest into a topic |
| `kb ingest youtube <url> --topic <slug>` | Extract a YouTube transcript and ingest into a topic |
| `kb ingest codebase <path> --topic <slug>` | Analyze a codebase and ingest artifacts into a topic |
| `kb ingest bookmarks <path> --topic <slug>` | Ingest a bookmark-cluster markdown file into a topic |
| `kb lint [<slug>]` | Check one topic for structural KB issues |
| `kb inspect smells` | List smell signals for symbols and files |
| `kb inspect dead-code` | List dead exports and orphan files |
| `kb inspect complexity` | Rank functions by cyclomatic complexity |
| `kb inspect blast-radius` | Rank symbols by blast radius |
| `kb inspect coupling` | Rank files by instability |
| `kb inspect symbol <name>` | Find symbols by case-insensitive substring |
| `kb inspect file <path>` | Resolve one source file by exact path |
| `kb inspect backlinks <name-or-path>` | Show inbound relations for a file or symbol |
| `kb inspect deps <name-or-path>` | Show outgoing relations for a file or symbol |
| `kb inspect circular-deps` | List detected circular dependency cycles |
| `kb search <query>` | Query a vault through QMD hybrid, lexical, or vector modes |
| `kb index` | Create or update a QMD collection for a topic |
| `kb version` | Print build version metadata |

### Command Notes

- `topic` subcommands share the root `--vault` flag for vault path resolution.
- `ingest` subcommands require `--topic <slug>` to identify the target topic.
- `ingest codebase` accepts `--include`, `--exclude`, `--semantic`, `--progress`, and `--log-format`.
- `ingest youtube` accepts `--stt` for OpenRouter speech-to-text fallback.
- `inspect` subcommands share `--vault`, `--topic`, and `--format` (`table`, `json`, `tsv`).
- `lint` accepts `--format`, `--save`, and optional positional `<slug>` or `--topic`.
- `search` supports `--lex`, `--vec`, `--limit`, `--min-score`, `--full`, `--all`, `--collection`, `--vault`, `--topic`, and `--format`.
- `index` supports `--vault`, `--topic`, `--name`, `--embed`, `--context`, and `--force-embed`.
- Hidden compatibility alias: `kb generate <path>` (delegates to codebase pipeline).

## Runtime Config

- `config.example.toml` documents the TOML keys currently supported by `internal/config`.
- `APP_CONFIG` overrides the config file path.
- `.env` is loaded automatically when present.
- `FIRECRAWL_API_KEY` and `FIRECRAWL_API_URL` configure the Firecrawl client for `ingest url`.
- `OPENROUTER_API_KEY`, `OPENROUTER_API_URL`, and `openrouter.stt_model` (TOML-only) configure the OpenRouter client for `ingest youtube --stt`.
- Generation, inspect, search, and index behavior is configured by CLI flags rather than TOML keys.

## Architecture Notes

- `internal/generate` is the codebase pipeline orchestration layer. Keep Cobra commands thin and push behavior into internal packages.
- The codebase pipeline is: scan -> adapter parse -> graph normalize -> metrics compute -> vault render -> vault write -> inspect/search/index read paths.
- `kb ingest file` routes through a converter registry (`internal/convert`) that matches file extensions to format-specific converters (PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, text, images with OCR).
- `kb ingest url` uses `internal/firecrawl` for web scraping, then writes through `internal/ingest`.
- `kb ingest youtube` uses `internal/youtube` for transcript extraction, with optional OpenRouter STT fallback.
- Raw KB documents must include frontmatter before being written; use `internal/frontmatter` helpers instead of hand-assembling YAML.
- `vault.RenderDocuments` returns markdown bodies that already include frontmatter. Base definitions are rendered separately and written as YAML `.base` files.
- `kb topic new` owns the topic skeleton under the vault root, including `raw/`, `wiki/`, `outputs/`, `bases/`, `CLAUDE.md`, `AGENTS.md`, and `log.md`.
- Codebase inspection commands operate on `raw/codebase/` beneath the resolved topic. Keep inspect behavior topic-aware rather than vault-global.

## Testing Notes

- Default to table-driven tests with focused helpers and `t.TempDir()` for filesystem isolation.
- Integration tests use the `integration` build tag and live next to the packages they exercise.
- CLI integration tests exercise real topic/ingest/lint/inspect flows instead of mocking Cobra wiring when the workflow itself is the behavior under test.
- QMD-related integration tests must isolate `HOME`, `XDG_CACHE_HOME`, and `XDG_CONFIG_HOME`.
- Treat failing tests as product bugs first. Fix production behavior instead of weakening assertions.
