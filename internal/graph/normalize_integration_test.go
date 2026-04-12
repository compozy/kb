//go:build integration

package graph

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/kb/internal/adapter"
	"github.com/user/kb/internal/models"
)

func TestNormalizeGraphMergesOverlappingImportsAcrossParsedFiles(t *testing.T) {
	rootDir := t.TempDir()

	writeFixtureFile(t, rootDir, "helper.go", `
package sample

import "fmt"

func helper() int {
	fmt.Println("helper")
	return 1
}
`)
	writeFixtureFile(t, rootDir, "run.go", `
package sample

import "fmt"

func Run() int {
	fmt.Println("run")
	return helper()
}
`)

	files := []models.ScannedSourceFile{
		{
			AbsolutePath: filepath.Join(rootDir, "helper.go"),
			RelativePath: "helper.go",
			Language:     models.LangGo,
		},
		{
			AbsolutePath: filepath.Join(rootDir, "run.go"),
			RelativePath: "run.go",
			Language:     models.LangGo,
		},
	}

	parsedFiles, err := (adapter.GoAdapter{}).ParseFiles(files, rootDir)
	if err != nil {
		t.Fatalf("parse files: %v", err)
	}

	snapshot := NormalizeGraph(rootDir, parsedFiles)

	if len(snapshot.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(snapshot.Files))
	}

	if len(snapshot.ExternalNodes) != 1 {
		t.Fatalf("expected shared fmt import to deduplicate to 1 external node, got %d", len(snapshot.ExternalNodes))
	}

	if snapshot.ExternalNodes[0].ID != "external:fmt" {
		t.Fatalf("external node ID = %q, want %q", snapshot.ExternalNodes[0].ID, "external:fmt")
	}

	helperFile := snapshot.Files[0]
	runFile := snapshot.Files[1]

	if helperFile.ID != "file:helper.go" || runFile.ID != "file:run.go" {
		t.Fatalf("unexpected file order: %#v", snapshot.Files)
	}

	if len(helperFile.SymbolIDs) == 0 || len(runFile.SymbolIDs) == 0 {
		t.Fatalf("expected symbol IDs to be attached to both files: %#v", snapshot.Files)
	}

	if !hasRelation(snapshot.Relations, "file:helper.go", "external:fmt", models.RelImports) {
		t.Fatal("expected helper.go import relation to fmt")
	}

	if !hasRelation(snapshot.Relations, "file:run.go", "external:fmt", models.RelImports) {
		t.Fatal("expected run.go import relation to fmt")
	}
}

func writeFixtureFile(t *testing.T, rootDir string, relativePath string, contents string) {
	t.Helper()

	absolutePath := filepath.Join(rootDir, relativePath)
	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func hasRelation(relations []models.RelationEdge, fromID string, toID string, relationType models.RelationType) bool {
	for _, relation := range relations {
		if relation.FromID == fromID && relation.ToID == toID && relation.Type == relationType {
			return true
		}
	}

	return false
}
