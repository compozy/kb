package okf

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/review"
)

func TestPromoteWritesConceptIndexAndLog(t *testing.T) {
	t.Parallel()

	t.Run("Should write promoted concept index and log entry", func(t *testing.T) {
		vaultPath := t.TempDir()
		sourceTopic := filepath.Join(vaultPath, "research")
		targetTopic := filepath.Join(vaultPath, "catalog")
		mkdirAll(t, filepath.Join(sourceTopic, "wiki", "concepts"))
		mkdirAll(t, targetTopic)
		writeFile(t, filepath.Join(sourceTopic, "CLAUDE.md"), "# Research\n")
		writeFile(t, filepath.Join(targetTopic, "CLAUDE.md"), "# Catalog\n")
		writeFile(t, filepath.Join(targetTopic, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Old Index\n")
		writeFile(t, filepath.Join(targetTopic, "log.md"), "# Directory Update Log\n")
		sourcePath := filepath.Join(sourceTopic, "wiki", "concepts", "Alpha Note.md")
		sourceContent := strings.Join([]string{
			"---",
			"title: Alpha Note",
			"type: wiki",
			"stage: compiled",
			"tags: [systems, alpha]",
			"---",
			"Alpha note explains the operational flow. See [[research/wiki/concepts/Beta Note|Beta]] and [[research/wiki/concepts/Alpha Note#details|details]].",
			"",
			"## Details",
		}, "\n")
		writeFile(t, sourcePath, sourceContent)

		result, err := Promote(context.Background(), PromoteInput{
			SourceDocPath: sourcePath,
			VaultPath:     vaultPath,
			TargetTopic: models.TopicInfo{
				Slug:     "catalog",
				Mode:     models.TopicModeOKF,
				RootPath: targetTopic,
			},
			Type:        "Playbook",
			Description: "Operational alpha note.",
			Types:       []string{"Playbook"},
			Clock: func() time.Time {
				return time.Date(2026, 6, 27, 10, 11, 12, 0, time.UTC)
			},
		})
		if err != nil {
			t.Fatalf("Promote returned error: %v", err)
		}
		if result.WrittenPath != "alpha-note.md" {
			t.Fatalf("written path = %q, want alpha-note.md", result.WrittenPath)
		}
		if result.LinksRewritten != 2 {
			t.Fatalf("links rewritten = %d, want 2", result.LinksRewritten)
		}
		if len(result.UnresolvedLinks) != 1 || result.UnresolvedLinks[0] != "beta-note.md" {
			t.Fatalf("unresolved links = %#v, want beta-note.md", result.UnresolvedLinks)
		}
		if got := readFile(t, sourcePath); got != sourceContent {
			t.Fatalf("source document changed:\n%s", got)
		}

		values, body := parseMarkdown(t, filepath.Join(targetTopic, "alpha-note.md"))
		for key, want := range map[string]any{
			"type":        "Playbook",
			"title":       "Alpha Note",
			"description": "Operational alpha note.",
			"timestamp":   "2026-06-27T10:11:12Z",
		} {
			if got := values[key]; got != want {
				t.Fatalf("frontmatter[%s] = %#v, want %#v", key, got, want)
			}
		}
		if _, ok := values["stage"]; ok {
			t.Fatalf("wiki stage leaked into OKF frontmatter: %#v", values)
		}
		if !strings.Contains(body, "[Beta](beta-note.md)") || !strings.Contains(body, "[details](alpha-note.md#details)") {
			t.Fatalf("body missing transformed links:\n%s", body)
		}

		index := readFile(t, filepath.Join(targetTopic, "index.md"))
		for _, fragment := range []string{"okf_version: \"0.1\"", "## Playbook", "[Alpha Note](alpha-note.md)", "Operational alpha note."} {
			if !strings.Contains(index, fragment) {
				t.Fatalf("index.md missing %q:\n%s", fragment, index)
			}
		}
		log := readFile(t, filepath.Join(targetTopic, "log.md"))
		for _, fragment := range []string{"## 2026-06-27", "**Creation**", "[Alpha Note](alpha-note.md)", "`wiki/concepts/Alpha Note.md`"} {
			if !strings.Contains(log, fragment) {
				t.Fatalf("log.md missing %q:\n%s", fragment, log)
			}
		}
	})
}

func TestPromoteRejectsNonOKFTargetBeforeWriting(t *testing.T) {
	t.Parallel()

	t.Run("Should reject non-OKF target before writing", func(t *testing.T) {
		targetTopic := t.TempDir()
		_, err := Promote(context.Background(), PromoteInput{
			SourceDocPath: "missing.md",
			TargetTopic: models.TopicInfo{
				Slug:     "wiki-topic",
				Mode:     models.TopicModeWiki,
				RootPath: targetTopic,
			},
			Type: "Reference",
		})
		if err == nil || !strings.Contains(err.Error(), "target topic must use mode okf") {
			t.Fatalf("error = %v, want non-OKF target rejection", err)
		}
		if entries, readErr := os.ReadDir(targetTopic); readErr != nil || len(entries) != 0 {
			t.Fatalf("target changed before rejection: entries=%d err=%v", len(entries), readErr)
		}
	})
}

func TestPromoteWarnsWhenSourceBodyHasNoSentenceFallback(t *testing.T) {
	t.Parallel()

	t.Run("Should warn when source body has no sentence fallback", func(t *testing.T) {
		vaultPath := t.TempDir()
		sourceTopic := filepath.Join(vaultPath, "research")
		targetTopic := filepath.Join(vaultPath, "catalog")
		mkdirAll(t, filepath.Join(sourceTopic, "wiki", "concepts"))
		mkdirAll(t, targetTopic)
		writeFile(t, filepath.Join(sourceTopic, "CLAUDE.md"), "# Research\n")
		writeFile(t, filepath.Join(targetTopic, "CLAUDE.md"), "# Catalog\n")
		writeFile(t, filepath.Join(targetTopic, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
		writeFile(t, filepath.Join(targetTopic, "log.md"), "# Directory Update Log\n")
		sourcePath := filepath.Join(sourceTopic, "wiki", "concepts", "No Sentence.md")
		writeFile(t, sourcePath, "---\ntitle: No Sentence\n---\nbody text without punctuation")

		result, err := Promote(context.Background(), PromoteInput{
			SourceDocPath: sourcePath,
			VaultPath:     vaultPath,
			TargetTopic: models.TopicInfo{
				Slug:     "catalog",
				Mode:     models.TopicModeOKF,
				RootPath: targetTopic,
			},
			Type: "Reference",
			Clock: func() time.Time {
				return time.Date(2026, 6, 27, 10, 11, 12, 0, time.UTC)
			},
		})
		if err != nil {
			t.Fatalf("Promote returned error: %v", err)
		}
		if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "source body has no sentence fallback") {
			t.Fatalf("warnings = %#v, want no sentence fallback warning", result.Warnings)
		}

		values, _ := parseMarkdown(t, filepath.Join(targetTopic, "no-sentence.md"))
		if got := values["description"]; got != "" {
			t.Fatalf("description = %#v, want empty string", got)
		}
	})
}

func TestCheckReportsConformanceAndStrictWarnings(t *testing.T) {
	t.Parallel()

	writeBundle := func(t *testing.T) string {
		t.Helper()
		bundle := t.TempDir()
		writeFile(t, filepath.Join(bundle, "CLAUDE.md"), "# Catalog\n")
		writeFile(t, filepath.Join(bundle, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
		writeFile(t, filepath.Join(bundle, "log.md"), "# Directory Update Log\n\n## 2026-06-27\n* **Creation**: Created bundle.\n")
		writeFile(t, filepath.Join(bundle, "good.md"), "---\ntype: Playbook\ntitle: Good\ndescription: Good concept.\ntimestamp: 2026-06-27T10:00:00Z\n---\nBody with [broken](missing.md).\n")
		writeFile(t, filepath.Join(bundle, "missing-type.md"), "---\ntitle: Missing\n---\nBody.\n")
		writeFile(t, filepath.Join(bundle, "missing-fields.md"), "---\ntype: Unknown\n---\nBody.\n")
		return bundle
	}

	testCases := []struct {
		name       string
		strict     bool
		assertions []expectedIssue
	}{
		{
			name:   "Should report lenient conformance diagnostics",
			strict: false,
			assertions: []expectedIssue{
				{severity: models.SeverityError, filePath: "missing-type.md", target: "type", messageContains: "concept frontmatter must include non-empty type"},
				{severity: models.SeverityWarning, filePath: "missing-fields.md", target: "title", messageContains: `producer field "title" is missing`},
				{severity: models.SeverityWarning, filePath: "missing-fields.md", target: "type", messageContains: `type "Unknown" is outside the configured OKF vocabulary`},
			},
		},
		{
			name:   "Should promote warnings to errors in strict mode",
			strict: true,
			assertions: []expectedIssue{
				{severity: models.SeverityError, filePath: "missing-fields.md", target: "title", messageContains: `producer field "title" is missing`},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bundle := writeBundle(t)
			issues, err := Check(context.Background(), bundle, CheckOptions{
				Types:  []string{"Playbook"},
				Strict: testCase.strict,
			})
			if err != nil {
				t.Fatalf("Check returned error: %v", err)
			}
			for _, assertion := range testCase.assertions {
				assertIssue(t, issues, assertion)
			}
			if !HasErrors(issues) {
				t.Fatal("HasErrors = false, want true for missing type")
			}
		})
	}
}

func TestCheckAllowsLenientExternalBundleTraits(t *testing.T) {
	t.Parallel()

	t.Run("Should allow lenient external bundle traits", func(t *testing.T) {
		bundle := t.TempDir()
		writeFile(t, filepath.Join(bundle, "index.md"), "# Index\n")
		writeFile(t, filepath.Join(bundle, "concept.md"), "---\ntype: External Type\n---\nSee [missing](missing.md).\n")

		issues, err := Check(context.Background(), bundle, CheckOptions{})
		if err != nil {
			t.Fatalf("Check returned error: %v", err)
		}
		if HasErrors(issues) {
			t.Fatalf("lenient external bundle should not have errors: %#v", issues)
		}
	})
}

func TestAllocateConceptPathReservesPath(t *testing.T) {
	t.Parallel()

	t.Run("Should reserve concept paths with exclusive creation", func(t *testing.T) {
		bundle := t.TempDir()
		firstRelativePath, firstAbsolutePath, firstFile, err := allocateConceptPath(bundle, "alpha")
		if err != nil {
			t.Fatalf("allocate first path: %v", err)
		}
		t.Cleanup(func() { _ = firstFile.Close() })
		if firstRelativePath != "alpha.md" {
			t.Fatalf("first relative path = %q, want alpha.md", firstRelativePath)
		}
		if _, err := os.Stat(firstAbsolutePath); err != nil {
			t.Fatalf("reserved first path does not exist: %v", err)
		}

		secondRelativePath, _, secondFile, err := allocateConceptPath(bundle, "alpha")
		if err != nil {
			t.Fatalf("allocate second path: %v", err)
		}
		t.Cleanup(func() { _ = secondFile.Close() })
		if secondRelativePath != "alpha-2.md" {
			t.Fatalf("second relative path = %q, want alpha-2.md", secondRelativePath)
		}
	})
}

func TestRegenerateIndexReturnsConceptParseErrors(t *testing.T) {
	t.Parallel()

	t.Run("Should return parse errors instead of dropping malformed concepts", func(t *testing.T) {
		bundle := t.TempDir()
		writeFile(t, filepath.Join(bundle, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
		writeFile(t, filepath.Join(bundle, "bad.md"), "---\ntype: Playbook\n\nBody without closed frontmatter.\n")

		err := RegenerateIndex(bundle)
		if err == nil {
			t.Fatal("expected malformed concept frontmatter error")
		}
		if !strings.Contains(err.Error(), "parse concept bad.md") {
			t.Fatalf("error = %q, want parse concept bad.md", err.Error())
		}
	})
}

type expectedIssue struct {
	severity        models.DiagnosticSeverity
	filePath        string
	target          string
	messageContains string
}

func assertIssue(t *testing.T, issues []models.LintIssue, expected expectedIssue) {
	t.Helper()
	for _, issue := range issues {
		if issue.Severity == expected.severity &&
			issue.FilePath == expected.filePath &&
			issue.Target == expected.target &&
			strings.Contains(issue.Message, expected.messageContains) {
			return
		}
	}
	t.Fatalf(
		"missing issue severity=%s file=%s target=%s message containing %q in %#v",
		expected.severity,
		expected.filePath,
		expected.target,
		expected.messageContains,
		issues,
	)
}

func parseMarkdown(t *testing.T, filePath string) (map[string]any, string) {
	t.Helper()
	values, body, err := frontmatter.Parse(readFile(t, filePath))
	if err != nil {
		t.Fatalf("parse %s: %v", filePath, err)
	}
	return values, body
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
}

func writeFile(t *testing.T, filePath string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatalf("mkdir parent for %s: %v", filePath, err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", filePath, err)
	}
}

func readFile(t *testing.T, filePath string) string {
	t.Helper()
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read %s: %v", filePath, err)
	}
	return string(content)
}

type stubSuggester struct {
	suggestion TypeSuggestion
}

func (s *stubSuggester) SuggestType(_ context.Context, _ TypeDocument) (TypeSuggestion, error) {
	return s.suggestion, nil
}

func decided(probs map[string]float64) *TypeSuggestion {
	suggestion := withRanking(TypeSuggestion{Status: decisions.StatusDecided, ReceiptKey: "receipt-1"}, probs)
	return &suggestion
}

func TestPromoteTypeSuggestion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		flagType    string
		types       []string
		suggestion  *TypeSuggestion
		wantType    string
		wantErr     string
		wantWarning string
		wantQueued  bool
	}{
		{name: "no type and no vocabulary", wantErr: "[okf].types is empty"},
		{name: "no type and no decision model", types: []string{"Playbook"}, wantErr: "needs the decision model"},
		{
			name: "confident suggestion is used", types: []string{"Playbook", "Reference"},
			suggestion: decided(map[string]float64{"Playbook": 0.9, "Reference": 0.06, "none": 0.04}),
			wantType:   "Playbook",
		},
		{
			name: "review band asks for --type and queues", types: []string{"Playbook", "Reference"},
			suggestion: decided(map[string]float64{"Playbook": 0.6, "Reference": 0.3, "none": 0.1}),
			wantErr:    "top: Playbook 0.60, Reference 0.30, none 0.10", wantQueued: true,
		},
		{
			name: "low confidence asks for --type without queueing", types: []string{"Playbook", "Reference"},
			suggestion: decided(map[string]float64{"Playbook": 0.4, "Reference": 0.35, "none": 0.25}),
			wantErr:    "no type reaches P ≥ 0.80",
		},
		{
			name: "none fits", types: []string{"Playbook"},
			suggestion: decided(map[string]float64{"Playbook": 0.1, "none": 0.9}),
			wantErr:    "no configured OKF type fits",
		},
		{
			name: "undecided suggestion", types: []string{"Playbook"},
			suggestion: &TypeSuggestion{Status: decisions.StatusUndecided, Reason: decisions.ReasonTimeout},
			wantErr:    "undecided (timeout)",
		},
		{
			name: "explicit type disagreeing at 0.8 warns", flagType: "Reference", types: []string{"Playbook", "Reference"},
			suggestion: decided(map[string]float64{"Playbook": 0.85, "Reference": 0.1, "none": 0.05}),
			wantType:   "Reference", wantWarning: `suggests type "Playbook" (P=0.85), not "Reference"`,
		},
		{
			name: "explicit type with weak disagreement stays quiet", flagType: "Reference", types: []string{"Playbook", "Reference"},
			suggestion: decided(map[string]float64{"Playbook": 0.7, "Reference": 0.2, "none": 0.1}),
			wantType:   "Reference",
		},
		{name: "explicit type without vocabulary", flagType: "Playbook", wantType: "Playbook"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vaultPath, sourceTopic, targetTopic, sourcePath := newPromoteFixture(t)
			input := PromoteInput{
				SourceDocPath: sourcePath,
				VaultPath:     vaultPath,
				TargetTopic:   models.TopicInfo{Slug: "catalog", Mode: models.TopicModeOKF, RootPath: targetTopic},
				Type:          tt.flagType,
				Types:         tt.types,
			}
			if tt.suggestion != nil {
				input.Suggester = &stubSuggester{suggestion: *tt.suggestion}
			}

			result, err := Promote(context.Background(), input)
			if tt.wantErr != "" {
				if err == nil || !errors.Is(err, ErrTypeRequired) || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Promote() error = %v, want ErrTypeRequired containing %q", err, tt.wantErr)
				}
				if _, statErr := os.Stat(filepath.Join(targetTopic, "alpha-note.md")); !os.IsNotExist(statErr) {
					t.Fatalf("concept written despite error: %v", statErr)
				}
				pending, loadErr := review.Open(sourceTopic, nil).Pending(review.QueueOKFType)
				if loadErr != nil {
					t.Fatalf("Pending: %v", loadErr)
				}
				if queued := len(pending) == 1; queued != tt.wantQueued {
					t.Fatalf("okf-type review items = %#v, want queued=%t", pending, tt.wantQueued)
				}
				if tt.wantQueued && (pending[0].Subject != "wiki/concepts/Alpha Note.md" || pending[0].Target != "catalog" || pending[0].ReceiptKey != "receipt-1" || !strings.Contains(err.Error(), pending[0].ID)) {
					t.Fatalf("review item = %#v, error = %v", pending[0], err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Promote() error = %v", err)
			}
			if result.Type != tt.wantType {
				t.Fatalf("type = %q, want %q", result.Type, tt.wantType)
			}
			values, _ := parseMarkdown(t, filepath.Join(targetTopic, result.WrittenPath))
			if values["type"] != tt.wantType {
				t.Fatalf("concept type = %#v, want %q", values["type"], tt.wantType)
			}
			if (tt.wantWarning == "") != (result.TypeWarning == "") || !strings.Contains(result.TypeWarning, tt.wantWarning) {
				t.Fatalf("type warning = %q, want %q", result.TypeWarning, tt.wantWarning)
			}
			if (tt.suggestion != nil) != (result.TypeSuggestion != nil) {
				t.Fatalf("type suggestion = %#v", result.TypeSuggestion)
			}
		})
	}
}

func TestEngineTypeSuggesterAsksOKFTypeChoice(t *testing.T) {
	t.Parallel()

	fake := fakes.NewOpenRouter(func(call fakes.Call, q fakes.Question) any {
		return fakes.Dist(map[string]float64{"Playbook": 0.7, "Reference": 0.2, "none": 0.1})
	}, nil)
	t.Cleanup(fake.Close)
	suggester := EngineTypeSuggester{
		Decider: newTestEngine(t, fake),
		Topic:   decisions.TopicRef{Slug: "research", Root: t.TempDir(), Thresholds: decisions.DefaultThresholds()},
		Options: TypeOptions([]string{"Playbook", "Reference", "Playbook"}, func(name string) string {
			if name == "Playbook" {
				return "A step-by-step procedure"
			}
			return ""
		}),
	}

	suggestion, err := suggester.SuggestType(context.Background(), TypeDocument{Subject: "wiki/concepts/Alpha.md", Title: "Alpha", Body: "1. Do this.\n2. Then that.\n"})
	if err != nil {
		t.Fatalf("SuggestType: %v", err)
	}
	if !suggestion.Decided() || suggestion.Type != "Playbook" || suggestion.Probability != 0.7 || len(suggestion.Top) != 3 || suggestion.ReceiptKey == "" {
		t.Fatalf("suggestion = %#v", suggestion)
	}

	calls := fake.Calls()
	if len(calls) != 1 {
		t.Fatalf("calls = %d, want 1", len(calls))
	}
	question, ok := calls[0].Questions["okf_type"]
	if !ok {
		t.Fatalf("questions = %#v, want okf_type", calls[0].Questions)
	}
	criteria, _ := question.Criteria.(map[string]any)
	if criteria["Playbook"] != "A step-by-step procedure" || criteria["Reference"] != "Reference" || criteria["none"] == nil || len(criteria) != 3 {
		t.Fatalf("okf_type options = %#v", criteria)
	}
	if state := calls[0].StateString(); !strings.Contains(state, `"title":"Alpha"`) || !strings.Contains(state, "Do this.") {
		t.Fatalf("state = %s", state)
	}
}

// TestOKFDecisionsExcludeAndExtras: a document matching decisions.exclude
// is never sent (promote needs --type, check reports no advisory finding),
// and the topic's extra okf_type questions ride along in the request.
func TestOKFDecisionsExcludeAndExtras(t *testing.T) {
	t.Parallel()
	fake := fakes.NewOpenRouter(func(call fakes.Call, q fakes.Question) any {
		if q.ID == "okf_type" {
			return fakes.Pick("Playbook", 0.9)
		}
		return nil
	}, nil)
	t.Cleanup(fake.Close)
	extra, err := questions.Parse([]byte(`{"id":"owner_okf","version":"1","purpose":"okf_type","guard":"G","questions":[{"id":"is_runbook","type":"noul","instructions":"Is it a runbook?","source":"topic owner"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	options := TypeOptions([]string{"Playbook", "Reference"}, nil)
	ref := decisions.TopicRef{Slug: "research", Root: t.TempDir(), Thresholds: decisions.DefaultThresholds(), Exclude: []string{"wiki/private/**"}}
	suggester := EngineTypeSuggester{Decider: newTestEngine(t, fake), Topic: ref, Options: options, Extras: []*questions.Bank{extra}}

	excluded, err := suggester.SuggestType(context.Background(), TypeDocument{Subject: "wiki/private/Secret.md", Title: "Secret", Body: "1. Do this."})
	if err != nil {
		t.Fatal(err)
	}
	if excluded.Status != decisions.StatusNotChecked || excluded.Reason != decisions.ReasonExcluded || len(fake.Calls()) != 0 {
		t.Fatalf("excluded suggestion = %#v (calls %d)", excluded, len(fake.Calls()))
	}
	_, err = chooseType(context.Background(), PromoteInput{Suggester: &stubSuggester{suggestion: excluded}}, typeContext{document: TypeDocument{Subject: "wiki/private/Secret.md"}})
	if !errors.Is(err, ErrTypeRequired) || !strings.Contains(err.Error(), "decisions.exclude") {
		t.Fatalf("chooseType error = %v", err)
	}

	judged, err := suggester.SuggestType(context.Background(), TypeDocument{Subject: "wiki/concepts/Alpha.md", Title: "Alpha", Body: "1. Do this."})
	if err != nil || !judged.Decided() || judged.Type != "Playbook" {
		t.Fatalf("judged = %#v, %v", judged, err)
	}
	calls := fake.Calls()
	if len(calls) != 1 {
		t.Fatalf("calls = %d", len(calls))
	}
	if _, ok := calls[0].Questions["is_runbook"]; !ok || len(calls[0].Questions) != 2 {
		t.Fatalf("questions = %#v", calls[0].Questions)
	}

	bundle := t.TempDir()
	writeFile(t, filepath.Join(bundle, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
	mkdirAll(t, filepath.Join(bundle, "private"))
	writeFile(t, filepath.Join(bundle, "private", "rotate-key.md"), "---\ntitle: Rotate Key\ntype: Reference\ndescription: Rotating cut errors by 40%.\ntimestamp: 2026-06-27T10:11:12Z\n---\n1. Create a key. 2. Restart workers.\n")
	before := len(fake.Calls())
	advisory := &AdvisoryOptions{Options: options, Decider: newTestEngine(t, fake), Topic: decisions.TopicRef{Slug: "catalog", Root: bundle, Exclude: []string{"private/**"}}}
	issues, err := Check(context.Background(), bundle, CheckOptions{Types: []string{"Playbook", "Reference"}, Advisory: advisory})
	if err != nil {
		t.Fatal(err)
	}
	for _, issue := range issues {
		if issue.Kind == models.LintIssueKindTypeMismatch || issue.Kind == models.LintIssueKindDescriptionUnsupported {
			t.Fatalf("an excluded concept got an advisory finding: %#v", issue)
		}
	}
	if len(fake.Calls()) != before {
		t.Fatalf("an excluded concept was sent: %d new calls", len(fake.Calls())-before)
	}
}

func TestCheckAdvisoryFindings(t *testing.T) {
	t.Parallel()

	fake := fakes.NewOpenRouter(func(call fakes.Call, q fakes.Question) any {
		switch q.ID {
		case "okf_type":
			return fakes.Pick("Playbook", 0.9)
		case "description_unsupported":
			return fakes.Noul(0.85)
		}
		return nil
	}, nil)
	t.Cleanup(fake.Close)

	bundle := t.TempDir()
	writeFile(t, filepath.Join(bundle, "index.md"), "---\nokf_version: \"0.1\"\n---\n# Index\n")
	conceptPath := filepath.Join(bundle, "rotate-key.md")
	writeFile(t, conceptPath, "---\ntitle: Rotate Key\ntype: Reference\ndescription: Rotating cut errors by 40%.\ntimestamp: 2026-06-27T10:11:12Z\n---\n1. Create a key. 2. Restart workers.\n")
	options := []TypeOption{{Name: "Playbook", Description: "Procedure"}, {Name: "Reference", Description: "Facts"}}

	check := func(decider Decider) []models.LintIssue {
		t.Helper()
		advisory := &AdvisoryOptions{Options: options, Decider: decider}
		if decider != nil {
			advisory.Topic = decisions.TopicRef{Slug: "catalog", Root: bundle, Thresholds: decisions.DefaultThresholds()}
		}
		issues, err := Check(context.Background(), bundle, CheckOptions{Types: []string{"Playbook", "Reference"}, Strict: true, Advisory: advisory})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}
		if HasErrors(issues) {
			t.Fatalf("advisory findings must never fail the check: %#v", issues)
		}
		return issues
	}
	advisoryKinds := func(issues []models.LintIssue) string {
		kinds := make([]string, 0)
		for _, issue := range issues {
			if issue.Kind == models.LintIssueKindTypeMismatch || issue.Kind == models.LintIssueKindDescriptionUnsupported {
				if issue.Severity != models.SeverityInfo || issue.FilePath != "rotate-key.md" {
					t.Fatalf("advisory issue = %#v", issue)
				}
				kinds = append(kinds, string(issue.Kind))
			}
		}
		sort.Strings(kinds)
		return strings.Join(kinds, ",")
	}

	if got := advisoryKinds(check(nil)); got != "" || len(fake.Calls()) != 0 {
		t.Fatalf("receipts-only check without receipts: findings %q, calls %d", got, len(fake.Calls()))
	}
	want := "description_unsupported,type_mismatch"
	if got := advisoryKinds(check(newTestEngine(t, fake))); got != want || len(fake.Calls()) != 2 {
		t.Fatalf("--decide findings = %q (calls %d), want %q after 2 calls", got, len(fake.Calls()), want)
	}
	if got := advisoryKinds(check(nil)); got != want || len(fake.Calls()) != 2 {
		t.Fatalf("receipts-only findings = %q (calls %d), want %q without new calls", got, len(fake.Calls()), want)
	}
	if got := advisoryKinds(check(newTestEngine(t, fake))); got != want || len(fake.Calls()) != 2 {
		t.Fatalf("second --decide findings = %q (calls %d), want cached answers", got, len(fake.Calls()))
	}

	writeFile(t, conceptPath, "---\ntitle: Rotate Key\ntype: Reference\ndescription: Rotating cut errors by 40%.\ntimestamp: 2026-06-27T10:11:12Z\n---\nA table of limits.\n")
	if got := advisoryKinds(check(nil)); got != "" {
		t.Fatalf("edited concept reused stale receipts: %q", got)
	}
}

func newPromoteFixture(t *testing.T) (vaultPath, sourceTopic, targetTopic, sourcePath string) {
	t.Helper()
	vaultPath = t.TempDir()
	sourceTopic = filepath.Join(vaultPath, "research")
	targetTopic = filepath.Join(vaultPath, "catalog")
	mkdirAll(t, filepath.Join(sourceTopic, "wiki", "concepts"))
	mkdirAll(t, targetTopic)
	writeFile(t, filepath.Join(sourceTopic, "CLAUDE.md"), "# Research\n")
	writeFile(t, filepath.Join(targetTopic, "CLAUDE.md"), "# Catalog\n")
	sourcePath = filepath.Join(sourceTopic, "wiki", "concepts", "Alpha Note.md")
	writeFile(t, sourcePath, "---\ntitle: Alpha Note\ntype: wiki\nstage: compiled\n---\nAlpha note explains the operational flow.\n")
	return vaultPath, sourceTopic, targetTopic, sourcePath
}

func newTestEngine(t *testing.T, fake *fakes.OpenRouter) *decisions.Engine {
	t.Helper()
	engine, err := decisions.New(decisions.Options{Config: config.Default().Decisions, APIKey: "test-key", APIURL: fake.URL})
	if err != nil {
		t.Fatalf("decisions.New: %v", err)
	}
	return engine
}
