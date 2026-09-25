package refs

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/frontmatter"
)

const (
	sourceX = `---
title: X
source_url: https://example.com/x
---
# X

Body of x.
`
	sourceY = `---
title: Y
---
Body of y.
`
	oldQuarantined = `---
title: Old
triage: quarantined
---
Old body.
`
	sourceIndex = `---
title: Source Index
---
# Sources

- [[raw/articles/x]] — X article
- [[raw/articles/x]] and [[raw/articles/y]]
- [X](../../raw/articles/x.md)
| [[x]] | note |
|---|---|
1. [[x]] and [[old]]
Plain paragraph about [[x]].
`
	article = `---
title: Article
sources:
  - "[[raw/articles/x]]"
  - "[[raw/articles/y]]"
related:
  - "[[x]]"
tags: [a, "[[x]]"]
---
Body mentions [[x|the X]] and [link](../../raw/articles/x.md).
` + "```\n[[x]] in code\n```\n"
	lockedArticle = `---
title: Locked
locked: true
sources:
  - "[[raw/articles/x]]"
---
Locked body.
`
)

type fixture struct {
	vault string
	topic string
}

func newFixture(t *testing.T) fixture {
	t.Helper()
	vault := t.TempDir()
	topic := filepath.Join(vault, "t")
	files := map[string]string{
		"CLAUDE.md":                          "# Topic\n",
		"raw/articles/x.md":                  sourceX,
		"raw/articles/y.md":                  sourceY,
		"raw/_quarantine/articles/old.md":    oldQuarantined,
		"wiki/index/Source Index.md":         sourceIndex,
		"wiki/concepts/Article.md":           article,
		"wiki/concepts/Locked.md":            lockedArticle,
		"wiki/concepts/Unrelated Thing.md":   "---\ntitle: Unrelated\n---\nNothing here.\n",
		"outputs/reports/2026-01-01-lint.md": "Report on [[x]].\n",
	}
	for rel, content := range files {
		writeFile(t, filepath.Join(topic, filepath.FromSlash(rel)), content)
	}
	return fixture{vault: vault, topic: topic}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readString(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func fixedNow() time.Time { return time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC) }

func TestScanFindsEveryReferenceKind(t *testing.T) {
	t.Parallel()
	fx := newFixture(t)

	refs, err := Scan(fx.vault, fx.topic, "raw/articles/x.md")
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	type key struct {
		file, kind, key string
		line            int
	}
	got := map[key]bool{}
	for _, ref := range refs {
		got[key{ref.File, ref.Kind, ref.Key, ref.Line}] = true
	}
	want := []key{
		{"wiki/index/Source Index.md", KindIndexLine, "", 6},
		{"wiki/index/Source Index.md", KindIndexLine, "", 7},
		{"wiki/index/Source Index.md", KindIndexLine, "", 8},
		{"wiki/index/Source Index.md", KindIndexLine, "", 9},
		{"wiki/index/Source Index.md", KindIndexLine, "", 11},
		{"wiki/index/Source Index.md", KindBody, "", 12},
		{"wiki/concepts/Article.md", KindFrontmatter, "sources", 4},
		{"wiki/concepts/Article.md", KindFrontmatter, "related", 7},
		{"wiki/concepts/Article.md", KindFrontmatter, "tags", 8},
		{"wiki/concepts/Article.md", KindBody, "", 10},
		{"wiki/concepts/Locked.md", KindFrontmatter, "sources", 5},
		{"outputs/reports/2026-01-01-lint.md", KindBody, "", 1},
	}
	for _, w := range want {
		if !got[w] {
			t.Errorf("missing ref %+v; got %+v", w, refs)
		}
	}
	bodyInArticle := 0
	for _, ref := range refs {
		if ref.File == "wiki/concepts/Article.md" && ref.Kind == KindBody {
			bodyInArticle++
		}
		if ref.File == "wiki/concepts/Article.md" && ref.Kind == KindFrontmatter && ref.Key == "sources" && strings.Contains(ref.Text, "raw/articles/y") {
			t.Errorf("reference to y reported as reference to x: %+v", ref)
		}
	}
	if bodyInArticle != 2 {
		t.Errorf("article body refs = %d, want 2 (code block excluded)", bodyInArticle)
	}
	if len(refs) != len(want)+1 {
		t.Errorf("refs = %d, want %d: %+v", len(refs), len(want)+1, refs)
	}
}

func TestScanRejectsUnknownTarget(t *testing.T) {
	t.Parallel()
	fx := newFixture(t)
	if _, err := Scan(fx.vault, fx.topic, "raw/articles/missing.md"); err == nil {
		t.Fatal("Scan of a missing target succeeded")
	}
}

func TestQuarantineAndRestoreRoundTrip(t *testing.T) {
	t.Parallel()
	fx := newFixture(t)
	indexPath := filepath.Join(fx.topic, "wiki/index/Source Index.md")
	articlePath := filepath.Join(fx.topic, "wiki/concepts/Article.md")
	lockedPath := filepath.Join(fx.topic, "wiki/concepts/Locked.md")
	originals := map[string]string{
		indexPath:   readString(t, indexPath),
		articlePath: readString(t, articlePath),
		lockedPath:  readString(t, lockedPath),
	}

	result, err := Quarantine(context.Background(), Options{
		VaultPath: fx.vault, TopicRoot: fx.topic, Path: "raw/articles/x.md", Reason: "off_topic", Now: fixedNow,
	})
	if err != nil {
		t.Fatalf("Quarantine: %v", err)
	}
	if result.NewPath != "raw/_quarantine/articles/x.md" || !strings.HasPrefix(result.ID, "q-") || len(result.ID) != 12 {
		t.Fatalf("result = %+v", result)
	}
	if _, err := os.Stat(filepath.Join(fx.topic, "raw/articles/x.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("original still present: %v", err)
	}
	moved := readString(t, filepath.Join(fx.topic, result.NewPath))
	values, _, err := frontmatter.Parse(moved)
	if err != nil {
		t.Fatal(err)
	}
	if values["triage"] != TriageQuarantined || values["triage_reason"] != "off_topic" {
		t.Fatalf("moved frontmatter = %v", values)
	}

	index := readString(t, indexPath)
	for _, gone := range []string{"- [[raw/articles/x]] — X article", "- [X](../../raw/articles/x.md)", "| [[x]] | note |", "1. [[x]] and [[old]]"} {
		if strings.Contains(index, gone) {
			t.Errorf("index still has %q", gone)
		}
	}
	for _, kept := range []string{"- [[raw/articles/x]] and [[raw/articles/y]]", "|---|---|", "Plain paragraph about [[x]]."} {
		if !strings.Contains(index, kept) {
			t.Errorf("index lost %q", kept)
		}
	}
	articleNow := readString(t, articlePath)
	if strings.Contains(articleNow, `- "[[raw/articles/x]]"`) || strings.Contains(articleNow, `- "[[x]]"`) {
		t.Errorf("article frontmatter items not removed:\n%s", articleNow)
	}
	if !strings.Contains(articleNow, `- "[[raw/articles/y]]"`) || !strings.Contains(articleNow, "Body mentions [[x|the X]]") {
		t.Errorf("article lost unrelated content:\n%s", articleNow)
	}
	if readString(t, lockedPath) != originals[lockedPath] {
		t.Error("locked document was edited")
	}

	edits := 0
	for _, row := range result.Removed {
		if row.Op == OpEdit {
			edits++
		}
	}
	if result.Removed[0].Op != OpMove || result.Removed[0].Original != "raw/articles/x.md" || edits != 6 {
		t.Fatalf("removed rows = %+v", result.Removed)
	}
	unremoved := map[string]bool{}
	for _, ref := range result.Unremoved {
		unremoved[ref.File+"#"+ref.Key] = true
	}
	if !unremoved["wiki/concepts/Article.md#tags"] || !unremoved["wiki/concepts/Locked.md#sources"] {
		t.Errorf("unremoved = %+v", result.Unremoved)
	}
	if len(result.BodyRefs) != 5 {
		t.Errorf("body refs = %+v, want 5 (mixed index line, index paragraph, two article links, report)", result.BodyRefs)
	}

	active, err := List(fx.topic)
	if err != nil || len(active) != 1 || active[0].ID != result.ID || active[0].Edits != 6 {
		t.Fatalf("List = %+v, %v", active, err)
	}

	restored, err := Restore(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, ID: result.ID, Now: fixedNow})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if len(restored.Manual) != 0 || len(restored.Restored) != 7 {
		t.Fatalf("restore = %+v", restored)
	}
	for path, want := range originals {
		if got := readString(t, path); got != want {
			t.Errorf("%s not byte-identical after restore:\n--- got\n%s\n--- want\n%s", path, got, want)
		}
	}
	back := readString(t, filepath.Join(fx.topic, "raw/articles/x.md"))
	if want := strings.Replace(sourceX, "source_url: https://example.com/x\n", "source_url: https://example.com/x\ntriage: kept\n", 1); back != want {
		t.Errorf("restored source differs beyond triage keys:\n%s", back)
	}
	if active, _ := List(fx.topic); len(active) != 0 {
		t.Errorf("quarantine still active after restore: %+v", active)
	}
}

func TestRestoreAfterEditReportsManual(t *testing.T) {
	t.Parallel()
	fx := newFixture(t)
	indexPath := filepath.Join(fx.topic, "wiki/index/Source Index.md")
	articlePath := filepath.Join(fx.topic, "wiki/concepts/Article.md")
	originalArticle := readString(t, articlePath)

	result, err := Quarantine(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, Path: "raw/articles/x.md", Reason: "thin"})
	if err != nil {
		t.Fatalf("Quarantine: %v", err)
	}
	edited := readString(t, indexPath) + "- a new line by the owner\n"
	writeFile(t, indexPath, edited)

	restored, err := Restore(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, Path: "raw/articles/x.md"})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if restored.ID != result.ID || len(restored.Manual) != 4 {
		t.Fatalf("manual = %+v", restored.Manual)
	}
	for _, entry := range restored.Manual {
		if entry.File != "wiki/index/Source Index.md" {
			t.Errorf("unexpected manual entry %+v", entry)
		}
	}
	if got := readString(t, indexPath); got != edited {
		t.Errorf("edited index was overwritten:\n%s", got)
	}
	if got := readString(t, articlePath); got != originalArticle {
		t.Errorf("article not restored:\n%s", got)
	}
}

