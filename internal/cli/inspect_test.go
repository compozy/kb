package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/vault"
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

	subcommands := []string{"smells", "dead-code", "complexity", "blast-radius", "coupling"}
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
		`symbol_name: "Alpha"`,
		`symbol_kind: "function"`,
		`source_path: "src/alpha.go"`,
		"blast_radius: 6",
		"centrality: 0.75",
		"external_reference_count: 2",
		"cyclomatic_complexity: 9",
		"loc: 21",
		"is_dead_export: true",
		`smells: ["bottleneck", "dead-export"]`,
		"---",
		"",
		"# Alpha",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/symbols/beta.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-symbol"`,
		`symbol_name: "Beta"`,
		`symbol_kind: "method"`,
		`source_path: "src/beta.go"`,
		"blast_radius: 3",
		"centrality: 0.25",
		"external_reference_count: 1",
		"cyclomatic_complexity: 4",
		"loc: 10",
		`smells: ["feature-envy"]`,
		"---",
		"",
		"# Beta",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/files/src/alpha.go.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-file"`,
		`source_path: "src/alpha.go"`,
		"afferent_coupling: 2",
		"efferent_coupling: 5",
		"instability: 0.7142857143",
		"has_circular_dependency: true",
		"is_orphan_file: true",
		`smells: ["orphan-file"]`,
		"---",
		"",
		"# File",
	}, "\n"))
	writeInspectMarkdown(t, topicPath, "raw/codebase/files/src/beta.go.md", strings.Join([]string{
		"---",
		`source_kind: "codebase-file"`,
		`source_path: "src/beta.go"`,
		"afferent_coupling: 3",
		"efferent_coupling: 1",
		"instability: 0.25",
		"has_circular_dependency: false",
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
