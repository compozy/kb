# Lint and Heal Procedure

Run `kb lint <topic-id> --save` for automated checks. The report is saved to `<topic>/outputs/reports/` and a log entry is auto-appended. `kb lint` never calls a model; besides the structural kinds it reads the decision records under `<topic>/.decisions/`:

| Kind | Severity | Fix |
| --- | --- | --- |
| `dead-link` / `frontmatter-dead-link` | error | create the target or correct the link (frontmatter lists included) |
| `link-to-quarantined` | warning | the target was quarantined: restore it (`kb review accept`) or rewrite the prose |
| `orphan` | warning | add incoming links; frontmatter relations count |
| `missing-source`, `format` | error | fix the reference or the frontmatter shape |
| `stale` / `needs-compile` | warning | recompile the article with the affecting source (Check 1) |
| `key-conflict` | warning | a kb-owned key holds a value kb did not write; keep it (kb leaves it alone) or delete it so kb can write it |
| `contradiction` | warning | open `contradicts` review item (Check 2) |
| `off-topic-kept`, `remove-pending`, `recapture-pending`, `pending-review` | warning / info | work the review queues (Check 7) |
| `contract-missing`, `contract-draft-pending` | warning / info | draft or import, then accept the selection contract (`selection-contract.md`) |
| `vocabulary-missing`, `criterion-missing`, `unclassified` | warning / info | `kb topic vocabulary --draft/--accept`, `kb classify <topic-id> [--only-missing]` |

This document covers the deeper **LLM-driven checks** that require reading articles and applying judgment. Run them periodically or after a batch of new content.

### Check 1: Stale content

`kb lint` reports `needs-compile` for every article older than a source whose `affects:` names it, and `stale` for the other cited sources. For each flagged article:

