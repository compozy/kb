package vault_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/vault"
)

func TestReadVaultSnapshotParsesFrontmatterAndClassifiesDocuments(t *testing.T) {
	t.Parallel()

	resolvedVault := createResolvedVault(t)

	writeMarkdownDocument(t, resolvedVault.TopicPath, "raw/codebase/symbols/example.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-symbol"`,
		`symbol_name: "normalizePath"`,
		`source_path: "src/path.ts"`,
		"start_line: 12",
		"---",
		"",
		"# Symbol",
	}, "\n"))
	writeMarkdownDocument(t, resolvedVault.TopicPath, "raw/codebase/files/src/path.ts.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-file"`,
		`source_path: "src/path.ts"`,
		"---",
		"",
		"# File",
	}, "\n"))
	writeMarkdownDocument(t, resolvedVault.TopicPath, "raw/codebase/indexes/directories/src.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-directory-index"`,
		`source_path: "src"`,
		"---",
		"",
		"# Directory",
	}, "\n"))
	writeMarkdownDocument(t, resolvedVault.TopicPath, "wiki/concepts/Overview.md", strings.Join([]string{
		"---",
		`type: "wiki"`,
		`title: "Overview"`,
		"---",
		"",
		"# Overview",
	}, "\n"))

	snapshot, err := vault.ReadVaultSnapshot(resolvedVault, vault.ReadVaultOptions{})
	if err != nil {
		t.Fatalf("ReadVaultSnapshot returned error: %v", err)
	}

	if len(snapshot.Symbols) != 1 {
		t.Fatalf("expected 1 symbol document, got %d", len(snapshot.Symbols))
	}
	if len(snapshot.Files) != 1 {
		t.Fatalf("expected 1 file document, got %d", len(snapshot.Files))
	}
	if len(snapshot.Directories) != 1 {
		t.Fatalf("expected 1 directory document, got %d", len(snapshot.Directories))
	}
	if len(snapshot.Wikis) != 1 {
		t.Fatalf("expected 1 wiki document, got %d", len(snapshot.Wikis))
	}

	if got := snapshot.Symbols[0].Frontmatter["symbol_name"]; got != "normalizePath" {
		t.Fatalf("symbol_name = %#v, want normalizePath", got)
	}
	if got := snapshot.Symbols[0].Frontmatter["source_kind"]; got != "codebase-symbol" {
		t.Fatalf("source_kind = %#v, want codebase-symbol", got)
	}
	if got := snapshot.Files[0].Frontmatter["source_kind"]; got != "codebase-file" {
		t.Fatalf("file source_kind = %#v, want codebase-file", got)
	}
}

func TestReadVaultSnapshotParsesRelationSections(t *testing.T) {
	t.Parallel()

	resolvedVault := createResolvedVault(t)

	writeMarkdownDocument(t, resolvedVault.TopicPath, "raw/codebase/symbols/example.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-symbol"`,
		`symbol_name: "greet"`,
		"---",
		"",
		"# Codebase Symbol: greet",
		"",
		"## Outgoing Relations",
		"- `calls` (semantic) -> [[demo-topic/raw/codebase/symbols/helper]]",
		"- `references` (syntactic) -> [[demo-topic/raw/codebase/files/src/helpers.ts]]",
		"",
		"## Backlinks",
		"- [[demo-topic/raw/codebase/symbols/caller]] via `calls` (semantic)",
		"- [[demo-topic/raw/codebase/files/src/index.ts]] via `imports` (syntactic)",
	}, "\n"))

	snapshot, err := vault.ReadVaultSnapshot(resolvedVault, vault.ReadVaultOptions{})
	if err != nil {
		t.Fatalf("ReadVaultSnapshot returned error: %v", err)
	}

	document := snapshot.Symbols[0]

	expectedOutgoing := []vault.VaultRelation{
		{Type: "calls", Confidence: "semantic", TargetPath: "demo-topic/raw/codebase/symbols/helper"},
		{Type: "references", Confidence: "syntactic", TargetPath: "demo-topic/raw/codebase/files/src/helpers.ts"},
	}
	if len(document.OutgoingRelations) != len(expectedOutgoing) {
		t.Fatalf("expected %d outgoing relations, got %d", len(expectedOutgoing), len(document.OutgoingRelations))
	}
	for index, relation := range expectedOutgoing {
		if document.OutgoingRelations[index] != relation {
			t.Fatalf("outgoing relation %d = %#v, want %#v", index, document.OutgoingRelations[index], relation)
		}
	}

	expectedBacklinks := []vault.VaultRelation{
		{TargetPath: "demo-topic/raw/codebase/symbols/caller", Type: "calls", Confidence: "semantic"},
		{TargetPath: "demo-topic/raw/codebase/files/src/index.ts", Type: "imports", Confidence: "syntactic"},
	}
	if len(document.Backlinks) != len(expectedBacklinks) {
		t.Fatalf("expected %d backlinks, got %d", len(expectedBacklinks), len(document.Backlinks))
	}
	for index, relation := range expectedBacklinks {
		if document.Backlinks[index] != relation {
			t.Fatalf("backlink %d = %#v, want %#v", index, document.Backlinks[index], relation)
		}
	}
}

