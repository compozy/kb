---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/generate/generate.go
line: 451
severity: major
author: coderabbitai[bot]
provider_ref: review:4586821980,nitpick_hash:55582f299925
review_hash: 55582f299925
source_review_id: "4586821980"
source_review_submitted_at: "2026-06-28T02:04:32Z"
---

# Issue 007: Don’t hardcode wiki mode in generated topic metadata.
## Review Comment

`linkFor` now dispatches through `LinkFormatterFor(topic)`, so forcing `Mode: wiki` here makes downstream renders emit wiki-style links even when the resolved topic is OKF. Carry the actual topic mode into `models.TopicMetadata` instead.

<!-- cr-comment:v1:6c680e125ef7fc3b8a530621 -->

## Triage

- Decision: `valid`
- Notes: `runner.createTopicMetadata` always sets `Mode: models.TopicModeWiki`. When generation targets an existing OKF topic, downstream renderers use wiki-style links despite the topic mode. Fix by carrying the resolved topic mode into generation metadata instead of hard-coding wiki. Regression coverage belongs in `internal/generate/generate_test.go`, which is outside the listed production files but is the canonical test suite for `resolveTarget`.
