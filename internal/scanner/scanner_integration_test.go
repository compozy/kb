//go:build integration

package scanner

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanWorkspaceIntegrationNestedProject(t *testing.T) {
	rootPath := t.TempDir()
	outputPath := filepath.Join(rootPath, "artifacts", "vault")

	writeTestFile(t, rootPath, ".gitignore", "web/src/generated/drop.ts\n")
	writeTestFile(t, rootPath, "pkg/service/service.go", "package service\n")
	writeTestFile(t, rootPath, "pkg/service/internal/helper.go", "package internal\n")
	writeTestFile(t, rootPath, "web/src/index.tsx", "export const Page = () => null;\n")
	writeTestFile(t, rootPath, "web/src/generated/drop.ts", "export const drop = true;\n")
	writeTestFile(t, rootPath, "web/src/generated/keep.ts", "export const keep = true;\n")
	writeTestFile(t, rootPath, "node_modules/library/index.ts", "export const ignored = true;\n")
	writeTestFile(t, rootPath, "artifacts/vault/index.ts", "export const ignored = true;\n")

	workspace, err := ScanWorkspace(rootPath, WithOutputPath(outputPath))
	if err != nil {
		t.Fatalf("ScanWorkspace(%s) returned error: %v", rootPath, err)
	}

	expectedPaths := []string{
		"pkg/service/internal/helper.go",
		"pkg/service/service.go",
		"web/src/generated/keep.ts",
		"web/src/index.tsx",
	}

	gotPaths := make([]string, 0, len(workspace.Files))
	for _, file := range workspace.Files {
		gotPaths = append(gotPaths, file.RelativePath)
	}

	if !reflect.DeepEqual(gotPaths, expectedPaths) {
		t.Fatalf("expected scanned paths %v, got %v", expectedPaths, gotPaths)
	}
}