1. Load the affecting sources (the issue's `target`) and recompile following *When updating an existing article* in `compilation-guide.md`.
2. For topics not yet linked (no `affects:`), compare the article's `updated:` date with each `sources:` entry's `scraped:` date.
3. Also flag articles where the topic has evolved rapidly (e.g., LLM model names, protocol versions) and the article has not been updated in 30+ days.

### Check 2: Inconsistencies across articles

Start from what kb already found: `kb lint` reports each open `contradicts` item as `contradiction`, and `kb review <topic-id> --queue contradiction` lists them with the two documents and the probability. `contradicts` never auto-applies. For each item, read both documents, decide which claim is right, fix the articles, then `kb review accept` (the relation is real) or `kb review reject` (it is not).

Then load groups of related articles (identified via shared `concepts`, relation lists such as `related`/`extends`, or wikilinks) and check for:

- Contradictory factual claims (e.g., "H100 has 80GB HBM3" vs "H100 has 80GB HBM2e")
- Inconsistent terminology (same concept called two different names across articles)
- Inconsistent formatting (some articles use tables, others prose, for the same kind of comparison)

Fix by picking the correct/canonical version and updating all affected articles.

### Check 3: Missing coverage

Scan all articles for wikilinks and identify targets that:

- Are referenced in 3+ articles
- Do not have their own article yet

These are strong candidates for new articles. For each:

1. Check whether relevant raw sources exist in `raw/`.
2. If yes, write the article (Procedure 1 in SKILL.md).
3. If no, ingest sources first (`kb ingest url/file`) or mark as a research gap in the topic's `CLAUDE.md`.

### Check 4: Format violations

Verify each article has:

- H1 title matching filename
- Lead paragraph
- Sources section at the bottom
- At least 5 wikilinks (outgoing)
- Frontmatter with all required fields (the `kb` CLI validates these automatically via `kb lint`)

Fix by rewriting or adding the missing elements.

### Check 5: Wikilink audit

Missing links are `kb link`'s job, not a manual grep:

1. Run `kb link <topic-id>` (incremental; `--all` relinks everything). It writes confident relations to frontmatter and, in `apply` mode, inserts body links at the first qualifying mention.
2. Work `kb review <topic-id> --queue link`: review-band links and, in `shadow` mode, proposed body insertions. Accept or reject; bulk-accept with `kb review accept --topic <topic-id> --all-purpose link --above 0.75` only after spot-checking.
3. Fix `frontmatter-dead-link` and `link-to-quarantined` issues from `kb lint`.

What still needs judgment, per article:

- Identify over-wikilinking (same term linked multiple times in close proximity)
- Identify wikilinks and relations to concepts that no longer match the linked article's actual content; kb never removes a link, so remove wrong ones by hand

### Check 6: Filed-back query absorption

Scan `<topic>/outputs/queries/` for recent query results. For each:

1. Identify the wiki articles listed under `informed_by:`.
2. Check whether the synthesis in the query result adds new insights not yet in those articles. `kb find <topic-id> "<the query's question>" --explain` shows which sources answer it; sources that answer it but are not cited by the `informed_by` articles are coverage gaps.
3. If yes, flag the articles for updates and absorb the insights on the next compile pass.
4. When a new article comes out of the query (a concept proposal accepted in `kb review --queue concept-proposal` creates a `stage: stub` article), compile the stub following `compilation-guide.md`.

This is the core compounding mechanism — query answers feeding back into the wiki.

### Check 7: Review queues and corpus hygiene

`kb lint` counts pending review items per queue (`pending-review`, `recapture-pending`, `remove-pending`) and lists kept sources judged off-topic (`off-topic-kept`). Work them with `kb review`:

1. `kb review <topic-id>` to see counts per queue; `--queue <name>` for one queue.
2. `recapture`: broken captures (`thin`, `error_page`, `paywall`, `not_an_article`, `no_speech`). Accept refetches and re-judges; a page still broken is quarantined.
3. `remove`: off-topic sources. Read titles and paths; a wrong item usually means the selection contract is too narrow, so fix the contract before accepting in bulk.
4. `gate`, `skip`, `link`, `concept-proposal`, `okf-type`: accept or reject item by item.

Nothing leaves `raw/` without an accept, and every verdict becomes a label for `kb review calibrate`. The full backfill order is in `cleaning-a-topic.md`.

## Lint report format

When running a manual lint pass, produce a report like:

```
LINT REPORT — <topic>/ — YYYY-MM-DD

DEAD LINKS (N)
  - [[Missing Article]] referenced in: Foo.md, Bar.md
    → SUGGEST: Create wiki/concepts/Missing Article.md
    → POTENTIAL SOURCES: raw/articles/relevant-source.md

ORPHAN ARTICLES (N)
  - Token Economics.md — 0 incoming links
    → SUGGEST: Add refs from Agent Infrastructure.md, Fine-Tuning.md

STALE CONTENT (N)
  - MCP article references "MCP spec v0.9" but raw/articles/mcp-spec.md is v1.2
    → UPDATE: Recompile with current spec

INCONSISTENCIES (N)
  - Hardware specs disagree: Agent Infrastructure.md vs Fine-Tuning.md
    → RESOLVE: Verify against authoritative source, pick canonical

MISSING COVERAGE (N)
  - "Inference Optimization" referenced in 4 articles, no article exists
    → SUGGEST: Create wiki/concepts/Inference Optimization.md

FORMAT VIOLATIONS (N)
  - Prompt Engineering Techniques.md — missing Sources section

FILED-BACK INSIGHTS (N)
  - outputs/queries/2026-04-02 memory vs context.md has synthesis not in Memory Systems.md
    → ABSORB: Update Memory Systems.md with the compaction tradeoffs insight
```

## Heal workflow

For each issue the lint report surfaces:

1. **Dead link + source available** → create the article (Procedure 1).
2. **Dead link + no source** → mark in topic `CLAUDE.md` research gaps, or rewrite the link.
3. **Link to quarantined** → restore the source with `kb review accept`, or rewrite the prose.
4. **Orphan** → run `kb link`, add incoming wikilinks, or delete if out-of-scope.
5. **Stale / needs-compile** → re-scrape source, recompile article.
6. **Inconsistency / contradiction** → find authoritative source, fix all affected articles, resolve the review item.
7. **Missing coverage** → ingest sources, write article.
8. **Format violation / key conflict** → fix formatting; leave or clear the conflicting key.
9. **Filed-back insight** → update affected wiki articles.
10. **Pending review** → work the queues (Check 7).

Run the cycle regularly. Each pass leaves the knowledge base in a better state than it found it.
