package review

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/compozy/kb/internal/decisions"
)

func TestStoreAppendRepairsTornTail(t *testing.T) {
	t.Parallel()

	pending := Item{ID: "rv-a", Queue: QueueGate, Purpose: "relevance", Subject: "raw/a.md", Status: StatusPending, Created: "2026-09-24T12:00:00Z"}
	pendingRow, err := json.Marshal(pending)
	if err != nil {
		t.Fatal(err)
	}
	label := Label{Subject: "raw/a.md", Purpose: "relevance", Verdict: VerdictPositive, Origin: OriginReview}
	labelRow, err := json.Marshal(label)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name string
		file string
		seed string
		act  func(*Store) error
		// check runs against a freshly opened store (a new process).
		check func(t *testing.T, root string)
	}{
		{
			name: "review row after torn tail survives reopen",
			file: QueueFile,
			seed: string(pendingRow) + "\n" + `{"id":`,
			act: func(s *Store) error {
				_, err := s.Resolve("rv-a", StatusRejected)
				return err
			},
			check: func(t *testing.T, root string) {
				t.Helper()
				item, ok, err := Open(root, fixedClock).Get("rv-a")
				if err != nil || !ok {
					t.Fatalf("Get(rv-a) = %v, %v", ok, err)
				}
				if item.Status != StatusRejected {
					t.Fatalf("status after reopen = %q, want %q", item.Status, StatusRejected)
				}
			},
		},
		{
			name: "complete review row without newline is kept",
			file: QueueFile,
			seed: string(pendingRow),
			act: func(s *Store) error {
				_, err := s.Add(Item{ID: "rv-b", Queue: QueueGate, Purpose: "relevance", Subject: "raw/b.md"})
				return err
			},
			check: func(t *testing.T, root string) {
				t.Helper()
				items, err := Open(root, fixedClock).Items("", "")
				if err != nil {
					t.Fatal(err)
				}
				if len(items) != 2 || items[0].ID != "rv-a" || items[1].ID != "rv-b" {
					t.Fatalf("items after reopen = %+v, want rv-a then rv-b", items)
				}
			},
		},
		{
			name: "label after torn tail survives reopen",
			file: LabelsFile,
			seed: string(labelRow) + "\n" + `{"subject":"raw/`,
			act: func(s *Store) error {
				return s.AddLabel(Label{Subject: "raw/b.md", Purpose: "relevance", Verdict: VerdictNegative})
			},
			check: func(t *testing.T, root string) {
				t.Helper()
				labels, err := LoadLabels(root)
				if err != nil {
					t.Fatal(err)
				}
				if len(labels) != 2 || labels[1].Subject != "raw/b.md" || labels[1].Verdict != VerdictNegative {
					t.Fatalf("labels after reopen = %+v, want the appended raw/b.md label", labels)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			path := filepath.Join(root, decisions.ReceiptsDir, tc.file)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(tc.seed), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := tc.act(Open(root, fixedClock)); err != nil {
				t.Fatal(err)
			}
			tc.check(t, root)
		})
	}
}
