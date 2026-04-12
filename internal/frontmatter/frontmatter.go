// Package frontmatter provides helpers for parsing and generating YAML frontmatter in markdown documents.
package frontmatter

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// DateLayout matches the KB schema date fields.
	DateLayout = "2006-01-02"
)

// ErrorKind identifies the failure class for frontmatter operations.
type ErrorKind string

const (
	ErrorKindMissingClosingDelimiter ErrorKind = "missing_closing_delimiter"
	ErrorKindInvalidYAML             ErrorKind = "invalid_yaml"
	ErrorKindUnsupportedValue        ErrorKind = "unsupported_value"
)

var yamlLinePattern = regexp.MustCompile(`line (\d+)`)

// Error reports structured parse and generation failures.
type Error struct {
	Kind ErrorKind
	Key  string
	Line int
	Err  error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{string(e.Kind)}
	if e.Key != "" {
		parts = append(parts, "key="+strconv.Quote(e.Key))
	}
	if e.Line > 0 {
		parts = append(parts, fmt.Sprintf("line=%d", e.Line))
	}
	if e.Err != nil {
		parts = append(parts, e.Err.Error())
	}

	return "frontmatter: " + strings.Join(parts, ": ")
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

// Parse extracts YAML frontmatter from markdown content and returns the decoded
// metadata plus the remaining markdown body. When the markdown has no leading
// frontmatter block, Parse returns an empty map and the original body.
func Parse(markdown string) (map[string]any, string, error) {
	if !hasOpeningDelimiter(markdown) {
		return map[string]any{}, markdown, nil
	}

	source, body, found := splitFrontmatter(markdown)
	if !found {
		return nil, "", &Error{
			Kind: ErrorKindMissingClosingDelimiter,
			Err:  fmt.Errorf("missing closing frontmatter delimiter"),
		}
	}

	if strings.TrimSpace(source) == "" {
		return map[string]any{}, body, nil
	}

	var values map[string]any
	if err := yaml.Unmarshal([]byte(source), &values); err != nil {
		return nil, "", &Error{
			Kind: ErrorKindInvalidYAML,
			Line: extractYAMLLine(err),
			Err:  err,
		}
	}

	return normalizeMap(values), body, nil
}

// Generate encodes frontmatter and prepends it to body. When values is empty,
// Generate returns body unchanged.
func Generate(values map[string]any, body string) (string, error) {
	if len(values) == 0 {
		return body, nil
	}

	node, err := buildNode(values, "")
	if err != nil {
		return "", err
	}

	encoded, err := yaml.Marshal(node)
	if err != nil {
		return "", fmt.Errorf("marshal frontmatter: %w", err)
	}

	return "---\n" + strings.TrimRight(string(encoded), "\n") + "\n---\n" + body, nil
}

// GetString returns the string value for key, or the zero value when missing.
func GetString(values map[string]any, key string) string {
	value, exists := values[key]
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

// GetStringSlice returns the string slice value for key, or nil when missing.
func GetStringSlice(values map[string]any, key string) []string {
	value, exists := values[key]
	if !exists {
		return nil
	}

	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...)
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil
			}
			items = append(items, text)
		}
		return items
	default:
		return nil
	}
}

// GetTime returns the time value for key, or the zero time when missing.
func GetTime(values map[string]any, key string) time.Time {
	value, exists := values[key]
	if !exists {
		return time.Time{}
	}

	switch typed := value.(type) {
	case time.Time:
		return typed
	case string:
		for _, layout := range []string{DateLayout, time.RFC3339, time.RFC3339Nano} {
			parsed, err := time.Parse(layout, typed)
			if err == nil {
				return parsed
			}
		}
		return time.Time{}
	default:
		return time.Time{}
	}
}

// GetBool returns the bool value for key, or false when missing.
func GetBool(values map[string]any, key string) bool {
	value, exists := values[key]
	if !exists {
		return false
	}

	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		parsed, err := strconv.ParseBool(typed)
		if err == nil {
			return parsed
		}
		return false
	default:
		return false
	}
}

func hasOpeningDelimiter(markdown string) bool {
	return strings.HasPrefix(markdown, "---\n") || strings.HasPrefix(markdown, "---\r\n") || markdown == "---"
}

