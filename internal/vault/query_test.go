package vault_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/kb/internal/vault"
)

func TestResolveVaultQueryFindsVaultByWalkingUp(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()
	repositoryRoot := filepath.Join(workspaceRoot, "repo")
	nestedCWD := filepath.Join(repositoryRoot, "packages", "cli", "src")
	vaultPath := filepath.Join(repositoryRoot, ".kb", "vault")
	topicPath := filepath.Join(vaultPath, "demo-topic")

	mkdirAll(t, topicPath)
	mkdirAll(t, nestedCWD)
	writeTopicMarker(t, topicPath)

	resolvedVault, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{CWD: nestedCWD})
	if err != nil {
		t.Fatalf("ResolveVaultQuery returned error: %v", err)
	}

	if resolvedVault.VaultPath != vaultPath {
		t.Fatalf("vault path = %q, want %q", resolvedVault.VaultPath, vaultPath)
	}
	if resolvedVault.TopicPath != topicPath {
		t.Fatalf("topic path = %q, want %q", resolvedVault.TopicPath, topicPath)
	}
	if resolvedVault.TopicSlug != "demo-topic" {
		t.Fatalf("topic slug = %q, want demo-topic", resolvedVault.TopicSlug)
	}
}

func TestResolveVaultQueryPrefersExplicitVault(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()
	cwd := filepath.Join(workspaceRoot, "outside")
	vaultPath := filepath.Join(workspaceRoot, "custom-vault")
	topicPath := filepath.Join(vaultPath, "chosen-topic")

	mkdirAll(t, cwd)
	mkdirAll(t, topicPath)
	writeTopicMarker(t, topicPath)

	resolvedVault, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{
		CWD:   cwd,
		Vault: vaultPath,
	})
	if err != nil {
		t.Fatalf("ResolveVaultQuery returned error: %v", err)
	}

	if resolvedVault.VaultPath != vaultPath {
		t.Fatalf("vault path = %q, want %q", resolvedVault.VaultPath, vaultPath)
	}
	if resolvedVault.TopicPath != topicPath {
		t.Fatalf("topic path = %q, want %q", resolvedVault.TopicPath, topicPath)
	}
	if resolvedVault.TopicSlug != "chosen-topic" {
		t.Fatalf("topic slug = %q, want chosen-topic", resolvedVault.TopicSlug)
	}
}

func TestResolveVaultQueryAutoResolvesSingleTopic(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()
	vaultPath := filepath.Join(workspaceRoot, ".kb", "vault")
	topicPath := filepath.Join(vaultPath, "single-topic")

	mkdirAll(t, topicPath)
	writeTopicMarker(t, topicPath)

	resolvedVault, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{Vault: vaultPath})
	if err != nil {
		t.Fatalf("ResolveVaultQuery returned error: %v", err)
	}

	if resolvedVault.TopicPath != topicPath {
		t.Fatalf("topic path = %q, want %q", resolvedVault.TopicPath, topicPath)
	}
	if resolvedVault.TopicSlug != "single-topic" {
		t.Fatalf("topic slug = %q, want single-topic", resolvedVault.TopicSlug)
	}
}

func TestResolveVaultQueryErrorsWhenMultipleTopicsExist(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()
	vaultPath := filepath.Join(workspaceRoot, ".kb", "vault")

	for _, topic := range []string{"alpha", "beta"} {
		topicPath := filepath.Join(vaultPath, topic)
		mkdirAll(t, topicPath)
		writeTopicMarker(t, topicPath)
	}

	_, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{Vault: vaultPath})
	if err == nil {
		t.Fatal("expected ResolveVaultQuery to fail when multiple topics exist")
	}
	if !strings.Contains(err.Error(), "multiple topics were found") {
		t.Fatalf("unexpected error message %q", err)
	}
	if !strings.Contains(err.Error(), "alpha, beta") {
		t.Fatalf("expected available topics in error message, got %q", err)
	}
}

func TestResolveVaultQueryUsesExplicitTopic(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()
	vaultPath := filepath.Join(workspaceRoot, ".kb", "vault")
	topicPath := filepath.Join(vaultPath, "explicit-topic")

	mkdirAll(t, topicPath)

	resolvedVault, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{
		Vault: vaultPath,
		Topic: "explicit-topic",
	})
	if err != nil {
		t.Fatalf("ResolveVaultQuery returned error: %v", err)
	}

	if resolvedVault.TopicPath != topicPath {
		t.Fatalf("topic path = %q, want %q", resolvedVault.TopicPath, topicPath)
	}
	if resolvedVault.TopicSlug != "explicit-topic" {
		t.Fatalf("topic slug = %q, want explicit-topic", resolvedVault.TopicSlug)
	}
}

func TestResolveVaultQueryErrorsWhenExplicitTopicIsMissing(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()
	vaultPath := filepath.Join(workspaceRoot, ".kb", "vault")

	mkdirAll(t, vaultPath)

	_, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{
		Vault: vaultPath,
		Topic: "missing-topic",
	})
	if err == nil {
		t.Fatal("expected ResolveVaultQuery to fail for a missing explicit topic")
	}
	if !strings.Contains(err.Error(), "Topic path was not found or is not a directory") {
		t.Fatalf("unexpected error message %q", err)
	}
}

func TestResolveVaultQueryErrorsClearlyWhenNoVaultIsFound(t *testing.T) {
	t.Parallel()

	cwd := t.TempDir()

	_, err := vault.ResolveVaultQuery(vault.VaultQueryOptions{CWD: cwd})
	if err == nil {
		t.Fatal("expected ResolveVaultQuery to fail when no vault exists")
	}
	if !strings.Contains(err.Error(), "unable to find a vault") {
		t.Fatalf("unexpected error message %q", err)
	}
}

func TestListAvailableTopicsReturnsSortedTopics(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()
	vaultPath := filepath.Join(workspaceRoot, ".kb", "vault")

	for _, topic := range []string{"zeta", "alpha"} {
		topicPath := filepath.Join(vaultPath, topic)
		mkdirAll(t, topicPath)
		writeTopicMarker(t, topicPath)
	}

	topics, err := vault.ListAvailableTopics(vault.VaultQueryOptions{Vault: vaultPath})
	if err != nil {
		t.Fatalf("ListAvailableTopics returned error: %v", err)
	}

	want := []string{"alpha", "zeta"}
	if len(topics) != len(want) {
		t.Fatalf("expected %d topics, got %d", len(want), len(topics))
	}
	for index := range want {
		if topics[index] != want[index] {
			t.Fatalf("topic %d = %q, want %q", index, topics[index], want[index])
		}
	}
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func writeTopicMarker(t *testing.T, topicPath string) {
	t.Helper()

	markerPath := filepath.Join(topicPath, "CLAUDE.md")
	if err := os.WriteFile(markerPath, []byte("# Topic\n"), 0o644); err != nil {
		t.Fatalf("write topic marker: %v", err)
	}
}
