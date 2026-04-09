package vault

import (
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/user/go-devstack/internal/models"
)

type smellEntry struct {
	FilePath     string
	Label        string
	RelativePath string
}

func buildStarterWikiArticles(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
	filesByLanguage map[string][]models.GraphFile,
	filesByDirectory map[string][]models.GraphFile,
	symbolsByDirectory map[string][]models.SymbolNode,
	symbolsByKind map[string][]models.SymbolNode,
	hotspotFiles []models.GraphFile,
) []starterWikiArticle {
	return []starterWikiArticle{
		createCodebaseOverviewArticle(topic, graph, metrics, filesByLanguage, filesByDirectory, hotspotFiles),
		createDirectoryMapArticle(topic, filesByDirectory, symbolsByDirectory, metrics),
		createSymbolTaxonomyArticle(topic, graph, symbolsByKind, filesByLanguage),
		createDependencyHotspotsArticle(topic, graph),
		createComplexityHotspotsArticle(topic, graph),
		createModuleHealthArticle(topic, graph, metrics),
		createDeadCodeReportArticle(topic, graph, metrics),
		createCodeSmellsArticle(topic, graph, metrics),
		createCircularDependenciesArticle(topic, metrics),
		createHighImpactSymbolsArticle(topic, graph, metrics),
	}
}

func makeWikiFrontmatter(
	topic models.TopicMetadata,
	article starterWikiArticle,
) map[string]interface{} {
	sources := make([]string, 0, len(article.Sources))
	for _, source := range article.Sources {
		sources = append(sources, toSourceWikiLink(topic, source, ""))
	}

	return map[string]interface{}{
		"created":   topic.Today,
		"domain":    topic.Domain,
		"generator": "kodebase",
		"sources":   sources,
		"stage":     "compiled",
		"tags":      []string{topic.Domain, "wiki", "codebase", "starter"},
		"title":     article.Title,
		"type":      "wiki",
		"updated":   topic.Today,
	}
}

func renderWikiArticle(topic models.TopicMetadata, article starterWikiArticle) models.RenderedDocument {
	return models.RenderedDocument{
		Kind:         models.DocWiki,
		ManagedArea:  models.AreaWikiConcept,
		RelativePath: GetWikiConceptPath(article.Title),
		Frontmatter:  makeWikiFrontmatter(topic, article),
		Body:         article.Body,
	}
}

func renderDashboard(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	articles []starterWikiArticle,
) models.RenderedDocument {
	directoryPaths := make(map[string]struct{}, len(graph.Files))
	languages := make(map[string]struct{}, len(graph.Files))

	for _, file := range graph.Files {
		directoryPaths[path.Dir(file.FilePath)] = struct{}{}
		languages[string(file.Language)] = struct{}{}
	}

	rawSourceCount := len(graph.Files) + len(graph.Symbols) + len(directoryPaths) + len(languages)

	lines := []string{
		fmt.Sprintf("# %s - Dashboard", topic.Title),
		"",
		"Landing page for the generated Karpathy-compatible codebase topic.",
		"",
		"## At a glance",
		fmt.Sprintf("- **Articles:** %d", len(articles)),
		fmt.Sprintf("- **Raw sources:** %d", rawSourceCount),
		fmt.Sprintf("- **Parsed files:** %d", len(graph.Files)),
		fmt.Sprintf("- **Parsed symbols:** %d", len(graph.Symbols)),
		fmt.Sprintf("- **Last updated:** %s", topic.Today),
		"",
		"## Starter articles",
	}

	for _, article := range articles {
		lines = append(lines, fmt.Sprintf(
			"- %s - %s",
			ToTopicWikiLink(topic.Slug, GetWikiConceptPath(article.Title), article.Title),
			article.Summary,
		))
	}

	lines = append(lines,
		"",
		"## Navigation",
		"- "+ToTopicWikiLink(topic.Slug, GetWikiIndexPath("Concept Index"), "Concept Index"),
		"- "+ToTopicWikiLink(topic.Slug, GetWikiIndexPath("Source Index"), "Source Index"),
	)

	return models.RenderedDocument{
		Kind:         models.DocIndex,
		ManagedArea:  models.AreaWikiIndex,
		RelativePath: GetWikiIndexPath("Dashboard"),
		Frontmatter: map[string]interface{}{
			"domain":  topic.Domain,
			"title":   "Dashboard",
			"type":    "index",
			"updated": topic.Today,
		},
		Body: strings.Join(lines, "\n"),
	}
}

