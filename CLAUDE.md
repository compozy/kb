# CLAUDE.md

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
| `internal/firecrawl` | Firecrawl REST client for `kb ingest url`, including the fresh-refetch options (`maxAge`, `waitFor`, `onlyMainContent`) |
| `internal/mediadl` | Shared `yt-dlp` + STT engine: audio/caption download, ffmpeg segmentation, and transcript extraction reused by `youtube` and `instagram` |
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

## Implementation Conventions

- Keep Cobra commands thin. Business logic belongs in packages like `internal/topic`, `internal/ingest`, `internal/lint`, `internal/generate`, `internal/vault`, and `internal/qmd`.
- `kb topic new` owns the topic skeleton under the selected vault root, including `raw/`, `wiki/`, `outputs/`, `bases/`, `CLAUDE.md`, `AGENTS.md`, `topic.yaml` (with an empty `contract:` block), and `log.md`.
- `kb ingest file`, `url`, `youtube`, `instagram`, and `bookmarks` should route through `internal/ingest` and the converter/client surfaces instead of duplicating write logic in CLI code.
- `kb ingest youtube` and `kb ingest instagram` are thin shims over `internal/mediadl`: each parses its platform URL and frontmatter, then delegates yt-dlp/STT extraction to the shared core. Add new yt-dlp-backed sources as `mediadl` consumers rather than duplicating the engine.
- `kb ingest codebase` is the supported codebase entrypoint. The hidden `kb generate` command remains as a compatibility wrapper and should stay thin.
- Codebase inspection commands operate on `raw/codebase/` beneath the resolved topic. Keep inspect behavior topic-aware rather than vault-global.
- Raw KB documents must include frontmatter before being written; use `internal/frontmatter` helpers instead of hand-assembling YAML.
- CLI commands share the root `--vault` flag from `internal/cli/root.go`. Reuse the vault helpers instead of defining duplicate per-command vault flags.
- Prefer native Go integrations for converters and clients. Optional OCR or remote-provider fallbacks must degrade cleanly when prerequisites are missing.
- QMD integrations parse JSON from stdout and treat stderr as progress or diagnostics.
- Decision-backed commands open an `internal/session` (which fails fast without `OPENROUTER_API_KEY`) and ask every question through `decisions.Engine.Decide`; never call OpenRouter directly. Code builds candidates, computes numbers, and applies thresholds; the decision model only judges text; the generation model only writes short literals, validated by code.
- kb-owned frontmatter keys (`decisions.OwnedKeys`) are written only through the `internal/corpus` writer: byte-preserving, value-hash ownership recorded in `.decisions/state.jsonl`, user-edited keys never overwritten (`skipped:user-key`), `locked: true` files skipped. Ingest extras cannot set owned keys.
- Quarantine and restore go through `internal/refs` so every index/frontmatter edit is in the quarantine ledger and restore is byte-exact. Nothing is deleted.
- `kb lint` never calls a model; it reads frontmatter, `topic.yaml`, and `.decisions/` records.
- Shadow mode never moves, skips, quarantines, or edits bodies; would-be outcomes go to receipts and the run summary. Relevance gates and body-link insertion default to shadow; quality gates and exact dedupe apply.
- `topic.yaml` is the only source of truth for the selection contract; the `## Selection contract` section of the topic `CLAUDE.md` is rendered from it on `kb topic new` and on every accept.

## CLI Surface

- `kb topic {new,list,info}` (`topic new --mode wiki|okf`)
- `kb topic contract <topic> --draft|--import-claude|--accept [--force] [--yes]`
- `kb topic vocabulary <topic> --draft|--accept`
- `kb promote <wiki-doc> --to <okf-topic> [--type <type>]` (type suggested by the decision model when omitted and `[okf].types` is set)
- `kb okf check <topic> [--strict] [--decide]`
- `kb ingest {url,file,youtube,channel,instagram,codebase,bookmarks}`; all but `codebase` run the gates and accept `--batch`, `--budget`, `--decisions shadow|apply`, `--force`; bulk (`channel`, URL lists) accepts `--rescue <id>`; `url` takes several URLs and `--from <file>`
- `kb classify <topic> [--only-missing] [--all]`
- `kb link <topic> [--all] [--dry-run] [--no-qmd]`
- `kb find <topic> "<question>" [--limit] [--kind] [--concept] [--min-depth] [--relevance] [--explain] [--json]`; `kb find --facets <topic>`
- `kb review <topic>` (= `review list`, `--queue`), `kb review accept|reject --topic <t> <id>... [--queue] [--all-purpose] [--above] [--topic-wide]`, `kb review import-links <topic>`, `kb review import-labels <topic> --from ...`, `kb review calibrate <topic> [--write]`
- `kb lint [<slug>]`
- `kb inspect {smells|dead-code|complexity|blast-radius|coupling|symbol|file|backlinks|deps|circular-deps}`
- `kb index`
- `kb search <query>`
- `kb version`
- `kb migrate transcripts`
- Hidden compatibility alias: `kb generate <path>`
- Decision-backed (need `OPENROUTER_API_KEY`): `ingest` except `codebase`, `classify`, `link`, `find`, `review accept|reject`, `topic contract --draft|--accept`, `topic vocabulary`, `promote` with `[okf].types`, `okf check --decide`. Everything else never calls a model.

