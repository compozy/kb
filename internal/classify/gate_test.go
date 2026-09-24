package classify

import (
	"context"
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/quality"
)

func questionIDs(call fakes.Call) []string {
	ids := make([]string, 0, len(call.Questions))
	for id := range call.Questions {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids
}

func TestJudgeGate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setup         func(t *testing.T, root string)
		path          string
		transcript    bool
		wantQuestions []string
		wantRole      string
		wantByPath    bool
		wantPOff      float64
		wantPKept     float64
		wantContract  bool
	}{
		{
			name:          "role and quality in one request",
			setup:         func(t *testing.T, root string) { setTestContract(t, root, nil) },
			path:          "raw/articles/a.md",
			wantQuestions: []string{NoulErrorPage, NoulPaywall, "role", NoulThin},
			wantRole:      RoleOffTopic,
			wantPOff:      0.7,
			wantPKept:     0.1 + 0.05 + 0.05,
			wantContract:  true,
		},
		{
			name:          "transcript adds no_speech_content",
			setup:         func(t *testing.T, root string) { setTestContract(t, root, nil) },
			path:          "raw/youtube/a.md",
			transcript:    true,
			wantQuestions: []string{NoulErrorPage, NoulNoSpeech, NoulPaywall, "role", NoulThin},
			wantRole:      RoleOffTopic,
			wantPOff:      0.7,
			wantPKept:     0.2,
			wantContract:  true,
		},
		{
			name:          "collected path omits role",
			setup:         func(t *testing.T, root string) { setTestContract(t, root, nil) },
			path:          "raw/audit/b016/page.md",
			wantQuestions: []string{NoulErrorPage, NoulPaywall, NoulThin},
			wantRole:      RoleCollectedOnPurpose,
			wantByPath:    true,
			wantPKept:     1,
		},
		{
			name: "relevance off omits role",
			setup: func(t *testing.T, root string) {
				setTestContract(t, root, nil)
				appendTopicYAML(t, root, "decisions:\n  relevance: off\n")
			},
			path:          "raw/audit/b016/page.md",
			wantQuestions: []string{NoulErrorPage, NoulPaywall, NoulThin},
		},
		{
			name:          "no contract still asks role against an empty scope",
			path:          "raw/articles/a.md",
			wantQuestions: []string{NoulErrorPage, NoulPaywall, "role", NoulThin},
			wantRole:      RoleOffTopic,
			wantPOff:      0.7,
			wantPKept:     0.2,
			wantContract:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fake := fakeServer(t, func(_ fakes.Call, q fakes.Question) any {
				switch q.ID {
				case "role":
					return fakes.Dist(map[string]float64{"core": 0.1, "adjacent": 0.05, "collected_on_purpose": 0.05, "general": 0.05, "off_topic": 0.7, "unknown": 0.05})
				case NoulErrorPage:
					return fakes.Noul(0.9)
				case NoulPaywall:
					return fakes.Noul(0.6)
				}
				return nil
			}, nil)
			vault, root := newTestTopic(t)
			if tt.setup != nil {
				tt.setup(t, root)
			}
			source(t, root, tt.path, "A page", "https://example.com/post", longBody("typed decisions", 100))
			s := openTestSession(t, vault, fake.URL)
			doc, err := corpus.ReadDocument(root, tt.path, corpus.KindSource)
			if err != nil {
				t.Fatal(err)
			}

			judgment, err := JudgeGate(context.Background(), s, doc, GateOptions{IsTranscript: tt.transcript})
			if err != nil {
				t.Fatalf("JudgeGate: %v", err)
			}
			calls := fake.Calls()
			if len(calls) != 1 {
				t.Fatalf("calls = %d, want exactly one request", len(calls))
			}
			if got := questionIDs(calls[0]); !slices.Equal(got, tt.wantQuestions) {
				t.Fatalf("questions = %v, want %v", got, tt.wantQuestions)
			}
			if got := strings.Contains(calls[0].StateString(), `"contract"`); got != tt.wantContract {
				t.Fatalf("state has contract = %v, want %v: %s", got, tt.wantContract, calls[0].StateString())
			}
			if judgment.Role != tt.wantRole || judgment.RoleByPath != tt.wantByPath {
				t.Fatalf("role = %q (by path %v), want %q (%v)", judgment.Role, judgment.RoleByPath, tt.wantRole, tt.wantByPath)
			}
			if math.Abs(judgment.POffTopic-tt.wantPOff) > 1e-9 || math.Abs(judgment.PKept-tt.wantPKept) > 1e-9 {
				t.Fatalf("P(off_topic) = %v, P(kept) = %v, want %v, %v", judgment.POffTopic, judgment.PKept, tt.wantPOff, tt.wantPKept)
			}
			if !judgment.QualityDecided || judgment.Quality[NoulErrorPage] != 0.9 || judgment.Quality[NoulPaywall] != 0.6 {
				t.Fatalf("quality = %v (decided %v)", judgment.Quality, judgment.QualityDecided)
			}
			if judgment.ReceiptKey == "" {
				t.Fatal("missing receipt key")
			}
			if reason, p := judgment.QualityReason(0.8); reason != ReasonErrorPage || p != 0.9 {
				t.Fatalf("QualityReason(0.8) = %q %v", reason, p)
			}

			// The same content is answered from the receipt cache.
			again, err := JudgeGate(context.Background(), s, doc, GateOptions{IsTranscript: tt.transcript})
			if err != nil || again.ReceiptKey != judgment.ReceiptKey || len(fake.Calls()) != 1 {
				t.Fatalf("second JudgeGate: err %v, key %q vs %q, calls %d", err, again.ReceiptKey, judgment.ReceiptKey, len(fake.Calls()))
			}
		})
	}
}

