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
	"github.com/compozy/kb/internal/questions"
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

// Decision context of the fixtures: receipts are judged under contract
// fixtureContract by testModel with the current built-in banks.
const (
	fixtureContract = "c1"
	testModel       = "test/jev"
)

// calibrateOptions calibrates in the fixtures' decision context under
// contractHash.
func calibrateOptions(contractHash string) CalibrateOptions {
	return CalibrateOptions{ContractHash: contractHash, Model: testModel, Now: fixedClock}
}

// inContext stamps a receipt with a decision context: contractHash,
// testModel, and bank — a built-in bank id or a "+"-joined composite —
// at the current built-in versions.
func inContext(receipt decisions.Receipt, contractHash, bank string) decisions.Receipt {
	ids := strings.Split(bank, "+")
	versions := make([]string, len(ids))
	for index, id := range ids {
		versions[index] = questions.MustLoad(id).Version
	}
	receipt.Contract, receipt.Model = contractHash, testModel
	receipt.Bank, receipt.BankVersion = bank, strings.Join(versions, "+")
	return receipt
}

// relevanceSubject returns subject i of the relevance fixture with its
// verdict and P(off_topic): even subjects are off-topic (negative, 0.85);
// subjects ≡ 1 mod 4 are relevant but look off-topic (0.72); the rest are
// relevant (0.2).
func relevanceSubject(i int) (subject, verdict string, p float64) {
	subject, verdict, p = fmt.Sprintf("raw/s%03d.md", i), VerdictPositive, 0.2
	switch {
	case i%2 == 0:
		verdict, p = VerdictNegative, 0.85
	case i%4 == 1:
		p = 0.72
	}
	return subject, verdict, p
}

// appendRelevanceReceipts judges the n fixture subjects again: one gate
// receipt per subject, keyed prefix+i, adjusted by mutate.
func appendRelevanceReceipts(t *testing.T, root string, n int, prefix string, mutate func(*decisions.Receipt)) {
	t.Helper()
	receipts := decisions.NewReceipts()
	for i := range n {
		subject, _, p := relevanceSubject(i)
		receipt := inContext(decisions.Receipt{
			Key: fmt.Sprintf("%s%03d", prefix, i), Subject: subject, Purpose: "relevance",
			Answers: map[string]json.RawMessage{"role": roleAnswer(p)},
		}, fixtureContract, "relevance+quality")
		if mutate != nil {
			mutate(&receipt)
		}
		appendReceipt(t, receipts, root, receipt)
	}
}

// relevanceFixture writes n relevance labels, each pinned to a gate receipt
// of the fixture decision context (see relevanceSubject).
func relevanceFixture(t *testing.T, n int) string {
	t.Helper()
	return relevanceFixtureWithOrigin(t, n, func(int) string { return "" })
}

