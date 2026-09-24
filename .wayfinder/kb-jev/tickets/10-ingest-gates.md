---
id: 10
title: Ingest gates: relevance, quality and duplicates
labels: [wayfinder:grilling]
status: closed
assignee: "claude"
blocked_by: [5, 6, 8]
---

## Question

Which gates run during ingest, where, and what happens to what they reject?

Decide: relevance gate for bulk ingest (channel, bookmarks) before any transcript or STT and its KEEP/SKIM/HIDE semantics; quality gate before write (paywall, 404, thin, music-only); ID dedupe for single YouTube/Instagram; near-duplicate vs new-version detection; quarantine location and how the user recovers items; per-command flags.

## Context

From session 4a04893e (2026-09-22, `~/Dev/courses/pedronauck/packages/research-enrich`):

- Cleaning an existing corpus paid little (polymarket-quant: 3% of kept docs flagged, $0.54 for 2,101 docs); the value is at the entry gate.
- Pattern for bands: a role `choice` (`core / adjacent / general / false_hit / boilerplate / unknown`) with drop = p(false_hit)+p(boilerplate). Measured on sports-tech: P 0.70/R 0.62 @0.5, 0.88/0.37 @0.7, 0.92/0.20 @0.8.
- Jev cannot spot a duplicate from one doc alone; duplicates need code-shortlisted pairs.
- Pedro's preference: never delete directly; mark keep/review/drop and sample-review first.

## Resolution

Decided by the agent on 2026-09-24 under Pedro's delegation ("the rest you decide on your own to create a spec"), grounded in the jev-engineering and typesafe-ai skills and the research on this map. Staged pipeline on every ingest: exact dedupe for all kinds, pre-fetch relevance for bulk (skip without fetching, rescuable), post-fetch quality + relevance, near-duplicate pairs with `new_version` → `supersedes`; quarantine under `raw/_quarantine/`; `--force`, `--decisions`, `--budget`. Detail: [spec §7](../spec.md).
