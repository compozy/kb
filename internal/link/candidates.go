package link

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/resolve"
)

// Candidate limits (spec §9.1).
const (
	// MaxCandidates caps the merged candidate list per document.
	MaxCandidates = 20
	// LexicalTopK is the number of BM25 neighbours kept per document.
	LexicalTopK = 15
	// VectorTopK is the number of optional semantic neighbours asked for.
	VectorTopK = 10
	// topTerms is the number of frequent body terms added to the BM25 query.
	topTerms = 10
	// sharedConcepts is the concept overlap that makes an article a
	// structural candidate.
	sharedConcepts = 2
)

// Candidate origins, recorded for dry runs and tests.
const (
	OriginMention   = "mention"
	OriginStructure = "structure"
	OriginLexical   = "bm25"
	OriginVector    = "vector"
)

// Neighbours proposes semantic neighbours of a document (optional qmd vector
// search, spec §9.1 step 4). It returns topic-relative article paths best
// first. Scores never enter bands. A nil Neighbours turns the source off.
type Neighbours interface {
	Neighbours(ctx context.Context, doc *corpus.Document, k int) ([]string, error)
}

// Candidate is one article proposed as a link target for a document.
type Candidate struct {
	// ID is the element name used in question ids (`c1`, `c2`, ...).
	ID        string
	Path      string
	Title     string
	Aliases   []string
	Criterion string
	Summary   string
	Origins   []string
}

// Generator proposes link candidates by code only: alias/title mentions,
// BM25 over articles, shared structure and optional semantic neighbours.
type Generator struct {
	articles   []*corpus.Document
	byKey      map[string]*corpus.Document
	dictionary *corpus.Dictionary
	index      *corpus.BM25
	neighbours Neighbours
}

// NewGenerator indexes articles (hub files are dropped) for candidate
// generation. neighbours may be nil.
func NewGenerator(articles []*corpus.Document, neighbours Neighbours) *Generator {
	kept := make([]*corpus.Document, 0, len(articles))
	for _, article := range articles {
		if !IsHub(article.Path) {
			kept = append(kept, article)
		}
	}
	generator := &Generator{
		articles:   kept,
		byKey:      make(map[string]*corpus.Document, len(kept)*2),
		dictionary: corpus.NewDictionary(kept),
		index:      corpus.NewBM25(kept, articleFields, corpus.DefaultWeights()),
		neighbours: neighbours,
	}
	for _, article := range kept {
		generator.byKey[strings.ToLower(strings.TrimSuffix(article.Path, ".md"))] = article
		if _, exists := generator.byKey[stemKey(article.Path)]; !exists {
			generator.byKey[stemKey(article.Path)] = article
		}
	}
	return generator
}

// Articles returns the indexed (non-hub) articles.
func (g *Generator) Articles() []*corpus.Document { return slices.Clone(g.articles) }

// Candidates returns at most MaxCandidates articles for doc, never doc
// itself, merged in this order: mentions (text order), shared structure
// (existing body links to articles, articles named in `concepts`, articles
// sharing ≥2 concepts or a `sources` entry), BM25 top 15, semantic
// neighbours. IDs are assigned in merged order.
func (g *Generator) Candidates(ctx context.Context, doc *corpus.Document) []Candidate {
	merged := make([]Candidate, 0, MaxCandidates)
	position := make(map[string]int)
	add := func(article *corpus.Document, origin string) {
		if article == nil || article.Path == doc.Path {
			return
		}
		if at, ok := position[article.Path]; ok {
			if !slices.Contains(merged[at].Origins, origin) {
				merged[at].Origins = append(merged[at].Origins, origin)
			}
			return
		}
		if len(merged) >= MaxCandidates {
			return
		}
		position[article.Path] = len(merged)
		merged = append(merged, Candidate{
			Path:      article.Path,
			Title:     article.Title,
			Aliases:   slices.Clone(article.Aliases),
			Criterion: article.Criterion(),
			Summary:   article.Summary(),
			Origins:   []string{origin},
		})
	}

	for _, mention := range g.dictionary.Mentions(doc.Body) {
		add(mention.Article, OriginMention)
	}
	for _, article := range g.structural(doc) {
		add(article, OriginStructure)
	}
	for _, hit := range g.lexical(doc) {
		add(hit, OriginLexical)
	}
	if g.neighbours != nil {
		if paths, err := g.neighbours.Neighbours(ctx, doc, VectorTopK); err == nil {
			for _, p := range paths {
				add(g.lookup(p), OriginVector)
			}
		}
	}

	for index := range merged {
		merged[index].ID = candidateID(index)
	}
	return merged
}

// Mentions returns the mentions in doc's body of the given candidates, in
// text order.
func (g *Generator) Mentions(doc *corpus.Document, candidates []Candidate) []corpus.Mention {
	wanted := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		wanted[candidate.Path] = struct{}{}
	}
	mentions := make([]corpus.Mention, 0)
	for _, mention := range g.dictionary.Mentions(doc.Body) {
		if _, ok := wanted[mention.Article.Path]; ok {
			mentions = append(mentions, mention)
		}
	}
	return mentions
}

// CandidateHash is the hash of the sorted candidate paths recorded in the
// state row as `link_candidates` (spec §9.4 incremental runs).
func CandidateHash(candidates []Candidate) string {
	paths := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		paths = append(paths, candidate.Path)
	}
	sort.Strings(paths)
	sum := sha256.Sum256([]byte(strings.Join(paths, "\n")))
	return hex.EncodeToString(sum[:])[:16]
}

