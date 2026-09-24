package classify

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/compozy/kb/internal/corpus"
)

// Report summarizes one classify run. Counts are per document unless named
// otherwise.
type Report struct {
	// Documents is the number of documents considered (after Paths).
	Documents int `json:"documents"`
	// Judged is the number of documents that went through the facet pass.
	Judged int `json:"judged"`
	// Decided is the number of judged documents whose every answer and
	// literal was decided (the rest stay unclassified for the next run).
	Decided int `json:"decided"`
	// UndecidedDocuments lists, sorted, the documents left with an
	// undecided answer or literal (spec §2.3: named in the run summary).
	UndecidedDocuments []string `json:"undecided_documents"`
	// Skipped counts documents not judged, by reason: unchanged (state row
	// current), complete (--only-missing and nothing missing), locked, stub.
	Skipped map[string]int `json:"skipped"`
	// Facets counts written keys (a key rewritten with the same value is not
	// counted).
	Facets map[string]int `json:"facets"`
	// Writes counts document write outcomes: written, unchanged,
	// skipped:locked, skipped:changed.
	Writes map[string]int `json:"writes"`
	// SkippedUserKeys counts keys left alone because the user owns them
	// (skipped:user-key), by key.
	SkippedUserKeys map[string]int `json:"skipped_user_keys"`
	// Undecided counts answers and literals that were not decided, by
	// "<question>:<reason>" (not-checked answers read
	// "<question>:not_checked:<reason>"). Undecided answers never write.
	Undecided map[string]int `json:"undecided"`
	// Queued counts review items added, by queue.
	Queued map[string]int `json:"queued"`
	// DroppedEntities counts generated entity names dropped by the verbatim
	// check or the cap.
	DroppedEntities int `json:"dropped_entities"`
	// InvalidLiterals counts generated literals rejected by code, by key.
	InvalidLiterals map[string]int `json:"invalid_literals"`
	// CriterionMissing lists articles without `criterion` whose option text
	// fell back to their summary or head.
	CriterionMissing []string `json:"criterion_missing"`
	// ArticlesChanged lists articles that gained aliases (they trigger the
	// link reverse pass).
	ArticlesChanged []string `json:"articles_changed"`
	// PrimaryNone counts sources whose primary concept was `none` in this
	// run.
	PrimaryNone int `json:"primary_none"`
	// PrimaryNoneTotal counts the topic's sources whose primary concept is
	// `none` (this run plus earlier runs), the concept proposal trigger.
	PrimaryNoneTotal int `json:"primary_none_total"`
	// NoVocabulary reports a topic without articles: no concept questions
	// were asked.
	NoVocabulary bool `json:"no_vocabulary"`
	// RelevanceOff reports decisions.relevance: off (no role questions).
	RelevanceOff bool `json:"relevance_off"`
}

func newReport() Report {
	return Report{
		Skipped:            map[string]int{},
		Facets:             map[string]int{},
		Writes:             map[string]int{},
		SkippedUserKeys:    map[string]int{},
		Undecided:          map[string]int{},
		Queued:             map[string]int{},
		InvalidLiterals:    map[string]int{},
		CriterionMissing:   []string{},
		ArticlesChanged:    []string{},
		UndecidedDocuments: []string{},
	}
}

// maxListedDocuments caps the undecided documents named in Lines.
const maxListedDocuments = 10

// UndecidedTotal is the number of undecided answers and literals.
func (r Report) UndecidedTotal() int {
	total := 0
	for _, count := range r.Undecided {
		total += count
	}
	return total
}

