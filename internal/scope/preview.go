package scope

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

// Preview defaults (spec §5.1).
const (
	// DefaultSampleSize is the source count above which the preview judges a
	// stratified sample.
	DefaultSampleSize = 500
	// DefaultTopN is the length of the highest-P(off_topic) list.
	DefaultTopN = 20
)

// Preview bands of one document.
const (
	BandKept       = "kept"
	BandReview     = "review"
	BandQuarantine = "quarantine"
	BandUndecided  = "undecided"
)

// PreviewOptions configures Preview.
type PreviewOptions struct {
	// SampleSize caps the judged sources (default 500); larger topics are
	// sampled with Stratify.
	SampleSize int
	// TopN is the length of the highest-P(off_topic) list (default 20).
	TopN int
	// OnEstimate, when set, is called with the planned judgments and their
	// estimated cost before any call is made.
	OnEstimate func(documents int, estimateUSD float64)
}

func (o PreviewOptions) withDefaults() PreviewOptions {
	if o.SampleSize <= 0 {
		o.SampleSize = DefaultSampleSize
	}
	if o.TopN <= 0 {
		o.TopN = DefaultTopN
	}
	return o
}

// BandCounts counts documents per preview band.
type BandCounts struct {
	Kept       int `json:"kept"`
	Review     int `json:"review"`
	Quarantine int `json:"quarantine"`
	Undecided  int `json:"undecided"`
}

func (b *BandCounts) add(band string) {
	switch band {
	case BandKept:
		b.Kept++
	case BandReview:
		b.Review++
	case BandQuarantine:
		b.Quarantine++
	default:
		b.Undecided++
	}
}

// FolderBands are the band counts of one raw/ subfolder.
type FolderBands struct {
	Folder string `json:"folder"`
	Judged int    `json:"judged"`
	BandCounts
}

// TopDocument is one of the highest-P(off_topic) documents.
type TopDocument struct {
	Path      string  `json:"path"`
	Title     string  `json:"title"`
	POffTopic float64 `json:"p_off_topic"`
	Band      string  `json:"band"`
	// Reason is the quality reason when quality drove the band.
	Reason string `json:"reason,omitempty"`
}

// LabelScore is the precision and recall of the would-be relevance
// quarantine against the topic's relevance labels (negative label = should
// be quarantined).
type LabelScore struct {
	// Labels counts the labeled subjects that were judged with a decided
	// role; Negatives counts the negative ones among them.
	Labels    int `json:"labels"`
	Negatives int `json:"negatives"`
	// Predicted counts labeled subjects at P(off_topic) ≥ relevance_quarantine.
	Predicted    int     `json:"predicted"`
	TruePositive int     `json:"true_positive"`
	Precision    float64 `json:"precision"`
	Recall       float64 `json:"recall"`
}

// PreviewReport is the outcome of an impact preview.
type PreviewReport struct {
	// Contract is the hash of the previewed draft.
	Contract string `json:"contract"`
	// Sources is the topic's source count; Judged the documents judged.
	Sources int  `json:"sources"`
	Judged  int  `json:"judged"`
	Sampled bool `json:"sampled"`
	// RelevanceOff reports topic.yaml decisions.relevance: off (quality only).
	RelevanceOff bool               `json:"relevance_off,omitempty"`
	Thresholds   map[string]float64 `json:"thresholds"`
	Bands        BandCounts         `json:"bands"`
	Folders      []FolderBands      `json:"folders"`
	Top          []TopDocument      `json:"top"`
	Labels       *LabelScore        `json:"labels,omitempty"`
	Undecided    map[string]int     `json:"undecided_reasons,omitempty"`
	EstimateUSD  float64            `json:"estimate_usd"`
	CostUSD      float64            `json:"cost_usd"`
	judgments    map[string]judgment
}

// judgment is the preview outcome of one document.
type judgment struct {
	doc    *corpus.Document
	gate   classify.GateJudgment
	band   string
	reason string
}

