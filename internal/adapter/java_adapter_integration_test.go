//go:build integration

package adapter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestJavaAdapterBuildsCrossFileImportAndCallRelations(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"modules/shared/src/com/shared/Helper.java": `
package com.shared;

public class Helper {
    public static void assist() {
    }
}
`,
		"modules/app/src/com/example/Runner.java": `
package com.example;

import com.shared.Helper;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	})

	helperFile := mustFindParsedFile(t, parsedFiles, "modules/shared/src/com/shared/Helper.java")
	runnerFile := mustFindParsedFile(t, parsedFiles, "modules/app/src/com/example/Runner.java")

	helperClass := mustFindSymbol(t, helperFile.Symbols, "Helper")
	assistMethod := mustFindSymbol(t, helperFile.Symbols, "assist")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	helperImportExternal := mustFindExternalNode(t, runnerFile.ExternalNodes, "com.shared.Helper")

	if !hasRelation(runnerFile.Relations, runnerFile.File.ID, helperImportExternal.ID, models.RelImports) {
		t.Fatal("expected import relation for com.shared.Helper")
	}
	if !hasRelation(runnerFile.Relations, runnerFile.File.ID, helperClass.ID, models.RelReferences) {
		t.Fatal("expected reference relation from Runner.java to Helper class symbol")
	}
	if !hasRelation(runnerFile.Relations, runMethod.ID, assistMethod.ID, models.RelCalls) {
		t.Fatal("expected cross-file call relation from run() to assist()")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runnerFile.File.ID,
		helperClass.ID,
		models.RelReferences,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic reference relation from Runner.java to Helper class symbol")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic call relation from run() to assist()")
	}
}

func TestJavaAdapterBuildsWildcardImportRelationsAcrossFiles(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"modules/shared/src/com/shared/Helper.java": `
package com.shared;

public class Helper {
    public static void assist() {
    }
}
`,
		"modules/shared/src/com/shared/Util.java": `
package com.shared;

public class Util {
    public static void noop() {
    }
}
`,
		"modules/app/src/com/example/Runner.java": `
package com.example;

import com.shared.*;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	})

	helperFile := mustFindParsedFile(t, parsedFiles, "modules/shared/src/com/shared/Helper.java")
	runnerFile := mustFindParsedFile(t, parsedFiles, "modules/app/src/com/example/Runner.java")

	helperClass := mustFindSymbol(t, helperFile.Symbols, "Helper")
	assistMethod := mustFindSymbol(t, helperFile.Symbols, "assist")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")

	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runnerFile.File.ID,
		helperClass.ID,
		models.RelReferences,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic wildcard reference relation from Runner.java to Helper class symbol")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic wildcard call relation from run() to assist()")
	}
}

func TestJavaAdapterOutputIsDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	sources := map[string]string{
		"src/com/example/Alpha.java": `
package com.example;

public class Alpha {
    public void a() {}
}
`,
		"src/com/example/Beta.java": `
package com.example;

import com.example.Alpha;

public class Beta {
    public void b() {
        new Alpha().a();
    }
}
`,
	}

	first := parseJavaSources(t, sources)
	second := parseJavaSources(t, sources)

	if !reflect.DeepEqual(first, second) {
		t.Fatal("expected deterministic Java adapter output across repeated parse runs")
	}
}

func TestJavaAdapterPartialMetadataFallsBackWithoutFailingIngest(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"src/com/example/Helper.java": `
package com.example;

public class Helper {
    public static void assist() {}
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.external.Missing;

public class Runner {
    public void run() {
        assist();
    }
}
`,
	})

	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
	helperFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Helper.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	assistMethod := mustFindSymbol(t, helperFile.Symbols, "assist")

	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSyntactic,
	) {
		t.Fatal("expected fallback syntactic call relation with partial metadata")
	}

	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if fallbackDiagnostic.Severity != models.SeverityWarning {
		t.Fatalf("fallback diagnostic severity = %q, want %q", fallbackDiagnostic.Severity, models.SeverityWarning)
	}
}

func TestJavaAdapterMissingWildcardPackageFallsBackWithoutFailingIngest(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"src/com/example/Runner.java": `
package com.example;

import com.unknown.*;

public class Runner {
    public void run() {}
}
`,
	})

	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if fallbackDiagnostic.Severity != models.SeverityWarning {
		t.Fatalf("fallback diagnostic severity = %q, want %q", fallbackDiagnostic.Severity, models.SeverityWarning)
	}
	if !hasRelation(runnerFile.Relations, runnerFile.File.ID, createExternalID("com.unknown.*"), models.RelImports) {
		t.Fatal("expected wildcard import relation to external node even when deep resolution falls back")
	}
}