// Lines renders the report for humans.
func (r Report) Lines() []string {
	lines := []string{
		fmt.Sprintf("classify: %d documents, %d judged%s", r.Documents, r.Judged, countsSuffix(r.Skipped, "skipped")),
		fmt.Sprintf("coverage: %d/%d judged documents fully decided", r.Decided, r.Judged),
	}
	if len(r.UndecidedDocuments) > 0 {
		lines = append(lines, fmt.Sprintf("undecided documents (%d, left unclassified): %s", len(r.UndecidedDocuments), corpus.ListPaths(r.UndecidedDocuments, maxListedDocuments)))
	}
	if len(r.Facets) > 0 {
		lines = append(lines, "written: "+joinCounts(r.Facets))
	}
	if len(r.Writes) > 0 {
		lines = append(lines, "writes: "+joinCounts(r.Writes))
	}
	if len(r.SkippedUserKeys) > 0 {
		lines = append(lines, "skipped:user-key: "+joinCounts(r.SkippedUserKeys)+" (the user owns these keys; lint reports key-conflict)")
	}
	if total := r.UndecidedTotal(); total > 0 {
		lines = append(lines, fmt.Sprintf("undecided: %d (%s); nothing was written for them, re-run to retry", total, joinCounts(r.Undecided)))
	}
	if len(r.Queued) > 0 {
		lines = append(lines, "review queue: "+joinCounts(r.Queued)+" added (see `kb review`)")
	}
	if r.DroppedEntities > 0 {
		lines = append(lines, fmt.Sprintf("entities: %d generated names dropped (not found verbatim in the body)", r.DroppedEntities))
	}
	if len(r.InvalidLiterals) > 0 {
		lines = append(lines, "invalid generated literals: "+joinCounts(r.InvalidLiterals))
	}
	if len(r.CriterionMissing) > 0 {
		lines = append(lines, fmt.Sprintf("criterion-missing: %d articles use their summary or opening as concept text: %s", len(r.CriterionMissing), strings.Join(r.CriterionMissing, ", ")))
	}
	if len(r.ArticlesChanged) > 0 {
		lines = append(lines, fmt.Sprintf("aliases: %d articles gained aliases: %s", len(r.ArticlesChanged), strings.Join(r.ArticlesChanged, ", ")))
	}
	if r.RelevanceOff {
		lines = append(lines, "relevance: off in topic.yaml (no role questions)")
	}
	if r.NoVocabulary {
		lines = append(lines, "no vocabulary: the topic has no articles, so no concepts were judged")
	}
	return lines
}

// Rows renders the report as metric/value rows for table output.
func (r Report) Rows() []map[string]any {
	rows := []map[string]any{
		{"metric": "documents", "value": r.Documents},
		{"metric": "judged", "value": r.Judged},
		{"metric": "decided", "value": r.Decided},
		{"metric": "undecided_documents", "value": len(r.UndecidedDocuments)},
	}
	add := func(prefix string, counts map[string]int) {
		for _, key := range slices.Sorted(maps.Keys(counts)) {
			rows = append(rows, map[string]any{"metric": prefix + key, "value": counts[key]})
		}
	}
	add("skipped:", r.Skipped)
	add("written:", r.Facets)
	add("writes:", r.Writes)
	add("skipped:user-key:", r.SkippedUserKeys)
	add("undecided:", r.Undecided)
	add("queued:", r.Queued)
	add("invalid:", r.InvalidLiterals)
	if r.DroppedEntities > 0 {
		rows = append(rows, map[string]any{"metric": "dropped_entities", "value": r.DroppedEntities})
	}
	if len(r.CriterionMissing) > 0 {
		rows = append(rows, map[string]any{"metric": "criterion_missing", "value": len(r.CriterionMissing)})
	}
	if len(r.ArticlesChanged) > 0 {
		rows = append(rows, map[string]any{"metric": "articles_changed", "value": len(r.ArticlesChanged)})
	}
	if r.NoVocabulary {
		rows = append(rows, map[string]any{"metric": "no_vocabulary", "value": true})
	}
	return rows
}

func joinCounts(counts map[string]int) string {
	parts := make([]string, 0, len(counts))
	for _, key := range slices.Sorted(maps.Keys(counts)) {
		parts = append(parts, fmt.Sprintf("%s %d", key, counts[key]))
	}
	return strings.Join(parts, ", ")
}

func countsSuffix(counts map[string]int, label string) string {
	if len(counts) == 0 {
		return ""
	}
	return ", " + label + " (" + joinCounts(counts) + ")"
}
