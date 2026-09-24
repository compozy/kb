package review

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/decisions"
)

func roleAnswer(offTopic float64) json.RawMessage {
	rest := (1 - offTopic) / 2
	raw, _ := json.Marshal(map[string]any{
		"type": "choice", "choice": "core", "confidence": 0.9,
		"probabilities": map[string]float64{"core": rest, "adjacent": rest, "off_topic": offTopic},
	})
	return raw
}

func noulAnswer(p float64) json.RawMessage {
	raw, _ := json.Marshal(map[string]any{"type": "noul", "noul": p})
	return raw
}

func appendReceipt(t *testing.T, store *decisions.Receipts, root string, receipt decisions.Receipt) {
	t.Helper()
	if receipt.Status == "" {
		receipt.Status = decisions.StatusDecided
	}
	if err := store.Append(root, receipt); err != nil {
		t.Fatal(err)
	}
}

// relevanceFixture writes n relevance labels with receipts: even subjects are
// off-topic (negative, P(off_topic)=0.85); subjects ≡ 1 mod 4 are relevant but
// look off-topic (0.72); the rest are relevant (0.2).
func relevanceFixture(t *testing.T, n int) string {
	t.Helper()
	return relevanceFixtureWithOrigin(t, n, func(int) string { return "" })
}

