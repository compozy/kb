package scope

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/review"
)

var defaultBandThresholds = map[string]float64{
	"relevance_quarantine": 0.8,
	"relevance_review":     0.5,
	"quality_apply":        0.8,
	"quality_review":       0.5,
}

func TestBandOf(t *testing.T) {
	t.Parallel()
	decided := func(pOff float64) classify.GateJudgment {
		return classify.GateJudgment{RoleAsked: true, RoleDecided: true, POffTopic: pOff, Quality: map[string]float64{}, QualityDecided: true}
	}
	tests := []struct {
		name       string
		gate       classify.GateJudgment
		wantBand   string
		wantReason string
	}{
		{name: "off topic at quarantine", gate: decided(0.8), wantBand: BandQuarantine, wantReason: classify.RoleOffTopic},
		{name: "off topic at review", gate: decided(0.5), wantBand: BandReview, wantReason: classify.RoleOffTopic},
		{name: "kept", gate: decided(0.1), wantBand: BandKept},
		{name: "code flag quarantines", gate: func() classify.GateJudgment {
			g := decided(0.1)
			g.Flags = []quality.Flag{{Code: quality.Thin}}
			return g
		}(), wantBand: BandQuarantine, wantReason: quality.Thin},
		{name: "quality noul at apply", gate: func() classify.GateJudgment {
			g := decided(0.1)
			g.Quality[classify.NoulPaywall] = 0.9
			return g
		}(), wantBand: BandQuarantine, wantReason: classify.ReasonPaywall},
		{name: "quality noul at review", gate: func() classify.GateJudgment {
			g := decided(0.1)
			g.Quality[classify.NoulThin] = 0.6
			return g
		}(), wantBand: BandReview, wantReason: classify.ReasonThin},
		{name: "role undecided", gate: classify.GateJudgment{RoleAsked: true, Quality: map[string]float64{}, QualityDecided: true}, wantBand: BandUndecided},
		{name: "quality undecided", gate: classify.GateJudgment{RoleAsked: true, RoleDecided: true, Quality: map[string]float64{}}, wantBand: BandUndecided},
		{name: "undecided role still quarantined by quality", gate: classify.GateJudgment{RoleAsked: true, Quality: map[string]float64{classify.NoulErrorPage: 0.95}}, wantBand: BandQuarantine, wantReason: classify.ReasonErrorPage},
		{name: "relevance off, quality clean", gate: classify.GateJudgment{Quality: map[string]float64{}, QualityDecided: true}, wantBand: BandKept},
		{name: "collected on purpose by path", gate: classify.GateJudgment{RoleDecided: true, RoleByPath: true, Quality: map[string]float64{}, QualityDecided: true}, wantBand: BandKept},
	}
	for _, tt := range tests {
		band, reason := bandOf(tt.gate, defaultBandThresholds)
		if band != tt.wantBand || reason != tt.wantReason {
			t.Errorf("%s: bandOf = (%q, %q), want (%q, %q)", tt.name, band, reason, tt.wantBand, tt.wantReason)
		}
	}
}

func TestScoreLabels(t *testing.T) {
	t.Parallel()
	judged := func(path string, pOff float64) judgment {
		return judgment{doc: &corpus.Document{Path: path}, gate: classify.GateJudgment{RoleAsked: true, RoleDecided: true, POffTopic: pOff}}
	}
	judgments := map[string]judgment{
		"raw/a.md": judged("raw/a.md", 0.9), // predicted, negative → TP
		"raw/b.md": judged("raw/b.md", 0.85),
		"raw/c.md": judged("raw/c.md", 0.1), // negative, missed → FN
		"raw/d.md": judged("raw/d.md", 0.2),
		"raw/e.md": {doc: &corpus.Document{Path: "raw/e.md"}, gate: classify.GateJudgment{RoleAsked: true}},
	}
	label := func(subject, purpose, verdict string) review.Label {
		return review.Label{Subject: subject, Purpose: purpose, Verdict: verdict}
	}
	labels := []review.Label{
		label("raw/a.md", review.PurposeRelevance, review.VerdictNegative),
		label("raw/b.md", review.PurposeRelevance, review.VerdictNegative),
		label("raw/b.md", review.PurposeRelevance, review.VerdictPositive), // latest wins: FP
		label("raw/c.md", review.PurposeRelevance, review.VerdictNegative),
		label("raw/d.md", review.PurposeRelevance, review.VerdictPositive),
		label("raw/d.md", review.PurposeLink, review.VerdictNegative),       // other purpose ignored
		label("raw/e.md", review.PurposeRelevance, review.VerdictNegative),  // undecided role ignored
		label("raw/zz.md", review.PurposeRelevance, review.VerdictNegative), // not judged
	}

	score, used := scoreLabels(judgments, labels, 0.8)
	want := &LabelScore{Labels: 4, Negatives: 2, Predicted: 2, TruePositive: 1, Precision: 0.5, Recall: 0.5}
	if !reflect.DeepEqual(score, want) {
		t.Fatalf("score = %+v, want %+v", score, want)
	}
	if !reflect.DeepEqual(used, []string{"raw/a.md", "raw/b.md", "raw/c.md", "raw/d.md"}) {
		t.Fatalf("used = %v", used)
	}
	if score, used := scoreLabels(judgments, nil, 0.8); score != nil || used != nil {
		t.Fatalf("no labels = %+v, %v", score, used)
	}
}

