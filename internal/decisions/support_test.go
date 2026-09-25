package decisions

import (
	"encoding/json"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
)

func TestDefaultThresholdsMatchConfigNames(t *testing.T) {
	t.Parallel()

	got := slices.Sorted(maps.Keys(DefaultThresholds()))
	want := slices.Sorted(slices.Values(config.KnownThresholdNames))
	if !slices.Equal(got, want) {
		t.Fatalf("DefaultThresholds keys %v != config.KnownThresholdNames %v", got, want)
	}
}

func TestThresholdsForPurpose(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		purpose     Purpose
		thresholds  Thresholds
		apply, rev  float64
		value       float64
		wantBand    Band
		description string
	}{
		{PurposeRelevance, nil, 0.8, 0.5, 0.79, BandReview, "nil thresholds use defaults"},
		{PurposeQuality, Thresholds{"quality_apply": 0.9}, 0.9, 0.5, 0.85, BandReview, "topic override wins"},
		{PurposeLink, nil, 0.85, 0.6, 0.85, BandApply, "apply is inclusive"},
		{PurposeMention, nil, 0.85, 0.6, 0.59, BandIgnore, "below review ignores"},
		{PurposeConcept, nil, 0.7, 0.3, 0.3, BandReview, "concept pair"},
		{PurposeClassify, nil, 0.5, 0.5, 0.49, BandIgnore, "classify fixed pair"},
		{PurposeDuplicate, nil, 0.8, 0.8, 0.8, BandApply, "duplicate single cut"},
		{PurposeFind, nil, 0.5, 0.5, 0.5, BandApply, "find keep"},
		{PurposeOKFType, nil, 0.8, 0.5, 0.6, BandReview, "okf type"},
		{PurposeContractCheck, nil, 0.7, 0.7, 0.69, BandIgnore, "contract conflict"},
		{Purpose("other"), nil, 0.8, 0.5, 0.9, BandApply, "unknown purpose"},
	}
	for _, tc := range testCases {
		t.Run("Should band "+string(tc.purpose)+": "+tc.description, func(t *testing.T) {
			t.Parallel()
			apply, review := tc.thresholds.For(tc.purpose)
			if apply != tc.apply || review != tc.rev {
				t.Fatalf("For = (%v, %v), want (%v, %v)", apply, review, tc.apply, tc.rev)
			}
			if got := BandFor(tc.value, apply, review); got != tc.wantBand {
				t.Fatalf("BandFor(%v) = %s, want %s", tc.value, got, tc.wantBand)
			}
		})
	}
}

func TestRedact(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		gone     []string
		kept     []string
		redacted bool
	}{
		{name: "OpenRouter key", input: "key sk-or-v1-abcdef0123456789abcdef here", gone: []string{"abcdef0123456789"}, kept: []string{"key", "here"}, redacted: true},
		{name: "OpenAI project key", input: "sk-proj-AbCdEf0123456789XyZ", gone: []string{"AbCdEf0123456789"}, redacted: true},
		{name: "GitHub token", input: "token ghp_abcdefghijklmnopqrstuvwxyz0123456789", gone: []string{"ghp_abc"}, redacted: true},
		{name: "Slack token", input: "xoxb-1234567890-abcdefghij", gone: []string{"1234567890-abcdefghij"}, redacted: true},
		{name: "AWS key id", input: "id AKIAIOSFODNN7EXAMPLE end", gone: []string{"AKIAIOSFODNN7EXAMPLE"}, redacted: true},
		{name: "bearer token", input: "Authorization: Bearer abcdefghijklmnop0123456789", gone: []string{"abcdefghijklmnop0123456789"}, kept: []string{"Bearer"}, redacted: true},
		{name: "PEM block", input: "a\n-----BEGIN RSA PRIVATE KEY-----\nMIIBOg\n-----END RSA PRIVATE KEY-----\nb", gone: []string{"MIIBOg", "BEGIN RSA"}, kept: []string{"a\n", "\nb"}, redacted: true},
		{name: "assignment", input: `api_key="abcd1234efgh5678" rest`, gone: []string{"abcd1234efgh5678"}, kept: []string{"api_key=", "rest"}, redacted: true},
		{name: "password assignment", input: "password: correcthorse", gone: []string{"correcthorse"}, redacted: true},
		{name: "business text stays", input: "token=invoice and a sketch of the design", kept: []string{"token=invoice", "sketch"}},
	}
	for _, tc := range testCases {
		t.Run("Should handle "+tc.name, func(t *testing.T) {
			t.Parallel()
			got := RedactString(tc.input)
			for _, secret := range tc.gone {
				if strings.Contains(got, secret) {
					t.Errorf("%q still contains %q", got, secret)
				}
			}
			for _, text := range tc.kept {
				if !strings.Contains(got, text) {
					t.Errorf("%q lost %q", got, text)
				}
			}
			if strings.Contains(got, Redacted) != tc.redacted {
				t.Errorf("%q redacted = %v, want %v", got, !tc.redacted, tc.redacted)
			}
		})
	}

	t.Run("Should redact nested values and secret-named fields", func(t *testing.T) {
		t.Parallel()
		type doc struct {
			Title  string            `json:"title"`
			Fields map[string]string `json:"fields"`
			List   []string          `json:"list"`
		}
		got := Redact(doc{
			Title:  "notes",
			Fields: map[string]string{"api_key": "short", "password": "pw", "note": "fine"},
			List:   []string{"sk-or-v1-0000000000000000000000"},
		})
		data, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		want := `{"fields":{"api_key":"[REDACTED]","note":"fine","password":"[REDACTED]"},"list":["[REDACTED]"],"title":"notes"}`
		if string(data) != want {
			t.Fatalf("Redact = %s, want %s", data, want)
		}
		if Redact("sk-or-v1-0000000000000000000000") != Redacted {
			t.Fatal("a string in must return a redacted string")
		}
	})
}

