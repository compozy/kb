---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/instagram/instagram_test.go
line: 160
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7mH,comment:PRRC_kwDOR-Fawc7Pt_er
---

# Issue 013: _📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Table-drive the extractor scenarios with `Should...` subtests.**

These cases all follow the same arrange/act/assert shape, but they are split into standalone tests instead of one table-driven subtest loop with the repo’s required naming pattern. As per coding guidelines, `**/*_test.go`: Default to table-driven tests with focused helpers and t.TempDir() for filesystem isolation. As per path instructions, `MUST use t.Run("Should...") pattern for ALL test cases`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/instagram/instagram_test.go` around lines 74 - 160, The extractor
tests are duplicated standalone cases and should be consolidated into a
table-driven test using the repo’s required t.Run("Should...") naming pattern.
Refactor the scenarios in TestExtractComposesCaptionAndTranscript,
TestExtractTranscriptOnlyWhenNoCaption,
TestExtractCaptionOnlyFallbackWhenTranscriptUnavailable,
TestExtractPropagatesErrorWhenNoCaptionAvailable, and
TestExtractRejectsInvalidURL into one table with focused arrange/act/assert
helpers, keeping the existing Extract and stubCore behavior covered under
separate Should... subtests.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:78183882179901ccf4ed94a6 -->

_Sources: Coding guidelines, Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The Instagram extractor scenarios share the same arrange/act/assert shape but are split across standalone tests. I will consolidate them into a table-driven `Should...` subtest matrix with focused assertions.
