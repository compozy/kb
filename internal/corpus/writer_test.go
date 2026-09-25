package corpus_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/frontmatter"
)

var fixedNow = func() time.Time { return time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC) }

const writerFixture = "---\n" +
	"# user comment\n" +
	"title: \"Source: One\"\n" +
	"tags:\n" +
	"    - kb\n" +
	"    - raw\n" +
	"source_url: https://example.com/a\n" +
	"notes: |\n" +
	"  keep me\n" +
	"\n" +
	"  exactly\n" +
	"---\n" +
	"# Body\n\nText with [[Link]].\n"

type writerEnv struct {
	root   string
	rel    string
	store  *corpus.StateStore
	writer *corpus.Writer
}

func newWriterEnv(t *testing.T, relative, content string) writerEnv {
	t.Helper()

	root := t.TempDir()
	writeTopicFile(t, root, relative, content)
	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("OpenState: %v", err)
	}

	return writerEnv{root: root, rel: relative, store: store, writer: corpus.NewWriter(root, store, fixedNow)}
}

func (env writerEnv) load(t *testing.T) *corpus.Document {
	t.Helper()

	document, err := corpus.ReadDocument(env.root, env.rel, corpus.KindSource)
	if err != nil {
		t.Fatalf("ReadDocument: %v", err)
	}
	return document
}

func (env writerEnv) read(t *testing.T) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(env.root, filepath.FromSlash(env.rel)))
	if err != nil {
		t.Fatalf("read %s: %v", env.rel, err)
	}
	return string(content)
}

func (env writerEnv) overwrite(t *testing.T, content string) {
	t.Helper()
	writeTopicFile(t, env.root, env.rel, content)
}

func (env writerEnv) apply(t *testing.T, doc *corpus.Document, updates map[string]any) corpus.WriteResult {
	t.Helper()

	result, err := env.writer.Apply(doc, updates, corpus.StateMeta{Contract: "contract-1", Banks: map[string]string{"classify": "1"}})
	if err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}
	return result
}

// assertNonOwnedPreserved checks that every non-owned key block, the
// leading comment and the body are byte-identical.
func assertNonOwnedPreserved(t *testing.T, before, after string) {
	t.Helper()

	nonOwned := func(text string) []frontmatter.KeyBlock {
		blocks, err := frontmatter.KeyBlocks(text)
		if err != nil {
			t.Fatalf("KeyBlocks: %v", err)
		}
		kept := make([]frontmatter.KeyBlock, 0, len(blocks))
		for _, block := range blocks {
			if !decisions.IsOwnedKey(block.Key) {
				kept = append(kept, block)
			}
		}
		return kept
	}
	if got, want := nonOwned(after), nonOwned(before); !reflect.DeepEqual(got, want) {
		t.Fatalf("non-owned keys changed:\n got %#v\nwant %#v", got, want)
	}
	_, beforeBody, _ := frontmatter.Split(before)
	_, afterBody, _ := frontmatter.Split(after)
	if beforeBody != afterBody {
		t.Fatalf("body changed: %q -> %q", beforeBody, afterBody)
	}
	if strings.Contains(before, "# user comment") && !strings.Contains(after, "# user comment") {
		t.Fatal("frontmatter comment lost")
	}
}

