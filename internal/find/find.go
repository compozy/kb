// Package find is kb's retrieval judge (spec §10): code proposes candidates
// (frontmatter facet filters, BM25, concept expansion, optional vectors),
// the decision model judges whether each candidate answers the question, and
// code ranks the survivors and names the reason every other candidate was
// dropped. Facets reports the corpus atlas counts without any model call.
package find

import (
	"context"
	"fmt"
	"maps"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/resolve"
	"github.com/compozy/kb/internal/session"
)

// Retrieval limits (spec §10).
const (
	// BankID is the question bank asked by find.
	BankID = "find"
	// LexicalTopK is the number of BM25 candidates.
	LexicalTopK = 60
	// VectorTopK is the number of optional vector candidates.
	VectorTopK = 30
	// ConceptSeeds is how many top BM25 hits contribute their `concepts`.
	ConceptSeeds = 10
	// BatchSize is the number of candidates judged per request.
	BatchSize = 20
	// HeadChars bounds each candidate's head in the state.
	HeadChars = 400
	// DefaultLimit is the default number of results.
	DefaultLimit = 10
)

// Drop reasons (spec §10 step 3). Undecided answers are reported as
// "undecided:<reason>".
const (
	ReasonLowProbability = "low-probability"
	ReasonBelowWindow    = "below-window"
	ReasonLimit          = "limit"
	ReasonFilteredFacet  = "filtered-facet"
	ReasonQuarantined    = "quarantined"
	reasonUndecided      = "undecided:"
)

// Question kinds of the find bank's question_kind choice.
const (
	KindDefinition = "definition"
	KindHowTo      = "how_to"
	KindComparison = "comparison"
	KindEvidence   = "evidence_or_data"
	KindOpinion    = "opinion"
	KindOther      = "other"
)

// KindGenres maps a question kind to the document genres that match it
// (kind_match = 1 in the rank formula). `other` matches nothing.
//
//	definition        reference_docs, article_or_essay
//	how_to            tutorial_or_guide
//	comparison        article_or_essay, paper
//	evidence_or_data  paper, dataset_or_benchmark
//	opinion           opinion_or_discussion
var KindGenres = map[string][]string{
	KindDefinition: {"reference_docs", "article_or_essay"},
	KindHowTo:      {"tutorial_or_guide"},
	KindComparison: {"article_or_essay", "paper"},
	KindEvidence:   {"paper", "dataset_or_benchmark"},
	KindOpinion:    {"opinion_or_discussion"},
}

// Vectors proposes candidates by vector similarity (optional qmd, spec §10
// step 1). It returns topic-relative paths best first; nil turns it off.
type Vectors interface {
	Search(ctx context.Context, query string, k int) ([]string, error)
}

// Query is one find request.
type Query struct {
	Text  string
	Limit int
	// Facet filters, applied to frontmatter before any model call.
	Kind      string
	Concept   string
	Relevance string
	MinDepth  float64
	Explain   bool
	// Vectors is the optional vector candidate source (nil = off).
	Vectors Vectors `json:"-"`
}

// Hit is one kept document.
type Hit struct {
	Path        string   `json:"path"`
	Title       string   `json:"title"`
	Rank        float64  `json:"rank"`
	PAnswers    float64  `json:"p_answers"`
	Specificity float64  `json:"specificity"`
	Kind        string   `json:"kind,omitempty"`
	Concepts    []string `json:"concepts,omitempty"`
	Depth       *float64 `json:"depth,omitempty"`
}

// Drop is one candidate that was not returned, with its named reason.
type Drop struct {
	Path     string   `json:"path"`
	Title    string   `json:"title"`
	Reason   string   `json:"reason"`
	PAnswers *float64 `json:"p_answers,omitempty"`
	Rank     *float64 `json:"rank,omitempty"`
}

// Result is the answer to one query.
type Result struct {
	Question     string `json:"question"`
	QuestionKind string `json:"question_kind,omitempty"`
	Candidates   int    `json:"candidates"`
	Hits         []Hit  `json:"results"`
	// Dropped lists every candidate not returned (filled only with Explain).
	Dropped []Drop `json:"dropped,omitempty"`
}

