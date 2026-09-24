package refs

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/resolve"
)

// LedgerFile is the quarantine ledger inside <topic>/.decisions/.
const LedgerFile = "quarantine-ledger.jsonl"

// Ledger operations.
const (
	OpMove    = "move"
	OpEdit    = "edit"
	OpRestore = "restore"
)

// Ledger line_or_key markers.
const (
	// MarkerIndex marks an index-line removal.
	MarkerIndex = "index"
	// MarkerMove marks the file move (and, on restore rows, the move back).
	MarkerMove = "move"
)

// Triage values written on quarantine and restore (spec §6).
const (
	TriageQuarantined = "quarantined"
	TriageKept        = "kept"
)

// ErrLocked is returned when the file to quarantine carries `locked: true`:
// kb never moves or edits a locked file.
var ErrLocked = errors.New("refs: file is locked (locked: true); not quarantined")

// LedgerEntry is one row of the quarantine ledger. Edit rows record the
// exact bytes removed (Before) and written instead (After, "" for a pure
// deletion) at 1-based Line, plus the whole-file sha256 before and after the
// edit. The move row records the quarantine path in File, the original path
// in Original, and the file hash before and after the triage keys were set.
type LedgerEntry struct {
	QuarantineID string `json:"quarantine_id"`
	Time         string `json:"time"`
	Op           string `json:"op"`
	File         string `json:"file"`
	Original     string `json:"original,omitempty"`
	LineOrKey    string `json:"line_or_key"`
	Line         int    `json:"line,omitempty"`
	Before       string `json:"before,omitempty"`
	After        string `json:"after,omitempty"`
	HashBefore   string `json:"hash_before,omitempty"`
	HashAfter    string `json:"hash_after,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

// Options configures Quarantine and Restore.
type Options struct {
	// VaultPath is the vault root (the parent of the topic directories).
	VaultPath string
	// TopicRoot is the topic directory.
	TopicRoot string
	// Path is the topic-relative file: `raw/...` for Quarantine; the original
	// or the quarantine path for Restore (used when ID is empty).
	Path string
	// Reason is the triage_reason written on quarantine (off_topic, thin, ...).
	Reason string
	// ID is the quarantine id; Quarantine derives one when empty.
	ID string
	// Writer writes the triage keys; one is opened on the topic when nil.
	Writer *corpus.Writer
	// State, when set, keeps the state row's written hashes in step with the
	// removal of items from kb-owned list keys, so a kb-written relation
	// stays kb-owned after quarantine and restore.
	State *corpus.StateStore
	// Now defaults to time.Now.
	Now func() time.Time
}

// Result reports a quarantine.
type Result struct {
	ID string
	// NewPath is the topic-relative quarantine path.
	NewPath string
	// Removed lists the ledger rows of every removal (move row first).
	Removed []LedgerEntry
	// BodyRefs are references left in place: body links and index lines that
	// also link to files that are not quarantined. Lint reports them as
	// link-to-quarantined.
	BodyRefs []Ref
	// Unremoved are frontmatter references that could not be removed line by
	// line (flow-style or scalar values, keys outside RemovableKeys, locked
	// documents); they are left for the owner.
	Unremoved []Ref
	// SkippedUserKeys lists triage keys the writer left alone because the
	// user owns them.
	SkippedUserKeys []string
}

// RestoreResult reports a restore.
type RestoreResult struct {
	ID string
	// Restored lists the ledger rows that were replayed.
	Restored []LedgerEntry
	// Manual lists edit rows that were not replayed because the touched file
	// changed since the quarantine; they are ready to print for manual repair.
	Manual []LedgerEntry
	// Path is the topic-relative path the file was moved back to.
	Path string
}

// Quarantined is one active quarantine.
type Quarantined struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Original string `json:"original"`
	Reason   string `json:"reason"`
	Time     string `json:"time"`
	Edits    int    `json:"edits"`
}

// Quarantine moves a source to raw/_quarantine/ and removes the references
// that would otherwise dangle (spec §7.1):
//  1. `triage: quarantined` and `triage_reason` are written through the
//     corpus writer (a locked file returns ErrLocked and nothing changes);
//  2. the file is moved to raw/_quarantine/<path without raw/> (refused when
//     the destination exists);
//  3. index lines in wiki/index/ whose every link points to quarantined files
//     are removed;
//  4. items of RemovableKeys frontmatter lists pointing to the file are
//     removed from other topic documents, deleting exactly their lines;
//  5. body links are left in place and returned in BodyRefs.
//
// Every change is appended to the quarantine ledger.
func Quarantine(ctx context.Context, opts Options) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	now := nowFunc(opts.Now)
	rel := cleanRel(opts.Path)
	newRel, ok := resolve.QuarantinePath(rel)
	if !ok {
		return Result{}, fmt.Errorf("refs: %q is not a raw/ source outside quarantine", opts.Path)
	}
	source := joinTopic(opts.TopicRoot, rel)
	dest := joinTopic(opts.TopicRoot, newRel)
	if _, err := os.Stat(dest); err == nil {
		return Result{}, fmt.Errorf("refs: quarantine destination %s already exists", newRel)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return Result{}, fmt.Errorf("refs: stat %s: %w", newRel, err)
	}

	doc, err := corpus.ReadDocument(opts.TopicRoot, rel, corpus.KindSource)
	if err != nil {
		return Result{}, fmt.Errorf("refs: %w", err)
	}
	if doc.Locked() {
		return Result{}, fmt.Errorf("%w: %s", ErrLocked, rel)
	}

	scans, err := scanTopic(opts.VaultPath, opts.TopicRoot, rel)
	if err != nil {
		return Result{}, err
	}

	writer, err := topicWriter(opts)
	if err != nil {
		return Result{}, err
	}
	hashBefore := hashString(doc.Raw)
	written, err := writer.Apply(doc, map[string]any{"triage": TriageQuarantined, "triage_reason": opts.Reason}, corpus.StateMeta{})
	if err != nil {
		return Result{}, fmt.Errorf("refs: mark %s quarantined: %w", rel, err)
	}
	switch written.Status {
	case corpus.StatusSkippedLocked:
		return Result{}, fmt.Errorf("%w: %s", ErrLocked, rel)
	case corpus.StatusSkippedChanged:
		return Result{}, fmt.Errorf("refs: %s changed while it was being quarantined; retry", rel)
	}
	moved, err := readFile(source)
	if err != nil {
		return Result{}, err
	}

	stamp := now().UTC().Format(time.RFC3339Nano)
	id := strings.TrimSpace(opts.ID)
	if id == "" {
		id = quarantineID(rel, stamp)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return Result{}, fmt.Errorf("refs: create %s: %w", filepath.Dir(newRel), err)
	}
	if err := os.Rename(source, dest); err != nil {
		return Result{}, fmt.Errorf("refs: move %s to %s: %w", rel, newRel, err)
	}

	result := Result{ID: id, NewPath: newRel, SkippedUserKeys: written.SkippedUserKeys}
	moveRow := LedgerEntry{
		QuarantineID: id, Time: stamp, Op: OpMove, File: newRel, Original: rel, LineOrKey: MarkerMove,
		HashBefore: hashBefore, HashAfter: hashString(moved), Reason: opts.Reason,
	}
	if err := appendLedger(opts.TopicRoot, moveRow); err != nil {
		return result, err
	}
	result.Removed = append(result.Removed, moveRow)

	for _, scan := range scans {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		rows, err := removeSpans(opts, scan, id, stamp)
		result.Removed = append(result.Removed, rows...)
		if err != nil {
			return result, err
		}
		for _, sp := range scan.spans {
			switch {
			case sp.removable:
			case sp.ref.Kind == KindFrontmatter:
				result.Unremoved = append(result.Unremoved, sp.ref)
			default:
				result.BodyRefs = append(result.BodyRefs, sp.ref)
			}
		}
	}
	return result, nil
}

// removeSpans deletes the removable spans of one file bottom-up (so every
// recorded line number is valid in the file state right before its edit) and
// appends one ledger row per deletion.
func removeSpans(opts Options, scan fileScan, id, stamp string) ([]LedgerEntry, error) {
	removable := make([]span, 0, len(scan.spans))
	for _, sp := range scan.spans {
		if sp.removable {
			removable = append(removable, sp)
		}
	}
	if len(removable) == 0 {
		return nil, nil
	}
	sort.Slice(removable, func(i, j int) bool { return removable[i].start > removable[j].start })

	current, err := readFile(scan.abs)
	if err != nil {
		return nil, err
	}
	if current != scan.content {
		return nil, fmt.Errorf("refs: %s changed during quarantine; its references were left in place", scan.rel)
	}

	rows := make([]LedgerEntry, 0, len(removable))
	for _, sp := range removable {
		lines := splitLines(current)
		before := strings.Join(lines[sp.start:sp.end], "")
		edited := strings.Join(lines[:sp.start], "") + strings.Join(lines[sp.end:], "")
		if err := writeAtomic(scan.abs, edited); err != nil {
			return rows, err
		}
		key := MarkerIndex
		if sp.ref.Kind == KindFrontmatter {
			key = sp.ref.Key
			if err := syncWritten(opts, scan.rel, key, current, edited); err != nil {
				return rows, err
			}
		}
		row := LedgerEntry{
			QuarantineID: id, Time: stamp, Op: OpEdit, File: scan.rel, LineOrKey: key, Line: sp.start + 1,
			Before: before, HashBefore: hashString(current), HashAfter: hashString(edited), Reason: opts.Reason,
		}
		if err := appendLedger(opts.TopicRoot, row); err != nil {
			return rows, err
		}
		rows = append(rows, row)
		current = edited
	}
	return rows, nil
}

// Restore replays the ledger of one quarantine in reverse: every removed
// span is put back at its recorded line when the touched file still hashes
// to the recorded hash_after (otherwise the row is returned in Manual and the
// file is not touched), the file is moved back, `triage: kept` is written and
// `triage_reason` removed through the corpus writer, and `restore` rows are
// appended. After a quarantine and restore with no edits in between, every
// touched index and frontmatter file is byte-identical to the original; the
// restored file itself differs only in its triage keys (`triage: kept`
// instead of no key or another value, `triage_reason` removed).
func Restore(ctx context.Context, opts Options) (RestoreResult, error) {
	if err := ctx.Err(); err != nil {
		return RestoreResult{}, err
	}
	now := nowFunc(opts.Now)
	entries, err := LoadLedger(opts.TopicRoot)
	if err != nil {
		return RestoreResult{}, err
	}
	move, ok := findActive(entries, strings.TrimSpace(opts.ID), cleanRel(opts.Path))
	if !ok {
		return RestoreResult{}, fmt.Errorf("refs: no active quarantine matches id %q path %q", opts.ID, opts.Path)
	}
	source := joinTopic(opts.TopicRoot, move.File)
	dest := joinTopic(opts.TopicRoot, move.Original)
	if _, err := os.Stat(source); err != nil {
		return RestoreResult{}, fmt.Errorf("refs: quarantined file %s: %w", move.File, err)
	}
	if _, err := os.Stat(dest); err == nil {
		return RestoreResult{}, fmt.Errorf("refs: cannot restore %s: %s already exists", move.File, move.Original)
	}

	stamp := now().UTC().Format(time.RFC3339Nano)
	result := RestoreResult{ID: move.QuarantineID, Path: move.Original}
	edits := make([]LedgerEntry, 0)
	for _, entry := range entries {
		if entry.QuarantineID == move.QuarantineID && entry.Op == OpEdit {
			edits = append(edits, entry)
		}
	}
	for _, entry := range slices.Backward(edits) {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		applied, err := replay(opts, entry, stamp)
		if err != nil {
			return result, err
		}
		if applied {
			result.Restored = append(result.Restored, entry)
		} else {
			result.Manual = append(result.Manual, entry)
		}
	}

	if err := os.Rename(source, dest); err != nil {
		return result, fmt.Errorf("refs: move %s back to %s: %w", move.File, move.Original, err)
	}
	writer, err := topicWriter(opts)
	if err != nil {
		return result, err
	}
	doc, err := corpus.ReadDocument(opts.TopicRoot, move.Original, corpus.KindSource)
	if err != nil {
		return result, fmt.Errorf("refs: %w", err)
	}
	if _, err := writer.Apply(doc, map[string]any{"triage": TriageKept, "triage_reason": nil}, corpus.StateMeta{}); err != nil {
		return result, fmt.Errorf("refs: mark %s kept: %w", move.Original, err)
	}
	row := LedgerEntry{
		QuarantineID: move.QuarantineID, Time: stamp, Op: OpRestore, File: move.Original, Original: move.File,
		LineOrKey: MarkerMove, Reason: move.Reason,
	}
	if err := appendLedger(opts.TopicRoot, row); err != nil {
		return result, err
	}
	result.Restored = append(result.Restored, move)
	return result, nil
}

// replay puts one removed span back. It returns false (without touching the
// file) when the file no longer hashes to the row's hash_after.
func replay(opts Options, entry LedgerEntry, stamp string) (bool, error) {
	abs := joinTopic(opts.TopicRoot, entry.File)
	current, err := readFile(abs)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if hashString(current) != entry.HashAfter {
		return false, nil
	}
	lines := splitLines(current)
	start := entry.Line - 1
	afterLines := splitLines(entry.After)
	end := start + len(afterLines)
	if start < 0 || end > len(lines) || strings.Join(lines[start:end], "") != entry.After {
		return false, nil
	}
	restored := strings.Join(lines[:start], "") + entry.Before + strings.Join(lines[end:], "")
	if err := writeAtomic(abs, restored); err != nil {
		return false, err
	}
	if entry.LineOrKey != MarkerIndex {
		if err := syncWritten(opts, entry.File, entry.LineOrKey, current, restored); err != nil {
			return true, err
		}
	}
	row := LedgerEntry{
		QuarantineID: entry.QuarantineID, Time: stamp, Op: OpRestore, File: entry.File, LineOrKey: entry.LineOrKey,
		Line: entry.Line, Before: entry.After, After: entry.Before,
		HashBefore: hashString(current), HashAfter: hashString(restored), Reason: entry.Reason,
	}
	return true, appendLedger(opts.TopicRoot, row)
}

// List returns the active quarantines of a topic (move rows without a
// matching restore of the move), oldest first.
func List(topicRoot string) ([]Quarantined, error) {
	entries, err := LoadLedger(topicRoot)
	if err != nil {
		return nil, err
	}
	restored := map[string]bool{}
	edits := map[string]int{}
	for _, entry := range entries {
		switch {
		case entry.Op == OpRestore && entry.LineOrKey == MarkerMove:
			restored[entry.QuarantineID] = true
		case entry.Op == OpEdit:
			edits[entry.QuarantineID]++
		}
	}
	out := make([]Quarantined, 0)
	for _, entry := range entries {
		if entry.Op != OpMove || restored[entry.QuarantineID] {
			continue
		}
		out = append(out, Quarantined{
			ID: entry.QuarantineID, Path: entry.File, Original: entry.Original,
			Reason: entry.Reason, Time: entry.Time, Edits: edits[entry.QuarantineID],
		})
	}
	return out, nil
}

// LoadLedger reads every well-formed ledger row in file order (a missing
// ledger is empty).
func LoadLedger(topicRoot string) ([]LedgerEntry, error) {
	file, err := os.Open(LedgerPath(topicRoot))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("refs: open ledger: %w", err)
	}
	defer func() { _ = file.Close() }()
	entries := make([]LedgerEntry, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var entry LedgerEntry
		if err := json.Unmarshal(line, &entry); err != nil || entry.QuarantineID == "" {
			continue
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("refs: read ledger: %w", err)
	}
	return entries, nil
}

// LedgerPath returns <topic>/.decisions/quarantine-ledger.jsonl.
func LedgerPath(topicRoot string) string {
	return filepath.Join(topicRoot, decisions.ReceiptsDir, LedgerFile)
}

// findActive returns the move row of the active quarantine matching id, or
// (when id is empty) whose original or quarantine path is rel.
func findActive(entries []LedgerEntry, id, rel string) (LedgerEntry, bool) {
	restored := map[string]bool{}
	for _, entry := range entries {
		if entry.Op == OpRestore && entry.LineOrKey == MarkerMove {
			restored[entry.QuarantineID] = true
		}
	}
	for _, entry := range slices.Backward(entries) {
		if entry.Op != OpMove || restored[entry.QuarantineID] {
			continue
		}
		if id != "" && entry.QuarantineID == id {
			return entry, true
		}
		if id == "" && rel != "" && (entry.Original == rel || entry.File == rel) {
			return entry, true
		}
	}
	return LedgerEntry{}, false
}

// syncWritten keeps state.written[key] equal to the value hash of an owned
// key after kb removed or restored one of its items, but only when kb owned
// the value before the edit.
func syncWritten(opts Options, rel, key, before, after string) error {
	if opts.State == nil || !decisions.IsOwnedKey(key) {
		return nil
	}
	row, ok := opts.State.Get(rel)
	if !ok || row.Written == nil {
		return nil
	}
	oldValues, _, err := frontmatter.Parse(before)
	if err != nil {
		return nil
	}
	if row.Written[key] != corpus.ValueHash(oldValues[key]) {
		return nil
	}
	newValues, _, err := frontmatter.Parse(after)
	if err != nil {
		return nil
	}
	updated := *row
	updated.Written = map[string]string{}
	for name, hash := range row.Written {
		updated.Written[name] = hash
	}
	updated.Written[key] = corpus.ValueHash(newValues[key])
	updated.Updated = nowFunc(opts.Now)().UTC().Format(time.RFC3339)
	if err := opts.State.Put(updated); err != nil {
		return fmt.Errorf("refs: update state of %s: %w", rel, err)
	}
	return nil
}

func topicWriter(opts Options) (*corpus.Writer, error) {
	if opts.Writer != nil {
		return opts.Writer, nil
	}
	store := opts.State
	if store == nil {
		opened, err := corpus.OpenState(opts.TopicRoot)
		if err != nil {
			return nil, fmt.Errorf("refs: %w", err)
		}
		store = opened
	}
	return corpus.NewWriter(opts.TopicRoot, store, opts.Now), nil
}

func quarantineID(rel, stamp string) string {
	sum := sha256.Sum256([]byte(rel + stamp))
	return "q-" + hex.EncodeToString(sum[:])[:10]
}

func appendLedger(topicRoot string, entry LedgerEntry) error {
	path := LedgerPath(topicRoot)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("refs: create %s: %w", filepath.Dir(path), err)
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("refs: encode ledger row: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("refs: open ledger: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		_ = file.Close()
		return fmt.Errorf("refs: append ledger: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("refs: close ledger: %w", err)
	}
	return nil
}

// writeAtomic replaces path with content through a temp file and rename,
// keeping the file mode.
func writeAtomic(path, content string) error {
	mode := fs.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	temp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("refs: write %s: %w", path, err)
	}
	name := temp.Name()
	if _, err := temp.WriteString(content); err != nil {
		_ = temp.Close()
		_ = os.Remove(name)
		return fmt.Errorf("refs: write %s: %w", path, err)
	}
	if err := temp.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("refs: write %s: %w", path, err)
	}
	if err := os.Chmod(name, mode); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("refs: write %s: %w", path, err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("refs: write %s: %w", path, err)
	}
	return nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("refs: read %s: %w", path, err)
	}
	return string(content), nil
}

func joinTopic(topicRoot, rel string) string {
	return filepath.Join(topicRoot, filepath.FromSlash(rel))
}

func hashString(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func nowFunc(now func() time.Time) func() time.Time {
	if now == nil {
		return time.Now
	}
	return now
}
