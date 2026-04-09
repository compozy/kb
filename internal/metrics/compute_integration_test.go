//go:build integration

package metrics

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/go-devstack/internal/adapter"
	"github.com/user/go-devstack/internal/graph"
	"github.com/user/go-devstack/internal/models"
)

func TestComputeMetricsIntegrationOnMultiDirectoryProject(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()
	sources := map[string]string{
		"core/math.ts": `
export function add(left: number, right: number): number {
	return left + right
}
`,
		"core/logger.ts": `
export function log(message: string): string {
	return message
}
`,
		"services/order.ts": `
import { add } from "../core/math"
import { log } from "../core/logger"

export function runOrder(): number {
	log("start")
	return add(20, 22)
}
`,
		"commands/run.ts": `
import { runOrder } from "../services/order"

export function main(): number {
	return runOrder()
}
`,
	}

	files := writeTSWorkspace(t, rootDir, sources)
	parsedFiles, err := (adapter.TSAdapter{}).ParseFiles(files, rootDir)
	if err != nil {
		t.Fatalf("parse files: %v", err)
	}

	snapshot := graph.NormalizeGraph(rootDir, parsedFiles)
	result := ComputeMetrics(snapshot)

	if len(result.CircularDependencies) != 0 {
		t.Fatalf("expected no circular dependencies, got %#v", result.CircularDependencies)
	}

	commandMetrics := result.Files["file:commands/run.ts"]
	if !commandMetrics.IsEntryPoint {
		t.Fatal("expected commands/run.ts to be treated as an entry point")
	}

	serviceMetrics := result.Files["file:services/order.ts"]
	if serviceMetrics.AfferentCoupling != 1 {
		t.Fatalf("service afferent coupling = %d, want 1", serviceMetrics.AfferentCoupling)
	}
	if serviceMetrics.EfferentCoupling != 2 {
		t.Fatalf("service efferent coupling = %d, want 2", serviceMetrics.EfferentCoupling)
	}
	if serviceMetrics.Instability != 0.6667 {
		t.Fatalf("service instability = %v, want 0.6667", serviceMetrics.Instability)
	}

	servicesDirectoryMetrics := result.Directories["services"]
	if servicesDirectoryMetrics.AfferentCoupling != 1 || servicesDirectoryMetrics.EfferentCoupling != 2 {
		t.Fatalf("services directory metrics = %#v, want afferent=1 efferent=2", servicesDirectoryMetrics)
	}

	runOrder := findSymbolByNameAndFile(t, snapshot.Symbols, "runOrder", "services/order.ts")
	runOrderMetrics := result.Symbols[runOrder.ID]
	if runOrderMetrics.IsDeadExport {
		t.Fatal("expected runOrder to be referenced and not marked as dead export")
	}
	if runOrderMetrics.DirectDependents < 1 {
		t.Fatalf("runOrder direct dependents = %d, want >= 1", runOrderMetrics.DirectDependents)
	}
}

func writeTSWorkspace(t *testing.T, rootDir string, sources map[string]string) []models.ScannedSourceFile {
	t.Helper()

	relativePaths := make([]string, 0, len(sources))
	for relativePath := range sources {
		relativePaths = append(relativePaths, relativePath)
	}
	sort.Strings(relativePaths)

	files := make([]models.ScannedSourceFile, 0, len(relativePaths))
	for _, relativePath := range relativePaths {
		absolutePath := filepath.Join(rootDir, relativePath)
		if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", relativePath, err)
		}
		if err := os.WriteFile(absolutePath, []byte(sources[relativePath]), 0o644); err != nil {
			t.Fatalf("write %s: %v", relativePath, err)
		}

		files = append(files, models.ScannedSourceFile{
			AbsolutePath: absolutePath,
			RelativePath: relativePath,
			Language:     models.LangTS,
		})
	}

	return files
}

func findSymbolByNameAndFile(t *testing.T, symbols []models.SymbolNode, name, filePath string) models.SymbolNode {
	t.Helper()

	for _, symbol := range symbols {
		if symbol.Name == name && symbol.FilePath == filePath {
			return symbol
		}
	}

	t.Fatalf("symbol %s in %s not found", name, filePath)
	return models.SymbolNode{}
}
