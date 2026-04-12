package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/user/kb/internal/output"
	"github.com/user/kb/internal/vault"
	"github.com/user/kb/internal/version"
)

func TestParseInspectOutputFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		value     string
		want      output.OutputFormat
		expectErr bool
	}{
		{name: "default", value: "", want: output.OutputFormatTable},
		{name: "table", value: "table", want: output.OutputFormatTable},
		{name: "json", value: "json", want: output.OutputFormatJSON},
		{name: "tsv", value: "tsv", want: output.OutputFormatTSV},
		{name: "invalid", value: "csv", expectErr: true},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseInspectOutputFormat(testCase.value)
			if testCase.expectErr {
				if err == nil {
					t.Fatalf("expected error for %q", testCase.value)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseInspectOutputFormat returned error: %v", err)
			}
			if got != testCase.want {
				t.Fatalf("format = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestInspectFrontmatterHelpers(t *testing.T) {
	t.Parallel()

	document := vault.VaultDocument{
		Frontmatter: map[string]any{
			"array_any":    []any{"a", "b"},
			"array_string": []string{"c", "d"},
			"bool_native":  true,
			"bool_string":  "true",
			"float_int":    int(4),
			"float_string": "3.5",
			"string_key":   "value",
		},
	}

	if got := inspectFrontmatterString(document, "string_key"); got != "value" {
		t.Fatalf("string = %q, want value", got)
	}
	if got := inspectFrontmatterString(document, "missing"); got != "" {
		t.Fatalf("missing string = %q, want empty", got)
	}

	if got := inspectFrontmatterStringArray(document, "array_string"); len(got) != 2 || got[0] != "c" || got[1] != "d" {
		t.Fatalf("unexpected []string result %#v", got)
	}
	if got := inspectFrontmatterStringArray(document, "array_any"); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("unexpected []any result %#v", got)
	}
	if got := inspectFrontmatterStringArray(document, "missing"); len(got) != 0 {
		t.Fatalf("missing string array = %#v, want empty", got)
	}

	if !inspectFrontmatterBool(document, "bool_native") {
		t.Fatal("expected native bool to parse as true")
	}
	if !inspectFrontmatterBool(document, "bool_string") {
		t.Fatal("expected string bool to parse as true")
	}
	if inspectFrontmatterBool(testVaultDocument(map[string]any{"value": "nope"}), "value") {
		t.Fatal("expected invalid bool string to parse as false")
	}

	intCases := []struct {
		name  string
		value any
		want  int
	}{
		{name: "int", value: int(1), want: 1},
		{name: "int8", value: int8(2), want: 2},
		{name: "int16", value: int16(3), want: 3},
		{name: "int32", value: int32(4), want: 4},
		{name: "int64", value: int64(5), want: 5},
		{name: "uint", value: uint(6), want: 6},
		{name: "uint8", value: uint8(7), want: 7},
		{name: "uint16", value: uint16(8), want: 8},
		{name: "uint32", value: uint32(9), want: 9},
		{name: "uint64", value: uint64(10), want: 10},
		{name: "float32", value: float32(11), want: 11},
		{name: "float64", value: float64(12), want: 12},
		{name: "string", value: "13", want: 13},
	}

	for _, testCase := range intCases {
		if got := inspectFrontmatterInt(testVaultDocument(map[string]any{"value": testCase.value}), "value"); got != testCase.want {
			t.Fatalf("%s int = %d, want %d", testCase.name, got, testCase.want)
		}
	}
	if got := inspectFrontmatterInt(testVaultDocument(map[string]any{"value": "bad"}), "value"); got != 0 {
		t.Fatalf("invalid int string = %d, want 0", got)
	}

	floatCases := []struct {
		name  string
		value any
		want  float64
	}{
		{name: "float64", value: float64(1.5), want: 1.5},
		{name: "float32", value: float32(2.5), want: 2.5},
		{name: "int", value: int(3), want: 3},
		{name: "int8", value: int8(4), want: 4},
		{name: "int16", value: int16(5), want: 5},
		{name: "int32", value: int32(6), want: 6},
		{name: "int64", value: int64(7), want: 7},
		{name: "uint", value: uint(8), want: 8},
		{name: "uint8", value: uint8(9), want: 9},
		{name: "uint16", value: uint16(10), want: 10},
		{name: "uint32", value: uint32(11), want: 11},
		{name: "uint64", value: uint64(12), want: 12},
		{name: "string", value: "13.5", want: 13.5},
	}

	for _, testCase := range floatCases {
		if got := inspectFrontmatterFloat(testVaultDocument(map[string]any{"value": testCase.value}), "value"); got != testCase.want {
			t.Fatalf("%s float = %v, want %v", testCase.name, got, testCase.want)
		}
	}
	if got := inspectFrontmatterFloat(testVaultDocument(map[string]any{"value": "bad"}), "value"); got != 0 {
		t.Fatalf("invalid float string = %v, want 0", got)
	}
}

func TestToSmellOutputFiltersByType(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"smells":      []string{"dead-export", "bottleneck"},
				"source_path": "src/a.go",
				"symbol_kind": "function",
				"symbol_name": "A",
			}),
		},
		Files: []vault.VaultDocument{
			testVaultDocument(map[string]any{
				"smells":      []string{"orphan-file"},
				"source_path": "src/b.go",
			}),
		},
	}

	output := toSmellOutput(snapshot, "dead-export")
	if len(output.Data) != 1 || output.Data[0]["name"] != "A" {
		t.Fatalf("unexpected filtered smell output %#v", output.Data)
	}
}

