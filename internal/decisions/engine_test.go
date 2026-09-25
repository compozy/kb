package decisions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/questions"
)

// decisionRequest is what the fake server decodes from each call.
type decisionRequest struct {
	Model     string                 `json:"model"`
	State     json.RawMessage        `json:"state"`
	Questions map[string]questions.Q `json:"questions"`

	raw    []byte
	header http.Header
}

// fakeServer replays recorded receipt shapes. respond is called with the
// 0-based call index; nil respond answers every question validly.
type fakeServer struct {
	t       *testing.T
	server  *httptest.Server
	respond func(w http.ResponseWriter, r *http.Request, call int, req decisionRequest)

	mu          sync.Mutex
	requests    []decisionRequest
	inflight    atomic.Int32
	maxInflight atomic.Int32
}

func newFakeServer(t *testing.T, respond func(w http.ResponseWriter, r *http.Request, call int, req decisionRequest)) *fakeServer {
	t.Helper()
	fake := &fakeServer{t: t, respond: respond}
	fake.server = httptest.NewServer(http.HandlerFunc(fake.serve))
	t.Cleanup(fake.server.Close)
	return fake
}

func (f *fakeServer) serve(w http.ResponseWriter, r *http.Request) {
	current := f.inflight.Add(1)
	defer f.inflight.Add(-1)
	for {
		seen := f.maxInflight.Load()
		if current <= seen || f.maxInflight.CompareAndSwap(seen, current) {
			break
		}
	}
	if r.URL.Path != "/api/alpha/decisions" || r.Method != http.MethodPost {
		http.Error(w, "unexpected route "+r.Method+" "+r.URL.Path, http.StatusNotFound)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req decisionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.raw = body
	req.header = r.Header.Clone()
	f.mu.Lock()
	call := len(f.requests)
	f.requests = append(f.requests, req)
	f.mu.Unlock()
	if f.respond == nil {
		writeJSON(w, http.StatusOK, validResponse(req))
		return
	}
	f.respond(w, r, call, req)
}

func (f *fakeServer) calls() []decisionRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.requests)
}

func (f *fakeServer) url() string { return f.server.URL + "/api" }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

// validAnswers answers every question with a well-formed receipt shape:
// noul 0.9; choice = first option (sorted) at 0.9; score = top level.
func validAnswers(req decisionRequest) map[string]any {
	answers := map[string]any{}
	for id, q := range req.Questions {
		switch q.Type {
		case questions.TypeNoul:
			answers[id] = map[string]any{"type": "noul", "noul": 0.9}
		case questions.TypeChoice:
			options := make([]string, 0)
			for option := range q.Criteria.(map[string]any) {
				options = append(options, option)
			}
			slices.Sort(options)
			probabilities := map[string]float64{}
			for index, option := range options {
				if index == 0 {
					probabilities[option] = 0.9
				} else {
					probabilities[option] = 0.1 / float64(len(options)-1)
				}
			}
			answers[id] = map[string]any{"type": "choice", "choice": options[0], "probabilities": probabilities, "confidence": 0.7}
		case questions.TypeScore:
			levels := len(q.Criteria.([]any))
			probabilities := map[string]float64{}
			for level := range levels {
				probabilities[strconv.Itoa(level)] = 0
			}
			probabilities[strconv.Itoa(levels-1)] = 1
			answers[id] = map[string]any{"type": "score", "score": float64(levels - 1), "probabilities": probabilities, "confidence": 0.95}
		}
	}
	return answers
}

func validResponse(req decisionRequest) map[string]any {
	return map[string]any{
		"id":       "gen-dec-1",
		"model":    "typesafe/jev-1.13-20260917",
		"provider": "TypeSafe",
		"answers":  validAnswers(req),
		"usage":    map[string]any{"input_tokens": 1000, "output_tokens": 10, "cost": 0.000042},
	}
}

type sleepRecorder struct {
	mu     sync.Mutex
	sleeps []time.Duration
}

func (s *sleepRecorder) sleep(ctx context.Context, d time.Duration) error {
	s.mu.Lock()
	s.sleeps = append(s.sleeps, d)
	s.mu.Unlock()
	return ctx.Err()
}

func (s *sleepRecorder) recorded() []time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.sleeps)
}

type engineSetup struct {
	cfg      func(*config.DecisionsConfig)
	budget   *Budget
	receipts *Receipts
}

