package review

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/frontmatter"
)

// Match keys accepted by ImportLabels.
const (
	MatchURL   = "url"
	MatchPath  = "path"
	MatchDOI   = "doi"
	MatchPMCID = "pmcid"
)

// Label identities (purpose and question ids) used by imports and calibration.
const (
	PurposeRelevance   = string(decisions.PurposeRelevance)
	QuestionRole       = "role"
	PurposeLink        = string(decisions.PurposeLink)
	QuestionShouldLink = "should_link"
	PurposeQuality     = string(decisions.PurposeQuality)
)

// OriginImportPrefix prefixes the origin of labels imported from a file.
const OriginImportPrefix = "import:"

// bodyHeadBytes bounds the body prefix searched for a DOI or PMCID.
const bodyHeadBytes = 4096

var (
	pmcidPattern = regexp.MustCompile(`(?i)\bPMC\d+\b`)
	doiPattern   = regexp.MustCompile(`(?i)\b10\.\d{4,9}/[^\s"'<>)\]]+`)
)

// DefaultKeepValues are the decision values read as "keep" when
// ImportOptions.KeepValues is empty.
var DefaultKeepValues = []string{"include", "true"}

// ImportOptions configures ImportLabels.
type ImportOptions struct {
	// From is the JSONL, JSON (array of objects) or CSV file to import.
	From string
	// IDField is the column holding the source identifier.
	IDField string
	// Match is how IDField is matched to a topic source: url|path|doi|pmcid.
	Match string
	// DecisionField is the column holding the screening decision.
	DecisionField string
	// KeepValues are the decision values meaning keep/relevant
	// (case-insensitive); every other value is a negative label.
	KeepValues []string
	// DecidedByField names the column copied to Label.DecidedBy; when empty
	// `screening_stage` or `decided_by` are used when present.
	DecidedByField string
	// Corpus is the loaded topic; loaded from the topic root when nil.
	Corpus *corpus.Corpus
	// Now defaults to time.Now.
	Now func() time.Time
}

// ImportReport summarizes an import.
type ImportReport struct {
	// Rows is the number of data rows read.
	Rows int
	// Matched counts rows matched to a topic source.
	Matched int
	// Unmatched counts rows with no matching source.
	Unmatched int
	// UnmatchedIDs lists the identifiers of unmatched rows in file order.
	UnmatchedIDs []string
	// Positive and Negative count the labels written by this import.
	Positive int
	Negative int
	// Duplicates counts matched rows already imported from the same file.
	Duplicates int
	// Origin is the origin recorded on every label.
	Origin string
}

