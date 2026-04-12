// Package scanner discovers source files in a workspace directory, respecting gitignore rules and custom include/exclude patterns.
package scanner

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"

	"github.com/user/go-devstack/internal/models"
)

var defaultIgnoredPatterns = []string{
	".git/",
	".kodebase/",
	".next/",
	".turbo/",
	"build/",
	"coverage/",
	"dist/",
	"node_modules/",
	"vendor/",
}

var hardSkippedDirectoryNames = map[string]struct{}{
	".git": {},
	".hg":  {},
	".svn": {},
}

// ScanOptions configures a workspace scan.
type ScanOptions struct {
	OutputPath      string
	IncludePatterns []string
	ExcludePatterns []string
}

// Option mutates ScanOptions for a scanner.
type Option func(*ScanOptions)

// Scanner discovers supported source files inside a repository workspace.
type Scanner struct {
	options ScanOptions
}

type ignoreRule struct {
	relativeDirectory string
	negated           bool
	matches           func(string) bool
}

// NewScanner constructs a scanner using functional options.
func NewScanner(opts ...Option) *Scanner {
	options := ScanOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return &Scanner{options: options}
}

// WithOutputPath excludes the generated output directory from scan results.
func WithOutputPath(path string) Option {
	return func(options *ScanOptions) {
		options.OutputPath = path
	}
}

// WithIncludePatterns configures user include patterns that re-include paths.
func WithIncludePatterns(patterns ...string) Option {
	return func(options *ScanOptions) {
		options.IncludePatterns = append([]string(nil), patterns...)
	}
}

// WithExcludePatterns configures user exclude patterns.
func WithExcludePatterns(patterns ...string) Option {
	return func(options *ScanOptions) {
		options.ExcludePatterns = append([]string(nil), patterns...)
	}
}

// ScanWorkspace is a convenience entrypoint for a one-off scan.
func ScanWorkspace(rootPath string, opts ...Option) (*models.ScannedWorkspace, error) {
	return NewScanner(opts...).ScanWorkspace(rootPath)
}

// ScanWorkspace scans a repository root and returns supported source files grouped by language.
func (s *Scanner) ScanWorkspace(rootPath string) (*models.ScannedWorkspace, error) {
	absoluteRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("resolve root path: %w", err)
	}

	rootInfo, err := os.Stat(absoluteRootPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("repository root does not exist or is not a directory: %s", absoluteRootPath)
		}

		return nil, fmt.Errorf("stat repository root: %w", err)
	}

	if !rootInfo.IsDir() {
		return nil, fmt.Errorf("repository root does not exist or is not a directory: %s", absoluteRootPath)
	}

	outputPath, err := resolveOutputPath(s.options.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("resolve output path: %w", err)
	}

	rules, err := collectIgnoreRules(absoluteRootPath, s.options)
	if err != nil {
		return nil, err
	}

	files := make([]models.ScannedSourceFile, 0)
	if err := filepath.WalkDir(absoluteRootPath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}

		if path == absoluteRootPath {
			return nil
		}

		if outputPath != "" && sameOrInsidePath(outputPath, path) {
			if entry.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		if entry.IsDir() {
			if shouldSkipHardDirectory(entry.Name()) {
				return fs.SkipDir
			}
		}

		relativePath, err := filepath.Rel(absoluteRootPath, path)
		if err != nil {
			return fmt.Errorf("derive relative path for %s: %w", path, err)
		}

		relativePath = filepath.ToSlash(relativePath)
		if entry.IsDir() {
			return nil
		}

		if isIgnored(relativePath, false, rules) {
			return nil
		}

		language, ok := supportedLanguage(relativePath)
		if !ok {
			return nil
		}

		files = append(files, models.ScannedSourceFile{
			AbsolutePath: path,
			RelativePath: relativePath,
			Language:     language,
		})

		return nil
	}); err != nil {
		return nil, err
	}

	sort.Slice(files, func(left, right int) bool {
		return files[left].RelativePath < files[right].RelativePath
	})

	filesByLanguage := make(map[models.SupportedLanguage][]models.ScannedSourceFile)
	for _, file := range files {
		filesByLanguage[file.Language] = append(filesByLanguage[file.Language], file)
	}

	return &models.ScannedWorkspace{
		Files:           files,
		FilesByLanguage: filesByLanguage,
	}, nil
}

func supportedLanguage(filePath string) (models.SupportedLanguage, bool) {
	switch {
	case strings.HasSuffix(filePath, ".d.ts"):
		return "", false
	case strings.HasSuffix(filePath, ".tsx"):
		return models.LangTSX, true
	case strings.HasSuffix(filePath, ".ts"):
		return models.LangTS, true
	case strings.HasSuffix(filePath, ".jsx"):
		return models.LangJSX, true
	case strings.HasSuffix(filePath, ".js"):
		return models.LangJS, true
	case strings.HasSuffix(filePath, ".go"):
		return models.LangGo, true
	default:
		return "", false
	}
}

func resolveOutputPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}

	resolvedPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return resolvedPath, nil
}

func sameOrInsidePath(rootPath string, candidatePath string) bool {
	relativePath, err := filepath.Rel(rootPath, candidatePath)
	if err != nil {
		return false
	}

	if relativePath == "." {
		return true
	}

	return relativePath != ".." && !strings.HasPrefix(relativePath, ".."+string(filepath.Separator))
}

