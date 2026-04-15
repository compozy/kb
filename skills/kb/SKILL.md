---
name: kb
description: Comprehensive skill for the `kb` CLI and the Karpathy Knowledge Base pattern. Covers the full KB lifecycle — topic scaffolding, multi-source ingestion (URLs, files, YouTube, bookmarks, codebases), wiki article compilation, cross-article querying with file-back, lint-and-heal passes, QMD indexing, and hybrid search. Also covers codebase-specific analysis via inspect commands for complexity, coupling, blast radius, dead code, circular dependencies, symbol/file lookups, backlinks, and code smells. Use when working with kb CLI commands, knowledge base workflows, code vault generation, code graph analysis, code metrics inspection, wiki compilation, or the ingest-compile-query-lint cycle. Do not use for general code review, linting, formatting, building Go projects, or writing application code.
---

# kb CLI and Knowledge Base Pattern

Build and maintain a self-compiling Obsidian markdown knowledge base using the `kb` CLI. The LLM reads raw sources, writes cross-linked wiki articles, files Q&A results back into the corpus, and runs lint-and-heal passes. The CLI also supports codebase ingestion with deep inspection commands for code quality, architecture health, and symbol relationships.

Each **topic** lives in its own top-level folder (e.g. `ai-harness/`) with `raw/`, `wiki/`, `outputs/`, `bases/` subtrees plus a topic-level `log.md` and `CLAUDE.md`. All topics share a single Obsidian vault at the repo root. Read `references/architecture.md` for the full rationale and the four-phase pipeline (ingest → compile → query → lint).

The topic's **`CLAUDE.md`** (symlinked to `AGENTS.md`) is the **schema document** — it tells the LLM the scope, conventions, current articles, and research gaps for that topic. Co-evolve it as the topic matures.

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

## Pattern Overview

Based on Andrej Karpathy's LLM Wiki pattern, the KB treats the LLM as a **compiler** that reads raw source documents and produces a structured, cross-linked markdown wiki. The four-phase loop:

1. **Ingest** — Scrape/curate sources via `kb` CLI → `raw/` (immutable staging)
2. **Compile** — LLM reads `raw/`, writes `wiki/concepts/` articles (3000-4000 words, dense wikilinks)
3. **Query** — Q&A against wiki → file answers to `outputs/queries/`, promote strong answers to wiki
4. **Lint** — Automated structural checks + LLM-driven semantic healing

Read `references/architecture.md` for the full rationale, context-window vs RAG tradeoffs, and multi-topic vault design.

## Related Skills

This skill orchestrates several companion skills for the LLM-driven phases:

