package adapter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/compozy/kb/internal/models"
)

const (
	tsParseErrorCode = "TS_PARSE_ERROR"

	tsSymbolKindFunction  = "function"
	tsSymbolKindMethod    = "method"
	tsSymbolKindClass     = "class"
	tsSymbolKindInterface = "interface"
	tsSymbolKindTypeAlias = "typeAlias"
	tsSymbolKindEnum      = "enum"
	tsSymbolKindVariable  = "variable"
)

var _ models.LanguageAdapter = (*TSAdapter)(nil)

var supportedJSImportExtensions = []string{".ts", ".tsx", ".js", ".jsx"}

// TSAdapter parses TypeScript, TSX, JavaScript, and JSX source files.
type TSAdapter struct{}

type tsCallTarget struct {
	exportName string
	localName  string
	kind       string
}

type tsImportBinding struct {
	importedName   string
	localName      string
	sourceFilePath string
	kind           string
}

type tsLocalExport struct {
	exportName string
	localName  string
}

type tsReExport struct {
	exportName     string
	importedName   string
	sourceFilePath string
}

type tsSymbolMatch struct {
	symbol          models.SymbolNode
	exportNames     map[string]struct{}
	callTargets     []tsCallTarget
	anchorStartLine int
}

type parsedTSFile struct {
	file           models.ScannedSourceFile
	fileNode       models.GraphFile
	symbolMatches  []tsSymbolMatch
	importBindings []tsImportBinding
	localExports   []tsLocalExport
	reExports      []tsReExport
	externalNodes  map[string]models.ExternalNode
	relations      []models.RelationEdge
	diagnostics    []models.StructuredDiagnostic
}

// Supports reports whether the adapter handles the provided language.
func (TSAdapter) Supports(language models.SupportedLanguage) bool {
	switch language {
	case models.LangTS, models.LangTSX, models.LangJS, models.LangJSX:
		return true
	default:
		return false
	}
}

// ParseFiles parses TS/JS source files into graph nodes, relations, and diagnostics.
func (adapter TSAdapter) ParseFiles(files []models.ScannedSourceFile, rootPath string) ([]models.ParsedFile, error) {
	return adapter.ParseFilesWithProgress(files, rootPath, nil)
}

