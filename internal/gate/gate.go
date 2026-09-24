// Package gate runs the ingest gates of spec §7 over sources about to enter
// a topic. Stages are code and predicates are decisions; only survivors pay
// the next stage:
//
//  1. exact dedupe (code): normalized URL, platform id, body hash against the
//     topic's sources, quarantined ones included;
//  2. pre-fetch relevance (bulk only): one request per ≤ 20 items of
//     metadata; skipped items are never fetched;
//  3. code quality checks (internal/quality);
//  4. one fresh refetch of URL sources that a code flag or a short body
//     marks as a possibly broken capture;
//  5. post-fetch quality + relevance: the shared gate judgment
//     (classify.JudgeGate) on the would-be document;
//  6. near-duplicate: BM25 and same-domain title candidates, one request.
//
// Quality gates and exact dedupe apply by default; relevance gates run in
// shadow until calibrated (session.RelevanceGateMode). Shadow never skips
// or quarantines: the would-be outcome is reported and the source is
// written with `triage: review` so the owner can label it.
package gate

import (
	"fmt"
	"strings"

	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

// Triage outcomes. Kept, review and quarantined are also `triage` values;
// skipped (stage 2) and duplicate-skipped (stage 1) never reach the disk.
const (
	TriageKept             = "kept"
	TriageReview           = "review"
	TriageQuarantined      = "quarantined"
	TriageSkipped          = "skipped"
	TriageDuplicateSkipped = "duplicate-skipped"
)

// Gate reasons that are not quality reasons.
const (
	ReasonOffTopic  = "off_topic"
	ReasonDuplicate = "duplicate"
)

// Stage names, for outcomes and review evidence.
const (
	StageDedupe        = "dedupe"
	StagePrefetch      = "prefetch"
	StageQuality       = "quality"
	StageRelevance     = "relevance"
	StageNearDuplicate = "near-duplicate"
)

// Outcome is what the gates decided about one source.
type Outcome struct {
	// Triage is kept, review, quarantined, skipped or duplicate-skipped.
	Triage string `json:"triage"`
	// Reason is the triage_reason (off_topic, thin, error_page, paywall,
	// not_an_article, no_speech, duplicate); empty when kept.
	Reason string `json:"triage_reason,omitempty"`
	// Stage is the stage that set Triage.
	Stage string `json:"stage,omitempty"`
	// Quality is the value for the `quality` key: the quality reason when a
	// quality check or question fired at or above the review threshold.
	Quality string `json:"quality,omitempty"`
	// Purpose, Question, Probability and ReceiptKey describe the decision
	// behind Triage (for review items and labels).
	Purpose     decisions.Purpose `json:"purpose,omitempty"`
	Question    string            `json:"question,omitempty"`
	Probability float64           `json:"probability,omitempty"`
	ReceiptKey  string            `json:"receipt_key,omitempty"`
	// DuplicateOf is the existing source a duplicate matched.
	DuplicateOf string `json:"duplicate_of,omitempty"`
	// DuplicateQuarantined reports that DuplicateOf is quarantined.
	DuplicateQuarantined bool `json:"duplicate_quarantined,omitempty"`
	// Shadow lists would-be outcomes a shadow gate did not apply, such as
	// "quarantined:off_topic (relevance shadow)".
	Shadow []string `json:"shadow,omitempty"`
	// Supersedes lists `[[target]]` wikilinks of older versions (stage 6).
	Supersedes []string `json:"supersedes,omitempty"`
	// Flags are the code quality flags of the kept body.
	Flags []quality.Flag `json:"flags,omitempty"`
	// Refetched reports that the stage-4 refetch body was kept.
	Refetched bool `json:"refetched,omitempty"`
	// Undecided lists "<question>:<reason>" answers that could not be used.
	Undecided []string `json:"undecided,omitempty"`
	// Evidence is one human-readable line for review items and logs.
	Evidence string `json:"evidence,omitempty"`
}

// Kept reports whether the source is written to raw/ (kept or review).
func (o Outcome) Kept() bool {
	return o.Triage == TriageKept || o.Triage == TriageReview
}

// severity orders triage outcomes: quarantine beats review beats kept.
func severity(triage string) int {
	switch triage {
	case TriageQuarantined:
		return 2
	case TriageReview:
		return 1
	default:
		return 0
	}
}

// verdict is one stage's proposed outcome before shadow handling.
type verdict struct {
	triage      string
	reason      string
	stage       string
	purpose     decisions.Purpose
	question    string
	probability float64
	receipt     string
	evidence    string
}

// merge applies v to o: in shadow a quarantine becomes a review plus a
// Shadow note; the more severe outcome wins (quality first on ties, because
// stages are merged in order).
func (o *Outcome) merge(v verdict, mode string) {
	if v.triage == "" || v.triage == TriageKept {
		return
	}
	triage := v.triage
	if mode != session.ModeApply && triage == TriageQuarantined {
		o.Shadow = append(o.Shadow, fmt.Sprintf("%s:%s (%s shadow)", TriageQuarantined, v.reason, v.stage))
		triage = TriageReview
	}
	if severity(triage) <= severity(o.Triage) {
		return
	}
	o.Triage, o.Reason, o.Stage = triage, v.reason, v.stage
	o.Purpose, o.Question, o.Probability, o.ReceiptKey = v.purpose, v.question, v.probability, v.receipt
	o.Evidence = v.evidence
}

// ReviewItem returns the `gate` review-queue item of an outcome written at
// path (review band, shadow would-be quarantine, or gate quarantine), and
// false for a kept source. Accepting it keeps the source (restoring it from
// quarantine), rejecting it quarantines it with Reason.
func ReviewItem(o Outcome, path, title string) (review.Item, bool) {
	if o.Triage != TriageReview && o.Triage != TriageQuarantined {
		return review.Item{}, false
	}
	reason := o.Reason
	if reason == "" {
		reason = ReasonOffTopic
	}
	purpose := o.Purpose
	if purpose == "" {
		purpose = decisions.PurposeRelevance
	}
	question := o.Question
	if question == "" {
		question = o.Stage
	}
	action := map[string]any{"reason": reason, "stage": o.Stage}
	if o.Triage == TriageQuarantined {
		action["quarantined"] = true
	}
	if len(o.Shadow) > 0 {
		action["shadow"] = strings.Join(o.Shadow, "; ")
	}
	if o.DuplicateOf != "" {
		action["duplicate_of"] = o.DuplicateOf
	}
	evidence := strings.TrimSpace(title + " — " + o.Triage + ":" + reason)
	if o.Evidence != "" {
		evidence += " (" + o.Evidence + ")"
	}
	return review.Item{
		Queue: review.QueueGate, Purpose: string(purpose), Subject: path, Question: question,
		Probability: o.Probability, ReceiptKey: o.ReceiptKey, Evidence: evidence, Action: action,
	}, true
}

// Tally counts the outcomes of a run.
type Tally struct {
	Kept        int `json:"kept"`
	Review      int `json:"review"`
	Quarantined int `json:"quarantined"`
	Skipped     int `json:"skipped"`
	Duplicates  int `json:"duplicates"`
	// Shadow counts sources with at least one would-be outcome not applied.
	Shadow int `json:"shadow"`
}

// Add counts one outcome.
func (t *Tally) Add(o Outcome) {
	switch o.Triage {
	case TriageKept:
		t.Kept++
	case TriageReview:
		t.Review++
	case TriageQuarantined:
		t.Quarantined++
	case TriageSkipped:
		t.Skipped++
	case TriageDuplicateSkipped:
		t.Duplicates++
	}
	if len(o.Shadow) > 0 {
		t.Shadow++
	}
}

// Line renders the tally: "kept 2, review 1, quarantined 0, skipped 3,
// duplicates 1".
func (t Tally) Line() string {
	line := fmt.Sprintf("kept %d, review %d, quarantined %d, skipped %d, duplicates %d",
		t.Kept, t.Review, t.Quarantined, t.Skipped, t.Duplicates)
	if t.Shadow > 0 {
		line += fmt.Sprintf(" (%d with shadow outcomes)", t.Shadow)
	}
	return line
}
