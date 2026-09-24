package session

import (
	"context"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/topic"
)

func newTopic(t *testing.T) (string, string) {
	t.Helper()
	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	return vault, info.RootPath
}

func testConfig(apiURL string) config.Config {
	cfg := config.Default()
	cfg.OpenRouter.APIKey = "test-key"
	cfg.OpenRouter.APIURL = apiURL
	return cfg
}

func TestOpenRequiresDecisionModel(t *testing.T) {
	t.Parallel()
	vault, _ := newTopic(t)
	cfg := config.Default()
	cfg.OpenRouter.APIKey = ""
	_, err := Open(Options{Config: cfg, VaultPath: vault, Topic: "demo", Command: "kb classify"})
	if !errors.Is(err, ErrDecisionsRequired) {
		t.Fatalf("Open() error = %v, want ErrDecisionsRequired", err)
	}
}

func TestModes(t *testing.T) {
	t.Parallel()
	accepted := &contract.Contract{Purpose: "Demo purpose.", Core: []string{"demo subjects"}}

	tests := []struct {
		name        string
		contract    *contract.Contract
		gates       string
		bodyMode    string
		flag        string
		calibration string
		relevance   string
		wantBody    string
		wantReason  string
		wantRelGate string
		wantQuality string
	}{
		{name: "defaults", wantBody: ModeShadow, wantRelGate: ModeShadow, wantQuality: ModeApply},
		{name: "contract without calibration stays shadow", contract: accepted, wantBody: ModeShadow, wantRelGate: ModeShadow, wantQuality: ModeApply},
		{name: "gates apply", contract: accepted, gates: "apply", wantBody: ModeShadow, wantRelGate: ModeApply, wantQuality: ModeApply},
		{name: "gates apply without contract stays shadow", gates: "apply", wantBody: ModeShadow, wantRelGate: ModeShadow, wantQuality: ModeApply},
		{name: "calibrated", contract: accepted, calibration: `{"purposes":{"relevance":{"dev_labels":24,"holdout_labels":9}}}`, wantBody: ModeShadow, wantRelGate: ModeApply, wantQuality: ModeApply},
		{name: "imported labels alone stay shadow", contract: accepted, calibration: `{"purposes":{"relevance":{"dev_labels":40,"holdout_labels":12,"imported_labels":40}}}`, wantBody: ModeShadow, wantRelGate: ModeShadow, wantQuality: ModeApply},
		{name: "under-calibrated", contract: accepted, calibration: `{"purposes":{"relevance":{"dev_labels":10,"holdout_labels":3}}}`, wantBody: ModeShadow, wantRelGate: ModeShadow, wantQuality: ModeApply},
		{name: "topic body apply", bodyMode: "apply", wantBody: ModeApply, wantRelGate: ModeShadow, wantQuality: ModeApply},
		{name: "flag shadow wins", contract: accepted, gates: "apply", bodyMode: "apply", flag: "shadow", wantBody: ModeShadow, wantRelGate: ModeShadow, wantQuality: ModeShadow},
		{name: "flag apply wins with a contract", contract: accepted, flag: "apply", wantBody: ModeApply, wantRelGate: ModeApply, wantQuality: ModeApply, wantReason: "set by --decisions"},
		{name: "flag apply without contract keeps relevance in shadow", flag: "apply", wantBody: ModeApply, wantRelGate: ModeShadow, wantQuality: ModeApply, wantReason: "--decisions=apply ignored for relevance: no accepted contract"},
		{name: "flag apply with relevance off keeps relevance in shadow", contract: accepted, relevance: "off", flag: "apply", wantBody: ModeApply, wantRelGate: ModeShadow, wantQuality: ModeApply, wantReason: "--decisions=apply ignored: relevance off"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vault, root := newTopic(t)
			if tt.contract != nil {
				if err := contract.SetContract(root, tt.contract); err != nil {
					t.Fatalf("SetContract: %v", err)
				}
			}
			if tt.gates != "" || tt.bodyMode != "" || tt.relevance != "" {
				writeDecisionSettings(t, root, tt.bodyMode, tt.gates, tt.relevance)
			}
			if tt.calibration != "" {
				if err := os.MkdirAll(filepath.Join(root, ".decisions"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(CalibrationPath(root), []byte(tt.calibration), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			s, err := Open(Options{Config: testConfig("http://127.0.0.1:1"), VaultPath: vault, Topic: "demo", Flags: Flags{Decisions: tt.flag}})
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			if got := s.BodyMode(); got != tt.wantBody {
				t.Errorf("BodyMode() = %q, want %q", got, tt.wantBody)
			}
			if got := s.RelevanceGateMode(); got != tt.wantRelGate {
				t.Errorf("RelevanceGateMode() = %q, want %q (%s)", got, tt.wantRelGate, s.GateModeReason())
			}
			if got := s.QualityGateMode(); got != tt.wantQuality {
				t.Errorf("QualityGateMode() = %q, want %q", got, tt.wantQuality)
			}
			if got := s.GateModeReason(); tt.wantReason != "" && !strings.HasPrefix(got, tt.wantReason) {
				t.Errorf("GateModeReason() = %q, want prefix %q", got, tt.wantReason)
			}
		})
	}
}

// TestOpenLoadsTopicExtraBanks: extra banks under .decisions/banks/ are
// loaded at Open and served per purpose; an invalid one fails the run with
// the bank path in the error (spec §4.3). decisions.exclude reaches the
// engine through Ref.
func TestOpenLoadsTopicExtraBanks(t *testing.T) {
	t.Parallel()
	const valid = `{"id":"owner_relevance","version":"1","purpose":"relevance","guard":"G","questions":[{"id":"press_release","type":"noul","instructions":"Is it a press release?","source":"topic owner"}]}`
	const colliding = `{"id":"owner_quality","version":"1","purpose":"quality","guard":"G","questions":[{"id":"thin_or_boilerplate","type":"noul","instructions":"Is it thin?","source":"topic owner"}]}`
	tests := []struct {
		name    string
		bank    string
		wantErr string
	}{
		{name: "valid extra bank", bank: valid},
		{name: "invalid extra bank fails the run", bank: colliding, wantErr: "invalid topic question bank"},
		{name: "malformed extra bank fails the run", bank: `{"id":`, wantErr: "invalid topic question bank"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			vault, root := newTopic(t)
			dir := filepath.Join(root, ".decisions", "banks")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "owner.json"), []byte(tt.bank), 0o644); err != nil {
				t.Fatal(err)
			}
			writeFileAppend(t, filepath.Join(root, "topic.yaml"), "decisions:\n  exclude:\n    - raw/private/**\n")
			s, err := Open(Options{Config: testConfig("http://127.0.0.1:1"), VaultPath: vault, Topic: "demo", Command: "kb classify"})
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) || !strings.Contains(err.Error(), "kb classify") {
					t.Fatalf("Open err = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			if got := s.ExtraBanks(decisions.PurposeRelevance); len(got) != 1 || got[0].ID != "owner_relevance" {
				t.Fatalf("relevance extras = %v", got)
			}
			if got := s.ExtraBanks(decisions.PurposeQuality); len(got) != 0 {
				t.Fatalf("quality extras = %v", got)
			}
			builtin := questions.MustLoad("relevance")
			role := builtin.MustQuestion("role", nil)
			bank, qs, err := s.WithExtras(decisions.PurposeRelevance, builtin, []questions.Q{role})
			if err != nil || bank.ID != "relevance+owner_relevance" || len(qs) != 2 || qs[1].ID != "press_release" {
				t.Fatalf("WithExtras = %v %v %v", bank, qs, err)
			}
			quality := questions.MustLoad("quality")
			if same, qs, err := s.WithExtras(decisions.PurposeQuality, quality, nil); err != nil || same != quality || len(qs) != 0 {
				t.Fatalf("WithExtras without extras must keep the bank: %v %v %v", same, qs, err)
			}
			if !s.Excluded("raw/private/a.md") || !s.Ref.Excluded("raw/private/x/y.md") || s.Excluded("raw/articles/a.md") {
				t.Fatal("decisions.exclude must reach the session ref")
			}
		})
	}
}

