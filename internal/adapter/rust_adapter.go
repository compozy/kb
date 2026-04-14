package adapter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/compozy/kb/internal/models"
)

const (
	rustParseErrorCode = "RUST_PARSE_ERROR"

	rustSymbolKindModule    = "module"
	rustSymbolKindFunction  = "function"
	rustSymbolKindMethod    = "method"
	rustSymbolKindStruct    = "struct"
	rustSymbolKindEnum      = "enum"
	rustSymbolKindTrait     = "trait"
	rustSymbolKindTypeAlias = "typeAlias"
	rustSymbolKindConst     = "const"
	rustSymbolKindStatic    = "static"
	rustSymbolKindUnion     = "union"
	rustSymbolKindMacro     = "macro"
)

var _ models.LanguageAdapter = (*RustAdapter)(nil)

// RustAdapter parses Rust source files into graph nodes, relations, and diagnostics.
type RustAdapter struct{}

type rustUseBinding struct {
	exportedOnly   bool
	importedName   string
	kind           string
	localName      string
	sourceFilePath string
}

type rustNamespaceBinding struct {
	exportedOnly   bool
	sourceFilePath string
}

type rustReExport struct {
	exportName     string
	importedName   string
	kind           string
	sourceFilePath string
}

type rustCallTarget struct {
	kind      string
	localName string
	path      []string
}

type rustSymbolMatch struct {
	callTargets []rustCallTarget
	exportNames map[string]struct{}
	symbol      models.SymbolNode
}

type parsedRustFile struct {
	diagnostics    []models.StructuredDiagnostic
	externalNodes  map[string]models.ExternalNode
	file           models.ScannedSourceFile
	fileNode       models.GraphFile
	importBindings []rustUseBinding
	reExports      []rustReExport
	relations      []models.RelationEdge
	symbolMatches  []rustSymbolMatch
}

type rustModuleInfo struct {
	crateDirAbs  string
	crateID      string
	crateName    string
	crateRootRel string
	file         models.ScannedSourceFile
	modulePath   []string
}

type rustModuleIndex struct {
	crateInfoByName map[string]rustModuleInfo
	fileInfoByPath  map[string]rustModuleInfo
	fileByModule    map[string]models.ScannedSourceFile
	fileByPath      map[string]models.ScannedSourceFile
}

type rustUseSpec struct {
	alias      string
	path       []string
	selfImport bool
	wildcard   bool
}

// Supports reports whether the adapter handles the provided language.
func (RustAdapter) Supports(language models.SupportedLanguage) bool {
	return language == models.LangRust
}

// ParseFiles parses Rust source files into graph nodes, relations, and diagnostics.
func (adapter RustAdapter) ParseFiles(files []models.ScannedSourceFile, rootPath string) ([]models.ParsedFile, error) {
	return adapter.ParseFilesWithProgress(files, rootPath, nil)
}

