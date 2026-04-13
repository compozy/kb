package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/qmd"
	"github.com/compozy/kb/internal/vault"
)

func TestIndexCommandResolvesTopicPathBeforeCallingQMD(t *testing.T) {
	originalClient := newIndexClient
	originalResolve := resolveIndexVaultQuery
	originalGetwd := indexGetwd
	t.Cleanup(func() {
		newIndexClient = originalClient
		resolveIndexVaultQuery = originalResolve
		indexGetwd = originalGetwd
	})

	var gotQuery vault.VaultQueryOptions
	var gotOptions qmd.IndexOptions
	indexGetwd = func() (string, error) {
		return "/workspace/repo", nil
	}
	resolveIndexVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		gotQuery = options
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/repo-topic",
			TopicSlug: "repo-topic",
		}, nil
	}
	newIndexClient = func() indexCommandClient {
		return fakeIndexClient{
			status: func(ctx context.Context) (qmd.IndexStatus, error) {
				return qmd.IndexStatus{}, nil
			},
			index: func(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error) {
				gotOptions = options
				return qmd.IndexResult{
					CollectionName: options.CollectionName,
					Status: qmd.IndexStatus{
						Collections: []qmd.CollectionInfo{
							{
								Name:      options.CollectionName,
								Documents: 12,
							},
						},
						TotalDocuments: 12,
					},
					UpdateResult: qmd.UpdateResult{
						Collections: 1,
						Indexed:     12,
					},
				}, nil
			},
		}
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"index"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotQuery != (vault.VaultQueryOptions{CWD: "/workspace/repo"}) {
		t.Fatalf("vault query = %#v, want cwd lookup", gotQuery)
	}
	if gotOptions.VaultPath != "/vault/repo-topic" {
		t.Fatalf("vault path = %q, want /vault/repo-topic", gotOptions.VaultPath)
	}
	if gotOptions.CollectionName != "repo-topic" {
		t.Fatalf("collection name = %q, want repo-topic", gotOptions.CollectionName)
	}
	if gotOptions.Operation != qmd.IndexOperationAdd {
		t.Fatalf("operation = %q, want add", gotOptions.Operation)
	}
	if gotOptions.Embed != true {
		t.Fatalf("embed = %t, want true", gotOptions.Embed)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if payload["collectionName"] != "repo-topic" {
		t.Fatalf("unexpected payload %#v", payload)
	}
}

func TestIndexCommandUpdatesExistingCollection(t *testing.T) {
	originalClient := newIndexClient
	originalResolve := resolveIndexVaultQuery
	originalGetwd := indexGetwd
	t.Cleanup(func() {
		newIndexClient = originalClient
		resolveIndexVaultQuery = originalResolve
		indexGetwd = originalGetwd
	})

	var gotOptions qmd.IndexOptions
	indexGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveIndexVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/repo-topic",
			TopicSlug: "repo-topic",
		}, nil
	}
	newIndexClient = func() indexCommandClient {
		return fakeIndexClient{
			status: func(ctx context.Context) (qmd.IndexStatus, error) {
				return qmd.IndexStatus{
					Collections: []qmd.CollectionInfo{{Name: "repo-topic"}},
				}, nil
			},
			index: func(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error) {
				gotOptions = options
				return qmd.IndexResult{CollectionName: options.CollectionName}, nil
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"index"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotOptions.Operation != qmd.IndexOperationUpdate {
		t.Fatalf("operation = %q, want update", gotOptions.Operation)
	}
}

func TestIndexCommandUsesTopicFlagForTopicScopedPathAndCollection(t *testing.T) {
	originalClient := newIndexClient
	originalResolve := resolveIndexVaultQuery
	originalGetwd := indexGetwd
	t.Cleanup(func() {
		newIndexClient = originalClient
		resolveIndexVaultQuery = originalResolve
		indexGetwd = originalGetwd
	})

	var gotQuery vault.VaultQueryOptions
	var gotOptions qmd.IndexOptions
	indexGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveIndexVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		gotQuery = options
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/systems-design",
			TopicSlug: "systems-design",
		}, nil
	}
	newIndexClient = func() indexCommandClient {
		return fakeIndexClient{
			status: func(ctx context.Context) (qmd.IndexStatus, error) {
				return qmd.IndexStatus{}, nil
			},
			index: func(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error) {
				gotOptions = options
				return qmd.IndexResult{CollectionName: options.CollectionName}, nil
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"index", "--topic", "systems-design"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if want := (vault.VaultQueryOptions{CWD: "/workspace/repo", Topic: "systems-design"}); gotQuery != want {
		t.Fatalf("vault query = %#v, want %#v", gotQuery, want)
	}
	if gotOptions.VaultPath != "/vault/systems-design" {
		t.Fatalf("vault path = %q, want /vault/systems-design", gotOptions.VaultPath)
	}
	if gotOptions.CollectionName != "systems-design" {
		t.Fatalf("collection name = %q, want systems-design", gotOptions.CollectionName)
	}
}

func TestIndexCommandHandlesQMDUnavailable(t *testing.T) {
	originalClient := newIndexClient
	originalResolve := resolveIndexVaultQuery
	originalGetwd := indexGetwd
	t.Cleanup(func() {
		newIndexClient = originalClient
		resolveIndexVaultQuery = originalResolve
		indexGetwd = originalGetwd
	})

	indexGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveIndexVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/repo-topic",
			TopicSlug: "repo-topic",
		}, nil
	}
	newIndexClient = func() indexCommandClient {
		return fakeIndexClient{
			status: func(ctx context.Context) (qmd.IndexStatus, error) {
				return qmd.IndexStatus{}, fmt.Errorf("wrapped: %w", qmd.ErrQMDUnavailable)
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"index"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected QMD unavailable error")
	}
	if !strings.Contains(err.Error(), "QMD is not available to kb") {
		t.Fatalf("unexpected error %q", err)
	}
	if !strings.Contains(err.Error(), qmd.InstallCommand) {
		t.Fatalf("expected install hint in %q", err)
	}
}

func TestIndexCommandHelpShowsFlags(t *testing.T) {
	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"index", "--help"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	for _, flag := range []string{"--embed", "--force-embed", "--context", "--name", "--vault", "--topic"} {
		if !strings.Contains(stdout.String(), flag) {
			t.Fatalf("expected help output to contain %q, got:\n%s", flag, stdout.String())
		}
	}
}

func TestIndexCommandIncludesEmbedFallbackMetadataInJSON(t *testing.T) {
	originalClient := newIndexClient
	originalResolve := resolveIndexVaultQuery
	originalGetwd := indexGetwd
	t.Cleanup(func() {
		newIndexClient = originalClient
		resolveIndexVaultQuery = originalResolve
		indexGetwd = originalGetwd
	})

	indexGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveIndexVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/repo-topic",
			TopicSlug: "repo-topic",
		}, nil
	}
	newIndexClient = func() indexCommandClient {
		return fakeIndexClient{
			status: func(ctx context.Context) (qmd.IndexStatus, error) {
				return qmd.IndexStatus{}, nil
			},
			index: func(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error) {
				return qmd.IndexResult{
					CollectionName: options.CollectionName,
					EmbedStatus:    qmd.EmbedStatusSkippedUnavailable,
					EmbedWarning:   "vector search is unavailable on this system; lexical indexing remains available",
					Status: qmd.IndexStatus{
						NeedsEmbedding: 3,
						TotalDocuments: 3,
					},
					UpdateResult: qmd.UpdateResult{
						Collections: 1,
						Indexed:     3,
					},
				}, nil
			},
		}
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"index"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if payload["embedStatus"] != string(qmd.EmbedStatusSkippedUnavailable) {
		t.Fatalf("embedStatus = %#v, want %q", payload["embedStatus"], qmd.EmbedStatusSkippedUnavailable)
	}
	if !strings.Contains(payload["embedWarning"].(string), "lexical indexing remains available") {
		t.Fatalf("embedWarning = %#v", payload["embedWarning"])
	}
}

type fakeIndexClient struct {
	status func(ctx context.Context) (qmd.IndexStatus, error)
	index  func(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error)
}

func (client fakeIndexClient) Status(ctx context.Context) (qmd.IndexStatus, error) {
	if client.status == nil {
		return qmd.IndexStatus{}, nil
	}
	return client.status(ctx)
}

func (client fakeIndexClient) Index(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error) {
	if client.index == nil {
		return qmd.IndexResult{}, nil
	}
	return client.index(ctx, options)
}
