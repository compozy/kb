# AGENTS.md

## Project Overview

`kb` is a single-binary Go CLI for building and maintaining topic-based knowledge bases in the Karpathy KB pattern. It handles the non-LLM workflow: topic scaffolding, multi-source ingestion, structural linting, codebase analysis, QMD indexing/search, and KB-oriented inspection commands. Semantic judgments (ingest gates, classification, linking, retrieval) go through a required decision model (Jev over OpenRouter); a generation model writes only short literals, never article prose.

**Reference TypeScript source:** `~/dev/projects/kodebase`

## Source of Truth

- KB pivot tech spec: `.compozy/tasks/kb-pivot/_techspec.md`
- KB pivot task tracker: `.compozy/tasks/kb-pivot/_tasks.md`
- KB pivot workflow memory: `.compozy/tasks/kb-pivot/memory/`
- Decision model spec (kb-jev, revision 3): `.wayfinder/kb-jev/spec.md`; cross-package contracts: `.claude/plans/kb-jev.md`; domain vocabulary: `CONTEXT.md`

## Critical Rules

- `make verify` is the non-negotiable completion gate: `fmt -> lint -> test -> build -> boundaries`.
- `make lint` must report zero findings.
- Use `go get` for dependency changes.
- Never use destructive git restore/reset/checkout/clean/rm commands without explicit approval.
- Prefer local repository inspection over web search for codebase questions.
- kb-written frontmatter keys carry no namespace prefix: they are plain (`summary`, `concepts`, `triage`, `genre`, ...) and ownership comes from `decisions.OwnedKeys` plus value hashes.

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
| `internal/firecrawl` | Firecrawl REST client for `kb ingest url`, including the fresh-refetch options (`maxAge`, `waitFor`, `onlyMainContent`) |
| `internal/mediadl` | Shared `yt-dlp` + STT engine reused by `youtube` and `instagram` |
| `internal/youtube` | YouTube URL parsing, channel/bulk ingestion, and frontmatter over the `mediadl` core |
| `internal/instagram` | Instagram reel/video URL parsing and caption+transcript composition over the `mediadl` core |
| `internal/frontmatter` | Shared frontmatter parsing and generation helpers |
| `internal/lint` | KB lint engine (structural and decision-workflow kinds; never calls a model) and report rendering |
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
| `internal/decisions` | Decision engine: `Engine.Decide`, transport to the OpenRouter decisions endpoint, receipt validation, receipts cache, banding, budget, concurrency, redaction, run summary, and the `OwnedKeys` list |
| `internal/generation` | Generation-model client for short literals over OpenRouter chat completions (strict JSON schema, reasoning disabled, fallback model), sharing the run budget and receipts |
| `internal/questions` | Embedded, versioned question banks, loader (guard prefix, schema check, bank hash), composite banks, topic extra banks, and regression cases |
| `internal/resolve` | Obsidian-style link resolution and link extraction shared by lint, refs, and link; quarantine path helpers |
| `internal/corpus` | Topic document loader, body and value hashes, BM25, alias dictionary and mention scan, section-aware excerpts, provenance state, `.decisions/state.jsonl`, and the byte-preserving owned-key writer |
| `internal/contract` | Selection contract parse/validate/hash/render, `CLAUDE.md` section import and render, `topic.yaml` settings (contract, drafts, decisions block), path globs |
| `internal/contract/topicyaml` | Node-level `topic.yaml` editor that keeps untouched keys, comments, and key order |
| `internal/session` | Per-run bundle for decision-backed commands: config, topic, settings, contract, engine, generation client, lazy corpus, state store, writer, review store, decision modes; fails fast without `OPENROUTER_API_KEY` |
| `internal/scope` | `kb topic contract` workflow: draft, `CLAUDE.md` import, self-check, impact preview, activation |
| `internal/classify` | Classification facets, generated literals (summary, entities, questions, criterion, aliases), the shared gate judgment, `recapture`/`remove` queues, concept proposals, vocabulary draft/accept |
| `internal/quality` | Code quality checks run before any quality question (`not_an_article`, `thin`, `error_page`) |
| `internal/gate` | Ingest gates: exact dedupe, pre-fetch relevance, code quality, refetch, post-fetch quality and relevance, near-duplicate, quarantine |
| `internal/link` | Link candidates, per-document judgment, write policy (frontmatter relations, body links), reverse pass |
| `internal/find` | Retrieval judge (candidates, judgment, ranking with named exclusions) and facet counts |
| `internal/review` | Review queue store, labels, `import-labels`, `import-links`, calibration with holdout, impact-preview ledger |
| `internal/refs` | Inbound-reference scan, quarantine and restore, `quarantine-ledger.jsonl` |
| `internal/actions` | `kb review accept\|reject` dispatcher over refs, link, gate, and classify |
| `internal/okf` | OKF bundle promote, conformance check, and decision-model type suggestion |
| `internal/fakes` | `httptest` fakes for the OpenRouter decisions and generation endpoints and Firecrawl, used by tests |