// ParseFilesWithProgress parses Rust files and reports one progress tick per file.
func (adapter RustAdapter) ParseFilesWithProgress(
	files []models.ScannedSourceFile,
	rootPath string,
	report func(models.ScannedSourceFile),
) ([]models.ParsedFile, error) {
	if len(files) == 0 {
		return []models.ParsedFile{}, nil
	}

	parser, err := newParser(rustLanguage())
	if err != nil {
		return nil, fmt.Errorf("create Rust parser: %w", err)
	}
	defer parser.Close()

	moduleIndex, err := buildRustModuleIndex(files, rootPath)
	if err != nil {
		return nil, err
	}

	orderedFiles := append([]models.ScannedSourceFile(nil), files...)
	sort.Slice(orderedFiles, func(i, j int) bool {
		return orderedFiles[i].RelativePath < orderedFiles[j].RelativePath
	})

	parsedEntries := make([]parsedRustFile, 0, len(orderedFiles))
	for _, file := range orderedFiles {
		if !adapter.Supports(file.Language) {
			return nil, fmt.Errorf("parse %s: unsupported language %q", file.RelativePath, file.Language)
		}

		entry, err := parseRustFile(parser, file, moduleIndex)
		if err != nil {
			return nil, err
		}
		parsedEntries = append(parsedEntries, entry)
		if report != nil {
			report(file)
		}
	}

	localSymbolsByFilePath := make(map[string]map[string][]string, len(parsedEntries))
	exportedSymbolsByFilePath := make(map[string]map[string][]string, len(parsedEntries))
	symbolByID := make(map[string]models.SymbolNode)

	for _, entry := range parsedEntries {
		localSymbols := make(map[string][]string)
		exportedSymbols := make(map[string][]string)

		for _, symbolMatch := range entry.symbolMatches {
			symbolByID[symbolMatch.symbol.ID] = symbolMatch.symbol
			localSymbols[symbolMatch.symbol.Name] = append(localSymbols[symbolMatch.symbol.Name], symbolMatch.symbol.ID)
			for exportName := range symbolMatch.exportNames {
				exportedSymbols[exportName] = append(exportedSymbols[exportName], symbolMatch.symbol.ID)
			}
		}

		localSymbolsByFilePath[entry.file.RelativePath] = localSymbols
		exportedSymbolsByFilePath[entry.file.RelativePath] = exportedSymbols
	}

	resolveRustReExports(parsedEntries, exportedSymbolsByFilePath)

	parsedFiles := make([]models.ParsedFile, 0, len(parsedEntries))
	for _, entry := range parsedEntries {
		relationKeys := make(map[string]struct{}, len(entry.relations))
		for _, relation := range entry.relations {
			relationKeys[relationKey(relation)] = struct{}{}
		}

		directBindings := make(map[string][]string)
		namespaceBindings := make(map[string]rustNamespaceBinding)

		for _, binding := range entry.importBindings {
			switch binding.kind {
			case "symbol":
				targetIDs := uniqueStringSlice(rustBindingSymbolIDs(
					binding,
					binding.importedName,
					localSymbolsByFilePath,
					exportedSymbolsByFilePath,
				))
				if len(targetIDs) == 1 {
					targetID := targetIDs[0]
					directBindings[binding.localName] = append(directBindings[binding.localName], targetID)
					pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
						FromID:     entry.fileNode.ID,
						ToID:       targetID,
						Type:       models.RelReferences,
						Confidence: models.ConfidenceSyntactic,
					})

					if symbol, exists := symbolByID[targetID]; exists && isRustNamespaceSymbol(symbol.SymbolKind) {
						namespaceBindings[binding.localName] = rustNamespaceBinding{
							exportedOnly:   false,
							sourceFilePath: binding.sourceFilePath,
						}
					}
				}
			case "glob":
				for exportName, targetIDs := range rustBindingSymbolIndex(
					binding,
					localSymbolsByFilePath,
					exportedSymbolsByFilePath,
				) {
					if resolved := uniqueStringSlice(targetIDs); len(resolved) == 1 {
						targetID := resolved[0]
						directBindings[exportName] = append(directBindings[exportName], targetID)
						pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
							FromID:     entry.fileNode.ID,
							ToID:       targetID,
							Type:       models.RelReferences,
							Confidence: models.ConfidenceSyntactic,
						})
					}
				}
			case "module":
				namespaceBindings[binding.localName] = rustNamespaceBinding{
					exportedOnly:   binding.exportedOnly,
					sourceFilePath: binding.sourceFilePath,
				}
			}
		}

		for _, symbolMatch := range entry.symbolMatches {
			if isRustNamespaceSymbol(symbolMatch.symbol.SymbolKind) {
				if _, exists := namespaceBindings[symbolMatch.symbol.Name]; exists {
					continue
				}
				namespaceBindings[symbolMatch.symbol.Name] = rustNamespaceBinding{
					exportedOnly:   false,
					sourceFilePath: entry.file.RelativePath,
				}
			}
		}

		for _, reExport := range entry.reExports {
			switch reExport.kind {
			case "symbol":
				targetIDs := uniqueStringSlice(exportedSymbolsByFilePath[reExport.sourceFilePath][reExport.importedName])
				if len(targetIDs) == 1 {
					pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
						FromID:     entry.fileNode.ID,
						ToID:       targetIDs[0],
						Type:       models.RelExports,
						Confidence: models.ConfidenceSyntactic,
					})
				}
			case "glob":
				for _, targetIDs := range exportedSymbolsByFilePath[reExport.sourceFilePath] {
					if resolved := uniqueStringSlice(targetIDs); len(resolved) == 1 {
						pushUniqueRelation(&entry.relations, relationKeys, models.RelationEdge{
							FromID:     entry.fileNode.ID,
							ToID:       resolved[0],
							Type:       models.RelExports,
							Confidence: models.ConfidenceSyntactic,
						})
					}
				}
			}
		}

		localSymbols := localSymbolsByFilePath[entry.file.RelativePath]
		for _, symbolMatch := range entry.symbolMatches {
			if symbolMatch.symbol.SymbolKind != rustSymbolKindFunction && symbolMatch.symbol.SymbolKind != rustSymbolKindMethod {
				continue
			}

			for _, callTarget := range symbolMatch.callTargets {
				targetID := resolveRustCallTarget(callTarget, localSymbols, directBindings, namespaceBindings, localSymbolsByFilePath, exportedSymbolsByFilePath)
				if targetID == "" || targetID == symbolMatch.symbol.ID {
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

	return parsedFiles, nil
}

func parseRustFile(parser *tree_sitter.Parser, file models.ScannedSourceFile, moduleIndex rustModuleIndex) (parsedRustFile, error) {
	source, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return parsedRustFile{}, fmt.Errorf("read Rust source %s: %w", file.RelativePath, err)
	}

	tree := parser.Parse(source, nil)
	if tree == nil {
		return parsedRustFile{}, fmt.Errorf("parse Rust source %s: nil syntax tree", file.RelativePath)
	}
	defer tree.Close()

	root := tree.RootNode()
	if root == nil {
		return parsedRustFile{}, fmt.Errorf("parse Rust source %s: missing root node", file.RelativePath)
	}

	fileNode := models.GraphFile{
		ID:        createFileID(file.RelativePath),
		NodeType:  "file",
		FilePath:  file.RelativePath,
		Language:  file.Language,
		ModuleDoc: extractRustLeadingModuleDoc(string(source)),
		SymbolIDs: []string{},
	}

	entry := parsedRustFile{
		diagnostics:    []models.StructuredDiagnostic{},
		externalNodes:  make(map[string]models.ExternalNode),
		file:           file,
		fileNode:       fileNode,
		importBindings: []rustUseBinding{},
		reExports:      []rustReExport{},
		relations:      []models.RelationEdge{},
		symbolMatches:  []rustSymbolMatch{},
	}

	if root.HasError() {
		entry.diagnostics = append(entry.diagnostics, createRustParseDiagnostic(file, root.ToSexp()))
		return entry, nil
	}

	info, exists := moduleIndex.fileInfoByPath[file.RelativePath]
	if !exists {
		return parsedRustFile{}, fmt.Errorf("parse %s: missing Rust module info", file.RelativePath)
	}

	collectRustDeclarations(&entry, root, source, info, moduleIndex)

	for _, symbolMatch := range entry.symbolMatches {
		entry.fileNode.SymbolIDs = append(entry.fileNode.SymbolIDs, symbolMatch.symbol.ID)
		entry.relations = append(entry.relations, models.RelationEdge{
			FromID:     entry.fileNode.ID,
			ToID:       symbolMatch.symbol.ID,
			Type:       models.RelContains,
			Confidence: models.ConfidenceSyntactic,
		})

		if symbolMatch.symbol.Exported {
			entry.relations = append(entry.relations, models.RelationEdge{
				FromID:     entry.fileNode.ID,
				ToID:       symbolMatch.symbol.ID,
				Type:       models.RelExports,
				Confidence: models.ConfidenceSyntactic,
			})
		}
	}

	return entry, nil
}

func collectRustDeclarations(
	entry *parsedRustFile,
	node *tree_sitter.Node,
	source []byte,
	info rustModuleInfo,
	moduleIndex rustModuleIndex,
) {
	for _, child := range namedChildren(node) {
		child := child

		switch child.Kind() {
		case "function_item":
			symbol := createRustSymbol(entry.file, &child, source, rustSymbolKindFunction, hasRustVisibility(&child), "")
			entry.symbolMatches = append(entry.symbolMatches, rustSymbolMatch{
				callTargets: collectRustCallTargets(&child, source),
				exportNames: exportNamesForRustSymbol(symbol),
				symbol:      symbol,
			})
		case "const_item":
			appendRustSymbol(entry, createRustSymbol(entry.file, &child, source, rustSymbolKindConst, hasRustVisibility(&child), ""))
		case "static_item":
			appendRustSymbol(entry, createRustSymbol(entry.file, &child, source, rustSymbolKindStatic, hasRustVisibility(&child), ""))
		case "struct_item":
			appendRustSymbol(entry, createRustSymbol(entry.file, &child, source, rustSymbolKindStruct, hasRustVisibility(&child), ""))
		case "enum_item":
			appendRustSymbol(entry, createRustSymbol(entry.file, &child, source, rustSymbolKindEnum, hasRustVisibility(&child), ""))
		case "trait_item":
			traitSymbol := createRustSymbol(entry.file, &child, source, rustSymbolKindTrait, hasRustVisibility(&child), "")
			appendRustSymbol(entry, traitSymbol)
			collectRustTraitMethods(entry, &child, source, traitSymbol.Exported)
		case "type_item":
			appendRustSymbol(entry, createRustSymbol(entry.file, &child, source, rustSymbolKindTypeAlias, hasRustVisibility(&child), ""))
		case "union_item":
			appendRustSymbol(entry, createRustSymbol(entry.file, &child, source, rustSymbolKindUnion, hasRustVisibility(&child), ""))
		case "macro_definition":
			appendRustSymbol(entry, createRustSymbol(entry.file, &child, source, rustSymbolKindMacro, hasRustMacroExport(&child, source), ""))
		case "mod_item":
			moduleSymbol := createRustSymbol(entry.file, &child, source, rustSymbolKindModule, hasRustVisibility(&child), "")
			appendRustSymbol(entry, moduleSymbol)
			moduleName := rustNodeName(&child, source)
			if body := child.ChildByFieldName("body"); body != nil {
				childInfo := info
				if moduleName != "" {
					childInfo.modulePath = append(append([]string(nil), info.modulePath...), moduleName)
				}
				collectRustDeclarations(entry, body, source, childInfo, moduleIndex)
				continue
			}

			if moduleName == "" {
				continue
			}
			targetFile, ok := moduleIndex.resolveChildModule(info, moduleName)
			if !ok {
				continue
			}

			entry.importBindings = append(entry.importBindings, rustUseBinding{
				exportedOnly:   false,
				kind:           "module",
				localName:      moduleName,
				sourceFilePath: targetFile.RelativePath,
			})
			entry.relations = append(entry.relations, models.RelationEdge{
				FromID:     entry.fileNode.ID,
				ToID:       createFileID(targetFile.RelativePath),
				Type:       models.RelImports,
				Confidence: models.ConfidenceSyntactic,
			})
		case "impl_item":
			collectRustImplMethods(entry, &child, source)
		case "use_declaration":
			collectRustUseDeclaration(entry, &child, source, info, moduleIndex)
		}
	}
}

func collectRustTraitMethods(entry *parsedRustFile, node *tree_sitter.Node, source []byte, exported bool) {
	body := node.ChildByFieldName("body")
	if body == nil {
		return
	}

	for _, child := range namedChildren(body) {
		child := child
		if child.Kind() != "function_signature_item" && child.Kind() != "function_item" {
			continue
		}

		symbol := createRustSymbol(entry.file, &child, source, rustSymbolKindMethod, exported, "")
		entry.symbolMatches = append(entry.symbolMatches, rustSymbolMatch{
			exportNames: exportNamesForRustSymbol(symbol),
			symbol:      symbol,
		})
	}
}

func collectRustImplMethods(entry *parsedRustFile, node *tree_sitter.Node, source []byte) {
	body := node.ChildByFieldName("body")
	if body == nil {
		return
	}

	for _, child := range namedChildren(body) {
		child := child
		if child.Kind() != "function_item" && child.Kind() != "function_signature_item" {
			continue
		}

		symbol := createRustSymbol(entry.file, &child, source, rustSymbolKindMethod, hasRustVisibility(&child), "")
		entry.symbolMatches = append(entry.symbolMatches, rustSymbolMatch{
			callTargets: collectRustCallTargets(&child, source),
			exportNames: exportNamesForRustSymbol(symbol),
			symbol:      symbol,
		})
	}
}

func collectRustUseDeclaration(
	entry *parsedRustFile,
	node *tree_sitter.Node,
	source []byte,
	info rustModuleInfo,
	moduleIndex rustModuleIndex,
) {
	argument := node.ChildByFieldName("argument")
	if argument == nil {
		return
	}

	specs := expandRustUseSpecs(textOf(argument, source))
	if len(specs) == 0 {
		return
	}

	isExportedUse := hasRustVisibility(node)
	for _, spec := range specs {
		resolved, ok := moduleIndex.resolveUseSpec(info, spec)
		if !ok {
			externalSource := rustUseSpecLabel(spec)
			externalID := createExternalID(externalSource)
			if _, exists := entry.externalNodes[externalID]; !exists {
				entry.externalNodes[externalID] = models.ExternalNode{
					ID:       externalID,
					NodeType: "external",
					Source:   externalSource,
					Label:    externalSource,
				}
			}

			entry.relations = append(entry.relations, models.RelationEdge{
				FromID:     entry.fileNode.ID,
				ToID:       externalID,
				Type:       models.RelImports,
				Confidence: models.ConfidenceSyntactic,
			})
			continue
		}

		entry.relations = append(entry.relations, models.RelationEdge{
			FromID:     entry.fileNode.ID,
			ToID:       createFileID(resolved.sourceFilePath),
			Type:       models.RelImports,
			Confidence: models.ConfidenceSyntactic,
		})

		switch resolved.kind {
		case "module":
			entry.importBindings = append(entry.importBindings, rustUseBinding{
				exportedOnly:   resolved.exportedOnly,
				kind:           "module",
				localName:      resolved.localName,
				sourceFilePath: resolved.sourceFilePath,
			})
		case "symbol":
			entry.importBindings = append(entry.importBindings, rustUseBinding{
				exportedOnly:   resolved.exportedOnly,
				kind:           "symbol",
				localName:      resolved.localName,
				importedName:   resolved.importedName,
				sourceFilePath: resolved.sourceFilePath,
			})
			if isExportedUse {
				entry.reExports = append(entry.reExports, rustReExport{
					exportName:     resolved.localName,
					importedName:   resolved.importedName,
					kind:           "symbol",
					sourceFilePath: resolved.sourceFilePath,
				})
			}
		case "glob":
			entry.importBindings = append(entry.importBindings, rustUseBinding{
				exportedOnly:   resolved.exportedOnly,
				kind:           "glob",
				sourceFilePath: resolved.sourceFilePath,
			})
			if isExportedUse {
				entry.reExports = append(entry.reExports, rustReExport{
					exportName:     "*",
					importedName:   "*",
					kind:           "glob",
					sourceFilePath: resolved.sourceFilePath,
				})
			}
		}
	}
}

type resolvedRustUseSpec struct {
	exportedOnly   bool
	importedName   string
	kind           string
	localName      string
	sourceFilePath string
}

func resolveRustCallTarget(
	callTarget rustCallTarget,
	localSymbols map[string][]string,
	directBindings map[string][]string,
	namespaceBindings map[string]rustNamespaceBinding,
	localSymbolsByFilePath map[string]map[string][]string,
	exportedSymbolsByFilePath map[string]map[string][]string,
) string {
	switch callTarget.kind {
	case "identifier":
		if ids := uniqueStringSlice(localSymbols[callTarget.localName]); len(ids) == 1 {
			return ids[0]
		}
		if ids := uniqueStringSlice(directBindings[callTarget.localName]); len(ids) == 1 {
			return ids[0]
		}
	case "namespace":
		if len(callTarget.path) < 2 {
			return ""
		}

		binding, exists := namespaceBindings[callTarget.path[0]]
		if !exists {
			return ""
		}

		targetName := callTarget.path[len(callTarget.path)-1]
		if ids := uniqueStringSlice(rustNamespaceTargetIDs(
			binding,
			targetName,
			localSymbolsByFilePath,
			exportedSymbolsByFilePath,
		)); len(ids) == 1 {
			return ids[0]
		}
	}

	return ""
}

func rustBindingSymbolIDs(
	binding rustUseBinding,
	name string,
	localSymbolsByFilePath map[string]map[string][]string,
	exportedSymbolsByFilePath map[string]map[string][]string,
) []string {
	return rustBindingSymbolIndex(binding, localSymbolsByFilePath, exportedSymbolsByFilePath)[name]
}

func rustBindingSymbolIndex(
	binding rustUseBinding,
	localSymbolsByFilePath map[string]map[string][]string,
	exportedSymbolsByFilePath map[string]map[string][]string,
) map[string][]string {
	if binding.exportedOnly {
		return exportedSymbolsByFilePath[binding.sourceFilePath]
	}

	return localSymbolsByFilePath[binding.sourceFilePath]
}

func rustNamespaceTargetIDs(
	binding rustNamespaceBinding,
	name string,
	localSymbolsByFilePath map[string]map[string][]string,
	exportedSymbolsByFilePath map[string]map[string][]string,
) []string {
	if binding.exportedOnly {
		return exportedSymbolsByFilePath[binding.sourceFilePath][name]
	}

	return localSymbolsByFilePath[binding.sourceFilePath][name]
}

func resolveRustReExports(entries []parsedRustFile, exportedSymbolsByFilePath map[string]map[string][]string) {
	changed := true
	for changed {
		changed = false
		for _, entry := range entries {
			exportMap := exportedSymbolsByFilePath[entry.file.RelativePath]
			if exportMap == nil {
				exportMap = make(map[string][]string)
				exportedSymbolsByFilePath[entry.file.RelativePath] = exportMap
			}

			for _, reExport := range entry.reExports {
				sourceExports := exportedSymbolsByFilePath[reExport.sourceFilePath]
				if sourceExports == nil {
					continue
				}

				if reExport.kind == "glob" {
					for exportName, ids := range sourceExports {
						merged := uniqueStringSlice(append(exportMap[exportName], ids...))
						if !equalStringSlices(exportMap[exportName], merged) {
							exportMap[exportName] = merged
							changed = true
						}
					}
					continue
				}

				merged := uniqueStringSlice(append(exportMap[reExport.exportName], sourceExports[reExport.importedName]...))
				if !equalStringSlices(exportMap[reExport.exportName], merged) {
					exportMap[reExport.exportName] = merged
					changed = true
				}
			}
		}
	}
}

func appendRustSymbol(entry *parsedRustFile, symbol models.SymbolNode) {
	entry.symbolMatches = append(entry.symbolMatches, rustSymbolMatch{
		exportNames: exportNamesForRustSymbol(symbol),
		symbol:      symbol,
	})
}

func exportNamesForRustSymbol(symbol models.SymbolNode) map[string]struct{} {
	if !symbol.Exported {
		return map[string]struct{}{}
	}

	return map[string]struct{}{symbol.Name: {}}
}

func createRustSymbol(
	file models.ScannedSourceFile,
	node *tree_sitter.Node,
	source []byte,
	symbolKind string,
	exported bool,
	forcedDoc string,
) models.SymbolNode {
	docComment := forcedDoc
	if docComment == "" {
		docComment = extractRustDocComment(node, source)
	}

	symbol := models.SymbolNode{
		NodeType:   "symbol",
		Name:       rustNodeName(node, source),
		SymbolKind: symbolKind,
		Language:   file.Language,
		FilePath:   file.RelativePath,
		StartLine:  int(node.StartPosition().Row) + 1,
		EndLine:    int(node.EndPosition().Row) + 1,
		Signature:  rustSignature(node, source, symbolKind),
		DocComment: docComment,
		Exported:   exported,
	}
	if complexity := computeRustCyclomaticComplexity(node, source, symbolKind); complexity > 0 {
		symbol.CyclomaticComplexity = complexity
	}

	symbol.ID = createSymbolID(symbol)
	return symbol
}

func rustNodeName(node *tree_sitter.Node, source []byte) string {
	if node == nil {
		return "anonymous"
	}

	nameNode := node.ChildByFieldName("name")
	name := textOf(nameNode, source)
	if name == "" {
		return "anonymous"
	}

	return name
}

func rustSignature(node *tree_sitter.Node, source []byte, symbolKind string) string {
	switch symbolKind {
	case rustSymbolKindModule:
		return "mod " + rustNodeName(node, source)
	default:
		firstLine := strings.TrimSpace(strings.Split(textOf(node, source), "\n")[0])
		return firstLine
	}
}

func extractRustDocComment(node *tree_sitter.Node, source []byte) string {
	commentParts := []string{}
	expectedRow := int(node.StartPosition().Row)

	for sibling := node.PrevSibling(); sibling != nil; sibling = sibling.PrevSibling() {
		if isRustCommentNode(sibling) {
			if expectedRow-int(sibling.EndPosition().Row) > 1 {
				break
			}

			normalized := normalizeRustDocComment(textOf(sibling, source))
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

func isRustCommentNode(node *tree_sitter.Node) bool {
	if node == nil {
		return false
	}

	switch node.Kind() {
	case "line_comment", "block_comment":
		return true
	default:
		return false
	}
}

func normalizeRustDocComment(rawComment string) string {
	trimmed := strings.TrimSpace(rawComment)
	switch {
	case strings.HasPrefix(trimmed, "///"):
		lines := strings.Split(trimmed, "\n")
		normalized := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "///") || strings.HasPrefix(line, "////") {
				return ""
			}
			line = strings.TrimSpace(strings.TrimPrefix(line, "///"))
			normalized = append(normalized, line)
		}
		return strings.TrimSpace(strings.Join(normalized, "\n"))
	case strings.HasPrefix(trimmed, "/**") && !strings.HasPrefix(trimmed, "/***"):
		return normalizeRustBlockDoc(trimmed, "/**")
	default:
		return ""
	}
}

func extractRustLeadingModuleDoc(sourceText string) string {
	remaining := strings.TrimLeft(sourceText, "\ufeff \t\r\n")
	parts := []string{}

	for {
		switch {
		case strings.HasPrefix(remaining, "//!"):
			lineEnd := strings.IndexByte(remaining, '\n')
			line := remaining
			if lineEnd >= 0 {
				line = remaining[:lineEnd]
				remaining = remaining[lineEnd+1:]
			} else {
				remaining = ""
			}
			line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "//!"))
			parts = append(parts, line)
		case strings.HasPrefix(remaining, "/*!"):
			end := strings.Index(remaining, "*/")
			if end < 0 {
				return strings.TrimSpace(strings.Join(parts, "\n"))
			}
			parts = append(parts, normalizeRustBlockDoc(remaining[:end+2], "/*!"))
			remaining = remaining[end+2:]
		case strings.TrimSpace(remaining) == "":
			return strings.TrimSpace(strings.Join(parts, "\n"))
		case remaining[0] == '\n' || remaining[0] == '\r' || remaining[0] == ' ' || remaining[0] == '\t':
			remaining = strings.TrimLeft(remaining, " \t\r\n")
		default:
			return strings.TrimSpace(strings.Join(parts, "\n"))
		}
	}
}

