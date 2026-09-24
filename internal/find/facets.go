package find

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/compozy/kb/internal/corpus"
)

// Facet bucket names.
const (
	// FacetNone buckets documents without a value for the facet.
	FacetNone = "(none)"
	// TopConcepts is the number of concepts listed by Facets.
	TopConcepts = 20
)

// Count is one facet value and the number of documents carrying it.
type Count struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// FacetReport is the corpus atlas of a topic (spec §8): document counts per
// facet value, computed from frontmatter only (no model call).
type FacetReport struct {
	Documents   int     `json:"documents"`
	Sources     int     `json:"sources"`
	Articles    int     `json:"articles"`
	Quarantined int     `json:"quarantined"`
	Genre       []Count `json:"genre"`
	Relevance   []Count `json:"relevance"`
	// Depth buckets the 0–3 depth expectation by its integer part.
	Depth    []Count `json:"depth"`
	Concepts []Count `json:"concepts"`
}

// Facets counts genre, relevance, depth (buckets 0–3) and the top concepts
// over docs. Quarantined documents (`triage: quarantined`) are counted in
// Quarantined only.
func Facets(docs []*corpus.Document) FacetReport {
	report := FacetReport{}
	genre, relevance, depth, concepts := map[string]int{}, map[string]int{}, map[string]int{}, map[string]int{}
	for _, doc := range docs {
		if strings.EqualFold(doc.Triage(), "quarantined") {
			report.Quarantined++
			continue
		}
		report.Documents++
		if doc.Kind == corpus.KindArticle {
			report.Articles++
		} else {
			report.Sources++
		}
		genre[valueOrNone(doc.Genre())]++
		relevance[valueOrNone(doc.Relevance())]++
		depth[depthBucket(doc)]++
		for _, concept := range doc.Concepts() {
			concepts[concept]++
		}
	}
	report.Genre = sortedCounts(genre, 0)
	report.Relevance = sortedCounts(relevance, 0)
	report.Depth = depthCounts(depth)
	report.Concepts = sortedCounts(concepts, TopConcepts)
	return report
}

func depthBucket(doc *corpus.Document) string {
	value, ok := doc.Depth()
	if !ok || math.IsNaN(value) {
		return FacetNone
	}
	return strconv.Itoa(int(math.Floor(clamp3(value))))
}

func depthCounts(counts map[string]int) []Count {
	result := make([]Count, 0, 5)
	for _, bucket := range []string{"0", "1", "2", "3", FacetNone} {
		result = append(result, Count{Value: bucket, Count: counts[bucket]})
	}
	return result
}

func sortedCounts(counts map[string]int, limit int) []Count {
	result := make([]Count, 0, len(counts))
	for value, count := range counts {
		result = append(result, Count{Value: value, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Value < result[j].Value
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}

func valueOrNone(value string) string {
	if strings.TrimSpace(value) == "" {
		return FacetNone
	}
	return value
}
