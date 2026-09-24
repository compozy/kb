package questions

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// wantBanks fixes the bank and question ids other packages code against.
var wantBanks = map[string]struct {
	purpose   string
	questions []string
}{
	"relevance":      {"relevance", []string{"role", "role_item_{id}"}},
	"quality":        {"quality", []string{"paywall_or_login", "error_or_placeholder_page", "thin_or_boilerplate", "no_speech_content"}},
	"duplicate":      {"duplicate", []string{"relation_{id}"}},
	"classify":       {"classify", []string{"kind", "depth", "enough_text_to_judge"}},
	"concept":        {"concept", []string{"primary_concept", "mentions_concept_{id}"}},
	"link":           {"link", []string{"should_link_{id}", "relation_{id}", "mention_sense_{n}", "affects_{id}"}},
	"find":           {"find", []string{"answers_{id}", "specificity_{id}", "question_kind"}},
	"okf_type":       {"okf_type", []string{"okf_type", "description_unsupported"}},
	"contract_check": {"contract_check", []string{"conflict_{n}"}},
}

func loadAllEmbedded(t *testing.T) []*Bank {
	t.Helper()
	ids := BankIDs()
	banks := make([]*Bank, 0, len(ids))
	for _, id := range ids {
		bank, err := Load(id)
		if err != nil {
			t.Fatalf("Load(%q): %v", id, err)
		}
		banks = append(banks, bank)
	}
	return banks
}

func TestEmbeddedBanksLoad(t *testing.T) {
	t.Parallel()

	gotIDs := BankIDs()
	wantIDs := make([]string, 0, len(wantBanks))
	for id := range wantBanks {
		wantIDs = append(wantIDs, id)
	}
	slices.Sort(wantIDs)
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("BankIDs() = %v, want %v", gotIDs, wantIDs)
	}
	for _, bank := range loadAllEmbedded(t) {
		t.Run("Should load bank "+bank.ID+" with fixed ids, guard and hash", func(t *testing.T) {
			t.Parallel()
			want := wantBanks[bank.ID]
			if bank.Purpose != want.purpose {
				t.Errorf("purpose = %q, want %q", bank.Purpose, want.purpose)
			}
			if !reflect.DeepEqual(bank.IDs(), want.questions) {
				t.Errorf("IDs() = %v, want %v", bank.IDs(), want.questions)
			}
			if bank.Guard != Guard {
				t.Errorf("guard = %q, want the package Guard", bank.Guard)
			}
			if bank.Version == "" || len(bank.Hash) != 64 {
				t.Errorf("version %q / hash %q not set", bank.Version, bank.Hash)
			}
			if again := MustLoad(bank.ID); again != bank {
				t.Error("Load should cache banks")
			}
		})
	}
}

// TestEmbeddedBanksFollowQuestionRules is the syntactic polarity guard of
// jev-engineering questions.md over every embedded bank.
func TestEmbeddedBanksFollowQuestionRules(t *testing.T) {
	t.Parallel()

	openers := []string{"Does ", "Is ", "Are ", "Has ", "Can "}
	forbidden := []string{"verify", "check", "ensure", "confirm", "make sure"}
	for _, bank := range loadAllEmbedded(t) {
		for _, id := range bank.IDs() {
			raw := bank.questions[id]
			t.Run(bank.ID+"/"+id, func(t *testing.T) {
				t.Parallel()
				if strings.TrimSpace(raw.Source) == "" {
					t.Error("source must be non-empty")
				}
				text, ok := raw.text.(string)
				if !ok {
					t.Fatalf("instructions must be a string, got %T", raw.text)
				}
				for sentence := range strings.SplitSeq(text, ". ") {
					lower := strings.ToLower(strings.TrimSpace(sentence))
					for _, word := range forbidden {
						if strings.HasPrefix(lower, word) {
							t.Errorf("sentence opens with forbidden %q: %q", word, sentence)
						}
					}
				}
				switch raw.Type {
				case TypeNoul:
					sentences := strings.Split(text, ". ")
					question := strings.TrimSpace(sentences[len(sentences)-1])
					if !strings.HasSuffix(question, "?") {
						t.Errorf("noul question must end with '?': %q", question)
					}
					if !slices.ContainsFunc(openers, func(opener string) bool { return strings.HasPrefix(question, opener) }) {
						t.Errorf("noul question must open with Does/Is/Are/Has/Can: %q", question)
					}
					criteria, _ := raw.Criteria.(map[string]any)
					if criteria["true"] == nil || criteria["false"] == nil {
						t.Error("embedded nouls document both true and false so polarity is explicit")
					}
				case TypeChoice:
					criteria, _ := raw.Criteria.(map[string]any)
					if !slices.ContainsFunc(ExitOptions, func(exit string) bool { _, ok := criteria[exit]; return ok }) {
						t.Errorf("choice has no exit option: %v", criteria)
					}
				}
				if strings.Contains(id, "{") {
					placeholder := placeholderPattern.FindString(id)
					if !strings.Contains(text, "`"+placeholder+"`") {
						t.Errorf("template question must name its element (%s) in the instructions", placeholder)
					}
				}
			})
		}
	}
}