func TestWriterWritesAbsentOwnedKeys(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	doc := env.load(t)
	result := env.apply(t, doc, map[string]any{
		"summary":  "A short summary.",
		"concepts": []string{"[[Agents]]", "[[Tools]]"},
		"depth":    2.5,
	})

	if result.Status != corpus.StatusWritten || !reflect.DeepEqual(result.Written, []string{"concepts", "depth", "summary"}) || len(result.SkippedUserKeys) != 0 {
		t.Fatalf("result = %#v", result)
	}
	after := env.read(t)
	assertNonOwnedPreserved(t, writerFixture, after)
	values, _, err := frontmatter.Parse(after)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if values["summary"] != "A short summary." || !reflect.DeepEqual(values["concepts"], []string{"[[Agents]]", "[[Tools]]"}) || values["depth"] != 2.5 {
		t.Fatalf("values = %#v", values)
	}
	if !strings.Contains(after, "    - '[[Agents]]'\n") {
		t.Fatalf("wikilinks must be quoted:\n%s", after)
	}

	row, ok := env.store.Get(env.rel)
	if !ok {
		t.Fatal("state row missing")
	}
	body := corpus.BodyHash("# Body\n\nText with [[Link]].\n")
	wantRow := corpus.StateRow{
		Path:     env.rel,
		BodyHash: body,
		Contract: "contract-1",
		Banks:    map[string]string{"classify": "1"},
		Written: map[string]string{
			"summary":  corpus.ValueHash("A short summary."),
			"concepts": corpus.ValueHash([]string{"[[Agents]]", "[[Tools]]"}),
			"depth":    corpus.ValueHash(2.5),
		},
		BankBody:     map[string]string{"classify": body},
		BankContract: map[string]string{"classify": "contract-1"},
		WrittenBody:  map[string]string{"summary": body, "concepts": body, "depth": body},
		Facts:        map[string]string{},
		Updated:      "2026-09-24T12:00:00Z",
	}
	if !reflect.DeepEqual(*row, wantRow) {
		t.Fatalf("row = %#v\nwant %#v", *row, wantRow)
	}
	if doc.Raw != after || doc.Summary() != "A short summary." {
		t.Fatal("document must be refreshed in place")
	}
}

func TestWriterUpdatesKbWrittenValues(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	doc := env.load(t)
	env.apply(t, doc, map[string]any{"summary": "first", "genre": "paper"})

	doc = env.load(t)
	result := env.apply(t, doc, map[string]any{"summary": "second", "genre": "paper"})
	if result.Status != corpus.StatusWritten || !reflect.DeepEqual(result.Written, []string{"summary"}) {
		t.Fatalf("result = %#v", result)
	}
	after := env.read(t)
	assertNonOwnedPreserved(t, writerFixture, after)
	if text, _ := frontmatter.TopLevelKeyText(after, "summary"); text != "summary: second\n" {
		t.Fatalf("summary block = %q", text)
	}
	if row, _ := env.store.Get(env.rel); row.Written["summary"] != corpus.ValueHash("second") || row.Written["genre"] != corpus.ValueHash("paper") {
		t.Fatalf("row written = %#v", row.Written)
	}
}

func TestWriterLeavesUserEditedValues(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	env.apply(t, env.load(t), map[string]any{"summary": "kb summary"})

	edited := strings.Replace(env.read(t), "summary: kb summary", "summary: my own words", 1)
	env.overwrite(t, edited)

	result := env.apply(t, env.load(t), map[string]any{"summary": "kb summary v2", "genre": "paper"})
	if !reflect.DeepEqual(result.SkippedUserKeys, []string{"summary"}) || !reflect.DeepEqual(result.Written, []string{"genre"}) {
		t.Fatalf("result = %#v", result)
	}
	after := env.read(t)
	if before, _ := frontmatter.TopLevelKeyText(edited, "summary"); !strings.Contains(after, before) {
		t.Fatalf("user summary changed:\n%s", after)
	}
}

func TestWriterSkipsUserKeysThatPredateKb(t *testing.T) {
	t.Parallel()

	content := strings.Replace(writerFixture, "source_url:", "summary: written by the user\nsource_url:", 1)
	env := newWriterEnv(t, "raw/articles/one.md", content)
	result := env.apply(t, env.load(t), map[string]any{"summary": "kb summary"})

	if result.Status != corpus.StatusUnchanged || !reflect.DeepEqual(result.SkippedUserKeys, []string{"summary"}) || len(result.Written) != 0 {
		t.Fatalf("result = %#v", result)
	}
	if env.read(t) != content {
		t.Fatal("file must be untouched")
	}
	row, ok := env.store.Get(env.rel)
	if !ok || row.Contract != "contract-1" || len(row.Written) != 0 {
		t.Fatalf("state row must record the classification without claiming keys: %#v", row)
	}
}

