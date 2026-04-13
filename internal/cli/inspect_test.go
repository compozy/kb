package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/vault"
)

func TestToSmellOutputListsSymbolsAndFilesWithSmells(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"smells":      []string{"dead-export"},
				"source_path": "src/a.go",
				"symbol_kind": "function",
				"symbol_name": "A",
			}),
			testVaultDocument(map[string]any{
				"smells":      []string{},
				"source_path": "src/b.go",
				"symbol_kind": "function",
				"symbol_name": "B",
			}),
		},
		Files: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"smells":      []string{"orphan-file"},
				"source_path": "src/file.go",
			}),
		},
	}

	output := toSmellOutput(snapshot, "")
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 smell rows, got %d", len(output.Data))
	}

	if output.Data[0]["kind"] != "file" || output.Data[0]["name"] != "src/file.go" {
		t.Fatalf("unexpected first smell row %#v", output.Data[0])
	}
	if output.Data[1]["kind"] != "symbol" || output.Data[1]["name"] != "A" {
		t.Fatalf("unexpected second smell row %#v", output.Data[1])
	}
}

func TestToDeadCodeOutputListsDeadExportsAndOrphanFiles(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"is_dead_export": true,
				"smells":         []string{"dead-export"},
				"source_path":    "src/export.go",
				"symbol_kind":    "function",
				"symbol_name":    "Exported",
			}),
		},
		Files: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"is_orphan_file": true,
				"smells":         []string{"orphan-file"},
				"source_path":    "src/orphan.go",
			}),
		},
	}

	output := toDeadCodeOutput(snapshot)
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 dead-code rows, got %d", len(output.Data))
	}

	if output.Data[0]["reason"] != "orphan-file" {
		t.Fatalf("unexpected first dead-code row %#v", output.Data[0])
	}
	if output.Data[1]["reason"] != "dead-export" {
		t.Fatalf("unexpected second dead-code row %#v", output.Data[1])
	}
}

func TestToComplexityOutputSortsByDescendingComplexity(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"blast_radius":          2,
				"cyclomatic_complexity": 5,
				"loc":                   20,
				"smells":                []string{"long-function"},
				"source_path":           "src/b.go",
				"symbol_kind":           "function",
				"symbol_name":           "B",
			}),
			testVaultDocument(map[string]any{
				"blast_radius":          1,
				"cyclomatic_complexity": 9,
				"loc":                   30,
				"smells":                []string{},
				"source_path":           "src/a.go",
				"symbol_kind":           "method",
				"symbol_name":           "A",
			}),
			testVaultDocument(map[string]any{
				"cyclomatic_complexity": 50,
				"source_path":           "src/c.go",
				"symbol_kind":           "struct",
				"symbol_name":           "Ignored",
			}),
		},
	}

	output := toComplexityOutput(snapshot, 20)
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 complexity rows, got %d", len(output.Data))
	}

	if output.Data[0]["symbol_name"] != "A" || output.Data[1]["symbol_name"] != "B" {
		t.Fatalf("unexpected complexity order %#v", output.Data)
	}
}

func TestToBlastRadiusOutputSortsByDescendingBlastRadius(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"blast_radius":             4,
				"centrality":               0.3,
				"external_reference_count": 1,
				"smells":                   []string{"bottleneck"},
				"source_path":              "src/b.go",
				"symbol_name":              "B",
			}),
			testVaultDocument(map[string]any{
				"blast_radius":             7,
				"centrality":               0.8,
				"external_reference_count": 3,
				"smells":                   []string{"high-blast-radius"},
				"source_path":              "src/a.go",
				"symbol_name":              "A",
			}),
		},
	}

	output := toBlastRadiusOutput(snapshot, 0, 0)
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 blast-radius rows, got %d", len(output.Data))
	}

	if output.Data[0]["symbol_name"] != "A" || output.Data[1]["symbol_name"] != "B" {
		t.Fatalf("unexpected blast-radius order %#v", output.Data)
	}
}

