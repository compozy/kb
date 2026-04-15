package adapter

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/compozy/kb/internal/models"
)

const (
	javaParseErrorCode         = "JAVA_PARSE_ERROR"
	javaResolutionFallbackCode = "JAVA_RESOLUTION_FALLBACK"
	javaModuleHintWarningCode  = "JAVA_MODULE_HINT_WARNING"

	javaDiagnosticDetailMaxBytes      = 16 * 1024
	javaFallbackDiagnosticMaxEntries  = 200
	javaModuleHintWarningMaxEntries   = 64
	javaDiagnosticTruncationPrefixKey = "meta:truncated"

	javaSymbolKindPackage   = "package"
	javaSymbolKindClass     = "class"
	javaSymbolKindInterface = "interface"
	javaSymbolKindEnum      = "enum"
	javaSymbolKindRecord    = "record"
	javaSymbolKindMethod    = "method"
)

var (
	_ models.LanguageAdapter = (*JavaAdapter)(nil)

	javaPackagePattern = regexp.MustCompile(`(?m)^\s*package\s+([A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)*)\s*;`)
	javaImportPattern  = regexp.MustCompile(`(?m)^\s*import\s+(static\s+)?([A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*|\.\*)*)\s*;`)

	javaGradleIncludeLinePattern = regexp.MustCompile(`(?m)^\s*include(?:\s*\(|\s+)(.+)$`)
	javaQuotedTokenPattern       = regexp.MustCompile(`["']([^"']+)["']`)
	javaGradleProjectDepPattern  = regexp.MustCompile(`project\(\s*["']:?([^"')]+)["']\s*\)`)
	javaMavenModulePattern       = regexp.MustCompile(`(?s)<module>\s*([^<]+)\s*</module>`)
	javaMavenDependencyPattern   = regexp.MustCompile(`(?s)<dependency>.*?<artifactId>\s*([^<]+)\s*</artifactId>.*?</dependency>`)
	javaMavenArtifactPattern     = regexp.MustCompile(`(?s)<artifactId>\s*([^<\s][^<]*)\s*</artifactId>`)
)

// JavaAdapter parses Java source files into graph nodes, relations, and diagnostics.
type JavaAdapter struct{}

type javaImportRef struct {
	importPath string
	isStatic   bool
	isWildcard bool
	simpleName string
}

type javaCallTarget struct {
	methodName string
	qualifier  string
}

type javaResolver interface {
	Resolve(
		entry parsedJavaFile,
		ctx javaResolutionContext,
		classSymbolByFQN map[string]string,
		localClasses map[string]string,
		importClassFQN map[string]string,
		wildcardClassFQNs map[string][]string,
		staticMethodImportIDs map[string][]string,
		ambiguousImportClassTargets map[string][]string,
		unresolved []javaUnresolvedRef,
	) ([]models.RelationEdge, []javaUnresolvedRef)
}

type javaDeepResolver struct{}

type javaSyntacticResolver struct{}

type javaSymbolMatch struct {
	symbol      models.SymbolNode
	ownerType   string
	callTargets []javaCallTarget
}

type javaResolutionContext struct {
	classSymbolByFQN       map[string]string
	classSymbolIDsByFQN    map[string][]string
	localClassFQNByFile    map[string]map[string]string
	topLevelClassFQNsByPkg map[string][]string
	methodIDsByClassFQN    map[string]map[string][]string
	methodIDsByPackage     map[string]map[string][]string
	ownerClassFQNByMethod  map[string]string
	moduleByClassSymbolID  map[string]string
}

type javaImportLookupIndexes struct {
	classSymbolByFQN            map[string]string
	importClassFQN              map[string]string
	wildcardClassFQNs           map[string][]string
	staticMethodImportIDs       map[string][]string
	ambiguousImportClassTargets map[string][]string
}

type javaModuleHints struct {
	moduleDependencies  map[string]map[string]struct{}
	fileModuleBySrcRoot map[string]string
	warnings            []string
}

type javaUnresolvedRef struct {
	callTarget   *javaCallTarget
	importRef    *javaImportRef
	reason       string
	relationType models.RelationType
	sourceID     string
	targetHint   string
}

type parsedJavaFile struct {
	packageName   string
	file          models.GraphFile
	symbolMatches []javaSymbolMatch
	externalNodes map[string]models.ExternalNode
	relations     []models.RelationEdge
	diagnostics   []models.StructuredDiagnostic
	imports       []javaImportRef
}

// Supports reports whether the adapter handles the provided language.
func (JavaAdapter) Supports(language models.SupportedLanguage) bool {
	return language == models.LangJava
}

// ParseFiles parses Java source files into graph nodes, relations, and diagnostics.
func (adapter JavaAdapter) ParseFiles(files []models.ScannedSourceFile, rootPath string) ([]models.ParsedFile, error) {
	return adapter.ParseFilesWithProgress(files, rootPath, nil)
}

