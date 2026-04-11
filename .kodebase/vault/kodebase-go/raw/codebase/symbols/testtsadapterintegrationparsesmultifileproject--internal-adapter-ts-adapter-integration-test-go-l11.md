---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 46
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 36
outgoing_relation_count: 4
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter_integration_test.go"
stage: "raw"
start_line: 11
symbol_kind: "function"
symbol_name: "TestTSAdapterIntegrationParsesMultiFileProject"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestTSAdapterIntegrationParsesMultiFileProject"
type: "source"
---

# Codebase Symbol: TestTSAdapterIntegrationParsesMultiFileProject

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_integration_test.go|internal/adapter/ts_adapter_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 36
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestTSAdapterIntegrationParsesMultiFileProject(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasrelation--internal-adapter-go-adapter-test-go-l405]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindsymbol--internal-adapter-go-adapter-test-go-l379]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindparsedfile--internal-adapter-ts-adapter-test-go-l377]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsetslikesources--internal-adapter-ts-adapter-test-go-l340]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_integration_test.go|internal/adapter/ts_adapter_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_integration_test.go|internal/adapter/ts_adapter_integration_test.go]] via `exports` (syntactic)
