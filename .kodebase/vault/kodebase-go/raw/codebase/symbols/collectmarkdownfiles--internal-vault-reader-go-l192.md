---
blast_radius: 1
centrality: 0.0579
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 221
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 30
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 192
symbol_kind: "function"
symbol_name: "collectMarkdownFiles"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: collectMarkdownFiles"
type: "source"
---

# Codebase Symbol: collectMarkdownFiles

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0579
- LOC: 30
- Dead export: false
- Smells: None

## Signature
```text
func collectMarkdownFiles(rootPath string) ([]string, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/readvaultsnapshot--internal-vault-reader-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
