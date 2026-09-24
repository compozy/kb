package resolve

import (
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Link kinds reported by BodyLinks.
const (
	KindWikilink = "wikilink"
	KindMarkdown = "markdown"
)

var (
	bodyWikilinkPattern     = regexp.MustCompile(`(!?)\[\[([^\[\]\n]+?)\]\]`)
	bodyMarkdownLinkPattern = regexp.MustCompile(`(!?)\[([^\[\]\n]*)\]\(\s*(<[^<>\n]+>|[^\s()<>]+)(?:\s+"[^"\n]*")?\s*\)`)
	schemePattern           = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.\-]*:`)
	stringWikilinkPattern   = regexp.MustCompile(`\[\[([^\[\]\n]+?)\]\]`)
)

// Link is one link found in a markdown body. Start and End are byte offsets
// of the whole link markup in the scanned text, so editors can replace
// text[Start:End].
type Link struct {
	// Target is the normalized link target (see Normalize); markdown link
	// destinations are URL-unescaped first.
	Target string
	// Display is the display text: the `|alias` part of a wikilink or the
	// label of a markdown link; "" when absent.
	Display string
	// Heading is the `#heading` (or `#^block`) part without the `#`, if any.
	Heading string
	Start   int
	End     int
	// Kind is KindWikilink or KindMarkdown.
	Kind string
	// Embed reports an embed (`![[x]]` or `![label](x.md)`).
	Embed bool
}

// Range is a half-open byte range [Start, End).
type Range struct {
	Start int
	End   int
}

// BodyLinks extracts wikilinks and internal markdown links from body in
// document order, skipping fenced code blocks and inline code. Markdown links
// are kept only when the destination has no URL scheme and either ends in
// `.md` or has no extension. A `[[label]](url)` construct is a markdown link
// whose label happens to be bracketed, not a wikilink.
func BodyLinks(body string) []Link {
	code := CodeRanges(body)
	links := make([]Link, 0)

	for _, match := range bodyWikilinkPattern.FindAllStringSubmatchIndex(body, -1) {
		start, end := match[0], match[1]
		if inRanges(code, start, end) {
			continue
		}
		after := body[end:]
		if strings.HasPrefix(after, "(") && strings.Contains(after, ")") {
			continue
		}
		inner := body[match[4]:match[5]]
		target, display, heading := splitWikilinkInner(inner)
		normalized := Normalize(target)
		if normalized == "" {
			continue
		}
		links = append(links, Link{
			Target:  normalized,
			Display: display,
			Heading: heading,
			Start:   start,
			End:     end,
			Kind:    KindWikilink,
			Embed:   match[3] > match[2],
		})
	}

	for _, match := range bodyMarkdownLinkPattern.FindAllStringSubmatchIndex(body, -1) {
		start, end := match[0], match[1]
		if inRanges(code, start, end) || overlapsLinks(links, start, end) {
			continue
		}
		destination := strings.TrimSuffix(strings.TrimPrefix(body[match[6]:match[7]], "<"), ">")
		target, heading, ok := internalMarkdownTarget(destination)
		if !ok {
			continue
		}
		links = append(links, Link{
			Target:  target,
			Display: body[match[4]:match[5]],
			Heading: heading,
			Start:   start,
			End:     end,
			Kind:    KindMarkdown,
			Embed:   match[3] > match[2],
		})
	}

	sort.SliceStable(links, func(i, j int) bool { return links[i].Start < links[j].Start })
	return links
}

// FrontmatterLinks returns, per frontmatter key, the normalized wikilink
// targets found in string values or string list values (for example
// `sources`, `related`, `concepts`, `affects`, `supersedes`, `informed_by`).
// Keys without links are omitted; nested maps are ignored.
func FrontmatterLinks(values map[string]any) map[string][]string {
	result := make(map[string][]string)
	for key, value := range values {
		var targets []string
		switch typed := value.(type) {
		case string:
			targets = stringLinks(typed)
		case []string:
			for _, item := range typed {
				targets = append(targets, stringLinks(item)...)
			}
		case []any:
			for _, item := range typed {
				if text, ok := item.(string); ok {
					targets = append(targets, stringLinks(text)...)
				}
			}
		}
		if len(targets) > 0 {
			result[key] = targets
		}
	}

	return result
}

