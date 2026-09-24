package scope

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/fakes"
)

func TestBuildPairs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		draft *contract.Contract
		want  []Pair
	}{
		{name: "no out_of_scope", draft: &contract.Contract{Purpose: "P", Core: []string{"A"}}, want: []Pair{}},
		{name: "core then adjacent against every excluded line", draft: &contract.Contract{
			Purpose: "P", Core: []string{"A"}, Adjacent: []string{"B"}, OutOfScope: []string{"X", "Y"},
		}, want: []Pair{
			{N: 1, KeptField: "core", Kept: "A", Excluded: "X"},
			{N: 2, KeptField: "core", Kept: "A", Excluded: "Y"},
			{N: 3, KeptField: "adjacent", Kept: "B", Excluded: "X"},
			{N: 4, KeptField: "adjacent", Kept: "B", Excluded: "Y"},
		}},
	}
	for _, tt := range tests {
		if got := BuildPairs(tt.draft); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: BuildPairs = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}

func TestSelfCheckFlagsPairsAtThreshold(t *testing.T) {
	t.Parallel()
	vault, _ := newTestTopic(t)
	fake := fakeServer(t, func(_ fakes.Call, q fakes.Question) any {
		switch q.ID {
		case "conflict_2":
			return fakes.Noul(0.7)
		case "conflict_3":
			return fakes.Noul(0.69)
		}
		return nil
	}, nil)
	s := openTestSession(t, vault, fake)
	draft := &contract.Contract{Purpose: "P", Core: []string{"Basketball studies"}, Adjacent: []string{"Rugby studies kept as comparison"}, OutOfScope: []string{"Stock picks", "Any sport other than basketball"}}

	got, err := SelfCheck(context.Background(), s, draft)
	if err != nil {
		t.Fatalf("SelfCheck: %v", err)
	}
	want := Conflicts{Pairs: 4, Threshold: 0.7, Conflicts: []Conflict{{Pair: Pair{N: 2, KeptField: "core", Kept: "Basketball studies", Excluded: "Any sport other than basketball"}, P: 0.7}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SelfCheck = %+v, want %+v", got, want)
	}
	if !got.Blocking() {
		t.Fatal("a conflict must block")
	}

	calls := fake.Calls()
	if len(calls) != 1 {
		t.Fatalf("calls = %d, want one request", len(calls))
	}
	var state struct {
		Pairs []struct {
			N        int    `json:"n"`
			Kept     string `json:"kept"`
			Excluded string `json:"excluded"`
		} `json:"pairs"`
	}
	if err := json.Unmarshal(calls[0].State, &state); err != nil {
		t.Fatal(err)
	}
	if len(state.Pairs) != 4 || state.Pairs[3].N != 4 || state.Pairs[3].Kept != "Rugby studies kept as comparison" || state.Pairs[3].Excluded != "Any sport other than basketball" {
		t.Fatalf("state = %s", calls[0].State)
	}
	if len(calls[0].Questions) != 4 || !strings.Contains(fmt.Sprint(calls[0].Questions["conflict_4"].Instructions), "pair `4`") {
		t.Fatalf("questions = %+v", calls[0].Questions)
	}
}

func TestSelfCheckSplitsLargeContracts(t *testing.T) {
	t.Parallel()
	vault, _ := newTestTopic(t)
	fake := fakeServer(t, nil, nil)
	s := openTestSession(t, vault, fake)
	draft := &contract.Contract{Purpose: "P"}
	for i := range 10 {
		draft.Core = append(draft.Core, fmt.Sprintf("Kept %d", i))
	}
	for i := range 5 {
		draft.OutOfScope = append(draft.OutOfScope, fmt.Sprintf("Excluded %d", i))
	}

	got, err := SelfCheck(context.Background(), s, draft)
	if err != nil {
		t.Fatalf("SelfCheck: %v", err)
	}
	if got.Pairs != 50 || len(got.Conflicts) != 0 || len(got.Undecided) != 0 {
		t.Fatalf("SelfCheck = %+v", got)
	}
	if calls := fake.Calls(); len(calls) != 2 {
		t.Fatalf("calls = %d, want 2 (50 questions over 48 per request)", len(calls))
	}
}
