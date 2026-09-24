package gate

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/firecrawl"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/session"
)

// Fetched is one fetched source about to be written, as seen by stages 3–6.
type Fetched struct {
	// Title and Markdown are the fetched title and body.
	Title    string
	Markdown string
	// Build returns the would-be document (path, frontmatter with
	// provenance, body) for a title and body. The gate calls it again after
	// a refetch, so the judged document is always the one that is written.
	Build func(title, markdown string) (*corpus.Document, error)
	// RequestedURL, FinalURL, StatusCode and SiteName are what the fetcher
	// reported (SiteName feeds the "title equals the site name" rule).
	RequestedURL string
	FinalURL     string
	StatusCode   int
	SiteName     string
	// Refetch performs the one fresh stage-4 scrape (URL sources only; nil
	// disables stage 4).
	Refetch func(ctx context.Context) (*firecrawl.ScrapeResult, error)
	// PlatformID is the `youtube:<id>` / `instagram:<shortcode>` identity
	// when the source URL does not carry it.
	PlatformID string
	// Transcript adds the `no_speech_content` question and disables the thin
	// rule.
	Transcript bool
	// PrefetchReview marks an item stage 2 put in the review band; Shadow
	// names a would-be pre-fetch skip that a shadow gate did not apply.
	PrefetchReview bool
	PrefetchShadow string
	// PrefetchPOffTopic is the pre-fetch P(off_topic) (review evidence).
	PrefetchPOffTopic float64
	// PrefetchReceipt is the pre-fetch receipt key.
	PrefetchReceipt string
}

// Checker runs the per-source stages of one ingest run over a topic. It is
// not safe for concurrent use; ingest serializes items.
type Checker struct {
	s         *session.Session
	index     *Index
	hostLines map[string]map[string]int
	// sources is the near-duplicate candidate pool: live sources that may
	// enter a model call (never a decisions.exclude match).
	sources []*corpus.Document
	near    *nearDuplicates
}

// NewChecker loads what the gates compare against: the topic's sources
// (dedupe, same-host lines, near-duplicate candidates) and the quarantined
// ones (dedupe only). Sources matching decisions.exclude are kept for exact
// dedupe but never become near-duplicate candidates (spec §14).
func NewChecker(s *session.Session) (*Checker, error) {
	c, err := s.Corpus()
	if err != nil {
		return nil, fmt.Errorf("gate: load corpus: %w", err)
	}
	sources := c.Sources()
	index, err := NewIndex(s.Root(), sources)
	if err != nil {
		return nil, err
	}
	return &Checker{
		s:         s,
		index:     index,
		hostLines: quality.HostLineCounts(sources),
		sources:   slices.DeleteFunc(slices.Clone(sources), func(doc *corpus.Document) bool { return s.Excluded(doc.Path) }),
		near:      newNearDuplicates(s),
	}, nil
}

// Skip returns the stage-1 outcome for ref when it duplicates an existing
// source and dedupe applies (quality gate mode apply): the source is not
// written. In shadow it returns false and Evaluate reports the would-be skip.
func (c *Checker) Skip(ref Ref) (Outcome, bool) {
	match, ok := c.index.Find(ref)
	if !ok || c.s.QualityGateMode() != session.ModeApply {
		return Outcome{}, false
	}
	return duplicateOutcome(match), true
}

func duplicateOutcome(match Match) Outcome {
	return Outcome{
		Triage: TriageDuplicateSkipped, Reason: ReasonDuplicate, Stage: StageDedupe,
		Purpose: decisions.PurposeDuplicate, Question: "exact:" + match.By, Probability: 1,
		DuplicateOf: match.Path, DuplicateQuarantined: match.Quarantined, Evidence: match.Describe(),
	}
}

// dedupe checks ref in Evaluate: in apply mode it returns the skip outcome;
// in shadow the source goes to review with a Shadow note.
func (c *Checker) dedupe(ref Ref, out *Outcome) (Outcome, bool) {
	match, ok := c.index.Find(ref)
	if !ok {
		return Outcome{}, false
	}
	if c.s.QualityGateMode() == session.ModeApply {
		return duplicateOutcome(match), true
	}
	out.Shadow = append(out.Shadow, TriageDuplicateSkipped+":"+match.Describe()+" (dedupe shadow)")
	out.merge(verdict{
		triage: TriageReview, reason: ReasonDuplicate, stage: StageDedupe, purpose: decisions.PurposeDuplicate,
		question: "exact:" + match.By, probability: 1, evidence: match.Describe(),
	}, session.ModeApply)
	out.DuplicateOf, out.DuplicateQuarantined = match.Path, match.Quarantined
	return Outcome{}, false
}

// Register adds a source written (or quarantined) by this run, so later
// items dedupe and near-dedupe against it. A source matching
// decisions.exclude is registered for exact dedupe only: it never becomes a
// near-duplicate candidate of a later item, since candidates are sent to
// the model (spec §14).
func (c *Checker) Register(doc *corpus.Document, quarantined bool) {
	c.index.Add(doc, quarantined)
	if !quarantined && !c.s.Excluded(doc.Path) {
		c.sources = append(c.sources, doc)
		c.near.reset()
	}
}

