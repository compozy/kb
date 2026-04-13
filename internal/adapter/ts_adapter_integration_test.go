//go:build integration

package adapter

import (
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestTSAdapterIntegrationParsesMultiFileProject(t *testing.T) {
	parsedFiles := parseTSLikeSources(t, map[string]string{
		"core/math.ts": `
export function add(left: number, right: number): number {
	return left + right
}
`,
		"core/barrel.ts": `
export { add } from "./math"
`,
		"app/main.ts": `
import { add } from "../core/barrel"

export function run(): number {
	return add(20, 22)
}
`,
	})

	barrelFile := mustFindParsedFile(t, parsedFiles, "core/barrel.ts")
	mathFile := mustFindParsedFile(t, parsedFiles, "core/math.ts")
	mainFile := mustFindParsedFile(t, parsedFiles, "app/main.ts")

	add := mustFindSymbol(t, mathFile.Symbols, "add")
	run := mustFindSymbol(t, mainFile.Symbols, "run")

	if !hasRelation(barrelFile.Relations, barrelFile.File.ID, add.ID, models.RelExports) {
		t.Fatal("expected barrel.ts to export add")
	}
	if !hasRelation(mainFile.Relations, mainFile.File.ID, add.ID, models.RelReferences) {
		t.Fatal("expected app/main.ts to reference add through barrel.ts")
	}
	if !hasRelation(mainFile.Relations, run.ID, add.ID, models.RelCalls) {
		t.Fatal("expected run() call to resolve to add() through barrel.ts")
	}
}
