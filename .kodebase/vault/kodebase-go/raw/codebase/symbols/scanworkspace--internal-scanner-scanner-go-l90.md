---
blast_radius: 12
centrality: 0.2501
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 92
exported: true
external_reference_count: 2
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 90
symbol_kind: "function"
symbol_name: "ScanWorkspace"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ScanWorkspace"
type: "source"
---

# Codebase Symbol: ScanWorkspace

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 12
- External references: 2
- Centrality: 0.2501
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func ScanWorkspace(rootPath string, opts ...Option) (*models.ScannedWorkspace, error) {
```

## Documentation
ScanWorkspace is a convenience entrypoint for a one-off scan.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newscanner--internal-scanner-scanner-go-l57]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceintegrationnestedproject--internal-scanner-scanner-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/scantestworkspace--internal-scanner-scanner-test-go-l218]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `exports` (syntactic)
