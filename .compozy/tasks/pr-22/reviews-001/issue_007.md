---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/config/config.go
line: 222
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6Mq7l-,comment:PRRC_kwDOR-Fawc7Pt_ed
---

# Issue 007: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Persist the trimmed caption language list.**

`Validate()` trims `CaptionLanguages` only for the emptiness check, but the raw slice is still forwarded downstream unchanged. A TOML value like `["orig", " pt "]` passes validation here and then reaches the ingest path as `" pt "`, which breaks caption selection because the CLI forwards `cfg.YouTube.CaptionLanguages` verbatim.




<details>
<summary>Possible fix</summary>

```diff
func (c *Config) applyDefaults() {
	if c == nil {
		return
	}
+	c.YouTube.CaptionLanguages = normalizeConfigStringList(c.YouTube.CaptionLanguages)
 	if strings.TrimSpace(c.Vault.Root) == "" {
 		c.Vault.Root = defaultVaultRoot
 	}
 	...
 	if len(c.YouTube.CaptionLanguages) == 0 {
 		c.YouTube.CaptionLanguages = []string{defaultYouTubeCaptionLanguage}
 	}
}
```
</details>


Also applies to: 440-442, 473-483

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/config/config.go` around lines 220 - 222, The Validate() flow for
YouTube.CaptionLanguages only trims values to check for emptiness, but leaves
the original slice unchanged, so trimmed languages are not persisted before
later use. Update the validation/normalization logic in config.go so the
validated cfg.YouTube.CaptionLanguages is rewritten with trimmed entries, and
make sure the same normalization is applied in the related YouTube config paths
referenced by Validate() and downstream CLI forwarding so the ingest path
receives the cleaned list.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:93c85f25dca3f2d9d2cc6771 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `applyDefaults` only defaults `YouTube.CaptionLanguages` when the slice is empty, and `YouTubeConfig.Validate` normalizes a copy for the emptiness check without persisting it. I will normalize the configured list during default application so downstream CLI forwarding receives trimmed values.