func newTestEngine(t *testing.T, fake *fakeServer, setup engineSetup) (*Engine, *sleepRecorder) {
	t.Helper()
	cfg := config.Default().Decisions
	if setup.cfg != nil {
		setup.cfg(&cfg)
	}
	recorder := &sleepRecorder{}
	engine, err := New(Options{
		Config:   cfg,
		APIKey:   "sk-or-v1-test-key-0000000000000000",
		APIURL:   fake.url(),
		Budget:   setup.budget,
		Receipts: setup.receipts,
		Sleep:    recorder.sleep,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return engine, recorder
}

func qualityRequest(t *testing.T, root string) Request {
	t.Helper()
	quality := questions.MustLoad("quality")
	classify := questions.MustLoad("classify")
	return Request{
		Topic:   TopicRef{Slug: "jev", Root: root, Contract: "contract-hash"},
		Purpose: PurposeQuality,
		Subject: "raw/articles/a.md",
		Bank:    quality,
		State: map[string]any{"document": map[string]any{
			"title": "A", "excerpt": "Full text of the article.",
			"provenance": map[string]any{"path": "raw/articles/a.md"},
		}},
		Questions: []questions.Q{
			quality.MustQuestion("paywall_or_login", nil),
			classify.MustQuestion("kind", nil),
			classify.MustQuestion("depth", nil),
		},
	}
}

func TestNewRequiresAPIKey(t *testing.T) {
	t.Parallel()

	_, err := New(Options{Config: config.Default().Decisions, APIKey: "  "})
	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("err = %v, want ErrNotConfigured", err)
	}
	if !strings.Contains(err.Error(), "OPENROUTER_API_KEY is not set") || !strings.Contains(err.Error(), "decision model is required") {
		t.Fatalf("message must name the missing key and the requirement: %q", err)
	}
	if _, err := New(Options{Config: config.DecisionsConfig{}, APIKey: "k"}); err == nil {
		t.Fatal("an invalid config must be rejected")
	}
}

func TestDecideValidAnswersWritesReceipt(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	engine, _ := newTestEngine(t, fake, engineSetup{})
	root := t.TempDir()
	req := qualityRequest(t, root)

	result, err := engine.Decide(context.Background(), req)
	if err != nil {
		t.Fatalf("Decide: %v", err)
	}
	if result.CacheHit || result.CostUSD != 0.000042 || result.ReceiptID == "" {
		t.Fatalf("unexpected result metadata: %+v", result)
	}

	t.Run("Should decide and band every answer type", func(t *testing.T) {
		noul := result.Answers["paywall_or_login"]
		if !noul.Decided() || *noul.Noul != 0.9 || noul.Band != BandApply || noul.Receipt != result.ReceiptID {
			t.Errorf("noul answer = %+v", noul)
		}
		choice := result.Answers["kind"]
		if !choice.Decided() || choice.Choice != "announcement_or_news" || *choice.Confidence != 0.7 || choice.Band != BandReview {
			t.Errorf("choice answer = %+v", choice)
		}
		score := result.Answers["depth"]
		if !score.Decided() || *score.Score != 3 || score.Band != BandApply || len(score.Probs) != 4 {
			t.Errorf("score answer = %+v", score)
		}
	})

	t.Run("Should send the pinned model, auth header and question map", func(t *testing.T) {
		calls := fake.calls()
		if len(calls) != 1 {
			t.Fatalf("calls = %d, want 1", len(calls))
		}
		call := calls[0]
		if call.Model != "typesafe/jev-1.13" {
			t.Errorf("model = %q", call.Model)
		}
		if got := call.header.Get("Authorization"); got != "Bearer sk-or-v1-test-key-0000000000000000" {
			t.Errorf("Authorization = %q", got)
		}
		if got := call.header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q", got)
		}
		if len(call.Questions) != 3 || call.Questions["kind"].Type != questions.TypeChoice {
			t.Errorf("questions = %+v", call.Questions)
		}
		if !strings.HasPrefix(call.Questions["paywall_or_login"].Instructions.(string), questions.Guard) {
			t.Error("guard must lead every question")
		}
	})

	t.Run("Should append one receipt row with the fixed schema", func(t *testing.T) {
		rows, err := LoadReceipts(root)
		if err != nil || len(rows) != 1 {
			t.Fatalf("rows = %d, err = %v", len(rows), err)
		}
		row := rows[0]
		bank := questions.MustLoad("quality")
		if row.Key != result.ReceiptID || row.Purpose != "quality" || row.Subject != "raw/articles/a.md" ||
			row.Bank != "quality" || row.BankVersion != bank.Version || row.BankHash != bank.Hash ||
			row.Contract != "contract-hash" || row.Model != "typesafe/jev-1.13" ||
			row.ReportedModel != "typesafe/jev-1.13-20260917" || row.Route != RouteDecisions ||
			row.ResponseID != "gen-dec-1" || row.Attempts != 1 || row.Status != StatusDecided || row.Reason != "" {
			t.Errorf("row = %+v", row)
		}
		if row.InputTokens == nil || *row.InputTokens != 1000 || row.Cost == nil || *row.Cost != 0.000042 || row.CostUnknown {
			t.Errorf("usage fields = tokens %v cost %v unknown %v", row.InputTokens, row.Cost, row.CostUnknown)
		}
		if !slices.Equal(row.Questions, []string{"depth", "kind", "paywall_or_login"}) || len(row.Answers) != 3 {
			t.Errorf("questions %v answers %d", row.Questions, len(row.Answers))
		}
		raw, err := os.ReadFile(ReceiptsPath(root))
		if err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{`"key"`, `"time"`, `"bank_version"`, `"reported_model"`, `"latency_ms"`, `"input_tokens"`, `"cost_unknown"`, `"answers"`} {
			if !strings.Contains(string(raw), key) {
				t.Errorf("receipt line misses %s", key)
			}
		}
		if strings.Contains(string(raw), `"band"`) {
			t.Error("bands must never be stored in receipts")
		}
	})
}

