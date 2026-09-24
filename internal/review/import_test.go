package review

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
)

func writeTopicFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixedClock() time.Time { return time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC) }

func screeningTopic(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "sports")
	writeTopicFile(t, root, "CLAUDE.md", "# Sports\n")
	writeTopicFile(t, root, "raw/papers/a.md", "---\ntitle: A\npmcid: PMC100\nsource_url: https://www.example.org/a/?utm_source=x\n---\nA body.\n")
	writeTopicFile(t, root, "raw/papers/b.md", "---\ntitle: B\nsource_url: https://pmc.ncbi.nlm.nih.gov/articles/PMC200/\n---\nB body.\n")
	writeTopicFile(t, root, "raw/papers/c.md", "---\ntitle: C\n---\nPMCID: pmc300 and doi 10.1234/abc.def\n")
	writeTopicFile(t, root, "wiki/concepts/Topic.md", "---\ntitle: Topic\npmcid: PMC999\n---\nArticle.\n")
	return root
}

func labelsBySubject(t *testing.T, root string) map[string]Label {
	t.Helper()
	labels, err := LoadLabels(root)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]Label{}
	for _, label := range labels {
		out[label.Subject+"|"+label.Target] = label
	}
	return out
}

func TestImportLabels(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		file          string
		content       string
		opts          ImportOptions
		wantVerdicts  map[string]string
		wantUnmatched []string
		wantDecidedBy map[string]string
	}{
		{
			name: "jsonl by pmcid from frontmatter, url and body head",
			file: "science-screening.jsonl",
			content: `{"pmcid":"PMC100","decision":"include","screening_stage":"title"}
{"pmcid":"PMC200","decision":"exclude","screening_stage":"abstract"}
{"pmcid":"300","decision":"INCLUDE"}
{"pmcid":"PMC999","decision":"include"}
{"pmcid":"PMC404","decision":"include"}
`,
			opts:          ImportOptions{IDField: "pmcid", Match: MatchPMCID, DecisionField: "decision", KeepValues: []string{"include", "true"}},
			wantVerdicts:  map[string]string{"raw/papers/a.md": VerdictPositive, "raw/papers/b.md": VerdictNegative, "raw/papers/c.md": VerdictPositive},
			wantUnmatched: []string{"PMC999", "PMC404"},
			wantDecidedBy: map[string]string{"raw/papers/a.md": "title", "raw/papers/b.md": "abstract"},
		},
		{
			name:          "csv by url with explicit decided-by field and boolean decisions",
			file:          "curation.csv",
			content:       "url,included,who\nhttp://example.org/a,true,owner\nhttps://pmc.ncbi.nlm.nih.gov/articles/PMC200,false,owner\nhttps://nowhere.example/x,true,owner\n",
			opts:          ImportOptions{IDField: "url", Match: MatchURL, DecisionField: "included", DecidedByField: "who"},
			wantVerdicts:  map[string]string{"raw/papers/a.md": VerdictPositive, "raw/papers/b.md": VerdictNegative},
			wantUnmatched: []string{"https://nowhere.example/x"},
			wantDecidedBy: map[string]string{"raw/papers/a.md": "owner"},
		},
		{
			name:          "json array by path and doi",
			file:          "decisions.json",
			content:       `[{"path":"raw/papers/c.md","keep":true},{"path":"raw/papers/missing.md","keep":false},{"path":"wiki/concepts/Topic.md","keep":true}]`,
			opts:          ImportOptions{IDField: "path", Match: MatchPath, DecisionField: "keep"},
			wantVerdicts:  map[string]string{"raw/papers/c.md": VerdictPositive},
			wantUnmatched: []string{"raw/papers/missing.md", "wiki/concepts/Topic.md"},
		},
		{
			name:         "doi match from body head",
			file:         "doi.jsonl",
			content:      `{"doi":"https://doi.org/10.1234/ABC.DEF","decision":"exclude"}` + "\n",
			opts:         ImportOptions{IDField: "doi", Match: MatchDOI, DecisionField: "decision"},
			wantVerdicts: map[string]string{"raw/papers/c.md": VerdictNegative},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := screeningTopic(t)
			from := filepath.Join(t.TempDir(), tc.file)
			if err := os.WriteFile(from, []byte(tc.content), 0o644); err != nil {
				t.Fatal(err)
			}
			opts := tc.opts
			opts.From = from
			opts.Now = fixedClock

			report, err := ImportLabels(root, opts)
			if err != nil {
				t.Fatalf("ImportLabels: %v", err)
			}
			if report.Matched != len(tc.wantVerdicts) || !slices.Equal(report.UnmatchedIDs, tc.wantUnmatched) {
				t.Fatalf("report = %+v", report)
			}
			got := labelsBySubject(t, root)
			if len(got) != len(tc.wantVerdicts) {
				t.Fatalf("labels = %+v", got)
			}
			for subject, verdict := range tc.wantVerdicts {
				label := got[subject+"|"]
				if label.Verdict != verdict || label.Purpose != PurposeRelevance || label.Question != QuestionRole || label.Origin != report.Origin {
					t.Errorf("label for %s = %+v", subject, label)
				}
			}
			for subject, who := range tc.wantDecidedBy {
				if got[subject+"|"].DecidedBy != who {
					t.Errorf("decided_by for %s = %q, want %q", subject, got[subject+"|"].DecidedBy, who)
				}
			}

			again, err := ImportLabels(root, opts)
			if err != nil {
				t.Fatal(err)
			}
			if again.Positive+again.Negative != 0 || again.Duplicates != len(tc.wantVerdicts) {
				t.Fatalf("re-import wrote labels: %+v", again)
			}
			all, _ := LoadLabels(root)
			if len(all) != len(tc.wantVerdicts) {
				t.Fatalf("labels after re-import = %d", len(all))
			}
		})
	}
}

