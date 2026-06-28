---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/config/config_test.go
line: 40
severity: minor
author: coderabbitai[bot]
provider_ref: review:4586821980,nitpick_hash:fc7f6cc4e9e2
review_hash: fc7f6cc4e9e2
source_review_id: "4586821980"
source_review_submitted_at: "2026-06-28T02:04:32Z"
---

# Issue 006: Wrap these new config cases in t.Run("Should...") subtests.
## Review Comment

The OKF assertions add meaningful coverage, but they extend `_test.go` cases without the required subtest structure, and the normalization scenario is a good fit for a small table-driven case. As per coding guidelines, `**/*_test.go`: "Default to table-driven tests with focused helpers and `t.TempDir()` for filesystem isolation." As per path instructions, `**/*_test.go`: "MUST use `t.Run("Should...")` pattern for ALL test cases."

Also applies to: 109-183

<!-- cr-comment:v1:46f8af9c987318fa272413f3 -->

_Sources: Coding guidelines, Path instructions_

## Triage

- Decision: `valid`
- Notes: The OKF default and normalization assertions were appended inside broad config tests rather than behavior-named subtests. Fix by wrapping those OKF-specific assertions in `t.Run("Should...")` blocks while preserving the existing config coverage.