func TestWriterRespectsLocked(t *testing.T) {
	t.Parallel()

	content := strings.Replace(writerFixture, "tags:", "locked: true\ntags:", 1)
	env := newWriterEnv(t, "raw/articles/one.md", content)
	doc := env.load(t)
	result := env.apply(t, doc, map[string]any{"summary": "x"})
	if result.Status != corpus.StatusSkippedLocked {
		t.Fatalf("result = %#v", result)
	}

	// Locked after load: the fresh file wins.
	env2 := newWriterEnv(t, "raw/articles/two.md", writerFixture)
	doc2 := env2.load(t)
	env2.overwrite(t, strings.Replace(writerFixture, "tags:", "locked: true\ntags:", 1))
	if result := env2.apply(t, doc2, map[string]any{"summary": "x"}); result.Status != corpus.StatusSkippedLocked {
		t.Fatalf("late lock result = %#v", result)
	}

	if env.read(t) != content {
		t.Fatal("locked file must be untouched")
	}
	if _, ok := env.store.Get(env.rel); ok {
		t.Fatal("locked file must get no state row")
	}
}

func TestWriterSkipsWhenBodyChanged(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	doc := env.load(t)
	changed := writerFixture + "More text.\n"
	env.overwrite(t, changed)

	result := env.apply(t, doc, map[string]any{"summary": "x"})
	if result.Status != corpus.StatusSkippedChanged {
		t.Fatalf("result = %#v", result)
	}
	if env.read(t) != changed {
		t.Fatal("file must be untouched")
	}

	bodyResult, err := env.writer.ApplyBody(doc, "# New body\n", nil, corpus.StateMeta{})
	if err != nil || bodyResult.Status != corpus.StatusSkippedChanged {
		t.Fatalf("ApplyBody = %#v, %v", bodyResult, err)
	}
}

func TestWriterRetriesOnNonOwnedKeyChange(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	doc := env.load(t)
	userEdit := strings.Replace(writerFixture, "source_url: https://example.com/a", "source_url: https://example.com/b\nrating: 5", 1)
	env.overwrite(t, userEdit)

	result := env.apply(t, doc, map[string]any{"summary": "x"})
	if result.Status != corpus.StatusWritten {
		t.Fatalf("result = %#v", result)
	}
	after := env.read(t)
	assertNonOwnedPreserved(t, userEdit, after)
	if !strings.Contains(after, "rating: 5\n") || !strings.Contains(after, "summary: x\n") {
		t.Fatalf("user edit or summary missing:\n%s", after)
	}
}

func TestWriterMergesAliases(t *testing.T) {
	t.Parallel()

	content := "---\ntitle: Artificial Intelligence\naliases:\n  - AI\n  - ia\n---\nbody\n"
	env := newWriterEnv(t, "wiki/concepts/Artificial Intelligence.md", content)
	doc := env.load(t)

	result := env.apply(t, doc, map[string]any{"aliases": []string{"IA", "Inteligência Artificial"}})
	if result.Status != corpus.StatusWritten || !reflect.DeepEqual(result.Written, []string{"aliases"}) {
		t.Fatalf("result = %#v", result)
	}
	values, _, err := frontmatter.Parse(env.read(t))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if want := []string{"AI", "ia", "Inteligência Artificial"}; !reflect.DeepEqual(values["aliases"], want) {
		t.Fatalf("aliases = %#v, want %#v", values["aliases"], want)
	}
	if !reflect.DeepEqual(doc.Aliases, []string{"AI", "ia", "Inteligência Artificial"}) {
		t.Fatalf("doc aliases = %#v", doc.Aliases)
	}

	again := env.apply(t, env.load(t), map[string]any{"aliases": []string{"ai"}})
	if again.Status != corpus.StatusUnchanged || len(again.Written) != 0 {
		t.Fatalf("merge of known alias = %#v", again)
	}
}

func TestWriterPreservesCRLF(t *testing.T) {
	t.Parallel()

	content := strings.ReplaceAll(writerFixture, "\n", "\r\n")
	env := newWriterEnv(t, "raw/articles/crlf.md", content)
	result := env.apply(t, env.load(t), map[string]any{"summary": "x", "concepts": []string{"[[A]]"}})
	if result.Status != corpus.StatusWritten {
		t.Fatalf("result = %#v", result)
	}
	after := env.read(t)
	if strings.Count(after, "\n") != strings.Count(after, "\r\n") {
		t.Fatalf("mixed line endings:\n%q", after)
	}
	assertNonOwnedPreserved(t, content, after)
	if !strings.Contains(after, "concepts:\r\n    - '[[A]]'\r\nsummary: x\r\n---\r\n") {
		t.Fatalf("CRLF keys missing:\n%q", after)
	}
}

