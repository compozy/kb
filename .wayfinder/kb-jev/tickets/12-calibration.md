---
id: 12
title: Calibration, evaluation and test strategy
labels: [wayfinder:grilling]
status: closed
assignee: "claude"
blocked_by: [9, 10]
---

## Question

How are confidence bands calibrated and the whole effort verified?

Decide: default band thresholds; using shadow mode on the `jev` topic against its existing wikilinks as a gold set; what shadow reports look like; test layering (fake HTTP only at the API boundary, integration tests over real vault flows); what `make verify` must cover; and a small PT/EN acceptance set for the generation model's aliases and summaries (its PT quality is unmeasured, see [xiaomi/mimo-v2.6-flash as the generation model](02-cheap-llm.md)).

## Context

From session 4a04893e (2026-09-22, `~/Dev/courses/pedronauck/packages/research-enrich`):

- Any change to the scope contract shifts precision by up to 30 points: the contract hash must be part of the cache key and receipt, and bands recalibrated after edits.
- The title-based gold set was itself wrong (5 of 6 "false positives" at p≥0.7 were label errors): gold sets need review too.
- `deepseek/deepseek-v4-flash` ($0.049/$0.098 per 1M, strict structured outputs, reasoning off) was the generator used there; include it in the PT/EN generator eval as fallback.

## Resolution

Decided by the agent on 2026-09-24 under Pedro's delegation ("the rest you decide on your own to create a spec"), grounded in the jev-engineering and typesafe-ai skills and the research on this map. Review queue whose accept/reject verdicts become labels (plus imported human wikilinks as link gold); `kb review calibrate` per purpose with a 30-label floor; measured default thresholds; fixtures, jevlive regression cases, integration flows against fake servers. Detail: [spec §12, §17](../spec.md).
