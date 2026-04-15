package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const javaPortfolioPlaybookPath = "../../.compozy/tasks/java-ingest-adapter/_java-portfolio-adoption-playbook.md"

func readJavaPortfolioPlaybook(t *testing.T) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Clean(javaPortfolioPlaybookPath))
	if err != nil {
		t.Fatalf("read java portfolio playbook: %v", err)
	}

	return string(content)
}

func TestJavaPortfolioPlaybookIncludesGovernanceAndContractRequirements(t *testing.T) {
	t.Parallel()

	playbook := readJavaPortfolioPlaybook(t)
	requiredFragments := []string{
		"<= 20%",
		"single-module Java library",
		"Spring-style service repository",
		"multi-module enterprise-style repository",
		">= 80%",
		">= 4/5",
		"`sourceType = \"codebase-file\"`",
		"`codebaseIngestResult`",
		"`GenerationSummary`",
		"`GenerationTimings`",
		"`topic`",
		"`summary`",
		"`timings`",
	}

	for _, fragment := range requiredFragments {
		if !strings.Contains(playbook, fragment) {
			t.Fatalf("playbook must contain fragment %q", fragment)
		}
	}
}

func TestJavaPortfolioPlaybookIncludesFallbackAndUnresolvedGuidance(t *testing.T) {
	t.Parallel()

	playbook := readJavaPortfolioPlaybook(t)
	requiredFragments := []string{
		"`JAVA_RESOLUTION_FALLBACK`",
		"`JAVA_PARSE_ERROR`",
		"`java_fallback_count`",
		"`java_unresolved_count`",
		"High fallback volume",
		"Troubleshooting Matrix",
	}

	for _, fragment := range requiredFragments {
		if !strings.Contains(playbook, fragment) {
			t.Fatalf("playbook must contain fallback guidance fragment %q", fragment)
		}
	}
}

func TestJavaPortfolioPlaybookCommandsAlignWithCurrentCLI(t *testing.T) {
	t.Parallel()

	playbook := readJavaPortfolioPlaybook(t)
	requiredFragments := []string{
		"kb topic new <topic-slug> \"<topic-title>\" <domain> --vault <portfolio-vault>",
		"kb ingest codebase <repo-path> \\",
		"--topic <topic-slug> \\",
		"--vault <portfolio-vault> \\",
		"--progress never \\",
		"--log-format json \\",
		"--dry-run",
		"kb lint <topic-slug> --vault <portfolio-vault> --format json",
	}

	for _, fragment := range requiredFragments {
		if !strings.Contains(playbook, fragment) {
			t.Fatalf("playbook command reference missing fragment %q", fragment)
		}
	}
}
