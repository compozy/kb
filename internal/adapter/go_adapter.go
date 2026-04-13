package adapter

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/compozy/kb/internal/models"
)

const (
	goParseErrorCode = "GO_PARSE_ERROR"

	goSymbolKindPackage   = "package"
	goSymbolKindFunction  = "function"
	goSymbolKindMethod    = "method"
	goSymbolKindStruct    = "struct"
	goSymbolKindInterface = "interface"
	goSymbolKindType      = "type"
)

var (
	_ models.LanguageAdapter = (*GoAdapter)(nil)

	leadingBlockCommentPattern = regexp.MustCompile(`(?s)^\s*/\*\*?[\s\S]*?\*/`)
	leadingLineCommentPattern  = regexp.MustCompile(`^(?:\s*//.*(?:\n|$))+`)
)

// GoAdapter parses Go source files into graph nodes and relations.
type GoAdapter struct{}

type goSymbolMatch struct {
	symbol          models.SymbolNode
	callTargetNames []string
}

type parsedGoFile struct {
	packageName   string
	file          models.GraphFile
	symbolMatches []goSymbolMatch
	externalNodes map[string]models.ExternalNode
	relations     []models.RelationEdge
	diagnostics   []models.StructuredDiagnostic
}

// Supports reports whether the adapter handles the provided language.
func (GoAdapter) Supports(language models.SupportedLanguage) bool {
	return language == models.LangGo
}

// ParseFiles parses Go source files into graph nodes, relations, and diagnostics.
func (adapter GoAdapter) ParseFiles(files []models.ScannedSourceFile, rootPath string) ([]models.ParsedFile, error) {
	return adapter.ParseFilesWithProgress(files, rootPath, nil)
}

// ParseFilesWithProgress parses Go files and reports one progress tick per file.
func (adapter GoAdapter) ParseFilesWithProgress(
	files []models.ScannedSourceFile,
	rootPath string,
	report func(models.ScannedSourceFile),
) ([]models.ParsedFile, error) {
	_ = rootPath

	parser, err := newParser(goLanguage())
	if err != nil {
		return nil, fmt.Errorf("create Go parser: %w", err)
	}
	defer parser.Close()

	orderedFiles := append([]models.ScannedSourceFile(nil), files...)
	sort.Slice(orderedFiles, func(i, j int) bool {
		return orderedFiles[i].RelativePath < orderedFiles[j].RelativePath
	})

	parsedEntries := make([]parsedGoFile, 0, len(orderedFiles))
	packageFunctions := make(map[string]map[string]string)

	for _, file := range orderedFiles {
		if !adapter.Supports(file.Language) {
			return nil, fmt.Errorf("parse %s: unsupported language %q", file.RelativePath, file.Language)
		}

		entry, entryErr := parseGoFile(parser, file)
		if entryErr != nil {
			return nil, entryErr
		}

		if len(entry.diagnostics) == 0 && entry.packageName != "" {
			symbolsByName := packageFunctions[entry.packageName]
			if symbolsByName == nil {
				symbolsByName = make(map[string]string)
				packageFunctions[entry.packageName] = symbolsByName
			}

			for _, symbolMatch := range entry.symbolMatches {
				if symbolMatch.symbol.SymbolKind != goSymbolKindFunction {
					continue
				}

				if _, exists := symbolsByName[symbolMatch.symbol.Name]; !exists {
					symbolsByName[symbolMatch.symbol.Name] = symbolMatch.symbol.ID
				}
			}
		}

		parsedEntries = append(parsedEntries, entry)
		if report != nil {
			report(file)
		}
	}

	parsedFiles := make([]models.ParsedFile, 0, len(parsedEntries))

	for _, entry := range parsedEntries {
		if len(entry.diagnostics) == 0 {
			symbolsByName := packageFunctions[entry.packageName]

			for _, symbolMatch := range entry.symbolMatches {
				if symbolMatch.symbol.SymbolKind != goSymbolKindFunction && symbolMatch.symbol.SymbolKind != goSymbolKindMethod {
					continue
				}

				for _, targetName := range symbolMatch.callTargetNames {
					targetID, exists := symbolsByName[targetName]
					if !exists {
						continue
					}

					entry.relations = append(entry.relations, models.RelationEdge{
						FromID:     symbolMatch.symbol.ID,
						ToID:       targetID,
						Type:       models.RelCalls,
						Confidence: models.ConfidenceSyntactic,
					})
				}
			}
		}

		symbols := make([]models.SymbolNode, 0, len(entry.symbolMatches))
		for _, symbolMatch := range entry.symbolMatches {
			symbols = append(symbols, symbolMatch.symbol)
		}

		parsedFiles = append(parsedFiles, models.ParsedFile{
			File:          entry.file,
			Symbols:       symbols,
			ExternalNodes: sortedExternalNodes(entry.externalNodes),
			Relations:     entry.relations,
			Diagnostics:   entry.diagnostics,
		})
	}

	return parsedFiles, nil
}

