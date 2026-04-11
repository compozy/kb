package models

import "testing"

func TestSourceKinds(t *testing.T) {
	t.Parallel()

	expected := []SourceKind{
		SourceKindArticle,
		SourceKindGitHubREADME,
		SourceKindYouTubeTranscript,
		SourceKindCodebaseFile,
		SourceKindCodebaseSymbol,
		SourceKindBookmarkCluster,
		SourceKindDocument,
	}

	assertUniqueNonEmptyKinds(t, SourceKinds(), expected)
}

func TestLintIssueKinds(t *testing.T) {
	t.Parallel()

	expected := []LintIssueKind{
		LintIssueKindDeadLink,
		LintIssueKindOrphan,
		LintIssueKindMissingSource,
		LintIssueKindStale,
		LintIssueKindFormat,
	}

	assertUniqueNonEmptyKinds(t, LintIssueKinds(), expected)
}

func TestConvertResultZeroValue(t *testing.T) {
	t.Parallel()

	var result ConvertResult

	if result.Markdown != "" {
		t.Fatalf("expected zero Markdown, got %q", result.Markdown)
	}

	if result.Title != "" {
		t.Fatalf("expected zero Title, got %q", result.Title)
	}

	if result.Metadata != nil {
		t.Fatalf("expected nil Metadata, got %#v", result.Metadata)
	}
}

func TestIngestResultZeroValue(t *testing.T) {
	t.Parallel()

	var result IngestResult

	if result.Topic != "" {
		t.Fatalf("expected zero Topic, got %q", result.Topic)
	}

	if result.SourceType != "" {
		t.Fatalf("expected zero SourceType, got %q", result.SourceType)
	}

	if result.FilePath != "" {
		t.Fatalf("expected zero FilePath, got %q", result.FilePath)
	}

	if result.Title != "" {
		t.Fatalf("expected zero Title, got %q", result.Title)
	}
}

func TestLintIssueZeroValue(t *testing.T) {
	t.Parallel()

	var issue LintIssue

	if issue.Kind != "" {
		t.Fatalf("expected zero Kind, got %q", issue.Kind)
	}

	if issue.Severity != "" {
		t.Fatalf("expected zero Severity, got %q", issue.Severity)
	}

	if issue.FilePath != "" {
		t.Fatalf("expected zero FilePath, got %q", issue.FilePath)
	}

	if issue.Message != "" {
		t.Fatalf("expected zero Message, got %q", issue.Message)
	}

	if issue.Target != "" {
		t.Fatalf("expected zero Target, got %q", issue.Target)
	}
}

func TestTopicInfoZeroValue(t *testing.T) {
	t.Parallel()

	var info TopicInfo

	if info.Slug != "" {
		t.Fatalf("expected zero Slug, got %q", info.Slug)
	}

	if info.Title != "" {
		t.Fatalf("expected zero Title, got %q", info.Title)
	}

	if info.Domain != "" {
		t.Fatalf("expected zero Domain, got %q", info.Domain)
	}

	if info.RootPath != "" {
		t.Fatalf("expected zero RootPath, got %q", info.RootPath)
	}

	if info.ArticleCount != 0 {
		t.Fatalf("expected zero ArticleCount, got %d", info.ArticleCount)
	}

	if info.SourceCount != 0 {
		t.Fatalf("expected zero SourceCount, got %d", info.SourceCount)
	}

	if info.LastLogEntry != "" {
		t.Fatalf("expected zero LastLogEntry, got %q", info.LastLogEntry)
	}
}

func assertUniqueNonEmptyKinds[T ~string](t *testing.T, got []T, expected []T) {
	t.Helper()

	if len(got) != len(expected) {
		t.Fatalf("expected %d kinds, got %d", len(expected), len(got))
	}

	seen := make(map[T]struct{}, len(got))
	for index, kind := range got {
		if kind != expected[index] {
			t.Fatalf("kind %d: expected %q, got %q", index, expected[index], kind)
		}

		if kind == "" {
			t.Fatalf("kind %d is empty", index)
		}

		if _, exists := seen[kind]; exists {
			t.Fatalf("kind %d duplicated value %q", index, kind)
		}

		seen[kind] = struct{}{}
	}
}
