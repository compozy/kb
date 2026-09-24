package scope

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/session"
)

// Acceptance errors.
var (
	// ErrNoDraft is returned by Accept when topic.yaml has no contract_draft.
	ErrNoDraft = errors.New("scope: topic.yaml has no contract_draft")
	// ErrConflicts is returned by Accept when the self-check found conflicts
	// and Force is off.
	ErrConflicts = errors.New("scope: the self-check found conflicting contract lines")
	// ErrNotConfirmed is returned by Accept when the activation was not
	// confirmed by typing the topic slug (and Yes is off).
	ErrNotConfirmed = errors.New("scope: contract not activated: confirmation did not match the topic slug")
	// ErrDraftChanged is returned by Accept when contract_draft changed while
	// the preview ran.
	ErrDraftChanged = errors.New("scope: contract_draft changed during the preview; re-run --accept")
)

// AcceptOptions configures Accept.
type AcceptOptions struct {
	// Force accepts despite self-check conflicts.
	Force bool
	// Yes activates without the interactive confirmation (scripts).
	Yes bool
	// Confirm asks the user to type the topic slug; it returns the typed
	// text. Nil without Yes means the draft is never activated.
	Confirm func(prompt string) (string, error)
	// Preview configures the impact preview.
	Preview PreviewOptions
	// OnSelfCheck and OnPreview, when set, receive each step's outcome as
	// soon as it is known (before the confirmation prompt).
	OnSelfCheck func(Conflicts)
	OnPreview   func(PreviewReport)
}

// AcceptReport is the outcome of Accept.
type AcceptReport struct {
	Draft     *contract.Contract `json:"draft"`
	SelfCheck Conflicts          `json:"self_check"`
	// Forced reports conflicts accepted through Force.
	Forced bool `json:"forced,omitempty"`
	// Preview is nil when the self-check blocked.
	Preview   *PreviewReport `json:"preview,omitempty"`
	Activated bool           `json:"activated"`
}

// Accept runs the acceptance of spec §5.1: the self-check (conflicts block
// unless Force), then the impact preview of the draft, then activation only
// after the topic slug is typed through Confirm (or Yes). Activation moves
// contract_draft to contract in topic.yaml and re-renders the CLAUDE.md
// `## Selection contract` section. The report is returned alongside every
// error with what was done.
func Accept(ctx context.Context, s *session.Session, opts AcceptOptions) (AcceptReport, error) {
	report := AcceptReport{}
	settings, err := contract.LoadSettings(s.Root())
	if err != nil {
		return report, fmt.Errorf("scope: accept: %w", err)
	}
	if settings.ContractDraft.Empty() {
		return report, fmt.Errorf("%w: run `kb topic contract %s --draft` or `kb topic contract %s --import-claude` first", ErrNoDraft, s.Topic.Slug, s.Topic.Slug)
	}
	draft := settings.ContractDraft.Normalized()
	report.Draft = &draft
	if err := draft.Validate(); err != nil {
		return report, fmt.Errorf("scope: accept: the draft is invalid; edit contract_draft in topic.yaml: %w", err)
	}

	checked, err := SelfCheck(ctx, s, &draft)
	report.SelfCheck = checked
	if err != nil {
		return report, err
	}
	if opts.OnSelfCheck != nil {
		opts.OnSelfCheck(checked)
	}
	if checked.Blocking() {
		if !opts.Force {
			return report, fmt.Errorf("%w (%d at P ≥ %.2f); edit contract_draft or pass --force", ErrConflicts, len(checked.Conflicts), checked.Threshold)
		}
		report.Forced = true
	}

	preview, err := Preview(ctx, s, &draft, opts.Preview)
	report.Preview = &preview
	if err != nil {
		return report, err
	}
	if opts.OnPreview != nil {
		opts.OnPreview(preview)
	}

	confirmed, err := confirmed(s.Topic.Slug, opts)
	if err != nil {
		return report, err
	}
	if !confirmed {
		return report, ErrNotConfirmed
	}
	if err := activate(s.Root(), draft.Hash()); err != nil {
		return report, err
	}
	report.Activated = true
	return report, nil
}

// confirmed reports whether activation was confirmed: Yes, or Confirm
// returning exactly the topic slug (surrounding spaces ignored).
func confirmed(slug string, opts AcceptOptions) (bool, error) {
	if opts.Yes {
		return true, nil
	}
	if opts.Confirm == nil {
		return false, nil
	}
	typed, err := opts.Confirm(fmt.Sprintf("Type the topic slug %q to activate this contract: ", slug))
	if err != nil {
		return false, fmt.Errorf("scope: accept: read confirmation: %w", err)
	}
	return strings.TrimSpace(typed) == slug, nil
}

// activate moves contract_draft to contract after checking it is still the
// previewed draft, then re-renders the CLAUDE.md contract section (when the
// topic has a CLAUDE.md).
func activate(topicRoot, previewedHash string) error {
	settings, err := contract.LoadSettings(topicRoot)
	if err != nil {
		return fmt.Errorf("scope: accept: %w", err)
	}
	if settings.ContractDraft.Hash() != previewedHash {
		return ErrDraftChanged
	}
	if _, err := contract.ActivateDraft(topicRoot); err != nil {
		return fmt.Errorf("scope: accept: %w", err)
	}
	if _, err := os.Stat(filepath.Join(topicRoot, claudeFile)); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err := contract.SyncClaude(topicRoot); err != nil {
		return fmt.Errorf("scope: accept: %w", err)
	}
	return nil
}
