//go:build integration

package vault_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/kb/internal/vault"
)

func TestWriteVaultIntegrationPersistsFullRenderedVault(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)

	result, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	})
	if err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	if result != countKinds(documents) {
		t.Fatalf("write result = %#v, want %#v", result, countKinds(documents))
	}

	for _, document := range documents {
		assertFileExists(t, filepath.Join(topic.TopicPath, filepath.FromSlash(document.RelativePath)))
	}
	for _, baseFile := range baseFiles {
		assertFileExists(t, filepath.Join(topic.TopicPath, filepath.FromSlash(baseFile.RelativePath)))
	}

	assertFileExists(t, filepath.Join(topic.TopicPath, "CLAUDE.md"))
	assertFileExists(t, filepath.Join(topic.TopicPath, "log.md"))

	entries, err := os.ReadDir(filepath.Join(topic.TopicPath, "wiki", "concepts"))
	if err != nil {
		t.Fatalf("read concept directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected generated concept files in integration write")
	}
}