func TestQuarantineRoundTripAfterUnterminatedLedger(t *testing.T) {
	t.Parallel()
	for _, seed := range []string{
		`{"quarantine_id":`,
		`{"quarantine_id":"prior","op":"restore","line_or_key":"move"}`,
	} {
		t.Run(seed, func(t *testing.T) {
			t.Parallel()
			fx := newFixture(t)
			writeFile(t, LedgerPath(fx.topic), seed)
			indexPath := filepath.Join(fx.topic, "wiki/index/Source Index.md")
			articlePath := filepath.Join(fx.topic, "wiki/concepts/Article.md")
			originalIndex := readString(t, indexPath)
			originalArticle := readString(t, articlePath)
			result, err := Quarantine(t.Context(), Options{
				VaultPath: fx.vault, TopicRoot: fx.topic, Path: "raw/articles/x.md", Reason: "thin", Now: fixedNow,
			})
			if err != nil {
				t.Fatal(err)
			}
			active, err := List(fx.topic)
			if err != nil || len(active) != 1 || active[0].ID != result.ID {
				t.Fatalf("quarantine lost after log reload: %+v, %v", active, err)
			}
			restored, err := Restore(t.Context(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, ID: result.ID})
			if err != nil || len(restored.Manual) != 0 {
				t.Fatalf("Restore = %+v, %v", restored, err)
			}
			if readString(t, indexPath) != originalIndex || readString(t, articlePath) != originalArticle {
				t.Fatal("restored references differ from original bytes")
			}
			if _, err := os.Stat(filepath.Join(fx.topic, "raw/articles/x.md")); err != nil {
				t.Fatalf("source was not restored: %v", err)
			}
		})
	}
}

