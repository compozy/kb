package link

import (
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

var testThresholds = Thresholds{Apply: 0.85, Review: 0.60, Mention: 0.85, Affects: 0.80}

func noul(p float64) decisions.Answer {
	return decisions.Answer{Type: "noul", Status: decisions.StatusDecided, Noul: &p, Receipt: "rk"}
}

func choice(option string, p, confidence float64) decisions.Answer {
	return decisions.Answer{
		Type: "choice", Status: decisions.StatusDecided, Choice: option,
		Probs: map[string]float64{option: p}, Confidence: &confidence, Receipt: "rk",
	}
}

func undecided() decisions.Answer {
	return decisions.Answer{Type: "noul", Status: decisions.StatusUndecided, Reason: decisions.ReasonTimeout}
}

// basePlan is a source with one candidate c1 (RAG) mentioned twice.
func basePlan(kind corpus.Kind) Plan {
	doc := testDoc("raw/articles/a.md", kind, nil, "RAG one. RAG two.")
	return Plan{
		Doc:        doc,
		Candidates: []Candidate{{ID: "c1", Path: "wiki/concepts/RAG.md", Title: "RAG"}},
		Mentions: []MentionRef{
			{N: 1, Candidate: "c1", Start: 0, End: 3, Text: "RAG"},
			{N: 2, Candidate: "c1", Start: 9, End: 12, Text: "RAG"},
		},
		Forms:      map[string]string{"wiki/concepts/RAG.md": "RAG"},
		Relations:  map[string][]string{},
		Present:    map[string]map[string]bool{},
		KBOwned:    map[string]bool{},
		BodyLinked: map[string]bool{},
	}
}

func TestDecide(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		kind    corpus.Kind
		mode    string
		mutate  func(*Plan)
		answers map[string]decisions.Answer
		updates map[string]any
		inserts []int // insertion starts
		queues  []string
		actions []map[string]any
		// questions, when set, are the producing question ids of the items.
		questions []string
		counters  [5]int // proposals, reviews, contradictions, demotions, undecided
	}{
		{
			name:    "apply band writes argmax relation",
			kind:    corpus.KindArticle,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("extends", 0.9, 0.88), "mention_sense_1": noul(0.1), "mention_sense_2": noul(0.1)},
			updates: map[string]any{"extends": []string{"[[RAG]]"}},
		},
		{
			name:    "none falls back to related",
			kind:    corpus.KindArticle,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("none", 0.7, 0.6), "mention_sense_1": noul(0.1), "mention_sense_2": noul(0.1)},
			updates: map[string]any{"related": []string{"[[RAG]]"}},
		},
		{
			name:    "low relation confidence falls back to related",
			kind:    corpus.KindArticle,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("prerequisite", 0.4, 0.3), "mention_sense_1": noul(0.1), "mention_sense_2": noul(0.1)},
			updates: map[string]any{"related": []string{"[[RAG]]"}},
		},
		{
			name:    "existing entries are kept and appended to",
			kind:    corpus.KindArticle,
			mutate:  func(p *Plan) { p.Relations["related"] = []string{"[[Other]]"} },
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.95), "relation_c1": choice("related", 0.9, 0.9), "mention_sense_1": noul(0.1), "mention_sense_2": noul(0.1)},
			updates: map[string]any{"related": []string{"[[Other]]", "[[RAG]]"}},
		},
		{
			name: "already related under another key is not re-added",
			kind: corpus.KindArticle,
			mutate: func(p *Plan) {
				p.Relations["related"] = []string{"[[RAG]]"}
				p.Present["related"] = map[string]bool{"wiki/concepts/RAG.md": true}
			},
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.95), "relation_c1": choice("extends", 0.9, 0.9), "mention_sense_1": noul(0.1), "mention_sense_2": noul(0.1)},
			updates: map[string]any{},
		},
		{
			name:     "contradicts never applies",
			kind:     corpus.KindArticle,
			answers:  map[string]decisions.Answer{"should_link_c1": noul(0.97), "relation_c1": choice("contradicts", 0.9, 0.9), "mention_sense_1": noul(0.99)},
			updates:  map[string]any{},
			queues:   []string{review.QueueContradiction},
			actions:  []map[string]any{{"relation": "contradicts", "target": "RAG"}},
			counters: [5]int{0, 0, 1, 0, 0},
		},
		{
			name:      "review band queues a link item",
			kind:      corpus.KindArticle,
			answers:   map[string]decisions.Answer{"should_link_c1": noul(0.7), "relation_c1": choice("example_of", 0.9, 0.9)},
			updates:   map[string]any{},
			queues:    []string{review.QueueLink},
			actions:   []map[string]any{{"relation": "example_of", "target": "RAG", "candidate": "c1"}},
			questions: []string{"should_link_c1"},
			counters:  [5]int{0, 1, 0, 0, 0},
		},
		{
			name:     "undecided relation in the apply band is never guessed",
			kind:     corpus.KindArticle,
			mode:     session.ModeApply,
			answers:  map[string]decisions.Answer{"should_link_c1": noul(0.95), "relation_c1": undecided(), "mention_sense_1": noul(0.99)},
			updates:  map[string]any{},
			counters: [5]int{0, 0, 0, 0, 1},
		},
		{
			name:     "undecided relation in the review band queues nothing",
			kind:     corpus.KindArticle,
			answers:  map[string]decisions.Answer{"should_link_c1": noul(0.7), "relation_c1": undecided()},
			updates:  map[string]any{},
			counters: [5]int{0, 0, 0, 0, 1},
		},
		{
			name:     "undecided mention_sense is counted and blocks a later insertion",
			kind:     corpus.KindArticle,
			mode:     session.ModeApply,
			answers:  map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("related", 0.9, 0.9), "mention_sense_1": undecided(), "mention_sense_2": noul(0.95)},
			updates:  map[string]any{"related": []string{"[[RAG]]"}},
			counters: [5]int{0, 0, 0, 0, 1},
		},
		{
			name:    "ignore band does nothing",
			kind:    corpus.KindArticle,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.2), "relation_c1": choice("related", 0.9, 0.9)},
			updates: map[string]any{},
		},
		{
			name: "kb-written relation below apply is demoted, not removed",
			kind: corpus.KindArticle,
			mutate: func(p *Plan) {
				p.Relations["extends"] = []string{"[[RAG]]"}
				p.Present["extends"] = map[string]bool{"wiki/concepts/RAG.md": true}
				p.KBOwned["extends"] = true
			},
			answers:  map[string]decisions.Answer{"should_link_c1": noul(0.3), "relation_c1": choice("extends", 0.9, 0.9)},
			updates:  map[string]any{},
			queues:   []string{review.QueueLink},
			actions:  []map[string]any{{"relation": "extends", "target": "RAG", "demote": true}},
			counters: [5]int{0, 0, 0, 1, 0},
		},
		{
			name: "human-written relation below apply is left alone",
			kind: corpus.KindArticle,
			mutate: func(p *Plan) {
				p.Relations["related"] = []string{"[[RAG]]"}
				p.Present["related"] = map[string]bool{"wiki/concepts/RAG.md": true}
			},
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.1), "relation_c1": choice("related", 0.9, 0.9)},
			updates: map[string]any{},
		},
		{
			name:    "affects on a source",
			kind:    corpus.KindSource,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.1), "affects_c1": noul(0.85)},
			updates: map[string]any{"affects": []string{"[[RAG]]"}},
		},
		{
			name:    "affects below threshold",
			kind:    corpus.KindSource,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.1), "affects_c1": noul(0.79)},
			updates: map[string]any{},
		},
		{
			name:     "undecided is counted, never a no",
			kind:     corpus.KindSource,
			answers:  map[string]decisions.Answer{"should_link_c1": undecided(), "affects_c1": undecided()},
			updates:  map[string]any{},
			counters: [5]int{0, 0, 0, 0, 2},
		},
		{
			name:    "apply mode inserts at first qualifying mention",
			kind:    corpus.KindArticle,
			mode:    session.ModeApply,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("related", 0.9, 0.9), "mention_sense_1": noul(0.5), "mention_sense_2": noul(0.9)},
			updates: map[string]any{"related": []string{"[[RAG]]"}},
			inserts: []int{9},
		},
		{
			name:     "shadow mode proposes instead of inserting",
			kind:     corpus.KindArticle,
			mode:     session.ModeShadow,
			answers:  map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("related", 0.9, 0.9), "mention_sense_1": noul(0.9)},
			updates:  map[string]any{"related": []string{"[[RAG]]"}},
			queues:   []string{review.QueueLink},
			actions:  []map[string]any{{"insert": true, "relation": "related", "target": "RAG", "mention_text": "RAG", "mention_start": 0}},
			counters: [5]int{1, 0, 0, 0, 0},
		},
		{
			name:    "no insertion below mention_sense",
			kind:    corpus.KindArticle,
			mode:    session.ModeApply,
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("related", 0.9, 0.9), "mention_sense_1": noul(0.84), "mention_sense_2": noul(0.2)},
			updates: map[string]any{"related": []string{"[[RAG]]"}},
		},
		{
			name:    "no insertion when the body already links the target",
			kind:    corpus.KindArticle,
			mode:    session.ModeApply,
			mutate:  func(p *Plan) { p.BodyLinked["wiki/concepts/RAG.md"] = true },
			answers: map[string]decisions.Answer{"should_link_c1": noul(0.9), "relation_c1": choice("related", 0.9, 0.9), "mention_sense_1": noul(0.99)},
			updates: map[string]any{"related": []string{"[[RAG]]"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			plan := basePlan(tt.kind)
			if tt.mutate != nil {
				tt.mutate(&plan)
			}
			mode := tt.mode
			if mode == "" {
				mode = session.ModeShadow
			}
			out := Decide(plan, tt.answers, testThresholds, mode)

			if !reflect.DeepEqual(out.Updates, tt.updates) {
				t.Fatalf("updates = %#v, want %#v", out.Updates, tt.updates)
			}
			starts := make([]int, 0)
			for _, insertion := range out.Insertions {
				starts = append(starts, insertion.Start)
			}
			if len(starts) != len(tt.inserts) || (len(starts) > 0 && !reflect.DeepEqual(starts, tt.inserts)) {
				t.Fatalf("insertions at %v, want %v", starts, tt.inserts)
			}
			if len(out.Items) != len(tt.queues) {
				t.Fatalf("items = %+v, want queues %v", out.Items, tt.queues)
			}
			for index, item := range out.Items {
				if item.Queue != tt.queues[index] || item.Subject != plan.Doc.Path || item.Target != "wiki/concepts/RAG.md" {
					t.Fatalf("item %d = %+v", index, item)
				}
				if !reflect.DeepEqual(item.Action, tt.actions[index]) {
					t.Fatalf("item %d action = %#v, want %#v", index, item.Action, tt.actions[index])
				}
				if tt.questions != nil {
					if item.Question != tt.questions[index] || item.ReceiptKey != "rk" {
						t.Fatalf("item %d question/receipt = %q/%q, want %q/rk", index, item.Question, item.ReceiptKey, tt.questions[index])
					}
					// The queue identity stays on the generic question so a
					// positional candidate id never re-asks a resolved item.
					if want := review.ItemID(item.Queue, item.Subject, item.Target, review.QuestionShouldLink); item.ID != want {
						t.Fatalf("item %d id = %q, want %q", index, item.ID, want)
					}
				}
			}
			got := [5]int{out.Proposals, out.Reviews, out.Contradictions, out.Demotions, out.Undecided}
			if got != tt.counters {
				t.Fatalf("counters (proposals, reviews, contradictions, demotions, undecided) = %v, want %v", got, tt.counters)
			}
		})
	}
}

