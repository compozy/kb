package frontmatter

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// KeyUpdate describes one top-level frontmatter key change for EditKeys.
// When Delete is true the key block is removed and Value is ignored.
type KeyUpdate struct {
	Key    string
	Value  any
	Delete bool
}

// KeyBlock is the raw text of one top-level frontmatter key: the key line
// plus its continuation lines (indented lines, column-0 list items, block
// scalar content), excluding trailing blank lines.
type KeyBlock struct {
	Key  string
	Text string
}

// frontmatterBounds locates the leading frontmatter. frontStart/frontEnd
// delimit the YAML text between the delimiter lines and bodyStart is the
// offset of the first byte after the closing delimiter line. opened reports
// whether markdown starts with an opening delimiter; found reports whether
// the closing delimiter exists.
func frontmatterBounds(markdown string) (frontStart, frontEnd, bodyStart int, opened, found bool) {
	if !hasOpeningDelimiter(markdown) {
		return 0, 0, 0, false, false
	}

	start := firstLineEnd(markdown)
	if start < 0 {
		return 0, 0, 0, true, false
	}

	lineStart := start
	for lineStart <= len(markdown) {
		lineEnd, next := lineBounds(markdown, lineStart)
		if strings.TrimSuffix(markdown[lineStart:lineEnd], "\r") == "---" {
			return start, lineStart, next, true, true
		}
		if next == len(markdown) {
			break
		}
		lineStart = next
	}

	return 0, 0, 0, true, false
}

// Split separates the leading frontmatter YAML text (without the delimiter
// lines) from the body. has is false when markdown has no complete leading
// frontmatter block (no opening delimiter or no closing delimiter); body is
// then the whole input.
func Split(markdown string) (front string, body string, has bool) {
	frontStart, frontEnd, bodyStart, _, found := frontmatterBounds(markdown)
	if !found {
		return "", markdown, false
	}

	return markdown[frontStart:frontEnd], markdown[bodyStart:], true
}

// KeyBlocks returns the raw text of every top-level key in the leading
// frontmatter, in document order. A document without frontmatter has no
// blocks. A missing closing delimiter is an error.
func KeyBlocks(markdown string) ([]KeyBlock, error) {
	frontStart, frontEnd, _, opened, found := frontmatterBounds(markdown)
	if !opened {
		return nil, nil
	}
	if !found {
		return nil, missingClosingDelimiterError()
	}

	segments := splitKeySegments(markdown[frontStart:frontEnd])
	blocks := make([]KeyBlock, 0, len(segments))
	for _, segment := range segments {
		if segment.key == "" {
			continue
		}
		blocks = append(blocks, KeyBlock{Key: segment.key, Text: segment.text})
	}

	return blocks, nil
}

// TopLevelKeyText returns the raw text of the top-level key block for key.
func TopLevelKeyText(markdown, key string) (string, bool) {
	blocks, err := KeyBlocks(markdown)
	if err != nil {
		return "", false
	}
	for _, block := range blocks {
		if block.Key == key {
			return block.Text, true
		}
	}

	return "", false
}

// EditKeys replaces, deletes or appends only the given top-level keys inside
// the leading frontmatter of markdown, without re-serializing any other YAML:
// comments, key order, quoting and indentation of untouched keys are kept
// byte for byte, and the body is never modified. Each updated key is encoded
// on its own with the same rules as Generate. New keys are appended at the end
// of the frontmatter in the order given. `\r\n` line endings are preserved.
// A document without frontmatter gets a new block (unless every update is a
// delete). The existing frontmatter must be valid YAML, and the edited result
// is re-validated before it is returned.
func EditKeys(markdown string, updates []KeyUpdate) (string, error) {
	for _, update := range updates {
		if strings.TrimSpace(update.Key) == "" || update.Key != strings.TrimSpace(update.Key) {
			return "", &Error{
				Kind: ErrorKindUnsupportedValue,
				Key:  update.Key,
				Err:  errors.New("key must be non-empty without surrounding whitespace"),
			}
		}
	}

	frontStart, frontEnd, bodyStart, opened, found := frontmatterBounds(markdown)
	if opened && !found {
		return "", missingClosingDelimiterError()
	}

	if !opened {
		return createFrontmatter(markdown, updates)
	}

	front := markdown[frontStart:frontEnd]
	if err := validateYAML(front); err != nil {
		return "", err
	}

	eol := "\n"
	if strings.HasPrefix(markdown, "---\r\n") {
		eol = "\r\n"
	}

	segments := splitKeySegments(front)
	for _, update := range updates {
		index := slices.IndexFunc(segments, func(segment keySegment) bool { return segment.key == update.Key })
		if update.Delete {
			if index >= 0 {
				segments = slices.Delete(segments, index, index+1)
			}
			continue
		}

		encoded, err := encodeKey(update.Key, update.Value, eol)
		if err != nil {
			return "", err
		}
		if index >= 0 {
			segments[index].text = encoded
			continue
		}

		if len(segments) > 0 {
			last := &segments[len(segments)-1]
			if last.text != "" && !strings.HasSuffix(last.text, "\n") {
				last.text += eol
			}
		}
		segments = append(segments, keySegment{key: update.Key, text: encoded})
	}

	var builder strings.Builder
	for _, segment := range segments {
		builder.WriteString(segment.text)
	}
	newFront := builder.String()
	if err := validateYAML(newFront); err != nil {
		return "", err
	}

	return markdown[:frontStart] + newFront + markdown[frontEnd:bodyStart] + markdown[bodyStart:], nil
}

