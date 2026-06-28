---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/vault/writer.go
line: 477
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHTN,comment:PRRC_kwDOR-Fawc7P2i9E
---

# Issue 020: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Pass `wiki/index` as the source directory for these OKF bridge links.**

This block is rendered into files under `wiki/index/`, but `linkFor(topic, "", ...)` makes the OKF formatter resolve from the bundle root. That produces `wiki/codebase/index/...` instead of `../codebase/index/...`, so the bridge navigation breaks in OKF bundles.
   

<details>
<summary>Proposed fix</summary>

```diff
-		"- " + linkFor(topic, "", GetWikiIndexPath(CodebaseDashboardTitle), CodebaseDashboardTitle),
-		"- " + linkFor(topic, "", GetWikiIndexPath(CodebaseConceptIndexTitle), CodebaseConceptIndexTitle),
-		"- " + linkFor(topic, "", GetWikiIndexPath(CodebaseSourceIndexTitle), CodebaseSourceIndexTitle),
+		"- " + linkFor(topic, "wiki/index", GetWikiIndexPath(CodebaseDashboardTitle), CodebaseDashboardTitle),
+		"- " + linkFor(topic, "wiki/index", GetWikiIndexPath(CodebaseConceptIndexTitle), CodebaseConceptIndexTitle),
+		"- " + linkFor(topic, "wiki/index", GetWikiIndexPath(CodebaseSourceIndexTitle), CodebaseSourceIndexTitle),
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		"- " + linkFor(topic, "wiki/index", GetWikiIndexPath(CodebaseDashboardTitle), CodebaseDashboardTitle),
		"- " + linkFor(topic, "wiki/index", GetWikiIndexPath(CodebaseConceptIndexTitle), CodebaseConceptIndexTitle),
		"- " + linkFor(topic, "wiki/index", GetWikiIndexPath(CodebaseSourceIndexTitle), CodebaseSourceIndexTitle),
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/vault/writer.go` around lines 475 - 477, The OKF bridge links in
writer.go are being resolved from the bundle root because linkFor is called with
an empty source directory, which breaks navigation for files rendered under
wiki/index/. Update the three linkFor calls in this block to pass wiki/index as
the source directory so GetWikiIndexPath(CodebaseDashboardTitle),
GetWikiIndexPath(CodebaseConceptIndexTitle), and
GetWikiIndexPath(CodebaseSourceIndexTitle) resolve correctly in OKF bundles.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e8adea6f301827daae62ae35 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: Topic index bridge links are rendered into `wiki/index/` files but call `linkFor` with an empty source directory. OKF links point at `wiki/codebase/...` instead of `../codebase/...`. Fix by passing `wiki/index` for these bridge links. Regression coverage is in `internal/vault/writer_test.go` because the bridge helper is unexported and `WriteVault` is the observable behavior.
