package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/models"
	kokf "github.com/compozy/kb/internal/okf"
)

func TestPromoteCommandResolvesTargetAndPrintsJSON(t *testing.T) {
	originalPromote := runPromote
	originalTopicInfo := runPromoteTopicInfo
	t.Cleanup(func() {
		runPromote = originalPromote
		runPromoteTopicInfo = originalTopicInfo
	})
	t.Setenv(kconfig.EnvConfigPath, writeCLIConfig(t, "[okf]\ntypes = [\"Playbook\"]\n"))

	var gotInput kokf.PromoteInput
	runPromoteTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{
			Slug:     slug,
			Mode:     models.TopicModeOKF,
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	runPromote = func(ctx context.Context, input kokf.PromoteInput) (kokf.ConceptResult, error) {
		gotInput = input
		return kokf.ConceptResult{
			WrittenPath:    "alpha.md",
			Type:           input.Type,
			LinksRewritten: 1,
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"promote", "research/wiki/concepts/Alpha.md",
		"--to", "catalog",
		"--type", "Playbook",
		"--description", "Alpha description.",
		"--vault", "/tmp/vault",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}
	if gotInput.VaultPath != "/tmp/vault" || gotInput.TargetTopic.Slug != "catalog" || gotInput.Type != "Playbook" {
		t.Fatalf("unexpected promote input: %#v", gotInput)
	}
	if gotInput.Description != "Alpha description." {
		t.Fatalf("description = %q", gotInput.Description)
	}
	if len(gotInput.Types) != 1 || gotInput.Types[0] != "Playbook" {
		t.Fatalf("types = %#v, want Playbook", gotInput.Types)
	}

	var result kokf.ConceptResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.WrittenPath != "alpha.md" || result.Type != "Playbook" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestOKFCheckCommandRendersIssuesAndFailsOnErrors(t *testing.T) {
	originalCheck := runOKFCheck
	originalTopicInfo := runOKFTopicInfo
	t.Cleanup(func() {
		runOKFCheck = originalCheck
		runOKFTopicInfo = originalTopicInfo
	})
	t.Setenv(kconfig.EnvConfigPath, writeCLIConfig(t, "[okf]\ntypes = [\"Playbook\"]\n"))

	runOKFTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{
			Slug:     slug,
			Mode:     models.TopicModeOKF,
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	runOKFCheck = func(ctx context.Context, bundlePath string, options kokf.CheckOptions) ([]models.LintIssue, error) {
		if bundlePath != "/tmp/vault/catalog" {
			return nil, fmt.Errorf("bundle path = %q", bundlePath)
		}
		if !options.Strict || len(options.Types) != 1 || options.Types[0] != "Playbook" {
			return nil, fmt.Errorf("unexpected options: %#v", options)
		}
		return []models.LintIssue{{
			Kind:     models.LintIssueKindFormat,
			Severity: models.SeverityError,
			FilePath: "bad.md",
			Target:   "type",
			Message:  "missing type",
		}}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"okf", "check", "catalog", "--strict", "--format", "json", "--vault", "/tmp/vault"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected okf check to fail on error issues")
	}
	if !strings.Contains(err.Error(), "found 1 issue") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout.String(), `"filePath": "bad.md"`) {
		t.Fatalf("stdout missing issue JSON:\n%s", stdout.String())
	}
}

func writeCLIConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "kb.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}
