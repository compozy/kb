---
blast_radius: 5
centrality: 0.16
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 217
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/query.go"
stage: "raw"
start_line: 210
symbol_kind: "function"
symbol_name: "isDirectoryPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: isDirectoryPath"
type: "source"
---

# Codebase Symbol: isDirectoryPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.16
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func isDirectoryPath(pathToCheck string) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/discovervaultpath--internal-vault-query-go-l28]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/ensuredirectory--internal-vault-query-go-l172]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `contains` (syntactic)
