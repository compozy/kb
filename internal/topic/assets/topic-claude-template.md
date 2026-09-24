# TOPIC_TITLE

**Topic scope:** one-paragraph description of what this topic covers.

**Domain:** `TOPIC_DOMAIN` — all notes in this topic use `domain: TOPIC_DOMAIN` in frontmatter.

This file is the **schema document** for the topic (per Karpathy's LLM Wiki pattern). It tells the LLM how this topic is structured, its conventions, and its current state. Co-evolve it as the topic matures. Follow the shared conventions in the [root CLAUDE.md](../CLAUDE.md) for architecture, frontmatter, lifecycle, and tools — this file captures only topic-specific context.

TOPIC_SELECTION_CONTRACT

## Decision workflow

kb judges relevance, classifies, and links documents against the selection contract above with a decision model (requires `OPENROUTER_API_KEY`). `topic.yaml` is the contract's only source of truth: the section above is re-rendered from it on every accept, so edit the contract through `kb topic contract`, never here.

- `kb topic contract TOPIC_SLUG --draft|--import-claude|--accept` — draft or import the contract into `contract_draft`, then accept it after the self-check and impact preview.
- `kb topic vocabulary TOPIC_SLUG --draft|--accept` — while the topic has no articles, propose concepts and create stub articles for them.
- `kb ingest ... --topic TOPIC_SLUG` — every ingest is gated (dedupe, quality, relevance), classified, and linked; gated sources move to `raw/_quarantine/`, never deleted.
- `kb classify TOPIC_SLUG` — fill the kb-owned frontmatter facets (`summary`, `concepts`, `relevance`, ...); `--only-missing` after new articles.
- `kb link TOPIC_SLUG` — write typed links (`related`, `extends`, ...) between sources and articles; run it after compiling an article instead of a manual backlink audit.
- `kb find TOPIC_SLUG "<question>"` — retrieve the documents that answer a question; `--explain` shows why others were dropped.
- `kb review TOPIC_SLUG` — accept or reject the decisions queued for review (`kb review accept|reject --topic TOPIC_SLUG <id>...`); read the `recapture` and `remove` queues before bulk-accepting.
- `kb lint TOPIC_SLUG` — reports `needs-compile` (articles older than a source that `affects` them), pending review items, and key conflicts; it never calls a model.

kb writes its own frontmatter keys (`summary`, `concepts`, `genre`, `triage`, `related`, `affects`, `criterion`, `aliases`, ...). Articles are still written by the agent. A kb-owned key you edit becomes yours and kb stops updating it; `locked: true` makes kb skip a file entirely.

## Audit log

See [log.md](log.md) for the chronological record of every ingest / compile / query / lint operation. Append an entry there after each operation (skill Procedure 7).

## Current wiki articles

_None yet. See the [karpathy-kb skill](../.claude/skills/karpathy-kb/SKILL.md) Procedure 3 for how to compile new articles._

## Corpus inventory

- `raw/articles/` — 0 files
- `raw/bookmarks/` — 0 clusters
- `raw/github/` — 0 snapshots
- `wiki/concepts/` — 0 articles
- `bases/` — 0 base files

## qmd collection (optional)

At small scale, the topic's Concept Index / Source Index provide sufficient navigation. Once the corpus grows (~20+ sources or ~50+ wiki articles), add a qmd collection for hybrid BM25/vector search. Indexed as collection `TOPIC_SLUG`. Re-index from this directory after changes:

```bash
qmd collection remove TOPIC_SLUG && qmd collection add . --name TOPIC_SLUG && qmd embed
```

See `.claude/skills/karpathy-kb/references/tooling-tips.md` for Obsidian plugin tips (Web Clipper, Dataview, Marp).

## Research gaps

Priority areas to cover first:

- Gap 1
- Gap 2
