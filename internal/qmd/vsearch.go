package qmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/compozy/kb/internal/corpus"
)

// Optional vector candidates (spec §9.1 step 4, §10 step 1).
//
// Command choice. The candidate query is
//
//	qmd query --json -n <k> --no-rerank -c <collection> 'vec: <single-line text>'
//
// a structured `vec:` query with reranking off, one subprocess per call.
// Plain `qmd vsearch` is not used although it is "vector only": it runs the
// local 1.7B query-expansion model first (2.4–6.6 s uncached per call in the
// qmd research), while a typed `vec:` line skips expansion and `--no-rerank`
// skips the 0.6B reranker, leaving only the embedding (≈1.2 s per call, of
// which ≈1.1 s is loading the embedding model; ≈1.8 s at 1k documents).
// A `qmd mcp --http` server held for the run would bring that to 0.08–0.76 s
// per call, but it needs a private port, health polling and a process
// lifecycle; the subprocess is kept for simplicity and because the decision
// model, not qmd, is the bottleneck of a find query. For `kb link` over a
// large topic the subprocess cost adds up (≈1.2 s × documents, at most
// maxConcurrentQueries at a time), which is why the CLI prints a note and
// offers --no-qmd.
//
// qmd scores are rank-shaped RRF values here and never enter bands: the
// adapters return paths only.

// maxConcurrentQueries bounds concurrent qmd subprocesses (each loads the
// embedding model).
const maxConcurrentQueries = 2

// maxVectorQueryRunes caps the embedded text (qmd truncates at 2048 tokens).
const maxVectorQueryRunes = 4000

// neighbourSlack is the number of extra hits asked for so self and hub
// matches can be dropped.
const neighbourSlack = 10

var indexPathPattern = regexp.MustCompile(`(?m)^Index:\s+(.+?)\s*$`)

// Candidates proposes topic documents by qmd vector similarity. It
// implements link.Neighbours and find.Vectors. Use OpenCandidates, which
// only returns one when the topic collection exists and is fresh.
type Candidates struct {
	client     *QMDClient
	collection string
	sem        chan struct{}
}

// OpenCandidates returns vector candidates over the qmd collection of a
// topic when they can be trusted, or nil and a one-line note saying why they
// are off. The collection is fresh when:
//
//   - the qmd binary resolves on PATH;
//   - `qmd status` lists a collection named collection;
//   - the index file reported by `qmd status` (`Index: <path>`) was modified
//     no earlier than the newest `*.md` file of the topic (dot directories
//     skipped), i.e. no topic file changed after the last `kb index`.
//
// The index mtime is index-wide (updating another collection also bumps it),
// so the rule catches edits made after the last index run, not every stale
// case; pending embeddings only shrink the vector result set.
func (client *QMDClient) OpenCandidates(ctx context.Context, collection, topicRoot string) (*Candidates, string) {
	collection = strings.TrimSpace(collection)
	if collection == "" {
		return nil, "qmd candidates: off (no collection name)"
	}
	if _, err := client.resolveBinary(); err != nil {
		return nil, "qmd candidates: off (qmd is not installed)"
	}
	stdout, _, err := client.run(ctx, client.statusCommand())
	if err != nil {
		return nil, "qmd candidates: off (qmd status failed)"
	}
	status, err := parseIndexStatus(stdout)
	if err != nil {
		return nil, "qmd candidates: off (qmd status output not understood)"
	}
	found := false
	for _, info := range status.Collections {
		if info.Name == collection {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Sprintf("qmd candidates: off (no qmd collection %q; run `kb index --topic %s`)", collection, collection)
	}
	match := indexPathPattern.FindStringSubmatch(cleanOutput(stdout))
	if match == nil {
		return nil, "qmd candidates: off (qmd status did not report the index path)"
	}
	indexInfo, err := os.Stat(match[1])
	if err != nil {
		return nil, "qmd candidates: off (qmd index file not found)"
	}
	newest, err := newestMarkdown(topicRoot)
	if err != nil {
		return nil, "qmd candidates: off (cannot read the topic files)"
	}
	if indexInfo.ModTime().Before(newest) {
		return nil, fmt.Sprintf("qmd candidates: off (qmd index is older than the topic files; run `kb index --topic %s`)", collection)
	}
	return &Candidates{client: client, collection: collection, sem: make(chan struct{}, maxConcurrentQueries)}, ""
}

// newestMarkdown returns the newest modification time of a *.md file under
// root, skipping dot directories.
func newestMarkdown(root string) (time.Time, error) {
	var newest time.Time
	err := filepath.WalkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if current != root && strings.HasPrefix(entry.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
		return nil
	})
	return newest, err
}

