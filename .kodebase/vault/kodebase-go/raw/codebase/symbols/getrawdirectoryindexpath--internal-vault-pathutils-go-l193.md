---
blast_radius: 6
centrality: 0.0827
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 200
exported: true
external_reference_count: 4
has_smells: false
incoming_relation_count: 6
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 193
symbol_kind: "function"
symbol_name: "GetRawDirectoryIndexPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GetRawDirectoryIndexPath"
type: "source"
---

# Codebase Symbol: GetRawDirectoryIndexPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 6
- External references: 4
- Centrality: 0.0827
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func GetRawDirectoryIndexPath(directoryPath string) string {
```

## Documentation
GetRawDirectoryIndexPath derives the vault document path for a raw directory index.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizedocumentpathsegment--internal-vault-pathutils-go-l250]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderrawdirectoryindex--internal-vault-render-go-l524]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
