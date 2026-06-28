---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/cli/okf_integration_test.go
line: 87
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHSX,comment:PRRC_kwDOR-Fawc7P2i8H
---

# Issue 002: _🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**The promotion log assertion does not validate the promoted entry.**

`frontmatter.DateLayout[:4]` is the literal `"2006"`, not the current year, and the `**Creation**` fallback is already present on a freshly scaffolded topic. A missing `InsertLogEntry` call would still pass here; assert a promotion-specific marker such as the promoted title or output path instead.

As per path instructions, "`**/*_test.go`: Ensure tests verify behavior outcomes, not just function calls`."

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/okf_integration_test.go` around lines 80 - 87, The promotion log
check in the test is too weak because it can pass without verifying that
InsertLogEntry actually recorded the promotion. Update the assertion in
okf_integration_test.go to look for a promotion-specific outcome in log.md,
using unique symbols like InsertLogEntry, frontmatter.DateLayout, and the
promoted topic/asset name, rather than the literal year prefix or the generic
**Creation** entry. Ensure the test fails unless the promotion entry itself is
present.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:dca60aa4fc7207a5540cf599 -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The integration test reads `log.md` but only checks the literal layout year prefix or the scaffold creation entry. Both can exist without a promotion log entry. Fix by asserting promotion-specific content such as the promoted title/path/source in the log.
