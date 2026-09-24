---
id: 9
title: Incremental linking pipeline
labels: [wayfinder:grilling]
status: closed
assignee: "claude"
blocked_by: [3, 5, 6, 7, 8]
---

## Question

How does kb link documents incrementally?

Decide: candidate sources and their merge (alias dictionary, qmd neighbors, shared sources/tags); the question set per document (mention sense, relation type, affects); body-insertion rules within write policy (d), already decided in [Decision engine seam and module design](06-decision-engine.md) (first literal mention, skip code/headings/existing links, one link per target, alias syntax for case/plural); the reverse pass when an article or alias is born; trigger points (end of `kb ingest`, and `kb link` over an already-ingested topic, which is required); the review queue and `kb link --review` UX; and when a `contradicts` decision escalates to the generation model.

## Context

From session 4a04893e (2026-09-22, `~/Dev/courses/pedronauck/packages/research-enrich`):

- For a closed concept vocabulary, a relative `choice` for the primary concept (≥0.3, with `none`) plus a `noul` per concept (≥0.7) beat noul-only (34/103 docs got no concept with noul-only; v3 left 1, avg 1.8 concepts/doc). Weigh this against the prior-art finding that noul-per-candidate wins recall for links.
- Shared extracted entities are another candidate source.

## Resolution

Decided by the agent on 2026-09-24 under Pedro's delegation ("the rest you decide on your own to create a spec"), grounded in the jev-engineering and typesafe-ai skills and the research on this map. Candidates from mention scan + BM25 + shared concepts/sources + optional qmd vectors (cap 20); one request per document with a noul per candidate, relation choice, mention-sense nouls and `affects`; write policy (d); contradicts always review; reverse pass on new article/alias; `kb link` incremental over existing topics. Detail: [spec §9](../spec.md).
