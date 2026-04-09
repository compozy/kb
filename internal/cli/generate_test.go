package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestGenerateCommandPassesFlagsAndPrintsJSONSummary(t *testing.T) {
	original := runGenerate
	t.Cleanup(func() {
		runGenerate = original
	})

	var got models.GenerateOptions
	runGenerate = func(ctx context.Context, opts models.GenerateOptions) (models.GenerationSummary, error) {
		got = opts
		return models.GenerationSummary{
			Command:               "generate",
			TopicSlug:             "fixture",
			FilesScanned:          2,
			FilesParsed:           2,
			RawDocumentsWritten:   8,
			WikiDocumentsWritten:  10,
			IndexDocumentsWritten: 3,
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{
		"generate",
		"/tmp/repo",
		"--output", "/tmp/vault",
		"--topic", "fixture",
		"--title", "Fixture Repo",
		"--domain", "demo",
		"--include", "src/**/*.go",
		"--include", "web/**/*.ts",
		"--exclude", "vendor/**",
		"--semantic",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	expectedOptions := models.GenerateOptions{
		RootPath:        "/tmp/repo",
		OutputPath:      "/tmp/vault",
		Topic:           "fixture",
		Title:           "Fixture Repo",
		Domain:          "demo",
		IncludePatterns: []string{"src/**/*.go", "web/**/*.ts"},
		ExcludePatterns: []string{"vendor/**"},
		Semantic:        true,
	}
	if !reflect.DeepEqual(got, expectedOptions) {
		t.Fatalf("Generate options = %#v, want %#v", got, expectedOptions)
	}

	var summary models.GenerationSummary
	if err := json.Unmarshal(stdout.Bytes(), &summary); err != nil {
		t.Fatalf("stdout did not contain JSON summary: %v\n%s", err, stdout.String())
	}
	if summary.Command != "generate" || summary.TopicSlug != "fixture" {
		t.Fatalf("unexpected summary payload: %#v", summary)
	}
}
