---
name: kb
description: "Operate the `kb` CLI to build and maintain topic-based Karpathy-style knowledge bases: gated ingestion, decision-model classification, linking, retrieval and review, wiki compilation, lint-and-heal, OKF bundles, QMD search, and codebase inspection. Use for kb commands and KB vault workflows; not for general code review, Go builds, or application code."
---

# kb CLI and Knowledge Base Pattern

`kb` maintains an Obsidian markdown knowledge base in the Karpathy LLM Wiki pattern. kb ingests sources, and its **decision model** gates, classifies and links them. You, the agent, compile wiki articles from the sources, answer questions from the wiki and file the answers back, and run lint-and-heal passes.

- **Topic.** Each topic is a folder in the vault: at the root (`go-best-practices/`) or nested by `[vault].topic_globs` (`harness/goclaw/`). `CLAUDE.md` is the topic marker and schema document (scope, conventions, current articles, gaps; `AGENTS.md` may symlink to it). `topic.yaml` holds `slug`, `title`, `domain`, `mode` and the **selection contract**. kb renders the contract into `CLAUDE.md` under `## Selection contract`; change it with `kb topic contract`, never in `CLAUDE.md`.
- **Mode.** `mode: wiki` (default; `raw/`, `wiki/`, `outputs/`, `bases/`, `log.md`, the ingest → compile → query → lint loop) or `mode: okf` (a flat, portable Open Knowledge Format bundle of typed concepts, filled by `kb promote`).
- **Division of labour.** The decision model only judges text (yes/no, choice, score) and a small generation model only writes short literals (summaries, criteria, aliases, drafts); code applies thresholds and writes files. You write and update article prose, decide what to ingest, draft or edit contracts, and give the human verdicts in `kb review`.

`references/architecture.md` has the rationale (compiler analogy, context window vs RAG, multi-topic vaults).

## Prerequisites

- `kb version` works. `qmd --version` for `kb index`/`kb search` and optional vector candidates (`npm install -g @tobilu/qmd`).
- `OPENROUTER_API_KEY` for every command that ingests (except `codebase`), classifies, links, finds, drafts or accepts a contract, builds a vocabulary, accepts or rejects review items, or suggests an OKF type. Those commands refuse to start without it. `lint`, `inspect`, `index`, `search`, `topic new|list|info`, `topic contract --import-claude`, `review list|import-links|import-labels|calibrate`, `find --facets`, `ingest codebase` and `migrate` never call a model.
- `FIRECRAWL_API_KEY` for `kb ingest url`; `yt-dlp` (and an STT key for `--transcribe auto|stt`) for YouTube and Instagram.
- A repo-root vault has `kb.toml` at its root:
  ```toml
  [vault]
  root = "."
  topic_globs = ["*", "harness/*"]
  ```

## Command map

| Intent | Command | Details |
|---|---|---|
| New wiki topic / OKF bundle | `kb topic new <slug> <title> <domain> [--mode okf]` | `okf-mode.md` |
| List topics, topic metadata | `kb topic list`, `kb topic info <topic>` | |
| Ingest a URL, URL list, file, YouTube video, channel, Instagram reel, bookmark cluster | `kb ingest url\|file\|youtube\|channel\|instagram\|bookmarks <input> --topic <topic>` | `cli-ingest-sources.md` |
| Write, import or accept the selection contract | `kb topic contract <topic> --draft\|--import-claude\|--accept` | `selection-contract.md` |
| Bootstrap concepts for a topic without articles | `kb topic vocabulary <topic> --draft`, then `--accept` | `decision-workflow.md` |
| Classify documents (facets, summaries, criteria, aliases) | `kb classify <topic> [--only-missing]` | `decision-workflow.md` |
| Write typed relations and backlinks | `kb link <topic> [--dry-run]` | `decision-workflow.md` |
| Find the documents that answer a question | `kb find <topic> "<question>" --explain` | `decision-workflow.md` |
| Count documents per genre, relevance, depth, concept | `kb find --facets <topic>` | |
| Work the review queues | `kb review <topic> [--queue q]`; `kb review accept\|reject --topic <topic> <id>...` | `decision-workflow.md` |
| Import screening files or human links as labels | `kb review import-labels <topic> --from <file> ...`; `kb review import-links <topic>` | `decision-workflow.md` |
| Calibrate thresholds | `kb review calibrate <topic> [--write]` | `decision-workflow.md` |
| Clean a topic that predates the decision model | five steps | `cleaning-a-topic.md` |
| Lint a topic | `kb lint <topic> [--save]` | `lint-procedure.md` |
| Index and search with QMD | `kb index --topic <topic>`; `kb search "<query>" --topic <topic> [--lex\|--vec] --format json` | `cli-search-index.md` |
| Distill a wiki concept into an OKF bundle | `kb promote <wiki-doc> --to <okf-topic> [--type <Type>]` | `okf-mode.md` |
| Check an OKF bundle | `kb okf check <okf-topic> [--strict] [--decide]` | `okf-mode.md` |
| Analyze a codebase, then inspect it | `kb ingest codebase <path> --topic <topic> --progress never`; `kb inspect <subcommand> --format json` | `cli-ingest-codebase.md`, `cli-inspect.md` |
| Move legacy transcripts | `kb migrate transcripts --topic <topic>` | `cli-ingest-sources.md` |

