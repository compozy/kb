package contract

import (
	"reflect"
	"strings"
	"testing"
)

func validContract() *Contract {
	return &Contract{
		Purpose:                 "How agents coordinate as a swarm.",
		Core:                    []string{"Swarm runtimes and patterns"},
		Adjacent:                []string{"Multi-agent RL, because methods transfer"},
		CollectedOnPurpose:      []string{"GitHub READMEs for named swarm products"},
		CollectedOnPurposePaths: []string{"raw/brand-audit-*/**"},
		OutOfScope:              []string{"Docker Swarm and other container orchestrators"},
	}
}

func TestContractValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(c *Contract)
		wantErr string
	}{
		{name: "Should accept a complete contract", mutate: func(*Contract) {}},
		{name: "Should accept a contract without optional lists", mutate: func(c *Contract) {
			c.Adjacent, c.CollectedOnPurpose, c.CollectedOnPurposePaths, c.OutOfScope = nil, nil, nil, nil
		}},
		{name: "Should require a purpose", mutate: func(c *Contract) { c.Purpose = "  " }, wantErr: "purpose is required"},
		{name: "Should require a core line", mutate: func(c *Contract) { c.Core = []string{" ", ""} }, wantErr: "core needs at least one line"},
		{
			name:    "Should reject a core line repeated in out_of_scope ignoring case and spaces",
			mutate:  func(c *Contract) { c.OutOfScope = append(c.OutOfScope, "  swarm RUNTIMES and patterns ") },
			wantErr: `"Swarm runtimes and patterns" appears in both core and out_of_scope`,
		},
		{
			name:    "Should reject an adjacent line repeated in out_of_scope",
			mutate:  func(c *Contract) { c.OutOfScope = []string{"multi-agent rl, because methods transfer"} },
			wantErr: "appears in both adjacent and out_of_scope",
		},
		{
			name:    "Should reject a glob that does not compile",
			mutate:  func(c *Contract) { c.CollectedOnPurposePaths = []string{"raw/[abc/**"} },
			wantErr: "collected_on_purpose_paths",
		},
		{
			name:    "Should reject an absolute glob",
			mutate:  func(c *Contract) { c.CollectedOnPurposePaths = []string{"/raw/**"} },
			wantErr: "topic-relative",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := validContract()
			tt.mutate(c)
			err := c.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}

	t.Run("Should reject a nil contract", func(t *testing.T) {
		t.Parallel()
		var c *Contract
		if err := c.Validate(); err == nil {
			t.Fatal("Validate() on nil contract returned nil")
		}
	})
}

func TestContractEmptyAndHash(t *testing.T) {
	t.Parallel()

	t.Run("Should treat nil and whitespace-only contracts as empty with no hash", func(t *testing.T) {
		t.Parallel()
		var nilContract *Contract
		blank := &Contract{Purpose: "  ", Core: []string{" "}}
		for _, c := range []*Contract{nilContract, {}, blank} {
			if !c.Empty() {
				t.Fatalf("Empty() = false for %#v", c)
			}
			if got := c.Hash(); got != "" {
				t.Fatalf("Hash() = %q, want empty", got)
			}
		}
	})

	t.Run("Should hash the trimmed contract deterministically", func(t *testing.T) {
		t.Parallel()
		base := validContract()
		padded := validContract()
		padded.Purpose = "  " + padded.Purpose + "\n"
		padded.Core = append([]string{" "}, " Swarm runtimes and patterns ")
		if base.Empty() {
			t.Fatal("Empty() = true for a complete contract")
		}
		if len(base.Hash()) != 64 {
			t.Fatalf("Hash() = %q, want 64 hex chars", base.Hash())
		}
		if base.Hash() != padded.Hash() {
			t.Fatalf("hash differs for whitespace-only changes: %s vs %s", base.Hash(), padded.Hash())
		}
	})

	t.Run("Should change the hash when any line changes", func(t *testing.T) {
		t.Parallel()
		base := validContract()
		changed := validContract()
		changed.OutOfScope = append(changed.OutOfScope, "Particle-swarm optimization with no LLM agents")
		if base.Hash() == changed.Hash() {
			t.Fatal("hash did not change after adding an out_of_scope line")
		}
		reordered := validContract()
		reordered.Core = []string{"b", "a"}
		ordered := validContract()
		ordered.Core = []string{"a", "b"}
		if reordered.Hash() == ordered.Hash() {
			t.Fatal("hash ignores line order")
		}
	})
}

