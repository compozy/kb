package okf

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/questions"
)

// AdvisoryOptions configures the advisory findings of `kb okf check` (spec
// §11): type_mismatch and description_unsupported. They are always
// SeverityInfo, so they never fail the command, strict or not, until a
// calibration exists for them.
//
// Answers are read from the bundle's receipts (`<bundle>/.decisions/
// receipts.jsonl`), keyed by a subject that carries a hash of the question
// state, so an edited concept never reuses a stale answer. When Decider is
// set (`kb okf check --decide`), a concept without a stored answer is asked
// once and the receipt it writes serves every later check; without Decider
// the check never calls a model.
type AdvisoryOptions struct {
	// Options is the type vocabulary with descriptions; empty skips
	// type_mismatch.
	Options []TypeOption
	// Threshold is the okf_type threshold (0 → DefaultTypeThreshold).
	Threshold float64
	// Decider computes missing answers; nil reads receipts only.
	Decider Decider
	// Topic is the bundle's topic reference for Decide (Root = bundle).
	Topic decisions.TopicRef
}

type conceptFile struct {
	absolutePath string
	relativePath string
}

type storedAnswers map[string]map[string]json.RawMessage

func advisoryIssues(ctx context.Context, bundlePath string, concepts []conceptFile, options AdvisoryOptions) ([]models.LintIssue, error) {
	threshold := options.Threshold
	if threshold <= 0 {
		threshold = DefaultTypeThreshold
	}
	stored, err := loadStoredAnswers(bundlePath)
	if err != nil {
		return nil, err
	}
	bank, err := questions.Load("okf_type")
	if err != nil {
		return nil, err
	}
	if options.Decider != nil && options.Topic.Root == "" {
		options.Topic.Root = bundlePath
	}

	issues := make([]models.LintIssue, 0)
	for _, concept := range concepts {
		content, err := os.ReadFile(concept.absolutePath)
		if err != nil {
			continue // reported by the conformance check
		}
		values, body, err := frontmatter.Parse(string(content))
		if err != nil || len(values) == 0 {
			continue
		}
		conceptType := strings.TrimSpace(frontmatter.GetString(values, "type"))
		title := strings.TrimSpace(frontmatter.GetString(values, "title"))
		description := strings.TrimSpace(frontmatter.GetString(values, "description"))

		if len(options.Options) > 0 && conceptType != "" {
			issue, err := typeMismatch(ctx, bank, stored, options, concept.relativePath, conceptType, title, body, threshold)
			if err != nil {
				return nil, err
			}
			if issue != nil {
				issues = append(issues, *issue)
			}
		}
		if description != "" {
			issue, err := descriptionUnsupported(ctx, bank, stored, options, concept.relativePath, description, body, threshold)
			if err != nil {
				return nil, err
			}
			if issue != nil {
				issues = append(issues, *issue)
			}
		}
	}
	return issues, nil
}

func typeMismatch(ctx context.Context, bank *questions.Bank, stored storedAnswers, options AdvisoryOptions, relativePath, conceptType, title, body string, threshold float64) (*models.LintIssue, error) {
	_, question, err := typeQuestion(options.Options)
	if err != nil {
		return nil, err
	}
	state := typeState(title, body)
	subject, err := advisorySubject(relativePath, bank, state, question)
	if err != nil {
		return nil, err
	}
	probs, found, err := answerProbabilities(ctx, stored, options, subject, bank, state, question)
	if err != nil || !found {
		return nil, err
	}
	suggestion := withRanking(TypeSuggestion{Status: decisions.StatusDecided}, probs)
	if !suggestion.Decided() || suggestion.Type == conceptType || suggestion.Probability < threshold {
		return nil, nil
	}
	message := fmt.Sprintf("the decision model suggests type %q (P=%.2f) for this %q concept (advisory)", suggestion.Type, suggestion.Probability, conceptType)
	if suggestion.Type == NoneType {
		message = fmt.Sprintf("no configured type fits this %q concept (P(none)=%.2f, advisory)", conceptType, suggestion.Probability)
	}
	issue := models.LintIssue{
		Kind:     models.LintIssueKindTypeMismatch,
		Severity: models.SeverityInfo,
		FilePath: relativePath,
		Target:   "type",
		Message:  message,
	}
	return &issue, nil
}