// ParseFilesWithProgress parses Java files and reports one progress tick per file.
func (adapter JavaAdapter) ParseFilesWithProgress(
	files []models.ScannedSourceFile,
	rootPath string,
	report func(models.ScannedSourceFile),
) ([]models.ParsedFile, error) {
	if len(files) == 0 {
		return []models.ParsedFile{}, nil
	}

	parser, err := newParser(javaLanguage())
	if err != nil {
		return nil, fmt.Errorf("create Java parser: %w", err)
	}
	defer parser.Close()

	orderedFiles := append([]models.ScannedSourceFile(nil), files...)
	sort.Slice(orderedFiles, func(i, j int) bool {
		return orderedFiles[i].RelativePath < orderedFiles[j].RelativePath
	})

	parsedEntries := make([]parsedJavaFile, 0, len(orderedFiles))
	for _, file := range orderedFiles {
		if !adapter.Supports(file.Language) {
			return nil, fmt.Errorf("parse %s: unsupported language %q", file.RelativePath, file.Language)
		}

		entry, parseErr := parseJavaFile(parser, file)
		if parseErr != nil {
			return nil, parseErr
		}

		parsedEntries = append(parsedEntries, entry)
		if report != nil {
			report(file)
		}
	}

	moduleHints := discoverJavaModuleHints(rootPath)
	if len(moduleHints.warnings) > 0 {
		parsedEntries[0].diagnostics = append(
			parsedEntries[0].diagnostics,
			createJavaModuleHintDiagnostic(parsedEntries[0].file, moduleHints.warnings),
		)
	}

	resolutionContext := buildJavaResolutionContext(parsedEntries, moduleHints)
	deepResolver := javaResolver(javaDeepResolver{})
	fallbackResolver := javaResolver(javaSyntacticResolver{})

	parsedFiles := make([]models.ParsedFile, 0, len(parsedEntries))
	for _, entry := range parsedEntries {
		relationKeys := make(map[string]struct{}, len(entry.relations))
		for _, relation := range entry.relations {
			relationKeys[relationKey(relation)] = struct{}{}
		}

		localClasses := resolutionContext.localClassFQNByFile[entry.file.FilePath]
		importIndexes := buildJavaImportLookupIndexes(
			entry.imports,
			resolutionContext.classSymbolIDsByFQN,
			resolutionContext.topLevelClassFQNsByPkg,
			resolutionContext.methodIDsByClassFQN,
			resolutionContext.moduleByClassSymbolID,
			moduleHints.moduleForFile(entry.file.FilePath),
			moduleHints.moduleDependencies,
		)

		deepRelations, unresolved := deepResolver.Resolve(
			entry,
			resolutionContext,
			importIndexes.classSymbolByFQN,
			localClasses,
			importIndexes.importClassFQN,
			importIndexes.wildcardClassFQNs,
			importIndexes.staticMethodImportIDs,
			importIndexes.ambiguousImportClassTargets,
			nil,
		)
		for _, relation := range deepRelations {
			pushUniqueRelation(&entry.relations, relationKeys, relation)
		}

		fallbackRelations, _ := fallbackResolver.Resolve(
			entry,
			resolutionContext,
			importIndexes.classSymbolByFQN,
			localClasses,
			importIndexes.importClassFQN,
			importIndexes.wildcardClassFQNs,
			importIndexes.staticMethodImportIDs,
			importIndexes.ambiguousImportClassTargets,
			unresolved,
		)
		for _, relation := range fallbackRelations {
			pushUniqueRelation(&entry.relations, relationKeys, relation)
		}

		if len(unresolved) > 0 {
			entry.diagnostics = append(entry.diagnostics, createJavaResolutionFallbackDiagnostic(entry.file, unresolved))
		}

		sortRelationEdges(entry.relations)
		sortJavaDiagnostics(entry.diagnostics)

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

	sort.Slice(parsedFiles, func(i, j int) bool {
		return parsedFiles[i].File.FilePath < parsedFiles[j].File.FilePath
	})

	return parsedFiles, nil
}

func parseJavaFile(parser *tree_sitter.Parser, file models.ScannedSourceFile) (parsedJavaFile, error) {
	source, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return parsedJavaFile{}, fmt.Errorf("read Java source %s: %w", file.RelativePath, err)
	}

	sourceText := string(source)
	fileID := createFileID(file.RelativePath)
	entry := parsedJavaFile{
		file: models.GraphFile{
			ID:        fileID,
			NodeType:  "file",
			FilePath:  file.RelativePath,
			Language:  file.Language,
			ModuleDoc: extractLeadingComment(sourceText),
			SymbolIDs: []string{},
		},
		externalNodes: map[string]models.ExternalNode{},
		relations:     []models.RelationEdge{},
		diagnostics:   []models.StructuredDiagnostic{},
		imports:       []javaImportRef{},
	}

	tree := parser.Parse(source, nil)
	if tree == nil {
		entry.diagnostics = append(entry.diagnostics, createJavaParseDiagnostic(file, "nil syntax tree"))
		return entry, nil
	}
	defer tree.Close()

	root := tree.RootNode()
	if root == nil {
		entry.diagnostics = append(entry.diagnostics, createJavaParseDiagnostic(file, "missing root node"))
		return entry, nil
	}

	if root.HasError() {
		entry.diagnostics = append(entry.diagnostics, createJavaParseDiagnostic(file, root.ToSexp()))
		return entry, nil
	}

	entry.packageName, _ = extractJavaPackage(sourceText)
	if entry.packageName != "" {
		packageSymbol := models.SymbolNode{
			NodeType:   "symbol",
			Name:       entry.packageName,
			SymbolKind: javaSymbolKindPackage,
			Language:   file.Language,
			FilePath:   file.RelativePath,
			StartLine:  1,
			EndLine:    1,
			Signature:  "package " + entry.packageName,
			Exported:   true,
		}
		packageSymbol.ID = createSymbolID(packageSymbol)
		entry.symbolMatches = append(entry.symbolMatches, javaSymbolMatch{symbol: packageSymbol})
	}

	entry.imports = extractJavaImports(sourceText)
	for _, importRef := range entry.imports {
		externalID := createExternalID(importRef.importPath)
		entry.externalNodes[externalID] = models.ExternalNode{
			ID:       externalID,
			NodeType: "external",
			Source:   importRef.importPath,
			Label:    importRef.importPath,
		}
		entry.relations = append(entry.relations, models.RelationEdge{
			FromID:     fileID,
			ToID:       externalID,
			Type:       models.RelImports,
			Confidence: models.ConfidenceSyntactic,
		})
	}

	for _, declaration := range namedChildren(root) {
		declaration := declaration
		switch declaration.Kind() {
		case "class_declaration":
			entry.symbolMatches = append(entry.symbolMatches, parseJavaTypeDeclaration(file, &declaration, source, javaSymbolKindClass, nil)...)
		case "interface_declaration":
			entry.symbolMatches = append(entry.symbolMatches, parseJavaTypeDeclaration(file, &declaration, source, javaSymbolKindInterface, nil)...)
		case "enum_declaration":
			entry.symbolMatches = append(entry.symbolMatches, parseJavaTypeDeclaration(file, &declaration, source, javaSymbolKindEnum, nil)...)
		case "record_declaration":
			entry.symbolMatches = append(entry.symbolMatches, parseJavaTypeDeclaration(file, &declaration, source, javaSymbolKindRecord, nil)...)
		}
	}

	sort.Slice(entry.symbolMatches, func(i, j int) bool {
		left := entry.symbolMatches[i].symbol
		right := entry.symbolMatches[j].symbol
		if left.StartLine == right.StartLine {
			return left.Name < right.Name
		}
		return left.StartLine < right.StartLine
	})

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

func parseJavaTypeDeclaration(
	file models.ScannedSourceFile,
	typeNode *tree_sitter.Node,
	source []byte,
	symbolKind string,
	ownerTypePath []string,
) []javaSymbolMatch {
	typeSimpleName := resolveJavaSymbolName(typeNode, source, symbolKind)
	qualifiedTypeName := javaJoinQualifiedName(ownerTypePath, typeSimpleName)
	typeSymbol := createJavaSymbol(file, typeNode, source, symbolKind, "", qualifiedTypeName)
	matches := []javaSymbolMatch{{symbol: typeSymbol}}

	bodyNode := typeNode.ChildByFieldName("body")
	if bodyNode == nil {
		return matches
	}

	nestedOwnerTypePath := append(append([]string(nil), ownerTypePath...), typeSimpleName)
	for _, member := range namedChildren(bodyNode) {
		member := member
		switch member.Kind() {
		case "method_declaration", "constructor_declaration":
			methodSymbol := createJavaSymbol(file, &member, source, javaSymbolKindMethod, "", "")
			matches = append(matches, javaSymbolMatch{
				symbol:      methodSymbol,
				ownerType:   qualifiedTypeName,
				callTargets: collectJavaCallTargets(&member, source),
			})
		case "class_declaration":
			matches = append(matches, parseJavaTypeDeclaration(file, &member, source, javaSymbolKindClass, nestedOwnerTypePath)...)
		case "interface_declaration":
			matches = append(matches, parseJavaTypeDeclaration(file, &member, source, javaSymbolKindInterface, nestedOwnerTypePath)...)
		case "enum_declaration":
			matches = append(matches, parseJavaTypeDeclaration(file, &member, source, javaSymbolKindEnum, nestedOwnerTypePath)...)
		case "record_declaration":
			matches = append(matches, parseJavaTypeDeclaration(file, &member, source, javaSymbolKindRecord, nestedOwnerTypePath)...)
		}
	}

	return matches
}

func createJavaSymbol(
	file models.ScannedSourceFile,
	node *tree_sitter.Node,
	source []byte,
	symbolKind string,
	fallbackDoc string,
	nameOverride string,
) models.SymbolNode {
	name := resolveJavaSymbolName(node, source, symbolKind)
	if strings.TrimSpace(nameOverride) != "" {
		name = strings.TrimSpace(nameOverride)
	}
	signature := formatJavaSignature(node, source, symbolKind, name)
	docComment := extractAttachedComment(node, source)
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
		Exported:   isJavaExported(node, source, symbolKind),
	}

	if complexity := computeJavaCyclomaticComplexity(node, source, symbolKind); complexity > 0 {
		symbol.CyclomaticComplexity = complexity
	}

	symbol.ID = createSymbolID(symbol)
	return symbol
}