func TestToBlastRadiusOutputRespectsMinimumAndTop(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Symbols: []vault.VaultDocument{
			testVaultDocument(map[string]any{"blast_radius": 1, "source_path": "src/a.go", "symbol_name": "A"}),
			testVaultDocument(map[string]any{"blast_radius": 5, "source_path": "src/b.go", "symbol_name": "B"}),
			testVaultDocument(map[string]any{"blast_radius": 3, "source_path": "src/c.go", "symbol_name": "C"}),
		},
	}

	output := toBlastRadiusOutput(snapshot, 2, 1)
	if len(output.Data) != 1 {
		t.Fatalf("expected 1 blast-radius row, got %d", len(output.Data))
	}
	if output.Data[0]["symbol_name"] != "B" {
		t.Fatalf("unexpected blast-radius row %#v", output.Data[0])
	}
}

func TestToCouplingOutputFiltersUnstableOnly(t *testing.T) {
	t.Parallel()

	snapshot := vault.VaultSnapshot{
		Files: []vault.VaultDocument{
			testVaultDocument(map[string]any{"instability": 0.8, "source_path": "src/a.go"}),
			testVaultDocument(map[string]any{"instability": 0.4, "source_path": "src/b.go"}),
		},
	}

	output := toCouplingOutput(snapshot, true)
	if len(output.Data) != 1 || output.Data[0]["source_path"] != "src/a.go" {
		t.Fatalf("unexpected unstable-only output %#v", output.Data)
	}
}

func TestInspectCommandValidationErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "complexity top",
			args: []string{"inspect", "complexity", "--top", "-1", "--vault", createInspectTestVault(t)},
			want: "--top must be >= 1",
		},
		{
			name: "blast min",
			args: []string{"inspect", "blast-radius", "--min", "-1", "--vault", createInspectTestVault(t)},
			want: "--min must be >= 0",
		},
		{
			name: "blast top",
			args: []string{"inspect", "blast-radius", "--top", "-1", "--vault", createInspectTestVault(t)},
			want: "--top must be >= 0",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command := newRootCommand()
			command.SetOut(new(bytes.Buffer))
			command.SetErr(new(bytes.Buffer))
			command.SetArgs(testCase.args)

			err := command.ExecuteContext(context.Background())
			if err == nil {
				t.Fatal("expected validation error")
			}
			if got := err.Error(); got == "" || !bytes.Contains([]byte(got), []byte(testCase.want)) {
				t.Fatalf("unexpected error %q, want substring %q", got, testCase.want)
			}
		})
	}
}

func TestInspectParentHelpListsAllSubcommands(t *testing.T) {
	t.Parallel()

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"inspect"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	helpText := stdout.String()
	for _, subcommand := range []string{"smells", "dead-code", "complexity", "blast-radius", "coupling", "symbol", "file", "backlinks", "deps", "circular-deps"} {
		if !bytes.Contains(stdout.Bytes(), []byte(subcommand)) {
			t.Fatalf("expected inspect help to list %s, got:\n%s", subcommand, helpText)
		}
	}
}

func TestVersionCommandPrintsBuildVersion(t *testing.T) {
	t.Parallel()

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"version"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if got := stdout.String(); got != version.String()+"\n" {
		t.Fatalf("unexpected version output %q", got)
	}
}
