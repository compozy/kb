---
name: kb
description: "Comprehensive skill for the `kb` CLI and the Karpathy Knowledge Base pattern. Covers the full KB lifecycle — topic scaffolding, multi-source ingestion (URLs, files, YouTube videos and channels, Instagram reels, bookmarks, codebases), wiki article compilation, cross-article querying with file-back, lint-and-heal passes, QMD indexing, and hybrid search. Also covers the OKF (Open Knowledge Format) dual-mode lifecycle: per-topic `mode: wiki|okf`, scaffolding portable OKF bundles, promoting compiled wiki concepts into a typed catalog with `kb promote`, the four-producer-field + relative-link contract, a local concept-type vocabulary, and OKF v0.1 conformance checking with `kb okf check`. Also covers codebase-specific analysis via inspect commands for complexity, coupling, blast radius, dead code, circular dependencies, symbol/file lookups, backlinks, and code smells. Use when working with kb CLI commands, knowledge base workflows, code vault generation, code graph analysis, code metrics inspection, wiki compilation, the ingest-compile-query-lint cycle, or the wiki→OKF distill loop (promote + conformance). Do not use for general code review, linting, formatting, building Go projects, or writing application code."
---

# kb CLI and Knowledge Base Pattern

Build and maintain a self-compiling Obsidian markdown knowledge base using the `kb` CLI. The LLM reads raw sources, writes cross-linked wiki articles, files Q&A results back into the corpus, and runs lint-and-heal passes. The CLI also supports codebase ingestion with deep inspection commands for code quality, architecture health, and symbol relationships.

Each **topic** lives in its own folder inside the Obsidian vault, either directly at the vault root (e.g. `go-best-practices/`) or nested by configured glob (e.g. `harness/goclaw/`). All topics share a single Obsidian vault, commonly the repo root. Read `references/architecture.md` for the full rationale and the four-phase pipeline (ingest → compile → query → lint).

Every topic has a lifecycle **`mode`** recorded in `topic.yaml` — **`wiki`** (the default) or **`okf`**:

- A **wiki topic** is the Karpathy research lab described throughout this skill. It contains `raw/`, `wiki/`, `outputs/`, and `bases/` subtrees plus topic-level `CLAUDE.md`, `topic.yaml`, and `log.md`, and follows the ingest → compile → query → lint loop.
- An **OKF topic** is a portable Open Knowledge Format catalog (a "bundle") that other people's agents and tools consume. It is **flat**: typed concept files at the bundle root, plus a generated `index.md`, a `log.md`, an OKF-flavored `CLAUDE.md`, and the `AGENTS.md` symlink — no `raw/`/`wiki/`/`outputs/` pyramid. See the **OKF Dual-Mode** section below and `references/okf-mode.md`.

Mode is opt-in per topic and invisible to existing users: absent/empty `mode` normalizes to `wiki`, so every existing topic behaves exactly as before.

The topic's **`CLAUDE.md`** is the **schema document** and topic marker — it tells the LLM the scope, conventions, current articles, and research gaps for that topic. `topic.yaml` is the structured source of truth for topic metadata (`slug`, `title`, `domain`). `AGENTS.md` may symlink to `CLAUDE.md` for Codex parity, but the valid-topic marker is `CLAUDE.md`.

## Prerequisites

1. Verify the `kb` binary is available:
   ```bash
   kb version
   ```
2. For search and index commands, verify QMD is installed:
   ```bash
   qmd --version
   # If missing: npm install -g @tobilu/qmd
   ```
3. Supported source languages for codebase analysis: TypeScript (`.ts`), TSX (`.tsx`), JavaScript (`.js`), JSX (`.jsx`), Go (`.go`).
4. For repo-root vaults, configure `kb.toml` at the vault root:
   ```toml
   [vault]
   root = "."
   topic_globs = ["*", "harness/*", "social-media/*"]
   ```

## Pattern Overview

Based on Andrej Karpathy's LLM Wiki pattern, the KB treats the LLM as a **compiler** that reads raw source documents and produces a structured, cross-linked markdown wiki. The four-phase loop:

1. **Ingest** — Scrape/curate sources via `kb` CLI → `raw/` (immutable staging)
2. **Compile** — LLM reads `raw/`, writes `wiki/concepts/` articles (3000-4000 words, dense wikilinks)
3. **Query** — Q&A against wiki → file answers to `outputs/queries/`, promote strong answers to wiki
4. **Lint** — Automated structural checks + LLM-driven semantic healing

Read `references/architecture.md` for the full rationale, context-window vs RAG tradeoffs, and multi-topic vault design.

## OKF Dual-Mode

`kb` supports two distinct knowledge lifecycles, selected per topic by `mode`:

- **`wiki`** (default) — a research **lab**: `ingest → compile → query → lint`. Obsidian `[[wikilinks]]`, the wiki frontmatter schema, the `raw/`+`wiki/`+`outputs/` pyramid. This is the lifecycle the rest of this skill documents.
- **`okf`** — a portable **catalog**: `declare → consume`. A flat bundle of typed concepts that conforms to the **Open Knowledge Format (OKF) v0.1**, meant to be shared and consumed by *other* people's agents and tools. Plain relative markdown links, the OKF producer-field contract, a generated `index.md`.

