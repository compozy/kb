package okf

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/review"
)

// SourceTopic returns the slug (vault-relative path) of the topic that holds
// the wiki document to promote: the type suggestion runs in that topic's
// decision session and its review queue.
func SourceTopic(vaultPath, sourceDocPath string) (string, error) {
	_, _, slug, _, err := resolveSourceDocument(PromoteInput{SourceDocPath: sourceDocPath, VaultPath: vaultPath})
	if err != nil {
		return "", err
	}
	if slug == "" {
		return "", fmt.Errorf("source document %q is not inside a topic", sourceDocPath)
	}
	return slug, nil
}

type typeContext struct {
	document TypeDocument
	// topicRoot is the source wiki topic, where okf-type review items go.
	topicRoot string
	clock     func() time.Time
}

type typeChoice struct {
	conceptType string
	suggestion  *TypeSuggestion
	warning     string
}

// chooseType settles the concept type (spec §11):
//   - with --type and a suggester, the type is kept and a suggestion that
//     disagrees at P ≥ threshold becomes a warning;
//   - without --type, a suggestion at P ≥ threshold (and not `none`) is used;
//     a review-band suggestion is queued as an okf-type review item and, like
//     every other outcome, fails asking for --type with the top types.
func chooseType(ctx context.Context, input PromoteInput, tc typeContext) (typeChoice, error) {
	explicit := strings.TrimSpace(input.Type)
	if input.Suggester == nil {
		return typeChoice{conceptType: explicit}, nil
	}
	threshold := input.TypeThreshold
	if threshold <= 0 {
		threshold = DefaultTypeThreshold
	}

	suggestion, err := input.Suggester.SuggestType(ctx, tc.document)
	if err != nil {
		return typeChoice{}, err
	}
	choice := typeChoice{conceptType: explicit, suggestion: &suggestion}

	if explicit != "" {
		if suggestion.Decided() && suggestion.Type != explicit && suggestion.Probability >= threshold {
			choice.warning = fmt.Sprintf("the decision model suggests type %q (P=%.2f), not %q", suggestion.Type, suggestion.Probability, explicit)
		}
		return choice, nil
	}

	switch {
	case suggestion.Status == decisions.StatusNotChecked && suggestion.Reason == decisions.ReasonExcluded:
		return typeChoice{}, fmt.Errorf("%w: %s matches decisions.exclude, so no type is suggested; pass --type", ErrTypeRequired, tc.document.Subject)
	case !suggestion.Decided():
		return typeChoice{}, fmt.Errorf("%w: the type suggestion is %s (%s); pass --type", ErrTypeRequired, suggestion.Status, suggestion.Reason)
	case suggestion.Type == NoneType && suggestion.Probability >= threshold:
		return typeChoice{}, fmt.Errorf("%w: no configured OKF type fits this document (P(none)=%.2f); pass --type", ErrTypeRequired, suggestion.Probability)
	case suggestion.Type != NoneType && suggestion.Probability >= threshold:
		choice.conceptType = suggestion.Type
		return choice, nil
	}

	queued := ""
	if suggestion.Type != NoneType && suggestion.Probability >= reviewTypeThreshold && tc.topicRoot != "" {
		id, err := queueTypeReview(input, tc, suggestion)
		if err != nil {
			return typeChoice{}, err
		}
		queued = fmt.Sprintf("; queued for review as %s", id)
	}
	return typeChoice{}, fmt.Errorf(
		"%w: no type reaches P ≥ %.2f (top: %s)%s; pass --type <type>",
		ErrTypeRequired, threshold, suggestion.TopText(), queued,
	)
}

// queueTypeReview records a review-band suggestion in the source topic's
// okf-type queue.
func queueTypeReview(input PromoteInput, tc typeContext, suggestion TypeSuggestion) (string, error) {
	store := review.Open(tc.topicRoot, tc.clock)
	item := review.Item{
		Queue:       review.QueueOKFType,
		Purpose:     string(decisions.PurposeOKFType),
		Subject:     tc.document.Subject,
		Target:      input.TargetTopic.Slug,
		Question:    questionOKFType,
		Probability: suggestion.Probability,
		ReceiptKey:  suggestion.ReceiptKey,
		Evidence:    fmt.Sprintf("promote to %s: top types %s", input.TargetTopic.Slug, suggestion.TopText()),
		Action: map[string]any{
			"to":   input.TargetTopic.Slug,
			"type": suggestion.Type,
			"top":  suggestion.Top,
		},
	}
	item.ID = review.ItemID(item.Queue, item.Subject, item.Target, item.Question)
	if _, err := store.Add(item); err != nil {
		return "", fmt.Errorf("queue okf-type review: %w", err)
	}
	return item.ID, nil
}