func TestDecideRejectsInvalidReceipts(t *testing.T) {
	t.Parallel()

	mutate := func(edit func(answers map[string]any)) func(w http.ResponseWriter, r *http.Request, call int, req decisionRequest) {
		return func(w http.ResponseWriter, _ *http.Request, _ int, req decisionRequest) {
			response := validResponse(req)
			edit(response["answers"].(map[string]any))
			writeJSON(w, http.StatusOK, response)
		}
	}
	testCases := []struct {
		name    string
		respond func(w http.ResponseWriter, r *http.Request, call int, req decisionRequest)
		want    map[string]string // question id → "" (decided) or reason
		wantRow Status
	}{
		{
			name: "Should mark every answer invalid when the body is malformed JSON",
			respond: func(w http.ResponseWriter, _ *http.Request, _ int, _ decisionRequest) {
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, `{"answers": {"paywall_or_login": `)
			},
			want:    map[string]string{"paywall_or_login": ReasonInvalidReceipt, "kind": ReasonInvalidReceipt, "depth": ReasonInvalidReceipt},
			wantRow: StatusUndecided,
		},
		{
			name: "Should reject a choice that is not the argmax",
			respond: mutate(func(answers map[string]any) {
				answers["kind"].(map[string]any)["choice"] = "paper"
			}),
			want:    map[string]string{"paywall_or_login": "", "kind": ReasonInvalidReceipt, "depth": ""},
			wantRow: StatusUndecided,
		},
		{
			name: "Should reject probabilities that do not sum to one",
			respond: mutate(func(answers map[string]any) {
				answers["kind"].(map[string]any)["probabilities"].(map[string]float64)["other"] = 0.5
			}),
			want:    map[string]string{"paywall_or_login": "", "kind": ReasonInvalidReceipt, "depth": ""},
			wantRow: StatusUndecided,
		},
		{
			name: "Should reject a score that is not its expectation",
			respond: mutate(func(answers map[string]any) {
				answers["depth"].(map[string]any)["score"] = 1.2
			}),
			want:    map[string]string{"paywall_or_login": "", "kind": "", "depth": ReasonInvalidReceipt},
			wantRow: StatusUndecided,
		},
		{
			name: "Should reject a noul outside [0,1] and a mismatched type",
			respond: mutate(func(answers map[string]any) {
				answers["paywall_or_login"] = map[string]any{"type": "noul", "noul": 1.3}
				answers["depth"] = map[string]any{"type": "noul", "noul": 0.2}
			}),
			want:    map[string]string{"paywall_or_login": ReasonInvalidReceipt, "kind": "", "depth": ReasonInvalidReceipt},
			wantRow: StatusUndecided,
		},
		{
			name: "Should reject a confidence outside [0,1]",
			respond: mutate(func(answers map[string]any) {
				answers["kind"].(map[string]any)["confidence"] = -0.1
			}),
			want:    map[string]string{"paywall_or_login": "", "kind": ReasonInvalidReceipt, "depth": ""},
			wantRow: StatusUndecided,
		},
		{
			name: "Should invalidate the whole set when an answer is missing",
			respond: mutate(func(answers map[string]any) {
				delete(answers, "depth")
			}),
			want:    map[string]string{"paywall_or_login": ReasonInvalidReceipt, "kind": ReasonInvalidReceipt, "depth": ReasonMissingAnswer},
			wantRow: StatusUndecided,
		},
		{
			name: "Should invalidate the whole set when an extra answer is present",
			respond: mutate(func(answers map[string]any) {
				answers["surprise"] = map[string]any{"type": "noul", "noul": 0.5}
			}),
			want:    map[string]string{"paywall_or_login": ReasonInvalidReceipt, "kind": ReasonInvalidReceipt, "depth": ReasonInvalidReceipt},
			wantRow: StatusUndecided,
		},
		{
			name: "Should accept rounding within tolerance and ties at the argmax",
			respond: mutate(func(answers map[string]any) {
				probabilities := answers["kind"].(map[string]any)["probabilities"].(map[string]float64)
				probabilities["announcement_or_news"] = 0.45
				probabilities["article_or_essay"] = 0.45
				probabilities["other"] = 0.11
				for key := range probabilities {
					if key != "announcement_or_news" && key != "article_or_essay" && key != "other" {
						probabilities[key] = 0
					}
				}
			}),
			want:    map[string]string{"paywall_or_login": "", "kind": "", "depth": ""},
			wantRow: StatusDecided,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fake := newFakeServer(t, tc.respond)
			engine, _ := newTestEngine(t, fake, engineSetup{})
			root := t.TempDir()
			result, err := engine.Decide(context.Background(), qualityRequest(t, root))
			if err != nil {
				t.Fatalf("Decide: %v", err)
			}
			for id, reason := range tc.want {
				answer := result.Answers[id]
				if reason == "" {
					if !answer.Decided() {
						t.Errorf("%s = %+v, want decided", id, answer)
					}
					continue
				}
				if answer.Status != StatusUndecided || answer.Reason != reason || answer.Band != "" {
					t.Errorf("%s = %+v, want undecided:%s", id, answer, reason)
				}
			}
			rows, err := LoadReceipts(root)
			if err != nil || len(rows) != 1 || rows[0].Status != tc.wantRow {
				t.Fatalf("rows = %+v, err = %v, want one %s row", rows, err, tc.wantRow)
			}
			if _, cached, _ := engine.Receipts().Lookup(root, rows[0].Key); cached != (tc.wantRow == StatusDecided) {
				t.Errorf("only decided rows may serve as cache (cached=%v)", cached)
			}
		})
	}
}

