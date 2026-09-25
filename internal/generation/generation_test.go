package generation

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/decisions"
)

const summarySchema = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["summary", "entities"],
  "properties": {
    "summary": {"type": "string"},
    "entities": {"type": "array", "items": {"type": "string"}}
  }
}`

type chatRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	ResponseFormat struct {
		Type       string `json:"type"`
		JSONSchema struct {
			Name   string          `json:"name"`
			Strict bool            `json:"strict"`
			Schema json.RawMessage `json:"schema"`
		} `json:"json_schema"`
	} `json:"response_format"`
	Reasoning   map[string]any `json:"reasoning"`
	Temperature *float64       `json:"temperature"`

	raw string
}

type fakeChat struct {
	server  *httptest.Server
	respond func(w http.ResponseWriter, call int, req chatRequest)

	mu       sync.Mutex
	requests []chatRequest
}

func newFakeChat(t *testing.T, respond func(w http.ResponseWriter, call int, req chatRequest)) *fakeChat {
	t.Helper()
	fake := &fakeChat{respond: respond}
	fake.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer test-key" {
			http.Error(w, "unexpected request", http.StatusNotFound)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.raw = string(body)
		fake.mu.Lock()
		call := len(fake.requests)
		fake.requests = append(fake.requests, req)
		fake.mu.Unlock()
		fake.respond(w, call, req)
	}))
	t.Cleanup(fake.server.Close)
	return fake
}

func (f *fakeChat) calls() []chatRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.requests)
}

func chatReply(w http.ResponseWriter, model, content string, reasoningTokens int, cost any) {
	usage := map[string]any{"prompt_tokens": 500, "completion_tokens": 40, "completion_tokens_details": map[string]any{"reasoning_tokens": reasoningTokens}}
	if cost != nil {
		usage["cost"] = cost
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":      "gen-1",
		"model":   model,
		"choices": []any{map[string]any{"message": map[string]any{"role": "assistant", "content": content}}},
		"usage":   usage,
	})
}

const goodOutput = `{"summary": "A short summary.", "entities": ["MCP", "A2A"]}`

func newTestClient(t *testing.T, fake *fakeChat, budget *decisions.Budget, receipts *decisions.Receipts) *Client {
	t.Helper()
	cfg := config.Default().Generation
	client, err := New(Options{
		Config:   cfg,
		APIKey:   "test-key",
		APIURL:   fake.server.URL + "/api",
		Budget:   budget,
		Receipts: receipts,
		Sleep:    func(ctx context.Context, _ time.Duration) error { return ctx.Err() },
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return client
}

func summaryRequest(root string) Request {
	return Request{
		Topic:      decisions.TopicRef{Slug: "jev", Root: root, Contract: "c1"},
		Kind:       "summary",
		Subject:    "raw/articles/a.md",
		System:     "Write a summary.",
		Prompt:     "Document: agents talk over MCP and A2A.",
		SchemaName: "summary",
		Schema:     json.RawMessage(summarySchema),
	}
}

func TestNewRequiresAPIKey(t *testing.T) {
	t.Parallel()

	if _, err := New(Options{Config: config.Default().Generation}); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("err = %v, want ErrNotConfigured", err)
	}
}

func TestGenerateSchemaValidOutput(t *testing.T) {
	t.Parallel()

	fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) {
		chatReply(w, req.Model+"-2026", goodOutput, 0, 0.0001)
	})
	budget := decisions.NewBudget(1)
	client := newTestClient(t, fake, budget, nil)
	root := t.TempDir()

	output, err := client.Generate(context.Background(), summaryRequest(root))
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if string(output) != `{"summary":"A short summary.","entities":["MCP","A2A"]}` {
		t.Fatalf("output = %s", output)
	}

	t.Run("Should request strict json_schema output with reasoning disabled", func(t *testing.T) {
		calls := fake.calls()
		if len(calls) != 1 {
			t.Fatalf("calls = %d", len(calls))
		}
		call := calls[0]
		if call.Model != "xiaomi/mimo-v2.6-flash" || call.ResponseFormat.Type != "json_schema" ||
			!call.ResponseFormat.JSONSchema.Strict || call.ResponseFormat.JSONSchema.Name != "summary" ||
			call.Reasoning["enabled"] != false || call.Temperature == nil || *call.Temperature != 0 {
			t.Fatalf("request = %s", call.raw)
		}
		if len(call.Messages) != 2 || call.Messages[0].Role != "system" || call.Messages[1].Content != "Document: agents talk over MCP and A2A." {
			t.Fatalf("messages = %+v", call.Messages)
		}
	})

	t.Run("Should charge the shared budget and log a generate receipt", func(t *testing.T) {
		if budget.Spent() != 0.0001 {
			t.Fatalf("spent = %v", budget.Spent())
		}
		rows, err := decisions.LoadReceipts(root)
		if err != nil || len(rows) != 1 {
			t.Fatalf("rows = %d, %v", len(rows), err)
		}
		row := rows[0]
		if row.Purpose != "generate:summary" || row.Status != decisions.StatusDecided || row.Route != decisions.RouteChatCompletions ||
			row.ReportedModel != "xiaomi/mimo-v2.6-flash-2026" || row.Attempts != 1 || row.Cost == nil || *row.Cost != 0.0001 ||
			string(row.Answers["output"]) != string(output) || row.Contract != "c1" {
			t.Fatalf("row = %+v", row)
		}
	})

	t.Run("Should serve an identical request from the cache", func(t *testing.T) {
		again, err := client.Generate(context.Background(), summaryRequest(root))
		if err != nil || string(again) != string(output) {
			t.Fatalf("again = %s, %v", again, err)
		}
		fresh := newTestClient(t, fake, budget, nil)
		fromDisk, err := fresh.Generate(context.Background(), summaryRequest(root))
		if err != nil || string(fromDisk) != string(output) {
			t.Fatalf("fromDisk = %s, %v", fromDisk, err)
		}
		if len(fake.calls()) != 1 {
			t.Fatalf("calls = %d, want 1", len(fake.calls()))
		}
		if got := client.Summary(); got.Calls != 1 || got.CacheHits != 1 || got.CostUSD != 0.0001 {
			t.Fatalf("summary = %+v", got)
		}
	})
}

func TestGenerateRejectsAndFallsBack(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		respond     func(w http.ResponseWriter, call int, req chatRequest)
		wantOutput  bool
		wantCalls   int
		wantReason  string
		wantModels  []string
		wantSummary decisions.GenerationSummary
	}{
		{
			name: "Should reject billed reasoning tokens and use the fallback model",
			respond: func(w http.ResponseWriter, _ int, req chatRequest) {
				if req.Model == "xiaomi/mimo-v2.6-flash" {
					chatReply(w, req.Model, goodOutput, 12, 0.0002)
					return
				}
				chatReply(w, req.Model, goodOutput, 0, 0.0001)
			},
			wantOutput:  true,
			wantCalls:   3,
			wantModels:  []string{"xiaomi/mimo-v2.6-flash", "xiaomi/mimo-v2.6-flash", "deepseek/deepseek-v4-flash"},
			wantSummary: decisions.GenerationSummary{Calls: 1, Fallbacks: 1, CostUSD: 0.0005},
		},
		{
			name: "Should fail after invalid JSON on every model",
			respond: func(w http.ResponseWriter, _ int, req chatRequest) {
				chatReply(w, req.Model, `{"summary": "unterminated`, 0, 0.0001)
			},
			wantCalls:   4,
			wantReason:  ReasonInvalidOutput,
			wantSummary: decisions.GenerationSummary{Calls: 1, Fallbacks: 1, Failures: 1, CostUSD: 0.0004},
		},
		{
			name: "Should reject a missing required key and an undeclared key",
			respond: func(w http.ResponseWriter, call int, req chatRequest) {
				if call%2 == 0 {
					chatReply(w, req.Model, `{"summary": "x"}`, 0, 0.0001)
					return
				}
				chatReply(w, req.Model, `{"summary": "x", "entities": [], "extra": 1}`, 0, 0.0001)
			},
			wantCalls:   4,
			wantReason:  ReasonInvalidOutput,
			wantSummary: decisions.GenerationSummary{Calls: 1, Fallbacks: 1, Failures: 1, CostUSD: 0.0004},
		},
		{
			name: "Should reject top-level and item type mismatches",
			respond: func(w http.ResponseWriter, call int, req chatRequest) {
				if call%2 == 0 {
					chatReply(w, req.Model, `{"summary": 3, "entities": []}`, 0, nil)
					return
				}
				chatReply(w, req.Model, `{"summary": "x", "entities": [1]}`, 0, nil)
			},
			wantCalls:   4,
			wantReason:  ReasonInvalidOutput,
			wantSummary: decisions.GenerationSummary{Calls: 1, Fallbacks: 1, Failures: 1, CostUnknown: 1},
		},
		{
			name: "Should retry a 503 and succeed",
			respond: func(w http.ResponseWriter, call int, req chatRequest) {
				if call == 0 {
					w.Header().Set("Retry-After", "1")
					http.Error(w, `{"error":{"message":"busy"}}`, http.StatusServiceUnavailable)
					return
				}
				chatReply(w, req.Model, goodOutput, 0, 0.0001)
			},
			wantOutput:  true,
			wantCalls:   2,
			wantSummary: decisions.GenerationSummary{Calls: 1, CostUSD: 0.0001},
		},
		{
			name: "Should move to the fallback on a model rejection",
			respond: func(w http.ResponseWriter, _ int, req chatRequest) {
				if req.Model == "xiaomi/mimo-v2.6-flash" {
					http.Error(w, `{"error":{"message":"json_schema unsupported"}}`, http.StatusBadRequest)
					return
				}
				chatReply(w, req.Model, goodOutput, 0, 0.0001)
			},
			wantOutput:  true,
			wantCalls:   2,
			wantSummary: decisions.GenerationSummary{Calls: 1, Fallbacks: 1, CostUSD: 0.0001},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fake := newFakeChat(t, tc.respond)
			client := newTestClient(t, fake, decisions.NewBudget(1), nil)
			root := t.TempDir()
			output, err := client.Generate(context.Background(), summaryRequest(root))
			if tc.wantOutput && (err != nil || output == nil) {
				t.Fatalf("Generate = %s, %v", output, err)
			}
			if !tc.wantOutput && (!errors.Is(err, ErrFailed) || !strings.Contains(err.Error(), tc.wantReason)) {
				t.Fatalf("err = %v, want ErrFailed with %s", err, tc.wantReason)
			}
			calls := fake.calls()
			if len(calls) != tc.wantCalls {
				t.Fatalf("calls = %d, want %d", len(calls), tc.wantCalls)
			}
			if tc.wantModels != nil {
				models := make([]string, len(calls))
				for index, call := range calls {
					models[index] = call.Model
				}
				if !slices.Equal(models, tc.wantModels) {
					t.Fatalf("models = %v, want %v", models, tc.wantModels)
				}
			}
			got := client.Summary()
			if got.Calls != tc.wantSummary.Calls || got.Fallbacks != tc.wantSummary.Fallbacks || got.Failures != tc.wantSummary.Failures ||
				got.CostUnknown != tc.wantSummary.CostUnknown || !near(got.CostUSD, tc.wantSummary.CostUSD) {
				t.Fatalf("summary = %+v, want %+v", got, tc.wantSummary)
			}
			rows, err := decisions.LoadReceipts(root)
			if err != nil || len(rows) != 1 || rows[0].Attempts != tc.wantCalls {
				t.Fatalf("rows = %+v, %v", rows, err)
			}
			if !tc.wantOutput && (rows[0].Status != decisions.StatusUndecided || rows[0].Reason != tc.wantReason) {
				t.Fatalf("row = %+v", rows[0])
			}
		})
	}
}

