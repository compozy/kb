package corpus

import (
	"net/url"
	"strings"
)

// Provenance returns the document's provenance state (spec §4.4): `path`
// (topic-relative), `source_kind`, `source_host` (host of `source_url`,
// lowercased, `www.` stripped), `ingest_batch` and `ingest_query`. Empty
// fields are omitted, never guessed.
func (d *Document) Provenance() map[string]any {
	provenance := make(map[string]any, 5)
	setIfPresent(provenance, "path", d.Path)
	setIfPresent(provenance, "source_kind", d.SourceKind())
	setIfPresent(provenance, "source_host", sourceHost(d.SourceURL()))
	setIfPresent(provenance, "ingest_batch", d.stringKey("ingest_batch"))
	setIfPresent(provenance, "ingest_query", d.stringKey("ingest_query"))

	return provenance
}

// DocumentState is the single state builder every decision purpose uses for
// a document (spec §4.4): `{"title", "excerpt", "provenance"}`, with the
// excerpt bounded by maxTokens and biased toward sections matching terms
// (see Excerpt).
func DocumentState(d *Document, maxTokens int, terms []string) map[string]any {
	return map[string]any{
		"title":      d.Title,
		"excerpt":    Excerpt(d.Body, maxTokens, terms),
		"provenance": d.Provenance(),
	}
}

func sourceHost(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Hostname() == "" {
		return ""
	}

	return strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
}

func setIfPresent(values map[string]any, key, value string) {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		values[key] = trimmed
	}
}
