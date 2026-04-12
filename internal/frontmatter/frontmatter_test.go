package frontmatter_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/user/kb/internal/frontmatter"
	"gopkg.in/yaml.v3"
)

type textValue string

func (v textValue) String() string {
	return string(v)
}

func TestParseValidFrontmatterWithSupportedTypes(t *testing.T) {
	t.Parallel()

	markdown := strings.Join([]string{
		"---",
		"title: Example Source",
		"issues_found: 3",
		"score: 1.5",
		"published: true",
		"scraped: 2026-04-11",
		"tags:",
		"  - kb",
		"  - raw",
		"---",
		"",
		"# Example",
	}, "\n")

	values, body, err := frontmatter.Parse(markdown)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if got := values["title"]; got != "Example Source" {
		t.Fatalf("title = %#v, want Example Source", got)
	}
	if got := values["issues_found"]; got != 3 {
		t.Fatalf("issues_found = %#v, want 3", got)
	}
	if got := values["score"]; got != 1.5 {
		t.Fatalf("score = %#v, want 1.5", got)
	}
	if got := values["published"]; got != true {
		t.Fatalf("published = %#v, want true", got)
	}

	scraped, ok := values["scraped"].(time.Time)
	if !ok {
		t.Fatalf("scraped type = %T, want time.Time", values["scraped"])
	}
	if scraped.Format(frontmatter.DateLayout) != "2026-04-11" {
		t.Fatalf("scraped = %s, want 2026-04-11", scraped.Format(frontmatter.DateLayout))
	}

	tags, ok := values["tags"].([]string)
	if !ok {
		t.Fatalf("tags type = %T, want []string", values["tags"])
	}
	if !reflect.DeepEqual(tags, []string{"kb", "raw"}) {
		t.Fatalf("tags = %#v, want [kb raw]", tags)
	}

	if body != "\n# Example" {
		t.Fatalf("body = %q, want %q", body, "\n# Example")
	}
}

func TestParseMarkdownWithoutFrontmatterReturnsOriginalBody(t *testing.T) {
	t.Parallel()

	body := "# No Frontmatter\n"
	values, gotBody, err := frontmatter.Parse(body)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("expected empty map, got %#v", values)
	}
	if gotBody != body {
		t.Fatalf("body = %q, want %q", gotBody, body)
	}
}

func TestParseEmptyFrontmatterReturnsEmptyMap(t *testing.T) {
	t.Parallel()

	values, body, err := frontmatter.Parse("---\n---\n# Body\n")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("expected empty map, got %#v", values)
	}
	if body != "# Body\n" {
		t.Fatalf("body = %q, want %q", body, "# Body\n")
	}
}

func TestParseMissingClosingDelimiterReturnsStructuredError(t *testing.T) {
	t.Parallel()

	_, _, err := frontmatter.Parse(strings.Join([]string{
		"---",
		"title: Broken",
		"# Body",
	}, "\n"))
	if err == nil {
		t.Fatal("expected Parse to fail")
	}

	var fmErr *frontmatter.Error
	if !errors.As(err, &fmErr) {
		t.Fatalf("expected *frontmatter.Error, got %T", err)
	}
	if fmErr.Kind != frontmatter.ErrorKindMissingClosingDelimiter {
		t.Fatalf("error kind = %q, want %q", fmErr.Kind, frontmatter.ErrorKindMissingClosingDelimiter)
	}
}

func TestParseInvalidYAMLReturnsStructuredErrorWithLine(t *testing.T) {
	t.Parallel()

	_, _, err := frontmatter.Parse(strings.Join([]string{
		"---",
		"title: Example",
		`broken: "unterminated`,
		"---",
		"# Body",
	}, "\n"))
	if err == nil {
		t.Fatal("expected Parse to fail")
	}

	var fmErr *frontmatter.Error
	if !errors.As(err, &fmErr) {
		t.Fatalf("expected *frontmatter.Error, got %T", err)
	}
	if fmErr.Kind != frontmatter.ErrorKindInvalidYAML {
		t.Fatalf("error kind = %q, want %q", fmErr.Kind, frontmatter.ErrorKindInvalidYAML)
	}
	if fmErr.Line != 2 {
		t.Fatalf("line = %d, want 2", fmErr.Line)
	}
}

