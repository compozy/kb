package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

var requiredCodebaseIngestResultKeys = []string{
	"topic",
	"sourceType",
	"filePath",
	"title",
	"summary",
}

var requiredGenerationSummaryKeys = []string{
	"command",
	"rootPath",
	"vaultPath",
	"topicPath",
	"topicSlug",
	"dryRun",
	"detectedLanguages",
	"selectedAdapters",
	"filesScanned",
	"filesParsed",
	"filesSkipped",
	"symbolsExtracted",
	"relationsEmitted",
	"rawDocumentsWritten",
	"wikiDocumentsWritten",
	"indexDocumentsWritten",
	"timings",
	"diagnostics",
}

var requiredGenerationTimingKeys = []string{
	"scanMillis",
	"selectAdaptersMillis",
	"parseMillis",
	"normalizeMillis",
	"metricsMillis",
	"renderMillis",
	"writeMillis",
	"totalMillis",
}

func decodeJSONMap(t *testing.T, payload []byte) map[string]any {
	t.Helper()

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("decode payload as map: %v", err)
	}

	return decoded
}

func assertCodebaseIngestContractShape(t *testing.T, payload map[string]any) {
	t.Helper()

	assertMapHasKeys(t, payload, requiredCodebaseIngestResultKeys)

	summary := requireMapValue(t, payload, "summary")
	assertMapHasKeys(t, summary, requiredGenerationSummaryKeys)

	timings := requireMapValue(t, summary, "timings")
	assertMapHasKeys(t, timings, requiredGenerationTimingKeys)
}

func assertCodebaseIngestContractSemantics(t *testing.T, payload map[string]any, expectDryRun bool) {
	t.Helper()

	summary := requireMapValue(t, payload, "summary")

	if got, want := requireStringValue(t, payload, "topic"), requireStringValue(t, summary, "topicSlug"); got != want {
		t.Fatalf("topic (%q) must match summary.topicSlug (%q)", got, want)
	}
	if got := requireStringValue(t, payload, "sourceType"); got != string(models.SourceKindCodebaseFile) {
		t.Fatalf("sourceType = %q, want %q", got, models.SourceKindCodebaseFile)
	}
	if got := requireBoolValue(t, summary, "dryRun"); got != expectDryRun {
		t.Fatalf("summary.dryRun = %t, want %t", got, expectDryRun)
	}

	rawWritten := requireIntValue(t, summary, "rawDocumentsWritten")
	wikiWritten := requireIntValue(t, summary, "wikiDocumentsWritten")
	indexWritten := requireIntValue(t, summary, "indexDocumentsWritten")

	if expectDryRun {
		if rawWritten != 0 || wikiWritten != 0 || indexWritten != 0 {
			t.Fatalf(
				"dry-run writes must stay zero, got raw=%d wiki=%d index=%d",
				rawWritten,
				wikiWritten,
				indexWritten,
			)
		}
		return
	}

	if rawWritten <= 0 {
		t.Fatalf("rawDocumentsWritten = %d, want > 0 on full ingest", rawWritten)
	}
}

func assertMapHasKeys(t *testing.T, values map[string]any, required []string) {
	t.Helper()

	for _, key := range required {
		if _, ok := values[key]; !ok {
			t.Fatalf("payload missing required key %q: %#v", key, values)
		}
	}
}

func requireMapValue(t *testing.T, values map[string]any, key string) map[string]any {
	t.Helper()

	raw, ok := values[key]
	if !ok {
		t.Fatalf("payload missing map key %q", key)
	}
	typed, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("key %q has type %T, want map[string]any", key, raw)
	}

	return typed
}

func requireStringValue(t *testing.T, values map[string]any, key string) string {
	t.Helper()

	raw, ok := values[key]
	if !ok {
		t.Fatalf("payload missing string key %q", key)
	}
	typed, ok := raw.(string)
	if !ok {
		t.Fatalf("key %q has type %T, want string", key, raw)
	}

	return typed
}

