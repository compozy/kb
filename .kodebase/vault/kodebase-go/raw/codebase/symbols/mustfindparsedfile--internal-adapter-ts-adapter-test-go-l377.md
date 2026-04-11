---
blast_radius: 4
centrality: 0.0939
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 388
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter_test.go"
stage: "raw"
start_line: 377
symbol_kind: "function"
symbol_name: "mustFindParsedFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: mustFindParsedFile"
type: "source"
---

# Codebase Symbol: mustFindParsedFile

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 4
- External references: 1
- Centrality: 0.0939
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func mustFindParsedFile(t *testing.T, parsedFiles []models.ParsedFile, relativePath string) models.ParsedFile {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtsadapterintegrationparsesmultifileproject--internal-adapter-ts-adapter-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterbuildsimportbindingsfordefaultnamedandnamespaceimports--internal-adapter-ts-adapter-test-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterhandlesnamedandstarreexports--internal-adapter-ts-adapter-test-go-l236]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsesjavascriptrequireimports--internal-adapter-ts-adapter-test-go-l99]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `contains` (syntactic)
