---
blast_radius: 12
centrality: 0.1433
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 190
exported: true
external_reference_count: 10
has_smells: true
incoming_relation_count: 12
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 186
symbol_kind: "function"
symbol_name: "GetRawSymbolDocumentPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GetRawSymbolDocumentPath"
type: "source"
---

# Codebase Symbol: GetRawSymbolDocumentPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 12
- External references: 10
- Centrality: 0.1433
- LOC: 5
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func GetRawSymbolDocumentPath(symbol models.SymbolNode) string {
```

## Documentation
GetRawSymbolDocumentPath derives the vault document path for a raw symbol snapshot.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/slugifysegment--internal-vault-pathutils-go-l87]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createdocumentlookup--internal-vault-render-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawdirectoryindex--internal-vault-render-go-l524]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawfiledocument--internal-vault-render-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawlanguageindex--internal-vault-render-go-l600]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdeadcodereportarticle--internal-vault-render-wiki-go-l622]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createhighimpactsymbolsarticle--internal-vault-render-wiki-go-l892]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
