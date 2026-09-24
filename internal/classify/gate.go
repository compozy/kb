// Package classify judges every document of a topic once on a fixed set of
// facets (spec §8) and writes them as kb-owned frontmatter: genre, depth,
// relevance, quality, concepts, plus the generated literals summary,
// entities, questions and, for articles, criterion and aliases. It also owns
// the shared gate judgment (relevance + quality in one request) that ingest
// gates, the contract impact preview and classify all ask, the backfill
// review queues (recapture, remove), concept proposals, and the vocabulary
// bootstrap for topics without articles (spec §5.2).
package classify

import (
	"context"
	"fmt"
	"slices"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/session"
)

// ExcerptTokens bounds the document excerpt of every classify and gate state
// (spec §7 stage 5: a section-aware excerpt of at most 12k tokens).
const ExcerptTokens = 12000

// Relevance roles (the `relevance` values).
const (
	RoleCore               = "core"
	RoleAdjacent           = "adjacent"
	RoleCollectedOnPurpose = "collected_on_purpose"
	RoleGeneral            = "general"
	RoleOffTopic           = "off_topic"
	RoleUnknown            = "unknown"
)

// Quality nouls of the quality bank.
const (
	NoulPaywall   = "paywall_or_login"
	NoulErrorPage = "error_or_placeholder_page"
	NoulThin      = "thin_or_boilerplate"
	NoulNoSpeech  = "no_speech_content"
)

// Quality reasons written to `quality` (and used as `triage_reason`).
const (
	ReasonPaywall      = "paywall"
	ReasonErrorPage    = quality.ErrorPage
	ReasonThin         = quality.Thin
	ReasonNotAnArticle = quality.NotAnArticle
	ReasonNoSpeech     = "no_speech"
)

// qualityNouls maps each quality noul to its reason, in the order reasons
// are preferred on ties.
var qualityNouls = []struct{ id, reason string }{
	{NoulErrorPage, ReasonErrorPage},
	{NoulPaywall, ReasonPaywall},
	{NoulThin, ReasonThin},
	{NoulNoSpeech, ReasonNoSpeech},
}

// Built-in banks. MustLoad cannot fail for embedded banks (the questions
// package tests load every one of them).
var (
	relevanceBank = questions.MustLoad("relevance")
	qualityBank   = questions.MustLoad("quality")
	classifyBank  = questions.MustLoad("classify")
	conceptBank   = questions.MustLoad("concept")
	gateBank      = questions.Compose(relevanceBank, qualityBank)
	facetBank     = questions.Compose(classifyBank, conceptBank)
)

// GateBank is the composite bank (relevance + quality) of the gate judgment.
func GateBank() *questions.Bank { return gateBank }

// FacetBank is the composite bank (classify + concept) of the classification
// request.
func FacetBank() *questions.Bank { return facetBank }

// SourceBanks are the banks recorded in the state row of a classified source
// (plan: `{relevance, quality, classify, concept}`).
func SourceBanks() []*questions.Bank {
	return []*questions.Bank{relevanceBank, qualityBank, classifyBank, conceptBank}
}

// ArticleBanks are the banks recorded in the state row of a classified
// article (plan: `{classify, concept}`).
func ArticleBanks() []*questions.Bank {
	return []*questions.Bank{classifyBank, conceptBank}
}

// IsTranscriptKind reports whether a source kind is a transcript, the only
// documents asked `no_speech_content`.
func IsTranscriptKind(sourceKind string) bool {
	switch models.SourceKind(sourceKind) {
	case models.SourceKindYouTubeTranscript, models.SourceKindInstagramVideo:
		return true
	}
	return false
}

// GateOptions configures JudgeGate.
type GateOptions struct {
	// IsTranscript adds the `no_speech_content` noul.
	IsTranscript bool
	// Flags are the code quality flags already computed by the caller
	// (quality.Check); they are carried into the judgment and win over the
	// nouls in QualityReason.
	Flags []quality.Flag
	// Contract overrides the session's active contract (the impact preview
	// judges a draft). Its hash replaces the session's in the cache key, so
	// the receipts are reused once the same draft is accepted.
	Contract *contract.Contract
}

// GateJudgment is the outcome of the shared gate request.
type GateJudgment struct {
	// Role is the argmax role, RoleCollectedOnPurpose when the path matches
	// `collected_on_purpose_paths` (no call), or "" when relevance is off or
	// the role answer is undecided.
	Role string `json:"role,omitempty"`
	// RoleAsked reports whether the `role` question was in the request.
	RoleAsked bool `json:"role_asked"`
	// RoleByPath reports that Role came from the path globs.
	RoleByPath bool `json:"role_by_path,omitempty"`
	// RoleProbs is the role distribution when decided.
	RoleProbs map[string]float64 `json:"role_probs,omitempty"`
	// POffTopic is P(off_topic).
	POffTopic float64 `json:"p_off_topic"`
	// PKept is P(core) + P(adjacent) + P(collected_on_purpose).
	PKept float64 `json:"p_kept"`
	// RoleDecided reports a usable role (decided answer or path glob).
	RoleDecided bool `json:"role_decided"`
	// Quality maps each decided quality noul id to P(yes).
	Quality map[string]float64 `json:"quality"`
	// QualityDecided reports that every quality noul was decided.
	QualityDecided bool `json:"quality_decided"`
	// Flags are the code flags passed in GateOptions.
	Flags []quality.Flag `json:"flags,omitempty"`
	// Undecided lists "<question>:<reason>" for every answer that was not
	// decided (not-checked answers read "<question>:not_checked:<reason>").
	Undecided []string `json:"undecided,omitempty"`
	// ReceiptKey is the receipts key of the request.
	ReceiptKey string `json:"receipt_key,omitempty"`
}

