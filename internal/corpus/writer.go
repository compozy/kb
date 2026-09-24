package corpus

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/resolve"
)

// Write statuses reported in WriteResult.Status.
const (
	StatusWritten        = "written"
	StatusUnchanged      = "unchanged"
	StatusSkippedLocked  = "skipped:locked"
	StatusSkippedChanged = "skipped:changed"
	// StatusSkippedUserKey is the run-summary label for WriteResult.SkippedUserKeys.
	StatusSkippedUserKey = "skipped:user-key"
)

// writeAttempts bounds the re-read-and-retry loop when the file changes
// between the read and the rename.
const writeAttempts = 2

// StateMeta is the classification context recorded in the state row.
type StateMeta struct {
	// Contract is the selection contract hash ("" keeps the row's value).
	Contract string
	// Banks maps question bank ids to versions; merged into the row. Each
	// listed bank is also recorded as judged on the document's body (the one
	// the write started from) under the contract.
	Banks map[string]string
	// Facts are per-document outcomes merged into the row; an empty value
	// deletes the fact.
	Facts map[string]string
}

// mergedListKeys are the relation lists kb only appends to (spec §9.3 "kb
// only adds"): a list the user owns is merged append-only instead of being
// skipped as a user key.
var mergedListKeys = append(slices.Clone(decisions.RelationKeys), "affects", "supersedes")

// IsMergedListKey reports whether kb merges new items into a user-owned
// value of key instead of leaving it alone (`aliases` and the relation
// lists).
func IsMergedListKey(key string) bool {
	return key == "aliases" || slices.Contains(mergedListKeys, key)
}

// WriteResult reports what Apply or ApplyBody did.
type WriteResult struct {
	// Written lists the keys whose value changed on disk, sorted.
	Written []string
	// SkippedUserKeys lists keys left alone because they belong to the user
	// (present and not last written by kb), sorted.
	SkippedUserKeys []string
	// BodyWritten reports whether ApplyBody replaced the body.
	BodyWritten bool
	// Status is StatusWritten, StatusUnchanged, StatusSkippedLocked or
	// StatusSkippedChanged.
	Status string
}