## CLI Commands

| Command | Purpose |
| --- | --- |
| `kb topic new <slug> <title> <domain>` | Scaffold a new knowledge base topic (`--mode wiki\|okf`) |
| `kb topic list` | List scaffolded topics |
| `kb topic info <slug>` | Show metadata for one topic |
| `kb topic contract <slug> --draft\|--import-claude\|--accept` | Draft, import, or accept (self-check + impact preview) the selection contract |
| `kb topic vocabulary <slug> --draft\|--accept` | Draft a concept vocabulary for a topic without articles; accept creates stub articles |
| `kb promote <wiki-doc> --to <okf-topic> [--type <type>]` | Promote a compiled wiki document into an OKF topic (type suggested when omitted) |
| `kb okf check <topic>` | Check an OKF topic for bundle conformance |
| `kb ingest url <url>... --topic <slug>` | Scrape one or more web URLs and ingest into a topic |
| `kb ingest file <path> --topic <slug>` | Convert a local file and ingest into a topic |
| `kb ingest youtube <url> --topic <slug>` | Extract a YouTube transcript and ingest into a topic |
| `kb ingest channel <url> --topic <slug>` | Bulk-extract a YouTube channel or playlist into a topic |
| `kb ingest instagram <url> --topic <slug>` | Extract an Instagram reel/video caption and transcript into a topic |
| `kb ingest codebase <path> --topic <slug>` | Analyze a codebase and ingest artifacts into a topic |
| `kb ingest bookmarks <path> --topic <slug>` | Ingest a bookmark-cluster markdown file into a topic |
| `kb classify <topic>` | Judge every document on the classification facets and write kb-owned frontmatter |
| `kb link <topic>` | Write typed relations and body links between documents and wiki articles |
| `kb find <topic> "<question>"` | Rank the documents that answer a question, with named exclusions |
| `kb find --facets <topic>` | Print facet counts without any model call |
| `kb review <topic>` | List pending review items (same as `kb review list <topic>`) |
| `kb review accept\|reject --topic <t> <id>...` | Act on or dismiss review items; every verdict becomes a label |
| `kb review import-links <topic>` | Import existing human links as positive link labels |
| `kb review import-labels <topic> --from <file>` | Import relevance labels from a screening or curation file |
| `kb review calibrate <topic>` | Measure precision/recall from labels and choose thresholds on a dev/holdout split |
| `kb lint [<slug>]` | Check one topic for structural and decision-workflow issues |
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
| `kb migrate transcripts --topic <slug>` | Move legacy `raw/transcripts/` files into `raw/youtube/` |

### Command Notes