func TestDecideTransportFailures(t *testing.T) {
	t.Parallel()

	statusThenOK := func(status int, header map[string]string, failures int) func(w http.ResponseWriter, r *http.Request, call int, req decisionRequest) {
		return func(w http.ResponseWriter, _ *http.Request, call int, req decisionRequest) {
			if call < failures {
				for key, value := range header {
					w.Header().Set(key, value)
				}
				writeJSON(w, status, map[string]any{"error": map[string]any{"code": status, "message": "upstream says no"}})
				return
			}
			writeJSON(w, http.StatusOK, validResponse(req))
		}
	}
	testCases := []struct {
		name         string
		respond      func(w http.ResponseWriter, r *http.Request, call int, req decisionRequest)
		cfg          func(*config.DecisionsConfig)
		wantReason   string // "" = decided
		wantCalls    int
		wantMinSleep time.Duration
	}{
		{
			name:         "Should honour Retry-After on 429 and then succeed",
			respond:      statusThenOK(http.StatusTooManyRequests, map[string]string{"Retry-After": "2"}, 1),
			wantCalls:    2,
			wantMinSleep: 2 * time.Second,
		},
		{
			name:         "Should honour retry-after-ms on 429",
			respond:      statusThenOK(http.StatusTooManyRequests, map[string]string{"retry-after-ms": "1500"}, 1),
			wantCalls:    2,
			wantMinSleep: 1500 * time.Millisecond,
		},
		{
			name:      "Should retry a 503 and then succeed",
			respond:   statusThenOK(http.StatusServiceUnavailable, nil, 1),
			wantCalls: 2,
		},
		{
			name:       "Should give up after the configured retries on 503",
			respond:    statusThenOK(http.StatusServiceUnavailable, nil, 10),
			wantReason: ReasonRetries,
			wantCalls:  3,
		},
		{
			name:       "Should give up on 529 with zero retries",
			respond:    statusThenOK(529, nil, 10),
			cfg:        func(c *config.DecisionsConfig) { c.Retries = 0 },
			wantReason: ReasonRetries,
			wantCalls:  1,
		},
		{
			name:       "Should not retry a 400",
			respond:    statusThenOK(http.StatusBadRequest, nil, 10),
			wantReason: ReasonHTTPStatus,
			wantCalls:  1,
		},
		{
			name:       "Should not retry a 413",
			respond:    statusThenOK(http.StatusRequestEntityTooLarge, nil, 10),
			wantReason: ReasonHTTPStatus,
			wantCalls:  1,
		},
		{
			name: "Should report a hang past the deadline as a timeout",
			respond: func(_ http.ResponseWriter, r *http.Request, _ int, _ decisionRequest) {
				select {
				case <-r.Context().Done():
				case <-time.After(5 * time.Second):
				}
			},
			cfg:        func(c *config.DecisionsConfig) { c.Deadline = "50ms"; c.Retries = 1 },
			wantReason: ReasonTimeout,
			wantCalls:  2,
		},
		{
			name: "Should decode a body preceded by keep-alive whitespace",
			respond: func(w http.ResponseWriter, _ *http.Request, _ int, req decisionRequest) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, "\n \n\t")
				w.(http.Flusher).Flush()
				time.Sleep(20 * time.Millisecond)
				_, _ = io.WriteString(w, " \n")
				_ = json.NewEncoder(w).Encode(validResponse(req))
			},
			wantCalls: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fake := newFakeServer(t, tc.respond)
			engine, recorder := newTestEngine(t, fake, engineSetup{cfg: tc.cfg})
			root := t.TempDir()
			result, err := engine.Decide(context.Background(), qualityRequest(t, root))
			if err != nil {
				t.Fatalf("Decide: %v", err)
			}
			for id, answer := range result.Answers {
				if tc.wantReason == "" && !answer.Decided() {
					t.Errorf("%s = %+v, want decided", id, answer)
				}
				if tc.wantReason != "" && (answer.Status != StatusUndecided || answer.Reason != tc.wantReason) {
					t.Errorf("%s = %+v, want undecided:%s", id, answer, tc.wantReason)
				}
			}
			if got := len(fake.calls()); got != tc.wantCalls {
				t.Errorf("calls = %d, want %d", got, tc.wantCalls)
			}
			if tc.wantMinSleep > 0 && !slices.ContainsFunc(recorder.recorded(), func(d time.Duration) bool { return d >= tc.wantMinSleep-100*time.Millisecond }) {
				t.Errorf("sleeps = %v, want one ≥ %v", recorder.recorded(), tc.wantMinSleep)
			}
			rows, err := LoadReceipts(root)
			if err != nil || len(rows) != 1 || rows[0].Attempts != tc.wantCalls {
				t.Fatalf("rows = %+v, err = %v, want one row with %d attempts", rows, err, tc.wantCalls)
			}
			if tc.wantReason != "" && (rows[0].Status != StatusUndecided || rows[0].Reason != tc.wantReason) {
				t.Errorf("row status = %s:%s", rows[0].Status, rows[0].Reason)
			}
		})
	}
}

