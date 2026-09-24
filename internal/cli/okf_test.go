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
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/models"
	kokf "github.com/compozy/kb/internal/okf"
	"github.com/compozy/kb/internal/review"
	ktopic "github.com/compozy/kb/internal/topic"
)

func TestPromoteCommandResolvesTargetAndPrintsJSON(t *testing.T) {
	t.Run("Should resolve target and print JSON", func(t *testing.T) {
		originalPromote := runPromote
		originalTopicInfo := runPromoteTopicInfo
		t.Cleanup(func() {
			runPromote = originalPromote
			runPromoteTopicInfo = originalTopicInfo
		})
		// A configured vocabulary makes promote open a decision session on
		// the source topic, so the vault and the decision model are real.
		vault := newPromoteVault(t)
		fake := fakes.NewOpenRouter(nil, nil)
		t.Cleanup(fake.Close)
		useDecisionConfig(t, fake, "[okf]\ntypes = [\"Playbook\"]\n")

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
			"--vault", vault,
		})

		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("ExecuteContext returned error: %v", err)
		}
		if gotInput.VaultPath != vault || gotInput.TargetTopic.Slug != "catalog" || gotInput.Type != "Playbook" {
			t.Fatalf("unexpected promote input: %#v", gotInput)
		}
		if gotInput.Suggester == nil || gotInput.TypeThreshold != 0.8 {
			t.Fatalf("type suggester = %#v, threshold = %v; want a decision-backed suggester at 0.8", gotInput.Suggester, gotInput.TypeThreshold)
		}
		if gotInput.SourceDocPath != "research/wiki/concepts/Alpha.md" {
			t.Fatalf("source doc path = %q, want research/wiki/concepts/Alpha.md", gotInput.SourceDocPath)
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
	})
}

func TestOKFCheckCommandRendersIssuesAndFailsOnErrors(t *testing.T) {
	t.Run("Should render issues and fail on errors", func(t *testing.T) {
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
	})
}

func TestOKFCheckCommandRejectsNonOKFTopic(t *testing.T) {
	t.Run("Should reject non-OKF topic", func(t *testing.T) {
		originalCheck := runOKFCheck
		originalTopicInfo := runOKFTopicInfo
		t.Cleanup(func() {
			runOKFCheck = originalCheck
			runOKFTopicInfo = originalTopicInfo
		})
		t.Setenv(kconfig.EnvConfigPath, writeCLIConfig(t, "[okf]\ntypes = [\"Playbook\"]\n"))

		checkCalled := false
		runOKFTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
			return models.TopicInfo{
				Slug:     slug,
				Mode:     models.TopicModeWiki,
				RootPath: filepath.Join(vaultPath, slug),
			}, nil
		}
		runOKFCheck = func(ctx context.Context, bundlePath string, options kokf.CheckOptions) ([]models.LintIssue, error) {
			checkCalled = true
			return nil, nil
		}

		command := newRootCommand()
		command.SetOut(new(bytes.Buffer))
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{"okf", "check", "research", "--vault", "/tmp/vault"})

		err := command.ExecuteContext(context.Background())
		if err == nil {
			t.Fatal("expected non-OKF topic rejection")
		}
		if err.Error() != `okf check: topic "research" is not an OKF topic` {
			t.Fatalf("error = %q, want non-OKF topic rejection", err.Error())
		}
		if checkCalled {
			t.Fatal("runOKFCheck was called for a non-OKF topic")
		}
	})
}

func TestPromoteCommandTypeSuggestion(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		noModel    bool
		args       []string
		answer     any
		wantType   string
		wantErr    string
		wantStderr string
		wantQueued bool
	}{
		{
			name:     "suggested type is used",
			config:   "[okf]\ntypes = [\"Playbook\", \"Reference\"]\n",
			answer:   fakes.Pick("Playbook", 0.9),
			wantType: "Playbook",
		},
		{
			name:       "explicit type that disagrees warns on stderr",
			config:     "[okf]\ntypes = [\"Playbook\", \"Reference\"]\n",
			args:       []string{"--type", "Reference"},
			answer:     fakes.Pick("Playbook", 0.9),
			wantType:   "Reference",
			wantStderr: `warning: the decision model suggests type "Playbook" (P=0.90), not "Reference"`,
		},
		{
			name:       "review band asks for --type",
			config:     "[okf]\ntypes = [\"Playbook\", \"Reference\"]\n",
			answer:     fakes.Dist(map[string]float64{"Playbook": 0.6, "Reference": 0.3, "none": 0.1}),
			wantErr:    "top: Playbook 0.60, Reference 0.30, none 0.10",
			wantQueued: true,
		},
		{
			name:    "no vocabulary needs --type",
			config:  "[okf]\ntypes = []\n",
			noModel: true,
			wantErr: "--type is required: [okf].types is empty",
		},
		{
			name:    "vocabulary needs the decision model",
			config:  "[okf]\ntypes = [\"Playbook\"]\n",
			noModel: true,
			args:    []string{"--type", "Playbook"},
			wantErr: "decision model is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vault := newPromoteVault(t)
			fake := fakes.NewOpenRouter(func(call fakes.Call, q fakes.Question) any { return tt.answer }, nil)
			t.Cleanup(fake.Close)
			if tt.noModel {
				t.Setenv(kconfig.EnvOpenRouterAPIKey, "")
				t.Setenv(kconfig.EnvConfigPath, writeCLIConfig(t, tt.config))
			} else {
				useDecisionConfig(t, fake, tt.config)
			}

			command := newRootCommand()
			var stdout, stderr bytes.Buffer
			command.SetOut(&stdout)
			command.SetErr(&stderr)
			command.SetArgs(append([]string{"promote", "research/wiki/concepts/Alpha.md", "--to", "catalog", "--vault", vault}, tt.args...))
			err := command.ExecuteContext(context.Background())

			pending, loadErr := review.Open(filepath.Join(vault, "research"), nil).Pending(review.QueueOKFType)
			if loadErr != nil {
				t.Fatalf("Pending: %v", loadErr)
			}
			if (len(pending) == 1) != tt.wantQueued {
				t.Fatalf("okf-type review items = %#v, want queued=%t", pending, tt.wantQueued)
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want %q", err, tt.wantErr)
				}
				if tt.noModel && len(fake.Calls()) != 0 {
					t.Fatalf("decision calls = %d, want none", len(fake.Calls()))
				}
				return
			}
			if err != nil {
				t.Fatalf("ExecuteContext: %v\nstderr: %s", err, stderr.String())
			}
			var result kokf.ConceptResult
			if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
				t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
			}
			if result.Type != tt.wantType || result.TypeSuggestion == nil {
				t.Fatalf("result = %#v, want type %q with a suggestion", result, tt.wantType)
			}
			if !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Fatalf("stderr = %q, want %q", stderr.String(), tt.wantStderr)
			}
			if len(fake.Calls()) != 1 {
				t.Fatalf("decision calls = %d, want 1", len(fake.Calls()))
			}
		})
	}
}