func normalizeRustBlockDoc(rawComment, prefix string) string {
	trimmed := strings.TrimSpace(rawComment)
	if !strings.HasPrefix(trimmed, prefix) {
		return ""
	}

	trimmed = strings.TrimPrefix(trimmed, prefix)
	trimmed = strings.TrimSuffix(trimmed, "*/")
	lines := strings.Split(trimmed, "\n")
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimLeft(line, " \t")
		line = strings.TrimPrefix(line, "*")
		normalized = append(normalized, strings.TrimSpace(line))
	}

	return strings.TrimSpace(strings.Join(normalized, "\n"))
}

func hasRustVisibility(node *tree_sitter.Node) bool {
	if node == nil {
		return false
	}

	for _, child := range namedChildren(node) {
		if child.Kind() == "visibility_modifier" {
			return true
		}
	}

	return false
}

func hasRustMacroExport(node *tree_sitter.Node, source []byte) bool {
	if hasRustVisibility(node) {
		return true
	}

	for sibling := node.PrevSibling(); sibling != nil; sibling = sibling.PrevSibling() {
		if sibling.IsNamed() {
			break
		}
		if sibling.Kind() != "attribute_item" {
			continue
		}
		if strings.Contains(textOf(sibling, source), "macro_export") {
			return true
		}
	}

	return false
}