// Writer applies kb-owned frontmatter keys (spec §6) to topic documents,
// byte-preserving everything else, and records each write in the state
// store. It is safe for concurrent use; writes to one document are
// serialized.
type Writer struct {
	topicRoot string
	store     *StateStore
	now       func() time.Time

	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

// NewWriter returns a writer for the topic at topicRoot. now defaults to
// time.Now.
func NewWriter(topicRoot string, store *StateStore, now func() time.Time) *Writer {
	if now == nil {
		now = time.Now
	}

	return &Writer{
		topicRoot: filepath.Clean(topicRoot),
		store:     store,
		now:       now,
		locks:     make(map[string]*sync.Mutex),
	}
}

// Apply writes updates (owned key → value; a nil value deletes the key) to
// doc. Rules, in order:
//   - any non-owned key is refused with an error;
//   - a document with `locked: true` is left alone (skipped:locked);
//   - the file is re-read: if its body changed since doc was loaded nothing
//     is written (skipped:changed); other non-owned edits are fine because
//     only owned key blocks are replaced in the fresh content;
//   - a key is written only if it is absent or its current value hash equals
//     the state row's written hash (kb wrote it and nobody edited it since);
//     otherwise it belongs to the user and is reported in SkippedUserKeys.
//     `aliases` is merged instead (case-insensitive dedupe, never removes),
//     even when user-owned, and so is a user-owned relation list (`related`,
//     `extends`, `prerequisite`, `example_of`, `contradicts`, `affects`,
//     `supersedes`): new targets are appended, deduplicated by normalized
//     link target, and existing items are never removed or reordered;
//   - unchanged values are not rewritten;
//   - the file is replaced atomically (temp file + rename, mode kept), and
//     only then the state row is appended with the new body hash, contract,
//     merged banks and merged written hashes. The banks of meta, and every
//     key that now holds the requested value, are recorded as derived from
//     the body the write started from; judgments recorded earlier keep their
//     own body hash, so a write never makes another command's answers look
//     current. A body replacement that only adds or changes wikilink markup
//     (link insertion) carries every judgment made on the old body over to
//     the new one. If the file changed between the read and the rename, the
//     whole step is retried once, then reported as skipped:changed.
//
// On success doc is refreshed in place with the written content.
func (w *Writer) Apply(doc *Document, updates map[string]any, meta StateMeta) (WriteResult, error) {
	return w.apply(doc, nil, updates, meta)
}

// ApplyBody is Apply plus a body replacement (link insertion, recapture).
// The body is replaced only while the file's body hash still equals
// doc.BodyHash; the same key rules apply to updates.
func (w *Writer) ApplyBody(doc *Document, newBody string, updates map[string]any, meta StateMeta) (WriteResult, error) {
	return w.apply(doc, &newBody, updates, meta)
}

func (w *Writer) apply(doc *Document, newBody *string, updates map[string]any, meta StateMeta) (WriteResult, error) {
	if doc == nil {
		return WriteResult{}, errors.New("corpus: document is required")
	}
	for key := range updates {
		if !decisions.IsOwnedKey(key) {
			return WriteResult{}, fmt.Errorf("corpus: refusing to write non-owned key %q to %s", key, doc.Path)
		}
	}

	unlock := w.lock(doc.Path)
	defer unlock()

	if doc.Locked() {
		return WriteResult{Status: StatusSkippedLocked}, nil
	}

	absolute := doc.AbsPath
	if absolute == "" {
		absolute = filepath.Join(w.topicRoot, filepath.FromSlash(doc.Path))
	}

	for range writeAttempts {
		result, retry, err := w.attempt(doc, absolute, newBody, updates, meta)
		if err != nil || !retry {
			return result, err
		}
	}

	return WriteResult{Status: StatusSkippedChanged}, nil
}

// attempt performs one read-plan-write cycle. retry is true when the file
// changed between the read and the rename.
func (w *Writer) attempt(doc *Document, absolute string, newBody *string, updates map[string]any, meta StateMeta) (WriteResult, bool, error) {
	content, err := os.ReadFile(absolute)
	if err != nil {
		return WriteResult{}, false, fmt.Errorf("corpus: read %s: %w", doc.Path, err)
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return WriteResult{}, false, fmt.Errorf("corpus: stat %s: %w", doc.Path, err)
	}

	current := string(content)
	values, body, err := frontmatter.Parse(current)
	if err != nil {
		return WriteResult{}, false, fmt.Errorf("corpus: parse %s: %w", doc.Path, err)
	}
	if frontmatter.GetBool(values, "locked") {
		return WriteResult{Status: StatusSkippedLocked}, false, nil
	}
	if BodyHash(body) != doc.BodyHash {
		return WriteResult{Status: StatusSkippedChanged}, false, nil
	}

	row, _ := w.store.Lookup(doc.Path, doc.BodyHash, w.exists)
	plan := planKeys(values, updates, row)
	result := WriteResult{Written: plan.written, SkippedUserKeys: plan.skipped}

	edited := current
	if len(plan.edits) > 0 {
		edited, err = frontmatter.EditKeys(current, plan.edits)
		if err != nil {
			return WriteResult{}, false, fmt.Errorf("corpus: edit %s: %w", doc.Path, err)
		}
	}
	if newBody != nil && *newBody != body {
		edited = edited[:len(edited)-len(body)] + *newBody
		result.BodyWritten = true
	}

	judged := BodyHash(body)
	if edited == current {
		result.Status = StatusUnchanged
		update := stateUpdate{values: values, previous: row, judged: judged, now: judged, carry: true, meta: meta, fresh: plan.fresh}
		if err := w.recordState(doc.Path, update); err != nil {
			return WriteResult{}, false, err
		}
		return result, false, nil
	}

	latest, err := os.ReadFile(absolute)
	if err != nil {
		return WriteResult{}, false, fmt.Errorf("corpus: re-read %s: %w", doc.Path, err)
	}
	if string(latest) != current {
		return WriteResult{}, true, nil
	}
	if err := writeFileAtomic(absolute, []byte(edited), info.Mode().Perm()); err != nil {
		return WriteResult{}, false, fmt.Errorf("corpus: write %s: %w", doc.Path, err)
	}

	writtenValues, writtenBody, err := frontmatter.Parse(edited)
	if err != nil {
		return WriteResult{}, false, fmt.Errorf("corpus: parse written %s: %w", doc.Path, err)
	}
	update := stateUpdate{
		values: writtenValues, previous: row, judged: judged, now: BodyHash(writtenBody),
		carry: linkMarkupOnly(body, writtenBody), meta: meta, keys: plan.written, fresh: plan.fresh, force: true,
	}
	if err := w.recordState(doc.Path, update); err != nil {
		return WriteResult{}, false, err
	}

	doc.Raw = edited
	doc.Body = writtenBody
	doc.BodyHash = BodyHash(writtenBody)
	doc.Frontmatter = writtenValues
	doc.Aliases = resolve.Aliases(writtenValues)
	if refreshed, err := os.Stat(absolute); err == nil {
		doc.ModTime = refreshed.ModTime()
	}

	result.Status = StatusWritten
	return result, false, nil
}

// stateUpdate is what one write records in the state row.
type stateUpdate struct {
	// values are the frontmatter values as parsed back from the file.
	values   map[string]any
	previous *StateRow
	// judged is the body hash the write started from (the body the caller's
	// answers and values were computed on); now is the body hash on disk
	// after the write.
	judged, now string
	// carry is set when now differs from judged by wikilink markup only:
	// judgments made on judged then hold for now as well.
	carry bool
	meta  StateMeta
	// keys are the owned keys just written or deleted.
	keys []string
	// fresh are the keys that now hold the value the caller asked for.
	fresh []string
	// force appends a row even when nothing changed.
	force bool
}

// recordState appends the state row for path. Unless update.force is set,
// nothing is appended when the stored row already holds the same values.
// Value hashes of update.keys come from update.values.
func (w *Writer) recordState(path string, update stateUpdate) error {
	row := newStateRow(path, update.previous)

	stamp := update.judged
	if update.carry {
		stamp = update.now
		for _, hashes := range []map[string]string{row.BankBody, row.WrittenBody} {
			for id, hash := range hashes {
				if hash == update.judged {
					hashes[id] = update.now
				}
			}
		}
	}
	row.BodyHash = update.now
	if update.meta.Contract != "" {
		row.Contract = update.meta.Contract
	}
	for bank, version := range update.meta.Banks {
		row.Banks[bank] = version
		row.BankBody[bank] = stamp
		row.BankContract[bank] = row.Contract
	}
	for name, value := range update.meta.Facts {
		if value == "" {
			delete(row.Facts, name)
			continue
		}
		row.Facts[name] = value
	}
	for _, key := range update.keys {
		value, present := update.values[key]
		if !present {
			delete(row.Written, key)
			delete(row.WrittenBody, key)
			continue
		}
		row.Written[key] = ValueHash(value)
	}
	for _, key := range update.fresh {
		row.WrittenBody[key] = stamp
	}

	if !update.force {
		if stored, ok := w.store.Get(path); ok && sameState(*stored, row) {
			return nil
		}
	}

	row.Updated = w.now().UTC().Format(time.RFC3339)
	if err := w.store.Put(row); err != nil {
		return fmt.Errorf("corpus: record state for %s: %w", path, err)
	}

	return nil
}

// newStateRow starts the next row of path from previous. A previous row
// that predates per-bank tracking judged every bank and key on its body hash
// under its contract: that is pinned before BodyHash moves.
func newStateRow(path string, previous *StateRow) StateRow {
	row := StateRow{
		Path: path, Banks: map[string]string{}, Written: map[string]string{},
		BankBody: map[string]string{}, BankContract: map[string]string{}, WrittenBody: map[string]string{}, Facts: map[string]string{},
	}
	if previous == nil {
		return row
	}
	row.Contract = previous.Contract
	maps.Copy(row.Banks, previous.Banks)
	maps.Copy(row.Written, previous.Written)
	maps.Copy(row.BankBody, previous.BankBody)
	maps.Copy(row.BankContract, previous.BankContract)
	maps.Copy(row.WrittenBody, previous.WrittenBody)
	maps.Copy(row.Facts, previous.Facts)
	for bank := range previous.Banks {
		if _, ok := row.BankBody[bank]; !ok {
			row.BankBody[bank] = previous.BodyHash
		}
		if _, ok := row.BankContract[bank]; !ok {
			row.BankContract[bank] = previous.Contract
		}
	}
	for key := range previous.Written {
		if _, ok := row.WrittenBody[key]; !ok {
			row.WrittenBody[key] = previous.BodyHash
		}
	}
	return row
}

// sameState reports whether two rows hold the same bookkeeping (Updated
// aside; nil and empty maps are equal).
func sameState(a, b StateRow) bool {
	return a.BodyHash == b.BodyHash && a.Contract == b.Contract &&
		maps.Equal(a.Banks, b.Banks) && maps.Equal(a.Written, b.Written) &&
		maps.Equal(a.BankBody, b.BankBody) && maps.Equal(a.BankContract, b.BankContract) &&
		maps.Equal(a.WrittenBody, b.WrittenBody) && maps.Equal(a.Facts, b.Facts)
}

// linkMarkupOnly reports whether after differs from before only by wikilink
// markup (`[[target|text]]` around text that was already there), the one
// kind of body edit kb makes that leaves the document's content unchanged.
func linkMarkupOnly(before, after string) bool {
	return before == after || flattenWikilinks(before) == flattenWikilinks(after)
}

func (w *Writer) exists(path string) bool {
	_, err := os.Stat(filepath.Join(w.topicRoot, filepath.FromSlash(path)))
	return !errors.Is(err, fs.ErrNotExist)
}

func (w *Writer) lock(path string) func() {
	w.mu.Lock()
	lock, ok := w.locks[path]
	if !ok {
		lock = &sync.Mutex{}
		w.locks[path] = lock
	}
	w.mu.Unlock()

	lock.Lock()
	return lock.Unlock
}

// keyPlan is the per-key outcome of the ownership rules.
type keyPlan struct {
	edits   []frontmatter.KeyUpdate
	written []string
	skipped []string
	// fresh lists the keys that hold the requested value after the write
	// (written, merged or already equal).
	fresh []string
}

func planKeys(values map[string]any, updates map[string]any, row *StateRow) keyPlan {
	plan := keyPlan{}
	for _, key := range slices.Sorted(maps.Keys(updates)) {
		desired := updates[key]
		current, present := values[key]
		present = present && !isEmptyValue(current)

		if key == "aliases" && desired != nil {
			plan.merge(key, current, present, mergeAliases(current, desired))
			continue
		}

		switch {
		case !present && desired == nil:
			continue
		case !present:
			plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Value: desired})
			plan.written = append(plan.written, key)
			plan.fresh = append(plan.fresh, key)
			continue
		case desired != nil && ValueHash(current) == ValueHash(desired):
			plan.fresh = append(plan.fresh, key)
			continue
		case row == nil || row.Written[key] != ValueHash(current):
			if merged, ok := mergeLinkList(key, current, desired); ok {
				plan.merge(key, current, present, merged)
				continue
			}
			plan.skipped = append(plan.skipped, key)
			continue
		}

		if desired == nil {
			plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Delete: true})
		} else {
			plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Value: desired})
			plan.fresh = append(plan.fresh, key)
		}
		plan.written = append(plan.written, key)
	}

	return plan
}

