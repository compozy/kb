---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/vault/render.go
line: 191
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHTK,comment:PRRC_kwDOR-Fawc7P2i8-
---

# Issue 017: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Pass the current document directory into `linkFor` for OKF mode.**

These call sites still use `fromDir == ""`, so nested raw documents render links as if they were at the topic root. In OKF mode that breaks navigation on GitHub—for example, a link inside `raw/codebase/files/commands/run.ts.md` to `raw/codebase/symbols/main--commands-run-ts-l1.md` needs a relative target like `../../symbols/main--commands-run-ts-l1.md`, not `raw/codebase/symbols/...`. Thread the source document’s directory through `linkForNode` instead of hard-coding the root. 






Also applies to: 224-224

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/vault/render.go` around lines 190 - 191, Update the OKF link
generation path so `linkFor` receives the current document directory instead of
an empty root path. Thread the source directory through `linkForNode` and any
related callers that currently pass fromDir as "", so
`LinkFormatterFor(...).Link(...)` can build correct relative links for nested
raw documents. Focus on the `linkFor` and `linkForNode` flow in `render.go` and
preserve existing behavior outside OKF mode.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7d09aebfa51c588d385156b7 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: Raw document relation helpers pass `fromDir == ""` into `linkFor`, so OKF links in nested raw documents are generated from the bundle root. Fix by threading the current rendered document directory into relation and source-link helpers.
