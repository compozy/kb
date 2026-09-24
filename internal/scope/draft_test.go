package scope

import (
	"context"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
)

func TestSummarizeScreening(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "outputs/datasets/science-screening.jsonl"), strings.Join([]string{
		`{"pmcid":"PMC1","title":"Jump IMU validity","included":true,"exclusion_reason":null}`,
		`{"pmcid":"PMC2","title":"Sprint GPS","included":true}`,
		`{"pmcid":"PMC3","title":"Knee surgery outcomes","included":false,"exclusion_reason":"treatment-without-measurement"}`,
		`{"pmcid":"PMC4","title":"Drug trial","included":false,"exclusion_reason":"treatment-without-measurement"}`,
		`{"pmcid":"PMC5","title":"Robot navigation","included":false,"exclusion_reason":"too-remote-transfer"}`,
		`not json`,
		``,
	}, "\n"))
	writeFile(t, filepath.Join(root, "outputs/collection/curation-decisions.json"), `[
  {"id":"s-1","title":"OracleTrust","reason":"Transaction signature verification, not oracle economics."},
  {"id":"s-2","title":"Building Trust","reason":"Earlier preprint of the same work."}
]`)
	writeFile(t, filepath.Join(root, "outputs/collection/exclusions.jsonl"), `{"url":"https://a","title":"A","reason":"duplicate-work"}
{"url":"https://b","title":"B","reason":"duplicate-work"}
{"url":"https://c","title":"C","reason":"secondary-host"}
`)
	writeFile(t, filepath.Join(root, "outputs/wearables/decision-screening.jsonl"), `{"pmcid":"PMC9","title":"ECG patch","decision":"include","reason":"hardware-method-transfer"}
{"pmcid":"PMC8","title":"Mouse study","decision":"exclude","reason":"unrelated-biology"}
`)
	writeFile(t, filepath.Join(root, "outputs/reports/summary.json"), `[{"title":"ignored"}]`)
	writeFile(t, filepath.Join(root, "outputs/.cache/screening.jsonl"), `{"title":"hidden","included":true}`)

	got, err := SummarizeScreening(root, 1)
	if err != nil {
		t.Fatalf("SummarizeScreening: %v", err)
	}
	want := []ScreeningSummary{
		{File: "outputs/collection/curation-decisions.json", Rows: 2, Decisions: []DecisionSummary{
			{Decision: DecisionExclude, Count: 2, Reasons: []ReasonCount{{Reason: "Earlier preprint of the same work.", Count: 1}}, Examples: []string{"OracleTrust", "Building Trust"}},
		}},
		{File: "outputs/collection/exclusions.jsonl", Rows: 3, Decisions: []DecisionSummary{
			{Decision: DecisionExclude, Count: 3, Reasons: []ReasonCount{{Reason: "duplicate-work", Count: 2}}, Examples: []string{"A", "B", "C"}},
		}},
		{File: "outputs/datasets/science-screening.jsonl", Rows: 5, Decisions: []DecisionSummary{
			{Decision: DecisionExclude, Count: 3, Reasons: []ReasonCount{{Reason: "treatment-without-measurement", Count: 2}}, Examples: []string{"Knee surgery outcomes", "Drug trial", "Robot navigation"}},
			{Decision: DecisionInclude, Count: 2, Reasons: []ReasonCount{}, Examples: []string{"Jump IMU validity", "Sprint GPS"}},
		}},
		{File: "outputs/wearables/decision-screening.jsonl", Rows: 2, Decisions: []DecisionSummary{
			{Decision: DecisionExclude, Count: 1, Reasons: []ReasonCount{{Reason: "unrelated-biology", Count: 1}}, Examples: []string{"Mouse study"}},
			{Decision: DecisionInclude, Count: 1, Reasons: []ReasonCount{{Reason: "hardware-method-transfer", Count: 1}}, Examples: []string{"ECG patch"}},
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SummarizeScreening =\n%#v\nwant\n%#v", got, want)
	}

	none, err := SummarizeScreening(t.TempDir(), 8)
	if err != nil || len(none) != 0 {
		t.Fatalf("topic without outputs = %v, %v", none, err)
	}
}

func TestIsScreeningFile(t *testing.T) {
	t.Parallel()
	tests := map[string]bool{
		"outputs/datasets/science-screening.jsonl":   true,
		"screening.jsonl":                            true,
		"outputs/collection/curation-decisions.json": true,
		"exclusions.jsonl":                           true,
		"science-screening.json":                     false,
		"summary.jsonl":                              false,
	}
	for name, want := range tests {
		if got := IsScreeningFile(name); got != want {
			t.Errorf("IsScreeningFile(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestScopeLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, markdown, want string
	}{
		{name: "template placeholder", markdown: "# T\n\n**Topic scope:** one-paragraph description of what this topic covers.\n", want: ""},
		{name: "bold topic scope", markdown: "# T\n\n**Topic scope:** Decision models and judges.\n", want: "Decision models and judges."},
		{name: "portuguese escopo", markdown: "# T\r\n\r\nEscopo: ciência de sensores.\r\n", want: "ciência de sensores."},
		{name: "missing", markdown: "# T\n\nNothing here.\n", want: ""},
	}
	for _, tt := range tests {
		if got := ScopeLine(tt.markdown); got != tt.want {
			t.Errorf("%s: ScopeLine = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestCleanDraft(t *testing.T) {
	t.Parallel()
	cleaned, dropped := CleanDraft(&contract.Contract{
		Purpose:                 " Purpose. ",
		Core:                    []string{"Jev", "jev ", "Judges"},
		Adjacent:                []string{"Evaluation"},
		CollectedOnPurposePaths: []string{"raw/audit/**", "raw/[bad"},
		OutOfScope:              []string{"Sports", "judges"},
	})
	want := &contract.Contract{
		Purpose:                 "Purpose.",
		Core:                    []string{"Jev", "Judges"},
		Adjacent:                []string{"Evaluation"},
		CollectedOnPurpose:      []string{},
		CollectedOnPurposePaths: []string{"raw/audit/**"},
		OutOfScope:              []string{"Sports"},
	}
	if !reflect.DeepEqual(cleaned, want) {
		t.Fatalf("CleanDraft = %#v, want %#v", cleaned, want)
	}
	if len(dropped) != 3 {
		t.Fatalf("dropped = %v, want 3 notes (duplicate, bad glob, kept line in out_of_scope)", dropped)
	}
	if err := cleaned.Validate(); err != nil {
		t.Fatalf("cleaned draft does not validate: %v", err)
	}
}

func TestCollectDraftInputsAndPrompt(t *testing.T) {
	t.Parallel()
	_, root := newTestTopic(t)
	writeFile(t, filepath.Join(root, "CLAUDE.md"), "# Demo\n\n**Topic scope:** Typed decision models.\n\n## Selection contract\n\n- **Purpose:** Old purpose.\n- **Core:** Jev\n- **Out of scope:** Sports\n\n## Other\n")
	for i := range 5 {
		writeSource(t, root, "raw/articles/a"+string(rune('0'+i))+".md", "Article "+string(rune('A'+i)), "Body.\n")
	}
	writeSource(t, root, "raw/audit/b.md", "Audit page", "Body.\n")
	writeSource(t, root, "wiki/concepts/Jev.md", "Jev", "Body.\n")
	writeFile(t, filepath.Join(root, "outputs/datasets/science-screening.jsonl"), `{"title":"Kept paper","included":true}`+"\n"+`{"title":"Dropped paper","included":false,"exclusion_reason":"off-scope"}`+"\n")
	c, err := corpus.Load(root, corpus.LoadOptions{})
	if err != nil {
		t.Fatal(err)
	}

	inputs, err := CollectDraftInputs(root, "Demo", "demo", c, DraftOptions{SourceSample: 3})
	if err != nil {
		t.Fatalf("CollectDraftInputs: %v", err)
	}
	if inputs.ScopeText != "Typed decision models." || inputs.ExistingContract == nil || inputs.ExistingContract.Purpose != "Old purpose." {
		t.Fatalf("scope = %q, existing = %+v", inputs.ScopeText, inputs.ExistingContract)
	}
	if !reflect.DeepEqual(inputs.Folders, []FolderCount{{Folder: "raw/articles", Files: 5}, {Folder: "raw/audit", Files: 1}}) {
		t.Fatalf("folders = %+v", inputs.Folders)
	}
	if inputs.Sources != 6 || len(inputs.SourceSample) != 3 || !reflect.DeepEqual(inputs.Articles, []string{"Jev"}) || len(inputs.Screening) != 1 {
		t.Fatalf("inputs = %+v", inputs)
	}
	sampledFolders := map[string]int{}
	for _, sample := range inputs.SourceSample {
		sampledFolders[sample.Folder]++
	}
	if sampledFolders["raw/audit"] != 1 {
		t.Fatalf("stratified sample misses raw/audit: %+v", inputs.SourceSample)
	}

	prompt := inputs.Prompt()
	for _, want := range append(slices.Clone(DraftRules),
		"Topic: Demo (domain: demo)",
		"Scope line from the topic CLAUDE.md: Typed decision models.",
		"- purpose: Old purpose.",
		"- raw/articles/: 5 files",
		"- raw/audit/: 1 files",
		"Source titles (3 of 6",
		"Audit page [raw/audit]",
		"- Jev",
		"outputs/datasets/science-screening.jsonl (2 rows)",
		"reason (1): off-scope",
		"example: Dropped paper",
	) {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt lacks %q:\n%s", want, prompt)
		}
	}
}

func TestDraftSavesValidatedDraft(t *testing.T) {
	t.Parallel()
	vault, root := newTestTopic(t)
	writeSource(t, root, "raw/articles/a.md", "Jev judges", "Body.\n")
	var prompts []string
	generated := map[string]any{
		"purpose":                    "Collect material on typed decision models.",
		"core":                       []string{"Typed decision models", "typed decision models"},
		"adjacent":                   []string{"LLM judges, kept as comparison"},
		"collected_on_purpose":       []string{"Vendor docs"},
		"collected_on_purpose_paths": []string{"raw/audit/**"},
		"out_of_scope":               []string{"Sports results with no decision-model content"},
	}
	fake := fakeServer(t, nil, func(schema, _, prompt string) any {
		prompts = append(prompts, schema+"\n"+prompt)
		return generated
	})
	s := openTestSession(t, vault, fake)

	draft, inputs, err := Draft(context.Background(), s, DraftOptions{})
	if err != nil {
		t.Fatalf("Draft: %v", err)
	}
	if len(prompts) != 1 || !strings.HasPrefix(prompts[0], "selection_contract\n") || !strings.Contains(prompts[0], DraftRules[2]) {
		t.Fatalf("generation prompts = %q", prompts)
	}
	if !reflect.DeepEqual(draft.Core, []string{"Typed decision models"}) || len(inputs.Dropped) != 1 {
		t.Fatalf("draft = %+v, dropped = %v", draft, inputs.Dropped)
	}
	settings, err := contract.LoadSettings(root)
	if err != nil {
		t.Fatal(err)
	}
	if settings.ContractDraft.Hash() != draft.Hash() || !settings.Contract.Empty() {
		t.Fatalf("topic.yaml draft = %+v, contract = %+v", settings.ContractDraft, settings.Contract)
	}

	// A new source changes the prompt, so the generation cache is not hit.
	generated["core"] = []string{}
	writeSource(t, root, "raw/articles/b.md", "Another source", "Body.\n")
	s.ReloadCorpus()
	if _, _, err := Draft(context.Background(), s, DraftOptions{}); err == nil || !strings.Contains(err.Error(), "failed validation") {
		t.Fatalf("invalid draft error = %v", err)
	}
	settings, _ = contract.LoadSettings(root)
	if settings.ContractDraft.Hash() != draft.Hash() {
		t.Fatal("an invalid draft replaced the saved one")
	}
}

func TestImportClaude(t *testing.T) {
	t.Parallel()
	_, root := newTestTopic(t)
	if _, err := ImportClaude(root); err == nil {
		t.Fatal("ImportClaude accepted a template section with no contract text")
	}

	claude := readFile(t, filepath.Join(root, "CLAUDE.md"))
	writeFile(t, filepath.Join(root, "CLAUDE.md"), contract.RenderIntoClaude(claude, testDraft))
	imported, err := ImportClaude(root)
	if err != nil {
		t.Fatalf("ImportClaude: %v", err)
	}
	if imported.Hash() != testDraft.Hash() {
		t.Fatalf("imported = %+v, want %+v", imported, testDraft)
	}
	settings, err := contract.LoadSettings(root)
	if err != nil {
		t.Fatal(err)
	}
	if settings.ContractDraft.Hash() != testDraft.Hash() || !settings.Contract.Empty() {
		t.Fatalf("settings = %+v", settings)
	}
}
