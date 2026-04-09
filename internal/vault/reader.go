package vault

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// VaultRelation describes a parsed relation link in a vault markdown document.
type VaultRelation struct {
	TargetPath string `json:"targetPath"`
	Type       string `json:"type"`
	Confidence string `json:"confidence"`
}

// VaultDocument is the parsed read-side representation of a vault markdown file.
type VaultDocument struct {
	RelativePath      string                 `json:"relativePath"`
	Frontmatter       map[string]interface{} `json:"frontmatter"`
	Body              string                 `json:"body"`
	Backlinks         []VaultRelation        `json:"backlinks"`
	OutgoingRelations []VaultRelation        `json:"outgoingRelations"`
}

// VaultSnapshot groups parsed vault documents by their source category.
type VaultSnapshot struct {
	VaultPath   string          `json:"vaultPath"`
	TopicSlug   string          `json:"topicSlug"`
	Symbols     []VaultDocument `json:"symbols"`
	Files       []VaultDocument `json:"files"`
	Directories []VaultDocument `json:"directories"`
	Wikis       []VaultDocument `json:"wikis"`
}

// ReadVaultOptions controls non-fatal reader behavior.
type ReadVaultOptions struct {
	Warn func(message string)
}

type vaultDocumentBucket string

const (
	vaultBucketSymbols     vaultDocumentBucket = "symbols"
	vaultBucketFiles       vaultDocumentBucket = "files"
	vaultBucketDirectories vaultDocumentBucket = "directories"
	vaultBucketWikis       vaultDocumentBucket = "wikis"
)

var (
	outgoingRelationPattern = regexp.MustCompile("- `(\\w+)` \\((\\w+)\\) -> \\[\\[([^\\]]+)\\]\\]")
	backlinkPattern         = regexp.MustCompile("- \\[\\[([^\\]]+)\\]\\] via `(\\w+)` \\((\\w+)\\)")
)

// ReadVaultSnapshot walks a topic directory and parses every managed markdown file.
func ReadVaultSnapshot(resolvedVault ResolvedVault, options ReadVaultOptions) (VaultSnapshot, error) {
	warn := options.Warn
	if warn == nil {
		warn = func(string) {}
	}

	snapshot := createEmptySnapshot(resolvedVault)

	markdownFiles, err := collectMarkdownFiles(resolvedVault.TopicPath)
	if err != nil {
		return VaultSnapshot{}, fmt.Errorf("read vault snapshot: %w", err)
	}

	for _, markdownFile := range markdownFiles {
		relativePath, err := filepath.Rel(resolvedVault.TopicPath, markdownFile)
		if err != nil {
			return VaultSnapshot{}, fmt.Errorf("read vault snapshot: derive relative path for %q: %w", markdownFile, err)
		}

		content, err := os.ReadFile(markdownFile)
		if err != nil {
			return VaultSnapshot{}, fmt.Errorf("read vault snapshot: read %q: %w", markdownFile, err)
		}

		document, ok := parseVaultDocument(string(content), ToPosixPath(relativePath), warn)
		if !ok {
			continue
		}

		switch classifyDocument(document.Frontmatter) {
		case vaultBucketSymbols:
			snapshot.Symbols = append(snapshot.Symbols, document)
		case vaultBucketFiles:
			snapshot.Files = append(snapshot.Files, document)
		case vaultBucketDirectories:
			snapshot.Directories = append(snapshot.Directories, document)
		default:
			snapshot.Wikis = append(snapshot.Wikis, document)
		}
	}

	sortVaultDocuments(snapshot.Symbols)
	sortVaultDocuments(snapshot.Files)
	sortVaultDocuments(snapshot.Directories)
	sortVaultDocuments(snapshot.Wikis)

	return snapshot, nil
}

// ExtractSection returns the markdown content under the named level-two heading.
func ExtractSection(body, heading string) string {
	trimmedHeading := strings.TrimSpace(heading)
	if trimmedHeading == "" {
		return ""
	}

	normalizedBody := strings.ReplaceAll(body, "\r\n", "\n")
	lines := strings.Split(normalizedBody, "\n")
	headingLine := "## " + trimmedHeading

	collecting := false
	collected := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if !collecting {
			if trimmedLine == headingLine {
				collecting = true
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "## ") {
			break
		}

		collected = append(collected, line)
	}

	return strings.TrimSpace(strings.Join(collected, "\n"))
}

// FindSymbolsByName returns symbol documents whose symbol_name frontmatter contains the query.
func FindSymbolsByName(snapshot VaultSnapshot, query string) []VaultDocument {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	if normalizedQuery == "" {
		return []VaultDocument{}
	}

	matches := make([]VaultDocument, 0, len(snapshot.Symbols))
	for _, document := range snapshot.Symbols {
		symbolName := strings.ToLower(frontmatterString(document.Frontmatter, "symbol_name"))
		if symbolName == "" || !strings.Contains(symbolName, normalizedQuery) {
			continue
		}

		matches = append(matches, document)
	}

	sort.Slice(matches, func(i, j int) bool {
		leftName := frontmatterString(matches[i].Frontmatter, "symbol_name")
		rightName := frontmatterString(matches[j].Frontmatter, "symbol_name")
		if leftName != rightName {
			return leftName < rightName
		}

		leftSource := frontmatterString(matches[i].Frontmatter, "source_path")
		rightSource := frontmatterString(matches[j].Frontmatter, "source_path")
		if leftSource != rightSource {
			return leftSource < rightSource
		}

		return frontmatterInt(matches[i].Frontmatter, "start_line") < frontmatterInt(matches[j].Frontmatter, "start_line")
	})

	return matches
}

