package classify

import (
	"maps"
	"reflect"
	"slices"
	"testing"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/quality"
)

func noul(p float64) decisions.Answer {
	return decisions.Answer{Type: "noul", Status: decisions.StatusDecided, Noul: &p}
}

func choice(option string, probs map[string]float64) decisions.Answer {
	return decisions.Answer{Type: "choice", Status: decisions.StatusDecided, Choice: option, Probs: probs}
}

func score(value float64) decisions.Answer {
	confidence := 1.0
	return decisions.Answer{Type: "score", Status: decisions.StatusDecided, Score: &value, Confidence: &confidence}
}

func failed(reason string) decisions.Answer {
	return decisions.Answer{Status: decisions.StatusUndecided, Reason: reason}
}

func testConcept(id, stem string) *concept {
	return &concept{id: id, title: stem, link: "[[" + stem + "]]", article: &corpus.Document{Path: "wiki/concepts/" + stem + ".md"}}
}

func baseAnswers() map[string]decisions.Answer {
	return map[string]decisions.Answer{
		"kind":                 choice("paper", map[string]float64{"paper": 0.9, "other": 0.1}),
		"enough_text_to_judge": noul(0.9),
		"depth":                score(2.26),
	}
}

func withAnswers(extra map[string]decisions.Answer) map[string]decisions.Answer {
	answers := baseAnswers()
	maps.Copy(answers, extra)
	return answers
}