Use path-relative topic ids for nested topics (`--topic harness/goclaw`). For parsing, use `--format json` where a command offers it and read stdout only; stderr carries diagnostics and run summaries (`output-formats.md`).

## Working with the decision model

These rules apply to every decision-backed command; `references/decision-workflow.md` has flags, thresholds, queue actions, calibration and the `.decisions/` files.

- **Budget.** `--budget <usd>` caps a run (default US$ 1). `--budget 0` uses cached answers only. Past the budget, items are `undecided:budget`: re-run to continue from the cache in `<topic>/.decisions/receipts.jsonl`.
- **Failures are never a "no".** `undecided:<reason>` and `not_checked:<reason>` leave a document as it was, and an undecided ingest judgment writes `triage: review` with `triage_reason: undecided`. Report these statuses from the run summary; re-run instead of concluding anything from them.
- **Shadow before apply.** Relevance gates and body-link insertion run in shadow (recorded as review items, nothing moved) until the topic has an accepted contract and a current calibration or `decisions.gates: apply` / `decisions.mode: apply`. Quality gates and exact dedupe apply from the start. Never switch a topic to apply without the user.
- **Nothing is deleted.** Gated sources move to `raw/_quarantine/` with a ledger of every index line and frontmatter entry removed, so `kb review accept` restores them byte for byte. Quarantine and restore only through `kb review`, never by moving files.
- **Frontmatter ownership.** kb writes its own keys (`triage`, `genre`, `depth`, `relevance`, `concepts`, `summary`, `criterion`, `entities`, `questions`, `quality`, relation lists, `aliases`, ...; `frontmatter-schemas.md`). A key you or the user edit becomes user-owned and kb stops writing it (lint: `key-conflict`); relation lists and `aliases` only get appended to. `locked: true` makes kb skip a file entirely. `topic.yaml` `decisions.exclude` globs keep files out of every model call.
- **Human verdicts.** Review items are yours or the user's to decide. Read items before bulk-accepting (`--above`, `--topic-wide`); several wrong `remove` items mean the contract is too narrow. Every verdict becomes a calibration label.

## KB maintenance procedures

### Procedure 1: Compile a wiki article

1. Read `references/compilation-guide.md` for length, style, wikilink density, sourcing and update rules.
2. Find candidate sources: `kb find <topic> "<what the article must cover>" --explain`, sources whose `concepts:` or `affects:` name the article, `kb search "<phrase>" --topic <topic>`, or `<topic>/wiki/index/Source Index.md`. For a stub (`stage: stub`), its `criterion` says what the article must cover.
3. Load the candidate sources fully, and `<topic>/wiki/index/Concept Index.md` for existing articles and wikilink targets.
4. Before drafting, present 3–5 key takeaways, the concepts the article introduces or updates, and any contradiction with existing articles; ask what to emphasize, unless the user asked for autonomous compilation.
5. Write `<topic>/wiki/concepts/<Article Title>.md` (Obsidian Flavored Markdown; frontmatter from `frontmatter-schemas.md`; 3000–4000 words; Sources section). Updates to an existing article use the `Current / Proposed / Reason / Source` diff format from `compilation-guide.md`.
6. Run `kb classify <topic> --only-missing` (criterion and aliases for the new article), then `kb link <topic>`, and work `kb review <topic> --queue link`. kb finds the mentions and writes the backlinks; you do not hand-audit them.
7. Update the topic indexes (Procedure 2) and the current-articles list in `<topic>/CLAUDE.md`; re-index with `kb index --topic <topic>` when the topic uses QMD.
8. Append `## [YYYY-MM-DD] compile | <Article Title> (<words> words, <N> sources)` to `<topic>/log.md`.