func renderConceptIndex(
	topic models.TopicMetadata,
	articles []starterWikiArticle,
) models.RenderedDocument {
	orderedArticles := append([]starterWikiArticle(nil), articles...)
	sort.Slice(orderedArticles, func(i, j int) bool {
		return orderedArticles[i].Title < orderedArticles[j].Title
	})

	rows := make([]string, 0, len(orderedArticles))
	for _, article := range orderedArticles {
		rows = append(rows, fmt.Sprintf(
			"| %s | %s |",
			ToTopicWikiLink(topic.Slug, GetWikiConceptPath(article.Title), article.Title),
			article.Summary,
		))
	}

	if len(rows) == 0 {
		rows = append(rows, "| None | No articles yet. |")
	}

	return models.RenderedDocument{
		Kind:         models.DocIndex,
		ManagedArea:  models.AreaWikiIndex,
		RelativePath: GetWikiIndexPath("Concept Index"),
		Frontmatter: map[string]interface{}{
			"domain":  topic.Domain,
			"title":   "Concept Index",
			"type":    "index",
			"updated": topic.Today,
		},
		Body: strings.Join([]string{
			fmt.Sprintf("# %s - Concept Index", topic.Title),
			"",
			"Alphabetical listing of every generated wiki article in this topic.",
			"",
			"| Article | Summary |",
			"| ------- | ------- |",
			strings.Join(rows, "\n"),
		}, "\n"),
	}
}

func renderSourceIndex(topic models.TopicMetadata, articles []starterWikiArticle) models.RenderedDocument {
	citedBySource := make(map[string][]string)

	for _, article := range articles {
		for _, source := range article.Sources {
			citedBySource[source] = append(citedBySource[source], article.Title)
		}
	}

	rows := make([]string, 0, len(citedBySource))
	for _, source := range sortedMapKeys(citedBySource) {
		titles := append([]string(nil), citedBySource[source]...)
		sort.Strings(titles)

		links := make([]string, 0, len(titles))
		for _, title := range titles {
			links = append(links, ToTopicWikiLink(topic.Slug, GetWikiConceptPath(title), title))
		}

		rows = append(rows, fmt.Sprintf(
			"| %s | %s |",
			toSourceWikiLink(topic, source, ""),
			strings.Join(links, ", "),
		))
	}

	if len(rows) == 0 {
		rows = append(rows, "| None | None |")
	}

	return models.RenderedDocument{
		Kind:         models.DocIndex,
		ManagedArea:  models.AreaWikiIndex,
		RelativePath: GetWikiIndexPath("Source Index"),
		Frontmatter: map[string]interface{}{
			"domain":  topic.Domain,
			"title":   "Source Index",
			"type":    "index",
			"updated": topic.Today,
		},
		Body: strings.Join([]string{
			fmt.Sprintf("# %s - Source Index", topic.Title),
			"",
			"Raw codebase snapshots currently cited by the starter wiki.",
			"",
			"| Source | Cited by |",
			"| ------ | -------- |",
			strings.Join(rows, "\n"),
		}, "\n"),
	}
}

func createCodebaseOverviewArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
	filesByLanguage map[string][]models.GraphFile,
	filesByDirectory map[string][]models.GraphFile,
	hotspotFiles []models.GraphFile,
) starterWikiArticle {
	languageRows := make([]string, 0, len(filesByLanguage))
	for _, language := range sortedMapKeys(filesByLanguage) {
		files := filesByLanguage[language]
		symbolCount := 0
		for _, symbol := range graph.Symbols {
			if string(symbol.Language) == language {
				symbolCount++
			}
		}

		languageRows = append(languageRows, fmt.Sprintf(
			"| %s | %d | %d |",
			toSourceWikiLink(topic, GetRawLanguageIndexPath(language), language),
			len(files),
			symbolCount,
		))
	}
	if len(languageRows) == 0 {
		languageRows = append(languageRows, "| None | 0 | 0 |")
	}

	directoryRows := []string{}
	directoryKeys := sortedMapKeys(filesByDirectory)
	if len(directoryKeys) > 8 {
		directoryKeys = directoryKeys[:8]
	}

	for _, directoryPath := range directoryKeys {
		files := filesByDirectory[directoryPath]
		directoryMetric := metrics.Directories[directoryPath]
		directoryRows = append(directoryRows, fmt.Sprintf(
			"- %s · %d files · instability=%s",
			toSourceWikiLink(topic, GetRawDirectoryIndexPath(directoryPath), directoryPath),
			len(files),
			strconv.FormatFloat(directoryMetric.Instability, 'f', -1, 64),
		))
	}

	hotspotLines := []string{"- No relation hotspots were detected."}
	if len(hotspotFiles) > 0 {
		hotspotLines = make([]string, 0, len(hotspotFiles))
		for _, file := range hotspotFiles {
			hotspotLines = append(hotspotLines, "- "+toSourceWikiLink(topic, GetRawFileDocumentPath(file.FilePath), file.FilePath))
		}
	}

	sources := make([]string, 0, len(filesByLanguage)+len(filesByDirectory)+len(hotspotFiles))
	for _, language := range sortedMapKeys(filesByLanguage) {
		sources = append(sources, GetRawLanguageIndexPath(language))
	}
	for _, directoryPath := range directoryKeys {
		sources = append(sources, GetRawDirectoryIndexPath(directoryPath))
	}
	for _, hotspot := range hotspotFiles {
		sources = append(sources, GetRawFileDocumentPath(hotspot.FilePath))
	}
	sources = uniqueStrings(sources)

	sourceLines := make([]string, 0, len(sources))
	for _, source := range sources {
		sourceLines = append(sourceLines, "- "+toSourceWikiLink(topic, source, ""))
	}

	bodyLines := []string{
		"# Codebase Overview",
		"",
		fmt.Sprintf(
			"%s currently compiles into a Karpathy-style topic where the codebase itself is staged in `raw/codebase/` and this starter wiki provides the first synthesized navigation layer. The corpus contains %d parsed source files, %d symbols, and %d extracted relations.",
			topic.Title,
			len(graph.Files),
			len(graph.Symbols),
			len(graph.Relations),
		),
		"",
		fmt.Sprintf(
			"Start with %s for coupling, %s for function-level complexity, and %s for likely cleanup candidates.",
			ToTopicWikiLink(topic.Slug, GetWikiConceptPath("Module Health"), "Module Health"),
			ToTopicWikiLink(topic.Slug, GetWikiConceptPath("Complexity Hotspots"), "Complexity Hotspots"),
			ToTopicWikiLink(topic.Slug, GetWikiConceptPath("Dead Code Report"), "Dead Code Report"),
		),
		"",
		"## Language Coverage",
		"",
		"| Language | Files | Symbols |",
		"| -------- | ----- | ------- |",
	}
	bodyLines = append(bodyLines, languageRows...)
	bodyLines = append(bodyLines,
		"",
		"## Directory Coverage",
	)
	bodyLines = append(bodyLines, directoryRows...)
	bodyLines = append(bodyLines,
		"",
		"## Relation Hotspots",
	)
	bodyLines = append(bodyLines, hotspotLines...)
	bodyLines = append(bodyLines,
		"",
		"## Sources and Further Reading",
	)
	bodyLines = append(bodyLines, sourceLines...)

	return starterWikiArticle{
		Title:   "Codebase Overview",
		Summary: "High-level inventory of repository shape, language coverage, and the main analysis entry points for complexity, coupling, and dead code.",
		Sources: sources,
		Body:    strings.Join(bodyLines, "\n"),
	}
}

func createDirectoryMapArticle(
	topic models.TopicMetadata,
	filesByDirectory map[string][]models.GraphFile,
	symbolsByDirectory map[string][]models.SymbolNode,
	metrics models.MetricsResult,
) starterWikiArticle {
	rows := make([]string, 0, len(filesByDirectory))
	sources := make([]string, 0, len(filesByDirectory))

	for _, directoryPath := range sortedMapKeys(filesByDirectory) {
		files := filesByDirectory[directoryPath]
		symbols := symbolsByDirectory[directoryPath]
		directoryMetric := metrics.Directories[directoryPath]
		rows = append(rows, fmt.Sprintf(
			"| %s | %d | %d | %s |",
			toSourceWikiLink(topic, GetRawDirectoryIndexPath(directoryPath), directoryPath),
			len(files),
			len(symbols),
			strconv.FormatFloat(directoryMetric.Instability, 'f', -1, 64),
		))
		sources = append(sources, GetRawDirectoryIndexPath(directoryPath))
	}

	if len(rows) == 0 {
		rows = append(rows, "| None | 0 | 0 | 0 |")
	}

	return starterWikiArticle{
		Title:   "Directory Map",
		Summary: "Directory-by-directory map of the codebase corpus, including file/symbol counts and directory-level instability rollups.",
		Sources: sources,
		Body: strings.Join([]string{
			"# Directory Map",
			"",
			"The directory layer is the fastest way to orient the codebase corpus before drilling into file or symbol snapshots. Each row below links directly to the raw staged view in `raw/codebase/indexes/directories/`.",
			"",
			"## Directory Inventory",
			"",
			"| Directory | Files | Symbols | Instability |",
			"| --------- | ----- | ------- | ----------- |",
			strings.Join(rows, "\n"),
			"",
			"## How To Use This Map",
			"",
			fmt.Sprintf(
				"Cross-check unstable directories against %s and hotspots against %s.",
				ToTopicWikiLink(topic.Slug, GetWikiConceptPath("Module Health"), "Module Health"),
				ToTopicWikiLink(topic.Slug, GetWikiConceptPath("Dependency Hotspots"), "Dependency Hotspots"),
			),
			"",
			"## Sources and Further Reading",
			renderSourceBulletList(topic, sources, "- No directories were extracted."),
		}, "\n"),
	}
}

func createSymbolTaxonomyArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	symbolsByKind map[string][]models.SymbolNode,
	filesByLanguage map[string][]models.GraphFile,
) starterWikiArticle {
	kindRows := make([]string, 0, len(symbolsByKind))
	for _, kind := range sortedMapKeys(symbolsByKind) {
		symbols := append([]models.SymbolNode(nil), symbolsByKind[kind]...)
		sortSymbolsByLocation(symbols)

		exampleLink := "None"
		if len(symbols) > 0 {
			exampleLink = toSourceWikiLink(topic, GetRawSymbolDocumentPath(symbols[0]), symbols[0].Name)
		}

		kindRows = append(kindRows, fmt.Sprintf("| `%s` | %d | %s |", kind, len(symbols), exampleLink))
	}
	if len(kindRows) == 0 {
		kindRows = append(kindRows, "| None | 0 | None |")
	}

	topSymbols := append([]models.SymbolNode(nil), graph.Symbols...)
	sortSymbolsByLocation(topSymbols)
	if len(topSymbols) > 10 {
		topSymbols = topSymbols[:10]
	}

	sources := make([]string, 0, len(filesByLanguage)+len(topSymbols))
	for _, language := range sortedMapKeys(filesByLanguage) {
		sources = append(sources, GetRawLanguageIndexPath(language))
	}
	for _, symbol := range topSymbols {
		sources = append(sources, GetRawSymbolDocumentPath(symbol))
	}
	sources = uniqueStrings(sources)

	return starterWikiArticle{
		Title:   "Symbol Taxonomy",
		Summary: "Inventory of top-level symbol kinds extracted from the repository, with representative raw snapshots for each category.",
		Sources: sources,
		Body: strings.Join([]string{
			"# Symbol Taxonomy",
			"",
			fmt.Sprintf("This topic currently tracks %d symbols across the supported languages. The goal of this starter taxonomy is to expose the shape of the codebase corpus before deeper article-by-article compilation begins.", len(graph.Symbols)),
			"",
			"## Symbol Kinds",
			"",
			"| Symbol kind | Count | Example |",
			"| ----------- | ----- | ------- |",
			strings.Join(kindRows, "\n"),
			"",
			"## Cross-References",
			"",
			fmt.Sprintf(
				"Use %s to find bottlenecks and %s to locate function-level smells.",
				ToTopicWikiLink(topic.Slug, GetWikiConceptPath("High-Impact Symbols"), "High-Impact Symbols"),
				ToTopicWikiLink(topic.Slug, GetWikiConceptPath("Code Smells"), "Code Smells"),
			),
			"",
			"## Sources and Further Reading",
			renderSourceBulletList(topic, sources, "- No symbol taxonomy sources were extracted."),
		}, "\n"),
	}
}

func createDependencyHotspotsArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
) starterWikiArticle {
	relationCounts := make(map[string]int)
	for _, relation := range graph.Relations {
		relationCounts[relation.FromID]++
		relationCounts[relation.ToID]++
	}

	type hotspot struct {
		File          models.GraphFile
		RelationCount int
	}

	hotspots := make([]hotspot, 0, len(graph.Files))
	for _, file := range graph.Files {
		hotspots = append(hotspots, hotspot{
			File:          file,
			RelationCount: relationCounts[file.ID],
		})
	}

	sort.Slice(hotspots, func(i, j int) bool {
		if hotspots[i].RelationCount != hotspots[j].RelationCount {
			return hotspots[i].RelationCount > hotspots[j].RelationCount
		}
		return hotspots[i].File.FilePath < hotspots[j].File.FilePath
	})
	if len(hotspots) > 10 {
		hotspots = hotspots[:10]
	}

	rows := []string{"| None | 0 |"}
	if len(hotspots) > 0 {
		rows = make([]string, 0, len(hotspots))
	}

	sources := make([]string, 0, len(hotspots))
	for _, entry := range hotspots {
		rows = append(rows, fmt.Sprintf(
			"| %s | %d |",
			toSourceWikiLink(topic, GetRawFileDocumentPath(entry.File.FilePath), entry.File.FilePath),
			entry.RelationCount,
		))
		sources = append(sources, GetRawFileDocumentPath(entry.File.FilePath))
	}

	return starterWikiArticle{
		Title:   "Dependency Hotspots",
		Summary: "Files with the highest relation density, useful as initial review targets and anchors for later compiled wiki articles.",
		Sources: sources,
		Body: strings.Join([]string{
			"# Dependency Hotspots",
			"",
			"Hotspots are file-level raw snapshots with the highest observed relation density in the normalized graph. They are the best starting points when deciding what to compile into deeper conceptual articles next.",
			"",
			"## Ranked Hotspots",
			"",
			"| File | Relation count |",
			"| ---- | -------------- |",
			strings.Join(rows, "\n"),
			"",
			"## Interpretation",
			"",
			fmt.Sprintf(
				"A high relation count usually indicates a coordination layer, a shared utility, or an entry point. Cross-check these files against %s to distinguish stable modules from unstable ones.",
				ToTopicWikiLink(topic.Slug, GetWikiConceptPath("Module Health"), "Module Health"),
			),
			"",
			"## Sources and Further Reading",
			renderSourceBulletList(topic, sources, "- No relation hotspots were extracted."),
		}, "\n"),
	}
}

func createComplexityHotspotsArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
) starterWikiArticle {
	functions := make([]models.SymbolNode, 0, len(graph.Symbols))
	for _, symbol := range graph.Symbols {
		if !isFunctionLike(symbol.SymbolKind) || symbol.CyclomaticComplexity == 0 {
			continue
		}
		functions = append(functions, symbol)
	}

	sort.Slice(functions, func(i, j int) bool {
		left := functions[i]
		right := functions[j]
		leftLOC := left.EndLine - left.StartLine + 1
		rightLOC := right.EndLine - right.StartLine + 1

		if left.CyclomaticComplexity != right.CyclomaticComplexity {
			return left.CyclomaticComplexity > right.CyclomaticComplexity
		}
		if leftLOC != rightLOC {
			return leftLOC > rightLOC
		}
		if left.FilePath != right.FilePath {
			return left.FilePath < right.FilePath
		}
		return left.Name < right.Name
	})
	if len(functions) > 20 {
		functions = functions[:20]
	}

	rows := []string{"| None | 0 | 0 | None |"}
	if len(functions) > 0 {
		rows = make([]string, 0, len(functions))
	}
	sources := make([]string, 0, len(functions)*2)
	for _, symbol := range functions {
		loc := maxInt(symbol.EndLine-symbol.StartLine+1, 0)
		rows = append(rows, fmt.Sprintf(
			"| %s | %d | %d | %s |",
			toSourceWikiLink(topic, GetRawSymbolDocumentPath(symbol), symbol.Name),
			maxInt(symbol.CyclomaticComplexity, 1),
			loc,
			toSourceWikiLink(topic, GetRawFileDocumentPath(symbol.FilePath), symbol.FilePath),
		))
		sources = append(sources, GetRawSymbolDocumentPath(symbol), GetRawFileDocumentPath(symbol.FilePath))
	}

	return starterWikiArticle{
		Title:   "Complexity Hotspots",
		Summary: "Top functions and methods ranked by cyclomatic complexity, useful for review, testing, and refactoring prioritization.",
		Sources: uniqueStrings(sources),
		Body: strings.Join([]string{
			"# Complexity Hotspots",
			"",
			"These functions have the highest measured cyclomatic complexity in the current codebase snapshot.",
			"",
			"## Top Functions",
			"",
			"| Symbol | Complexity | LOC | File |",
			"| ------ | ---------- | --- | ---- |",
			strings.Join(rows, "\n"),
			"",
			"## Cross-References",
			"",
			fmt.Sprintf(
				"Compare these hotspots against %s to distinguish locally complex functions from high-blast-radius functions.",
				ToTopicWikiLink(topic.Slug, GetWikiConceptPath("High-Impact Symbols"), "High-Impact Symbols"),
			),
		}, "\n"),
	}
}

func createDeadCodeReportArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
) starterWikiArticle {
	deadExportsByDirectory := make(map[string][]string)
	orphanFilesByDirectory := make(map[string][]string)

	for _, symbol := range graph.Symbols {
		metric, exists := metrics.Symbols[symbol.ID]
		if !exists || !metric.IsDeadExport {
			continue
		}

		directoryPath := path.Dir(symbol.FilePath)
		deadExportsByDirectory[directoryPath] = append(deadExportsByDirectory[directoryPath],
			"- "+toSourceWikiLink(topic, GetRawSymbolDocumentPath(symbol), symbol.Name),
		)
	}

	for _, file := range graph.Files {
		metric, exists := metrics.Files[file.ID]
		if !exists || !metric.IsOrphanFile {
			continue
		}

		directoryPath := path.Dir(file.FilePath)
		orphanFilesByDirectory[directoryPath] = append(orphanFilesByDirectory[directoryPath],
			"- "+toSourceWikiLink(topic, GetRawFileDocumentPath(file.FilePath), file.FilePath),
		)
	}

	for directoryPath, links := range deadExportsByDirectory {
		deadExportsByDirectory[directoryPath] = sortStrings(links)
	}
	for directoryPath, links := range orphanFilesByDirectory {
		orphanFilesByDirectory[directoryPath] = sortStrings(links)
	}

	sources := []string{}
	for _, symbol := range graph.Symbols {
		if metric, exists := metrics.Symbols[symbol.ID]; exists && metric.IsDeadExport {
			sources = append(sources, GetRawSymbolDocumentPath(symbol))
		}
	}
	for _, file := range graph.Files {
		if metric, exists := metrics.Files[file.ID]; exists && metric.IsOrphanFile {
			sources = append(sources, GetRawFileDocumentPath(file.FilePath))
		}
	}

	bodyLines := []string{
		"# Dead Code Report",
		"",
		"This report groups likely dead code candidates by directory for easier cleanup review.",
		"",
		"## Dead Exports",
	}
	bodyLines = append(bodyLines, renderGroupedLinks(deadExportsByDirectory, "No dead exports detected.")...)
	bodyLines = append(bodyLines, "## Orphan Files")
	bodyLines = append(bodyLines, renderGroupedLinks(orphanFilesByDirectory, "No orphan files detected.")...)

	return starterWikiArticle{
		Title:   "Dead Code Report",
		Summary: "Dead exports and orphan files grouped by directory, highlighting likely cleanup opportunities in the current snapshot.",
		Sources: uniqueStrings(sources),
		Body:    strings.Join(bodyLines, "\n"),
	}
}

func createModuleHealthArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
) starterWikiArticle {
	type fileEntry struct {
		File   models.GraphFile
		Metric models.FileMetrics
	}

	fileEntries := make([]fileEntry, 0, len(graph.Files))
	for _, file := range graph.Files {
		fileEntries = append(fileEntries, fileEntry{
			File:   file,
			Metric: metrics.Files[file.ID],
		})
	}

	sort.Slice(fileEntries, func(i, j int) bool {
		left := fileEntries[i]
		right := fileEntries[j]

		if left.Metric.Instability != right.Metric.Instability {
			return left.Metric.Instability > right.Metric.Instability
		}
		if left.Metric.EfferentCoupling != right.Metric.EfferentCoupling {
			return left.Metric.EfferentCoupling > right.Metric.EfferentCoupling
		}
		return left.File.FilePath < right.File.FilePath
	})

	fileRows := make([]string, 0, len(fileEntries))
	for _, entry := range fileEntries {
		fileRows = append(fileRows, fmt.Sprintf(
			"| %s | %d | %d | %s |",
			toSourceWikiLink(topic, GetRawFileDocumentPath(entry.File.FilePath), entry.File.FilePath),
			entry.Metric.AfferentCoupling,
			entry.Metric.EfferentCoupling,
			strconv.FormatFloat(entry.Metric.Instability, 'f', -1, 64),
		))
	}
	if len(fileRows) == 0 {
		fileRows = append(fileRows, "| None | 0 | 0 | 0 |")
	}

	directoryRows := make([]string, 0, len(metrics.Directories))
	for _, directoryPath := range sortedMapKeys(metrics.Directories) {
		directoryMetric := metrics.Directories[directoryPath]
		directoryRows = append(directoryRows, fmt.Sprintf(
			"| %s | %d | %d | %s |",
			toSourceWikiLink(topic, GetRawDirectoryIndexPath(directoryPath), directoryPath),
			directoryMetric.AfferentCoupling,
			directoryMetric.EfferentCoupling,
			strconv.FormatFloat(directoryMetric.Instability, 'f', -1, 64),
		))
	}
	if len(directoryRows) == 0 {
		directoryRows = append(directoryRows, "| None | 0 | 0 | 0 |")
	}

	sources := make([]string, 0, len(graph.Files)+len(metrics.Directories))
	for _, file := range graph.Files {
		sources = append(sources, GetRawFileDocumentPath(file.FilePath))
	}
	for _, directoryPath := range sortedMapKeys(metrics.Directories) {
		sources = append(sources, GetRawDirectoryIndexPath(directoryPath))
	}

	return starterWikiArticle{
		Title:   "Module Health",
		Summary: "File-level and directory-level coupling metrics ranked by instability to highlight volatile modules.",
		Sources: uniqueStrings(sources),
		Body: strings.Join([]string{
			"# Module Health",
			"",
			"Instability highlights modules with high outgoing dependence relative to incoming dependence.",
			"",
			"## Files",
			"",
			"| File | Ca | Ce | Instability |",
			"| ---- | -- | -- | ----------- |",
			strings.Join(fileRows, "\n"),
			"",
			"## Directories",
			"",
			"| Directory | Ca | Ce | Instability |",
			"| --------- | -- | -- | ----------- |",
			strings.Join(directoryRows, "\n"),
		}, "\n"),
	}
}

func createCodeSmellsArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
) starterWikiArticle {
	smellEntries := make(map[string][]smellEntry)

	for _, symbol := range graph.Symbols {
		metric, exists := metrics.Symbols[symbol.ID]
		if !exists || len(metric.Smells) == 0 {
			continue
		}

		for _, smell := range metric.Smells {
			smellEntries[smell] = append(smellEntries[smell], smellEntry{
				FilePath:     symbol.FilePath,
				Label:        symbol.Name,
				RelativePath: GetRawSymbolDocumentPath(symbol),
			})
		}
	}

	for _, file := range graph.Files {
		metric, exists := metrics.Files[file.ID]
		if !exists || len(metric.Smells) == 0 {
			continue
		}

		for _, smell := range metric.Smells {
			smellEntries[smell] = append(smellEntries[smell], smellEntry{
				FilePath:     file.FilePath,
				Label:        file.FilePath,
				RelativePath: GetRawFileDocumentPath(file.FilePath),
			})
		}
	}

	sections := []string{"# Code Smells", ""}
	for _, smell := range sortedMapKeys(smellEntries) {
		entries := append([]smellEntry(nil), smellEntries[smell]...)
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].FilePath != entries[j].FilePath {
				return entries[i].FilePath < entries[j].FilePath
			}
			return entries[i].Label < entries[j].Label
		})

		sections = append(sections, "## "+smell)
		for _, entry := range entries {
			sections = append(sections, fmt.Sprintf(
				"- %s · %s",
				toSourceWikiLink(topic, entry.RelativePath, entry.Label),
				entry.FilePath,
			))
		}
		sections = append(sections, "")
	}
	if len(smellEntries) == 0 {
		sections = append(sections, "No smells detected.")
	}

	sources := []string{}
	for _, symbol := range graph.Symbols {
		if metric, exists := metrics.Symbols[symbol.ID]; exists && len(metric.Smells) > 0 {
			sources = append(sources, GetRawSymbolDocumentPath(symbol))
		}
	}
	for _, file := range graph.Files {
		if metric, exists := metrics.Files[file.ID]; exists && len(metric.Smells) > 0 {
			sources = append(sources, GetRawFileDocumentPath(file.FilePath))
		}
	}

	return starterWikiArticle{
		Title:   "Code Smells",
		Summary: "All file-level and symbol-level smells grouped by smell type for triage and refactoring prioritization.",
		Sources: uniqueStrings(sources),
		Body:    strings.Join(sections, "\n"),
	}
}

