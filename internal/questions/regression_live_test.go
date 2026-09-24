//go:build jevlive

// The regression cases run against the real decisions route only with
// `go test -tags jevlive ./internal/questions/` and OPENROUTER_API_KEY set.
// Receipts are reused by hash across runs through a persistent cache topic
// root: $KB_JEVLIVE_CACHE when set, else <os.UserCacheDir()>/kb/jevlive/<bank>.
// Delete that directory to force fresh calls (for example after a model or
// bank change, which changes the cache key anyway). Spend is capped at
// US$ 0.50 per run.
package questions_test

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/questions"
)

func TestRegressionCasesLive(t *testing.T) {
	cfg := config.Default()
	config.ApplyEnvOverrides(&cfg)
	if cfg.OpenRouter.APIKey == "" {
		t.Skip("OPENROUTER_API_KEY is not set")
	}
	cacheRoot := os.Getenv("KB_JEVLIVE_CACHE")
	if cacheRoot == "" {
		base, err := os.UserCacheDir()
		if err != nil {
			t.Fatalf("resolve cache dir: %v", err)
		}
		cacheRoot = filepath.Join(base, "kb", "jevlive")
	}
	engine, err := decisions.New(decisions.Options{
		Config: cfg.Decisions,
		APIKey: cfg.OpenRouter.APIKey,
		APIURL: cfg.OpenRouter.APIURL,
		Budget: decisions.NewBudget(0.5),
	})
	if err != nil {
		t.Fatalf("decisions.New: %v", err)
	}
	files, err := questions.LoadRegressionFiles(filepath.Join("testdata", "regression"))
	if err != nil {
		t.Fatal(err)
	}
	for _, bankID := range questions.BankIDs() {
		file := files[bankID]
		for _, c := range file.Cases {
			t.Run(bankID+"/"+c.Name, func(t *testing.T) {
				bank, q, state, err := file.Build(c)
				if err != nil {
					t.Fatal(err)
				}
				result, err := engine.Decide(context.Background(), decisions.Request{
					Topic:     decisions.TopicRef{Slug: "jevlive-" + bankID, Root: filepath.Join(cacheRoot, bankID)},
					Purpose:   decisions.Purpose(bank.Purpose),
					Subject:   c.Name,
					Bank:      bank,
					State:     state,
					Questions: []questions.Q{q},
				})
				if err != nil {
					t.Fatalf("Decide: %v", err)
				}
				answer := result.Answers[q.ID]
				if !answer.Decided() {
					t.Fatalf("answer %s:%s", answer.Status, answer.Reason)
				}
				got := observed(answer)
				if !slices.Contains(c.Allowed, got) {
					t.Errorf("answer %q (probs %v, noul %v), allowed %v", got, answer.Probs, answer.Noul, c.Allowed)
				}
			})
		}
	}
	for _, line := range engine.Summary().Lines() {
		t.Log(line)
	}
}

func observed(answer decisions.Answer) string {
	switch {
	case answer.Noul != nil:
		if *answer.Noul >= 0.5 {
			return "yes"
		}
		return "no"
	case answer.Score != nil:
		return strconv.Itoa(int(math.Round(*answer.Score)))
	default:
		return answer.Choice
	}
}
