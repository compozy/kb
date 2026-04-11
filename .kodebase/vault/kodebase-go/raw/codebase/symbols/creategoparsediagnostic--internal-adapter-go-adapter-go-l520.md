---
blast_radius: 2
centrality: 0.0573
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 530
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 520
symbol_kind: "function"
symbol_name: "createGoParseDiagnostic"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createGoParseDiagnostic"
type: "source"
---

# Codebase Symbol: createGoParseDiagnostic

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0573
- LOC: 11
- Dead export: false
- Smells: None

## Signature
```text
func createGoParseDiagnostic(file models.ScannedSourceFile, detail string) models.StructuredDiagnostic {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