func TestWriterUnchangedValuesAreNotRewritten(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	env.apply(t, env.load(t), map[string]any{"summary": "same"})
	before := env.read(t)
	beforeState, err := os.ReadFile(env.store.Path())
	if err != nil {
		t.Fatalf("read state: %v", err)
	}

	result := env.apply(t, env.load(t), map[string]any{"summary": "same"})
	if result.Status != corpus.StatusUnchanged || len(result.Written) != 0 {
		t.Fatalf("result = %#v", result)
	}
	if env.read(t) != before {
		t.Fatal("unchanged value must not be rewritten")
	}
	afterState, err := os.ReadFile(env.store.Path())
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	if string(afterState) != string(beforeState) {
		t.Fatal("identical state must not append a row")
	}
}

func TestWriterDeletesKbWrittenKeys(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	env.apply(t, env.load(t), map[string]any{"triage": "review", "triage_reason": "thin"})

	result := env.apply(t, env.load(t), map[string]any{"triage": "kept", "triage_reason": nil})
	if result.Status != corpus.StatusWritten || !reflect.DeepEqual(result.Written, []string{"triage", "triage_reason"}) {
		t.Fatalf("result = %#v", result)
	}
	after := env.read(t)
	if strings.Contains(after, "triage_reason") || !strings.Contains(after, "triage: kept\n") {
		t.Fatalf("delete failed:\n%s", after)
	}
	if row, _ := env.store.Get(env.rel); row.Written["triage_reason"] != "" {
		t.Fatalf("deleted key still recorded: %#v", row.Written)
	}
	assertNonOwnedPreserved(t, writerFixture, after)
}

func TestWriterRefusesNonOwnedKeys(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	for _, key := range []string{"title", "locked", "ai_summary"} {
		if _, err := env.writer.Apply(env.load(t), map[string]any{key: "x"}, corpus.StateMeta{}); err == nil {
			t.Fatalf("Apply(%q) returned nil error", key)
		}
	}
	if env.read(t) != writerFixture {
		t.Fatal("refused write must not touch the file")
	}
	if _, err := env.writer.Apply(nil, nil, corpus.StateMeta{}); err == nil {
		t.Fatal("Apply(nil doc) returned nil error")
	}
}

func TestWriterWritesStateAfterFile(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	if err := os.MkdirAll(env.store.Path(), 0o755); err != nil {
		t.Fatalf("block state file: %v", err)
	}

	_, err := env.writer.Apply(env.load(t), map[string]any{"summary": "x"}, corpus.StateMeta{})
	if err == nil {
		t.Fatal("Apply must report the failed state append")
	}
	if !strings.Contains(env.read(t), "summary: x\n") {
		t.Fatal("the file must be written before the state row")
	}
}

func TestWriterFollowsRenamedFiles(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/old.md", writerFixture)
	env.apply(t, env.load(t), map[string]any{"summary": "kb"})

	oldPath := filepath.Join(env.root, "raw", "articles", "old.md")
	newPath := filepath.Join(env.root, "raw", "articles", "new.md")
	if err := os.Rename(oldPath, newPath); err != nil {
		t.Fatalf("rename: %v", err)
	}
	renamed := writerEnv{root: env.root, rel: "raw/articles/new.md", store: env.store, writer: env.writer}

	result := renamed.apply(t, renamed.load(t), map[string]any{"summary": "kb v2"})
	if result.Status != corpus.StatusWritten || !reflect.DeepEqual(result.Written, []string{"summary"}) {
		t.Fatalf("renamed result = %#v", result)
	}
	row, ok := env.store.Get("raw/articles/new.md")
	if !ok || row.Contract != "contract-1" || row.Written["summary"] != corpus.ValueHash("kb v2") {
		t.Fatalf("renamed row = %#v", row)
	}
}

