---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 2
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0.6667
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 28
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
symbol_count: 22
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/render_wiki.go"
type: "source"
---

# Codebase File: internal/vault/render_wiki.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 2
- Instability: 0.6667
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-render-wiki-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/smellentry--internal-vault-render-wiki-go-l13|smellEntry (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/buildstarterwikiarticles--internal-vault-render-wiki-go-l19|buildStarterWikiArticles (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/makewikifrontmatter--internal-vault-render-wiki-go-l43|makeWikiFrontmatter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderwikiarticle--internal-vault-render-wiki-go-l65|renderWikiArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderdashboard--internal-vault-render-wiki-go-l75|renderDashboard (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderconceptindex--internal-vault-render-wiki-go-l134|renderConceptIndex (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178|renderSourceIndex (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230|createCodebaseOverviewArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349|createDirectoryMapArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405|createSymbolTaxonomyArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471|createDependencyHotspotsArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547|createComplexityHotspotsArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createdeadcodereportarticle--internal-vault-render-wiki-go-l622|createDeadCodeReportArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692|createModuleHealthArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784|createCodeSmellsArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createcirculardependenciesarticle--internal-vault-render-wiki-go-l865|createCircularDependenciesArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createhighimpactsymbolsarticle--internal-vault-render-wiki-go-l892|createHighImpactSymbolsArticle (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/rendergroupedlinks--internal-vault-render-wiki-go-l964|renderGroupedLinks (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/rendersourcebulletlist--internal-vault-render-wiki-go-l978|renderSourceBulletList (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/uniquestrings--internal-vault-render-wiki-go-l991|uniqueStrings (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sortstrings--internal-vault-render-wiki-go-l1005|sortStrings (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildstarterwikiarticles--internal-vault-render-wiki-go-l19]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcirculardependenciesarticle--internal-vault-render-wiki-go-l865]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdeadcodereportarticle--internal-vault-render-wiki-go-l622]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createhighimpactsymbolsarticle--internal-vault-render-wiki-go-l892]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/makewikifrontmatter--internal-vault-render-wiki-go-l43]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderconceptindex--internal-vault-render-wiki-go-l134]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderdashboard--internal-vault-render-wiki-go-l75]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendergroupedlinks--internal-vault-render-wiki-go-l964]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendersourcebulletlist--internal-vault-render-wiki-go-l978]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderwikiarticle--internal-vault-render-wiki-go-l65]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/smellentry--internal-vault-render-wiki-go-l13]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortstrings--internal-vault-render-wiki-go-l1005]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquestrings--internal-vault-render-wiki-go-l991]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-render-wiki-go-l1]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `path`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strconv`
- `imports` (syntactic) -> `strings`

## Backlinks
None
