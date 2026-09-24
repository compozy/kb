package corpus

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	flattenWikilinkPattern = regexp.MustCompile(`!?\[\[([^\[\]|\n]+?)(?:\|([^\[\]\n]*?))?\]\]`)
	atxHeadingPattern      = regexp.MustCompile(`^ {0,3}#{1,6}(?:\s|$)`)
)

// truncatedMarker ends a head section that alone exceeds the excerpt budget.
const truncatedMarker = "[... truncated ...]"

// EstimateTokens estimates the token count of text as runes/4, rounded up.
// It is a deliberately cheap, provider-independent approximation used only to
// bound state size.
func EstimateTokens(text string) int {
	return (utf8.RuneCountInString(text) + 3) / 4
}

// Excerpt returns a section-aware excerpt of body within about maxTokens
// tokens (EstimateTokens). A body that fits is returned unchanged. Otherwise
// the body is split into sections at ATX headings (outside fenced code) and
// kept in document order by priority: the head (the text before the first
// heading plus the first section), then every section containing one of terms
// (case-insensitive), then the remaining sections, while the budget allows.
// Each run of omitted sections is replaced by `[... N sections omitted ...]`;
// a head section larger than the whole budget is cut and ends with
// `[... truncated ...]`. The tail is never dropped without a marker.
func Excerpt(body string, maxTokens int, terms []string) string {
	if maxTokens <= 0 || EstimateTokens(body) <= maxTokens {
		return body
	}

	sections := splitSections(body)
	headCount := 1
	if len(sections) > 1 && !startsWithHeading(sections[0]) {
		headCount = 2
	}

	budget := maxTokens
	keep := make([]bool, len(sections))
	texts := append([]string(nil), sections...)
	for index := range headCount {
		cost := EstimateTokens(texts[index]) + markerReserve
		if cost <= budget {
			keep[index] = true
			budget -= cost
			continue
		}
		if index == 0 {
			texts[index] = truncateRunes(texts[index], max(0, budget-markerReserve*2)*4) + "\n" + truncatedMarker + "\n"
			keep[index] = true
			budget = 0
		}
		break
	}

	loweredTerms := make([]string, 0, len(terms))
	for _, term := range terms {
		if trimmed := strings.ToLower(strings.TrimSpace(term)); trimmed != "" {
			loweredTerms = append(loweredTerms, trimmed)
		}
	}
	for _, wantMatch := range []bool{true, false} {
		for index := range sections {
			if keep[index] || sectionMatches(sections[index], loweredTerms) != wantMatch {
				continue
			}
			cost := EstimateTokens(sections[index]) + markerReserve
			if cost > budget {
				continue
			}
			keep[index] = true
			budget -= cost
		}
	}

	var builder strings.Builder
	omitted := 0
	flush := func() {
		if omitted == 0 {
			return
		}
		ensureTrailingNewline(&builder)
		label := "sections"
		if omitted == 1 {
			label = "section"
		}
		fmt.Fprintf(&builder, "[... %d %s omitted ...]\n", omitted, label)
		omitted = 0
	}
	for index, section := range texts {
		if !keep[index] {
			omitted++
			continue
		}
		flush()
		ensureTrailingNewline(&builder)
		builder.WriteString(section)
	}
	flush()

	return builder.String()
}

// markerReserve is the token allowance kept per section for an omission
// marker that may follow it.
const markerReserve = 8

// Head flattens wikilinks (`[[a|b]]` → `b`, `[[a]]` → `a`) and returns at most
// maxChars runes of the trimmed result.
func Head(body string, maxChars int) string {
	flattened := flattenWikilinkPattern.ReplaceAllStringFunc(body, func(match string) string {
		groups := flattenWikilinkPattern.FindStringSubmatch(match)
		if strings.TrimSpace(groups[2]) != "" {
			return groups[2]
		}
		target := groups[1]
		if before, _, found := strings.Cut(target, "#"); found && strings.TrimSpace(before) != "" {
			target = before
		}
		return strings.TrimSpace(target)
	})

	return strings.TrimSpace(truncateRunes(strings.TrimSpace(flattened), maxChars))
}

func splitSections(body string) []string {
	sections := make([]string, 0, 8)
	var current strings.Builder
	inFence := false
	fence := ""

	for _, line := range splitKeepNewlines(body) {
		trimmed := strings.TrimSpace(line)
		if marker := fenceOpener(trimmed); marker != "" {
			switch {
			case !inFence:
				inFence, fence = true, marker
			case strings.HasPrefix(trimmed, fence) && strings.Trim(trimmed, fence[:1]) == "":
				inFence = false
			}
		}
		if !inFence && atxHeadingPattern.MatchString(line) && current.Len() > 0 {
			sections = append(sections, current.String())
			current.Reset()
		}
		current.WriteString(line)
	}
	if current.Len() > 0 {
		sections = append(sections, current.String())
	}

	return sections
}

func fenceOpener(trimmed string) string {
	for _, marker := range []string{"```", "~~~"} {
		if strings.HasPrefix(trimmed, marker) {
			return marker
		}
	}

	return ""
}

func splitKeepNewlines(text string) []string {
	lines := make([]string, 0, strings.Count(text, "\n")+1)
	for text != "" {
		index := strings.IndexByte(text, '\n')
		if index < 0 {
			lines = append(lines, text)
			break
		}
		lines = append(lines, text[:index+1])
		text = text[index+1:]
	}

	return lines
}

func startsWithHeading(section string) bool {
	return atxHeadingPattern.MatchString(section)
}

func sectionMatches(section string, terms []string) bool {
	if len(terms) == 0 {
		return false
	}
	lowered := strings.ToLower(section)
	for _, term := range terms {
		if strings.Contains(lowered, term) {
			return true
		}
	}

	return false
}

func truncateRunes(text string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	count := 0
	for index := range text {
		if count == maxRunes {
			return text[:index]
		}
		count++
	}

	return text
}

func ensureTrailingNewline(builder *strings.Builder) {
	if builder.Len() == 0 {
		return
	}
	if !strings.HasSuffix(builder.String(), "\n") {
		builder.WriteString("\n")
	}
}
