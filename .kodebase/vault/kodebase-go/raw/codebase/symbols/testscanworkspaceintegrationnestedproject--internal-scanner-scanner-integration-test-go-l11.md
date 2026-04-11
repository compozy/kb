---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 44
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 34
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner_integration_test.go"
stage: "raw"
start_line: 11
symbol_kind: "function"
symbol_name: "TestScanWorkspaceIntegrationNestedProject"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestScanWorkspaceIntegrationNestedProject"
type: "source"
---

# Codebase Symbol: TestScanWorkspaceIntegrationNestedProject

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner_integration_test.go|internal/scanner/scanner_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 34
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestScanWorkspaceIntegrationNestedProject(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l90]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withoutputpath--internal-scanner-scanner-go-l69]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetestfile--internal-scanner-scanner-test-go-l259]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_integration_test.go|internal/scanner/scanner_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_integration_test.go|internal/scanner/scanner_integration_test.go]] via `exports` (syntactic)