func parseGoFile(parser *tree_sitter.Parser, file models.ScannedSourceFile) (parsedGoFile, error) {
	source, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return parsedGoFile{}, fmt.Errorf("read Go source %s: %w", file.RelativePath, err)
	}

	tree := parser.Parse(source, nil)
	if tree == nil {
		return parsedGoFile{}, fmt.Errorf("parse Go source %s: nil syntax tree", file.RelativePath)
	}
	defer tree.Close()

	moduleDoc := extractLeadingComment(string(source))
	fileID := createFileID(file.RelativePath)
	fileNode := models.GraphFile{
		ID:        fileID,
		NodeType:  "file",
		FilePath:  file.RelativePath,
		Language:  file.Language,
		ModuleDoc: moduleDoc,
		SymbolIDs: []string{},
	}

	root := tree.RootNode()
	if root == nil {
		return parsedGoFile{}, fmt.Errorf("parse Go source %s: missing root node", file.RelativePath)
	}

	if root.HasError() {
		return parsedGoFile{
			file: fileNode,
			diagnostics: []models.StructuredDiagnostic{
				createGoParseDiagnostic(file, root.ToSexp()),
			},
			externalNodes: map[string]models.ExternalNode{},
			relations:     []models.RelationEdge{},
		}, nil
	}

	entry := parsedGoFile{
		file:          fileNode,
		externalNodes: make(map[string]models.ExternalNode),
		relations:     []models.RelationEdge{},
		diagnostics:   []models.StructuredDiagnostic{},
	}

	for _, child := range namedChildren(root) {
		child := child

		switch child.Kind() {
		case "package_clause":
			symbol := createGoSymbol(file, &child, source, goSymbolKindPackage, moduleDoc)
			entry.packageName = symbol.Name
			entry.symbolMatches = append(entry.symbolMatches, goSymbolMatch{symbol: symbol})
		case "type_declaration":
			for _, declaredType := range namedChildren(&child) {
				declaredType := declaredType
				if declaredType.Kind() != "type_spec" && declaredType.Kind() != "type_alias" {
					continue
				}

				symbol := createGoSymbol(file, &declaredType, source, getGoSymbolKind(&declaredType), "")
				entry.symbolMatches = append(entry.symbolMatches, goSymbolMatch{symbol: symbol})
			}
		case "function_declaration":
			symbol := createGoSymbol(file, &child, source, goSymbolKindFunction, "")
			entry.symbolMatches = append(entry.symbolMatches, goSymbolMatch{
				symbol:          symbol,
				callTargetNames: extractCallTargetNames(&child, source),
			})
		case "method_declaration":
			symbol := createGoSymbol(file, &child, source, goSymbolKindMethod, "")
			entry.symbolMatches = append(entry.symbolMatches, goSymbolMatch{
				symbol:          symbol,
				callTargetNames: extractCallTargetNames(&child, source),
			})
		case "import_declaration":
			extractImports(fileID, &child, source, entry.externalNodes, &entry.relations)
		}
	}

	for _, symbolMatch := range entry.symbolMatches {
		entry.file.SymbolIDs = append(entry.file.SymbolIDs, symbolMatch.symbol.ID)
		entry.relations = append(entry.relations, models.RelationEdge{
			FromID:     fileID,
			ToID:       symbolMatch.symbol.ID,
			Type:       models.RelContains,
			Confidence: models.ConfidenceSyntactic,
		})

		if symbolMatch.symbol.Exported {
			entry.relations = append(entry.relations, models.RelationEdge{
				FromID:     fileID,
				ToID:       symbolMatch.symbol.ID,
				Type:       models.RelExports,
				Confidence: models.ConfidenceSyntactic,
			})
		}
	}

	return entry, nil
}

