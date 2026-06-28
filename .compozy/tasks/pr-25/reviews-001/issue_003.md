---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/cli/okf_test.go
line: 75
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHSb,comment:PRRC_kwDOR-Fawc7P2i8K
---

# Issue 003: _🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Assert that the positional source document reaches `kokf.Promote`.**

This test only checks flags and config-derived fields. If `args[0]` stopped being forwarded, the command would be broken at runtime and this would still pass.

As per path instructions, "`**/*_test.go`: Focus on critical paths: parsing`" and "`Ensure tests verify behavior outcomes, not just function calls`."





<details>
<summary>Suggested assertion</summary>

```diff
 	if gotInput.VaultPath != "/tmp/vault" || gotInput.TargetTopic.Slug != "catalog" || gotInput.Type != "Playbook" {
 		t.Fatalf("unexpected promote input: %#v", gotInput)
 	}
+	if gotInput.SourceDocPath != "research/wiki/concepts/Alpha.md" {
+		t.Fatalf("source doc path = %q, want research/wiki/concepts/Alpha.md", gotInput.SourceDocPath)
+	}
 	if gotInput.Description != "Alpha description." {
 		t.Fatalf("description = %q", gotInput.Description)
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	command.SetArgs([]string{
		"promote", "research/wiki/concepts/Alpha.md",
		"--to", "catalog",
		"--type", "Playbook",
		"--description", "Alpha description.",
		"--vault", "/tmp/vault",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}
	if gotInput.VaultPath != "/tmp/vault" || gotInput.TargetTopic.Slug != "catalog" || gotInput.Type != "Playbook" {
		t.Fatalf("unexpected promote input: %#v", gotInput)
	}
	if gotInput.SourceDocPath != "research/wiki/concepts/Alpha.md" {
		t.Fatalf("source doc path = %q, want research/wiki/concepts/Alpha.md", gotInput.SourceDocPath)
	}
	if gotInput.Description != "Alpha description." {
		t.Fatalf("description = %q", gotInput.Description)
	}
	if len(gotInput.Types) != 1 || gotInput.Types[0] != "Playbook" {
		t.Fatalf("types = %#v, want Playbook", gotInput.Types)
	}

	var result kokf.ConceptResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.WrittenPath != "alpha.md" || result.Type != "Playbook" {
		t.Fatalf("unexpected result: %#v", result)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/okf_test.go` around lines 48 - 75, The test for the promote
command currently verifies flags and output but does not confirm the positional
source document is forwarded into kokf.Promote. Update the assertion around the
promote path in okf_test.go to capture and verify the input passed to Promot(e)
includes the source argument from command args[0]
(research/wiki/concepts/Alpha.md), so a regression in positional-argument
plumbing will fail the test.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fa18721bdfae2345e4e6824c -->

_Source: Path instructions_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `TestPromoteCommandResolvesTargetAndPrintsJSON` captures `kokf.PromoteInput` but does not assert `SourceDocPath`, so losing the positional argument would not fail the test. Fix by asserting the exact source document path forwarded from `args[0]`.
