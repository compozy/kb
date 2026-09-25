package decisions

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/compozy/kb/internal/jsonl"
)

// ReceiptsDir is the per-topic directory holding decision records.
const ReceiptsDir = ".decisions"

// ReceiptsFile is the append-only receipts log inside ReceiptsDir.
const ReceiptsFile = "receipts.jsonl"

// Routes recorded in receipts.
const (
	RouteDecisions       = "openrouter:alpha/decisions"
	RouteChatCompletions = "openrouter:v1/chat/completions"
)

// Receipt is one row of <topic>/.decisions/receipts.jsonl: one network call
// (one decision request, or one generation request). Answers hold the raw
// answers as returned by the model, never bands: bands are recomputed from
// the current thresholds on every read. Generation rows use purpose
// "generate:<kind>" and answers {"output": <json>}.
type Receipt struct {
	Key           string                     `json:"key"`
	Time          string                     `json:"time"`
	Purpose       string                     `json:"purpose"`
	Subject       string                     `json:"subject"`
	Bank          string                     `json:"bank"`
	BankVersion   string                     `json:"bank_version"`
	BankHash      string                     `json:"bank_hash"`
	Contract      string                     `json:"contract"`
	Model         string                     `json:"model"`
	ReportedModel string                     `json:"reported_model"`
	Route         string                     `json:"route"`
	ResponseID    string                     `json:"response_id,omitempty"`
	Attempts      int                        `json:"attempts"`
	LatencyMS     int64                      `json:"latency_ms"`
	InputTokens   *int64                     `json:"input_tokens"`
	Cost          *float64                   `json:"cost"`
	CostUnknown   bool                       `json:"cost_unknown"`
	Status        Status                     `json:"status"`
	Reason        string                     `json:"reason,omitempty"`
	Questions     []string                   `json:"questions"`
	Answers       map[string]json.RawMessage `json:"answers"`
}

func (r Receipt) clone() Receipt {
	out := r
	out.Questions = slices.Clone(r.Questions)
	if r.InputTokens != nil {
		tokens := *r.InputTokens
		out.InputTokens = &tokens
	}
	if r.Cost != nil {
		cost := *r.Cost
		out.Cost = &cost
	}
	if r.Answers != nil {
		out.Answers = make(map[string]json.RawMessage, len(r.Answers))
		for id, raw := range r.Answers {
			out.Answers[id] = slices.Clone(raw)
		}
	}
	return out
}

// ReceiptsPath returns the receipts log path of a topic.
func ReceiptsPath(topicRoot string) string {
	return filepath.Join(topicRoot, ReceiptsDir, ReceiptsFile)
}

// LoadReceipts reads every well-formed row of a topic's receipts log in file
// order. A missing log yields no rows; malformed lines (for example a line
// cut by a crash) are skipped. Readers that need one row per key take the
// last one.
func LoadReceipts(topicRoot string) ([]Receipt, error) {
	file, err := os.Open(ReceiptsPath(topicRoot))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("open receipts: %w", err)
	}
	defer func() { _ = file.Close() }()

	var rows []Receipt
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 64*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var row Receipt
		if err := json.Unmarshal(line, &row); err != nil || row.Key == "" {
			continue
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read receipts: %w", err)
	}
	return rows, nil
}

// Receipts is the per-topic receipts log plus its in-memory cache. The
// engine and the generation client share one instance (Engine.Receipts) so
// appends to a topic's log are serialized. Each topic's log is loaded lazily,
// once; only rows with status decided serve as cache (last one per key wins).
// A topic with an empty root is cached in memory only.
type Receipts struct {
	mu     sync.Mutex
	topics map[string]map[string]Receipt
}

// NewReceipts returns an empty store.
func NewReceipts() *Receipts {
	return &Receipts{topics: map[string]map[string]Receipt{}}
}

// Lookup returns a copy of the cached decided row for key.
func (s *Receipts) Lookup(topicRoot, key string) (Receipt, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.topicLocked(topicRoot)
	if err != nil {
		return Receipt{}, false, err
	}
	row, ok := rows[key]
	if !ok {
		return Receipt{}, false, nil
	}
	return row.clone(), true, nil
}

// Append writes row as one whole JSON line to the topic's log and, when the
// row is decided, makes it the cached row for its key.
func (s *Receipts) Append(topicRoot string, row Receipt) error {
	line, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("encode receipt: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.topicLocked(topicRoot)
	if err != nil {
		return err
	}
	if topicRoot != "" {
		if err := jsonl.Append(ReceiptsPath(topicRoot), line); err != nil {
			return fmt.Errorf("append receipt: %w", err)
		}
	}
	if row.Status == StatusDecided {
		rows[row.Key] = row.clone()
	}
	return nil
}

func (s *Receipts) topicLocked(topicRoot string) (map[string]Receipt, error) {
	if rows, ok := s.topics[topicRoot]; ok {
		return rows, nil
	}
	rows := map[string]Receipt{}
	if topicRoot != "" {
		loaded, err := LoadReceipts(topicRoot)
		if err != nil {
			return nil, err
		}
		for _, row := range loaded {
			if row.Status == StatusDecided {
				rows[row.Key] = row
			}
		}
	}
	s.topics[topicRoot] = rows
	return rows, nil
}
