package vault

import (
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/compozy/kb/internal/models"
)

type starterWikiArticle struct {
	Body    string
	Sources []string
	Summary string
	Title   string
}

func RenderDocuments(
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
	topic models.TopicMetadata,
) []models.RenderedDocument {
	documentLookup := createDocumentLookup(graph)
	externalNodeLookup := createExternalNodeLookup(graph.ExternalNodes)
	relationsBySource := groupRelationsBySource(graph.Relations)
	relationsByTarget := groupRelationsByTarget(graph.Relations)
	symbolsByFile := groupSymbolsByFile(graph.Symbols)
	filesByDirectory := groupFilesByDirectory(graph.Files)
	symbolsByDirectory := groupSymbolsByDirectory(graph.Symbols)
	filesByLanguage := groupFilesByLanguage(graph.Files)
	symbolsByLanguage := groupSymbolsByLanguage(graph.Symbols)
	symbolsByKind := groupSymbolsByKind(graph.Symbols)

	rawDocuments := make([]models.RenderedDocument, 0, len(graph.Files)+len(graph.Symbols)+len(filesByDirectory)+len(filesByLanguage))

	files := append([]models.GraphFile(nil), graph.Files...)
	sort.Slice(files, func(i, j int) bool {
		return files[i].FilePath < files[j].FilePath
	})

	for _, file := range files {
		fileSymbols := append([]models.SymbolNode(nil), symbolsByFile[file.FilePath]...)
		sortSymbolsByLocation(fileSymbols)

		rawDocuments = append(rawDocuments, renderRawFileDocument(
			topic,
			file,
			fileSymbols,
			relationsBySource[file.ID],
			relationsByTarget[file.ID],
			defaultFileMetrics(metrics.Files[file.ID]),
			documentLookup,
			externalNodeLookup,
		))
	}

	symbols := append([]models.SymbolNode(nil), graph.Symbols...)
	sortSymbolsByLocation(symbols)

	for _, symbol := range symbols {
		rawDocuments = append(rawDocuments, renderRawSymbolDocument(
			topic,
			symbol,
			relationsBySource[symbol.ID],
			relationsByTarget[symbol.ID],
			defaultSymbolMetrics(symbol, metrics.Symbols[symbol.ID]),
			documentLookup,
			externalNodeLookup,
		))
	}

	for _, directoryPath := range sortedMapKeys(filesByDirectory) {
		rawDocuments = append(rawDocuments, renderRawDirectoryIndex(
			topic,
			directoryPath,
			filesByDirectory[directoryPath],
			symbolsByDirectory[directoryPath],
			defaultDirectoryMetrics(metrics.Directories[directoryPath]),
		))
	}

	for _, language := range sortedMapKeys(filesByLanguage) {
		rawDocuments = append(rawDocuments, renderRawLanguageIndex(
			topic,
			language,
			filesByLanguage[language],
			symbolsByLanguage[language],
		))
	}

	hotspotFiles := topRelationHotspotFiles(graph.Files, relationsBySource, relationsByTarget)
	starterArticles := buildStarterWikiArticles(
		topic,
		graph,
		metrics,
		filesByLanguage,
		filesByDirectory,
		symbolsByDirectory,
		symbolsByKind,
		hotspotFiles,
	)

	wikiConcepts := make([]models.RenderedDocument, 0, len(starterArticles))
	for _, article := range starterArticles {
		wikiConcepts = append(wikiConcepts, renderWikiArticle(topic, article))
	}

	documents := make([]models.RenderedDocument, 0, len(rawDocuments)+len(wikiConcepts)+3)
	documents = append(documents, rawDocuments...)
	documents = append(documents, wikiConcepts...)
	documents = append(documents,
		renderDashboard(topic, graph, starterArticles),
		renderConceptIndex(topic, starterArticles),
		renderSourceIndex(topic, starterArticles),
	)

	for index := range documents {
		documents[index].Body = renderMarkdownDocument(documents[index])
	}

	sort.Slice(documents, func(i, j int) bool {
		return documents[i].RelativePath < documents[j].RelativePath
	})

	return documents
}