### Procedure 2: Maintain topic indexes

After adding, renaming or removing a wiki article, update `<topic>/wiki/index/Dashboard.md` (counts, featured sections, Base embeds), `Concept Index.md` (alphabetical row with a one-line summary) and `Source Index.md` (every cited source with a wikilink back to the article).

### Procedure 3: Query the wiki and file back the answer

1. Read the topic's `Concept Index.md` first, then locate articles: the index at small scale (< 30 articles), otherwise `kb find <topic> "<question>" --explain`, plus `kb search` or `grep` for exact terms. Answer from the wiki, never from general knowledge; a contradiction between the two is itself a finding.
2. Read the identified articles in full, following one level of `[[wikilinks]]` when relevant. Pull a raw source only to verify an ambiguous claim.
3. Answer with a `[[wikilink]]` citation for every factual claim, note agreements and disagreements between articles, and state gaps explicitly ("the wiki has no article on X"). Match the format to the question: prose, comparison table, numbered steps, Known / Open questions / Gaps, or a diagram (`tooling-tips.md` for Marp).
4. Save the answer to `<topic>/outputs/queries/<YYYY-MM-DD> <Question Slug>.md` (`type: output`, `stage: query`, `informed_by:` wikilinks) and call out insights to absorb on the next compile.
5. When the answer contradicts or extends an article, recompile it (Procedure 1). When the synthesis is durable (a comparison, a trade-off analysis, a new concept), promote it to `wiki/concepts/` with Procedure 1 standards.
6. Append `## [YYYY-MM-DD] query | <Question Slug>` (and `## [YYYY-MM-DD] promote | <Title>` when promoted) to `log.md`.

### Procedure 4: Lint and heal

Run `kb lint <topic> --save` (never calls a model; the report goes to `<topic>/outputs/reports/`). Propose each fix as a diff before applying; do not batch-apply. `references/lint-procedure.md` maps every issue kind to its fix (dead links, orphans, missing sources, `needs-compile`, `link-to-quarantined`, `key-conflict`, the review-queue counts, missing contract, vocabulary or criterion, `unclassified`, format) and covers the semantic checks lint cannot do. Append `## [YYYY-MM-DD] lint | <N> issues found, <M> fixed` to `log.md`.

### Procedure 5: Log entries

