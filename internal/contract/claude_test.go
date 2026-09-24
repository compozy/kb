package contract

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func readTestdata(t *testing.T, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read testdata %s: %v", name, err)
	}
	return string(content)
}

func TestParseClaudeSectionVaultSnippets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		file           string
		purposePrefix  string
		core           int
		adjacent       int
		collected      int
		outOfScope     []string
		lastCoreSuffix string
	}{
		{
			file:          "agent-swarm.md",
			purposePrefix: "How many LLM agents coordinate as a swarm",
			core:          1,
			adjacent:      1,
			collected:     1,
			outOfScope: []string{
				"Docker Swarm and other container orchestrators",
				"numerical particle-swarm optimization with no LLM agents",
				"biological swarms with no LLM agents",
				"step-by-step coding lessons that show one agent and never a swarm runtime",
				"no-code recipes that never hand work between two or more agents.",
			},
		},
		{
			file:           "wearable-hardware.md",
			purposePrefix:  "Determine which body-worn",
			core:           8,
			adjacent:       4,
			collected:      1,
			lastCoreSuffix: "camera or event-camera ball spin and trajectory.",
			outOfScope: []string{
				"Clinical, biological or behavioral studies whose contribution is an endpoint, treatment or cohort result with no sensing, materials, power, radio or fabrication method",
				"vehicle, building or robot navigation that uses no IMU, event camera or radio-localization method",
				"sports-product market surveys with no hardware detail.",
			},
		},
		{
			file:          "polymarket-quant.md",
			purposePrefix: "How to price, quote, execute",
			core:          1,
			// The parenthesised "(Kelly, drawdown control, ...)" stays inside one item.
			adjacent:  1,
			collected: 1,
			outOfScope: []string{
				"Cross-chain routing and signature schemes",
				"personality or time-travel papers",
				"profit listicles",
				"factor investing, ESG or sustainable investing, fund-performance league tables, and stock-picking or sector-rotation strategies with no carry-over rule.",
			},
		},
		{file: "branding.md", purposePrefix: "Gather how brands are created", core: 5, adjacent: 1, collected: 1},
		{file: "sports-tech.md", purposePrefix: "Decide which US basketball", core: 1, adjacent: 1, collected: 1},
		{file: "yc-companies.md", purposePrefix: "Keep a local mirror", core: 1, adjacent: 2, collected: 1},
	}
	for _, tt := range tests {
		t.Run("Should parse "+tt.file, func(t *testing.T) {
			t.Parallel()
			parsed, err := ParseClaudeSection(readTestdata(t, tt.file))
			if err != nil {
				t.Fatalf("ParseClaudeSection returned error: %v", err)
			}
			if !strings.HasPrefix(parsed.Purpose, tt.purposePrefix) {
				t.Fatalf("purpose = %q, want prefix %q", parsed.Purpose, tt.purposePrefix)
			}
			if len(parsed.Core) != tt.core || len(parsed.Adjacent) != tt.adjacent || len(parsed.CollectedOnPurpose) != tt.collected {
				t.Fatalf("counts core=%d adjacent=%d collected=%d, want %d/%d/%d\n%#v",
					len(parsed.Core), len(parsed.Adjacent), len(parsed.CollectedOnPurpose), tt.core, tt.adjacent, tt.collected, parsed)
			}
			if tt.outOfScope != nil && !reflect.DeepEqual(parsed.OutOfScope, tt.outOfScope) {
				t.Fatalf("out_of_scope = %#v, want %#v", parsed.OutOfScope, tt.outOfScope)
			}
			if len(parsed.OutOfScope) == 0 {
				t.Fatal("out_of_scope is empty")
			}
			if tt.lastCoreSuffix != "" && parsed.Core[len(parsed.Core)-1] != tt.lastCoreSuffix {
				t.Fatalf("last core item = %q, want %q", parsed.Core[len(parsed.Core)-1], tt.lastCoreSuffix)
			}
			if err := parsed.Validate(); err != nil {
				t.Fatalf("imported contract does not validate: %v", err)
			}
		})
	}
}