func resolveJavaSymbolName(node *tree_sitter.Node, source []byte, symbolKind string) string {
	if symbolKind == javaSymbolKindPackage {
		return textOf(node.ChildByFieldName("name"), source)
	}

	name := textOf(node.ChildByFieldName("name"), source)
	if name != "" {
		return name
	}

	for _, child := range namedChildren(node) {
		child := child
		if child.Kind() == "identifier" {
			if identifier := textOf(&child, source); identifier != "" {
				return identifier
			}
		}
	}

	return "anonymous"
}

func formatJavaSignature(node *tree_sitter.Node, source []byte, symbolKind string, name string) string {
	switch symbolKind {
	case javaSymbolKindPackage:
		return "package " + name
	case javaSymbolKindClass:
		return "class " + name
	case javaSymbolKindInterface:
		return "interface " + name
	case javaSymbolKindEnum:
		return "enum " + name
	case javaSymbolKindRecord:
		return "record " + name
	default:
		firstLine := strings.TrimSpace(strings.Split(textOf(node, source), "\n")[0])
		if firstLine == "" {
			return "method " + name
		}
		return firstLine
	}
}

func isJavaExported(node *tree_sitter.Node, source []byte, symbolKind string) bool {
	if symbolKind == javaSymbolKindPackage {
		return true
	}

	for _, child := range namedChildren(node) {
		child := child
		if child.Kind() != "modifiers" {
			continue
		}
		if strings.Contains(textOf(&child, source), "public") {
			return true
		}
	}

	return false
}

func collectJavaCallTargets(node *tree_sitter.Node, source []byte) []javaCallTarget {
	targets := []javaCallTarget{}
	for _, methodInvocation := range collectNodesByKind(node, "method_invocation") {
		methodInvocation := methodInvocation

		methodName := textOf(methodInvocation.ChildByFieldName("name"), source)
		qualifier := textOf(methodInvocation.ChildByFieldName("object"), source)
		if methodName == "" {
			methodName = resolveJavaMethodInvocationName(&methodInvocation, source)
		}

		if qualifier == "" {
			qualifier = resolveJavaMethodInvocationQualifier(&methodInvocation, source)
		}

		qualifier = normalizeJavaQualifier(qualifier)
		if methodName == "" {
			continue
		}

		targets = append(targets, javaCallTarget{
			methodName: methodName,
			qualifier:  qualifier,
		})
	}

	return targets
}

func resolveJavaMethodInvocationName(node *tree_sitter.Node, source []byte) string {
	if node == nil {
		return ""
	}

	text := strings.TrimSpace(textOf(node, source))
	openParen := strings.IndexByte(text, '(')
	if openParen < 0 {
		return ""
	}
	prefix := strings.TrimSpace(text[:openParen])
	if prefix == "" {
		return ""
	}

	parts := strings.Split(prefix, ".")
	return javaLastIdentifierSegment(parts[len(parts)-1])
}

func resolveJavaMethodInvocationQualifier(node *tree_sitter.Node, source []byte) string {
	if node == nil {
		return ""
	}

	text := strings.TrimSpace(textOf(node, source))
	openParen := strings.IndexByte(text, '(')
	if openParen < 0 {
		return ""
	}
	prefix := strings.TrimSpace(text[:openParen])
	lastDot := strings.LastIndex(prefix, ".")
	if lastDot <= 0 {
		return ""
	}

	qualifier := strings.TrimSpace(prefix[:lastDot])
	if qualifier == "" {
		return ""
	}

	return qualifier
}

func javaLastIdentifierSegment(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	for _, segment := range strings.Split(trimmed, ".") {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		trimmed = segment
	}

	return strings.Trim(trimmed, "[]")
}

func normalizeJavaQualifier(value string) string {
	trimmed := strings.TrimSpace(strings.Trim(value, "[]"))
	if trimmed == "" {
		return ""
	}

	parts := strings.Split(trimmed, ".")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		segment := javaLastIdentifierSegment(part)
		if segment == "" || segment == "this" || segment == "super" {
			continue
		}
		filtered = append(filtered, segment)
	}
	if len(filtered) == 0 {
		return ""
	}

	return strings.Join(filtered, ".")
}

func resolveJavaCallTarget(
	callTarget javaCallTarget,
	packageName string,
	localClasses map[string]string,
	importClassFQN map[string]string,
	wildcardClassFQNs map[string][]string,
	staticMethodImportIDs map[string][]string,
	ambiguousImportClassTargets map[string][]string,
	methodIDsByClassFQN map[string]map[string][]string,
	methodIDsByPackage map[string]map[string][]string,
) []string {
	if callTarget.methodName == "" {
		return nil
	}

	if callTarget.qualifier == "" {
		if ids := staticMethodImportIDs[callTarget.methodName]; len(ids) > 0 {
			return ids
		}
		return methodIDsByPackage[packageName][callTarget.methodName]
	}
	if len(ambiguousImportClassTargets[callTarget.qualifier]) > 0 {
		return nil
	}
	headQualifier, _ := javaSplitQualifiedHead(callTarget.qualifier)
	if len(ambiguousImportClassTargets[headQualifier]) > 0 {
		return nil
	}

	if classFQN, exists := importClassFQN[callTarget.qualifier]; exists {
		return methodIDsByClassFQN[classFQN][callTarget.methodName]
	}
	if classFQNs := wildcardClassFQNs[callTarget.qualifier]; len(classFQNs) > 0 {
		targetIDs := []string{}
		for _, classFQN := range classFQNs {
			targetIDs = append(targetIDs, methodIDsByClassFQN[classFQN][callTarget.methodName]...)
		}
		return uniqueStringSlice(targetIDs)
	}

	if classFQN, exists := localClasses[callTarget.qualifier]; exists {
		return methodIDsByClassFQN[classFQN][callTarget.methodName]
	}

	classFQN := javaQualifiedName(packageName, callTarget.qualifier)
	return methodIDsByClassFQN[classFQN][callTarget.methodName]
}

func extractJavaPackage(source string) (string, int) {
	matches := javaPackagePattern.FindStringSubmatchIndex(source)
	if len(matches) < 4 {
		return "", 0
	}

	packageName := source[matches[2]:matches[3]]
	line := 1 + strings.Count(source[:matches[0]], "\n")
	return packageName, line
}

func extractJavaImports(source string) []javaImportRef {
	matches := javaImportPattern.FindAllStringSubmatch(source, -1)
	if len(matches) == 0 {
		return []javaImportRef{}
	}

	imports := make([]javaImportRef, 0, len(matches))
	for _, match := range matches {
		importPath := ""
		isStatic := false
		if len(match) > 1 && strings.TrimSpace(match[1]) != "" {
			isStatic = true
		}
		if len(match) > 2 {
			importPath = strings.TrimSpace(match[2])
		}
		if importPath == "" {
			continue
		}

		isWildcard := strings.HasSuffix(importPath, ".*")
		simpleName := ""
		pathSegments := strings.Split(importPath, ".")
		if len(pathSegments) > 0 && !isWildcard {
			simpleName = pathSegments[len(pathSegments)-1]
		}

		imports = append(imports, javaImportRef{
			importPath: importPath,
			isStatic:   isStatic,
			isWildcard: isWildcard,
			simpleName: simpleName,
		})
	}

	return imports
}

func javaQualifiedName(packageName string, symbolName string) string {
	if symbolName == "" {
		return ""
	}
	if packageName == "" {
		return symbolName
	}
	return packageName + "." + symbolName
}

func javaJoinQualifiedName(prefix []string, name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	if len(prefix) == 0 {
		return trimmed
	}

	segments := append(append([]string{}, prefix...), trimmed)
	return strings.Join(segments, ".")
}

