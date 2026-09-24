package questions

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestCompose(t *testing.T) {
	t.Parallel()
	relevance := MustLoad("relevance")
	quality := MustLoad("quality")

	composite := Compose(relevance, nil, quality)
	if composite.ID != "relevance+quality" {
		t.Fatalf("ID = %q", composite.ID)
	}
	if composite.Version != relevance.Version+"+"+quality.Version {
		t.Fatalf("Version = %q", composite.Version)
	}
	if composite.Purpose != relevance.Purpose {
		t.Fatalf("Purpose = %q, want the first member's", composite.Purpose)
	}
	hasher := sha256.New()
	for _, member := range []*Bank{relevance, quality} {
		hasher.Write([]byte(member.Hash))
		hasher.Write([]byte{0})
	}
	if want := hex.EncodeToString(hasher.Sum(nil)); composite.Hash != want {
		t.Fatalf("Hash = %q, want %q", composite.Hash, want)
	}
	if reversed := Compose(quality, relevance); reversed.Hash == composite.Hash {
		t.Fatal("member order must change the composite hash")
	}

	wantIDs := append(relevance.IDs(), quality.IDs()...)
	if got := composite.IDs(); !reflect.DeepEqual(got, wantIDs) {
		t.Fatalf("IDs() = %v, want %v", got, wantIDs)
	}

	tests := []struct {
		name   string
		member *Bank
		id     string
		vars   map[string]string
	}{
		{name: "first member question", member: relevance, id: "role"},
		{name: "templated first member question", member: relevance, id: "role_item_{id}", vars: map[string]string{"id": "i1"}},
		{name: "second member question", member: quality, id: "no_speech_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := composite.Question(tt.id, tt.vars)
			if err != nil {
				t.Fatalf("Question(%q): %v", tt.id, err)
			}
			want, err := tt.member.Question(tt.id, tt.vars)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("composite question differs from the owning member's:\n got %#v\nwant %#v", got, want)
			}
		})
	}
	if _, err := composite.Question("missing", nil); err == nil {
		t.Fatal("unknown question must fail")
	}
}

func TestComposeWithCriteriaKeepsMemberExits(t *testing.T) {
	t.Parallel()
	concept := MustLoad("concept")
	composite := Compose(MustLoad("classify"), concept)
	q, err := composite.Question("primary_concept", nil)
	if err != nil {
		t.Fatal(err)
	}
	injected := composite.WithCriteria(q, map[string]string{"c01": "Topic — criterion"})
	criteria, ok := injected.Criteria.(map[string]any)
	if !ok {
		t.Fatalf("criteria = %T", injected.Criteria)
	}
	if _, ok := criteria["none"]; !ok {
		t.Fatalf("WithCriteria dropped the member's exit option: %v", criteria)
	}
	if _, ok := criteria["c01"]; !ok {
		t.Fatalf("WithCriteria lost the injected option: %v", criteria)
	}
}

func TestComposeEdgeCases(t *testing.T) {
	t.Parallel()
	if Compose() != nil || Compose(nil) != nil {
		t.Fatal("Compose of no bank must be nil")
	}
	quality := MustLoad("quality")
	if Compose(quality) != quality {
		t.Fatal("Compose of one bank must return it unchanged")
	}
}