// merge plans writing the merged value of key unless it is already current.
func (plan *keyPlan) merge(key string, current any, present bool, merged []string) {
	plan.fresh = append(plan.fresh, key)
	if present && ValueHash(current) == ValueHash(merged) {
		return
	}
	plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Value: merged})
	plan.written = append(plan.written, key)
}

// mergeLinkList merges desired into a user-owned relation list: every
// current item stays in place and order, and each desired item whose
// normalized link target (case-insensitive) is not listed yet is appended.
// It reports false for other keys, a deletion, or a current value that is
// not a list of strings (kb never rewrites a shape it does not know).
func mergeLinkList(key string, current, desired any) ([]string, bool) {
	if desired == nil || !slices.Contains(mergedListKeys, key) {
		return nil, false
	}
	items, ok := stringList(current)
	if !ok {
		return nil, false
	}
	merged := slices.Clone(items)
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		seen[linkKey(item)] = struct{}{}
	}
	for _, item := range stringValues(desired) {
		k := linkKey(item)
		if k == "" {
			continue
		}
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		merged = append(merged, strings.TrimSpace(item))
	}
	return merged, true
}

// linkKey is the dedupe key of a relation list item: its normalized link
// target, case-insensitive.
func linkKey(item string) string {
	return strings.ToLower(resolve.Normalize(item))
}

// stringList returns a YAML list whose every item is a string.
func stringList(value any) ([]string, bool) {
	switch typed := value.(type) {
	case []string:
		return typed, true
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil, false
			}
			items = append(items, text)
		}
		return items, true
	default:
		return nil, false
	}
}

// mergeAliases keeps every existing alias in order and appends new ones not
// already present (case-insensitive).
func mergeAliases(current, desired any) []string {
	merged := make([]string, 0)
	seen := make(map[string]struct{})
	add := func(alias string) {
		trimmed := strings.TrimSpace(alias)
		if trimmed == "" {
			return
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		merged = append(merged, trimmed)
	}
	for _, alias := range stringValues(current) {
		add(alias)
	}
	for _, alias := range stringValues(desired) {
		add(alias)
	}

	return merged
}

func stringValues(value any) []string {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case []string:
		return typed
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				values = append(values, text)
			}
		}
		return values
	default:
		return nil
	}
}

func isEmptyValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(typed) == ""
	default:
		return false
	}
}