- `topic` subcommands share the root `--vault` flag for vault path resolution; `topic new --mode okf` creates a root-level OKF bundle instead of the wiki/raw scaffold.
- `topic contract` takes exactly one of `--draft` (generation model), `--import-claude` (no model), or `--accept` (self-check, impact preview, then activation after typing the slug; `--force` passes self-check conflicts, `--yes` skips the prompt). `topic vocabulary` takes `--draft` or `--accept`.
- `promote` requires `--to <okf-topic>`; `--type` is optional when `[okf].types` is set (the decision model suggests it, applies it at P ≥ `okf_type`, otherwise queues an `okf-type` review item and asks for `--type`), and required otherwise. `--description` overrides the generated concept description.
- `okf check` accepts `--strict`, `--format table|json|tsv`, and `--decide` (ask the decision model for advisory `type_mismatch` / `description_unsupported` findings missing from receipts).
- `ingest` subcommands require `--topic <slug>` to identify the target topic.
- Every `ingest` subcommand except `codebase` runs the gates, then classifies and links new sources, and accepts `--batch <name>`, `--budget <usd>`, `--decisions shadow|apply`, and `--force` (skip dedupe and gates). Bulk subcommands (`channel`, URL lists) accept `--rescue <id>` for items skipped before fetch. `ingest url` accepts several URLs and `--from <file>`.
- `classify` accepts `--only-missing`, `--all`, `--budget`, `--decisions`, `--format`. `link` accepts `--all`, `--dry-run`, `--no-qmd`, `--budget`, `--decisions`. `find` accepts `--limit`, `--kind`, `--concept`, `--min-depth`, `--relevance`, `--explain`, `--json`, `--no-qmd`, `--budget`, and `--facets`.
- `review` / `review list` accept `--queue gate|skip|recapture|remove|link|contradiction|concept-proposal|okf-type` and `--format table|json`. `review accept|reject` take `--topic <t>` and ids, or bulk selectors `--queue`, `--all-purpose <purpose>`, `--above <p>`, `--topic-wide`. `review import-labels` takes `--from`, `--id-field`, `--match url|path|doi|pmcid`, `--decision-field`, `--keep-values`, `--decided-by-field`. `review calibrate` accepts `--write`.
- Decision-backed commands (`ingest` except `codebase`, `classify`, `link`, `find`, `review accept|reject`, `topic contract --draft|--accept`, `topic vocabulary`, `promote` with `[okf].types`, `okf check --decide`) fail fast without `OPENROUTER_API_KEY`. `version`, `topic new|list|info`, `topic contract --import-claude`, `review list|import-links|import-labels|calibrate`, `lint`, `inspect`, `index`, `search`, `ingest codebase`, and `migrate` never call a model.
- `ingest codebase` accepts `--include`, `--exclude`, `--semantic`, `--progress`, and `--log-format`.
- `ingest youtube` accepts `--transcribe captions|auto|stt`; `yt-dlp` is required for YouTube metadata, captions, and audio extraction.
- `inspect` subcommands share `--vault`, `--topic`, and `--format` (`table`, `json`, `tsv`).
- `lint` accepts `--format`, `--save`, and optional positional `<slug>` or `--topic`. It never calls a model; decision-workflow kinds (`frontmatter-dead-link`, `link-to-quarantined`, `needs-compile`, `key-conflict`, `contradiction`, `off-topic-kept`, `criterion-missing`, `unclassified`, `contract-missing`, `contract-draft-pending`, `vocabulary-missing`, `pending-review`, `recapture-pending`, `remove-pending`) come from frontmatter, `topic.yaml`, and `.decisions/` records.
- `search` supports `--lex`, `--vec`, `--limit`, `--min-score`, `--full`, `--all`, `--collection`, `--vault`, `--topic`, and `--format`.
- `index` supports `--vault`, `--topic`, `--name`, `--embed`, `--context`, and `--force-embed`.
- Hidden compatibility alias: `kb generate <path>` (delegates to codebase pipeline).

## Runtime Config

- `config.example.toml` documents the TOML keys currently supported by `internal/config`.
- `APP_CONFIG` overrides the config file path.
- `.env` is loaded automatically when present.
- `FIRECRAWL_API_KEY` and `FIRECRAWL_API_URL` configure the Firecrawl client for `ingest url`.
- `[okf].types` configures the optional OKF type vocabulary used by `promote` and `okf check`; `[okf.type_descriptions]` gives the option text per type for decision-model type suggestions.
- `OPENROUTER_API_KEY` (or `[openrouter].api_key`) is **required** by every decision-backed command; `OPENROUTER_API_URL` / `[openrouter].api_url` is shared by STT, decisions, and generation.
- `[decisions]` configures the decision model: `model` (default `typesafe/jev-1.13`, env `KB_DECISIONS_MODEL`), `deadline`, `retries`, `concurrency`, `max_state_bytes`, `budget_usd` (per-run ceiling shared with generation; `--budget` overrides), and `mode` (default body-link insertion mode, `shadow`).
- `[decisions.thresholds]` holds named thresholds (`link_apply`, `link_review`, `mention_sense`, `affects`, `relevance_quarantine`, `relevance_review`, `relevance_fetch`, `quality_apply`, `quality_review`, `duplicate`, `primary_concept`, `concept_noul`, `find_keep`, `find_window`, `okf_type`, `contract_conflict`); `topic.yaml` `decisions.thresholds` overrides them per topic.
- `[generation]` configures the generation model: `model` (default `xiaomi/mimo-v2.6-flash`, env `KB_GENERATION_MODEL`), `fallback_model`, `deadline`, `retries`, `summary_language`.
- `[firecrawl].refetch_wait_ms` (default 3000) and `[firecrawl].refetch_only_main_content` (default false) configure the fresh refetch the quality gate runs before judging a thin or broken capture.
- Per topic, `topic.yaml` carries `contract`, `contract_draft`, `vocabulary_draft`, and `decisions: {mode, gates, relevance, thresholds, exclude}`. Relevance gates stay in shadow until the topic has an accepted contract and either a calibration with ≥ 30 review relevance labels or `decisions.gates: apply`; quality gates and dedupe apply by default.
- `OPENAI_API_KEY`, `OPENAI_API_URL`, `STT_PROVIDER`, and `STT_MODEL` configure the default OpenAI STT provider for `ingest youtube --transcribe auto|stt`.
- `openrouter.stt_model` (TOML-only) selects the model of the optional OpenRouter STT provider when `stt.provider = "openrouter"`.
- Generation, inspect, search, and index behavior is configured by CLI flags rather than TOML keys.

