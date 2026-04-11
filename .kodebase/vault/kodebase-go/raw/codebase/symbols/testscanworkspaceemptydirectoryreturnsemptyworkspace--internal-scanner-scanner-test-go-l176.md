---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 194
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 1
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner_test.go"
stage: "raw"
start_line: 176
symbol_kind: "function"
symbol_name: "TestScanWorkspaceEmptyDirectoryReturnsEmptyWorkspace"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestScanWorkspaceEmptyDirectoryReturnsEmptyWorkspace"
type: "source"
---

# Codebase Symbol: TestScanWorkspaceEmptyDirectoryReturnsEmptyWorkspace

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 19
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestScanWorkspaceEmptyDirectoryReturnsEmptyWorkspace(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scantestworkspace--internal-scanner-scanner-test-go-l218]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `exports` (syntactic)
