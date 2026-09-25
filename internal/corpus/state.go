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
//
// BodyHash and Contract describe the document at kb's last write (rename
// detection follows BodyHash). Freshness is tracked per judgment instead: a
// write by one command (link, review, ingest) must not make another
// command's answers look current. BankBody and BankContract record, per
// question bank, the body and contract its answers were computed on;
// WrittenBody records, per owned key, the body its value was derived from;
// Facts holds small per-document outcomes other runs aggregate (for example
// the primary concept). Rows written before these maps existed read as if
// every bank and key was judged on BodyHash under Contract.
type StateRow struct {
	Path         string            `json:"path"`
	BodyHash     string            `json:"body_hash"`
	Contract     string            `json:"contract"`
	Banks        map[string]string `json:"banks"`
	Written      map[string]string `json:"written"`
	BankBody     map[string]string `json:"bank_body,omitempty"`
	BankContract map[string]string `json:"bank_contract,omitempty"`
	WrittenBody  map[string]string `json:"written_body,omitempty"`
	Facts        map[string]string `json:"facts,omitempty"`
	Updated      string            `json:"updated"`
}

func (row StateRow) clone() StateRow {
	row.Banks = maps.Clone(row.Banks)
	row.Written = maps.Clone(row.Written)
	row.BankBody = maps.Clone(row.BankBody)
	row.BankContract = maps.Clone(row.BankContract)
	row.WrittenBody = maps.Clone(row.WrittenBody)
	row.Facts = maps.Clone(row.Facts)
	return row
}

// JudgedBody returns the body hash the answers of bank were computed on
// (BodyHash for rows that predate per-bank tracking).
func (row *StateRow) JudgedBody(bank string) string {
	if hash, ok := row.BankBody[bank]; ok {
		return hash
	}
	return row.BodyHash
}

// JudgedContract returns the contract hash bank was judged under (Contract
// for rows that predate per-bank tracking).
func (row *StateRow) JudgedContract(bank string) string {
	if hash, ok := row.BankContract[bank]; ok {
		return hash
	}
	return row.Contract
}

// KeyBody returns the body hash the value of the owned key was derived from
// (BodyHash for rows that predate per-key tracking).
func (row *StateRow) KeyBody(key string) string {
	if hash, ok := row.WrittenBody[key]; ok {
		return hash
	}
	return row.BodyHash
}

// StateStore is the append-only document state store. The last row per path
// wins. Methods are safe for concurrent use; appends and compaction coordinate
// across independent stores and processes. Reads use the store's snapshot;
// open a new store to observe another writer's updates.
type StateStore struct {
	mu   sync.Mutex
	path string
	rows map[string]StateRow
}

// OpenState reads `<topicRoot>/.decisions/state.jsonl`. A missing file is an
// empty store. A torn last line (a crash mid-append) is ignored; any other
// malformed line is an error.
func OpenState(topicRoot string) (*StateStore, error) {
	path := filepath.Join(topicRoot, filepath.FromSlash(StateFile))
	rows, _, err := readState(path)
	if err != nil {
		return nil, err
	}
	return &StateStore{path: path, rows: rows}, nil
}

// stateTail describes only the file just read. Mutations must obtain it again
// under the file lock instead of reusing offsets from an earlier snapshot.
type stateTail struct {
	tornAt   int64
	needsSep bool
}