func TestContractQuestionState(t *testing.T) {
	t.Parallel()

	t.Run("Should expose the five text fields and never the path globs", func(t *testing.T) {
		t.Parallel()
		state := validContract().QuestionState()
		want := map[string]any{
			"purpose":              "How agents coordinate as a swarm.",
			"core":                 []string{"Swarm runtimes and patterns"},
			"adjacent":             []string{"Multi-agent RL, because methods transfer"},
			"collected_on_purpose": []string{"GitHub READMEs for named swarm products"},
			"out_of_scope":         []string{"Docker Swarm and other container orchestrators"},
		}
		if !reflect.DeepEqual(state, want) {
			t.Fatalf("QuestionState() = %#v, want %#v", state, want)
		}
	})

	t.Run("Should return empty lists for a nil contract", func(t *testing.T) {
		t.Parallel()
		var c *Contract
		state := c.QuestionState()
		if got, ok := state["core"].([]string); !ok || got == nil || len(got) != 0 {
			t.Fatalf("core = %#v, want empty non-nil list", state["core"])
		}
	})
}

func TestMatchPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"raw/brand-audit-*/**", "raw/brand-audit-2026-09-06/B016-catapult/firecrawl.md", true},
		{"raw/brand-audit-*/**", "raw/brand-audit-2026-09-06/index.md", true},
		{"raw/brand-audit-*/**", "raw/articles/brand-audit.md", false},
		{"raw/acquisition/**", "raw/acquisition/logs/2026/run.jsonl", true},
		{"raw/acquisition/**", "raw/acquisitions/run.md", false},
		{"raw/*.md", "raw/a.md", true},
		{"raw/*.md", "raw/sub/a.md", false},
		{"raw/**/*.md", "raw/a.md", true},
		{"raw/**/*.md", "raw/x/y/z.md", true},
		{"raw/**/*.md", "raw/x/y/z.pdf", false},
		{"**/README.md", "raw/github/repo/README.md", true},
		{"raw/?.md", "raw/é.md", true},
		{"raw/?.md", "raw/ab.md", false},
		{"raw/[ab].md", "raw/b.md", true},
		{"raw/articles/x.md", "./raw/articles/x.md", true},
		{"raw/articles/x.md", `raw\articles\x.md`, true},
		{"raw/[abc/**", "raw/a/x.md", false},
		{"raw/a**/x.md", "raw/ab/x.md", false},
		{"", "raw/a.md", false},
	}
	for _, tt := range tests {
		t.Run(tt.pattern+" ~ "+tt.path, func(t *testing.T) {
			t.Parallel()
			if got := MatchPath(tt.pattern, tt.path); got != tt.want {
				t.Fatalf("MatchPath(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
			}
		})
	}
}

func TestValidatePattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		pattern string
		wantErr bool
	}{
		{"raw/**", false},
		{"raw/brand-audit-*/**", false},
		{"**/*.md", false},
		{"", true},
		{"/raw/**", true},
		{"raw//x", true},
		{"raw/../x", true},
		{"raw/a**", true},
		{"raw/[x", true},
	}
	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			if err := ValidatePattern(tt.pattern); (err != nil) != tt.wantErr {
				t.Fatalf("ValidatePattern(%q) = %v, wantErr %v", tt.pattern, err, tt.wantErr)
			}
		})
	}
}

func TestContractMatchesCollectedPath(t *testing.T) {
	t.Parallel()

	t.Run("Should match any configured glob and nothing on a nil contract", func(t *testing.T) {
		t.Parallel()
		c := validContract()
		c.CollectedOnPurposePaths = append(c.CollectedOnPurposePaths, " raw/acquisition/** ")
		if !c.MatchesCollectedPath("raw/acquisition/x.md") {
			t.Fatal("expected raw/acquisition/x.md to match")
		}
		if !c.MatchesCollectedPath("raw/brand-audit-1/a.md") {
			t.Fatal("expected raw/brand-audit-1/a.md to match")
		}
		if c.MatchesCollectedPath("raw/articles/a.md") {
			t.Fatal("raw/articles/a.md must not match")
		}
		var nilContract *Contract
		if nilContract.MatchesCollectedPath("raw/acquisition/x.md") {
			t.Fatal("nil contract must not match")
		}
	})
}
