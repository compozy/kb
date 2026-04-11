---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 327
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 22
outgoing_relation_count: 1
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter_test.go"
stage: "raw"
start_line: 306
symbol_kind: "function"
symbol_name: "TestTSAdapterProducesDiagnosticsForParseErrors"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestTSAdapterProducesDiagnosticsForParseErrors"
type: "source"
---

# Codebase Symbol: TestTSAdapterProducesDiagnosticsForParseErrors

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 22
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestTSAdapterProducesDiagnosticsForParseErrors(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesingletslikefile--internal-adapter-ts-adapter-test-go-l329]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `exports` (syntactic)