kb appends `log.md` entries for `ingest` (one per run, with gate counts and cost), `lint --save` and `kb promote` (in the OKF bundle's log). Add the others yourself as one H2 line each, newest last, so the log stays greppable:

```markdown
## [YYYY-MM-DD] <compile|query|promote|split> | <short description>
```

`promote` here is the wiki-internal query → wiki promotion, not `kb promote`. `grep "^## \[" <topic>/log.md | tail -10` shows recent activity.

### Procedure 6: Distill a wiki concept into an OKF bundle

1. Target an existing `mode: okf` topic (or create one with `kb topic new ... --mode okf`); `kb promote --to` must be an OKF bundle.
2. Promote a compiled `<topic>/wiki/concepts/<Article>.md`: `kb promote "<path>" --to <okf-topic> [--type "<Type>"] [--description "<one line>"]`. Without `--type`, the decision model suggests one from `[okf].types`; a weak suggestion lists candidates and asks for `--type`.
3. Promote the targets listed in `unresolvedLinks` when those links should resolve.
4. Run `kb okf check <okf-topic>` (`--strict` in CI). No manual log entry is needed.

`kb promote` is mechanical: it remaps frontmatter, rewrites `[[wikilinks]]` to relative links, regenerates `index.md` and leaves the source untouched; condensing prose is your job. `references/okf-mode.md` has the producer-field contract and conformance rules.

### Procedure 7: Set up the decision workflow on a topic

- **New topic**, after the first ingests: `kb topic contract <topic> --draft` (or hand-write `contract_draft` following `selection-contract.md`), edit it, then `kb topic contract <topic> --accept`. Show the impact preview (bands, per-folder counts, the 20 most off-topic documents) to the user before the slug is typed; only the user decides whether the contract describes the collection. From then on every ingest is gated, classified and linked; check `kb review <topic>` after bulk ingests.
- **Topic that predates the decision model**: follow `references/cleaning-a-topic.md` (contract, labels, classify and vocabulary, the `recapture` then `remove` queues, link).
- **Calibration**: once the topic has ≥ 30 relevance labels given in `kb review` on the dev split and ≥ 10 on holdout, run `kb review calibrate <topic>`, show dev and holdout precision, and `--write` only with the user's approval. A later contract change makes the calibration stale and relevance gates return to shadow until it is redone.

## Codebase analysis

Run `kb ingest codebase <path> --topic <topic> --progress never` before any `kb inspect` command on that topic; parse its stdout JSON (`topicSlug`, `vaultPath`, `topicPath`, `diagnostics`) and treat stderr stage logs as diagnostics, not failures. Codebase snapshots never enter the decision pipeline. Inspect subcommands (`smells`, `dead-code`, `complexity`, `blast-radius`, `coupling`, `circular-deps`, `symbol`, `file`, `backlinks`, `deps`) take `--format json`, `--vault` and `--topic`; `inspect file` needs the exact source-relative path (`src/config.ts`). Flags and schemas: `cli-ingest-codebase.md`, `cli-inspect.md`.

## Constraints

- Change the selection contract only through `kb topic contract` or `topic.yaml` `contract_draft`; the `CLAUDE.md` section is overwritten on accept.
- Do not hand-edit kb-owned frontmatter keys unless you mean to take them over, and never edit files under `<topic>/.decisions/`.
- Do not move, delete or restore gated sources by hand; use `kb review accept|reject`.
- Do not set `decisions.gates: apply` or `decisions.mode: apply` without an impact preview, a calibration report and the user's approval.
- Hand-edit neither an OKF bundle's `index.md` (kb regenerates it) nor links inside OKF concepts into `[[wikilinks]]` or absolute paths; OKF uses relative markdown links.
- `kb search` takes at most one of `--lex` / `--vec`; `kb index` rejects `--force-embed` with `--embed=false`.

## References

| Read | When |
|---|---|
| `references/decision-workflow.md` | Flags of the decision commands, gates, modes and thresholds, review-queue actions, calibration, topic extra questions, `.decisions/` files |
| `references/selection-contract.md` | Drafting, importing or editing a contract |
| `references/cleaning-a-topic.md` | Bringing an existing topic under the decision workflow |
| `references/compilation-guide.md` | Writing or updating a wiki article |
| `references/frontmatter-schemas.md` | Frontmatter of any document type, `topic.yaml`, kb-owned keys and ownership rules |
| `references/lint-procedure.md` | Healing lint issues and the semantic checks |
| `references/cli-ingest-sources.md` | Source-specific ingest flags and configuration (URLs, files, YouTube, channels, Instagram, bookmarks) |
| `references/okf-mode.md` | OKF bundles, `kb promote`, `kb okf check` |
| `references/cli-search-index.md` | `kb index` and `kb search` |
| `references/cli-ingest-codebase.md`, `references/cli-inspect.md` | Codebase ingestion and inspection |
| `references/output-formats.md` | Output shapes and empty results |
| `references/error-handling.md` | Any error message or run-summary status |
| `references/architecture.md` | Why the pattern works, scale, multi-topic vaults |
| `references/tooling-tips.md` | Obsidian plugins, Web Clipper, Marp, Dataview |

Companion skills for the LLM-driven phases: [obsidian-markdown](https://github.com/pedronauck/skills/tree/main/skills/obsidian-markdown) (article syntax), [obsidian-bases](https://github.com/pedronauck/skills/tree/main/skills/obsidian-bases) (`.base` dashboards), [obsidian-cli](https://github.com/pedronauck/skills/tree/main/skills/obsidian-cli) (the running vault).
