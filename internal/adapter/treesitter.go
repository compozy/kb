package adapter

import (
	"errors"
	"fmt"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

var errNilLanguage = errors.New("tree-sitter language is nil")

func goLanguage() *tree_sitter.Language {
	return tree_sitter.NewLanguage(tree_sitter_go.Language())
}

func typeScriptLanguage() *tree_sitter.Language {
	return tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())
}

func javaScriptLanguage() *tree_sitter.Language {
	return tree_sitter.NewLanguage(tree_sitter_javascript.Language())
}

func newParser(language *tree_sitter.Language) (*tree_sitter.Parser, error) {
	if language == nil {
		return nil, errNilLanguage
	}

	parser := tree_sitter.NewParser()
	if err := parser.SetLanguage(language); err != nil {
		parser.Close()
		return nil, fmt.Errorf("set parser language: %w", err)
	}

	return parser, nil
}