// ImportLabels imports relevance labels from a topic's own screening or
// curation file (spec §12.2). Each row is matched to a source by the chosen
// key; a decision in KeepValues is a positive (keep) label, anything else a
// negative one, with purpose `relevance`, question `role` and origin
// `import:<file base name>@<sha256 of the file>`. Re-importing the same file
// writes nothing new; when one subject appears on several rows the last row
// wins. Every row must contain the selected ID and decision fields; missing
// fields reject the import before any labels are written.
func ImportLabels(topicRoot string, opts ImportOptions) (ImportReport, error) {
	if strings.TrimSpace(opts.IDField) == "" || strings.TrimSpace(opts.DecisionField) == "" {
		return ImportReport{}, errors.New("review: import-labels needs --id-field and --decision-field")
	}
	match := strings.ToLower(strings.TrimSpace(opts.Match))
	if !slices.Contains([]string{MatchURL, MatchPath, MatchDOI, MatchPMCID}, match) {
		return ImportReport{}, fmt.Errorf("review: --match must be url|path|doi|pmcid, got %q", opts.Match)
	}
	data, err := os.ReadFile(opts.From)
	if err != nil {
		return ImportReport{}, fmt.Errorf("review: read %s: %w", opts.From, err)
	}
	rows, err := parseRows(opts.From, data)
	if err != nil {
		return ImportReport{}, err
	}
	for index, row := range rows {
		for _, field := range []string{opts.IDField, opts.DecisionField} {
			if _, present := row[field]; !present {
				return ImportReport{}, fmt.Errorf("review: %s row %d: missing field %q", opts.From, index+1, field)
			}
		}
	}
	c := opts.Corpus
	if c == nil {
		if c, err = corpus.Load(topicRoot, corpus.LoadOptions{}); err != nil {
			return ImportReport{}, fmt.Errorf("review: %w", err)
		}
	}
	keep := opts.KeepValues
	if len(keep) == 0 {
		keep = DefaultKeepValues
	}

	sum := sha256.Sum256(data)
	report := ImportReport{Rows: len(rows), Origin: OriginImportPrefix + filepath.Base(opts.From) + "@" + hex.EncodeToString(sum[:])}
	matcher := newSourceMatcher(c, match)

	bySubject := map[string]Label{}
	order := make([]string, 0)
	for _, row := range rows {
		id := strings.TrimSpace(row[opts.IDField])
		subject := matcher.match(id)
		if subject == "" {
			report.Unmatched++
			report.UnmatchedIDs = append(report.UnmatchedIDs, id)
			continue
		}
		report.Matched++
		verdict := VerdictNegative
		if containsFold(keep, row[opts.DecisionField]) {
			verdict = VerdictPositive
		}
		if _, seen := bySubject[subject]; !seen {
			order = append(order, subject)
		}
		bySubject[subject] = Label{
			Subject:   subject,
			Purpose:   PurposeRelevance,
			Question:  QuestionRole,
			Verdict:   verdict,
			Origin:    report.Origin,
			DecidedBy: decidedBy(row, opts.DecidedByField),
		}
	}

	existing, err := LoadLabels(topicRoot)
	if err != nil {
		return report, err
	}
	seen := map[string]bool{}
	for _, label := range existing {
		seen[label.Origin+"\x00"+label.Subject+"\x00"+label.Purpose] = true
	}
	store := Open(topicRoot, opts.Now)
	for _, subject := range order {
		label := bySubject[subject]
		if seen[label.Origin+"\x00"+label.Subject+"\x00"+label.Purpose] {
			report.Duplicates++
			continue
		}
		if err := store.AddLabel(label); err != nil {
			return report, err
		}
		if label.Verdict == VerdictPositive {
			report.Positive++
		} else {
			report.Negative++
		}
	}
	return report, nil
}

func decidedBy(row map[string]string, field string) string {
	if field = strings.TrimSpace(field); field != "" {
		return strings.TrimSpace(row[field])
	}
	for _, candidate := range []string{"screening_stage", "decided_by"} {
		if value := strings.TrimSpace(row[candidate]); value != "" {
			return value
		}
	}
	return ""
}

func containsFold(values []string, value string) bool {
	value = strings.TrimSpace(value)
	return slices.ContainsFunc(values, func(candidate string) bool {
		return strings.EqualFold(strings.TrimSpace(candidate), value)
	})
}

// parseRows reads a CSV (by extension), a JSON array of objects (.json) or
// JSONL file into string maps.
func parseRows(name string, data []byte) ([]map[string]string, error) {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".csv":
		return parseCSV(name, data)
	case ".json":
		var objects []map[string]any
		if err := json.Unmarshal(data, &objects); err != nil {
			return nil, fmt.Errorf("review: parse %s: %w", name, err)
		}
		rows := make([]map[string]string, 0, len(objects))
		for _, object := range objects {
			rows = append(rows, stringifyRow(object))
		}
		return rows, nil
	default:
		rows := make([]map[string]string, 0)
		for number, line := range bytes.Split(data, []byte("\n")) {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			var object map[string]any
			if err := json.Unmarshal(line, &object); err != nil {
				return nil, fmt.Errorf("review: parse %s line %d: %w", name, number+1, err)
			}
			rows = append(rows, stringifyRow(object))
		}
		return rows, nil
	}
}

func parseCSV(name string, data []byte) ([]map[string]string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if errors.Is(err, io.EOF) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("review: parse %s header: %w", name, err)
	}
	for index := range header {
		header[index] = strings.TrimSpace(strings.TrimPrefix(header[index], "\ufeff"))
	}
	rows := make([]map[string]string, 0)
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			return rows, nil
		}
		if err != nil {
			return nil, fmt.Errorf("review: parse %s: %w", name, err)
		}
		row := make(map[string]string, len(header))
		for index, column := range header {
			if index < len(record) {
				row[column] = record[index]
			}
		}
		rows = append(rows, row)
	}
}