func requireBoolValue(t *testing.T, values map[string]any, key string) bool {
	t.Helper()

	raw, ok := values[key]
	if !ok {
		t.Fatalf("payload missing bool key %q", key)
	}
	typed, ok := raw.(bool)
	if !ok {
		t.Fatalf("key %q has type %T, want bool", key, raw)
	}

	return typed
}

func requireIntValue(t *testing.T, values map[string]any, key string) int {
	t.Helper()

	raw, ok := values[key]
	if !ok {
		t.Fatalf("payload missing numeric key %q", key)
	}

	switch typed := raw.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case int64:
		return int(typed)
	default:
		t.Fatalf("key %q has type %T, want number", key, raw)
	}

	return 0
}

func writeJavaMultiModuleCodebaseFixture(t *testing.T, repoRoot string) {
	t.Helper()

	files := map[string]string{
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
		"shared-a/src/main/java/com/acme/shareda/Helper.java": strings.Join([]string{
			"package com.acme.shareda;",
			"",
			"public class Helper {",
			"    public static int assist() {",
			"        return 1;",
			"    }",
			"}",
			"",
		}, "\n"),
		"shared-b/src/main/java/com/acme/sharedb/Helper.java": strings.Join([]string{
			"package com.acme.sharedb;",
			"",
			"public class Helper {",
			"    public static int assist() {",
			"        return 2;",
			"    }",
			"}",
			"",
		}, "\n"),
		"shared-b/src/main/java/com/acme/sharedb/Tooling.java": strings.Join([]string{
			"package com.acme.sharedb;",
			"",
			"public class Tooling {",
			"    public static void noop() {}",
			"}",
			"",
		}, "\n"),
		"shared-b/src/main/java/com/acme/sharedb/Outer.java": strings.Join([]string{
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
		}, "\n"),
		"app/src/main/java/com/acme/app/Runner.java": strings.Join([]string{
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
		}, "\n"),
		"app/src/main/java/com/acme/app/AppMain.java": strings.Join([]string{
			"package com.acme.app;",
			"",
			"public class AppMain {",
			"    public int execute() {",
			"        return new Runner().run();",
			"    }",
			"}",
			"",
		}, "\n"),
	}

	for relativePath, content := range files {
		writeFixtureSourceFile(t, repoRoot, relativePath, content)
	}
}

func validateJavaCodebaseSummary(summary models.GenerationSummary, minFiles int, minSymbols int) error {
	if !containsLanguage(summary.DetectedLanguages, "java") {
		return fmt.Errorf("expected detected languages to include java, got %#v", summary.DetectedLanguages)
	}
	if summary.FilesScanned < minFiles {
		return fmt.Errorf("files scanned = %d, want >= %d", summary.FilesScanned, minFiles)
	}
	if summary.SymbolsExtracted < minSymbols {
		return fmt.Errorf("symbols extracted = %d, want >= %d", summary.SymbolsExtracted, minSymbols)
	}

	return nil
}

func containsLanguage(languages []string, want string) bool {
	for _, language := range languages {
		if strings.EqualFold(strings.TrimSpace(language), strings.TrimSpace(want)) {
			return true
		}
	}

	return false
}

func assertJavaCodebaseSummary(t *testing.T, summary models.GenerationSummary, minFiles int, minSymbols int) {
	t.Helper()

	if err := validateJavaCodebaseSummary(summary, minFiles, minSymbols); err != nil {
		t.Fatal(err)
	}
}

func writeFixtureSourceFile(t *testing.T, rootPath string, relativePath string, content string) {
	t.Helper()

	targetPath := filepath.Join(rootPath, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		t.Fatalf("create parent directory for %q: %v", targetPath, err)
	}
	if err := os.WriteFile(targetPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture file %q: %v", targetPath, err)
	}
}

