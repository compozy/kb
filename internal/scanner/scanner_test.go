package scanner

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestScanWorkspaceRoutesSupportedFilesByLanguage(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	outputPath := filepath.Join(rootPath, "generated", "vault")

	writeTestFile(t, rootPath, "src/index.ts", "export const index = true;\n")
	writeTestFile(t, rootPath, "src/component.tsx", "export const Component = () => null;\n")
	writeTestFile(t, rootPath, "src/script.js", "export const script = true;\n")
	writeTestFile(t, rootPath, "src/view.jsx", "export const View = () => null;\n")
	writeTestFile(t, rootPath, "go/main.go", "package main\n")
	writeTestFile(t, rootPath, "src/types.d.ts", "export type Value = string;\n")
	writeTestFile(t, rootPath, "README.md", "# ignored\n")
	writeTestFile(t, rootPath, filepath.Join("generated", "vault", "index.ts"), "export const ignored = true;\n")

	workspace := scanTestWorkspace(t, rootPath, WithOutputPath(outputPath))

	expectedPaths := []string{
		"go/main.go",
		"src/component.tsx",
		"src/index.ts",
		"src/script.js",
		"src/view.jsx",
	}

	if got := scannedPaths(workspace.Files); !reflect.DeepEqual(got, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, got)
	}

	expectedGroups := map[string]int{
		"go":  1,
		"js":  1,
		"jsx": 1,
		"ts":  1,
		"tsx": 1,
	}

	if got := groupedCounts(workspace); !reflect.DeepEqual(got, expectedGroups) {
		t.Fatalf("expected grouped counts %v, got %v", expectedGroups, got)
	}
}

func TestScanWorkspaceExcludesNodeModulesByDefault(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, "src/index.ts", "export const index = true;\n")
	writeTestFile(t, rootPath, "node_modules/library/index.ts", "export const ignored = true;\n")

	workspace := scanTestWorkspace(t, rootPath)

	expectedPaths := []string{"src/index.ts"}
	if got := scannedPaths(workspace.Files); !reflect.DeepEqual(got, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, got)
	}
}

func TestScanWorkspaceExcludesGitDirectoryByDefault(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, "src/index.ts", "export const index = true;\n")
	writeTestFile(t, rootPath, ".git/objects/ignored.ts", "export const ignored = true;\n")

	workspace := scanTestWorkspace(t, rootPath)

	expectedPaths := []string{"src/index.ts"}
	if got := scannedPaths(workspace.Files); !reflect.DeepEqual(got, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, got)
	}
}

func TestScanWorkspaceRespectsGitIgnorePatterns(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, ".gitignore", "src/ignored.ts\n")
	writeTestFile(t, rootPath, "src/index.ts", "export const keep = true;\n")
	writeTestFile(t, rootPath, "src/ignored.ts", "export const ignored = true;\n")

	workspace := scanTestWorkspace(t, rootPath)

	expectedPaths := []string{"src/index.ts"}
	if got := scannedPaths(workspace.Files); !reflect.DeepEqual(got, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, got)
	}
}

func TestScanWorkspaceRespectsNestedGitIgnorePatterns(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, "src/.gitignore", "generated/drop.ts\n")
	writeTestFile(t, rootPath, "src/generated/.gitignore", "nested.ts\n")
	writeTestFile(t, rootPath, "src/generated/drop.ts", "export const drop = true;\n")
	writeTestFile(t, rootPath, "src/generated/keep.ts", "export const keep = true;\n")
	writeTestFile(t, rootPath, "src/generated/nested.ts", "export const nested = true;\n")

	workspace := scanTestWorkspace(t, rootPath)

	expectedPaths := []string{"src/generated/keep.ts"}
	if got := scannedPaths(workspace.Files); !reflect.DeepEqual(got, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, got)
	}
}

func TestScanWorkspaceIncludePatternRestrictsResults(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, "src/index.ts", "export const index = true;\n")
	writeTestFile(t, rootPath, "src/other.ts", "export const other = true;\n")
	writeTestFile(t, rootPath, "node_modules/kept/index.ts", "export const kept = true;\n")

	workspace := scanTestWorkspace(t, rootPath, WithIncludePatterns("node_modules/kept/index.ts"))

	expectedPaths := []string{
		"node_modules/kept/index.ts",
		"src/index.ts",
		"src/other.ts",
	}
	if got := scannedPaths(workspace.Files); !reflect.DeepEqual(got, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, got)
	}
}

