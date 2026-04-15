package adapter

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestJavaAdapterSupportsOnlyJava(t *testing.T) {
	t.Parallel()

	adapter := JavaAdapter{}
	if !adapter.Supports(models.LangJava) {
		t.Fatal("expected JavaAdapter to support Java")
	}

	for _, language := range []models.SupportedLanguage{
		models.LangTS,
		models.LangTSX,
		models.LangJS,
		models.LangJSX,
		models.LangGo,
		models.LangRust,
	} {
		if adapter.Supports(language) {
			t.Fatalf("expected JavaAdapter to reject %q", language)
		}
	}
}

func TestJavaAdapterParsesPackageClassAndMethodSymbols(t *testing.T) {
	t.Parallel()

	parsed := parseSingleJavaFile(t, "src/com/example/Runner.java", `
package com.example;

public class Runner {
    public void run() {
    }
}
`)

	packageSymbol := mustFindSymbol(t, parsed.Symbols, "com.example")
	if packageSymbol.SymbolKind != javaSymbolKindPackage {
		t.Fatalf("package symbol kind = %q, want %q", packageSymbol.SymbolKind, javaSymbolKindPackage)
	}
	if packageSymbol.Signature != "package com.example" {
		t.Fatalf("package signature = %q", packageSymbol.Signature)
	}

	classSymbol := mustFindSymbol(t, parsed.Symbols, "Runner")
	if classSymbol.SymbolKind != javaSymbolKindClass {
		t.Fatalf("class symbol kind = %q, want %q", classSymbol.SymbolKind, javaSymbolKindClass)
	}
	if classSymbol.Signature != "class Runner" {
		t.Fatalf("class signature = %q", classSymbol.Signature)
	}
	if !classSymbol.Exported {
		t.Fatal("expected class symbol to be exported")
	}

	methodSymbol := mustFindSymbol(t, parsed.Symbols, "run")
	if methodSymbol.SymbolKind != javaSymbolKindMethod {
		t.Fatalf("method symbol kind = %q, want %q", methodSymbol.SymbolKind, javaSymbolKindMethod)
	}
	if !strings.Contains(methodSymbol.Signature, "void run(") {
		t.Fatalf("method signature = %q", methodSymbol.Signature)
	}
	if !methodSymbol.Exported {
		t.Fatal("expected method symbol to be exported")
	}

	if !hasRelation(parsed.Relations, parsed.File.ID, classSymbol.ID, models.RelContains) {
		t.Fatal("expected file contains relation for class symbol")
	}
	if !hasRelation(parsed.Relations, parsed.File.ID, methodSymbol.ID, models.RelContains) {
		t.Fatal("expected file contains relation for method symbol")
	}
}

func TestJavaAdapterExtractsImportRelationsAndExternalNodes(t *testing.T) {
	t.Parallel()

	parsed := parseSingleJavaFile(t, "src/com/example/Runner.java", `
package com.example;

import java.util.List;
import static com.example.Helper.assist;

public class Runner {
    public void run(List<String> values) {
        assist();
    }
}
`)

	if len(parsed.ExternalNodes) != 2 {
		t.Fatalf("expected 2 external nodes, got %d", len(parsed.ExternalNodes))
	}

	javaUtilList := mustFindExternalNode(t, parsed.ExternalNodes, "java.util.List")
	if !hasRelation(parsed.Relations, parsed.File.ID, javaUtilList.ID, models.RelImports) {
		t.Fatal("expected import relation for java.util.List")
	}

	staticImport := mustFindExternalNode(t, parsed.ExternalNodes, "com.example.Helper.assist")
	if !hasRelation(parsed.Relations, parsed.File.ID, staticImport.ID, models.RelImports) {
		t.Fatal("expected import relation for static import")
	}
}

func TestJavaAdapterProducesDiagnosticsForParseErrors(t *testing.T) {
	t.Parallel()

	parsed := parseSingleJavaFile(t, "src/com/example/Broken.java", `
package com.example;

public class Broken {
    public void run( {
}
`)

	if len(parsed.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(parsed.Diagnostics))
	}

	diagnostic := parsed.Diagnostics[0]
	if diagnostic.Code != javaParseErrorCode {
		t.Fatalf("diagnostic code = %q, want %q", diagnostic.Code, javaParseErrorCode)
	}
	if diagnostic.Stage != models.StageParse {
		t.Fatalf("diagnostic stage = %q, want %q", diagnostic.Stage, models.StageParse)
	}
	if len(parsed.Symbols) != 0 {
		t.Fatalf("expected no symbols on parse error, got %d", len(parsed.Symbols))
	}
}

func TestJavaAdapterParseFilesWithProgressReportsPerFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	files := []models.ScannedSourceFile{
		writeJavaSource(t, dir, "src/com/example/First.java", `
package com.example;

public class First {
    public void run() {}
}
`),
		writeJavaSource(t, dir, "src/com/example/Second.java", `
package com.example;

public class Second {
    public void run() {}
}
`),
	}

	reported := []string{}
	parsedFiles, err := (JavaAdapter{}).ParseFilesWithProgress(files, dir, func(file models.ScannedSourceFile) {
		reported = append(reported, file.RelativePath)
	})
	if err != nil {
		t.Fatalf("ParseFilesWithProgress() error = %v", err)
	}

	if len(parsedFiles) != len(files) {
		t.Fatalf("expected %d parsed files, got %d", len(files), len(parsedFiles))
	}
	if len(reported) != len(files) {
		t.Fatalf("expected %d progress ticks, got %d", len(files), len(reported))
	}
}

