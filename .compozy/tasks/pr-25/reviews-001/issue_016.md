---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/vault/pathutils_test.go
line: 199
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4586821980,nitpick_hash:4f089a4d674b
review_hash: 4f089a4d674b
source_review_id: "4586821980"
source_review_submitted_at: "2026-06-28T02:04:32Z"
---

# Issue 016: Use t.Run("Should...") consistently in the new test cases.
## Review Comment

This block mixes ad-hoc assertions with subtests, and the subtest names do not follow the required `Should...` pattern. Split the wiki/OKF/selector checks into named subtests so failures stay consistent with the repo’s test contract. As per path instructions, `**/*_test.go` “MUST use t.Run("Should...") pattern for ALL test cases.”

<!-- cr-comment:v1:8850e596056339e4bad336c4 -->

_Source: Path instructions_

## Triage

- Decision: `valid`
- Notes: The new link formatter test cases use names like `same directory` and include some direct assertions outside behavior-named subtests. Fix by using `t.Run("Should...")` names consistently for the formatter scenarios and related selector checks.
