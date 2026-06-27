---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/youtube/channel.go
line: 72
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7mT,comment:PRRC_kwDOR-Fawc7Pt_e2
---

# Issue 019: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Reject only actual YouTube hosts here.**

`strings.Contains(host, "youtube.com")` also matches domains like `notyoutube.com` and `youtube.com.evil.test`, so `NormalizeChannelURL` can accept unsupported URLs and pass them to yt-dlp instead of failing fast.

<details>
<summary>Suggested fix</summary>

```diff
-	if !strings.Contains(host, "youtube.com") {
+	if host != "youtube.com" && !strings.HasSuffix(host, ".youtube.com") {
 		return "", &Error{
 			Kind:    ErrorKindInvalidURL,
 			URL:     raw,
 			Message: fmt.Sprintf("expected a youtube.com channel or playlist URL: %q", raw),
 		}
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	if host != "youtube.com" && !strings.HasSuffix(host, ".youtube.com") {
		return "", &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     raw,
			Message: fmt.Sprintf("expected a youtube.com channel or playlist URL: %q", raw),
		}
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/youtube/channel.go` around lines 66 - 72, The host check in
NormalizeChannelURL is too permissive because strings.Contains(host,
"youtube.com") accepts lookalike domains. Tighten the validation in
NormalizeChannelURL by checking the parsed host against actual YouTube hosts
only, using host-specific matching for youtube.com and its legitimate subdomains
(and rejecting unrelated domains such as notyoutube.com or
youtube.com.evil.test). Keep the existing ErrorKindInvalidURL path and message
behavior for rejected URLs.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1ab54e3567bd9b863b80fdb5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `NormalizeChannelURL` uses `strings.Contains(host, "youtube.com")`, which accepts lookalike domains such as `notyoutube.com` and `youtube.com.evil.test`. I will replace it with exact-host/subdomain matching for real YouTube hosts and extend the URL table.