func TestToCouplingOutputSortsByInstability(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Files: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"afferent_coupling":       1,
				"efferent_coupling":       4,
				"has_circular_dependency": true,
				"instability":             0.8,
				"smells":                  []string{"god-file"},
				"source_path":             "src/a.go",
			}),
			testVaultDocument(map[string]any{
				"afferent_coupling":       2,
				"efferent_coupling":       1,
				"has_circular_dependency": false,
				"instability":             0.3,
				"smells":                  []string{},
				"source_path":             "src/b.go",
			}),
		},
	}

	output := toCouplingOutput(snapshot, false)
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 coupling rows, got %d", len(output.Data))
	}

	if output.Data[0]["source_path"] != "src/a.go" || output.Data[1]["source_path"] != "src/b.go" {
		t.Fatalf("unexpected coupling order %#v", output.Data)
	}
}

func TestToSymbolLookupOutputReturnsDetailForSingleMatch(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			{
				RelativePath: "raw/codebase/symbols/alpha.md",
				Frontmatter: map[string]any{
					"blast_radius":             4,
					"centrality":               0.75,
					"cyclomatic_complexity":    9,
					"end_line":                 18,
					"exported":                 true,
					"external_reference_count": 2,
					"is_dead_export":           false,
					"is_long_function":         true,
					"language":                 "go",
					"loc":                      12,
					"smells":                   []string{"bottleneck"},
					"source_path":              "src/alpha.go",
					"start_line":               7,
					"symbol_kind":              "function",
					"symbol_name":              "Alpha",
				},
				Body: strings.Join([]string{
					"# Codebase Symbol: Alpha",
					"",
					"## Signature",
					"```text",
					"func Alpha(input string) string",
					"```",
				}, "\n"),
				OutgoingRelations: []vault.VaultRelation{
					{TargetPath: "demo-topic/raw/codebase/symbols/beta", Type: "calls", Confidence: "semantic"},
				},
				Backlinks: []vault.VaultRelation{
					{TargetPath: "demo-topic/raw/codebase/symbols/gamma", Type: "references", Confidence: "syntactic"},
				},
			},
		},
	}

	output, err := toSymbolLookupOutput(snapshot, "alpha")
	if err != nil {
		t.Fatalf("toSymbolLookupOutput returned error: %v", err)
	}

	if got := detailOutputValue(t, output, "signature"); got != "func Alpha(input string) string" {
		t.Fatalf("signature = %#v, want func Alpha(input string) string", got)
	}
	if got := detailOutputValue(t, output, "symbol_name"); got != "Alpha" {
		t.Fatalf("symbol_name = %#v, want Alpha", got)
	}
	if got := detailOutputValue(t, output, "cyclomatic_complexity"); got != 9 {
		t.Fatalf("cyclomatic_complexity = %#v, want 9", got)
	}

	outgoingRelations, ok := detailOutputValue(t, output, "outgoing_relations").([]map[string]any)
	if !ok || len(outgoingRelations) != 1 {
		t.Fatalf("unexpected outgoing_relations %#v", detailOutputValue(t, output, "outgoing_relations"))
	}
	if outgoingRelations[0]["target_path"] != "demo-topic/raw/codebase/symbols/beta" {
		t.Fatalf("unexpected outgoing relation %#v", outgoingRelations[0])
	}
}

func TestToSymbolLookupOutputReturnsSummaryForMultipleMatches(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			{Frontmatter: map[string]any{
				"language":    "go",
				"source_path": "src/zeta.go",
				"start_line":  18,
				"symbol_kind": "function",
				"symbol_name": "AlphaWorker",
			}},
			{Frontmatter: map[string]any{
				"language":    "go",
				"source_path": "src/alpha.go",
				"start_line":  4,
				"symbol_kind": "method",
				"symbol_name": "AlphaThing",
			}},
		},
	}

	output, err := toSymbolLookupOutput(snapshot, "alpha")
	if err != nil {
		t.Fatalf("toSymbolLookupOutput returned error: %v", err)
	}
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 summary rows, got %d", len(output.Data))
	}
	if output.Data[0]["symbol_name"] != "AlphaThing" || output.Data[1]["symbol_name"] != "AlphaWorker" {
		t.Fatalf("unexpected symbol summary order %#v", output.Data)
	}
}

