//go:build integration

package qmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// Quarantined sources and decision records never reach a search result:
// kb's collection mask keeps them out of the index, so a valid match ranked
// below many strongly matching quarantined files is still returned.
func TestQMDCollectionExcludesQuarantineAndDecisions(t *testing.T) {
	qmdPath, err := exec.LookPath(DefaultBinaryPath)
	if err != nil {
		t.Skip("qmd is not installed on PATH")
	}

	cacheRoot := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheRoot)
	t.Setenv("XDG_CONFIG_HOME", cacheRoot)
	t.Setenv("HOME", cacheRoot)

	topicRoot := t.TempDir()
	for _, dir := range []string{"raw/_quarantine", "raw/articles", ".decisions"} {
		if err := os.MkdirAll(filepath.Join(topicRoot, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for index := range 12 {
		writeMarkdownFile(t, topicRoot, fmt.Sprintf("raw/_quarantine/junk-%02d.md", index),
			"# Zebrafish zebrafish\n\nZebrafish zebrafish zebrafish zebrafish zebrafish.\n")
	}
	writeMarkdownFile(t, topicRoot, ".decisions/record.md", "# Zebrafish\n\nZebrafish zebrafish zebrafish.\n")
	writeMarkdownFile(t, topicRoot, "raw/articles/valid.md", "# Fish models\n\nA short note on the zebrafish as a model organism, among other topics.\n")

	indexName := sanitizeTestIdentifier(t.Name())
	client := NewClient(WithBinaryPath(qmdPath), WithIndexName(indexName))
	result, err := client.Index(context.Background(), IndexOptions{
		Operation:      IndexOperationAdd,
		VaultPath:      topicRoot,
		CollectionName: indexName,
	})
	if err != nil {
		t.Fatalf("Index returned error: %v", err)
	}
	if result.UpdateResult.Indexed != 1 {
		t.Fatalf("indexed = %d, want only the valid source (%#v)", result.UpdateResult.Indexed, result.UpdateResult)
	}

	// A quarantine made after the collection exists is still excluded on update.
	writeMarkdownFile(t, topicRoot, "raw/_quarantine/late.md", "# Zebrafish\n\nZebrafish zebrafish zebrafish.\n")
	if _, err := client.Index(context.Background(), IndexOptions{Operation: IndexOperationUpdate, CollectionName: indexName}); err != nil {
		t.Fatalf("Index update returned error: %v", err)
	}

	results, err := client.Search(context.Background(), SearchOptions{
		Query:      "zebrafish",
		Mode:       SearchModeLexical,
		Collection: indexName,
		Limit:      1,
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 || !strings.Contains(results[0].Path, "/raw/articles/valid.md") {
		t.Fatalf("results = %#v, want only the valid source", results)
	}

	// A collection created before the mask existed still indexes the
	// quarantine; over-fetching and filtering keep the valid match within
	// the limit.
	legacy := indexName + "-legacy"
	if _, _, err := client.run(context.Background(), commandSpec{
		label: "collection add (legacy)",
		args:  client.baseArgs("collection", "add", topicRoot, "--name", legacy),
	}); err != nil {
		t.Fatalf("legacy collection add: %v", err)
	}
	results, err = client.Search(context.Background(), SearchOptions{
		Query:      "zebrafish",
		Mode:       SearchModeLexical,
		Collection: legacy,
		Limit:      1,
	})
	if err != nil {
		t.Fatalf("legacy Search returned error: %v", err)
	}
	if len(results) != 1 || !strings.Contains(results[0].Path, "/raw/articles/valid.md") {
		t.Fatalf("legacy results = %#v, want the valid source despite 13 quarantined hits ranked ahead", results)
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