func TestJudgeGateUndecidedNeverDecides(t *testing.T) {
	t.Parallel()
	fake := fakeServer(t, func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "role" {
			// A malformed answer (distribution does not sum to 1) is an
			// invalid receipt, never coerced.
			return map[string]any{"type": "choice", "choice": "off_topic", "probabilities": map[string]float64{"off_topic": 0.4}, "confidence": 0.1}
		}
		return nil
	}, nil)
	vault, root := newTestTopic(t)
	setTestContract(t, root, nil)
	source(t, root, "raw/articles/a.md", "A page", "", longBody("x", 50))
	s := openTestSession(t, vault, fake.URL)
	doc, err := corpus.ReadDocument(root, "raw/articles/a.md", corpus.KindSource)
	if err != nil {
		t.Fatal(err)
	}
	judgment, err := JudgeGate(context.Background(), s, doc, GateOptions{})
	if err != nil {
		t.Fatalf("JudgeGate: %v", err)
	}
	if judgment.RoleDecided || judgment.Role != "" || judgment.POffTopic != 0 {
		t.Fatalf("undecided role must stay empty: %+v", judgment)
	}
	if len(judgment.Undecided) == 0 || !strings.HasPrefix(judgment.Undecided[0], "role:") {
		t.Fatalf("Undecided = %v", judgment.Undecided)
	}
}

func TestQualityReason(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		judgment   GateJudgment
		threshold  float64
		wantReason string
		wantP      float64
	}{
		{name: "nothing fired", judgment: GateJudgment{Quality: map[string]float64{NoulPaywall: 0.2}}, threshold: 0.5},
		{name: "highest noul wins", judgment: GateJudgment{Quality: map[string]float64{NoulPaywall: 0.85, NoulThin: 0.95}}, threshold: 0.8, wantReason: ReasonThin, wantP: 0.95},
		{name: "review band", judgment: GateJudgment{Quality: map[string]float64{NoulPaywall: 0.6}}, threshold: 0.5, wantReason: ReasonPaywall, wantP: 0.6},
		{name: "no speech", judgment: GateJudgment{Quality: map[string]float64{NoulNoSpeech: 0.9}}, threshold: 0.8, wantReason: ReasonNoSpeech, wantP: 0.9},
		{name: "code flag wins", judgment: GateJudgment{Flags: []quality.Flag{{Code: quality.NotAnArticle}}, Quality: map[string]float64{NoulThin: 0.99}}, threshold: 0.8, wantReason: ReasonNotAnArticle, wantP: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reason, p := tt.judgment.QualityReason(tt.threshold)
			if reason != tt.wantReason || p != tt.wantP {
				t.Fatalf("QualityReason = %q %v, want %q %v", reason, p, tt.wantReason, tt.wantP)
			}
		})
	}
	g := GateJudgment{Flags: []quality.Flag{{Code: quality.Thin}}}
	if got := g.QualityQuestion(quality.Thin); got != "code:thin" {
		t.Fatalf("QualityQuestion = %q", got)
	}
	if got := (GateJudgment{}).QualityQuestion(ReasonPaywall); got != NoulPaywall {
		t.Fatalf("QualityQuestion = %q", got)
	}
}

func TestIsTranscriptKind(t *testing.T) {
	t.Parallel()
	for kind, want := range map[models.SourceKind]bool{
		models.SourceKindYouTubeTranscript: true,
		models.SourceKindInstagramVideo:    true,
		models.SourceKindArticle:           false,
		models.SourceKindDocument:          false,
	} {
		if got := IsTranscriptKind(string(kind)); got != want {
			t.Errorf("IsTranscriptKind(%q) = %v", kind, got)
		}
	}
}
