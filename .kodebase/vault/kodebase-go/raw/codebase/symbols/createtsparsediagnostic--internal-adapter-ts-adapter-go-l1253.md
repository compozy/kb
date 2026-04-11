---
blast_radius: 2
centrality: 0.0542
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 1263
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
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1253
symbol_kind: "function"
symbol_name: "createTSParseDiagnostic"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createTSParseDiagnostic"
type: "source"
---

# Codebase Symbol: createTSParseDiagnostic

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0542
- LOC: 11
- Dead export: false
- Smells: None

## Signature
```text
func createTSParseDiagnostic(file models.ScannedSourceFile, detail string) models.StructuredDiagnostic {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
