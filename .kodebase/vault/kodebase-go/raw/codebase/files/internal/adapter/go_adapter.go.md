---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0.5
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 45
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
symbol_count: 33
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/adapter/go_adapter.go"
type: "source"
---

# Codebase File: internal/adapter/go_adapter.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 1
- Instability: 0.5
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-go-adapter-go-l1|adapter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/goadapter--internal-adapter-go-adapter-go-l35|GoAdapter (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/gosymbolmatch--internal-adapter-go-adapter-go-l37|goSymbolMatch (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/parsedgofile--internal-adapter-go-adapter-go-l42|parsedGoFile (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/supports--internal-adapter-go-adapter-go-l52|Supports (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsefiles--internal-adapter-go-adapter-go-l57|ParseFiles (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62|ParseFilesWithProgress (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161|parseGoFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractimports--internal-adapter-go-adapter-go-l264|extractImports (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/getgosymbolkind--internal-adapter-go-adapter-go-l302|getGoSymbolKind (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/creategosymbol--internal-adapter-go-adapter-go-l318|createGoSymbol (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolvegosymbolname--internal-adapter-go-adapter-go-l354|resolveGoSymbolName (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/formatgosignature--internal-adapter-go-adapter-go-l372|formatGoSignature (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractgodoccomment--internal-adapter-go-adapter-go-l381|extractGoDocComment (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractattachedcomment--internal-adapter-go-adapter-go-l395|extractAttachedComment (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractcalltargetnames--internal-adapter-go-adapter-go-l429|extractCallTargetNames (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/computecyclomaticcomplexity--internal-adapter-go-adapter-go-l465|computeCyclomaticComplexity (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/operatortext--internal-adapter-go-adapter-go-l503|operatorText (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/creategoparsediagnostic--internal-adapter-go-adapter-go-l520|createGoParseDiagnostic (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sortedexternalnodes--internal-adapter-go-adapter-go-l532|sortedExternalNodes (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551|walkNamed (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566|namedChildren (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/collectnodesbykind--internal-adapter-go-adapter-go-l577|collectNodesByKind (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588|textOf (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-adapter-go-adapter-go-l595|extractLeadingComment (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-adapter-go-adapter-go-l607|normalizeComment (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizelinecomment--internal-adapter-go-adapter-go-l631|normalizeLineComment (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/stripquotes--internal-adapter-go-adapter-go-l648|stripQuotes (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isgoexported--internal-adapter-go-adapter-go-l662|isGoExported (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createfileid--internal-adapter-go-adapter-go-l669|createFileID (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createexternalid--internal-adapter-go-adapter-go-l673|createExternalID (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-adapter-go-adapter-go-l677|createSymbolID (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/slugifysegment--internal-adapter-go-adapter-go-l688|slugifySegment (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-go-adapter-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectnodesbykind--internal-adapter-go-adapter-go-l577]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computecyclomaticcomplexity--internal-adapter-go-adapter-go-l465]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalid--internal-adapter-go-adapter-go-l673]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-adapter-go-adapter-go-l669]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/creategoparsediagnostic--internal-adapter-go-adapter-go-l520]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/creategosymbol--internal-adapter-go-adapter-go-l318]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-adapter-go-adapter-go-l677]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractattachedcomment--internal-adapter-go-adapter-go-l395]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractcalltargetnames--internal-adapter-go-adapter-go-l429]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractgodoccomment--internal-adapter-go-adapter-go-l381]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractimports--internal-adapter-go-adapter-go-l264]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-adapter-go-adapter-go-l595]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatgosignature--internal-adapter-go-adapter-go-l372]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getgosymbolkind--internal-adapter-go-adapter-go-l302]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/goadapter--internal-adapter-go-adapter-go-l35]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/gosymbolmatch--internal-adapter-go-adapter-go-l37]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isgoexported--internal-adapter-go-adapter-go-l662]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizecomment--internal-adapter-go-adapter-go-l607]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizelinecomment--internal-adapter-go-adapter-go-l631]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/operatortext--internal-adapter-go-adapter-go-l503]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsedgofile--internal-adapter-go-adapter-go-l42]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefiles--internal-adapter-go-adapter-go-l57]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvegosymbolname--internal-adapter-go-adapter-go-l354]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/slugifysegment--internal-adapter-go-adapter-go-l688]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedexternalnodes--internal-adapter-go-adapter-go-l532]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripquotes--internal-adapter-go-adapter-go-l648]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supports--internal-adapter-go-adapter-go-l52]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/goadapter--internal-adapter-go-adapter-go-l35]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefiles--internal-adapter-go-adapter-go-l57]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supports--internal-adapter-go-adapter-go-l52]]
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `tree_sitter (github.com/tree-sitter/go-tree-sitter)`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `regexp`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `unicode`

## Backlinks
None