// Judged is a candidate with its answers.
type Judged struct {
	Doc *corpus.Document
	// PAnswers is P(yes) of answers_{id}; valid when Decided.
	PAnswers float64
	// Specificity is the expected level (0–3) of specificity_{id}; 0 when
	// undecided.
	Specificity float64
	Decided     bool
	// Reason is "undecided:<reason>" when not Decided.
	Reason string
}

// Run answers q over the session topic.
func Run(ctx context.Context, s *session.Session, q Query) (Result, error) {
	text := strings.TrimSpace(q.Text)
	if text == "" {
		return Result{}, fmt.Errorf("find: a question is required")
	}
	if q.Limit <= 0 {
		q.Limit = DefaultLimit
	}
	bank, err := questions.Load(BankID)
	if err != nil {
		return Result{}, fmt.Errorf("find: %w", err)
	}
	c, err := s.Corpus()
	if err != nil {
		return Result{}, fmt.Errorf("find: %w", err)
	}

	candidates, dropped := Candidates(ctx, c.Documents(), q)
	result := Result{Question: text, Candidates: len(candidates)}

	judgedAll, kindProbs, err := judge(ctx, s, bank, text, candidates)
	if err != nil {
		return Result{}, err
	}
	result.QuestionKind = argmax(kindProbs)
	hits, rankDrops := Rank(judgedAll, result.QuestionKind, s.Threshold("find_keep"), s.Threshold("find_window"), q.Limit)
	result.Hits = hits
	if q.Explain {
		result.Dropped = append(dropped, rankDrops...)
	}
	return result, nil
}

// Candidates proposes the documents to judge, by code only: facet filters
// first, then BM25 (title, aliases, summary, questions, entities, head) top
// 60, the articles named in `concepts` of the top 10 BM25 hits, and optional
// vector hits. Quarantined (`triage: quarantined`) and facet-filtered
// documents never consume a candidate slot; those that ranked among the BM25
// hits are returned as drops with their reason.
func Candidates(ctx context.Context, docs []*corpus.Document, q Query) ([]*corpus.Document, []Drop) {
	index := corpus.NewBM25(docs, retrievalFields, corpus.DefaultWeights())
	byKey := make(map[string]*corpus.Document, len(docs)*2)
	for _, doc := range docs {
		key := strings.ToLower(strings.TrimSuffix(doc.Path, ".md"))
		byKey[key] = doc
		stem := pathStem(key)
		if _, exists := byKey[stem]; !exists {
			byKey[stem] = doc
		}
	}

	candidates := make([]*corpus.Document, 0, LexicalTopK)
	seen := make(map[string]bool)
	dropped := make([]Drop, 0)
	lexical := 0
	add := func(doc *corpus.Document) bool {
		if doc == nil || seen[doc.Path] {
			return false
		}
		seen[doc.Path] = true
		if reason := exclusion(doc, q); reason != "" {
			dropped = append(dropped, Drop{Path: doc.Path, Title: doc.Title, Reason: reason})
			return false
		}
		candidates = append(candidates, doc)
		return true
	}

	hits := index.Search(q.Text, 0)
	top := make([]*corpus.Document, 0, ConceptSeeds)
	for _, hit := range hits {
		if lexical >= LexicalTopK {
			break
		}
		if add(hit.Doc) {
			lexical++
			if len(top) < ConceptSeeds {
				top = append(top, hit.Doc)
			}
		}
	}
	for _, doc := range top {
		for _, concept := range doc.Concepts() {
			key := strings.ToLower(resolve.Normalize(concept))
			target, ok := byKey[key]
			if !ok {
				target = byKey[pathStem(key)]
			}
			add(target)
		}
	}
	if q.Vectors != nil {
		if paths, err := q.Vectors.Search(ctx, q.Text, VectorTopK); err == nil {
			for _, p := range paths {
				add(byKey[strings.ToLower(strings.TrimSuffix(resolve.Normalize(p), ".md"))])
			}
		}
	}
	return candidates, dropped
}