// ParseFilesWithProgress parses TS/JS files and reports one progress tick per file.
func (adapter TSAdapter) ParseFilesWithProgress(
	files []models.ScannedSourceFile,
	rootPath string,
	report func(models.ScannedSourceFile),
) ([]models.ParsedFile, error) {
	_ = rootPath

	if len(files) == 0 {
		return []models.ParsedFile{}, nil
	}

	orderedFiles := append([]models.ScannedSourceFile(nil), files...)
	sort.Slice(orderedFiles, func(i, j int) bool {
		return orderedFiles[i].RelativePath < orderedFiles[j].RelativePath
	})

	fileByAbsolutePath := make(map[string]models.ScannedSourceFile, len(orderedFiles))
	for _, file := range orderedFiles {
		normalizedPath, err := normalizeAbsolutePath(file.AbsolutePath)
		if err != nil {
			return nil, fmt.Errorf("normalize path %s: %w", file.RelativePath, err)
		}
		fileByAbsolutePath[normalizedPath] = file
	}

	parsedEntries := make([]parsedTSFile, 0, len(orderedFiles))
	for _, file := range orderedFiles {
		if !adapter.Supports(file.Language) {
			return nil, fmt.Errorf("parse %s: unsupported language %q", file.RelativePath, file.Language)
		}

		entry, err := parseTSFile(file, fileByAbsolutePath)
		if err != nil {
			return nil, err
		}
		parsedEntries = append(parsedEntries, entry)
		if report != nil {
			report(file)
		}
	}

	localSymbolsByFilePath := make(map[string]map[string]string, len(parsedEntries))
	exportedSymbolsByFilePath := make(map[string]map[string]string, len(parsedEntries))
	for _, entry := range parsedEntries {
		localSymbolsByName := make(map[string]string)
		exportedSymbolsByName := make(map[string]string)

		for _, symbolMatch := range entry.symbolMatches {
			localSymbolsByName[symbolMatch.symbol.Name] = symbolMatch.symbol.ID
			for exportName := range symbolMatch.exportNames {
				exportedSymbolsByName[exportName] = symbolMatch.symbol.ID
			}
		}

		localSymbolsByFilePath[entry.file.RelativePath] = localSymbolsByName
		exportedSymbolsByFilePath[entry.file.RelativePath] = exportedSymbolsByName
	}

	resolveLocalExports(parsedEntries, exportedSymbolsByFilePath, localSymbolsByFilePath)
	resolveReExports(parsedEntries, exportedSymbolsByFilePath)

	parsedFiles := make([]models.ParsedFile, 0, len(parsedEntries))
	for _, entry := range parsedEntries {
		relationKeys := make(map[string]struct{}, len(entry.relations))
		for _, relation := range entry.relations {
			relationKeys[relationKey(relation)] = struct{}{}
		}

		localSymbolsByName := localSymbolsByFilePath[entry.file.RelativePath]
		directBindings := make(map[string]tsImportBinding)
		namespaceBindings := make(map[string]tsImportBinding)

		for _, binding := range entry.importBindings {
			switch binding.kind {
			case "namespace":
				namespaceBindings[binding.localName] = binding
				for exportName, targetID := range exportedSymbolsByFilePath[binding.sourceFilePath] {
					if exportName == "default" || targetID == "" {
						continue
					}

					pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
						FromID:     entry.fileNode.ID,
						ToID:       targetID,
						Type:       models.RelReferences,
						Confidence: models.ConfidenceSyntactic,
					})
				}
			default:
				directBindings[binding.localName] = binding
				targetID := exportedSymbolsByFilePath[binding.sourceFilePath][binding.importedName]
				if targetID == "" {
					continue
				}

				pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
					FromID:     entry.fileNode.ID,
					ToID:       targetID,
					Type:       models.RelReferences,
					Confidence: models.ConfidenceSyntactic,
				})
			}
		}

		for _, reExport := range entry.reExports {
			if reExport.importedName == "*" {
				for exportName, targetID := range exportedSymbolsByFilePath[reExport.sourceFilePath] {
					if exportName == "default" || targetID == "" {
						continue
					}

					pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
						FromID:     entry.fileNode.ID,
						ToID:       targetID,
						Type:       models.RelExports,
						Confidence: models.ConfidenceSyntactic,
					})
				}
				continue
			}

			targetID := exportedSymbolsByFilePath[reExport.sourceFilePath][reExport.importedName]
			if targetID == "" {
				continue
			}

			pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
				FromID:     entry.fileNode.ID,
				ToID:       targetID,
				Type:       models.RelExports,
				Confidence: models.ConfidenceSyntactic,
			})
		}

		for _, symbolMatch := range entry.symbolMatches {
			for _, callTarget := range symbolMatch.callTargets {
				if callTarget.kind == "identifier" {
					if targetID := localSymbolsByName[callTarget.localName]; targetID != "" && targetID != symbolMatch.symbol.ID {
						pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
							FromID:     symbolMatch.symbol.ID,
							ToID:       targetID,
							Type:       models.RelCalls,
							Confidence: models.ConfidenceSyntactic,
						})
						continue
					}

					binding, exists := directBindings[callTarget.localName]
					if !exists {
						continue
					}

					targetID := exportedSymbolsByFilePath[binding.sourceFilePath][binding.importedName]
					if targetID == "" {
						continue
					}

					pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
						FromID:     symbolMatch.symbol.ID,
						ToID:       targetID,
						Type:       models.RelCalls,
						Confidence: models.ConfidenceSyntactic,
					})
					continue
				}

				binding, exists := namespaceBindings[callTarget.localName]
				if !exists || callTarget.exportName == "" {
					continue
				}

				targetID := exportedSymbolsByFilePath[binding.sourceFilePath][callTarget.exportName]
				if targetID == "" {
					continue
				}

				pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
					FromID:     symbolMatch.symbol.ID,
					ToID:       targetID,
					Type:       models.RelCalls,
					Confidence: models.ConfidenceSyntactic,
				})
			}
		}

		symbols := make([]models.SymbolNode, 0, len(entry.symbolMatches))
		for _, symbolMatch := range entry.symbolMatches {
			symbols = append(symbols, symbolMatch.symbol)
		}

		parsedFiles = append(parsedFiles, models.ParsedFile{
			File:          entry.fileNode,
			Symbols:       symbols,
			ExternalNodes: sortedExternalNodes(entry.externalNodes),
			Relations:     entry.relations,
			Diagnostics:   entry.diagnostics,
		})
	}

	sort.Slice(parsedFiles, func(i, j int) bool {
		return parsedFiles[i].File.FilePath < parsedFiles[j].File.FilePath
	})

	return parsedFiles, nil
}