func createEmptySnapshot(resolvedVault ResolvedVault) VaultSnapshot {
	return VaultSnapshot{
		VaultPath:   resolvedVault.VaultPath,
		TopicSlug:   resolvedVault.TopicSlug,
		Symbols:     []VaultDocument{},
		Files:       []VaultDocument{},
		Directories: []VaultDocument{},
		Wikis:       []VaultDocument{},
	}
}

func collectMarkdownFiles(rootPath string) ([]string, error) {
	info, err := os.Stat(rootPath)
	switch {
	case os.IsNotExist(err):
		return []string{}, nil
	case err != nil:
		return nil, fmt.Errorf("stat topic path %q: %w", rootPath, err)
	case !info.IsDir():
		return nil, fmt.Errorf("topic path is not a directory: %s", rootPath)
	}

	files := make([]string, 0)
	if err := filepath.WalkDir(rootPath, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			files = append(files, currentPath)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("scan markdown files under %q: %w", rootPath, err)
	}

	sort.Strings(files)
	return files, nil
}

func parseVaultDocument(markdown, relativePath string, warn func(string)) (VaultDocument, bool) {
	frontmatterMatch := frontmatterBlockPattern.FindStringSubmatch(markdown)
	if frontmatterMatch == nil {
		return VaultDocument{}, false
	}

	var parsedFrontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatterMatch[1]), &parsedFrontmatter); err != nil {
		warn(fmt.Sprintf("Skipping %s: malformed YAML frontmatter (%s).", relativePath, err))
		return VaultDocument{}, false
	}
	if len(parsedFrontmatter) == 0 {
		warn(fmt.Sprintf("Skipping %s: frontmatter must parse into an object.", relativePath))
		return VaultDocument{}, false
	}

	normalizedFrontmatter := normalizeFrontmatterMap(parsedFrontmatter)
	body := markdown[len(frontmatterMatch[0]):]

	return VaultDocument{
		RelativePath:      relativePath,
		Frontmatter:       normalizedFrontmatter,
		Body:              body,
		OutgoingRelations: parseRelations(ExtractSection(body, "Outgoing Relations"), outgoingRelationPattern),
		Backlinks:         parseBacklinks(ExtractSection(body, "Backlinks")),
	}, true
}

func classifyDocument(frontmatter map[string]interface{}) vaultDocumentBucket {
	switch frontmatterString(frontmatter, "source_kind") {
	case "codebase-symbol":
		return vaultBucketSymbols
	case "codebase-file":
		return vaultBucketFiles
	case "codebase-directory-index":
		return vaultBucketDirectories
	default:
		return vaultBucketWikis
	}
}

func parseRelations(section string, pattern *regexp.Regexp) []VaultRelation {
	if strings.TrimSpace(section) == "" {
		return []VaultRelation{}
	}

	matches := pattern.FindAllStringSubmatch(section, -1)
	relations := make([]VaultRelation, 0, len(matches))
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		relations = append(relations, VaultRelation{
			Type:       match[1],
			Confidence: match[2],
			TargetPath: match[3],
		})
	}

	return relations
}

func parseBacklinks(section string) []VaultRelation {
	if strings.TrimSpace(section) == "" {
		return []VaultRelation{}
	}

	matches := backlinkPattern.FindAllStringSubmatch(section, -1)
	relations := make([]VaultRelation, 0, len(matches))
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		relations = append(relations, VaultRelation{
			TargetPath: match[1],
			Type:       match[2],
			Confidence: match[3],
		})
	}

	return relations
}

func sortVaultDocuments(documents []VaultDocument) {
	sort.Slice(documents, func(i, j int) bool {
		return documents[i].RelativePath < documents[j].RelativePath
	})
}

func normalizeFrontmatterMap(values map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(values))
	for key, value := range values {
		normalized[key] = normalizeFrontmatterValue(value)
	}
	return normalized
}

func normalizeFrontmatterValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return normalizeFrontmatterMap(typed)
	case map[interface{}]interface{}:
		normalized := make(map[string]interface{}, len(typed))
		for key, nestedValue := range typed {
			normalized[fmt.Sprint(key)] = normalizeFrontmatterValue(nestedValue)
		}
		return normalized
	case []interface{}:
		normalizedValues := make([]interface{}, len(typed))
		stringValues := make([]string, len(typed))
		allStrings := true

		for index, item := range typed {
			normalizedItem := normalizeFrontmatterValue(item)
			normalizedValues[index] = normalizedItem

			stringValue, ok := normalizedItem.(string)
			if !ok {
				allStrings = false
				continue
			}

			stringValues[index] = stringValue
		}

		if allStrings {
			return stringValues
		}

		return normalizedValues
	default:
		return typed
	}
}

func frontmatterString(frontmatter map[string]interface{}, key string) string {
	value, exists := frontmatter[key]
	if !exists {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	default:
		return fmt.Sprint(typed)
	}
}

func frontmatterInt(frontmatter map[string]interface{}, key string) int {
	value, exists := frontmatter[key]
	if !exists {
		return 0
	}

	switch typed := value.(type) {
	case int:
		return typed
	case int8:
		return int(typed)
	case int16:
		return int(typed)
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case uint:
		return int(typed)
	case uint8:
		return int(typed)
	case uint16:
		return int(typed)
	case uint32:
		return int(typed)
	case uint64:
		return int(typed)
	case float32:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}
