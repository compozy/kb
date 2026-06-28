---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/topic/topic.go
line: 440
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHTB,comment:PRRC_kwDOR-Fawc7P2i82
---

# Issue 013: _🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Keep the standard topic directories for OKF topics.**

Short-circuiting here means `kb topic new --mode okf` no longer creates `raw/`, `wiki/`, `outputs/`, or `bases/`. That breaks the topic skeleton contract and leaves `raw/codebase/` missing until some later command recreates it. Keep the normal scaffold and layer the OKF-specific files on top.
    
As per coding guidelines, `kb topic new` owns the topic skeleton under the vault root, including `raw/`, `wiki/`, `outputs/`, `bases/`, `CLAUDE.md`, `AGENTS.md`, and `log.md`. Based on learnings, codebase inspection commands operate on `raw/codebase/` beneath the resolved topic.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/topic/topic.go` around lines 437 - 440, The createTopicSkeleton
helper is incorrectly short-circuiting for TopicModeOKF, which skips the
standard topic scaffold. Update createTopicSkeleton so it always creates the
normal topic directories and base files owned by kb topic new, then apply any
OKF-specific additions on top; keep the scaffold generation for raw, wiki,
outputs, bases, CLAUDE.md, AGENTS.md, and log.md, and ensure raw/codebase
remains present for later inspection commands.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:bcdc5e0433253d84d0d11502 -->

_Sources: Coding guidelines, Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `createTopicSkeleton` returns after creating only the OKF root directory for OKF topics, skipping the standard `raw/`, `wiki/`, `outputs/`, and `bases/` directories owned by `kb topic new`. Fix by always creating the standard skeleton, then layering OKF-specific files/templates on top.
