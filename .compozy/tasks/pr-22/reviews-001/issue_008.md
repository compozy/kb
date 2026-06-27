---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/config/config_test.go
line: 328
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7l8,comment:PRRC_kwDOR-Fawc7Pt_eZ
---

# Issue 008: _🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Assert the `youtube.caption_languages` error explicitly.**

This new case still passes on any validation failure, so it does not prove the caption-language branch is the one being exercised. Check that the returned error contains `youtube.caption_languages`.



As per path instructions, `**/*_test.go`: MUST have specific error assertions and ensure tests can fail when business logic changes.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/config/config_test.go` around lines 325 - 328, The new empty YouTube
caption languages case in Config validation is too broad because it passes on
any error; update the test to explicitly assert that the failure comes from the
youtube.caption_languages validation branch. In the Config test case around the
YouTube caption language mutation, check the returned error text for
youtube.caption_languages so the test will fail if business logic changes and a
different validation path is triggered.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:75e33a1c5ce46d1c9c9c6b24 -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The empty caption-language validation case only checks for a non-nil error. I will add a case-specific error substring assertion for `youtube.caption_languages`.