func shouldSkipHardDirectory(name string) bool {
	_, exists := hardSkippedDirectoryNames[name]
	return exists
}

func collectIgnoreRules(rootPath string, options ScanOptions) ([]ignoreRule, error) {
	rules := buildRules(".", defaultIgnoredPatterns)

	gitIgnorePaths, err := collectGitIgnorePaths(rootPath)
	if err != nil {
		return nil, err
	}

	for _, relativePath := range gitIgnorePaths {
		fileRules, err := readGitIgnoreRules(rootPath, relativePath)
		if err != nil {
			return nil, err
		}

		rules = append(rules, fileRules...)
	}

	rules = append(rules, buildUserRules(options)...)
	return rules, nil
}

func collectGitIgnorePaths(rootPath string) ([]string, error) {
	paths := make([]string, 0)
	if err := filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}

		if path == rootPath {
			return nil
		}

		if entry.IsDir() && shouldSkipHardDirectory(entry.Name()) {
			return fs.SkipDir
		}

		if entry.IsDir() || entry.Name() != ".gitignore" {
			return nil
		}

		relativePath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return fmt.Errorf("derive relative gitignore path for %s: %w", path, err)
		}

		paths = append(paths, filepath.ToSlash(relativePath))
		return nil
	}); err != nil {
		return nil, err
	}

	sort.Strings(paths)
	return paths, nil
}

func readGitIgnoreRules(rootPath string, relativePath string) ([]ignoreRule, error) {
	contents, err := os.ReadFile(filepath.Join(rootPath, filepath.FromSlash(relativePath)))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", relativePath, err)
	}

	relativeDirectory := filepath.ToSlash(filepath.Dir(relativePath))
	return buildRules(relativeDirectory, strings.Split(string(contents), "\n")), nil
}

func buildUserRules(options ScanOptions) []ignoreRule {
	patterns := make([]string, 0, len(options.ExcludePatterns)+len(options.IncludePatterns))
	for _, pattern := range options.ExcludePatterns {
		normalizedPattern := normalizePattern(pattern)
		if normalizedPattern == "" {
			continue
		}

		patterns = append(patterns, strings.TrimLeft(normalizedPattern, "!"))
	}

	for _, pattern := range options.IncludePatterns {
		normalizedPattern := normalizePattern(pattern)
		if normalizedPattern == "" {
			continue
		}

		if !strings.HasPrefix(normalizedPattern, "!") {
			normalizedPattern = "!" + normalizedPattern
		}

		patterns = append(patterns, normalizedPattern)
	}

	return buildRules(".", patterns)
}

func normalizePattern(pattern string) string {
	normalizedPattern := filepath.ToSlash(strings.TrimSpace(pattern))
	for strings.HasPrefix(normalizedPattern, "./") {
		normalizedPattern = strings.TrimPrefix(normalizedPattern, "./")
	}

	return normalizedPattern
}

func buildRules(relativeDirectory string, patterns []string) []ignoreRule {
	rules := make([]ignoreRule, 0, len(patterns))
	for _, pattern := range patterns {
		if !shouldKeepPattern(pattern) {
			continue
		}

		negated := isNegatedPattern(pattern)
		if negated {
			matcher := ignore.CompileIgnoreLines("*", pattern)
			rules = append(rules, ignoreRule{
				relativeDirectory: relativeDirectory,
				negated:           true,
				matches: func(path string) bool {
					return !matcher.MatchesPath(path)
				},
			})
			continue
		}

		matcher := ignore.CompileIgnoreLines(pattern)
		rules = append(rules, ignoreRule{
			relativeDirectory: relativeDirectory,
			matches:           matcher.MatchesPath,
		})
	}

	return rules
}

func shouldKeepPattern(pattern string) bool {
	trimmed := strings.TrimSpace(strings.TrimRight(pattern, "\r"))
	if trimmed == "" {
		return false
	}

	return !strings.HasPrefix(trimmed, "#")
}

func isNegatedPattern(pattern string) bool {
	trimmed := strings.TrimSpace(strings.TrimRight(pattern, "\r"))
	return strings.HasPrefix(trimmed, "!") && !strings.HasPrefix(trimmed, `\!`)
}

func isIgnored(relativePath string, isDirectory bool, rules []ignoreRule) bool {
	pathToCheck := relativePath
	if isDirectory && !strings.HasSuffix(pathToCheck, "/") {
		pathToCheck += "/"
	}

	ignored := false
	for _, rule := range rules {
		scopedPath, ok := scopePath(rule.relativeDirectory, pathToCheck)
		if !ok {
			continue
		}

		if rule.matches(scopedPath) {
			ignored = !rule.negated
		}
	}

	return ignored
}

func scopePath(relativeDirectory string, relativePath string) (string, bool) {
	if relativeDirectory == "." {
		return relativePath, true
	}

	directoryPrefix := relativeDirectory + "/"
	if !strings.HasPrefix(relativePath, directoryPrefix) {
		return "", false
	}

	return strings.TrimPrefix(relativePath, directoryPrefix), true
}
