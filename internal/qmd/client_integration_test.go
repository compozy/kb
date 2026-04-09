//go:build integration

package qmd

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestQMDClientIndexesAndSearchesTempVault(t *testing.T) {
	qmdPath, err := exec.LookPath(DefaultBinaryPath)
	if err != nil {
		t.Skip("qmd is not installed on PATH")
	}

	cacheRoot := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheRoot)
	t.Setenv("XDG_CONFIG_HOME", cacheRoot)
	t.Setenv("HOME", cacheRoot)

	vaultRoot := t.TempDir()
	writeMarkdownFile(t, vaultRoot, "auth.md", "# Authentication\n\nToken refresh flow documentation.\n")

	indexName := sanitizeTestIdentifier(t.Name())
	collectionName := indexName

	client := NewClient(
		WithBinaryPath(qmdPath),
		WithIndexName(indexName),
	)

	result, err := client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperationAdd,
		VaultPath:      vaultRoot,
		CollectionName: collectionName,
		Embed:          false,
	})
	if err != nil {
		t.Fatalf("Index returned error: %v", err)
	}
	if result.UpdateResult.Indexed == 0 {
		t.Fatalf("Index result = %#v, want indexed documents", result.UpdateResult)
	}

	results, err := client.Search(context.Background(), SearchOptions{
		Query:      "Token refresh flow",
		Mode:       SearchModeLexical,
		Collection: collectionName,
		Limit:      5,
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("Search returned no results, want at least one hit")
	}
	if results[0].Title != "Authentication" {
		t.Fatalf("first search result = %#v, want Authentication title", results[0])
	}
}

func writeMarkdownFile(t *testing.T, rootPath, name, body string) {
	t.Helper()

	path := rootPath + "/" + name
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) returned error: %v", path, err)
	}
}

func sanitizeTestIdentifier(value string) string {
	replacer := strings.NewReplacer("/", "-", " ", "-", "_", "-")
	return strings.ToLower(replacer.Replace(value))
}
