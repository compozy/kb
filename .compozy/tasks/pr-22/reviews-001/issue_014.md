---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/logger/logger_test.go
line: 31
severity: minor
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:a044bc5db4c1
review_hash: a044bc5db4c1
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 014: Assert the configured level instead of only successful construction.
## Review Comment

`wantLevel` is never checked, so every non-error row passes even if `New` ignores `tc.level` and always returns the same logger. Emitting a record per case and asserting the `"level"` field would make this table cover the behavior it names; while touching it, please also switch the subtest names to the required `Should...` form. As per path instructions, `MUST test meaningful business logic, not trivial operations`, `Ensure tests verify behavior outcomes, not just function calls`, and `MUST use t.Run("Should...") pattern for ALL test cases`.

<!-- cr-comment:v1:78b82ef844c14c87b049c045 -->

_Source: Path instructions_

## Triage

- Decision: `valid`
- Notes: `TestNew` records `wantLevel` but never verifies the emitted log level, so a logger ignoring the requested level would pass. I will emit records per case, assert enabled/disabled level behavior, and rename the subtests to `Should...`.