func TestCanonicalJSONIsStable(t *testing.T) {
	t.Parallel()

	a, err := HashJSON(map[string]any{"b": 1, "a": []any{"<x>", 2.5}})
	if err != nil {
		t.Fatal(err)
	}
	b, err := HashJSON(struct {
		A []any `json:"a"`
		B int   `json:"b"`
	}{A: []any{"<x>", 2.5}, B: 1})
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatal("equal values must hash equally regardless of key order or Go type")
	}
	data, _ := CanonicalJSON(map[string]any{"z": "<a&b>", "a": json.Number("1.50")})
	if string(data) != `{"a":1.50,"z":"<a&b>"}` {
		t.Fatalf("CanonicalJSON = %s", data)
	}
	if _, err := CanonicalJSON(map[string]any{"f": func() {}}); err == nil {
		t.Fatal("non-JSON values must fail")
	}
}

func TestBudget(t *testing.T) {
	t.Parallel()

	budget := NewBudget(0.01)
	budget.Charge(0.004)
	budget.Charge(-1)
	if budget.Spent() != 0.004 || budget.Exceeded() || budget.Remaining() != 0.006 || budget.Limit() != 0.01 {
		t.Fatalf("spent %v remaining %v exceeded %v", budget.Spent(), budget.Remaining(), budget.Exceeded())
	}
	budget.Charge(0.007)
	if !budget.Exceeded() || budget.Remaining() != 0 {
		t.Fatalf("spent %v remaining %v", budget.Spent(), budget.Remaining())
	}
	if !NewBudget(0).Exceeded() {
		t.Fatal("a zero budget allows no calls (cache-only run)")
	}
	var unlimited *Budget
	unlimited.Charge(5)
	if unlimited.Exceeded() || unlimited.Spent() != 0 {
		t.Fatal("a nil budget is unlimited")
	}
}

func TestReceiptsRetryLoadAfterReadError(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := ReceiptsPath(root)
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	store := NewReceipts()
	if err := store.Append(root, Receipt{Key: "new", Status: StatusDecided}); err == nil {
		t.Fatal("append to an unreadable log must fail")
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"key":"existing","status":"decided"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	row, ok, err := store.Lookup(root, "existing")
	if err != nil || !ok || row.Key != "existing" {
		t.Fatalf("cache failed to reload repaired log: %+v, %v, %v", row, ok, err)
	}
}

func TestParseCostAndRetryAfter(t *testing.T) {
	t.Parallel()

	costs := []struct {
		raw  string
		want *float64
	}{
		{`0.0002`, new(0.0002)},
		{`"0.0003"`, new(0.0003)},
		{`0`, new(float64(0))},
		{``, nil},
		{`null`, nil},
		{`-1`, nil},
		{`"abc"`, nil},
	}
	for _, tc := range costs {
		got := ParseCost(json.RawMessage(tc.raw))
		if (got == nil) != (tc.want == nil) || (got != nil && *got != *tc.want) {
			t.Errorf("ParseCost(%q) = %v, want %v", tc.raw, got, tc.want)
		}
	}

	now := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	waits := []struct {
		header http.Header
		want   time.Duration
		ok     bool
	}{
		{http.Header{"Retry-After": {"3"}}, 3 * time.Second, true},
		{http.Header{"Retry-After-Ms": {"250"}}, 250 * time.Millisecond, true},
		{http.Header{"Retry-After": {now.Add(4 * time.Second).Format(http.TimeFormat)}}, 4 * time.Second, true},
		{http.Header{"Retry-After": {"600"}}, maxRetryAfter, true},
		{http.Header{"Retry-After": {"1e20"}}, maxRetryAfter, true},
		{http.Header{"Retry-After-Ms": {"1e20"}}, maxRetryAfter, true},
		{http.Header{"Retry-After": {"soon"}}, 0, false},
		{http.Header{}, 0, false},
	}
	for _, tc := range waits {
		got, ok := retryAfter(tc.header, now)
		if got != tc.want || ok != tc.ok {
			t.Errorf("retryAfter(%v) = %v, %v; want %v, %v", tc.header, got, ok, tc.want, tc.ok)
		}
	}
	for retry := 1; retry <= 6; retry++ {
		if d := backoff(retry); d < 375*time.Millisecond || d > backoffCap*5/4 {
			t.Errorf("backoff(%d) = %v out of range", retry, d)
		}
	}
}
