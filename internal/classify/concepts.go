package classify

import (
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/compozy/kb/internal/corpus"
)

// Concept-question limits (spec §8).
const (
	// MaxConceptOptions is the number of articles offered to primary_concept;
	// with the `none` exit the choice stays within Jev's 255 options.
	MaxConceptOptions = 254
	// MaxConceptCandidates caps the mentions_concept_{id} nouls per document.
	MaxConceptCandidates = 20
	// criterionFallbackChars bounds the article head used as option text when
	// an article has neither criterion nor summary.
	criterionFallbackChars = 300
	// candidateQueryChars bounds the body head used as BM25 query.
	candidateQueryChars = 1000
)

// concept is one article of the topic vocabulary.
type concept struct {
	// id is the short option key (c01, c02, ...).
	id      string
	article *corpus.Document
	// title and aliases are snapshots taken when the vocabulary is built:
	// articles are rewritten in place by the writer while other documents
	// are being judged, so option texts never read the live document.
	title   string
	aliases []string
	// criterion is the option text: the article's criterion, else its
	// summary, else its head.
	criterion string
	// fallback reports that criterion is not a real `criterion`.
	fallback bool
	// link is the frontmatter value written to `concepts` (`[[stem]]`).
	link string
}

// optionText is the primary_concept option description: title — aliases —
// criterion (spec §5.2: the criterion, never the summary, when present).
func (c *concept) optionText() string {
	parts := []string{c.title}
	if len(c.aliases) > 0 {
		parts = append(parts, strings.Join(c.aliases, ", "))
	}
	if c.criterion != "" {
		parts = append(parts, c.criterion)
	}
	return strings.Join(parts, " — ")
}

// stateEntry is the concept as it appears in the state's `concepts` list,
// which the mentions_concept_{id} nouls refer to.
func (c *concept) stateEntry() map[string]any {
	entry := map[string]any{"id": c.id, "title": c.title, "criterion": c.criterion}
	if len(c.aliases) > 0 {
		entry["aliases"] = slices.Clone(c.aliases)
	}
	return entry
}

// vocabulary is the topic's concept vocabulary: its articles (spec §5.2).
type vocabulary struct {
	concepts    []*concept
	byPath      map[string]*concept
	index       *corpus.BM25
	dict        *corpus.Dictionary
	noCriterion []string
}

// newVocabulary builds the vocabulary over articles (sorted by path, so ids
// are stable while the article set is). linkTarget renders the `concepts`
// value of an article.
func newVocabulary(articles []*corpus.Document, linkTarget func(*corpus.Document) string) *vocabulary {
	sorted := append([]*corpus.Document(nil), articles...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Path < sorted[j].Path })
	v := &vocabulary{byPath: map[string]*concept{}}
	for index, article := range sorted {
		c := &concept{
			id:      fmt.Sprintf("c%02d", index+1),
			article: article,
			title:   article.Title,
			aliases: slices.Clone(article.Aliases),
			link:    linkTarget(article),
		}
		switch {
		case article.Criterion() != "":
			c.criterion = article.Criterion()
		case article.Summary() != "":
			c.criterion, c.fallback = article.Summary(), true
		default:
			c.criterion, c.fallback = corpus.Head(article.Body, criterionFallbackChars), true
		}
		if c.fallback {
			v.noCriterion = append(v.noCriterion, article.Path)
		}
		v.concepts = append(v.concepts, c)
		v.byPath[article.Path] = c
	}
	v.index = corpus.NewBM25(sorted, corpus.DefaultFields, corpus.DefaultWeights())
	v.dict = corpus.NewDictionary(sorted)
	return v
}

// empty reports a topic without articles (bootstrap: no concept questions).
func (v *vocabulary) empty() bool { return v == nil || len(v.concepts) == 0 }

// criterionMissing lists the articles whose option text fell back to their
// summary or head.
func (v *vocabulary) criterionMissing() []string {
	return append([]string(nil), v.noCriterion...)
}