// relevanceFixtureWithOrigin is relevanceFixture with the label origin of
// subject i given by origin.
func relevanceFixtureWithOrigin(t *testing.T, n int, origin func(i int) string) string {
	t.Helper()
	root := t.TempDir()
	appendRelevanceReceipts(t, root, n, "k", nil)
	store := Open(root, fixedClock)
	for i := range n {
		subject, verdict, _ := relevanceSubject(i)
		key := fmt.Sprintf("k%03d", i)
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

	report, err := Calibrate(root, calibrateOptions(fixtureContract))
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
	report, err := Calibrate(few, calibrateOptions(fixtureContract))
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
			report, err := Calibrate(root, calibrateOptions(fixtureContract))
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

// TestCalibrateAcrossContractChange: labels are joined only to scores of the
// active decision context. After a contract change the labels given under
// the old contract are unjoined (their receipts ask different questions)
// until the subjects are judged again under the new contract; previewed
// holdout labels are then demoted to dev.
func TestCalibrateAcrossContractChange(t *testing.T) {
	t.Parallel()
	root := relevanceFixture(t, 200)
	holdout := make([]string, 0)
	for i := range 200 {
		if subject, _, _ := relevanceSubject(i); !DevBucket(subject) {
			holdout = append(holdout, subject)
		}
	}
	shown := holdout[:5]
	if err := AppendPreview(root, Preview{Contract: fixtureContract, Subjects: shown}); err != nil {
		t.Fatal(err)
	}

	base, err := Calibrate(root, calibrateOptions(fixtureContract))
	if err != nil {
		t.Fatal(err)
	}
	baseRel := purposeReport(t, base, PurposeRelevance)
	if baseRel.Demoted != 0 || baseRel.DevLabels+baseRel.HoldoutLabels != 200 || !baseRel.Recommended {
		t.Fatalf("under the labels' contract: %+v", baseRel)
	}

	// The contract changes: every label still pins a receipt of the old
	// contract, so none is joined and nothing can be recommended or written.
	changed, err := Calibrate(root, calibrateOptions("c2"))
	if err != nil {
		t.Fatal(err)
	}
	rel := purposeReport(t, changed, PurposeRelevance)
	if rel.Unjoined != 200 || rel.DevLabels+rel.HoldoutLabels != 0 || rel.Recommended || rel.Demoted != len(shown) {
		t.Fatalf("after the contract change: %+v", rel)
	}
	if err := WriteCalibration(root, changed); !errors.Is(err, ErrNotEnoughLabels) {
		t.Fatalf("WriteCalibration with only old-contract scores = %v, want ErrNotEnoughLabels", err)
	}

	// The subjects are judged again under the new contract: current scores
	// exist, labels join them, and the previewed holdout labels move to dev.
	appendRelevanceReceipts(t, root, 200, "n", func(r *decisions.Receipt) { r.Contract = "c2" })
	rejudged, err := Calibrate(root, calibrateOptions("c2"))
	if err != nil {
		t.Fatal(err)
	}
	rel = purposeReport(t, rejudged, PurposeRelevance)
	if rel.Unjoined != 0 || rel.Demoted != len(shown) || !rel.Recommended ||
		rel.HoldoutLabels != baseRel.HoldoutLabels-len(shown) || rel.DevLabels != baseRel.DevLabels+len(shown) {
		t.Fatalf("after re-judging: unjoined %d demoted %d dev %d holdout %d (base dev %d holdout %d)",
			rel.Unjoined, rel.Demoted, rel.DevLabels, rel.HoldoutLabels, baseRel.DevLabels, baseRel.HoldoutLabels)
	}
	if err := WriteCalibration(root, rejudged); err != nil {
		t.Fatal(err)
	}
	record := readCalibrationRecord(t, root)
	entry := record.Purposes[PurposeRelevance]
	if entry.Contract != "c2" || entry.Model != testModel || entry.Banks["relevance"] != questions.MustLoad("relevance").Version {
		t.Fatalf("calibration context = %q %q %v", entry.Contract, entry.Model, entry.Banks)
	}
}

// TestCalibrateJoinsOnlyTheActiveDecisionContext: a receipt from another
// contract, model or bank version is never joined, even when a label pins it.
func TestCalibrateJoinsOnlyTheActiveDecisionContext(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		mutate     func(*decisions.Receipt)
		wantJoined int
	}{
		{name: "active context", wantJoined: 200},
		{name: "other contract", mutate: func(r *decisions.Receipt) { r.Contract = "c0" }},
		{name: "other model", mutate: func(r *decisions.Receipt) { r.Model = "other/model" }},
		{name: "stale relevance bank", mutate: func(r *decisions.Receipt) { r.BankVersion = "2020-01-01.1+" + questions.MustLoad("quality").Version }},
		{name: "bank without relevance", mutate: func(r *decisions.Receipt) { r.Bank, r.BankVersion = "quality", questions.MustLoad("quality").Version }},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			appendRelevanceReceipts(t, root, 200, "k", tc.mutate)
			store := Open(root, fixedClock)
			for i := range 200 {
				subject, verdict, _ := relevanceSubject(i)
				if err := store.AddLabel(Label{Subject: subject, Purpose: PurposeRelevance, Question: QuestionRole, Verdict: verdict, ReceiptKey: fmt.Sprintf("k%03d", i)}); err != nil {
					t.Fatal(err)
				}
			}
			report, err := Calibrate(root, calibrateOptions(fixtureContract))
			if err != nil {
				t.Fatal(err)
			}
			rel := purposeReport(t, report, PurposeRelevance)
			if joined := rel.DevLabels + rel.HoldoutLabels; joined != tc.wantJoined || rel.Unjoined != 200-tc.wantJoined {
				t.Fatalf("joined %d unjoined %d, want %d joined", joined, rel.Unjoined, tc.wantJoined)
			}
			if rel.Recommended != (tc.wantJoined > 0) {
				t.Fatalf("recommended = %v with %d joined labels", rel.Recommended, tc.wantJoined)
			}
		})
	}
}