func renderFrontmatter(frontmatter map[string]interface{}) string {
	lines := []string{"---"}
	keys := sortedMapKeys(frontmatter)

	for _, key := range keys {
		value := frontmatter[key]

		switch typed := value.(type) {
		case []string:
			lines = append(lines, key+":")
			for _, item := range typed {
				lines = append(lines, "  - "+strconv.Quote(item))
			}
		case string:
			lines = append(lines, key+": "+strconv.Quote(typed))
		case bool:
			lines = append(lines, key+": "+strconv.FormatBool(typed))
		case int:
			lines = append(lines, key+": "+strconv.Itoa(typed))
		case int64:
			lines = append(lines, key+": "+strconv.FormatInt(typed, 10))
		case float64:
			lines = append(lines, key+": "+strconv.FormatFloat(typed, 'f', -1, 64))
		case float32:
			lines = append(lines, key+": "+strconv.FormatFloat(float64(typed), 'f', -1, 32))
		default:
			lines = append(lines, key+": "+strconv.Quote(fmt.Sprint(typed)))
		}
	}

	lines = append(lines, "---")
	return strings.Join(lines, "\n")
}

func renderMarkdownDocument(document models.RenderedDocument) string {
	return strings.Join([]string{
		renderFrontmatter(document.Frontmatter),
		"",
		strings.TrimRight(document.Body, "\n"),
		"",
	}, "\n")
}

func toSourceWikiLink(topic models.TopicMetadata, relativePath, label string) string {
	return ToTopicWikiLink(topic.Slug, relativePath, label)
}

func createDocumentLookup(graph models.GraphSnapshot) map[string]string {
	lookup := make(map[string]string, len(graph.Files)+len(graph.Symbols))

	for _, file := range graph.Files {
		lookup[file.ID] = GetRawFileDocumentPath(file.FilePath)
	}

	for _, symbol := range graph.Symbols {
		lookup[symbol.ID] = GetRawSymbolDocumentPath(symbol)
	}

	return lookup
}

func createExternalNodeLookup(externalNodes []models.ExternalNode) map[string]models.ExternalNode {
	lookup := make(map[string]models.ExternalNode, len(externalNodes))
	for _, node := range externalNodes {
		lookup[node.ID] = node
	}
	return lookup
}

func linkForNode(
	topic models.TopicMetadata,
	nodeID string,
	documentLookup map[string]string,
	externalNodes map[string]models.ExternalNode,
	fallbackLabel string,
) string {
	if documentPath, exists := documentLookup[nodeID]; exists {
		return toSourceWikiLink(topic, documentPath, fallbackLabel)
	}

	if externalNode, exists := externalNodes[nodeID]; exists {
		return "`" + externalNode.Label + "`"
	}

	if fallbackLabel != "" {
		return "`" + fallbackLabel + "`"
	}

	return "`" + nodeID + "`"
}

func renderRelationList(
	topic models.TopicMetadata,
	relations []models.RelationEdge,
	documentLookup map[string]string,
	externalNodes map[string]models.ExternalNode,
) []string {
	if len(relations) == 0 {
		return []string{"None"}
	}

	orderedRelations := append([]models.RelationEdge(nil), relations...)
	sort.Slice(orderedRelations, func(i, j int) bool {
		left := orderedRelations[i]
		right := orderedRelations[j]

		if left.Type != right.Type {
			return left.Type < right.Type
		}
		if left.Confidence != right.Confidence {
			return left.Confidence < right.Confidence
		}
		if left.ToID != right.ToID {
			return left.ToID < right.ToID
		}
		return left.FromID < right.FromID
	})

	lines := make([]string, 0, len(orderedRelations))
	for _, relation := range orderedRelations {
		targetLabel := ""
		if strings.HasPrefix(relation.ToID, "file:") {
			targetLabel = strings.TrimPrefix(relation.ToID, "file:")
		}

		target := linkForNode(topic, relation.ToID, documentLookup, externalNodes, targetLabel)
		lines = append(lines, fmt.Sprintf("- `%s` (%s) -> %s", relation.Type, relation.Confidence, target))
	}

	return lines
}

