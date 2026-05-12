package models

import (
	"reflect"
	"testing"
)

func TestSupportedLanguages(t *testing.T) {
	t.Parallel()

	expected := []SupportedLanguage{LangTS, LangTSX, LangJS, LangJSX, LangGo, LangRust, LangJava}

	languages := SupportedLanguages()
	if len(languages) != len(expected) {
		t.Fatalf("expected %d supported languages, got %d", len(expected), len(languages))
	}

	for index, language := range languages {
		if language != expected[index] {
			t.Fatalf("language %d: expected %q, got %q", index, expected[index], language)
		}

		if language == "" {
			t.Fatalf("language %d is empty", index)
		}
	}
}

func TestSupportedLanguageNames(t *testing.T) {
	t.Parallel()

	expected := []string{"ts", "tsx", "js", "jsx", "go", "rust", "java"}
	if got := SupportedLanguageNames(); !reflect.DeepEqual(got, expected) {
		t.Fatalf("SupportedLanguageNames() = %#v, want %#v", got, expected)
	}
}

func TestDocumentAndBaseConstants(t *testing.T) {
	t.Parallel()

	if DocRaw != "raw" || DocWiki != "wiki" || DocIndex != "index" {
		t.Fatalf("unexpected document kind constants: %q %q %q", DocRaw, DocWiki, DocIndex)
	}

	if AreaRawCodebase != "raw-codebase" || AreaWikiConcept != "wiki-concept" || AreaWikiIndex != "wiki-index" {
		t.Fatalf("unexpected managed area constants: %q %q %q", AreaRawCodebase, AreaWikiConcept, AreaWikiIndex)
	}

	if ViewTable != "table" || ViewCards != "cards" || ViewList != "list" {
		t.Fatalf("unexpected base view constants: %q %q %q", ViewTable, ViewCards, ViewList)
	}
}