func TestParseClaudeSectionVariants(t *testing.T) {
	t.Parallel()

	t.Run("Should return ErrNoSection when the section is missing", func(t *testing.T) {
		t.Parallel()
		_, err := ParseClaudeSection("# Topic\n\n## Audit log\n\n- **Purpose:** not here\n")
		if !errors.Is(err, ErrNoSection) {
			t.Fatalf("err = %v, want ErrNoSection", err)
		}
	})

	t.Run("Should accept Central, plain Adjacent, sub-bullets and wrapped lines", func(t *testing.T) {
		t.Parallel()
		markdown := strings.Join([]string{
			"# Topic",
			"",
			"## Selection contract",
			"",
			"- **Purpose**: Track things",
			"  that wrap.",
			"- **Central:**",
			"  - first (a; b)",
			"  - second",
			"- **Adjacent:** kept one; kept two",
			"- **Collected on purpose (paths):** `raw/audit/**`; `raw/acquisition/**`",
			"- **Out of scope:** none",
			"- **Unknown label:** ignored",
			"  - ignored child",
			"",
			"Trailing paragraph that is not a bullet.",
			"",
			"## Next",
			"",
			"- **Core:** outside the section",
			"",
		}, "\n")
		parsed, err := ParseClaudeSection(markdown)
		if err != nil {
			t.Fatalf("ParseClaudeSection returned error: %v", err)
		}
		want := Contract{
			Purpose:                 "Track things that wrap.",
			Core:                    []string{"first (a; b)", "second"},
			Adjacent:                []string{"kept one", "kept two"},
			CollectedOnPurpose:      []string{},
			CollectedOnPurposePaths: []string{"raw/audit/**", "raw/acquisition/**"},
			OutOfScope:              []string{},
		}
		if !reflect.DeepEqual(*parsed, want) {
			t.Fatalf("parsed = %#v, want %#v", *parsed, want)
		}
	})

	t.Run("Should parse the placeholder section as an empty contract", func(t *testing.T) {
		t.Parallel()
		parsed, err := ParseClaudeSection("# T\n\n" + RenderSection(nil))
		if err != nil {
			t.Fatalf("ParseClaudeSection returned error: %v", err)
		}
		if !parsed.Empty() {
			t.Fatalf("placeholder parsed as %#v, want empty", parsed)
		}
	})
}

func TestRenderSection(t *testing.T) {
	t.Parallel()

	t.Run("Should render every bullet with the generated-section note", func(t *testing.T) {
		t.Parallel()
		c := validContract()
		c.OutOfScope = []string{"a", "b (x; y)"}
		got := RenderSection(c)
		want := strings.Join([]string{
			"## Selection contract",
			"",
			sectionNote,
			"",
			"- **Purpose:** How agents coordinate as a swarm.",
			"- **Core:** Swarm runtimes and patterns",
			"- **Adjacent (keep):** Multi-agent RL, because methods transfer",
			"- **Collected on purpose:** GitHub READMEs for named swarm products",
			"- **Collected on purpose (paths):** `raw/brand-audit-*/**`",
			"- **Out of scope:** a; b (x; y)",
			"",
		}, "\n")
		if got != want {
			t.Fatalf("RenderSection() =\n%s\nwant\n%s", got, want)
		}
	})

	t.Run("Should render placeholders pointing at kb topic contract for an empty contract", func(t *testing.T) {
		t.Parallel()
		got := RenderSection(&Contract{})
		for _, fragment := range []string{"## Selection contract", "kb topic contract <topic> --draft", "--import-claude", "- **Out of scope:** _not set_"} {
			if !strings.Contains(got, fragment) {
				t.Fatalf("empty render missing %q:\n%s", fragment, got)
			}
		}
		if got != RenderSection(nil) {
			t.Fatal("nil and empty contracts render differently")
		}
	})

	t.Run("Should round-trip through ParseClaudeSection including items with semicolons", func(t *testing.T) {
		t.Parallel()
		c := validContract()
		c.Core = []string{"alpha; with a top-level semicolon", "beta"}
		c.CollectedOnPurposePaths = []string{"raw/a/**", "raw/b/*.md"}
		parsed, err := ParseClaudeSection(RenderIntoClaude("# T\n", c))
		if err != nil {
			t.Fatalf("ParseClaudeSection returned error: %v", err)
		}
		if want := c.Normalized(); !reflect.DeepEqual(*parsed, want) {
			t.Fatalf("round trip = %#v, want %#v", *parsed, want)
		}
		if parsed.Hash() != c.Hash() {
			t.Fatal("round trip changed the contract hash")
		}
	})

	t.Run("Should round-trip every vault snippet", func(t *testing.T) {
		t.Parallel()
		for _, file := range []string{"agent-swarm.md", "branding.md", "sports-tech.md", "wearable-hardware.md", "polymarket-quant.md", "yc-companies.md"} {
			imported, err := ParseClaudeSection(readTestdata(t, file))
			if err != nil {
				t.Fatalf("%s: %v", file, err)
			}
			again, err := ParseClaudeSection(RenderIntoClaude(readTestdata(t, file), imported))
			if err != nil {
				t.Fatalf("%s: reparse: %v", file, err)
			}
			if again.Hash() != imported.Hash() {
				t.Fatalf("%s: re-rendered contract differs:\n%#v\n%#v", file, imported, again)
			}
		}
	})
}

