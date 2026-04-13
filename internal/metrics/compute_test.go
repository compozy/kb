package metrics

import (
	"reflect"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestComputeMetricsReturnsEmptyResultForEmptyGraph(t *testing.T) {
	t.Parallel()

	result := ComputeMetrics(models.GraphSnapshot{RootPath: "/workspace"})

	if result.CircularDependencies == nil {
		t.Fatal("expected circular dependencies slice to be initialized")
	}
	if result.Directories == nil {
		t.Fatal("expected directories map to be initialized")
	}
	if result.Files == nil {
		t.Fatal("expected files map to be initialized")
	}
	if result.Symbols == nil {
		t.Fatal("expected symbols map to be initialized")
	}
	if len(result.CircularDependencies) != 0 || len(result.Directories) != 0 || len(result.Files) != 0 || len(result.Symbols) != 0 {
		t.Fatalf("expected empty metrics result, got %#v", result)
	}
}

func TestComputeMetricsSingleSymbolHasZeroBlastRadius(t *testing.T) {
	t.Parallel()

	symbol := symbolNode("single.ts", "single", "function", false, 1, 10, 0)
	graph := snapshot(
		[]models.GraphFile{fileNode("single.ts", symbol.ID)},
		[]models.SymbolNode{symbol},
	)

	result := ComputeMetrics(graph)
	metric := result.Symbols[symbol.ID]

	if metric.BlastRadius != 0 {
		t.Fatalf("blast radius = %d, want 0", metric.BlastRadius)
	}
}

func TestComputeMetricsCountsBlastRadiusAcrossDependents(t *testing.T) {
	t.Parallel()

	target := symbolNode("target.ts", "target", "function", true, 1, 10, 0)
	sourceA := symbolNode("a.ts", "useA", "function", false, 1, 5, 0)
	sourceB := symbolNode("b.ts", "useB", "function", false, 1, 5, 0)
	sourceC := symbolNode("c.ts", "useC", "function", false, 1, 5, 0)

	graph := snapshot(
		[]models.GraphFile{
			fileNode("target.ts", target.ID),
			fileNode("a.ts", sourceA.ID),
			fileNode("b.ts", sourceB.ID),
			fileNode("c.ts", sourceC.ID),
		},
		[]models.SymbolNode{target, sourceA, sourceB, sourceC},
		relation(sourceA.ID, target.ID, models.RelCalls),
		relation(sourceB.ID, target.ID, models.RelReferences),
		relation(sourceC.ID, target.ID, models.RelCalls),
	)

	result := ComputeMetrics(graph)
	metric := result.Symbols[target.ID]

	if metric.BlastRadius < 3 {
		t.Fatalf("blast radius = %d, want >= 3", metric.BlastRadius)
	}
	if metric.DirectDependents != 3 {
		t.Fatalf("direct dependents = %d, want 3", metric.DirectDependents)
	}
}

func TestComputeMetricsCountsFileEfferentCoupling(t *testing.T) {
	t.Parallel()

	graph := snapshot(
		[]models.GraphFile{
			fileNode("a.ts"),
			fileNode("b.ts"),
			fileNode("c.ts"),
		},
		nil,
		relation("file:a.ts", "file:b.ts", models.RelImports),
		relation("file:a.ts", "file:c.ts", models.RelImports),
	)

	result := ComputeMetrics(graph)
	metric := result.Files["file:a.ts"]

	if metric.EfferentCoupling != 2 {
		t.Fatalf("efferent coupling = %d, want 2", metric.EfferentCoupling)
	}
}

func TestComputeMetricsCountsFileAfferentCoupling(t *testing.T) {
	t.Parallel()

	graph := snapshot(
		[]models.GraphFile{
			fileNode("shared.ts"),
			fileNode("a.ts"),
			fileNode("b.ts"),
			fileNode("c.ts"),
			fileNode("d.ts"),
		},
		nil,
		relation("file:a.ts", "file:shared.ts", models.RelImports),
		relation("file:b.ts", "file:shared.ts", models.RelImports),
		relation("file:c.ts", "file:shared.ts", models.RelImports),
		relation("file:d.ts", "file:shared.ts", models.RelImports),
	)

	result := ComputeMetrics(graph)
	metric := result.Files["file:shared.ts"]

	if metric.AfferentCoupling != 4 {
		t.Fatalf("afferent coupling = %d, want 4", metric.AfferentCoupling)
	}
}

func TestComputeMetricsComputesBalancedInstability(t *testing.T) {
	t.Parallel()

	graph := snapshot(
		[]models.GraphFile{
			fileNode("a.ts"),
			fileNode("mid.ts"),
			fileNode("c.ts"),
		},
		nil,
		relation("file:a.ts", "file:mid.ts", models.RelImports),
		relation("file:mid.ts", "file:c.ts", models.RelImports),
	)

	result := ComputeMetrics(graph)
	metric := result.Files["file:mid.ts"]

	if metric.Instability != 0.5 {
		t.Fatalf("instability = %v, want 0.5", metric.Instability)
	}
}

func TestComputeMetricsFlagsDeadExports(t *testing.T) {
	t.Parallel()

	symbol := symbolNode("util.ts", "UnusedExport", "function", true, 1, 8, 0)
	graph := snapshot(
		[]models.GraphFile{fileNode("util.ts", symbol.ID)},
		[]models.SymbolNode{symbol},
	)

	result := ComputeMetrics(graph)
	metric := result.Symbols[symbol.ID]

	if !metric.IsDeadExport {
		t.Fatal("expected dead export flag to be true")
	}
	if !containsString(metric.Smells, "dead-export") {
		t.Fatalf("expected dead-export smell, got %#v", metric.Smells)
	}
}

func TestComputeMetricsFlagsLongFunctions(t *testing.T) {
	t.Parallel()

	symbol := symbolNode("worker.ts", "process", "function", false, 1, 75, 0)
	graph := snapshot(
		[]models.GraphFile{fileNode("worker.ts", symbol.ID)},
		[]models.SymbolNode{symbol},
	)

	result := ComputeMetrics(graph)
	metric := result.Symbols[symbol.ID]

	if !metric.IsLongFunction {
		t.Fatal("expected long function flag to be true")
	}
	if !containsString(metric.Smells, "long-function") {
		t.Fatalf("expected long-function smell, got %#v", metric.Smells)
	}
}

func TestComputeMetricsDetectsCircularDependencies(t *testing.T) {
	t.Parallel()

	graph := snapshot(
		[]models.GraphFile{
			fileNode("a.ts"),
			fileNode("b.ts"),
			fileNode("c.ts"),
		},
		nil,
		relation("file:a.ts", "file:b.ts", models.RelImports),
		relation("file:b.ts", "file:c.ts", models.RelImports),
		relation("file:c.ts", "file:a.ts", models.RelImports),
	)

	result := ComputeMetrics(graph)

	expectedCycle := [][]string{{"a.ts", "b.ts", "c.ts"}}
	if !reflect.DeepEqual(result.CircularDependencies, expectedCycle) {
		t.Fatalf("circular dependencies = %#v, want %#v", result.CircularDependencies, expectedCycle)
	}

	for _, fileID := range []string{"file:a.ts", "file:b.ts", "file:c.ts"} {
		if !result.Files[fileID].HasCircularDependency {
			t.Fatalf("expected %s to be marked as circular", fileID)
		}
	}
}

func TestFindCircularDependencyGroupsMergesOverlappingCyclesIntoSingleComponent(t *testing.T) {
	t.Parallel()

	got := FindCircularDependencyGroups(map[string][]string{
		"d.ts": {"a.ts"},
		"a.ts": {"b.ts"},
		"c.ts": {"a.ts", "d.ts"},
		"b.ts": {"c.ts"},
	})

	want := [][]string{{"a.ts", "b.ts", "c.ts", "d.ts"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("circular dependency groups = %#v, want %#v", got, want)
	}
}

func TestFindCircularDependencyGroupsReturnsStableSortedComponents(t *testing.T) {
	t.Parallel()

	got := FindCircularDependencyGroups(map[string][]string{
		"d.ts": {"c.ts"},
		"c.ts": {"d.ts"},
		"b.ts": {"a.ts"},
		"a.ts": {"b.ts"},
	})

	want := [][]string{
		{"a.ts", "b.ts"},
		{"c.ts", "d.ts"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("circular dependency groups = %#v, want %#v", got, want)
	}
}

func TestComputeMetricsReturnsNoCircularDependenciesForAcyclicGraph(t *testing.T) {
	t.Parallel()

	graph := snapshot(
		[]models.GraphFile{
			fileNode("a.ts"),
			fileNode("b.ts"),
			fileNode("c.ts"),
		},
		nil,
		relation("file:a.ts", "file:b.ts", models.RelImports),
		relation("file:b.ts", "file:c.ts", models.RelImports),
	)

	result := ComputeMetrics(graph)

	if len(result.CircularDependencies) != 0 {
		t.Fatalf("expected no circular dependencies, got %#v", result.CircularDependencies)
	}
}

func TestComputeMetricsAggregatesDirectoryMetrics(t *testing.T) {
	t.Parallel()

	graph := snapshot(
		[]models.GraphFile{
			fileNode("app/main.ts"),
			fileNode("lib/util.ts"),
			fileNode("lib/extra.ts"),
			fileNode("core/base.ts"),
		},
		nil,
		relation("file:app/main.ts", "file:lib/util.ts", models.RelImports),
		relation("file:app/main.ts", "file:lib/extra.ts", models.RelImports),
		relation("file:lib/util.ts", "file:core/base.ts", models.RelImports),
	)

	result := ComputeMetrics(graph)

	if metric := result.Directories["app"]; metric.AfferentCoupling != 0 || metric.EfferentCoupling != 2 || metric.Instability != 1 {
		t.Fatalf("app metrics = %#v, want afferent=0 efferent=2 instability=1", metric)
	}
	if metric := result.Directories["lib"]; metric.AfferentCoupling != 1 || metric.EfferentCoupling != 1 || metric.Instability != 0.5 {
		t.Fatalf("lib metrics = %#v, want afferent=1 efferent=1 instability=0.5", metric)
	}
	if metric := result.Directories["core"]; metric.AfferentCoupling != 1 || metric.EfferentCoupling != 0 || metric.Instability != 0 {
		t.Fatalf("core metrics = %#v, want afferent=1 efferent=0 instability=0", metric)
	}
}

func snapshot(files []models.GraphFile, symbols []models.SymbolNode, relations ...models.RelationEdge) models.GraphSnapshot {
	if files == nil {
		files = []models.GraphFile{}
	}
	if symbols == nil {
		symbols = []models.SymbolNode{}
	}
	if relations == nil {
		relations = []models.RelationEdge{}
	}

	return models.GraphSnapshot{
		RootPath:      "/workspace",
		Files:         files,
		Symbols:       symbols,
		ExternalNodes: []models.ExternalNode{},
		Relations:     relations,
		Diagnostics:   []models.StructuredDiagnostic{},
	}
}

func fileNode(filePath string, symbolIDs ...string) models.GraphFile {
	return models.GraphFile{
		ID:        "file:" + filePath,
		NodeType:  "file",
		FilePath:  filePath,
		Language:  models.LangTS,
		SymbolIDs: append([]string(nil), symbolIDs...),
	}
}

func symbolNode(
	filePath, name, symbolKind string,
	exported bool,
	startLine, endLine, cyclomaticComplexity int,
) models.SymbolNode {
	return models.SymbolNode{
		ID:                   "symbol:" + filePath + ":" + name,
		NodeType:             "symbol",
		Name:                 name,
		SymbolKind:           symbolKind,
		Language:             models.LangTS,
		FilePath:             filePath,
		StartLine:            startLine,
		EndLine:              endLine,
		Exported:             exported,
		CyclomaticComplexity: cyclomaticComplexity,
	}
}

func relation(fromID, toID string, relationType models.RelationType) models.RelationEdge {
	return models.RelationEdge{
		FromID:     fromID,
		ToID:       toID,
		Type:       relationType,
		Confidence: models.ConfidenceSyntactic,
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