func TestCalibrateJoinsLinkAndQualityLabels(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	receipts := decisions.NewReceipts()
	store := Open(root, fixedClock)
	const subject = "wiki/concepts/A.md"

	appendReceipt(t, receipts, root, inContext(decisions.Receipt{Key: "L0", Subject: subject, Purpose: "link", Answers: map[string]json.RawMessage{"should_link_c1": noulAnswer(0.8)}}, "c0", "link"))
	appendReceipt(t, receipts, root, inContext(decisions.Receipt{Key: "L1", Subject: subject, Purpose: "link", Answers: map[string]json.RawMessage{"should_link_c1": noulAnswer(0.9), "should_link_c2": noulAnswer(0.3)}}, fixtureContract, "link"))
	// The shape the link producer writes: the producing question and its
	// receipt, with the queue identity on the generic question name.
	item := Item{
		ID:    ItemID(QueueLink, subject, "wiki/concepts/B.md", QuestionShouldLink),
		Queue: QueueLink, Purpose: PurposeLink, Subject: subject, Target: "wiki/concepts/B.md", Question: "should_link_c2", ReceiptKey: "L1",
	}
	if _, err := store.Add(item); err != nil {
		t.Fatal(err)
	}
	labels := []Label{
		{Subject: subject, Target: "wiki/concepts/C.md", Purpose: PurposeLink, Question: "should_link_c1", Verdict: VerdictPositive, ReceiptKey: "L1"},
		{Subject: subject, Target: "wiki/concepts/B.md", Purpose: PurposeLink, Question: QuestionShouldLink, Verdict: VerdictNegative, Origin: OriginImportLinks},
		{Subject: subject, Target: "wiki/concepts/Z.md", Purpose: PurposeLink, Question: QuestionShouldLink, Verdict: VerdictPositive, Origin: OriginImportLinks},
		{Subject: "raw/q.md", Purpose: PurposeQuality, Verdict: VerdictNegative},
		// Pinned to a receipt of another contract: not joined.
		{Subject: subject, Target: "wiki/concepts/D.md", Purpose: PurposeLink, Question: "should_link_c1", Verdict: VerdictPositive, ReceiptKey: "L0"},
		// A positional question id without its receipt is ambiguous: another
		// receipt of the subject may ask c1 about a different target.
		{Subject: subject, Target: "wiki/concepts/E.md", Purpose: PurposeLink, Question: "should_link_c1", Verdict: VerdictPositive},
	}
	appendReceipt(t, receipts, root, inContext(decisions.Receipt{Key: "Q1", Subject: "raw/q.md", Purpose: "relevance+quality", Answers: map[string]json.RawMessage{"thin_or_boilerplate": noulAnswer(0.4), "error_or_placeholder_page": noulAnswer(0.95), "role": roleAnswer(0.1)}}, fixtureContract, "relevance+quality"))
	for _, label := range labels {
		if err := store.AddLabel(label); err != nil {
			t.Fatal(err)
		}
	}

	context := decisionContext{contract: fixtureContract, model: testModel, banks: map[string]string{}}
	specs := map[string]purposeSpec{}
	for _, spec := range calibratedPurposes {
		context.banks[spec.bank] = questions.MustLoad(spec.bank).Version
		specs[spec.purpose] = spec
	}
	joiner := newReceiptJoiner(mustReceipts(t, root), mustItems(t, store), context)
	testCases := []struct {
		label  Label
		want   float64
		joined bool
	}{
		{label: labels[0], want: 0.9, joined: true},
		{label: labels[1], want: 0.3, joined: true},
		{label: labels[2], joined: false},
		{label: labels[3], want: 0.95, joined: true},
		{label: labels[4], joined: false},
		{label: labels[5], joined: false},
	}
	for _, tc := range testCases {
		got, ok := joiner.quantity(specs[tc.label.Purpose], tc.label)
		if ok != tc.joined || got != tc.want {
			t.Errorf("quantity(%+v) = %.2f %v, want %.2f %v", tc.label, got, ok, tc.want, tc.joined)
		}
	}

	report, err := Calibrate(root, calibrateOptions(fixtureContract))
	if err != nil {
		t.Fatal(err)
	}
	link := purposeReport(t, report, PurposeLink)
	if link.Unjoined != 3 || link.DevLabels+link.HoldoutLabels != 2 || link.ImportedLabels != 1 {
		t.Fatalf("link report = %+v", link)
	}
}

func readCalibrationRecord(t *testing.T, root string) calibrationRecord {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, decisions.ReceiptsDir, CalibrationFile))
	if err != nil {
		t.Fatal(err)
	}
	var record calibrationRecord
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	return record
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
	report, err := Calibrate(root, calibrateOptions(fixtureContract))
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
	if !ok || entry.DevLabels+entry.HoldoutLabels != 200 || entry.Thresholds["relevance_quarantine"] != 0.75 || entry.Holdout.N != entry.HoldoutLabels || entry.Contract != fixtureContract || entry.Model != testModel || entry.Banks["relevance"] != questions.MustLoad("relevance").Version {
		t.Fatalf("record = %+v", record)
	}
	if record.Time != "2026-09-24T12:00:00Z" {
		t.Fatalf("time = %q", record.Time)
	}
}