func descriptionUnsupported(ctx context.Context, bank *questions.Bank, stored storedAnswers, options AdvisoryOptions, relativePath, description, body string, threshold float64) (*models.LintIssue, error) {
	question, err := bank.Question(questionDescriptionUnsupported, nil)
	if err != nil {
		return nil, err
	}
	state := descriptionState(description, body)
	subject, err := advisorySubject(relativePath, bank, state, question)
	if err != nil {
		return nil, err
	}
	probs, found, err := answerProbabilities(ctx, stored, options, subject, bank, state, question)
	if err != nil || !found {
		return nil, err
	}
	p := probs[""]
	if p < threshold {
		return nil, nil
	}
	issue := models.LintIssue{
		Kind:     models.LintIssueKindDescriptionUnsupported,
		Severity: models.SeverityInfo,
		FilePath: relativePath,
		Target:   "description",
		Message:  fmt.Sprintf("the description states something the body does not support (P=%.2f, advisory)", p),
	}
	return &issue, nil
}

// answerProbabilities returns the stored answer of question for subject, or
// asks it when a Decider is configured. A noul answer is returned as
// {"": P(yes)}; a choice as its option distribution. found is false when no
// usable answer exists.
func answerProbabilities(ctx context.Context, stored storedAnswers, options AdvisoryOptions, subject string, bank *questions.Bank, state any, question questions.Q) (map[string]float64, bool, error) {
	if raw, ok := stored[subject][question.ID]; ok {
		if answer, ok := decodeRawAnswer(raw); ok {
			return rawProbabilities(answer)
		}
	}
	if options.Decider == nil {
		return nil, false, nil
	}
	result, err := options.Decider.Decide(ctx, decisions.Request{
		Topic:     options.Topic,
		Purpose:   decisions.PurposeOKFType,
		Subject:   subject,
		Bank:      bank,
		State:     state,
		Questions: []questions.Q{question},
	})
	if err != nil {
		return nil, false, err
	}
	answer := result.Answers[question.ID]
	if !answer.Decided() {
		return nil, false, nil
	}
	if answer.Noul != nil {
		return map[string]float64{"": *answer.Noul}, true, nil
	}
	return answer.Probs, len(answer.Probs) > 0, nil
}

func rawProbabilities(answer rawAnswer) (map[string]float64, bool, error) {
	if answer.Noul != nil {
		return map[string]float64{"": *answer.Noul}, true, nil
	}
	return answer.Probabilities, len(answer.Probabilities) > 0, nil
}

// advisorySubject names a concept question in receipts: the concept path
// plus a short hash of the bank, the state and the question, so an answer is
// only reused for the exact content and options it judged.
func advisorySubject(relativePath string, bank *questions.Bank, state any, question questions.Q) (string, error) {
	hash, err := decisions.HashJSON(map[string]any{
		"bank":      bank.Hash,
		"state":     state,
		"question":  question.ID,
		"criteria":  question.Criteria,
		"questions": question.Instructions,
	})
	if err != nil {
		return "", fmt.Errorf("hash advisory state: %w", err)
	}
	return relativePath + "#" + hash[:12], nil
}

// loadStoredAnswers indexes the bundle's decided okf_type receipts by
// subject and question id (the last row wins).
func loadStoredAnswers(bundlePath string) (storedAnswers, error) {
	rows, err := decisions.LoadReceipts(bundlePath)
	if err != nil {
		return nil, err
	}
	stored := make(storedAnswers)
	for _, row := range rows {
		if row.Status != decisions.StatusDecided || row.Purpose != string(decisions.PurposeOKFType) {
			continue
		}
		if stored[row.Subject] == nil {
			stored[row.Subject] = make(map[string]json.RawMessage)
		}
		maps.Copy(stored[row.Subject], row.Answers)
	}
	return stored, nil
}
