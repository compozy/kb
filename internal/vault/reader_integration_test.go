//go:build integration

package vault_test

import (
	"context"
	"testing"

	"github.com/compozy/kb/internal/vault"
)

func TestReadVaultSnapshotRoundTripsWriterOutput(t *testing.T) {
	t.Parallel()

	topic, graph, documents, baseFiles := testWriteVaultInputs(t)

	if _, err := vault.WriteVault(context.Background(), vault.WriteVaultOptions{
		Topic:     topic,
		Graph:     graph,
		Documents: documents,
		BaseFiles: baseFiles,
	}); err != nil {
		t.Fatalf("WriteVault returned error: %v", err)
	}

	snapshot, err := vault.ReadVaultSnapshot(vault.ResolvedVault{
		VaultPath: topic.VaultPath,
		TopicPath: topic.TopicPath,
		TopicSlug: topic.Slug,
	}, vault.ReadVaultOptions{})
	if err != nil {
		t.Fatalf("ReadVaultSnapshot returned error: %v", err)
	}

	if len(snapshot.Files) != 4 {
		t.Fatalf("expected 4 file documents, got %d", len(snapshot.Files))
	}
	if len(snapshot.Symbols) != 4 {
		t.Fatalf("expected 4 symbol documents, got %d", len(snapshot.Symbols))
	}
	if len(snapshot.Directories) != 2 {
		t.Fatalf("expected 2 directory documents, got %d", len(snapshot.Directories))
	}
	if len(snapshot.Wikis) != 14 {
		t.Fatalf("expected 14 wiki/default documents, got %d", len(snapshot.Wikis))
	}

	matches := vault.FindSymbolsByName(snapshot, "alpha")
	if len(matches) != 1 {
		t.Fatalf("expected 1 alpha symbol match, got %d", len(matches))
	}
	if len(matches[0].OutgoingRelations) == 0 {
		t.Fatal("expected alpha symbol to retain parsed outgoing relations after round-trip")
	}
	if len(matches[0].Backlinks) == 0 {
		t.Fatal("expected alpha symbol to retain parsed backlinks after round-trip")
	}
}
