---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/okf/okf_test.go
line: 224
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHS2,comment:PRRC_kwDOR-Fawc7P2i8o
---

# Issue 012: _🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Assert diagnostic messages in `assertIssue`.**

Matching only severity/file/target lets the tests pass when the checker reports the right field for the wrong reason. Add an expected message substring so these tests fail on behavior changes, not just missing entries. As per path instructions, tests must have specific error assertions and verify behavior outcomes.

<details>
<summary>Proposed helper shape</summary>

```diff
-func assertIssue(t *testing.T, issues []models.LintIssue, severity models.DiagnosticSeverity, filePath, target string) {
+func assertIssue(
+	t *testing.T,
+	issues []models.LintIssue,
+	severity models.DiagnosticSeverity,
+	filePath, target, messageContains string,
+) {
 	t.Helper()
 	for _, issue := range issues {
-		if issue.Severity == severity && issue.FilePath == filePath && issue.Target == target {
+		if issue.Severity == severity &&
+			issue.FilePath == filePath &&
+			issue.Target == target &&
+			strings.Contains(issue.Message, messageContains) {
 			return
 		}
 	}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/okf/okf_test.go` around lines 216 - 224, The test helper assertIssue
currently matches only severity, file path, and target, so it can miss
regressions in the actual diagnostic text. Update assertIssue in okf_test.go to
also accept an expected message substring and verify it against each
models.LintIssue’s message/diagnostic text before returning. Then update the
callers to pass the expected message so the tests assert the exact behavior
using the existing assertIssue helper.
```

</details>

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:868612d702d0a12b68003f1e -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `assertIssue` matches severity, file, and target but not the diagnostic message. A checker regression could report the right target for the wrong reason. Fix by requiring an expected message substring and updating callers.