func renderBacklinkList(
	topic models.TopicMetadata,
	relations []models.RelationEdge,
	documentLookup map[string]string,
	externalNodes map[string]models.ExternalNode,
) []string {
	if len(relations) == 0 {
		return []string{"None"}
	}

	orderedRelations := append([]models.RelationEdge(nil), relations...)
	sort.Slice(orderedRelations, func(i, j int) bool {
		left := orderedRelations[i]
		right := orderedRelations[j]

		if left.Type != right.Type {
			return left.Type < right.Type
		}
		if left.Confidence != right.Confidence {
			return left.Confidence < right.Confidence
		}
		return left.FromID < right.FromID
	})

	lines := make([]string, 0, len(orderedRelations))
	for _, relation := range orderedRelations {
		sourceLabel := ""
		if strings.HasPrefix(relation.FromID, "file:") {
			sourceLabel = strings.TrimPrefix(relation.FromID, "file:")
		}

		source := linkForNode(topic, relation.FromID, documentLookup, externalNodes, sourceLabel)
		lines = append(lines, fmt.Sprintf("- %s via `%s` (%s)", source, relation.Type, relation.Confidence))
	}

	return lines
}

func renderSmellList(smells []string) string {
	if len(smells) == 0 {
		return "None"
	}

	quoted := make([]string, 0, len(smells))
	for _, smell := range smells {
		quoted = append(quoted, "`"+smell+"`")
	}
	return strings.Join(quoted, ", ")
}

func renderRawFileDocument(
	topic models.TopicMetadata,
	file models.GraphFile,
	fileSymbols []models.SymbolNode,
	outgoingRelations []models.RelationEdge,
	incomingRelations []models.RelationEdge,
	fileMetrics models.FileMetrics,
	documentLookup map[string]string,
	externalNodes map[string]models.ExternalNode,
) models.RenderedDocument {
	symbolLines := []string{"None"}
	if len(fileSymbols) > 0 {
		symbolLines = make([]string, 0, len(fileSymbols))
		for _, symbol := range fileSymbols {
			symbolLink := toSourceWikiLink(
				topic,
				GetRawSymbolDocumentPath(symbol),
				fmt.Sprintf("%s (%s)", symbol.Name, symbol.SymbolKind),
			)
			symbolLines = append(symbolLines, fmt.Sprintf("- %s · exported=%t", symbolLink, symbol.Exported))
		}
	}

	moduleDoc := strings.TrimSpace(file.ModuleDoc)
	if moduleDoc == "" {
		moduleDoc = "None"
	}

	sections := []string{
		fmt.Sprintf("# Codebase File: %s", file.FilePath),
		"",
		"Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.",
		"",
		"## Language",
		"`" + string(file.Language) + "`",
		"",
		"## Static Analysis",
		fmt.Sprintf("- Afferent coupling: %d", fileMetrics.AfferentCoupling),
		fmt.Sprintf("- Efferent coupling: %d", fileMetrics.EfferentCoupling),
		fmt.Sprintf("- Instability: %s", strconv.FormatFloat(fileMetrics.Instability, 'f', -1, 64)),
		fmt.Sprintf("- Entry point: %t", fileMetrics.IsEntryPoint),
		fmt.Sprintf("- Circular dependency: %t", fileMetrics.HasCircularDependency),
		fmt.Sprintf("- Smells: %s", renderSmellList(fileMetrics.Smells)),
		"",
		"## Module Notes",
		moduleDoc,
		"",
		"## Symbols",
	}
	sections = append(sections, symbolLines...)
	sections = append(sections,
		"",
		"## Outgoing Relations",
	)
	sections = append(sections, renderRelationList(topic, outgoingRelations, documentLookup, externalNodes)...)
	sections = append(sections,
		"",
		"## Backlinks",
	)
	sections = append(sections, renderBacklinkList(topic, incomingRelations, documentLookup, externalNodes)...)

	return models.RenderedDocument{
		Kind:         models.DocRaw,
		ManagedArea:  models.AreaRawCodebase,
		RelativePath: GetRawFileDocumentPath(file.FilePath),
		Frontmatter: map[string]interface{}{
			"afferent_coupling":       fileMetrics.AfferentCoupling,
			"domain":                  topic.Domain,
			"efferent_coupling":       fileMetrics.EfferentCoupling,
			"has_circular_dependency": fileMetrics.HasCircularDependency,
			"has_smells":              len(fileMetrics.Smells) > 0,
			"incoming_relation_count": len(incomingRelations),
			"instability":             fileMetrics.Instability,
			"is_god_file":             fileMetrics.IsGodFile,
			"is_orphan_file":          fileMetrics.IsOrphanFile,
			"language":                string(file.Language),
			"outgoing_relation_count": len(outgoingRelations),
			"smells":                  append([]string(nil), fileMetrics.Smells...),
			"scraped":                 topic.Today,
			"source_kind":             "codebase-file",
			"source_path":             file.FilePath,
			"stage":                   "raw",
			"symbol_count":            len(fileSymbols),
			"tags":                    []string{topic.Domain, "raw", "codebase", "file", string(file.Language)},
			"title":                   fmt.Sprintf("Codebase File: %s", file.FilePath),
			"type":                    "source",
		},
		Body: strings.Join(sections, "\n"),
	}
}

