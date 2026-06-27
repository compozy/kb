---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/cli/ingest_channel.go
line: 69
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:cf97d49eabb7
review_hash: cf97d49eabb7
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 002: Move the bulk-ingest workflow out of the Cobra handler.
## Review Comment

This `RunE` now owns normalization, resume filtering, bulk extraction orchestration, per-video ingest, and summary mutation. That makes the command itself the workflow, which is harder to reuse and test than a package-level service/helper. As per coding guidelines, `{cmd/kb/**/*.go,internal/cli/**/*.go}`: `Keep Cobra commands thin and push behavior into internal packages, particularly in internal/cli`.

<!-- cr-comment:v1:90455b1bc334cf6531bd844a -->

_Source: Coding guidelines_

## Triage

- Decision: `valid`
- Notes: `newIngestChannelCommand` currently contains the channel normalization, listing, resume filtering, bulk extraction, per-video ingest, and summary handling directly inside `RunE`. I will extract the workflow into package-level helpers so the Cobra handler only resolves flags/config and delegates execution.
