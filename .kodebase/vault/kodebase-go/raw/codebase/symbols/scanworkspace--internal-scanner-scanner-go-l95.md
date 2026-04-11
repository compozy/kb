---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 191
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 97
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 95
symbol_kind: "method"
symbol_name: "ScanWorkspace"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: ScanWorkspace"
type: "source"
---

# Codebase Symbol: ScanWorkspace

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 97
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func (s *Scanner) ScanWorkspace(rootPath string) (*models.ScannedWorkspace, error) {
```

## Documentation
ScanWorkspace scans a repository root and returns supported source files grouped by language.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectignorerules--internal-scanner-scanner-go-l243]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveoutputpath--internal-scanner-scanner-go-l212]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `exports` (syntactic)