func TestMentionInsertionPlacement(t *testing.T) {
	t.Parallel()
	body := strings.Join([]string{
		"# Vector Database overview",
		"",
		"```",
		"Vector Database in code",
		"```",
		"",
		"Inline `Vector Database` code and a [[Vector Database|link]] already.",
		"",
		"A real vector database mention. Another Vector Database later.",
	}, "\n")
	target := article("Vector Database", nil, "x")
	doc := testDoc("raw/articles/a.md", corpus.KindSource, nil, body)
	generator := NewGenerator([]*corpus.Document{target}, nil)
	candidates := []Candidate{{ID: "c1", Path: target.Path}}
	mentions := generator.Mentions(doc, candidates)
	if len(mentions) != 2 {
		t.Fatalf("mentions = %+v, want 2 plain-text mentions", mentions)
	}

	plan := Plan{
		Doc: doc, Candidates: candidates, Forms: map[string]string{target.Path: "Vector Database"},
		Relations: map[string][]string{}, Present: map[string]map[string]bool{}, KBOwned: map[string]bool{}, BodyLinked: map[string]bool{},
	}
	for index, mention := range mentions {
		plan.Mentions = append(plan.Mentions, MentionRef{N: index + 1, Candidate: "c1", Start: mention.Start, End: mention.End, Text: body[mention.Start:mention.End]})
	}
	answers := map[string]decisions.Answer{
		"should_link_c1": noul(0.9), "relation_c1": choice("related", 0.9, 0.9),
		"mention_sense_1": noul(0.9), "mention_sense_2": noul(0.9), "affects_c1": noul(0),
	}
	out := Decide(plan, answers, testThresholds, session.ModeApply)
	if len(out.Insertions) != 1 {
		t.Fatalf("insertions = %+v, want one per target", out.Insertions)
	}
	got := InsertLinks(body, out.Insertions)
	want := strings.Replace(body, "A real vector database mention.", "A real [[Vector Database|vector database]] mention.", 1)
	if got != want {
		t.Fatalf("InsertLinks:\n%s\nwant:\n%s", got, want)
	}
	if strings.Count(got, "[[Vector Database|") != 2 {
		t.Fatalf("expected the existing link plus exactly one inserted link:\n%s", got)
	}
}

func TestRelationOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		answer decisions.Answer
		want   string
	}{
		{name: "argmax", answer: choice("example_of", 0.8, 0.7), want: "example_of"},
		{name: "none", answer: choice("none", 0.8, 0.7), want: RelationRelated},
		{name: "low confidence", answer: choice("extends", 0.4, 0.2), want: RelationRelated},
		{name: "contradicts keeps its name", answer: choice("contradicts", 0.4, 0.2), want: RelationContradicts},
		{name: "undecided", answer: undecided(), want: RelationRelated},
	}
	for _, tt := range tests {
		if got, _ := relationOf(tt.answer); got != tt.want {
			t.Errorf("%s: relationOf = %q, want %q", tt.name, got, tt.want)
		}
	}
}
