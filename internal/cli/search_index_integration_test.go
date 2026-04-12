//go:build integration

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	kgenerate "github.com/user/kb/internal/generate"
	"github.com/user/kb/internal/models"
	"github.com/user/kb/internal/qmd"
)

func TestSearchCommandReturnsResultsAgainstIndexedVault(t *testing.T) {
	if _, err := exec.LookPath(qmd.DefaultBinaryPath); err != nil {
		t.Skip("qmd is not installed on PATH")
	}

	cacheRoot := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheRoot)
	t.Setenv("XDG_CONFIG_HOME", cacheRoot)
	t.Setenv("HOME", cacheRoot)

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

	collectionName := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))

	indexCommand := newRootCommand()
	indexCommand.SetOut(new(bytes.Buffer))
	indexCommand.SetErr(new(bytes.Buffer))
	indexCommand.SetArgs([]string{
		"index",
		"--vault", summary.VaultPath,
		"--topic", summary.TopicSlug,
		"--name", collectionName,
		"--embed=false",
		"--context", "CLI integration test",
	})

	if err := indexCommand.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("index command returned error: %v", err)
	}

	searchCommand := newRootCommand()
	var stdout bytes.Buffer
	searchCommand.SetOut(&stdout)
	searchCommand.SetErr(new(bytes.Buffer))
	searchCommand.SetArgs([]string{
		"search",
		"Hello",
		"--lex",
		"--collection", collectionName,
		"--format", "json",
	})

	if err := searchCommand.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("search command returned error: %v", err)
	}

	var results []map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &results); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if len(results) == 0 {
		t.Fatal("expected search results, got none")
	}
	if !strings.Contains(strings.ToLower(stdout.String()), "hello") {
		t.Fatalf("expected output to mention fixture content, got:\n%s", stdout.String())
	}
}