func TestWriteJavaMultiModuleCodebaseFixtureCreatesDeterministicLayout(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	writeJavaMultiModuleCodebaseFixture(t, repoRoot)

	expected := map[string]string{
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
		"shared-a/src/main/java/com/acme/shareda/Helper.java": strings.Join([]string{
			"package com.acme.shareda;",
			"",
			"public class Helper {",
			"    public static int assist() {",
			"        return 1;",
			"    }",
			"}",
			"",
		}, "\n"),
		"shared-b/src/main/java/com/acme/sharedb/Helper.java": strings.Join([]string{
			"package com.acme.sharedb;",
			"",
			"public class Helper {",
			"    public static int assist() {",
			"        return 2;",
			"    }",
			"}",
			"",
		}, "\n"),
		"shared-b/src/main/java/com/acme/sharedb/Tooling.java": strings.Join([]string{
			"package com.acme.sharedb;",
			"",
			"public class Tooling {",
			"    public static void noop() {}",
			"}",
			"",
		}, "\n"),
		"shared-b/src/main/java/com/acme/sharedb/Outer.java": strings.Join([]string{
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
		}, "\n"),
		"app/src/main/java/com/acme/app/Runner.java": strings.Join([]string{
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
		}, "\n"),
		"app/src/main/java/com/acme/app/AppMain.java": strings.Join([]string{
			"package com.acme.app;",
			"",
			"public class AppMain {",
			"    public int execute() {",
			"        return new Runner().run();",
			"    }",
			"}",
			"",
		}, "\n"),
	}

	for relativePath, want := range expected {
		got, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(relativePath)))
		if err != nil {
			t.Fatalf("read %q: %v", relativePath, err)
		}
		if string(got) != want {
			t.Fatalf("fixture content for %q mismatch\nwant:\n%s\n\ngot:\n%s", relativePath, want, string(got))
		}
	}
}

func TestValidateJavaCodebaseSummary(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		summary   models.GenerationSummary
		minFiles  int
		minSymbol int
		wantErr   string
	}{
		{
			name: "valid summary",
			summary: models.GenerationSummary{
				DetectedLanguages: []string{"go", "java"},
				FilesScanned:      3,
				SymbolsExtracted:  6,
			},
			minFiles:  3,
			minSymbol: 4,
		},
		{
			name: "missing java language",
			summary: models.GenerationSummary{
				DetectedLanguages: []string{"go"},
				FilesScanned:      3,
				SymbolsExtracted:  6,
			},
			minFiles:  3,
			minSymbol: 4,
			wantErr:   "include java",
		},
		{
			name: "insufficient files",
			summary: models.GenerationSummary{
				DetectedLanguages: []string{"java"},
				FilesScanned:      1,
				SymbolsExtracted:  6,
			},
			minFiles:  2,
			minSymbol: 4,
			wantErr:   "files scanned",
		},
		{
			name: "insufficient symbols",
			summary: models.GenerationSummary{
				DetectedLanguages: []string{"java"},
				FilesScanned:      3,
				SymbolsExtracted:  1,
			},
			minFiles:  2,
			minSymbol: 4,
			wantErr:   "symbols extracted",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := validateJavaCodebaseSummary(testCase.summary, testCase.minFiles, testCase.minSymbol)
			if testCase.wantErr == "" {
				if err != nil {
					t.Fatalf("validateJavaCodebaseSummary() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("validateJavaCodebaseSummary() expected error containing %q", testCase.wantErr)
			}
			if !strings.Contains(err.Error(), testCase.wantErr) {
				t.Fatalf("validateJavaCodebaseSummary() error = %q, want substring %q", err.Error(), testCase.wantErr)
			}
		})
	}
}

func TestAssertJavaCodebaseSummaryPassesForValidInput(t *testing.T) {
	t.Parallel()

	assertJavaCodebaseSummary(t, models.GenerationSummary{
		DetectedLanguages: []string{"java"},
		FilesScanned:      3,
		SymbolsExtracted:  3,
	}, 3, 3)
}