func computeJavaCyclomaticComplexity(node *tree_sitter.Node, source []byte, symbolKind string) int {
	if symbolKind != javaSymbolKindMethod {
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

		switch current.Kind() {
		case "if_statement", "for_statement", "enhanced_for_statement", "while_statement", "do_statement", "catch_clause", "switch_label":
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

func createJavaParseDiagnostic(file models.ScannedSourceFile, detail string) models.StructuredDiagnostic {
	return models.StructuredDiagnostic{
		Code:     javaParseErrorCode,
		Severity: models.SeverityError,
		Stage:    models.StageParse,
		Message:  "Failed to parse Java source file",
		FilePath: file.RelativePath,
		Language: file.Language,
		Detail:   detail,
	}
}

func buildJavaResolutionContext(entries []parsedJavaFile, moduleHints javaModuleHints) javaResolutionContext {
	context := javaResolutionContext{
		classSymbolByFQN:       map[string]string{},
		classSymbolIDsByFQN:    map[string][]string{},
		localClassFQNByFile:    map[string]map[string]string{},
		topLevelClassFQNsByPkg: map[string][]string{},
		methodIDsByClassFQN:    map[string]map[string][]string{},
		methodIDsByPackage:     map[string]map[string][]string{},
		ownerClassFQNByMethod:  map[string]string{},
		moduleByClassSymbolID:  map[string]string{},
	}

	for _, entry := range entries {
		if hasJavaErrorDiagnostics(entry.diagnostics) {
			continue
		}

		classesByName := map[string]string{}
		classFQNsBySimpleName := map[string][]string{}
		for _, symbolMatch := range entry.symbolMatches {
			switch symbolMatch.symbol.SymbolKind {
			case javaSymbolKindClass, javaSymbolKindInterface, javaSymbolKindEnum, javaSymbolKindRecord:
				fqn := javaQualifiedName(entry.packageName, symbolMatch.symbol.Name)
				classesByName[symbolMatch.symbol.Name] = fqn
				context.classSymbolByFQN[fqn] = symbolMatch.symbol.ID
				context.classSymbolIDsByFQN[fqn] = append(context.classSymbolIDsByFQN[fqn], symbolMatch.symbol.ID)
				moduleName := moduleHints.moduleForFile(entry.file.FilePath)
				if moduleName != "" {
					context.moduleByClassSymbolID[symbolMatch.symbol.ID] = moduleName
				}
				if !strings.Contains(symbolMatch.symbol.Name, ".") {
					context.topLevelClassFQNsByPkg[entry.packageName] = append(
						context.topLevelClassFQNsByPkg[entry.packageName],
						fqn,
					)
				}
				simpleName := javaLastIdentifierSegment(symbolMatch.symbol.Name)
				classFQNsBySimpleName[simpleName] = append(classFQNsBySimpleName[simpleName], fqn)
			}
		}
		for simpleName, fqns := range classFQNsBySimpleName {
			uniqueFQNs := uniqueStringSlice(fqns)
			if len(uniqueFQNs) == 1 {
				classesByName[simpleName] = uniqueFQNs[0]
			}
		}
		context.localClassFQNByFile[entry.file.FilePath] = classesByName

		for _, symbolMatch := range entry.symbolMatches {
			if symbolMatch.symbol.SymbolKind != javaSymbolKindMethod {
				continue
			}

			classFQN := javaQualifiedName(entry.packageName, symbolMatch.ownerType)
			methodsByName := context.methodIDsByClassFQN[classFQN]
			if methodsByName == nil {
				methodsByName = map[string][]string{}
				context.methodIDsByClassFQN[classFQN] = methodsByName
			}
			methodsByName[symbolMatch.symbol.Name] = append(methodsByName[symbolMatch.symbol.Name], symbolMatch.symbol.ID)

			pkgMethodsByName := context.methodIDsByPackage[entry.packageName]
			if pkgMethodsByName == nil {
				pkgMethodsByName = map[string][]string{}
				context.methodIDsByPackage[entry.packageName] = pkgMethodsByName
			}
			pkgMethodsByName[symbolMatch.symbol.Name] = append(pkgMethodsByName[symbolMatch.symbol.Name], symbolMatch.symbol.ID)
			context.ownerClassFQNByMethod[symbolMatch.symbol.ID] = classFQN
		}
	}

	for classFQN, methodsByName := range context.methodIDsByClassFQN {
		for methodName, ids := range methodsByName {
			context.methodIDsByClassFQN[classFQN][methodName] = uniqueStringSlice(ids)
		}
	}
	for packageName, methodsByName := range context.methodIDsByPackage {
		for methodName, ids := range methodsByName {
			context.methodIDsByPackage[packageName][methodName] = uniqueStringSlice(ids)
		}
	}
	for packageName, fqns := range context.topLevelClassFQNsByPkg {
		context.topLevelClassFQNsByPkg[packageName] = uniqueStringSlice(fqns)
	}
	for classFQN, symbolIDs := range context.classSymbolIDsByFQN {
		context.classSymbolIDsByFQN[classFQN] = uniqueStringSlice(symbolIDs)
	}

	return context
}

func buildJavaImportLookupIndexes(
	imports []javaImportRef,
	classSymbolIDsByFQN map[string][]string,
	topLevelClassFQNsByPkg map[string][]string,
	methodIDsByClassFQN map[string]map[string][]string,
	moduleByClassSymbolID map[string]string,
	currentModule string,
	moduleDependencies map[string]map[string]struct{},
) javaImportLookupIndexes {
	classSymbolByFQN := map[string]string{}
	for classFQN, symbolIDs := range classSymbolIDsByFQN {
		classSymbolID := selectPreferredJavaClassSymbolID(
			symbolIDs,
			moduleByClassSymbolID,
			currentModule,
			moduleDependencies,
		)
		if classSymbolID == "" {
			continue
		}
		classSymbolByFQN[classFQN] = classSymbolID
	}

	importClassFQN := map[string]string{}
	wildcardClassFQNs := map[string][]string{}
	staticMethodImportIDs := map[string][]string{}
	importClassCandidates := map[string][]string{}
	for _, importRef := range imports {
		if importRef.importPath == "" {
			continue
		}
		if importRef.isWildcard {
			if importRef.isStatic {
				continue
			}
			packagePath := strings.TrimSuffix(importRef.importPath, ".*")
			for _, classFQN := range topLevelClassFQNsByPkg[packagePath] {
				simpleName := javaLastIdentifierSegment(classFQN)
				if simpleName == "" {
					continue
				}
				wildcardClassFQNs[simpleName] = append(wildcardClassFQNs[simpleName], classFQN)
			}
			continue
		}

		if importRef.simpleName != "" {
			if _, resolved := classSymbolByFQN[importRef.importPath]; resolved {
				importClassCandidates[importRef.simpleName] = append(
					importClassCandidates[importRef.simpleName],
					importRef.importPath,
				)
				qualifiedTypeName := javaTypeQualifierFromFQN(importRef.importPath)
				if qualifiedTypeName != "" {
					importClassCandidates[qualifiedTypeName] = append(
						importClassCandidates[qualifiedTypeName],
						importRef.importPath,
					)
				}
			}
		}

		if !importRef.isStatic || importRef.simpleName == "" {
			continue
		}

		pathSegments := strings.Split(importRef.importPath, ".")
		if len(pathSegments) < 2 {
			continue
		}

		classFQN := strings.Join(pathSegments[:len(pathSegments)-1], ".")
		methodName := pathSegments[len(pathSegments)-1]
		staticMethodImportIDs[methodName] = uniqueStringSlice(append(
			staticMethodImportIDs[methodName],
			methodIDsByClassFQN[classFQN][methodName]...,
		))
	}

	for simpleName, classFQNs := range wildcardClassFQNs {
		candidates := uniqueStringSlice(classFQNs)
		preferred := preferJavaClassFQNsForModule(
			candidates,
			classSymbolByFQN,
			moduleByClassSymbolID,
			currentModule,
			moduleDependencies,
		)
		if len(preferred) > 0 {
			wildcardClassFQNs[simpleName] = preferred
			continue
		}
		wildcardClassFQNs[simpleName] = candidates
	}

	ambiguousImportClassTargets := map[string][]string{}
	for importKey, candidates := range importClassCandidates {
		uniqueCandidates := uniqueStringSlice(candidates)
		preferredCandidates := preferJavaClassFQNsForModule(
			uniqueCandidates,
			classSymbolByFQN,
			moduleByClassSymbolID,
			currentModule,
			moduleDependencies,
		)
		if len(preferredCandidates) > 0 {
			uniqueCandidates = preferredCandidates
		}
		if len(uniqueCandidates) == 1 {
			importClassFQN[importKey] = uniqueCandidates[0]
			continue
		}
		ambiguousImportClassTargets[importKey] = uniqueCandidates
	}

	return javaImportLookupIndexes{
		classSymbolByFQN:            classSymbolByFQN,
		importClassFQN:              importClassFQN,
		wildcardClassFQNs:           wildcardClassFQNs,
		staticMethodImportIDs:       staticMethodImportIDs,
		ambiguousImportClassTargets: ambiguousImportClassTargets,
	}
}

func (javaDeepResolver) Resolve(
	entry parsedJavaFile,
	ctx javaResolutionContext,
	classSymbolByFQN map[string]string,
	localClasses map[string]string,
	importClassFQN map[string]string,
	wildcardClassFQNs map[string][]string,
	staticMethodImportIDs map[string][]string,
	ambiguousImportClassTargets map[string][]string,
	unused []javaUnresolvedRef,
) ([]models.RelationEdge, []javaUnresolvedRef) {
	_ = unused
	relations := []models.RelationEdge{}
	unresolved := []javaUnresolvedRef{}

	relationKeys := map[string]struct{}{}
	for _, importRef := range entry.imports {
		importRef := importRef
		resolvedRelations, targetHint, reason, ok := resolveJavaDeepImport(
			entry.file.ID,
			importRef,
			classSymbolByFQN,
			ctx.topLevelClassFQNsByPkg,
			ctx.methodIDsByClassFQN,
		)
		if ok {
			for _, relation := range resolvedRelations {
				pushUniqueRelation(&relations, relationKeys, relation)
			}
			continue
		}

		unresolved = append(unresolved, javaUnresolvedRef{
			importRef:    &importRef,
			reason:       reason,
			relationType: models.RelReferences,
			sourceID:     entry.file.ID,
			targetHint:   targetHint,
		})
	}

	for _, symbolMatch := range entry.symbolMatches {
		if symbolMatch.symbol.SymbolKind != javaSymbolKindMethod {
			continue
		}

		ownerClassFQN := ctx.ownerClassFQNByMethod[symbolMatch.symbol.ID]
		for _, callTarget := range symbolMatch.callTargets {
			callTarget := callTarget
			targetID, reason := resolveJavaDeepCallTarget(
				callTarget,
				ownerClassFQN,
				entry.packageName,
				localClasses,
				importClassFQN,
				wildcardClassFQNs,
				staticMethodImportIDs,
				ambiguousImportClassTargets,
				ctx.methodIDsByClassFQN,
			)
			if targetID == "" {
				unresolved = append(unresolved, javaUnresolvedRef{
					callTarget:   &callTarget,
					reason:       reason,
					relationType: models.RelCalls,
					sourceID:     symbolMatch.symbol.ID,
					targetHint:   formatJavaCallTargetHint(callTarget),
				})
				continue
			}
			if targetID == symbolMatch.symbol.ID {
				continue
			}

			pushUniqueRelation(&relations, relationKeys, models.RelationEdge{
				FromID:     symbolMatch.symbol.ID,
				ToID:       targetID,
				Type:       models.RelCalls,
				Confidence: models.ConfidenceSemantic,
			})
		}
	}

	return relations, unresolved
}

func (javaSyntacticResolver) Resolve(
	entry parsedJavaFile,
	ctx javaResolutionContext,
	classSymbolByFQN map[string]string,
	localClasses map[string]string,
	importClassFQN map[string]string,
	wildcardClassFQNs map[string][]string,
	staticMethodImportIDs map[string][]string,
	ambiguousImportClassTargets map[string][]string,
	unresolved []javaUnresolvedRef,
) ([]models.RelationEdge, []javaUnresolvedRef) {
	_ = wildcardClassFQNs
	relations := []models.RelationEdge{}
	relationKeys := map[string]struct{}{}

	for _, unresolvedRef := range unresolved {
		if unresolvedRef.relationType == models.RelReferences && unresolvedRef.importRef != nil {
			relation, ok := resolveJavaSyntacticImport(entry.file.ID, *unresolvedRef.importRef, classSymbolByFQN)
			if !ok {
				continue
			}
			pushUniqueRelation(&relations, relationKeys, relation)
			continue
		}

		if unresolvedRef.relationType != models.RelCalls || unresolvedRef.callTarget == nil {
			continue
		}
		targetIDs := resolveJavaCallTarget(
			*unresolvedRef.callTarget,
			entry.packageName,
			localClasses,
			importClassFQN,
			wildcardClassFQNs,
			staticMethodImportIDs,
			ambiguousImportClassTargets,
			ctx.methodIDsByClassFQN,
			ctx.methodIDsByPackage,
		)
		targetIDs = uniqueStringSlice(targetIDs)
		if len(targetIDs) != 1 || targetIDs[0] == unresolvedRef.sourceID {
			continue
		}

		pushUniqueRelation(&relations, relationKeys, models.RelationEdge{
			FromID:     unresolvedRef.sourceID,
			ToID:       targetIDs[0],
			Type:       models.RelCalls,
			Confidence: models.ConfidenceSyntactic,
		})
	}

	return relations, nil
}
func resolveJavaDeepImport(
	sourceID string,
	importRef javaImportRef,
	classSymbolByFQN map[string]string,
	topLevelClassFQNsByPkg map[string][]string,
	methodIDsByClassFQN map[string]map[string][]string,
) ([]models.RelationEdge, string, string, bool) {
	if importRef.importPath == "" {
		return nil, "", "missing-import-path", false
	}
	if importRef.isWildcard {
		packagePath := strings.TrimSuffix(importRef.importPath, ".*")
		classFQNs := topLevelClassFQNsByPkg[packagePath]
		if len(classFQNs) == 0 {
			return nil, importRef.importPath, "missing-wildcard-package", false
		}

		relations := make([]models.RelationEdge, 0, len(classFQNs))
		for _, classFQN := range classFQNs {
			targetClassID, resolved := classSymbolByFQN[classFQN]
			if !resolved {
				continue
			}
			relations = append(relations, models.RelationEdge{
				FromID:     sourceID,
				ToID:       targetClassID,
				Type:       models.RelReferences,
				Confidence: models.ConfidenceSemantic,
			})
		}
		if len(relations) == 0 {
			return nil, importRef.importPath, "missing-wildcard-symbols", false
		}

		return relations, importRef.importPath, "", true
	}

	if importRef.isStatic {
		pathSegments := strings.Split(importRef.importPath, ".")
		if len(pathSegments) < 2 {
			return nil, importRef.importPath, "invalid-static-import", false
		}
		classFQN := strings.Join(pathSegments[:len(pathSegments)-1], ".")
		methodName := pathSegments[len(pathSegments)-1]
		targetIDs := uniqueStringSlice(methodIDsByClassFQN[classFQN][methodName])
		if len(targetIDs) != 1 {
			if len(targetIDs) == 0 {
				return nil, importRef.importPath, "missing-static-target", false
			}
			return nil, importRef.importPath, "ambiguous-static-target", false
		}
		return []models.RelationEdge{{
			FromID:     sourceID,
			ToID:       targetIDs[0],
			Type:       models.RelReferences,
			Confidence: models.ConfidenceSemantic,
		}}, importRef.importPath, "", true
	}

	targetClassID, resolved := classSymbolByFQN[importRef.importPath]
	if !resolved {
		return nil, importRef.importPath, "missing-class-symbol", false
	}

	return []models.RelationEdge{{
		FromID:     sourceID,
		ToID:       targetClassID,
		Type:       models.RelReferences,
		Confidence: models.ConfidenceSemantic,
	}}, importRef.importPath, "", true
}

func resolveJavaSyntacticImport(
	sourceID string,
	importRef javaImportRef,
	classSymbolByFQN map[string]string,
) (models.RelationEdge, bool) {
	if importRef.importPath == "" || importRef.isWildcard {
		return models.RelationEdge{}, false
	}
	targetClassID, resolved := classSymbolByFQN[importRef.importPath]
	if !resolved {
		return models.RelationEdge{}, false
	}

	return models.RelationEdge{
		FromID:     sourceID,
		ToID:       targetClassID,
		Type:       models.RelReferences,
		Confidence: models.ConfidenceSyntactic,
	}, true
}

func resolveJavaDeepCallTarget(
	callTarget javaCallTarget,
	ownerClassFQN string,
	packageName string,
	localClasses map[string]string,
	importClassFQN map[string]string,
	wildcardClassFQNs map[string][]string,
	staticMethodImportIDs map[string][]string,
	ambiguousImportClassTargets map[string][]string,
	methodIDsByClassFQN map[string]map[string][]string,
) (string, string) {
	if callTarget.methodName == "" {
		return "", "missing-method-name"
	}

	if callTarget.qualifier == "" {
		if ids := uniqueStringSlice(staticMethodImportIDs[callTarget.methodName]); len(ids) > 0 {
			if len(ids) == 1 {
				return ids[0], ""
			}
			return "", "ambiguous-static-call-target"
		}
		if ownerClassFQN == "" {
			return "", "missing-owner-class"
		}

		ownerMethodIDs := uniqueStringSlice(methodIDsByClassFQN[ownerClassFQN][callTarget.methodName])
		if len(ownerMethodIDs) == 1 {
			return ownerMethodIDs[0], ""
		}
		if len(ownerMethodIDs) > 1 {
			return "", "ambiguous-owner-method"
		}
		return "", "missing-qualifier-metadata"
	}

	classCandidates := resolveJavaClassCandidates(
		callTarget.qualifier,
		packageName,
		localClasses,
		importClassFQN,
		wildcardClassFQNs,
		ambiguousImportClassTargets,
		methodIDsByClassFQN,
	)
	if len(classCandidates) == 0 {
		if len(ambiguousImportClassTargets[callTarget.qualifier]) > 0 {
			return "", "ambiguous-import-class"
		}
		headQualifier, _ := javaSplitQualifiedHead(callTarget.qualifier)
		if len(ambiguousImportClassTargets[headQualifier]) > 0 {
			return "", "ambiguous-import-class"
		}
		return "", "unresolved-qualifier"
	}

	classCandidates = uniqueStringSlice(classCandidates)

	targetIDs := []string{}
	for _, classFQN := range classCandidates {
		targetIDs = append(targetIDs, methodIDsByClassFQN[classFQN][callTarget.methodName]...)
	}
	targetIDs = uniqueStringSlice(targetIDs)
	if len(targetIDs) == 1 {
		return targetIDs[0], ""
	}
	if len(targetIDs) > 1 {
		return "", "ambiguous-qualified-target"
	}

	return "", "missing-qualified-method"
}

func resolveJavaClassCandidates(
	qualifier string,
	packageName string,
	localClasses map[string]string,
	importClassFQN map[string]string,
	wildcardClassFQNs map[string][]string,
	ambiguousImportClassTargets map[string][]string,
	methodIDsByClassFQN map[string]map[string][]string,
) []string {
	candidates := []string{}
	if qualifier == "" {
		return candidates
	}
	if len(ambiguousImportClassTargets[qualifier]) > 0 {
		return candidates
	}

	if classFQN, exists := importClassFQN[qualifier]; exists {
		candidates = append(candidates, classFQN)
	}
	if classFQN, exists := localClasses[qualifier]; exists {
		candidates = append(candidates, classFQN)
	}
	candidates = append(candidates, wildcardClassFQNs[qualifier]...)

	headQualifier, tailQualifier := javaSplitQualifiedHead(qualifier)
	if headQualifier != "" && tailQualifier != "" {
		if len(ambiguousImportClassTargets[headQualifier]) > 0 {
			return []string{}
		}
		if classFQN, exists := importClassFQN[headQualifier]; exists {
			candidates = append(candidates, classFQN+"."+tailQualifier)
		}
		if classFQN, exists := localClasses[headQualifier]; exists {
			candidates = append(candidates, classFQN+"."+tailQualifier)
		}
		for _, classFQN := range wildcardClassFQNs[headQualifier] {
			candidates = append(candidates, classFQN+"."+tailQualifier)
		}
	}

	candidates = append(candidates, javaQualifiedName(packageName, qualifier))

	filtered := make([]string, 0, len(candidates))
	for _, classFQN := range uniqueStringSlice(candidates) {
		if _, exists := methodIDsByClassFQN[classFQN]; !exists {
			continue
		}
		filtered = append(filtered, classFQN)
	}

	return filtered
}

func javaSplitQualifiedHead(value string) (string, string) {
	dotIdx := strings.Index(value, ".")
	if dotIdx <= 0 || dotIdx >= len(value)-1 {
		return "", ""
	}
	return value[:dotIdx], value[dotIdx+1:]
}

func javaTypeQualifierFromFQN(typeFQN string) string {
	segments := strings.Split(strings.TrimSpace(typeFQN), ".")
	if len(segments) == 0 {
		return ""
	}

	firstTypeIdx := -1
	for idx, segment := range segments {
		if segment == "" {
			continue
		}
		firstRune := rune(segment[0])
		if firstRune >= 'A' && firstRune <= 'Z' {
			firstTypeIdx = idx
			break
		}
	}
	if firstTypeIdx == -1 || firstTypeIdx >= len(segments) {
		return ""
	}

	return strings.Join(segments[firstTypeIdx:], ".")
}

func formatJavaCallTargetHint(callTarget javaCallTarget) string {
	if callTarget.qualifier == "" {
		return callTarget.methodName
	}
	return callTarget.qualifier + "." + callTarget.methodName
}

func createJavaResolutionFallbackDiagnostic(
	file models.GraphFile,
	unresolved []javaUnresolvedRef,
) models.StructuredDiagnostic {
	sorted := append([]javaUnresolvedRef(nil), unresolved...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].sourceID != sorted[j].sourceID {
			return sorted[i].sourceID < sorted[j].sourceID
		}
		if sorted[i].targetHint != sorted[j].targetHint {
			return sorted[i].targetHint < sorted[j].targetHint
		}
		return sorted[i].reason < sorted[j].reason
	})

	detailParts := make([]string, 0, minInt(len(sorted), javaFallbackDiagnosticMaxEntries))
	detailLength := 0
	omittedCount := 0
	for _, unresolvedRef := range sorted {
		relationLabel := string(unresolvedRef.relationType)
		if relationLabel == "" {
			relationLabel = "relation"
		}
		part := fmt.Sprintf(
			"%s:%s (%s)",
			relationLabel,
			unresolvedRef.targetHint,
			unresolvedRef.reason,
		)
		separatorLength := 0
		if len(detailParts) > 0 {
			separatorLength = 2
		}
		if len(detailParts) >= javaFallbackDiagnosticMaxEntries ||
			detailLength+separatorLength+len(part) > javaDiagnosticDetailMaxBytes {
			omittedCount++
			continue
		}
		detailParts = append(detailParts, part)
		detailLength += separatorLength + len(part)
	}
	if omittedCount > 0 {
		detailParts = append(detailParts, fmt.Sprintf("%s (%d entries omitted)", javaDiagnosticTruncationPrefixKey, omittedCount))
	}

	return models.StructuredDiagnostic{
		Code:     javaResolutionFallbackCode,
		Severity: models.SeverityWarning,
		Stage:    models.StageParse,
		Message:  "Deep Java resolution fallback applied",
		FilePath: file.FilePath,
		Language: file.Language,
		Detail:   joinDiagnosticPartsWithinLimit(detailParts, javaDiagnosticDetailMaxBytes),
	}
}

func sortJavaDiagnostics(diagnostics []models.StructuredDiagnostic) {
	sort.Slice(diagnostics, func(i, j int) bool {
		left := diagnostics[i]
		right := diagnostics[j]
		if left.Severity != right.Severity {
			return left.Severity < right.Severity
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		if left.FilePath != right.FilePath {
			return left.FilePath < right.FilePath
		}
		return left.Detail < right.Detail
	})
}

func sortRelationEdges(relations []models.RelationEdge) {
	sort.Slice(relations, func(i, j int) bool {
		left := relations[i]
		right := relations[j]

		if left.FromID != right.FromID {
			return left.FromID < right.FromID
		}
		if left.Type != right.Type {
			return left.Type < right.Type
		}
		if left.ToID != right.ToID {
			return left.ToID < right.ToID
		}
		return left.Confidence < right.Confidence
	})
}

func createJavaModuleHintDiagnostic(file models.GraphFile, warnings []string) models.StructuredDiagnostic {
	details := uniqueStringSlice(warnings)
	omittedCount := 0
	if len(details) > javaModuleHintWarningMaxEntries {
		omittedCount = len(details) - javaModuleHintWarningMaxEntries
		details = append([]string{}, details[:javaModuleHintWarningMaxEntries]...)
	}
	if omittedCount > 0 {
		details = append(details, fmt.Sprintf("%s (%d warnings omitted)", javaDiagnosticTruncationPrefixKey, omittedCount))
	}
	return models.StructuredDiagnostic{
		Code:     javaModuleHintWarningCode,
		Severity: models.SeverityWarning,
		Stage:    models.StageParse,
		Message:  "Java module metadata hints parsed with warnings",
		FilePath: file.FilePath,
		Language: file.Language,
		Detail:   joinDiagnosticPartsWithinLimit(details, javaDiagnosticDetailMaxBytes),
	}
}

func joinDiagnosticPartsWithinLimit(parts []string, maxBytes int) string {
	if len(parts) == 0 || maxBytes <= 0 {
		return ""
	}

	builder := strings.Builder{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if builder.Len() == 0 {
			if len(part) > maxBytes {
				return part[:maxBytes]
			}
			builder.WriteString(part)
			continue
		}

		required := 2 + len(part)
		if builder.Len()+required > maxBytes {
			break
		}
		builder.WriteString("; ")
		builder.WriteString(part)
	}

	return builder.String()
}

func minInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}

func hasJavaErrorDiagnostics(diagnostics []models.StructuredDiagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == models.SeverityError {
			return true
		}
	}

	return false
}