func TestGenerateProducesValidDelimitedYAMLWithSortedKeys(t *testing.T) {
	t.Parallel()

	values := map[string]any{
		"zeta":  "last",
		"alpha": "first",
		"tags":  []string{"kb", "raw"},
		"date":  mustDate(t, "2026-04-11"),
	}

	rendered, err := frontmatter.Generate(values, "body")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if !strings.HasPrefix(rendered, "---\n") {
		t.Fatalf("generated output missing opening delimiter:\n%s", rendered)
	}
	if !strings.Contains(rendered, "\n---\nbody") {
		t.Fatalf("generated output missing closing delimiter or body:\n%s", rendered)
	}

	alphaIndex := strings.Index(rendered, "alpha:")
	dateIndex := strings.Index(rendered, "date:")
	tagsIndex := strings.Index(rendered, "tags:")
	zetaIndex := strings.Index(rendered, "zeta:")
	if alphaIndex >= dateIndex || dateIndex >= tagsIndex || tagsIndex >= zetaIndex {
		t.Fatalf("frontmatter keys are not sorted:\n%s", rendered)
	}

	parsed, body := mustDecodeGeneratedFrontmatter(t, rendered)
	if got := parsed["alpha"]; got != "first" {
		t.Fatalf("alpha = %#v, want first", got)
	}
	if got := parsed["zeta"]; got != "last" {
		t.Fatalf("zeta = %#v, want last", got)
	}
	if body != "body" {
		t.Fatalf("body = %q, want %q", body, "body")
	}
}

func TestGenerateEmptyMapProducesNoFrontmatterPrefix(t *testing.T) {
	t.Parallel()

	body := "\n# Body\n"
	rendered, err := frontmatter.Generate(map[string]any{}, body)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if rendered != body {
		t.Fatalf("rendered = %q, want %q", rendered, body)
	}
}

func TestRoundTripSchemaVariants(t *testing.T) {
	t.Parallel()

	dateA := mustDate(t, "2026-04-11")
	dateB := mustDate(t, "2026-04-12")

	cases := map[string]map[string]any{
		"source-article": {
			"title":       "Descriptive Title",
			"type":        "source",
			"stage":       "raw",
			"domain":      "kb",
			"source_kind": "article",
			"source_url":  "https://example.com/article",
			"scraped":     dateA,
			"tags":        []string{"kb", "raw", "article"},
		},
		"source-bookmarks": {
			"title":       "KB Bookmarks Retrieval",
			"type":        "source",
			"stage":       "raw",
			"domain":      "kb",
			"source_kind": "bookmark-cluster",
			"status":      "seeded",
			"created":     dateA,
			"updated":     dateB,
			"source_urls": []string{
				"https://twitter.com/example/1",
				"https://twitter.com/example/2",
			},
			"tags": []string{"kb", "bookmarks", "raw"},
		},
		"wiki": {
			"title":   "Retriever Design",
			"type":    "wiki",
			"stage":   "compiled",
			"domain":  "kb",
			"tags":    []string{"kb", "wiki", "retrieval"},
			"created": dateA,
			"updated": dateB,
			"sources": []string{"[[Source File Name]]", "[[Another Source]]"},
		},
		"output-query": {
			"title":       "Answer Retrieval Questions",
			"type":        "output",
			"stage":       "query",
			"domain":      "kb",
			"tags":        []string{"kb", "output", "query"},
			"created":     dateA,
			"updated":     dateB,
			"informed_by": []string{"[[Wiki Article 1]]", "[[Wiki Article 2]]"},
		},
		"output-lint-report": {
			"title":        "Lint Report 2026-04-11",
			"type":         "output",
			"stage":        "lint-report",
			"domain":       "kb",
			"tags":         []string{"kb", "output", "lint-report"},
			"created":      dateA,
			"issues_found": 2,
			"issues_fixed": 1,
		},
		"index": {
			"title":   "Dashboard",
			"type":    "index",
			"domain":  "kb",
			"updated": dateB,
		},
	}

	for name, values := range cases {
		values := values
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			body := "\n# Body\n"
			rendered, err := frontmatter.Generate(values, body)
			if err != nil {
				t.Fatalf("Generate returned error: %v", err)
			}

			parsed, parsedBody, err := frontmatter.Parse(rendered)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}

			if !reflect.DeepEqual(parsed, values) {
				t.Fatalf("parsed values = %#v, want %#v", parsed, values)
			}
			if parsedBody != body {
				t.Fatalf("body = %q, want %q", parsedBody, body)
			}
		})
	}
}

