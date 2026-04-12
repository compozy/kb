package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/user/kb/internal/qmd"
	"github.com/user/kb/internal/vault"
)

func TestSearchCommandDefaultsToHybridMode(t *testing.T) {
	originalClient := newSearchClient
	originalResolve := resolveSearchVaultQuery
	originalGetwd := searchGetwd
	t.Cleanup(func() {
		newSearchClient = originalClient
		resolveSearchVaultQuery = originalResolve
		searchGetwd = originalGetwd
	})

	var gotOptions qmd.SearchOptions
	var gotQuery vault.VaultQueryOptions
	searchGetwd = func() (string, error) {
		return "/workspace/repo", nil
	}
	resolveSearchVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		gotQuery = options
		return vault.ResolvedVault{
			VaultPath: "/vault",
			TopicPath: "/vault/repo-topic",
			TopicSlug: "repo-topic",
		}, nil
	}
	newSearchClient = func() searchCommandClient {
		return fakeSearchClient{
			search: func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
				gotOptions = options
				return []qmd.SearchResult{
					{
						Path:    "repo-topic/raw/codebase/symbols/branchy.md",
						Score:   0.91,
						Snippet: "Branchy snippet",
					},
				}, nil
			},
		}
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "branchy"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotOptions.Mode != qmd.SearchModeHybrid {
		t.Fatalf("search mode = %q, want hybrid", gotOptions.Mode)
	}
	if gotOptions.Collection != "repo-topic" {
		t.Fatalf("collection = %q, want repo-topic", gotOptions.Collection)
	}
	if gotOptions.Limit != 10 {
		t.Fatalf("limit = %d, want 10", gotOptions.Limit)
	}
	if gotOptions.Query != "branchy" {
		t.Fatalf("query = %q, want branchy", gotOptions.Query)
	}
	if gotOptions.MinScore != nil {
		t.Fatalf("minScore = %#v, want nil", gotOptions.MinScore)
	}
	if gotQuery != (vault.VaultQueryOptions{CWD: "/workspace/repo"}) {
		t.Fatalf("vault query = %#v, want cwd lookup", gotQuery)
	}
}