func splitFrontmatter(markdown string) (string, string, bool) {
	start := firstLineEnd(markdown)
	if start < 0 {
		return "", "", false
	}

	lineStart := start
	for lineStart <= len(markdown) {
		lineEnd, next := lineBounds(markdown, lineStart)
		if strings.TrimSuffix(markdown[lineStart:lineEnd], "\r") == "---" {
			return markdown[start:lineStart], markdown[next:], true
		}
		if next == len(markdown) {
			break
		}
		lineStart = next
	}

	return "", "", false
}

func firstLineEnd(markdown string) int {
	if strings.HasPrefix(markdown, "---\r\n") {
		return len("---\r\n")
	}
	if strings.HasPrefix(markdown, "---\n") {
		return len("---\n")
	}
	return -1
}

func lineBounds(markdown string, start int) (int, int) {
	if start >= len(markdown) {
		return len(markdown), len(markdown)
	}

	index := strings.IndexByte(markdown[start:], '\n')
	if index < 0 {
		return len(markdown), len(markdown)
	}

	lineEnd := start + index
	return lineEnd, lineEnd + 1
}

func normalizeMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return map[string]any{}
	}

	normalized := make(map[string]any, len(values))
	for key, value := range values {
		normalized[key] = normalizeValue(value)
	}

	return normalized
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return normalizeMap(typed)
	case map[any]any:
		normalized := make(map[string]any, len(typed))
		for key, nested := range typed {
			normalized[fmt.Sprint(key)] = normalizeValue(nested)
		}
		return normalized
	case []any:
		values := make([]any, len(typed))
		texts := make([]string, 0, len(typed))
		allStrings := true

		for index, item := range typed {
			normalized := normalizeValue(item)
			values[index] = normalized

			text, ok := normalized.(string)
			if !ok {
				allStrings = false
				continue
			}
			texts = append(texts, text)
		}

		if allStrings {
			return texts
		}

		return values
	default:
		return typed
	}
}

func buildNode(value any, key string) (*yaml.Node, error) {
	if value == nil {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	}

	switch typed := value.(type) {
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: typed}, nil
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(typed)}, nil
	case time.Time:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!timestamp", Value: formatTime(typed)}, nil
	case map[string]any:
		return buildMappingNode(typed)
	case []string:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range typed {
			child, err := buildNode(item, key)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, child)
		}
		return node, nil
	case []any:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range typed {
			child, err := buildNode(item, key)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, child)
		}
		return node, nil
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	}

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(rv.Int(), 10)}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatUint(rv.Uint(), 10)}, nil
	case reflect.Float32, reflect.Float64:
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!float",
			Value: strconv.FormatFloat(rv.Float(), 'f', -1, 64),
		}, nil
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, &Error{
				Kind: ErrorKindUnsupportedValue,
				Key:  key,
				Err:  fmt.Errorf("map keys must be strings, got %s", rv.Type().Key()),
			}
		}

		values := make(map[string]any, rv.Len())
		for _, mapKey := range rv.MapKeys() {
			values[mapKey.String()] = rv.MapIndex(mapKey).Interface()
		}
		return buildMappingNode(values)
	case reflect.Slice, reflect.Array:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for index := 0; index < rv.Len(); index++ {
			child, err := buildNode(rv.Index(index).Interface(), key)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, child)
		}
		return node, nil
	default:
		return nil, &Error{
			Kind: ErrorKindUnsupportedValue,
			Key:  key,
			Err:  fmt.Errorf("unsupported value type %T", value),
		}
	}
}

func buildMappingNode(values map[string]any) (*yaml.Node, error) {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		valueNode, err := buildNode(values[key], key)
		if err != nil {
			return nil, err
		}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
			valueNode,
		)
	}

	return node, nil
}

func formatTime(value time.Time) string {
	if value.Hour() == 0 && value.Minute() == 0 && value.Second() == 0 && value.Nanosecond() == 0 {
		return value.Format(DateLayout)
	}

	return value.Format(time.RFC3339Nano)
}

func extractYAMLLine(err error) int {
	matches := yamlLinePattern.FindStringSubmatch(err.Error())
	if len(matches) != 2 {
		return 0
	}

	line, convErr := strconv.Atoi(matches[1])
	if convErr != nil {
		return 0
	}

	return line
}
