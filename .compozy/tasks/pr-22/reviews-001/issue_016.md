---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/mediadl/transcription.go
line: 123
severity: major
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:ee127a088db0
review_hash: ee127a088db0
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 016: Forced STT can currently succeed without a transcriber.
## Review Comment

The early return on Line 123 makes the new nil-transcriber branch on Lines 129-130 unreachable when `result.TranscriptionPolicy == TranscriptionPolicySTT` and `extractor.stt == nil`, so `Extract(..., {TranscriptionPolicy: STT})` can return `(result, nil)` with no transcript at all.

[suggested fix]

<!-- cr-comment:v1:737ad3a64e7f6babba4c8921 -->

## Triage

- Decision: `valid`
- Notes: With `TranscriptionPolicySTT` and no transcriber, `extractSTTFromYTDLP` returns `(result, nil)` because `shouldAttemptSTT` is false and `transcriptErr` is nil. I will return an explicit STT configuration error for forced STT without a transcriber and cover it with a regression.
