package review

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/compozy/kb/internal/decisions"
)

// PreviewsFile is the impact-preview ledger inside <topic>/.decisions/.
const PreviewsFile = "previews.jsonl"

// Preview records the documents a contract impact preview showed (or used
// to estimate precision) under one contract hash. Calibration demotes the
// labels of those subjects to dev once the contract changes (spec §12.3).
type Preview struct {
	Time     string   `json:"time"`
	Contract string   `json:"contract"`
	Subjects []string `json:"subjects"`
}

// AppendPreview appends p to previews.jsonl; an empty Time is set to now.
func AppendPreview(topicRoot string, p Preview) error {
	if p.Time == "" {
		p.Time = time.Now().UTC().Format(time.RFC3339)
	}
	return appendJSONL(filepath.Join(topicRoot, decisions.ReceiptsDir, PreviewsFile), p)
}

// LoadPreviews reads previews.jsonl (missing → empty; malformed lines are
// skipped).
func LoadPreviews(topicRoot string) ([]Preview, error) {
	var previews []Preview
	err := readJSONL(filepath.Join(topicRoot, decisions.ReceiptsDir, PreviewsFile), func(line []byte) error {
		var p Preview
		if json.Unmarshal(line, &p) == nil {
			previews = append(previews, p)
		}
		return nil
	})
	return previews, err
}
