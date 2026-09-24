package decisions

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"

	"github.com/compozy/kb/internal/questions"
)

// Receipt validation tolerances. Probabilities and scores arrive rounded to
// two decimals, so distributions may drift from 1 and the reported score
// from its expectation; beyond these tolerances the answer is rejected, never
// coerced.
const (
	sumTolerance    = 0.02
	argmaxTolerance = 0.01
	scoreTolerance  = 0.02
	halfRounding    = 0.005
	valueEpsilon    = 1e-9
)

// preparedQuestion is a request question with its option set resolved once.
type preparedQuestion struct {
	q       questions.Q
	options []string
	levels  int
}

func prepareQuestions(list []questions.Q) ([]preparedQuestion, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("%w: no questions", ErrInvalidRequest)
	}
	prepared := make([]preparedQuestion, 0, len(list))
	seen := make(map[string]bool, len(list))
	for _, q := range list {
		if q.ID == "" {
			return nil, fmt.Errorf("%w: question without id", ErrInvalidRequest)
		}
		if seen[q.ID] {
			return nil, fmt.Errorf("%w: duplicate question id %q", ErrInvalidRequest, q.ID)
		}
		seen[q.ID] = true
		item := preparedQuestion{q: q}
		switch q.Type {
		case questions.TypeNoul:
		case questions.TypeChoice:
			options, err := choiceOptions(q.Criteria)
			if err != nil {
				return nil, fmt.Errorf("%w: question %q: %v", ErrInvalidRequest, q.ID, err)
			}
			item.options = options
		case questions.TypeScore:
			levels, err := scoreLevels(q.Criteria)
			if err != nil {
				return nil, fmt.Errorf("%w: question %q: %v", ErrInvalidRequest, q.ID, err)
			}
			item.levels = levels
		default:
			return nil, fmt.Errorf("%w: question %q has type %q", ErrInvalidRequest, q.ID, q.Type)
		}
		prepared = append(prepared, item)
	}
	slices.SortFunc(prepared, func(a, b preparedQuestion) int { return cmp.Compare(a.q.ID, b.q.ID) })
	return prepared, nil
}

func choiceOptions(criteria any) ([]string, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return nil, fmt.Errorf("encode criteria: %w", err)
	}
	var options map[string]json.RawMessage
	if err := json.Unmarshal(data, &options); err != nil || len(options) < 2 || len(options) > 255 {
		return nil, errors.New("choice criteria must map 2 to 255 options")
	}
	names := make([]string, 0, len(options))
	for name := range options {
		names = append(names, name)
	}
	slices.Sort(names)
	return names, nil
}

func scoreLevels(criteria any) (int, error) {
	data, err := json.Marshal(criteria)
	if err != nil {
		return 0, fmt.Errorf("encode criteria: %w", err)
	}
	var levels []json.RawMessage
	if err := json.Unmarshal(data, &levels); err != nil || len(levels) < 2 || len(levels) > 10 {
		return 0, errors.New("score criteria must list 2 to 10 levels")
	}
	return len(levels), nil
}

// rawAnswer is the decoded answer as returned by the decisions endpoint.
type rawAnswer struct {
	Type          string             `json:"type"`
	Noul          *float64           `json:"noul"`
	Choice        *string            `json:"choice"`
	Score         *float64           `json:"score"`
	Probabilities map[string]float64 `json:"probabilities"`
	Confidence    *float64           `json:"confidence"`
}

// validateAnswer checks one raw answer against its question and returns the
// decided (unbanded) answer, or an error naming the first violated rule.
func validateAnswer(pq preparedQuestion, raw json.RawMessage) (Answer, error) {
	var answer rawAnswer
	if err := json.Unmarshal(raw, &answer); err != nil {
		return Answer{}, fmt.Errorf("decode answer: %w", err)
	}
	if answer.Type != pq.q.Type {
		return Answer{}, fmt.Errorf("answer type %q does not match question type %q", answer.Type, pq.q.Type)
	}
	switch pq.q.Type {
	case questions.TypeNoul:
		if answer.Noul == nil || !unit(*answer.Noul) {
			return Answer{}, errors.New("noul missing or outside [0,1]")
		}
		return Answer{Type: pq.q.Type, Status: StatusDecided, Noul: ptr(*answer.Noul)}, nil
	case questions.TypeChoice:
		return validateChoice(pq, answer)
	default:
		return validateScore(pq, answer)
	}
}

func validateDistribution(probabilities map[string]float64, keys []string) error {
	if len(probabilities) != len(keys) {
		return fmt.Errorf("probabilities cover %d keys, want %d", len(probabilities), len(keys))
	}
	sum := 0.0
	for _, key := range keys {
		p, ok := probabilities[key]
		if !ok {
			return fmt.Errorf("probabilities miss %q", key)
		}
		if !unit(p) {
			return fmt.Errorf("probability of %q outside [0,1]", key)
		}
		sum += p
	}
	if math.Abs(sum-1) > sumTolerance+valueEpsilon {
		return fmt.Errorf("probabilities sum to %.4f", sum)
	}
	return nil
}

