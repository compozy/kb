//go:build integration

package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/vault"
)

func TestGenerateIntegrationBuildsVaultFromFixtureRepository(t *testing.T) {
	t.Parallel()

	fixtureRoot := filepath.Join("testdata", "fixture-go-repo")
	outputRoot := filepath.Join(t.TempDir(), "vault")
	generator := newRunner()

	summary, err := generator.Generate(context.Background(), models.GenerateOptions{
		RootPath:  fixtureRoot,
		VaultPath: outputRoot,
		TopicSlug: "fixture-go-repo",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 2 {
		t.Fatalf("FilesScanned = %d, want 2", summary.FilesScanned)
	}
	if summary.FilesParsed != 2 {
		t.Fatalf("FilesParsed = %d, want 2", summary.FilesParsed)
	}
	if summary.FilesSkipped != 0 {
		t.Fatalf("FilesSkipped = %d, want 0", summary.FilesSkipped)
	}
	if summary.SymbolsExtracted != 4 {
		t.Fatalf("SymbolsExtracted = %d, want 4", summary.SymbolsExtracted)
	}
	if summary.RawDocumentsWritten != 9 {
		t.Fatalf("RawDocumentsWritten = %d, want 9", summary.RawDocumentsWritten)
	}
	if summary.WikiDocumentsWritten != 10 {
		t.Fatalf("WikiDocumentsWritten = %d, want 10", summary.WikiDocumentsWritten)
	}
	if summary.IndexDocumentsWritten != 3 {
		t.Fatalf("IndexDocumentsWritten = %d, want 3", summary.IndexDocumentsWritten)
	}

	expectedPaths := []string{
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "main.go.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "internal", "greeter", "greeter.go.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "symbols", "hello--internal-greeter-greeter-go-l6.md"),
		filepath.Join(summary.TopicPath, filepath.FromSlash(vault.GetWikiConceptPath("Codebase Overview"))),
		filepath.Join(summary.TopicPath, filepath.FromSlash(vault.GetWikiIndexPath(vault.CodebaseDashboardTitle))),
		filepath.Join(summary.TopicPath, "bases", "module-health.base"),
		filepath.Join(summary.TopicPath, "CLAUDE.md"),
	}

	for _, expectedPath := range expectedPaths {
		if _, err := os.Stat(expectedPath); err != nil {
			t.Fatalf("expected generated path %s: %v", expectedPath, err)
		}
	}

	logContent, err := os.ReadFile(filepath.Join(summary.TopicPath, "log.md"))
	if err != nil {
		t.Fatalf("read log.md: %v", err)
	}
	if got := string(logContent); !containsAll(got,
		"## [",
		"ingest | codebase (2 files, 4 symbols)",
	) {
		t.Fatalf("expected codebase ingest log entry, got:\n%s", got)
	}
}

func TestGenerateIntegrationScansRepositoryNestedInsideVaultButOutsideTopic(t *testing.T) {
	t.Parallel()

	vaultRoot := t.TempDir()
	repoRoot := filepath.Join(vaultRoot, ".resources", "nested-go-repo")
	if err := os.MkdirAll(filepath.Join(repoRoot, "internal", "greeter"), 0o755); err != nil {
		t.Fatalf("create nested repo: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "go.mod"), []byte("module example.com/nested-go-repo\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "main.go"), []byte(strings.Join([]string{
		"package main",
		"",
		"import \"example.com/nested-go-repo/internal/greeter\"",
		"",
		"func main() {",
		"\tgreeter.Hello()",
		"}",
		"",
	}, "\n")), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "internal", "greeter", "greeter.go"), []byte(strings.Join([]string{
		"package greeter",
		"",
		"func Hello() string {",
		"\treturn \"hello\"",
		"}",
		"",
	}, "\n")), 0o644); err != nil {
		t.Fatalf("write greeter.go: %v", err)
	}

	summary, err := newRunner().Generate(context.Background(), models.GenerateOptions{
		RootPath:  repoRoot,
		VaultPath: vaultRoot,
		TopicSlug: "nested-go-repo",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 2 {
		t.Fatalf("FilesScanned = %d, want 2", summary.FilesScanned)
	}
	if summary.FilesParsed != 2 {
		t.Fatalf("FilesParsed = %d, want 2", summary.FilesParsed)
	}
	if _, err := os.Stat(filepath.Join(vaultRoot, "nested-go-repo", "raw", "codebase", "files", "main.go.md")); err != nil {
		t.Fatalf("expected nested-vault ingest output: %v", err)
	}
}

func TestGenerateIntegrationBuildsVaultFromRustWorkspace(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	writeRustWorkspaceFixture(t, repoRoot)

	outputRoot := filepath.Join(t.TempDir(), "vault")
	summary, err := newRunner().Generate(context.Background(), models.GenerateOptions{
		RootPath:  repoRoot,
		VaultPath: outputRoot,
		TopicSlug: "fixture-rust-workspace",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 3 {
		t.Fatalf("FilesScanned = %d, want 3", summary.FilesScanned)
	}
	if summary.FilesParsed != 3 {
		t.Fatalf("FilesParsed = %d, want 3", summary.FilesParsed)
	}
	if !containsAll(strings.Join(summary.DetectedLanguages, ","), "rust") {
		t.Fatalf("expected rust in detected languages, got %#v", summary.DetectedLanguages)
	}
	if summary.SymbolsExtracted == 0 {
		t.Fatalf("SymbolsExtracted = %d, want > 0", summary.SymbolsExtracted)
	}

	expectedPaths := []string{
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "crates", "core", "src", "lib.rs.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "crates", "core", "src", "util.rs.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "crates", "app", "src", "lib.rs.md"),
	}
	for _, expectedPath := range expectedPaths {
		if _, err := os.Stat(expectedPath); err != nil {
			t.Fatalf("expected generated path %s: %v", expectedPath, err)
		}
	}
}

func TestGenerateIntegrationDryRunSelectsJavaAdapterForMixedWorkspace(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	writeFixtureFile(t, repoRoot, "go.mod", "module example.com/mixed-repo\n\ngo 1.22\n")
	writeFixtureFile(t, repoRoot, "main.go", strings.Join([]string{
		"package main",
		"",
		"func main() {}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "src/App.java", strings.Join([]string{
		"package src;",
		"",
		"public class App {",
		"    public static void main(String[] args) {",
		"    }",
		"}",
		"",
	}, "\n"))

	summary, err := newRunner().Generate(context.Background(), models.GenerateOptions{
		RootPath: repoRoot,
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if !summary.DryRun {
		t.Fatalf("DryRun = %t, want true", summary.DryRun)
	}
	if got, want := summary.DetectedLanguages, []string{"go", "java"}; !strings.EqualFold(strings.Join(got, ","), strings.Join(want, ",")) {
		t.Fatalf("DetectedLanguages = %#v, want %#v", got, want)
	}
	if got, want := summary.SelectedAdapters, []string{"adapter.GoAdapter", "adapter.JavaAdapter"}; !strings.EqualFold(strings.Join(got, ","), strings.Join(want, ",")) {
		t.Fatalf("SelectedAdapters = %#v, want %#v", got, want)
	}
}

func TestGenerateIntegrationBuildsVaultFromJavaPhase2Workspace(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	writeJavaPhase2Fixture(t, repoRoot)

	outputRoot := filepath.Join(t.TempDir(), "vault")
	summary, err := newRunner().Generate(context.Background(), models.GenerateOptions{
		RootPath:  repoRoot,
		VaultPath: outputRoot,
		TopicSlug: "fixture-java-phase2",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if summary.FilesScanned != 6 {
		t.Fatalf("FilesScanned = %d, want 6", summary.FilesScanned)
	}
	if summary.FilesParsed != 6 {
		t.Fatalf("FilesParsed = %d, want 6", summary.FilesParsed)
	}
	if !containsAll(strings.Join(summary.DetectedLanguages, ","), "java") {
		t.Fatalf("expected java in detected languages, got %#v", summary.DetectedLanguages)
	}
	if got, want := strings.Join(summary.SelectedAdapters, ","), "adapter.JavaAdapter"; !strings.EqualFold(got, want) {
		t.Fatalf("SelectedAdapters = %#v, want [%s]", summary.SelectedAdapters, want)
	}

	for _, expectedPath := range []string{
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "shared-a", "src", "main", "java", "com", "acme", "shareda", "Helper.java.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "shared-b", "src", "main", "java", "com", "acme", "sharedb", "Helper.java.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "shared-b", "src", "main", "java", "com", "acme", "sharedb", "Outer.java.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "shared-b", "src", "main", "java", "com", "acme", "sharedb", "Tooling.java.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "app", "src", "main", "java", "com", "acme", "app", "Runner.java.md"),
		filepath.Join(summary.TopicPath, "raw", "codebase", "files", "app", "src", "main", "java", "com", "acme", "app", "AppMain.java.md"),
	} {
		if _, err := os.Stat(expectedPath); err != nil {
			t.Fatalf("expected generated path %s: %v", expectedPath, err)
		}
	}
}

func TestGenerateIntegrationJavaIngestPerformanceBudget(t *testing.T) {
	policy := canonicalJavaBenchmarkPolicy()
	enforcePerformanceBudget := os.Getenv("ENABLE_PERF_BUDGET") == "1"

	goRepoRoot := t.TempDir()
	writeGoBaselineFixture(t, goRepoRoot)

	baselineDuration := measureMedianGenerateDuration(t, benchmarkGenerateOptions(goRepoRoot), policy.RepeatCount)

	if baselineDuration <= 0 {
		t.Fatalf("baseline duration must be positive, got %s", baselineDuration)
	}

	for _, fixture := range canonicalJavaBenchmarkFixtures() {
		javaRepoRoot := t.TempDir()
		writeJavaBenchmarkFixtureByProfile(t, javaRepoRoot, fixture.Profile)

		javaDuration := measureMedianGenerateDuration(t, benchmarkGenerateOptions(javaRepoRoot), policy.RepeatCount)
		overhead := float64(javaDuration) / float64(baselineDuration)
		t.Logf(
			"java ingest budget sample: profile=%s baseline=%s java=%s overhead=%.2f%% budget=%.2f%% samples=%d",
			fixture.Profile,
			baselineDuration,
			javaDuration,
			(overhead-1)*100,
			(policy.OverheadBudget-1)*100,
			policy.RepeatCount,
		)
		if overhead > policy.OverheadBudget && enforcePerformanceBudget {
			t.Fatalf(
				"profile %s java ingest overhead %.2f%% exceeds budget %.2f%% (baseline=%s java=%s)",
				fixture.Profile,
				(overhead-1)*100,
				(policy.OverheadBudget-1)*100,
				baselineDuration,
				javaDuration,
			)
		}
		if overhead > policy.OverheadBudget && !enforcePerformanceBudget {
			t.Logf(
				"java ingest overhead exceeded budget but enforcement is disabled (set ENABLE_PERF_BUDGET=1 to enforce): profile=%s overhead=%.2f%% budget=%.2f%%",
				fixture.Profile,
				(overhead-1)*100,
				(policy.OverheadBudget-1)*100,
			)
		}
	}
}

func TestGenerateIntegrationJavaHighFallbackVolumeCapsDiagnosticPayload(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	writeFixtureFile(t, repoRoot, "src/main/java/com/acme/shared/Helper.java", strings.Join([]string{
		"package com.acme.shared;",
		"",
		"public class Helper {",
		"    public static void assist() {}",
		"}",
		"",
	}, "\n"))

	var runnerBody strings.Builder
	runnerBody.WriteString("package com.acme.app;\n\n")
	runnerBody.WriteString("public class Runner {\n")
	runnerBody.WriteString("    public void run() {\n")
	for index := 0; index < 320; index++ {
		runnerBody.WriteString(fmt.Sprintf("        Missing%03d.assist();\n", index))
	}
	runnerBody.WriteString("    }\n")
	runnerBody.WriteString("}\n")
	writeFixtureFile(t, repoRoot, "src/main/java/com/acme/app/Runner.java", runnerBody.String())

	summary, err := newRunner().Generate(context.Background(), models.GenerateOptions{
		RootPath: repoRoot,
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	var fallbackDiagnostic models.StructuredDiagnostic
	foundFallback := false
	for _, diagnostic := range summary.Diagnostics {
		if diagnostic.Code != "JAVA_RESOLUTION_FALLBACK" {
			continue
		}
		fallbackDiagnostic = diagnostic
		foundFallback = true
		break
	}
	if !foundFallback {
		t.Fatalf("expected JAVA_RESOLUTION_FALLBACK diagnostic in summary diagnostics: %#v", summary.Diagnostics)
	}
	if !strings.Contains(fallbackDiagnostic.Detail, "meta:truncated") {
		t.Fatalf("expected capped fallback diagnostic detail with truncation metadata, got %q", fallbackDiagnostic.Detail)
	}
}

func BenchmarkGenerateIntegrationGoBaselineDryRun(b *testing.B) {
	repoRoot := b.TempDir()
	writeGoBaselineFixture(b, repoRoot)

	ctx := context.Background()
	generator := newRunner()
	options := benchmarkGenerateOptions(repoRoot)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		if _, err := generator.Generate(ctx, options); err != nil {
			b.Fatalf("Generate baseline dry-run failed: %v", err)
		}
	}
}

func BenchmarkGenerateIntegrationJavaCanonicalDryRun(b *testing.B) {
	for _, fixture := range canonicalJavaBenchmarkFixtures() {
		fixture := fixture
		b.Run(string(fixture.Profile), func(b *testing.B) {
			repoRoot := b.TempDir()
			writeJavaBenchmarkFixtureByProfile(b, repoRoot, fixture.Profile)

			ctx := context.Background()
			generator := newRunner()
			options := benchmarkGenerateOptions(repoRoot)

			b.ReportAllocs()
			b.ResetTimer()

			for idx := 0; idx < b.N; idx++ {
				if _, err := generator.Generate(ctx, options); err != nil {
					b.Fatalf("Generate java dry-run failed for profile %s: %v", fixture.Profile, err)
				}
			}
		})
	}
}

func measureMedianGenerateDuration(t testing.TB, options models.GenerateOptions, samples int) time.Duration {
	t.Helper()

	if samples <= 0 {
		t.Fatalf("samples must be > 0, got %d", samples)
	}

	durations := make([]time.Duration, 0, samples)
	ctx := context.Background()
	repoRoot := options.RootPath

	// Warm cache to reduce first-run penalty noise before budget comparison.
	if _, err := newRunner().Generate(ctx, options); err != nil {
		t.Fatalf("warm-up dry-run failed for %q: %v", repoRoot, err)
	}

	for sample := 0; sample < samples; sample++ {
		startedAt := time.Now()
		if _, err := newRunner().Generate(ctx, options); err != nil {
			t.Fatalf("dry-run sample %d failed for %q: %v", sample+1, repoRoot, err)
		}
		durations = append(durations, time.Since(startedAt))
	}

	median, err := medianDurationFromSamples(durations)
	if err != nil {
		t.Fatalf("measureMedianGenerateDuration median extraction failed for %q: %v", repoRoot, err)
	}

	return median
}

func writeGoBaselineFixture(t testing.TB, repoRoot string) {
	t.Helper()

	writeFixtureFile(t, repoRoot, "go.mod", "module example.com/bench-go\n\ngo 1.22\n")
	for idx := 0; idx < 24; idx++ {
		writeFixtureFile(t, repoRoot, fmt.Sprintf("pkg/service_%02d/service_%02d.go", idx, idx), strings.Join([]string{
			fmt.Sprintf("package service_%02d", idx),
			"",
			fmt.Sprintf("func Compute%02d(input int) int {", idx),
			"\tvalue := input + 1",
			"\treturn value * 2",
			"}",
			"",
		}, "\n"))
	}
}

func writeJavaBenchmarkFixtureByProfile(
	t testing.TB,
	repoRoot string,
	profile javaBenchmarkProfile,
) {
	t.Helper()

	switch profile {
	case javaBenchmarkProfileSingleModuleLibrary:
		writeJavaSingleModuleLibraryBenchmarkFixture(t, repoRoot)
	case javaBenchmarkProfileSpringService:
		writeJavaSpringServiceBenchmarkFixture(t, repoRoot)
	case javaBenchmarkProfileMultiModuleEnterprise:
		writeJavaMultiModuleEnterpriseBenchmarkFixture(t, repoRoot)
	default:
		t.Fatalf("unsupported java benchmark profile %q", profile)
	}
}

func writeJavaSingleModuleLibraryBenchmarkFixture(t testing.TB, repoRoot string) {
	t.Helper()

	for idx := 0; idx < 18; idx++ {
		writeFixtureFile(t, repoRoot, fmt.Sprintf("src/main/java/com/acme/library/LibraryMath%02d.java", idx), strings.Join([]string{
			"package com.acme.library;",
			"",
			fmt.Sprintf("public class LibraryMath%02d {", idx),
			"    public int scale(int input) {",
			"        return input * 2;",
			"    }",
			"}",
			"",
		}, "\n"))
	}
}

func writeJavaSpringServiceBenchmarkFixture(t testing.TB, repoRoot string) {
	t.Helper()

	for idx := 0; idx < 8; idx++ {
		writeFixtureFile(t, repoRoot, fmt.Sprintf("src/main/java/com/acme/service/repository/OrderRepository%02d.java", idx), strings.Join([]string{
			"package com.acme.service.repository;",
			"",
			fmt.Sprintf("public class OrderRepository%02d {", idx),
			"    public int findTotal() {",
			"        return 21;",
			"    }",
			"}",
			"",
		}, "\n"))
		writeFixtureFile(t, repoRoot, fmt.Sprintf("src/main/java/com/acme/service/service/OrderService%02d.java", idx), strings.Join([]string{
			"package com.acme.service.service;",
			"",
			fmt.Sprintf("import com.acme.service.repository.OrderRepository%02d;", idx),
			"",
			fmt.Sprintf("public class OrderService%02d {", idx),
			fmt.Sprintf("    private final OrderRepository%02d repository = new OrderRepository%02d();", idx, idx),
			"    public int execute() {",
			"        return repository.findTotal();",
			"    }",
			"}",
			"",
		}, "\n"))
		writeFixtureFile(t, repoRoot, fmt.Sprintf("src/main/java/com/acme/service/controller/OrderController%02d.java", idx), strings.Join([]string{
			"package com.acme.service.controller;",
			"",
			fmt.Sprintf("import com.acme.service.service.OrderService%02d;", idx),
			"",
			fmt.Sprintf("public class OrderController%02d {", idx),
			fmt.Sprintf("    private final OrderService%02d service = new OrderService%02d();", idx, idx),
			"    public int get() {",
			"        return service.execute();",
			"    }",
			"}",
			"",
		}, "\n"))
	}
}

func writeJavaMultiModuleEnterpriseBenchmarkFixture(t testing.TB, repoRoot string) {
	t.Helper()

	writeFixtureFile(t, repoRoot, "settings.gradle", strings.Join([]string{
		`rootProject.name = "bench-java"`,
		`include("shared", "app")`,
		"",
	}, "\n"))
	for idx := 0; idx < 12; idx++ {
		writeFixtureFile(t, repoRoot, fmt.Sprintf("shared/src/main/java/com/acme/shared/SharedMath%02d.java", idx), strings.Join([]string{
			"package com.acme.shared;",
			"",
			fmt.Sprintf("public class SharedMath%02d {", idx),
			"    public static int add(int left, int right) {",
			"        return left + right;",
			"    }",
			"}",
			"",
		}, "\n"))
		writeFixtureFile(t, repoRoot, fmt.Sprintf("app/src/main/java/com/acme/app/AppMain%02d.java", idx), strings.Join([]string{
			"package com.acme.app;",
			"",
			fmt.Sprintf("import com.acme.shared.SharedMath%02d;", idx),
			"",
			fmt.Sprintf("public class AppMain%02d {", idx),
			"    public int run() {",
			fmt.Sprintf("        return SharedMath%02d.add(20, 22);", idx),
			"    }",
			"}",
			"",
		}, "\n"))
	}
}

func writeJavaPhase2Fixture(t testing.TB, repoRoot string) {
	t.Helper()

	writeFixtureFile(t, repoRoot, "settings.gradle", strings.Join([]string{
		`rootProject.name = "atlas"`,
		`include("shared-a", "shared-b", "app")`,
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "app/build.gradle", strings.Join([]string{
		"dependencies {",
		`    implementation(project(":shared-b"))`,
		"}",
		"",
	}, "\n"))

	writeFixtureFile(t, repoRoot, "shared-a/src/main/java/com/acme/shareda/Helper.java", strings.Join([]string{
		"package com.acme.shareda;",
		"",
		"public class Helper {",
		"    public static int assist() {",
		"        return 1;",
		"    }",
		"}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "shared-b/src/main/java/com/acme/sharedb/Helper.java", strings.Join([]string{
		"package com.acme.sharedb;",
		"",
		"public class Helper {",
		"    public static int assist() {",
		"        return 2;",
		"    }",
		"}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "shared-b/src/main/java/com/acme/sharedb/Outer.java", strings.Join([]string{
		"package com.acme.sharedb;",
		"",
		"public class Outer {",
		"    public static class Inner {",
		"        public static int assistNested() {",
		"            return 40;",
		"        }",
		"    }",
		"}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "shared-b/src/main/java/com/acme/sharedb/Tooling.java", strings.Join([]string{
		"package com.acme.sharedb;",
		"",
		"public class Tooling {",
		"    public static void noop() {}",
		"}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "app/src/main/java/com/acme/app/Runner.java", strings.Join([]string{
		"package com.acme.app;",
		"",
		"import com.acme.shareda.Helper;",
		"import com.acme.sharedb.Helper;",
		"import com.acme.sharedb.*;",
		"import com.acme.sharedb.Outer.Inner;",
		"",
		"public class Runner {",
		"    public int run() {",
		"        return Helper.assist() + Inner.assistNested();",
		"    }",
		"    public void smoke() {",
		"        Tooling.noop();",
		"    }",
		"}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "app/src/main/java/com/acme/app/AppMain.java", strings.Join([]string{
		"package com.acme.app;",
		"",
		"public class AppMain {",
		"    public int execute() {",
		"        return new Runner().run();",
		"    }",
		"}",
		"",
	}, "\n"))
}

func writeRustWorkspaceFixture(t *testing.T, repoRoot string) {
	t.Helper()

	writeFixtureFile(t, repoRoot, "Cargo.toml", strings.Join([]string{
		"[workspace]",
		`members = ["crates/core", "crates/app"]`,
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/core/Cargo.toml", strings.Join([]string{
		"[package]",
		`name = "openfang-core"`,
		`version = "0.1.0"`,
		`edition = "2021"`,
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/core/src/lib.rs", strings.Join([]string{
		"pub mod util;",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/core/src/util.rs", strings.Join([]string{
		"pub fn helper() {}",
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/app/Cargo.toml", strings.Join([]string{
		"[package]",
		`name = "openfang-app"`,
		`version = "0.1.0"`,
		`edition = "2021"`,
		"",
	}, "\n"))
	writeFixtureFile(t, repoRoot, "crates/app/src/lib.rs", strings.Join([]string{
		"use openfang_core::util::helper;",
		"",
		"pub fn run() {",
		"\thelper();",
		"}",
		"",
	}, "\n"))
}

func writeFixtureFile(t testing.TB, rootPath, relativePath, content string) {
	t.Helper()

	absolutePath := filepath.Join(rootPath, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(value, part) {
			return false
		}
	}

	return true
}
