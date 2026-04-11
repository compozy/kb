---
blast_radius: 2
centrality: 0.0965
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 527
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 520
symbol_kind: "function"
symbol_name: "hasManagedGenerator"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: hasManagedGenerator"
type: "source"
---

# Codebase Symbol: hasManagedGenerator

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0965
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func hasManagedGenerator(content string) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/removemanagedwikiconcepts--internal-vault-writer-go-l489]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
