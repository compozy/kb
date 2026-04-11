---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 53
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 41
outgoing_relation_count: 5
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner_test.go"
stage: "raw"
start_line: 13
symbol_kind: "function"
symbol_name: "TestScanWorkspaceRoutesSupportedFilesByLanguage"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestScanWorkspaceRoutesSupportedFilesByLanguage"
type: "source"
---

# Codebase Symbol: TestScanWorkspaceRoutesSupportedFilesByLanguage

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 41
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestScanWorkspaceRoutesSupportedFilesByLanguage(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withoutputpath--internal-scanner-scanner-go-l69]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupedcounts--internal-scanner-scanner-test-go-l229]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scannedpaths--internal-scanner-scanner-test-go-l250]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scantestworkspace--internal-scanner-scanner-test-go-l218]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetestfile--internal-scanner-scanner-test-go-l259]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `exports` (syntactic)