// Preview runs the impact preview of a draft (spec §5.1): the shared gate
// judgment (classify.JudgeGate with GateOptions.Contract = draft, so the
// receipts are reused by `kb classify` once the draft is accepted unchanged)
// over every source, or a stratified sample of opts.SampleSize by raw/
// subfolder. Each document lands in one band: quarantine (P(off_topic) ≥
// relevance_quarantine or a quality reason at quality_apply, code flags
// included), review (P(off_topic) ≥ relevance_review or a quality reason at
// quality_review), undecided (an answer is missing and nothing fired) or
// kept. When the topic has relevance labels, the would-be relevance
// quarantine is scored against them. The preview is recorded in
// previews.jsonl (the top list plus every label used) so a later contract
// change demotes those labels to dev.
func Preview(ctx context.Context, s *session.Session, draft *contract.Contract, opts PreviewOptions) (PreviewReport, error) {
	opts = opts.withDefaults()
	normalized := draft.Normalized()
	report := PreviewReport{
		Contract:     normalized.Hash(),
		RelevanceOff: !s.RelevanceEnabled(),
		Thresholds: map[string]float64{
			"relevance_quarantine": s.Threshold("relevance_quarantine"),
			"relevance_review":     s.Threshold("relevance_review"),
			"quality_apply":        s.Threshold("quality_apply"),
			"quality_review":       s.Threshold("quality_review"),
		},
		Undecided: map[string]int{},
		judgments: map[string]judgment{},
	}
	c, err := s.Corpus()
	if err != nil {
		return report, fmt.Errorf("scope: preview: load corpus: %w", err)
	}
	sources := c.Sources()
	report.Sources = len(sources)
	docs := Stratify(sources, opts.SampleSize)
	report.Sampled = len(docs) < len(sources)
	report.EstimateUSD = float64(len(docs)) * CostPerDocumentUSD
	if opts.OnEstimate != nil {
		opts.OnEstimate(len(docs), report.EstimateUSD)
	}

	hostLines := quality.HostLineCounts(sources)
	costBefore := s.Engine.Summary().CostUSD
	var mu sync.Mutex
	err = forEach(ctx, s.Engine.Concurrency(), docs, func(ctx context.Context, doc *corpus.Document) error {
		gate, err := classify.JudgeGate(ctx, s, doc, classify.GateOptions{
			IsTranscript: classify.IsTranscriptKind(doc.SourceKind()),
			Flags:        classify.CodeFlags(doc, hostLines),
			Contract:     &normalized,
		})
		if err != nil {
			return err
		}
		band, reason := bandOf(gate, report.Thresholds)
		mu.Lock()
		report.judgments[doc.Path] = judgment{doc: doc, gate: gate, band: band, reason: reason}
		mu.Unlock()
		return nil
	})
	report.CostUSD = s.Engine.Summary().CostUSD - costBefore
	if err != nil {
		return report, fmt.Errorf("scope: preview: %w", err)
	}
	report.Judged = len(report.judgments)
	report.tally(opts.TopN)

	labels, err := review.LoadLabels(s.Root())
	if err != nil {
		return report, fmt.Errorf("scope: preview: %w", err)
	}
	score, used := scoreLabels(report.judgments, labels, report.Thresholds["relevance_quarantine"])
	report.Labels = score

	subjects := make([]string, 0, len(report.Top)+len(used))
	for _, top := range report.Top {
		subjects = append(subjects, top.Path)
	}
	subjects = append(subjects, used...)
	slices.Sort(subjects)
	subjects = slices.Compact(subjects)
	if err := review.AppendPreview(s.Root(), review.Preview{
		Time:     s.Now().UTC().Format(time.RFC3339),
		Contract: report.Contract,
		Subjects: subjects,
	}); err != nil {
		return report, fmt.Errorf("scope: preview: %w", err)
	}
	return report, nil
}