The two are not derived from each other; the high-leverage move is **distilling research into the catalog** — turning a finished wiki concept into a typed OKF concept. `kb` makes that one command (`kb promote`), mechanical and non-LLM, then validates the result (`kb okf check`).

**OKF bundle layout** (what `kb topic new --mode okf` scaffolds):

```
<bundle>/
  CLAUDE.md        # OKF-flavored schema/marker (topic marker)
  AGENTS.md        # symlink → CLAUDE.md
  topic.yaml       # mode: okf
  index.md         # generated; frontmatter okf_version: "0.1"; concepts grouped by type
  log.md           # "# Directory Update Log"; ISO-date (## YYYY-MM-DD) headings, newest first
  <concept>.md     # typed concept files, flat at the bundle root (added by promote/authoring)
```

There is no `raw/`/`wiki/`/`outputs/`/`bases/` in OKF mode. `index.md` and `log.md` are **auto-maintained by `kb`** (regenerated on scaffold and on every `promote`).

**OKF concept contract** — every concept `.md` carries four producer fields in frontmatter (emitted alphabetically): `description`, `timestamp` (RFC3339, UTC), `title`, `type`, plus optional `tags`. `type` is the only field OKF v0.1 strictly requires; the other three and the relative-link style are two deliberate, documented deviations matching Google's reference tooling. Concept bodies use **relative markdown links** (`[label](other-concept.md)`), never `[[wikilinks]]`.

Two deviations from the written OKF spec are intentional (ADR-002): emit the four producer fields (spec mandates only `type`) and emit relative links (spec recommends absolute `/path.md`, which breaks GitHub rendering).

Read `references/okf-mode.md` for the full bundle layout, frontmatter contract, promote transform rules, conformance ruleset, and the concept-path/link model.

## Related Skills

This skill orchestrates several companion skills for the LLM-driven phases:

- **[obsidian-markdown](https://github.com/pedronauck/skills/tree/main/skills/obsidian-markdown)** — author wiki articles with valid Obsidian Flavored Markdown (wikilinks, callouts, embeds, properties).
- **[obsidian-bases](https://github.com/pedronauck/skills/tree/main/skills/obsidian-bases)** — create `.base` files under `<topic>/bases/` for dashboard views, filters, and formulas.
- **[obsidian-cli](https://github.com/pedronauck/skills/tree/main/skills/obsidian-cli)** — interact with the running Obsidian vault from the command line (open notes, search, refresh indexes).

## kb CLI Quick Reference

### Topic management

```bash
kb topic new <slug> <title> <domain>              # scaffold a new wiki topic (mode: wiki, default)
kb topic new <slug> <title> <domain> --mode okf   # scaffold a flat OKF bundle (mode: okf)
kb topic list                                      # list all topics in the vault
kb topic info <topic-id>                           # topic metadata (counts, last log entry)
```

`--mode` accepts `wiki` (default) or `okf`. An OKF topic scaffolds the flat bundle layout (see the OKF Dual-Mode section), not the wiki pyramid.

### OKF bundles (promote + conformance)

```bash
kb promote <wiki-doc> --to <okf-topic> --type <Type>                  # distill a wiki concept into an OKF bundle
kb promote <wiki-doc> --to <okf-topic> --type <Type> --description "…" # override the generated description
kb okf check <okf-topic>                                               # validate OKF v0.1 conformance (lenient)
kb okf check <okf-topic> --strict --format json                       # promote local-standard warnings to errors (CI gate)
```

`kb promote` is **mechanical, non-LLM, and non-destructive**: it reads a compiled wiki document, remaps its frontmatter to the OKF producer contract, rewrites `[[wikilinks]]` to relative markdown links, writes a new typed concept at the bundle root, regenerates the bundle's `index.md`, and inserts a newest-first `log.md` entry. **The source wiki document is left untouched.** Both `--to` and `--type` are required, and `--to` must resolve to a `mode: okf` topic or `promote` errors before writing. The concept filename is a slug of the source document's base name (collisions get a `-2`, `-3` suffix), so its filename and every inbound link share one canonical key. `kb promote` emits the JSON `ConceptResult` (`writtenPath`, `type`, `linksRewritten`, `unresolvedLinks`, `warnings`).

`kb okf check` is **lenient by default** per OKF §9 (tolerates broken cross-links, unknown `type` values, missing optional fields) so externally produced bundles pass. It emits diagnostics as `severity, kind, filePath, target, message` and exits non-zero when any **error** is present (and on **warnings** too under `--strict`). Hard errors include a concept with a missing/empty `type` and unparseable frontmatter; local-standard **warnings** include a missing producer field (`title`/`description`/`timestamp`) and a `type` outside the configured vocabulary. The type vocabulary is the local standard in `kb.toml`:

```toml
[okf]
# Local OKF concept type vocabulary. Empty (the default) means `kb okf check`
# never warns about unknown types until you opt into a local standard.
types = ["Voice Profile", "Offer", "Playbook"]
```

When `[okf].types` is non-empty, a `--type` (or existing concept type) outside the list is a warning — `--strict` turns it into a CI-failing error, preventing type drift (`Voice Profile` vs `voice-profile`).

### Ingestion (auto-generates frontmatter, auto-appends to log.md)

```bash
kb ingest url <url> --topic <topic-id>        # scrape a web URL via Firecrawl
kb ingest file <path> --topic <topic-id>      # convert local file (PDF, DOCX, EPUB, HTML, images w/OCR, etc.)
kb ingest youtube <url> --topic <topic-id>    # extract a single YouTube transcript -> raw/youtube/
kb ingest channel <url> --topic <topic-id>    # bulk-extract a YouTube channel/playlist -> raw/youtube/
kb ingest instagram <url> --topic <topic-id>  # Instagram reel/video (caption + transcript) -> raw/instagram/
kb ingest bookmarks <path> --topic <topic-id> # ingest a bookmark-cluster markdown file
kb ingest codebase <path> --topic <topic-id>  # analyze a codebase into raw/codebase/
```

Use path-relative topic identifiers for nested topics, e.g. `--topic harness/goclaw`.

YouTube extraction uses `raw/youtube/` as the canonical transcript directory. `kb` requires `yt-dlp` for metadata, captions, and audio extraction. Install or update `yt-dlp` when public captions fail before treating the issue as only a proxy/cookie problem:

```toml
[youtube]
yt_dlp_path = "yt-dlp"
proxy = "http://127.0.0.1:8080"
cookies_file = "/path/to/youtube-cookies.txt"
user_agent = "Mozilla/5.0 ..."
retry_attempts = 3
retry_backoff = "1s"
transcription = "captions" # captions | auto | stt
```

Use `kb ingest youtube <url> --topic <topic-id> --transcribe auto` to use manual captions when present and STT when only automatic captions or no captions are available. Use `--transcribe stt` to force STT. The old shortcut flag is intentionally unsupported.

YouTube transcript frontmatter includes video metadata from `yt-dlp` for filtering and ordering, including `view_count`, `like_count`, `comment_count`, `upload_date`, `duration`, `duration_string`, `channel`, `channel_id`, `uploader_id`, `channel_follower_count`, `categories`, `youtube_tags`, `language`, `live_status`, `was_live`, and `chapter_count`. Use `youtube_tags` for video keywords because `tags` is KB taxonomy; see `references/frontmatter-schemas.md` for null, list, and count behavior.

`YOUTUBE_YT_DLP_PATH`, `YOUTUBE_PROXY`, `YOUTUBE_COOKIES_FILE`, and `YOUTUBE_USER_AGENT` override the matching TOML values for local runs. The CLI reports blocked caption/audio requests as `network_blocked`; after confirming `yt-dlp` is installed and current, treat that as a network or auth configuration issue, not a missing-transcript issue.

#### YouTube channels and playlists

`kb ingest channel <url> --topic <topic-id>` bulk-extracts a YouTube channel or playlist into `raw/youtube/`, writing one document per upload with the same `youtube-transcript` frontmatter as single videos. Flags: `--limit <n>` (newest n; `0` = all), `--all` (every upload, overrides `--limit`), `--concurrency <n>`, `--throttle <dur>`, `--dry-run` (resolve and list videos without ingesting), and `--transcribe captions|auto|stt`. The worker-pool size, inter-request throttle, adaptive-backoff ceiling, and per-video retries default from `[youtube].bulk_concurrency`, `[youtube].bulk_throttle`, `[youtube].bulk_backoff_max`, and `[youtube].bulk_retries`. Single-video URLs are rejected — use `kb ingest youtube` for those.

### Instagram ingestion

`kb ingest instagram <url> --topic <topic-id>` ingests Instagram reels, video posts (`/p/`), and IGTV (`/tv/`) into `raw/instagram/` using the same `yt-dlp` + STT engine as YouTube. The document **body** is the post caption under a `## Caption` heading followed by the spoken transcript under `## Transcript`. When the audio yields no transcript (e.g. a music-only reel) it degrades to a caption-only document with `transcript_source: none` and no `## Transcript` section.

```toml
[instagram]
yt_dlp_path = "yt-dlp"
proxy = ""
cookies_file = "/path/to/instagram-cookies.txt" # separate from YouTube — Instagram needs its own session
user_agent = ""
transcription = "auto" # captions | auto | stt (default auto: reels rarely carry caption tracks)
retry_attempts = 3
retry_backoff = "1s"
```

`[instagram]` mirrors the `[youtube]` yt-dlp knobs but keeps **cookies separate** because the two platforms need distinct authenticated sessions. Public reels often work without cookies, but Instagram rate-limits aggressively — set `cookies_file` (and optionally `proxy`) for reliable access. The default `--transcribe` is `auto` (manual caption when present, otherwise Whisper STT). Whole-profile bulk ingestion is out of scope because yt-dlp's `instagram:user` extractor is currently broken. Blocked requests are reported as `network_blocked`; remediate with `[instagram].cookies_file` / `[instagram].proxy` or a trusted network.

### Layout migrations

```bash
kb migrate transcripts --topic <topic-id> # move legacy raw/transcripts/*.md into raw/youtube/
```

### Codebase inspection

```bash
kb inspect smells [--type <smell-type>] --format json
kb inspect dead-code --format json
kb inspect complexity [--top N] --format json
kb inspect blast-radius [--min N] [--top N] --format json
kb inspect coupling [--unstable] --format json
kb inspect circular-deps --format json
kb inspect symbol <name> --format json
kb inspect file <path> --format json
kb inspect backlinks <name-or-path> --format json
kb inspect deps <name-or-path> --format json
```

### Structural linting

```bash
kb lint [<topic-id>] [--save]             # dead links, orphans, missing sources, format violations, stale content
```

### Indexing and search (requires QMD)

```bash
kb index --topic <topic-id>               # create or update QMD collection
kb search "<query>" --topic <topic-id>    # hybrid BM25 + vector search
kb search "<query>" --lex --topic <topic-id>  # keyword-only search
kb search "<query>" --vec --topic <topic-id>  # vector-only search
```

After running `kb ingest` or `kb lint --save`, the CLI auto-appends entries to `<topic>/log.md`. `kb promote` (wiki→OKF) also auto-maintains the **OKF bundle's** `log.md` and `index.md`. Manual log entries are still needed for compile, query, the wiki-internal query→wiki promotion, and split operations (Procedure 5).

## Command Dispatch

Map the user's intent to the correct command:

| Intent | Command |
|--------|---------|
| Scaffold a new wiki topic | `kb topic new <slug> <title> <domain>` |
| Scaffold a new OKF bundle | `kb topic new <slug> <title> <domain> --mode okf` |
| Distill a wiki concept into an OKF bundle | `kb promote <wiki-doc> --to <okf-topic> --type <Type>` |
| Check an OKF bundle for conformance | `kb okf check <okf-topic>` |
| Gate OKF conformance in CI | `kb okf check <okf-topic> --strict --format json` |
| List all topics | `kb topic list` |
| Scrape a web URL | `kb ingest url <url> --topic <topic-id>` |
| Ingest a local file (PDF, DOCX, etc.) | `kb ingest file <path> --topic <topic-id>` |
| Extract a YouTube transcript | `kb ingest youtube <url> --topic <topic-id>` |
| Bulk-ingest a YouTube channel/playlist | `kb ingest channel <url> --topic <topic-id>` |
| Extract an Instagram reel/video | `kb ingest instagram <url> --topic <topic-id>` |
| Migrate legacy transcripts | `kb migrate transcripts --topic <topic-id>` |
| Ingest bookmark clusters | `kb ingest bookmarks <path> --topic <topic-id>` |
| Analyze a codebase | `kb ingest codebase <path> --topic <topic-id> --progress never` |
| Find code smells | `kb inspect smells --format json` |
| Find dead exports and orphan files | `kb inspect dead-code --format json` |
| Rank functions by complexity | `kb inspect complexity --format json` |
| Find high-impact symbols (blast radius) | `kb inspect blast-radius --min 5 --format json` |
| Find unstable files (coupling) | `kb inspect coupling --unstable --format json` |
| Find circular imports | `kb inspect circular-deps --format json` |
| Look up a specific symbol | `kb inspect symbol <name> --format json` |
| Look up a specific file | `kb inspect file <path> --format json` |
| Find what depends on X (incoming refs) | `kb inspect backlinks <name-or-path> --format json` |
| Find what X depends on (outgoing deps) | `kb inspect deps <name-or-path> --format json` |
| Run structural lint | `kb lint <topic-id> --save` |
| Index vault for search | `kb index --topic <topic-id>` |
| Search the knowledge base | `kb search "<query>" --topic <topic-id> --format json` |

## Codebase Analysis Workflow

For codebase-specific analysis, the `kb ingest codebase` command must run before any inspect command.

**Workflow A -- Code Analysis (no QMD required):**
```
kb ingest codebase <path> --topic <topic-id> --> kb inspect <subcommand>
```

**Workflow B -- Full Pipeline (requires QMD):**
```
kb ingest codebase <path> --topic <topic-id> --> kb index --> kb search <query>
```

On first run, `kb ingest codebase` uses the discovered vault from `kb.toml` or legacy `.kb/vault/`. If no vault marker is discoverable from the current working directory, it falls back to bootstrapping under `<path>/.kb/vault/<topic-id>/`. Later commands auto-discover `kb.toml` or `.kb/vault/` from the current directory; otherwise pass `--vault <path>`.

### Ingest a Codebase

```bash
kb ingest codebase <path> --topic <topic-id> --progress never
```

Always use `--progress never` in agent contexts to prevent TTY progress bars from corrupting stdout.
Use `--title` and `--domain` only when bootstrapping a missing topic.

Parse the JSON output from stdout to extract key values:
- `topicSlug` -- the topic identifier for later commands
- `vaultPath` -- absolute path to the vault root
- `topicPath` -- absolute path to the topic directory
- `filesScanned`, `filesParsed`, `symbolsExtracted` -- summary statistics
- `diagnostics` -- check for warnings or errors

Stderr carries structured stage logs. Do not treat stderr content as failure evidence.

Key flags:
- `--vault <dir>` -- override vault root location
- `--output <dir>` -- deprecated alias for `--vault`
- `--topic <topic-id>` -- target topic id inside the vault; nested topics use a relative path such as `harness/goclaw`
- `--title <value>` -- bootstrap-only topic title override
- `--domain <value>` -- bootstrap-only topic domain override
- `--include <pattern>` -- re-include paths that would otherwise be ignored (repeatable)
- `--exclude <pattern>` -- exclude additional paths from scanning (repeatable)
- `--semantic` -- enable semantic analysis when adapters support it

Read `references/cli-ingest-codebase.md` for the full flag table and output schema.

### Inspect the Vault

Run inspect subcommands to analyze code quality and architecture.

**Shared flags for all inspect subcommands:**
- `--format json` -- always use JSON for programmatic parsing
- `--vault <path>` -- explicit vault root (omit to auto-discover from cwd)
- `--topic <topic-id>` -- explicit topic id, including nested relative paths such as `harness/goclaw` (omit if only one topic exists)

#### Tabular Subcommands

These return a list of rows sorted by the primary metric:

1. **smells** -- List symbols and files with detected code smells.
   ```
   kb inspect smells --format json
   kb inspect smells --type high-complexity --format json
   ```

2. **dead-code** -- List dead exports and orphan files.
   ```
   kb inspect dead-code --format json
   ```

3. **complexity** -- Rank functions/methods by cyclomatic complexity. Default top 20.
   ```
   kb inspect complexity --format json
   kb inspect complexity --top 50 --format json
   ```

4. **blast-radius** -- Rank symbols by transitive dependent count.
   ```
   kb inspect blast-radius --format json
   kb inspect blast-radius --min 10 --top 20 --format json
   ```

5. **coupling** -- Rank files by instability (Ce / (Ca + Ce)).
   ```
   kb inspect coupling --format json
   kb inspect coupling --unstable --format json
   ```

6. **circular-deps** -- List files participating in circular import chains.
   ```
   kb inspect circular-deps --format json
   ```

#### Detail Lookup Subcommands

These return field-value pairs for a single matched entity:

7. **symbol \<name\>** -- Case-insensitive substring match. Returns detail fields for a single match, or a summary table for multiple matches.
   ```
   kb inspect symbol parseConfig --format json
   ```

8. **file \<path\>** -- Exact source path lookup. Use the source-relative path as stored in vault frontmatter.
   ```
   kb inspect file src/config.ts --format json
   ```

#### Relation Subcommands

These return relation edges (`target_path`, `type`, `confidence`):

9. **backlinks \<name-or-path\>** -- Incoming references. Accepts a symbol name or file path.
   ```
   kb inspect backlinks parseConfig --format json
   ```

10. **deps \<name-or-path\>** -- Outgoing dependencies. Accepts a symbol name or file path.
    ```
    kb inspect deps src/config.ts --format json
    ```

Read `references/cli-inspect.md` for all column schemas and flag details.

### Index the Vault

Index the vault content into QMD for search. This step requires QMD on PATH.

```bash
kb index --topic <topic-id>
```

The command is idempotent: it checks whether the collection already exists and chooses `add` (create) or `update` (refresh) automatically.

Key flags:
- `--embed` (default true) -- run embedding after syncing files
- `--force-embed` -- force re-embedding all documents
- `--context <text>` -- attach human context to improve search relevance
- `--name <name>` -- override the derived collection name

Read `references/cli-search-index.md` for the full output schema.

### Search the Vault

Search indexed vault content with QMD. Requires a prior `kb index` run.

```bash
kb search "<query>" --topic <topic-id> --format json
```

**Search modes:**
- Hybrid (default) -- combines lexical and vector search
- Lexical (`--lex`) -- BM25 keyword search only
- Vector (`--vec`) -- embedding-based semantic search

The `--lex` and `--vec` flags are mutually exclusive. Omit both for hybrid mode.

Key flags:
- `--limit N` (default 10) -- maximum results
- `--min-score N` -- minimum relevance threshold
- `--full` -- return full document content instead of snippets
- `--all` -- return all matches above the minimum score

Read `references/cli-search-index.md` for full details.

## KB Maintenance Procedures

### Procedure 1: Compile a wiki article

1. Read `references/compilation-guide.md` to anchor on length, style, wikilink density, and sourcing rules.
2. Identify candidate sources via `kb search "<topic phrase>" --topic <topic-id>` or read `<topic>/wiki/index/Source Index.md`.
3. Load the candidate raw sources fully into context.
4. Load `<topic>/wiki/index/Concept Index.md` for orientation on existing articles and wikilink targets (including in other topics).
5. **Surface takeaways BEFORE drafting.** Present to the user: 3-5 key takeaways from the sources, the entities/concepts this article will introduce or update, and anything that contradicts existing wiki articles. Ask: *"Anything specific to emphasize or de-emphasize?"* Wait for the response. Skip this step only if the user has explicitly asked for autonomous compilation.
6. Write the article to `<topic>/wiki/concepts/<Article Title>.md` following the [obsidian-markdown skill](https://github.com/pedronauck/skills/tree/main/skills/obsidian-markdown) for wikilink, callout, and frontmatter syntax. Use the frontmatter schema from `references/frontmatter-schemas.md`. Target 3000-4000 words with a Sources section, wikilinks to related articles, and code or diagram blocks where applicable.
7. **Backlink audit -- do not skip.** Grep every existing article in `<topic>/wiki/concepts/` for mentions of the new article's title, aliases, or core entities. For each match, add a `[[New Article]]` wikilink at the first mention (and one later occurrence). This is the step most commonly skipped -- a compounding wiki depends on bidirectional links.
   ```bash
   grep -rln "<new article title or key term>" <topic>/wiki/concepts/
   ```
8. Update the topic's indexes (Procedure 2).
9. Update `<topic>/CLAUDE.md` current-articles list.
10. Re-index the topic's collection: `kb index --topic <topic-id>`.
11. Append an entry to `<topic>/log.md` (Procedure 5) -- e.g., `## [YYYY-MM-DD] compile | <Article Title> (<word_count> words, <N> sources)`.

When **updating an existing article** (rather than writing new), use the `Current / Proposed / Reason / Source` diff format and contradiction-sweep workflow described in `references/compilation-guide.md`.

### Procedure 2: Maintain topic indexes

After adding, renaming, or removing any wiki article:

1. `<topic>/wiki/index/Dashboard.md` -- update article count, total word count, featured sections, and any Obsidian Base embeds (use the [obsidian-bases skill](https://github.com/pedronauck/skills/tree/main/skills/obsidian-bases) to author `.base` files and embed them).
2. `<topic>/wiki/index/Concept Index.md` -- insert/update the article row alphabetically with its one-line summary.
3. `<topic>/wiki/index/Source Index.md` -- for each new article, append rows for every source it cites, with a wikilink back to the article.
4. Optionally refresh the live view in Obsidian with the [obsidian-cli skill](https://github.com/pedronauck/skills/tree/main/skills/obsidian-cli) (`obsidian open <path>`, `obsidian search <query>`).

### Procedure 3: Query the wiki and file back the answer

A query has two phases: **Phase A** produces the answer by reading the wiki (never from general knowledge); **Phase B** files the answer back so the exploration compounds.

**Precondition:** Identify which topic(s) the question belongs to. If the question spans topics, load each topic's Concept Index.

#### Phase A -- Answer from the wiki

1. **Read the topic's Concept Index first** (`<topic>/wiki/index/Concept Index.md`). Scan the full index to identify candidate articles. Do NOT answer from general knowledge -- the wiki is the source of truth, even when the answer seems obvious. A contradiction between the wiki and general knowledge is itself valuable signal.
2. **Locate relevant articles.** At small scale (<30 articles), the index is enough. At larger scale, supplement with `kb search "<phrase>" --topic <topic-id>`. Also grep the topic for keywords: `grep -rl "<keyword>" <topic>/wiki/concepts/`.
3. **Read the identified articles in full.** Follow one level of `[[wikilinks]]` when targets look relevant to the question. Stop at one hop -- deeper traversal wastes context.
4. **(Optional) Pull in raw sources** if an article's claim is ambiguous and its `sources:` frontmatter points at a specific raw file worth verifying.
5. **Synthesize the answer** with these properties:
   - Grounded in the wiki articles you just read -- every factual claim traces back to a `[[Wiki Article]]` citation.
   - Notes **agreements and disagreements** between articles when they exist.
   - Flags **gaps explicitly**: "The wiki has no article on X" or "[[Article Y]] does not yet cover Z".
   - Suggests follow-up **ingest targets** or open questions.
6. **Match format to question type:**
   - Factual → prose with inline `[[wikilink]]` citations.
   - Comparison → table with rows per alternative, citations in cells.
   - How-it-works → numbered steps with citations.
   - What-do-we-know-about-X → structured summary with "Known", "Open questions", "Gaps".
   - Visual → ASCII/Mermaid diagram, Marp deck (see `references/tooling-tips.md`), or matplotlib chart.

#### Phase B -- File back the answer

7. **Save the answer** to `<topic>/outputs/queries/<YYYY-MM-DD> <Question Slug>.md` with frontmatter: `type: output`, `stage: query`, `informed_by: ["[[Article 1]]", "[[Article 2]]"]`. See `references/frontmatter-schemas.md` for the full schema.
8. In the body, list which wiki articles informed the answer under `informed_by:` (as wikilinks) and call out new insights that should be absorbed back into those articles on the next compile pass.
9. When a filed-back insight contradicts or extends an article's claims, **recompile the affected articles** (Procedure 1).
10. **Promote to wiki when the synthesis is durable.** If the answer is a first-class reference (a comparison table, a trade-off analysis, a new concept synthesized from multiple articles), copy it to `<topic>/wiki/concepts/<Title>.md` following Procedure 1 standards and update the indexes (Procedure 2). Karpathy's pattern treats strong query answers as wiki citizens, not secondary artifacts.
11. **Append to `<topic>/log.md`** (Procedure 5) -- e.g., `## [YYYY-MM-DD] query | <Question Slug>` plus a second line `## [YYYY-MM-DD] promote | <Title>` if promoted.

**Anti-patterns to avoid:**

- **Answering from memory** -- always read the wiki pages. The wiki may contradict what you think you know.
- **No citations** -- every factual claim must trace back to a `[[wikilink]]`.
- **Skipping the save** -- good query answers compound the wiki's value. Always file to `outputs/queries/`; promote when durable.
- **Silent gaps** -- surface missing coverage explicitly so the next ingest pass can fill it.

### Procedure 4: Lint and heal

Run structural lint via the `kb` CLI:

```bash
kb lint <topic-id> --save
```

This checks dead wikilinks, orphan articles, missing source references, format violations, and stale content, saving a dated report to `<topic>/outputs/reports/`. For each issue, **propose the fix with a diff before applying** -- do not batch-apply changes:

- **Dead wikilink** -- either create the missing article (Procedure 1) or rewrite the wikilink to point at an existing article.
- **Orphan article** -- add incoming wikilinks from at least one related article, or remove the article if it is outside the topic's scope.
- **Missing source file** -- an article's `sources:` frontmatter references a file absent from `raw/`. Either re-ingest (`kb ingest url/file`) or correct the reference.
- **Stale content** -- article's `updated:` date is older than its source's `scraped:` date. Recompile with current sources.
- **Format violation** -- fix missing frontmatter fields, H1 title, lead paragraph, or Sources section.

For deeper LLM-driven self-healing checks (inconsistencies across articles, missing coverage, wikilink audits, filed-back query absorption), read `references/lint-procedure.md`.

After the heal pass, append `## [YYYY-MM-DD] lint | <N> issues found, <M> fixed` to `<topic>/log.md`.

### Procedure 5: Append to log.md

The `kb` CLI auto-appends log entries for `ingest`, `lint --save`, and `kb promote` (the latter writes to the **target OKF bundle's** `log.md`). Manual entries are needed for **compile**, **query**, **promote** (the wiki-internal query→wiki promotion of Procedure 3, distinct from the `kb promote` CLI command), and **split** operations.

**Format** -- each entry is a single H2 heading with a consistent prefix so the log stays grep-able:

```markdown
## [YYYY-MM-DD] <op> | <short description>
```

Where `<op>` is one of `compile`, `query`, `promote`, or `split` (ingest and lint are handled by `kb`).

**Examples:**

```markdown
## [2026-04-04] compile | Transformer Architecture (3847 words, 6 sources)
## [2026-04-04] query | 2026-04-04 flash-attention-vs-paged-attention.md
## [2026-04-04] promote | FlashAttention vs PagedAttention (from query)
## [2026-04-05] split | "Inference Optimization" → KV Cache, Speculative Decoding
```

Optionally add a body paragraph under each entry with more context (key findings, source urls, decisions made). Keep entries terse -- the log is for skimming, not prose.

**Quick recent-activity check** -- the consistent prefix lets unix tools query the log:

```bash
grep "^## \[" <topic>/log.md | tail -10                  # last 10 events
grep "^## \[.*compile" <topic>/log.md | wc -l            # total compiles
grep "^## \[2026-04" <topic>/log.md                      # April 2026 events
```

Keep `log.md` at the topic root (not inside `wiki/` or `outputs/`) so it sits alongside `CLAUDE.md` as a first-class topic artifact.

### Procedure 6: Distill a wiki concept into an OKF bundle

The wiki→OKF distill loop turns finished research into a portable, typed catalog. It is mechanical and CLI-driven; `kb` does the structural work and any intelligent rewriting/condensing is left to you or an external agent.

1. **Ensure a target OKF bundle exists.** If not, scaffold one: `kb topic new <slug> <title> <domain> --mode okf`. The target of `promote` **must** be a `mode: okf` topic.
2. **Pick the source.** Choose a compiled, durable `<topic>/wiki/concepts/<Article>.md` worth declaring in the catalog (promote operates on a single document; ingest-stage `raw/` files are not the intended source).
3. **Choose the type.** Pass `--type <Type>` from your local vocabulary (`[okf].types` in `kb.toml`). Off-vocabulary types are allowed but warn under `kb okf check`; extend the vocabulary by editing `kb.toml` rather than inventing drifted variants.
4. **Promote.**
   ```bash
   kb promote "<topic>/wiki/concepts/<Article>.md" --to <okf-topic> --type "<Type>" --description "<one-line summary>"
   ```
   `kb` writes the new concept (four producer fields + relative links) at the bundle root, regenerates `index.md`, and appends the bundle's `log.md`. The source wiki document is untouched. Omit `--description` to let `kb` fall back to the first body sentence (it warns if the body has none).
5. **Resolve links.** Inspect `unresolvedLinks` in the JSON result — these are `[[wikilinks]]` whose target concept has not been promoted yet (the link is still emitted and tolerated per §9). Promote those targets too when you want the link to resolve.
6. **Verify conformance.** Run `kb okf check <okf-topic>`; fix any errors and any warnings you care about. Use `--strict` in CI to enforce the producer fields and the type vocabulary.
7. **No manual log entry is needed** — `kb promote` already appended the OKF bundle's `log.md`.

Read `references/okf-mode.md` for the frontmatter remap table, the wikilink→markdown transform rules, the concept-path/collision model, and the full conformance ruleset.

## Output Format Selection

All `inspect` and `search` commands support `--format`:
- **json** -- always use for programmatic parsing
- **table** -- human-readable aligned columns (default)
- **tsv** -- tab-separated for piping to Unix tools

The `ingest codebase` and `index` commands always output JSON to stdout.

Read `references/output-formats.md` for format examples and empty result handling.

## Error Handling

### CLI Errors

| Error | Recovery |
|-------|----------|
| `unable to find a vault from <path>` | Run `kb ingest codebase <path> --topic <topic-id>` first, or re-run with `--vault <path>` if the vault lives elsewhere |
| YouTube caption extraction fails | Install/update `yt-dlp` or set `[youtube].yt_dlp_path`; for confirmed `network_blocked`, then configure proxy/cookies or use a trusted network |
| `QMD is not available` | Run `npm install -g @tobilu/qmd` |
| `no topics were found` | Run `kb ingest codebase` or `kb topic new` to populate the vault |
| `multiple topics were found` | Re-run with `--topic <topic-id>` |
| `--title and --domain are bootstrap-only` | Remove those flags when re-ingesting an existing topic |
| `no symbols matched "<query>"` | Use `inspect smells` or `inspect complexity` to discover valid names |
| `no file matched "<path>"` | Use exact source-relative path from vault frontmatter (e.g. `src/config.ts` not `./src/config.ts`) |
| `promote: target topic must use mode okf` | `--to` points at a wiki topic. Pass an existing `mode: okf` topic, or scaffold one with `kb topic new <slug> <title> <domain> --mode okf` |
| `promote: --type is required` / `required flag(s) "to", "type" not set` | Pass both `--to <okf-topic>` and `--type <Type>` |
| `promote: source document not found` | Use a vault-relative path to an existing wiki document, e.g. `<topic>/wiki/concepts/<Article>.md` |
| `okf check: found N issue(s)` (non-zero exit) | Read the diagnostics rows; fix `severity=error` concepts (missing/empty `type`, bad frontmatter). Under `--strict`, warnings (missing producer fields, off-vocabulary types) also fail |

### KB Workflow Errors

| Error | Recovery |
|-------|----------|
| `kb` not found | Install the `kb` binary and ensure it is on PATH. Verify with `kb version` |
| Topic not found | Run `kb topic list` to see available topics, or scaffold with `kb topic new` |
| Article exceeds 4000 words | Extract a sub-topic into its own article and wikilink to it |
| Cross-topic wikilink ambiguity | Disambiguate with full path: `[[other-topic/wiki/concepts/Article Name\|Display Name]]` |
| `log.md` missing in existing topic | Run any write operation that autocures the topic skeleton, or create manually and backfill from git: `git log --format='## [%ad] <op> \| %s' --date=short <topic>/` |

Read `references/error-handling.md` for the full error catalog with causes and recovery steps.

## Constraints

### MUST DO
- Run `kb ingest codebase` before any inspect command on that topic
- Use `--format json` when parsing output programmatically
- Use `--progress never` when running `kb ingest codebase` in a non-interactive context
- Parse stdout only for command output; treat stderr as diagnostics
- Use the `topicSlug` from ingest output for subsequent `--topic` flags; for nested topics this is the relative path
- Read `references/compilation-guide.md` before writing wiki articles
- Run backlink audits after every article compile (Procedure 1, step 7)
- File query answers to `outputs/queries/` (Procedure 3)
- Append manual log entries for compile, query, the wiki-internal query→wiki promotion, and split operations
- Use an existing `mode: okf` topic (or scaffold one with `--mode okf`) as the `--to` target of `kb promote`
- Pass both `--to` and `--type` to `kb promote`; keep `--type` within `[okf].types` to avoid type drift
- Run `kb okf check` after promoting; use `--strict` in CI to enforce producer fields and the type vocabulary

### MUST NOT DO
- Pass both `--lex` and `--vec` to `search`
- Pass `--force-embed` with `--embed=false` to `index`
- Treat stderr content as failure evidence for `kb ingest codebase`
- Assume vault location without running ingest or checking for `kb.toml` / `.kb/vault/`
- Use relative paths like `./src/config.ts` for `inspect file` -- use `src/config.ts` instead
- Answer wiki queries from general knowledge -- the wiki is the source of truth
- Skip the backlink audit when compiling articles
- Batch-apply lint fixes without proposing diffs first
- Promote into a wiki topic — `kb promote --to` must be a `mode: okf` topic
- Hand-edit an OKF bundle's `index.md` (it is regenerated by `kb promote`) or use `[[wikilinks]]`/absolute links in OKF concepts (use relative markdown links)
- Expect `kb promote` to rewrite or condense prose — it is mechanical and non-LLM; do any intelligent distillation yourself or with an external agent