// options returns the primary_concept options for doc: every article except
// doc itself, or, beyond MaxConceptOptions, the best BM25 matches for the
// document's query topped up in path order.
func (v *vocabulary) options(doc *corpus.Document) []*concept {
	eligible := make([]*concept, 0, len(v.concepts))
	for _, c := range v.concepts {
		if c.article.Path != doc.Path {
			eligible = append(eligible, c)
		}
	}
	if len(eligible) <= MaxConceptOptions {
		return eligible
	}
	chosen := make([]*concept, 0, MaxConceptOptions)
	seen := map[string]bool{doc.Path: true}
	for _, hit := range v.index.Search(candidateQuery(doc), 0) {
		if len(chosen) == MaxConceptOptions {
			break
		}
		if seen[hit.Doc.Path] {
			continue
		}
		seen[hit.Doc.Path] = true
		chosen = append(chosen, v.byPath[hit.Doc.Path])
	}
	for _, c := range eligible {
		if len(chosen) == MaxConceptOptions {
			break
		}
		if !seen[c.article.Path] {
			seen[c.article.Path] = true
			chosen = append(chosen, c)
		}
	}
	sort.Slice(chosen, func(i, j int) bool { return chosen[i].id < chosen[j].id })
	return chosen
}

// candidates returns the ≤ MaxConceptCandidates concepts asked
// mentions_concept_{id} for doc, built by code in one pass: first the
// articles whose title or alias occurs in the body (most occurrences first),
// then the BM25 top matches for the document's title + summary + head.
//
// The spec also lists "primary_concept distribution top 5" as a candidate
// source; that distribution is only known after the classification request,
// and classify asks one request per document, so it is not a candidate
// source here. The primary concept itself still enters `concepts` through
// its own threshold, which covers the case the distribution source exists
// for (the main subject missing from the mention/BM25 set).
func (v *vocabulary) candidates(doc *corpus.Document) []*concept {
	counts := map[string]int{}
	first := map[string]int{}
	for position, mention := range v.dict.Mentions(doc.Body) {
		p := mention.Article.Path
		if p == doc.Path {
			continue
		}
		if _, seen := first[p]; !seen {
			first[p] = position
		}
		counts[p]++
	}
	mentioned := make([]string, 0, len(counts))
	for p := range counts {
		mentioned = append(mentioned, p)
	}
	sort.Slice(mentioned, func(i, j int) bool {
		if counts[mentioned[i]] != counts[mentioned[j]] {
			return counts[mentioned[i]] > counts[mentioned[j]]
		}
		return first[mentioned[i]] < first[mentioned[j]]
	})

	out := make([]*concept, 0, MaxConceptCandidates)
	seen := map[string]bool{doc.Path: true}
	add := func(p string) {
		if len(out) >= MaxConceptCandidates || seen[p] {
			return
		}
		if c := v.byPath[p]; c != nil {
			seen[p] = true
			out = append(out, c)
		}
	}
	for _, p := range mentioned {
		add(p)
	}
	for _, hit := range v.index.Search(candidateQuery(doc), MaxConceptCandidates+1) {
		add(hit.Doc.Path)
	}
	return out
}

// candidateQuery is the BM25 query of a document: title, summary and the
// head of its body.
func candidateQuery(doc *corpus.Document) string {
	return strings.Join([]string{doc.Title, doc.Summary(), corpus.Head(doc.Body, candidateQueryChars)}, "\n")
}

// linkStem renders the `concepts` target of an article as a quoted wikilink:
// the file stem, or the vault-relative path without `.md` when the stem is
// ambiguous (plan: relation target form).
func linkStem(doc *corpus.Document, ambiguous func(stem string) bool, vaultRel string) string {
	stem := strings.TrimSuffix(path.Base(doc.Path), path.Ext(doc.Path))
	if ambiguous != nil && ambiguous(stem) {
		target := strings.TrimSuffix(doc.Path, path.Ext(doc.Path))
		if vaultRel != "" && vaultRel != "." {
			target = vaultRel + "/" + target
		}
		return "[[" + target + "]]"
	}
	return "[[" + stem + "]]"
}
