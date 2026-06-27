---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/cli/ingest_instagram_test.go
line: 29
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:17409ec17b90
review_hash: 17409ec17b90
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 003: Align these CLI tests with the repo test harness conventions.
## Review Comment

These scenarios are still one-off top-level tests, and both command invocations hard-code `/tmp/vault`. Please move them behind `t.Run("Should...")` tables/helpers and feed a per-test vault from `t.TempDir()` so later filesystem behavior changes do not couple the cases together.

As per coding guidelines, `**/*_test.go`: `Default to table-driven tests with focused helpers and t.TempDir() for filesystem isolation`; as per path instructions, `**/*_test.go`: `MUST use t.Run("Should...") pattern for ALL test cases`.

<!-- cr-comment:v1:cfc76049dedd77db6c32d38f -->

_Sources: Coding guidelines, Path instructions_

## Triage

- Decision: `valid`
- Notes: The Instagram CLI command scenarios are top-level tests and hard-code `/tmp/vault`. I will move the scenarios behind `t.Run("Should...")` subtests and use `t.TempDir()` for the vault path.