// CodeRanges returns the byte ranges of fenced code blocks (``` or ~~~,
// unterminated fences run to the end) and inline code spans in text.
func CodeRanges(text string) []Range {
	ranges := make([]Range, 0)
	offset := 0
	fenceChar := byte(0)
	fenceLen := 0
	fenceStart := 0
	proseStart := 0

	for offset < len(text) {
		lineEnd := strings.IndexByte(text[offset:], '\n')
		next := len(text)
		if lineEnd >= 0 {
			next = offset + lineEnd + 1
		}
		line := strings.TrimRight(text[offset:next], "\r\n")
		char, length := fenceMarker(line)

		switch {
		case fenceChar == 0 && length > 0:
			ranges = append(ranges, inlineCodeRanges(text, proseStart, offset)...)
			fenceChar, fenceLen, fenceStart = char, length, offset
		case fenceChar != 0 && char == fenceChar && length >= fenceLen && strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(line), string(char))) == "":
			ranges = append(ranges, Range{Start: fenceStart, End: next})
			fenceChar = 0
			proseStart = next
		}
		offset = next
	}

	if fenceChar != 0 {
		ranges = append(ranges, Range{Start: fenceStart, End: len(text)})
	} else {
		ranges = append(ranges, inlineCodeRanges(text, proseStart, len(text))...)
	}

	sort.Slice(ranges, func(i, j int) bool { return ranges[i].Start < ranges[j].Start })
	return ranges
}

// InRanges reports whether [start, end) overlaps any range.
func InRanges(ranges []Range, start, end int) bool {
	return inRanges(ranges, start, end)
}

func fenceMarker(line string) (byte, int) {
	indent := len(line) - len(strings.TrimLeft(line, " "))
	if indent > 3 {
		return 0, 0
	}
	trimmed := line[indent:]
	if trimmed == "" || (trimmed[0] != '`' && trimmed[0] != '~') {
		return 0, 0
	}
	char := trimmed[0]
	length := len(trimmed) - len(strings.TrimLeft(trimmed, string(char)))
	if length < 3 {
		return 0, 0
	}
	if char == '`' && strings.Contains(trimmed[length:], "`") {
		return 0, 0
	}

	return char, length
}

func inlineCodeRanges(text string, start, end int) []Range {
	ranges := make([]Range, 0)
	index := start
	for index < end {
		if text[index] != '`' {
			index++
			continue
		}
		runLen := backtickRun(text, index, end)
		closing := findBacktickRun(text, index+runLen, end, runLen)
		if closing < 0 {
			index += runLen
			continue
		}
		ranges = append(ranges, Range{Start: index, End: closing + runLen})
		index = closing + runLen
	}

	return ranges
}

func backtickRun(text string, index, end int) int {
	length := 0
	for index+length < end && text[index+length] == '`' {
		length++
	}

	return length
}

func findBacktickRun(text string, from, end, want int) int {
	index := from
	for index < end {
		if text[index] != '`' {
			index++
			continue
		}
		length := backtickRun(text, index, end)
		if length == want {
			return index
		}
		index += length
	}

	return -1
}

func inRanges(ranges []Range, start, end int) bool {
	for _, current := range ranges {
		if start < current.End && current.Start < end {
			return true
		}
	}

	return false
}

func overlapsLinks(links []Link, start, end int) bool {
	for _, link := range links {
		if start < link.End && link.Start < end {
			return true
		}
	}

	return false
}

func splitWikilinkInner(inner string) (target, display, heading string) {
	target = inner
	if before, after, found := strings.Cut(target, "|"); found {
		target, display = before, strings.TrimSpace(after)
	}
	if before, after, found := strings.Cut(target, "#"); found {
		target, heading = before, strings.TrimSpace(after)
	}

	return target, display, heading
}

func internalMarkdownTarget(destination string) (string, string, bool) {
	destination = strings.TrimSpace(destination)
	if destination == "" || strings.HasPrefix(destination, "#") || schemePattern.MatchString(destination) {
		return "", "", false
	}

	heading := ""
	if before, after, found := strings.Cut(destination, "#"); found {
		destination, heading = before, after
	}
	if unescaped, err := url.PathUnescape(destination); err == nil {
		destination = unescaped
	}

	extension := strings.ToLower(path.Ext(destination))
	if extension != "" && extension != ".md" {
		return "", "", false
	}

	target := Normalize(destination)
	if target == "" {
		return "", "", false
	}

	return target, heading, true
}

func stringLinks(text string) []string {
	matches := stringWikilinkPattern.FindAllStringSubmatch(text, -1)
	targets := make([]string, 0, len(matches))
	for _, match := range matches {
		target, _, _ := splitWikilinkInner(match[1])
		if normalized := Normalize(target); normalized != "" {
			targets = append(targets, normalized)
		}
	}

	return targets
}