func collectRustCallTargets(node *tree_sitter.Node, source []byte) []rustCallTarget {
	body := node.ChildByFieldName("body")
	if body == nil {
		return nil
	}

	targets := []rustCallTarget{}
	walkNamed(body, func(current *tree_sitter.Node) bool {
		if current == nil {
			return false
		}

		if current != body && current.Kind() == "closure_expression" {
			return false
		}

		if current.Kind() != "call_expression" {
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
				targets = append(targets, rustCallTarget{
					kind:      "identifier",
					localName: localName,
				})
			}
		case "generic_function":
			targets = append(targets, collectRustCallTargetsFromFunctionNode(functionNode.ChildByFieldName("function"), source)...)
		default:
			targets = append(targets, collectRustCallTargetsFromFunctionNode(functionNode, source)...)
		}

		return true
	})

	return targets
}

func collectRustCallTargetsFromFunctionNode(node *tree_sitter.Node, source []byte) []rustCallTarget {
	if node == nil {
		return nil
	}

	switch node.Kind() {
	case "identifier":
		localName := textOf(node, source)
		if localName == "" {
			return nil
		}
		return []rustCallTarget{{
			kind:      "identifier",
			localName: localName,
		}}
	case "scoped_identifier", "scoped_type_identifier":
		path := splitRustPath(textOf(node, source))
		if len(path) < 2 {
			return nil
		}
		return []rustCallTarget{{
			kind: "namespace",
			path: path,
		}}
	default:
		return nil
	}
}