func createSymbolFrontmatter(
	topic models.TopicMetadata,
	symbol models.SymbolNode,
	symbolMetrics models.SymbolMetrics,
	incomingRelations []models.RelationEdge,
	outgoingRelations []models.RelationEdge,
) map[string]interface{} {
	frontmatter := map[string]interface{}{
		"blast_radius":             symbolMetrics.BlastRadius,
		"centrality":               symbolMetrics.Centrality,
		"domain":                   topic.Domain,
		"end_line":                 symbol.EndLine,
		"exported":                 symbol.Exported,
		"external_reference_count": symbolMetrics.ExternalReferenceCount,
		"has_smells":               len(symbolMetrics.Smells) > 0,
		"incoming_relation_count":  len(incomingRelations),
		"is_dead_export":           symbolMetrics.IsDeadExport,
		"language":                 string(symbol.Language),
		"loc":                      symbolMetrics.LOC,
		"outgoing_relation_count":  len(outgoingRelations),
		"scraped":                  topic.Today,
		"smells":                   append([]string(nil), symbolMetrics.Smells...),
		"source_kind":              "codebase-symbol",
		"source_path":              symbol.FilePath,
		"stage":                    "raw",
		"start_line":               symbol.StartLine,
		"symbol_kind":              symbol.SymbolKind,
		"symbol_name":              symbol.Name,
		"tags": []string{
			topic.Domain,
			"raw",
			"codebase",
			"symbol",
			string(symbol.Language),
			symbol.SymbolKind,
		},
		"title": fmt.Sprintf("Codebase Symbol: %s", symbol.Name),
		"type":  "source",
	}

	if isFunctionLike(symbol.SymbolKind) {
		frontmatter["cyclomatic_complexity"] = maxInt(symbol.CyclomaticComplexity, 1)
		frontmatter["is_long_function"] = symbolMetrics.IsLongFunction
	}

	return frontmatter
}

