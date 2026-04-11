---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 158
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 4
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner_test.go"
stage: "raw"
start_line: 144
symbol_kind: "function"
symbol_name: "TestScanWorkspaceExcludePatternRemovesMatches"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestScanWorkspaceExcludePatternRemovesMatches"
type: "source"
---

# Codebase Symbol: TestScanWorkspaceExcludePatternRemovesMatches

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 15
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestScanWorkspaceExcludePatternRemovesMatches(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withexcludepatterns--internal-scanner-scanner-go-l83]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scannedpaths--internal-scanner-scanner-test-go-l250]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scantestworkspace--internal-scanner-scanner-test-go-l218]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetestfile--internal-scanner-scanner-test-go-l259]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `exports` (syntactic)