func TestToSymbolLookupOutputReturnsDescriptiveErrorForUnknownName(t *testing.T) {
	t.Parallel()

	_, err := toSymbolLookupOutput(vault.VaultSnapshot{}, "missing")
	if err == nil {
		t.Fatal("expected symbol lookup error")
	}
	if !strings.Contains(err.Error(), `no symbols matched "missing"`) {
		t.Fatalf("unexpected error %q", err)
	}
}

func TestToFileLookupOutputIncludesContainedSymbolsAndMetrics(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Files: []vault.VaultDocument{
			{
				RelativePath: "raw/codebase/files/src/alpha.go.md",
				Frontmatter: map[string]any{
					"afferent_coupling":       2,
					"efferent_coupling":       5,
					"has_circular_dependency": true,
					"instability":             0.7143,
					"is_god_file":             false,
					"is_orphan_file":          false,
					"language":                "go",
					"smells":                  []string{"god-file"},
					"source_path":             "src/alpha.go",
					"symbol_count":            2,
				},
				OutgoingRelations: []vault.VaultRelation{
					{TargetPath: "demo-topic/raw/codebase/files/src/beta.go", Type: "imports", Confidence: "semantic"},
				},
				Backlinks: []vault.VaultRelation{
					{TargetPath: "demo-topic/raw/codebase/files/src/main.go", Type: "imports", Confidence: "semantic"},
				},
			},
		},
		Symbols: []vault.VaultDocument{
			{Frontmatter: map[string]any{
				"source_path": "src/alpha.go",
				"start_line":  10,
				"symbol_kind": "function",
				"symbol_name": "Alpha",
			}},
			{Frontmatter: map[string]any{
				"source_path": "src/alpha.go",
				"start_line":  20,
				"symbol_kind": "method",
				"symbol_name": "Beta",
			}},
		},
	}

	output, err := toFileLookupOutput(snapshot, "src/alpha.go")
	if err != nil {
		t.Fatalf("toFileLookupOutput returned error: %v", err)
	}

	if got := detailOutputValue(t, output, "source_path"); got != "src/alpha.go" {
		t.Fatalf("source_path = %#v, want src/alpha.go", got)
	}
	if got := detailOutputValue(t, output, "instability"); got != 0.7143 {
		t.Fatalf("instability = %#v, want 0.7143", got)
	}

	symbols := detailOutputStringSlice(t, output, "symbols")
	if want := []string{"Alpha (function)", "Beta (method)"}; !reflect.DeepEqual(symbols, want) {
		t.Fatalf("symbols = %#v, want %#v", symbols, want)
	}
}

func TestToFileLookupOutputReturnsDescriptiveErrorForUnknownPath(t *testing.T) {
	t.Parallel()

	_, err := toFileLookupOutput(vault.VaultSnapshot{}, "missing.go")
	if err == nil {
		t.Fatal("expected file lookup error")
	}
	if !strings.Contains(err.Error(), `no file matched "missing.go"`) {
		t.Fatalf("unexpected error %q", err)
	}
}

func TestToBacklinksOutputListsIncomingReferencesForSymbol(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			{
				Frontmatter: map[string]any{
					"symbol_name": "Alpha",
				},
				Backlinks: []vault.VaultRelation{
					{TargetPath: "demo-topic/raw/codebase/symbols/beta", Type: "calls", Confidence: "semantic"},
					{TargetPath: "demo-topic/raw/codebase/files/src/main.go", Type: "references", Confidence: "syntactic"},
				},
			},
		},
	}

	output, err := toBacklinksOutput(snapshot, "alpha")
	if err != nil {
		t.Fatalf("toBacklinksOutput returned error: %v", err)
	}
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 backlink rows, got %d", len(output.Data))
	}
	if output.Data[0]["target_path"] != "demo-topic/raw/codebase/files/src/main.go" {
		t.Fatalf("unexpected backlink order %#v", output.Data)
	}
}

