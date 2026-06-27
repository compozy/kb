---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/ingest/youtube.go
line: 50
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7l_,comment:PRRC_kwDOR-Fawc7Pt_ee
---

# Issue 010: _🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Surface malformed frontmatter instead of skipping it.**

Line 47-Line 50 silently drop any file whose frontmatter no longer parses. That makes resumable channel ingest treat an already-written YouTube document as unseen, so the same video can be ingested again instead of surfacing the corrupted file.

<details>
<summary>🩹 Suggested fix</summary>

```diff
-		values, _, err := frontmatter.Parse(string(data))
-		if err != nil {
-			continue
-		}
+		values, _, err := frontmatter.Parse(string(data))
+		if err != nil {
+			return nil, fmt.Errorf("existing youtube ids: parse %q: %w", path, err)
+		}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		values, _, err := frontmatter.Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("existing youtube ids: parse %q: %w", path, err)
		}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/ingest/youtube.go` around lines 47 - 50, The frontmatter parsing in
the YouTube ingest flow is silently skipping corrupted documents, which can make
already-written videos look unseen. Update the frontmatter.Parse handling in the
ingest logic to surface the parse failure instead of continuing, using the
existing YouTube ingest path that reads values from the parsed frontmatter.
Return or propagate the error from this branch so resumable ingestion fails
loudly on malformed files rather than reprocessing them.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:58eb32c6703927b9bf428ff2 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `ExistingYouTubeVideoIDs` currently continues when `frontmatter.Parse` fails, which can make a corrupted raw YouTube document invisible to resume logic. I will return a contextual parse error and add regression coverage for malformed frontmatter.
