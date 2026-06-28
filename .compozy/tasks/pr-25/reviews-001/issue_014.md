---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/topic/topic.go
line: 552
severity: major
author: coderabbitai[bot]
provider_ref: review:4586821980,nitpick_hash:062b4b1f33ec
review_hash: 062b4b1f33ec
source_review_id: "4586821980"
source_review_submitted_at: "2026-06-28T02:04:32Z"
---

# Issue 014: Preserve mode when rewriting topic.yaml.
## Review Comment

`WriteMetadataFile` always passes `""`, so the YAML omits `mode`. `readTopicMetadata` then falls back to wiki, which will silently downgrade an OKF topic the next time a caller rewrites metadata through this helper.

Also applies to: 560-566

<!-- cr-comment:v1:7bc4a93f80bde71e91e1364d -->

## Triage

- Decision: `valid`
- Notes: `topic.WriteMetadataFile` rewrites `topic.yaml` with an empty mode, so any caller that rewrites an OKF topic metadata file silently downgrades it to the wiki default on the next read. Fix by preserving the existing metadata mode when rewriting, falling back to wiki only when no mode exists.
