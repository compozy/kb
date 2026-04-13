package topic

import "embed"

var (
	// topicTemplateFS keeps scaffolding templates with the topic package so
	// `kb topic new` does not depend on external skill directory layout.
	//
	//go:embed assets/*.md
	topicTemplateFS embed.FS
)
