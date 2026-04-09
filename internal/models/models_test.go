package models

import "testing"

func TestSupportedLanguages(t *testing.T) {
	t.Parallel()

	expected := []SupportedLanguage{LangTS, LangTSX, LangJS, LangJSX, LangGo}

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
