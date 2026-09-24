package gate

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/compozy/kb/internal/decisions"
)

// SkippedFile is the log of items skipped before fetch, under
// <topic>/.decisions/.
const SkippedFile = "skipped.jsonl"

// ErrSkippedNotFound is returned when a skipped id has no row.
var ErrSkippedNotFound = errors.New("gate: skipped item not found")

// SkippedRow is one row of skipped.jsonl (plan schema). Rescue appends a new
// row with the same id and Rescued true; readers keep the last row per id.
type SkippedRow struct {
	ID          string  `json:"id"`
	Time        string  `json:"time"`
	Batch       string  `json:"batch,omitempty"`
	Query       string  `json:"query,omitempty"`
	URL         string  `json:"url"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Channel     string  `json:"channel,omitempty"`
	Date        string  `json:"date,omitempty"`
	POffTopic   float64 `json:"p_off_topic"`
	ReceiptKey  string  `json:"receipt_key,omitempty"`
	Rescued     bool    `json:"rescued"`
}

// SkippedID derives the stable id of a skipped item from its URL:
// "sk-" + the first 10 hex digits of sha256(NormalizeURL(url)).
func SkippedID(rawURL string) string {
	key := NormalizeURL(rawURL)
	if key == "" {
		key = strings.TrimSpace(rawURL)
	}
	sum := sha256.Sum256([]byte(key))
	return "sk-" + hex.EncodeToString(sum[:])[:10]
}

// SkippedPath returns <topic>/.decisions/skipped.jsonl.
func SkippedPath(topicRoot string) string {
	return filepath.Join(topicRoot, decisions.ReceiptsDir, SkippedFile)
}

var skippedMu sync.Mutex

// AppendSkipped appends one row to skipped.jsonl.
func AppendSkipped(topicRoot string, row SkippedRow) error {
	skippedMu.Lock()
	defer skippedMu.Unlock()
	path := SkippedPath(topicRoot)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("gate: skipped log: %w", err)
	}
	line, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("gate: skipped log: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("gate: skipped log: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		_ = file.Close()
		return fmt.Errorf("gate: skipped log: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("gate: skipped log: %w", err)
	}
	return nil
}

// LoadSkipped reads skipped.jsonl and returns the last row per id, in first
// appearance order (a missing file is empty; torn lines are ignored).
func LoadSkipped(topicRoot string) ([]SkippedRow, error) {
	file, err := os.Open(SkippedPath(topicRoot))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gate: skipped log: %w", err)
	}
	defer func() { _ = file.Close() }()
	rows := map[string]SkippedRow{}
	order := make([]string, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var row SkippedRow
		if json.Unmarshal(line, &row) != nil || row.ID == "" {
			continue
		}
		if _, seen := rows[row.ID]; !seen {
			order = append(order, row.ID)
		}
		rows[row.ID] = row
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("gate: skipped log: %w", err)
	}
	out := make([]SkippedRow, 0, len(order))
	for _, id := range order {
		out = append(out, rows[id])
	}
	return out, nil
}

// FindSkipped returns the latest row of one skipped id.
func FindSkipped(topicRoot, id string) (SkippedRow, error) {
	rows, err := LoadSkipped(topicRoot)
	if err != nil {
		return SkippedRow{}, err
	}
	id = strings.TrimSpace(id)
	for _, row := range rows {
		if row.ID == id {
			return row, nil
		}
	}
	return SkippedRow{}, fmt.Errorf("%w: %q", ErrSkippedNotFound, id)
}

// MarkRescued appends the row again with Rescued true.
func MarkRescued(topicRoot string, row SkippedRow, stamp string) error {
	row.Rescued = true
	row.Time = stamp
	return AppendSkipped(topicRoot, row)
}