func TestRenderIntoClaude(t *testing.T) {
	t.Parallel()

	c := validContract()
	section := RenderSection(c)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Should insert before the first H2 heading",
			input: "# Title\n\n**Topic scope:** x.\n\n## Audit log\n\nbody\n",
			want:  "# Title\n\n**Topic scope:** x.\n\n" + section + "\n## Audit log\n\nbody\n",
		},
		{
			name:  "Should add a blank line when the intro does not end with one",
			input: "# Title\nintro\n## Audit log\n",
			want:  "# Title\nintro\n\n" + section + "\n## Audit log\n",
		},
		{
			name:  "Should append at the end when there is no H2",
			input: "# Title\n\nintro",
			want:  "# Title\n\nintro\n\n" + section,
		},
		{
			name:  "Should replace an existing section in the middle",
			input: "# Title\n\n## Selection contract\n\n- **Purpose:** old\n\n## Audit log\n\nkept\n",
			want:  "# Title\n\n" + section + "\n## Audit log\n\nkept\n",
		},
		{
			name:  "Should replace an existing section at the end",
			input: "# Title\n\n## Research gaps\n\n- Gap\n\n## Selection contract\n\n- **Purpose:** old\n- **Core:** old\n",
			want:  "# Title\n\n## Research gaps\n\n- Gap\n\n" + section,
		},
		{
			name:  "Should ignore headings inside fenced code blocks",
			input: "# Title\n\n```md\n## Not a heading\n```\n\n## Audit log\n",
			want:  "# Title\n\n```md\n## Not a heading\n```\n\n" + section + "\n## Audit log\n",
		},
		{
			name:  "Should preserve CRLF bytes outside the section",
			input: "# Title\r\n\r\n## Selection contract\r\n\r\n- **Purpose:** old\r\n\r\n## Audit log\r\nkept\r\n",
			want:  "# Title\r\n\r\n" + section + "\n## Audit log\r\nkept\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := RenderIntoClaude(tt.input, c)
			if got != tt.want {
				t.Fatalf("RenderIntoClaude() =\n%q\nwant\n%q", got, tt.want)
			}
			if again := RenderIntoClaude(got, c); again != got {
				t.Fatalf("RenderIntoClaude is not idempotent:\n%q\n%q", got, again)
			}
		})
	}
}

func TestSyncClaude(t *testing.T) {
	t.Parallel()

	t.Run("Should re-render the section from the active contract and keep the rest", func(t *testing.T) {
		t.Parallel()
		topicRoot := t.TempDir()
		claudePath := filepath.Join(topicRoot, "CLAUDE.md")
		original := "# Topic\n\nIntro.\n\n## Selection contract\n\n- **Purpose:** old\n\n## Audit log\n\nKeep me.\n"
		writeTestFile(t, claudePath, original)
		if err := SetContract(topicRoot, validContract()); err != nil {
			t.Fatalf("SetContract returned error: %v", err)
		}

		if err := SyncClaude(topicRoot); err != nil {
			t.Fatalf("SyncClaude returned error: %v", err)
		}
		got := readTestFile(t, claudePath)
		want := "# Topic\n\nIntro.\n\n" + RenderSection(validContract()) + "\n## Audit log\n\nKeep me.\n"
		if got != want {
			t.Fatalf("CLAUDE.md =\n%s\nwant\n%s", got, want)
		}
	})

	t.Run("Should render placeholders when no contract is accepted", func(t *testing.T) {
		t.Parallel()
		topicRoot := t.TempDir()
		claudePath := filepath.Join(topicRoot, "CLAUDE.md")
		writeTestFile(t, claudePath, "# Topic\n\n## Audit log\n")
		if err := SyncClaude(topicRoot); err != nil {
			t.Fatalf("SyncClaude returned error: %v", err)
		}
		if got := readTestFile(t, claudePath); !strings.Contains(got, "_not set") {
			t.Fatalf("CLAUDE.md missing placeholder section:\n%s", got)
		}
	})

	t.Run("Should fail when CLAUDE.md is missing", func(t *testing.T) {
		t.Parallel()
		if err := SyncClaude(t.TempDir()); err == nil {
			t.Fatal("SyncClaude without CLAUDE.md returned nil")
		}
	})
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}