func readState(path string) (map[string]StateRow, stateTail, error) {
	rows := make(map[string]StateRow)
	tail := stateTail{tornAt: -1}
	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return rows, tail, nil
	}
	if err != nil {
		return nil, tail, fmt.Errorf("corpus: read state %q: %w", path, err)
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
				tail.tornAt = lineStart
				continue
			}
			return nil, tail, fmt.Errorf("corpus: state %q line %d is malformed", path, index+1)
		}
		rows[row.Path] = row
	}
	tail.needsSep = tail.tornAt < 0 && len(content) > 0 && content[len(content)-1] != '\n'
	return rows, tail, nil
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
func (s *StateStore) Put(row StateRow) (err error) {
	if row.Path == "" {
		return errors.New("corpus: state row path is required")
	}
	encoded, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("corpus: encode state row: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	unlock, err := s.lockFile()
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, unlock()) }()

	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("corpus: open state %q: %w", s.path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("corpus: close state %q: %w", s.path, closeErr))
		}
	}()
	needsSep, err := s.repairTail(file)
	if err != nil {
		return err
	}
	line := make([]byte, 0, len(encoded)+2)
	if needsSep {
		line = append(line, '\n')
	}
	line = append(append(line, encoded...), '\n')
	if _, err := file.Write(line); err != nil {
		return fmt.Errorf("corpus: append state %q: %w", s.path, err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("corpus: sync state %q: %w", s.path, err)
	}
	s.rows[row.Path] = row.clone()
	return nil
}

// repairTail runs under the file lock and inspects the current tail. Normal
// appends read only one byte; recovering an incomplete write revalidates the
// log before truncating, preserving strict rejection of malformed middle rows.
func (s *StateStore) repairTail(file *os.File) (bool, error) {
	info, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("corpus: stat state: %w", err)
	}
	if info.Size() == 0 {
		return false, nil
	}
	var last [1]byte
	if _, err := file.ReadAt(last[:], info.Size()-1); err != nil {
		return false, fmt.Errorf("corpus: read state tail: %w", err)
	}
	if last[0] == '\n' {
		return false, nil
	}
	_, tail, err := readState(s.path)
	if err != nil {
		return false, err
	}
	if tail.tornAt >= 0 {
		if err := file.Truncate(tail.tornAt); err != nil {
			return false, fmt.Errorf("corpus: drop torn state line: %w", err)
		}
	}
	return tail.needsSep, nil
}

// Compact atomically rewrites the current on-disk state with only the latest
// row per path, sorted by path, then refreshes this store's snapshot.
func (s *StateStore) Compact() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	unlock, err := s.lockFile()
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, unlock()) }()

	rows, _, err := readState(s.path)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	for _, path := range slices.Sorted(maps.Keys(rows)) {
		encoded, err := json.Marshal(rows[path])
		if err != nil {
			return fmt.Errorf("corpus: encode state row: %w", err)
		}
		buffer.Write(encoded)
		buffer.WriteByte('\n')
	}
	if err := writeFileAtomic(s.path, buffer.Bytes(), 0o644); err != nil {
		return fmt.Errorf("corpus: compact state: %w", err)
	}
	s.rows = rows
	return nil
}

// lockFile uses a separate file because compaction replaces the state file's
// inode. The OS releases the lock if a writer exits unexpectedly; the lock file
// itself remains so every writer continues to lock the same inode.
func (s *StateStore) lockFile() (func() error, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return nil, fmt.Errorf("corpus: create state directory: %w", err)
	}
	file, err := os.OpenFile(s.path+".lock", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("corpus: open state lock: %w", err)
	}
	if err := lockStateFile(file); err != nil {
		return nil, fmt.Errorf("corpus: lock state: %w", errors.Join(err, file.Close()))
	}
	return func() error {
		if err := errors.Join(unlockStateFile(file), file.Close()); err != nil {
			return fmt.Errorf("corpus: unlock state: %w", err)
		}
		return nil
	}, nil
}

// Unclassified reports whether doc needs a (re)classification: it has no
// state row, or any bank in banks was recorded with another version, judged
// on another body, or judged under another contract hash (when contract is
// non-empty). Writes by other commands never refresh these per-bank records,
// so a body edit followed by `kb link` still leaves doc unclassified.
func Unclassified(doc *Document, row *StateRow, contract string, banks map[string]string) bool {
	if doc == nil || row == nil {
		return true
	}
	if len(banks) == 0 {
		return row.BodyHash != doc.BodyHash || (contract != "" && row.Contract != contract)
	}
	for bank, version := range banks {
		if row.Banks[bank] != version || row.JudgedBody(bank) != doc.BodyHash {
			return true
		}
		if contract != "" && row.JudgedContract(bank) != contract {
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