func TestDiscoverJavaModuleHintsParsesGradleAndMavenSignals(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeRepositoryFile(t, root, "settings.gradle", strings.Join([]string{
		`rootProject.name = "atlas"`,
		`include("shared", "app")`,
		"",
	}, "\n"))
	writeRepositoryFile(t, root, "app/build.gradle", strings.Join([]string{
		"dependencies {",
		`    implementation(project(":shared"))`,
		"}",
		"",
	}, "\n"))
	writeRepositoryFile(t, root, "pom.xml", strings.Join([]string{
		"<project>",
		"  <modules>",
		"    <module>shared</module>",
		"    <module>app</module>",
		"  </modules>",
		"</project>",
		"",
	}, "\n"))
	writeRepositoryFile(t, root, "shared/pom.xml", strings.Join([]string{
		"<project>",
		"  <artifactId>shared</artifactId>",
		"</project>",
		"",
	}, "\n"))
	writeRepositoryFile(t, root, "app/pom.xml", strings.Join([]string{
		"<project>",
		"  <artifactId>app</artifactId>",
		"  <dependencies>",
		"    <dependency><artifactId>shared</artifactId></dependency>",
		"  </dependencies>",
		"</project>",
		"",
	}, "\n"))

	hints := discoverJavaModuleHints(root)
	if got := hints.moduleForFile("app/src/main/java/com/acme/app/Runner.java"); got != "app" {
		t.Fatalf("moduleForFile(app) = %q, want app", got)
	}
	if got := hints.moduleForFile("shared/src/main/java/com/acme/shared/SharedMath.java"); got != "shared" {
		t.Fatalf("moduleForFile(shared) = %q, want shared", got)
	}
	if _, ok := hints.moduleDependencies["app"]["shared"]; !ok {
		t.Fatalf("expected app module dependency to include shared, got %#v", hints.moduleDependencies)
	}
	if len(hints.warnings) != 0 {
		t.Fatalf("expected no metadata warnings, got %#v", hints.warnings)
	}
}

func TestParseMavenModulePomSignalsIgnoresParentArtifactID(t *testing.T) {
	t.Parallel()

	content := strings.Join([]string{
		"<project>",
		"  <parent>",
		"    <groupId>com.acme</groupId>",
		"    <artifactId>platform-parent</artifactId>",
		"    <version>1.0.0</version>",
		"  </parent>",
		"  <artifactId>billing-service</artifactId>",
		"  <dependencies>",
		"    <dependency><artifactId>shared-kernel</artifactId></dependency>",
		"  </dependencies>",
		"</project>",
	}, "\n")

	artifactID, dependencies, malformed := parseMavenModulePomSignals(content)
	if artifactID != "billing-service" {
		t.Fatalf("artifactID = %q, want billing-service", artifactID)
	}
	if len(dependencies) != 1 || dependencies[0] != "shared-kernel" {
		t.Fatalf("dependencies = %#v, want [shared-kernel]", dependencies)
	}
	if malformed {
		t.Fatal("expected valid pom metadata to avoid malformed flag")
	}
}

func TestJavaAdapterMissingMetadataKeepsResolutionStable(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"src/com/shared/Helper.java": `
package com.shared;

public class Helper {
    public static void assist() {}
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.shared.Helper;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	})

	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
	helperFile := mustFindParsedFile(t, parsedFiles, "src/com/shared/Helper.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	assistMethod := mustFindSymbol(t, helperFile.Symbols, "assist")

	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic call relation without module metadata")
	}

	for _, diagnostic := range runnerFile.Diagnostics {
		if diagnostic.Code == javaModuleHintWarningCode {
			t.Fatalf("did not expect module metadata warning without metadata files: %#v", runnerFile.Diagnostics)
		}
	}
}

func TestJavaAdapterMalformedMetadataEmitsWarningWithoutFailingParse(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSourcesWithRepositoryFiles(
		t,
		map[string]string{
			"settings.gradle": `include(`,
		},
		map[string]string{
			"app/src/main/java/com/acme/app/Runner.java": `
package com.acme.app;

public class Runner {
    public void run() {}
}
`,
		},
	)

	runnerFile := mustFindParsedFile(t, parsedFiles, "app/src/main/java/com/acme/app/Runner.java")
	moduleDiag := mustFindDiagnostic(t, runnerFile.Diagnostics, javaModuleHintWarningCode)
	if moduleDiag.Severity != models.SeverityWarning {
		t.Fatalf("module metadata diagnostic severity = %q, want %q", moduleDiag.Severity, models.SeverityWarning)
	}
	if !strings.Contains(moduleDiag.Detail, "include declarations malformed") {
		t.Fatalf("module metadata diagnostic detail = %q", moduleDiag.Detail)
	}
}
func TestJavaAdapterPrefersDeepResolutionForImportsAndCalls(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"src/com/shared/Helper.java": `
package com.shared;

public class Helper {
    public static void assist() {
    }
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.shared.Helper;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	})

	helperFile := mustFindParsedFile(t, parsedFiles, "src/com/shared/Helper.java")
	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
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
		t.Fatal("expected deep semantic reference relation from Runner.java to Helper class symbol")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected deep semantic call relation from run() to assist()")
	}
}