// exclusion returns the drop reason of a document the query must not judge.
func exclusion(doc *corpus.Document, q Query) string {
	if strings.EqualFold(doc.Triage(), "quarantined") || strings.HasPrefix(doc.Path, resolve.QuarantineDir+"/") {
		return ReasonQuarantined
	}
	if q.Kind != "" && !strings.EqualFold(doc.Genre(), strings.TrimSpace(q.Kind)) {
		return ReasonFilteredFacet
	}
	if q.Relevance != "" && !strings.EqualFold(doc.Relevance(), strings.TrimSpace(q.Relevance)) {
		return ReasonFilteredFacet
	}
	if q.MinDepth > 0 {
		if depth, ok := doc.Depth(); !ok || depth < q.MinDepth {
			return ReasonFilteredFacet
		}
	}
	if q.Concept != "" && !hasConcept(doc, q.Concept) {
		return ReasonFilteredFacet
	}
	return ""
}

func hasConcept(doc *corpus.Document, concept string) bool {
	want := strings.ToLower(resolve.Normalize(concept))
	if want == "" {
		return true
	}
	for _, target := range doc.Concepts() {
		lowered := strings.ToLower(target)
		if lowered == want || pathStem(lowered) == pathStem(want) {
			return true
		}
	}
	return false
}

// judge asks the find bank over candidates in requests of BatchSize. It
// returns every candidate with its answers and the summed question_kind
// distribution across requests.
func judge(ctx context.Context, s *session.Session, bank *questions.Bank, question string, candidates []*corpus.Document) ([]Judged, map[string]float64, error) {
	all := make([]Judged, 0, len(candidates))
	kinds := map[string]float64{}
	for start := 0; start < len(candidates); start += BatchSize {
		batch := candidates[start:min(start+BatchSize, len(candidates))]
		state := map[string]any{"question": question}
		entries := make([]map[string]any, 0, len(batch))
		qs := make([]questions.Q, 0, len(batch)*2+1)
		for index, doc := range batch {
			id := "d" + strconv.Itoa(index+1)
			entry := map[string]any{"id": id, "title": doc.Title, "head": corpus.Head(doc.Body, HeadChars)}
			if genre := doc.Genre(); genre != "" {
				entry["kind"] = genre
			}
			if concepts := doc.Concepts(); len(concepts) > 0 {
				entry["concepts"] = concepts
			}
			if summary := doc.Summary(); summary != "" {
				entry["summary"] = summary
			}
			entries = append(entries, entry)
			for _, template := range []string{"answers_{id}", "specificity_{id}"} {
				q, err := bank.Question(template, map[string]string{"id": id})
				if err != nil {
					return nil, nil, fmt.Errorf("find: %w", err)
				}
				qs = append(qs, q)
			}
		}
		kindQ, err := bank.Question("question_kind", nil)
		if err != nil {
			return nil, nil, fmt.Errorf("find: %w", err)
		}
		qs = append(qs, kindQ)
		state["candidates"] = entries

		subject := "find:" + question
		res, err := s.Engine.Decide(ctx, s.Request(decisions.PurposeFind, subject, bank, state, qs))
		if err != nil {
			return nil, nil, fmt.Errorf("find: %w", err)
		}
		if kind := res.Answers["question_kind"]; kind.Decided() {
			for option, p := range kind.Probs {
				kinds[option] += p
			}
		}
		for index, doc := range batch {
			id := "d" + strconv.Itoa(index+1)
			answer := res.Answers["answers_"+id]
			item := Judged{Doc: doc}
			if p, ok := answer.P(""); ok {
				item.PAnswers, item.Decided = p, true
			} else {
				item.Reason = reasonUndecided + undecidedReason(answer)
			}
			if spec := res.Answers["specificity_"+id]; spec.Decided() && spec.Score != nil {
				item.Specificity = *spec.Score
			}
			all = append(all, item)
		}
	}
	return all, kinds, nil
}