func TestQuarantineRefusals(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		setup   func(t *testing.T, fx fixture)
		path    string
		wantErr error
	}{
		{
			name: "locked source",
			setup: func(t *testing.T, fx fixture) {
				writeFile(t, filepath.Join(fx.topic, "raw/articles/x.md"), "---\ntitle: X\nlocked: true\n---\nBody.\n")
			},
			path:    "raw/articles/x.md",
			wantErr: ErrLocked,
		},
		{
			name: "destination exists",
			setup: func(t *testing.T, fx fixture) {
				writeFile(t, filepath.Join(fx.topic, "raw/_quarantine/articles/x.md"), "taken\n")
			},
			path: "raw/articles/x.md",
		},
		{name: "not under raw", path: "wiki/concepts/Article.md"},
		{name: "already quarantined", path: "raw/_quarantine/articles/old.md"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fx := newFixture(t)
			if tc.setup != nil {
				tc.setup(t, fx)
			}
			before := readString(t, filepath.Join(fx.topic, "wiki/index/Source Index.md"))
			_, err := Quarantine(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, Path: tc.path, Reason: "off_topic"})
			if err == nil {
				t.Fatal("Quarantine succeeded")
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if got := readString(t, filepath.Join(fx.topic, "wiki/index/Source Index.md")); got != before {
				t.Error("index edited by a refused quarantine")
			}
			if _, err := os.Stat(LedgerPath(fx.topic)); !errors.Is(err, os.ErrNotExist) {
				t.Error("ledger written by a refused quarantine")
			}
		})
	}
}

