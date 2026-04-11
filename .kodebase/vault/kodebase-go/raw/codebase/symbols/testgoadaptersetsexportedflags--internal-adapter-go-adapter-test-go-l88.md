---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 106
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter_test.go"
stage: "raw"
start_line: 88
symbol_kind: "function"
symbol_name: "TestGoAdapterSetsExportedFlags"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestGoAdapterSetsExportedFlags"
type: "source"
---

# Codebase Symbol: TestGoAdapterSetsExportedFlags

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 19
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestGoAdapterSetsExportedFlags(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mustfindsymbol--internal-adapter-go-adapter-test-go-l379]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsesinglegofile--internal-adapter-go-adapter-test-go-l331]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] via `exports` (syntactic)
