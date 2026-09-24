package corpus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"sync"

	"github.com/compozy/kb/internal/resolve"
)

// StateFile is the topic-relative path of the document state store.
const StateFile = resolve.DecisionsDir + "/state.jsonl"

// StateRow is one row of `<topic>/.decisions/state.jsonl` (spec §6): the
// bookkeeping kb keeps per document instead of writing it to frontmatter.
type StateRow struct {
	Path     string            `json:"path"`
	BodyHash string            `json:"body_hash"`
	Contract string            `json:"contract"`
	Banks    map[string]string `json:"banks"`
	Written  map[string]string `json:"written"`
	Updated  string            `json:"updated"`
}

func (row StateRow) clone() StateRow {
	row.Banks = maps.Clone(row.Banks)
	row.Written = maps.Clone(row.Written)
	return row
}

// StateStore is the append-only document state store. The last row per path
// wins. It is safe for concurrent use.
type StateStore struct {
	mu   sync.Mutex
	path string
	rows map[string]StateRow
	// needsSep is set when the file ends with a complete row but no newline,
	// so the next append starts on a fresh line.
	needsSep bool
	// tornAt is the offset of a torn (unparsable) last line, or -1. The next
	// append truncates it away so it never ends up in the middle of the file.
	tornAt int64
}

// OpenState reads `<topicRoot>/.decisions/state.jsonl`. A missing file is an
// empty store. A torn last line (a crash mid-append) is ignored; any other
// malformed line is an error.
func OpenState(topicRoot string) (*StateStore, error) {
	store := &StateStore{
		path:   filepath.Join(topicRoot, filepath.FromSlash(StateFile)),
		rows:   make(map[string]StateRow),
		tornAt: -1,
	}

	content, err := os.ReadFile(store.path)
	if errors.Is(err, fs.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("corpus: read state %q: %w", store.path, err)
	}

	lines := bytes.Split(content, []byte("\n"))
	offset := int64(0)
	for index, line := range lines {
		lineStart := offset
		offset += int64(len(line)) + 1
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var row StateRow
		if err := json.Unmarshal(line, &row); err != nil || row.Path == "" {
			if index == len(lines)-1 {
				store.tornAt = lineStart
				continue
			}
			return nil, fmt.Errorf("corpus: state %q line %d is malformed", store.path, index+1)
		}
		store.rows[row.Path] = row
	}
	if store.tornAt < 0 && len(content) > 0 && content[len(content)-1] != '\n' {
		store.needsSep = true
	}

	return store, nil
}

// Path returns the absolute path of the state file.
func (s *StateStore) Path() string { return s.path }

// Get returns the latest row for the topic-relative path.
func (s *StateStore) Get(path string) (*StateRow, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row, ok := s.rows[path]
	if !ok {
		return nil, false
	}
	cloned := row.clone()
	return &cloned, true
}

// Rows returns the latest row of every path, sorted by path.
func (s *StateStore) Rows() []StateRow {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows := make([]StateRow, 0, len(s.rows))
	for _, path := range slices.Sorted(maps.Keys(s.rows)) {
		rows = append(rows, s.rows[path].clone())
	}

	return rows
}

// FindByBodyHash returns the latest rows whose body hash is h, sorted by
// path.
func (s *StateStore) FindByBodyHash(h string) []StateRow {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows := make([]StateRow, 0)
	for _, row := range s.rows {
		if row.BodyHash == h {
			rows = append(rows, row.clone())
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Path < rows[j].Path })

	return rows
}

// Lookup returns the row of path. When path has no row it performs rename
// detection: the first row with the same body hash whose own path no longer
// exists (exists reports false) follows the file and is returned with Path
// set to path. The store itself is not modified.
func (s *StateStore) Lookup(path, bodyHash string, exists func(path string) bool) (*StateRow, bool) {
	if row, ok := s.Get(path); ok {
		return row, true
	}
	if bodyHash == "" {
		return nil, false
	}
	for _, row := range s.FindByBodyHash(bodyHash) {
		if exists != nil && exists(row.Path) {
			continue
		}
		row.Path = path
		return &row, true
	}

	return nil, false
}

// Put appends row as one whole JSON line and makes it the latest row for its
// path.
func (s *StateStore) Put(row StateRow) error {
	if row.Path == "" {
		return errors.New("corpus: state row path is required")
	}
	encoded, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("corpus: encode state row: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("corpus: create state directory: %w", err)
	}
	if s.tornAt >= 0 {
		if err := os.Truncate(s.path, s.tornAt); err != nil {
			return fmt.Errorf("corpus: drop torn state line: %w", err)
		}
		s.tornAt = -1
	}
	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("corpus: open state %q: %w", s.path, err)
	}

	line := make([]byte, 0, len(encoded)+2)
	if s.needsSep {
		line = append(line, '\n')
	}
	line = append(line, encoded...)
	line = append(line, '\n')
	if _, err := file.Write(line); err != nil {
		_ = file.Close()
		return fmt.Errorf("corpus: append state %q: %w", s.path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("corpus: close state %q: %w", s.path, err)
	}

	s.needsSep = false
	s.rows[row.Path] = row.clone()
	return nil
}

// Compact atomically rewrites the state file with only the latest row per
// path, sorted by path.
func (s *StateStore) Compact() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var buffer bytes.Buffer
	for _, path := range slices.Sorted(maps.Keys(s.rows)) {
		encoded, err := json.Marshal(s.rows[path])
		if err != nil {
			return fmt.Errorf("corpus: encode state row: %w", err)
		}
		buffer.Write(encoded)
		buffer.WriteByte('\n')
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("corpus: create state directory: %w", err)
	}
	if err := writeFileAtomic(s.path, buffer.Bytes(), 0o644); err != nil {
		return fmt.Errorf("corpus: compact state: %w", err)
	}
	s.needsSep = false
	s.tornAt = -1

	return nil
}

// Unclassified reports whether doc needs a (re)classification: it has no
// state row, its body changed since the row was written, or the row was
// written under another contract hash (when contract is non-empty) or
// another version of any bank in banks.
func Unclassified(doc *Document, row *StateRow, contract string, banks map[string]string) bool {
	if doc == nil || row == nil {
		return true
	}
	if row.BodyHash != doc.BodyHash {
		return true
	}
	if contract != "" && row.Contract != contract {
		return true
	}
	for bank, version := range banks {
		if row.Banks[bank] != version {
			return true
		}
	}

	return false
}

// writeFileAtomic writes data to a temp file in the target directory and
// renames it over path.
func writeFileAtomic(path string, data []byte, mode fs.FileMode) error {
	temp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".kb-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tempPath := temp.Name()
	cleanup := func() { _ = os.Remove(tempPath) }

	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		cleanup()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := temp.Chmod(mode); err != nil {
		_ = temp.Close()
		cleanup()
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := temp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		cleanup()
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}
