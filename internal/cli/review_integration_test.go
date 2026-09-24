//go:build integration

package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

func roleReceipt(t *testing.T, offTopic float64) json.RawMessage {
	t.Helper()
	rest := (1 - offTopic) / 2
	raw, err := json.Marshal(map[string]any{
		"type": "choice", "choice": "core", "confidence": 0.9,
		"probabilities": map[string]float64{"core": rest, "adjacent": rest, "off_topic": offTopic},
	})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func calibrateJSON(t *testing.T, vaultPath string) review.PurposeReport {
	t.Helper()
	stdout, _, err := runReviewCLI(t, "review", "calibrate", "sports", "--format", "json", "--vault", vaultPath)
	if err != nil {
		t.Fatalf("calibrate: %v", err)
	}
	var report review.Report
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("calibrate json: %v\n%s", err, stdout)
	}
	for _, entry := range report.Purposes {
		if entry.Purpose == review.PurposeRelevance {
			return entry
		}
	}
	t.Fatal("no relevance entry")
	return review.PurposeReport{}
}

// TestReviewImportAndHoldoutIntegration is the spec §17 "import and holdout"
// flow: import-labels on a screening JSONL matched by pmcid reports
// unmatched rows; calibrate chooses the threshold on dev buckets only and
// prints holdout metrics; editing the contract after a preview demotes the
// previewed labels to dev.
func TestReviewImportAndHoldoutIntegration(t *testing.T) {
	vaultPath, root := newReviewVault(t)
	v1 := &contract.Contract{Purpose: "Sports technology research", Core: []string{"wearables for athletes"}, OutOfScope: []string{"radar"}}
	if err := contract.SetContract(root, v1); err != nil {
		t.Fatal(err)
	}

	receipts := decisions.NewReceipts()
	var screening strings.Builder
	holdoutNegatives := make([]string, 0)
	for i := range 150 {
		subject := fmt.Sprintf("raw/papers/p%03d.md", i)
		writeReviewFile(t, filepath.Join(root, subject), fmt.Sprintf("---\ntitle: Paper %d\npmcid: PMC%d\n---\nBody %d.\n", i, 1000+i, i))
		dev := review.DevBucket(subject)
		decision, p := "include", 0.2
		switch {
		case i%2 == 0:
			decision, p = "exclude", 0.85
			if !dev {
				holdoutNegatives = append(holdoutNegatives, subject)
			}
		case i%4 == 1 && dev:
			p = 0.72
		case i%4 == 1:
			p = 0.78 // would push the choice to 0.80 if holdout leaked into the sweep
		}
		if err := receipts.Append(root, decisions.Receipt{
			Key: fmt.Sprintf("r%03d", i), Subject: subject, Purpose: "relevance", Status: decisions.StatusDecided,
			Answers: map[string]json.RawMessage{"role": roleReceipt(t, p)},
		}); err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(&screening, "{\"pmcid\":\"PMC%d\",\"decision\":%q,\"screening_stage\":\"title\"}\n", 1000+i, decision)
	}
	screening.WriteString("{\"pmcid\":\"PMC9\",\"decision\":\"include\"}\n{\"pmcid\":\"PMC8\",\"decision\":\"exclude\"}\n")
	from := filepath.Join(t.TempDir(), "science-screening.jsonl")
	writeReviewFile(t, from, screening.String())

	stdout, _, err := runReviewCLI(t, "review", "import-labels", "sports", "--from", from, "--id-field", "pmcid",
		"--match", "pmcid", "--decision-field", "decision", "--keep-values", "include", "--vault", vaultPath)
	if err != nil {
		t.Fatalf("import-labels: %v", err)
	}
	if !strings.Contains(stdout, "matched: 150, unmatched: 2") || !strings.Contains(stdout, "  PMC9\n  PMC8") {
		t.Fatalf("import-labels output:\n%s", stdout)
	}

	// Imported labels alone never set a threshold (spec §12.2, §20): the
	// metrics are reported, the purpose is flagged, nothing is recommended.
	importedOnly := calibrateJSON(t, vaultPath)
	if importedOnly.Recommended || !importedOnly.ImportedOnly || importedOnly.ReviewDevLabels != 0 || importedOnly.DevCurrent.N == 0 {
		t.Fatalf("imported labels only = %+v", importedOnly)
	}
	if table, _, err := runReviewCLI(t, "review", "calibrate", "sports", "--vault", vaultPath); err != nil || !strings.Contains(table, "imported labels only") {
		t.Fatalf("calibrate table with imported labels only: %v\n%s", err, table)
	}
	if _, _, err := runReviewCLI(t, "review", "calibrate", "sports", "--write", "--vault", vaultPath); err != nil {
		t.Fatalf("calibrate --write: %v", err)
	}
	if record, err := session.ReadCalibration(root); err != nil || record != nil {
		t.Fatalf("imported labels only must write no calibration: %+v %v", record, err)
	}

	// Labels given in kb review on other sources reach the dev floor.
	store := review.Open(root, nil)
	reviewDev := 0
	for i := 150; reviewDev < review.MinDevLabels; i++ {
		subject := fmt.Sprintf("raw/papers/p%03d.md", i)
		if !review.DevBucket(subject) {
			continue
		}
		verdict, p := review.VerdictPositive, 0.2
		if i%2 == 0 {
			verdict, p = review.VerdictNegative, 0.85
		}
		key := fmt.Sprintf("r%03d", i)
		if err := receipts.Append(root, decisions.Receipt{
			Key: key, Subject: subject, Purpose: "relevance", Status: decisions.StatusDecided,
			Answers: map[string]json.RawMessage{"role": roleReceipt(t, p)},
		}); err != nil {
			t.Fatal(err)
		}
		if err := store.AddLabel(review.Label{Subject: subject, Purpose: review.PurposeRelevance, Question: review.QuestionRole, Verdict: verdict, ReceiptKey: key, Origin: review.OriginReview}); err != nil {
			t.Fatal(err)
		}
		reviewDev++
	}

	rel := calibrateJSON(t, vaultPath)
	if rel.ImportedOnly || rel.ReviewDevLabels != review.MinDevLabels {
		t.Fatalf("with review labels = %+v", rel)
	}
	if !rel.Recommended || rel.Chosen != 0.75 {
		t.Fatalf("chosen on dev = %+v", rel)
	}
	if rel.HoldoutLabels < review.MinHoldoutLabels || rel.HoldoutChosen.N != rel.HoldoutLabels || rel.HoldoutChosen.Precision >= 1 {
		t.Fatalf("holdout metrics = %+v (labels %d)", rel.HoldoutChosen, rel.HoldoutLabels)
	}
	if rel.ImportedLabels != 150 || rel.Demoted != 0 {
		t.Fatalf("imported %d demoted %d", rel.ImportedLabels, rel.Demoted)
	}
	table, _, err := runReviewCLI(t, "review", "calibrate", "sports", "--vault", vaultPath)
	if err != nil || !strings.Contains(table, "holdout (chosen)") || !strings.Contains(table, "relevance_quarantine = 0.75") {
		t.Fatalf("calibrate table: %v\n%s", err, table)
	}

	shown := holdoutNegatives[:12]
	if err := review.AppendPreview(root, review.Preview{Contract: v1.Hash(), Subjects: shown}); err != nil {
		t.Fatal(err)
	}
	if again := calibrateJSON(t, vaultPath); again.Demoted != 0 || again.HoldoutLabels != rel.HoldoutLabels {
		t.Fatalf("preview under the active contract demoted labels: %+v", again)
	}
	v2 := &contract.Contract{Purpose: v1.Purpose, Core: []string{"wearables for athletes", "sports analytics"}, OutOfScope: v1.OutOfScope}
	if err := contract.SetContract(root, v2); err != nil {
		t.Fatal(err)
	}
	edited := calibrateJSON(t, vaultPath)
	if edited.Demoted != len(shown) || edited.HoldoutLabels != rel.HoldoutLabels-len(shown) || edited.DevLabels != rel.DevLabels+len(shown) {
		t.Fatalf("after contract edit: demoted %d dev %d holdout %d (before dev %d holdout %d)",
			edited.Demoted, edited.DevLabels, edited.HoldoutLabels, rel.DevLabels, rel.HoldoutLabels)
	}

	if _, _, err := runReviewCLI(t, "review", "calibrate", "sports", "--write", "--vault", vaultPath); err != nil {
		t.Fatalf("calibrate --write: %v", err)
	}
	settings, err := contract.LoadSettings(root)
	if err != nil {
		t.Fatal(err)
	}
	if settings.Decisions.Thresholds["relevance_quarantine"] != 0.75 || settings.Contract.Hash() != v2.Hash() {
		t.Fatalf("topic.yaml after --write: %+v", settings.Decisions)
	}
	record, err := session.ReadCalibration(root)
	if err != nil || record == nil {
		t.Fatalf("calibration.json: %+v %v", record, err)
	}
	entry := record.Purposes[review.PurposeRelevance]
	if entry.DevLabels != edited.DevLabels || entry.HoldoutLabels != edited.HoldoutLabels || entry.ImportedLabels != 150 {
		t.Fatalf("calibration entry = %+v", entry)
	}
}