// Evaluate runs stages 3–6 over one fetched source and returns the outcome
// and the document to write (the refetched body when stage 4 kept it). A
// code flag that survives the refetch quarantines without any decision call
// (in apply mode). The error is fatal only (auth, cancellation, I/O).
func (c *Checker) Evaluate(ctx context.Context, in Fetched) (Outcome, *corpus.Document, error) {
	doc, err := in.Build(in.Title, in.Markdown)
	if err != nil {
		return Outcome{}, nil, err
	}
	out := Outcome{Triage: TriageKept}
	if in.PrefetchReview {
		out.merge(verdict{
			triage: TriageReview, reason: ReasonOffTopic, stage: StagePrefetch, purpose: decisions.PurposeRelevance,
			question: "role_item", probability: in.PrefetchPOffTopic, receipt: in.PrefetchReceipt,
			evidence: fmt.Sprintf("pre-fetch P(off_topic) %.2f", in.PrefetchPOffTopic),
		}, session.ModeApply)
	}
	if in.PrefetchShadow != "" {
		out.Shadow = append(out.Shadow, fmt.Sprintf("%s:%s (%s shadow)", TriageSkipped, ReasonOffTopic, StagePrefetch))
	}

	if skip, ok := c.dedupe(Ref{URL: doc.SourceURL(), PlatformID: in.PlatformID, BodyHash: doc.BodyHash}, &out); ok {
		return skip, doc, nil
	}

	fetch := fetchInfo{requested: in.RequestedURL, final: in.FinalURL, status: in.StatusCode, site: in.SiteName}
	flags := c.flags(doc, fetch, in.Transcript)
	if NeedsRefetch(flags, c.words(doc), in.Refetch != nil) {
		refetched, info, ok := c.refetch(ctx, in, doc, &out)
		if ok {
			doc, fetch, out.Refetched = refetched, info, true
			if skip, dup := c.dedupe(Ref{BodyHash: doc.BodyHash}, &out); dup {
				skip.Refetched = true
				return skip, doc, nil
			}
		}
		flags = c.flags(doc, fetch, in.Transcript)
	}
	out.Flags = flags

	qualityMode := c.s.QualityGateMode()
	if len(flags) > 0 {
		out.Quality = flags[0].Code
		out.merge(verdict{
			triage: TriageQuarantined, reason: flags[0].Code, stage: StageQuality, purpose: decisions.PurposeQuality,
			question: "code:" + flags[0].Code, probability: 1, evidence: flags[0].Detail,
		}, qualityMode)
		if qualityMode == session.ModeApply {
			return out, doc, nil
		}
	}

	if c.s.Excluded(doc.Path) {
		// decisions.exclude keeps the source out of every call (spec §14):
		// no judgment, no near-duplicate request. It is kept without one.
		out.Excluded = true
		if out.Evidence == "" {
			out.Evidence = ExcludedNote
		}
		return out, doc, nil
	}

	judgment, err := classify.JudgeGate(ctx, c.s, doc, classify.GateOptions{IsTranscript: in.Transcript, Flags: flags})
	if err != nil {
		return out, doc, fmt.Errorf("gate: %w", err)
	}
	out.Undecided = append(out.Undecided, judgment.Undecided...)
	if len(flags) == 0 {
		out.mergeQuality(judgment, c.s.Threshold("quality_apply"), c.s.Threshold("quality_review"), qualityMode)
	}
	out.mergeRelevance(judgment, c.s.Threshold("relevance_quarantine"), c.s.Threshold("relevance_review"), c.s.RelevanceGateMode())
	out.mergeUndecided(judgment, len(flags) > 0)

	if out.Triage == TriageQuarantined {
		return out, doc, nil
	}
	if err := c.near.check(ctx, doc, c.sources, &out, qualityMode); err != nil {
		return out, doc, err
	}
	return out, doc, nil
}

// NeedsRefetch is the stage-4 trigger: a URL source (refetchable) with a
// code flag or fewer than quality.MinWords content words.
func NeedsRefetch(flags []quality.Flag, words int, refetchable bool) bool {
	return refetchable && (len(flags) > 0 || words < quality.MinWords)
}

// mergeQuality bands the quality nouls: ≥ apply → quarantine, ≥ review →
// review, with the noul's reason; `quality` records either.
func (o *Outcome) mergeQuality(j classify.GateJudgment, apply, reviewAt float64, mode string) {
	triage := TriageQuarantined
	reason, p := j.QualityReason(apply)
	if reason == "" {
		triage = TriageReview
		reason, p = j.QualityReason(reviewAt)
	}
	if reason == "" {
		return
	}
	o.Quality = reason
	o.merge(verdict{
		triage: triage, reason: reason, stage: StageQuality, purpose: decisions.PurposeQuality,
		question: j.QualityQuestion(reason), probability: p, receipt: j.ReceiptKey,
		evidence: fmt.Sprintf("P(%s) %.2f", j.QualityQuestion(reason), p),
	}, mode)
}

