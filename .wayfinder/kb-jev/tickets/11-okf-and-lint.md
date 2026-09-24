---
id: 11
title: OKF type suggestion and semantic lint
labels: [wayfinder:grilling]
status: closed
assignee: "claude"
blocked_by: [6, 7]
---

## Question

How do decisions show up in `kb promote`, `kb okf check`, and `kb lint`?

Decide: `--type` suggestion/validation against `[okf].types`; description-fits-body check; the new lint issue kinds (semantic stale, contradiction, missing link, off-topic), their severity, and how they stay read-only. Also: lint today ignores frontmatter wikilinks and `aliases`, and resolves links Obsidian would not (see [Obsidian conventions for aliases and typed relations](04-obsidian-conventions.md)); decide whether relations join dead-link/orphan checks and how link resolution aligns with Obsidian.

## Context

From session 4a04893e (2026-09-22, `~/Dev/courses/pedronauck/packages/research-enrich`):

- A "description fits the body" flag misfired (0.94 on a correct summary): such checks need calibration before becoming lint issues.
- Deleting a topic left 285 dead wikilinks in other topics, fixed by hand.

## Resolution

Decided by the agent on 2026-09-24 under Pedro's delegation ("the rest you decide on your own to create a spec"), grounded in the jev-engineering and typesafe-ai skills and the research on this map. `kb promote --type` optional with suggestion; OKF advisory findings; lint never calls a model and gains frontmatter-dead-link, needs-compile, contradiction, pending-review, off-topic-kept, unclassified, contract kinds. Detail: [spec §11](../spec.md).
