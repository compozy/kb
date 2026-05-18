package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestHomebrewFormulaReleaseContract(t *testing.T) {
	t.Parallel()

	root := repoRoot(t)
	configText := readRepoFile(t, root, ".goreleaser.yml")

	var cfg map[string]any
	if err := yaml.Unmarshal([]byte(configText), &cfg); err != nil {
		t.Fatalf("unmarshal .goreleaser.yml: %v", err)
	}

	if _, ok := cfg["homebrew_casks"]; ok {
		t.Fatal("expected GoReleaser to publish Homebrew formulas, not casks")
	}

	brews := sliceAt(t, cfg, "brews")
	if len(brews) != 1 {
		t.Fatalf("brews len = %d, want 1", len(brews))
	}
	formula := asMap(t, brews[0], "brews[0]")
	if got, want := stringAt(t, formula, "name"), "kb"; got != want {
		t.Fatalf("brews[0].name = %q, want %q", got, want)
	}
	if !stringSliceContains(sliceAt(t, formula, "ids"), "kb-archive") {
		t.Fatalf("brews[0].ids = %#v, want kb-archive", formula["ids"])
	}
	if got, want := stringAt(t, formula, "directory"), "Formula"; got != want {
		t.Fatalf("brews[0].directory = %q, want %q", got, want)
	}

	repository := asMap(t, formula["repository"], "brews[0].repository")
	if got, want := stringAt(t, repository, "owner"), "compozy"; got != want {
		t.Fatalf("brews[0].repository.owner = %q, want %q", got, want)
	}
	if got, want := stringAt(t, repository, "name"), "homebrew-kb"; got != want {
		t.Fatalf("brews[0].repository.name = %q, want %q", got, want)
	}

	assertArchiveNotWrapped(t, cfg, "kb-archive")
}

func TestHomebrewInstallDocsUseFormulaTap(t *testing.T) {
	t.Parallel()

	const installCommand = "brew install compozy/kb/kb"
	forbidden := []string{
		"brew install --cask",
		"homebrew_casks:",
		"directory: Casks",
	}

	root := repoRoot(t)
	paths := []string{
		"README.md",
		".goreleaser.yml",
	}

	for _, path := range paths {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			content := readRepoFile(t, root, path)
			if !strings.Contains(content, installCommand) {
				t.Fatalf("expected %s to contain %q", path, installCommand)
			}
			for _, retired := range forbidden {
				if strings.Contains(content, retired) {
					t.Fatalf("expected %s to avoid retired Homebrew contract %q", path, retired)
				}
			}
		})
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".goreleaser.yml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repo root with .goreleaser.yml not found")
		}
		dir = parent
	}
}

func readRepoFile(t *testing.T, root string, path ...string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(append([]string{root}, path...)...))
	if err != nil {
		t.Fatalf("read %s: %v", filepath.Join(path...), err)
	}
	return string(content)
}

func sliceAt(t *testing.T, src map[string]any, key string) []any {
	t.Helper()

	value, ok := src[key]
	if !ok {
		t.Fatalf("%s missing", key)
	}
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("%s type = %T, want []any", key, value)
	}
	return items
}

func asMap(t *testing.T, value any, label string) map[string]any {
	t.Helper()

	item, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s type = %T, want map[string]any", label, value)
	}
	return item
}

func stringAt(t *testing.T, src map[string]any, key string) string {
	t.Helper()

	value, ok := src[key]
	if !ok {
		t.Fatalf("%s missing", key)
	}
	text, ok := value.(string)
	if !ok {
		t.Fatalf("%s type = %T, want string", key, value)
	}
	return text
}

func stringSliceContains(values []any, want string) bool {
	for _, value := range values {
		if text, ok := value.(string); ok && text == want {
			return true
		}
	}
	return false
}

func assertArchiveNotWrapped(t *testing.T, cfg map[string]any, archiveID string) {
	t.Helper()

	for _, item := range sliceAt(t, cfg, "archives") {
		archive := asMap(t, item, "archives[]")
		if stringAt(t, archive, "id") != archiveID {
			continue
		}
		if wrapped, ok := archive["wrap_in_directory"].(bool); ok && wrapped {
			t.Fatalf("archive %q must keep the binary at the archive root for Homebrew formula installs", archiveID)
		}
		return
	}
	t.Fatalf("archive %q not found", archiveID)
}
