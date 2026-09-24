---
id: 8
title: Topic scope and vocabulary lifecycle
labels: [wayfinder:grilling]
status: closed
assignee: "claude"
blocked_by: [2, 7]
---

## Question

Where does a topic's scope and vocabulary come from, and how does it grow?

The decision model is literal: relevance and tagging only work against a written scope and a written vocabulary. Decide: where scope lives (topic.yaml, CLAUDE.md, a new file), what the tag vocabulary is built from (article titles, aliases, existing tags), when the generation model proposes aliases and new tags, and how the user approves them.

## Context

From session 4a04893e (2026-09-22, `~/Dev/courses/pedronauck/packages/research-enrich`):

- A one-line scope failed (49/160 flagged off-topic, only 30–40% truly junk). What worked: a `## Selection contract` section in each topic's CLAUDE.md, in English, with Purpose / Core / Adjacent (keep) / Collected on purpose / Out of scope.
- Agent-written contracts are unstable (precision swung up to 30 points between rewrites; internal contradictions make Jev follow whichever line matches best), so the owner must approve the contract and something should detect contradictions in it.
- Extracted entities had duplicates (`MCP` vs `Model Context Protocol`) and noise; aliases/canonicalization is needed.

## Resolution

Decided by the agent on 2026-09-24 under Pedro's delegation ("the rest you decide on your own to create a spec"), grounded in the jev-engineering and typesafe-ai skills and the research on this map. Five-part selection contract in `topic.yaml` (purpose/core/adjacent/collected_on_purpose/out_of_scope), drafted by the generation model but active only after user `--accept`, with a contradiction self-check; concept vocabulary = wiki articles (title + aliases + summary); aliases generated and validated; new concepts only as review proposals. Detail: [spec §5](../spec.md).
