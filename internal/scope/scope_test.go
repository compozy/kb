package scope

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/compozy/kb/internal/corpus"
)

func TestFolder(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"raw/articles/a.md":       "raw/articles",
		"raw/audit-2026/b/c.md":   "raw/audit-2026",
		"raw/top.md":              RootFolder,
		"wiki/concepts/Answer.md": RootFolder,
	}
	for input, want := range tests {
		if got := Folder(input); got != want {
			t.Errorf("Folder(%q) = %q, want %q", input, got, want)
		}
	}
}

// folderSources returns n documents per folder, named raw/<folder>/<i>.md.
func folderSources(counts map[string]int) []*corpus.Document {
	var paths []string
	for folder, n := range counts {
		for i := range n {
			paths = append(paths, fmt.Sprintf("raw/%s/%03d.md", folder, i))
		}
	}
	return docs(paths...)
}

func perFolder(sample []*corpus.Document) map[string]int {
	counts := map[string]int{}
	for _, doc := range sample {
		counts[Folder(doc.Path)]++
	}
	return counts
}

func TestStratify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		counts map[string]int
		n      int
		want   map[string]int
	}{
		{name: "smaller than n returns everything", counts: map[string]int{"a": 3, "b": 2}, n: 10, want: map[string]int{"raw/a": 3, "raw/b": 2}},
		{name: "proportional quotas", counts: map[string]int{"a": 800, "b": 150, "c": 50}, n: 100, want: map[string]int{"raw/a": 80, "raw/b": 15, "raw/c": 5}},
		{name: "largest remainder fills the rest", counts: map[string]int{"a": 5, "b": 5, "c": 5}, n: 10, want: map[string]int{"raw/a": 4, "raw/b": 3, "raw/c": 3}},
		{name: "every folder gets one slot", counts: map[string]int{"big": 997, "tiny": 1, "small": 2}, n: 10, want: map[string]int{"raw/big": 8, "raw/tiny": 1, "raw/small": 1}},
		{name: "more folders than slots keeps the largest remainders", counts: map[string]int{"a": 10, "b": 5, "c": 1, "d": 1}, n: 2, want: map[string]int{"raw/a": 1, "raw/b": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sources := folderSources(tt.counts)
			sample := Stratify(sources, tt.n)
			if got := perFolder(sample); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("per folder = %v, want %v", got, tt.want)
			}
			again := Stratify(sources, tt.n)
			if !reflect.DeepEqual(sample, again) {
				t.Fatal("Stratify is not deterministic")
			}
			for i := 1; i < len(sample); i++ {
				if sample[i-1].Path >= sample[i].Path {
					t.Fatalf("sample not sorted by path at %d: %s >= %s", i, sample[i-1].Path, sample[i].Path)
				}
			}
		})
	}
}

func TestStratifySpreadsInsideAFolder(t *testing.T) {
	t.Parallel()
	sample := Stratify(folderSources(map[string]int{"a": 10}), 2)
	got := []string{sample[0].Path, sample[1].Path}
	want := []string{"raw/a/002.md", "raw/a/007.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sample = %v, want %v", got, want)
	}
}

func TestForEachStopsOnFirstError(t *testing.T) {
	t.Parallel()
	boom := errors.New("boom")
	var ran atomic.Int32
	err := forEach(context.Background(), 2, []int{1, 2, 3, 4, 5, 6, 7, 8}, func(_ context.Context, item int) error {
		ran.Add(1)
		if item == 1 {
			return boom
		}
		return nil
	})
	if !errors.Is(err, boom) {
		t.Fatalf("forEach error = %v, want boom", err)
	}
	if err := forEach(context.Background(), 4, []int{1, 2, 3}, func(context.Context, int) error { return nil }); err != nil {
		t.Fatalf("forEach = %v", err)
	}
}
