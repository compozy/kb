// Package corpus loads a topic's documents (sources and wiki articles) and
// provides everything decisions need to reason about them without a model:
// body hashes, value hashes, BM25 retrieval, the alias dictionary and mention
// scan, section-aware excerpts, provenance state, the `.decisions/state.jsonl`
// store and the byte-preserving owned-key writer.
package corpus

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/resolve"
)

// Kind classifies a corpus document.
type Kind string

// Document kinds.
const (
	KindSource  Kind = "source"
	KindArticle Kind = "article"
)

// Document is one loaded topic document.
type Document struct {
	// Path is the topic-relative slash path (with `.md`).
	Path string
	// AbsPath is the absolute filesystem path.
	AbsPath string
	// Title is the frontmatter title, else the first `# ` heading, else the
	// file stem.
	Title string
	Kind  Kind
	// Frontmatter holds the parsed frontmatter values.
	Frontmatter map[string]any
	// Raw is the full file text as read.
	Raw string
	// Body is the text after the frontmatter.
	Body string
	// BodyHash is BodyHash(Body).
	BodyHash string
	// Aliases are the Obsidian `aliases` values.
	Aliases []string
	ModTime time.Time
}

// Summary returns the `summary` key.
func (d *Document) Summary() string { return d.stringKey("summary") }

// Criterion returns the `criterion` key.
func (d *Document) Criterion() string { return d.stringKey("criterion") }

// Concepts returns the normalized wikilink targets of `concepts`.
func (d *Document) Concepts() []string { return d.linkList("concepts") }

// Entities returns the `entities` list.
func (d *Document) Entities() []string { return d.stringList("entities") }

// Questions returns the `questions` list.
func (d *Document) Questions() []string { return d.stringList("questions") }

// Genre returns the `genre` key.
func (d *Document) Genre() string { return d.stringKey("genre") }

// Relevance returns the `relevance` key.
func (d *Document) Relevance() string { return d.stringKey("relevance") }

// Triage returns the `triage` key.
func (d *Document) Triage() string { return d.stringKey("triage") }

// SourceURL returns the `source_url` key.
func (d *Document) SourceURL() string { return d.stringKey("source_url") }

// SourceKind returns the `source_kind` key.
func (d *Document) SourceKind() string { return d.stringKey("source_kind") }

// Locked reports whether the user set `locked: true`.
func (d *Document) Locked() bool {
	return frontmatter.GetBool(d.Frontmatter, "locked")
}

// Depth returns the numeric `depth` key.
func (d *Document) Depth() (float64, bool) {
	switch typed := d.Frontmatter["depth"].(type) {
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float64:
		return typed, true
	case string:
		value, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0, false
		}
		return value, true
	default:
		return 0, false
	}
}

func (d *Document) stringKey(key string) string {
	if d == nil || d.Frontmatter == nil {
		return ""
	}
	if value, exists := d.Frontmatter[key]; !exists || value == nil {
		return ""
	}

	return strings.TrimSpace(frontmatter.GetString(d.Frontmatter, key))
}

func (d *Document) stringList(key string) []string {
	if d == nil || d.Frontmatter == nil {
		return nil
	}
	raw := frontmatter.GetStringSlice(d.Frontmatter, key)
	if single, ok := d.Frontmatter[key].(string); ok {
		raw = []string{single}
	}

	values := make([]string, 0, len(raw))
	for _, item := range raw {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			values = append(values, trimmed)
		}
	}

	return values
}

func (d *Document) linkList(key string) []string {
	items := d.stringList(key)
	targets := make([]string, 0, len(items))
	for _, item := range items {
		if target := resolve.Normalize(item); target != "" {
			targets = append(targets, target)
		}
	}

	return targets
}

// LoadOptions configures Load.
type LoadOptions struct {
	// Exclude lists topic-relative globs (see MatchGlob) of documents to skip.
	Exclude []string
}

// Skip records a document that could not be loaded.
type Skip struct {
	Path   string
	Reason string
}

// Corpus is the loaded set of topic documents.
type Corpus struct {
	// Root is the topic root directory.
	Root string
	// Skipped lists documents skipped because they could not be read or
	// their frontmatter did not parse.
	Skipped []Skip

	documents []*Document
	byPath    map[string]*Document
	byHash    map[string][]*Document
}

// Load reads the topic at topicRoot. Sources are `raw/**/*.md` except
// `raw/_quarantine/**`, `raw/codebase/**` and documents whose `source_kind`
// starts with `codebase-`; articles are `wiki/concepts/**/*.md`. Hub files
// (`log.md`, `CLAUDE.md`, `AGENTS.md`, `wiki/index/**`, `outputs/**`,
// `.decisions/**`, `bases/**`) fall outside both trees and dot directories are
// never walked, so neither is ever loaded. A
// document whose frontmatter does not parse is skipped and recorded in
// Corpus.Skipped; it never fails the load.
func Load(topicRoot string, opts LoadOptions) (*Corpus, error) {
	root := filepath.Clean(topicRoot)
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("corpus: stat topic %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("corpus: topic %q is not a directory", root)
	}

	corpus := &Corpus{
		Root:   root,
		byPath: make(map[string]*Document),
		byHash: make(map[string][]*Document),
	}

	walkErr := filepath.WalkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %q: %w", current, walkErr)
		}
		relative, err := filepath.Rel(root, current)
		if err != nil {
			return fmt.Errorf("relative path for %q: %w", current, err)
		}
		relative = filepath.ToSlash(relative)

		if entry.IsDir() {
			if relative != "." && skipDir(relative) {
				return filepath.SkipDir
			}
			return nil
		}

		kind, ok := documentKind(relative)
		if !ok || excluded(opts.Exclude, relative) {
			return nil
		}

		document, reason := loadDocument(current, relative, kind)
		if reason != "" {
			corpus.Skipped = append(corpus.Skipped, Skip{Path: relative, Reason: reason})
			return nil
		}
		if kind == KindSource && strings.HasPrefix(document.SourceKind(), "codebase-") {
			return nil
		}
		corpus.add(document)
		return nil
	})
	if walkErr != nil {
		return nil, fmt.Errorf("corpus: %w", walkErr)
	}

	sort.Slice(corpus.documents, func(i, j int) bool {
		return corpus.documents[i].Path < corpus.documents[j].Path
	})
	sort.Slice(corpus.Skipped, func(i, j int) bool {
		return corpus.Skipped[i].Path < corpus.Skipped[j].Path
	})

	return corpus, nil
}

