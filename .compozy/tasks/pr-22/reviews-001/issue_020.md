---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/youtube/channel_test.go
line: 15
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:239a97ef291b
review_hash: 239a97ef291b
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 020: Align these scenarios with the repo’s required subtest shape.
## Review Comment

Most of the new coverage is expressed as standalone top-level tests, and the table-driven case does not use `Should...` subtest names. Please fold these scenarios into `t.Run("Should...")` cases so the file follows the repo’s required test structure consistently. As per coding guidelines, `**/*_test.go`: `Default to table-driven tests with focused helpers and t.TempDir() for filesystem isolation`; and as per path instructions, `MUST use t.Run("Should...") pattern for ALL test cases`.

Also applies to: 59-251

<!-- cr-comment:v1:cf0eef2c1834f613ff2e88c4 -->

_Sources: Coding guidelines, Path instructions_

## Triage

- Decision: `valid`
- Notes: `internal/youtube/channel_test.go` has standalone test cases and raw subtest names rather than `Should...` scenarios. I will fold the channel behavior checks into `Should...` subtests while keeping the existing assertions.
