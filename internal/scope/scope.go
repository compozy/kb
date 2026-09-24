// Package scope orchestrates the selection-contract workflow of `kb topic
// contract` (spec §5.1): drafting a contract with the generation model,
// importing one from the topic CLAUDE.md, the self-check of kept lines
// against out-of-scope lines, the impact preview (a shadow gate run of the
// draft over the topic's own sources) and activation after an explicit
// confirmation.
//
// internal/contract stays free of the decision stack (it only parses,
// validates, hashes and writes); everything that needs a session, the gate
// judgment or the review ledgers lives here.
package scope

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"sync"

	"github.com/compozy/kb/internal/corpus"
)

// RootFolder names sources stored directly under raw/ (no subfolder).
const RootFolder = "raw"

// CostPerDocumentUSD is the measured decision cost of one gate judgment
// (spec §5.1: about US$ 0.0003 per document), used for estimates only.
const CostPerDocumentUSD = 0.0003

// Folder returns the `raw/<subfolder>` stratum of a topic-relative source
// path, or RootFolder for a file directly under raw/.
func Folder(topicRel string) string {
	rest, ok := strings.CutPrefix(topicRel, "raw/")
	if !ok {
		return RootFolder
	}
	first, _, nested := strings.Cut(rest, "/")
	if !nested {
		return RootFolder
	}
	return "raw/" + first
}

// Stratify returns a deterministic sample of at most n documents, stratified
// by raw/ subfolder: every folder gets a quota proportional to its size (at
// least one per folder while n allows it, largest remainders first, ties by
// folder name) and each quota is filled with evenly spaced documents in path
// order. The result is sorted by path. When len(docs) ≤ n (or n ≤ 0) every
// document is returned.
func Stratify(docs []*corpus.Document, n int) []*corpus.Document {
	sorted := slices.Clone(docs)
	slices.SortFunc(sorted, func(a, b *corpus.Document) int { return cmp.Compare(a.Path, b.Path) })
	if n <= 0 || len(sorted) <= n {
		return sorted
	}

	groups := map[string][]*corpus.Document{}
	for _, doc := range sorted {
		folder := Folder(doc.Path)
		groups[folder] = append(groups[folder], doc)
	}
	names := make([]string, 0, len(groups))
	sizes := make(map[string]int, len(groups))
	for name, members := range groups {
		names = append(names, name)
		sizes[name] = len(members)
	}
	slices.Sort(names)
	quotas := allocate(names, sizes, len(sorted), n)

	sample := make([]*corpus.Document, 0, n)
	for _, name := range names {
		sample = append(sample, evenly(groups[name], quotas[name])...)
	}
	slices.SortFunc(sample, func(a, b *corpus.Document) int { return cmp.Compare(a.Path, b.Path) })
	return sample
}

// allocate splits n over the named groups proportionally to their sizes
// (largest remainder method), giving every group at least one slot while n
// covers all groups and never more than its size.
func allocate(names []string, sizes map[string]int, total, n int) map[string]int {
	quotas := make(map[string]int, len(names))
	remainders := make(map[string]float64, len(names))
	sum := 0
	for _, name := range names {
		exact := float64(n) * float64(sizes[name]) / float64(total)
		quota := min(int(exact), sizes[name])
		if quota == 0 && n >= len(names) {
			quota = 1
		}
		quotas[name] = quota
		remainders[name] = exact - float64(int(exact))
		sum += quota
	}
	for sum < n {
		best := ""
		for _, name := range names {
			if quotas[name] >= sizes[name] {
				continue
			}
			if best == "" || remainders[name] > remainders[best] {
				best = name
			}
		}
		if best == "" {
			break
		}
		quotas[best]++
		remainders[best] = -1
		sum++
	}
	for sum > n {
		// Only reachable through the one-per-group floor: take slots back from
		// the group with the largest quota (first by name on ties).
		best := ""
		for _, name := range names {
			if quotas[name] > 1 && (best == "" || quotas[name] > quotas[best]) {
				best = name
			}
		}
		if best == "" {
			break
		}
		quotas[best]--
		sum--
	}
	return quotas
}

// evenly picks q evenly spaced members (the middle of each of q equal
// slices), preserving order.
func evenly(members []*corpus.Document, q int) []*corpus.Document {
	if q >= len(members) {
		return members
	}
	picked := make([]*corpus.Document, 0, q)
	for i := range q {
		picked = append(picked, members[(2*i+1)*len(members)/(2*q)])
	}
	return picked
}

// forEach runs fn over items with at most workers goroutines and returns
// the first error (remaining items are skipped once one fails).
func forEach[T any](ctx context.Context, workers int, items []T, fn func(context.Context, T) error) error {
	if len(items) == 0 {
		return nil
	}
	workers = max(1, min(workers, len(items)))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	next := make(chan T)
	var (
		wg       sync.WaitGroup
		once     sync.Once
		firstErr error
	)
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for item := range next {
				if err := fn(ctx, item); err != nil {
					once.Do(func() {
						firstErr = err
						cancel()
					})
				}
			}
		}()
	}
feed:
	for _, item := range items {
		select {
		case <-ctx.Done():
			break feed
		case next <- item:
		}
	}
	close(next)
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}
	return ctx.Err()
}
