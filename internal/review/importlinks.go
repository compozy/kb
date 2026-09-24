package review

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/resolve"
)

// InsertedLinksFile is the log of body links inserted by `kb link`, inside
// <topic>/.decisions/. Rows: {time, subject, target, text, mode}.
const InsertedLinksFile = "inserted-links.jsonl"

// insertedLink is one row of inserted-links.jsonl.
type insertedLink struct {
	Time    string `json:"time"`
	Subject string `json:"subject"`
	Target  string `json:"target"`
	Text    string `json:"text"`
	Mode    string `json:"mode"`
}

// importLinkKeys are the frontmatter lists whose links count as human links
// unless kb wrote them.
var importLinkKeys = append([]string{"sources"}, decisions.RelationKeys...)

// ImportLinks turns the existing links between topic documents into positive
// `link` labels (question `should_link`, subject = linking document, target =
// linked document, origin `import-links`), which gives calibration a gold
// set (spec §12.2). Counted: body wikilinks and markdown links plus the
// `sources` and relation frontmatter lists. Excluded: links kb inserted
// (inserted-links.jsonl rows not in shadow mode) and relation lists kb wrote
// that nobody edited since (state row `written` hash equals the current
// value hash). Links to quarantined or non-corpus files and self links are
// ignored. A subject/target pair that already has a link label is skipped,
// so re-running writes nothing new. Returns the number of labels written.
func ImportLinks(vaultPath, topicRoot string, c *corpus.Corpus) (int, error) {
	if c == nil {
		loaded, err := corpus.Load(topicRoot, corpus.LoadOptions{})
		if err != nil {
			return 0, fmt.Errorf("review: %w", err)
		}
		c = loaded
	}
	index, _, err := resolve.LoadVaultIndex(vaultPath, topicRoot)
	if err != nil {
		return 0, fmt.Errorf("review: %w", err)
	}
	state, err := corpus.OpenState(topicRoot)
	if err != nil {
		return 0, fmt.Errorf("review: %w", err)
	}
	slug := filepath.Base(filepath.Clean(topicRoot))
	vaultPathOf := func(rel string) string { return path.Join(slug, rel) }

	resolveTarget := func(subject, target string) string {
		file := index.ResolveFrom(vaultPathOf(subject), target)
		if file == nil || !file.InTopic || file.Quarantined || file.TopicRel == subject || c.ByPath(file.TopicRel) == nil {
			return ""
		}
		return file.TopicRel
	}

	inserted := map[string]bool{}
	err = readJSONL(filepath.Join(topicRoot, decisions.ReceiptsDir, InsertedLinksFile), func(line []byte) error {
		var row insertedLink
		if json.Unmarshal(line, &row) != nil || row.Subject == "" || strings.EqualFold(row.Mode, "shadow") {
			return nil
		}
		if target := resolveTarget(row.Subject, row.Target); target != "" {
			inserted[row.Subject+"\x00"+target] = true
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	existing, err := LoadLabels(topicRoot)
	if err != nil {
		return 0, err
	}
	seen := map[string]bool{}
	for _, label := range existing {
		if label.Purpose == PurposeLink {
			seen[label.Subject+"\x00"+label.Target] = true
		}
	}

	store := Open(topicRoot, nil)
	written := 0
	for _, doc := range c.Documents() {
		targets := make([]string, 0)
		for _, link := range resolve.BodyLinks(doc.Body) {
			targets = append(targets, link.Target)
		}
		row, _ := state.Get(doc.Path)
		links := resolve.FrontmatterLinks(doc.Frontmatter)
		for _, key := range importLinkKeys {
			if len(links[key]) == 0 || kbWrote(row, key, doc.Frontmatter[key]) {
				continue
			}
			targets = append(targets, links[key]...)
		}
		for _, raw := range targets {
			target := resolveTarget(doc.Path, raw)
			pair := doc.Path + "\x00" + target
			if target == "" || inserted[pair] || seen[pair] {
				continue
			}
			seen[pair] = true
			if err := store.AddLabel(Label{
				Subject:  doc.Path,
				Purpose:  PurposeLink,
				Question: QuestionShouldLink,
				Target:   target,
				Verdict:  VerdictPositive,
				Origin:   OriginImportLinks,
			}); err != nil {
				return written, err
			}
			written++
		}
	}
	return written, nil
}

// kbWrote reports whether the state row records key as written by kb with
// the current value.
func kbWrote(row *corpus.StateRow, key string, value any) bool {
	if row == nil || !slices.Contains(decisions.RelationKeys, key) {
		return false
	}
	hash, ok := row.Written[key]
	return ok && hash == corpus.ValueHash(value)
}
