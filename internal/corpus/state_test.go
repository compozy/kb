package corpus_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/compozy/kb/internal/corpus"
)

func TestStateStore(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("OpenState(empty) returned error: %v", err)
	}
	if _, ok := store.Get("raw/a.md"); ok {
		t.Fatal("empty store returned a row")
	}

	first := corpus.StateRow{Path: "raw/a.md", BodyHash: "h1", Contract: "c1", Banks: map[string]string{"classify": "1"}, Written: map[string]string{"summary": "v1"}, Updated: "t1"}
	second := corpus.StateRow{Path: "raw/a.md", BodyHash: "h2", Contract: "c1", Banks: map[string]string{"classify": "2"}, Written: map[string]string{}, Updated: "t2"}
	other := corpus.StateRow{Path: "raw/b.md", BodyHash: "h2", Updated: "t3"}
	for _, row := range []corpus.StateRow{first, second, other} {
		if err := store.Put(row); err != nil {
			t.Fatalf("Put returned error: %v", err)
		}
	}
	if err := store.Put(corpus.StateRow{}); err == nil {
		t.Fatal("Put without path returned nil error")
	}

	row, ok := store.Get("raw/a.md")
	if !ok || !reflect.DeepEqual(*row, second) {
		t.Fatalf("Get = %#v, want %#v", row, second)
	}
	row.Banks["classify"] = "mutated"
	if again, _ := store.Get("raw/a.md"); again.Banks["classify"] != "2" {
		t.Fatal("Get must return a copy")
	}

	if got := store.FindByBodyHash("h2"); len(got) != 2 || got[0].Path != "raw/a.md" || got[1].Path != "raw/b.md" {
		t.Fatalf("FindByBodyHash = %#v", got)
	}

	reopened, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("reopen returned error: %v", err)
	}
	if got := reopened.Rows(); !reflect.DeepEqual(got, []corpus.StateRow{second, other}) {
		t.Fatalf("Rows after reopen = %#v", got)
	}

	content, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	if lines := strings.Count(string(content), "\n"); lines != 3 {
		t.Fatalf("state lines = %d, want 3 (append-only)", lines)
	}

	if err := reopened.Compact(); err != nil {
		t.Fatalf("Compact returned error: %v", err)
	}
	compacted, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatalf("read compacted state: %v", err)
	}
	if lines := strings.Count(string(compacted), "\n"); lines != 2 {
		t.Fatalf("compacted lines = %d, want 2", lines)
	}
	afterCompact, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("reopen after compact: %v", err)
	}
	if got := afterCompact.Rows(); !reflect.DeepEqual(got, []corpus.StateRow{second, other}) {
		t.Fatalf("Rows after compact = %#v", got)
	}
}

func TestStateStoreTornAndMalformedLines(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	statePath := writeTopicFile(t, root, corpus.StateFile, `{"path":"raw/a.md","body_hash":"h1"}`+"\n"+`{"path":"raw/b.md","body_ha`)
	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("OpenState(torn) returned error: %v", err)
	}
	if got := store.Rows(); len(got) != 1 || got[0].Path != "raw/a.md" {
		t.Fatalf("Rows(torn) = %#v", got)
	}
	if err := store.Put(corpus.StateRow{Path: "raw/c.md", BodyHash: "h3"}); err != nil {
		t.Fatalf("Put after torn line: %v", err)
	}
	reopened, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("reopen after torn append: %v", err)
	}
	if got := reopened.Rows(); len(got) != 2 || got[1].Path != "raw/c.md" {
		t.Fatalf("Rows after torn append = %#v", got)
	}

	if err := os.WriteFile(statePath, []byte("not json\n"+`{"path":"raw/a.md"}`+"\n"), 0o644); err != nil {
		t.Fatalf("write malformed state: %v", err)
	}
	if _, err := corpus.OpenState(root); err == nil {
		t.Fatal("OpenState(malformed middle line) returned nil error")
	}
}