func parseTSFile(
	file models.ScannedSourceFile,
	fileByAbsolutePath map[string]models.ScannedSourceFile,
) (parsedTSFile, error) {
	source, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return parsedTSFile{}, fmt.Errorf("read TS/JS source %s: %w", file.RelativePath, err)
	}

	language := selectTSLanguage(file.Language)
	if language == nil {
		return parsedTSFile{}, fmt.Errorf("select parser language for %s: unsupported language %q", file.RelativePath, file.Language)
	}

	parser, err := newParser(language)
	if err != nil {
		return parsedTSFile{}, fmt.Errorf("create parser for %s: %w", file.RelativePath, err)
	}
	defer parser.Close()

	tree := parser.Parse(source, nil)
	if tree == nil {
		return parsedTSFile{}, fmt.Errorf("parse TS/JS source %s: nil syntax tree", file.RelativePath)
	}
	defer tree.Close()

	fileID := createFileID(file.RelativePath)
	entry := parsedTSFile{
		file: file,
		fileNode: models.GraphFile{
			ID:        fileID,
			NodeType:  "file",
			FilePath:  file.RelativePath,
			Language:  file.Language,
			ModuleDoc: extractLeadingComment(string(source)),
			SymbolIDs: []string{},
		},
		externalNodes: map[string]models.ExternalNode{},
		relations:     []models.RelationEdge{},
		diagnostics:   []models.StructuredDiagnostic{},
	}

	root := tree.RootNode()
	if root == nil {
		return parsedTSFile{}, fmt.Errorf("parse TS/JS source %s: missing root node", file.RelativePath)
	}

	if root.HasError() {
		entry.diagnostics = append(entry.diagnostics, createTSParseDiagnostic(file, root.ToSexp()))
		return entry, nil
	}

	for _, child := range namedChildren(root) {
		child := child

		switch child.Kind() {
		case "function_declaration":
			entry.symbolMatches = append(entry.symbolMatches, createTSSymbolMatch(file, &child, &child, source, nil))
		case "class_declaration":
			entry.symbolMatches = append(entry.symbolMatches, createTSSymbolMatch(file, &child, &child, source, nil))
			entry.symbolMatches = append(entry.symbolMatches, extractClassMethodSymbols(file, &child, source)...)
		case "interface_declaration", "type_alias_declaration", "enum_declaration":
			entry.symbolMatches = append(entry.symbolMatches, createTSSymbolMatch(file, &child, &child, source, nil))
		case "lexical_declaration", "variable_declaration":
			entry.symbolMatches = append(entry.symbolMatches, extractVariableSymbols(file, &child, &child, source, nil)...)
			extractRequireBindings(file, &entry, &child, source, fileByAbsolutePath)
		case "import_statement":
			extractTSImports(file, &entry, &child, source, fileByAbsolutePath)
		case "export_statement":
			extractTSExports(file, &entry, &child, source, fileByAbsolutePath)
		case "expression_statement":
			extractCommonJSExports(file, &entry, &child, source)
		}
	}

	applyLocalExportState(&entry)
	return entry, nil
}

func selectTSLanguage(language models.SupportedLanguage) *tree_sitter.Language {
	switch language {
	case models.LangTS:
		return typeScriptLanguage()
	case models.LangTSX:
		return tsxLanguage()
	case models.LangJS, models.LangJSX:
		return javaScriptLanguage()
	default:
		return nil
	}
}

func extractTSImports(
	file models.ScannedSourceFile,
	entry *parsedTSFile,
	node *tree_sitter.Node,
	source []byte,
	fileByAbsolutePath map[string]models.ScannedSourceFile,
) {
	sourceNode := node.ChildByFieldName("source")
	moduleSpecifier := stripQuotes(textOf(sourceNode, source))
	if moduleSpecifier == "" {
		return
	}

	resolvedFile, resolved := resolveRelativeImportFile(file, moduleSpecifier, fileByAbsolutePath)
	targetID := ""
	if resolved {
		targetID = createFileID(resolvedFile.RelativePath)
	} else {
		targetID = createExternalID(moduleSpecifier)
		entry.externalNodes[targetID] = models.ExternalNode{
			ID:       targetID,
			NodeType: "external",
			Source:   moduleSpecifier,
			Label:    moduleSpecifier,
		}
	}

	entry.relations = append(entry.relations, models.RelationEdge{
		FromID:     entry.fileNode.ID,
		ToID:       targetID,
		Type:       models.RelImports,
		Confidence: models.ConfidenceSyntactic,
	})

	if !resolved {
		return
	}

	importClause := node.NamedChild(0)
	if importClause == nil || importClause.Kind() != "import_clause" {
		return
	}

	for _, clauseChild := range namedChildren(importClause) {
		clauseChild := clauseChild

		switch clauseChild.Kind() {
		case "identifier":
			entry.importBindings = append(entry.importBindings, tsImportBinding{
				importedName:   "default",
				localName:      textOf(&clauseChild, source),
				sourceFilePath: resolvedFile.RelativePath,
				kind:           "default",
			})
		case "namespace_import":
			namespaceName := textOf(clauseChild.NamedChild(0), source)
			if namespaceName == "" {
				continue
			}
			entry.importBindings = append(entry.importBindings, tsImportBinding{
				importedName:   "*",
				localName:      namespaceName,
				sourceFilePath: resolvedFile.RelativePath,
				kind:           "namespace",
			})
		case "named_imports":
			for _, importSpecifier := range namedChildren(&clauseChild) {
				importSpecifier := importSpecifier
				if importSpecifier.Kind() != "import_specifier" {
					continue
				}

				importedName := textOf(importSpecifier.ChildByFieldName("name"), source)
				localName := importedName
				if alias := textOf(importSpecifier.ChildByFieldName("alias"), source); alias != "" {
					localName = alias
				}

				if importedName == "" || localName == "" {
					continue
				}

				entry.importBindings = append(entry.importBindings, tsImportBinding{
					importedName:   importedName,
					localName:      localName,
					sourceFilePath: resolvedFile.RelativePath,
					kind:           "named",
				})
			}
		}
	}
}

