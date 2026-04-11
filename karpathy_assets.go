package repoassets

import "embed"

// KarpathyKBTemplates exposes the scaffold templates shipped in the repo.
var (
	//go:embed all:.agents/skills/karpathy-kb/assets/*.md
	KarpathyKBTemplates embed.FS
)