func extractImports(
	fileID string,
	importDeclaration *tree_sitter.Node,
	source []byte,
	externalNodes map[string]models.ExternalNode,
	relations *[]models.RelationEdge,
) {
	for _, importSpec := range collectNodesByKind(importDeclaration, "import_spec") {
		importSpec := importSpec

		pathNode := importSpec.ChildByFieldName("path")
		importPath := stripQuotes(textOf(pathNode, source))
		if importPath == "" {
			continue
		}

		externalID := createExternalID(importPath)
		importName := textOf(importSpec.ChildByFieldName("name"), source)
		label := importPath
		if importName != "" {
			label = fmt.Sprintf("%s (%s)", importName, importPath)
		}

		externalNodes[externalID] = models.ExternalNode{
			ID:       externalID,
			NodeType: "external",
			Source:   importPath,
			Label:    label,
		}
		*relations = append(*relations, models.RelationEdge{
			FromID:     fileID,
			ToID:       externalID,
			Type:       models.RelImports,
			Confidence: models.ConfidenceSyntactic,
		})
	}
}

func getGoSymbolKind(typeSpec *tree_sitter.Node) string {
	typeNode := typeSpec.ChildByFieldName("type")
	if typeNode == nil {
		return goSymbolKindType
	}

	switch typeNode.Kind() {
	case "struct_type":
		return goSymbolKindStruct
	case "interface_type":
		return goSymbolKindInterface
	default:
		return goSymbolKindType
	}
}

func createGoSymbol(
	file models.ScannedSourceFile,
	node *tree_sitter.Node,
	source []byte,
	symbolKind string,
	fallbackDoc string,
) models.SymbolNode {
	name := resolveGoSymbolName(node, source, symbolKind)
	signature := formatGoSignature(node, source, symbolKind, name)
	docComment := extractGoDocComment(node, source)
	if docComment == "" {
		docComment = fallbackDoc
	}

	symbol := models.SymbolNode{
		NodeType:   "symbol",
		Name:       name,
		SymbolKind: symbolKind,
		Language:   file.Language,
		FilePath:   file.RelativePath,
		StartLine:  int(node.StartPosition().Row) + 1,
		EndLine:    int(node.EndPosition().Row) + 1,
		Signature:  signature,
		DocComment: docComment,
		Exported:   isGoExported(name),
	}

	if complexity := computeCyclomaticComplexity(node, source, symbolKind); complexity > 0 {
		symbol.CyclomaticComplexity = complexity
	}

	symbol.ID = createSymbolID(symbol)

	return symbol
}

func resolveGoSymbolName(node *tree_sitter.Node, source []byte, symbolKind string) string {
	if symbolKind == goSymbolKindPackage {
		firstNamedChild := node.NamedChild(0)
		if firstNamedChild != nil {
			return textOf(firstNamedChild, source)
		}
		return "anonymous"
	}

	nameNode := node.ChildByFieldName("name")
	name := textOf(nameNode, source)
	if name == "" {
		return "anonymous"
	}

	return name
}

