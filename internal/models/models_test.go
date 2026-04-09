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
