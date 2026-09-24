package corpus

import (
	"math"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// BM25 parameters.
const (
	bm25K1 = 1.2
	bm25B  = 0.75
	// HeadChars is the number of body characters indexed by DefaultFields.
	HeadChars = 2000
)

// stopwords is a small English + Portuguese list dropped by Tokenize.
var stopwords = func() map[string]struct{} {
	words := strings.Fields(`
		a an and are as at be but by for from has have in into is it its of on or
		that the their this to was were will with what which who how why when
		not no do does can than then there these those they we you your our
		o os as um uma uns umas de do da dos das em no na nos nas por para com
		que se e ou ao aos à às é são ser foi como mais mas não sem sobre entre
		seu sua seus suas ele ela eles elas isso isto este esta esse essa
	`)
	set := make(map[string]struct{}, len(words))
	for _, word := range words {
		set[word] = struct{}{}
	}
	return set
}()

// Tokenize lowercases text, splits it on every rune that is not a Unicode
// letter or digit, drops stopwords and keeps tokens of at least 2 runes.
func Tokenize(text string) []string {
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	tokens := fields[:0]
	for _, field := range fields {
		if utf8.RuneCountInString(field) < 2 {
			continue
		}
		if _, stop := stopwords[field]; stop {
			continue
		}
		tokens = append(tokens, field)
	}

	return tokens
}

// Hit is one BM25 search result.
type Hit struct {
	Doc   *Document
	Score float64
}

// BM25 is an in-memory BM25F index over documents: each named field has its
// own length normalization and a weight applied to its term frequency.
type BM25 struct {
	docs      []*Document
	weights   map[string]float64
	fieldTF   []map[string]map[string]int
	fieldLen  []map[string]int
	avgLen    map[string]float64
	docFreq   map[string]int
	docsTotal int
	// fieldNames is sorted so score sums are deterministic.
	fieldNames []string
}

// NewBM25 indexes docs. fields extracts the text of every named field of a
// document; weights gives each field's weight (fields missing from weights
// weigh 1). Parameters: k1 = 1.2, b = 0.75.
func NewBM25(docs []*Document, fields func(*Document) map[string]string, weights map[string]float64) *BM25 {
	index := &BM25{
		docs:      append([]*Document(nil), docs...),
		weights:   weights,
		fieldTF:   make([]map[string]map[string]int, len(docs)),
		fieldLen:  make([]map[string]int, len(docs)),
		avgLen:    make(map[string]float64),
		docFreq:   make(map[string]int),
		docsTotal: len(docs),
	}

	totalLen := make(map[string]int)
	for position, document := range index.docs {
		index.fieldTF[position] = make(map[string]map[string]int)
		index.fieldLen[position] = make(map[string]int)
		seen := make(map[string]struct{})
		for name, text := range fields(document) {
			tokens := Tokenize(text)
			frequencies := make(map[string]int, len(tokens))
			for _, token := range tokens {
				frequencies[token]++
				seen[token] = struct{}{}
			}
			index.fieldTF[position][name] = frequencies
			index.fieldLen[position][name] = len(tokens)
			totalLen[name] += len(tokens)
		}
		for token := range seen {
			index.docFreq[token]++
		}
	}
	for name, total := range totalLen {
		index.fieldNames = append(index.fieldNames, name)
		if index.docsTotal > 0 {
			index.avgLen[name] = float64(total) / float64(index.docsTotal)
		}
	}
	sort.Strings(index.fieldNames)

	return index
}

// Search returns the k best-scoring documents for query (all matches when
// k <= 0), highest score first, ties by path. Documents that match no query
// term are not returned.
func (index *BM25) Search(query string, k int) []Hit {
	terms := uniqueTokens(Tokenize(query))
	if len(terms) == 0 || index.docsTotal == 0 {
		return nil
	}

	hits := make([]Hit, 0)
	for position, document := range index.docs {
		score := 0.0
		for _, term := range terms {
			df := index.docFreq[term]
			if df == 0 {
				continue
			}
			weighted := 0.0
			for _, name := range index.fieldNames {
				tf := index.fieldTF[position][name][term]
				if tf == 0 {
					continue
				}
				weight, ok := index.weights[name]
				if !ok {
					weight = 1
				}
				norm := 1.0
				if average := index.avgLen[name]; average > 0 {
					norm = 1 - bm25B + bm25B*float64(index.fieldLen[position][name])/average
				}
				weighted += weight * float64(tf) / norm
			}
			if weighted == 0 {
				continue
			}
			idf := math.Log(1 + (float64(index.docsTotal)-float64(df)+0.5)/(float64(df)+0.5))
			score += idf * weighted / (bm25K1 + weighted)
		}
		if score > 0 {
			hits = append(hits, Hit{Doc: document, Score: score})
		}
	}

	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		return hits[i].Doc.Path < hits[j].Doc.Path
	})
	if k > 0 && len(hits) > k {
		hits = hits[:k]
	}

	return hits
}

// DefaultFields extracts the standard retrieval fields of a document: title,
// aliases, criterion, summary, entities, questions and head (the first
// HeadChars characters of the body with wikilinks flattened).
func DefaultFields(d *Document) map[string]string {
	return map[string]string{
		"title":     d.Title,
		"aliases":   strings.Join(d.Aliases, "\n"),
		"criterion": d.Criterion(),
		"summary":   d.Summary(),
		"entities":  strings.Join(d.Entities(), "\n"),
		"questions": strings.Join(d.Questions(), "\n"),
		"head":      Head(d.Body, HeadChars),
	}
}

// DefaultWeights are the field weights used with DefaultFields.
func DefaultWeights() map[string]float64 {
	return map[string]float64{
		"title":     3,
		"aliases":   3,
		"criterion": 2,
		"summary":   1.5,
		"entities":  1.5,
		"questions": 1.2,
		"head":      1,
	}
}

func uniqueTokens(tokens []string) []string {
	seen := make(map[string]struct{}, len(tokens))
	unique := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		unique = append(unique, token)
	}

	return unique
}