func TestDecideAuthFailuresAreFatalAndSticky(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusForbidden} {
		t.Run(fmt.Sprintf("Should stop the run on HTTP %d", status), func(t *testing.T) {
			t.Parallel()
			fake := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request, _ int, _ decisionRequest) {
				writeJSON(w, status, map[string]any{"error": map[string]any{"code": status, "message": "Insufficient credits"}})
			})
			engine, _ := newTestEngine(t, fake, engineSetup{})
			_, err := engine.Decide(context.Background(), qualityRequest(t, t.TempDir()))
			if !errors.Is(err, ErrAuth) || !strings.Contains(err.Error(), strconv.Itoa(status)) || !strings.Contains(err.Error(), "Insufficient credits") {
				t.Fatalf("err = %v, want ErrAuth with status and message", err)
			}
			_, again := engine.Decide(context.Background(), qualityRequest(t, t.TempDir()))
			if !errors.Is(again, ErrAuth) {
				t.Fatalf("later Decide err = %v, want ErrAuth", again)
			}
			if got := len(fake.calls()); got != 1 {
				t.Fatalf("calls = %d, want 1 (later calls must fail fast)", got)
			}
		})
	}
}

func TestDecideStateTooLargeMakesNoCall(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	engine, _ := newTestEngine(t, fake, engineSetup{cfg: func(c *config.DecisionsConfig) { c.MaxStateBytes = 1024 }})
	req := qualityRequest(t, t.TempDir())
	req.State = map[string]any{"document": map[string]any{"excerpt": strings.Repeat("word ", 400)}}
	result, err := engine.Decide(context.Background(), req)
	if err != nil {
		t.Fatalf("Decide: %v", err)
	}
	for id, answer := range result.Answers {
		if answer.Status != StatusNotChecked || answer.Reason != ReasonStateTooLarge {
			t.Errorf("%s = %+v, want not_checked:state_too_large", id, answer)
		}
	}
	if len(fake.calls()) != 0 {
		t.Fatal("an oversize state must not be sent")
	}
}