func TestJavaAdapterResolvesNestedAndTopLevelRelationsAcrossFiles(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"modules/shared/src/com/shared/Helper.java": `
package com.shared;

public class Helper {
    public static void assistTopLevel() {
    }
}
`,
		"modules/shared/src/com/shared/Outer.java": `
package com.shared;

public class Outer {
    public static class Inner {
        public static void assistNested() {
        }
    }
}
`,
		"modules/app/src/com/example/Runner.java": `
package com.example;

import com.shared.Helper;
import com.shared.Outer.Inner;

public class Runner {
    public void run() {
        Helper.assistTopLevel();
        Inner.assistNested();
    }
}
`,
	})

	helperFile := mustFindParsedFile(t, parsedFiles, "modules/shared/src/com/shared/Helper.java")
	outerFile := mustFindParsedFile(t, parsedFiles, "modules/shared/src/com/shared/Outer.java")
	runnerFile := mustFindParsedFile(t, parsedFiles, "modules/app/src/com/example/Runner.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	topLevelHelperClass := mustFindSymbol(t, helperFile.Symbols, "Helper")
	topLevelAssistMethod := mustFindSymbol(t, helperFile.Symbols, "assistTopLevel")
	nestedInnerClass := mustFindSymbol(t, outerFile.Symbols, "Outer.Inner")
	nestedAssistMethod := mustFindSymbol(t, outerFile.Symbols, "assistNested")

	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runnerFile.File.ID,
		topLevelHelperClass.ID,
		models.RelReferences,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected top-level semantic reference relation from Runner.java to Helper")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		topLevelAssistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected top-level semantic call relation from run() to assistTopLevel()")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runnerFile.File.ID,
		nestedInnerClass.ID,
		models.RelReferences,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected nested semantic reference relation from Runner.java to Outer.Inner")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		nestedAssistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected nested semantic call relation from run() to assistNested()")
	}
}

func TestJavaAdapterAmbiguousImportsRemainStableAcrossRuns(t *testing.T) {
	t.Parallel()

	sources := map[string]string{
		"src/com/alpha/Helper.java": `
package com.alpha;

public class Helper {
    public static void assist() {}
}
`,
		"src/com/beta/Helper.java": `
package com.beta;

public class Helper {
    public static void assist() {}
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.alpha.Helper;
import com.beta.Helper;
import static com.alpha.Helper.assist;
import static com.beta.Helper.assist;

public class Runner {
    public void run() {
        Helper.assist();
        assist();
    }
}
`,
	}

	first := parseJavaSources(t, sources)
	second := parseJavaSources(t, sources)
	if !reflect.DeepEqual(first, second) {
		t.Fatal("expected deterministic Java adapter output across repeated ambiguous import runs")
	}

	runnerFile := mustFindParsedFile(t, first, "src/com/example/Runner.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	for _, relation := range runnerFile.Relations {
		if relation.FromID == runMethod.ID && relation.Type == models.RelCalls {
			t.Fatalf("did not expect ambiguous run() call relations, got %+v", relation)
		}
	}

	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if fallbackDiagnostic.Severity != models.SeverityWarning {
		t.Fatalf("fallback diagnostic severity = %q, want %q", fallbackDiagnostic.Severity, models.SeverityWarning)
	}
}

