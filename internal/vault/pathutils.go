// Package vault provides path resolution, document rendering, reading, writing, and query helpers for Obsidian-compatible knowledge vaults.
package vault

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/user/go-devstack/internal/models"
)

var windowsDriveRootPattern = regexp.MustCompile(`^[A-Za-z]:/?$`)

const wikiConceptFilePrefix = "Kodebase - "

// ToPosixPath normalizes path separators to forward slashes and trims trailing slashes.
func ToPosixPath(value string) string {
	if value == "" {
		return ""
	}

	normalized := strings.ReplaceAll(value, `\`, `/`)
	if windowsDriveRootPattern.MatchString(normalized) {
		return normalized[:2] + "/"
	}

	for len(normalized) > 1 && strings.HasSuffix(normalized, "/") {
		normalized = strings.TrimSuffix(normalized, "/")
	}

	if normalized == "" {
		return "/"
	}

	return normalized
}

// IsPathInside reports whether targetPath is the same as or nested under parentPath.
func IsPathInside(parentPath, targetPath string) bool {
	parent := cleanComparablePath(parentPath)
	target := cleanComparablePath(targetPath)

	switch {
	case parent == "":
		return target == ""
	case target == "":
		return false
	}

	if parent == target {
		return true
	}

	parentDrive, parentParts, parentAbsolute := splitComparablePath(parent)
	targetDrive, targetParts, targetAbsolute := splitComparablePath(target)

	if parentAbsolute != targetAbsolute {
		return false
	}

	if !strings.EqualFold(parentDrive, targetDrive) {
		return false
	}

	if len(parentParts) > len(targetParts) {
		return false
	}

	for index := range parentParts {
		if parentParts[index] != targetParts[index] {
			return false
		}
	}

	return true
}

// CreateFileID creates a stable file node identifier from a file path.
func CreateFileID(filePath string) string {
	return "file:" + ToPosixPath(filePath)
}

// CreateExternalID creates a stable external node identifier from a module source.
func CreateExternalID(source string) string {
	return "external:" + source
}

// SlugifySegment converts a free-form segment into a filesystem-friendly slug.
func SlugifySegment(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return "item"
	}

	var builder strings.Builder
	lastWasDash := false

	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastWasDash = false
			continue
		}

		if lastWasDash {
			continue
		}

		builder.WriteByte('-')
		lastWasDash = true
	}

	slug := strings.Trim(builder.String(), "-")
	if slug == "" {
		return "item"
	}

	return slug
}

// HumanizeSlug converts a hyphenated slug into a title-cased label.
func HumanizeSlug(value string) string {
	if value == "" {
		return ""
	}

	segments := strings.Split(value, "-")
	humanized := make([]string, 0, len(segments))

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		humanized = append(humanized, strings.ToUpper(segment[:1])+segment[1:])
	}

	return strings.Join(humanized, " ")
}

// DeriveTopicSlug derives the topic slug from a root path.
func DeriveTopicSlug(rootPath string) string {
	normalized := ToPosixPath(rootPath)
	if normalized == "" {
		return SlugifySegment("")
	}

	base := strings.TrimRight(normalized, "/")
	if base == "" {
		base = normalized
	}

	return SlugifySegment(path.Base(base))
}

// DeriveTopicTitle converts a topic slug into a human-readable title.
func DeriveTopicTitle(topicSlug string) string {
	if title := HumanizeSlug(topicSlug); title != "" {
		return title
	}

	return "Knowledge Base"
}

// DeriveTopicDomain returns the topic domain identifier for a topic slug.
func DeriveTopicDomain(topicSlug string) string {
	return topicSlug
}

// CreateSymbolID creates a stable symbol identifier from a symbol node.
func CreateSymbolID(symbol models.SymbolNode) string {
	return strings.Join([]string{
		"symbol",
		ToPosixPath(symbol.FilePath),
		SlugifySegment(symbol.Name),
		SlugifySegment(symbol.SymbolKind),
		fmt.Sprintf("%d", symbol.StartLine),
		fmt.Sprintf("%d", symbol.EndLine),
	}, ":")
}

// GetRawFileDocumentPath derives the vault document path for a raw file snapshot.
func GetRawFileDocumentPath(filePath string) string {
	return fmt.Sprintf("raw/codebase/files/%s.md", normalizeDocumentPathSegment(filePath))
}

// GetRawSymbolDocumentPath derives the vault document path for a raw symbol snapshot.
func GetRawSymbolDocumentPath(symbol models.SymbolNode) string {
	fileSlug := SlugifySegment(strings.ReplaceAll(ToPosixPath(symbol.FilePath), ".", "-"))
	symbolSlug := fmt.Sprintf("%s--%s-l%d", SlugifySegment(symbol.Name), fileSlug, symbol.StartLine)
	return fmt.Sprintf("raw/codebase/symbols/%s.md", symbolSlug)
}

// GetRawDirectoryIndexPath derives the vault document path for a raw directory index.
func GetRawDirectoryIndexPath(directoryPath string) string {
	normalized := normalizeDocumentPathSegment(directoryPath)
	if normalized == "" || normalized == "." {
		return "raw/codebase/indexes/directories/root.md"
	}

	return fmt.Sprintf("raw/codebase/indexes/directories/%s.md", normalized)
}

// GetRawLanguageIndexPath derives the vault document path for a raw language index.
func GetRawLanguageIndexPath(language string) string {
	return fmt.Sprintf("raw/codebase/indexes/languages/%s.md", strings.TrimSpace(language))
}

// GetWikiConceptPath derives the vault document path for a generated wiki concept article.
func GetWikiConceptPath(articleTitle string) string {
	normalizedTitle := strings.TrimSpace(articleTitle)
	if !strings.HasPrefix(normalizedTitle, wikiConceptFilePrefix) {
		normalizedTitle = wikiConceptFilePrefix + normalizedTitle
	}

	return fmt.Sprintf("wiki/concepts/%s.md", normalizedTitle)
}

// GetWikiIndexPath derives the vault document path for a generated wiki index page.
func GetWikiIndexPath(indexTitle string) string {
	return fmt.Sprintf("wiki/index/%s.md", indexTitle)
}

// GetBaseFilePath derives the vault document path for an Obsidian Base definition.
func GetBaseFilePath(baseName string) string {
	return fmt.Sprintf("bases/%s.base", baseName)
}

// StripMarkdownExtension removes a trailing .md extension from a document path.
func StripMarkdownExtension(documentPath string) string {
	return strings.TrimSuffix(documentPath, ".md")
}

// ToTopicWikiLink formats a topic-scoped Obsidian wiki-link.
func ToTopicWikiLink(topicSlug, documentPath, label string) string {
	target := fmt.Sprintf("%s/%s", topicSlug, StripMarkdownExtension(documentPath))
	if label != "" {
		return fmt.Sprintf("[[%s|%s]]", target, label)
	}

	return fmt.Sprintf("[[%s]]", target)
}

func stripWikiConceptFilePrefix(articleTitle string) string {
	normalizedTitle := strings.TrimSpace(articleTitle)
	return strings.TrimSpace(strings.TrimPrefix(normalizedTitle, wikiConceptFilePrefix))
}

func cleanComparablePath(value string) string {
	normalized := ToPosixPath(value)
	if normalized == "" || normalized == "/" || windowsDriveRootPattern.MatchString(normalized) {
		return normalized
	}

	return ToPosixPath(path.Clean(normalized))
}

func hasWindowsDrivePrefix(value string) bool {
	return len(value) >= 2 && ((value[0] >= 'A' && value[0] <= 'Z') || (value[0] >= 'a' && value[0] <= 'z')) && value[1] == ':'
}

func normalizeDocumentPathSegment(value string) string {
	return strings.TrimLeft(ToPosixPath(value), "/")
}

func splitComparablePath(value string) (drive string, parts []string, absolute bool) {
	normalized := value
	if hasWindowsDrivePrefix(normalized) {
		drive = normalized[:2]
		normalized = strings.TrimPrefix(normalized[2:], "/")
		absolute = true
	} else {
		absolute = strings.HasPrefix(normalized, "/")
		normalized = strings.TrimPrefix(normalized, "/")
	}

	if normalized == "" || normalized == "." {
		return drive, nil, absolute
	}

	rawParts := strings.Split(normalized, "/")
	parts = make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		if part == "" || part == "." {
			continue
		}
		parts = append(parts, part)
	}

	return drive, parts, absolute
}
