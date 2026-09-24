package gate

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/resolve"
)

// trackingParams are query parameters that identify a campaign or a share
// action, never the page; NormalizeURL drops them. Parameters starting with
// `utm_` are dropped as well.
var trackingParams = map[string]struct{}{
	"fbclid": {}, "gclid": {}, "dclid": {}, "msclkid": {}, "yclid": {}, "twclid": {},
	"igshid": {}, "igsh": {}, "mc_cid": {}, "mc_eid": {}, "ref": {}, "ref_src": {}, "ref_url": {},
	"si": {}, "feature": {}, "_hsenc": {}, "_hsmi": {}, "mkt_tok": {}, "spm": {},
}

var (
	youtubeIDPattern     = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)
	instagramCodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{5,}$`)
)

// NormalizeURL returns the dedupe key of a URL: host lowercased without
// `www.`, the path without a trailing slash, the query without tracking
// parameters (`utm_*`, `fbclid`, `gclid`, `si`, `igshid`, ...) and sorted,
// no fragment. The scheme is dropped, so the http and https captures of one
// page are the same source. It returns "" for anything that is not an
// absolute http(s) URL.
func NormalizeURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" {
		return ""
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
	default:
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	if port := parsed.Port(); port != "" && port != "80" && port != "443" {
		host += ":" + port
	}
	cleanPath := strings.TrimRight(parsed.EscapedPath(), "/")
	query := parsed.Query()
	for key := range query {
		lower := strings.ToLower(key)
		if _, tracking := trackingParams[lower]; tracking || strings.HasPrefix(lower, "utm_") {
			query.Del(key)
		}
	}
	key := host + cleanPath
	if encoded := query.Encode(); encoded != "" { // Encode sorts by key
		key += "?" + encoded
	}
	return key
}

// PlatformID returns the platform identity of a media URL:
// `youtube:<video id>` for watch, youtu.be, shorts, embed and live URLs, and
// `instagram:<shortcode>` for post, reel and IGTV URLs; "" otherwise.
func PlatformID(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" {
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	host = strings.TrimPrefix(host, "m.")
	segments := strings.FieldsFunc(parsed.Path, func(r rune) bool { return r == '/' })
	switch {
	case host == "youtu.be":
		if len(segments) > 0 && youtubeIDPattern.MatchString(segments[0]) {
			return "youtube:" + segments[0]
		}
	case host == "youtube.com" || strings.HasSuffix(host, ".youtube.com") || host == "youtube-nocookie.com":
		if id := parsed.Query().Get("v"); youtubeIDPattern.MatchString(id) {
			return "youtube:" + id
		}
		if len(segments) >= 2 && slices.Contains([]string{"shorts", "embed", "live", "v"}, segments[0]) && youtubeIDPattern.MatchString(segments[1]) {
			return "youtube:" + segments[1]
		}
	case host == "instagram.com" || strings.HasSuffix(host, ".instagram.com"):
		for index := 0; index+1 < len(segments); index++ {
			if slices.Contains([]string{"p", "reel", "reels", "tv"}, segments[index]) && instagramCodePattern.MatchString(segments[index+1]) {
				return "instagram:" + segments[index+1]
			}
		}
	}
	return ""
}

// Ref identifies a source for exact dedupe. Empty fields are not compared.
type Ref struct {
	// URL is the source URL (normalized with NormalizeURL).
	URL string
	// PlatformID is a `youtube:<id>` / `instagram:<shortcode>` identity; when
	// empty it is derived from URL.
	PlatformID string
	// BodyHash is corpus.BodyHash of the body about to be written.
	BodyHash string
}

// Match is an existing source equal to a new one.
type Match struct {
	// Path is the topic-relative path of the existing source (the quarantine
	// path for a quarantined one).
	Path string
	// Quarantined reports that the existing source is in raw/_quarantine/.
	Quarantined bool
	// By names what matched: url, platform_id or body_hash.
	By string
}

// Describe renders the match for reports: "duplicate of raw/x.md" or
// "duplicate of quarantined raw/_quarantine/x.md".
func (m Match) Describe() string {
	if m.Quarantined {
		return "duplicate of quarantined " + m.Path + " (" + m.By + ")"
	}
	return "duplicate of " + m.Path + " (" + m.By + ")"
}

// Index is the exact-dedupe index of a topic's sources, quarantined ones
// included (stage 1).
type Index struct {
	urls      map[string]Match
	platforms map[string]Match
	hashes    map[string]Match
}

// NewIndex builds the dedupe index from loaded sources plus every markdown
// file under raw/_quarantine/ (which corpus.Load never loads).
func NewIndex(topicRoot string, sources []*corpus.Document) (*Index, error) {
	index := &Index{urls: map[string]Match{}, platforms: map[string]Match{}, hashes: map[string]Match{}}
	for _, doc := range sources {
		index.Add(doc, false)
	}
	quarantined, err := quarantinedDocuments(topicRoot)
	if err != nil {
		return nil, err
	}
	for _, doc := range quarantined {
		index.Add(doc, true)
	}
	return index, nil
}

// Add registers one source (also used for sources written earlier in the
// same run, so a URL listed twice is ingested once).
func (x *Index) Add(doc *corpus.Document, quarantined bool) {
	match := Match{Path: doc.Path, Quarantined: quarantined}
	sourceURL := doc.SourceURL()
	if key := NormalizeURL(sourceURL); key != "" {
		setOnce(x.urls, key, match, "url")
	}
	for _, id := range documentPlatformIDs(doc) {
		setOnce(x.platforms, id, match, "platform_id")
	}
	if strings.TrimSpace(doc.Body) != "" {
		setOnce(x.hashes, doc.BodyHash, match, "body_hash")
	}
}

func setOnce(values map[string]Match, key string, match Match, by string) {
	if _, ok := values[key]; ok {
		return
	}
	match.By = by
	values[key] = match
}

// Find returns the existing source equal to ref, checking the platform id,
// then the normalized URL, then the body hash.
func (x *Index) Find(ref Ref) (Match, bool) {
	platform := ref.PlatformID
	if platform == "" {
		platform = PlatformID(ref.URL)
	}
	if platform != "" {
		if match, ok := x.platforms[platform]; ok {
			return match, true
		}
	}
	if key := NormalizeURL(ref.URL); key != "" {
		if match, ok := x.urls[key]; ok {
			return match, true
		}
	}
	if ref.BodyHash != "" && ref.BodyHash != corpus.BodyHash("") {
		if match, ok := x.hashes[ref.BodyHash]; ok {
			return match, true
		}
	}
	return Match{}, false
}

// documentPlatformIDs returns the platform identities of a stored source:
// from its `video_id` / `shortcode` frontmatter and from its source_url.
func documentPlatformIDs(doc *corpus.Document) []string {
	ids := make([]string, 0, 2)
	if id := strings.TrimSpace(frontmatter.GetString(doc.Frontmatter, "video_id")); youtubeIDPattern.MatchString(id) {
		ids = append(ids, "youtube:"+id)
	}
	if code := strings.TrimSpace(frontmatter.GetString(doc.Frontmatter, "shortcode")); code != "" && instagramCodePattern.MatchString(code) {
		ids = append(ids, "instagram:"+code)
	}
	if id := PlatformID(doc.SourceURL()); id != "" && !slices.Contains(ids, id) {
		ids = append(ids, id)
	}
	return ids
}

// quarantinedDocuments reads every markdown file under raw/_quarantine/.
// Unreadable files are skipped: dedupe is best effort for them.
func quarantinedDocuments(topicRoot string) ([]*corpus.Document, error) {
	dir := filepath.Join(topicRoot, filepath.FromSlash(resolve.QuarantineDir))
	docs := make([]*corpus.Document, 0)
	err := filepath.WalkDir(dir, func(current string, entry fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return filepath.SkipDir
			}
			return err
		}
		if entry.IsDir() || !strings.EqualFold(path.Ext(entry.Name()), ".md") {
			return nil
		}
		rel, err := filepath.Rel(topicRoot, current)
		if err != nil {
			return err
		}
		doc, err := corpus.ReadDocument(topicRoot, filepath.ToSlash(rel), corpus.KindSource)
		if err == nil {
			docs = append(docs, doc)
		}
		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("gate: read quarantine: %w", err)
	}
	return docs, nil
}
