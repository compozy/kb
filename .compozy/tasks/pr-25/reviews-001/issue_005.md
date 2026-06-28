---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/cli/topic_test.go
line: 140
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHSu,comment:PRRC_kwDOR-Fawc7P2i8i
---

# Issue 005: _🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Tighten the invalid-mode assertion.**

`strings.Contains(err.Error(), "invalid --mode")` can still pass on the wrong failure path. Compare the full message or use the shared error helper so this only passes on the intended validation error.

As per path instructions, "`**/*_test.go`: MUST have specific error assertions (ErrorContains, ErrorAs)`."

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/topic_test.go` around lines 127 - 140, The invalid-mode test in
TestTopicNewCommandRejectsInvalidMode is too loose because
strings.Contains(err.Error(), "invalid --mode") can match unrelated failures;
update the assertion to use the shared error-checking helper or compare the full
validation message from newRootCommand()/command.ExecuteContext so only the
intended --mode validation error passes. Keep the test anchored on the specific
invalid-mode path by asserting against the exact error content rather than a
substring.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:b0d03294055de1a0067118fe -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The invalid-mode assertion only checks `strings.Contains(err.Error(), "invalid --mode")`, which could pass for the wrong validation path. Fix by asserting the complete expected validation message.