func TestReadVaultSnapshotSkipsMalformedYAMLAndWarns(t *testing.T) {
	t.Parallel()

	resolvedVault := createResolvedVault(t)

	writeMarkdownDocument(t, resolvedVault.TopicPath, "raw/codebase/symbols/broken.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-symbol"`,
		`symbol_name: "broken`,
		"---",
		"# Broken",
	}, "\n"))

	warnings := make([]string, 0, 1)
	snapshot, err := vault.ReadVaultSnapshot(resolvedVault, vault.ReadVaultOptions{
		Warn: func(message string) {
			warnings = append(warnings, message)
		},
	})
	if err != nil {
		t.Fatalf("ReadVaultSnapshot returned error: %v", err)
	}

	if len(snapshot.Symbols) != 0 {
		t.Fatalf("expected malformed document to be skipped, got %d symbol documents", len(snapshot.Symbols))
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if !strings.Contains(warnings[0], "malformed YAML frontmatter") {
		t.Fatalf("unexpected warning message %q", warnings[0])
	}
}

func TestExtractSectionReturnsHeadingBody(t *testing.T) {
	t.Parallel()

	body := strings.Join([]string{
		"# Document",
		"",
		"## Overview",
		"Summary line",
		"Second line",
		"",
		"## Details",
		"More detail",
	}, "\n")

	got := vault.ExtractSection(body, "Overview")
	want := "Summary line\nSecond line"
	if got != want {
		t.Fatalf("ExtractSection returned %q, want %q", got, want)
	}
}

func TestExtractSectionReturnsEmptyStringWhenMissing(t *testing.T) {
	t.Parallel()

	if got := vault.ExtractSection("# Document\n", "Missing"); got != "" {
		t.Fatalf("ExtractSection returned %q, want empty string", got)
	}
}

func TestFindSymbolsByNameUsesCaseInsensitivePartialMatch(t *testing.T) {
	t.Parallel()

	resolvedVault := createResolvedVault(t)

	for relativePath, symbolName := range map[string]string{
		"raw/codebase/symbols/normalize.md":   "normalizePath",
		"raw/codebase/symbols/denormalize.md": "deNormalizeValue",
		"raw/codebase/symbols/render.md":      "renderOutput",
	} {
		writeMarkdownDocument(t, resolvedVault.TopicPath, relativePath, strings.Join([]string{
			"---",
			`source_kind: "codebase-symbol"`,
			`symbol_name: "` + symbolName + `"`,
			`source_path: "src/example.ts"`,
			"---",
			"",
			"# Symbol",
		}, "\n"))
	}

	snapshot, err := vault.ReadVaultSnapshot(resolvedVault, vault.ReadVaultOptions{})
	if err != nil {
		t.Fatalf("ReadVaultSnapshot returned error: %v", err)
	}

	matches := vault.FindSymbolsByName(snapshot, "normalize")
	if len(matches) != 2 {
		t.Fatalf("expected 2 symbol matches, got %d", len(matches))
	}

	gotNames := []string{
		matches[0].Frontmatter["symbol_name"].(string),
		matches[1].Frontmatter["symbol_name"].(string),
	}
	wantNames := []string{"deNormalizeValue", "normalizePath"}
	for index := range wantNames {
		if gotNames[index] != wantNames[index] {
			t.Fatalf("match %d = %q, want %q", index, gotNames[index], wantNames[index])
		}
	}

	if matches := vault.FindSymbolsByName(snapshot, "missing"); len(matches) != 0 {
		t.Fatalf("expected no matches for missing query, got %d", len(matches))
	}
}

func createResolvedVault(t *testing.T) vault.ResolvedVault {
	t.Helper()

	workspaceRoot := t.TempDir()
	vaultPath := filepath.Join(workspaceRoot, "vault")
	resolvedVault := vault.ResolvedVault{
		VaultPath: vaultPath,
		TopicPath: filepath.Join(vaultPath, "demo-topic"),
		TopicSlug: "demo-topic",
	}

	if err := os.MkdirAll(resolvedVault.TopicPath, 0o755); err != nil {
		t.Fatalf("create topic path: %v", err)
	}

	return resolvedVault
}

func writeMarkdownDocument(t *testing.T, topicPath, relativePath, content string) {
	t.Helper()

	absolutePath := filepath.Join(topicPath, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("create parent directory for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}