func TestImportLabelsRejectsBadOptions(t *testing.T) {
	t.Parallel()
	root := screeningTopic(t)
	for _, opts := range []ImportOptions{
		{From: "x.jsonl", IDField: "id", Match: "isbn", DecisionField: "d"},
		{From: "x.jsonl", Match: MatchURL, DecisionField: "d"},
		{From: filepath.Join(root, "missing.jsonl"), IDField: "id", Match: MatchURL, DecisionField: "d"},
	} {
		if _, err := ImportLabels(root, opts); err == nil {
			t.Errorf("ImportLabels(%+v) succeeded", opts)
		}
	}
}

func TestNormalizeURL(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"https://www.Example.org/a/?utm_source=x#frag": "example.org/a",
		"http://example.org:80/a":                      "example.org/a",
		"https://example.org/a?b=2&a=1":                "example.org/a?a=1&b=2",
		"https://example.org:8443/a":                   "example.org:8443/a",
		"":                                             "",
	}
	for input, want := range testCases {
		if got := NormalizeURL(input); got != want {
			t.Errorf("NormalizeURL(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestImportLinks(t *testing.T) {
	t.Parallel()
	vault := t.TempDir()
	root := filepath.Join(vault, "jev")
	writeTopicFile(t, root, "CLAUDE.md", "# Jev\n")
	writeTopicFile(t, root, "raw/articles/src.md", "---\ntitle: Src\n---\nSource.\n")
	writeTopicFile(t, root, "raw/_quarantine/articles/gone.md", "---\ntitle: Gone\n---\nGone.\n")
	for _, name := range []string{"B", "D", "E", "F"} {
		writeTopicFile(t, root, "wiki/concepts/"+name+".md", "---\ntitle: "+name+"\n---\n"+name+" body.\n")
	}
	writeTopicFile(t, root, "wiki/concepts/A.md", `---
title: A
sources:
  - "[[raw/articles/src]]"
---
See [[B]] and [[E|the E]] and [[gone]] and [[A]] and [[Missing]].
`+"```\n[[F]]\n```\n")
	writeTopicFile(t, root, "wiki/concepts/F.md", "---\ntitle: F\nextends:\n  - \"[[B]]\"\n---\nF links [[B]] twice [[B]].\n")

	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatal(err)
	}
	c, err := corpus.Load(root, corpus.LoadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	writer := corpus.NewWriter(root, store, fixedClock)
	doc := c.ByPath("wiki/concepts/A.md")
	result, err := writer.Apply(doc, map[string]any{"related": []string{"[[D]]"}}, corpus.StateMeta{})
	if err != nil || result.Status != corpus.StatusWritten {
		t.Fatalf("kb relation write = %+v, %v", result, err)
	}
	inserted, _ := json.Marshal(insertedLink{Time: "t", Subject: "wiki/concepts/A.md", Target: "wiki/concepts/E.md", Text: "the E", Mode: "apply"})
	shadow, _ := json.Marshal(insertedLink{Time: "t", Subject: "wiki/concepts/F.md", Target: "B", Mode: "shadow"})
	writeTopicFile(t, root, filepath.Join(decisions.ReceiptsDir, InsertedLinksFile), string(inserted)+"\n"+string(shadow)+"\n")

	c, err = corpus.Load(root, corpus.LoadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	count, err := ImportLinks(vault, root, c)
	if err != nil {
		t.Fatalf("ImportLinks: %v", err)
	}
	got := labelsBySubject(t, root)
	want := []string{
		"wiki/concepts/A.md|wiki/concepts/B.md",
		"wiki/concepts/A.md|raw/articles/src.md",
		"wiki/concepts/F.md|wiki/concepts/B.md",
	}
	if count != len(want) || len(got) != len(want) {
		t.Fatalf("count = %d, labels = %+v", count, got)
	}
	for _, key := range want {
		label, ok := got[key]
		if !ok || label.Purpose != PurposeLink || label.Question != QuestionShouldLink || label.Verdict != VerdictPositive || label.Origin != OriginImportLinks {
			t.Errorf("label %s = %+v (present %v)", key, label, ok)
		}
	}

	again, err := ImportLinks(vault, root, nil)
	if err != nil || again != 0 {
		t.Fatalf("second import = %d, %v", again, err)
	}
}

func TestPreviewsRoundTrip(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := AppendPreview(root, Preview{Contract: "c1", Subjects: []string{"raw/a.md"}}); err != nil {
		t.Fatal(err)
	}
	if err := AppendPreview(root, Preview{Time: "2026-01-01T00:00:00Z", Contract: "c2", Subjects: []string{"raw/b.md"}}); err != nil {
		t.Fatal(err)
	}
	previews, err := LoadPreviews(root)
	if err != nil || len(previews) != 2 || previews[0].Time == "" || previews[1].Contract != "c2" {
		t.Fatalf("previews = %+v, %v", previews, err)
	}
}
