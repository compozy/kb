---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/okf/okf.go
line: 580
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHS7,comment:PRRC_kwDOR-Fawc7P2i8s
---

# Issue 010: _🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Return parse errors instead of dropping concepts from the regenerated index.**

A malformed concept frontmatter currently makes `loadConcepts` silently skip that file, so `RegenerateIndex` can rewrite `index.md` without existing concepts. Preserve the error and only skip files that truly have no frontmatter.

<details>
<summary>Proposed fix</summary>

```diff
 		values, _, err := frontmatter.Parse(string(content))
-		if err != nil || len(values) == 0 {
+		if err != nil {
+			return fmt.Errorf("parse concept %s: %w", relativePath, err)
+		}
+		if len(values) == 0 {
 			return nil
 		}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		values, _, err := frontmatter.Parse(string(content))
		if err != nil {
			return fmt.Errorf("parse concept %s: %w", relativePath, err)
		}
		if len(values) == 0 {
			return nil
		}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/okf/okf.go` around lines 577 - 580, In loadConcepts, the
frontmatter.Parse error is currently being swallowed by returning nil, which
lets malformed concept files disappear from the regenerated index. Update the
logic around frontmatter.Parse so it preserves and returns parse errors, and
only skips files when there is truly no frontmatter (for example, an empty
values result without an error). Make sure the handling in loadConcepts still
supports RegenerateIndex but does not silently drop concepts from malformed
files.
```

</details>

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:bb73ea14d5a07e58bd7ee4db -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `loadConcepts` currently treats `frontmatter.Parse` errors the same as absent frontmatter and silently skips the malformed file. `RegenerateIndex` can then drop an existing concept. Fix by returning parse errors and only skipping files that parse successfully with no frontmatter values.