func extractTSExports(
	file models.ScannedSourceFile,
	entry *parsedTSFile,
	node *tree_sitter.Node,
	source []byte,
	fileByAbsolutePath map[string]models.ScannedSourceFile,
) {
	if declaration := node.ChildByFieldName("declaration"); declaration != nil {
		exportNames := map[string]struct{}{}
		if isDefaultExport(node, source) {
			exportNames["default"] = struct{}{}
		} else {
			for _, localName := range declarationExportNames(declaration, source) {
				exportNames[localName] = struct{}{}
			}
		}

		switch declaration.Kind() {
		case "function_declaration", "class_declaration", "interface_declaration", "type_alias_declaration", "enum_declaration":
			entry.symbolMatches = append(entry.symbolMatches, createTSSymbolMatch(file, declaration, node, source, exportNames))
			if declaration.Kind() == "class_declaration" {
				entry.symbolMatches = append(entry.symbolMatches, extractClassMethodSymbols(file, declaration, source)...)
			}
		case "lexical_declaration", "variable_declaration":
			entry.symbolMatches = append(entry.symbolMatches, extractVariableSymbols(file, declaration, node, source, exportNames)...)
			extractRequireBindings(file, entry, declaration, source, fileByAbsolutePath)
		}
		return
	}

	sourceNode := node.ChildByFieldName("source")
	moduleSpecifier := stripQuotes(textOf(sourceNode, source))
	resolvedFile, resolved := resolveRelativeImportFile(file, moduleSpecifier, fileByAbsolutePath)
	if moduleSpecifier != "" {
		targetID := ""
		if resolved {
			targetID = createFileID(resolvedFile.RelativePath)
		} else {
			targetID = createExternalID(moduleSpecifier)
			entry.externalNodes[targetID] = models.ExternalNode{
				ID:       targetID,
				NodeType: "external",
				Source:   moduleSpecifier,
				Label:    moduleSpecifier,
			}
		}

		entry.relations = append(entry.relations, models.RelationEdge{
			FromID:     entry.fileNode.ID,
			ToID:       targetID,
			Type:       models.RelImports,
			Confidence: models.ConfidenceSyntactic,
		})
	}

	exportClause := findChildByKind(node, "export_clause")
	if exportClause == nil {
		if resolved {
			entry.reExports = append(entry.reExports, tsReExport{
				exportName:     "*",
				importedName:   "*",
				sourceFilePath: resolvedFile.RelativePath,
			})
		}
		return
	}

	for _, exportSpecifier := range namedChildren(exportClause) {
		exportSpecifier := exportSpecifier
		if exportSpecifier.Kind() != "export_specifier" {
			continue
		}

		localName := textOf(exportSpecifier.ChildByFieldName("name"), source)
		exportName := localName
		if alias := textOf(exportSpecifier.ChildByFieldName("alias"), source); alias != "" {
			exportName = alias
		}

		if localName == "" || exportName == "" {
			continue
		}

		if resolved {
			entry.reExports = append(entry.reExports, tsReExport{
				exportName:     exportName,
				importedName:   localName,
				sourceFilePath: resolvedFile.RelativePath,
			})
			continue
		}

		entry.localExports = append(entry.localExports, tsLocalExport{
			exportName: exportName,
			localName:  localName,
		})
	}
}

func extractCommonJSExports(
	file models.ScannedSourceFile,
	entry *parsedTSFile,
	node *tree_sitter.Node,
	source []byte,
) {
	assignment := node.NamedChild(0)
	if assignment == nil || assignment.Kind() != "assignment_expression" {
		return
	}

	left := assignment.ChildByFieldName("left")
	right := assignment.ChildByFieldName("right")
	if left == nil || right == nil {
		return
	}

	objectName, propertyName, isModuleExports, ok := matchCommonJSExportTarget(left, source)
	if !ok {
		return
	}

	switch {
	case isModuleExports && propertyName == "":
		switch right.Kind() {
		case "identifier":
			localName := textOf(right, source)
			if localName != "" {
				entry.localExports = append(entry.localExports, tsLocalExport{
					localName:  localName,
					exportName: "default",
				})
			}
		case "function_expression", "class":
			entry.symbolMatches = append(entry.symbolMatches, createTSSymbolMatch(file, right, node, source, map[string]struct{}{"default": {}}))
		}
	case propertyName != "":
		switch right.Kind() {
		case "identifier":
			localName := textOf(right, source)
			if localName != "" {
				entry.localExports = append(entry.localExports, tsLocalExport{
					localName:  localName,
					exportName: propertyName,
				})
			}
		}
	case objectName == "exports" && propertyName != "":
		switch right.Kind() {
		case "identifier":
			localName := textOf(right, source)
			if localName != "" {
				entry.localExports = append(entry.localExports, tsLocalExport{
					localName:  localName,
					exportName: propertyName,
				})
			}
		}
	}
}