func renderRawSymbolDocument(
	topic models.TopicMetadata,
	symbol models.SymbolNode,
	outgoingRelations []models.RelationEdge,
	incomingRelations []models.RelationEdge,
	symbolMetrics models.SymbolMetrics,
	documentLookup map[string]string,
	externalNodes map[string]models.ExternalNode,
) models.RenderedDocument {
	fileLink := toSourceWikiLink(topic, GetRawFileDocumentPath(symbol.FilePath), symbol.FilePath)

	signatureBlock := "None"
	if strings.TrimSpace(symbol.Signature) != "" {
		signatureBlock = strings.Join([]string{"```text", symbol.Signature, "```"}, "\n")
	}

	staticAnalysisLines := []string{
		fmt.Sprintf("- Blast radius: %d", symbolMetrics.BlastRadius),
		fmt.Sprintf("- External references: %d", symbolMetrics.ExternalReferenceCount),
		fmt.Sprintf("- Centrality: %s", strconv.FormatFloat(symbolMetrics.Centrality, 'f', -1, 64)),
		fmt.Sprintf("- LOC: %d", symbolMetrics.LOC),
		fmt.Sprintf("- Dead export: %t", symbolMetrics.IsDeadExport),
		fmt.Sprintf("- Smells: %s", renderSmellList(symbolMetrics.Smells)),
	}

	if isFunctionLike(symbol.SymbolKind) {
		prefix := []string{
			fmt.Sprintf("- Cyclomatic complexity: %d", maxInt(symbol.CyclomaticComplexity, 1)),
			fmt.Sprintf("- Long function: %t", symbolMetrics.IsLongFunction),
		}
		staticAnalysisLines = append(prefix, staticAnalysisLines...)
	}

	docComment := strings.TrimSpace(symbol.DocComment)
	if docComment == "" {
		docComment = "None"
	}

	sections := []string{
		fmt.Sprintf("# Codebase Symbol: %s", symbol.Name),
		"",
		fmt.Sprintf("Source file: %s", fileLink),
		"",
		"## Kind",
		"`" + symbol.SymbolKind + "`",
		"",
		"## Static Analysis",
	}
	sections = append(sections, staticAnalysisLines...)
	sections = append(sections,
		"",
		"## Signature",
		signatureBlock,
		"",
		"## Documentation",
		docComment,
		"",
		"## Outgoing Relations",
	)
	sections = append(sections, renderRelationList(topic, outgoingRelations, documentLookup, externalNodes)...)
	sections = append(sections,
		"",
		"## Backlinks",
	)
	sections = append(sections, renderBacklinkList(topic, incomingRelations, documentLookup, externalNodes)...)

	return models.RenderedDocument{
		Kind:         models.DocRaw,
		ManagedArea:  models.AreaRawCodebase,
		RelativePath: GetRawSymbolDocumentPath(symbol),
		Frontmatter:  createSymbolFrontmatter(topic, symbol, symbolMetrics, incomingRelations, outgoingRelations),
		Body:         strings.Join(sections, "\n"),
	}
}

func renderRawDirectoryIndex(
	topic models.TopicMetadata,
	directoryPath string,
	files []models.GraphFile,
	symbols []models.SymbolNode,
	directoryMetrics models.DirectoryMetrics,
) models.RenderedDocument {
	orderedFiles := append([]models.GraphFile(nil), files...)
	sort.Slice(orderedFiles, func(i, j int) bool {
		return orderedFiles[i].FilePath < orderedFiles[j].FilePath
	})

	fileLinks := []string{"None"}
	if len(orderedFiles) > 0 {
		fileLinks = make([]string, 0, len(orderedFiles))
		for _, file := range orderedFiles {
			fileLinks = append(fileLinks, "- "+toSourceWikiLink(topic, GetRawFileDocumentPath(file.FilePath), file.FilePath))
		}
	}

	orderedSymbols := append([]models.SymbolNode(nil), symbols...)
	sortSymbolsByLocation(orderedSymbols)

	symbolLinks := []string{"None"}
	if len(orderedSymbols) > 0 {
		symbolLinks = make([]string, 0, len(orderedSymbols))
		for _, symbol := range orderedSymbols {
			symbolLinks = append(symbolLinks, "- "+toSourceWikiLink(
				topic,
				GetRawSymbolDocumentPath(symbol),
				fmt.Sprintf("%s (%s)", symbol.Name, symbol.SymbolKind),
			))
		}
	}

	sections := []string{
		fmt.Sprintf("# Directory Snapshot: %s", directoryPath),
		"",
		"Generated directory-level inventory for the codebase corpus.",
		"",
		"## Static Analysis",
		fmt.Sprintf("- Afferent coupling: %d", directoryMetrics.AfferentCoupling),
		fmt.Sprintf("- Efferent coupling: %d", directoryMetrics.EfferentCoupling),
		fmt.Sprintf("- Instability: %s", strconv.FormatFloat(directoryMetrics.Instability, 'f', -1, 64)),
		"",
		"## Files",
	}
	sections = append(sections, fileLinks...)
	sections = append(sections,
		"",
		"## Symbols",
	)
	sections = append(sections, symbolLinks...)

	return models.RenderedDocument{
		Kind:         models.DocRaw,
		ManagedArea:  models.AreaRawCodebase,
		RelativePath: GetRawDirectoryIndexPath(directoryPath),
		Frontmatter: map[string]interface{}{
			"afferent_coupling": directoryMetrics.AfferentCoupling,
			"domain":            topic.Domain,
			"efferent_coupling": directoryMetrics.EfferentCoupling,
			"file_count":        len(files),
			"instability":       directoryMetrics.Instability,
			"scraped":           topic.Today,
			"source_kind":       "codebase-directory-index",
			"source_path":       directoryPath,
			"stage":             "raw",
			"symbol_count":      len(symbols),
			"tags":              []string{topic.Domain, "raw", "codebase", "directory-index"},
			"title":             fmt.Sprintf("Directory Snapshot: %s", directoryPath),
			"type":              "source",
		},
		Body: strings.Join(sections, "\n"),
	}
}