func TestWriterApplyBody(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	doc := env.load(t)
	newBody := "# Body\n\nText with [[Link]] and [[Agents|agents]].\n"

	result, err := env.writer.ApplyBody(doc, newBody, map[string]any{"related": []string{"[[Agents]]"}}, corpus.StateMeta{})
	if err != nil {
		t.Fatalf("ApplyBody: %v", err)
	}
	if result.Status != corpus.StatusWritten || !result.BodyWritten || !reflect.DeepEqual(result.Written, []string{"related"}) {
		t.Fatalf("result = %#v", result)
	}
	after := env.read(t)
	if _, body, _ := frontmatter.Split(after); body != newBody {
		t.Fatalf("body = %q", body)
	}
	if doc.BodyHash != corpus.BodyHash(newBody) {
		t.Fatal("document body hash must be refreshed")
	}
	if row, _ := env.store.Get(env.rel); row.BodyHash != corpus.BodyHash(newBody) {
		t.Fatalf("row body hash = %q", row.BodyHash)
	}

	same, err := env.writer.ApplyBody(doc, newBody, nil, corpus.StateMeta{})
	if err != nil || same.Status != corpus.StatusUnchanged || same.BodyWritten {
		t.Fatalf("same body = %#v, %v", same, err)
	}
}

func TestWriterPreservesFileMode(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	absolute := filepath.Join(env.root, "raw", "articles", "one.md")
	if err := os.Chmod(absolute, 0o600); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	env.apply(t, env.load(t), map[string]any{"summary": "x"})
	info, err := os.Stat(absolute)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %v, want 0600", info.Mode().Perm())
	}
	entries, err := os.ReadDir(filepath.Dir(absolute))
	if err != nil || len(entries) != 1 {
		t.Fatalf("temp files left behind: %v %v", entries, err)
	}
}