func TestToDependencyOutputListsOutgoingDependenciesForSymbol(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			{
				Frontmatter: map[string]any{
					"symbol_name": "Alpha",
				},
				OutgoingRelations: []vault.VaultRelation{
					{TargetPath: "demo-topic/raw/codebase/symbols/beta", Type: "calls", Confidence: "semantic"},
				},
			},
		},
	}

	output, err := toDependencyOutput(snapshot, "alpha")
	if err != nil {
		t.Fatalf("toDependencyOutput returned error: %v", err)
	}
	if len(output.Data) != 1 {
		t.Fatalf("expected 1 dependency row, got %d", len(output.Data))
	}
	if output.Data[0]["type"] != "calls" {
		t.Fatalf("unexpected dependency row %#v", output.Data[0])
	}
}

func TestToDependencyOutputSupportsExactFilePathLookup(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Files: []vault.VaultDocument{
			{
				Frontmatter: map[string]any{
					"source_path": "src/alpha.go",
				},
				OutgoingRelations: []vault.VaultRelation{
					{TargetPath: "demo-topic/raw/codebase/files/src/beta.go", Type: "imports", Confidence: "semantic"},
				},
			},
		},
	}

	output, err := toDependencyOutput(snapshot, "src/alpha.go")
	if err != nil {
		t.Fatalf("toDependencyOutput returned error: %v", err)
	}
	if len(output.Data) != 1 || output.Data[0]["target_path"] != "demo-topic/raw/codebase/files/src/beta.go" {
		t.Fatalf("unexpected file dependency rows %#v", output.Data)
	}
}

func TestToCircularDepsOutputListsFilesWithCircularDependencyFlag(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		TopicSlug: "demo-topic",
		Files: []vault.VaultDocument{
			testFileDocumentForCycle("raw/codebase/files/src/a.go.md", "src/a.go", map[string]any{
				"afferent_coupling":       1,
				"efferent_coupling":       2,
				"has_circular_dependency": true,
				"instability":             0.6667,
				"smells":                  []string{"god-file"},
			}),
			testFileDocumentForCycle("raw/codebase/files/src/b.go.md", "src/b.go", map[string]any{
				"afferent_coupling":       2,
				"efferent_coupling":       1,
				"has_circular_dependency": true,
				"instability":             0.3333,
				"smells":                  []string{"orphan-file"},
			}),
			testFileDocumentForCycle("raw/codebase/files/src/c.go.md", "src/c.go", map[string]any{
				"has_circular_dependency": false,
			}),
		},
	}

	output := toCircularDepsOutput(snapshot)
	if len(output.Data) != 2 {
		t.Fatalf("expected 2 circular dependency rows, got %d", len(output.Data))
	}
	if want := []string{"source_path", "afferent_coupling", "efferent_coupling", "instability", "smells"}; !reflect.DeepEqual(output.Columns, want) {
		t.Fatalf("columns = %#v, want %#v", output.Columns, want)
	}
	if got := []string{output.Data[0]["source_path"].(string), output.Data[1]["source_path"].(string)}; !reflect.DeepEqual(got, []string{"src/a.go", "src/b.go"}) {
		t.Fatalf("source paths = %#v, want [src/a.go src/b.go]", got)
	}
	if got := output.Data[0]["smells"].([]string); !reflect.DeepEqual(got, []string{"god-file"}) {
		t.Fatalf("unexpected first smells %#v", got)
	}
}