func extractRequireBindings(
	file models.ScannedSourceFile,
	entry *parsedTSFile,
	node *tree_sitter.Node,
	source []byte,
	fileByAbsolutePath map[string]models.ScannedSourceFile,
) {
	for _, declarator := range collectNodesByKind(node, "variable_declarator") {
		declarator := declarator
		valueNode := declarator.ChildByFieldName("value")
		if valueNode == nil || valueNode.Kind() != "call_expression" {
			continue
		}

		functionNode := valueNode.ChildByFieldName("function")
		if functionNode == nil || functionNode.Kind() != "identifier" || textOf(functionNode, source) != "require" {
			continue
		}

		moduleSpecifier := ""
		argumentsNode := valueNode.ChildByFieldName("arguments")
		if argumentsNode != nil {
			for _, argumentChild := range namedChildren(argumentsNode) {
				argumentChild := argumentChild
				if argumentChild.Kind() == "string" {
					moduleSpecifier = stripQuotes(textOf(&argumentChild, source))
					break
				}
			}
		}

		if moduleSpecifier == "" {
			continue
		}

		resolvedFile, resolved := resolveRelativeImportFile(file, moduleSpecifier, fileByAbsolutePath)
		targetID := ""
		if resolved {
			targetID = createFileID(resolvedFile.RelativePath)
		} else {
			targetID = createExternalID(moduleSpecifier)
			entry.externalNodes[targetID] = models.ExternalNode{
				ID:       targetID,
				NodeType: "external",
				Source:   moduleSpecifier,
				Label:    moduleSpecifier,
			}
		}

		entry.relations = append(entry.relations, models.RelationEdge{
			FromID:     entry.fileNode.ID,
			ToID:       targetID,
			Type:       models.RelImports,
			Confidence: models.ConfidenceSyntactic,
		})

		if !resolved {
			continue
		}

		nameNode := declarator.ChildByFieldName("name")
		if nameNode == nil {
			continue
		}

		switch nameNode.Kind() {
		case "identifier":
			localName := textOf(nameNode, source)
			if localName == "" {
				continue
			}

			entry.importBindings = append(entry.importBindings,
				tsImportBinding{
					importedName:   "default",
					localName:      localName,
					sourceFilePath: resolvedFile.RelativePath,
					kind:           "default",
				},
				tsImportBinding{
					importedName:   "*",
					localName:      localName,
					sourceFilePath: resolvedFile.RelativePath,
					kind:           "namespace",
				},
			)
		case "object_pattern":
			for _, patternChild := range namedChildren(nameNode) {
				patternChild := patternChild

				switch patternChild.Kind() {
				case "pair_pattern":
					importedName := textOf(patternChild.ChildByFieldName("key"), source)
					localName := textOf(patternChild.ChildByFieldName("value"), source)
					if importedName == "" || localName == "" {
						continue
					}

					entry.importBindings = append(entry.importBindings, tsImportBinding{
						importedName:   importedName,
						localName:      localName,
						sourceFilePath: resolvedFile.RelativePath,
						kind:           "named",
					})
				case "shorthand_property_identifier_pattern":
					importedName := textOf(&patternChild, source)
					if importedName == "" {
						continue
					}

					entry.importBindings = append(entry.importBindings, tsImportBinding{
						importedName:   importedName,
						localName:      importedName,
						sourceFilePath: resolvedFile.RelativePath,
						kind:           "named",
					})
				}
			}
		}
	}
}

func createTSSymbolMatch(
	file models.ScannedSourceFile,
	node *tree_sitter.Node,
	anchorNode *tree_sitter.Node,
	source []byte,
	exportNames map[string]struct{},
) tsSymbolMatch {
	symbol := createTSSymbol(file, node, anchorNode, source, exportNames)
	return tsSymbolMatch{
		symbol:          symbol,
		exportNames:     cloneStringSet(exportNames),
		callTargets:     collectTSCallTargets(node, source),
		anchorStartLine: int(anchorNode.StartPosition().Row) + 1,
	}
}

func extractClassMethodSymbols(
	file models.ScannedSourceFile,
	classNode *tree_sitter.Node,
	source []byte,
) []tsSymbolMatch {
	body := classNode.ChildByFieldName("body")
	if body == nil {
		return nil
	}

	methods := []tsSymbolMatch{}
	for _, child := range namedChildren(body) {
		child := child
		if child.Kind() != "method_definition" {
			continue
		}

		methods = append(methods, createTSSymbolMatch(file, &child, &child, source, nil))
	}

	return methods
}

func extractVariableSymbols(
	file models.ScannedSourceFile,
	declarationNode *tree_sitter.Node,
	anchorNode *tree_sitter.Node,
	source []byte,
	exportNames map[string]struct{},
) []tsSymbolMatch {
	symbols := []tsSymbolMatch{}

	for _, declarator := range collectNodesByKind(declarationNode, "variable_declarator") {
		declarator := declarator
		if variableName := resolveTSVariableName(&declarator, source); variableName == "" {
			continue
		}

		symbols = append(symbols, createTSSymbolMatch(file, &declarator, anchorNode, source, exportNames))
	}

	return symbols
}

func createTSSymbol(
	file models.ScannedSourceFile,
	node *tree_sitter.Node,
	anchorNode *tree_sitter.Node,
	source []byte,
	exportNames map[string]struct{},
) models.SymbolNode {
	symbolKind := getTSSymbolKind(node)
	name := resolveTSSymbolName(node, source, symbolKind)
	signature := formatTSSignature(node, source, symbolKind, name)
	docComment := extractTSDocComment(node, anchorNode, source)

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
		Exported:   len(exportNames) > 0,
	}

	if complexity := computeTSCyclomaticComplexity(node, source, symbolKind); complexity > 0 {
		symbol.CyclomaticComplexity = complexity
	}

	symbol.ID = createSymbolID(symbol)
	return symbol
}

