---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/instagram/instagram_test.go
line: 149
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7mI,comment:PRRC_kwDOR-Fawc7Pt_es
---

# Issue 012: _🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Assert the specific error kind in the negative-path tests.**

Both checks currently pass on any failure, including the wrong one. Use `errors.As` plus the expected `mediadl.ErrorKind...` so the tests pin the intended behavior. As per path instructions, `MUST have specific error assertions (ErrorContains, ErrorAs)`.

   


Also applies to: 157-158

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/instagram/instagram_test.go` around lines 148 - 149, The
negative-path checks in instagram_test should assert the exact error type
instead of only checking that an error occurred, so the tests fail on the wrong
error. Update the Extract-based assertions to use errors.As against the expected
mediadl error type and verify the specific mediadl.ErrorKind value in both
failing cases, including the other check at the additional referenced spot, so
the behavior is pinned to the intended transcript/caption absence path.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c0c1e8c107630909fa5bcf49 -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The Instagram negative-path tests only assert that an error occurred. I will assert the concrete `mediadl.ErrorKind` with `errors.As` for the transcript-unavailable and invalid-URL paths.
