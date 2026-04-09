package adapter

import (
	"errors"
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func TestLanguagesInitialize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		load func() *tree_sitter.Language
	}{
		{
			name: "go",
			load: func() *tree_sitter.Language { return goLanguage() },
		},
		{
			name: "typescript",
			load: func() *tree_sitter.Language { return typeScriptLanguage() },
		},
		{
			name: "tsx",
			load: func() *tree_sitter.Language { return tsxLanguage() },
		},
		{
			name: "javascript",
			load: func() *tree_sitter.Language { return javaScriptLanguage() },
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			language := tt.load()
			if language == nil {
				t.Fatalf("expected non-nil tree-sitter language")
			}

			if language.AbiVersion() == 0 {
				t.Fatalf("expected loaded language ABI version")
			}
		})
	}
}

func TestParsersParseTrivialSources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		language func() *tree_sitter.Language
		source   []byte
		wantKind string
	}{
		{
			name:     "go",
			language: func() *tree_sitter.Language { return goLanguage() },
			source:   []byte("package main\n"),
			wantKind: "source_file",
		},
		{
			name:     "typescript",
			language: func() *tree_sitter.Language { return typeScriptLanguage() },
			source:   []byte("const x = 1;\n"),
			wantKind: "program",
		},
		{
			name:     "tsx",
			language: func() *tree_sitter.Language { return tsxLanguage() },
			source:   []byte("export default function App(){ return <div /> }\n"),
			wantKind: "program",
		},
		{
			name:     "javascript",
			language: func() *tree_sitter.Language { return javaScriptLanguage() },
			source:   []byte("const x = 1;\n"),
			wantKind: "program",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parser, err := newParser(tt.language())
			if err != nil {
				t.Fatalf("newParser() error = %v", err)
			}
			defer parser.Close()

			tree := parser.Parse(tt.source, nil)
			if tree == nil {
				t.Fatal("expected non-nil parse tree")
			}
			defer tree.Close()

			root := tree.RootNode()
			if root == nil {
				t.Fatal("expected non-nil root node")
			}

			if root.Kind() != tt.wantKind {
				t.Fatalf("root kind = %q, want %q", root.Kind(), tt.wantKind)
			}

			if root.HasError() {
				t.Fatalf("expected valid parse tree for %s source", tt.name)
			}
		})
	}
}

func TestNewParserRejectsNilLanguage(t *testing.T) {
	t.Parallel()

	parser, err := newParser(nil)
	if !errors.Is(err, errNilLanguage) {
		t.Fatalf("newParser(nil) error = %v, want %v", err, errNilLanguage)
	}

	if parser != nil {
		t.Fatal("expected nil parser when language is nil")
	}
}
