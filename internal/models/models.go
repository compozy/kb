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

// DocumentKind identifies the rendered markdown document bucket.
type DocumentKind string

const (
	// DocRaw marks a raw source snapshot document.
	DocRaw DocumentKind = "raw"
	// DocWiki marks a compiled wiki concept article.
	DocWiki DocumentKind = "wiki"
	// DocIndex marks a generated index page.
	DocIndex DocumentKind = "index"
)

// ManagedArea identifies the managed subtree within a generated topic.
type ManagedArea string

const (
	// AreaRawCodebase stores generated raw source snapshots.
	AreaRawCodebase ManagedArea = "raw-codebase"
	// AreaWikiConcept stores generated wiki concept articles.
	AreaWikiConcept ManagedArea = "wiki-concept"
	// AreaWikiIndex stores generated wiki index pages.
	AreaWikiIndex ManagedArea = "wiki-index"
)

// BaseViewType identifies the Obsidian Base view mode.
type BaseViewType string

const (
	// ViewTable renders a table-based Base view.
	ViewTable BaseViewType = "table"
	// ViewCards renders a card-based Base view.
	ViewCards BaseViewType = "cards"
	// ViewList renders a list-based Base view.
	ViewList BaseViewType = "list"
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

// SymbolMetrics stores computed metrics for an individual symbol.
type SymbolMetrics struct {
	BlastRadius            int      `json:"blastRadius"`
	Centrality             float64  `json:"centrality"`
	DirectDependents       int      `json:"directDependents"`
	ExternalReferenceCount int      `json:"externalReferenceCount"`
	IsDeadExport           bool     `json:"isDeadExport"`
	IsLongFunction         bool     `json:"isLongFunction"`
	LOC                    int      `json:"loc"`
	Smells                 []string `json:"smells"`
}

// FileMetrics stores computed metrics for an individual file.
type FileMetrics struct {
	AfferentCoupling      int      `json:"afferentCoupling"`
	EfferentCoupling      int      `json:"efferentCoupling"`
	HasCircularDependency bool     `json:"hasCircularDependency"`
	Instability           float64  `json:"instability"`
	IsEntryPoint          bool     `json:"isEntryPoint"`
	IsGodFile             bool     `json:"isGodFile"`
	IsOrphanFile          bool     `json:"isOrphanFile"`
	Smells                []string `json:"smells"`
}

// DirectoryMetrics stores aggregated metrics for a directory.
type DirectoryMetrics struct {
	AfferentCoupling int     `json:"afferentCoupling"`
	EfferentCoupling int     `json:"efferentCoupling"`
	Instability      float64 `json:"instability"`
}

// MetricsResult contains every computed metrics view for a graph snapshot.
type MetricsResult struct {
	// CircularDependencies stores cyclic file groups, one per strongly connected
	// component with more than one file.
	CircularDependencies [][]string                  `json:"circularDependencies"`
	Directories          map[string]DirectoryMetrics `json:"directories"`
	Files                map[string]FileMetrics      `json:"files"`
	Symbols              map[string]SymbolMetrics    `json:"symbols"`
}

// RenderedDocument is the in-memory representation of a generated markdown file.
type RenderedDocument struct {
	Kind         DocumentKind           `json:"kind"`
	ManagedArea  ManagedArea            `json:"managedArea"`
	RelativePath string                 `json:"relativePath"`
	Frontmatter  map[string]interface{} `json:"frontmatter"`
	Body         string                 `json:"body"`
}

// TopicMetadata captures the derived topic information for a vault render.
type TopicMetadata struct {
	RootPath  string `json:"rootPath"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Domain    string `json:"domain"`
	Today     string `json:"today"`
	VaultPath string `json:"vaultPath"`
	TopicPath string `json:"topicPath"`
}

// GenerateOptions configures a full knowledge-base generation run.
type GenerateOptions struct {
	RootPath        string   `json:"rootPath"`
	VaultPath       string   `json:"vaultPath,omitempty"`
	TopicSlug       string   `json:"topicSlug,omitempty"`
	Title           string   `json:"title,omitempty"`
	Domain          string   `json:"domain,omitempty"`
	IncludePatterns []string `json:"includePatterns,omitempty"`
	ExcludePatterns []string `json:"excludePatterns,omitempty"`
	Semantic        bool     `json:"semantic,omitempty"`
}

// GenerationTimings reports the elapsed wall-clock time for each pipeline stage.
type GenerationTimings struct {
	ScanMillis           int64 `json:"scanMillis"`
	SelectAdaptersMillis int64 `json:"selectAdaptersMillis"`
	ParseMillis          int64 `json:"parseMillis"`
	NormalizeMillis      int64 `json:"normalizeMillis"`
	MetricsMillis        int64 `json:"metricsMillis"`
	RenderMillis         int64 `json:"renderMillis"`
	WriteMillis          int64 `json:"writeMillis"`
	TotalMillis          int64 `json:"totalMillis"`
}

// GenerationSummary reports the outcome of a generation run.
type GenerationSummary struct {
	Command               string                 `json:"command"`
	RootPath              string                 `json:"rootPath"`
	VaultPath             string                 `json:"vaultPath"`
	TopicPath             string                 `json:"topicPath"`
	TopicSlug             string                 `json:"topicSlug"`
	FilesScanned          int                    `json:"filesScanned"`
	FilesParsed           int                    `json:"filesParsed"`
	FilesSkipped          int                    `json:"filesSkipped"`
	SymbolsExtracted      int                    `json:"symbolsExtracted"`
	RelationsEmitted      int                    `json:"relationsEmitted"`
	RawDocumentsWritten   int                    `json:"rawDocumentsWritten"`
	WikiDocumentsWritten  int                    `json:"wikiDocumentsWritten"`
	IndexDocumentsWritten int                    `json:"indexDocumentsWritten"`
	Timings               GenerationTimings      `json:"timings"`
	Diagnostics           []StructuredDiagnostic `json:"diagnostics"`
}

// BaseFilter is a recursive Obsidian Base filter tree.
type BaseFilter struct {
	Expression string       `json:"expression,omitempty"`
	And        []BaseFilter `json:"and,omitempty"`
	Or         []BaseFilter `json:"or,omitempty"`
	Not        *BaseFilter  `json:"not,omitempty"`
}

// BaseProperty configures the display metadata for a Base property.
type BaseProperty struct {
	DisplayName string `json:"displayName"`
}

// BaseGroupBy configures the grouping rule for a Base view.
type BaseGroupBy struct {
	Direction string `json:"direction"`
	Property  string `json:"property"`
}

// BaseView configures a single Obsidian Base view.
type BaseView struct {
	Filters   *BaseFilter       `json:"filters,omitempty"`
	GroupBy   *BaseGroupBy      `json:"groupBy,omitempty"`
	Name      string            `json:"name"`
	Order     []string          `json:"order"`
	Summaries map[string]string `json:"summaries,omitempty"`
	Type      BaseViewType      `json:"type"`
}

// BaseDefinition is the persisted definition of an Obsidian Base file.
type BaseDefinition struct {
	Filters    *BaseFilter             `json:"filters,omitempty"`
	Formulas   map[string]string       `json:"formulas,omitempty"`
	Properties map[string]BaseProperty `json:"properties,omitempty"`
	Views      []BaseView              `json:"views"`
}

// BaseFile describes one generated Obsidian Base file.
type BaseFile struct {
	Definition   BaseDefinition `json:"definition"`
	RelativePath string         `json:"relativePath"`
}