func validateConfidence(confidence *float64) error {
	if confidence == nil || !unit(*confidence) {
		return errors.New("confidence missing or outside [0,1]")
	}
	return nil
}

func validateChoice(pq preparedQuestion, answer rawAnswer) (Answer, error) {
	if answer.Choice == nil || !slices.Contains(pq.options, *answer.Choice) {
		return Answer{}, errors.New("choice missing or not among the options")
	}
	if err := validateDistribution(answer.Probabilities, pq.options); err != nil {
		return Answer{}, err
	}
	best := 0.0
	for _, p := range answer.Probabilities {
		best = math.Max(best, p)
	}
	if answer.Probabilities[*answer.Choice] < best-argmaxTolerance-valueEpsilon {
		return Answer{}, fmt.Errorf("choice %q is not the argmax", *answer.Choice)
	}
	if err := validateConfidence(answer.Confidence); err != nil {
		return Answer{}, err
	}
	return Answer{
		Type:       pq.q.Type,
		Status:     StatusDecided,
		Choice:     *answer.Choice,
		Probs:      copyProbs(answer.Probabilities),
		Confidence: ptr(*answer.Confidence),
	}, nil
}

func validateScore(pq preparedQuestion, answer rawAnswer) (Answer, error) {
	keys := make([]string, pq.levels)
	for level := range pq.levels {
		keys[level] = strconv.Itoa(level)
	}
	if err := validateDistribution(answer.Probabilities, keys); err != nil {
		return Answer{}, err
	}
	if answer.Score == nil || math.IsNaN(*answer.Score) || *answer.Score < -valueEpsilon || *answer.Score > float64(pq.levels-1)+valueEpsilon {
		return Answer{}, fmt.Errorf("score missing or outside [0,%d]", pq.levels-1)
	}
	expected := 0.0
	for level, key := range keys {
		expected += float64(level) * answer.Probabilities[key]
	}
	// Each rounded probability may be off by half a unit, weighted by its
	// level, and the score itself is rounded too.
	sumLevels := float64(pq.levels*(pq.levels-1)) / 2
	tolerance := scoreTolerance + sumLevels*halfRounding + halfRounding
	if math.Abs(*answer.Score-expected) > tolerance {
		return Answer{}, fmt.Errorf("score %.4f is not the expectation %.4f", *answer.Score, expected)
	}
	if err := validateConfidence(answer.Confidence); err != nil {
		return Answer{}, err
	}
	return Answer{
		Type:       pq.q.Type,
		Status:     StatusDecided,
		Score:      ptr(*answer.Score),
		Probs:      copyProbs(answer.Probabilities),
		Confidence: ptr(*answer.Confidence),
	}, nil
}

// validateAnswers validates a whole answer set. Missing ids make the request
// incomplete: missing answers become undecided:missing_answer and every other
// answer undecided:invalid_receipt. Extra ids make the whole set invalid. A
// complete set is validated per answer. ok reports whether every answer is
// decided.
func validateAnswers(prepared []preparedQuestion, answers map[string]json.RawMessage) (map[string]Answer, bool) {
	out := make(map[string]Answer, len(prepared))
	complete := len(answers) == len(prepared)
	for _, pq := range prepared {
		raw, ok := answers[pq.q.ID]
		if !ok || isJSONNull(raw) {
			complete = false
		}
	}
	if !complete {
		for _, pq := range prepared {
			reason := ReasonInvalidReceipt
			if raw, ok := answers[pq.q.ID]; !ok || isJSONNull(raw) {
				reason = ReasonMissingAnswer
			}
			out[pq.q.ID] = undecided(pq.q.Type, reason)
		}
		return out, false
	}
	allDecided := true
	for _, pq := range prepared {
		answer, err := validateAnswer(pq, answers[pq.q.ID])
		if err != nil {
			answer = undecided(pq.q.Type, ReasonInvalidReceipt)
			allDecided = false
		}
		out[pq.q.ID] = answer
	}
	return out, allDecided
}

func isJSONNull(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null"))
}

func undecided(questionType, reason string) Answer {
	return Answer{Type: questionType, Status: StatusUndecided, Reason: reason}
}

func notChecked(questionType, reason string) Answer {
	return Answer{Type: questionType, Status: StatusNotChecked, Reason: reason}
}

func unit(value float64) bool {
	return !math.IsNaN(value) && value >= 0 && value <= 1
}

func ptr(value float64) *float64 {
	return &value
}

func copyProbs(probabilities map[string]float64) map[string]float64 {
	return maps.Clone(probabilities)
}
