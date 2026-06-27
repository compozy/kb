---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/instagram/instagram.go
line: 81
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7mL,comment:PRRC_kwDOR-Fawc7Pt_eu
---

# Issue 011: _🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Backfill `Metadata.URL` with the canonical Instagram URL before returning.**

`parseInstagramURL` already normalizes accepted inputs, including scheme-less ones, but that canonical URL is dropped if yt-dlp omits `webpage_url`. `internal/cli/ingest_instagram.go` then falls back to the raw CLI arg, so a successful ingest can persist an invalid `source_url` into frontmatter via `internal/ingest/ingest.go`.

   

<details>
<summary>Proposed fix</summary>

```diff
 	result, err := extractor.core.Extract(ctx, parsed, options)
 	if err != nil {
 		if mediadl.IsTranscriptUnavailable(err) && result != nil && captionOf(result) != "" {
+			if strings.TrimSpace(result.Metadata.URL) == "" {
+				result.Metadata.URL = parsed.CanonicalURL
+			}
 			result.Markdown = composeBody(captionOf(result), "")
 			result.Source = transcriptSourceNone
 			result.Language = ""
 			result.CaptionKind = ""
 			return result, nil
@@
 		return nil, err
 	}
 
+	if strings.TrimSpace(result.Metadata.URL) == "" {
+		result.Metadata.URL = parsed.CanonicalURL
+	}
 	result.Markdown = composeBody(captionOf(result), result.Markdown)
 	return result, nil
 }
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	result, err := extractor.core.Extract(ctx, parsed, options)
	if err != nil {
		if mediadl.IsTranscriptUnavailable(err) && result != nil && captionOf(result) != "" {
			if strings.TrimSpace(result.Metadata.URL) == "" {
				result.Metadata.URL = parsed.CanonicalURL
			}
			result.Markdown = composeBody(captionOf(result), "")
			result.Source = transcriptSourceNone
			result.Language = ""
			result.CaptionKind = ""
			return result, nil
		}
		return nil, err
	}

	if strings.TrimSpace(result.Metadata.URL) == "" {
		result.Metadata.URL = parsed.CanonicalURL
	}
	result.Markdown = composeBody(captionOf(result), result.Markdown)
	return result, nil
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/instagram/instagram.go` around lines 68 - 81, Backfill the canonical
Instagram URL into Metadata.URL before returning from the Instagram extraction
flow. In the code path around extractor.core.Extract and the post-processing in
internal/instagram/instagram.go, preserve the normalized URL from
parseInstagramURL and assign it to result.Metadata.URL when yt-dlp does not
provide webpage_url. Make sure this value survives both the successful return
path and the transcript-unavailable fallback so internal/cli/ingest_instagram.go
and internal/ingest/ingest.go always receive a valid canonical source URL.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fe2533f8d0c955a8587ab9ec -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `parseInstagramURL` computes a canonical URL, but `Extractor.Extract` does not backfill `result.Metadata.URL` when the media backend omits it. I will set the canonical URL on both the successful result path and the caption-only fallback path when `Metadata.URL` is blank.
