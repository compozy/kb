---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 40
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 30
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter_integration_test.go"
stage: "raw"
start_line: 11
symbol_kind: "function"
symbol_name: "TestGoAdapterBuildsCrossFileCallRelations"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestGoAdapterBuildsCrossFileCallRelations"
type: "source"
---

# Codebase Symbol: TestGoAdapterBuildsCrossFileCallRelations

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_integration_test.go|internal/adapter/go_adapter_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 30
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestGoAdapterBuildsCrossFileCallRelations(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasrelation--internal-adapter-go-adapter-test-go-l405]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindsymbol--internal-adapter-go-adapter-test-go-l379]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsegosources--internal-adapter-go-adapter-test-go-l342]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_integration_test.go|internal/adapter/go_adapter_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_integration_test.go|internal/adapter/go_adapter_integration_test.go]] via `exports` (syntactic)