// bandOf places one gate judgment in a preview band.
func bandOf(g classify.GateJudgment, thresholds map[string]float64) (string, string) {
	offTopic := g.RoleAsked && g.RoleDecided && !g.RoleByPath
	if offTopic && g.POffTopic >= thresholds["relevance_quarantine"] {
		return BandQuarantine, classify.RoleOffTopic
	}
	if reason, _ := g.QualityReason(thresholds["quality_apply"]); reason != "" {
		return BandQuarantine, reason
	}
	if offTopic && g.POffTopic >= thresholds["relevance_review"] {
		return BandReview, classify.RoleOffTopic
	}
	if reason, _ := g.QualityReason(thresholds["quality_review"]); reason != "" {
		return BandReview, reason
	}
	if (g.RoleAsked && !g.RoleDecided) || !g.QualityDecided {
		return BandUndecided, ""
	}
	return BandKept, ""
}

// tally fills the band counts, folder table, undecided reasons and top list.
func (r *PreviewReport) tally(topN int) {
	folders := map[string]*FolderBands{}
	all := make([]judgment, 0, len(r.judgments))
	for _, j := range r.judgments {
		all = append(all, j)
	}
	slices.SortFunc(all, func(a, b judgment) int { return cmp.Compare(a.doc.Path, b.doc.Path) })
	for _, j := range all {
		r.Bands.add(j.band)
		name := Folder(j.doc.Path)
		folder := folders[name]
		if folder == nil {
			folder = &FolderBands{Folder: name}
			folders[name] = folder
		}
		folder.Judged++
		folder.add(j.band)
		for _, undecided := range j.gate.Undecided {
			r.Undecided[undecided]++
		}
	}
	r.Folders = make([]FolderBands, 0, len(folders))
	for _, folder := range folders {
		r.Folders = append(r.Folders, *folder)
	}
	slices.SortFunc(r.Folders, func(a, b FolderBands) int {
		return cmp.Or(cmp.Compare(b.Judged, a.Judged), cmp.Compare(a.Folder, b.Folder))
	})

	ranked := slices.Clone(all)
	slices.SortStableFunc(ranked, func(a, b judgment) int {
		return cmp.Or(cmp.Compare(b.gate.POffTopic, a.gate.POffTopic), cmp.Compare(a.doc.Path, b.doc.Path))
	})
	r.Top = make([]TopDocument, 0, topN)
	for _, j := range ranked {
		if len(r.Top) == topN {
			break
		}
		if !j.gate.RoleDecided || j.gate.RoleByPath {
			continue
		}
		top := TopDocument{Path: j.doc.Path, Title: j.doc.Title, POffTopic: j.gate.POffTopic, Band: j.band}
		if j.reason != classify.RoleOffTopic {
			top.Reason = j.reason
		}
		r.Top = append(r.Top, top)
	}
	if len(r.Undecided) == 0 {
		r.Undecided = nil
	}
}

// scoreLabels scores the would-be relevance quarantine (P(off_topic) ≥
// threshold) against the latest relevance label of every judged subject
// with a decided role (negative = should be quarantined). It returns nil
// without such labels, and the subjects used.
func scoreLabels(judgments map[string]judgment, labels []review.Label, threshold float64) (*LabelScore, []string) {
	latest := map[string]review.Label{}
	for _, label := range labels {
		if label.Purpose == review.PurposeRelevance {
			latest[label.Subject] = label
		}
	}
	score := &LabelScore{}
	used := make([]string, 0)
	for subject, label := range latest {
		j, ok := judgments[subject]
		if !ok || !j.gate.RoleDecided || j.gate.RoleByPath {
			continue
		}
		used = append(used, subject)
		score.Labels++
		negative := label.Verdict == review.VerdictNegative
		predicted := j.gate.POffTopic >= threshold
		if negative {
			score.Negatives++
		}
		if predicted {
			score.Predicted++
			if negative {
				score.TruePositive++
			}
		}
	}
	if score.Labels == 0 {
		return nil, nil
	}
	if score.Predicted > 0 {
		score.Precision = float64(score.TruePositive) / float64(score.Predicted)
	}
	if score.Negatives > 0 {
		score.Recall = float64(score.TruePositive) / float64(score.Negatives)
	}
	slices.Sort(used)
	return score, used
}
