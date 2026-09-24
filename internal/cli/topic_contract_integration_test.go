//go:build integration

package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/review"
)

// runContractCLI runs the root command with stdin and returns both streams
// and the error (it does not fail the test on error).
func runContractCLI(t *testing.T, stdin string, args ...string) (string, string, error) {
	t.Helper()
	command := newRootCommand()
	var stdout, stderr bytes.Buffer
	command.SetIn(strings.NewReader(stdin))
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs(args)
	err := command.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

// holdoutSubject returns the first raw/articles/labeled-<n>.md path whose
// calibration bucket is holdout (never dev by construction).
func holdoutSubject(t *testing.T) string {
	t.Helper()
	for n := range 100 {
		subject := fmt.Sprintf("raw/articles/labeled-%d.md", n)
		if !review.DevBucket(subject) {
			return subject
		}
	}
	t.Fatal("no holdout subject in 100 tries")
	return ""
}

// TestCLIIntegrationTopicContractImpactPreview runs spec §17 *impact
// preview* through the real commands and session: import from CLAUDE.md,
// a self-check conflict that blocks without --force, an --accept that
// prints the preview and is not activated without the typed slug, an
// --accept --yes that activates, a following `kb classify` that reuses the
// preview receipts (no new role or quality questions), and a contract change
// that demotes the previewed labels to dev.
func TestCLIIntegrationTopicContractImpactPreview(t *testing.T) {
	var conflicts atomic.Bool
	conflicts.Store(true)
	fake := fakes.NewOpenRouter(func(call fakes.Call, q fakes.Question) any {
		switch {
		case strings.HasPrefix(q.ID, "conflict_"):
			if conflicts.Load() {
				return fakes.Noul(0.9)
			}
			return fakes.Noul(0.1)
		case q.ID == "role":
			if strings.Contains(call.StateString(), "Offtopic") {
				return fakes.Pick("off_topic", 0.95)
			}
			return fakes.Pick("core", 0.9)
		case q.ID == "kind":
			return fakes.Pick("article_or_essay", 0.9)
		case q.ID == "enough_text_to_judge":
			return fakes.Noul(0.9)
		}
		return nil
	}, func(schemaName, _, _ string) any {
		if schemaName == "document_literals" {
			return map[string]any{"summary": "Notes about typed decision models.", "entities": []string{}, "questions": []string{"What is a typed decision?"}}
		}
		return map[string]any{}
	})
	t.Cleanup(fake.Close)

	configPath := filepath.Join(t.TempDir(), "kb.toml")
	writeFile(t, configPath, "")
	t.Setenv("APP_CONFIG", configPath)
	t.Setenv("OPENROUTER_API_KEY", "test-key")
	t.Setenv("OPENROUTER_API_URL", fake.URL)

	vaultRoot := t.TempDir()
	info := scaffoldTopicForIntegration(t, vaultRoot, "demo", "Demo", "demo")
	root := info.RootPath
	labeled := holdoutSubject(t)
	sources := []string{"raw/articles/core-1.md", "raw/articles/core-2.md", "raw/articles/offtopic.md", "raw/audit/page.md", labeled}
	for _, rel := range sources {
		title := "Typed decision notes " + filepath.Base(rel)
		if strings.Contains(rel, "offtopic") || rel == labeled {
			title = "Offtopic sports results " + filepath.Base(rel)
		}
		writeMarkdownDocument(t, root, rel, map[string]any{"title": title, "type": "source", "stage": "raw", "source_kind": "document"},
			"# "+title+"\n\nThe document discusses its subject with enough words to judge it.\n")
	}
	store := review.Open(root, nil)
	for subject, verdict := range map[string]string{labeled: review.VerdictNegative, "raw/articles/core-1.md": review.VerdictPositive} {
		if err := store.AddLabel(review.Label{Subject: subject, Purpose: review.PurposeRelevance, Verdict: verdict, Origin: "import:screening.jsonl@abc"}); err != nil {
			t.Fatal(err)
		}
	}

	// --import-claude: the CLAUDE.md section becomes contract_draft.
	claudePath := filepath.Join(root, "CLAUDE.md")
	claude, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}
	section := "## Selection contract\n\n- **Purpose:** Collect material on typed decision models.\n- **Core:** Typed decision models\n- **Adjacent (keep):** LLM judges, kept as comparison\n- **Out of scope:** Sports results with no decision-model content\n"
	writeFile(t, claudePath, strings.Replace(string(claude), contract.RenderSection(nil), section, 1))
	stdout, _ := runCLIWithStreams(t, "topic", "contract", "demo", "--import-claude", "--vault", vaultRoot)
	if !strings.Contains(stdout, "contract_draft:") || !strings.Contains(stdout, "- Typed decision models") {
		t.Fatalf("--import-claude output:\n%s", stdout)
	}
	settings, err := contract.LoadSettings(root)
	if err != nil {
		t.Fatal(err)
	}
	draftHash := settings.ContractDraft.Hash()
	if draftHash == "" || !settings.Contract.Empty() {
		t.Fatalf("after import: %+v", settings)
	}

	// Self-check conflicts block without --force, before any preview call.
	stdout, stderr, err := runContractCLI(t, "", "topic", "contract", "demo", "--accept", "--yes", "--vault", vaultRoot)
	if err == nil || !strings.Contains(err.Error(), "self-check") {
		t.Fatalf("--accept with conflicts: err = %v\n%s", err, stderr)
	}
	if !strings.Contains(stdout, "self-check: 2 pairs, 2 conflicts") || fake.QuestionCount("role") != 0 {
		t.Fatalf("conflict output (role questions %d):\n%s", fake.QuestionCount("role"), stdout)
	}

	// The owner qualifies the out_of_scope line (a new draft, so the cached
	// self-check answers no longer apply). Without the typed slug the
	// preview runs but nothing is activated.
	conflicts.Store(false)
	edited := settings.ContractDraft.Normalized()
	edited.OutOfScope = []string{"Sports results with no decision-model or judging content"}
	if err := contract.SaveContractDraft(root, &edited); err != nil {
		t.Fatal(err)
	}
	draftHash = edited.Hash()
	stdout, stderr, err = runContractCLI(t, "wrong\n", "topic", "contract", "demo", "--accept", "--vault", vaultRoot)
	if err == nil || !strings.Contains(err.Error(), "not activated") {
		t.Fatalf("--accept with a wrong confirmation: err = %v\n%s", err, stderr)
	}
	for _, want := range []string{
		"self-check: 2 pairs, 0 conflicts",
		"impact preview over all 5 sources",
		"quarantine-would-be",
		"per raw/ folder:",
		"raw/audit",
		"highest P(off_topic) (5):",
		"raw/articles/offtopic.md",
		"against 2 relevance labels (1 negative): would-be quarantine precision 1.00 (1/1), recall 1.00 (1/1)",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("preview output lacks %q:\n%s", want, stdout)
		}
	}
	if !strings.Contains(stderr, "impact preview: judging 5 sources, estimated US$ 0.0015") || !strings.Contains(stderr, `Type the topic slug "demo"`) {
		t.Fatalf("stderr:\n%s", stderr)
	}
	if settings, _ = contract.LoadSettings(root); !settings.Contract.Empty() {
		t.Fatal("contract activated without confirmation")
	}
	roleQuestions := fake.QuestionCount("role")
	if roleQuestions != len(sources) {
		t.Fatalf("role questions = %d, want %d", roleQuestions, len(sources))
	}

	// --yes activates; the preview is served from the receipts.
	stdout, stderr, err = runContractCLI(t, "", "topic", "contract", "demo", "--accept", "--yes", "--vault", vaultRoot)
	if err != nil {
		t.Fatalf("--accept --yes: %v\n%s", err, stderr)
	}
	if !strings.Contains(stderr, "contract activated for demo") || !strings.Contains(stdout, "quarantine-would-be") {
		t.Fatalf("--accept --yes output:\n%s\n%s", stdout, stderr)
	}
	settings, err = contract.LoadSettings(root)
	if err != nil {
		t.Fatal(err)
	}
	if settings.Contract.Hash() != draftHash || settings.ContractDraft != nil {
		t.Fatalf("after activation: %+v", settings)
	}
	if rendered, _ := os.ReadFile(claudePath); !strings.Contains(string(rendered), "Collect material on typed decision models.") || !strings.Contains(string(rendered), "Rendered by kb") {
		t.Fatalf("CLAUDE.md not re-rendered:\n%s", rendered)
	}

	// kb classify reuses the preview receipts: only the facet requests are new.
	qualityQuestions := fake.QuestionCount("paywall_or_login")
	kindQuestions := fake.QuestionCount("kind")
	if got := fake.QuestionCount("role"); got != roleQuestions {
		t.Fatalf("the second preview asked %d new role questions", got-roleQuestions)
	}
	runCLI(t, "classify", "demo", "--vault", vaultRoot)
	if got := fake.QuestionCount("role"); got != roleQuestions {
		t.Fatalf("classify asked %d new role questions after acceptance", got-roleQuestions)
	}
	if got := fake.QuestionCount("paywall_or_login"); got != qualityQuestions {
		t.Fatalf("classify asked %d new quality questions after acceptance", got-qualityQuestions)
	}
	if got := fake.QuestionCount("kind"); got != kindQuestions+len(sources) {
		t.Fatalf("classify kind questions = %d, want %d new facet requests", got-kindQuestions, len(sources))
	}

	// The previewed labels count as holdout under the previewed contract and
	// are demoted to dev once the contract changes. Both labeled subjects
	// were used by the preview; `labeled` is holdout by construction.
	holdout := 0
	for _, subject := range []string{labeled, "raw/articles/core-1.md"} {
		if !review.DevBucket(subject) {
			holdout++
		}
	}
	report, err := review.Calibrate(root, review.CalibrateOptions{ContractHash: draftHash})
	if err != nil {
		t.Fatal(err)
	}
	if relevance := report.Purposes[0]; relevance.Purpose != review.PurposeRelevance || relevance.Demoted != 0 || relevance.HoldoutLabels != holdout {
		t.Fatalf("calibration under the previewed contract = %+v", relevance)
	}
	changed := settings.Contract.Normalized()
	changed.Purpose = "Collect material on typed decision models and their judges."
	if err := contract.SetContract(root, &changed); err != nil {
		t.Fatal(err)
	}
	report, err = review.Calibrate(root, review.CalibrateOptions{ContractHash: changed.Hash()})
	if err != nil {
		t.Fatal(err)
	}
	if relevance := report.Purposes[0]; relevance.Demoted != holdout || relevance.HoldoutLabels != 0 {
		t.Fatalf("calibration after a contract change = %+v", relevance)
	}
}
