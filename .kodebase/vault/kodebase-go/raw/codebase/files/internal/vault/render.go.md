---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 2
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0.5
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 40
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/vault/render.go"
stage: "raw"
symbol_count: 33
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/render.go"
type: "source"
---

# Codebase File: internal/vault/render.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 2
- Instability: 0.5
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault--internal-vault-render-go-l1|vault (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/starterwikiarticle--internal-vault-render-go-l13|starterWikiArticle (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20|RenderDocuments (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/renderfrontmatter--internal-vault-render-go-l130|renderFrontmatter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/rendermarkdowndocument--internal-vault-render-go-l164|renderMarkdownDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173|toSourceWikiLink (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createdocumentlookup--internal-vault-render-go-l177|createDocumentLookup (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createexternalnodelookup--internal-vault-render-go-l191|createExternalNodeLookup (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/linkfornode--internal-vault-render-go-l199|linkForNode (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderrelationlist--internal-vault-render-go-l221|renderRelationList (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderbacklinklist--internal-vault-render-go-l262|renderBacklinkList (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/rendersmelllist--internal-vault-render-go-l300|renderSmellList (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderrawfiledocument--internal-vault-render-go-l312|renderRawFileDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createsymbolfrontmatter--internal-vault-render-go-l402|createSymbolFrontmatter (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449|renderRawSymbolDocument (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderrawdirectoryindex--internal-vault-render-go-l524|renderRawDirectoryIndex (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/renderrawlanguageindex--internal-vault-render-go-l600|renderRawLanguageIndex (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/defaultfilemetrics--internal-vault-render-go-l667|defaultFileMetrics (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/defaultsymbolmetrics--internal-vault-render-go-l671|defaultSymbolMetrics (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/defaultdirectorymetrics--internal-vault-render-go-l678|defaultDirectoryMetrics (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/toprelationhotspotfiles--internal-vault-render-go-l682|topRelationHotspotFiles (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/grouprelationsbysource--internal-vault-render-go-l719|groupRelationsBySource (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/grouprelationsbytarget--internal-vault-render-go-l727|groupRelationsByTarget (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupsymbolsbyfile--internal-vault-render-go-l735|groupSymbolsByFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupfilesbydirectory--internal-vault-render-go-l743|groupFilesByDirectory (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupsymbolsbydirectory--internal-vault-render-go-l751|groupSymbolsByDirectory (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupfilesbylanguage--internal-vault-render-go-l759|groupFilesByLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupsymbolsbylanguage--internal-vault-render-go-l767|groupSymbolsByLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupsymbolsbykind--internal-vault-render-go-l775|groupSymbolsByKind (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sortsymbolsbylocation--internal-vault-render-go-l783|sortSymbolsByLocation (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sortedmapkeys--internal-vault-render-go-l798|sortedMapKeys (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isfunctionlike--internal-vault-render-go-l807|isFunctionLike (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/maxint--internal-vault-render-go-l811|maxInt (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdocumentlookup--internal-vault-render-go-l177]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalnodelookup--internal-vault-render-go-l191]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymbolfrontmatter--internal-vault-render-go-l402]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/defaultdirectorymetrics--internal-vault-render-go-l678]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/defaultfilemetrics--internal-vault-render-go-l667]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/defaultsymbolmetrics--internal-vault-render-go-l671]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupfilesbydirectory--internal-vault-render-go-l743]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupfilesbylanguage--internal-vault-render-go-l759]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/grouprelationsbysource--internal-vault-render-go-l719]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/grouprelationsbytarget--internal-vault-render-go-l727]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbydirectory--internal-vault-render-go-l751]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbyfile--internal-vault-render-go-l735]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbykind--internal-vault-render-go-l775]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbylanguage--internal-vault-render-go-l767]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isfunctionlike--internal-vault-render-go-l807]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/linkfornode--internal-vault-render-go-l199]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/maxint--internal-vault-render-go-l811]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbacklinklist--internal-vault-render-go-l262]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderfrontmatter--internal-vault-render-go-l130]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendermarkdowndocument--internal-vault-render-go-l164]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawdirectoryindex--internal-vault-render-go-l524]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawfiledocument--internal-vault-render-go-l312]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawlanguageindex--internal-vault-render-go-l600]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrelationlist--internal-vault-render-go-l221]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendersmelllist--internal-vault-render-go-l300]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedmapkeys--internal-vault-render-go-l798]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortsymbolsbylocation--internal-vault-render-go-l783]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/starterwikiarticle--internal-vault-render-go-l13]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toprelationhotspotfiles--internal-vault-render-go-l682]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault--internal-vault-render-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `path`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strconv`
- `imports` (syntactic) -> `strings`

## Backlinks
None
