---
blast_radius: 15
centrality: 0.1783
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 183
exported: true
external_reference_count: 13
has_smells: true
incoming_relation_count: 15
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 181
symbol_kind: "function"
symbol_name: "GetRawFileDocumentPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GetRawFileDocumentPath"
type: "source"
---

# Codebase Symbol: GetRawFileDocumentPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 15
- External references: 13
- Centrality: 0.1783
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func GetRawFileDocumentPath(filePath string) string {
```

## Documentation
GetRawFileDocumentPath derives the vault document path for a raw file snapshot.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizedocumentpathsegment--internal-vault-pathutils-go-l250]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createdocumentlookup--internal-vault-render-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawdirectoryindex--internal-vault-render-go-l524]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawfiledocument--internal-vault-render-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawlanguageindex--internal-vault-render-go-l600]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcirculardependenciesarticle--internal-vault-render-wiki-go-l865]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdeadcodereportarticle--internal-vault-render-wiki-go-l622]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createhighimpactsymbolsarticle--internal-vault-render-wiki-go-l892]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