// QualityReason returns the reason code would act on at threshold: the first
// code flag (probability 1), else the quality noul with the highest P(yes) ≥
// threshold. It returns "" when nothing fired.
func (g GateJudgment) QualityReason(threshold float64) (string, float64) {
	if len(g.Flags) > 0 {
		return g.Flags[0].Code, 1
	}
	best, bestP := "", -1.0
	for _, noul := range qualityNouls {
		p, ok := g.Quality[noul.id]
		if !ok || p < threshold || p <= bestP {
			continue
		}
		best, bestP = noul.reason, p
	}
	if best == "" {
		return "", 0
	}
	return best, bestP
}

// QualityQuestion returns the question id behind a reason ("code:<flag>" for
// a code flag), for review items.
func (g GateJudgment) QualityQuestion(reason string) string {
	for _, flag := range g.Flags {
		if flag.Code == reason {
			return "code:" + flag.Code
		}
	}
	for _, noul := range qualityNouls {
		if noul.reason == reason {
			return noul.id
		}
	}
	return reason
}

// JudgeGate asks the shared gate request for one document: the `role`
// choice against the contract (omitted when relevance is off in topic.yaml
// or the path matches `collected_on_purpose_paths`) plus the quality nouls
// (`no_speech_content` only for transcripts), over the state
// `{contract, document}` built by corpus.DocumentState. The contract is only
// in the state when `role` is asked; without an accepted contract it is `{}`
// (relevance then runs in shadow, decided by the caller's mode). Ingest
// gates, the impact preview and classify all call this, so the same content
// reuses one receipt. The error is fatal only (auth, cancellation, invalid
// request); failed answers are listed in Undecided.
func JudgeGate(ctx context.Context, s *session.Session, doc *corpus.Document, opts GateOptions) (GateJudgment, error) {
	judgment := GateJudgment{Quality: map[string]float64{}, Flags: slices.Clone(opts.Flags)}
	active := s.Contract
	ref := s.Ref
	if opts.Contract != nil {
		normalized := opts.Contract.Normalized()
		active = &normalized
		ref.Contract = normalized.Hash()
	}

	askRole := s.RelevanceEnabled()
	if askRole && active.MatchesCollectedPath(doc.Path) {
		askRole = false
		judgment.Role = RoleCollectedOnPurpose
		judgment.RoleByPath = true
		judgment.RoleDecided = true
		judgment.PKept = 1
	}

	state := map[string]any{"document": corpus.DocumentState(doc, ExcerptTokens, nil)}
	qs := make([]questions.Q, 0, 5)
	purpose := decisions.PurposeQuality
	if askRole {
		state["contract"] = contractState(active)
		qs = append(qs, gateBank.MustQuestion("role", nil))
		purpose = decisions.PurposeRelevance
		judgment.RoleAsked = true
	}
	for _, noul := range qualityNouls {
		if noul.id == NoulNoSpeech && !opts.IsTranscript {
			continue
		}
		qs = append(qs, gateBank.MustQuestion(noul.id, nil))
	}

	// Topic extra questions (spec §4.3) ride along; their answers are only
	// recorded in receipts.
	bank := gateBank
	var err error
	if askRole {
		if bank, qs, err = s.WithExtras(decisions.PurposeRelevance, bank, qs); err != nil {
			return judgment, fmt.Errorf("classify: gate judgment for %s: %w", doc.Path, err)
		}
	}
	if bank, qs, err = s.WithExtras(decisions.PurposeQuality, bank, qs); err != nil {
		return judgment, fmt.Errorf("classify: gate judgment for %s: %w", doc.Path, err)
	}
	request := decisions.Request{Topic: ref, Purpose: purpose, Subject: doc.Path, Bank: bank, State: state, Questions: qs}
	result, err := s.Engine.Decide(ctx, request)
	if err != nil {
		return judgment, fmt.Errorf("classify: gate judgment for %s: %w", doc.Path, err)
	}
	judgment.ReceiptKey = result.ReceiptID

	if askRole {
		answer := result.Answers["role"]
		if answer.Decided() {
			judgment.RoleDecided = true
			judgment.Role = answer.Choice
			judgment.RoleProbs = answer.Probs
			judgment.POffTopic = answer.Probs[RoleOffTopic]
			judgment.PKept = answer.Probs[RoleCore] + answer.Probs[RoleAdjacent] + answer.Probs[RoleCollectedOnPurpose]
		} else {
			judgment.Undecided = append(judgment.Undecided, "role:"+answerReason(answer))
		}
	}
	judgment.QualityDecided = true
	for _, q := range qs {
		if q.ID == "role" {
			continue
		}
		answer := result.Answers[q.ID]
		if p, ok := answer.P(""); ok {
			judgment.Quality[q.ID] = p
			continue
		}
		judgment.QualityDecided = false
		judgment.Undecided = append(judgment.Undecided, q.ID+":"+answerReason(answer))
	}
	return judgment, nil
}

// contractState is the contract fragment of a state: `{}` without a
// contract (plan: relevance then runs in shadow against an empty scope).
func contractState(c *contract.Contract) map[string]any {
	if c.Empty() {
		return map[string]any{}
	}
	return c.QuestionState()
}

// answerReason names why an answer is not usable.
func answerReason(answer decisions.Answer) string {
	reason := answer.Reason
	if reason == "" {
		reason = decisions.ReasonMissingAnswer
	}
	if answer.Status == decisions.StatusNotChecked {
		return string(decisions.StatusNotChecked) + ":" + reason
	}
	return reason
}
