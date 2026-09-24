---
id: 7
title: Frontmatter schema for decisions, vocabulary and relations
labels: [wayfinder:grilling]
status: closed
assignee: "claude"
blocked_by: [4]
---

## Question

Which frontmatter fields does this effort add, on which document kinds (source, article, OKF concept), and how does lint validate them?

Candidates on the table: vocabulary `tags`, `content_kind`, `depth`, `relevance`, `language`, `summary`, `aliases`, `affects`, and relations (`related`, `extends`, `prerequisite`, `contradicts`). Decide names, value shapes, which are written by code vs model, whether probabilities are stored, and how `mergeExtraFrontmatter`'s reserved keys change.

## Context

From session 4a04893e (2026-09-22, `~/Dev/courses/pedronauck/packages/research-enrich`):

- Namespacing: prefix `ai_*` came from a real collision (711 docs already had `summary`). Merge without YAML round-trip, refuse to write if body or any other key changed, atomic write, `ai_hash` of the body to skip docs changed since the run.
- Pedro's standing preference: `raw/` may receive enrichment frontmatter (sidecar-only was vetoed).
- Generator wrote in English because Jev performs better in EN: decide the language of `summary`.
- Low-reliability flags (e.g. `summary_unsupported` at 0.94 on a correct summary) stayed out of frontmatter, only in raw probabilities.

## Resolution

*Superseded in part by revision 3 below (no `ai_` prefix).* Decided by the agent on 2026-09-24 under Pedro's delegation ("the rest you decide on your own to create a spec"), grounded in the jev-engineering and typesafe-ai skills and the research on this map. All kb-written keys are flat top-level `ai_*` (status, gate_reason, kind, depth, relevance, concepts, summary, relation lists, affects, supersedes, hash, contract, bank) plus Obsidian's `aliases` and a user `locked`; no probabilities in frontmatter; byte-preserving atomic writes that refuse on changed body; lint scans frontmatter links and resolves like Obsidian. Detail: [spec §6](../spec.md).

### Revision 3 (2026-09-24, Pedro)

Pedro vetoed the `ai_` prefix. Keys are plain: `triage`, `triage_reason`, `genre`, `depth`, `relevance`, `concepts`, `summary`, `criterion`, `entities`, `questions`, `quality`, `related`/`extends`/`prerequisite`/`example_of`/`contradicts`, `affects`, `supersedes`, plus `aliases` and `locked`. `status` and `kind` are avoided because the research vault already uses them (735 and 1,677 documents). Body hash, contract and bank versions move to `<topic>/.decisions/state.jsonl`. Collisions are handled by the owned-key list plus a value-hash check: kb only overwrites a value it wrote itself and nobody edited; otherwise it skips and lint reports `key-conflict`. Detail: [spec §6](../spec.md).
