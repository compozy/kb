package okf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/questions"
)

// NoneType is the exit option of the type choice: no configured type fits.
const NoneType = "none"

// DefaultTypeThreshold is the P(argmax) a type suggestion needs before kb
// uses it or reports a disagreement (spec §12.3 default `okf_type`).
const DefaultTypeThreshold = 0.8

// reviewTypeThreshold is the lower edge of the okf_type review band.
const reviewTypeThreshold = 0.5

// Question ids of the okf_type bank.
const (
	questionOKFType                = "okf_type"
	questionDescriptionUnsupported = "description_unsupported"
)

// excerptTokens bounds the document text sent with a type or description
// question; the excerpt is section-aware (corpus.Excerpt).
const excerptTokens = 3000

// topTypes is how many ranked types a suggestion reports.
const topTypes = 3

// ErrTypeRequired is returned by Promote when no type can be chosen: no
// --type and either no vocabulary or no suggestion in the apply band.
var ErrTypeRequired = errors.New("--type is required")

// Decider asks the decision model; *decisions.Engine satisfies it.
type Decider interface {
	Decide(ctx context.Context, r decisions.Request) (decisions.Result, error)
}

// TypeOption is one OKF type offered to the decision model: the type name is
// the option and Description its option text.
type TypeOption struct {
	Name        string
	Description string
}

// TypeScore is one ranked type of a suggestion.
type TypeScore struct {
	Type        string  `json:"type"`
	Probability float64 `json:"probability"`
}

// TypeSuggestion is the decision model's answer to "which configured type
// describes this document". Type is the argmax option (NoneType when no
// type fits); Status is not decided when the model gave no usable answer.
type TypeSuggestion struct {
	Status      decisions.Status `json:"status"`
	Reason      string           `json:"reason,omitempty"`
	Type        string           `json:"type,omitempty"`
	Probability float64          `json:"probability"`
	Top         []TypeScore      `json:"top,omitempty"`
	ReceiptKey  string           `json:"receipt,omitempty"`
}

// Decided reports whether the suggestion carries a usable answer.
func (s TypeSuggestion) Decided() bool { return s.Status == decisions.StatusDecided }

// TopText renders the ranked types as "a 0.55, b 0.30, none 0.10".
func (s TypeSuggestion) TopText() string {
	parts := make([]string, 0, len(s.Top))
	for _, score := range s.Top {
		parts = append(parts, fmt.Sprintf("%s %.2f", score.Type, score.Probability))
	}
	return strings.Join(parts, ", ")
}

// TypeDocument is the document a type is suggested for. Subject names it in
// receipts and review items (topic-relative path).
type TypeDocument struct {
	Subject string
	Title   string
	Body    string
}

// TypeSuggester suggests the OKF type of a document.
type TypeSuggester interface {
	SuggestType(ctx context.Context, doc TypeDocument) (TypeSuggestion, error)
}

// EngineTypeSuggester asks the okf_type choice over Options through a
// decision engine for one topic (receipts and cache live under Topic.Root).
type EngineTypeSuggester struct {
	Decider Decider
	Topic   decisions.TopicRef
	Options []TypeOption
	// Extras are the topic's extra okf_type question banks (spec §4.3):
	// their questions ride along in the same request and their answers are
	// only recorded in receipts.
	Extras []*questions.Bank
}

// SuggestType asks one okf_type choice with state {document: {title,
// excerpt}}.
func (s EngineTypeSuggester) SuggestType(ctx context.Context, doc TypeDocument) (TypeSuggestion, error) {
	if s.Decider == nil {
		return TypeSuggestion{}, errors.New("okf: type suggester has no decision engine")
	}
	bank, question, err := typeQuestion(s.Options)
	if err != nil {
		return TypeSuggestion{}, err
	}
	bank, qs, err := withExtras(bank, []questions.Q{question}, s.Extras)
	if err != nil {
		return TypeSuggestion{}, err
	}
	result, err := s.Decider.Decide(ctx, decisions.Request{
		Topic:     s.Topic,
		Purpose:   decisions.PurposeOKFType,
		Subject:   doc.Subject,
		Bank:      bank,
		State:     typeState(doc.Title, doc.Body),
		Questions: qs,
	})
	if err != nil {
		return TypeSuggestion{}, fmt.Errorf("okf: suggest type: %w", err)
	}
	return suggestionFromAnswer(result.Answers[questionOKFType]), nil
}