func discoverJavaModuleHints(rootPath string) javaModuleHints {
	hints := javaModuleHints{
		moduleDependencies:  map[string]map[string]struct{}{},
		fileModuleBySrcRoot: map[string]string{},
		warnings:            []string{},
	}
	if strings.TrimSpace(rootPath) == "" {
		return hints
	}

	resolvedRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return hints
	}

	loadGradleModuleHints(resolvedRootPath, &hints)
	loadMavenModuleHints(resolvedRootPath, &hints)
	return hints
}

func (h javaModuleHints) moduleForFile(relativePath string) string {
	cleanPath := filepath.ToSlash(strings.TrimSpace(relativePath))
	if cleanPath == "" {
		return ""
	}

	bestRoot := ""
	bestModule := ""
	for sourceRoot, moduleName := range h.fileModuleBySrcRoot {
		if sourceRoot == "" || moduleName == "" {
			continue
		}
		if cleanPath != sourceRoot && !strings.HasPrefix(cleanPath, sourceRoot+"/") {
			continue
		}
		if len(sourceRoot) > len(bestRoot) || (len(sourceRoot) == len(bestRoot) && sourceRoot < bestRoot) {
			bestRoot = sourceRoot
			bestModule = moduleName
		}
	}

	return bestModule
}

func loadGradleModuleHints(rootPath string, hints *javaModuleHints) {
	if hints == nil {
		return
	}

	for _, settingsFile := range []string{"settings.gradle", "settings.gradle.kts"} {
		settingsPath := filepath.Join(rootPath, settingsFile)
		content, err := os.ReadFile(settingsPath)
		if err != nil {
			if !os.IsNotExist(err) {
				hints.warnings = append(hints.warnings, fmt.Sprintf("%s: %v", settingsFile, err))
			}
			continue
		}

		modulePaths := parseGradleIncludedModules(string(content))
		if strings.Contains(string(content), "include") && len(modulePaths) == 0 {
			hints.warnings = append(hints.warnings, fmt.Sprintf("%s: include declarations malformed", settingsFile))
		}
		for _, modulePath := range modulePaths {
			moduleName := javaModuleNameFromPath(modulePath)
			addJavaModuleSourceRoots(hints, moduleName, modulePath)

			for _, buildFile := range []string{"build.gradle", "build.gradle.kts"} {
				buildPath := filepath.Join(rootPath, filepath.FromSlash(modulePath), buildFile)
				buildContent, readErr := os.ReadFile(buildPath)
				if readErr != nil {
					continue
				}

				dependencies := parseGradleProjectDependencies(string(buildContent))
				if strings.Contains(string(buildContent), "project(") && len(dependencies) == 0 {
					hints.warnings = append(
						hints.warnings,
						fmt.Sprintf("%s: project() dependencies malformed", filepath.ToSlash(strings.TrimPrefix(buildPath, rootPath+"/"))),
					)
				}
				for _, dependencyPath := range dependencies {
					dependencyModule := javaModuleNameFromPath(dependencyPath)
					if dependencyModule == "" {
						continue
					}
					addJavaModuleDependency(hints, moduleName, dependencyModule)
				}
			}
		}
	}
}

