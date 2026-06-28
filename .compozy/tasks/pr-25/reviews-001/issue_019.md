---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/vault/render_wiki.go
line: 117
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHTJ,comment:PRRC_kwDOR-Fawc7P2i89
---

# Issue 019: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Use each rendered document’s own directory when building OKF links.**

All of these calls pass `""` as `fromDir`, so OKF-mode links in `wiki/codebase/index/*.md` and `wiki/codebase/concepts/*.md` are rendered as topic-rooted paths like `wiki/codebase/...`. GitHub resolves markdown targets relative to the current file, so those links break; e.g. dashboard → concept should be `../concepts/...`, and concept → concept links should usually be sibling-relative. Compute `fromDir` from the document being rendered before calling `linkFor`.






Also applies to: 147-147, 194-194, 314-316, 395-396, 461-462, 538-538, 616-616

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/vault/render_wiki.go` around lines 108 - 117, The OKF link
generation in the wiki renderer is using an empty fromDir, so links in rendered
markdown are rooted at the topic instead of relative to the current document.
Update the linkFor calls in the wiki rendering flow to pass the rendered
document’s own directory as fromDir before building links, especially for the
codebase index and concept pages. Use the current document context from the
render function that emits wiki pages so dashboard/index-to-concept and
concept-to-concept links resolve with correct relative paths. Apply the same
fromDir fix to every linkFor invocation in this render path, including the
navigation and article link sections.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1d183f40075d3c882bcdf343 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: Wiki renderers call `linkFor` with an empty source directory while rendering files under `wiki/codebase/index/` and `wiki/codebase/concepts/`. OKF markdown links therefore resolve from the bundle root. Fix by computing each rendered document's directory and passing it into link generation.
