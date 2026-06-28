---
provider: coderabbit
pr: "25"
round: 1
round_created_at: 2026-06-28T02:05:08.396169Z
status: resolved
file: internal/cli/okf.go
line: 74
author: coderabbitai[bot]
provider_ref: thread:PRRT_kwDOR-Fawc6MxHSr,comment:PRRC_kwDOR-Fawc7P2i8f
---

# Issue 001: _🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_
## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Reject non-OKF topics before running bundle checks.**

`runOKFTopicInfo` gives you `TopicInfo.Mode`, but this command ignores it and always hands `RootPath` to `kokf.Check`. Since `kokf.Check` only validates that the path is a directory, `kb okf check <wiki-topic>` will lint the wrong tree instead of failing fast.






<details>
<summary>Suggested fix</summary>

```diff
 import (
 	"fmt"
 	"os"
 
 	"github.com/spf13/cobra"
 
+	"github.com/compozy/kb/internal/models"
 	kokf "github.com/compozy/kb/internal/okf"
 	"github.com/compozy/kb/internal/output"
 	ktopic "github.com/compozy/kb/internal/topic"
 )
@@
 	topicInfo, err := runOKFTopicInfo(vaultPath, topicSlug)
 	if err != nil {
 		return fmt.Errorf("okf check: %w", err)
 	}
+	if topicInfo.Mode != models.TopicModeOKF {
+		return fmt.Errorf("okf check: topic %q is not an OKF topic", topicSlug)
+	}
 
 	issues, err := runOKFCheck(commandContext(cmd), topicInfo.RootPath, kokf.CheckOptions{
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	topicInfo, err := runOKFTopicInfo(vaultPath, topicSlug)
	if err != nil {
		return fmt.Errorf("okf check: %w", err)
	}
	if topicInfo.Mode != models.TopicModeOKF {
		return fmt.Errorf("okf check: topic %q is not an OKF topic", topicSlug)
	}

	issues, err := runOKFCheck(commandContext(cmd), topicInfo.RootPath, kokf.CheckOptions{
		Types:  cfg.OKF.Types,
		Strict: options.Strict,
	})
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/okf.go` around lines 66 - 74, The okf check command currently
ignores TopicInfo.Mode from runOKFTopicInfo and always passes TopicInfo.RootPath
into kokf.Check, so non-OKF topics can be linted as directories instead of being
rejected. Update the okf check flow in the command handler to inspect the
returned TopicInfo.Mode before calling kokf.Check, fail fast for non-OKF topics,
and only run the bundle check when the topic is actually an OKF topic.
```

</details>

<!-- fingerprinting:phantom:poseidon:grasshopper -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2ddd87dea620484e1aac90b5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: `runOKFCheckCommand` resolves `TopicInfo` but currently calls `kokf.Check` without checking `TopicInfo.Mode`. A wiki topic directory can therefore be checked as if it were an OKF bundle. Fix by rejecting non-OKF topics before invoking the checker and add CLI coverage for that path.
