---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/mediadl/ytdlp.go
line: 727
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7mR,comment:PRRC_kwDOR-Fawc7Pt_e1
---

# Issue 018: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Default extraction can reject real captions when yt-dlp omits `language`.**

With no preferred languages, this path forces `orig`, but candidates only become `Original` when `info.Language` or a `-orig` key identifies them. A payload with manual subtitles like `{"subtitles":{"en":[...]}}` and an empty `language` field now falls through to `transcript_unavailable` even though a usable caption exists. The updated happy-path fixtures in `internal/mediadl/ytdlp_test.go` adding `"language":"en"` / `"pt"` are a strong sign this path is now metadata-dependent.

   


Also applies to: 730-765, 767-805

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/mediadl/ytdlp.go` around lines 707 - 727, The default-caption
selection in selectYTDLPCaption is too dependent on info.Language, so manual
subtitles can be missed when yt-dlp omits language metadata. Update the matching
logic in ytDLPCaptionCandidates / matchingYTDLPCaptionCandidates /
ytDLPCaptionLess so a no-preferences request still accepts usable manual
subtitle tracks even when info.Language is empty, rather than requiring an
orig/original classification. Keep the fallback behavior in selectYTDLPCaption
but ensure it can pick a real caption candidate from subtitles-only payloads
like those covered by the ytdlp tests.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:67f55926dd92c41b96415408 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: Default caption preferences normalize to `orig`, but manual subtitle tracks are only considered original when yt-dlp supplies `info.language` or a `-orig` key. I will allow usable manual subtitles to satisfy the default original-language preference when yt-dlp gives no original-language metadata, and add a regression.
