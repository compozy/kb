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
	// Banks maps question bank ids to versions; merged into the row.
	Banks map[string]string
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
//     even when user-owned;
//   - unchanged values are not rewritten;
//   - the file is replaced atomically (temp file + rename, mode kept), and
//     only then the state row is appended with the new body hash, contract,
//     merged banks and merged written hashes. If the file changed between
//     the read and the rename, the whole step is retried once, then reported
//     as skipped:changed.
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

	if edited == current {
		result.Status = StatusUnchanged
		if err := w.recordState(doc.Path, values, row, BodyHash(body), meta, nil, false); err != nil {
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
	if err := w.recordState(doc.Path, writtenValues, row, BodyHash(writtenBody), meta, plan.written, true); err != nil {
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

// recordState appends the state row for path. Unless force is set, nothing is
// appended when the stored row already holds the same values. keys are the
// owned keys just written or deleted; their value hashes come from values as
// parsed back from the written file.
func (w *Writer) recordState(path string, values map[string]any, previous *StateRow, bodyHash string, meta StateMeta, keys []string, force bool) error {
	row := StateRow{Path: path, Banks: map[string]string{}, Written: map[string]string{}}
	if previous != nil {
		row.Contract = previous.Contract
		maps.Copy(row.Banks, previous.Banks)
		maps.Copy(row.Written, previous.Written)
	}
	row.BodyHash = bodyHash
	if meta.Contract != "" {
		row.Contract = meta.Contract
	}
	maps.Copy(row.Banks, meta.Banks)
	for _, key := range keys {
		value, present := values[key]
		if !present {
			delete(row.Written, key)
			continue
		}
		row.Written[key] = ValueHash(value)
	}

	if !force {
		if stored, ok := w.store.Get(path); ok && stored.BodyHash == row.BodyHash && stored.Contract == row.Contract && maps.Equal(stored.Banks, row.Banks) && maps.Equal(stored.Written, row.Written) {
			return nil
		}
	}

	row.Updated = w.now().UTC().Format(time.RFC3339)
	if err := w.store.Put(row); err != nil {
		return fmt.Errorf("corpus: record state for %s: %w", path, err)
	}

	return nil
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
}

func planKeys(values map[string]any, updates map[string]any, row *StateRow) keyPlan {
	plan := keyPlan{}
	for _, key := range slices.Sorted(maps.Keys(updates)) {
		desired := updates[key]
		current, present := values[key]
		present = present && !isEmptyValue(current)

		if key == "aliases" && desired != nil {
			merged := mergeAliases(current, desired)
			if present && ValueHash(current) == ValueHash(merged) {
				continue
			}
			plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Value: merged})
			plan.written = append(plan.written, key)
			continue
		}

		switch {
		case !present && desired == nil:
			continue
		case !present:
			plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Value: desired})
			plan.written = append(plan.written, key)
			continue
		case desired != nil && ValueHash(current) == ValueHash(desired):
			continue
		case row == nil || row.Written[key] != ValueHash(current):
			plan.skipped = append(plan.skipped, key)
			continue
		}

		if desired == nil {
			plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Delete: true})
		} else {
			plan.edits = append(plan.edits, frontmatter.KeyUpdate{Key: key, Value: desired})
		}
		plan.written = append(plan.written, key)
	}

	return plan
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