func TestJavaAdapterFallsBackToSyntacticResolutionWithDiagnostic(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"src/com/example/Helper.java": `
package com.example;

public class Helper {
    public static void assist() {
    }
}
`,
		"src/com/example/Runner.java": `
package com.example;

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
		t.Fatal("expected fallback syntactic call relation from run() to assist()")
	}
	if hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("did not expect semantic call relation for unresolved deep call target")
	}

	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if fallbackDiagnostic.Severity != models.SeverityWarning {
		t.Fatalf("fallback diagnostic severity = %q, want %q", fallbackDiagnostic.Severity, models.SeverityWarning)
	}
	if fallbackDiagnostic.Stage != models.StageParse {
		t.Fatalf("fallback diagnostic stage = %q, want %q", fallbackDiagnostic.Stage, models.StageParse)
	}
	if !strings.Contains(fallbackDiagnostic.Detail, "calls:assist") {
		t.Fatalf("fallback diagnostic detail = %q, want call target hint", fallbackDiagnostic.Detail)
	}
}

func TestJavaAdapterModelsNestedTypeOwnershipDeterministically(t *testing.T) {
	t.Parallel()

	parsed := parseSingleJavaFile(t, "src/com/example/Outer.java", `
package com.example;

public class Outer {
    public static class Inner {
        public void ping() {
        }
    }
}
`)

	outerClass := mustFindSymbol(t, parsed.Symbols, "Outer")
	innerClass := mustFindSymbol(t, parsed.Symbols, "Outer.Inner")
	innerMethod := mustFindSymbol(t, parsed.Symbols, "ping")

	if outerClass.SymbolKind != javaSymbolKindClass {
		t.Fatalf("outer class symbol kind = %q, want %q", outerClass.SymbolKind, javaSymbolKindClass)
	}
	if innerClass.SymbolKind != javaSymbolKindClass {
		t.Fatalf("inner class symbol kind = %q, want %q", innerClass.SymbolKind, javaSymbolKindClass)
	}
	if innerClass.Signature != "class Outer.Inner" {
		t.Fatalf("inner class signature = %q, want class Outer.Inner", innerClass.Signature)
	}
	if innerMethod.SymbolKind != javaSymbolKindMethod {
		t.Fatalf("inner method symbol kind = %q, want %q", innerMethod.SymbolKind, javaSymbolKindMethod)
	}
	if !hasRelation(parsed.Relations, parsed.File.ID, innerClass.ID, models.RelContains) {
		t.Fatal("expected file contains relation for nested class symbol")
	}
	if !hasRelation(parsed.Relations, parsed.File.ID, innerMethod.ID, models.RelContains) {
		t.Fatal("expected file contains relation for nested method symbol")
	}
}

func TestJavaAdapterResolvesOuterInnerReferences(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"src/com/shared/Outer.java": `
package com.shared;

public class Outer {
    public static class Inner {
        public static void assist() {
        }
    }
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.shared.Outer.Inner;

public class Runner {
    public void run() {
        Inner.assist();
    }
}
`,
	})

	outerFile := mustFindParsedFile(t, parsedFiles, "src/com/shared/Outer.java")
	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
	nestedClass := mustFindSymbol(t, outerFile.Symbols, "Outer.Inner")
	assistMethod := mustFindSymbol(t, outerFile.Symbols, "assist")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")

	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runnerFile.File.ID,
		nestedClass.ID,
		models.RelReferences,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic reference relation from Runner.java to Outer.Inner class symbol")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic call relation from run() to Outer.Inner.assist()")
	}
}

func TestJavaAdapterResolvesWildcardImportsForReferencesAndCalls(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
		"src/com/shared/Helper.java": `
package com.shared;

public class Helper {
    public static void assist() {
    }
}
`,
		"src/com/shared/Util.java": `
package com.shared;

public class Util {
    public static void noop() {
    }
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.shared.*;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	})

	helperFile := mustFindParsedFile(t, parsedFiles, "src/com/shared/Helper.java")
	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
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
		t.Fatal("expected semantic reference relation from wildcard import to Helper class symbol")
	}
	if !hasRelationWithConfidence(
		runnerFile.Relations,
		runMethod.ID,
		assistMethod.ID,
		models.RelCalls,
		models.ConfidenceSemantic,
	) {
		t.Fatal("expected semantic call relation via wildcard import resolution")
	}
}

func TestJavaAdapterEmitsFallbackDiagnosticForUnresolvedWildcardImport(t *testing.T) {
	t.Parallel()

	parsed := parseSingleJavaFile(t, "src/com/example/Runner.java", `
package com.example;

import com.missing.*;

public class Runner {
    public void run() {
    }
}
`)

	fallbackDiagnostic := mustFindDiagnostic(t, parsed.Diagnostics, javaResolutionFallbackCode)
	if !strings.Contains(fallbackDiagnostic.Detail, "references:com.missing.* (missing-wildcard-package)") {
		t.Fatalf("fallback diagnostic detail = %q, want unresolved wildcard import reason", fallbackDiagnostic.Detail)
	}
}