func computeRustCyclomaticComplexity(node *tree_sitter.Node, source []byte, symbolKind string) int {
	if symbolKind != rustSymbolKindFunction && symbolKind != rustSymbolKindMethod {
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

		if current != body && current.Kind() == "closure_expression" {
			return false
		}

		switch current.Kind() {
		case "if_expression", "for_expression", "while_expression", "loop_expression", "match_arm":
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

func createRustParseDiagnostic(file models.ScannedSourceFile, detail string) models.StructuredDiagnostic {
	return models.StructuredDiagnostic{
		Code:     rustParseErrorCode,
		Severity: models.SeverityError,
		Stage:    models.StageParse,
		Message:  "Failed to parse Rust source file",
		FilePath: file.RelativePath,
		Language: file.Language,
		Detail:   detail,
	}
}

func isRustNamespaceSymbol(kind string) bool {
	switch kind {
	case rustSymbolKindStruct, rustSymbolKindEnum, rustSymbolKindTrait, rustSymbolKindTypeAlias, rustSymbolKindUnion:
		return true
	default:
		return false
	}
}

func buildRustModuleIndex(files []models.ScannedSourceFile, rootPath string) (rustModuleIndex, error) {
	rootAbs, err := filepath.Abs(rootPath)
	if err != nil {
		return rustModuleIndex{}, fmt.Errorf("resolve Rust root path: %w", err)
	}

	fileByPath := make(map[string]models.ScannedSourceFile, len(files))
	for _, file := range files {
		fileByPath[file.RelativePath] = file
	}

	index := rustModuleIndex{
		crateInfoByName: make(map[string]rustModuleInfo),
		fileInfoByPath:  make(map[string]rustModuleInfo, len(files)),
		fileByModule:    make(map[string]models.ScannedSourceFile, len(files)),
		fileByPath:      fileByPath,
	}
	crateNames := make(map[string]string)

	for _, file := range files {
		info, err := detectRustModuleInfo(file, rootAbs, fileByPath, crateNames)
		if err != nil {
			return rustModuleIndex{}, err
		}

		index.fileInfoByPath[file.RelativePath] = info
		index.fileByModule[rustModuleKey(info.crateID, info.modulePath)] = file
		if file.RelativePath == info.crateRootRel && info.crateName != "" {
			index.crateInfoByName[info.crateName] = info
		}
	}

	return index, nil
}

func detectRustModuleInfo(
	file models.ScannedSourceFile,
	rootAbs string,
	fileByPath map[string]models.ScannedSourceFile,
	crateNames map[string]string,
) (rustModuleInfo, error) {
	crateDirAbs, err := findRustCrateDir(filepath.Dir(file.AbsolutePath), rootAbs)
	if err != nil {
		return rustModuleInfo{}, err
	}

	relativeToCrate, err := filepath.Rel(crateDirAbs, file.AbsolutePath)
	if err != nil {
		return rustModuleInfo{}, fmt.Errorf("resolve crate-relative path for %s: %w", file.RelativePath, err)
	}
	relativeToCrate = filepath.ToSlash(relativeToCrate)

	rootFileRel := file.RelativePath
	if directBinRoot(relativeToCrate) {
		rootFileRel = file.RelativePath
	} else {
		for _, candidate := range []string{"src/lib.rs", "src/main.rs"} {
			absoluteCandidate := filepath.Join(crateDirAbs, filepath.FromSlash(candidate))
			if _, err := os.Stat(absoluteCandidate); err == nil {
				relativeCandidate, relErr := filepath.Rel(rootAbs, absoluteCandidate)
				if relErr != nil {
					return rustModuleInfo{}, fmt.Errorf("resolve crate root file for %s: %w", file.RelativePath, relErr)
				}
				rootFileRel = filepath.ToSlash(relativeCandidate)
				break
			}
		}
	}

	rootDirRel := pathDir(rootFileRel)
	modulePath := []string{}
	if file.RelativePath != rootFileRel {
		rootDirAbs := filepath.Join(rootAbs, filepath.FromSlash(rootDirRel))
		relativeToRootDir, err := filepath.Rel(rootDirAbs, file.AbsolutePath)
		if err != nil {
			return rustModuleInfo{}, fmt.Errorf("resolve module path for %s: %w", file.RelativePath, err)
		}
		modulePath = moduleSegmentsFromRelativePath(filepath.ToSlash(relativeToRootDir))
	}

	crateName := crateNames[crateDirAbs]
	if crateName == "" {
		crateName, err = readRustCrateName(crateDirAbs)
		if err != nil {
			return rustModuleInfo{}, err
		}
		crateNames[crateDirAbs] = crateName
	}

	return rustModuleInfo{
		crateDirAbs:  crateDirAbs,
		crateID:      rootFileRel,
		crateName:    crateName,
		crateRootRel: rootFileRel,
		file:         file,
		modulePath:   modulePath,
	}, nil
}

func findRustCrateDir(startDir string, rootAbs string) (string, error) {
	current := filepath.Clean(startDir)
	root := filepath.Clean(rootAbs)

	for {
		cargoPath := filepath.Join(current, "Cargo.toml")
		if _, err := os.Stat(cargoPath); err == nil {
			return current, nil
		}

		if current == root {
			return root, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return root, nil
		}
		current = parent
	}
}

func directBinRoot(relativeToCrate string) bool {
	if !strings.HasPrefix(relativeToCrate, "src/bin/") {
		return false
	}

	rest := strings.TrimPrefix(relativeToCrate, "src/bin/")
	return rest != "" && !strings.Contains(rest, "/")
}

func pathDir(relativePath string) string {
	dir := filepath.ToSlash(filepath.Dir(relativePath))
	if dir == "." {
		return ""
	}
	return dir
}

func moduleSegmentsFromRelativePath(relativePath string) []string {
	normalized := strings.TrimPrefix(relativePath, "./")
	normalized = strings.TrimPrefix(normalized, "src/")
	normalized = strings.TrimPrefix(normalized, "bin/")
	normalized = strings.TrimSuffix(normalized, ".rs")
	normalized = strings.TrimSuffix(normalized, "/mod")
	normalized = strings.Trim(normalized, "/")
	if normalized == "" {
		return nil
	}
	return strings.Split(normalized, "/")
}

func rustModuleKey(crateID string, modulePath []string) string {
	if len(modulePath) == 0 {
		return crateID + "::"
	}
	return crateID + "::" + strings.Join(modulePath, "::")
}

func (index rustModuleIndex) resolveChildModule(info rustModuleInfo, moduleName string) (models.ScannedSourceFile, bool) {
	modulePath := append(append([]string(nil), info.modulePath...), moduleName)
	file, ok := index.fileByModule[rustModuleKey(info.crateID, modulePath)]
	return file, ok
}

func (index rustModuleIndex) resolveUseSpec(info rustModuleInfo, spec rustUseSpec) (resolvedRustUseSpec, bool) {
	segments := append([]string(nil), spec.path...)
	if len(segments) == 0 {
		return resolvedRustUseSpec{}, false
	}

	targetCrateID := info.crateID
	baseModulePath := []string{}
	switch segments[0] {
	case "crate":
		segments = segments[1:]
	case "self":
		baseModulePath = append(baseModulePath, info.modulePath...)
		segments = segments[1:]
	case "super":
		baseModulePath = append(baseModulePath, info.modulePath...)
		for len(segments) > 0 && segments[0] == "super" {
			if len(baseModulePath) > 0 {
				baseModulePath = baseModulePath[:len(baseModulePath)-1]
			}
			segments = segments[1:]
		}
	default:
		if crateInfo, exists := index.crateInfoByName[segments[0]]; exists {
			targetCrateID = crateInfo.crateID
			segments = segments[1:]
		}
	}
	exportedOnly := targetCrateID != info.crateID

	segments = append(baseModulePath, segments...)
	if len(segments) == 0 {
		if file, ok := index.fileByModule[rustModuleKey(targetCrateID, nil)]; ok {
			localName := spec.alias
			if localName == "" {
				if crateInfo, exists := index.fileInfoByPath[file.RelativePath]; exists && crateInfo.crateName != "" {
					localName = crateInfo.crateName
				} else {
					localName = pathBaseWithoutExt(file.RelativePath)
				}
			}
			return resolvedRustUseSpec{
				exportedOnly:   exportedOnly,
				kind:           "module",
				localName:      localName,
				sourceFilePath: file.RelativePath,
			}, true
		}
		return resolvedRustUseSpec{}, false
	}

	if spec.selfImport {
		if file, ok := index.fileByModule[rustModuleKey(targetCrateID, segments)]; ok {
			localName := spec.alias
			if localName == "" {
				localName = segments[len(segments)-1]
			}
			return resolvedRustUseSpec{
				exportedOnly:   exportedOnly,
				kind:           "module",
				localName:      localName,
				sourceFilePath: file.RelativePath,
			}, true
		}
		return resolvedRustUseSpec{}, false
	}

	if spec.wildcard {
		if file, ok := index.fileByModule[rustModuleKey(targetCrateID, segments)]; ok {
			return resolvedRustUseSpec{
				exportedOnly:   exportedOnly,
				kind:           "glob",
				localName:      "*",
				sourceFilePath: file.RelativePath,
			}, true
		}
		return resolvedRustUseSpec{}, false
	}

	if file, ok := index.fileByModule[rustModuleKey(targetCrateID, segments)]; ok {
		localName := spec.alias
		if localName == "" {
			localName = segments[len(segments)-1]
		}
		return resolvedRustUseSpec{
			exportedOnly:   exportedOnly,
			kind:           "module",
			localName:      localName,
			sourceFilePath: file.RelativePath,
		}, true
	}

	if len(segments) >= 2 {
		modulePath := segments[:len(segments)-1]
		importedName := segments[len(segments)-1]
		if file, ok := index.fileByModule[rustModuleKey(targetCrateID, modulePath)]; ok {
			localName := spec.alias
			if localName == "" {
				localName = importedName
			}
			return resolvedRustUseSpec{
				exportedOnly:   exportedOnly,
				importedName:   importedName,
				kind:           "symbol",
				localName:      localName,
				sourceFilePath: file.RelativePath,
			}, true
		}
	}

	return resolvedRustUseSpec{}, false
}

func expandRustUseSpecs(argument string) []rustUseSpec {
	return expandRustUseSpec(nil, strings.TrimSpace(argument))
}

func expandRustUseSpec(prefix []string, value string) []rustUseSpec {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	if parts := splitTopLevelComma(value); len(parts) > 1 {
		specs := make([]rustUseSpec, 0, len(parts))
		for _, part := range parts {
			specs = append(specs, expandRustUseSpec(prefix, part)...)
		}
		return specs
	}

	if open := strings.IndexByte(value, '{'); open >= 0 {
		close := matchingBraceIndex(value, open)
		if close > open {
			head := strings.TrimSpace(strings.TrimSuffix(value[:open], "::"))
			base := append([]string(nil), prefix...)
			if head != "" {
				base = append(base, splitRustPath(head)...)
			}

			inner := value[open+1 : close]
			specs := []rustUseSpec{}
			for _, part := range splitTopLevelComma(inner) {
				specs = append(specs, expandRustUseSpec(base, part)...)
			}
			return specs
		}
	}

	pathPart, alias := splitRustUseAlias(value)
	pathSegments := append([]string(nil), prefix...)
	pathSegments = append(pathSegments, splitRustPath(pathPart)...)
	if len(pathSegments) == 0 {
		return nil
	}

	spec := rustUseSpec{
		alias: alias,
		path:  pathSegments,
	}
	last := pathSegments[len(pathSegments)-1]
	switch last {
	case "*":
		spec.wildcard = true
		spec.path = pathSegments[:len(pathSegments)-1]
	case "self":
		spec.selfImport = true
		spec.path = pathSegments[:len(pathSegments)-1]
	}

	return []rustUseSpec{spec}
}

func splitTopLevelComma(value string) []string {
	depth := 0
	start := 0
	parts := []string{}

	for index, r := range value {
		switch r {
		case '{':
			depth++
		case '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				part := strings.TrimSpace(value[start:index])
				if part != "" {
					parts = append(parts, part)
				}
				start = index + 1
			}
		}
	}

	if tail := strings.TrimSpace(value[start:]); tail != "" {
		parts = append(parts, tail)
	}

	return parts
}

func matchingBraceIndex(value string, open int) int {
	depth := 0
	for index := open; index < len(value); index++ {
		switch value[index] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return index
			}
		}
	}
	return -1
}