func TestGenerateValidateHookKeepsRejectedOutputOutOfTheCache(t *testing.T) {
	t.Parallel()

	const emptySummary = `{"summary": "", "entities": []}`
	requireSummary := func(raw json.RawMessage) error {
		var out struct {
			Summary string `json:"summary"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			return err
		}
		return ValidateLiteral(out.Summary, 400)
	}

	t.Run("Should retry, fall back, record invalid_output and call again on the next run", func(t *testing.T) {
		t.Parallel()
		var valid atomic.Bool
		fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) {
			if valid.Load() {
				chatReply(w, req.Model, goodOutput, 0, 0.0001)
				return
			}
			chatReply(w, req.Model, emptySummary, 0, 0.0001)
		})
		client := newTestClient(t, fake, decisions.NewBudget(1), nil)
		root := t.TempDir()
		req := summaryRequest(root)
		req.Validate = requireSummary

		_, err := client.Generate(context.Background(), req)
		if !errors.Is(err, ErrFailed) || !errors.Is(err, ErrInvalidOutput) {
			t.Fatalf("err = %v, want ErrFailed wrapping ErrInvalidOutput", err)
		}
		calls := fake.calls()
		models := make([]string, len(calls))
		for index, call := range calls {
			models[index] = call.Model
		}
		wantModels := []string{"xiaomi/mimo-v2.6-flash", "xiaomi/mimo-v2.6-flash", "deepseek/deepseek-v4-flash", "deepseek/deepseek-v4-flash"}
		if !slices.Equal(models, wantModels) {
			t.Fatalf("models = %v, want one retry then the fallback %v", models, wantModels)
		}
		rows, err := decisions.LoadReceipts(root)
		if err != nil || len(rows) != 1 || rows[0].Status != decisions.StatusUndecided || rows[0].Reason != ReasonInvalidOutput {
			t.Fatalf("rows = %+v, %v, want one undecided:invalid_output row", rows, err)
		}

		valid.Store(true)
		output, err := newTestClient(t, fake, decisions.NewBudget(1), nil).Generate(context.Background(), req)
		if err != nil || !strings.Contains(string(output), "A short summary.") {
			t.Fatalf("second run = %s, %v", output, err)
		}
		if got := len(fake.calls()); got != len(wantModels)+1 {
			t.Fatalf("calls = %d, want a new call on the second run", got)
		}
	})

	t.Run("Should not serve a cached output that fails the hook", func(t *testing.T) {
		t.Parallel()
		var valid atomic.Bool
		fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) {
			if valid.Load() {
				chatReply(w, req.Model, goodOutput, 0, 0.0001)
				return
			}
			chatReply(w, req.Model, emptySummary, 0, 0.0001)
		})
		root := t.TempDir()
		// A row decided before the caller validated its literals.
		if _, err := newTestClient(t, fake, decisions.NewBudget(1), nil).Generate(context.Background(), summaryRequest(root)); err != nil {
			t.Fatalf("unvalidated Generate: %v", err)
		}
		valid.Store(true)
		client := newTestClient(t, fake, decisions.NewBudget(1), nil)
		req := summaryRequest(root)
		req.Validate = requireSummary
		output, err := client.Generate(context.Background(), req)
		if err != nil || !strings.Contains(string(output), "A short summary.") {
			t.Fatalf("Generate = %s, %v", output, err)
		}
		if got := client.Summary(); got.CacheHits != 0 || got.Calls != 1 || len(fake.calls()) != 2 {
			t.Fatalf("summary = %+v, calls = %d, want a fresh call instead of the cached rejected output", got, len(fake.calls()))
		}
	})
}

func near(a, b float64) bool {
	diff := a - b
	return diff < 1e-12 && diff > -1e-12
}

func TestGenerateBudgetAuthAndRedaction(t *testing.T) {
	t.Parallel()

	t.Run("Should refuse when the budget is exhausted", func(t *testing.T) {
		t.Parallel()
		fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) {
			chatReply(w, req.Model, goodOutput, 0, 0.002)
		})
		budget := decisions.NewBudget(0.001)
		client := newTestClient(t, fake, budget, nil)
		if _, err := client.Generate(context.Background(), summaryRequest(t.TempDir())); err != nil {
			t.Fatalf("first Generate: %v", err)
		}
		second := summaryRequest(t.TempDir())
		second.Prompt = "another document"
		if _, err := client.Generate(context.Background(), second); !errors.Is(err, ErrBudget) {
			t.Fatalf("err = %v, want ErrBudget", err)
		}
		if len(fake.calls()) != 1 {
			t.Fatalf("calls = %d, want 1", len(fake.calls()))
		}
	})

	t.Run("Should stop the run on 401 and fail fast afterwards", func(t *testing.T) {
		t.Parallel()
		fake := newFakeChat(t, func(w http.ResponseWriter, _ int, _ chatRequest) {
			http.Error(w, `{"error":{"code":401,"message":"No auth credentials found"}}`, http.StatusUnauthorized)
		})
		client := newTestClient(t, fake, decisions.NewBudget(1), nil)
		root := t.TempDir()
		_, err := client.Generate(context.Background(), summaryRequest(root))
		if !errors.Is(err, decisions.ErrAuth) || !strings.Contains(err.Error(), "No auth credentials found") {
			t.Fatalf("err = %v", err)
		}
		if _, err := client.Generate(context.Background(), summaryRequest(t.TempDir())); !errors.Is(err, decisions.ErrAuth) {
			t.Fatalf("second err = %v", err)
		}
		if len(fake.calls()) != 1 {
			t.Fatalf("calls = %d", len(fake.calls()))
		}
		rows, err := decisions.LoadReceipts(root)
		if err != nil || len(rows) != 1 || rows[0].Reason != ReasonHTTPStatus || !rows[0].CostUnknown {
			t.Fatalf("fatal request missing from receipts: %+v, %v", rows, err)
		}
		if summary := client.Summary(); summary.Calls != 1 || summary.CostUnknown != 1 || summary.Failures != 1 {
			t.Fatalf("fatal request missing from summary: %+v", summary)
		}
	})

	t.Run("Should redact secrets in the prompt", func(t *testing.T) {
		t.Parallel()
		fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) {
			chatReply(w, req.Model, goodOutput, 0, 0.0001)
		})
		client := newTestClient(t, fake, decisions.NewBudget(1), nil)
		req := summaryRequest(t.TempDir())
		req.Prompt = "Config: OPENROUTER_API_KEY=sk-or-v1-0123456789abcdef0123456789"
		if _, err := client.Generate(context.Background(), req); err != nil {
			t.Fatal(err)
		}
		if body := fake.calls()[0].raw; strings.Contains(body, "0123456789abcdef") || !strings.Contains(body, decisions.Redacted) {
			t.Fatalf("body = %s", body)
		}
	})

	t.Run("Should refuse an excluded subject without a call", func(t *testing.T) {
		t.Parallel()
		fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) { chatReply(w, req.Model, goodOutput, 0, 0.0001) })
		client := newTestClient(t, fake, decisions.NewBudget(1), nil)
		excluded := summaryRequest(t.TempDir())
		excluded.Topic.Exclude = []string{"raw/articles/**"}
		if _, err := client.Generate(context.Background(), excluded); !errors.Is(err, ErrExcluded) {
			t.Fatalf("err = %v, want ErrExcluded", err)
		}
		if len(fake.calls()) != 0 {
			t.Fatal("an excluded subject must not be sent")
		}
		other := summaryRequest(t.TempDir())
		other.Topic.Exclude = []string{"raw/private/**"}
		if _, err := client.Generate(context.Background(), other); err != nil {
			t.Fatalf("non-excluded Generate: %v", err)
		}
	})

	t.Run("Should reject invalid requests without a call", func(t *testing.T) {
		t.Parallel()
		fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) { chatReply(w, req.Model, goodOutput, 0, 0.0001) })
		client := newTestClient(t, fake, decisions.NewBudget(1), nil)
		for _, mutate := range []func(*Request){
			func(r *Request) { r.Kind = "" },
			func(r *Request) { r.Prompt = " " },
			func(r *Request) { r.SchemaName = "" },
			func(r *Request) { r.Schema = json.RawMessage(`{"type":"array"}`) },
			func(r *Request) { r.Schema = json.RawMessage(`not json`) },
		} {
			req := summaryRequest(t.TempDir())
			mutate(&req)
			if _, err := client.Generate(context.Background(), req); !errors.Is(err, ErrInvalidRequest) {
				t.Errorf("err = %v, want ErrInvalidRequest", err)
			}
		}
		if len(fake.calls()) != 0 {
			t.Fatal("invalid requests must not be sent")
		}
	})
}

func TestRetryAfterCapsLargeValues(t *testing.T) {
	t.Parallel()
	for _, header := range []http.Header{
		{"Retry-After": {"1e20"}},
		{"Retry-After-Ms": {"1e20"}},
	} {
		wait, ok := retryAfter(header, time.Now())
		if !ok || wait != maxRetryAfter {
			t.Errorf("retryAfter(%v) = %v, %v; want %v, true", header, wait, ok, maxRetryAfter)
		}
	}
}

func TestGenerateAccountsForMalformedEnvelope(t *testing.T) {
	t.Parallel()
	fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "gen-invalid", "model": req.Model, "choices": "invalid",
			"usage": map[string]any{"cost": 0.25, "prompt_tokens": 500},
		})
	})
	budget := decisions.NewBudget(0.1)
	client := newTestClient(t, fake, budget, nil)
	root := t.TempDir()
	if _, err := client.Generate(t.Context(), summaryRequest(root)); !errors.Is(err, ErrBudget) {
		t.Fatalf("Generate = %v, want budget stop after invalid billed output", err)
	}
	if len(fake.calls()) != 1 || budget.Spent() != 0.25 {
		t.Fatalf("calls=%d, spent=%v; want 1 call and 0.25 spent", len(fake.calls()), budget.Spent())
	}
	rows, err := decisions.LoadReceipts(root)
	if err != nil || len(rows) != 1 || rows[0].Cost == nil || *rows[0].Cost != 0.25 || rows[0].CostUnknown || rows[0].ResponseID != "gen-invalid" {
		t.Fatalf("invalid output lost reported cost: %+v, %v", rows, err)
	}
	if summary := client.Summary(); summary.Calls != 1 || summary.CostUSD != 0.25 || summary.CostUnknown != 0 || summary.Failures != 1 {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestGenerateRecordsCanceledNetworkCall(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	fake := newFakeChat(t, func(_ http.ResponseWriter, _ int, _ chatRequest) { cancel() })
	client := newTestClient(t, fake, decisions.NewBudget(1), nil)
	root := t.TempDir()
	if _, err := client.Generate(ctx, summaryRequest(root)); !errors.Is(err, context.Canceled) {
		t.Fatalf("Generate = %v, want cancellation", err)
	}
	rows, err := decisions.LoadReceipts(root)
	if err != nil || len(rows) != 1 || rows[0].Reason != decisions.ReasonContextCanceled || !rows[0].CostUnknown || rows[0].Attempts != 1 {
		t.Fatalf("canceled request missing from receipts: %+v, %v", rows, err)
	}
	if summary := client.Summary(); summary.Calls != 1 || summary.CostUnknown != 1 || summary.Failures != 1 {
		t.Fatalf("canceled request missing from summary: %+v", summary)
	}
}

func TestShareBudgetAndReceiptsWithEngine(t *testing.T) {
	t.Parallel()

	engine, err := decisions.New(decisions.Options{Config: config.Default().Decisions, APIKey: "k"})
	if err != nil {
		t.Fatal(err)
	}
	fake := newFakeChat(t, func(w http.ResponseWriter, _ int, req chatRequest) { chatReply(w, req.Model, goodOutput, 0, 0.25) })
	client := newTestClient(t, fake, engine.Budget(), engine.Receipts())
	root := t.TempDir()
	if _, err := client.Generate(context.Background(), summaryRequest(root)); err != nil {
		t.Fatal(err)
	}
	if engine.Budget().Spent() != 0.25 {
		t.Fatalf("engine budget spent = %v", engine.Budget().Spent())
	}
	summary := engine.Summary()
	summary.Generation = client.Summary()
	if lines := strings.Join(summary.Lines(), "\n"); !strings.Contains(lines, "generation: 1 call, 0 cache hits, US$ 0.2500") {
		t.Fatalf("lines = %s", lines)
	}
}

func TestValidateLiteralAndDedupe(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value   string
		maxLen  int
		wantErr string
	}{
		{"Model Context Protocol", 40, ""},
		{"Protocolo de Contexto", 21, ""},
		{"", 10, "empty"},
		{"  ", 10, "empty"},
		{" padded", 10, "whitespace"},
		{"two\nlines", 20, "several lines"},
		{"ação ação ação", 5, "more than 5"},
		{"unbounded", 0, ""},
	}
	for _, tc := range testCases {
		err := ValidateLiteral(tc.value, tc.maxLen)
		if tc.wantErr == "" && err != nil {
			t.Errorf("ValidateLiteral(%q) = %v", tc.value, err)
		}
		if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
			t.Errorf("ValidateLiteral(%q) = %v, want %q", tc.value, err, tc.wantErr)
		}
	}
	got := DedupeStrings([]string{" MCP ", "mcp", "", "A2A", "Model Context Protocol", "a2a"})
	if !slices.Equal(got, []string{"MCP", "A2A", "Model Context Protocol"}) {
		t.Fatalf("DedupeStrings = %v", got)
	}
}
