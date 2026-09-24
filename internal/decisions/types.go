// Package decisions is kb's decision engine: one entry point (Engine.Decide)
// that owns transport to the decision model, receipt validation, the
// per-topic receipts cache, banding, budget, concurrency and run summaries.
package decisions

import "github.com/compozy/kb/internal/questions"

// Purpose names why a decision is asked; it selects thresholds and groups the
// run summary.
type Purpose string

// Purposes asked by kb.
const (
	PurposeRelevance     Purpose = "relevance"
	PurposeQuality       Purpose = "quality"
	PurposeDuplicate     Purpose = "duplicate"
	PurposeClassify      Purpose = "classify"
	PurposeConcept       Purpose = "concept"
	PurposeLink          Purpose = "link"
	PurposeMention       Purpose = "mention"
	PurposeFind          Purpose = "find"
	PurposeOKFType       Purpose = "okf_type"
	PurposeContractCheck Purpose = "contract_check"
)

// Status is the outcome of one answer. Failure is never a "no".
type Status string

// Answer statuses.
const (
	StatusDecided    Status = "decided"
	StatusUndecided  Status = "undecided"
	StatusNotChecked Status = "not_checked"
)

// Band is the confidence band of a decided answer.
type Band string

// Confidence bands.
const (
	BandApply  Band = "apply"
	BandReview Band = "review"
	BandIgnore Band = "ignore"
)

// Undecided / not-checked reasons.
const (
	ReasonTimeout         = "timeout"
	ReasonRetries         = "retries"
	ReasonBudget          = "budget"
	ReasonInvalidReceipt  = "invalid_receipt"
	ReasonStateTooLarge   = "state_too_large"
	ReasonHTTPStatus      = "http_status"
	ReasonMissingAnswer   = "missing_answer"
	ReasonExcluded        = "excluded"
	ReasonContextCanceled = "canceled"
)

// TopicRef identifies the topic a decision belongs to: receipts live under
// Root/.decisions/, and Contract (hash) enters every cache key.
type TopicRef struct {
	Slug       string
	Root       string
	Contract   string
	Thresholds Thresholds
}

// Request is one decision request: one state, many questions.
type Request struct {
	Topic     TopicRef
	Purpose   Purpose
	Subject   string
	Bank      *questions.Bank
	State     any
	Questions []questions.Q
}

// Answer is one banded answer. Noul/Score/Confidence are pointers because
// absence differs from zero.
type Answer struct {
	Type       string             `json:"type"`
	Status     Status             `json:"status"`
	Reason     string             `json:"reason,omitempty"`
	Band       Band               `json:"band,omitempty"`
	Noul       *float64           `json:"noul,omitempty"`
	Choice     string             `json:"choice,omitempty"`
	Score      *float64           `json:"score,omitempty"`
	Probs      map[string]float64 `json:"probabilities,omitempty"`
	Confidence *float64           `json:"confidence,omitempty"`
}

// Decided reports whether the answer carries a usable judgment.
func (a Answer) Decided() bool { return a.Status == StatusDecided }

// P returns P(yes) for a decided noul, the probability of option for a
// decided choice, or 0 with ok=false.
func (a Answer) P(option string) (float64, bool) {
	if !a.Decided() {
		return 0, false
	}
	if a.Noul != nil && option == "" {
		return *a.Noul, true
	}
	if a.Probs != nil {
		p, ok := a.Probs[option]
		return p, ok
	}
	return 0, false
}

// Result carries the answers of one request, keyed by question id.
type Result struct {
	Answers   map[string]Answer
	ReceiptID string
	CostUSD   float64
	CacheHit  bool
}

// OwnedKeys are the frontmatter keys kb writes (spec §6). The writer,
// ingest's extra-frontmatter guard and lint all use this list. `locked` is
// user-owned and never written.
var OwnedKeys = []string{
	"triage", "triage_reason", "genre", "depth", "relevance", "concepts",
	"summary", "criterion", "entities", "questions", "quality",
	"ingest_batch", "ingest_query",
	"related", "extends", "prerequisite", "example_of", "contradicts",
	"affects", "supersedes", "aliases", "recaptured",
}

// RelationKeys are the typed relation lists written by link.
var RelationKeys = []string{"related", "extends", "prerequisite", "example_of", "contradicts"}

// IsOwnedKey reports whether key is kb-owned.
func IsOwnedKey(key string) bool {
	for _, owned := range OwnedKeys {
		if owned == key {
			return true
		}
	}
	return false
}

// Thresholds maps named thresholds (link_apply, relevance_quarantine, ...)
// to values. Missing names fall back to DefaultThresholds.
type Thresholds map[string]float64

// DefaultThresholds are the measured defaults of spec §12.3.
func DefaultThresholds() Thresholds {
	return Thresholds{
		"link_apply":           0.85,
		"link_review":          0.60,
		"mention_sense":        0.85,
		"affects":              0.80,
		"relevance_quarantine": 0.80,
		"relevance_review":     0.50,
		"relevance_fetch":      0.50,
		"quality_apply":        0.80,
		"quality_review":       0.50,
		"duplicate":            0.80,
		"primary_concept":      0.30,
		"concept_noul":         0.70,
		"find_keep":            0.50,
		"find_window":          0.35,
		"okf_type":             0.80,
		"contract_conflict":    0.70,
	}
}

// Get returns the named threshold, falling back to the default.
func (t Thresholds) Get(name string) float64 {
	if value, ok := t[name]; ok {
		return value
	}
	return DefaultThresholds()[name]
}

// Merge returns a copy of t with every entry of override applied.
func (t Thresholds) Merge(override map[string]float64) Thresholds {
	merged := Thresholds{}
	for key, value := range DefaultThresholds() {
		merged[key] = value
	}
	for key, value := range t {
		merged[key] = value
	}
	for key, value := range override {
		merged[key] = value
	}
	return merged
}

// BandFor bands a named gate quantity: value ≥ apply → apply, ≥ review →
// review, else ignore.
func BandFor(value, apply, review float64) Band {
	switch {
	case value >= apply:
		return BandApply
	case value >= review:
		return BandReview
	default:
		return BandIgnore
	}
}
