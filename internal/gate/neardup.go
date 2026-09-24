package gate

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/link"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/resolve"
	"github.com/compozy/kb/internal/session"
)

// Near-duplicate candidate limits (spec §7 stage 6).
const (
	// NearBM25TopK is the number of BM25 candidates.
	NearBM25TopK = 5
	// NearTitleSimilarity is the same-domain title similarity floor.
	NearTitleSimilarity = 0.6
	// nearMaxCandidates caps the candidates of one request.
	nearMaxCandidates = 10
	// nearHeadChars is the body prefix indexed and shown per candidate.
	nearHeadChars = 2000
	// nearExcerptTokens bounds the new document's excerpt.
	nearExcerptTokens = 4000
)

// Near-duplicate relation options of the duplicate bank.
const (
	RelationSameContent = "same_content"
	RelationNewVersion  = "new_version"
)

var duplicateBank = questions.MustLoad("duplicate")

// nearDuplicates holds the lazily built BM25 index and link targets of one
// run.
type nearDuplicates struct {
	s       *session.Session
	index   *corpus.BM25
	indexed []*corpus.Document
	targets *link.Targets
	loaded  bool
}

func newNearDuplicates(s *session.Session) *nearDuplicates {
	return &nearDuplicates{s: s}
}

func (n *nearDuplicates) reset() { n.index, n.indexed = nil, nil }

// NearCandidate is one near-duplicate candidate proposed by code.
type NearCandidate struct {
	Doc *corpus.Document
	// Origin is "bm25" or "title".
	Origin string
}

// TitleSimilarity is the token Jaccard similarity of two titles (tokens as
// corpus.Tokenize: lowercased letters/digits, stopwords dropped).
func TitleSimilarity(a, b string) float64 {
	left, right := tokenSet(a), tokenSet(b)
	if len(left) == 0 || len(right) == 0 {
		return 0
	}
	shared := 0
	for token := range left {
		if _, ok := right[token]; ok {
			shared++
		}
	}
	return float64(shared) / float64(len(left)+len(right)-shared)
}

func tokenSet(text string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, token := range corpus.Tokenize(text) {
		set[token] = struct{}{}
	}
	return set
}

// NearCandidates proposes the stage-6 candidates for doc: the BM25 top 5
// over title + first 2k characters of every source, plus sources of the
// same host whose title similarity is ≥ 0.6; the document itself excluded,
// at most 10.
func NearCandidates(doc *corpus.Document, sources []*corpus.Document, index *corpus.BM25) []NearCandidate {
	out := make([]NearCandidate, 0, nearMaxCandidates)
	seen := map[string]bool{doc.Path: true}
	if index != nil {
		for _, hit := range index.Search(doc.Title+"\n"+corpus.Head(doc.Body, nearHeadChars), NearBM25TopK+1) {
			if seen[hit.Doc.Path] || len(out) >= NearBM25TopK {
				continue
			}
			seen[hit.Doc.Path] = true
			out = append(out, NearCandidate{Doc: hit.Doc, Origin: "bm25"})
		}
	}
	host, _ := doc.Provenance()["source_host"].(string)
	if host == "" {
		return out
	}
	for _, source := range sources {
		if len(out) >= nearMaxCandidates || seen[source.Path] {
			continue
		}
		if other, _ := source.Provenance()["source_host"].(string); other != host {
			continue
		}
		if TitleSimilarity(doc.Title, source.Title) >= NearTitleSimilarity {
			seen[source.Path] = true
			out = append(out, NearCandidate{Doc: source, Origin: "title"})
		}
	}
	return out
}

func nearFields(d *corpus.Document) map[string]string {
	return map[string]string{"title": d.Title, "head": corpus.Head(d.Body, nearHeadChars)}
}

// check asks one request with a `relation_{id}` choice per candidate:
// P(same_content) ≥ duplicate → quarantine `duplicate` (in mode);
// P(new_version) ≥ duplicate → keep and record `supersedes`.
func (n *nearDuplicates) check(ctx context.Context, doc *corpus.Document, sources []*corpus.Document, out *Outcome, mode string) error {
	if len(sources) == 0 {
		return nil
	}
	if n.index == nil || len(n.indexed) != len(sources) {
		n.index = corpus.NewBM25(sources, nearFields, map[string]float64{"title": 2, "head": 1})
		n.indexed = sources
	}
	candidates := NearCandidates(doc, sources, n.index)
	if len(candidates) == 0 {
		return nil
	}

	entries := make([]map[string]any, 0, len(candidates))
	qs := make([]questions.Q, 0, len(candidates))
	for index, candidate := range candidates {
		id := "c" + strconv.Itoa(index+1)
		entries = append(entries, map[string]any{
			"id": id, "title": candidate.Doc.Title,
			"excerpt":    corpus.Head(candidate.Doc.Body, nearHeadChars),
			"provenance": candidate.Doc.Provenance(),
		})
		qs = append(qs, duplicateBank.MustQuestion("relation_{id}", map[string]string{"id": id}))
	}
	state := map[string]any{"document": corpus.DocumentState(doc, nearExcerptTokens, nil), "candidates": entries}
	result, err := n.s.Engine.Decide(ctx, n.s.Request(decisions.PurposeDuplicate, doc.Path, duplicateBank, state, qs))
	if err != nil {
		return fmt.Errorf("gate: near-duplicate %s: %w", doc.Path, err)
	}

	threshold := n.s.Threshold("duplicate")
	best, bestP := -1, 0.0
	for index := range candidates {
		answer := result.Answers["relation_c"+strconv.Itoa(index+1)]
		if !answer.Decided() {
			out.Undecided = append(out.Undecided, "relation_c"+strconv.Itoa(index+1)+":"+answer.Reason)
			continue
		}
		if p := answer.Probs[RelationSameContent]; p >= threshold && p > bestP {
			best, bestP = index, p
		}
	}
	if best >= 0 {
		match := candidates[best].Doc.Path
		out.merge(verdict{
			triage: TriageQuarantined, reason: ReasonDuplicate, stage: StageNearDuplicate, purpose: decisions.PurposeDuplicate,
			question: "relation_c" + strconv.Itoa(best+1), probability: bestP, receipt: result.ReceiptID,
			evidence: fmt.Sprintf("same content as %s (%.2f)", match, bestP),
		}, mode)
		out.DuplicateOf = match
		return nil
	}
	for index, candidate := range candidates {
		answer := result.Answers["relation_c"+strconv.Itoa(index+1)]
		if p, ok := answer.P(RelationNewVersion); ok && p >= threshold {
			out.Supersedes = append(out.Supersedes, n.wikilink(candidate.Doc.Path))
		}
	}
	return nil
}

// wikilink returns `[[<form>]]` for a topic source, using link's target form
// (stem, or the vault path when the stem is ambiguous).
func (n *nearDuplicates) wikilink(topicRel string) string {
	if !n.loaded {
		n.loaded = true
		if index, files, err := resolve.LoadVaultIndex(n.s.VaultPath, n.s.Root()); err == nil {
			prefix, relErr := filepath.Rel(filepath.Clean(n.s.VaultPath), filepath.Clean(n.s.Root()))
			if relErr != nil || strings.HasPrefix(prefix, "..") {
				prefix = filepath.Base(n.s.Root())
			}
			n.targets = link.NewTargets(index, files, filepath.ToSlash(prefix))
		}
	}
	return n.targets.Wikilink(topicRel)
}