func stringifyRow(object map[string]any) map[string]string {
	row := make(map[string]string, len(object))
	for key, value := range object {
		switch typed := value.(type) {
		case nil:
			row[key] = ""
		case string:
			row[key] = typed
		case bool:
			row[key] = strconv.FormatBool(typed)
		case float64:
			row[key] = strconv.FormatFloat(typed, 'f', -1, 64)
		default:
			encoded, _ := json.Marshal(typed)
			row[key] = string(encoded)
		}
	}
	return row
}

// sourceMatcher maps identifiers to topic source paths.
type sourceMatcher struct {
	kind  string
	c     *corpus.Corpus
	index map[string]string
}

func newSourceMatcher(c *corpus.Corpus, kind string) sourceMatcher {
	m := sourceMatcher{kind: kind, c: c, index: map[string]string{}}
	add := func(key, subject string) {
		if key == "" {
			return
		}
		if _, taken := m.index[key]; !taken {
			m.index[key] = subject
		}
	}
	sources := c.Sources()
	sort.Slice(sources, func(i, j int) bool { return sources[i].Path < sources[j].Path })
	for _, doc := range sources {
		switch kind {
		case MatchURL:
			add(NormalizeURL(doc.SourceURL()), doc.Path)
		case MatchDOI, MatchPMCID:
			key := kind
			add(normalizeID(kind, frontmatter.GetString(doc.Frontmatter, key)), doc.Path)
			head := doc.Body
			if len(head) > bodyHeadBytes {
				head = head[:bodyHeadBytes]
			}
			for _, text := range []string{doc.SourceURL(), head} {
				for _, found := range extractIDs(kind, text) {
					add(found, doc.Path)
				}
			}
		}
	}
	return m
}

func (m sourceMatcher) match(id string) string {
	if id == "" {
		return ""
	}
	switch m.kind {
	case MatchPath:
		if doc := m.c.ByPath(strings.TrimPrefix(filepath.ToSlash(id), "./")); doc != nil && doc.Kind == corpus.KindSource {
			return doc.Path
		}
		return ""
	case MatchURL:
		return m.index[NormalizeURL(id)]
	default:
		return m.index[normalizeID(m.kind, id)]
	}
}

func extractIDs(kind, text string) []string {
	pattern := doiPattern
	if kind == MatchPMCID {
		pattern = pmcidPattern
	}
	found := pattern.FindAllString(text, -1)
	out := make([]string, 0, len(found))
	for _, value := range found {
		out = append(out, normalizeID(kind, value))
	}
	return out
}

func normalizeID(kind, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if kind == MatchPMCID {
		upper := strings.ToUpper(value)
		if !strings.HasPrefix(upper, "PMC") {
			upper = "PMC" + upper
		}
		return upper
	}
	lower := strings.ToLower(value)
	for _, prefix := range []string{"https://doi.org/", "http://doi.org/", "https://dx.doi.org/", "http://dx.doi.org/", "doi:"} {
		lower = strings.TrimPrefix(lower, prefix)
	}
	return strings.TrimRight(strings.TrimSpace(lower), ".,;")
}

// NormalizeURL reduces a URL to a comparison key: scheme, `www.`, default
// ports, fragment, trailing slash and utm_* parameters are dropped, the host
// is lower-cased and the remaining query parameters are sorted.
func NormalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return strings.TrimRight(strings.ToLower(raw), "/")
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	if port := parsed.Port(); port != "" && port != "80" && port != "443" {
		host += ":" + port
	}
	query := parsed.Query()
	for key := range query {
		if strings.HasPrefix(strings.ToLower(key), "utm_") {
			query.Del(key)
		}
	}
	key := host + strings.TrimRight(parsed.EscapedPath(), "/")
	if encoded := query.Encode(); encoded != "" {
		key += "?" + encoded
	}
	return key
}
