package models

import (
	"context"
	"io"
)

// Converter transforms a source into Markdown content.
type Converter interface {
	// Accepts reports whether the converter supports the given file extension
	// and/or MIME type.
	Accepts(ext string, mimeType string) bool
	// Convert reads from the source and produces Markdown content plus metadata.
	Convert(ctx context.Context, input ConvertInput) (*ConvertResult, error)
}

// ConvertInput carries the source content and metadata needed for conversion.
type ConvertInput struct {
	Reader   io.ReadSeeker  `json:"-"`
	FilePath string         `json:"filePath,omitempty"`
	URL      string         `json:"url,omitempty"`
	Options  map[string]any `json:"options,omitempty"`
}

// ConvertResult contains the Markdown output and metadata from a conversion.
type ConvertResult struct {
	Markdown string         `json:"markdown,omitempty"`
	Title    string         `json:"title,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// SourceKind identifies the source category for ingested KB content.
type SourceKind string

const (
	// SourceKindArticle marks a general article source.
	SourceKindArticle SourceKind = "article"
	// SourceKindGitHubREADME marks an ingested GitHub README.
	SourceKindGitHubREADME SourceKind = "github-readme"
	// SourceKindYouTubeTranscript marks an ingested YouTube transcript.
	SourceKindYouTubeTranscript SourceKind = "youtube-transcript"
	// SourceKindCodebaseFile marks a codebase file snapshot.
	SourceKindCodebaseFile SourceKind = "codebase-file"
	// SourceKindCodebaseSymbol marks a codebase symbol snapshot.
	SourceKindCodebaseSymbol SourceKind = "codebase-symbol"
	// SourceKindBookmarkCluster marks an ingested bookmark cluster.
	SourceKindBookmarkCluster SourceKind = "bookmark-cluster"
	// SourceKindDocument marks a general uploaded document.
	SourceKindDocument SourceKind = "document"
)

// SourceKinds returns every source kind in stable order.
func SourceKinds() []SourceKind {
	return []SourceKind{
		SourceKindArticle,
		SourceKindGitHubREADME,
		SourceKindYouTubeTranscript,
		SourceKindCodebaseFile,
		SourceKindCodebaseSymbol,
		SourceKindBookmarkCluster,
		SourceKindDocument,
	}
}

// IngestResult represents a successfully ingested source.
type IngestResult struct {
	Topic      string     `json:"topic"`
	SourceType SourceKind `json:"sourceType"`
	FilePath   string     `json:"filePath"`
	Title      string     `json:"title"`
}

// LintIssueKind identifies the structural lint issue category.
type LintIssueKind string

const (
	// LintIssueKindDeadLink marks a dead wikilink or reference.
	LintIssueKindDeadLink LintIssueKind = "dead-link"
	// LintIssueKindOrphan marks content with no inbound references.
	LintIssueKindOrphan LintIssueKind = "orphan"
	// LintIssueKindMissingSource marks missing referenced source material.
	LintIssueKindMissingSource LintIssueKind = "missing-source"
	// LintIssueKindStale marks content that is older than its source material.
	LintIssueKindStale LintIssueKind = "stale"
	// LintIssueKindFormat marks frontmatter or structural format violations.
	LintIssueKindFormat LintIssueKind = "format"
)

// LintIssueKinds returns every lint issue kind in stable order.
func LintIssueKinds() []LintIssueKind {
	return []LintIssueKind{
		LintIssueKindDeadLink,
		LintIssueKindOrphan,
		LintIssueKindMissingSource,
		LintIssueKindStale,
		LintIssueKindFormat,
	}
}

// LintIssue represents a single structural problem found in the vault.
type LintIssue struct {
	Kind     LintIssueKind      `json:"kind"`
	Severity DiagnosticSeverity `json:"severity"`
	FilePath string             `json:"filePath,omitempty"`
	Message  string             `json:"message"`
	Target   string             `json:"target,omitempty"`
}

// TopicInfo captures topic metadata for list and info operations.
type TopicInfo struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Domain       string `json:"domain"`
	RootPath     string `json:"rootPath"`
	ArticleCount int    `json:"articleCount"`
	SourceCount  int    `json:"sourceCount"`
	LastLogEntry string `json:"lastLogEntry,omitempty"`
}
