---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/mediadl/mediadl_test.go
line: 13
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:0390e1fe786f
review_hash: 0390e1fe786f
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 015: Restructure these helper tests into table-driven subtests.
## Review Comment

This file is mostly case matrices (`ParseTranscriptionPolicy`, language matching, status classification), but every scenario is still a standalone top-level test, and the invalid-policy path only checks `err != nil`. Please collapse these into tables behind `t.Run("Should...")` and assert the expected error text for negative cases so failures stay behavior-specific.

As per coding guidelines, `**/*_test.go`: `Default to table-driven tests with focused helpers and t.TempDir() for filesystem isolation`; as per path instructions, `**/*_test.go`: `MUST use t.Run("Should...") pattern for ALL test cases` and `MUST have specific error assertions (ErrorContains, ErrorAs)`.

<!-- cr-comment:v1:8b7e95db9f4f0f54c3511180 -->

_Sources: Coding guidelines, Path instructions_

## Triage

- Decision: `valid`
- Notes: The mediadl helper tests are standalone assertions and the invalid transcription policy case only checks for any error. I will table-drive the cases under `Should...` subtests and assert the expected invalid-policy error text.