func TestWriterConcurrentDocuments(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	for index := range 8 {
		writeTopicFile(t, root, fmt.Sprintf("raw/articles/doc-%d.md", index), writerFixture)
	}
	store, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("OpenState: %v", err)
	}
	writer := corpus.NewWriter(root, store, fixedNow)
	loaded, err := corpus.Load(root, corpus.LoadOptions{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	var group sync.WaitGroup
	for index, doc := range loaded.Sources() {
		group.Go(func() {
			if _, err := writer.Apply(doc, map[string]any{"summary": fmt.Sprintf("s%d", index)}, corpus.StateMeta{}); err != nil {
				t.Errorf("Apply: %v", err)
			}
		})
	}
	group.Wait()

	reopened, err := corpus.OpenState(root)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if got := len(reopened.Rows()); got != 8 {
		t.Fatalf("rows = %d, want 8", got)
	}
}

func TestWriterRecordsJudgmentsPerBank(t *testing.T) {
	t.Parallel()

	classifyMeta := corpus.StateMeta{Contract: "contract-1", Banks: map[string]string{"classify": "1"}}
	classifyBanks := map[string]string{"classify": "1"}

	tests := []struct {
		name string
		// edit runs after the classify write, before the second write.
		edit func(t *testing.T, env writerEnv)
		// body, when set, replaces the body in the second write.
		body     func(old string) string
		updates  map[string]any
		meta     corpus.StateMeta
		contract string
		// wantUnclassified is corpus.Unclassified for the classify bank
		// after the second write.
		wantUnclassified bool
		// wantSummaryFresh reports whether the summary is recorded as
		// derived from the current body.
		wantSummaryFresh bool
	}{
		{
			name: "link after a user body edit leaves the document unclassified",
			edit: func(t *testing.T, env writerEnv) {
				env.overwrite(t, strings.Replace(env.read(t), "Text with [[Link]].", "Text with [[Link]], edited by the user.", 1))
			},
			updates:          map[string]any{"related": []string{"[[Agents]]"}},
			meta:             corpus.StateMeta{Contract: "contract-1", Banks: map[string]string{"link": "1"}},
			contract:         "contract-1",
			wantUnclassified: true,
		},
		{
			name:             "link body insertion keeps the classification current",
			body:             func(old string) string { return strings.Replace(old, "Text with", "[[Text|Text]] with", 1) },
			updates:          map[string]any{"related": []string{"[[Text]]"}},
			meta:             corpus.StateMeta{Contract: "contract-1", Banks: map[string]string{"link": "1"}},
			contract:         "contract-1",
			wantSummaryFresh: true,
		},
		{
			name:             "a content-changing body replacement leaves the document unclassified",
			body:             func(string) string { return "# Recaptured\n\nNew text.\n" },
			updates:          map[string]any{"recaptured": "2026-09-24"},
			meta:             corpus.StateMeta{Contract: "contract-1"},
			contract:         "contract-1",
			wantUnclassified: true,
		},
		{
			name:             "a contract recorded by link does not refresh classify",
			updates:          map[string]any{"related": []string{"[[Agents]]"}},
			meta:             corpus.StateMeta{Contract: "contract-9", Banks: map[string]string{"link": "1"}},
			contract:         "contract-9",
			wantUnclassified: true,
			wantSummaryFresh: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
			if _, err := env.writer.Apply(env.load(t), map[string]any{"summary": "A summary."}, classifyMeta); err != nil {
				t.Fatal(err)
			}
			if tc.edit != nil {
				tc.edit(t, env)
			}
			doc := env.load(t)
			var (
				result corpus.WriteResult
				err    error
			)
			if tc.body != nil {
				result, err = env.writer.ApplyBody(doc, tc.body(doc.Body), tc.updates, tc.meta)
			} else {
				result, err = env.writer.Apply(doc, tc.updates, tc.meta)
			}
			if err != nil || result.Status != corpus.StatusWritten {
				t.Fatalf("second write = %#v, %v", result, err)
			}
			doc = env.load(t)
			row, ok := env.store.Get(env.rel)
			if !ok {
				t.Fatal("state row missing")
			}
			if row.BodyHash != doc.BodyHash {
				t.Fatal("the row body hash must follow the file (rename detection)")
			}
			if got := corpus.Unclassified(doc, row, tc.contract, classifyBanks); got != tc.wantUnclassified {
				t.Fatalf("Unclassified = %v, want %v (row %#v)", got, tc.wantUnclassified, row)
			}
			if got := row.KeyBody("summary") == doc.BodyHash; got != tc.wantSummaryFresh {
				t.Fatalf("summary fresh = %v, want %v", got, tc.wantSummaryFresh)
			}
			for bank := range tc.meta.Banks {
				if row.JudgedBody(bank) != doc.BodyHash || row.JudgedContract(bank) != tc.meta.Contract {
					t.Fatalf("bank %s of the second write must be current: %#v", bank, row)
				}
			}
		})
	}
}

func TestWriterPinsLegacyRowJudgments(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	oldBody := corpus.BodyHash("# Body\n\nText with [[Link]].\n")
	legacy := corpus.StateRow{
		Path: env.rel, BodyHash: oldBody, Contract: "contract-1",
		Banks: map[string]string{"classify": "1"}, Written: map[string]string{}, Updated: "t0",
	}
	if err := env.store.Put(legacy); err != nil {
		t.Fatal(err)
	}
	env.overwrite(t, strings.Replace(writerFixture, "Text with", "Edited text with", 1))
	doc := env.load(t)
	if _, err := env.writer.Apply(doc, map[string]any{"related": []string{"[[Agents]]"}}, corpus.StateMeta{Contract: "contract-2", Banks: map[string]string{"link": "1"}}); err != nil {
		t.Fatal(err)
	}
	row, _ := env.store.Get(env.rel)
	if row.JudgedBody("classify") != oldBody || row.JudgedContract("classify") != "contract-1" {
		t.Fatalf("a legacy classify judgment must stay pinned to its body and contract: %#v", row)
	}
	if !corpus.Unclassified(env.load(t), row, "", map[string]string{"classify": "1"}) {
		t.Fatal("a legacy row must not look classified after an unrelated write")
	}
}

