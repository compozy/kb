package repohealth

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

// FindPackageBoundaryViolations checks every internal package outside cli for
// imports of the CLI layer. It includes tests and platform-specific files but
// excludes test fixtures and vendored sources. Read and parse failures are errors.
func FindPackageBoundaryViolations(root string) ([]PathViolation, error) {
	const cliPackage = "github.com/compozy/kb/internal/cli"
	internal := filepath.Join(root, "internal")
	var violations []PathViolation
	err := filepath.WalkDir(internal, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path == filepath.Join(internal, "cli") || entry.Name() == "testdata" || entry.Name() == "vendor" || strings.HasPrefix(entry.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return fmt.Errorf("parse import in %s: %w", path, err)
			}
			if importPath == cliPackage || strings.HasPrefix(importPath, cliPackage+"/") {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				violations = append(violations, PathViolation{Path: filepath.ToSlash(rel), Reason: "imports " + importPath})
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("check package boundaries: %w", err)
	}
	return violations, nil
}