func createCircularDependenciesArticle(
	topic models.TopicMetadata,
	metrics models.MetricsResult,
) starterWikiArticle {
	sources := make([]string, 0)
	cycles := []string{"No circular dependencies detected."}

	if len(metrics.CircularDependencies) > 0 {
		cycles = make([]string, 0, len(metrics.CircularDependencies))
		for index, cycle := range metrics.CircularDependencies {
			links := make([]string, 0, len(cycle))
			for _, filePath := range cycle {
				sources = append(sources, GetRawFileDocumentPath(filePath))
				links = append(links, toSourceWikiLink(topic, GetRawFileDocumentPath(filePath), filePath))
			}
			cycles = append(cycles, fmt.Sprintf("%d. %s", index+1, strings.Join(links, " -> ")))
		}
	}

	return starterWikiArticle{
		Title:   "Circular Dependencies",
		Summary: "Detected file-level import cycles, listed as ordered loops with links back to the raw file snapshots.",
		Sources: uniqueStrings(sources),
		Body:    strings.Join(append([]string{"# Circular Dependencies", ""}, cycles...), "\n"),
	}
}

func createHighImpactSymbolsArticle(
	topic models.TopicMetadata,
	graph models.GraphSnapshot,
	metrics models.MetricsResult,
) starterWikiArticle {
	type symbolEntry struct {
		Symbol models.SymbolNode
		Metric models.SymbolMetrics
	}

	symbolEntries := make([]symbolEntry, 0, len(graph.Symbols))
	for _, symbol := range graph.Symbols {
		symbolEntries = append(symbolEntries, symbolEntry{
			Symbol: symbol,
			Metric: metrics.Symbols[symbol.ID],
		})
	}

	sort.Slice(symbolEntries, func(i, j int) bool {
		left := symbolEntries[i]
		right := symbolEntries[j]

		if left.Metric.BlastRadius != right.Metric.BlastRadius {
			return left.Metric.BlastRadius > right.Metric.BlastRadius
		}
		if left.Metric.DirectDependents != right.Metric.DirectDependents {
			return left.Metric.DirectDependents > right.Metric.DirectDependents
		}
		if left.Symbol.FilePath != right.Symbol.FilePath {
			return left.Symbol.FilePath < right.Symbol.FilePath
		}
		return left.Symbol.Name < right.Symbol.Name
	})
	if len(symbolEntries) > 20 {
		symbolEntries = symbolEntries[:20]
	}

	rows := []string{"| None | 0 | 0 | None |"}
	if len(symbolEntries) > 0 {
		rows = make([]string, 0, len(symbolEntries))
	}
	sources := make([]string, 0, len(symbolEntries)*2)

	for _, entry := range symbolEntries {
		rows = append(rows, fmt.Sprintf(
			"| %s | %d | %d | %s |",
			toSourceWikiLink(topic, GetRawSymbolDocumentPath(entry.Symbol), entry.Symbol.Name),
			entry.Metric.BlastRadius,
			entry.Metric.DirectDependents,
			toSourceWikiLink(topic, GetRawFileDocumentPath(entry.Symbol.FilePath), entry.Symbol.FilePath),
		))
		sources = append(sources, GetRawSymbolDocumentPath(entry.Symbol), GetRawFileDocumentPath(entry.Symbol.FilePath))
	}

	return starterWikiArticle{
		Title:   "High-Impact Symbols",
		Summary: "Top symbols by blast radius, highlighting declarations with the largest downstream dependency surface.",
		Sources: uniqueStrings(sources),
		Body: strings.Join([]string{
			"# High-Impact Symbols",
			"",
			"Blast radius counts both direct and transitive dependents reachable through call and reference chains.",
			"",
			"## Top Symbols",
			"",
			"| Symbol | Blast Radius | Direct Dependents | File |",
			"| ------ | ------------ | ----------------- | ---- |",
			strings.Join(rows, "\n"),
		}, "\n"),
	}
}

func renderGroupedLinks(groups map[string][]string, emptyMessage string) []string {
	if len(groups) == 0 {
		return []string{emptyMessage}
	}

	lines := make([]string, 0)
	for _, directoryPath := range sortedMapKeys(groups) {
		lines = append(lines, "### "+directoryPath)
		lines = append(lines, groups[directoryPath]...)
		lines = append(lines, "")
	}
	return lines
}

func renderSourceBulletList(topic models.TopicMetadata, sources []string, emptyMessage string) string {
	if len(sources) == 0 {
		return emptyMessage
	}

	lines := make([]string, 0, len(sources))
	for _, source := range uniqueStrings(sources) {
		lines = append(lines, "- "+toSourceWikiLink(topic, source, ""))
	}

	return strings.Join(lines, "\n")
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func sortStrings(values []string) []string {
	ordered := append([]string(nil), values...)
	sort.Strings(ordered)
	return ordered
}
