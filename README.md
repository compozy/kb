<div align="center">

# kb

### Build and maintain topic-based knowledge bases.

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/pedronauck/kodebase-go/ci.yaml?branch=main&label=CI)](https://github.com/pedronauck/kodebase-go/actions)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8.svg)](https://go.dev/)

[Install](#install) &#8226; [See It Work](#see-it-work) &#8226; [Features](#features) &#8226; [Decision Model](#decision-model) &#8226; [Commands](#commands) &#8226; [Contributing](#contributing)

</div>

---

`kb` is a single-binary Go CLI for building and maintaining topic-based knowledge bases in the [Karpathy KB](https://github.com/karpathy) pattern. It handles the non-LLM workflow: topic scaffolding, multi-source ingestion (URLs, files, YouTube, Instagram, codebases, bookmarks), ingest gates, classification, linking and retrieval through a typed [decision model](#decision-model), structural linting, codebase analysis, QMD indexing/search, and KB-oriented inspection commands. Writing wiki articles stays in your agent layer.

No SaaS. Just markdown. The only remote judgments are typed decision-model calls, recorded per topic and replayable.

---

## Install

> [!NOTE]
> `kb` works even better with its companion skill in [`skills/`](skills/). Install it with `npx skills add https://github.com/compozy/kb --skill kb`.

### Homebrew

```bash
brew install compozy/kb/kb
```

### npm

```bash
npm install -g @compozy/kb
```

### Go

```bash
go install github.com/compozy/kb/cmd/kb@latest
```

### Build from source

```bash
git clone https://github.com/compozy/kb.git
cd kb
make build
# binary is at bin/kb
```

**Optional** -- for semantic search capabilities:

```bash
npm install -g @tobilu/qmd
```

> [!NOTE]
> **Requirements:** Go >= 1.24. `OPENROUTER_API_KEY` is required by every command that ingests, classifies, links, finds, reviews or promotes (see [Decision model](#decision-model)). The `search` and `index` commands require [QMD](https://github.com/tobilu/qmd) to be installed separately. The `ingest url` command requires a [Firecrawl](https://firecrawl.dev) API key. The `ingest youtube` command requires [yt-dlp](https://github.com/yt-dlp/yt-dlp) for captions and audio extraction. STT transcription with `--transcribe auto` or `--transcribe stt` uses the configured `[stt]` provider; OpenAI audio transcriptions are the default and require `OPENAI_API_KEY`. Long audio is segmented with `ffmpeg`.

<details>
<summary><strong>What it touches</strong></summary>

- **Creates files** in `.kb/vault/` inside the target repository (or a custom `--vault` path)
- **Reads** source files in the target repository (never modifies them)
- **Network calls** -- `ingest url` calls the Firecrawl API; `ingest youtube`/`channel`/`instagram` call the platform through `yt-dlp`; `--transcribe auto|stt` also calls the configured STT provider. `ingest` (except `codebase`), `classify`, `link`, `find`, `review accept`, `topic contract --draft|--accept`, `topic vocabulary`, `promote` (when `[okf].types` is set) and `okf check --decide` send document excerpts to OpenRouter (decision and generation models). `topic list|info|new`, `lint`, `inspect`, `index`, `search`, `ingest codebase` and `version` never call a model.
- **Writes per topic** -- decision records under `<topic>/.decisions/`, kb-owned frontmatter keys on sources and articles, and quarantined sources under `<topic>/raw/_quarantine/`. Nothing is deleted.
- **No telemetry** -- nothing is sent anywhere
- **Uninstall:** Remove the `kb` binary from your `PATH` and delete the `.kb/` directory

</details>

---

## See It Work

Create a topic and ingest content from multiple sources:

```bash
# scaffold a new topic
$ kb topic new rust-lang "Rust Language" programming
{
  "slug": "rust-lang",
  "title": "Rust Language",
  "domain": "programming"
}

# scaffold an OKF bundle topic
$ kb topic new rust-catalog "Rust Catalog" programming --mode okf

# ingest a web article
$ kb ingest url https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html --topic rust-lang

# ingest a local PDF
$ kb ingest file ./rust-reference.pdf --topic rust-lang

# ingest a YouTube video transcript
$ kb ingest youtube https://www.youtube.com/watch?v=... --topic rust-lang

# ingest a codebase snapshot with full analysis
# first run bootstraps ./my-rust-project/.kb/vault/rust-lang automatically
$ kb ingest codebase ./my-rust-project --topic rust-lang --progress never

# lint the topic for structural issues
$ kb lint rust-lang

# promote a compiled wiki document into the OKF bundle
$ kb promote .kb/vault/rust-lang/wiki/concepts/Ownership.md --to rust-catalog --type Concept

# check the OKF bundle
$ kb okf check rust-catalog --strict
```

Judge, classify, link and query a topic with the decision model (requires `OPENROUTER_API_KEY`):

```bash
# write the topic's selection contract, check it against the topic's own sources, activate it
$ kb topic contract rust-lang --draft
$ kb topic contract rust-lang --accept

# fill facets (genre, depth, relevance, concepts, summary, ...) and typed links
$ kb classify rust-lang
$ kb link rust-lang

# ask which documents answer a question, with the reason every other candidate was dropped
$ kb find rust-lang "how does the borrow checker handle two-phase borrows?" --explain

# work through what landed in the review band
$ kb review rust-lang
```

Analyze codebase snapshots from the terminal:

```bash
$ kb inspect complexity --top 5

 symbol_name       | cyclomatic_complexity | loc | source_path
 computeMetrics    | 12                    | 89  | src/compute-metrics.ts
 parseTypeScript   | 9                     | 67  | src/adapters/typescript.ts
 normalizeGraph    | 8                     | 45  | src/normalize-graph.ts
 renderDocuments   | 7                     | 112 | src/render-documents.ts
 scanWorkspace     | 6                     | 54  | src/scan-workspace.ts
```

```bash
$ kb inspect dead-code

 kind   | name           | source_path                | reason
 symbol | oldHelper      | src/utils.ts               | dead-export
 file   | unused-cfg.ts  | src/config/unused-cfg.ts   | orphan-file
```

Search your vault with natural language (requires QMD):

```bash
$ kb search "error handling patterns" --limit 3

 title                  | score | path
 Error Handling Guide   | 0.89  | wiki/concepts/Error Handling.md
 src/error-boundary.ts  | 0.74  | raw/codebase/files/src/error-boundary.ts.md
 handleError            | 0.61  | raw/codebase/symbols/handleError--src-utils-l42.md
```

---

## Features

**Topic-based knowledge bases** -- `kb` organizes knowledge into topics, each with its own `raw/`, `wiki/`, `outputs/`, and `bases/` directories. Use `kb topic new` for manual scaffolding, or let `kb ingest codebase` bootstrap a new topic directly on first run.

**Multi-source ingestion** -- Ingest web articles via Firecrawl, local files (PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, plain text, images with OCR), YouTube captions or STT transcripts, codebases, and bookmark clusters. Each source type goes through a converter registry that normalizes content to frontmatter-annotated markdown.

**Codebase analysis** -- Point `kb ingest codebase` at a repository and it generates an [Obsidian](https://obsidian.md) vault layer with every symbol, file, and dependency relationship mapped into interconnected markdown notes. It computes cyclomatic complexity, blast radius, coupling, instability, and dead code detection, then compiles wiki articles and interactive [Base](https://obsidian.md/blog/bases/) views.

**Decision model** -- A typed decision model (Jev over OpenRouter) answers yes/no, choice and score questions about each document for about US$ 0.0003 per document. kb uses it to gate ingests against a written selection contract, classify every document into filterable frontmatter facets, write typed links between sources and articles, and rank documents that answer a question. Every call is recorded as a replayable receipt; thresholds can be re-cut without new calls; nothing is deleted.

**Structural linting** -- `kb lint` checks topics for missing frontmatter, broken wikilinks, orphaned files, and other structural health issues. Reports can be saved as markdown in `outputs/reports/`.

**10 inspect subcommands** -- Query your codebase like a database, from the terminal. Rank functions by complexity, find dead exports nobody imports, trace dependency chains, detect circular imports. Three output formats (table, JSON, TSV), zero external dependencies.

**Obsidian-native output** -- Every generated note uses wikilinks, YAML frontmatter, and backlink-aware cross-references. Base views give you filterable, sortable tables inside Obsidian -- Symbol Explorer, Complexity Hotspots, Danger Zone, Module Health, and more.

**Semantic search** -- Index your vault with QMD for hybrid, lexical, or vector search across all documentation. Useful for onboarding, architecture review, or feeding context to LLMs.

**AI-friendly output** -- Every vault includes `CLAUDE.md` and `AGENTS.md` as schema documents. The structured markdown output is designed for direct consumption by LLMs -- frontmatter metadata, consistent structure, and explicit cross-references.

**Single binary** -- No runtime dependencies. One `kb` binary handles everything. Built in Go for fast startup and low memory usage.

---

## Why kb

Most code analysis tools give you a dashboard. `kb` gives you a knowledge base.

**SonarQube and CodeClimate** tell you your code has problems. `kb` tells you which problems, where they connect, and gives you a structured workspace to reason about them. The output is markdown you own, not a SaaS dashboard that disappears when you cancel the subscription.

**Sourcegraph** is excellent for code search across repositories. `kb` is for understanding a single repository deeply -- its architecture, its coupling patterns, its risk surface -- and building a persistent knowledge artifact around that understanding.

**Obsidian and Notion** are great note-taking tools. `kb` automates the scaffolding and ingestion so you start with structure instead of a blank page, then extend the vault with your own notes and analysis.

The key difference: `kb` outputs compound. A SonarQube scan from last month is stale data. A `kb` vault from last month is a knowledge base you've been building on for 30 days.

---

## Decision Model

kb makes every semantic judgment (is this source on topic, is this page broken, which article does this document discuss, should these two documents link, does this document answer that question) with a **decision model**: a model that answers typed questions (yes/no, choice, score) with probabilities and generates no text. The default is `typesafe/jev-1.13` through OpenRouter's decisions endpoint (`POST <openrouter.api_url>/alpha/decisions`). A small **generation model** writes the few short literals kb needs (summaries, criteria, aliases, entities, questions, contract and vocabulary drafts) under a strict JSON schema; it never edits article prose. Code does everything else: it proposes candidates, computes every number, applies thresholds and writes files.

- **Bands.** Every answer falls in a confidence band, `apply`, `review` or `ignore`, cut by named thresholds. Review-band outcomes go to the topic's [review queue](#review-queues-and-labels).
- **Failures stay visible.** A timeout, exhausted retries, an invalid receipt or an exceeded budget yields `undecided:<reason>`; an oversized state yields `not_checked:state_too_large`. Neither is read as a "no": the document keeps its previous state and the run summary lists it.
- **Receipts.** Each call is appended to `<topic>/.decisions/receipts.jsonl` with its full answer distributions, model, cost and latency. The cache key covers model, question bank version, contract hash and state, so re-running a command costs nothing for unchanged documents, and re-cutting a threshold never needs a new call.
- **Run summary.** Every decision-backed command prints calls, cache hits, cost, decided/undecided/not_checked counts, bands per purpose and the mode of each gate.

### Requirement and configuration

The decision model is **required**. These commands refuse to start without `OPENROUTER_API_KEY` (or `[openrouter].api_key`) and name what is missing:

`kb ingest` (every subcommand except `codebase`), `kb classify`, `kb link`, `kb find`, `kb review accept|reject`, `kb topic contract --draft|--accept`, `kb topic vocabulary`, `kb promote` (when `[okf].types` is set) and `kb okf check --decide`.

These never call a model and run without it: `kb version`, `kb topic new|list|info`, `kb topic contract --import-claude`, `kb review list|import-links|import-labels|calibrate`, `kb lint`, `kb inspect`, `kb index`, `kb search`, `kb ingest codebase`, `kb migrate`.

```toml
[openrouter]                  # api_key/api_url are shared by STT, [decisions] and [generation]
api_key = ""                  # or OPENROUTER_API_KEY

[decisions]
model = "typesafe/jev-1.13"   # KB_DECISIONS_MODEL overrides
deadline = "20s"              # total per attempt
retries = 2                   # on 408/429/5xx; 401/402/403 stop the run
concurrency = 8
max_state_bytes = 98304
budget_usd = 1.0              # per-run ceiling shared with [generation]; --budget overrides
mode = "shadow"               # default body-link insertion mode

[decisions.thresholds]        # see "Gates, modes and thresholds"; topic.yaml can override
link_apply = 0.85

[generation]
model = "xiaomi/mimo-v2.6-flash"              # KB_GENERATION_MODEL overrides
fallback_model = "deepseek/deepseek-v4-flash"
deadline = "90s"
summary_language = "en"
```

**Privacy.** Judged content (a section-aware excerpt of each document plus its provenance) is sent to OpenRouter. The decisions route is listed as no-retention; the generation model follows its provider's data policy, and kb prints a notice naming both models the first time a topic makes a call. Obvious secrets (API keys, tokens, private keys) are redacted before any call. `topic.yaml` `decisions.exclude` globs keep files out of every call, and codebase snapshot documents never enter a decision pipeline.

**Cost.** Measured on real topics:

| Operation | Cost | Notes |
| --- | --- | --- |
| Gate and classification decisions | ≈ US$ 0.0003 per document | one request per document |
| `kb link` | ≈ US$ 0.0007 per document | measured on a 47-document topic |
| Generation (summary, entities, questions, aliases) | ≈ US$ 0.00015–0.00026 per document | only for documents that lack them or changed |
| `kb find` | ≈ US$ 0.001–0.002 per query | about 1 s wall time |
| Contract impact preview | ≈ US$ 0.15 | up to 500 sources, reused by the next `kb classify` |

A 1,000-document topic costs about US$ 0.30 to classify, US$ 0.30 to link and US$ 0.15 in generation. The default US$ 1 budget covers it in one or two runs; past the budget, remaining items are `undecided:budget` and the next run resumes through the cache.

### Selection contract

Relevance is judged against a **selection contract** stored in `topic.yaml` (the only source of truth) and rendered into the topic `CLAUDE.md` under `## Selection contract`:

```yaml
contract:
  purpose: "What this topic is for, in one or two sentences."
  core: ["Subjects that are the reason the topic exists"]
  adjacent: ["Neighbouring subjects that are kept, each with the reason it is kept"]
  collected_on_purpose: ["Material kept on purpose even if it looks off-topic (SDK docs, READMEs, ...)"]
  collected_on_purpose_paths: ["raw/brand-audit-*/**"]   # optional topic-relative globs
  out_of_scope: ["Clear junk only, each qualified so it cannot match a kept line"]
```

- `kb topic contract <topic> --draft` asks the generation model for a draft from the topic's own collection; `--import-claude` copies an existing `## Selection contract` section of the topic `CLAUDE.md`. Both write `contract_draft:` and never activate it.
- `--accept` runs a **self-check** (a kept line that an `out_of_scope` line would exclude blocks acceptance unless `--force`) and an **impact preview**: the relevance and quality questions over every source (a stratified sample of 500 when larger), printing counts per band and per `raw/` folder, the 20 most off-topic documents, and precision/recall against imported labels when there are any. The draft becomes active only after you type the topic slug (or pass `--yes`).
- Documents under `collected_on_purpose_paths` get `relevance: collected_on_purpose` from code, without a relevance call, and are never quarantined for relevance.
- `decisions.relevance: off` in `topic.yaml` disables relevance judgments for topics defined by construction (a directory mirror, one note per company); quality gates and classification still run.

The drafting rules the generation model follows, and that agents drafting by hand should follow, are in [`skills/kb/references/selection-contract.md`](skills/kb/references/selection-contract.md).

### Ingest gates

Every `kb ingest` subcommand except `codebase` runs the same staged pipeline; only survivors pay for the next stage.

1. **Exact dedupe** (code): normalized URL, platform video id and body hash against existing sources. `--force` skips the gates.
2. **Pre-fetch relevance** (bulk only: `channel`, `bookmarks`, URL lists): title and description are judged against the contract before anything is fetched. P(off_topic) ≥ 0.8 skips the item: it is recorded in `.decisions/skipped.jsonl`, never fetched, and can be rescued with `--rescue <id>`.
3. **Code quality checks**: `not_an_article` (root or listing URL, redirect to the site home, title equal to the site name), `thin` (< 300 words after navigation lines), `error_page` (HTTP ≥ 400 or a 404-like title).
4. **Refetch** (URL sources): a flagged or thin capture is scraped once more with Firecrawl `maxAge: 0`, `waitFor` (`[firecrawl].refetch_wait_ms`, default 3000) and `onlyMainContent` (`[firecrawl].refetch_only_main_content`, default false); the longer body is kept.
5. **Post-fetch quality and relevance**: `paywall_or_login`, `error_or_placeholder_page`, `thin_or_boilerplate`, `no_speech_content` (transcripts only) and the relevance role over the real content. ≥ 0.8 quarantines; 0.5–0.8 writes the source with `triage: review`.
6. **Near-duplicate**: `same_content` ≥ 0.8 quarantines with `triage_reason: duplicate`; `new_version` ≥ 0.8 keeps the source and writes `supersedes`.
7. **Write**, with `ingest_batch` and `ingest_query` provenance, then classify and link the new source in the same run.

The `log.md` entry of every ingest lists kept, review, quarantined and skipped counts and the cost.

### Gates, modes and thresholds

A **decision mode** says whether an outcome may move files or edit bodies: `shadow` records what would happen (in receipts and the run summary) and changes nothing; `apply` acts.

| What | Default | Becomes `apply` when | Override for one run |
| --- | --- | --- | --- |
| Exact dedupe and quality gates | `apply` | always | `--decisions shadow` |
| Relevance gates (skip, quarantine) | `shadow` | the topic has an accepted contract **and** either a stored calibration with ≥ 30 relevance labels given in `kb review` (imported labels do not count) or `decisions.gates: apply` in `topic.yaml` | `--decisions shadow\|apply` |
| Body-link insertion by `kb link` | `shadow` | `[decisions].mode`, then `topic.yaml` `decisions.mode`, is `apply` | `--decisions shadow\|apply` |

Frontmatter relations are written in both modes; `contradicts` never auto-applies.

Default thresholds (`[decisions.thresholds]`, overridden per topic by `topic.yaml` `decisions.thresholds`, which `kb review calibrate --write` fills):

| Threshold | Default | Meaning |
| --- | --- | --- |
| `relevance_quarantine` | 0.80 | P(off_topic) that quarantines (or skips before fetch) |
| `relevance_review` | 0.50 | P(off_topic) that sends a source to review |
| `relevance_fetch` | 0.50 | P(core + adjacent + collected on purpose) that fetches a bulk item |
| `quality_apply` / `quality_review` | 0.80 / 0.50 | quality question that quarantines / sends to review |
| `duplicate` | 0.80 | `same_content` or `new_version` |
| `link_apply` / `link_review` | 0.85 / 0.60 | `should_link` that writes a relation / queues it |
| `mention_sense` | 0.85 | a body mention refers to the linked article |
| `affects` | 0.80 | a source should change what an article says |
| `primary_concept` / `concept_noul` | 0.30 / 0.70 | concept written to `concepts` |
| `find_keep` / `find_window` | 0.50 / 0.35 | `kb find` keeps a candidate / maximum distance from the best |
| `okf_type` | 0.80 | `kb promote` type suggestion used |
| `contract_conflict` | 0.70 | self-check conflict that blocks `--accept` |

The `topic.yaml` decisions block:

```yaml
decisions:
  mode: shadow                  # body-link insertion: shadow | apply
  gates: shadow                 # relevance gates: shadow | apply (default: shadow until calibrated)
  relevance: on                 # off for topics whose membership is defined by construction
  thresholds:
    relevance_quarantine: 0.85
  exclude: ["raw/private/**"]   # never sent to any model
```

### Frontmatter written by kb

kb writes plain keys, flat at the top level, so Obsidian, Dataview and Bases read them directly. Probabilities stay in receipts; frontmatter holds banded outcomes only.

| Key | On | Shape | Written by |
| --- | --- | --- | --- |
| `triage` | source | `kept` \| `review` \| `quarantined` | gate |
| `triage_reason` | source (not kept) | `off_topic`, `paywall`, `error_page`, `thin`, `not_an_article`, `no_speech`, `duplicate` | gate |
| `genre` | source, article | `paper`, `article_or_essay`, `tutorial_or_guide`, `reference_docs`, `announcement_or_news`, `opinion_or_discussion`, `talk_or_interview`, `repository_or_code`, `dataset_or_benchmark`, `other` | classify |
| `depth` | source | 0–3 (0 mention, 1 overview, 2 detailed, 3 primary), one decimal | classify |
| `relevance` | source | `core` \| `adjacent` \| `collected_on_purpose` \| `general` \| `off_topic` \| `unknown` | classify |
| `concepts` | source, article | quoted wikilinks to articles | classify |
| `summary` | source, article | ≤ 400 characters | generation |
| `criterion` | article | ≤ 300 characters: what a document must discuss to cover the concept | generation |
| `entities` | source, article | ≤ 15 names copied verbatim from the body | generation, verified by code |
| `questions` | source | 3–6 questions the document answers | generation |
| `quality` | source | `thin`, `error_page`, `paywall`, `not_an_article`, `no_speech` | gate, classify |
| `ingest_batch`, `ingest_query` | source | the ingest run and the query, URL or file that produced the item | ingest |
| `related`, `extends`, `prerequisite`, `example_of`, `contradicts` | source, article | quoted wikilinks | link |
| `affects` | source | quoted wikilinks to articles the source should change | link |
| `supersedes` | source | quoted wikilinks to older sources | gate (near-duplicate) |
| `aliases` | article | strings; kb merges, never removes | classify |
| `recaptured` | source | date of an in-place recapture | review (recapture queue) |
| `locked` | any | `true` stops kb from writing anything to the file | you |

`status` and `kind` stay yours: kb uses `triage` and `genre` instead. **Ownership rule:** kb writes a key only when it is absent or its value is still exactly what kb last wrote (tracked by value hash in `.decisions/state.jsonl`). A key you edited, or one your vault already used, belongs to you: kb leaves it, the run summary reports `skipped:user-key`, and `kb lint` reports `key-conflict`. Relation lists (`related`, `extends`, `prerequisite`, `example_of`, `contradicts`, `affects`, `supersedes`) and `aliases` are the exception: kb only appends new targets to your list, never removes or reorders your items, and lint does not report them as conflicts. Writes are byte-preserving: kb replaces or appends only its own keys and refuses to write when the file changed since it was read (`skipped:changed`).

### Files under `<topic>/.decisions/`

Records and caches, never the graph (the graph is wikilinks and frontmatter). JSONL files are append-only.

| File | Holds |
| --- | --- |
| `receipts.jsonl` | every decision and generation call: key, purpose, subject, bank version, contract hash, model, cost, latency, status, raw answers |
| `state.jsonl` | one row per document: body hash, contract hash, bank versions, the value hash of every key kb wrote, and the body (and contract) each bank and key was judged on, so a write by one command never makes another command's judgment look current |
| `review.jsonl` | review queue items and their status changes |
| `labels.jsonl` | human verdicts and imported labels, with their origin |
| `calibration.json` | the last `kb review calibrate --write`: date, label counts, dev and holdout metrics, and per purpose the decision context (contract hash, model, bank versions) the scores were joined under |
| `previews.jsonl` | the documents shown in each contract impact preview |
| `skipped.jsonl` | bulk items skipped before fetch, rescuable with `--rescue <id>` |
| `quarantine-ledger.jsonl` | every move and index/frontmatter edit made by quarantine, replayed on restore |
| `inserted-links.jsonl` | body links inserted by `kb link` |
| `banks/*.json` | optional extra questions for this topic (same schema as the built-in banks, which they never replace) |

### Quarantine and restore

Gated sources move to `<topic>/raw/_quarantine/<path under raw/>` with `triage: quarantined` and `triage_reason`. Nothing is deleted. Quarantined files are excluded from classify, link, find and qmd indexing (`kb index` creates collections with a mask that skips `raw/_quarantine/` and `.decisions/`; a collection created before this release keeps its old pattern, so kb filters those hits and `kb index` warns: re-create it with `qmd collection remove <name>` then `kb index`). Before the move, kb removes the index lines in `wiki/index/*.md` whose every link points at quarantined files and the `sources:` and relation list entries that point at the file, recording each edit in `quarantine-ledger.jsonl`. Body links in other documents stay; lint reports them as `link-to-quarantined` instead of `dead-link`. `kb review accept` on a quarantined item moves the file back and replays the ledger in reverse; an entry whose file changed since is printed for manual repair instead of being forced.

### Review queues and labels

`kb review <topic>` (same as `kb review list <topic>`) lists pending items per queue with their evidence and probability.

| Queue | Holds |
| --- | --- |
| `gate` | ingested sources in the review band (`triage: review`) |
| `skip` | bulk items skipped before fetch; accept rescues and ingests them |
| `recapture` | sources with a quality flag (`thin`, `error_page`, `paywall`, `not_an_article`, `no_speech`); accept refetches fresh and re-judges: a good page replaces the body in place (`recaptured: <date>`), a still-broken one is quarantined, a source without `source_url` is only quarantined |
| `remove` | sources with P(off_topic) ≥ `relevance_review` and no quality flag; accept quarantines with `off_topic` |
| `link` | link decisions in the review band and body insertions proposed in shadow mode; accept writes them |
| `contradiction` | `contradicts` relations, which never auto-apply |
| `concept-proposal` | concept titles proposed when ≥ 5 sources match no article; accept creates a stub article |
| `okf-type` | `kb promote` type suggestions in the review band |

```bash
kb review <topic> [--queue <queue>] [--format table|json]
kb review accept --topic <topic> <id>...                        # act on items
kb review reject --topic <topic> <id>...                        # dismiss items
kb review accept --topic <topic> --queue remove --above 0.9     # bulk, by queue and probability
kb review accept --topic <topic> --all-purpose link --above 0.75
kb review accept --topic <topic> --queue recapture --topic-wide
```

Nothing leaves `raw/` without an accept. Every accept and reject is stored as a **label**. Existing human links become positive link labels with `kb review import-links <topic>` (links kb inserted, and relations kb wrote that nobody edited, are excluded). A topic's own screening or curation file becomes relevance labels with `kb review import-labels`; imported labels keep their origin and are reported separately.

### Calibration with a holdout

`kb review calibrate <topic>` measures precision, recall, coverage and review rate for the `relevance`, `link` and `quality` purposes from labels only, at the current thresholds and over a sweep. Labels are joined only to scores of the active decision context (contract, decision model, bank version); labels whose receipts belong to another context count as unjoined until the subjects are judged again, and a stored calibration stops enabling relevance `apply` once the contract, model or relevance bank changes. Labels are split by `sha256(subject) mod 10`: buckets 0–6 are **dev** (thresholds are chosen there, targeting 0.90 precision) and 7–9 are **holdout** (only reported). Labels shown in a contract's impact preview become dev-only once that contract changes. `--write` stores the chosen thresholds in `topic.yaml` `decisions.thresholds` and the record in `.decisions/calibration.json`; a purpose with fewer than 30 dev or 10 holdout labels reports "not enough labels" and nothing is written for it.

### Cleaning an existing topic

`kb classify` and `kb link` bring a topic that predates the decision model in line; both are incremental and resume through the cache. Cleaning never quarantines on its own: it fills the `recapture` and `remove` queues, and lint reports `off-topic-kept`.

1. `kb topic contract <topic> --import-claude` (or `--draft`), edit `contract_draft` in `topic.yaml`, then `--accept` and read the impact preview.
2. `kb review import-labels <topic> --from <file> ...` when the topic has its own screening files.
3. `kb classify <topic>`. For a topic without articles: `kb topic vocabulary <topic> --draft`, edit, `--accept`, then `kb classify <topic> --only-missing`.
4. `kb review <topic> --queue recapture` and accept or reject; then the same for `--queue remove`.
5. `kb link <topic>`.

Measured on a 34-topic research vault: US$ 2.03 for 8,588 sources (relevance and quality) plus about 160 Firecrawl credits to refetch 101 pages. The step-by-step procedure for agents is in [`skills/kb/references/cleaning-a-topic.md`](skills/kb/references/cleaning-a-topic.md).

---

## Commands

### `kb topic`

Scaffold and manage knowledge base topics.

```bash
kb topic new <slug> <title> <domain> [--mode wiki|okf]   # Create a new topic
kb topic list                           # List all topics in the vault
kb topic info <slug>                    # Show metadata for a topic
kb topic contract <slug> --draft | --import-claude | --accept [--force] [--yes] [--budget <usd>]
kb topic vocabulary <slug> --draft | --accept [--budget <usd>]
```

`--mode wiki` is the default scaffold; it writes an empty `contract:` block to `topic.yaml` and renders a `## Selection contract` section into the topic `CLAUDE.md`. `--mode okf` creates a root-level Open Knowledge File bundle with `index.md`, `log.md`, and OKF authoring guidance.

`kb topic contract` manages the [selection contract](#selection-contract): `--draft` asks the generation model for a draft, `--import-claude` copies the `CLAUDE.md` section (no model call), and `--accept` runs the self-check and impact preview, then activates the draft after you type the topic slug. `--force` continues past self-check conflicts; `--yes` skips the confirmation in scripts. Accepting re-renders the `CLAUDE.md` section from `topic.yaml`.

`kb topic vocabulary` bootstraps a concept vocabulary for a topic with sources but no articles: `--draft` proposes 12–40 concepts (title and criterion) from the contract and the source summaries into `topic.yaml` `vocabulary_draft:`, printing how many summaries mention each; `--accept` creates one stub article (`stage: stub`) per concept under `wiki/concepts/`. Run `kb classify <slug> --only-missing` afterwards to fill `concepts`.

### `kb promote`

Promote a compiled wiki document into an OKF topic without modifying the source document.

```bash
kb promote <wiki-doc> --to <okf-topic> [--type <type>] [--description <text>] [--budget <usd>]
```

`--type` is optional when `[okf].types` is configured: the decision model picks a type from that vocabulary (option text from `[okf.type_descriptions]`). A suggestion at P ≥ 0.8 is used; a weaker one prints the top candidates, is queued in the source topic's `okf-type` review queue, and asks for `--type`. With `--type`, kb warns when a confident suggestion disagrees. Without a vocabulary, `--type` is required and no model is called.

### `kb okf`

Check an OKF topic for bundle conformance.

```bash
kb okf check <topic> [--strict] [--format table|json|tsv] [--decide] [--budget <usd>]
```

Stored type suggestions add advisory `type_mismatch` and `description_unsupported` findings, never errors. `--decide` asks the decision model for findings missing from receipts (requires `[decisions]`).

### `kb ingest`

Ingest source material into a topic. `url`, `file`, `youtube`, `instagram`, and `bookmarks` require an existing topic; `channel` creates it by default; `codebase` can bootstrap one on first run. Every subcommand except `codebase` runs the [ingest gates](#ingest-gates), then classifies and links the new sources.

```bash
kb ingest url <url>... --topic <slug> [--from <file>]   # Scrape one or more web URLs (requires Firecrawl)
kb ingest file <path> --topic <slug>              # Convert and ingest a local file
kb ingest youtube <url> --topic <slug> [--transcribe captions|auto|stt] [--sub-langs orig,pt]
kb ingest channel <url> --topic <slug> [--limit N | --all] [--dry-run]
kb ingest instagram <url> --topic <slug> [--transcribe captions|auto|stt]
kb ingest codebase <path> --topic <slug>          # Analyze a codebase and bootstrap the topic if missing
kb ingest bookmarks <path> --topic <slug>         # Ingest a bookmark-cluster markdown file
```

**Supported file formats** for `ingest file`: PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, plain text (`.txt`, `.md`), and images (PNG, JPG, TIFF, BMP, GIF -- with optional OCR via Tesseract).

Decision flags shared by every subcommand except `codebase`:

| Flag | Description |
| --- | --- |
| `--batch <name>` | Value written to `ingest_batch` (default `<command>-<YYYY-MM-DD>-<run id>`) |
| `--budget <usd>` | Spend ceiling for this run (default `[decisions].budget_usd`) |
| `--decisions shadow\|apply` | Gate and body-link mode for this run (see [modes](#gates-modes-and-thresholds)) |
| `--force` | Skip dedupe and gates; the source is still classified and linked |
| `--rescue <id>` | Bulk only (`channel`, URL lists): fetch an item the pre-fetch gate skipped (ids from `kb review <slug> --queue skip`) |

| Ingest subcommand | `--topic` | Additional flags |
| --- | --- | --- |
| `url` | required | one or more URLs, `--from <file>` (one URL per line) |
| `file` | required | -- |
| `youtube` | required | `--transcribe captions\|auto\|stt` (default: `captions`), `--sub-langs`/`--lang` (caption languages; default: `orig`) |
| `channel` | required | `--limit`, `--all`, `--concurrency`, `--throttle`, `--dry-run`, `--create-topic`, `--transcribe`, `--sub-langs`/`--lang` |
| `instagram` | required | `--transcribe captions\|auto\|stt` (default: `auto`) |
| `codebase` | required | `--vault`, `--output` (deprecated alias), `--title`, `--domain`, `--include`, `--exclude`, `--semantic`, `--progress`, `--log-format` |
| `bookmarks` | required | -- |

### `kb classify`

Judge every document of a topic on the classification facets and write them as kb-owned frontmatter.

```bash
kb classify <topic> [--only-missing] [--all] [--budget <usd>] [--decisions shadow|apply] [--format table|json|tsv]
```

Writes `genre`, `depth`, `relevance`, `quality`, `concepts`, `summary`, `entities` and `questions` on sources, and `genre`, `concepts`, `summary`, `entities`, `criterion` and `aliases` on articles. Incremental: a document whose body, contract and question banks did not change is skipped without a call. `--only-missing` judges only documents without a state record or with a missing facet key; `--all` re-judges everything. Broken captures go to the `recapture` queue and off-topic sources to the `remove` queue; classify never quarantines. When ≥ 5 sources match no article, up to 5 concept proposals are queued. A topic without articles gets no `concepts` until `kb topic vocabulary` runs.

### `kb link`

Link a topic's documents to its wiki articles.

```bash
kb link <topic> [--all] [--dry-run] [--no-qmd] [--budget <usd>] [--decisions shadow|apply]
```

Code proposes up to 20 candidates per document (title and alias mentions, BM25, shared `concepts` or `sources`, and qmd vector neighbours when the topic's qmd collection is fresh); one decision request per document judges each candidate. `should_link` ≥ 0.85 writes the typed relation (`related`, `extends`, `prerequisite`, `example_of`) in frontmatter in both modes; in `apply` mode kb also inserts `[[File name|matched text]]` at the first qualifying mention in the body. `affects` ≥ 0.8 marks the articles a source should change. `contradicts` and review-band decisions go to the review queue. kb only adds links: a kb-written relation that later drops below the apply band moves to review, it is never deleted. Incremental unless `--all`; `--dry-run` prints the changes without writing files, state or review items; `--no-qmd` skips vector candidates.

### `kb find`

Find the documents of a topic that answer a question, ranked, with the reason every other candidate was dropped.

```bash
kb find <topic> "<question>" [--limit 10] [--kind <genre>] [--concept <title>] [--min-depth N] [--relevance <role>] [--explain] [--json] [--no-qmd] [--budget <usd>]
kb find --facets <topic>
```

Candidates come from frontmatter filters, BM25 over title, aliases, `summary`, `questions`, `entities` and excerpt, the `concepts` of the top hits, and qmd vectors when fresh. The decision model judges whether each candidate directly helps answer the question; code ranks by that probability, specificity, genre match and depth, keeps candidates with P ≥ 0.5 within 0.35 of the best, and names every drop (`low-probability`, `below-window`, `limit`, `filtered-facet`, `quarantined`) under `--explain`. About 1 s and US$ 0.001–0.002 per query. `--facets` prints counts per `genre`, `relevance`, `depth` and top `concepts` without any model call. `kb search` remains the raw qmd interface.

### `kb review`

Review queued decisions, import labels and calibrate thresholds. See [Review queues and labels](#review-queues-and-labels).

```bash
kb review <topic> [--queue <queue>] [--format table|json]         # same as `kb review list <topic>`
kb review accept --topic <topic> <id>... [--queue <queue>] [--all-purpose <purpose>] [--above <p>] [--topic-wide]
kb review reject --topic <topic> <id>... [--queue <queue>] [--all-purpose <purpose>] [--above <p>] [--topic-wide]
kb review import-links <topic>
kb review import-labels <topic> --from <file> --id-field <field> --match url|path|doi|pmcid --decision-field <field> [--keep-values include,true] [--decided-by-field <field>]
kb review calibrate <topic> [--write] [--format table|json]
```

Queues: `gate`, `skip`, `recapture`, `remove`, `link`, `contradiction`, `concept-proposal`, `okf-type`. `import-labels` reads JSONL, JSON or CSV, matches each row to a source by the chosen key and lists unmatched rows. `list`, `import-links`, `import-labels` and `calibrate` read files only and need no decision model.

### `kb lint`

Check a topic for structural KB issues.

```bash
kb lint [<slug>] [--format table|json|tsv] [--save] [--topic <slug>]
```

`--save` writes a markdown report to `outputs/reports/<date>-lint.md`.

### `kb inspect`

Query codebase vault data using frontmatter and extracted metrics.

```bash
kb inspect <subcommand> [options]
```

| Category | Subcommand      | Description                                          |
| -------- | --------------- | ---------------------------------------------------- |
| Metrics  | `smells`        | List symbols and files with detected code smells     |
| Metrics  | `dead-code`     | List dead exports and orphan files                   |
| Metrics  | `complexity`    | Rank functions by cyclomatic complexity              |
| Metrics  | `blast-radius`  | Rank symbols by blast radius (transitive dependents) |
| Metrics  | `coupling`      | Rank files by instability (efferent/afferent)        |
| Graph    | `backlinks`     | Show what references a given symbol                  |
| Graph    | `deps`          | Show outgoing relations for a file                   |
| Graph    | `circular-deps` | List files participating in circular dependencies    |
| Lookup   | `symbol`        | Fuzzy-match and detail view of a symbol              |
| Lookup   | `file`          | Exact lookup of a file by source path                |

Shared flags: `--format table|json|tsv` (default: `table`), `--vault <path>`, `--topic <slug>`.

### `kb lint` issue kinds

`kb lint` is read-only and never calls a model: it reads frontmatter, `topic.yaml` and the records under `.decisions/`. Links resolve the way Obsidian resolves them (topic-relative path, vault-relative path or file stem; titles and aliases never resolve a link), and frontmatter wikilinks count for orphan detection.

| Kind | Severity | Reported when |
| --- | --- | --- |
| `dead-link` | error | a body link resolves nowhere |
| `frontmatter-dead-link` | error | a frontmatter wikilink (`sources`, `concepts`, relations, ...) resolves nowhere |
| `link-to-quarantined` | warning | a link points at a file in `raw/_quarantine/`; restore it with `kb review accept` or edit the link |
| `orphan` | warning | no inbound links, body or frontmatter |
| `missing-source` | error | an article's `sources:` entry is missing from `raw/` |
| `stale` | warning | an article is older than a cited source's `scraped` date |
| `needs-compile` | warning | a source with `affects: [[A]]` was scraped after article A's `updated` (replaces `stale` for those pairs) |
| `format` | error / warning | missing required frontmatter, or a kb-owned key with the wrong shape (size limits are warnings) |
| `key-conflict` | warning | a kb-owned key holds a value kb did not write or that was edited since; kb leaves it alone |
| `off-topic-kept` | warning | a kept source has `relevance: off_topic` |
| `contradiction` | warning | an open `contradicts` review item |
| `contract-missing` | warning | the topic has sources and no accepted contract (relevance gates stay in shadow) |
| `contract-draft-pending` | info | `contract_draft` waits for `kb topic contract --accept` |
| `vocabulary-missing` | warning | the topic has sources and no articles; run `kb topic vocabulary --draft` |
| `criterion-missing` | warning | an article has no `criterion`; run `kb classify --only-missing` |
| `unclassified` | info | documents without a state record, or classified under another body, contract or bank version (one aggregated issue) |
| `pending-review` | info | pending items per review queue |
| `recapture-pending` / `remove-pending` | info | pending items in the two cleaning queues |

`contract-missing`, `vocabulary-missing`, `criterion-missing` and `unclassified` apply only to topics that adopted the decision workflow (records under `.decisions/`, an accepted contract, a contract draft or a vocabulary draft); other topics keep their structural lint.

### `kb search`

Semantic search across vault documents. Requires [QMD](https://github.com/tobilu/qmd).

```bash
kb search <query> [options]
```

| Flag           | Default | Description                              |
| -------------- | ------- | ---------------------------------------- |
| `--lex`        | false   | Use BM25 keyword search only             |
| `--vec`        | false   | Use vector similarity search only        |
| `--limit`      | 10      | Maximum results to return                |
| `--full`       | false   | Show full document instead of snippet    |
| `--min-score`  | --      | Minimum similarity threshold             |
| `--all`        | false   | Return all matches above threshold       |
| `--collection` | --      | Explicit QMD collection name             |
| `--format`     | table   | Output format: `table`, `json`, or `tsv` |

Default mode is hybrid (BM25 + vector similarity). Use `--lex` or `--vec` to restrict.
On hosts where QMD cannot run vector search because `sqlite-vec` is unavailable, the default hybrid mode falls back to lexical search automatically.

### `kb index`

Create or update a QMD collection for semantic indexing. Requires [QMD](https://github.com/tobilu/qmd).

```bash
kb index [options]
```

| Flag            | Default | Description                                |
| --------------- | ------- | ------------------------------------------ |
| `--vault`       | --      | Vault root path                            |
| `--topic`       | --      | Topic slug inside vault                    |
| `--name`        | --      | Override derived QMD collection name       |
| `--embed`       | true    | Run embedding after syncing files          |
| `--context`     | --      | Attach context to improve search relevance |
| `--force-embed` | false   | Force re-embedding of all documents        |

If QMD is installed without `sqlite-vec` support, `kb index` still syncs the collection and reports `embedStatus: "skipped_unavailable"` in JSON output so lexical search remains usable. Use `--force-embed` only when vector embedding is required on the current host.

### `kb version`

Print build version metadata.

---

## What Gets Generated

```text
.kb/vault/
  <topic-slug>/
    CLAUDE.md                    # Schema document for LLMs (with the rendered ## Selection contract)
    AGENTS.md                    # Agent-facing project reference
    topic.yaml                   # Topic metadata, selection contract, drafts, decisions settings
    log.md                       # Append-only operation log
    .decisions/                  # Receipts, state, review queue, labels, calibration, ledgers
    raw/
      _quarantine/               # Gated sources, recoverable with kb review accept
      codebase/                  # Machine-generated codebase snapshot
        files/                   #   One markdown note per source file
        symbols/                 #   One markdown note per extracted symbol
        indexes/
          directories/           #   Directory-level inventories
          languages/             #   Language-level inventories
      <ingested-sources>.md      # Ingested articles, transcripts, documents
    wiki/
      concepts/                  # Synthesized wiki articles (codebase topics)
        Codebase Overview.md
        Directory Map.md
        Symbol Taxonomy.md
        Dependency Hotspots.md
        Complexity Hotspots.md
        Module Health.md
        Dead Code Report.md
        Code Smells.md
        Circular Dependencies.md
        High-Impact Symbols.md
      index/
        Dashboard.md             # Landing page
        Concept Index.md         # Article listing
        Source Index.md          # Reverse index
    outputs/                     # Your analysis outputs (preserved)
      reports/                   # Lint reports (when --save is used)
    bases/                       # Obsidian Base views
```

- **`raw/`** -- Machine-generated source snapshots and ingested documents. Codebase snapshots refresh on every run; ingested documents are append-only. kb adds its own frontmatter keys to ingested documents (see [Frontmatter written by kb](#frontmatter-written-by-kb)) and never deletes them; gated sources move to `raw/_quarantine/`.
- **`.decisions/`** -- Decision records and caches for the topic (see [Files under `<topic>/.decisions/`](#files-under-topicdecisions)). Labels, the review queue and the quarantine ledger live only here, so keep it with the topic.
- **`wiki/`** -- Starter articles synthesized from codebase metrics. Managed areas refresh; your additions are preserved.
- **`outputs/`** -- A place for your own briefings, queries, diagrams, and reports. Never touched by `kb`.
- **`bases/`** -- Obsidian Base `.base` files for interactive table/card/list views of metrics.

---

## Supported Languages (Codebase Analysis)

| Language   | Extensions    | Parser        | Relation Confidence |
| ---------- | ------------- | ------------- | ------------------- |
| TypeScript | `.ts`, `.tsx` | `tree-sitter` | syntactic           |
| JavaScript | `.js`, `.jsx` | `tree-sitter` | syntactic           |
| Go         | `.go`         | `tree-sitter` | syntactic           |

Want to add a language? See [CONTRIBUTING.md](CONTRIBUTING.md#adding-a-new-language-adapter).

---

## Supported File Formats (Ingest)

| Format     | Extensions                        | Notes                                  |
| ---------- | --------------------------------- | -------------------------------------- |
| PDF        | `.pdf`                            | Native text extraction via pdfcpu      |
| DOCX       | `.docx`                           | XML-based extraction                   |
| XLSX       | `.xlsx`                           | Sheet-to-markdown table conversion     |
| PPTX       | `.pptx`                           | Slide text extraction                  |
| EPUB       | `.epub`                           | Chapter extraction with HTML-to-MD     |
| HTML       | `.html`, `.htm`                   | HTML-to-markdown conversion            |
| CSV        | `.csv`                            | Table conversion                       |
| JSON       | `.json`                           | Pretty-printed code block              |
| XML        | `.xml`                            | Pretty-printed code block              |
| Plain text | `.txt`, `.md`                     | Pass-through                           |
| Images     | `.png`, `.jpg`, `.tiff`, `.bmp`, `.gif` | Optional OCR via Tesseract       |

---

## Tracked Relations

| Relation     | Description                                  |
| ------------ | -------------------------------------------- |
| `imports`    | File imports another file or external module |
| `exports`    | File exports a symbol                        |
| `calls`      | Function or method calls another symbol      |
| `references` | Code references a symbol                     |
| `declares`   | File declares a symbol                       |
| `contains`   | File structurally contains a symbol          |

---

## Detected Code Smells

| Smell               | Scope  | Condition                                              |
| ------------------- | ------ | ------------------------------------------------------ |
| `dead-export`       | symbol | Exported but never referenced from outside its file    |
| `long-function`     | symbol | Function with > 50 LOC or cyclomatic complexity > 10   |
| `high-blast-radius` | symbol | More than 20 transitive dependents                     |
| `bottleneck`        | symbol | Betweenness centrality > 0.1                           |
| `feature-envy`      | symbol | More references to another file's symbols than its own |
| `god-file`          | file   | More than 15 symbols or efferent coupling > 10         |
| `orphan-file`       | file   | Zero afferent coupling and not an entry point          |

---

## Excluded Paths

The scanner:

- Always skips `.git`, `.hg`, `.svn`, symlinks, and the configured vault root itself
- Applies default convenience ignores for `vendor/`, `.turbo/`, `.next/`, `node_modules/`, `dist/`, `build/`, and `coverage/`
- Respects `.gitignore` files found from the scan root downward, including nested ones
- Applies `--exclude` patterns after repository ignore rules
- Applies `--include` patterns last as explicit re-includes

`--include` and `--exclude` use `.gitignore`-style patterns evaluated against paths relative to the scan root.

---

## Configuration

`kb` is primarily configured through CLI flags. Optional runtime configuration is loaded from a TOML file and environment variables.

| Variable | Source | Description |
| --- | --- | --- |
| `APP_CONFIG` | env | Path to TOML config file |
| `FIRECRAWL_API_KEY` | env / TOML | Firecrawl API key for `ingest url` |
| `FIRECRAWL_API_URL` | env / TOML | Firecrawl API endpoint |
| `OPENAI_API_KEY` | env / TOML | OpenAI STT API key for the default `openai` provider |
| `OPENAI_API_URL` | env / TOML | OpenAI-compatible STT API base URL |
| `STT_PROVIDER` | env / TOML | STT provider: `openai` or `openrouter` |
| `STT_MODEL` | env / TOML | STT model override, for example `gpt-4o-transcribe` |
| `OPENROUTER_API_KEY` | env / TOML | **Required** by decision-backed commands (decision and generation models); also the STT key when `stt.provider = "openrouter"` |
| `OPENROUTER_API_URL` | env / TOML | OpenRouter API endpoint |
| `KB_DECISIONS_MODEL` | env / TOML | Decision model override (`[decisions].model`, default `typesafe/jev-1.13`) |
| `KB_GENERATION_MODEL` | env / TOML | Generation model override (`[generation].model`, default `xiaomi/mimo-v2.6-flash`) |
| `YOUTUBE_YT_DLP_PATH` | env / TOML | `yt-dlp` executable path for `ingest youtube` |
| `YOUTUBE_PROXY` | env / TOML | Proxy URL passed to `yt-dlp` |
| `YOUTUBE_COOKIES_FILE` | env / TOML | Netscape cookies file passed to `yt-dlp` |
| `YOUTUBE_USER_AGENT` | env / TOML | User-Agent passed to `yt-dlp` |
| `YOUTUBE_CAPTION_LANGUAGES` | env / TOML | Comma-separated caption language preference list, for example `orig,pt,en` |

`ingest youtube` supports three transcription policies:
`captions` uses YouTube captions only, `stt` forces audio transcription, and `auto` uses manual captions when present and falls back to STT when only automatic captions exist or captions are unavailable. STT writes provenance frontmatter such as `transcript_source`, `stt_provider`, and `stt_model`.

Caption selection defaults to `caption_languages = ["orig"]`, so non-English videos use their native/original caption track instead of YouTube's machine-translated English. `orig` resolves per video from yt-dlp metadata, with `<lang>-orig` automatic tracks preferred when manual original-language captions are unavailable. Override per command with `--sub-langs` or `--lang`; CLI values override `YOUTUBE_CAPTION_LANGUAGES`, which overrides `[youtube].caption_languages`. Automatic translated captions are disabled unless `[youtube].allow_translated_captions = true`, because translated tracks use YouTube's translation endpoint and are more aggressively rate-limited.

YouTube ingests also persist video metadata from `yt-dlp` in the transcript frontmatter for sorting and filtering: engagement counts (`view_count`, `like_count`, `comment_count`), publication context (`upload_date`, `duration`, `duration_string`, `channel`, `channel_id`, `uploader_id`, `channel_follower_count`), and classification fields (`categories`, `youtube_tags`, `language`, `live_status`, `was_live`, `chapter_count`). Hidden or omitted scalar fields are written as `null`; `categories` and `youtube_tags` are empty lists when absent; `chapter_count` is `0` when no chapter metadata is returned. Video keywords use `youtube_tags` because `tags` remains the KB taxonomy field, and only the chapter count is stored, not full chapter content.

STT has real cost and latency. Prefer `captions` for normal ingestion, and use `auto` or `stt` when YouTube captions are unavailable or unsuitable.

TOML-only sections for the decision workflow: `[decisions]` (model, deadline, retries, concurrency, `max_state_bytes`, `budget_usd`, `mode`), `[decisions.thresholds]` (named thresholds, see [Gates, modes and thresholds](#gates-modes-and-thresholds)), `[generation]` (model, `fallback_model`, deadline, retries, `summary_language`), `[firecrawl].refetch_wait_ms` and `[firecrawl].refetch_only_main_content` (the fresh refetch of the quality gate), and `[okf.type_descriptions]` (option text per `[okf].types` entry for `kb promote` type suggestions). Per-topic overrides live in `topic.yaml` under `decisions:`.

See [`config.example.toml`](config.example.toml) for the full TOML schema.

---

## Development

**Prerequisites:** [Go](https://go.dev) >= 1.24

```bash
git clone https://github.com/pedronauck/kodebase-go.git
cd kodebase-go
make verify    # format + lint + test + build + boundaries
```

| Command                | Description                              |
| ---------------------- | ---------------------------------------- |
| `make fmt`             | Format all Go files with gofmt           |
| `make lint`            | Run golangci-lint with zero tolerance    |
| `make test`            | Unit tests with race detector            |
| `make test-integration`| Unit + integration tests                 |
| `make build`           | Build binary to `bin/kb`                 |
| `make verify`          | fmt -> lint -> test -> build -> boundaries|
| `make deps`            | Run `go mod tidy`                        |

See [CONTRIBUTING.md](CONTRIBUTING.md) for code style, testing requirements, and how to add a new language adapter.

---

## Contributing

`kb` is MIT-licensed and built in the open. We welcome contributions of all kinds:

- **Language adapters** -- Add support for Python, Rust, Java, or any language with a tree-sitter grammar
- **File converters** -- Add support for new file formats in the converter registry
- **New code smell detectors** -- The metrics engine is designed to be extended
- **Wiki article templates** -- Better starter articles mean better vaults out of the box
- **Bug reports and feature requests** -- [Open an issue](https://github.com/pedronauck/kodebase-go/issues), we read them all

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

---

## License

MIT -- see [LICENSE](LICENSE).