func formatGoSignature(node *tree_sitter.Node, source []byte, symbolKind string, name string) string {
	if symbolKind == goSymbolKindPackage {
		return fmt.Sprintf("package %s", name)
	}

	firstLine := strings.TrimSpace(strings.Split(textOf(node, source), "\n")[0])
	return firstLine
}

func extractGoDocComment(node *tree_sitter.Node, source []byte) string {
	comment := extractAttachedComment(node, source)
	if comment != "" {
		return comment
	}

	parent := node.Parent()
	if parent != nil && parent.Kind() == "type_declaration" {
		return extractAttachedComment(parent, source)
	}

	return ""
}

func extractAttachedComment(node *tree_sitter.Node, source []byte) string {
	commentParts := []string{}
	expectedRow := int(node.StartPosition().Row)

	for sibling := node.PrevSibling(); sibling != nil; sibling = sibling.PrevSibling() {
		if sibling.Kind() == "comment" {
			if expectedRow-int(sibling.EndPosition().Row) > 1 {
				break
			}

			normalized := normalizeComment(textOf(sibling, source))
			if normalized != "" {
				commentParts = append(commentParts, normalized)
			}
			expectedRow = int(sibling.StartPosition().Row)
			continue
		}

		if sibling.IsNamed() {
			break
		}
	}

	if len(commentParts) == 0 {
		return ""
	}

	for left, right := 0, len(commentParts)-1; left < right; left, right = left+1, right-1 {
		commentParts[left], commentParts[right] = commentParts[right], commentParts[left]
	}

	return strings.Join(commentParts, "\n")
}

func extractCallTargetNames(node *tree_sitter.Node, source []byte) []string {
	body := node.ChildByFieldName("body")
	if body == nil {
		return nil
	}

	targets := []string{}
	walkNamed(body, func(current *tree_sitter.Node) bool {
		if current == nil {
			return false
		}

		if current != body && current.Kind() == "func_literal" {
			return false
		}

		if current.Kind() != "call_expression" {
			return true
		}

		functionNode := current.ChildByFieldName("function")
		if functionNode == nil || functionNode.Kind() != "identifier" {
			return true
		}

		targetName := textOf(functionNode, source)
		if targetName != "" {
			targets = append(targets, targetName)
		}

		return true
	})

	return targets
}

func computeCyclomaticComplexity(node *tree_sitter.Node, source []byte, symbolKind string) int {
	if symbolKind != goSymbolKindFunction && symbolKind != goSymbolKindMethod {
		return 0
	}

	body := node.ChildByFieldName("body")
	if body == nil {
		return 1
	}

	complexity := 1
	walkNamed(body, func(current *tree_sitter.Node) bool {
		if current == nil {
			return false
		}

		if current != body && current.Kind() == "func_literal" {
			return false
		}

		switch current.Kind() {
		case "if_statement", "for_statement", "expression_switch_statement", "type_switch_statement", "expression_case", "type_case", "select_statement":
			if current != body {
				complexity++
			}
		case "binary_expression":
			switch operatorText(current, source) {
			case "&&", "||":
				complexity++
			}
		}

		return true
	})

	return complexity
}

func operatorText(node *tree_sitter.Node, source []byte) string {
	operatorNode := node.ChildByFieldName("operator")
	if operatorNode != nil {
		return textOf(operatorNode, source)
	}

	for childIndex := uint(0); childIndex < node.ChildCount(); childIndex++ {
		if node.FieldNameForChild(uint32(childIndex)) != "operator" {
			continue
		}

		return textOf(node.Child(childIndex), source)
	}

	return ""
}

func createGoParseDiagnostic(file models.ScannedSourceFile, detail string) models.StructuredDiagnostic {
	return models.StructuredDiagnostic{
		Code:     goParseErrorCode,
		Severity: models.SeverityError,
		Stage:    models.StageParse,
		Message:  "Failed to parse Go source file",
		FilePath: file.RelativePath,
		Language: file.Language,
		Detail:   detail,
	}
}

