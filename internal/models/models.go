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

// RelationType identifies a graph edge type between parsed nodes.
type RelationType string

const (
	// RelImports indicates an import edge from a file to another file or module.
	RelImports RelationType = "imports"
	// RelExports indicates a file exports a symbol.
	RelExports RelationType = "exports"
	// RelCalls indicates one symbol calls another symbol.
	RelCalls RelationType = "calls"
	// RelReferences indicates a non-call reference edge.
	RelReferences RelationType = "references"
	// RelDeclares indicates a declaration edge.
	RelDeclares RelationType = "declares"
	// RelContains indicates containment from a file to a symbol.
	RelContains RelationType = "contains"
)

// RelationConfidence captures how strongly a relation is known.
type RelationConfidence string

const (
	// ConfidenceSemantic indicates a semantically resolved edge.
	ConfidenceSemantic RelationConfidence = "semantic"
	// ConfidenceSyntactic indicates a purely syntactic edge.
	ConfidenceSyntactic RelationConfidence = "syntactic"
)

// DiagnosticSeverity indicates the severity of a structured diagnostic.
type DiagnosticSeverity string

const (
	// SeverityWarning is used for recoverable issues.
	SeverityWarning DiagnosticSeverity = "warning"
	// SeverityError is used for blocking issues.
	SeverityError DiagnosticSeverity = "error"
)

// DiagnosticStage indicates the pipeline stage that produced a diagnostic.
type DiagnosticStage string

const (
	// StageScan indicates scan-time diagnostics.
	StageScan DiagnosticStage = "scan"
	// StageParse indicates parse-time diagnostics.
	StageParse DiagnosticStage = "parse"
	// StageRender indicates render-time diagnostics.
	StageRender DiagnosticStage = "render"
	// StageWrite indicates write-time diagnostics.
	StageWrite DiagnosticStage = "write"
	// StageValidate indicates validation-time diagnostics.
	StageValidate DiagnosticStage = "validate"
)

// StructuredDiagnostic is a machine-readable issue emitted during processing.
type StructuredDiagnostic struct {
	Code     string             `json:"code"`
	Severity DiagnosticSeverity `json:"severity"`
	Stage    DiagnosticStage    `json:"stage"`
	Message  string             `json:"message"`
	FilePath string             `json:"filePath,omitempty"`
	Language SupportedLanguage  `json:"language,omitempty"`
	Detail   string             `json:"detail,omitempty"`
}

// GraphFile represents a parsed source file in the graph snapshot.
type GraphFile struct {
	ID        string            `json:"id"`
	NodeType  string            `json:"nodeType"`
	FilePath  string            `json:"filePath"`
	Language  SupportedLanguage `json:"language"`
	ModuleDoc string            `json:"moduleDoc,omitempty"`
	SymbolIDs []string          `json:"symbolIds"`
}

// SymbolNode represents a symbol declaration discovered in a file.
type SymbolNode struct {
	ID                   string            `json:"id"`
	NodeType             string            `json:"nodeType"`
	Name                 string            `json:"name"`
	SymbolKind           string            `json:"symbolKind"`
	Language             SupportedLanguage `json:"language"`
	FilePath             string            `json:"filePath"`
	StartLine            int               `json:"startLine"`
	EndLine              int               `json:"endLine"`
	Signature            string            `json:"signature,omitempty"`
	DocComment           string            `json:"docComment,omitempty"`
	Exported             bool              `json:"exported"`
	CyclomaticComplexity int               `json:"cyclomaticComplexity,omitempty"`
}

// ExternalNode represents an imported module or package outside the graph.
type ExternalNode struct {
	ID       string `json:"id"`
	NodeType string `json:"nodeType"`
	Source   string `json:"source"`
	Label    string `json:"label"`
}

// RelationEdge connects two graph nodes.
type RelationEdge struct {
	FromID     string             `json:"fromId"`
	ToID       string             `json:"toId"`
	Type       RelationType       `json:"type"`
	Confidence RelationConfidence `json:"confidence"`
}

// ParsedFile is the adapter output for one source file.
type ParsedFile struct {
	File          GraphFile              `json:"file"`
	Symbols       []SymbolNode           `json:"symbols"`
	ExternalNodes []ExternalNode         `json:"externalNodes"`
	Relations     []RelationEdge         `json:"relations"`
	Diagnostics   []StructuredDiagnostic `json:"diagnostics"`
}

// GraphSnapshot is the merged graph output across all parsed files.
type GraphSnapshot struct {
	RootPath      string                 `json:"rootPath"`
	Files         []GraphFile            `json:"files"`
	Symbols       []SymbolNode           `json:"symbols"`
	ExternalNodes []ExternalNode         `json:"externalNodes"`
	Relations     []RelationEdge         `json:"relations"`
	Diagnostics   []StructuredDiagnostic `json:"diagnostics"`
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

// LanguageAdapter parses source files for a specific language into graph nodes.
type LanguageAdapter interface {
	Supports(lang SupportedLanguage) bool
	ParseFiles(files []ScannedSourceFile, rootPath string) ([]ParsedFile, error)
}