func TestQuarantineKeepsKBOwnedRelationsOwned(t *testing.T) {
	t.Parallel()
	fx := newFixture(t)
	store, err := corpus.OpenState(fx.topic)
	if err != nil {
		t.Fatal(err)
	}
	writer := corpus.NewWriter(fx.topic, store, fixedNow)
	rel := "wiki/concepts/Unrelated Thing.md"
	doc, err := corpus.ReadDocument(fx.topic, rel, corpus.KindArticle)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Apply(doc, map[string]any{"related": []string{"[[x]]", "[[y]]"}}, corpus.StateMeta{}); err != nil {
		t.Fatal(err)
	}
	original := readString(t, filepath.Join(fx.topic, filepath.FromSlash(rel)))

	ownedHash := func() (string, string) {
		row, ok := store.Get(rel)
		if !ok {
			t.Fatal("no state row")
		}
		values, _, err := frontmatter.Parse(readString(t, filepath.Join(fx.topic, filepath.FromSlash(rel))))
		if err != nil {
			t.Fatal(err)
		}
		return row.Written["related"], corpus.ValueHash(values["related"])
	}

	result, err := Quarantine(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, Path: "raw/articles/x.md", Reason: "off_topic", Writer: writer, State: store})
	if err != nil {
		t.Fatal(err)
	}
	if written, current := ownedHash(); written != current {
		t.Fatalf("after quarantine written hash %s != current %s", written, current)
	}
	if _, err := Restore(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, ID: result.ID, Writer: writer, State: store}); err != nil {
		t.Fatal(err)
	}
	if written, current := ownedHash(); written != current {
		t.Fatalf("after restore written hash %s != current %s", written, current)
	}
	if got := readString(t, filepath.Join(fx.topic, filepath.FromSlash(rel))); got != original {
		t.Errorf("relation file not byte-identical:\n%s", got)
	}
}

// ledgerFault decides, for the n-th ledger append (1-based) of entry,
// whether the row still reaches the ledger and whether the append fails.
type ledgerFault func(n int, entry LedgerEntry) (write, fail bool)

func faultyAppend(fault ledgerFault) func(string, LedgerEntry) error {
	n := 0
	return func(topicRoot string, entry LedgerEntry) error {
		n++
		write, fail := fault(n, entry)
		if write {
			if err := appendLedger(topicRoot, entry); err != nil {
				return err
			}
		}
		if fail {
			return errors.New("injected ledger failure")
		}
		return nil
	}
}

// onEdit fails the k-th edit-row append, writing the row first when write.
func onEdit(k int, write bool) ledgerFault {
	edits := 0
	return func(_ int, entry LedgerEntry) (bool, bool) {
		if entry.Op != OpEdit {
			return true, false
		}
		edits++
		if edits == k {
			return write, true
		}
		return true, false
	}
}

