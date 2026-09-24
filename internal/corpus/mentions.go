package corpus

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/resolve"
)

// maxSentenceRunes bounds Mention.Sentence.
const maxSentenceRunes = 300

var (
	mentionURLPattern      = regexp.MustCompile(`(?i)(?:\b[a-z][a-z0-9+.\-]*://|\bwww\.)[^\s<>()\[\]]+`)
	mentionWikilinkPattern = regexp.MustCompile(`!?\[\[[^\]\n]*\]\]`)
	mentionMarkdownPattern = regexp.MustCompile(`!?\[[^\]\n]*\]\([^)\n]*\)`)
	mentionHeadingPattern  = regexp.MustCompile(`(?m)^ {0,3}#{1,6}(?:[ \t].*)?$`)
)

// Mention is one occurrence of an article title or alias in a body.
type Mention struct {
	Article *Document
	// Term is the dictionary term (title or alias) that matched.
	Term string
	// Start and End are byte offsets of the matched text in the scanned text.
	Start int
	End   int
	// Sentence is the sentence around the match, trimmed and bounded to about
	// 300 characters.
	Sentence string
}

// Dictionary is the alias dictionary of a topic's articles: each article's
// title and `aliases`.
type Dictionary struct {
	byFirst map[string][]dictionaryTerm
}

type dictionaryTerm struct {
	article *Document
	term    string
	words   []string
	gaps    []string
	// exact terms (two-letter acronyms like "AI") match case-sensitively so
	// that ordinary words ("ai" in Portuguese) are not taken for them.
	exact bool
}

// NewDictionary builds the dictionary over articles. Terms shorter than 3
// runes are skipped unless they are all-caps acronyms of at least 2 runes;
// such short acronyms match case-sensitively, every other term
// case-insensitively.
func NewDictionary(articles []*Document) *Dictionary {
	dictionary := &Dictionary{byFirst: make(map[string][]dictionaryTerm)}
	for _, article := range articles {
		seen := make(map[string]struct{})
		for _, term := range append([]string{article.Title}, article.Aliases...) {
			term = strings.TrimSpace(term)
			if term == "" {
				continue
			}
			short := utf8.RuneCountInString(term) < 3
			if short && !isAcronym(term) {
				continue
			}
			spans := wordSpans(term)
			if len(spans) == 0 {
				continue
			}
			key := strings.ToLower(term)
			if short {
				key = term
			}
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}

			entry := dictionaryTerm{article: article, term: term, exact: short}
			for index, span := range spans {
				entry.words = append(entry.words, foldWord(term[span.Start:span.End], short))
				if index > 0 {
					entry.gaps = append(entry.gaps, normalizeGap(term[spans[index-1].End:span.Start]))
				}
			}
			dictionary.byFirst[entry.words[0]] = append(dictionary.byFirst[entry.words[0]], entry)
		}
	}

	return dictionary
}

// Mentions scans text for dictionary terms at word boundaries, skipping
// leading frontmatter, fenced and inline code, heading lines, URLs, and
// existing wikilinks and markdown links. At each position the longest term
// wins; overlapping shorter terms are not reported. Mentions are returned in
// text order, one per article per occurrence.
func (d *Dictionary) Mentions(text string) []Mention {
	if d == nil || len(d.byFirst) == 0 {
		return nil
	}

	masked := maskedRanges(text)
	spans := wordSpans(text)
	mentions := make([]Mention, 0)

	for index := 0; index < len(spans); {
		span := spans[index]
		if resolve.InRanges(masked, span.Start, span.End) {
			index++
			continue
		}

		word := text[span.Start:span.End]
		candidates := d.byFirst[strings.ToLower(word)]
		if lowered := strings.ToLower(word); lowered != word {
			candidates = append(append([]dictionaryTerm(nil), candidates...), d.byFirst[word]...)
		}
		matchedWords := 0
		var matched []dictionaryTerm
		for _, candidate := range candidates {
			if len(candidate.words) < matchedWords || !termMatchesAt(text, spans, index, candidate, masked) {
				continue
			}
			if len(candidate.words) > matchedWords {
				matchedWords = len(candidate.words)
				matched = matched[:0]
			}
			matched = append(matched, candidate)
		}
		if matchedWords == 0 {
			index++
			continue
		}

		start, end := span.Start, spans[index+matchedWords-1].End
		sentence := sentenceAround(text, start, end)
		seenArticles := make(map[*Document]struct{}, len(matched))
		for _, term := range matched {
			if _, exists := seenArticles[term.article]; exists {
				continue
			}
			seenArticles[term.article] = struct{}{}
			mentions = append(mentions, Mention{Article: term.article, Term: term.term, Start: start, End: end, Sentence: sentence})
		}
		index += matchedWords
	}

	return mentions
}

