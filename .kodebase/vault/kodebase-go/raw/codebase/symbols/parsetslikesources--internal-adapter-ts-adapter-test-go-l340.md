---
blast_radius: 10
centrality: 0.2288
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 375
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: false
language: "go"
loc: 36
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter_test.go"
stage: "raw"
start_line: 340
symbol_kind: "function"
symbol_name: "parseTSLikeSources"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseTSLikeSources"
type: "source"
---

# Codebase Symbol: parseTSLikeSources

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 10
- External references: 1
- Centrality: 0.2288
- LOC: 36
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseTSLikeSources(t *testing.T, sources map[string]string) []models.ParsedFile {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/languageforpath--internal-adapter-ts-adapter-test-go-l390]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtsadapterintegrationparsesmultifileproject--internal-adapter-ts-adapter-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsesingletslikefile--internal-adapter-ts-adapter-test-go-l329]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterbuildsimportbindingsfordefaultnamedandnamespaceimports--internal-adapter-ts-adapter-test-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterhandlesnamedandstarreexports--internal-adapter-ts-adapter-test-go-l236]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsesjavascriptrequireimports--internal-adapter-ts-adapter-test-go-l99]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `contains` (syntactic)