func undecidedReason(answer decisions.Answer) string {
	if answer.Reason != "" {
		return answer.Reason
	}
	if answer.Status != "" {
		return string(answer.Status)
	}
	return decisions.ReasonMissingAnswer
}

// Score computes the rank of spec §10 step 3:
//
//	rank = p_answers × (1 + 0.25·specificity/3) × (1 + 0.2·kind_match) × (1 + 0.1·depth/3)
//
// specificity and depth are 0–3 expected levels; kindMatch is 0 or 1.
func Score(pAnswers, specificity float64, kindMatch bool, depth float64) float64 {
	match := 0.0
	if kindMatch {
		match = 1
	}
	return pAnswers * (1 + 0.25*clamp3(specificity)/3) * (1 + 0.2*match) * (1 + 0.1*clamp3(depth)/3)
}

// KindMatches reports whether a document genre matches a question kind
// (see KindGenres).
func KindMatches(questionKind, genre string) bool {
	return genre != "" && slices.Contains(KindGenres[questionKind], strings.ToLower(strings.TrimSpace(genre)))
}

// Rank keeps judged candidates with p_answers ≥ keep and within window of
// the best p_answers, orders them by Score and returns the top limit. Every
// other candidate becomes a Drop with its reason.
func Rank(items []Judged, questionKind string, keep, window float64, limit int) ([]Hit, []Drop) {
	best := 0.0
	for _, item := range items {
		if item.Decided {
			best = math.Max(best, item.PAnswers)
		}
	}
	hits := make([]Hit, 0)
	drops := make([]Drop, 0)
	for _, item := range items {
		doc := item.Doc
		if !item.Decided {
			drops = append(drops, Drop{Path: doc.Path, Title: doc.Title, Reason: item.Reason})
			continue
		}
		depth, hasDepth := doc.Depth()
		rank := Score(item.PAnswers, item.Specificity, KindMatches(questionKind, doc.Genre()), depth)
		p := item.PAnswers
		switch {
		case item.PAnswers < keep:
			drops = append(drops, Drop{Path: doc.Path, Title: doc.Title, Reason: ReasonLowProbability, PAnswers: &p, Rank: &rank})
			continue
		case item.PAnswers < best-window:
			drops = append(drops, Drop{Path: doc.Path, Title: doc.Title, Reason: ReasonBelowWindow, PAnswers: &p, Rank: &rank})
			continue
		}
		hit := Hit{
			Path: doc.Path, Title: doc.Title, Rank: rank, PAnswers: item.PAnswers, Specificity: item.Specificity,
			Kind: doc.Genre(), Concepts: doc.Concepts(),
		}
		if hasDepth {
			hit.Depth = &depth
		}
		hits = append(hits, hit)
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].Rank != hits[j].Rank {
			return hits[i].Rank > hits[j].Rank
		}
		return hits[i].Path < hits[j].Path
	})
	if limit > 0 && len(hits) > limit {
		for _, hit := range hits[limit:] {
			p, rank := hit.PAnswers, hit.Rank
			drops = append(drops, Drop{Path: hit.Path, Title: hit.Title, Reason: ReasonLimit, PAnswers: &p, Rank: &rank})
		}
		hits = hits[:limit]
	}
	return hits, drops
}

// retrievalFields are the BM25 fields of §10: title, aliases, summary,
// questions, entities and the head of the body.
func retrievalFields(d *corpus.Document) map[string]string {
	fields := corpus.DefaultFields(d)
	delete(fields, "criterion")
	return fields
}

func argmax(probs map[string]float64) string {
	best, bestP := "", -1.0
	for _, option := range slices.Sorted(maps.Keys(probs)) {
		if probs[option] > bestP {
			best, bestP = option, probs[option]
		}
	}
	return best
}

func clamp3(value float64) float64 {
	return math.Max(0, math.Min(3, value))
}

func pathStem(key string) string {
	if index := strings.LastIndex(key, "/"); index >= 0 {
		return key[index+1:]
	}
	return key
}