// TestDecideExcludedSubjectMakesNoCall: decisions.exclude keeps a subject
// out of every call (spec §14): every answer is not_checked:excluded, no
// request reaches the server and no receipt is written.
func TestDecideExcludedSubjectMakesNoCall(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		exclude  []string
		subject  string
		excluded bool
	}{
		{name: "double star folder", exclude: []string{"raw/private/**"}, subject: "raw/private/notes/a.md", excluded: true},
		{name: "single segment star", exclude: []string{"raw/*/secret-*.md"}, subject: "raw/articles/secret-plan.md", excluded: true},
		{name: "leading double star", exclude: []string{"**/drafts/**"}, subject: "wiki/drafts/x.md", excluded: true},
		{name: "no match", exclude: []string{"raw/private/**"}, subject: "raw/articles/a.md"},
		{name: "star stays in its segment", exclude: []string{"raw/*.md"}, subject: "raw/articles/a.md"},
		{name: "non-path subject", exclude: []string{"raw/**"}, subject: "prefetch:url-2026-09-24"},
		{name: "no globs", subject: "raw/private/a.md"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fake := newFakeServer(t, nil)
			engine, _ := newTestEngine(t, fake, engineSetup{})
			root := t.TempDir()
			req := qualityRequest(t, root)
			req.Topic.Exclude = tc.exclude
			req.Subject = tc.subject
			result, err := engine.Decide(context.Background(), req)
			if err != nil {
				t.Fatalf("Decide: %v", err)
			}
			if !tc.excluded {
				if len(fake.calls()) != 1 {
					t.Fatalf("calls = %d, want 1", len(fake.calls()))
				}
				return
			}
			if len(result.Answers) != len(req.Questions) {
				t.Fatalf("answers = %+v", result.Answers)
			}
			for id, answer := range result.Answers {
				if answer.Status != StatusNotChecked || answer.Reason != ReasonExcluded || answer.Decided() {
					t.Errorf("%s = %+v, want not_checked:excluded", id, answer)
				}
			}
			if len(fake.calls()) != 0 {
				t.Fatal("an excluded subject must not be sent")
			}
			if rows, err := LoadReceipts(root); err != nil || len(rows) != 0 {
				t.Fatalf("receipts = %+v, %v", rows, err)
			}
			if summary := engine.Summary(); summary.Calls != 0 {
				t.Fatalf("summary calls = %d", summary.Calls)
			}
		})
	}
}

func TestDecideStopsAtBudget(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	budget := NewBudget(0.00004)
	engine, _ := newTestEngine(t, fake, engineSetup{budget: budget})
	root := t.TempDir()

	first, err := engine.Decide(context.Background(), qualityRequest(t, root))
	if err != nil || !first.Answers["kind"].Decided() {
		t.Fatalf("first Decide = %+v, %v", first, err)
	}
	if !budget.Exceeded() || budget.Spent() != 0.000042 {
		t.Fatalf("budget spent %v exceeded %v", budget.Spent(), budget.Exceeded())
	}
	second := qualityRequest(t, root)
	second.Subject = "raw/articles/b.md"
	second.State = map[string]any{"document": map[string]any{"title": "B"}}
	result, err := engine.Decide(context.Background(), second)
	if err != nil {
		t.Fatalf("Decide: %v", err)
	}
	for id, answer := range result.Answers {
		if answer.Status != StatusUndecided || answer.Reason != ReasonBudget {
			t.Errorf("%s = %+v, want undecided:budget", id, answer)
		}
	}
	if len(fake.calls()) != 1 {
		t.Fatalf("calls = %d, want 1", len(fake.calls()))
	}
	cached, err := engine.Decide(context.Background(), qualityRequest(t, root))
	if err != nil || !cached.CacheHit {
		t.Fatalf("a cache hit must still be served past the budget: %+v, %v", cached, err)
	}
}