func TestToCircularDepsOutputFallsBackToSCCDetectionForLegacyVaults(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		TopicSlug: "demo-topic",
		Files: []vault.VaultDocument{
			testFileDocumentForCycle("raw/codebase/files/src/a.go.md", "src/a.go", nil, "demo-topic/raw/codebase/files/src/b.go"),
			testFileDocumentForCycle("raw/codebase/files/src/b.go.md", "src/b.go", nil, "demo-topic/raw/codebase/files/src/a.go"),
			testFileDocumentForCycle("raw/codebase/files/src/c.go.md", "src/c.go", nil, "demo-topic/raw/codebase/files/src/d.go"),
			testFileDocumentForCycle("raw/codebase/files/src/d.go.md", "src/d.go", nil, "demo-topic/raw/codebase/files/src/e.go"),
			testFileDocumentForCycle("raw/codebase/files/src/e.go.md", "src/e.go", nil, "demo-topic/raw/codebase/files/src/c.go"),
		},
	}

	output := toCircularDepsOutput(snapshot)
	if len(output.Data) != 5 {
		t.Fatalf("expected 5 circular dependency rows, got %d", len(output.Data))
	}

	got := make([]string, 0, len(output.Data))
	for _, row := range output.Data {
		got = append(got, row["source_path"].(string))
	}

	want := []string{"src/a.go", "src/b.go", "src/c.go", "src/d.go", "src/e.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("source paths = %#v, want %#v", got, want)
	}
}

func TestToCircularDepsOutputShowsMessageWhenNoCycles(t *testing.T) {
	t.Parallel()

	output := toCircularDepsOutput(vault.VaultSnapshot{TopicSlug: "demo-topic"})
	if len(output.Data) != 1 {
		t.Fatalf("expected single message row, got %d", len(output.Data))
	}
	if output.Data[0]["message"] != "no circular dependencies found" {
		t.Fatalf("unexpected empty-cycle output %#v", output.Data[0])
	}
}

func TestInspectCommandJSONFormatProducesValidJSON(t *testing.T) {
	t.Parallel()

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"inspect", "complexity", "--format", "json", "--vault", createInspectTestVault(t)})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	var decoded []map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("expected valid JSON output, got %v\n%s", err, stdout.String())
	}
	if len(decoded) == 0 {
		t.Fatal("expected non-empty complexity JSON output")
	}
}

func TestInspectCommandTSVFormatProducesHeaderAndRows(t *testing.T) {
	t.Parallel()

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"inspect", "coupling", "--format", "tsv", "--vault", createInspectTestVault(t)})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	rendered := stdout.String()
	if !strings.HasPrefix(rendered, "source_path\tafferent_coupling\tefferent_coupling\tinstability\thas_circular_dependency\tsmells\n") {
		t.Fatalf("unexpected TSV header:\n%s", rendered)
	}
	if !strings.Contains(rendered, "src/alpha.go") {
		t.Fatalf("expected TSV output to include a file row, got:\n%s", rendered)
	}
}

func TestResolveInspectContextReadsCodebaseSubtree(t *testing.T) {
	originalResolve := resolveInspectVaultQuery
	originalRead := readInspectVaultSnapshot
	originalGetwd := inspectGetwd
	t.Cleanup(func() {
		resolveInspectVaultQuery = originalResolve
		readInspectVaultSnapshot = originalRead
		inspectGetwd = originalGetwd
	})

	var gotQuery vault.VaultQueryOptions
	var gotResolved vault.ResolvedVault

	inspectGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveInspectVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		gotQuery = options
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/demo-topic",
			TopicSlug: "demo-topic",
		}, nil
	}
	readInspectVaultSnapshot = func(resolvedVault vault.ResolvedVault) (vault.VaultSnapshot, error) {
		gotResolved = resolvedVault
		return vault.VaultSnapshot{}, nil
	}

	_, err := resolveInspectContext(&inspectSharedOptions{
		Format: "json",
		Topic:  "demo-topic",
	})
	if err != nil {
		t.Fatalf("resolveInspectContext returned error: %v", err)
	}

	if want := (vault.VaultQueryOptions{CWD: "/workspace/repo", Topic: "demo-topic"}); gotQuery != want {
		t.Fatalf("vault query = %#v, want %#v", gotQuery, want)
	}
	if gotResolved.TopicPath != "/vault/demo-topic/raw/codebase" {
		t.Fatalf("inspect topic path = %q, want /vault/demo-topic/raw/codebase", gotResolved.TopicPath)
	}
	if gotResolved.TopicSlug != "demo-topic" {
		t.Fatalf("inspect topic slug = %q, want demo-topic", gotResolved.TopicSlug)
	}
}

