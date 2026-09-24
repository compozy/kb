---
id: 0
title: kb-jev — decision model and cheap LLM in kb
labels: [wayfinder:map]
status: closed
assignee: ""
---

## Destination

An implementation spec, ready to hand off, that makes a decision model (Jev) a required core of kb, plus a cheap generation model: ingest gates, a classification system (facets), incremental linking, retrieval (`kb find`), OKF type suggestion, and lint over stored decisions. Reached: [spec.md](spec.md).

## Notes

- Domain: kb, the Go CLI in this repo. Vocabulary lives in `CONTEXT.md`: use Link, Relation, Affects, Candidate, Decision, Decision model, Confidence band, Decision mode.
- Grilling tickets: call the `grilling` and `domain-modeling` skills. Design tickets also consult `codebase-design` and `golang-master`.
- Division of labour: code builds candidates and applies results; the decision model judges text (choice / noul / score); the generation model only writes text (aliases, summaries, contradiction rewrites).
- Standing decisions from the charting session:
  - Markdown (wikilinks + frontmatter) is the source of truth for the graph; any index under `.kb/` is derived and rebuildable.
  - Code and config use provider-neutral names (`decisions`); Jev is only the default model.
  - Write policy (d): apply-band relations always go to frontmatter; kb also inserts `[[link]]` at the first literal mention when the decision model confirms its sense; review-band items only go to a review list; shadow is the default decision mode until a topic opts into apply.
  - Generation model: OpenRouter `xiaomi/mimo-v2.6-flash`.
  - Nothing a decision rejects is deleted; rejected items are quarantined.
  - Every feature that writes links, relations, or enrichment must also run over an already-ingested topic, not only at ingest time (Pedro, 2026-09-22).
  - The decision model is required (Pedro, 2026-09-24): content commands refuse to start without it; failures during a run are visible `undecided` statuses, never silent fallbacks. This replaces the earlier "everything optional" decision.
  - Pedro's pain, in priority: choose what to ingest with high precision; find the right documents later; a complete classification system between them.
- Background: `~/Dev/courses/pedronauck/research/jev` (Jev topic KB), this conversation's tweet survey, and the research-enrich session (Claude Code session 4a04893e, 2026-09-22; code in `~/Dev/courses/pedronauck/packages/research-enrich`), whose findings are attached as `## Context` on the affected tickets.
- Pedro's standing preferences: explain simply; never delete directly (mark and sample-review first); `raw/` may receive enrichment frontmatter.

## Decisions so far

<!-- one line per closed ticket -->

- [xiaomi/mimo-v2.6-flash as the generation model](tickets/02-cheap-llm.md): slug live and cheap ($0.14/$0.28 per 1M, JSON schema, reasoning must be disabled explicitly); PT quality unmeasured, keep model configurable, `mimo-v2.5` as fallback.
- [Obsidian conventions for aliases and typed relations](tickets/04-obsidian-conventions.md): `aliases` as a plain string list; one flat top-level list per relation with quoted wikilinks to file name/path (backlinks + untyped graph, Dataview reads them, Breadcrumbs needs user-configured fields); kb lint ignores frontmatter links and aliases today.
- [Decision model API contract](tickets/01-decision-model-api.md): `POST /api/alpha/decisions` with `{model, state, questions}`; choice is single-label (≤255 options, returns confidence), noul returns only p(yes), score 2–10 levels; state ≤32k; pin `typesafe/jev-1.13`; ~15 s timeout, retry 408/429/5xx, failure = undecided; `usage.cost` per call; no retention on OpenRouter.
- [Prior art for decision-model linking in LLM wikis](tickets/03-prior-art.md): code proposes candidates, the model judges; one noul per candidate beats one big choice for recall (s1m 0.83 vs 0.57–0.72); probabilities rank but aren't calibrated, so calibrate on labels and store question hash + model version; copy jev-the-janitor's write-back (namespaced frontmatter, atomic byte-preserving rewrite, quarantine, review file, `locked: true`).
- [qmd as the candidate generator](tickets/05-qmd-candidates.md): no doc-to-doc neighbors and default modes are too slow (5–29 s/call); candidates come from a pure-Go alias scan plus in-memory BM25, with qmd vector neighbors (`--no-rerank`, via `qmd mcp --http`, ~0.1–0.8 s/doc) as an optional layer after `update`+`embed`; qmd scores never feed confidence bands.
- [Decision engine seam and module design](tickets/06-decision-engine.md): deep `internal/decisions` module with one `Decide` entry that owns transport, chunking, cache, banding, budget and mode; per-purpose thresholds; mode from config → `topic.yaml` → flag; receipts in `<topic>/.decisions/receipts.jsonl` (cache + calibration); `[decisions]` + `[generation]` config; US$ 1 default budget per run; separate `internal/generation`. Also fixed write policy (d).
- [Frontmatter schema for decisions, vocabulary and relations](tickets/07-frontmatter-schema.md): flat unprefixed keys (revision 3 dropped `ai_`; `triage`/`genre` avoid the vault's `status`/`kind`) + `aliases` + `locked`; hash/contract/bank in `.decisions/state.jsonl`; owned-key list + value-hash check for collisions; byte-preserving writes; lint reads frontmatter links like Obsidian.
- [Topic scope and vocabulary lifecycle](tickets/08-topic-vocabulary.md): user-accepted five-part selection contract in `topic.yaml`; articles are the concept vocabulary; generated aliases; concept proposals via review.
- [Incremental linking pipeline](tickets/09-linking-pipeline.md): code candidates (≤20) → one Jev request per document → write policy (d); reverse pass; `kb link` over existing topics.
- [Ingest gates: relevance, quality and duplicates](tickets/10-ingest-gates.md): dedupe → pre-fetch relevance (bulk) → quality + relevance → near-duplicate → write; quarantine, never delete.
- [OKF type suggestion and semantic lint](tickets/11-okf-and-lint.md): optional `--type` with suggestion; lint never calls a model, reads stored decisions.
- [Calibration, evaluation and test strategy](tickets/12-calibration.md): review verdicts become labels; `kb review calibrate`; fixtures + jevlive regression + integration flows.
- [Write the kb-jev spec](tickets/13-write-spec.md): [spec.md](spec.md); also adds retrieval (`kb find`) and facets, which Pedro brought into scope.
- [Revision 2 of the spec](spec.md#revision-2-what-changed-and-why): checked against the research-vault run (session 4a04893e, 2026-09-22/24); adds impact preview, contract import and drafting rules, provenance in state, `criterion`, vocabulary bootstrap, refetch before quarantine, quarantine ledger, recapture/remove queues, label import with holdout, shadow-by-default relevance gates at 0.8, `entities`/`questions`.
- [Revision 3 of the spec](spec.md#revision-3-no-ai_-prefix): Pedro vetoed the `ai_` prefix; plain keys (`triage`, `genre`, `summary`, `concepts`, `related`...), bookkeeping in `.decisions/state.jsonl`, collisions via owned-key list + value-hash check (§6).

## Not yet specified

Nothing: every earlier fog patch was either decided in the spec (budgets §4.1, backfill §15, skills and templates §16, concurrent edits §6, privacy §14, rescue of skipped items §7) or ruled out of scope below.

## Out of scope

- Choosing captions vs STT with the decision model: no public evidence.
- An external graph store (Neo4j) or the Obsidian API: markdown is the graph; Obsidian only renders it.
- Cross-topic retrieval and entity normalization: find, link and vocabulary are topic-scoped in this effort.
- A topic remove/rename command that repairs links elsewhere.
- Implementing the spec: this map ends at the spec.
