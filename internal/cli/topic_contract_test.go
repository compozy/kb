package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/session"
)

func TestTopicContractNeedsExactlyOneAction(t *testing.T) {
	tests := [][]string{
		{"topic", "contract", "demo"},
		{"topic", "contract", "demo", "--draft", "--accept"},
		{"topic", "contract", "demo", "--draft", "--import-claude", "--accept"},
		{"topic", "contract", "demo", "--import-claude", "--yes"},
		{"topic", "contract", "demo", "--draft", "--force"},
	}
	original := openSession
	t.Cleanup(func() { openSession = original })
	openSession = func(*cobra.Command, string, string, session.Flags) (*session.Session, error) {
		t.Fatal("an invalid flag set opened a session")
		return nil, nil
	}
	for _, args := range tests {
		if _, _, err := runLinkFindRoot(t, args...); err == nil {
			t.Errorf("%v: want an error", args)
		}
	}
}

func TestTopicContractImportClaudeNeedsNoDecisionModel(t *testing.T) {
	vault := newLinkFindVault(t, nil)
	root := filepath.Join(vault, "demo")
	claudePath := filepath.Join(root, "CLAUDE.md")
	claude, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}
	section := "## Selection contract\n\n- **Purpose:** Collect material on typed decision models.\n- **Central:** Typed decision models; LLM judges\n- **Adjacent (keep):** Evaluation methods, kept as comparison\n- **Out of scope:** Sports results with no decision-model content\n"
	// Replace the scaffolded placeholder section with a hand-written one.
	withSection := strings.Replace(string(claude), contract.RenderSection(nil), section, 1)
	if withSection == string(claude) {
		t.Fatal("placeholder section not found in the topic CLAUDE.md")
	}
	if err := os.WriteFile(claudePath, []byte(withSection), 0o644); err != nil {
		t.Fatal(err)
	}
	original := openSession
	t.Cleanup(func() { openSession = original })
	openSession = func(*cobra.Command, string, string, session.Flags) (*session.Session, error) {
		t.Fatal("--import-claude opened a decision session")
		return nil, nil
	}

	stdout, stderr, err := runLinkFindRoot(t, "topic", "contract", "demo", "--import-claude", "--vault", vault)
	if err != nil {
		t.Fatalf("--import-claude: %v\n%s", err, stderr)
	}
	for _, want := range []string{"contract_draft:", "purpose: Collect material on typed decision models.", "- LLM judges", "## Selection contract"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout lacks %q:\n%s", want, stdout)
		}
	}
	if !strings.Contains(stderr, "kb topic contract demo --accept") {
		t.Fatalf("stderr = %s", stderr)
	}
	settings, err := contract.LoadSettings(root)
	if err != nil {
		t.Fatal(err)
	}
	if settings.ContractDraft == nil || len(settings.ContractDraft.Core) != 2 || !settings.Contract.Empty() {
		t.Fatalf("settings = %+v", settings)
	}
}