func renderRawLanguageIndex(
	topic models.TopicMetadata,
	language string,
	files []models.GraphFile,
	symbols []models.SymbolNode,
) models.RenderedDocument {
	orderedFiles := append([]models.GraphFile(nil), files...)
	sort.Slice(orderedFiles, func(i, j int) bool {
		return orderedFiles[i].FilePath < orderedFiles[j].FilePath
	})

	fileLinks := []string{"None"}
	if len(orderedFiles) > 0 {
		fileLinks = make([]string, 0, len(orderedFiles))
		for _, file := range orderedFiles {
			fileLinks = append(fileLinks, "- "+toSourceWikiLink(topic, GetRawFileDocumentPath(file.FilePath), file.FilePath))
		}
	}

	orderedSymbols := append([]models.SymbolNode(nil), symbols...)
	sortSymbolsByLocation(orderedSymbols)

	symbolLinks := []string{"None"}
	if len(orderedSymbols) > 0 {
		symbolLinks = make([]string, 0, len(orderedSymbols))
		for _, symbol := range orderedSymbols {
			symbolLinks = append(symbolLinks, "- "+toSourceWikiLink(
				topic,
				GetRawSymbolDocumentPath(symbol),
				fmt.Sprintf("%s (%s)", symbol.Name, symbol.SymbolKind),
			))
		}
	}

	sections := []string{
		fmt.Sprintf("# Language Snapshot: %s", language),
		"",
		"Generated language-level inventory for the codebase corpus.",
		"",
		"## Files",
	}
	sections = append(sections, fileLinks...)
	sections = append(sections,
		"",
		"## Symbols",
	)
	sections = append(sections, symbolLinks...)

	return models.RenderedDocument{
		Kind:         models.DocRaw,
		ManagedArea:  models.AreaRawCodebase,
		RelativePath: GetRawLanguageIndexPath(language),
		Frontmatter: map[string]interface{}{
			"domain":       topic.Domain,
			"file_count":   len(files),
			"language":     language,
			"scraped":      topic.Today,
			"source_kind":  "codebase-language-index",
			"stage":        "raw",
			"symbol_count": len(symbols),
			"tags":         []string{topic.Domain, "raw", "codebase", "language-index", language},
			"title":        fmt.Sprintf("Language Snapshot: %s", language),
			"type":         "source",
		},
		Body: strings.Join(sections, "\n"),
	}
}

func defaultFileMetrics(metric models.FileMetrics) models.FileMetrics {
	return metric
}

func defaultSymbolMetrics(symbol models.SymbolNode, metric models.SymbolMetrics) models.SymbolMetrics {
	if metric.LOC == 0 {
		metric.LOC = maxInt(symbol.EndLine-symbol.StartLine+1, 0)
	}
	return metric
}