func TestAccessorsReturnValuesAndZeroValues(t *testing.T) {
	t.Parallel()

	date := mustDate(t, "2026-04-11")
	values := map[string]any{
		"title":        "Dashboard",
		"tags":         []string{"kb", "index"},
		"updated":      date,
		"published":    true,
		"publishedRaw": "true",
		"dateRaw":      "2026-04-12",
	}

	if got := frontmatter.GetString(values, "title"); got != "Dashboard" {
		t.Fatalf("GetString(title) = %q, want Dashboard", got)
	}
	if got := frontmatter.GetString(values, "missing"); got != "" {
		t.Fatalf("GetString(missing) = %q, want empty string", got)
	}

	if got := frontmatter.GetStringSlice(values, "tags"); !reflect.DeepEqual(got, []string{"kb", "index"}) {
		t.Fatalf("GetStringSlice(tags) = %#v, want [kb index]", got)
	}
	if got := frontmatter.GetStringSlice(values, "missing"); got != nil {
		t.Fatalf("GetStringSlice(missing) = %#v, want nil", got)
	}

	if got := frontmatter.GetTime(values, "updated"); !got.Equal(date) {
		t.Fatalf("GetTime(updated) = %s, want %s", got, date)
	}
	if got := frontmatter.GetTime(values, "dateRaw"); got.Format(frontmatter.DateLayout) != "2026-04-12" {
		t.Fatalf("GetTime(dateRaw) = %s, want 2026-04-12", got.Format(frontmatter.DateLayout))
	}
	if got := frontmatter.GetTime(values, "missing"); !got.IsZero() {
		t.Fatalf("GetTime(missing) = %s, want zero time", got)
	}

	if got := frontmatter.GetBool(values, "published"); !got {
		t.Fatal("GetBool(published) = false, want true")
	}
	if got := frontmatter.GetBool(values, "publishedRaw"); !got {
		t.Fatal("GetBool(publishedRaw) = false, want true")
	}
	if got := frontmatter.GetBool(values, "missing"); got {
		t.Fatal("GetBool(missing) = true, want false")
	}
}

func TestGenerateRejectsUnsupportedValueTypes(t *testing.T) {
	t.Parallel()

	_, err := frontmatter.Generate(map[string]any{
		"bad": struct{ Name string }{Name: "unsupported"},
	}, "")
	if err == nil {
		t.Fatal("expected Generate to fail")
	}

	var fmErr *frontmatter.Error
	if !errors.As(err, &fmErr) {
		t.Fatalf("expected *frontmatter.Error, got %T", err)
	}
	if fmErr.Kind != frontmatter.ErrorKindUnsupportedValue {
		t.Fatalf("error kind = %q, want %q", fmErr.Kind, frontmatter.ErrorKindUnsupportedValue)
	}
	if fmErr.Key != "bad" {
		t.Fatalf("error key = %q, want bad", fmErr.Key)
	}
}

func TestParseSupportsCRLFAndNoTrailingBodyNewline(t *testing.T) {
	t.Parallel()

	values, body, err := frontmatter.Parse("---\r\ntitle: Example\r\n---")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got := values["title"]; got != "Example" {
		t.Fatalf("title = %#v, want Example", got)
	}
	if body != "" {
		t.Fatalf("body = %q, want empty", body)
	}
}

func TestParseNormalizesNestedMapsAndMixedSlices(t *testing.T) {
	t.Parallel()

	values, _, err := frontmatter.Parse(strings.Join([]string{
		"---",
		"meta:",
		"  source: article",
		"  count: 2",
		"mixed:",
		"  - alpha",
		"  - 2",
		"---",
	}, "\n"))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	meta, ok := values["meta"].(map[string]any)
	if !ok {
		t.Fatalf("meta type = %T, want map[string]any", values["meta"])
	}
	if got := meta["source"]; got != "article" {
		t.Fatalf("meta[source] = %#v, want article", got)
	}
	if got := meta["count"]; got != 2 {
		t.Fatalf("meta[count] = %#v, want 2", got)
	}

	mixed, ok := values["mixed"].([]any)
	if !ok {
		t.Fatalf("mixed type = %T, want []any", values["mixed"])
	}
	if !reflect.DeepEqual(mixed, []any{"alpha", 2}) {
		t.Fatalf("mixed = %#v, want [alpha 2]", mixed)
	}
}

func TestGenerateRoundTripsGenericValueShapes(t *testing.T) {
	t.Parallel()

	type customMap map[string]any

	values := map[string]any{
		"count":    uint(4),
		"ratio":    float32(1.25),
		"items":    []int{1, 2, 3},
		"mixed":    []any{"alpha", 2, true},
		"metadata": customMap{"kind": "report", "enabled": true},
		"optional": nil,
		"updated":  time.Date(2026, 4, 11, 12, 30, 0, 0, time.UTC),
	}

	rendered, err := frontmatter.Generate(values, "")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	parsed, body, err := frontmatter.Parse(rendered)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if body != "" {
		t.Fatalf("body = %q, want empty", body)
	}

	if got := parsed["count"]; got != 4 {
		t.Fatalf("count = %#v, want 4", got)
	}
	if got := parsed["ratio"]; got != 1.25 {
		t.Fatalf("ratio = %#v, want 1.25", got)
	}
	if got := parsed["items"]; !reflect.DeepEqual(got, []any{1, 2, 3}) {
		t.Fatalf("items = %#v, want [1 2 3]", got)
	}
	if got := parsed["mixed"]; !reflect.DeepEqual(got, []any{"alpha", 2, true}) {
		t.Fatalf("mixed = %#v, want [alpha 2 true]", got)
	}

	metadata, ok := parsed["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata type = %T, want map[string]any", parsed["metadata"])
	}
	if got := metadata["kind"]; got != "report" {
		t.Fatalf("metadata[kind] = %#v, want report", got)
	}
	if got := metadata["enabled"]; got != true {
		t.Fatalf("metadata[enabled] = %#v, want true", got)
	}
	if got := parsed["optional"]; got != nil {
		t.Fatalf("optional = %#v, want nil", got)
	}

	updated, ok := parsed["updated"].(time.Time)
	if !ok {
		t.Fatalf("updated type = %T, want time.Time", parsed["updated"])
	}
	if updated.Format(time.RFC3339) != "2026-04-11T12:30:00Z" {
		t.Fatalf("updated = %s, want 2026-04-11T12:30:00Z", updated.Format(time.RFC3339))
	}
}