func termMatchesAt(text string, spans []resolve.Range, index int, term dictionaryTerm, masked []resolve.Range) bool {
	if index+len(term.words) > len(spans) {
		return false
	}
	for offset, want := range term.words {
		span := spans[index+offset]
		if foldWord(text[span.Start:span.End], term.exact) != want {
			return false
		}
		if offset > 0 && normalizeGap(text[spans[index+offset-1].End:span.Start]) != term.gaps[offset-1] {
			return false
		}
	}

	last := spans[index+len(term.words)-1]
	return !resolve.InRanges(masked, spans[index].Start, last.End)
}

func maskedRanges(text string) []resolve.Range {
	ranges := resolve.CodeRanges(text)
	if _, body, has := frontmatter.Split(text); has {
		ranges = append(ranges, resolve.Range{Start: 0, End: len(text) - len(body)})
	}
	for _, pattern := range []*regexp.Regexp{mentionHeadingPattern, mentionURLPattern, mentionWikilinkPattern, mentionMarkdownPattern} {
		for _, match := range pattern.FindAllStringIndex(text, -1) {
			ranges = append(ranges, resolve.Range{Start: match[0], End: match[1]})
		}
	}

	return ranges
}

func wordSpans(text string) []resolve.Range {
	spans := make([]resolve.Range, 0, len(text)/5)
	start := -1
	for index, r := range text {
		isWord := unicode.IsLetter(r) || unicode.IsDigit(r)
		switch {
		case isWord && start < 0:
			start = index
		case !isWord && start >= 0:
			spans = append(spans, resolve.Range{Start: start, End: index})
			start = -1
		}
	}
	if start >= 0 {
		spans = append(spans, resolve.Range{Start: start, End: len(text)})
	}

	return spans
}

func foldWord(word string, exact bool) string {
	if exact {
		return word
	}

	return strings.ToLower(word)
}

// normalizeGap canonicalizes the separator between two words: whitespace
// runs collapse and case folds, so "Chain-of-Thought" matches
// "chain-of-thought" and a soft-wrapped "Machine\nLearning" matches "Machine
// Learning". A paragraph break never joins words of one term.
func normalizeGap(gap string) string {
	if strings.Count(gap, "\n") > 1 {
		return "\x00"
	}

	return strings.ToLower(strings.Join(strings.Fields(gap), " "))
}

func isAcronym(term string) bool {
	if utf8.RuneCountInString(term) < 2 {
		return false
	}
	hasLetter := false
	for _, r := range term {
		switch {
		case unicode.IsUpper(r):
			hasLetter = true
		case unicode.IsDigit(r):
		default:
			return false
		}
	}

	return hasLetter
}

func sentenceAround(text string, start, end int) string {
	sentenceStart := 0
	for index := start - 1; index >= 0; index-- {
		if text[index] == '\n' {
			sentenceStart = index + 1
			break
		}
		if isSentenceEnd(text, index) {
			sentenceStart = index + 1
			break
		}
	}

	sentenceEnd := len(text)
	for index := end; index < len(text); index++ {
		if text[index] == '\n' {
			sentenceEnd = index
			break
		}
		if isSentenceEnd(text, index) {
			sentenceEnd = index + 1
			break
		}
	}

	sentence := text[sentenceStart:sentenceEnd]
	if utf8.RuneCountInString(sentence) <= maxSentenceRunes {
		return strings.TrimSpace(sentence)
	}

	half := (maxSentenceRunes - utf8.RuneCountInString(text[start:end])) / 2
	windowStart := start
	for count := 0; windowStart > sentenceStart && count < half; count++ {
		_, size := utf8.DecodeLastRuneInString(text[sentenceStart:windowStart])
		windowStart -= size
	}
	windowEnd := end
	for count := 0; windowEnd < sentenceEnd && count < half; count++ {
		_, size := utf8.DecodeRuneInString(text[windowEnd:sentenceEnd])
		windowEnd += size
	}

	return strings.TrimSpace(text[windowStart:windowEnd])
}

func isSentenceEnd(text string, index int) bool {
	switch text[index] {
	case '.', '!', '?':
		return index+1 >= len(text) || text[index+1] == ' ' || text[index+1] == '\t' || text[index+1] == '\n' || text[index+1] == '\r'
	default:
		return false
	}
}