## Architecture Notes

- `internal/generate` is the codebase pipeline orchestration layer. Keep Cobra commands thin and push behavior into internal packages.
- The codebase pipeline is: scan -> adapter parse -> graph normalize -> metrics compute -> vault render -> vault write -> inspect/search/index read paths.
- `kb ingest file` routes through a converter registry (`internal/convert`) that matches file extensions to format-specific converters (PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, text, images with OCR).
- `kb ingest url` uses `internal/firecrawl` for web scraping, then writes through `internal/ingest`.
- `kb ingest youtube` uses `internal/youtube` for `yt-dlp` caption/audio extraction, with OpenAI STT by default and OpenRouter as an optional provider.
- Raw KB documents must include frontmatter before being written; use `internal/frontmatter` helpers instead of hand-assembling YAML.
- `vault.RenderDocuments` returns markdown bodies that already include frontmatter. Base definitions are rendered separately and written as YAML `.base` files.
- `kb topic new` owns the topic skeleton under the vault root, including `raw/`, `wiki/`, `outputs/`, `bases/`, `CLAUDE.md`, `AGENTS.md`, `topic.yaml` (with an empty `contract:` block), and `log.md`. `topic.yaml` is the only source of truth for the selection contract; the `## Selection contract` section of the topic `CLAUDE.md` is rendered from it.
- Codebase inspection commands operate on `raw/codebase/` beneath the resolved topic. Keep inspect behavior topic-aware rather than vault-global.
- Decision-backed commands open an `internal/session` and ask every question through `session`/`decisions.Engine.Decide`; never call OpenRouter directly. Code builds candidates and applies thresholds; the decision model only judges text; the generation model only writes short literals, validated by code.
- kb-owned frontmatter keys (`decisions.OwnedKeys`) are written only through the `internal/corpus` writer, which keeps bytes it does not own, records value hashes in `.decisions/state.jsonl`, never overwrites a user-edited key, and skips files with `locked: true`. Ingest extras cannot set owned keys.
- Quarantine moves files to `raw/_quarantine/` through `internal/refs`, which records every index/frontmatter edit in the quarantine ledger so restore is byte-exact. Nothing is deleted.
- `kb lint` never calls a model; it reads frontmatter, `topic.yaml`, and `.decisions/` records.
- Shadow mode never moves, skips, quarantines, or edits bodies; would-be outcomes go to receipts and the run summary.

## Testing Notes

- Default to table-driven tests with focused helpers and `t.TempDir()` for filesystem isolation.
- Integration tests use the `integration` build tag and live next to the packages they exercise.
- CLI integration tests exercise real topic/ingest/lint/inspect flows instead of mocking Cobra wiring when the workflow itself is the behavior under test.
- QMD-related integration tests must isolate `HOME`, `XDG_CACHE_HOME`, and `XDG_CONFIG_HOME`.
- Decision-backed tests use `internal/fakes` (fake OpenRouter decisions and generation, fake Firecrawl) at the HTTP boundary; never hit the network in unit or integration tests.
- The live question-bank regression suite runs only with `go test -tags jevlive ./internal/questions/` and `OPENROUTER_API_KEY` set; receipts are cached under `$KB_JEVLIVE_CACHE` (or the user cache dir) and spend is capped at US$ 0.50 per run.
- Treat failing tests as product bugs first. Fix production behavior instead of weakening assertions.