func TestQuestionSubstitutesPlaceholdersAndAppliesGuard(t *testing.T) {
	t.Parallel()

	link := MustLoad("link")
	testCases := []struct {
		name       string
		templateID string
		vars       map[string]string
		wantID     string
		wantText   string
		wantErr    string
	}{
		{
			name:       "Should substitute id placeholder in id and instructions",
			templateID: "should_link_{id}",
			vars:       map[string]string{"id": "c7"},
			wantID:     "should_link_c7",
			wantText:   "Evaluate only candidate `c7` in `candidates`.",
		},
		{
			name:       "Should substitute n placeholder",
			templateID: "mention_sense_{n}",
			vars:       map[string]string{"n": "3"},
			wantID:     "mention_sense_3",
			wantText:   "Evaluate only mention `3` in `mentions`.",
		},
		{
			name:       "Should reject a missing placeholder value",
			templateID: "affects_{id}",
			vars:       map[string]string{"n": "1"},
			wantErr:    "needs a value for {id}",
		},
		{
			name:       "Should reject an unknown question",
			templateID: "nope",
			wantErr:    "unknown question",
		},
		{
			name:       "Should reject braces inside a value",
			templateID: "affects_{id}",
			vars:       map[string]string{"id": "{x}"},
			wantErr:    "must not contain braces",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q, err := link.Question(tc.templateID, tc.vars)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Question: %v", err)
			}
			if q.ID != tc.wantID {
				t.Errorf("ID = %q, want %q", q.ID, tc.wantID)
			}
			text, _ := q.Instructions.(string)
			if !strings.HasPrefix(text, Guard+" ") {
				t.Errorf("instructions must start with the guard: %q", text)
			}
			if !strings.Contains(text, tc.wantText) || strings.Contains(text, "{") {
				t.Errorf("instructions = %q, want substituted %q", text, tc.wantText)
			}
		})
	}

	t.Run("Should hand out deep copies", func(t *testing.T) {
		t.Parallel()
		first := MustLoad("duplicate").MustQuestion("relation_{id}", map[string]string{"id": "a"})
		first.Criteria.(map[string]any)["same_content"] = "mutated"
		second := MustLoad("duplicate").MustQuestion("relation_{id}", map[string]string{"id": "a"})
		if second.Criteria.(map[string]any)["same_content"] == "mutated" {
			t.Fatal("mutating a returned question changed the bank")
		}
	})
}