func TestInspectCommandReturnsDescriptiveErrorForMissingVault(t *testing.T) {
	t.Parallel()

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"inspect", "smells", "--vault", filepath.Join(t.TempDir(), "missing-vault")})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected missing vault error")
	}
	if !strings.Contains(err.Error(), "Vault path was not found or is not a directory") {
		t.Fatalf("unexpected error message %q", err)
	}
}

func TestInspectSubcommandsRespondToHelp(t *testing.T) {
	t.Parallel()

	subcommands := []string{"smells", "dead-code", "complexity", "blast-radius", "coupling", "symbol", "file", "backlinks", "deps", "circular-deps"}
	for _, subcommand := range subcommands {
		subcommand := subcommand
		t.Run(subcommand, func(t *testing.T) {
			t.Parallel()

			command := newRootCommand()
			var stdout bytes.Buffer
			command.SetOut(&stdout)
			command.SetErr(new(bytes.Buffer))
			command.SetArgs([]string{"inspect", subcommand, "--help"})

			if err := command.ExecuteContext(context.Background()); err != nil {
				t.Fatalf("ExecuteContext returned error: %v", err)
			}

			if !strings.Contains(stdout.String(), "Usage:") {
				t.Fatalf("expected help output for %s, got:\n%s", subcommand, stdout.String())
			}
			if !strings.Contains(stdout.String(), "--topic") {
				t.Fatalf("expected help output for %s to include --topic, got:\n%s", subcommand, stdout.String())
			}
		})
	}
}

func testVaultDocument(frontmatter map[string]any) vault.VaultDocument {
	return vault.VaultDocument{Frontmatter: frontmatter}
}