func loadMavenModuleHints(rootPath string, hints *javaModuleHints) {
	if hints == nil {
		return
	}

	rootPomPath := filepath.Join(rootPath, "pom.xml")
	rootPomContent, err := os.ReadFile(rootPomPath)
	if err != nil {
		if !os.IsNotExist(err) {
			hints.warnings = append(hints.warnings, fmt.Sprintf("pom.xml: %v", err))
		}
		return
	}

	modulePaths, dependencyArtifacts, malformed := parseMavenPomSignals(string(rootPomContent))
	if malformed {
		hints.warnings = append(hints.warnings, "pom.xml: malformed module/dependency metadata")
	}
	if len(modulePaths) == 0 {
		modulePaths = []string{"."}
	}

	artifactToModule := map[string]string{}
	moduleDependencyArtifacts := map[string][]string{}
	for _, modulePath := range modulePaths {
		modulePath = normalizeJavaModulePath(modulePath)
		moduleName := javaModuleNameFromPath(modulePath)
		addJavaModuleSourceRoots(hints, moduleName, modulePath)
		artifactToModule[moduleName] = moduleName

		modulePomPath := filepath.Join(rootPath, filepath.FromSlash(modulePath), "pom.xml")
		modulePomContent, readErr := os.ReadFile(modulePomPath)
		if readErr != nil {
			if modulePath == "." {
				moduleDependencyArtifacts[moduleName] = append(moduleDependencyArtifacts[moduleName], dependencyArtifacts...)
			}
			continue
		}

		artifactID, moduleDeps, moduleMalformed := parseMavenModulePomSignals(string(modulePomContent))
		if moduleMalformed {
			hints.warnings = append(
				hints.warnings,
				fmt.Sprintf("%s: malformed module/dependency metadata", filepath.ToSlash(filepath.Join(modulePath, "pom.xml"))),
			)
		}
		if artifactID != "" {
			artifactToModule[artifactID] = moduleName
		}
		moduleDependencyArtifacts[moduleName] = append(moduleDependencyArtifacts[moduleName], moduleDeps...)
	}

	for moduleName, artifacts := range moduleDependencyArtifacts {
		for _, artifact := range uniqueStringSlice(artifacts) {
			dependencyModule := artifactToModule[artifact]
			if dependencyModule == "" || dependencyModule == moduleName {
				continue
			}
			addJavaModuleDependency(hints, moduleName, dependencyModule)
		}
	}
}