func TestApplyGuardKeepsInstructionShape(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input any
		want  any
	}{
		{name: "Should prefix a string", input: "Does x?", want: "G Does x?"},
		{name: "Should add a guard key to an object", input: map[string]any{"task": "Does x?"}, want: map[string]any{"guard": "G", "task": "Does x?"}},
		{name: "Should prepend to an array", input: []any{"Does x?"}, want: []any{"G", "Does x?"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := applyGuard("G", tc.input); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("applyGuard = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestWithCriteriaKeepsExitOptions(t *testing.T) {
	t.Parallel()

	concept := MustLoad("concept")
	base := concept.MustQuestion("primary_concept", nil)
	injected := concept.WithCriteria(base, map[string]string{
		"a1": "Agent protocols — A2A, MCP — documents comparing agent-to-agent protocols",
	})
	criteria, ok := injected.Criteria.(map[string]any)
	if !ok {
		t.Fatalf("criteria type = %T", injected.Criteria)
	}
	if _, ok := criteria["a1"]; !ok {
		t.Error("injected option missing")
	}
	if _, ok := criteria["none"]; !ok {
		t.Error("bank exit option none must be merged in")
	}
	if _, ok := base.Criteria.(map[string]any)["a1"]; ok {
		t.Error("WithCriteria must not mutate the input question")
	}
}

func TestParseRejectsInvalidBanks(t *testing.T) {
	t.Parallel()

	const header = `"id":"x","version":"1","purpose":"p","guard":"G",`
	testCases := []struct {
		name    string
		bank    string
		wantErr string
	}{
		{"missing guard", `{"id":"x","version":"1","purpose":"p","questions":[{"id":"a","type":"noul","instructions":"Is it?","source":"s"}]}`, "guard is required"},
		{"missing purpose", `{"id":"x","version":"1","guard":"G","questions":[{"id":"a","type":"noul","instructions":"Is it?","source":"s"}]}`, "purpose is required"},
		{"no questions", `{` + header + `"questions":[]}`, "at least one question"},
		{"missing source", `{` + header + `"questions":[{"id":"a","type":"noul","instructions":"Is it?"}]}`, "source is required"},
		{"duplicate id", `{` + header + `"questions":[{"id":"a","type":"noul","instructions":"Is it?","source":"s"},{"id":"a","type":"noul","instructions":"Is it?","source":"s"}]}`, "duplicate question id"},
		{"bad type", `{` + header + `"questions":[{"id":"a","type":"bool","instructions":"Is it?","source":"s"}]}`, "type must be noul, choice or score"},
		{"bad placeholder", `{` + header + `"questions":[{"id":"a_{x}","type":"noul","instructions":"Is it?","source":"s"}]}`, "placeholders"},
		{"empty instructions", `{` + header + `"questions":[{"id":"a","type":"noul","instructions":" ","source":"s"}]}`, "instructions are empty"},
		{"noul bad criteria key", `{` + header + `"questions":[{"id":"a","type":"noul","instructions":"Is it?","criteria":{"yes":"y"},"source":"s"}]}`, "true or false"},
		{"choice without exit", `{` + header + `"questions":[{"id":"a","type":"choice","instructions":"Which?","criteria":{"x":"1","y":"2"},"source":"s"}]}`, "exit option"},
		{"choice with one option", `{` + header + `"questions":[{"id":"a","type":"choice","instructions":"Which?","criteria":{"none":"n"},"source":"s"}]}`, "2 to 255 options"},
		{"score with one level", `{` + header + `"questions":[{"id":"a","type":"score","instructions":"How?","criteria":["only"],"source":"s"}]}`, "2 to 10 levels"},
		{"score with eleven levels", `{` + header + `"questions":[{"id":"a","type":"score","instructions":"How?","criteria":["0","1","2","3","4","5","6","7","8","9","10"],"source":"s"}]}`, "2 to 10 levels"},
		{"injected options on noul", `{` + header + `"questions":[{"id":"a","type":"noul","injected_options":true,"instructions":"Is it?","source":"s"}]}`, "injected_options"},
		{"unknown field", `{` + header + `"extra":1,"questions":[{"id":"a","type":"noul","instructions":"Is it?","source":"s"}]}`, "unknown field"},
	}
	for _, tc := range testCases {
		t.Run("Should reject "+tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := Parse([]byte(tc.bank))
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Parse err = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestParseHashIsCanonical(t *testing.T) {
	t.Parallel()

	a, err := Parse([]byte(`{"id":"x","version":"1","purpose":"p","guard":"G","questions":[{"id":"a","type":"noul","instructions":"Is it?","source":"s"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	b, err := Parse([]byte(`{"guard":"G","purpose":"p","version":"1","id":"x",
		"questions":[{"source":"s","instructions":"Is it?","type":"noul","id":"a"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	c, err := Parse([]byte(`{"id":"x","version":"1","purpose":"p","guard":"G","questions":[{"id":"a","type":"noul","instructions":"Is it now?","source":"s"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if a.Hash != b.Hash {
		t.Error("key order and whitespace must not change the hash")
	}
	if a.Hash == c.Hash {
		t.Error("a wording change must change the hash")
	}
}

func TestLoadTopicExtras(t *testing.T) {
	t.Parallel()

	const extra = `{"id":"extra_quality","version":"1","purpose":"quality","guard":"G","questions":[{"id":"%s","type":"noul","instructions":"Is it a press release?","source":"topic owner"}]}`
	testCases := []struct {
		name    string
		files   map[string]string
		wantIDs []string
		wantErr string
	}{
		{name: "Should return nothing without a banks directory"},
		{
			name:    "Should load a valid extra bank",
			files:   map[string]string{"extra.json": strings.Replace(extra, "%s", "press_release", 1), "notes.txt": "ignored"},
			wantIDs: []string{"extra_quality"},
		},
		{
			name:    "Should reject an id colliding with a built-in question of the same purpose",
			files:   map[string]string{"extra.json": strings.Replace(extra, "%s", "thin_or_boilerplate", 1)},
			wantErr: "collides with a built-in quality question",
		},
		{
			name:    "Should reject an extra bank without purpose",
			files:   map[string]string{"extra.json": `{"id":"e","version":"1","guard":"G","questions":[{"id":"a","type":"noul","instructions":"Is it?","source":"s"}]}`},
			wantErr: "purpose is required",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			if len(tc.files) > 0 {
				dir := filepath.Join(root, ".decisions", "banks")
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				for name, content := range tc.files {
					if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}
			banks, err := LoadTopicExtras(root)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadTopicExtras: %v", err)
			}
			gotIDs := make([]string, 0, len(banks))
			for _, bank := range banks {
				gotIDs = append(gotIDs, bank.ID)
			}
			if len(gotIDs) != len(tc.wantIDs) || (len(gotIDs) > 0 && !reflect.DeepEqual(gotIDs, tc.wantIDs)) {
				t.Fatalf("banks = %v, want %v", gotIDs, tc.wantIDs)
			}
		})
	}
}