// ExcludedNote is the evidence line of a source kept without a judgment
// because it matches decisions.exclude.
const ExcludedNote = "decisions.exclude: kept without a judgment (no decision call)"

// mergeUndecided puts a source whose gate judgment failed in the review
// band with triage_reason `undecided` (spec §2 principle 3: a failed
// judgment is never a "no", and never a "kept" either): the `role` answer
// when it was asked, or a quality noul when no code flag already decided
// quality. Answers not checked because of decisions.exclude are no
// judgment, not a failure. Any other review or quarantine verdict wins.
func (o *Outcome) mergeUndecided(j classify.GateJudgment, flagged bool) {
	failed := make([]string, 0, len(j.Undecided))
	purpose := decisions.PurposeQuality
	for _, entry := range j.Undecided {
		if strings.HasSuffix(entry, ":"+decisions.ReasonExcluded) {
			continue
		}
		question, _, _ := strings.Cut(entry, ":")
		if question == "role" {
			if !j.RoleAsked {
				continue
			}
			purpose = decisions.PurposeRelevance
		} else if flagged {
			continue
		}
		failed = append(failed, entry)
	}
	if len(failed) == 0 {
		return
	}
	question, _, _ := strings.Cut(failed[0], ":")
	o.merge(verdict{
		triage: TriageReview, reason: ReasonUndecided, stage: StageUndecided, purpose: purpose,
		question: question, receipt: j.ReceiptKey, evidence: "undecided: " + strings.Join(failed, ", "),
	}, session.ModeApply)
}

// mergeRelevance bands P(off_topic): ≥ quarantine → quarantine off_topic,
// ≥ review → review. A role from the path globs, relevance off or an
// undecided role never moves anything.
func (o *Outcome) mergeRelevance(j classify.GateJudgment, quarantine, reviewAt float64, mode string) {
	if !j.RoleAsked || !j.RoleDecided {
		return
	}
	var triage string
	switch {
	case j.POffTopic >= quarantine:
		triage = TriageQuarantined
	case j.POffTopic >= reviewAt:
		triage = TriageReview
	default:
		return
	}
	o.merge(verdict{
		triage: triage, reason: ReasonOffTopic, stage: StageRelevance, purpose: decisions.PurposeRelevance,
		question: "role", probability: j.POffTopic, receipt: j.ReceiptKey,
		evidence: fmt.Sprintf("P(off_topic) %.2f", j.POffTopic),
	}, mode)
}

// fetchInfo is what the fetcher reported about the kept body.
type fetchInfo struct {
	requested, final, site string
	status                 int
}

func (c *Checker) flags(doc *corpus.Document, fetch fetchInfo, transcript bool) []quality.Flag {
	host, _ := doc.Provenance()["source_host"].(string)
	return quality.Check(quality.Input{
		Title:        doc.Title,
		Body:         doc.Body,
		SourceURL:    doc.SourceURL(),
		FinalURL:     fetch.final,
		RequestedURL: fetch.requested,
		StatusCode:   fetch.status,
		SiteName:     fetch.site,
		HostLines:    c.hostLines[host],
		SkipThin:     !webCapture(doc, transcript),
	})
}

func (c *Checker) words(doc *corpus.Document) int {
	host, _ := doc.Provenance()["source_host"].(string)
	return quality.ContentWords(doc.Body, c.hostLines[host])
}

// refetch runs the one stage-4 scrape and returns the refetched document
// when its body is longer than the current one. The refetched document
// carries the title the refetch reported (the first capture's title only
// when the refetch has none), so the quality checks judge the title and the
// body of the same fetch.
func (c *Checker) refetch(ctx context.Context, in Fetched, doc *corpus.Document, out *Outcome) (*corpus.Document, fetchInfo, bool) {
	result, err := in.Refetch(ctx)
	if err != nil || result == nil {
		reason := "empty result"
		if err != nil {
			reason = err.Error()
		}
		out.Undecided = append(out.Undecided, "refetch:"+reason)
		return nil, fetchInfo{}, false
	}
	if len(strings.TrimSpace(result.Markdown)) <= len(strings.TrimSpace(doc.Body)) {
		return nil, fetchInfo{}, false
	}
	refetched, err := in.Build(firstNonEmpty(result.Title, in.Title), result.Markdown)
	if err != nil {
		out.Undecided = append(out.Undecided, "refetch:"+err.Error())
		return nil, fetchInfo{}, false
	}
	return refetched, fetchInfo{requested: in.RequestedURL, final: result.FinalURL, status: result.StatusCode, site: firstNonEmpty(result.SiteName, in.SiteName)}, true
}

// webCapture reports whether the thin rule applies: an http(s) capture that
// is not a transcript or a bookmark cluster (same rule as classify).
func webCapture(doc *corpus.Document, transcript bool) bool {
	sourceURL := strings.ToLower(doc.SourceURL())
	if !strings.HasPrefix(sourceURL, "http://") && !strings.HasPrefix(sourceURL, "https://") {
		return false
	}
	return !transcript && !classify.IsTranscriptKind(doc.SourceKind()) && doc.SourceKind() != string(models.SourceKindBookmarkCluster)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
