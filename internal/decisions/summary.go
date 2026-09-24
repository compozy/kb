package decisions

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
)

// Summary is the run summary every command prints: calls, cache hits, cost,
// and per-purpose outcomes. Engine.Summary fills everything but Generation;
// callers that also generate set Generation from generation.Client.Summary
// before printing Lines.
type Summary struct {
	// Calls counts decision requests that reached the network.
	Calls int
	// CacheHits counts decision requests served from receipts.
	CacheHits int
	// CostUSD sums the reported cost of decision calls.
	CostUSD float64
	// CostUnknown counts decision calls with no reported cost (failed calls
	// included); their spend is never assumed to be zero.
	CostUnknown int
	// Purposes holds per-purpose answer counts.
	Purposes map[Purpose]PurposeSummary
	// Generation is the generation client's part of the run.
	Generation GenerationSummary
}

// PurposeSummary counts the answers of one purpose.
type PurposeSummary struct {
	Decided    int
	Bands      map[Band]int
	Undecided  map[string]int // by reason
	NotChecked map[string]int // by reason
}

// GenerationSummary is the generation client's contribution to a run.
type GenerationSummary struct {
	Calls       int
	CacheHits   int
	Fallbacks   int
	Failures    int
	CostUSD     float64
	CostUnknown int
}

// Lines renders the summary as printable lines, for example
//
//	decisions: 12 calls, 3 cache hits, US$ 0.0041
//	relevance: 10 decided (apply 2, review 3, ignore 5), 1 undecided:timeout
func (s Summary) Lines() []string {
	head := fmt.Sprintf("decisions: %s, %s, US$ %.4f", plural(s.Calls, "call"), plural(s.CacheHits, "cache hit"), s.CostUSD)
	if s.CostUnknown > 0 {
		head += fmt.Sprintf(", %s with unknown cost", plural(s.CostUnknown, "call"))
	}
	lines := []string{head}
	for _, purpose := range slices.Sorted(maps.Keys(s.Purposes)) {
		lines = append(lines, string(purpose)+": "+s.Purposes[purpose].line())
	}
	if g := s.Generation; g.Calls > 0 || g.CacheHits > 0 || g.Failures > 0 {
		line := fmt.Sprintf("generation: %s, %s, US$ %.4f, %s, %s",
			plural(g.Calls, "call"), plural(g.CacheHits, "cache hit"), g.CostUSD,
			plural(g.Fallbacks, "fallback"), plural(g.Failures, "failure"))
		if g.CostUnknown > 0 {
			line += fmt.Sprintf(", %s with unknown cost", plural(g.CostUnknown, "call"))
		}
		lines = append(lines, line)
	}
	return lines
}

func (p PurposeSummary) line() string {
	parts := []string{fmt.Sprintf("%d decided", p.Decided)}
	if p.Decided > 0 {
		parts[0] += fmt.Sprintf(" (apply %d, review %d, ignore %d)", p.Bands[BandApply], p.Bands[BandReview], p.Bands[BandIgnore])
	}
	for _, reason := range slices.Sorted(maps.Keys(p.Undecided)) {
		parts = append(parts, fmt.Sprintf("%d undecided:%s", p.Undecided[reason], reason))
	}
	for _, reason := range slices.Sorted(maps.Keys(p.NotChecked)) {
		parts = append(parts, fmt.Sprintf("%d not_checked:%s", p.NotChecked[reason], reason))
	}
	return strings.Join(parts, ", ")
}

func plural(count int, noun string) string {
	if count == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", count, noun)
}

// stats accumulates the engine's run summary.
type stats struct {
	mu      sync.Mutex
	summary Summary
}

func (s *stats) call(cost *float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.summary.Calls++
	if cost == nil {
		s.summary.CostUnknown++
		return
	}
	s.summary.CostUSD += *cost
}

func (s *stats) cacheHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.summary.CacheHits++
}

func (s *stats) answers(purpose Purpose, answers map[string]Answer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.summary.Purposes == nil {
		s.summary.Purposes = map[Purpose]PurposeSummary{}
	}
	entry := s.summary.Purposes[purpose]
	for _, answer := range answers {
		switch answer.Status {
		case StatusDecided:
			entry.Decided++
			if entry.Bands == nil {
				entry.Bands = map[Band]int{}
			}
			entry.Bands[answer.Band]++
		case StatusNotChecked:
			if entry.NotChecked == nil {
				entry.NotChecked = map[string]int{}
			}
			entry.NotChecked[answer.Reason]++
		default:
			if entry.Undecided == nil {
				entry.Undecided = map[string]int{}
			}
			entry.Undecided[answer.Reason]++
		}
	}
	s.summary.Purposes[purpose] = entry
}

func (s *stats) snapshot() Summary {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.summary
	out.Purposes = make(map[Purpose]PurposeSummary, len(s.summary.Purposes))
	for purpose, entry := range s.summary.Purposes {
		out.Purposes[purpose] = PurposeSummary{
			Decided:    entry.Decided,
			Bands:      maps.Clone(entry.Bands),
			Undecided:  maps.Clone(entry.Undecided),
			NotChecked: maps.Clone(entry.NotChecked),
		}
	}
	return out
}
