package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

// newReviewVault scaffolds a vault with one topic and isolates config
// discovery. Tests using it must not run in parallel (t.Setenv).
func newReviewVault(t *testing.T) (string, string) {
	t.Helper()
	vaultPath := t.TempDir()
	info, err := topic.New(vaultPath, "sports", "Sports", "sports")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	configPath := filepath.Join(t.TempDir(), "kb.toml")
	if err := os.WriteFile(configPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("APP_CONFIG", configPath)
	t.Setenv("OPENROUTER_API_KEY", "")
	return vaultPath, info.RootPath
}

func runReviewCLI(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	command := newRootCommand()
	var stdout, stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), stderr.String(), err
}

func writeReviewFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReviewListWithoutDecisionModel(t *testing.T) {
	vaultPath, root := newReviewVault(t)
	writeReviewFile(t, filepath.Join(root, "raw/articles/a.md"), "---\ntitle: Article A\n---\nBody.\n")
	store := review.Open(root, nil)
	for _, item := range []review.Item{
		{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/articles/a.md", Question: "role", Probability: 0.7, Evidence: "off topic"},
		{Queue: review.QueueLink, Purpose: "link", Subject: "raw/articles/a.md", Target: "wiki/concepts/B.md", Question: "should_link_c1", Probability: 0.65},
	} {
		if _, err := store.Add(item); err != nil {
			t.Fatal(err)
		}
	}

	testCases := []struct {
		name     string
		args     []string
		contains []string
		absent   []string
	}{
		{
			name:     "bare topic lists every queue",
			args:     []string{"review", "sports", "--vault", vaultPath},
			contains: []string{"remove", "link", "Article A (raw/articles/a.md)", "should_link_c1", "0.70", "concept-proposal"},
		},
		{
			name:     "list subcommand filters one queue",
			args:     []string{"review", "list", "sports", "--queue", "link", "--vault", vaultPath},
			contains: []string{"should_link_c1"},
			absent:   []string{"off topic"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := runReviewCLI(t, tc.args...)
			if err != nil {
				t.Fatalf("review: %v", err)
			}
			for _, want := range tc.contains {
				if !strings.Contains(stdout, want) {
					t.Errorf("output lacks %q:\n%s", want, stdout)
				}
			}
			for _, unwanted := range tc.absent {
				if strings.Contains(stdout, unwanted) {
					t.Errorf("output has %q:\n%s", unwanted, stdout)
				}
			}
		})
	}

	stdout, _, err := runReviewCLI(t, "review", "list", "sports", "--format", "json", "--vault", vaultPath)
	if err != nil {
		t.Fatal(err)
	}
	var view review.ListView
	if err := json.Unmarshal([]byte(stdout), &view); err != nil {
		t.Fatalf("json output: %v\n%s", err, stdout)
	}
	if view.Counts[review.QueueRemove] != 1 || len(view.Items) != 2 || view.Items[0].Title != "Article A" {
		t.Fatalf("view = %+v", view)
	}

	if _, _, err := runReviewCLI(t, "review", "list", "sports", "--queue", "nope", "--vault", vaultPath); err == nil {
		t.Fatal("unknown queue accepted")
	}
}

func TestReviewImportLabelsAndCalibrateCommands(t *testing.T) {
	vaultPath, root := newReviewVault(t)
	writeReviewFile(t, filepath.Join(root, "raw/papers/a.md"), "---\ntitle: A\npmcid: PMC1\n---\nA.\n")
	from := filepath.Join(t.TempDir(), "screening.jsonl")
	writeReviewFile(t, from, `{"pmcid":"PMC1","decision":"include"}`+"\n"+`{"pmcid":"PMC2","decision":"exclude"}`+"\n")

	stdout, _, err := runReviewCLI(t, "review", "import-labels", "sports", "--from", from, "--id-field", "pmcid",
		"--match", "pmcid", "--decision-field", "decision", "--keep-values", "include,true", "--vault", vaultPath)
	if err != nil {
		t.Fatalf("import-labels: %v", err)
	}
	for _, want := range []string{"matched: 1, unmatched: 1", "1 positive, 0 negative", "PMC2", "origin: import:screening.jsonl@"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("import-labels output lacks %q:\n%s", want, stdout)
		}
	}
	if _, _, err := runReviewCLI(t, "review", "import-labels", "sports", "--from", from, "--vault", vaultPath); err == nil {
		t.Error("import-labels without required flags succeeded")
	}

	stdout, stderr, err := runReviewCLI(t, "review", "calibrate", "sports", "--write", "--vault", vaultPath)
	if err != nil {
		t.Fatalf("calibrate: %v", err)
	}
	if !strings.Contains(stdout, "not enough labels") || !strings.Contains(stderr, "nothing written") {
		t.Fatalf("calibrate output:\n%s\n%s", stdout, stderr)
	}
	if record, err := session.ReadCalibration(root); err != nil || record != nil {
		t.Fatalf("calibration written: %+v %v", record, err)
	}

	stdout, _, err = runReviewCLI(t, "review", "import-links", "sports", "--vault", vaultPath)
	if err != nil || !strings.Contains(stdout, "link labels into sports") {
		t.Fatalf("import-links: %v\n%s", err, stdout)
	}
}

// TestReviewCalibrationRecordReadBySession pins the calibration.json contract
// between review (writer) and session (reader of the default-mode rule).
func TestReviewCalibrationRecordReadBySession(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	report := review.Report{Time: "2026-09-24T12:00:00Z", Purposes: []review.PurposeReport{{
		Purpose: "relevance", ApplyName: "relevance_quarantine", Recommended: true, Chosen: 0.75,
		DevLabels: 40, HoldoutLabels: 15, ImportedLabels: 20,
		DevChosen:     review.Metrics{N: 40, Precision: 0.95, Recall: 0.8, Coverage: 1, ReviewRate: 0.1},
		HoldoutChosen: review.Metrics{N: 15, Precision: 0.9, Recall: 0.7, Coverage: 1, ReviewRate: 0.2},
	}}}
	if err := review.WriteCalibration(root, report); err != nil {
		t.Fatal(err)
	}
	record, err := session.ReadCalibration(root)
	if err != nil || record == nil {
		t.Fatalf("ReadCalibration = %+v, %v", record, err)
	}
	entry := record.Purposes["relevance"]
	if entry.DevLabels != 40 || entry.HoldoutLabels != 15 || entry.ImportedLabels != 20 ||
		entry.Thresholds["relevance_quarantine"] != 0.75 || entry.Holdout.Precision != 0.9 || entry.Dev.N != 40 {
		t.Fatalf("entry = %+v", entry)
	}
}