func writeFileAppend(t *testing.T, path, text string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, []byte(text)...), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeDecisionSettings(t *testing.T, root, mode, gates, relevance string) {
	t.Helper()
	path := filepath.Join(root, "topic.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	block := "decisions:\n"
	if mode != "" {
		block += "  mode: " + mode + "\n"
	}
	if gates != "" {
		block += "  gates: " + gates + "\n"
	}
	if relevance != "" {
		block += "  relevance: " + relevance + "\n"
	}
	if err := os.WriteFile(path, append(data, []byte(block)...), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestSessionDecidesThroughFakeAndSharesBudget(t *testing.T) {
	t.Parallel()
	fake := fakes.NewOpenRouter(func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "paywall_or_login" {
			return fakes.Noul(0.91)
		}
		return nil
	}, func(string, string, string) any { return map[string]any{"summary": "A short summary."} })
	t.Cleanup(fake.Close)

	vault, _ := newTopic(t)
	s, err := Open(Options{Config: testConfig(fake.URL), VaultPath: vault, Topic: "demo", Flags: Flags{BudgetUSD: 0.5}})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	bank := questions.MustLoad("quality")
	q, err := bank.Question("paywall_or_login", nil)
	if err != nil {
		t.Fatal(err)
	}
	state := map[string]any{"document": map[string]any{"title": "x", "excerpt": "y"}}
	result, err := s.Engine.Decide(context.Background(), s.Request(decisions.PurposeQuality, "raw/articles/x.md", bank, state, []questions.Q{q}))
	if err != nil {
		t.Fatalf("Decide: %v", err)
	}
	answer := result.Answers["paywall_or_login"]
	if p, ok := answer.P(""); !ok || p != 0.91 || answer.Band != decisions.BandApply {
		t.Fatalf("answer = %+v, want decided 0.91 apply", answer)
	}
	again, err := s.Engine.Decide(context.Background(), s.Request(decisions.PurposeQuality, "raw/articles/x.md", bank, state, []questions.Q{q}))
	if err != nil || !again.CacheHit {
		t.Fatalf("second Decide cache hit = %v, err = %v", again.CacheHit, err)
	}
	if got := len(fake.Calls()); got != 1 {
		t.Fatalf("fake calls = %d, want 1", got)
	}
	if s.Engine.Budget().Limit() != 0.5 {
		t.Fatalf("budget limit = %v, want 0.5", s.Engine.Budget().Limit())
	}
	if len(s.SummaryLines()) == 0 {
		t.Fatal("empty summary")
	}
}

// TestOpenBudgetFlag: an explicit --budget 0 is honored as a cache-only run
// (no request reaches the endpoint) instead of falling back to
// [decisions].budget_usd; an unset flag uses the configured ceiling; a
// negative, NaN or infinite amount is rejected before the session opens.
func TestOpenBudgetFlag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		flags     Flags
		wantLimit float64
		wantErr   bool
	}{
		{name: "Should use the configured ceiling when --budget is not given", flags: Flags{}, wantLimit: config.Default().Decisions.BudgetUSD},
		{name: "Should honor an explicit zero budget as cache-only", flags: Flags{BudgetUSD: 0, BudgetSet: true}, wantLimit: 0},
		{name: "Should honor an explicit positive budget", flags: Flags{BudgetUSD: 0.25, BudgetSet: true}, wantLimit: 0.25},
		{name: "Should reject a negative budget", flags: Flags{BudgetUSD: -1, BudgetSet: true}, wantErr: true},
		{name: "Should reject a NaN budget", flags: Flags{BudgetUSD: math.NaN(), BudgetSet: true}, wantErr: true},
		{name: "Should reject an infinite budget", flags: Flags{BudgetUSD: math.Inf(1), BudgetSet: true}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fake := fakes.NewOpenRouter(nil, nil)
			t.Cleanup(fake.Close)
			vault, _ := newTopic(t)
			s, err := Open(Options{Config: testConfig(fake.URL), VaultPath: vault, Topic: "demo", Command: "kb classify", Flags: tt.flags})
			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), "--budget") {
					t.Fatalf("Open() error = %v, want a --budget error", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			if got := s.Engine.Budget().Limit(); got != tt.wantLimit {
				t.Fatalf("budget limit = %v, want %v", got, tt.wantLimit)
			}
			if tt.wantLimit != 0 {
				return
			}
			bank := questions.MustLoad("quality")
			q, err := bank.Question("paywall_or_login", nil)
			if err != nil {
				t.Fatal(err)
			}
			state := map[string]any{"document": map[string]any{"title": "x", "excerpt": "y"}}
			result, err := s.Engine.Decide(context.Background(), s.Request(decisions.PurposeQuality, "raw/articles/x.md", bank, state, []questions.Q{q}))
			if err != nil {
				t.Fatalf("Decide: %v", err)
			}
			if answer := result.Answers["paywall_or_login"]; answer.Decided() || answer.Reason != decisions.ReasonBudget {
				t.Fatalf("answer = %+v, want undecided:budget", answer)
			}
			if calls := len(fake.Calls()); calls != 0 {
				t.Fatalf("a zero budget must make no request, got %d", calls)
			}
		})
	}
}