func defaultDirectoryMetrics(metric models.DirectoryMetrics) models.DirectoryMetrics {
	return metric
}

func topRelationHotspotFiles(
	files []models.GraphFile,
	relationsBySource map[string][]models.RelationEdge,
	relationsByTarget map[string][]models.RelationEdge,
) []models.GraphFile {
	type hotspot struct {
		File          models.GraphFile
		RelationCount int
	}

	entries := make([]hotspot, 0, len(files))
	for _, file := range files {
		entries = append(entries, hotspot{
			File:          file,
			RelationCount: len(relationsBySource[file.ID]) + len(relationsByTarget[file.ID]),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].RelationCount != entries[j].RelationCount {
			return entries[i].RelationCount > entries[j].RelationCount
		}
		return entries[i].File.FilePath < entries[j].File.FilePath
	})

	if len(entries) > 5 {
		entries = entries[:5]
	}

	hotspots := make([]models.GraphFile, 0, len(entries))
	for _, entry := range entries {
		hotspots = append(hotspots, entry.File)
	}

	return hotspots
}

func groupRelationsBySource(relations []models.RelationEdge) map[string][]models.RelationEdge {
	grouped := make(map[string][]models.RelationEdge)
	for _, relation := range relations {
		grouped[relation.FromID] = append(grouped[relation.FromID], relation)
	}
	return grouped
}

func groupRelationsByTarget(relations []models.RelationEdge) map[string][]models.RelationEdge {
	grouped := make(map[string][]models.RelationEdge)
	for _, relation := range relations {
		grouped[relation.ToID] = append(grouped[relation.ToID], relation)
	}
	return grouped
}

func groupSymbolsByFile(symbols []models.SymbolNode) map[string][]models.SymbolNode {
	grouped := make(map[string][]models.SymbolNode)
	for _, symbol := range symbols {
		grouped[symbol.FilePath] = append(grouped[symbol.FilePath], symbol)
	}
	return grouped
}

func groupFilesByDirectory(files []models.GraphFile) map[string][]models.GraphFile {
	grouped := make(map[string][]models.GraphFile)
	for _, file := range files {
		grouped[path.Dir(file.FilePath)] = append(grouped[path.Dir(file.FilePath)], file)
	}
	return grouped
}

func groupSymbolsByDirectory(symbols []models.SymbolNode) map[string][]models.SymbolNode {
	grouped := make(map[string][]models.SymbolNode)
	for _, symbol := range symbols {
		grouped[path.Dir(symbol.FilePath)] = append(grouped[path.Dir(symbol.FilePath)], symbol)
	}
	return grouped
}

func groupFilesByLanguage(files []models.GraphFile) map[string][]models.GraphFile {
	grouped := make(map[string][]models.GraphFile)
	for _, file := range files {
		grouped[string(file.Language)] = append(grouped[string(file.Language)], file)
	}
	return grouped
}

func groupSymbolsByLanguage(symbols []models.SymbolNode) map[string][]models.SymbolNode {
	grouped := make(map[string][]models.SymbolNode)
	for _, symbol := range symbols {
		grouped[string(symbol.Language)] = append(grouped[string(symbol.Language)], symbol)
	}
	return grouped
}

func groupSymbolsByKind(symbols []models.SymbolNode) map[string][]models.SymbolNode {
	grouped := make(map[string][]models.SymbolNode)
	for _, symbol := range symbols {
		grouped[symbol.SymbolKind] = append(grouped[symbol.SymbolKind], symbol)
	}
	return grouped
}

func sortSymbolsByLocation(symbols []models.SymbolNode) {
	sort.Slice(symbols, func(i, j int) bool {
		left := symbols[i]
		right := symbols[j]

		if left.FilePath != right.FilePath {
			return left.FilePath < right.FilePath
		}
		if left.StartLine != right.StartLine {
			return left.StartLine < right.StartLine
		}
		return left.Name < right.Name
	})
}

func sortedMapKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func isFunctionLike(symbolKind string) bool {
	return symbolKind == "function" || symbolKind == "method"
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}