// relevanceFixtureWithOrigin is relevanceFixture with the label origin of
// subject i given by origin.
func relevanceFixtureWithOrigin(t *testing.T, n int, origin func(i int) string) string {
	t.Helper()
	root := t.TempDir()
	receipts := decisions.NewReceipts()
	store := Open(root, fixedClock)
	for i := range n {
		subject := fmt.Sprintf("raw/s%03d.md", i)
		verdict, p := VerdictPositive, 0.2
		switch {
		case i%2 == 0:
			verdict, p = VerdictNegative, 0.85
		case i%4 == 1:
			p = 0.72
		}
		key := fmt.Sprintf("k%03d", i)
		appendReceipt(t, receipts, root, decisions.Receipt{Key: key, Subject: subject, Purpose: "relevance+quality", Answers: map[string]json.RawMessage{"role": roleAnswer(p)}})
		if err := store.AddLabel(Label{Subject: subject, Purpose: PurposeRelevance, Question: QuestionRole, Verdict: verdict, ReceiptKey: key, Origin: origin(i)}); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func purposeReport(t *testing.T, report Report, purpose string) PurposeReport {
	t.Helper()
	for _, entry := range report.Purposes {
		if entry.Purpose == purpose {
			return entry
		}
	}
	t.Fatalf("no %s in report", purpose)
	return PurposeReport{}
}

func TestDevBucketIsDeterministic(t *testing.T) {
	t.Parallel()
	dev := 0
	for i := range 1000 {
		subject := fmt.Sprintf("raw/x%d.md", i)
		first := DevBucket(subject)
		if DevBucket(subject) != first {
			t.Fatalf("DevBucket(%q) not stable", subject)
		}
		if first {
			dev++
		}
	}
	if dev < 600 || dev > 800 {
		t.Fatalf("dev share = %d/1000, want about 700", dev)
	}
}

func TestCalibrateChoosesThresholdOnDev(t *testing.T) {
	t.Parallel()
	root := relevanceFixture(t, 200)

	report, err := Calibrate(root, CalibrateOptions{ContractHash: "c1", Now: fixedClock})
	if err != nil {
		t.Fatalf("Calibrate: %v", err)
	}
	rel := purposeReport(t, report, PurposeRelevance)
	if !rel.Recommended || rel.Chosen != 0.75 {
		t.Fatalf("relevance = %+v", rel)
	}
	if rel.DevLabels+rel.HoldoutLabels != 200 || rel.DevLabels < MinDevLabels || rel.HoldoutLabels < MinHoldoutLabels {
		t.Fatalf("label counts dev %d holdout %d", rel.DevLabels, rel.HoldoutLabels)
	}
	if rel.DevChosen.Precision != 1 || rel.HoldoutChosen.Precision != 1 || rel.HoldoutChosen.Recall != 1 {
		t.Fatalf("chosen metrics dev %+v holdout %+v", rel.DevChosen, rel.HoldoutChosen)
	}
	if rel.DevCurrent.Precision != 1 || rel.Current != 0.8 {
		t.Fatalf("current = %.2f %+v", rel.Current, rel.DevCurrent)
	}
	if rel.HoldoutCurrent.N != rel.HoldoutLabels || rel.HoldoutCurrent.ReviewRate == 0 {
		t.Fatalf("holdout current = %+v", rel.HoldoutCurrent)
	}
	for _, point := range rel.Sweep {
		if point.Threshold < 0.75 && point.Dev.Precision >= TargetPrecision {
			t.Fatalf("sweep point %+v should be below target precision", point)
		}
	}
	if link := purposeReport(t, report, PurposeLink); link.Recommended || link.Note == "" {
		t.Fatalf("link without labels = %+v", link)
	}
}

func TestCalibrateFloorsAndF05Fallback(t *testing.T) {
	t.Parallel()

	few := relevanceFixture(t, 20)
	report, err := Calibrate(few, CalibrateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if rel := purposeReport(t, report, PurposeRelevance); rel.Recommended || rel.Note == "" {
		t.Fatalf("20 labels recommended: %+v", rel)
	}
	if err := WriteCalibration(few, report); !errors.Is(err, ErrNotEnoughLabels) {
		t.Fatalf("WriteCalibration = %v, want ErrNotEnoughLabels", err)
	}
	if _, err := os.Stat(filepath.Join(few, decisions.ReceiptsDir, CalibrationFile)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("calibration.json written without enough labels")
	}

	sweep := make([]SweepPoint, 0)
	for step := 10; step <= 19; step++ {
		threshold := float64(step) / 20
		precision := 0.6 + float64(step-10)*0.02
		recall := 1 - float64(step-10)*0.1
		sweep = append(sweep, SweepPoint{Threshold: threshold, Dev: Metrics{Predicted: 10, Precision: precision, Recall: recall}})
	}
	chosen, ok := chooseThreshold(sweep)
	if !ok || chosen != 0.65 { // F0.5 peaks at precision 0.66, recall 0.70
		t.Fatalf("F0.5 fallback chose %.2f %v", chosen, ok)
	}
	if _, ok := chooseThreshold([]SweepPoint{{Threshold: 0.5, Dev: Metrics{}}}); ok {
		t.Fatal("threshold chosen without true positives")
	}
}

// TestCalibrateImportedLabelsNeverSetThresholdsAlone: the dev floor must be
// met by labels given in kb review; imported labels alone are measured and
// flagged, never recommended or written (spec §12.2, §12.3, §20).
func TestCalibrateImportedLabelsNeverSetThresholdsAlone(t *testing.T) {
	t.Parallel()
	imported := OriginImportPrefix + "screening.jsonl@abc"
	testCases := []struct {
		name        string
		origin      func(i int) string
		wantRec     bool
		wantFlagged bool
	}{
		{name: "imported only", origin: func(int) string { return imported }, wantFlagged: true},
		{name: "import-links only", origin: func(int) string { return OriginImportLinks }, wantFlagged: true},
		{name: "mostly imported, review below floor", origin: func(i int) string {
			if i < 20 {
				return OriginReview
			}
			return imported
		}, wantFlagged: true},
		{name: "review labels", origin: func(int) string { return OriginReview }, wantRec: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := relevanceFixtureWithOrigin(t, 200, tc.origin)
			report, err := Calibrate(root, CalibrateOptions{ContractHash: "c1", Now: fixedClock})
			if err != nil {
				t.Fatal(err)
			}
			rel := purposeReport(t, report, PurposeRelevance)
			if rel.Recommended != tc.wantRec || rel.ImportedOnly != tc.wantFlagged {
				t.Fatalf("recommended %v imported-only %v, want %v %v (%+v)", rel.Recommended, rel.ImportedOnly, tc.wantRec, tc.wantFlagged, rel)
			}
			if rel.DevCurrent.N < MinDevLabels || rel.HoldoutCurrent.N < MinHoldoutLabels {
				t.Fatalf("metrics must still be reported: dev %+v holdout %+v", rel.DevCurrent, rel.HoldoutCurrent)
			}
			err = WriteCalibration(root, report)
			if tc.wantRec {
				if err != nil {
					t.Fatalf("WriteCalibration: %v", err)
				}
				return
			}
			if !errors.Is(err, ErrNotEnoughLabels) {
				t.Fatalf("WriteCalibration = %v, want ErrNotEnoughLabels", err)
			}
			if !strings.Contains(rel.Note, "imported labels only") {
				t.Fatalf("note = %q", rel.Note)
			}
		})
	}
}

func TestCalibratePreviewDemotesToDev(t *testing.T) {
	t.Parallel()
	root := relevanceFixture(t, 200)
	holdout := make([]string, 0)
	for i := range 200 {
		subject := fmt.Sprintf("raw/s%03d.md", i)
		if !DevBucket(subject) {
			holdout = append(holdout, subject)
		}
	}
	shown := holdout[:5]
	if err := AppendPreview(root, Preview{Contract: "old", Subjects: shown}); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		contract    string
		wantDemoted int
	}{
		{contract: "new", wantDemoted: 5},
		{contract: "old", wantDemoted: 0},
	}
	base, err := Calibrate(root, CalibrateOptions{ContractHash: "old"})
	if err != nil {
		t.Fatal(err)
	}
	baseRel := purposeReport(t, base, PurposeRelevance)
	for _, tc := range testCases {
		report, err := Calibrate(root, CalibrateOptions{ContractHash: tc.contract})
		if err != nil {
			t.Fatal(err)
		}
		rel := purposeReport(t, report, PurposeRelevance)
		if rel.Demoted != tc.wantDemoted || rel.HoldoutLabels != baseRel.HoldoutLabels-tc.wantDemoted || rel.DevLabels != baseRel.DevLabels+tc.wantDemoted {
			t.Errorf("contract %s: demoted %d dev %d holdout %d (base dev %d holdout %d)", tc.contract, rel.Demoted, rel.DevLabels, rel.HoldoutLabels, baseRel.DevLabels, baseRel.HoldoutLabels)
		}
	}
}

func TestCalibrateJoinsLinkAndQualityLabels(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	receipts := decisions.NewReceipts()
	store := Open(root, fixedClock)

	appendReceipt(t, receipts, root, decisions.Receipt{Key: "L1", Subject: "wiki/concepts/A.md", Purpose: "link", Answers: map[string]json.RawMessage{"should_link_c1": noulAnswer(0.9), "should_link_c2": noulAnswer(0.3)}})
	if _, err := store.Add(Item{Queue: QueueLink, Purpose: PurposeLink, Subject: "wiki/concepts/A.md", Target: "wiki/concepts/B.md", Question: "should_link_c2", ReceiptKey: "L1"}); err != nil {
		t.Fatal(err)
	}
	labels := []Label{
		{Subject: "wiki/concepts/A.md", Target: "wiki/concepts/C.md", Purpose: PurposeLink, Question: "should_link_c1", Verdict: VerdictPositive, ReceiptKey: "L1"},
		{Subject: "wiki/concepts/A.md", Target: "wiki/concepts/B.md", Purpose: PurposeLink, Question: QuestionShouldLink, Verdict: VerdictNegative, Origin: OriginImportLinks},
		{Subject: "wiki/concepts/A.md", Target: "wiki/concepts/Z.md", Purpose: PurposeLink, Question: QuestionShouldLink, Verdict: VerdictPositive, Origin: OriginImportLinks},
		{Subject: "raw/q.md", Purpose: PurposeQuality, Verdict: VerdictNegative},
	}
	appendReceipt(t, receipts, root, decisions.Receipt{Key: "Q1", Subject: "raw/q.md", Purpose: "relevance+quality", Answers: map[string]json.RawMessage{"thin_or_boilerplate": noulAnswer(0.4), "error_or_placeholder_page": noulAnswer(0.95), "role": roleAnswer(0.1)}})
	for _, label := range labels {
		if err := store.AddLabel(label); err != nil {
			t.Fatal(err)
		}
	}

	joiner := newReceiptJoiner(mustReceipts(t, root), mustItems(t, store))
	testCases := []struct {
		label  Label
		want   float64
		joined bool
	}{
		{label: labels[0], want: 0.9, joined: true},
		{label: labels[1], want: 0.3, joined: true},
		{label: labels[2], joined: false},
		{label: labels[3], want: 0.95, joined: true},
	}
	for _, tc := range testCases {
		got, ok := joiner.quantity(tc.label.Purpose, tc.label)
		if ok != tc.joined || got != tc.want {
			t.Errorf("quantity(%+v) = %.2f %v, want %.2f %v", tc.label, got, ok, tc.want, tc.joined)
		}
	}

	report, err := Calibrate(root, CalibrateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	link := purposeReport(t, report, PurposeLink)
	if link.Unjoined != 1 || link.DevLabels+link.HoldoutLabels != 2 || link.ImportedLabels != 1 {
		t.Fatalf("link report = %+v", link)
	}
}

func mustReceipts(t *testing.T, root string) []decisions.Receipt {
	t.Helper()
	rows, err := decisions.LoadReceipts(root)
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

func mustItems(t *testing.T, store *Store) []Item {
	t.Helper()
	items, err := store.Items("", "")
	if err != nil {
		t.Fatal(err)
	}
	return items
}

func TestWriteCalibrationStoresThresholdsAndRecord(t *testing.T) {
	t.Parallel()
	root := relevanceFixture(t, 200)
	if err := os.WriteFile(contract.SettingsPath(root), []byte("# topic settings\ndecisions:\n  mode: shadow\n  thresholds:\n    link_apply: 0.9\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Calibrate(root, CalibrateOptions{Now: fixedClock})
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteCalibration(root, report); err != nil {
		t.Fatalf("WriteCalibration: %v", err)
	}
	settings, err := contract.LoadSettings(root)
	if err != nil {
		t.Fatal(err)
	}
	if settings.Decisions.Thresholds["relevance_quarantine"] != 0.75 || settings.Decisions.Thresholds["link_apply"] != 0.9 || settings.Decisions.Mode != "shadow" {
		t.Fatalf("settings = %+v", settings.Decisions)
	}
	data, err := os.ReadFile(filepath.Join(root, decisions.ReceiptsDir, CalibrationFile))
	if err != nil {
		t.Fatal(err)
	}
	var record calibrationRecord
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	entry, ok := record.Purposes[PurposeRelevance]
	if !ok || entry.DevLabels+entry.HoldoutLabels != 200 || entry.Thresholds["relevance_quarantine"] != 0.75 || entry.Holdout.N != entry.HoldoutLabels {
		t.Fatalf("record = %+v", record)
	}
	if record.Time != "2026-09-24T12:00:00Z" {
		t.Fatalf("time = %q", record.Time)
	}
}
