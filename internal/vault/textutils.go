package vault

import (
	"regexp"
	"strings"
)

var (
	leadingBlockCommentPattern = regexp.MustCompile(`(?s)^\s*/\*\*?[\s\S]*?\*/`)
	leadingLineCommentPattern  = regexp.MustCompile(`^(?:\s*//.*(?:\n|$))+`)
	blockCommentStartPattern   = regexp.MustCompile(`^/\*+\s*`)
	blockCommentEndPattern     = regexp.MustCompile(`\s*\*/$`)
	blockCommentLinePattern    = regexp.MustCompile(`^\s*\*?\s?`)
	lineCommentPrefixPattern   = regexp.MustCompile(`^\s*//\s?`)
)

// NormalizeComment strips Go/TS comment delimiters while preserving the comment text.
func NormalizeComment(rawComment string) string {
	trimmed := strings.TrimSpace(rawComment)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "/*") {
		trimmed = blockCommentStartPattern.ReplaceAllString(trimmed, "")
		trimmed = blockCommentEndPattern.ReplaceAllString(trimmed, "")

		lines := strings.Split(trimmed, "\n")
		for index, line := range lines {
			lines[index] = strings.TrimRight(blockCommentLinePattern.ReplaceAllString(line, ""), " \t\r")
		}

		return strings.TrimSpace(strings.Join(lines, "\n"))
	}

	return normalizeLineCommentBlock(trimmed)
}

// ExtractLeadingComment returns the first leading block or line comment from source text.
func ExtractLeadingComment(sourceText string) string {
	if blockMatch := leadingBlockCommentPattern.FindString(sourceText); blockMatch != "" {
		if normalized := NormalizeComment(blockMatch); normalized != "" {
			return normalized
		}
		return ""
	}

	if lineMatch := leadingLineCommentPattern.FindString(sourceText); lineMatch != "" {
		if normalized := normalizeLineCommentBlock(lineMatch); normalized != "" {
			return normalized
		}
	}

	return ""
}

// StripQuotes removes a single leading and trailing quote character when present.
func StripQuotes(value string) string {
	if value == "" {
		return ""
	}

	if isQuote(value[0]) {
		value = value[1:]
	}

	if value == "" {
		return ""
	}

	if isQuote(value[len(value)-1]) {
		value = value[:len(value)-1]
	}

	return value
}

func normalizeLineCommentBlock(rawComment string) string {
	lines := strings.Split(rawComment, "\n")
	for index, line := range lines {
		lines[index] = strings.TrimRight(lineCommentPrefixPattern.ReplaceAllString(line, ""), " \t\r")
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func isQuote(char byte) bool {
	return char == '\'' || char == '"' || char == '`'
}
