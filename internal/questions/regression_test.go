package questions

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// RegressionFile is one testdata/regression/<bank>.json file: short states
// with the answers a correct model may give. StateDefaults are merged into
// every case state (case keys win). Exported for the jevlive test in
// package questions_test.
type RegressionFile struct {
	Bank          string           `json:"bank"`
	StateDefaults map[string]any   `json:"state_defaults"`
	Cases         []RegressionCase `json:"cases"`
}

// RegressionCase is one state + question + allowed answers. Allowed holds
// "yes"/"no" for nouls (P(yes) ≥ 0.5 is yes), option ids for choices, and
// rounded level numbers ("0".."n-1") for scores. Criteria, when set, are
// injected with WithCriteria (concept and OKF type options).
type RegressionCase struct {
	Name     string            `json:"name"`
	Question string            `json:"question"`
	Vars     map[string]string `json:"vars,omitempty"`
	Criteria map[string]any    `json:"criteria,omitempty"`
	State    map[string]any    `json:"state"`
	Allowed  []string          `json:"allowed"`
}

// LoadRegressionFiles reads every regression file in dir, keyed by bank id.
func LoadRegressionFiles(dir string) (map[string]RegressionFile, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	files := make(map[string]RegressionFile, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var file RegressionFile
		if err := json.Unmarshal(data, &file); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		files[file.Bank] = file
	}
	return files, nil
}

// Build returns the bank, the API-ready question and the merged state of a
// case.
func (f RegressionFile) Build(c RegressionCase) (*Bank, Q, map[string]any, error) {
	bank, err := Load(f.Bank)
	if err != nil {
		return nil, Q{}, nil, err
	}
	q, err := bank.Question(c.Question, c.Vars)
	if err != nil {
		return nil, Q{}, nil, err
	}
	if c.Criteria != nil {
		q = bank.WithCriteria(q, c.Criteria)
	}
	state := maps.Clone(f.StateDefaults)
	if state == nil {
		state = map[string]any{}
	}
	maps.Copy(state, c.State)
	return bank, q, state, nil
}

// AllowedValues lists the answers a case may legally allow for q.
func AllowedValues(q Q) []string {
	switch q.Type {
	case TypeNoul:
		return []string{"yes", "no"}
	case TypeChoice:
		criteria, _ := q.Criteria.(map[string]any)
		return slices.Sorted(maps.Keys(criteria))
	default:
		levels, _ := q.Criteria.([]any)
		out := make([]string, len(levels))
		for index := range levels {
			out[index] = strconv.Itoa(index)
		}
		return out
	}
}

func TestRegressionFilesAreWellFormed(t *testing.T) {
	t.Parallel()

	files, err := LoadRegressionFiles(filepath.Join("testdata", "regression"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := slices.Sorted(maps.Keys(files)), BankIDs(); !slices.Equal(got, want) {
		t.Fatalf("regression files for %v, want one per bank %v", got, want)
	}
	for bankID, file := range files {
		t.Run(bankID, func(t *testing.T) {
			t.Parallel()
			if n := len(file.Cases); n < 8 || n > 20 {
				t.Fatalf("%d cases, want 8 to 20", n)
			}
			covered := map[string]bool{}
			names := map[string]bool{}
			for _, c := range file.Cases {
				if c.Name == "" || names[c.Name] {
					t.Errorf("case name %q empty or duplicated", c.Name)
				}
				names[c.Name] = true
				_, q, state, err := file.Build(c)
				if err != nil {
					t.Errorf("%s: %v", c.Name, err)
					continue
				}
				covered[c.Question] = true
				if len(state) == 0 || len(c.Allowed) == 0 {
					t.Errorf("%s: state and allowed answers are required", c.Name)
				}
				legal := AllowedValues(q)
				for _, allowed := range c.Allowed {
					if !slices.Contains(legal, allowed) {
						t.Errorf("%s: allowed %q is not one of %v", c.Name, allowed, legal)
					}
				}
				if text, _ := q.Instructions.(string); strings.Contains(text, "{") {
					t.Errorf("%s: unresolved placeholder in %q", c.Name, text)
				}
			}
			bank := MustLoad(bankID)
			for _, id := range bank.IDs() {
				if !covered[id] {
					t.Errorf("question %q has no regression case", id)
				}
			}
		})
	}
}
