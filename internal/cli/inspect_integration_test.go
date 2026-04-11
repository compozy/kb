//go:build integration

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	kgenerate "github.com/user/go-devstack/internal/generate"
	"github.com/user/go-devstack/internal/models"
)

func TestInspectCommandsAgainstGeneratedFixtureVault(t *testing.T) {
	fixtureRoot := filepath.Join("..", "generate", "testdata", "fixture-go-repo")
	vaultRoot := filepath.Join(t.TempDir(), "vault")

	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	defer slog.SetDefault(previous)

	summary, err := kgenerate.Generate(context.Background(), models.GenerateOptions{
		RootPath:  fixtureRoot,
		VaultPath: vaultRoot,
		TopicSlug: "fixture-go-repo",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	testCases := []struct {
		name          string
		args          []string
		minRows       int
		expectMessage string
	}{
		{name: "smells", args: []string{"inspect", "smells", "--format", "json", "--vault", summary.VaultPath}, minRows: 0},
		{name: "dead-code", args: []string{"inspect", "dead-code", "--format", "json", "--vault", summary.VaultPath}, minRows: 0},
		{name: "complexity", args: []string{"inspect", "complexity", "--format", "json", "--vault", summary.VaultPath}, minRows: 1},
		{name: "blast-radius", args: []string{"inspect", "blast-radius", "--format", "json", "--vault", summary.VaultPath}, minRows: 1},
		{name: "coupling", args: []string{"inspect", "coupling", "--format", "json", "--vault", summary.VaultPath}, minRows: 1},
		{name: "symbol", args: []string{"inspect", "symbol", "Hello", "--format", "json", "--vault", summary.VaultPath}, minRows: 1},
		{name: "file", args: []string{"inspect", "file", "internal/greeter/greeter.go", "--format", "json", "--vault", summary.VaultPath}, minRows: 1},
		{name: "backlinks", args: []string{"inspect", "backlinks", "Hello", "--format", "json", "--vault", summary.VaultPath}, minRows: 1},
		{name: "deps", args: []string{"inspect", "deps", "internal/greeter/greeter.go", "--format", "json", "--vault", summary.VaultPath}, minRows: 1},
		{name: "circular-deps", args: []string{"inspect", "circular-deps", "--format", "json", "--vault", summary.VaultPath}, minRows: 1, expectMessage: "no circular dependencies found"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			command := newRootCommand()
			var stdout bytes.Buffer
			command.SetOut(&stdout)
			command.SetErr(new(bytes.Buffer))
			command.SetArgs(testCase.args)

			if err := command.ExecuteContext(context.Background()); err != nil {
				t.Fatalf("ExecuteContext returned error: %v", err)
			}

			var decoded []map[string]any
			if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
				t.Fatalf("expected valid JSON output, got %v\n%s", err, stdout.String())
			}
			if len(decoded) < testCase.minRows {
				t.Fatalf("expected at least %d rows, got %d", testCase.minRows, len(decoded))
			}
			if testCase.expectMessage != "" {
				if len(decoded) != 1 || decoded[0]["message"] != testCase.expectMessage {
					t.Fatalf("expected message %q, got %#v", testCase.expectMessage, decoded)
				}
			}
		})
	}
}