// TypeOptions builds the offered types from the vocabulary and a description
// lookup (nil or "" falls back to the type name).
func TypeOptions(types []string, describe func(string) string) []TypeOption {
	options := make([]TypeOption, 0, len(types))
	seen := make(map[string]struct{}, len(types))
	for _, name := range types {
		name = strings.TrimSpace(name)
		if name == "" || name == NoneType {
			continue
		}
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		description := name
		if describe != nil {
			if text := strings.TrimSpace(describe(name)); text != "" {
				description = text
			}
		}
		options = append(options, TypeOption{Name: name, Description: description})
	}
	return options
}

func typeQuestion(options []TypeOption) (*questions.Bank, questions.Q, error) {
	if len(options) == 0 {
		return nil, questions.Q{}, errors.New("okf: no OKF types configured to choose from")
	}
	bank, err := questions.Load("okf_type")
	if err != nil {
		return nil, questions.Q{}, err
	}
	question, err := bank.Question(questionOKFType, nil)
	if err != nil {
		return nil, questions.Q{}, err
	}
	return bank, bank.WithCriteria(question, typeCriteria(options)), nil
}

func typeCriteria(options []TypeOption) map[string]any {
	criteria := make(map[string]any, len(options))
	for _, option := range options {
		criteria[option.Name] = option.Description
	}
	return criteria
}

func typeState(title, body string) map[string]any {
	return map[string]any{"document": map[string]any{
		"title":   strings.TrimSpace(title),
		"excerpt": corpus.Excerpt(body, excerptTokens, nil),
	}}
}

func descriptionState(description, body string) map[string]any {
	return map[string]any{"concept": map[string]any{
		"description": strings.TrimSpace(description),
		"body":        corpus.Excerpt(body, excerptTokens, nil),
	}}
}

func suggestionFromAnswer(answer decisions.Answer) TypeSuggestion {
	suggestion := TypeSuggestion{Status: answer.Status, Reason: answer.Reason, ReceiptKey: answer.Receipt}
	if suggestion.Status == "" {
		suggestion.Status = decisions.StatusUndecided
		suggestion.Reason = decisions.ReasonMissingAnswer
	}
	if !answer.Decided() {
		return suggestion
	}
	return withRanking(suggestion, answer.Probs)
}

func withRanking(suggestion TypeSuggestion, probs map[string]float64) TypeSuggestion {
	ranked := make([]TypeScore, 0, len(probs))
	for _, name := range slices.Sorted(maps.Keys(probs)) {
		ranked = append(ranked, TypeScore{Type: name, Probability: probs[name]})
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].Probability > ranked[j].Probability })
	if len(ranked) == 0 {
		suggestion.Status = decisions.StatusUndecided
		suggestion.Reason = decisions.ReasonMissingAnswer
		return suggestion
	}
	suggestion.Type = ranked[0].Type
	suggestion.Probability = ranked[0].Probability
	suggestion.Top = ranked[:min(len(ranked), topTypes)]
	return suggestion
}

// rawAnswer is the stored shape of one answer in a receipt.
type rawAnswer struct {
	Type          string             `json:"type"`
	Noul          *float64           `json:"noul"`
	Probabilities map[string]float64 `json:"probabilities"`
}

func decodeRawAnswer(raw json.RawMessage) (rawAnswer, bool) {
	var answer rawAnswer
	if err := json.Unmarshal(raw, &answer); err != nil {
		return rawAnswer{}, false
	}
	return answer, true
}

// withExtras appends the plain questions of the topic's extra okf_type banks
// to qs and composes their bank after bank (unchanged without extras, so the
// cache key only moves when a topic adds questions).
func withExtras(bank *questions.Bank, qs []questions.Q, extras []*questions.Bank) (*questions.Bank, []questions.Q, error) {
	extraBank, extra, err := questions.Instantiate(extras, nil)
	if err != nil {
		return bank, qs, fmt.Errorf("okf: topic extra okf_type questions: %w", err)
	}
	if extraBank == nil {
		return bank, qs, nil
	}
	return questions.Compose(bank, extraBank), append(qs, extra...), nil
}
