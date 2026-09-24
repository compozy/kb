// Package review owns a topic's review queue, its labels (human verdicts,
// including imported ones) and calibration with a holdout (spec §12).
package review

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/decisions"
)

// Queue names (spec §12.1).
const (
	QueueGate            = "gate"
	QueueSkip            = "skip"
	QueueLink            = "link"
	QueueContradiction   = "contradiction"
	QueueConceptProposal = "concept-proposal"
	QueueOKFType         = "okf-type"
	QueueRecapture       = "recapture"
	QueueRemove          = "remove"
)

// Queues lists every queue in display order.
var Queues = []string{QueueGate, QueueSkip, QueueRecapture, QueueRemove, QueueLink, QueueContradiction, QueueConceptProposal, QueueOKFType}

// Item statuses.
const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)

// Files under <topic>/.decisions/.
const (
	QueueFile  = "review.jsonl"
	LabelsFile = "labels.jsonl"
)

// Item is one review-queue entry. Status changes are appended as new rows
// with the same ID; readers keep the last row per ID.
type Item struct {
	ID          string         `json:"id"`
	Queue       string         `json:"queue"`
	Purpose     string         `json:"purpose"`
	Subject     string         `json:"subject"`
	Target      string         `json:"target,omitempty"`
	Question    string         `json:"question,omitempty"`
	Probability float64        `json:"probability"`
	ReceiptKey  string         `json:"receipt_key,omitempty"`
	Evidence    string         `json:"evidence,omitempty"`
	Action      map[string]any `json:"action,omitempty"`
	Status      string         `json:"status"`
	Created     string         `json:"created"`
	Resolved    string         `json:"resolved,omitempty"`
}

// ItemID derives the stable id of an item from its identity fields.
func ItemID(queue, subject, target, question string) string {
	sum := sha256.Sum256([]byte(queue + "\x00" + subject + "\x00" + target + "\x00" + question))
	return "rv-" + hex.EncodeToString(sum[:])[:10]
}

// Label is a human verdict on a decision (spec §12.2).
type Label struct {
	Time       string `json:"time"`
	Subject    string `json:"subject"`
	Purpose    string `json:"purpose"`
	Question   string `json:"question,omitempty"`
	Target     string `json:"target,omitempty"`
	Verdict    string `json:"verdict"` // positive | negative
	ReceiptKey string `json:"receipt_key,omitempty"`
	Origin     string `json:"origin"` // review | import:<file>@<sha256> | import-links
	DecidedBy  string `json:"decided_by,omitempty"`
	ItemID     string `json:"item_id,omitempty"`
}

// Label verdicts and origins.
const (
	VerdictPositive   = "positive"
	VerdictNegative   = "negative"
	OriginReview      = "review"
	OriginImportLinks = "import-links"
)

// Store is a topic's review queue and label log. Safe for concurrent use.
type Store struct {
	root string
	now  func() time.Time

	mu     sync.Mutex
	items  map[string]Item
	order  []string
	loaded bool
}

// Open returns the review store of a topic (files are read lazily).
func Open(topicRoot string, now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{root: topicRoot, now: now}
}

// Dir returns <topic>/.decisions.
func (s *Store) Dir() string { return filepath.Join(s.root, decisions.ReceiptsDir) }

func (s *Store) load() error {
	if s.loaded {
		return nil
	}
	s.items = map[string]Item{}
	s.order = nil
	err := readJSONL(filepath.Join(s.Dir(), QueueFile), func(line []byte) error {
		var item Item
		if err := json.Unmarshal(line, &item); err != nil {
			return nil // a torn line is ignored
		}
		if _, seen := s.items[item.ID]; !seen {
			s.order = append(s.order, item.ID)
		}
		s.items[item.ID] = item
		return nil
	})
	if err != nil {
		return err
	}
	s.loaded = true
	return nil
}

// Add enqueues an item as pending. An item whose ID already exists (pending
// or resolved) is left alone, so a verdict is never re-asked; returns whether
// a row was written.
func (s *Store) Add(item Item) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.load(); err != nil {
		return false, err
	}
	if item.ID == "" {
		item.ID = ItemID(item.Queue, item.Subject, item.Target, item.Question)
	}
	if _, ok := s.items[item.ID]; ok {
		return false, nil
	}
	item.Status = StatusPending
	if item.Created == "" {
		item.Created = s.now().UTC().Format(time.RFC3339)
	}
	if err := s.appendItem(item); err != nil {
		return false, err
	}
	return true, nil
}

