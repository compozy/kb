---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/okf/official_integration_test.go
line: 16
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHSx,comment:PRRC_kwDOR-Fawc7P2i8k
---

# Issue 008: _📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Use `Should...` subtest names for official bundle cases.**

The subtest names are currently just bundle IDs. Rename them to the required behavior-oriented pattern. As per path instructions, `MUST use t.Run("Should...") pattern for ALL test cases`.

<details>
<summary>Proposed fix</summary>

```diff
-		t.Run(bundle, func(t *testing.T) {
+		t.Run("Should pass lenient conformance for "+bundle, func(t *testing.T) {
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	for _, bundle := range []string{"ga4", "stackoverflow", "crypto_bitcoin"} {
		bundle := bundle
		t.Run("Should pass lenient conformance for "+bundle, func(t *testing.T) {
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/okf/official_integration_test.go` around lines 14 - 16, The official
integration subtests are using raw bundle IDs as names instead of the required
behavior-oriented `Should...` pattern. Update the `t.Run` call in
`official_integration_test.go` so each case name follows the `Should...`
convention while still iterating over the same bundle list, keeping the test
logic in `t.Run` and the surrounding loop unchanged.
```

</details>

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a3d3e1c78f85240461833e5b -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: Official bundle integration subtests use raw bundle IDs as test names. The repo convention requires behavior-oriented `Should...` subtest names. Fix by renaming the subtests while keeping the same cases and assertions.
