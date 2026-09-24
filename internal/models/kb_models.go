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
	// SourceKindInstagramVideo marks an ingested Instagram reel/video transcript.
	SourceKindInstagramVideo SourceKind = "instagram-video"
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
		SourceKindInstagramVideo,
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
	// LintIssueKindJavaDiagnosticGovernance marks Java diagnostic governance threshold outcomes.
	LintIssueKindJavaDiagnosticGovernance LintIssueKind = "java-diagnostic-governance"

	// LintIssueKindFrontmatterDeadLink marks a frontmatter wikilink that resolves nowhere.
	LintIssueKindFrontmatterDeadLink LintIssueKind = "frontmatter-dead-link"
	// LintIssueKindLinkToQuarantined marks a link whose target is a quarantined source.
	LintIssueKindLinkToQuarantined LintIssueKind = "link-to-quarantined"
	// LintIssueKindNeedsCompile marks an article older than a source that affects it.
	LintIssueKindNeedsCompile LintIssueKind = "needs-compile"
	// LintIssueKindKeyConflict marks a kb-owned key whose value kb did not write.
	LintIssueKindKeyConflict LintIssueKind = "key-conflict"
	// LintIssueKindContradiction marks an open contradiction review item.
	LintIssueKindContradiction LintIssueKind = "contradiction"
	// LintIssueKindOffTopicKept marks a source judged off-topic that is still kept.
	LintIssueKindOffTopicKept LintIssueKind = "off-topic-kept"
	// LintIssueKindCriterionMissing marks an article without a criterion.
	LintIssueKindCriterionMissing LintIssueKind = "criterion-missing"
	// LintIssueKindUnclassified marks documents without a current classification state.
	LintIssueKindUnclassified LintIssueKind = "unclassified"
	// LintIssueKindContractMissing marks a topic with sources and no accepted contract.
	LintIssueKindContractMissing LintIssueKind = "contract-missing"
	// LintIssueKindContractDraftPending marks a contract draft awaiting acceptance.
	LintIssueKindContractDraftPending LintIssueKind = "contract-draft-pending"
	// LintIssueKindVocabularyMissing marks a topic with sources but no concept vocabulary.
	LintIssueKindVocabularyMissing LintIssueKind = "vocabulary-missing"
	// LintIssueKindPendingReview counts pending review-queue items per queue.
	LintIssueKindPendingReview LintIssueKind = "pending-review"
	// LintIssueKindRecapturePending counts pending items of the recapture queue.
	LintIssueKindRecapturePending LintIssueKind = "recapture-pending"
	// LintIssueKindRemovePending counts pending items of the remove queue.
	LintIssueKindRemovePending LintIssueKind = "remove-pending"

	// LintIssueKindTypeMismatch is the advisory `kb okf check` finding for a
	// concept whose stored type suggestion disagrees with its `type`.
	LintIssueKindTypeMismatch LintIssueKind = "type_mismatch"
	// LintIssueKindDescriptionUnsupported is the advisory `kb okf check`
	// finding for a description the concept body does not support.
	LintIssueKindDescriptionUnsupported LintIssueKind = "description_unsupported"
)

// SeverityInfo marks advisory findings (queue counts, pending drafts,
// uncalibrated OKF suggestions) that never fail a command.
const SeverityInfo DiagnosticSeverity = "info"

// LintIssueKinds returns every lint issue kind in stable order.
func LintIssueKinds() []LintIssueKind {
	return []LintIssueKind{
		LintIssueKindDeadLink,
		LintIssueKindFrontmatterDeadLink,
		LintIssueKindLinkToQuarantined,
		LintIssueKindOrphan,
		LintIssueKindMissingSource,
		LintIssueKindStale,
		LintIssueKindNeedsCompile,
		LintIssueKindFormat,
		LintIssueKindKeyConflict,
		LintIssueKindContradiction,
		LintIssueKindOffTopicKept,
		LintIssueKindCriterionMissing,
		LintIssueKindUnclassified,
		LintIssueKindContractMissing,
		LintIssueKindContractDraftPending,
		LintIssueKindVocabularyMissing,
		LintIssueKindPendingReview,
		LintIssueKindRecapturePending,
		LintIssueKindRemovePending,
		LintIssueKindJavaDiagnosticGovernance,
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
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Domain       string    `json:"domain"`
	Mode         TopicMode `json:"mode"`
	RootPath     string    `json:"rootPath"`
	ArticleCount int       `json:"articleCount"`
	SourceCount  int       `json:"sourceCount"`
	LastLogEntry string    `json:"lastLogEntry,omitempty"`
}
