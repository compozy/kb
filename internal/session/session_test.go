package session

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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
		wantBody    string
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
		{name: "flag apply wins", flag: "apply", wantBody: ModeApply, wantRelGate: ModeApply, wantQuality: ModeApply},
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
			if tt.gates != "" || tt.bodyMode != "" {
				writeDecisionSettings(t, root, tt.bodyMode, tt.gates)
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
		})
	}
}

func writeDecisionSettings(t *testing.T, root, mode, gates string) {
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
