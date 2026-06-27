---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/mediadl/transcription.go
line: 151
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7mO,comment:PRRC_kwDOR-Fawc7Pt_ey
---

# Issue 017: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Do not classify every downstream STT failure as `transcript_unavailable`.**

Once this path is entered from caption fallback, `transcriptErr` already carries `ErrorKindTranscriptUnavailable`. Joining it into audio-download and generic STT errors means `mediadl.IsTranscriptUnavailable(err)` still returns true, so `internal/instagram/instagram.go` will convert real backend/config/network failures into a silent caption-only success.

<details>
<summary>Diff</summary>

```diff
 	audio, audioErr := extractor.ytDLP.downloadAudio(ctx, parsed.CanonicalURL, sttConfig.AudioFormat)
 	if audioErr != nil {
-		return result, errors.Join(transcriptErr, fmt.Errorf("media stt: %w", audioErr))
+		return result, fmt.Errorf("media stt: %w", audioErr)
 	}
@@
 		if isEmptyTranscriptionError(sttErr) {
 			return result, errors.Join(transcriptErr, &Error{
 				Kind:    ErrorKindTranscriptUnavailable,
@@
 			})
 		}
-		return result, errors.Join(transcriptErr, fmt.Errorf("media stt: %w", sttErr))
+		return result, fmt.Errorf("media stt: %w", sttErr)
 	}
```
</details>

As per path instructions, `internal/instagram/**/*.go`: `caption-plus-transcript body composition, and caption-only fallback when audio yields no transcript`.

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		return result, fmt.Errorf("media stt: %w", audioErr)
	}
	defer audio.Cleanup()

	markdown, sttErr := extractor.transcribeAudioPath(ctx, audio.Path, audio.Format, result.Metadata.Duration)
	if sttErr != nil {
		if isEmptyTranscriptionError(sttErr) {
			return result, errors.Join(transcriptErr, &Error{
				Kind:    ErrorKindTranscriptUnavailable,
				URL:     parsed.CanonicalURL,
				VideoID: parsed.VideoID,
				Message: "speech-to-text produced no transcript",
				Err:     sttErr,
			})
		}
		return result, fmt.Errorf("media stt: %w", sttErr)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/mediadl/transcription.go` around lines 136 - 151, The fallback path
in transcribeAudioPath currently preserves transcriptErr when joining
audio-download or generic STT failures, which causes
mediadl.IsTranscriptUnavailable to misclassify real backend/config/network
errors. Update the transcript assembly in transcription.go so only the
empty-transcript case carries ErrorKindTranscriptUnavailable, and ensure
audioErr/sttErr failures are returned without the transcript-unavailable marker;
use transcribeAudioPath, isEmptyTranscriptionError, and the existing
ErrorKindTranscriptUnavailable handling to locate the fix.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:93208c55b606e21515d2c0a8 -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: In the caption-to-STT fallback path, audio download and generic STT failures are joined with the original transcript-unavailable error, so callers can misclassify backend failures as caption absence. I will keep `ErrorKindTranscriptUnavailable` only for empty-transcript fallback and return audio/STT failures without that marker.