func createInspectTestVault(t *testing.T) string {
	t.Helper()

	vaultPath := filepath.Join(t.TempDir(), "vault")
	topicPath := filepath.Join(vaultPath, "demo-topic")
	mkdirAll(t, topicPath)

	writeInspectMarkdown(t, topicPath, "CLAUDE.md", "# Topic\n")
	writeInspectMarkdown(t, topicPath, "raw/codebase/symbols/alpha.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-symbol"`,
		`language: "go"`,
		`symbol_name: "Alpha"`,
		`symbol_kind: "function"`,
		`source_path: "src/alpha.go"`,
		"exported: true",
		"start_line: 9",
		"end_line: 21",
		"blast_radius: 6",
		"centrality: 0.75",
		"external_reference_count: 2",
		"cyclomatic_complexity: 9",
		"loc: 21",
		"is_dead_export: true",
		"is_long_function: true",
		`smells: ["bottleneck", "dead-export"]`,
		"---",
		"",
		"# Alpha",
		"",
		"## Signature",
		"```text",
		"func Alpha(name string) string",
		"```",
		"",
		"## Outgoing Relations",
		"- `calls` (semantic) -> [[demo-topic/raw/codebase/symbols/beta]]",
		"",
		"## Backlinks",
		"- [[demo-topic/raw/codebase/symbols/gamma]] via `references` (syntactic)",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/symbols/beta.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-symbol"`,
		`language: "go"`,
		`symbol_name: "Beta"`,
		`symbol_kind: "method"`,
		`source_path: "src/beta.go"`,
		"exported: false",
		"start_line: 14",
		"end_line: 18",
		"blast_radius: 3",
		"centrality: 0.25",
		"external_reference_count: 1",
		"cyclomatic_complexity: 4",
		"loc: 10",
		`smells: ["feature-envy"]`,
		"---",
		"",
		"# Beta",
		"",
		"## Signature",
		"```text",
		"func (Service) Beta()",
		"```",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/symbols/gamma.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-symbol"`,
		`language: "go"`,
		`symbol_name: "Gamma"`,
		`symbol_kind: "function"`,
		`source_path: "src/gamma.go"`,
		"start_line: 3",
		"end_line: 6",
		"blast_radius: 1",
		"centrality: 0.1",
		"external_reference_count: 0",
		"loc: 4",
		`smells: []`,
		"---",
		"",
		"# Gamma",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/files/src/alpha.go.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-file"`,
		`language: "go"`,
		`source_path: "src/alpha.go"`,
		"symbol_count: 1",
		"afferent_coupling: 2",
		"efferent_coupling: 5",
		"instability: 0.7142857143",
		"has_circular_dependency: true",
		"is_god_file: false",
		"is_orphan_file: true",
		`smells: ["orphan-file"]`,
		"---",
		"",
		"# File",
		"",
		"## Symbols",
		"- [[demo-topic/raw/codebase/symbols/alpha|Alpha (function)]] · exported=true",
		"",
		"## Outgoing Relations",
		"- `imports` (semantic) -> [[demo-topic/raw/codebase/files/src/beta.go]]",
		"",
		"## Backlinks",
		"- [[demo-topic/raw/codebase/files/src/main.go]] via `imports` (semantic)",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/files/src/beta.go.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-file"`,
		`language: "go"`,
		`source_path: "src/beta.go"`,
		"symbol_count: 1",
		"afferent_coupling: 3",
		"efferent_coupling: 1",
		"instability: 0.25",
		"has_circular_dependency: false",
		"is_god_file: false",
		"is_orphan_file: false",
		`smells: []`,
		"---",
		"",
		"# File",
		"",
		"## Symbols",
		"- [[demo-topic/raw/codebase/symbols/beta|Beta (method)]] · exported=false",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/files/src/main.go.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-file"`,
		`language: "go"`,
		`source_path: "src/main.go"`,
		"symbol_count: 0",
		"afferent_coupling: 0",
		"efferent_coupling: 1",
		"instability: 1",
		"has_circular_dependency: false",
		"is_god_file: false",
		"is_orphan_file: false",
		`smells: []`,
		"---",
		"",
		"# File",
	}, "\n"))

	return vaultPath
}

func writeInspectMarkdown(t *testing.T, topicPath, relativePath, content string) {
	t.Helper()

	fullPath := filepath.Join(topicPath, filepath.FromSlash(relativePath))
	mkdirAll(t, filepath.Dir(fullPath))
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", fullPath, err)
	}
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func detailOutputValue(t *testing.T, output inspectOutput, field string) any {
	t.Helper()

	for _, row := range output.Data {
		if row["field"] == field {
			return row["value"]
		}
	}

	t.Fatalf("field %q not found in output %#v", field, output.Data)
	return nil
}

func detailOutputStringSlice(t *testing.T, output inspectOutput, field string) []string {
	t.Helper()

	value := detailOutputValue(t, output, field)
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		values := make([]string, 0, len(typed))
		for _, entry := range typed {
			if text, ok := entry.(string); ok {
				values = append(values, text)
			}
		}
		return values
	default:
		t.Fatalf("field %q was not a string slice: %#v", field, value)
		return nil
	}
}

func testFileDocumentForCycle(relativePath, sourcePath string, frontmatter map[string]any, targets ...string) vault.VaultDocument {
	relations := make([]vault.VaultRelation, 0, len(targets))
	for _, target := range targets {
		relations = append(relations, vault.VaultRelation{
			TargetPath: target,
			Type:       "imports",
			Confidence: "semantic",
		})
	}

	mergedFrontmatter := map[string]any{"source_path": sourcePath}
	for key, value := range frontmatter {
		mergedFrontmatter[key] = value
	}

	return vault.VaultDocument{
		RelativePath:      relativePath,
		Frontmatter:       mergedFrontmatter,
		OutgoingRelations: relations,
	}
}
