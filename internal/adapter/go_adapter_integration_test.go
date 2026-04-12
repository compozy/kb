//go:build integration

package adapter

import (
	"testing"

	"github.com/user/kb/internal/models"
)

func TestGoAdapterBuildsCrossFileCallRelations(t *testing.T) {
	t.Parallel()

	parsedFiles := parseGoSources(t, map[string]string{
		"helper.go": `
package sample

func Helper() {}
`,
		"run.go": `
package sample

func Run() {
	Helper()
}
`,
	})

	if len(parsedFiles) != 2 {
		t.Fatalf("expected 2 parsed files, got %d", len(parsedFiles))
	}

	runFile := parsedFiles[1]
	helper := mustFindSymbol(t, parsedFiles[0].Symbols, "Helper")
	run := mustFindSymbol(t, runFile.Symbols, "Run")

	if !hasRelation(runFile.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected cross-file call relation from Run to Helper")
	}
}