func TestGenerateRejectsMapsWithNonStringKeys(t *testing.T) {
	t.Parallel()

	_, err := frontmatter.Generate(map[string]any{
		"badmap": map[int]string{1: "x"},
	}, "")
	if err == nil {
		t.Fatal("expected Generate to fail")
	}

	var fmErr *frontmatter.Error
	if !errors.As(err, &fmErr) {
		t.Fatalf("expected *frontmatter.Error, got %T", err)
	}
	if fmErr.Kind != frontmatter.ErrorKindUnsupportedValue {
		t.Fatalf("error kind = %q, want %q", fmErr.Kind, frontmatter.ErrorKindUnsupportedValue)
	}
	if fmErr.Key != "badmap" {
		t.Fatalf("error key = %q, want badmap", fmErr.Key)
	}
}

func TestAccessorsHandleCoercionsAndInvalidShapes(t *testing.T) {
	t.Parallel()

	values := map[string]any{
		"name":      textValue("stringer"),
		"sliceAny":  []any{"a", "b"},
		"sliceBad":  []any{"a", 1},
		"timeText":  "2026-04-11T12:30:00Z",
		"boolBad":   "not-bool",
		"boolFalse": "false",
	}

	if got := frontmatter.GetString(values, "name"); got != "stringer" {
		t.Fatalf("GetString(name) = %q, want stringer", got)
	}
	if got := frontmatter.GetStringSlice(values, "sliceAny"); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("GetStringSlice(sliceAny) = %#v, want [a b]", got)
	}
	if got := frontmatter.GetStringSlice(values, "sliceBad"); got != nil {
		t.Fatalf("GetStringSlice(sliceBad) = %#v, want nil", got)
	}
	if got := frontmatter.GetTime(values, "timeText"); got.Format(time.RFC3339) != "2026-04-11T12:30:00Z" {
		t.Fatalf("GetTime(timeText) = %s, want 2026-04-11T12:30:00Z", got.Format(time.RFC3339))
	}
	if got := frontmatter.GetBool(values, "boolBad"); got {
		t.Fatal("GetBool(boolBad) = true, want false")
	}
	if got := frontmatter.GetBool(values, "boolFalse"); got {
		t.Fatal("GetBool(boolFalse) = true, want false")
	}
}

func TestErrorStringAndUnwrap(t *testing.T) {
	t.Parallel()

	root := errors.New("boom")
	err := &frontmatter.Error{
		Kind: frontmatter.ErrorKindInvalidYAML,
		Key:  "title",
		Line: 3,
		Err:  root,
	}

	if !errors.Is(err, root) {
		t.Fatal("expected wrapped error to be discoverable with errors.Is")
	}
	message := err.Error()
	if !strings.Contains(message, "invalid_yaml") || !strings.Contains(message, "key=\"title\"") || !strings.Contains(message, "line=3") {
		t.Fatalf("unexpected error string %q", message)
	}
}

func mustDate(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(frontmatter.DateLayout, value)
	if err != nil {
		t.Fatalf("parse date %q: %v", value, err)
	}

	return parsed
}

func mustDecodeGeneratedFrontmatter(t *testing.T, markdown string) (map[string]any, string) {
	t.Helper()

	if !strings.HasPrefix(markdown, "---\n") {
		t.Fatalf("missing opening delimiter:\n%s", markdown)
	}

	remainder := strings.TrimPrefix(markdown, "---\n")
	index := strings.Index(remainder, "\n---\n")
	if index < 0 {
		t.Fatalf("missing closing delimiter:\n%s", markdown)
	}

	var parsed map[string]any
	if err := yaml.Unmarshal([]byte(remainder[:index]), &parsed); err != nil {
		t.Fatalf("generated YAML is invalid: %v", err)
	}

	return parsed, remainder[index+5:]
}
