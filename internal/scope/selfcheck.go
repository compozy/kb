package scope

import (
	"context"
	"fmt"
	"strconv"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/session"
)

// contractCheckBank is the self-check bank (embedded; MustLoad cannot fail).
var contractCheckBank = questions.MustLoad("contract_check")

// ConflictThreshold names the self-check threshold (default 0.7).
const ConflictThreshold = "contract_conflict"

// Pair is one (kept line, out_of_scope line) pair of the self-check.
type Pair struct {
	// N numbers the pair from 1; the question id is `conflict_<N>`.
	N int `json:"n"`
	// KeptField is "core" or "adjacent".
	KeptField string `json:"kept_field"`
	Kept      string `json:"kept"`
	Excluded  string `json:"excluded"`
}

// Conflict is a pair whose P(conflict) reached the threshold.
type Conflict struct {
	Pair
	P float64 `json:"p"`
}

// Conflicts is the self-check outcome.
type Conflicts struct {
	Pairs     int        `json:"pairs"`
	Threshold float64    `json:"threshold"`
	Conflicts []Conflict `json:"conflicts"`
	// Undecided lists pairs without a usable answer ("conflict_<n>:<reason>").
	Undecided []string `json:"undecided,omitempty"`
}

// Blocking reports whether the self-check blocks acceptance.
func (c Conflicts) Blocking() bool { return len(c.Conflicts) > 0 }

// BuildPairs returns one pair per (core or adjacent line, out_of_scope line)
// of the contract, numbered from 1 in contract order (core first).
func BuildPairs(c *contract.Contract) []Pair {
	n := c.Normalized()
	pairs := make([]Pair, 0, (len(n.Core)+len(n.Adjacent))*len(n.OutOfScope))
	for _, group := range []struct {
		field string
		lines []string
	}{{"core", n.Core}, {"adjacent", n.Adjacent}} {
		for _, kept := range group.lines {
			for _, excluded := range n.OutOfScope {
				pairs = append(pairs, Pair{N: len(pairs) + 1, KeptField: group.field, Kept: kept, Excluded: excluded})
			}
		}
	}
	return pairs
}

// SelfCheck asks, in one request (the engine splits it over several when
// there are more than 48 pairs), whether a document that fits each kept line
// also fits each out_of_scope line: bank `contract_check`, questions
// `conflict_<n>`, state `{pairs: [{n, kept, excluded}]}`. Pairs at P ≥
// contract_conflict (0.7) are conflicts. The draft's hash keys the receipts.
func SelfCheck(ctx context.Context, s *session.Session, draft *contract.Contract) (Conflicts, error) {
	pairs := BuildPairs(draft)
	result := Conflicts{Pairs: len(pairs), Threshold: s.Threshold(ConflictThreshold), Conflicts: []Conflict{}}
	if len(pairs) == 0 {
		return result, nil
	}

	statePairs := make([]map[string]any, 0, len(pairs))
	qs := make([]questions.Q, 0, len(pairs))
	for _, pair := range pairs {
		n := strconv.Itoa(pair.N)
		statePairs = append(statePairs, map[string]any{"n": pair.N, "kept": pair.Kept, "excluded": pair.Excluded})
		q, err := contractCheckBank.Question("conflict_{n}", map[string]string{"n": n})
		if err != nil {
			return result, fmt.Errorf("scope: self-check: %w", err)
		}
		qs = append(qs, q)
	}
	ref := s.Ref
	ref.Contract = draft.Hash()
	answers, err := s.Engine.Decide(ctx, decisions.Request{
		Topic:     ref,
		Purpose:   decisions.PurposeContractCheck,
		Subject:   "contract:" + ref.Contract,
		Bank:      contractCheckBank,
		State:     map[string]any{"pairs": statePairs},
		Questions: qs,
	})
	if err != nil {
		return result, fmt.Errorf("scope: self-check: %w", err)
	}
	for _, pair := range pairs {
		id := "conflict_" + strconv.Itoa(pair.N)
		answer := answers.Answers[id]
		p, ok := answer.P("")
		if !ok {
			reason := answer.Reason
			if reason == "" {
				reason = string(answer.Status)
			}
			if reason == "" {
				reason = decisions.ReasonMissingAnswer
			}
			result.Undecided = append(result.Undecided, id+":"+reason)
			continue
		}
		if p >= result.Threshold {
			result.Conflicts = append(result.Conflicts, Conflict{Pair: pair, P: p})
		}
	}
	return result, nil
}