// Query runs one no-rerank vector query and returns up to limit
// collection-relative paths, best first, without duplicates.
func (c *Candidates) Query(ctx context.Context, text string, limit int) ([]string, error) {
	text = singleLine(text)
	if text == "" {
		return nil, errors.New("qmd vector query: text is required")
	}
	if limit <= 0 {
		return nil, errors.New("qmd vector query: limit must be positive")
	}
	select {
	case c.sem <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	defer func() { <-c.sem }()

	stdout, _, err := c.client.run(ctx, commandSpec{
		label: "query (vector candidates)",
		args:  c.client.baseArgs("query", "--json", "-n", strconv.Itoa(limit), "--no-rerank", "-c", c.collection, "vec: "+text),
	})
	if err != nil {
		return nil, err
	}
	return parseCandidatePaths(stdout, c.collection)
}

// Neighbours implements link.Neighbours: the top k wiki articles closest to
// the document's title, summary and opening text, the document itself and
// non-article paths excluded.
func (c *Candidates) Neighbours(ctx context.Context, doc *corpus.Document, k int) ([]string, error) {
	paths, err := c.Query(ctx, neighbourText(doc), k+neighbourSlack)
	if err != nil {
		return nil, err
	}
	kept := make([]string, 0, k)
	for _, p := range paths {
		if len(kept) == k {
			break
		}
		if p == doc.Path || !strings.HasPrefix(p, "wiki/concepts/") {
			continue
		}
		kept = append(kept, p)
	}
	return kept, nil
}

// Search implements find.Vectors: the top k topic documents closest to the
// question.
func (c *Candidates) Search(ctx context.Context, query string, k int) ([]string, error) {
	return c.Query(ctx, query, k)
}

// neighbourText is the query text of a document: title, summary and the
// start of the body.
func neighbourText(doc *corpus.Document) string {
	return strings.Join([]string{doc.Title, doc.Summary(), doc.Body}, " ")
}

// singleLine collapses whitespace (structured query lines must be single
// line) and caps the length.
func singleLine(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	runes := []rune(text)
	if len(runes) > maxVectorQueryRunes {
		text = string(runes[:maxVectorQueryRunes])
	}
	return text
}

// parseCandidatePaths maps qmd JSON hits to collection-relative paths:
// `qmd://<collection>/` and any `?index=` suffix are stripped and the path is
// URL-decoded. Quarantined sources and decision records (ExcludedPath) are
// dropped.
func parseCandidatePaths(stdout, collection string) ([]string, error) {
	var hits []searchResultPayload
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &hits); err != nil {
		return nil, fmt.Errorf("qmd vector query: parse JSON output: %w", err)
	}
	seen := map[string]bool{}
	paths := make([]string, 0, len(hits))
	for _, hit := range hits {
		p := candidatePath(firstNonEmpty(hit.File, hit.FilePath, hit.DisplayPath), collection)
		if p == "" || seen[p] || hit.excluded() || ExcludedPath(p) {
			continue
		}
		seen[p] = true
		paths = append(paths, p)
	}
	return paths, nil
}

func candidatePath(raw, collection string) string {
	p := strings.TrimSpace(raw)
	if index := strings.Index(p, "?index="); index >= 0 {
		p = p[:index]
	}
	p = strings.TrimPrefix(p, "qmd://")
	p = strings.TrimPrefix(p, collection+"/")
	if strings.Contains(p, "%") {
		if decoded, err := url.PathUnescape(p); err == nil {
			p = decoded
		}
	}
	return strings.TrimPrefix(filepath.ToSlash(p), "/")
}