func getTSSymbolKind(node *tree_sitter.Node) string {
	switch node.Kind() {
	case "function_declaration", "function_expression", "arrow_function":
		return tsSymbolKindFunction
	case "method_definition":
		return tsSymbolKindMethod
	case "class_declaration", "class":
		return tsSymbolKindClass
	case "interface_declaration":
		return tsSymbolKindInterface
	case "type_alias_declaration":
		return tsSymbolKindTypeAlias
	case "enum_declaration":
		return tsSymbolKindEnum
	case "variable_declarator":
		return tsSymbolKindVariable
	default:
		return "symbol"
	}
}

func resolveTSSymbolName(node *tree_sitter.Node, source []byte, symbolKind string) string {
	if node == nil {
		return "anonymous"
	}

	switch symbolKind {
	case tsSymbolKindMethod:
		nameNode := node.ChildByFieldName("name")
		if name := textOf(nameNode, source); name != "" {
			return name
		}
	case tsSymbolKindClass, tsSymbolKindInterface, tsSymbolKindTypeAlias, tsSymbolKindEnum:
		nameNode := node.ChildByFieldName("name")
		if name := textOf(nameNode, source); name != "" {
			return name
		}
	case tsSymbolKindVariable:
		if name := resolveTSVariableName(node, source); name != "" {
			return name
		}
	default:
		nameNode := node.ChildByFieldName("name")
		if name := textOf(nameNode, source); name != "" {
			return name
		}
	}

	if firstNamedChild := node.NamedChild(0); firstNamedChild != nil {
		if name := textOf(firstNamedChild, source); name != "" {
			return name
		}
	}

	return "anonymous"
}

func resolveTSVariableName(node *tree_sitter.Node, source []byte) string {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return ""
	}

	if nameNode.Kind() == "identifier" {
		return textOf(nameNode, source)
	}

	return ""
}

func formatTSSignature(node *tree_sitter.Node, source []byte, symbolKind, name string) string {
	switch symbolKind {
	case tsSymbolKindFunction:
		parameters := textOf(node.ChildByFieldName("parameters"), source)
		returnType := formatTSReturnType(node, source)
		if returnType != "" {
			return fmt.Sprintf("function %s%s: %s", name, parameters, returnType)
		}
		return fmt.Sprintf("function %s%s", name, parameters)
	case tsSymbolKindMethod:
		parameters := textOf(node.ChildByFieldName("parameters"), source)
		returnType := formatTSReturnType(node, source)
		if returnType != "" {
			return fmt.Sprintf("method %s%s: %s", name, parameters, returnType)
		}
		return fmt.Sprintf("method %s%s", name, parameters)
	case tsSymbolKindClass:
		return fmt.Sprintf("class %s", name)
	case tsSymbolKindInterface:
		return fmt.Sprintf("interface %s", name)
	case tsSymbolKindTypeAlias:
		typeNode := node.NamedChild(1)
		if typeNode != nil {
			return fmt.Sprintf("type %s = %s", name, textOf(typeNode, source))
		}
		return fmt.Sprintf("type %s", name)
	case tsSymbolKindEnum:
		return fmt.Sprintf("enum %s", name)
	case tsSymbolKindVariable:
		line := strings.TrimSpace(strings.Split(textOf(node, source), "\n")[0])
		if line != "" {
			return "const " + name + formatTSVariableTypeSuffix(node, source)
		}
		return "const " + name
	default:
		return strings.TrimSpace(strings.Split(textOf(node, source), "\n")[0])
	}
}

func formatTSReturnType(node *tree_sitter.Node, source []byte) string {
	typeNode := node.ChildByFieldName("return_type")
	if typeNode == nil {
		typeNode = node.ChildByFieldName("type")
	}
	if typeNode == nil {
		for _, child := range namedChildren(node) {
			child := child
			if child.Kind() == "type_annotation" {
				return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(textOf(&child, source)), ":"))
			}
		}
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(textOf(typeNode, source)), ":"))
}

func formatTSVariableTypeSuffix(node *tree_sitter.Node, source []byte) string {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return ""
	}

	for _, child := range namedChildren(nameNode) {
		child := child
		if child.Kind() != "type_annotation" {
			continue
		}

		return ": " + strings.TrimPrefix(strings.TrimSpace(textOf(&child, source)), ":")
	}

	return ""
}

func extractTSDocComment(node *tree_sitter.Node, anchorNode *tree_sitter.Node, source []byte) string {
	if comment := extractAttachedComment(node, source); comment != "" {
		return comment
	}

	if anchorNode != nil && anchorNode != node {
		if comment := extractAttachedComment(anchorNode, source); comment != "" {
			return comment
		}
	}

	return ""
}