func TestStateStoreLookupFollowsRenames(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("OpenState: %v", err)
	}
	for _, row := range []corpus.StateRow{
		{Path: "raw/old.md", BodyHash: "moved", Written: map[string]string{"summary": "v"}},
		{Path: "raw/still-here.md", BodyHash: "copy"},
	} {
		if err := store.Put(row); err != nil {
			t.Fatalf("Put: %v", err)
		}
	}
	exists := func(path string) bool { return path == "raw/still-here.md" || path == "raw/new.md" }

	row, ok := store.Lookup("raw/new.md", "moved", exists)
	if !ok || row.Path != "raw/new.md" || row.Written["summary"] != "v" {
		t.Fatalf("Lookup(renamed) = %#v, %v", row, ok)
	}
	if _, ok := store.Lookup("raw/copy.md", "copy", exists); ok {
		t.Fatal("Lookup must not follow a row whose file still exists")
	}
	if _, ok := store.Lookup("raw/none.md", "", exists); ok {
		t.Fatal("Lookup without body hash must not match")
	}
	if row, ok := store.Lookup("raw/old.md", "other", exists); !ok || row.BodyHash != "moved" {
		t.Fatalf("Lookup(direct) = %#v, %v", row, ok)
	}
	if _, ok := store.Get("raw/new.md"); ok {
		t.Fatal("Lookup must not modify the store")
	}
}

func TestStateStoreConcurrentPut(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("OpenState: %v", err)
	}

	var group sync.WaitGroup
	for index := range 20 {
		group.Go(func() {
			if err := store.Put(corpus.StateRow{Path: filepath.ToSlash(filepath.Join("raw", string(rune('a'+index))+".md")), BodyHash: "h"}); err != nil {
				t.Errorf("Put: %v", err)
			}
		})
	}
	group.Wait()

	reopened, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if got := len(reopened.Rows()); got != 20 {
		t.Fatalf("rows = %d, want 20", got)
	}
}

func TestUnclassified(t *testing.T) {
	t.Parallel()

	doc := &corpus.Document{Path: "raw/a.md", BodyHash: "h"}
	row := &corpus.StateRow{Path: "raw/a.md", BodyHash: "h", Contract: "c", Banks: map[string]string{"classify": "1", "extra": "9"}}
	// Another command wrote the document after its body changed: the row's
	// body hash and contract are current, the classify judgment is not.
	staleJudgment := &corpus.StateRow{
		Path: "raw/a.md", BodyHash: "h", Contract: "c", Banks: map[string]string{"classify": "1", "link": "1"},
		BankBody: map[string]string{"classify": "old", "link": "h"}, BankContract: map[string]string{"classify": "c", "link": "c"},
	}
	oldContractJudgment := &corpus.StateRow{
		Path: "raw/a.md", BodyHash: "h", Contract: "c", Banks: map[string]string{"classify": "1"},
		BankBody: map[string]string{"classify": "h"}, BankContract: map[string]string{"classify": "b"},
	}
	tests := []struct {
		name     string
		doc      *corpus.Document
		row      *corpus.StateRow
		contract string
		banks    map[string]string
		want     bool
	}{
		{name: "no row", doc: doc, want: true},
		{name: "current", doc: doc, row: row, contract: "c", banks: map[string]string{"classify": "1"}},
		{name: "no contract required", doc: doc, row: row, banks: map[string]string{"classify": "1"}},
		{name: "stale body", doc: &corpus.Document{BodyHash: "other"}, row: row, contract: "c", want: true},
		{name: "contract changed", doc: doc, row: row, contract: "d", want: true},
		{name: "bank version changed", doc: doc, row: row, contract: "c", banks: map[string]string{"classify": "2"}, want: true},
		{name: "bank missing", doc: doc, row: row, contract: "c", banks: map[string]string{"link": "1"}, want: true},
		{name: "bank judged on an older body", doc: doc, row: staleJudgment, contract: "c", banks: map[string]string{"classify": "1"}, want: true},
		{name: "other bank judged on the current body", doc: doc, row: staleJudgment, contract: "c", banks: map[string]string{"link": "1"}},
		{name: "bank judged under an older contract", doc: doc, row: oldContractJudgment, contract: "c", banks: map[string]string{"classify": "1"}, want: true},
	}
	for _, tc := range tests {
		if got := corpus.Unclassified(tc.doc, tc.row, tc.contract, tc.banks); got != tc.want {
			t.Errorf("%s: Unclassified = %v, want %v", tc.name, got, tc.want)
		}
	}
}