func parseGradleIncludedModules(content string) []string {
	matches := javaGradleIncludeLinePattern.FindAllStringSubmatch(content, -1)
	modules := []string{}
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		tokens := javaQuotedTokenPattern.FindAllStringSubmatch(match[1], -1)
		for _, token := range tokens {
			if len(token) < 2 {
				continue
			}
			modulePath := normalizeJavaModulePath(token[1])
			if modulePath == "" {
				continue
			}
			modules = append(modules, modulePath)
		}
	}

	return uniqueStringSlice(modules)
}

func parseGradleProjectDependencies(content string) []string {
	matches := javaGradleProjectDepPattern.FindAllStringSubmatch(content, -1)
	modules := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		modulePath := normalizeJavaModulePath(match[1])
		if modulePath == "" {
			continue
		}
		modules = append(modules, modulePath)
	}

	return uniqueStringSlice(modules)
}

func parseMavenPomSignals(content string) ([]string, []string, bool) {
	moduleMatches := javaMavenModulePattern.FindAllStringSubmatch(content, -1)
	modules := []string{}
	for _, match := range moduleMatches {
		if len(match) < 2 {
			continue
		}
		modulePath := normalizeJavaModulePath(match[1])
		if modulePath == "" {
			continue
		}
		modules = append(modules, modulePath)
	}

	dependencies := parseMavenDependencyArtifacts(content)
	malformed := strings.Contains(content, "<project") &&
		((strings.Contains(content, "<modules>") && len(modules) == 0) ||
			(strings.Contains(content, "<dependency>") && len(dependencies) == 0))
	return uniqueStringSlice(modules), dependencies, malformed
}