func collectTSCallTargets(node *tree_sitter.Node, source []byte) []tsCallTarget {
	targets := []tsCallTarget{}

	walkNamed(node, func(current *tree_sitter.Node) bool {
		if current == nil || current.Kind() != "call_expression" {
			return true
		}

		functionNode := current.ChildByFieldName("function")
		if functionNode == nil {
			return true
		}

		switch functionNode.Kind() {
		case "identifier":
			localName := textOf(functionNode, source)
			if localName != "" {
				targets = append(targets, tsCallTarget{
					localName: localName,
					kind:      "identifier",
				})
			}
		case "member_expression":
			objectNode := functionNode.ChildByFieldName("object")
			propertyNode := functionNode.ChildByFieldName("property")
			if objectNode == nil || propertyNode == nil {
				return true
			}
			if objectNode.Kind() != "identifier" || propertyNode.Kind() != "property_identifier" {
				return true
			}

			localName := textOf(objectNode, source)
			exportName := textOf(propertyNode, source)
			if localName == "" || exportName == "" {
				return true
			}

			targets = append(targets, tsCallTarget{
				localName:  localName,
				exportName: exportName,
				kind:       "namespace",
			})
		}

		return true
	})

	return targets
}

func computeTSCyclomaticComplexity(node *tree_sitter.Node, source []byte, symbolKind string) int {
	if symbolKind != tsSymbolKindFunction && symbolKind != tsSymbolKindMethod {
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

		if current != body && isTSExecutableBoundary(current.Kind()) {
			return false
		}

		switch current.Kind() {
		case "catch_clause", "ternary_expression", "do_statement", "for_in_statement", "for_statement", "if_statement", "switch_case", "while_statement":
			if current != body {
				complexity++
			}
		case "binary_expression":
			switch operatorText(current, source) {
			case "&&", "||", "??":
				complexity++
			}
		}

		return true
	})

	return complexity
}

func isTSExecutableBoundary(kind string) bool {
	switch kind {
	case "arrow_function", "function_declaration", "function_expression", "generator_function", "generator_function_declaration", "method_definition":
		return true
	default:
		return false
	}
}

func applyLocalExportState(entry *parsedTSFile) {
	for index := range entry.symbolMatches {
		exportNames := entry.symbolMatches[index].exportNames
		if exportNames == nil {
			exportNames = map[string]struct{}{}
		}

		for _, localExport := range entry.localExports {
			if localExport.localName != entry.symbolMatches[index].symbol.Name {
				continue
			}
			exportNames[localExport.exportName] = struct{}{}
		}

		entry.symbolMatches[index].exportNames = exportNames
		entry.symbolMatches[index].symbol.Exported = len(exportNames) > 0
		entry.fileNode.SymbolIDs = append(entry.fileNode.SymbolIDs, entry.symbolMatches[index].symbol.ID)
		entry.relations = append(entry.relations, models.RelationEdge{
			FromID:     entry.fileNode.ID,
			ToID:       entry.symbolMatches[index].symbol.ID,
			Type:       models.RelContains,
			Confidence: models.ConfidenceSyntactic,
		})

		if entry.symbolMatches[index].symbol.Exported {
			entry.relations = append(entry.relations, models.RelationEdge{
				FromID:     entry.fileNode.ID,
				ToID:       entry.symbolMatches[index].symbol.ID,
				Type:       models.RelExports,
				Confidence: models.ConfidenceSyntactic,
			})
		}
	}
}

func resolveLocalExports(
	parsedEntries []parsedTSFile,
	exportedSymbolsByFilePath map[string]map[string]string,
	localSymbolsByFilePath map[string]map[string]string,
) {
	for _, entry := range parsedEntries {
		exportMap := exportedSymbolsByFilePath[entry.file.RelativePath]
		localSymbolsByName := localSymbolsByFilePath[entry.file.RelativePath]

		for _, localExport := range entry.localExports {
			targetID := localSymbolsByName[localExport.localName]
			if targetID == "" {
				continue
			}
			exportMap[localExport.exportName] = targetID
		}
	}
}

func resolveReExports(parsedEntries []parsedTSFile, exportedSymbolsByFilePath map[string]map[string]string) {
	changed := true
	for changed {
		changed = false

		for _, entry := range parsedEntries {
			exportMap := exportedSymbolsByFilePath[entry.file.RelativePath]
			for _, reExport := range entry.reExports {
				sourceExports := exportedSymbolsByFilePath[reExport.sourceFilePath]
				if reExport.importedName == "*" {
					for exportName, targetID := range sourceExports {
						if exportName == "default" || targetID == "" {
							continue
						}
						if _, exists := exportMap[exportName]; exists {
							continue
						}
						exportMap[exportName] = targetID
						changed = true
					}
					continue
				}

				targetID := sourceExports[reExport.importedName]
				if targetID == "" {
					continue
				}
				if existingID, exists := exportMap[reExport.exportName]; exists && existingID == targetID {
					continue
				}
				exportMap[reExport.exportName] = targetID
				changed = true
			}
		}
	}
}

