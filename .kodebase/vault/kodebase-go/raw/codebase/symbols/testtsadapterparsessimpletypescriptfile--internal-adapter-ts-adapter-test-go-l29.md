---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 76
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 48
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter_test.go"
stage: "raw"
start_line: 29
symbol_kind: "function"
symbol_name: "TestTSAdapterParsesSimpleTypeScriptFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestTSAdapterParsesSimpleTypeScriptFile"
type: "source"
---

# Codebase Symbol: TestTSAdapterParsesSimpleTypeScriptFile

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 48
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestTSAdapterParsesSimpleTypeScriptFile(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasrelation--internal-adapter-go-adapter-test-go-l405]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindsymbol--internal-adapter-go-adapter-test-go-l379]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesingletslikefile--internal-adapter-ts-adapter-test-go-l329]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `exports` (syntactic)