func TestDecideCacheHits(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	root := t.TempDir()
	engine, _ := newTestEngine(t, fake, engineSetup{})
	first, err := engine.Decide(context.Background(), qualityRequest(t, root))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Should serve an identical request from the cache as a copy", func(t *testing.T) {
		second, err := engine.Decide(context.Background(), qualityRequest(t, root))
		if err != nil {
			t.Fatal(err)
		}
		if !second.CacheHit || second.CostUSD != 0 || second.ReceiptID != first.ReceiptID {
			t.Fatalf("second = %+v", second)
		}
		second.Answers["kind"].Probs["paper"] = 0.99
		third, _ := engine.Decide(context.Background(), qualityRequest(t, root))
		if third.Answers["kind"].Probs["paper"] == 0.99 {
			t.Fatal("cache hits must be deep copies")
		}
	})

	t.Run("Should reload receipts lazily and re-band with current thresholds", func(t *testing.T) {
		fresh, _ := newTestEngine(t, fake, engineSetup{})
		req := qualityRequest(t, root)
		req.Topic.Thresholds = Thresholds{"quality_apply": 0.95}
		result, err := fresh.Decide(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		if !result.CacheHit {
			t.Fatal("a new engine must reuse receipts from disk")
		}
		if got := result.Answers["paywall_or_login"].Band; got != BandReview {
			t.Fatalf("band = %s, want review after raising quality_apply", got)
		}
	})

	t.Run("Should miss the cache when the contract changes", func(t *testing.T) {
		req := qualityRequest(t, root)
		req.Topic.Contract = "other-contract"
		result, err := engine.Decide(context.Background(), req)
		if err != nil || result.CacheHit {
			t.Fatalf("result = %+v, %v", result, err)
		}
	})

	if got := len(fake.calls()); got != 2 {
		t.Fatalf("calls = %d, want 2", got)
	}
	summary := engine.Summary()
	if summary.Calls != 2 || summary.CacheHits != 2 {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestDecideBatchesLargeQuestionSets(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	engine, _ := newTestEngine(t, fake, engineSetup{})
	concept := questions.MustLoad("concept")
	list := make([]questions.Q, 0, 50)
	for index := range 50 {
		list = append(list, concept.MustQuestion("mentions_concept_{id}", map[string]string{"id": fmt.Sprintf("c%02d", index)}))
	}
	result, err := engine.Decide(context.Background(), Request{
		Topic: TopicRef{Root: t.TempDir()}, Purpose: PurposeConcept, Bank: concept,
		State: map[string]any{"document": map[string]any{"title": "x"}}, Questions: list,
	})
	if err != nil {
		t.Fatal(err)
	}
	calls := fake.calls()
	if len(calls) != 2 || len(calls[0].Questions)+len(calls[1].Questions) != 50 || len(calls[0].Questions) != MaxQuestionsPerRequest {
		t.Fatalf("batches = %d", len(calls))
	}
	if len(result.Answers) != 50 || result.CostUSD != 0.000084 {
		t.Fatalf("answers %d cost %v", len(result.Answers), result.CostUSD)
	}
	if result.Answers["mentions_concept_c00"].Receipt == result.Answers["mentions_concept_c49"].Receipt {
		t.Fatal("answers of different batches must carry their own receipt keys")
	}
	if string(calls[0].State) != string(calls[1].State) {
		t.Fatal("every batch must carry the same state")
	}
}

func TestDecideRedactsSecretsBeforeSending(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	engine, _ := newTestEngine(t, fake, engineSetup{})
	req := qualityRequest(t, t.TempDir())
	secrets := []string{
		"sk-or-v1-0123456789abcdef0123456789abcdef",
		"ghp_abcdefghijklmnopqrstuvwxyz0123456789",
		"AKIAIOSFODNN7EXAMPLE",
		"MIIEvQIBADANBgkqhkiG9w0BAQEFAASC",
		"supersecretvalue123",
	}
	req.State = map[string]any{"document": map[string]any{
		"title": "Setup notes",
		"excerpt": "Use key " + secrets[0] + " or " + secrets[1] + " with " + secrets[2] +
			"\n-----BEGIN PRIVATE KEY-----\n" + secrets[3] + "\n-----END PRIVATE KEY-----\napi_key=" + secrets[4],
		"password": "hunter2",
	}}
	if _, err := engine.Decide(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	body := string(fake.calls()[0].raw)
	for _, secret := range append(secrets, "hunter2") {
		if strings.Contains(body, secret) {
			t.Errorf("secret %q reached the decision endpoint", secret)
		}
	}
	if !strings.Contains(body, Redacted) || !strings.Contains(body, "Setup notes") {
		t.Error("redaction must keep ordinary text and mark removed secrets")
	}
}

func TestDecideFatalAndInvalidRequests(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	engine, _ := newTestEngine(t, fake, engineSetup{})

	t.Run("Should return the context error when canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := engine.Decide(ctx, qualityRequest(t, t.TempDir())); !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("Should reject duplicate question ids", func(t *testing.T) {
		req := qualityRequest(t, t.TempDir())
		req.Questions = append(req.Questions, req.Questions[0])
		if _, err := engine.Decide(context.Background(), req); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("Should reject a request without questions", func(t *testing.T) {
		req := qualityRequest(t, t.TempDir())
		req.Questions = nil
		if _, err := engine.Decide(context.Background(), req); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("err = %v", err)
		}
	})
	if len(fake.calls()) != 0 {
		t.Fatal("invalid requests must not reach the network")
	}
}

func TestDecideRespectsConcurrency(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request, _ int, req decisionRequest) {
		time.Sleep(20 * time.Millisecond)
		writeJSON(w, http.StatusOK, validResponse(req))
	})
	engine, _ := newTestEngine(t, fake, engineSetup{cfg: func(c *config.DecisionsConfig) { c.Concurrency = 2 }})
	if engine.Concurrency() != 2 || engine.Model() != "typesafe/jev-1.13" {
		t.Fatalf("Concurrency %d Model %q", engine.Concurrency(), engine.Model())
	}
	root := t.TempDir()
	var wg sync.WaitGroup
	for index := range 8 {
		wg.Go(func() {
			req := qualityRequest(t, root)
			req.State = map[string]any{"n": index}
			if _, err := engine.Decide(context.Background(), req); err != nil {
				t.Errorf("Decide: %v", err)
			}
		})
	}
	wg.Wait()
	if got := fake.maxInflight.Load(); got > 2 {
		t.Fatalf("max in-flight = %d, want ≤ 2", got)
	}
	rows, err := LoadReceipts(root)
	if err != nil || len(rows) != 8 {
		t.Fatalf("rows = %d, err = %v (appends must be whole lines)", len(rows), err)
	}
}

func TestSummaryLines(t *testing.T) {
	t.Parallel()

	fake := newFakeServer(t, nil)
	engine, _ := newTestEngine(t, fake, engineSetup{cfg: func(c *config.DecisionsConfig) { c.MaxStateBytes = 1024 }})
	root := t.TempDir()
	if _, err := engine.Decide(context.Background(), qualityRequest(t, root)); err != nil {
		t.Fatal(err)
	}
	if _, err := engine.Decide(context.Background(), qualityRequest(t, root)); err != nil {
		t.Fatal(err)
	}
	big := qualityRequest(t, root)
	big.State = strings.Repeat("x", 2000)
	if _, err := engine.Decide(context.Background(), big); err != nil {
		t.Fatal(err)
	}
	summary := engine.Summary()
	summary.Generation = GenerationSummary{Calls: 1, CostUSD: 0.001, Fallbacks: 1}
	got := strings.Join(summary.Lines(), "\n")
	for _, want := range []string{
		"decisions: 1 call, 1 cache hit, US$ 0.0000",
		"quality: 6 decided (apply 4, review 2, ignore 0), 3 not_checked:state_too_large",
		"generation: 1 call, 0 cache hits, US$ 0.0010, 1 fallback, 0 failures",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("summary lines miss %q:\n%s", want, got)
		}
	}
}