func TestWriterMergesFacts(t *testing.T) {
	t.Parallel()

	env := newWriterEnv(t, "raw/articles/one.md", writerFixture)
	apply := func(facts map[string]string) *corpus.StateRow {
		t.Helper()
		if _, err := env.writer.Apply(env.load(t), nil, corpus.StateMeta{Facts: facts}); err != nil {
			t.Fatal(err)
		}
		row, _ := env.store.Get(env.rel)
		return row
	}
	if row := apply(map[string]string{"primary_concept": "none", "other": "x"}); row.Facts["primary_concept"] != "none" || row.Facts["other"] != "x" {
		t.Fatalf("facts = %#v", row.Facts)
	}
	if row := apply(map[string]string{"primary_concept": ""}); row.Facts["primary_concept"] != "" || row.Facts["other"] != "x" {
		t.Fatalf("an empty fact must delete only that fact: %#v", row.Facts)
	}
}

func TestWriterMergesUserRelationLists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		frontmatter string
		key         string
		desired     any
		want        any
		wantWritten bool
		wantSkipped bool
	}{
		{
			name:        "new targets are appended after the user's items",
			frontmatter: "related:\n  - \"[[Zeta]]\"\n  - my own note\n  - \"[[Alpha|the alpha]]\"\n",
			key:         "related",
			desired:     []string{"[[Alpha]]", "[[Beta]]", "[[beta]]"},
			want:        []string{"[[Zeta]]", "my own note", "[[Alpha|the alpha]]", "[[Beta]]"},
			wantWritten: true,
		},
		{
			name:        "a desired value without the user's items never removes them",
			frontmatter: "affects:\n  - \"[[wiki/concepts/Gamma.md]]\"\n",
			key:         "affects",
			desired:     []string{"[[Delta]]"},
			want:        []string{"[[wiki/concepts/Gamma.md]]", "[[Delta]]"},
			wantWritten: true,
		},
		{
			name:        "nothing new leaves the list alone",
			frontmatter: "supersedes:\n  - \"[[Old]]\"\n",
			key:         "supersedes",
			desired:     []string{"[[old]]"},
			want:        []string{"[[Old]]"},
		},
		{
			name:        "a deletion of a user list is refused",
			frontmatter: "extends:\n  - \"[[Old]]\"\n",
			key:         "extends",
			desired:     nil,
			want:        []string{"[[Old]]"},
			wantSkipped: true,
		},
		{
			name:        "a scalar user value is not rewritten",
			frontmatter: "related: see the index\n",
			key:         "related",
			desired:     []string{"[[Beta]]"},
			want:        "see the index",
			wantSkipped: true,
		},
		{
			name:        "concepts keep the conflict rule",
			frontmatter: "concepts:\n  - \"[[Mine]]\"\n",
			key:         "concepts",
			desired:     []string{"[[Mine]]", "[[Theirs]]"},
			want:        []string{"[[Mine]]"},
			wantSkipped: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := newWriterEnv(t, "raw/articles/one.md", "---\ntitle: One\n"+tc.frontmatter+"---\nbody\n")
			result := env.apply(t, env.load(t), map[string]any{tc.key: tc.desired})
			if got := len(result.Written) == 1; got != tc.wantWritten {
				t.Fatalf("written = %v, want %v", result.Written, tc.wantWritten)
			}
			if got := len(result.SkippedUserKeys) == 1; got != tc.wantSkipped {
				t.Fatalf("skipped = %v, want %v", result.SkippedUserKeys, tc.wantSkipped)
			}
			values, _, err := frontmatter.Parse(env.read(t))
			if err != nil {
				t.Fatal(err)
			}
			if fmt.Sprint(values[tc.key]) != fmt.Sprint(tc.want) {
				t.Fatalf("%s = %#v, want %#v", tc.key, values[tc.key], tc.want)
			}
			if !tc.wantWritten {
				return
			}
			row, _ := env.store.Get(env.rel)
			if row.Written[tc.key] != corpus.ValueHash(values[tc.key]) {
				t.Fatal("the merged value hash must be recorded so later kb appends keep working")
			}
			next := append(frontmatter.GetStringSlice(values, tc.key), "[[Epsilon]]")
			if again := env.apply(t, env.load(t), map[string]any{tc.key: next}); len(again.Written) != 1 {
				t.Fatalf("a later kb append = %#v", again)
			}
		})
	}
}