// Documents returns every loaded document sorted by path.
func (c *Corpus) Documents() []*Document {
	return append([]*Document(nil), c.documents...)
}

// Sources returns the source documents sorted by path.
func (c *Corpus) Sources() []*Document { return c.ofKind(KindSource) }

// Articles returns the article documents sorted by path.
func (c *Corpus) Articles() []*Document { return c.ofKind(KindArticle) }

// ByPath returns the document at the topic-relative path p (with or without
// `.md`), or nil.
func (c *Corpus) ByPath(p string) *Document {
	cleaned := strings.Trim(filepath.ToSlash(strings.TrimSpace(p)), "/")
	if document := c.byPath[cleaned]; document != nil {
		return document
	}

	return c.byPath[cleaned+".md"]
}

// ByBodyHash returns the documents whose body hash is h, sorted by path.
func (c *Corpus) ByBodyHash(h string) []*Document {
	return append([]*Document(nil), c.byHash[h]...)
}

func (c *Corpus) ofKind(kind Kind) []*Document {
	documents := make([]*Document, 0, len(c.documents))
	for _, document := range c.documents {
		if document.Kind == kind {
			documents = append(documents, document)
		}
	}

	return documents
}

func (c *Corpus) add(document *Document) {
	c.documents = append(c.documents, document)
	c.byPath[document.Path] = document
	c.byHash[document.BodyHash] = append(c.byHash[document.BodyHash], document)
}

// ReadDocument loads one topic document from disk. relative is the
// topic-relative path.
func ReadDocument(topicRoot, relative string, kind Kind) (*Document, error) {
	relative = strings.Trim(filepath.ToSlash(relative), "/")
	document, reason := loadDocument(filepath.Join(topicRoot, filepath.FromSlash(relative)), relative, kind)
	if reason != "" {
		return nil, fmt.Errorf("corpus: read %q: %s", relative, reason)
	}

	return document, nil
}

func loadDocument(absolute, relative string, kind Kind) (*Document, string) {
	content, err := os.ReadFile(absolute)
	if err != nil {
		return nil, fmt.Sprintf("read: %v", err)
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return nil, fmt.Sprintf("stat: %v", err)
	}

	raw := string(content)
	values, body, err := frontmatter.Parse(raw)
	if err != nil {
		var fmErr *frontmatter.Error
		if errors.As(err, &fmErr) {
			return nil, fmt.Sprintf("frontmatter: %s", fmErr.Kind)
		}
		return nil, fmt.Sprintf("frontmatter: %v", err)
	}

	document := &Document{
		Path:        relative,
		AbsPath:     absolute,
		Kind:        kind,
		Frontmatter: values,
		Raw:         raw,
		Body:        body,
		BodyHash:    BodyHash(body),
		Aliases:     resolve.Aliases(values),
		ModTime:     info.ModTime(),
	}
	document.Title = documentTitle(values, body, relative)

	return document, ""
}

func documentTitle(values map[string]any, body, relative string) string {
	if values["title"] != nil {
		if title := strings.TrimSpace(frontmatter.GetString(values, "title")); title != "" {
			return title
		}
	}
	for line := range strings.SplitSeq(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if heading, ok := strings.CutPrefix(trimmed, "# "); ok && strings.TrimSpace(heading) != "" {
			return strings.TrimSpace(heading)
		}
	}

	return strings.TrimSuffix(path.Base(relative), path.Ext(relative))
}

func documentKind(relative string) (Kind, bool) {
	if !strings.EqualFold(path.Ext(relative), ".md") {
		return "", false
	}
	switch {
	case strings.HasPrefix(relative, "raw/"):
		return KindSource, true
	case strings.HasPrefix(relative, "wiki/concepts/"):
		return KindArticle, true
	default:
		return "", false
	}
}

func skipDir(relative string) bool {
	if strings.HasPrefix(path.Base(relative), ".") {
		return true
	}
	switch relative {
	case resolve.QuarantineDir, "raw/codebase", "wiki/index", "wiki/codebase", "outputs", "bases":
		return true
	}

	return false
}

// ListPaths renders document paths for a run summary: at most limit of
// them, comma-separated, then "(+N more)" when some were left out.
func ListPaths(paths []string, limit int) string {
	if limit <= 0 || len(paths) <= limit {
		return strings.Join(paths, ", ")
	}
	return fmt.Sprintf("%s (+%d more)", strings.Join(paths[:limit], ", "), len(paths)-limit)
}

func excluded(patterns []string, relative string) bool {
	for _, pattern := range patterns {
		if MatchGlob(pattern, relative) {
			return true
		}
	}

	return false
}
