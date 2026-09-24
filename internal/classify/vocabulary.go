package classify

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/session"
)

// Vocabulary bootstrap limits (spec §5.2).
const (
	MinVocabulary = 12
	MaxVocabulary = 40
	// maxVocabularySources caps the source lines given to the vocabulary
	// draft; larger topics are sampled evenly by path so the prompt stays
	// within the generation model's context.
	maxVocabularySources = 300
	// maxTitleWords bounds a generated concept title.
	maxTitleWords = 8
)

// ConceptsDir is the topic-relative directory of wiki concept articles.
const ConceptsDir = "wiki/concepts"

// ErrArticleExists is returned by CreateStubArticle when the article file is
// already there; it is never overwritten.
var ErrArticleExists = errors.New("classify: article already exists")

// ErrNoVocabularyDraft is returned by AcceptVocabulary without a draft.
var ErrNoVocabularyDraft = errors.New("classify: topic.yaml has no vocabulary_draft; run `kb topic vocabulary <topic> --draft` first")

// DraftVocabulary asks the generation model for a closed list of 12–40
// concepts (title + criterion) from the contract and one line per source
// (title + summary), validates it by code (unique titles, criterion rules,
// no collision with existing articles), counts for each concept how many
// source summaries mention its title (code, no call) and saves the list
// under topic.yaml `vocabulary_draft:`.
func DraftVocabulary(ctx context.Context, s *session.Session) ([]contract.VocabularyItem, error) {
	c, err := s.Corpus()
	if err != nil {
		return nil, fmt.Errorf("classify: vocabulary draft: %w", err)
	}
	sources := c.Sources()
	if len(sources) == 0 {
		return nil, errors.New("classify: vocabulary draft: the topic has no sources")
	}

	lines := make([]string, 0, maxVocabularySources)
	for _, doc := range sampleEvenly(sources, maxVocabularySources) {
		line := "- " + doc.Title
		if summary := doc.Summary(); summary != "" {
			line += ": " + summary
		}
		lines = append(lines, line)
	}
	prompt := strings.Join([]string{
		fmt.Sprintf("Propose the concept vocabulary of the knowledge-base topic below: a closed list of %d to %d concepts that together cover what its sources discuss. Each concept becomes a wiki article that documents are classified into.", MinVocabulary, MaxVocabulary),
		"- title: the concept's name as a reader would look it up, 1 to 8 words, no dates.",
		"- criterion: " + criterionRules,
		"Concepts must not overlap: two concepts never share the same subject.",
		"",
		"Topic: " + s.Topic.Title + " (domain: " + s.Topic.Domain + ")",
		contractLines(s.Contract),
		"",
		fmt.Sprintf("Sources (%d of %d, title: summary):", len(lines), len(sources)),
		strings.Join(lines, "\n"),
	}, "\n")
	proposed, err := generateConcepts(ctx, s, KindVocabulary, s.Topic.Slug, prompt)
	if err != nil {
		return nil, err
	}

	taken := map[string]bool{}
	for _, article := range c.Articles() {
		for _, name := range append([]string{article.Title}, article.Aliases...) {
			taken[strings.ToLower(strings.TrimSpace(name))] = true
		}
	}
	items := make([]contract.VocabularyItem, 0, len(proposed))
	for _, concept := range validConcepts(proposed, taken, MaxVocabulary) {
		items = append(items, contract.VocabularyItem{
			Title:     concept.Title,
			Criterion: concept.Criterion,
			Mentions:  summaryMentions(sources, concept.Title),
		})
	}
	if len(items) < MinVocabulary {
		return nil, fmt.Errorf("classify: vocabulary draft: the generation model returned %d valid concepts (of %d proposed), need %d to %d; re-run --draft", len(items), len(proposed), MinVocabulary, MaxVocabulary)
	}
	if err := contract.SaveVocabularyDraft(s.Root(), items); err != nil {
		return nil, fmt.Errorf("classify: vocabulary draft: %w", err)
	}
	return items, nil
}

// AcceptVocabulary creates one stub article per `vocabulary_draft:` item
// (skipping articles that already exist), clears the draft and returns the
// topic-relative paths it created.
func AcceptVocabulary(s *session.Session) ([]string, error) {
	settings, err := contract.LoadSettings(s.Root())
	if err != nil {
		return nil, fmt.Errorf("classify: vocabulary accept: %w", err)
	}
	if len(settings.VocabularyDraft) == 0 {
		return nil, ErrNoVocabularyDraft
	}
	now := s.Now()
	created := make([]string, 0, len(settings.VocabularyDraft))
	for _, item := range settings.VocabularyDraft {
		rel, err := CreateStubArticle(s.Root(), s.Topic.Domain, item.Title, item.Criterion, now)
		if errors.Is(err, ErrArticleExists) {
			continue
		}
		if err != nil {
			return created, fmt.Errorf("classify: vocabulary accept: %w", err)
		}
		created = append(created, rel)
	}
	if err := contract.ClearVocabularyDraft(s.Root()); err != nil {
		return created, fmt.Errorf("classify: vocabulary accept: %w", err)
	}
	s.Settings.VocabularyDraft = nil
	s.ReloadCorpus()
	return created, nil
}