## Runtime Config Notes

- `config.example.toml` documents every TOML section currently accepted by `internal/config`: `[app]`, `[log]`, `[vault]`, `[firecrawl]`, `[openrouter]`, `[stt]`, `[youtube]`, `[instagram]`, `[okf]`, `[decisions]` (plus `[decisions.thresholds]`), and `[generation]`.
- `APP_CONFIG` selects the TOML file path.
- `kb ingest youtube` requires `yt-dlp` for metadata, captions, and audio extraction. `[youtube].yt_dlp_path`, `[youtube].proxy`, `[youtube].cookies_file`, `[youtube].user_agent`, `[youtube].retry_attempts`, `[youtube].retry_backoff`, and `[youtube].transcription` configure that path.
- YouTube transcription policy is `captions|auto|stt`: `captions` uses YouTube captions only, `auto` uses manual captions when present and STT when only automatic captions or no captions are available, and `stt` forces audio transcription.
- `kb ingest instagram <url>` ingests reels, video posts, and IGTV through the same `yt-dlp` + STT engine. `[instagram]` mirrors the `[youtube]` yt-dlp knobs (`yt_dlp_path`, `proxy`, `cookies_file`, `user_agent`, `retry_attempts`, `retry_backoff`, `transcription`) with cookies kept separate because Instagram needs its own session. The default `[instagram].transcription` is `auto` since reels rarely carry caption tracks; the document body is the post caption followed by the transcript, degrading to caption-only when audio yields no transcript. Whole-profile bulk ingestion is out of scope (yt-dlp's `instagram:user` extractor is currently broken).
- OpenAI is the default STT provider through `/v1/audio/transcriptions`. Configure it with `[stt]` plus `OPENAI_API_KEY`, `OPENAI_API_URL`, `STT_PROVIDER`, and `STT_MODEL`.
- OpenRouter is an optional STT provider when `stt.provider = "openrouter"`; configure it with `OPENROUTER_API_KEY`, `OPENROUTER_API_URL`, and `openrouter.stt_model`.
- `OPENROUTER_API_KEY` (or `[openrouter].api_key`) is **required** by every decision-backed command; the key and `[openrouter].api_url` are shared by STT, decisions, and generation.
- `[decisions]`: `model` (default `typesafe/jev-1.13`, env `KB_DECISIONS_MODEL`), `deadline` (20s), `retries` (2), `concurrency` (8), `max_state_bytes` (98304), `budget_usd` (1.0, shared with generation; `--budget` overrides), `mode` (default `shadow`, body-link insertion).
- `[decisions.thresholds]`: `link_apply`, `link_review`, `mention_sense`, `affects`, `relevance_quarantine`, `relevance_review`, `relevance_fetch`, `quality_apply`, `quality_review`, `duplicate`, `primary_concept`, `concept_noul`, `find_keep`, `find_window`, `okf_type`, `contract_conflict`; `topic.yaml` `decisions.thresholds` overrides them per topic and `kb review calibrate --write` writes there.
- `[generation]`: `model` (default `xiaomi/mimo-v2.6-flash`, env `KB_GENERATION_MODEL`), `fallback_model` (`deepseek/deepseek-v4-flash`), `deadline` (90s), `retries`, `summary_language` (`en`).
- `[firecrawl].refetch_wait_ms` (3000) and `[firecrawl].refetch_only_main_content` (false) configure the fresh refetch (`maxAge: 0`) run before a quality gate judges a thin or broken capture.
- `topic.yaml` carries `contract`, `contract_draft`, `vocabulary_draft`, and `decisions: {mode, gates, relevance, thresholds, exclude}`. Relevance gates stay in shadow until the topic has an accepted contract and either a calibration with ≥ 30 review relevance labels or `decisions.gates: apply`.
- Long STT audio is segmented with `ffmpeg`; keep `[stt].ffmpeg_path`, `[stt].chunk_duration`, `[stt].max_chunk_bytes`, and `[stt].concurrency` aligned with provider upload limits.
- `[okf].types` configures the optional OKF type vocabulary used by `kb promote` and `kb okf check`; `[okf.type_descriptions]` gives the option text per type for decision-model type suggestions.

## Testing Conventions

- Default to table-driven tests with focused helpers and `t.TempDir()` for filesystem isolation.
- Keep integration tests co-located with the package under test behind `//go:build integration`.
- CLI integration tests should exercise real topic/ingest/lint/inspect flows instead of mocking Cobra wiring when the workflow itself is the behavior under test.
- QMD-related integration tests must isolate `HOME`, `XDG_CACHE_HOME`, and `XDG_CONFIG_HOME`.
- Decision-backed tests use `internal/fakes` (fake OpenRouter decisions and generation, fake Firecrawl) at the HTTP boundary; unit and integration tests never hit the network.
- The live question-bank regression suite runs only with `go test -tags jevlive ./internal/questions/` and `OPENROUTER_API_KEY` set; receipts are cached under `$KB_JEVLIVE_CACHE` (or the user cache dir) and spend is capped at US$ 0.50 per run.
- Treat test failures as behavior bugs first; do not weaken assertions to fit broken behavior.
