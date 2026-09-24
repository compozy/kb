// Package quality holds the code quality checks that run before any quality
// question (spec §7 stage 3): properties of the URL, the HTTP status, the
// title and the text length that code decides exactly and for free. The gate
// (ingest) and classify (backfill) both use it.
package quality

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/compozy/kb/internal/corpus"
)

// Flag codes (the `quality` / `triage_reason` values they map to).
const (
	NotAnArticle = "not_an_article"
	Thin         = "thin"
	ErrorPage    = "error_page"
)

// MinWords is the body length, in words after removing navigation-like
// lines, under which a document is flagged thin.
const MinWords = 300

// HostRepeatThreshold is the number of same-host sources a line must appear
// on to count as site chrome.
const HostRepeatThreshold = 3

// Flag is one code quality finding.
type Flag struct {
	// Code is NotAnArticle, Thin or ErrorPage.
	Code string `json:"code"`
	// Detail says which rule fired, for run summaries and review evidence.
	Detail string `json:"detail"`
}

// Input is everything the checks read. Empty fields disable the checks that
// need them (a document without a URL is never not_an_article by URL).
type Input struct {
	Title string
	Body  string
	// SourceURL is the stored source URL, used as the final URL when FinalURL
	// is empty.
	SourceURL string
	// FinalURL is the URL after redirects, as reported by the fetcher.
	FinalURL string
	// RequestedURL is the URL the fetch was asked for.
	RequestedURL string
	// StatusCode is the HTTP status reported by the fetcher (0 = unknown).
	StatusCode int
	// SiteName is the site name reported by the fetcher; when empty it is
	// derived from the host (first label without `www.`).
	SiteName string
	// HostLines maps a normalized line (NormalizeLine) to the number of
	// same-host sources it appears on (see HostLineCounts); nil disables the
	// repeated-line rule.
	HostLines map[string]int
	// SkipThin disables the thin rule (for inputs whose length is not a
	// capture property, such as transcripts or local notes).
	SkipThin bool
}

var errorTitlePattern = regexp.MustCompile(`(?i)(^|[^\p{L}\p{N}])(404|page not found|not found|access denied)([^\p{L}\p{N}]|$)`)

// Check runs every code check and returns the flags that fired, in the fixed
// order not_an_article, error_page, thin (the first one is the reason code
// writes).
func Check(in Input) []Flag {
	var flags []Flag
	if detail := notAnArticle(in); detail != "" {
		flags = append(flags, Flag{Code: NotAnArticle, Detail: detail})
	}
	if detail := errorPage(in); detail != "" {
		flags = append(flags, Flag{Code: ErrorPage, Detail: detail})
	}
	if !in.SkipThin {
		if words := ContentWords(in.Body, in.HostLines); words < MinWords {
			flags = append(flags, Flag{Code: Thin, Detail: thinDetail(words)})
		}
	}
	return flags
}

func thinDetail(words int) string {
	return "body has " + strconv.Itoa(words) + " words after removing navigation (< " + strconv.Itoa(MinWords) + ")"
}

func notAnArticle(in Input) string {
	final := strings.TrimSpace(in.FinalURL)
	if final == "" {
		final = strings.TrimSpace(in.SourceURL)
	}
	finalURL := parseHTTPURL(final)
	requestedURL := parseHTTPURL(in.RequestedURL)

	if finalURL != nil && isRootPath(finalURL.Path) {
		if requestedURL != nil && !isRootPath(requestedURL.Path) {
			return "redirect dropped the requested path " + requestedURL.Path
		}
		return "URL is a site root (" + final + ")"
	}

	title := normalizeTitle(in.Title)
	if title == "" {
		return ""
	}
	site := normalizeTitle(in.SiteName)
	if site == "" && finalURL != nil {
		site = normalizeTitle(siteLabel(finalURL.Hostname()))
	}
	if site != "" && title == site {
		return "title is the site name alone"
	}
	return ""
}

func errorPage(in Input) string {
	if in.StatusCode >= 400 {
		return "HTTP status " + strconv.Itoa(in.StatusCode)
	}
	if errorTitlePattern.MatchString(in.Title) {
		return "title looks like an error page"
	}
	return ""
}

// ContentWords counts the words of body that are not navigation: lines made
// only of links, only table separators, or repeated on at least
// HostRepeatThreshold same-host sources (hostLines) do not count.
func ContentWords(body string, hostLines map[string]int) int {
	words := 0
	for line := range strings.SplitSeq(body, "\n") {
		if NavigationLine(line, hostLines) {
			continue
		}
		for field := range strings.FieldsSeq(line) {
			if strings.ContainsFunc(field, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }) {
				words++
			}
		}
	}
	return words
}

var (
	markdownLinkPattern = regexp.MustCompile(`!?\[[^\]\n]*\]\([^)\n]*\)`)
	wikilinkPattern     = regexp.MustCompile(`!?\[\[[^\]\n]*\]\]`)
	bareURLPattern      = regexp.MustCompile(`(?i)(?:\b[a-z][a-z0-9+.\-]*://|\bwww\.)\S+`)
)

// NavigationLine reports whether a line is navigation-like: only links (plus
// bullets and punctuation), only a `|`-table separator, or repeated on at
// least HostRepeatThreshold same-host sources.
func NavigationLine(line string, hostLines map[string]int) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	if isTableSeparator(trimmed) {
		return true
	}
	if hostLines != nil && hostLines[NormalizeLine(trimmed)] >= HostRepeatThreshold {
		return true
	}
	stripped := markdownLinkPattern.ReplaceAllString(trimmed, "")
	stripped = wikilinkPattern.ReplaceAllString(stripped, "")
	stripped = bareURLPattern.ReplaceAllString(stripped, "")
	if stripped == trimmed {
		return false
	}
	return !strings.ContainsFunc(stripped, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) })
}

func isTableSeparator(line string) bool {
	if !strings.Contains(line, "-") && !strings.Contains(line, "|") {
		return false
	}
	for _, r := range line {
		switch r {
		case '|', '-', ':', ' ', '\t', '+', '=':
		default:
			return false
		}
	}
	return true
}

// NormalizeLine is the key HostLines uses: trimmed, lowercased, whitespace
// collapsed.
func NormalizeLine(line string) string {
	return strings.Join(strings.Fields(strings.ToLower(line)), " ")
}

// HostLineCounts counts, per source host, on how many sources of that host
// each normalized non-empty line appears (a line repeated inside one source
// counts once). Documents without a source host are ignored.
func HostLineCounts(docs []*corpus.Document) map[string]map[string]int {
	counts := map[string]map[string]int{}
	for _, doc := range docs {
		host, _ := doc.Provenance()["source_host"].(string)
		if host == "" {
			continue
		}
		perHost := counts[host]
		if perHost == nil {
			perHost = map[string]int{}
			counts[host] = perHost
		}
		seen := map[string]struct{}{}
		for line := range strings.SplitSeq(doc.Body, "\n") {
			key := NormalizeLine(line)
			if key == "" {
				continue
			}
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			perHost[key]++
		}
	}
	return counts
}

func parseHTTPURL(raw string) *url.URL {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil
	}
	return parsed
}

func isRootPath(p string) bool {
	switch strings.ToLower(p) {
	case "", "/", "/index.html", "/index.htm":
		return true
	}
	return false
}

func siteLabel(host string) string {
	host = strings.TrimPrefix(strings.ToLower(host), "www.")
	label, _, _ := strings.Cut(host, ".")
	return label
}

func normalizeTitle(title string) string {
	return strings.Join(strings.FieldsFunc(strings.ToLower(title), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}), " ")
}
