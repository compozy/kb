---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/adapter/java_adapter_test.go
line: 1071
severity: major
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:a773909c2f69
review_hash: a773909c2f69
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 001: Wrap these cases in t.Run("Should...") subtests.
## Review Comment

These functions still execute their assertions directly at the top level. Please wrap each scenario in a named subtest to match the repository's required test shape. As per path instructions, `**/*_test.go`: MUST use `t.Run("Should...")` pattern for ALL test cases.

Also applies to: 1103-1129

<!-- cr-comment:v1:81fc480b6369a6bfada4b14d -->

_Source: Path instructions_

## Triage

- Decision: `valid`
- Notes: The diagnostics/module/helper utility assertions in `internal/adapter/java_adapter_test.go` execute directly in top-level test functions without `t.Run("Should...")` subtests. I will wrap each scenario in named `Should...` subtests while preserving the current assertions.