func TestJavaAdapterHandlesAmbiguousSimpleNameImportsDeterministically(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
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

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	})

	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
	alphaFile := mustFindParsedFile(t, parsedFiles, "src/com/alpha/Helper.java")
	betaFile := mustFindParsedFile(t, parsedFiles, "src/com/beta/Helper.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	alphaAssist := mustFindSymbol(t, alphaFile.Symbols, "assist")
	betaAssist := mustFindSymbol(t, betaFile.Symbols, "assist")

	if hasRelation(runnerFile.Relations, runMethod.ID, alphaAssist.ID, models.RelCalls) {
		t.Fatal("did not expect call relation to com.alpha.Helper.assist under ambiguous simple-name imports")
	}
	if hasRelation(runnerFile.Relations, runMethod.ID, betaAssist.ID, models.RelCalls) {
		t.Fatal("did not expect call relation to com.beta.Helper.assist under ambiguous simple-name imports")
	}

	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if !strings.Contains(fallbackDiagnostic.Detail, "calls:Helper.assist (ambiguous-import-class)") {
		t.Fatalf("fallback diagnostic detail = %q, want ambiguous-import-class reason", fallbackDiagnostic.Detail)
	}
}

func TestJavaAdapterTreatsStaticImportConflictsAsAmbiguous(t *testing.T) {
	t.Parallel()

	parsedFiles := parseJavaSources(t, map[string]string{
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

import static com.alpha.Helper.assist;
import static com.beta.Helper.assist;

public class Runner {
    public void run() {
        assist();
    }
}
`,
	})

	runnerFile := mustFindParsedFile(t, parsedFiles, "src/com/example/Runner.java")
	alphaFile := mustFindParsedFile(t, parsedFiles, "src/com/alpha/Helper.java")
	betaFile := mustFindParsedFile(t, parsedFiles, "src/com/beta/Helper.java")
	runMethod := mustFindSymbol(t, runnerFile.Symbols, "run")
	alphaAssist := mustFindSymbol(t, alphaFile.Symbols, "assist")
	betaAssist := mustFindSymbol(t, betaFile.Symbols, "assist")

	if hasRelation(runnerFile.Relations, runMethod.ID, alphaAssist.ID, models.RelCalls) {
		t.Fatal("did not expect call relation to com.alpha.Helper.assist under static import conflict")
	}
	if hasRelation(runnerFile.Relations, runMethod.ID, betaAssist.ID, models.RelCalls) {
		t.Fatal("did not expect call relation to com.beta.Helper.assist under static import conflict")
	}

	fallbackDiagnostic := mustFindDiagnostic(t, runnerFile.Diagnostics, javaResolutionFallbackCode)
	if !strings.Contains(fallbackDiagnostic.Detail, "calls:assist (ambiguous-static-call-target)") {
		t.Fatalf("fallback diagnostic detail = %q, want ambiguous-static-call-target reason", fallbackDiagnostic.Detail)
	}
}

func TestJavaAdapterNestedTypeOutputIsDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	sources := map[string]string{
		"src/com/shared/Outer.java": `
package com.shared;

public class Outer {
    public static class Inner {
        public static void assist() {
        }
    }
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.shared.Outer.Inner;

public class Runner {
    public void run() {
        Inner.assist();
    }
}
`,
	}

	first := parseJavaSources(t, sources)
	second := parseJavaSources(t, sources)

	if len(first) != len(second) {
		t.Fatalf("parse lengths differ: first=%d second=%d", len(first), len(second))
	}
	for idx := range first {
		if first[idx].File.FilePath != second[idx].File.FilePath {
			t.Fatalf("file order mismatch at %d: first=%q second=%q", idx, first[idx].File.FilePath, second[idx].File.FilePath)
		}
		if !reflect.DeepEqual(first[idx].Relations, second[idx].Relations) {
			t.Fatalf("relations differ for %s across repeated parse runs", first[idx].File.FilePath)
		}
	}
}

func TestJavaAdapterWildcardImportOutputIsDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()

	sources := map[string]string{
		"src/com/shared/Helper.java": `
package com.shared;

public class Helper {
    public static void assist() {
    }
}
`,
		"src/com/shared/Util.java": `
package com.shared;

public class Util {
    public static void noop() {
    }
}
`,
		"src/com/example/Runner.java": `
package com.example;

import com.shared.*;

public class Runner {
    public void run() {
        Helper.assist();
    }
}
`,
	}

	first := parseJavaSources(t, sources)
	second := parseJavaSources(t, sources)

	if len(first) != len(second) {
		t.Fatalf("parse lengths differ: first=%d second=%d", len(first), len(second))
	}
	for idx := range first {
		if first[idx].File.FilePath != second[idx].File.FilePath {
			t.Fatalf("file order mismatch at %d: first=%q second=%q", idx, first[idx].File.FilePath, second[idx].File.FilePath)
		}
		if !reflect.DeepEqual(first[idx].Relations, second[idx].Relations) {
			t.Fatalf("wildcard relations differ for %s across repeated parse runs", first[idx].File.FilePath)
		}
	}
}

func TestResolveJavaDeepImportBranches(t *testing.T) {
	t.Parallel()

	t.Run("wildcard import unresolved", func(t *testing.T) {
		t.Parallel()
		_, hint, reason, ok := resolveJavaDeepImport(
			"file:runner",
			javaImportRef{importPath: "com.example.*", isWildcard: true},
			map[string]string{},
			map[string][]string{},
			map[string]map[string][]string{},
		)
		if ok {
			t.Fatal("expected unresolved wildcard import")
		}
		if hint != "com.example.*" || reason != "missing-wildcard-package" {
			t.Fatalf("unexpected unresolved payload: hint=%q reason=%q", hint, reason)
		}
	})

	t.Run("wildcard import resolves package classes", func(t *testing.T) {
		t.Parallel()
		relations, _, _, ok := resolveJavaDeepImport(
			"file:runner",
			javaImportRef{importPath: "com.example.*", isWildcard: true},
			map[string]string{
				"com.example.Helper": "sym:helper",
				"com.example.Util":   "sym:util",
			},
			map[string][]string{
				"com.example": {"com.example.Util", "com.example.Helper"},
			},
			map[string]map[string][]string{},
		)
		if !ok {
			t.Fatal("expected wildcard import to resolve package classes")
		}
		if len(relations) != 2 {
			t.Fatalf("expected 2 wildcard relations, got %d", len(relations))
		}
		if relations[0].ToID != "sym:util" || relations[1].ToID != "sym:helper" {
			t.Fatalf("unexpected wildcard relation targets: %#v", relations)
		}
	})

	t.Run("semantic static import", func(t *testing.T) {
		t.Parallel()
		relation, _, _, ok := resolveJavaDeepImport(
			"file:runner",
			javaImportRef{importPath: "com.example.Helper.assist", isStatic: true},
			map[string]string{},
			map[string][]string{},
			map[string]map[string][]string{
				"com.example.Helper": {
					"assist": {"sym:assist"},
				},
			},
		)
		if !ok {
			t.Fatal("expected static import to resolve semantically")
		}
		if len(relation) != 1 || relation[0].ToID != "sym:assist" || relation[0].Confidence != models.ConfidenceSemantic {
			t.Fatalf("unexpected deep relation: %#v", relation)
		}
	})

	t.Run("missing class symbol unresolved", func(t *testing.T) {
		t.Parallel()
		_, hint, reason, ok := resolveJavaDeepImport(
			"file:runner",
			javaImportRef{importPath: "com.example.Missing"},
			map[string]string{},
			map[string][]string{},
			map[string]map[string][]string{},
		)
		if ok {
			t.Fatal("expected unresolved class import")
		}
		if hint != "com.example.Missing" || reason != "missing-class-symbol" {
			t.Fatalf("unexpected unresolved payload: hint=%q reason=%q", hint, reason)
		}
	})
}

func TestResolveJavaDeepCallTargetBranches(t *testing.T) {
	t.Parallel()

	methodIDsByClassFQN := map[string]map[string][]string{
		"com.example.Runner": {"run": {"sym:run"}},
		"com.example.Helper": {"assist": {"sym:assist"}},
		"com.shared.Helper":  {"assist": {"sym:assist_shared"}},
	}

	t.Run("unqualified owner method", func(t *testing.T) {
		t.Parallel()
		targetID, reason := resolveJavaDeepCallTarget(
			javaCallTarget{methodName: "run"},
			"com.example.Runner",
			"com.example",
			map[string]string{},
			map[string]string{},
			map[string][]string{},
			map[string][]string{},
			map[string][]string{},
			methodIDsByClassFQN,
		)
		if targetID != "sym:run" || reason != "" {
			t.Fatalf("unexpected deep target resolution: target=%q reason=%q", targetID, reason)
		}
	})

	t.Run("unqualified requires metadata", func(t *testing.T) {
		t.Parallel()
		targetID, reason := resolveJavaDeepCallTarget(
			javaCallTarget{methodName: "assist"},
			"",
			"com.example",
			map[string]string{},
			map[string]string{},
			map[string][]string{},
			map[string][]string{},
			map[string][]string{},
			methodIDsByClassFQN,
		)
		if targetID != "" || reason != "missing-owner-class" {
			t.Fatalf("unexpected deep target resolution: target=%q reason=%q", targetID, reason)
		}
	})

	t.Run("qualified ambiguous target", func(t *testing.T) {
		t.Parallel()
		targetID, reason := resolveJavaDeepCallTarget(
			javaCallTarget{methodName: "assist", qualifier: "Helper"},
			"com.example.Runner",
			"com.example",
			map[string]string{"Helper": "com.shared.Helper"},
			map[string]string{"Helper": "com.example.Helper"},
			map[string][]string{},
			map[string][]string{},
			map[string][]string{},
			methodIDsByClassFQN,
		)
		if targetID != "" || reason != "ambiguous-qualified-target" {
			t.Fatalf("unexpected deep target resolution: target=%q reason=%q", targetID, reason)
		}
	})
}

func TestJavaNestedResolutionHelperBranches(t *testing.T) {
	t.Parallel()

	if got := javaTypeQualifierFromFQN("com.shared.Outer.Inner"); got != "Outer.Inner" {
		t.Fatalf("javaTypeQualifierFromFQN() = %q, want Outer.Inner", got)
	}
	if got := javaTypeQualifierFromFQN("com.shared.helper"); got != "" {
		t.Fatalf("javaTypeQualifierFromFQN() lowercase type = %q, want empty", got)
	}

	head, tail := javaSplitQualifiedHead("Outer.Inner")
	if head != "Outer" || tail != "Inner" {
		t.Fatalf("javaSplitQualifiedHead() = (%q, %q), want (Outer, Inner)", head, tail)
	}
	head, tail = javaSplitQualifiedHead("Single")
	if head != "" || tail != "" {
		t.Fatalf("javaSplitQualifiedHead() without dot = (%q, %q), want empty", head, tail)
	}

	classCandidates := resolveJavaClassCandidates(
		"Outer.Inner",
		"com.shared",
		map[string]string{
			"Outer":       "com.shared.Outer",
			"Outer.Inner": "com.shared.Outer.Inner",
		},
		map[string]string{
			"Outer":       "com.shared.Outer",
			"Outer.Inner": "com.shared.Outer.Inner",
		},
		map[string][]string{},
		map[string][]string{},
		map[string]map[string][]string{
			"com.shared.Outer.Inner": {"assist": {"sym:assist"}},
		},
	)
	if len(classCandidates) != 1 || classCandidates[0] != "com.shared.Outer.Inner" {
		t.Fatalf("resolveJavaClassCandidates() = %#v, want [com.shared.Outer.Inner]", classCandidates)
	}

	if got := normalizeJavaQualifier(" this . Outer . Inner "); got != "Outer.Inner" {
		t.Fatalf("normalizeJavaQualifier() = %q, want Outer.Inner", got)
	}
	if got := normalizeJavaQualifier("super"); got != "" {
		t.Fatalf("normalizeJavaQualifier(super) = %q, want empty", got)
	}
}

func TestResolveJavaMethodInvocationFallbackParsing(t *testing.T) {
	t.Parallel()

	parser, err := newParser(javaLanguage())
	if err != nil {
		t.Fatalf("newParser() error = %v", err)
	}
	defer parser.Close()

	source := []byte(`
package com.example;

public class Runner {
    public void run() {
        com.example.Helper.assist();
    }
}
`)
	tree := parser.Parse(source, nil)
	if tree == nil {
		t.Fatal("expected syntax tree")
	}
	defer tree.Close()

	root := tree.RootNode()
	methodInvocations := collectNodesByKind(root, "method_invocation")
	if len(methodInvocations) != 1 {
		t.Fatalf("expected 1 method invocation, got %d", len(methodInvocations))
	}
	methodInvocation := methodInvocations[0]

	methodName := resolveJavaMethodInvocationName(&methodInvocation, source)
	qualifier := resolveJavaMethodInvocationQualifier(&methodInvocation, source)
	if methodName != "assist" {
		t.Fatalf("method name = %q, want assist", methodName)
	}
	if qualifier != "com.example.Helper" {
		t.Fatalf("qualifier = %q, want com.example.Helper", qualifier)
	}
}

func TestSortJavaDiagnosticsDeterministicOrder(t *testing.T) {
	t.Parallel()

	diagnostics := []models.StructuredDiagnostic{
		{Code: javaParseErrorCode, Severity: models.SeverityError, FilePath: "b.java", Detail: "z"},
		{Code: javaResolutionFallbackCode, Severity: models.SeverityWarning, FilePath: "a.java", Detail: "b"},
		{Code: javaResolutionFallbackCode, Severity: models.SeverityWarning, FilePath: "a.java", Detail: "a"},
	}

	sortJavaDiagnostics(diagnostics)
	if diagnostics[0].Code != javaParseErrorCode {
		t.Fatalf("first diagnostic code = %q, want %q", diagnostics[0].Code, javaParseErrorCode)
	}
	if diagnostics[1].Detail != "a" || diagnostics[2].Detail != "b" {
		t.Fatalf("expected fallback diagnostics sorted by detail, got %#v", diagnostics)
	}
}

func TestCreateJavaResolutionFallbackDiagnosticTruncatesHighVolumeDetailsDeterministically(t *testing.T) {
	t.Parallel()

	unresolved := make([]javaUnresolvedRef, 0, javaFallbackDiagnosticMaxEntries+25)
	for index := 0; index < javaFallbackDiagnosticMaxEntries+25; index++ {
		unresolved = append(unresolved, javaUnresolvedRef{
			relationType: models.RelCalls,
			targetHint:   fmt.Sprintf("target-%03d", index),
			reason:       "missing-qualified-method",
		})
	}

	diagnostic := createJavaResolutionFallbackDiagnostic(models.GraphFile{
		FilePath: "src/com/example/Runner.java",
		Language: models.LangJava,
	}, unresolved)
	if diagnostic.Code != javaResolutionFallbackCode {
		t.Fatalf("diagnostic code = %q, want %q", diagnostic.Code, javaResolutionFallbackCode)
	}
	if !strings.Contains(diagnostic.Detail, "calls:target-000 (missing-qualified-method)") {
		t.Fatalf("expected deterministic first fallback segment, got %q", diagnostic.Detail)
	}
	if strings.Contains(diagnostic.Detail, "calls:target-224 (missing-qualified-method)") {
		t.Fatalf("expected capped fallback diagnostic detail, got %q", diagnostic.Detail)
	}

	truncationSegment := fmt.Sprintf("%s (25 entries omitted)", javaDiagnosticTruncationPrefixKey)
	if !strings.Contains(diagnostic.Detail, truncationSegment) {
		t.Fatalf("expected truncation metadata %q, got %q", truncationSegment, diagnostic.Detail)
	}
}

func TestCreateJavaModuleHintDiagnosticTruncatesWarningPayloadDeterministically(t *testing.T) {
	t.Parallel()

	warnings := make([]string, 0, javaModuleHintWarningMaxEntries+7)
	for index := 0; index < javaModuleHintWarningMaxEntries+7; index++ {
		warnings = append(warnings, fmt.Sprintf("warning-%03d", index))
	}

	diagnostic := createJavaModuleHintDiagnostic(models.GraphFile{
		FilePath: "src/com/example/Runner.java",
		Language: models.LangJava,
	}, warnings)
	if diagnostic.Code != javaModuleHintWarningCode {
		t.Fatalf("diagnostic code = %q, want %q", diagnostic.Code, javaModuleHintWarningCode)
	}
	if !strings.Contains(diagnostic.Detail, "warning-000") {
		t.Fatalf("expected warning payload to include first warning, got %q", diagnostic.Detail)
	}
	if strings.Contains(diagnostic.Detail, fmt.Sprintf("warning-%03d", javaModuleHintWarningMaxEntries+6)) {
		t.Fatalf("expected warning payload to truncate overflow warnings, got %q", diagnostic.Detail)
	}

	truncationSegment := fmt.Sprintf("%s (7 warnings omitted)", javaDiagnosticTruncationPrefixKey)
	if !strings.Contains(diagnostic.Detail, truncationSegment) {
		t.Fatalf("expected warning truncation metadata %q, got %q", truncationSegment, diagnostic.Detail)
	}
}

func TestJavaHelperResolutionUtilities(t *testing.T) {
	t.Parallel()

	if got := javaQualifiedName("com.example", "Runner"); got != "com.example.Runner" {
		t.Fatalf("javaQualifiedName() = %q, want com.example.Runner", got)
	}
	if got := javaQualifiedName("", "Runner"); got != "Runner" {
		t.Fatalf("javaQualifiedName() with empty package = %q, want Runner", got)
	}

	targetIDs := resolveJavaCallTarget(
		javaCallTarget{methodName: "assist", qualifier: "Helper"},
		"com.example",
		map[string]string{"Helper": "com.example.Helper"},
		map[string]string{},
		map[string][]string{},
		map[string][]string{},
		map[string][]string{},
		map[string]map[string][]string{
			"com.example.Helper": {"assist": {"sym:assist"}},
		},
		map[string]map[string][]string{},
	)
	if len(targetIDs) != 1 || targetIDs[0] != "sym:assist" {
		t.Fatalf("resolveJavaCallTarget() = %#v, want [sym:assist]", targetIDs)
	}
}

func TestJavaSymbolNameSignatureAndComplexityHelpers(t *testing.T) {
	t.Parallel()

	parser, err := newParser(javaLanguage())
	if err != nil {
		t.Fatalf("newParser() error = %v", err)
	}
	defer parser.Close()

	source := []byte(`
package com.example;

public class Runner {
    public void run() {
        if (true && true) {
            Helper.assist();
        }
    }
}
`)
	tree := parser.Parse(source, nil)
	if tree == nil {
		t.Fatal("expected syntax tree")
	}
	defer tree.Close()

	root := tree.RootNode()
	classNodes := collectNodesByKind(root, "class_declaration")
	if len(classNodes) != 1 {
		t.Fatalf("expected one class declaration, got %d", len(classNodes))
	}
	classNode := classNodes[0]

	methodNodes := collectNodesByKind(root, "method_declaration")
	if len(methodNodes) != 1 {
		t.Fatalf("expected one method declaration, got %d", len(methodNodes))
	}
	methodNode := methodNodes[0]

	if got := resolveJavaSymbolName(&classNode, source, javaSymbolKindClass); got != "Runner" {
		t.Fatalf("resolveJavaSymbolName(class) = %q, want Runner", got)
	}
	if got := formatJavaSignature(&classNode, source, javaSymbolKindClass, "Runner"); got != "class Runner" {
		t.Fatalf("formatJavaSignature(class) = %q, want class Runner", got)
	}
	if got := formatJavaSignature(&methodNode, source, javaSymbolKindMethod, "run"); !strings.Contains(got, "void run(") {
		t.Fatalf("formatJavaSignature(method) = %q, want method signature", got)
	}
	if got := computeJavaCyclomaticComplexity(&methodNode, source, javaSymbolKindMethod); got < 3 {
		t.Fatalf("computeJavaCyclomaticComplexity() = %d, want >= 3", got)
	}
}

func TestResolveJavaCallTargetBranches(t *testing.T) {
	t.Parallel()

	methodIDsByClassFQN := map[string]map[string][]string{
		"com.example.Helper":      {"assist": {"sym:assist"}},
		"com.example.Outer.Inner": {"assist": {"sym:nested-assist"}},
	}

	if ids := resolveJavaCallTarget(
		javaCallTarget{},
		"com.example",
		map[string]string{},
		map[string]string{},
		map[string][]string{},
		map[string][]string{},
		map[string][]string{},
		methodIDsByClassFQN,
		map[string]map[string][]string{},
	); ids != nil {
		t.Fatalf("resolveJavaCallTarget() for empty method should return nil, got %#v", ids)
	}

	if ids := resolveJavaCallTarget(
		javaCallTarget{methodName: "assist"},
		"com.example",
		map[string]string{},
		map[string]string{},
		map[string][]string{},
		map[string][]string{"assist": {"sym:assist"}},
		map[string][]string{},
		methodIDsByClassFQN,
		map[string]map[string][]string{
			"com.example": {"assist": {"sym:package-assist"}},
		},
	); len(ids) != 1 || ids[0] != "sym:assist" {
		t.Fatalf("resolveJavaCallTarget() static import branch = %#v, want [sym:assist]", ids)
	}

	if ids := resolveJavaCallTarget(
		javaCallTarget{methodName: "assist", qualifier: "Outer.Inner"},
		"com.example",
		map[string]string{"Outer": "com.example.Outer"},
		map[string]string{},
		map[string][]string{},
		map[string][]string{},
		map[string][]string{},
		methodIDsByClassFQN,
		map[string]map[string][]string{},
	); len(ids) != 1 || ids[0] != "sym:nested-assist" {
		t.Fatalf("resolveJavaCallTarget() nested qualifier branch = %#v, want [sym:nested-assist]", ids)
	}
}

func TestResolveJavaSyntacticImportBranches(t *testing.T) {
	t.Parallel()

	if _, ok := resolveJavaSyntacticImport(
		"file:runner",
		javaImportRef{importPath: "com.example.*", isWildcard: true},
		map[string]string{"com.example.Helper": "sym:helper"},
	); ok {
		t.Fatal("expected wildcard syntactic import to be unresolved")
	}

	relation, ok := resolveJavaSyntacticImport(
		"file:runner",
		javaImportRef{importPath: "com.example.Helper"},
		map[string]string{"com.example.Helper": "sym:helper"},
	)
	if !ok {
		t.Fatal("expected class syntactic import to resolve")
	}
	if relation.FromID != "file:runner" || relation.ToID != "sym:helper" {
		t.Fatalf("unexpected syntactic import relation: %#v", relation)
	}
}

func parseSingleJavaFile(t *testing.T, relativePath string, source string) models.ParsedFile {
	t.Helper()

	parsedFiles := parseJavaSources(t, map[string]string{relativePath: source})
	if len(parsedFiles) != 1 {
		t.Fatalf("expected 1 parsed file, got %d", len(parsedFiles))
	}

	return parsedFiles[0]
}

func parseJavaSources(t *testing.T, sources map[string]string) []models.ParsedFile {
	t.Helper()

	dir := t.TempDir()
	paths := make([]string, 0, len(sources))
	for relativePath := range sources {
		paths = append(paths, relativePath)
	}
	sort.Strings(paths)

	files := make([]models.ScannedSourceFile, 0, len(paths))
	for _, relativePath := range paths {
		files = append(files, writeJavaSource(t, dir, relativePath, sources[relativePath]))
	}

	parsedFiles, err := (JavaAdapter{}).ParseFiles(files, dir)
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	return parsedFiles
}

func parseJavaSourcesWithRepositoryFiles(
	t *testing.T,
	repositoryFiles map[string]string,
	javaSources map[string]string,
) []models.ParsedFile {
	t.Helper()

	dir := t.TempDir()
	for relativePath, content := range repositoryFiles {
		writeRepositoryFile(t, dir, relativePath, content)
	}

	paths := make([]string, 0, len(javaSources))
	for relativePath := range javaSources {
		paths = append(paths, relativePath)
	}
	sort.Strings(paths)

	files := make([]models.ScannedSourceFile, 0, len(paths))
	for _, relativePath := range paths {
		files = append(files, writeJavaSource(t, dir, relativePath, javaSources[relativePath]))
	}

	parsedFiles, err := (JavaAdapter{}).ParseFiles(files, dir)
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}
	return parsedFiles
}

func writeJavaSource(t *testing.T, root string, relativePath string, source string) models.ScannedSourceFile {
	t.Helper()

	absolutePath := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", relativePath, err)
	}

	content := strings.TrimLeft(source, "\n")
	if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}

	return models.ScannedSourceFile{
		AbsolutePath: absolutePath,
		RelativePath: relativePath,
		Language:     models.LangJava,
	}
}

func writeRepositoryFile(t *testing.T, root string, relativePath string, content string) {
	t.Helper()

	absolutePath := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", relativePath, err)
	}
	if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func hasRelationWithConfidence(
	relations []models.RelationEdge,
	fromID string,
	toID string,
	relationType models.RelationType,
	confidence models.RelationConfidence,
) bool {
	for _, relation := range relations {
		if relation.FromID == fromID &&
			relation.ToID == toID &&
			relation.Type == relationType &&
			relation.Confidence == confidence {
			return true
		}
	}

	return false
}

func mustFindDiagnostic(
	t *testing.T,
	diagnostics []models.StructuredDiagnostic,
	code string,
) models.StructuredDiagnostic {
	t.Helper()

	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return diagnostic
		}
	}

	t.Fatalf("missing diagnostic %q", code)
	return models.StructuredDiagnostic{}
}
