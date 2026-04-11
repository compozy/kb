---
blast_radius: 5
centrality: 0.1586
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 338
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter_test.go"
stage: "raw"
start_line: 329
symbol_kind: "function"
symbol_name: "parseSingleTSLikeFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseSingleTSLikeFile"
type: "source"
---

# Codebase Symbol: parseSingleTSLikeFile

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.1586
- LOC: 10
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseSingleTSLikeFile(t *testing.T, relativePath string, source string) models.ParsedFile {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsetslikesources--internal-adapter-ts-adapter-test-go-l340]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtsadaptercomputescyclomaticcomplexity--internal-adapter-ts-adapter-test-go-l285]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterextractsclassandmethodsymbols--internal-adapter-ts-adapter-test-go-l140]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsessimpletypescriptfile--internal-adapter-ts-adapter-test-go-l29]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsestsxcomponent--internal-adapter-ts-adapter-test-go-l78]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterproducesdiagnosticsforparseerrors--internal-adapter-ts-adapter-test-go-l306]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `contains` (syntactic)