func TestSearchCommandUsesLexicalMode(t *testing.T) {
	originalClient := newSearchClient
	originalResolve := resolveSearchVaultQuery
	originalGetwd := searchGetwd
	t.Cleanup(func() {
		newSearchClient = originalClient
		resolveSearchVaultQuery = originalResolve
		searchGetwd = originalGetwd
	})

	var gotOptions qmd.SearchOptions
	searchGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveSearchVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{TopicSlug: "repo-topic"}, nil
	}
	newSearchClient = func() searchCommandClient {
		return fakeSearchClient{
			search: func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
				gotOptions = options
				return nil, nil
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "branchy", "--lex"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotOptions.Mode != qmd.SearchModeLexical {
		t.Fatalf("search mode = %q, want lexical", gotOptions.Mode)
	}
}

func TestSearchCommandUsesVectorMode(t *testing.T) {
	originalClient := newSearchClient
	originalResolve := resolveSearchVaultQuery
	originalGetwd := searchGetwd
	t.Cleanup(func() {
		newSearchClient = originalClient
		resolveSearchVaultQuery = originalResolve
		searchGetwd = originalGetwd
	})

	var gotOptions qmd.SearchOptions
	searchGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveSearchVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{TopicSlug: "repo-topic"}, nil
	}
	newSearchClient = func() searchCommandClient {
		return fakeSearchClient{
			search: func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
				gotOptions = options
				return nil, nil
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "branchy", "--vec"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotOptions.Mode != qmd.SearchModeVector {
		t.Fatalf("search mode = %q, want vector", gotOptions.Mode)
	}
}

func TestSearchCommandDisplaysPathScoreAndPreview(t *testing.T) {
	originalClient := newSearchClient
	originalResolve := resolveSearchVaultQuery
	originalGetwd := searchGetwd
	t.Cleanup(func() {
		newSearchClient = originalClient
		resolveSearchVaultQuery = originalResolve
		searchGetwd = originalGetwd
	})

	searchGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveSearchVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{TopicSlug: "repo-topic"}, nil
	}
	newSearchClient = func() searchCommandClient {
		return fakeSearchClient{
			search: func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
				return []qmd.SearchResult{
					{
						Path:    "repo-topic/raw/codebase/symbols/branchy.md",
						Score:   0.91,
						Snippet: "Branchy snippet",
					},
				}, nil
			},
		}
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "branchy"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	rendered := stdout.String()
	for _, fragment := range []string{"path", "score", "preview", "repo-topic/raw/codebase/symbols/branchy.md", "0.91", "Branchy snippet"} {
		if !strings.Contains(rendered, fragment) {
			t.Fatalf("expected output to contain %q, got:\n%s", fragment, rendered)
		}
	}
}

func TestSearchCommandPassesLimitFlag(t *testing.T) {
	originalClient := newSearchClient
	originalResolve := resolveSearchVaultQuery
	originalGetwd := searchGetwd
	t.Cleanup(func() {
		newSearchClient = originalClient
		resolveSearchVaultQuery = originalResolve
		searchGetwd = originalGetwd
	})

	var gotOptions qmd.SearchOptions
	searchGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveSearchVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{TopicSlug: "repo-topic"}, nil
	}
	newSearchClient = func() searchCommandClient {
		return fakeSearchClient{
			search: func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
				gotOptions = options
				return nil, nil
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "branchy", "--limit", "5"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotOptions.Limit != 5 {
		t.Fatalf("limit = %d, want 5", gotOptions.Limit)
	}
}

func TestSearchCommandUsesTopicFlagWhenDerivingCollection(t *testing.T) {
	originalClient := newSearchClient
	originalResolve := resolveSearchVaultQuery
	originalGetwd := searchGetwd
	t.Cleanup(func() {
		newSearchClient = originalClient
		resolveSearchVaultQuery = originalResolve
		searchGetwd = originalGetwd
	})

	var gotOptions qmd.SearchOptions
	var gotQuery vault.VaultQueryOptions
	searchGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveSearchVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		gotQuery = options
		return vault.ResolvedVault{TopicSlug: "systems-design"}, nil
	}
	newSearchClient = func() searchCommandClient {
		return fakeSearchClient{
			search: func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
				gotOptions = options
				return nil, nil
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "branchy", "--topic", "systems-design"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if want := (vault.VaultQueryOptions{CWD: "/workspace/repo", Topic: "systems-design"}); gotQuery != want {
		t.Fatalf("vault query = %#v, want %#v", gotQuery, want)
	}
	if gotOptions.Collection != "systems-design" {
		t.Fatalf("collection = %q, want systems-design", gotOptions.Collection)
	}
}

func TestSearchCommandHandlesQMDUnavailable(t *testing.T) {
	originalClient := newSearchClient
	originalResolve := resolveSearchVaultQuery
	originalGetwd := searchGetwd
	t.Cleanup(func() {
		newSearchClient = originalClient
		resolveSearchVaultQuery = originalResolve
		searchGetwd = originalGetwd
	})

	searchGetwd = func() (string, error) { return "/workspace/repo", nil }
	resolveSearchVaultQuery = func(options vault.VaultQueryOptions) (vault.ResolvedVault, error) {
		return vault.ResolvedVault{TopicSlug: "repo-topic"}, nil
	}
	newSearchClient = func() searchCommandClient {
		return fakeSearchClient{
			search: func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
				return nil, fmt.Errorf("boom: %w", qmd.ErrQMDUnavailable)
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "branchy"})

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

func TestSearchCommandHelpShowsFlags(t *testing.T) {
	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"search", "--help"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	for _, flag := range []string{"--lex", "--vec", "--limit", "--min-score", "--full", "--all", "--collection", "--topic", "--vault"} {
		if !strings.Contains(stdout.String(), flag) {
			t.Fatalf("expected help output to contain %q, got:\n%s", flag, stdout.String())
		}
	}
}

type fakeSearchClient struct {
	search func(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error)
}

func (client fakeSearchClient) Search(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
	return client.search(ctx, options)
}