func TestScanWorkspaceExcludePatternRemovesMatches(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, "src/index.ts", "export const index = true;\n")
	writeTestFile(t, rootPath, "src/excluded.ts", "export const excluded = true;\n")

	workspace := scanTestWorkspace(t, rootPath, WithExcludePatterns("src/excluded.ts"))

	expectedPaths := []string{"src/index.ts"}
	if got := scannedPaths(workspace.Files); !reflect.DeepEqual(got, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, got)
	}
}

func TestScanWorkspaceIgnoresUnsupportedExtensions(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, "src/index.py", "print('ignored')\n")
	writeTestFile(t, rootPath, "src/index.rb", "puts 'ignored'\n")
	writeTestFile(t, rootPath, "src/index.md", "# ignored\n")

	workspace := scanTestWorkspace(t, rootPath)

	if len(workspace.Files) != 0 {
		t.Fatalf("expected no scanned files, got %d", len(workspace.Files))
	}
}

func TestScanWorkspaceEmptyDirectoryReturnsEmptyWorkspace(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	workspace := scanTestWorkspace(t, rootPath)

	if len(workspace.Files) != 0 {
		t.Fatalf("expected no scanned files, got %d", len(workspace.Files))
	}

	if workspace.FilesByLanguage == nil {
		t.Fatal("expected FilesByLanguage to be initialized")
	}

	if len(workspace.FilesByLanguage) != 0 {
		t.Fatalf("expected no grouped files, got %d groups", len(workspace.FilesByLanguage))
	}
}

func TestScanWorkspaceGroupsFilesByLanguage(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()

	writeTestFile(t, rootPath, "src/a.ts", "export const a = true;\n")
	writeTestFile(t, rootPath, "src/b.ts", "export const b = true;\n")
	writeTestFile(t, rootPath, "src/c.js", "export const c = true;\n")

	workspace := scanTestWorkspace(t, rootPath)

	groupedPaths := groupedPaths(workspace)
	expected := map[string][]string{
		"js": {"src/c.js"},
		"ts": {"src/a.ts", "src/b.ts"},
	}

	if !reflect.DeepEqual(groupedPaths, expected) {
		t.Fatalf("expected grouped paths %v, got %v", expected, groupedPaths)
	}
}

func scanTestWorkspace(t *testing.T, rootPath string, opts ...Option) *models.ScannedWorkspace {
	t.Helper()

	workspace, err := ScanWorkspace(rootPath, opts...)
	if err != nil {
		t.Fatalf("ScanWorkspace(%s) returned error: %v", rootPath, err)
	}

	return workspace
}

func groupedCounts(workspace *models.ScannedWorkspace) map[string]int {
	counts := make(map[string]int, len(workspace.FilesByLanguage))
	for language, files := range workspace.FilesByLanguage {
		counts[string(language)] = len(files)
	}
	return counts
}

func groupedPaths(workspace *models.ScannedWorkspace) map[string][]string {
	grouped := make(map[string][]string, len(workspace.FilesByLanguage))
	for language, files := range workspace.FilesByLanguage {
		paths := make([]string, 0, len(files))
		for _, file := range files {
			paths = append(paths, file.RelativePath)
		}
		sort.Strings(paths)
		grouped[string(language)] = paths
	}
	return grouped
}

func scannedPaths(files []models.ScannedSourceFile) []string {
	paths := make([]string, 0, len(files))
	for _, file := range files {
		paths = append(paths, file.RelativePath)
	}
	sort.Strings(paths)
	return paths
}

func writeTestFile(t *testing.T, rootPath string, relativePath string, contents string) {
	t.Helper()

	absolutePath := filepath.Join(rootPath, relativePath)
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s) returned error: %v", filepath.Dir(absolutePath), err)
	}

	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) returned error: %v", absolutePath, err)
	}
}