func parseMavenModulePomSignals(content string) (string, []string, bool) {
	artifactID := ""
	match := javaMavenArtifactPattern.FindStringSubmatch(content)
	if len(match) > 1 {
		artifactID = strings.TrimSpace(match[1])
	}
	dependencies := parseMavenDependencyArtifacts(content)
	malformed := strings.Contains(content, "<project") &&
		strings.Contains(content, "<dependency>") &&
		len(dependencies) == 0
	return artifactID, dependencies, malformed
}

func parseMavenDependencyArtifacts(content string) []string {
	matches := javaMavenDependencyPattern.FindAllStringSubmatch(content, -1)
	artifacts := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		artifactID := strings.TrimSpace(match[1])
		if artifactID == "" {
			continue
		}
		artifacts = append(artifacts, artifactID)
	}

	return uniqueStringSlice(artifacts)
}

func addJavaModuleSourceRoots(hints *javaModuleHints, moduleName string, modulePath string) {
	if hints == nil || moduleName == "" {
		return
	}

	modulePath = normalizeJavaModulePath(modulePath)
	for _, sourceRoot := range []string{
		normalizeJavaModulePath(filepath.ToSlash(filepath.Join(modulePath, "src", "main", "java"))),
		normalizeJavaModulePath(filepath.ToSlash(filepath.Join(modulePath, "src", "test", "java"))),
		normalizeJavaModulePath(filepath.ToSlash(filepath.Join(modulePath, "src"))),
	} {
		if sourceRoot == "" {
			continue
		}
		if existing, exists := hints.fileModuleBySrcRoot[sourceRoot]; exists && existing != moduleName {
			hints.warnings = append(
				hints.warnings,
				fmt.Sprintf("conflicting module mapping for %s: %s vs %s", sourceRoot, existing, moduleName),
			)
			continue
		}
		hints.fileModuleBySrcRoot[sourceRoot] = moduleName
	}
}

func addJavaModuleDependency(hints *javaModuleHints, fromModule string, toModule string) {
	if hints == nil || fromModule == "" || toModule == "" || fromModule == toModule {
		return
	}

	if hints.moduleDependencies[fromModule] == nil {
		hints.moduleDependencies[fromModule] = map[string]struct{}{}
	}
	hints.moduleDependencies[fromModule][toModule] = struct{}{}
}

func javaModuleNameFromPath(modulePath string) string {
	modulePath = normalizeJavaModulePath(modulePath)
	if modulePath == "" || modulePath == "." {
		return "root"
	}

	segments := strings.Split(modulePath, "/")
	if len(segments) == 0 {
		return "root"
	}
	return strings.TrimSpace(segments[len(segments)-1])
}

func normalizeJavaModulePath(value string) string {
	normalized := filepath.ToSlash(strings.TrimSpace(value))
	normalized = strings.Trim(normalized, `"'`)
	normalized = strings.TrimSpace(normalized)
	if normalized == "" {
		return ""
	}
	for strings.HasPrefix(normalized, ":") {
		normalized = strings.TrimPrefix(normalized, ":")
	}
	normalized = strings.ReplaceAll(normalized, ":", "/")
	normalized = strings.TrimPrefix(normalized, "./")
	normalized = strings.TrimSpace(normalized)
	if normalized == "" {
		return "."
	}
	return normalized
}

func selectPreferredJavaClassSymbolID(
	candidates []string,
	moduleByClassSymbolID map[string]string,
	currentModule string,
	moduleDependencies map[string]map[string]struct{},
) string {
	if len(candidates) == 0 {
		return ""
	}
	if currentModule == "" {
		return candidates[0]
	}

	preferred := []string{}
	for _, candidate := range candidates {
		targetModule := moduleByClassSymbolID[candidate]
		if javaModuleDependencyAllowed(currentModule, targetModule, moduleDependencies) {
			preferred = append(preferred, candidate)
		}
	}
	if len(preferred) == 0 {
		return candidates[0]
	}
	return preferred[0]
}

func preferJavaClassFQNsForModule(
	candidates []string,
	classSymbolByFQN map[string]string,
	moduleByClassSymbolID map[string]string,
	currentModule string,
	moduleDependencies map[string]map[string]struct{},
) []string {
	if len(candidates) == 0 || currentModule == "" {
		return nil
	}

	preferred := []string{}
	for _, classFQN := range candidates {
		symbolID := classSymbolByFQN[classFQN]
		targetModule := moduleByClassSymbolID[symbolID]
		if !javaModuleDependencyAllowed(currentModule, targetModule, moduleDependencies) {
			continue
		}
		preferred = append(preferred, classFQN)
	}

	return uniqueStringSlice(preferred)
}

func javaModuleDependencyAllowed(
	currentModule string,
	targetModule string,
	moduleDependencies map[string]map[string]struct{},
) bool {
	if currentModule == "" || targetModule == "" {
		return false
	}
	if currentModule == targetModule {
		return true
	}
	_, allowed := moduleDependencies[currentModule][targetModule]
	return allowed
}
