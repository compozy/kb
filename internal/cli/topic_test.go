package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestRootCommandUsesKBNameAndKnowledgeBaseDescription(t *testing.T) {
	command := newRootCommand()

	if command.Use != "kb" {
		t.Fatalf("root Use = %q, want kb", command.Use)
	}
	if !strings.Contains(command.Short, "knowledge") {
		t.Fatalf("root Short = %q, want knowledge-base description", command.Short)
	}
	if !strings.Contains(command.Long, "topic-based knowledge vaults") {
		t.Fatalf("root Long = %q, want topic-based KB description", command.Long)
	}
	if command.PersistentFlags().Lookup(rootVaultFlagName) == nil {
		t.Fatalf("expected persistent %q flag on root command", rootVaultFlagName)
	}
}

func TestTopicNewCommandPassesArgsAndPrintsJSON(t *testing.T) {
	originalNew := runTopicNew
	originalGetwd := topicGetwd
	t.Cleanup(func() {
		runTopicNew = originalNew
		topicGetwd = originalGetwd
	})

	var gotVault string
	var gotArgs []string
	runTopicNew = func(vaultPath, slug, title, domain string) (models.TopicInfo, error) {
		gotVault = vaultPath
		gotArgs = []string{slug, title, domain}
		return models.TopicInfo{
			Slug:     slug,
			Title:    title,
			Domain:   domain,
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	topicGetwd = func() (string, error) {
		return "/workspace/repo", nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"topic", "new", "distributed-systems", "Distributed Systems", "distributed", "--vault", "/tmp/vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotVault != "/tmp/vault" {
		t.Fatalf("vault path = %q, want /tmp/vault", gotVault)
	}
	if strings.Join(gotArgs, "|") != "distributed-systems|Distributed Systems|distributed" {
		t.Fatalf("topic new args = %#v", gotArgs)
	}

	var info models.TopicInfo
	if err := json.Unmarshal(stdout.Bytes(), &info); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if info.Slug != "distributed-systems" || info.Title != "Distributed Systems" || info.Domain != "distributed" {
		t.Fatalf("unexpected topic info payload: %#v", info)
	}
}

func TestTopicNewCommandRequiresThreeArgs(t *testing.T) {
	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"topic", "new", "distributed-systems"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected missing args error")
	}
	if !strings.Contains(err.Error(), "accepts 3 arg(s)") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTopicListCommandUsesVaultFlagAndFormatsOutput(t *testing.T) {
	originalList := runTopicList
	originalGetwd := topicGetwd
	t.Cleanup(func() {
		runTopicList = originalList
		topicGetwd = originalGetwd
	})

	var gotVault string
	runTopicList = func(vaultPath string) ([]models.TopicInfo, error) {
		gotVault = vaultPath
		return []models.TopicInfo{
			{
				Slug:         "distributed-systems",
				Title:        "Distributed Systems",
				Domain:       "distributed",
				ArticleCount: 4,
				SourceCount:  9,
			},
		}, nil
	}
	topicGetwd = func() (string, error) {
		return "/workspace/repo", nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"topic", "list", "--vault", "/tmp/vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotVault != "/tmp/vault" {
		t.Fatalf("vault path = %q, want /tmp/vault", gotVault)
	}

	rendered := stdout.String()
	for _, fragment := range []string{"slug", "title", "domain", "sources", "articles", "distributed-systems", "Distributed Systems", "9", "4"} {
		if !strings.Contains(rendered, fragment) {
			t.Fatalf("expected output to contain %q, got:\n%s", fragment, rendered)
		}
	}
}

func TestTopicListCommandDiscoversVaultFromCWD(t *testing.T) {
	originalList := runTopicList
	originalGetwd := topicGetwd
	t.Cleanup(func() {
		runTopicList = originalList
		topicGetwd = originalGetwd
	})

	repositoryRoot := t.TempDir()
	vaultPath := filepath.Join(repositoryRoot, ".kodebase", "vault")
	if err := os.MkdirAll(vaultPath, 0o755); err != nil {
		t.Fatalf("create vault path: %v", err)
	}

	nestedPath := filepath.Join(repositoryRoot, "workspace", "repo")
	if err := os.MkdirAll(nestedPath, 0o755); err != nil {
		t.Fatalf("create nested cwd: %v", err)
	}

	var gotVault string
	runTopicList = func(resolvedVaultPath string) ([]models.TopicInfo, error) {
		gotVault = resolvedVaultPath
		return nil, nil
	}
	topicGetwd = func() (string, error) {
		return nestedPath, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"topic", "list"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotVault != vaultPath {
		t.Fatalf("discovered vault path = %q, want %q", gotVault, vaultPath)
	}
}

func TestTopicInfoCommandRequiresOneArg(t *testing.T) {
	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"topic", "info"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected missing slug error")
	}
	if !strings.Contains(err.Error(), "accepts 1 arg(s)") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTopicInfoCommandPrintsJSON(t *testing.T) {
	originalInfo := runTopicInfo
	originalGetwd := topicGetwd
	t.Cleanup(func() {
		runTopicInfo = originalInfo
		topicGetwd = originalGetwd
	})

	var gotVault string
	var gotSlug string
	runTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		gotVault = vaultPath
		gotSlug = slug
		return models.TopicInfo{
			Slug:         slug,
			Title:        "Distributed Systems",
			Domain:       "distributed",
			RootPath:     filepath.Join(vaultPath, slug),
			ArticleCount: 4,
			SourceCount:  9,
		}, nil
	}
	topicGetwd = func() (string, error) {
		return "/workspace/repo", nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"topic", "info", "distributed-systems", "--vault", "/tmp/vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotVault != "/tmp/vault" {
		t.Fatalf("vault path = %q, want /tmp/vault", gotVault)
	}
	if gotSlug != "distributed-systems" {
		t.Fatalf("slug = %q, want distributed-systems", gotSlug)
	}

	var info models.TopicInfo
	if err := json.Unmarshal(stdout.Bytes(), &info); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if info.Slug != "distributed-systems" || info.SourceCount != 9 || info.ArticleCount != 4 {
		t.Fatalf("unexpected topic info payload: %#v", info)
	}
}
