//go:build mage

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
)

const (
	golangciLintVersion   = "v2.11.4"
	gotestsumVersion      = "v1.13.0"
	goplsModernizeVersion = "v0.22.0"
	binDir                = "bin"
	cliBinary             = "kb"
	versionPackage        = "github.com/compozy/kb/internal/version"
)

var Default = Verify

func Deps() error {
	return sh.RunV("go", "mod", "tidy")
}

func Fmt() error {
	files, err := goFiles(".")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	args := append([]string{"-w"}, files...)
	return sh.RunV("gofmt", args...)
}

func Lint() error {
	if err := runGolangCILint(); err != nil {
		return err
	}
	return Modernize()
}

func runGolangCILint() error {
	args := []string{"run", "./..."}
	if hasPinnedTool("golangci-lint", golangciLintVersion) {
		return sh.RunV("golangci-lint", args...)
	}

	goRunArgs := append(
		[]string{"run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@" + golangciLintVersion},
		args...,
	)
	return sh.RunV("go", goRunArgs...)
}

// Modernize runs gopls' modernize analyzer to enforce idiomatic Go
// (any, min/max, range-over-int, slices/maps helpers, strings.Cut, ...).
// CGO stays enabled because the tree-sitter adapters require it to type-check.
func Modernize() error {
	return sh.RunV(
		"go",
		"run",
		"golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@"+goplsModernizeVersion,
		"./...",
	)
}

// hasPinnedTool reports whether an executable named name is on PATH and reports
// the pinned version. It runs `<name> version` and matches the output so a broken
// or mismatched shim (for example a mise shim with no version configured) is
// rejected in favor of `go run <module>@<version>`.
func hasPinnedTool(name string, wantVersion string) bool {
	path, err := exec.LookPath(name)
	if err != nil {
		return false
	}
	output, err := exec.Command(path, "version").CombinedOutput()
	if err != nil {
		return false
	}
	versionToken := "version " + strings.TrimPrefix(wantVersion, "v")
	return bytes.Contains(output, []byte(versionToken))
}

// Test runs unit tests only (no integration tag).
func Test() error {
	return runGoTests("-race", "-parallel=4", "./...")
}

// TestIntegration runs all tests including integration tests.
func TestIntegration() error {
	return runGoTests("-race", "-parallel=4", "-tags", "integration", "./...")
}

func Build() error {
	return buildGo()
}

func buildGo() error {
	ldflags := buildLDFlags()
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return err
	}
	if err := sh.RunV("go", "build", "-ldflags", ldflags, "./..."); err != nil {
		return err
	}
	out := filepath.Join(binDir, cliBinary)
	return sh.RunV("go", "build", "-ldflags", ldflags, "-o", out, "./cmd/"+cliBinary)
}

// Boundaries verifies that package import rules are not violated.
// Rules: no package may import cli/.
func Boundaries() error {
	forbidden := []struct {
		importer string
		imported string
	}{
		{"internal/config", "internal/cli"},
		{"internal/logger", "internal/cli"},
		{"internal/version", "internal/cli"},
		{"internal/kodebase", "internal/cli"},
	}

	violations := 0
	for _, rule := range forbidden {
		importerDir := rule.importer
		if _, err := os.Stat(importerDir); os.IsNotExist(err) {
			continue
		}
		importPath := "github.com/compozy/kb/" + rule.imported
		cmd := exec.Command("grep", "-r", "--include=*.go", "-l", importPath, importerDir)
		out, err := cmd.Output()
		if err != nil {
			continue // grep returns exit 1 when no match — that's good
		}
		if len(strings.TrimSpace(string(out))) > 0 {
			fmt.Printf("VIOLATION: %s imports %s\n", rule.importer, rule.imported)
			for _, f := range strings.Split(strings.TrimSpace(string(out)), "\n") {
				fmt.Printf("  %s\n", f)
			}
			violations++
		}
	}

	if violations > 0 {
		return fmt.Errorf("found %d boundary violations", violations)
	}
	fmt.Println("OK: all package boundaries respected")
	return nil
}

func Verify() error {
	steps := []func() error{
		Fmt,
		Lint,
		Test,
		buildGo,
		Boundaries,
	}

	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}

	return nil
}

func goFiles(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if path != root && (name == "vendor" || strings.HasPrefix(name, ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func buildLDFlags() string {
	version := gitOutput("describe", "--tags", "--always", "--dirty")
	if version == "" {
		version = "dev"
	}

	commit := gitOutput("rev-parse", "--short", "HEAD")
	if commit == "" {
		commit = "unknown"
	}

	buildDate := time.Now().UTC().Format(time.RFC3339)

	return strings.Join([]string{
		"-X " + versionPackage + ".Version=" + version,
		"-X " + versionPackage + ".Commit=" + commit,
		"-X " + versionPackage + ".Date=" + buildDate,
	}, " ")
}

func gitOutput(args ...string) string {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func runGoTests(testArgs ...string) error {
	// Always go run the pinned gotestsum so a broken or mismatched shim on PATH
	// cannot break the test gate. The build is cached after the first run.
	args := append([]string{"run", "gotest.tools/gotestsum@" + gotestsumVersion, "--format", "pkgname", "--"}, testArgs...)
	return sh.RunV("go", args...)
}