func createFrontmatter(markdown string, updates []KeyUpdate) (string, error) {
	eol := "\n"
	if index := strings.IndexByte(markdown, '\n'); index > 0 && markdown[index-1] == '\r' {
		eol = "\r\n"
	}

	segments := make([]keySegment, 0, len(updates))
	for _, update := range updates {
		index := slices.IndexFunc(segments, func(segment keySegment) bool { return segment.key == update.Key })
		if update.Delete {
			if index >= 0 {
				segments = slices.Delete(segments, index, index+1)
			}
			continue
		}
		encoded, err := encodeKey(update.Key, update.Value, eol)
		if err != nil {
			return "", err
		}
		if index >= 0 {
			segments[index].text = encoded
			continue
		}
		segments = append(segments, keySegment{key: update.Key, text: encoded})
	}
	if len(segments) == 0 {
		return markdown, nil
	}

	var builder strings.Builder
	for _, segment := range segments {
		builder.WriteString(segment.text)
	}

	return "---" + eol + builder.String() + "---" + eol + markdown, nil
}

func encodeKey(key string, value any, eol string) (string, error) {
	valueNode, err := buildNode(value, key)
	if err != nil {
		return "", err
	}

	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
			valueNode,
		},
	}
	encoded, err := yaml.Marshal(node)
	if err != nil {
		return "", &Error{Kind: ErrorKindUnsupportedValue, Key: key, Err: fmt.Errorf("marshal key: %w", err)}
	}

	text := strings.TrimRight(string(encoded), "\n") + "\n"
	if eol != "\n" {
		text = strings.ReplaceAll(text, "\n", eol)
	}

	return text, nil
}

func validateYAML(front string) error {
	if strings.TrimSpace(front) == "" {
		return nil
	}

	var values map[string]any
	if err := yaml.Unmarshal([]byte(front), &values); err != nil {
		return &Error{Kind: ErrorKindInvalidYAML, Line: extractYAMLLine(err), Err: err}
	}

	return nil
}

func missingClosingDelimiterError() error {
	return &Error{
		Kind: ErrorKindMissingClosingDelimiter,
		Err:  errors.New("missing closing frontmatter delimiter"),
	}
}

// keySegment is a run of raw frontmatter text. Segments with a key hold one
// top-level key block; segments without a key hold comments, blank lines or
// other text between blocks.
type keySegment struct {
	key  string
	text string
}

func splitKeySegments(front string) []keySegment {
	lines := splitLinesKeepEnds(front)
	segments := make([]keySegment, 0, len(lines))

	var filler strings.Builder
	flushFiller := func() {
		if filler.Len() == 0 {
			return
		}
		segments = append(segments, keySegment{text: filler.String()})
		filler.Reset()
	}

	for index := 0; index < len(lines); {
		key, ok := topLevelKey(lines[index])
		if !ok {
			filler.WriteString(lines[index])
			index++
			continue
		}

		end := index + 1
		lastContent := index
		for end < len(lines) {
			line := strings.TrimRight(lines[end], "\r\n")
			if isBlankLine(line) {
				end++
				continue
			}
			if isContinuationLine(line) {
				lastContent = end
				end++
				continue
			}
			break
		}

		flushFiller()
		var block strings.Builder
		for _, line := range lines[index : lastContent+1] {
			block.WriteString(line)
		}
		segments = append(segments, keySegment{key: key, text: block.String()})
		index = lastContent + 1
	}
	flushFiller()

	return segments
}

func splitLinesKeepEnds(text string) []string {
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

func isBlankLine(line string) bool {
	return strings.TrimSpace(line) == ""
}

// isContinuationLine reports whether a non-blank line belongs to the previous
// top-level key: indented content or a column-0 block sequence item.
func isContinuationLine(line string) bool {
	if line == "" {
		return false
	}
	switch line[0] {
	case ' ', '\t':
		return true
	case '-':
		return line == "-" || strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "-\t")
	default:
		return false
	}
}

// topLevelKey parses a column-0 mapping key line (plain, double-quoted or
// single-quoted key followed by ':' and a space or end of line).
func topLevelKey(rawLine string) (string, bool) {
	line := strings.TrimRight(rawLine, "\r\n")
	if line == "" {
		return "", false
	}
	switch line[0] {
	case ' ', '\t', '#', '-', '[', '{', '?', '&', '*', '!', '|', '>', '%', '@', '`':
		return "", false
	}

	if line[0] == '"' || line[0] == '\'' {
		key, rest, ok := parseQuotedKey(line)
		if !ok {
			return "", false
		}
		rest = strings.TrimLeft(rest, " \t")
		if !strings.HasPrefix(rest, ":") || !keyTerminator(rest[1:]) {
			return "", false
		}
		return key, true
	}

	for index := 0; index < len(line); index++ {
		if line[index] != ':' {
			continue
		}
		if !keyTerminator(line[index+1:]) {
			continue
		}
		key := strings.TrimRight(line[:index], " \t")
		if key == "" {
			return "", false
		}
		return key, true
	}

	return "", false
}

func keyTerminator(rest string) bool {
	return rest == "" || rest[0] == ' ' || rest[0] == '\t'
}

func parseQuotedKey(line string) (string, string, bool) {
	quote := line[0]
	if quote == '\'' {
		var builder strings.Builder
		for index := 1; index < len(line); index++ {
			if line[index] != '\'' {
				builder.WriteByte(line[index])
				continue
			}
			if index+1 < len(line) && line[index+1] == '\'' {
				builder.WriteByte('\'')
				index++
				continue
			}
			return builder.String(), line[index+1:], true
		}
		return "", "", false
	}

	for index := 1; index < len(line); index++ {
		switch line[index] {
		case '\\':
			index++
		case '"':
			key, err := strconv.Unquote(line[:index+1])
			if err != nil {
				key = line[1:index]
			}
			return key, line[index+1:], true
		}
	}

	return "", "", false
}