func TestOKFCheckCommandAdvisoryFindings(t *testing.T) {
	vault := newPromoteVault(t)
	fake := fakes.NewOpenRouter(func(call fakes.Call, q fakes.Question) any {
		if q.ID == "okf_type" {
			return fakes.Pick("Playbook", 0.9)
		}
		return fakes.Noul(0.9)
	}, nil)
	t.Cleanup(fake.Close)
	useDecisionConfig(t, fake, "[okf]\ntypes = [\"Playbook\", \"Reference\"]\n")
	concept := "---\ntitle: Rotate Key\ntype: Reference\ndescription: Rotating cut errors by 40%.\ntimestamp: 2026-06-27T10:11:12Z\n---\n1. Create a key. 2. Restart workers.\n"
	if err := os.WriteFile(filepath.Join(vault, "catalog", "rotate-key.md"), []byte(concept), 0o644); err != nil {
		t.Fatalf("write concept: %v", err)
	}

	run := func(extra ...string) []models.LintIssue {
		t.Helper()
		command := newRootCommand()
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetErr(new(bytes.Buffer))
		command.SetArgs(append([]string{"okf", "check", "catalog", "--strict", "--format", "json", "--vault", vault}, extra...))
		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("okf check %v: %v\n%s", extra, err, stdout.String())
		}
		var issues []models.LintIssue
		if err := json.Unmarshal(stdout.Bytes(), &issues); err != nil {
			t.Fatalf("decode issues: %v\n%s", err, stdout.String())
		}
		return issues
	}
	advisory := func(issues []models.LintIssue) int {
		count := 0
		for _, issue := range issues {
			if issue.Kind == models.LintIssueKindTypeMismatch || issue.Kind == models.LintIssueKindDescriptionUnsupported {
				if issue.Severity != models.SeverityInfo {
					t.Fatalf("advisory finding must be info: %#v", issue)
				}
				count++
			}
		}
		return count
	}

	if got := advisory(run()); got != 0 || len(fake.Calls()) != 0 {
		t.Fatalf("without --decide: %d findings, %d calls; want none of either", got, len(fake.Calls()))
	}
	if got := advisory(run("--decide")); got != 2 || len(fake.Calls()) != 2 {
		t.Fatalf("with --decide: %d findings, %d calls; want 2 and 2", got, len(fake.Calls()))
	}
	if got := advisory(run()); got != 2 || len(fake.Calls()) != 2 {
		t.Fatalf("from receipts: %d findings, %d calls; want 2 findings without new calls", got, len(fake.Calls()))
	}
}

// newPromoteVault creates a vault with a wiki topic "research" holding
// wiki/concepts/Alpha.md and an OKF topic "catalog".
func newPromoteVault(t *testing.T) string {
	t.Helper()
	vault := t.TempDir()
	if _, err := ktopic.New(vault, "research", "Research", "ops"); err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	if _, err := ktopic.NewWithMode(vault, "catalog", "Catalog", "ops", models.TopicModeOKF); err != nil {
		t.Fatalf("topic.NewWithMode: %v", err)
	}
	source := "---\ntitle: Alpha\ntype: wiki\nstage: compiled\n---\n1. Create a new key. 2. Restart the workers.\n"
	if err := os.WriteFile(filepath.Join(vault, "research", "wiki", "concepts", "Alpha.md"), []byte(source), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	return vault
}

// useDecisionConfig points the CLI config and env at a fake OpenRouter.
func useDecisionConfig(t *testing.T, fake *fakes.OpenRouter, extra string) {
	t.Helper()
	t.Setenv(kconfig.EnvOpenRouterAPIKey, "test-key")
	t.Setenv(kconfig.EnvOpenRouterAPIURL, fake.URL)
	t.Setenv(kconfig.EnvConfigPath, writeCLIConfig(t, extra))
}

func writeCLIConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "kb.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}
