//go:build integration

package find

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
)

// answersByTitle answers answers_{id} and specificity_{id} from the
// candidate title found in the call state.
func answersByTitle(p map[string]float64, spec map[string]int) fakes.DecideFunc {
	return func(call fakes.Call, q fakes.Question) any {
		var state struct {
			Candidates []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			} `json:"candidates"`
		}
		if err := json.Unmarshal(call.State, &state); err != nil {
			return nil
		}
		if q.ID == "question_kind" {
			return fakes.Pick(KindHowTo, 0.9)
		}
		for _, candidate := range state.Candidates {
			switch q.ID {
			case "answers_" + candidate.ID:
				return fakes.Noul(p[candidate.Title])
			case "specificity_" + candidate.ID:
				return fakes.Level(spec[candidate.Title])
			}
		}
		return nil
	}
}

func TestFindRanksWithNamedExclusions(t *testing.T) {
	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	docs := map[string]string{
		"raw/articles/guide.md":    "---\ntitle: How to tune HNSW indexes\ngenre: tutorial_or_guide\nrelevance: core\ndepth: 2\n---\nTune HNSW with ef and M.\n",
		"raw/articles/paper.md":    "---\ntitle: HNSW benchmark paper\ngenre: paper\nrelevance: core\n---\nWe measure HNSW recall.\n",
		"raw/articles/history.md":  "---\ntitle: A history of HNSW\nrelevance: core\n---\nHNSW was introduced in 2016.\n",
		"raw/articles/spam.md":     "---\ntitle: HNSW tuning spam\ntriage: quarantined\nrelevance: core\n---\nHNSW HNSW tune tune.\n",
		"raw/articles/adjacent.md": "---\ntitle: HNSW tuning notes\nrelevance: adjacent\n---\nTune HNSW.\n",
	}
	for rel, content := range docs {
		path := filepath.Join(info.RootPath, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	fake := fakes.NewOpenRouter(answersByTitle(
		map[string]float64{"How to tune HNSW indexes": 0.9, "HNSW benchmark paper": 0.85, "A history of HNSW": 0.3},
		map[string]int{"How to tune HNSW indexes": 3, "HNSW benchmark paper": 1},
	), nil)
	t.Cleanup(fake.Close)

	cfg := config.Default()
	cfg.OpenRouter.APIKey = "test-key"
	cfg.OpenRouter.APIURL = fake.URL
	s, err := session.Open(session.Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb find"})
	if err != nil {
		t.Fatalf("session.Open: %v", err)
	}

	result, err := Run(context.Background(), s, Query{Text: "how do I tune HNSW", Relevance: "core", Explain: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	hits := make([]string, 0, len(result.Hits))
	for _, hit := range result.Hits {
		hits = append(hits, hit.Path)
	}
	if want := []string{"raw/articles/guide.md", "raw/articles/paper.md"}; !reflect.DeepEqual(hits, want) {
		t.Fatalf("hits = %v, want %v", hits, want)
	}
	if result.QuestionKind != KindHowTo {
		t.Fatalf("question kind = %q", result.QuestionKind)
	}
	reasons := map[string]string{}
	for _, drop := range result.Dropped {
		reasons[drop.Path] = drop.Reason
	}
	want := map[string]string{
		"raw/articles/spam.md":     ReasonQuarantined,
		"raw/articles/adjacent.md": ReasonFilteredFacet,
		"raw/articles/history.md":  ReasonLowProbability,
	}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("dropped = %v, want %v", reasons, want)
	}
	for _, call := range fake.Calls() {
		state := call.StateString()
		if strings.Contains(state, "HNSW tuning spam") || strings.Contains(state, "HNSW tuning notes") {
			t.Fatalf("an excluded document reached the model: %s", state)
		}
	}
	if len(fake.Calls()) != 1 {
		t.Fatalf("calls = %d, want 1 request for 3 candidates", len(fake.Calls()))
	}
}

func TestFindJudgesBatchesConcurrentlyInOrder(t *testing.T) {
	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	const docs = 2*BatchSize + 5
	for index := range docs {
		rel := fmt.Sprintf("raw/articles/note-%02d.md", index)
		content := fmt.Sprintf("---\ntitle: HNSW note %02d\n---\nTune HNSW indexes, note %02d.\n", index, index)
		path := filepath.Join(info.RootPath, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Every request waits until a second request is in flight (or a
	// timeout): sequential batches would time out.
	var (
		mu       sync.Mutex
		seen     = map[string]bool{}
		overlap  = make(chan struct{})
		closed   bool
		timedOut bool
	)
	byTitle := answersByTitle(map[string]float64{}, map[string]int{})
	fake := fakes.NewOpenRouter(func(call fakes.Call, q fakes.Question) any {
		key := call.StateString()
		mu.Lock()
		first := !seen[key]
		seen[key] = true
		if len(seen) >= 2 && !closed {
			closed = true
			close(overlap)
		}
		mu.Unlock()
		if first {
			select {
			case <-overlap:
			case <-time.After(3 * time.Second):
				mu.Lock()
				timedOut = true
				mu.Unlock()
			}
		}
		if strings.HasPrefix(q.ID, "answers_") {
			// note-07 gets a malformed answer: undecided, never a "no".
			if strings.Contains(key, "HNSW note 07") && q.ID == "answers_"+idOf(t, call, "HNSW note 07") {
				return map[string]any{"type": "choice", "choice": "yes"}
			}
			return fakes.Noul(0.9)
		}
		return byTitle(call, q)
	}, nil)
	t.Cleanup(fake.Close)

	cfg := config.Default()
	cfg.OpenRouter.APIKey = "test-key"
	cfg.OpenRouter.APIURL = fake.URL
	cfg.Decisions.Concurrency = 4
	s, err := session.Open(session.Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb find"})
	if err != nil {
		t.Fatalf("session.Open: %v", err)
	}

	result, err := Run(context.Background(), s, Query{Text: "tune HNSW", Limit: 100})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if timedOut {
		t.Fatal("find judged its batches one after another; they must run concurrently")
	}
	if len(fake.Calls()) != 3 {
		t.Fatalf("calls = %d, want 3 batches for %d candidates", len(fake.Calls()), docs)
	}
	if result.Candidates != docs || result.Decided != docs-1 || !reflect.DeepEqual(result.Undecided, []string{"raw/articles/note-07.md"}) {
		t.Fatalf("coverage: candidates %d, decided %d, undecided %v", result.Candidates, result.Decided, result.Undecided)
	}
	joined := strings.Join(result.Lines(), "\n")
	if !strings.Contains(joined, fmt.Sprintf("coverage: %d/%d candidates decided", docs-1, docs)) || !strings.Contains(joined, "undecided candidates (1, not ranked; re-run to retry through the cache): raw/articles/note-07.md") {
		t.Fatalf("lines:\n%s", joined)
	}
	// Equal ranks are ordered by path, whatever order the batches finished.
	for index := 1; index < len(result.Hits); index++ {
		if result.Hits[index-1].Path > result.Hits[index].Path {
			t.Fatalf("hits out of order: %s before %s", result.Hits[index-1].Path, result.Hits[index].Path)
		}
	}
	if len(result.Hits) != docs-1 {
		t.Fatalf("hits = %d, want %d", len(result.Hits), docs-1)
	}
}

// idOf returns the candidate id of the candidate titled title in the call.
func idOf(t *testing.T, call fakes.Call, title string) string {
	t.Helper()
	var state struct {
		Candidates []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(call.State, &state); err != nil {
		t.Errorf("state: %v", err)
		return ""
	}
	for _, candidate := range state.Candidates {
		if candidate.Title == title {
			return candidate.ID
		}
	}
	return ""
}