func resolveRelativeImportFile(
	file models.ScannedSourceFile,
	moduleSpecifier string,
	fileByAbsolutePath map[string]models.ScannedSourceFile,
) (models.ScannedSourceFile, bool) {
	if !strings.HasPrefix(moduleSpecifier, ".") {
		return models.ScannedSourceFile{}, false
	}

	resolvedSpecifierPath, err := filepath.Abs(filepath.Join(filepath.Dir(file.AbsolutePath), moduleSpecifier))
	if err != nil {
		return models.ScannedSourceFile{}, false
	}

	candidatePaths := map[string]struct{}{
		filepath.Clean(resolvedSpecifierPath): {},
	}
	specifierExtension := filepath.Ext(resolvedSpecifierPath)
	if specifierExtension != "" {
		specifierStem := strings.TrimSuffix(resolvedSpecifierPath, specifierExtension)
		for _, extension := range supportedJSImportExtensions {
			candidatePaths[filepath.Clean(specifierStem+extension)] = struct{}{}
		}
	} else {
		for _, extension := range supportedJSImportExtensions {
			candidatePaths[filepath.Clean(resolvedSpecifierPath+extension)] = struct{}{}
			candidatePaths[filepath.Clean(filepath.Join(resolvedSpecifierPath, "index"+extension))] = struct{}{}
		}
	}

	orderedPaths := make([]string, 0, len(candidatePaths))
	for candidatePath := range candidatePaths {
		orderedPaths = append(orderedPaths, candidatePath)
	}
	sort.Strings(orderedPaths)

	for _, candidatePath := range orderedPaths {
		resolvedFile, exists := fileByAbsolutePath[candidatePath]
		if exists {
			return resolvedFile, true
		}
	}

	return models.ScannedSourceFile{}, false
}

func createTSParseDiagnostic(file models.ScannedSourceFile, detail string) models.StructuredDiagnostic {
	return models.StructuredDiagnostic{
		Code:     tsParseErrorCode,
		Severity: models.SeverityError,
		Stage:    models.StageParse,
		Message:  "Failed to parse TypeScript source file",
		FilePath: file.RelativePath,
		Language: file.Language,
		Detail:   detail,
	}
}

func declarationExportNames(node *tree_sitter.Node, source []byte) []string {
	switch node.Kind() {
	case "lexical_declaration", "variable_declaration":
		names := []string{}
		for _, declarator := range collectNodesByKind(node, "variable_declarator") {
			declarator := declarator
			if name := resolveTSVariableName(&declarator, source); name != "" {
				names = append(names, name)
			}
		}
		return names
	default:
		name := resolveTSSymbolName(node, source, getTSSymbolKind(node))
		if name == "" || name == "anonymous" {
			return nil
		}
		return []string{name}
	}
}

func isDefaultExport(node *tree_sitter.Node, source []byte) bool {
	text := strings.TrimSpace(textOf(node, source))
	return strings.HasPrefix(text, "export default")
}

func matchCommonJSExportTarget(
	node *tree_sitter.Node,
	source []byte,
) (objectName string, propertyName string, isModuleExports bool, ok bool) {
	if node == nil || node.Kind() != "member_expression" {
		return "", "", false, false
	}

	objectNode := node.ChildByFieldName("object")
	propertyNode := node.ChildByFieldName("property")
	if objectNode == nil || propertyNode == nil {
		return "", "", false, false
	}

	switch objectNode.Kind() {
	case "identifier":
		objectName = textOf(objectNode, source)
		propertyName = textOf(propertyNode, source)
		if objectName == "exports" && propertyName != "" {
			return objectName, propertyName, false, true
		}
	case "member_expression":
		nestedObject := objectNode.ChildByFieldName("object")
		nestedProperty := objectNode.ChildByFieldName("property")
		if nestedObject == nil || nestedProperty == nil {
			return "", "", false, false
		}

		if nestedObject.Kind() == "identifier" && nestedProperty.Kind() == "property_identifier" &&
			textOf(nestedObject, source) == "module" && textOf(nestedProperty, source) == "exports" {
			propertyName = textOf(propertyNode, source)
			return "module.exports", propertyName, true, true
		}
	}

	if objectNode.Kind() == "identifier" && propertyNode.Kind() == "property_identifier" &&
		textOf(objectNode, source) == "module" && textOf(propertyNode, source) == "exports" {
		return "module", "", true, true
	}

	return "", "", false, false
}

func findChildByKind(node *tree_sitter.Node, kind string) *tree_sitter.Node {
	if node == nil {
		return nil
	}

	for _, child := range namedChildren(node) {
		child := child
		if child.Kind() == kind {
			return &child
		}
	}

	return nil
}

func relationKey(relation models.RelationEdge) string {
	return strings.Join([]string{
		relation.FromID,
		string(relation.Type),
		relation.ToID,
		string(relation.Confidence),
	}, ":")
}

func pushUniqueRelation(
	relations *[]models.RelationEdge,
	relationKeys map[string]struct{},
	relation models.RelationEdge,
) {
	key := relationKey(relation)
	if _, exists := relationKeys[key]; exists {
		return
	}

	relationKeys[key] = struct{}{}
	*relations = append(*relations, relation)
}

func normalizeAbsolutePath(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Clean(absolutePath), nil
}

func cloneStringSet(source map[string]struct{}) map[string]struct{} {
	if len(source) == 0 {
		return map[string]struct{}{}
	}

	cloned := make(map[string]struct{}, len(source))
	for key := range source {
		cloned[key] = struct{}{}
	}

	return cloned
}