// CreateStubArticle writes `wiki/concepts/<Title>.md` (file name = title
// sanitized for the filesystem) with the stub frontmatter of the plan:
// title, type wiki, stage stub, domain, tags [domain, wiki, stub], created,
// updated, sources [] and criterion. It returns the topic-relative path, and
// ErrArticleExists (with that path) when the file is already there.
func CreateStubArticle(topicRoot, domain, title, criterion string, now time.Time) (string, error) {
	title = strings.TrimSpace(title)
	name := sanitizeFileName(title)
	if name == "" {
		return "", fmt.Errorf("classify: stub article: title %q has no usable characters", title)
	}
	rel := path.Join(ConceptsDir, name+".md")
	absolute := filepath.Join(topicRoot, filepath.FromSlash(rel))
	if _, err := os.Stat(absolute); err == nil {
		return rel, fmt.Errorf("%w: %s", ErrArticleExists, rel)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("classify: stub article %s: %w", rel, err)
	}

	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	values := map[string]any{
		"title":   title,
		"type":    "wiki",
		"stage":   "stub",
		"domain":  domain,
		"tags":    []string{domain, "wiki", "stub"},
		"created": day,
		"updated": day,
		"sources": []string{},
	}
	if criterion = strings.TrimSpace(criterion); criterion != "" {
		values["criterion"] = criterion
	}
	content, err := frontmatter.Generate(values, "# "+title+"\n")
	if err != nil {
		return "", fmt.Errorf("classify: stub article %s: %w", rel, err)
	}
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		return "", fmt.Errorf("classify: stub article %s: %w", rel, err)
	}
	file, err := os.OpenFile(absolute, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if errors.Is(err, os.ErrExist) {
		return rel, fmt.Errorf("%w: %s", ErrArticleExists, rel)
	}
	if err != nil {
		return "", fmt.Errorf("classify: stub article %s: %w", rel, err)
	}
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		return "", fmt.Errorf("classify: stub article %s: %w", rel, err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("classify: stub article %s: %w", rel, err)
	}
	return rel, nil
}

// isStub reports an article created by CreateStubArticle and not compiled
// yet.
func isStub(doc *corpus.Document) bool {
	return strings.EqualFold(strings.TrimSpace(frontmatter.GetString(doc.Frontmatter, "stage")), "stub")
}

// sanitizeFileName turns a title into a file name: characters that are
// invalid on common filesystems or special in Obsidian links (`/ \ : * ? " <
// > | # ^ [ ]`) and control characters become spaces, whitespace collapses,
// leading/trailing dots and spaces go, and the result is capped at 120 runes.
func sanitizeFileName(title string) string {
	mapped := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f || strings.ContainsRune(`/\:*?"<>|#^[]`, r) {
			return ' '
		}
		return r
	}, title)
	name := strings.Trim(strings.Join(strings.Fields(mapped), " "), ". ")
	if runes := []rune(name); len(runes) > MaxTitleChars {
		name = strings.Trim(string(runes[:MaxTitleChars]), ". ")
	}
	return name
}

// validConcepts keeps generated concepts with a valid title (1–8 words, a
// literal, a usable file name), a valid criterion, and a title that is new
// (case-insensitive) against taken and the earlier items; at most limit.
func validConcepts(proposed []proposedConcept, taken map[string]bool, limit int) []proposedConcept {
	seen := map[string]bool{}
	kept := make([]proposedConcept, 0, len(proposed))
	for _, concept := range proposed {
		concept.Title = strings.TrimSpace(concept.Title)
		concept.Criterion = strings.TrimSpace(concept.Criterion)
		key := strings.ToLower(concept.Title)
		words := len(strings.Fields(concept.Title))
		switch {
		case len(kept) >= limit:
			return kept
		case validateTitle(concept.Title) != nil, words < 1 || words > maxTitleWords:
		case validateCriterion(concept.Criterion, concept.Title) != nil:
		case seen[key], taken[key], seen[strings.ToLower(sanitizeFileName(concept.Title))]:
		default:
			seen[key] = true
			seen[strings.ToLower(sanitizeFileName(concept.Title))] = true
			kept = append(kept, concept)
		}
	}
	return kept
}

// summaryMentions counts the sources whose summary mentions title
// (case-insensitive).
func summaryMentions(sources []*corpus.Document, title string) int {
	needle := strings.ToLower(strings.TrimSpace(title))
	if needle == "" {
		return 0
	}
	count := 0
	for _, doc := range sources {
		if strings.Contains(strings.ToLower(doc.Summary()), needle) {
			count++
		}
	}
	return count
}

// sampleEvenly returns at most limit documents spread evenly over docs
// (which are sorted by path), deterministically.
func sampleEvenly(docs []*corpus.Document, limit int) []*corpus.Document {
	if len(docs) <= limit {
		return docs
	}
	out := make([]*corpus.Document, 0, limit)
	for index := range limit {
		out = append(out, docs[index*len(docs)/limit])
	}
	return out
}

// contractLines renders the contract for generation prompts.
func contractLines(c *contract.Contract) string {
	if c.Empty() {
		return "Selection contract: none accepted yet."
	}
	n := c.Normalized()
	parts := []string{"Selection contract:", "- purpose: " + n.Purpose}
	add := func(label string, values []string) {
		if len(values) > 0 {
			parts = append(parts, "- "+label+": "+strings.Join(values, "; "))
		}
	}
	add("core", n.Core)
	add("adjacent", n.Adjacent)
	add("collected on purpose", n.CollectedOnPurpose)
	add("out of scope", n.OutOfScope)
	return strings.Join(parts, "\n")
}