func TestWaitPauseRereadsAnExtendedDeadline(t *testing.T) {
	t.Parallel()

	clock := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	var engine *Engine
	var slept []time.Duration
	engine, err := New(Options{
		Config: config.Default().Decisions,
		APIKey: "k",
		Now:    func() time.Time { return clock },
		Sleep: func(ctx context.Context, d time.Duration) error {
			slept = append(slept, d)
			clock = clock.Add(d)
			if len(slept) == 1 {
				// Another worker hits a longer 429 while this one sleeps.
				engine.pause(700 * time.Millisecond)
			}
			return ctx.Err()
		},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	start := clock
	engine.pause(200 * time.Millisecond)
	if err := engine.waitPause(context.Background()); err != nil {
		t.Fatalf("waitPause: %v", err)
	}
	if waited := clock.Sub(start); waited < 900*time.Millisecond {
		t.Fatalf("waited %v (sleeps %v), want at least the extended 900ms", waited, slept)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	engine.pause(time.Second)
	if err := engine.waitPause(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("waitPause on a cancelled context = %v, want context.Canceled", err)
	}
}

func TestDecideHonoursExtendedSharedPause(t *testing.T) {
	t.Parallel()

	const (
		shortWait = 200 * time.Millisecond
		longWait  = 800 * time.Millisecond
	)
	var (
		mu       sync.Mutex
		extendAt time.Time
		retries  []time.Time
		arrived  atomic.Int32
	)
	bothArrived := make(chan struct{})
	rateLimited := func(w http.ResponseWriter, wait time.Duration) {
		w.Header().Set("retry-after-ms", strconv.FormatInt(wait.Milliseconds(), 10))
		writeJSON(w, http.StatusTooManyRequests, map[string]any{"error": map[string]any{"code": 429, "message": "slow down"}})
	}
	fake := newFakeServer(t, func(w http.ResponseWriter, _ *http.Request, call int, req decisionRequest) {
		arrival := time.Now()
		if call >= 2 {
			mu.Lock()
			retries = append(retries, arrival)
			mu.Unlock()
			writeJSON(w, http.StatusOK, validResponse(req))
			return
		}
		if arrived.Add(1) == 2 {
			close(bothArrived)
		}
		select {
		case <-bothArrived:
		case <-time.After(5 * time.Second):
			t.Errorf("the first two requests were not in flight together")
		}
		if call == 0 {
			rateLimited(w, shortWait)
			return
		}
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		extendAt = time.Now()
		mu.Unlock()
		rateLimited(w, longWait)
	})
	cfg := config.Default().Decisions
	cfg.Concurrency = 2
	engine, err := New(Options{Config: cfg, APIKey: "sk-or-v1-test-key-0000000000000000", APIURL: fake.url()})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	root := t.TempDir()
	var wg sync.WaitGroup
	for index := range 2 {
		wg.Go(func() {
			req := qualityRequest(t, root)
			req.State = map[string]any{"n": index}
			result, err := engine.Decide(context.Background(), req)
			if err != nil {
				t.Errorf("Decide: %v", err)
				return
			}
			for id, answer := range result.Answers {
				if !answer.Decided() {
					t.Errorf("%s = %+v, want decided after the pause", id, answer)
				}
			}
		})
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if extendAt.IsZero() || len(retries) != 2 {
		t.Fatalf("extendAt=%v retries=%d, want the pause extended and both requests retried", extendAt, len(retries))
	}
	deadline := extendAt.Add(longWait)
	for _, at := range retries {
		if at.Before(deadline) {
			t.Errorf("a retry reached the server %v before the extended pause ended", deadline.Sub(at))
		}
	}
}
