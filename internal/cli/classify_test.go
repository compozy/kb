package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/output"
)

func TestClassifyAndVocabularyFlagValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "only-missing with all", args: []string{"classify", "demo", "--only-missing", "--all"}, wantErr: "--only-missing and --all cannot be combined"},
		{name: "classify bad format", args: []string{"classify", "demo", "--format", "xml"}, wantErr: "xml"},
		{name: "classify needs a topic", args: []string{"classify"}, wantErr: "accepts 1 arg"},
		{name: "vocabulary without mode", args: []string{"topic", "vocabulary", "demo"}, wantErr: "exactly one of --draft or --accept"},
		{name: "vocabulary with both modes", args: []string{"topic", "vocabulary", "demo", "--draft", "--accept"}, wantErr: "exactly one of --draft or --accept"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			command := newRootCommand()
			command.SetOut(&bytes.Buffer{})
			command.SetErr(&bytes.Buffer{})
			command.SetArgs(tt.args)
			err := command.ExecuteContext(context.Background())
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestWriteClassifyReport(t *testing.T) {
	t.Parallel()
	report := classify.Report{
		Documents: 3, Judged: 2,
		Skipped:   map[string]int{"unchanged": 1},
		Facets:    map[string]int{"genre": 2},
		Undecided: map[string]int{"kind:timeout": 1},
	}

	var table bytes.Buffer
	if err := writeClassifyReport(&table, output.OutputFormatTable, report); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"metric", "judged", "skipped:unchanged", "written:genre", "undecided:kind:timeout"} {
		if !strings.Contains(table.String(), want) {
			t.Fatalf("table misses %q:\n%s", want, table.String())
		}
	}

	var encoded bytes.Buffer
	if err := writeClassifyReport(&encoded, output.OutputFormatJSON, report); err != nil {
		t.Fatal(err)
	}
	var decoded classify.Report
	if err := json.Unmarshal(encoded.Bytes(), &decoded); err != nil {
		t.Fatalf("json output: %v\n%s", err, encoded.String())
	}
	if decoded.Judged != 2 || decoded.Facets["genre"] != 2 {
		t.Fatalf("decoded = %+v", decoded)
	}
}