func TestMapFacets(t *testing.T) {
	t.Parallel()
	c1, c2, c3 := testConcept("c01", "Alpha"), testConcept("c02", "Beta"), testConcept("c03", "Gamma")
	options := []*concept{c1, c2, c3}
	decidedGate := func(role string, pOff float64, q map[string]float64) *GateJudgment {
		return &GateJudgment{Role: role, RoleAsked: true, RoleDecided: true, POffTopic: pOff, Quality: q, QualityDecided: true, ReceiptKey: "k"}
	}

	tests := []struct {
		name          string
		in            facetInput
		wantUpdates   map[string]any
		wantComplete  bool
		wantUndecided []string
		wantRecapture string
		wantRemove    bool
		wantNone      bool
	}{
		{
			name:         "decided source writes every facet",
			in:           facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleCore, 0.05, map[string]float64{NoulPaywall: 0.1}), answers: baseAnswers()},
			wantUpdates:  map[string]any{"relevance": RoleCore, "genre": "paper", "depth": 2.3},
			wantComplete: true,
		},
		{
			name:         "not enough text skips depth",
			in:           facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleAdjacent, 0, nil), answers: withAnswers(map[string]decisions.Answer{"enough_text_to_judge": noul(0.3)})},
			wantUpdates:  map[string]any{"relevance": RoleAdjacent, "genre": "paper"},
			wantComplete: true,
		},
		{
			name:          "undecided depth never writes",
			in:            facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleCore, 0, nil), answers: withAnswers(map[string]decisions.Answer{"depth": failed("timeout")})},
			wantUpdates:   map[string]any{"relevance": RoleCore, "genre": "paper"},
			wantUndecided: []string{"depth:timeout"},
		},
		{
			name:          "undecided kind never writes genre",
			in:            facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleCore, 0, nil), answers: withAnswers(map[string]decisions.Answer{"kind": failed("budget")})},
			wantUpdates:   map[string]any{"relevance": RoleCore, "depth": 2.3},
			wantUndecided: []string{"kind:budget"},
		},
		{
			name:          "undecided role writes unknown only when absent, never off_topic",
			in:            facetInput{source: true, relevanceOn: true, gate: &GateJudgment{RoleAsked: true, Undecided: []string{"role:timeout"}, QualityDecided: true}, answers: baseAnswers()},
			wantUpdates:   map[string]any{"relevance": RoleUnknown, "genre": "paper", "depth": 2.3},
			wantUndecided: []string{"role:timeout"},
		},
		{
			name:          "undecided role keeps an existing relevance",
			in:            facetInput{source: true, relevanceOn: true, hasRelevance: true, gate: &GateJudgment{RoleAsked: true, Undecided: []string{"role:timeout"}, QualityDecided: true}, answers: baseAnswers()},
			wantUpdates:   map[string]any{"genre": "paper", "depth": 2.3},
			wantUndecided: []string{"role:timeout"},
		},
		{
			name:         "relevance off writes no relevance",
			in:           facetInput{source: true, relevanceOn: false, gate: &GateJudgment{QualityDecided: true}, answers: baseAnswers()},
			wantUpdates:  map[string]any{"genre": "paper", "depth": 2.3},
			wantComplete: true,
		},
		{
			name:         "collected path",
			in:           facetInput{source: true, relevanceOn: true, gate: &GateJudgment{Role: RoleCollectedOnPurpose, RoleByPath: true, RoleDecided: true, PKept: 1, QualityDecided: true}, answers: baseAnswers()},
			wantUpdates:  map[string]any{"relevance": RoleCollectedOnPurpose, "genre": "paper", "depth": 2.3},
			wantComplete: true,
		},
		{
			name:          "quality noul at apply writes quality and queues recapture",
			in:            facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleOffTopic, 0.9, map[string]float64{NoulErrorPage: 0.92}), answers: baseAnswers()},
			wantUpdates:   map[string]any{"relevance": RoleOffTopic, "quality": ReasonErrorPage, "genre": "paper", "depth": 2.3},
			wantComplete:  true,
			wantRecapture: ReasonErrorPage,
		},
		{
			name:          "quality noul in review band queues recapture without writing",
			in:            facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleCore, 0, map[string]float64{NoulPaywall: 0.6}), answers: baseAnswers()},
			wantUpdates:   map[string]any{"relevance": RoleCore, "genre": "paper", "depth": 2.3},
			wantComplete:  true,
			wantRecapture: ReasonPaywall,
		},
		{
			name:          "code flag writes quality",
			in:            facetInput{source: true, relevanceOn: true, gate: &GateJudgment{Role: RoleCore, RoleAsked: true, RoleDecided: true, QualityDecided: true, Flags: []quality.Flag{{Code: quality.NotAnArticle}}}, answers: baseAnswers()},
			wantUpdates:   map[string]any{"relevance": RoleCore, "quality": ReasonNotAnArticle, "genre": "paper", "depth": 2.3},
			wantComplete:  true,
			wantRecapture: ReasonNotAnArticle,
		},
		{
			name:         "off topic above review queues remove",
			in:           facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleOffTopic, 0.55, map[string]float64{NoulThin: 0.1}), answers: baseAnswers()},
			wantUpdates:  map[string]any{"relevance": RoleOffTopic, "genre": "paper", "depth": 2.3},
			wantComplete: true,
			wantRemove:   true,
		},
		{
			name:         "off topic below review is not queued",
			in:           facetInput{source: true, relevanceOn: true, gate: decidedGate(RoleGeneral, 0.4, nil), answers: baseAnswers()},
			wantUpdates:  map[string]any{"relevance": RoleGeneral, "genre": "paper", "depth": 2.3},
			wantComplete: true,
		},
		{
			name: "concepts from primary and nouls",
			in: facetInput{source: true, relevanceOn: false, gate: &GateJudgment{QualityDecided: true}, options: options, candidates: []*concept{c2, c3},
				answers: withAnswers(map[string]decisions.Answer{
					"primary_concept":      choice("c01", map[string]float64{"c01": 0.35, "c02": 0.3, "c03": 0.2, "none": 0.15}),
					"mentions_concept_c02": noul(0.71),
					"mentions_concept_c03": noul(0.69),
				})},
			wantUpdates:  map[string]any{"genre": "paper", "depth": 2.3, "concepts": []string{"[[Alpha]]", "[[Beta]]"}},
			wantComplete: true,
		},
		{
			name: "primary below threshold and none",
			in: facetInput{source: true, relevanceOn: false, gate: &GateJudgment{QualityDecided: true}, options: options, candidates: []*concept{c1},
				answers: withAnswers(map[string]decisions.Answer{
					"primary_concept":      choice("none", map[string]float64{"none": 0.6, "c01": 0.2, "c02": 0.1, "c03": 0.1}),
					"mentions_concept_c01": noul(0.2),
				})},
			wantUpdates:  map[string]any{"genre": "paper", "depth": 2.3, "concepts": []string{}},
			wantComplete: true,
			wantNone:     true,
		},
		{
			name: "primary argmax under threshold is not written",
			in: facetInput{source: true, relevanceOn: false, gate: &GateJudgment{QualityDecided: true}, options: options,
				answers: withAnswers(map[string]decisions.Answer{
					"primary_concept": choice("c02", map[string]float64{"c01": 0.26, "c02": 0.28, "c03": 0.26, "none": 0.2}),
				})},
			wantUpdates:  map[string]any{"genre": "paper", "depth": 2.3, "concepts": []string{}},
			wantComplete: true,
		},
		{
			name: "undecided concept answer writes no concepts",
			in: facetInput{source: true, relevanceOn: false, gate: &GateJudgment{QualityDecided: true}, options: options, candidates: []*concept{c2},
				answers: withAnswers(map[string]decisions.Answer{
					"primary_concept":      choice("c01", map[string]float64{"c01": 0.9, "none": 0.1}),
					"mentions_concept_c02": {Status: decisions.StatusNotChecked, Reason: "state_too_large"},
				})},
			wantUpdates:   map[string]any{"genre": "paper", "depth": 2.3},
			wantUndecided: []string{"mentions_concept_c02:not_checked:state_too_large"},
		},
		{
			name:         "article asks kind and concepts only",
			in:           facetInput{answers: map[string]decisions.Answer{"kind": choice("reference_docs", nil), "primary_concept": choice("none", map[string]float64{"none": 1})}, options: options},
			wantUpdates:  map[string]any{"genre": "reference_docs", "concepts": []string{}},
			wantComplete: true,
			wantNone:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.in.thresholds = decisions.DefaultThresholds()
			out := mapFacets(tt.in)
			if !reflect.DeepEqual(out.updates, tt.wantUpdates) {
				t.Fatalf("updates = %#v\nwant %#v", out.updates, tt.wantUpdates)
			}
			wantComplete := tt.wantComplete && len(tt.wantUndecided) == 0
			if out.complete != wantComplete {
				t.Fatalf("complete = %v, want %v", out.complete, wantComplete)
			}
			if !slices.Equal(out.undecided, tt.wantUndecided) {
				t.Fatalf("undecided = %v, want %v", out.undecided, tt.wantUndecided)
			}
			gotRecapture := ""
			if out.recapture != nil {
				gotRecapture = out.recapture.reason
			}
			if gotRecapture != tt.wantRecapture {
				t.Fatalf("recapture = %q, want %q", gotRecapture, tt.wantRecapture)
			}
			if (out.remove != nil) != tt.wantRemove {
				t.Fatalf("remove = %+v, want %v", out.remove, tt.wantRemove)
			}
			if out.primaryNone != tt.wantNone {
				t.Fatalf("primaryNone = %v, want %v", out.primaryNone, tt.wantNone)
			}
		})
	}
}