// IsHub reports whether a topic-relative path is a hub file that is never a
// link candidate or a linked document (`log.md`, indexes, `CLAUDE.md`, ...).
func IsHub(topicRel string) bool {
	stem := stemKey(topicRel)
	switch stem {
	case "index", "dashboard":
		return true
	case "log", "claude", "agents", "readme":
		// Topic-root bookkeeping files; an article may be named "Agents".
		return !strings.Contains(strings.Trim(topicRel, "/"), "/")
	}
	return strings.HasSuffix(stem, " index") || strings.HasSuffix(stem, "-index") || strings.HasSuffix(stem, "_index")
}

func (g *Generator) lookup(target string) *corpus.Document {
	key := strings.ToLower(resolve.Normalize(target))
	if key == "" {
		return nil
	}
	if article, ok := g.byKey[key]; ok {
		return article
	}
	if !strings.Contains(key, "/") {
		return nil
	}
	// A vault-relative path (`<topic>/wiki/concepts/X`) or a longer path:
	// match on the topic-relative suffix.
	for _, article := range g.articles {
		articleKey := strings.ToLower(strings.TrimSuffix(article.Path, ".md"))
		if strings.HasSuffix(key, "/"+articleKey) {
			return article
		}
	}
	return nil
}

func (g *Generator) structural(doc *corpus.Document) []*corpus.Document {
	found := make(map[string]*corpus.Document)
	for _, link := range resolve.BodyLinks(doc.Body) {
		if article := g.lookup(link.Target); article != nil {
			found[article.Path] = article
		}
	}

	docConcepts := keySet(doc.Concepts())
	for key := range docConcepts {
		if article := g.lookup(key); article != nil {
			found[article.Path] = article
		}
	}
	docSources := keySet(linkTargets(doc.Frontmatter, "sources"))
	selfKeys := []string{strings.ToLower(strings.TrimSuffix(doc.Path, ".md")), stemKey(doc.Path)}

	for _, article := range g.articles {
		if overlap(docConcepts, keySet(article.Concepts())) >= sharedConcepts {
			found[article.Path] = article
			continue
		}
		articleSources := keySet(linkTargets(article.Frontmatter, "sources"))
		if overlap(docSources, articleSources) > 0 {
			found[article.Path] = article
			continue
		}
		if doc.Kind == corpus.KindSource {
			for _, key := range selfKeys {
				if _, ok := articleSources[key]; ok {
					found[article.Path] = article
					break
				}
			}
		}
	}

	paths := make([]string, 0, len(found))
	for p := range found {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	articles := make([]*corpus.Document, 0, len(paths))
	for _, p := range paths {
		articles = append(articles, found[p])
	}
	return articles
}

func (g *Generator) lexical(doc *corpus.Document) []*corpus.Document {
	query := strings.Join(append([]string{doc.Title, doc.Summary(), strings.Join(doc.Entities(), " ")}, frequentTerms(doc.Body, topTerms)...), " ")
	hits := g.index.Search(query, LexicalTopK+1)
	articles := make([]*corpus.Document, 0, len(hits))
	for _, hit := range hits {
		if hit.Doc.Path == doc.Path {
			continue
		}
		articles = append(articles, hit.Doc)
		if len(articles) == LexicalTopK {
			break
		}
	}
	return articles
}

// articleFields are the BM25 fields of §9.1: title, aliases, criterion,
// summary, entities and the first 2k characters of the body.
func articleFields(d *corpus.Document) map[string]string {
	fields := corpus.DefaultFields(d)
	delete(fields, "questions")
	return fields
}

func frequentTerms(body string, n int) []string {
	counts := make(map[string]int)
	for _, token := range corpus.Tokenize(corpus.Head(body, 20000)) {
		counts[token]++
	}
	terms := make([]string, 0, len(counts))
	for term := range counts {
		terms = append(terms, term)
	}
	sort.Slice(terms, func(i, j int) bool {
		if counts[terms[i]] != counts[terms[j]] {
			return counts[terms[i]] > counts[terms[j]]
		}
		return terms[i] < terms[j]
	})
	if len(terms) > n {
		terms = terms[:n]
	}
	return terms
}

func linkTargets(values map[string]any, key string) []string {
	targets := make([]string, 0)
	for _, item := range stringList(values, key) {
		if target := resolve.Normalize(item); target != "" {
			targets = append(targets, target)
		}
	}
	return targets
}

// keySet lowercases targets into comparison keys: the full normalized target
// and, for paths, also the stem.
func keySet(targets []string) map[string]struct{} {
	set := make(map[string]struct{}, len(targets)*2)
	for _, target := range targets {
		lowered := strings.ToLower(resolve.Normalize(target))
		if lowered == "" {
			continue
		}
		set[lowered] = struct{}{}
		set[path.Base(lowered)] = struct{}{}
	}
	return set
}

func overlap(left, right map[string]struct{}) int {
	count := 0
	for key := range left {
		if strings.Contains(key, "/") {
			continue // count stems only, so a path and its stem are one entry
		}
		if _, ok := right[key]; ok {
			count++
		}
	}
	return count
}

func stringList(values map[string]any, key string) []string {
	if single, ok := values[key].(string); ok {
		return []string{single}
	}
	return frontmatter.GetStringSlice(values, key)
}

func candidateID(index int) string {
	return "c" + strconv.Itoa(index+1)
}