func TestJavaAdapterModuleHintsImproveAmbiguousImportResolution(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSourcesWithRepositoryFiles(
		t,
		map[string]string{
			"settings.gradle": strings.Join([]string{
				`rootProject.name = "atlas"`,
				`include("shared-a", "shared-b", "app")`,
				"",
			}, "\n"),
			"app/build.gradle": strings.Join([]string{
				"dependencies {",
				`    implementation(project(":shared-b"))`,
				"}",
				"",
			}, "\n"),
		},
		map[string]string{
			"shared-a/src/main/java/com/acme/alpha/Helper.java": `
package com.acme.alpha;

public class Helper {
    public static void assist() {}
}
`,
			"shared-b/src/main/java/com/acme/beta/Helper.java": `
package com.acme.beta;

public class Helper {
    public static void assist() {}
}
`,
			"app/src/main/java/com/acme/app/Runner.java": `
package com.acme.app;

import com.acme.alpha.Helper;
import com.acme.beta.Helper;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
		},
	)

	runnerFile := mustFindParsedFile(t, parsedFiles, "app/src/main/java/com/acme/app/Runner.java")
	alphaFile := mustFindParsedFile(t, parsedFiles, "shared-a/src/main/java/com/acme/alpha/Helper.java")
	betaFile := mustFindParsedFile(t, parsedFiles, "shared-b/src/main/java/com/acme/beta/Helper.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	alphaAssist := mustFindSymbol(t, alphaFile.Symbols, "assist")
	betaAssist := mustFindSymbol(t, betaFile.Symbols, "assist")

	if hasRelation(runnerFile.Relations, runMethod.ID, alphaAssist.ID, models.RelCalls) {
		t.Fatal("did not expect module-hinted call relation to resolve to shared-a helper")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		betaAssist.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected module-hinted semantic call relation to shared-b helper")
	}
}

func TestJavaAdapterWithoutModuleHintsKeepsAmbiguousFallbackPath(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"shared-a/src/main/java/com/acme/alpha/Helper.java": `
package com.acme.alpha;

public class Helper {
    public static void assist() {}
}
`,
		"shared-b/src/main/java/com/acme/beta/Helper.java": `
package com.acme.beta;

public class Helper {
    public static void assist() {}
}
`,
		"app/src/main/java/com/acme/app/Runner.java": `
package com.acme.app;

import com.acme.alpha.Helper;
import com.acme.beta.Helper;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	})

	runnerFile := mustFindParsedFile(t, parsedFiles, "app/src/main/java/com/acme/app/Runner.java")
	alphaFile := mustFindParsedFile(t, parsedFiles, "shared-a/src/main/java/com/acme/alpha/Helper.java")
	betaFile := mustFindParsedFile(t, parsedFiles, "shared-b/src/main/java/com/acme/beta/Helper.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	alphaAssist := mustFindSymbol(t, alphaFile.Symbols, "assist")
	betaAssist := mustFindSymbol(t, betaFile.Symbols, "assist")

	if hasRelation(runnerFile.Relations, runMethod.ID, alphaAssist.ID, models.RelCalls) {
		t.Fatal("did not expect fallback path to emit alpha call relation under ambiguous imports")
	}
	if hasRelation(runnerFile.Relations, runMethod.ID, betaAssist.ID, models.RelCalls) {
		t.Fatal("did not expect fallback path to emit beta call relation under ambiguous imports")
	}

	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if !strings.Contains(fallbackDiagnostic.Detail, "calls:Helper.assist (ambiguous-import-class)") {
		t.Fatalf("fallback diagnostic detail = %q", fallbackDiagnostic.Detail)
	}
}

func TestJavaAdapterPhase2EnterpriseScenarioRegression(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSourcesWithRepositoryFiles(
		t,
		map[string]string{
			"settings.gradle": strings.Join([]string{
				`rootProject.name = "atlas"`,
				`include("shared-a", "shared-b", "app")`,
				"",
			}, "\n"),
			"app/build.gradle": strings.Join([]string{
				"dependencies {",
				`    implementation(project(":shared-b"))`,
				"}",
				"",
			}, "\n"),
		},
		map[string]string{
			"shared-a/src/main/java/com/acme/shareda/Helper.java": `
package com.acme.shareda;

public class Helper {
    public static void assist() {}
}
`,
			"shared-b/src/main/java/com/acme/sharedb/Helper.java": `
package com.acme.sharedb;

public class Helper {
    public static void assist() {}
}
`,
			"shared-b/src/main/java/com/acme/sharedb/Tooling.java": `
package com.acme.sharedb;

public class Tooling {
    public static void noop() {}
}
`,
			"shared-b/src/main/java/com/acme/sharedb/Outer.java": `
package com.acme.sharedb;

public class Outer {
    public static class Inner {
        public static void assistNested() {}
    }
}
`,
			"app/src/main/java/com/acme/app/Runner.java": `
package com.acme.app;

import com.acme.shareda.Helper;
import com.acme.sharedb.Helper;
import com.acme.sharedb.*;
import com.acme.sharedb.Outer.Inner;
import com.acme.missing.*;

public class Runner {
    public void run() {
        Helper.assist();
        Inner.assistNested();
        Tooling.noop();
    }
}
`,
		},
	)

	runnerFile := mustFindParsedFile(t, parsedFiles, "app/src/main/java/com/acme/app/Runner.java")
	sharedBHelper := mustFindParsedFile(t, parsedFiles, "shared-b/src/main/java/com/acme/sharedb/Helper.java")
	outerFile := mustFindParsedFile(t, parsedFiles, "shared-b/src/main/java/com/acme/sharedb/Outer.java")
	toolingFile := mustFindParsedFile(t, parsedFiles, "shared-b/src/main/java/com/acme/sharedb/Tooling.java")

	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	betaAssist := mustFindSymbol(t, sharedBHelper.Symbols, "assist")
	innerClass := mustFindSymbol(t, outerFile.Symbols, "Outer.Inner")
	nestedAssist := mustFindSymbol(t, outerFile.Symbols, "assistNested")
	toolingClass := mustFindSymbol(t, toolingFile.Symbols, "Tooling")
	toolingNoop := mustFindSymbol(t, toolingFile.Symbols, "noop")

	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		betaAssist.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected module-assisted semantic call relation to shared-b Helper.assist")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runnerFile.File.ID,
		innerClass.ID,
		models.RelReferences,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic nested type reference relation to Outer.Inner")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		nestedAssist.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic nested call relation to Outer.Inner.assistNested")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runnerFile.File.ID,
		toolingClass.ID,
		models.RelReferences,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic wildcard reference relation to Tooling")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		toolingNoop.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic wildcard call relation to Tooling.noop")
	}

	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if !strings.Contains(fallbackDiagnostic.Detail, "references:com.acme.missing.* (missing-wildcard-package)") {
		t.Fatalf("fallback diagnostic detail = %q, want missing wildcard package reason", fallbackDiagnostic.Detail)
	}
}
