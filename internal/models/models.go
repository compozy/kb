package models

// SupportedLanguage identifies a source language handled by the pipeline.
type SupportedLanguage string

const (
	// LangTS represents TypeScript source files.
	LangTS SupportedLanguage = "ts"
	// LangTSX represents TSX source files.
	LangTSX SupportedLanguage = "tsx"
	// LangJS represents JavaScript source files.
	LangJS SupportedLanguage = "js"
	// LangJSX represents JSX source files.
	LangJSX SupportedLanguage = "jsx"
	// LangGo represents Go source files.
	LangGo SupportedLanguage = "go"
)

// SupportedLanguages returns every supported language constant in stable order.
func SupportedLanguages() []SupportedLanguage {
	return []SupportedLanguage{LangTS, LangTSX, LangJS, LangJSX, LangGo}
}

// ScannedSourceFile describes a supported source file discovered in a workspace scan.
type ScannedSourceFile struct {
	AbsolutePath string            `json:"absolutePath"`
	RelativePath string            `json:"relativePath"`
	Language     SupportedLanguage `json:"language"`
}

// ScannedWorkspace contains the files discovered during workspace scanning.
type ScannedWorkspace struct {
	Files           []ScannedSourceFile                       `json:"files"`
	FilesByLanguage map[SupportedLanguage][]ScannedSourceFile `json:"filesByLanguage"`
}