- **[obsidian-markdown](https://github.com/pedronauck/skills/tree/main/skills/obsidian-markdown)** — author wiki articles with valid Obsidian Flavored Markdown (wikilinks, callouts, embeds, properties).
- **[obsidian-bases](https://github.com/pedronauck/skills/tree/main/skills/obsidian-bases)** — create `.base` files under `<topic>/bases/` for dashboard views, filters, and formulas.
- **[obsidian-cli](https://github.com/pedronauck/skills/tree/main/skills/obsidian-cli)** — interact with the running Obsidian vault from the command line (open notes, search, refresh indexes).

## kb CLI Quick Reference

### Topic management

```bash
kb topic new <slug> <title> <domain>     # scaffold a new topic
kb topic list                             # list all topics in the vault
kb topic info <slug>                      # topic metadata (counts, last log entry)
```

### Ingestion (auto-generates frontmatter, auto-appends to log.md)

```bash
kb ingest url <url> --topic <slug>        # scrape a web URL via Firecrawl
kb ingest file <path> --topic <slug>      # convert local file (PDF, DOCX, EPUB, HTML, images w/OCR, etc.)
kb ingest youtube <url> --topic <slug>    # extract YouTube transcript
kb ingest bookmarks <path> --topic <slug> # ingest a bookmark-cluster markdown file
kb ingest codebase <path> --topic <slug>  # analyze a codebase into raw/codebase/
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
kb lint [<slug>] [--save]                 # dead links, orphans, missing sources, format violations, stale content
```

### Indexing and search (requires QMD)

```bash
kb index --topic <slug>                   # create or update QMD collection
kb search "<query>" --topic <slug>        # hybrid BM25 + vector search
kb search "<query>" --lex --topic <slug>  # keyword-only search
kb search "<query>" --vec --topic <slug>  # vector-only search
```

After running `kb ingest` or `kb lint --save`, the CLI auto-appends entries to `<topic>/log.md`. Manual log entries are still needed for compile, query, promote, and split operations (Procedure 5).

## Command Dispatch

Map the user's intent to the correct command:

| Intent | Command |
|--------|---------|
| Scaffold a new topic | `kb topic new <slug> <title> <domain>` |
| List all topics | `kb topic list` |
| Scrape a web URL | `kb ingest url <url> --topic <slug>` |
| Ingest a local file (PDF, DOCX, etc.) | `kb ingest file <path> --topic <slug>` |
| Extract a YouTube transcript | `kb ingest youtube <url> --topic <slug>` |
| Ingest bookmark clusters | `kb ingest bookmarks <path> --topic <slug>` |
| Analyze a codebase | `kb ingest codebase <path> --topic <slug> --progress never` |
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
| Run structural lint | `kb lint <slug> --save` |
| Index vault for search | `kb index --topic <slug>` |
| Search the knowledge base | `kb search "<query>" --topic <slug> --format json` |

## Codebase Analysis Workflow

For codebase-specific analysis, the `kb ingest codebase` command must run before any inspect command.

**Workflow A -- Code Analysis (no QMD required):**
```
kb ingest codebase <path> --topic <slug> --> kb inspect <subcommand>
```

**Workflow B -- Full Pipeline (requires QMD):**
```
kb ingest codebase <path> --topic <slug> --> kb index --> kb search <query>
```

On first run, `kb ingest codebase` bootstraps the topic under `<path>/.kb/vault/<topic-slug>/` by default. Later commands auto-discover this vault only when they run from inside that repository tree; otherwise pass `--vault <path>`.

### Ingest a Codebase

```bash
kb ingest codebase <path> --topic <slug> --progress never
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
- `--topic <slug>` -- target topic slug inside the vault
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
- `--topic <slug>` -- explicit topic slug (omit if only one topic exists)

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
kb index --topic <slug>
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
kb search "<query>" --topic <slug> --format json
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
2. Identify candidate sources via `kb search "<topic phrase>" --topic <slug>` or read `<topic>/wiki/index/Source Index.md`.
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
10. Re-index the topic's collection: `kb index --topic <slug>`.
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
2. **Locate relevant articles.** At small scale (<30 articles), the index is enough. At larger scale, supplement with `kb search "<phrase>" --topic <slug>`. Also grep the topic for keywords: `grep -rl "<keyword>" <topic>/wiki/concepts/`.
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
kb lint <slug> --save
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

The `kb` CLI auto-appends log entries for `ingest` and `lint --save` operations. Manual entries are needed for **compile**, **query**, **promote**, and **split** operations.

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
| `unable to find a vault from <path>` | Run `kb ingest codebase <path> --topic <slug>` first, or re-run with `--vault <path>` if the vault lives elsewhere |
| `QMD is not available` | Run `npm install -g @tobilu/qmd` |
| `no topics were found` | Run `kb ingest codebase` or `kb topic new` to populate the vault |
| `multiple topics were found` | Re-run with `--topic <slug>` |
| `--title and --domain are bootstrap-only` | Remove those flags when re-ingesting an existing topic |
| `no symbols matched "<query>"` | Use `inspect smells` or `inspect complexity` to discover valid names |
| `no file matched "<path>"` | Use exact source-relative path from vault frontmatter (e.g. `src/config.ts` not `./src/config.ts`) |

### KB Workflow Errors

| Error | Recovery |
|-------|----------|
| `kb` not found | Install the `kb` binary and ensure it is on PATH. Verify with `kb version` |
| Topic not found | Run `kb topic list` to see available topics, or scaffold with `kb topic new` |
| Article exceeds 4000 words | Extract a sub-topic into its own article and wikilink to it |
| Cross-topic wikilink ambiguity | Disambiguate with full path: `[[other-topic/wiki/concepts/Article Name\|Display Name]]` |
| `log.md` missing in existing topic | Create manually and backfill from git: `git log --format='## [%ad] <op> \| %s' --date=short <topic>/` |

Read `references/error-handling.md` for the full error catalog with causes and recovery steps.

## Constraints

### MUST DO
- Run `kb ingest codebase` before any inspect command on that topic
- Use `--format json` when parsing output programmatically
- Use `--progress never` when running `kb ingest codebase` in a non-interactive context
- Parse stdout only for command output; treat stderr as diagnostics
- Use the `topicSlug` from ingest output for subsequent `--topic` flags
- Read `references/compilation-guide.md` before writing wiki articles
- Run backlink audits after every article compile (Procedure 1, step 7)
- File query answers to `outputs/queries/` (Procedure 3)
- Append manual log entries for compile, query, promote, and split operations

### MUST NOT DO
- Pass both `--lex` and `--vec` to `search`
- Pass `--force-embed` with `--embed=false` to `index`
- Treat stderr content as failure evidence for `kb ingest codebase`
- Assume vault location without running ingest or checking for `.kb/vault/`
- Use relative paths like `./src/config.ts` for `inspect file` -- use `src/config.ts` instead
- Answer wiki queries from general knowledge -- the wiki is the source of truth
- Skip the backlink audit when compiling articles
- Batch-apply lint fixes without proposing diffs first
