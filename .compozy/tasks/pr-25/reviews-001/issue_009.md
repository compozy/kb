---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/okf/okf.go
line: 128
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHS5,comment:PRRC_kwDOR-Fawc7P2i8r
---

# Issue 009: _🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Make concept path allocation atomic.**

`allocateConceptPath` checks for a missing file with `os.Stat`, but `Promote` writes later with `os.WriteFile`. Two concurrent promotions for the same key can both select the same path and the last writer wins, while index/log updates race around the overwritten concept. Reserve/write the file with `O_CREATE|O_EXCL` or retry the full path-dependent transform on collision.





Also applies to: 389-398

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/okf/okf.go` around lines 94 - 128, Make concept path allocation
atomic in Promote/allocateConceptPath: the current os.Stat check can race with
the later os.WriteFile in Promot e, letting two concurrent promotions pick the
same concept path. Update allocateConceptPath to reserve the file with exclusive
creation (O_CREATE|O_EXCL) or otherwise atomically claim the path before
writing, and on collision retry the path selection and dependent transform using
conceptKey, allocateConceptPath, and os.WriteFile.
```

</details>

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:86eeac6d5e42ab419a6ac1a8 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `allocateConceptPath` uses `os.Stat` and returns a path that `Promote` writes later with `os.WriteFile`, so concurrent promotions can pick the same filename and overwrite each other. Fix by atomically reserving the path with exclusive creation and writing the generated document through that reserved file.
