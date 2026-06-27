---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/cli/ingest_test.go
line: 1199
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:1e91ed1ed0dd
review_hash: 1e91ed1ed0dd
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 005: Check the specific rejection reason here.
## Review Comment

This only asserts that some error came back, so an unrelated failure would still satisfy the test. Please assert the invalid-URL diagnostic from `ingest channel` as well. As per path instructions, `MUST have specific error assertions (ErrorContains, ErrorAs)`.

<!-- cr-comment:v1:e357a4c1db5d06ad2cb09293 -->

_Source: Path instructions_

---

## Triage

- Decision: `valid`
- Notes: `TestIngestChannelCommandRejectsVideoURL` only checks for any error. I will assert the specific single-video rejection diagnostic so unrelated failures cannot satisfy the test.