// previewTopic writes a topic with three core sources, one off-topic, one
// paywalled one under raw/news, and relevance labels.
func previewTopic(t *testing.T) (vault, root string) {
	t.Helper()
	vault, root = newTestTopic(t)
	for _, name := range []string{"core-1", "core-2", "core-3"} {
		writeSource(t, root, "raw/articles/"+name+".md", "Jev "+name, "Typed decision models answer questions.\n")
	}
	writeSource(t, root, "raw/articles/offtopic.md", "Offtopic sports results", "Scores of last night.\n")
	writeSource(t, root, "raw/news/paywalled.md", "Paywalled story", "Subscribe to read.\n")
	store := review.Open(root, nil)
	for subject, verdict := range map[string]string{
		"raw/articles/offtopic.md": review.VerdictNegative,
		"raw/articles/core-1.md":   review.VerdictNegative,
		"raw/articles/core-2.md":   review.VerdictPositive,
	} {
		if err := store.AddLabel(review.Label{Subject: subject, Purpose: review.PurposeRelevance, Verdict: verdict, Origin: "import:screening.jsonl@abc"}); err != nil {
			t.Fatal(err)
		}
	}
	return vault, root
}

func previewDecide(call fakes.Call, q fakes.Question) any {
	if q.ID == classify.NoulPaywall && strings.Contains(call.StateString(), "Paywalled") {
		return fakes.Noul(0.9)
	}
	return roleByTitle(call, q)
}

func TestPreviewBandsTopListLabelsAndLedger(t *testing.T) {
	t.Parallel()
	vault, root := previewTopic(t)
	fake := fakeServer(t, previewDecide, nil)
	s := openTestSession(t, vault, fake)
	var estimated int

	report, err := Preview(context.Background(), s, testDraft, PreviewOptions{OnEstimate: func(n int, usd float64) {
		estimated = n
		if usd != float64(n)*CostPerDocumentUSD {
			t.Errorf("estimate = %v for %d documents", usd, n)
		}
	}})
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	if estimated != 5 || report.Sources != 5 || report.Judged != 5 || report.Sampled || report.Contract != testDraft.Hash() {
		t.Fatalf("report header = %+v (estimated %d)", report, estimated)
	}
	if report.Bands != (BandCounts{Kept: 3, Quarantine: 2}) {
		t.Fatalf("bands = %+v", report.Bands)
	}
	wantFolders := []FolderBands{
		{Folder: "raw/articles", Judged: 4, BandCounts: BandCounts{Kept: 3, Quarantine: 1}},
		{Folder: "raw/news", Judged: 1, BandCounts: BandCounts{Quarantine: 1}},
	}
	if !reflect.DeepEqual(report.Folders, wantFolders) {
		t.Fatalf("folders = %+v", report.Folders)
	}
	if len(report.Top) != 5 || report.Top[0].Path != "raw/articles/offtopic.md" || report.Top[0].Band != BandQuarantine || report.Top[0].Reason != "" {
		t.Fatalf("top = %+v", report.Top)
	}
	for _, top := range report.Top {
		if top.Path == "raw/news/paywalled.md" && top.Reason != classify.ReasonPaywall {
			t.Fatalf("paywalled top entry = %+v", top)
		}
	}
	wantLabels := &LabelScore{Labels: 3, Negatives: 2, Predicted: 1, TruePositive: 1, Precision: 1, Recall: 0.5}
	if !reflect.DeepEqual(report.Labels, wantLabels) {
		t.Fatalf("labels = %+v, want %+v", report.Labels, wantLabels)
	}
	if report.CostUSD <= 0 {
		t.Fatalf("cost = %v, want the fake's reported cost", report.CostUSD)
	}

	for _, call := range fake.Calls() {
		if _, ok := call.Questions["role"]; ok && !strings.Contains(call.StateString(), testDraft.Purpose) {
			t.Fatalf("role judged without the draft in the state: %s", call.StateString())
		}
	}
	previews, err := review.LoadPreviews(root)
	if err != nil || len(previews) != 1 {
		t.Fatalf("previews = %+v, %v", previews, err)
	}
	if previews[0].Contract != testDraft.Hash() || len(previews[0].Subjects) != 5 {
		t.Fatalf("preview ledger = %+v", previews[0])
	}

	// A second preview of the same draft is served from the receipts.
	before := len(fake.Calls())
	if _, err := Preview(context.Background(), s, testDraft, PreviewOptions{}); err != nil {
		t.Fatalf("second Preview: %v", err)
	}
	if after := len(fake.Calls()); after != before {
		t.Fatalf("second preview made %d new calls, want 0", after-before)
	}
}

func TestPreviewSamplesLargeTopics(t *testing.T) {
	t.Parallel()
	vault, _ := previewTopic(t)
	fake := fakeServer(t, previewDecide, nil)
	s := openTestSession(t, vault, fake)

	report, err := Preview(context.Background(), s, testDraft, PreviewOptions{SampleSize: 2, TopN: 1})
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	if !report.Sampled || report.Judged != 2 || len(fake.Calls()) != 2 || len(report.Top) != 1 {
		t.Fatalf("report = %+v, calls = %d", report, len(fake.Calls()))
	}
	folders := map[string]int{}
	for _, folder := range report.Folders {
		folders[folder.Folder] = folder.Judged
	}
	if !reflect.DeepEqual(folders, map[string]int{"raw/articles": 1, "raw/news": 1}) {
		t.Fatalf("sample per folder = %v", folders)
	}
}
