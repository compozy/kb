---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/vault/render_test.go
line: 290
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHTF,comment:PRRC_kwDOR-Fawc7P2i85
---

# Issue 018: _🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Assert the real relative targets for nested OKF documents.**

These expectations currently lock in the broken root-scoped links. From `raw/codebase/files/commands/run.ts.md`, the symbol link should resolve relative to that directory (for this fixture, `../../symbols/main--commands-run-ts-l1.md`), and from `wiki/codebase/index/Codebase Dashboard.md` the concept link should be `../concepts/Codebase%20Overview.md`. As per path instructions, tests should verify behavior outcomes and not be weakened to fit broken behavior.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/vault/render_test.go` around lines 266 - 290, The OKF link
assertions in TestRenderDocumentsUseOKFMarkdownLinkSyntax are currently
expecting root-scoped paths instead of the true relative targets. Update the
test to verify the actual nested relative links produced by
vault.RenderDocuments for the raw/codebase/files/commands/run.ts.md document and
the wiki/codebase/index/Codebase Dashboard.md document, using the symbol and
concept link outputs from the fixture rather than the broken root-relative
paths. Keep the checks focused on the rendered document bodies and ensure they
assert the correct relative markdown destinations.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:962b78f82ef77b6e60e8c4ce -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `TestRenderDocumentsUseOKFMarkdownLinkSyntax` currently asserts root-scoped OKF paths, which locks in the broken behavior. Fix the production link generation and update the test to assert true document-relative markdown targets.