// Resolve marks an item accepted or rejected.
func (s *Store) Resolve(id, status string) (Item, error) {
	if status != StatusAccepted && status != StatusRejected {
		return Item{}, fmt.Errorf("review: invalid status %q", status)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.load(); err != nil {
		return Item{}, err
	}
	item, ok := s.items[id]
	if !ok {
		return Item{}, fmt.Errorf("review: unknown item %q", id)
	}
	item.Status = status
	item.Resolved = s.now().UTC().Format(time.RFC3339)
	return item, s.appendItem(item)
}

// Get returns an item by id.
func (s *Store) Get(id string) (Item, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.load(); err != nil {
		return Item{}, false, err
	}
	item, ok := s.items[id]
	return item, ok, nil
}

// Items returns every item (latest row per id) in insertion order, filtered
// by queue ("" = all) and status ("" = all).
func (s *Store) Items(queue, status string) ([]Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.load(); err != nil {
		return nil, err
	}
	out := make([]Item, 0, len(s.order))
	for _, id := range s.order {
		item := s.items[id]
		if queue != "" && item.Queue != queue {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		out = append(out, item)
	}
	return out, nil
}

// Pending returns pending items of one queue ("" = all), highest probability
// first.
func (s *Store) Pending(queue string) ([]Item, error) {
	items, err := s.Items(queue, StatusPending)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Probability > items[j].Probability })
	return items, nil
}

// PendingCounts counts pending items per queue.
func (s *Store) PendingCounts() (map[string]int, error) {
	items, err := s.Items("", StatusPending)
	if err != nil {
		return nil, err
	}
	counts := map[string]int{}
	for _, item := range items {
		counts[item.Queue]++
	}
	return counts, nil
}

func (s *Store) appendItem(item Item) error {
	if err := appendJSONL(filepath.Join(s.Dir(), QueueFile), item); err != nil {
		return err
	}
	if _, seen := s.items[item.ID]; !seen {
		s.order = append(s.order, item.ID)
	}
	s.items[item.ID] = item
	return nil
}

// AddLabel appends a label to labels.jsonl.
func (s *Store) AddLabel(label Label) error {
	if label.Verdict != VerdictPositive && label.Verdict != VerdictNegative {
		return fmt.Errorf("review: invalid label verdict %q", label.Verdict)
	}
	if strings.TrimSpace(label.Origin) == "" {
		label.Origin = OriginReview
	}
	if label.Time == "" {
		label.Time = s.now().UTC().Format(time.RFC3339)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return appendJSONL(filepath.Join(s.Dir(), LabelsFile), label)
}

// Labels reads every label of the topic.
func (s *Store) Labels() ([]Label, error) {
	return LoadLabels(s.root)
}

// LoadLabels reads <topic>/.decisions/labels.jsonl (missing → empty).
func LoadLabels(topicRoot string) ([]Label, error) {
	var labels []Label
	err := readJSONL(filepath.Join(topicRoot, decisions.ReceiptsDir, LabelsFile), func(line []byte) error {
		var label Label
		if err := json.Unmarshal(line, &label); err != nil {
			return nil
		}
		labels = append(labels, label)
		return nil
	})
	return labels, err
}

func readJSONL(path string, fn func([]byte) error) error {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("review: open %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		if err := fn(line); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("review: read %s: %w", path, err)
	}
	return nil
}

func appendJSONL(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("review: create %s: %w", filepath.Dir(path), err)
	}
	line, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("review: encode row: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("review: open %s: %w", path, err)
	}
	prefix, err := repairTail(file)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("review: repair %s: %w", path, err)
	}
	row := make([]byte, 0, len(prefix)+len(line)+1)
	row = append(append(append(row, prefix...), line...), '\n')
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		_ = file.Close()
		return fmt.Errorf("review: append %s: %w", path, err)
	}
	if _, err := file.Write(row); err != nil {
		_ = file.Close()
		return fmt.Errorf("review: append %s: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("review: sync %s: %w", path, err)
	}
	return file.Close()
}

// repairTail makes the file end on a row boundary before an append, the way
// the corpus state store recovers from a crash mid-append: a final row
// without its newline is kept (the returned prefix supplies the separator)
// when it is a complete JSON value, and truncated away when it is torn, so
// the next row is never glued to — and discarded with — the torn bytes.
func repairTail(file *os.File) ([]byte, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	if size == 0 {
		return nil, nil
	}
	last := make([]byte, 1)
	if _, err := file.ReadAt(last, size-1); err != nil {
		return nil, err
	}
	if last[0] == '\n' {
		return nil, nil
	}
	start, err := lastLineStart(file, size)
	if err != nil {
		return nil, err
	}
	tail := make([]byte, size-start)
	if _, err := file.ReadAt(tail, start); err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(tail)
	if len(trimmed) == 0 || json.Valid(trimmed) {
		return []byte{'\n'}, nil
	}
	if err := file.Truncate(start); err != nil {
		return nil, err
	}
	return nil, nil
}

// lastLineStart returns the offset just past the last newline before size
// (0 when the file has none), scanning backwards in chunks.
func lastLineStart(file *os.File, size int64) (int64, error) {
	const chunk = 64 * 1024
	buf := make([]byte, chunk)
	end := size
	for end > 0 {
		begin := max(end-chunk, 0)
		part := buf[:end-begin]
		if _, err := file.ReadAt(part, begin); err != nil {
			return 0, err
		}
		if i := bytes.LastIndexByte(part, '\n'); i >= 0 {
			return begin + int64(i) + 1, nil
		}
		end = begin
	}
	return 0, nil
}
