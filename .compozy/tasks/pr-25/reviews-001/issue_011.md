---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/okf/okf_test.go
line: 15
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4586821980,nitpick_hash:799c4d0503e6
review_hash: 799c4d0503e6
source_review_id: "4586821980"
source_review_submitted_at: "2026-06-28T02:04:32Z"
---

# Issue 011: Wrap test cases in Should... subtests.
## Review Comment

These scenarios are direct top-level test bodies; the repo test instructions require `t.Run("Should...")` for all test cases. The repeated strict/lenient check cases are also good candidates for a small table. As per path instructions, `MUST use t.Run("Should...") pattern for ALL test cases`; as per coding guidelines, default to table-driven tests.

<!-- cr-comment:v1:aed6da973f24c45c676f882e -->

_Sources: Coding guidelines, Path instructions_

## Triage

- Decision: `valid`
- Notes: New OKF tests include direct top-level scenarios and repeated strict/lenient checks that are not wrapped in `t.Run("Should...")` cases. Fix by adding behavior-named subtests and keeping assertions focused on real OKF behavior.
