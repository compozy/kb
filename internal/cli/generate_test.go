package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	kgenerate "github.com/user/kb/internal/generate"
	"github.com/user/kb/internal/models"
)

func TestGenerateCommandPassesFlagsAndPrintsJSONSummary(t *testing.T) {
	original := runGenerate
	t.Cleanup(func() {
		runGenerate = original
	})

	var got models.GenerateOptions
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
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
		"--vault", "/tmp/vault",
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
		VaultPath:       "/tmp/vault",
		TopicSlug:       "fixture",
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

func TestGenerateCommandWritesTextEventsToStderrByDefault(t *testing.T) {
	original := runGenerate
	t.Cleanup(func() {
		runGenerate = original
	})

	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		observer.ObserveGenerateEvent(ctx, kgenerate.Event{Kind: kgenerate.EventStageStarted, Stage: "scan"})
		observer.ObserveGenerateEvent(ctx, kgenerate.Event{Kind: kgenerate.EventStageCompleted, Stage: "scan", DurationMillis: 7})
		return models.GenerationSummary{Command: "generate", TopicSlug: "fixture"}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{"generate", "/tmp/repo"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if !strings.Contains(stderr.String(), "generate: scan started") {
		t.Fatalf("expected text stage start on stderr, got:\n%s", stderr.String())
	}
	if !strings.Contains(stderr.String(), "generate: scan completed in 7ms") {
		t.Fatalf("expected text stage completion on stderr, got:\n%s", stderr.String())
	}
}

func TestGenerateCommandWritesJSONEventsWhenRequested(t *testing.T) {
	original := runGenerate
	t.Cleanup(func() {
		runGenerate = original
	})

	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		observer.ObserveGenerateEvent(ctx, kgenerate.Event{
			Kind:   kgenerate.EventStageStarted,
			Stage:  "parse",
			Total:  2,
			Fields: map[string]any{"adapter_count": 1},
		})
		observer.ObserveGenerateEvent(ctx, kgenerate.Event{
			Kind:      kgenerate.EventStageProgress,
			Stage:     "parse",
			Completed: 1,
			Total:     2,
			Fields:    map[string]any{"path": "main.go"},
		})
		return models.GenerationSummary{Command: "generate", TopicSlug: "fixture"}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{"generate", "/tmp/repo", "--log-format", "json"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(stderr.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected two JSON event lines, got %d\n%s", len(lines), stderr.String())
	}

	var started kgenerate.Event
	if err := json.Unmarshal([]byte(lines[0]), &started); err != nil {
		t.Fatalf("stderr line 1 was not JSON: %v\n%s", err, lines[0])
	}
	if started.Kind != kgenerate.EventStageStarted || started.Stage != "parse" || started.Total != 2 {
		t.Fatalf("unexpected started event: %#v", started)
	}

	var progress kgenerate.Event
	if err := json.Unmarshal([]byte(lines[1]), &progress); err != nil {
		t.Fatalf("stderr line 2 was not JSON: %v\n%s", err, lines[1])
	}
	if progress.Kind != kgenerate.EventStageProgress || progress.Completed != 1 || progress.Total != 2 {
		t.Fatalf("unexpected progress event: %#v", progress)
	}
}

func TestGenerateCommandSupportsDeprecatedOutputAlias(t *testing.T) {
	original := runGenerate
	t.Cleanup(func() {
		runGenerate = original
	})

	var got models.GenerateOptions
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		got = opts
		return models.GenerationSummary{Command: "generate", TopicSlug: "fixture"}, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"generate", "/tmp/repo", "--output", "/tmp/legacy-vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if got.VaultPath != "/tmp/legacy-vault" {
		t.Fatalf("VaultPath = %q, want /tmp/legacy-vault", got.VaultPath)
	}
}

func TestGenerateTextObserverReportsCompletedCountsAndFailures(t *testing.T) {
	var stderr bytes.Buffer

	observer := &generateTextObserver{
		writer:       &stderr,
		liveProgress: false,
	}

	observer.ObserveGenerateEvent(context.Background(), kgenerate.Event{
		Kind:  kgenerate.EventStageStarted,
		Stage: "write",
		Total: 4,
	})
	observer.ObserveGenerateEvent(context.Background(), kgenerate.Event{
		Kind:      kgenerate.EventStageProgress,
		Stage:     "write",
		Completed: 2,
		Total:     4,
	})
	observer.ObserveGenerateEvent(context.Background(), kgenerate.Event{
		Kind:           kgenerate.EventStageCompleted,
		Stage:          "write",
		Completed:      4,
		Total:          4,
		DurationMillis: 12,
	})
	observer.ObserveGenerateEvent(context.Background(), kgenerate.Event{
		Kind:           kgenerate.EventStageFailed,
		Stage:          "select_adapters",
		DurationMillis: 5,
		Error:          "boom",
	})

	output := stderr.String()
	if !strings.Contains(output, "generate: write started") {
		t.Fatalf("expected write start message, got:\n%s", output)
	}
	if !strings.Contains(output, "generate: write completed in 12ms (4/4)") {
		t.Fatalf("expected write completion counts, got:\n%s", output)
	}
	if !strings.Contains(output, "generate: select adapters failed after 5ms: boom") {
		t.Fatalf("expected failure message, got:\n%s", output)
	}
}

func TestGenerateCommandRejectsInvalidProgressMode(t *testing.T) {
	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"generate", "/tmp/repo", "--progress", "loud"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected invalid progress mode error")
	}
	if !strings.Contains(err.Error(), "unsupported progress mode") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateCommandRejectsInvalidLogFormat(t *testing.T) {
	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"generate", "/tmp/repo", "--log-format", "pretty"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected invalid log format error")
	}
	if !strings.Contains(err.Error(), "unsupported log format") {
		t.Fatalf("unexpected error: %v", err)
	}
}