// TestQuarantineRecordingFailureIsRecoverable: recovery information is
// durable before every mutation, so a quarantine whose ledger append fails
// leaves nothing moved without a ledger row and Restore brings every file
// back byte-identical afterwards.
func TestQuarantineRecordingFailureIsRecoverable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		// setup prepares the failure; it returns the appender for the
		// quarantine (nil = the real ledger file).
		setup func(t *testing.T, fx fixture) func(string, LedgerEntry) error
		// moved reports whether the source is expected in raw/_quarantine/
		// after the failed quarantine.
		moved bool
		// active reports whether a quarantine is left active to restore.
		active bool
		// retry makes the ledger writable again and runs the full cycle.
		retry func(t *testing.T, fx fixture)
	}{
		{
			name: "unwritable ledger fails before the move",
			setup: func(t *testing.T, fx fixture) func(string, LedgerEntry) error {
				if os.Geteuid() == 0 {
					t.Skip("root ignores file permissions")
				}
				writeFile(t, LedgerPath(fx.topic), "")
				if err := os.Chmod(LedgerPath(fx.topic), 0o444); err != nil {
					t.Fatal(err)
				}
				return nil
			},
			retry: func(t *testing.T, fx fixture) {
				if err := os.Chmod(LedgerPath(fx.topic), 0o644); err != nil {
					t.Fatal(err)
				}
				result, err := Quarantine(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, Path: "raw/articles/x.md", Reason: "off_topic"})
				if err != nil {
					t.Fatalf("retried Quarantine: %v", err)
				}
				if _, err := Restore(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, ID: result.ID}); err != nil {
					t.Fatalf("Restore after retry: %v", err)
				}
			},
		},
		{
			name: "move row recorded but the move and its rollback row fail",
			setup: func(*testing.T, fixture) func(string, LedgerEntry) error {
				return faultyAppend(func(_ int, entry LedgerEntry) (bool, bool) {
					switch entry.Op {
					case OpMove:
						return true, true // on disk, but reported as failed
					case OpRestore:
						return false, true
					}
					return true, false
				})
			},
			active: true,
		},
		{
			name:   "second reference edit cannot be recorded",
			setup:  func(*testing.T, fixture) func(string, LedgerEntry) error { return faultyAppend(onEdit(2, false)) },
			moved:  true,
			active: true,
		},
		{
			name:   "reference edit recorded but never performed",
			setup:  func(*testing.T, fixture) func(string, LedgerEntry) error { return faultyAppend(onEdit(1, true)) },
			moved:  true,
			active: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fx := newFixture(t)
			touched := []string{"wiki/index/Source Index.md", "wiki/concepts/Article.md", "wiki/concepts/Locked.md"}
			originals := map[string]string{}
			for _, rel := range touched {
				originals[rel] = readString(t, filepath.Join(fx.topic, filepath.FromSlash(rel)))
			}
			source := filepath.Join(fx.topic, "raw/articles/x.md")
			quarantined := filepath.Join(fx.topic, "raw/_quarantine/articles/x.md")

			opts := Options{VaultPath: fx.vault, TopicRoot: fx.topic, Path: "raw/articles/x.md", Reason: "off_topic", Now: fixedNow}
			opts.appendRow = tc.setup(t, fx)
			if _, err := Quarantine(context.Background(), opts); err == nil {
				t.Fatal("Quarantine succeeded despite the ledger failure")
			}

			_, movedErr := os.Stat(quarantined)
			if gotMoved := movedErr == nil; gotMoved != tc.moved {
				t.Fatalf("source moved = %v, want %v", gotMoved, tc.moved)
			}
			if !tc.moved {
				values, _, err := frontmatter.Parse(readString(t, source))
				if err != nil {
					t.Fatal(err)
				}
				if values["triage"] == TriageQuarantined {
					t.Errorf("unmoved source still marked quarantined: %v", values)
				}
			}
			active, err := List(fx.topic)
			if err != nil {
				t.Fatal(err)
			}
			if tc.moved && len(active) != 1 {
				t.Fatalf("moved source has no active ledger row: %+v", active)
			}
			if (len(active) == 1) != tc.active {
				t.Fatalf("active quarantines = %+v, want active=%v", active, tc.active)
			}

			if tc.active {
				restored, err := Restore(context.Background(), Options{VaultPath: fx.vault, TopicRoot: fx.topic, ID: active[0].ID, Now: fixedNow})
				if err != nil {
					t.Fatalf("Restore after failed quarantine: %v", err)
				}
				if len(restored.Manual) != 0 {
					t.Errorf("restore left manual repairs: %+v", restored.Manual)
				}
				if left, _ := List(fx.topic); len(left) != 0 {
					t.Errorf("quarantine still active after restore: %+v", left)
				}
			}
			if tc.retry != nil {
				tc.retry(t, fx)
			}

			if _, err := os.Stat(quarantined); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("source left in quarantine: %v", err)
			}
			if _, err := os.Stat(source); err != nil {
				t.Errorf("source not at its original path: %v", err)
			}
			for _, rel := range touched {
				if got := readString(t, filepath.Join(fx.topic, filepath.FromSlash(rel))); got != originals[rel] {
					t.Errorf("%s not byte-identical after recovery:\n%s", rel, got)
				}
			}
		})
	}
}

func TestIsIndexEntry(t *testing.T) {
	t.Parallel()
	testCases := map[string]bool{
		"- [[x]]":        true,
		"  * [[x]]":      true,
		"12. [[x]]":      true,
		"3) [[x]]":       true,
		"| [[x]] | a |":  true,
		"|---|:---:|":    false,
		"Text [[x]]":     false,
		"## [[x]]":       false,
		"2026 [[x]] was": false,
	}
	for line, want := range testCases {
		if got := isIndexEntry(line); got != want {
			t.Errorf("isIndexEntry(%q) = %v, want %v", line, got, want)
		}
	}
}