func sortedExternalNodes(externalNodes map[string]models.ExternalNode) []models.ExternalNode {
	if len(externalNodes) == 0 {
		return []models.ExternalNode{}
	}

	ids := make([]string, 0, len(externalNodes))
	for id := range externalNodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	nodes := make([]models.ExternalNode, 0, len(ids))
	for _, id := range ids {
		nodes = append(nodes, externalNodes[id])
	}

	return nodes
}

func walkNamed(node *tree_sitter.Node, visit func(*tree_sitter.Node) bool) {
	if node == nil {
		return
	}

	if !visit(node) {
		return
	}

	for _, child := range namedChildren(node) {
		child := child
		walkNamed(&child, visit)
	}
}

func namedChildren(node *tree_sitter.Node) []tree_sitter.Node {
	if node == nil {
		return nil
	}

	cursor := node.Walk()
	defer cursor.Close()

	return node.NamedChildren(cursor)
}

func collectNodesByKind(node *tree_sitter.Node, targetKind string) []tree_sitter.Node {
	nodes := []tree_sitter.Node{}
	walkNamed(node, func(current *tree_sitter.Node) bool {
		if current != nil && current.Kind() == targetKind {
			nodes = append(nodes, *current)
		}
		return true
	})
	return nodes
}

func textOf(node *tree_sitter.Node, source []byte) string {
	if node == nil {
		return ""
	}
	return node.Utf8Text(source)
}

func extractLeadingComment(sourceText string) string {
	if blockMatch := leadingBlockCommentPattern.FindString(sourceText); blockMatch != "" {
		return normalizeComment(blockMatch)
	}

	if lineMatch := leadingLineCommentPattern.FindString(sourceText); lineMatch != "" {
		return normalizeLineComment(lineMatch)
	}

	return ""
}

func normalizeComment(rawComment string) string {
	trimmed := strings.TrimSpace(rawComment)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "/*") {
		trimmed = strings.TrimPrefix(trimmed, "/**")
		trimmed = strings.TrimPrefix(trimmed, "/*")
		trimmed = strings.TrimSuffix(trimmed, "*/")

		lines := strings.Split(trimmed, "\n")
		for index, line := range lines {
			line = strings.TrimLeft(line, " \t")
			line = strings.TrimPrefix(line, "*")
			lines[index] = strings.TrimSpace(line)
		}

		return strings.TrimSpace(strings.Join(lines, "\n"))
	}

	return normalizeLineComment(trimmed)
}

func normalizeLineComment(rawComment string) string {
	lines := strings.Split(rawComment, "\n")
	normalized := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		normalized = append(normalized, line)
	}

	return strings.TrimSpace(strings.Join(normalized, "\n"))
}

func stripQuotes(value string) string {
	if len(value) < 2 {
		return value
	}

	first := value[0]
	last := value[len(value)-1]
	if (first == '\'' || first == '"' || first == '`') && first == last {
		return value[1 : len(value)-1]
	}

	return value
}

func isGoExported(name string) bool {
	for _, r := range name {
		return unicode.IsUpper(r)
	}
	return false
}

func createFileID(filePath string) string {
	return "file:" + filePath
}

func createExternalID(source string) string {
	return "external:" + source
}

func createSymbolID(symbol models.SymbolNode) string {
	return strings.Join([]string{
		"symbol",
		symbol.FilePath,
		slugifySegment(symbol.Name),
		slugifySegment(symbol.SymbolKind),
		fmt.Sprintf("%d", symbol.StartLine),
		fmt.Sprintf("%d", symbol.EndLine),
	}, ":")
}

func slugifySegment(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return "item"
	}

	var builder strings.Builder
	lastWasDash := false

	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastWasDash = false
			continue
		}

		if lastWasDash {
			continue
		}

		builder.WriteByte('-')
		lastWasDash = true
	}

	slug := strings.Trim(builder.String(), "-")
	if slug == "" {
		return "item"
	}

	return slug
}