func splitRustUseAlias(value string) (string, string) {
	parts := strings.Split(value, " as ")
	if len(parts) != 2 {
		return strings.TrimSpace(value), ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func splitRustPath(value string) []string {
	normalized := strings.TrimSpace(value)
	normalized = strings.TrimPrefix(normalized, "::")
	normalized = strings.TrimSuffix(normalized, "::")
	if normalized == "" {
		return nil
	}

	raw := strings.Split(normalized, "::")
	segments := make([]string, 0, len(raw))
	for _, segment := range raw {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		segments = append(segments, segment)
	}
	return segments
}

func rustUseSpecLabel(spec rustUseSpec) string {
	label := strings.Join(spec.path, "::")
	switch {
	case spec.wildcard:
		if label == "" {
			return "*"
		}
		return label + "::*"
	case spec.selfImport:
		if label == "" {
			return "self"
		}
		return label + "::self"
	default:
		return label
	}
}

func uniqueStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	sort.Strings(unique)
	return unique
}

func equalStringSlices(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

type cargoManifest struct {
	Lib struct {
		Name string `toml:"name"`
	} `toml:"lib"`
	Package struct {
		Name string `toml:"name"`
	} `toml:"package"`
}

func readRustCrateName(crateDirAbs string) (string, error) {
	manifestPath := filepath.Join(crateDirAbs, "Cargo.toml")
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return filepath.Base(crateDirAbs), nil
		}
		return "", fmt.Errorf("read Cargo.toml in %s: %w", crateDirAbs, err)
	}

	manifest := cargoManifest{}
	if err := toml.Unmarshal(content, &manifest); err != nil {
		return "", fmt.Errorf("parse Cargo.toml in %s: %w", crateDirAbs, err)
	}

	name := strings.TrimSpace(manifest.Lib.Name)
	if name == "" {
		name = strings.TrimSpace(manifest.Package.Name)
	}
	if name == "" {
		name = filepath.Base(crateDirAbs)
	}

	return strings.ReplaceAll(name, "-", "_"), nil
}

func pathBaseWithoutExt(relativePath string) string {
	base := filepath.Base(relativePath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}
