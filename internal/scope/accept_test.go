package scope

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/fakes"
)

func TestConfirmed(t *testing.T) {
	t.Parallel()
	answer := func(text string, err error) func(string) (string, error) {
		return func(prompt string) (string, error) {
			if !strings.Contains(prompt, `"demo"`) {
				t.Errorf("prompt %q does not name the slug", prompt)
			}
			return text, err
		}
	}
	readErr := errors.New("closed")
	tests := []struct {
		name    string
		opts    AcceptOptions
		want    bool
		wantErr error
	}{
		{name: "yes skips the prompt", opts: AcceptOptions{Yes: true, Confirm: answer("", readErr)}, want: true},
		{name: "no prompt means no activation", opts: AcceptOptions{}},
		{name: "exact slug", opts: AcceptOptions{Confirm: answer("demo", nil)}, want: true},
		{name: "slug with newline and spaces", opts: AcceptOptions{Confirm: answer("  demo \n", nil)}, want: true},
		{name: "wrong text", opts: AcceptOptions{Confirm: answer("yes\n", nil)}},
		{name: "case matters", opts: AcceptOptions{Confirm: answer("Demo", nil)}},
		{name: "empty input", opts: AcceptOptions{Confirm: answer("", nil)}},
		{name: "read error", opts: AcceptOptions{Confirm: answer("", readErr)}, wantErr: readErr},
	}
	for _, tt := range tests {
		got, err := confirmed("demo", tt.opts)
		if got != tt.want || !errors.Is(err, tt.wantErr) {
			t.Errorf("%s: confirmed = (%v, %v), want (%v, %v)", tt.name, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestAcceptWithoutDraftPointsAtDraftAndImport(t *testing.T) {
	t.Parallel()
	vault, _ := newTestTopic(t)
	fake := fakeServer(t, nil, nil)
	s := openTestSession(t, vault, fake)

	_, err := Accept(context.Background(), s, AcceptOptions{Yes: true})
	if !errors.Is(err, ErrNoDraft) || !strings.Contains(err.Error(), "--draft") || !strings.Contains(err.Error(), "--import-claude") {
		t.Fatalf("Accept error = %v", err)
	}
	if len(fake.Calls()) != 0 {
		t.Fatal("Accept without a draft called the model")
	}
}

func TestAcceptFlow(t *testing.T) {
	t.Parallel()
	conflicting := func(call fakes.Call, q fakes.Question) any {
		if strings.HasPrefix(q.ID, "conflict_") {
			return fakes.Noul(0.9)
		}
		return roleByTitle(call, q)
	}
	tests := []struct {
		name          string
		decide        fakes.DecideFunc
		opts          AcceptOptions
		wantErr       error
		wantPreview   bool
		wantActivated bool
		wantForced    bool
	}{
		{name: "conflicts block before the preview", decide: conflicting, opts: AcceptOptions{Yes: true}, wantErr: ErrConflicts},
		{name: "force previews and activates despite conflicts", decide: conflicting, opts: AcceptOptions{Force: true, Yes: true}, wantPreview: true, wantActivated: true, wantForced: true},
		{name: "wrong confirmation previews but does not activate", decide: roleByTitle, opts: AcceptOptions{Confirm: func(string) (string, error) { return "nope\n", nil }}, wantErr: ErrNotConfirmed, wantPreview: true},
		{name: "typed slug activates", decide: roleByTitle, opts: AcceptOptions{Confirm: func(string) (string, error) { return "demo\n", nil }}, wantPreview: true, wantActivated: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vault, root := newTestTopic(t)
			writeSource(t, root, "raw/articles/a.md", "Jev notes", "Typed decision models.\n")
			writeSource(t, root, "raw/articles/b.md", "Offtopic scores", "Last night.\n")
			if err := contract.SaveContractDraft(root, testDraft); err != nil {
				t.Fatal(err)
			}
			fake := fakeServer(t, tt.decide, nil)
			s := openTestSession(t, vault, fake)
			var sawSelfCheck, sawPreview bool
			tt.opts.OnSelfCheck = func(Conflicts) { sawSelfCheck = true }
			tt.opts.OnPreview = func(PreviewReport) { sawPreview = true }

			report, err := Accept(context.Background(), s, tt.opts)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Accept error = %v, want %v", err, tt.wantErr)
			}
			if !sawSelfCheck || report.SelfCheck.Pairs != 2 {
				t.Fatalf("self-check not reported: %+v", report.SelfCheck)
			}
			if (report.Preview != nil) != tt.wantPreview || sawPreview != tt.wantPreview {
				t.Fatalf("preview = %v (callback %v), want %v", report.Preview != nil, sawPreview, tt.wantPreview)
			}
			if report.Activated != tt.wantActivated || report.Forced != tt.wantForced {
				t.Fatalf("activated = %v forced = %v, want %v %v", report.Activated, report.Forced, tt.wantActivated, tt.wantForced)
			}

			settings, err := contract.LoadSettings(root)
			if err != nil {
				t.Fatal(err)
			}
			claude := readFile(t, filepath.Join(root, "CLAUDE.md"))
			if tt.wantActivated {
				if settings.Contract.Hash() != testDraft.Hash() || settings.ContractDraft != nil {
					t.Fatalf("after activation: contract %+v, draft %+v", settings.Contract, settings.ContractDraft)
				}
				if !strings.Contains(claude, testDraft.Purpose) {
					t.Fatalf("CLAUDE.md not re-rendered:\n%s", claude)
				}
				return
			}
			if !settings.Contract.Empty() || settings.ContractDraft.Hash() != testDraft.Hash() {
				t.Fatalf("without activation: contract %+v, draft %+v", settings.Contract, settings.ContractDraft)
			}
			if strings.Contains(claude, testDraft.Purpose) {
				t.Fatal("CLAUDE.md changed without activation")
			}
		})
	}
}

func TestAcceptRefusesADraftEditedDuringThePreview(t *testing.T) {
	t.Parallel()
	vault, root := newTestTopic(t)
	writeSource(t, root, "raw/articles/a.md", "Jev notes", "Typed decision models.\n")
	if err := contract.SaveContractDraft(root, testDraft); err != nil {
		t.Fatal(err)
	}
	fake := fakeServer(t, roleByTitle, nil)
	s := openTestSession(t, vault, fake)
	edited := testDraft.Normalized()
	edited.Purpose = "Something else entirely."

	_, err := Accept(context.Background(), s, AcceptOptions{Confirm: func(string) (string, error) {
		if err := contract.SaveContractDraft(root, &edited); err != nil {
			t.Fatal(err)
		}
		return "demo", nil
	}})
	if !errors.Is(err, ErrDraftChanged) {
		t.Fatalf("Accept error = %v, want ErrDraftChanged", err)
	}
	settings, _ := contract.LoadSettings(root)
	if !settings.Contract.Empty() {
		t.Fatal("an edited draft was activated")
	}
}
