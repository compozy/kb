---
blast_radius: 7
centrality: 0.1385
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 257
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 8
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner_test.go"
stage: "raw"
start_line: 250
symbol_kind: "function"
symbol_name: "scannedPaths"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: scannedPaths"
type: "source"
---

# Codebase Symbol: scannedPaths

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.1385
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func scannedPaths(files []models.ScannedSourceFile) []string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludepatternremovesmatches--internal-scanner-scanner-test-go-l144]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesgitdirectorybydefault--internal-scanner-scanner-test-go-l71]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesnodemodulesbydefault--internal-scanner-scanner-test-go-l55]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceincludepatternrestrictsresults--internal-scanner-scanner-test-go-l123]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsgitignorepatterns--internal-scanner-scanner-test-go-l87]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsnestedgitignorepatterns--internal-scanner-scanner-test-go-l104]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceroutessupportedfilesbylanguage--internal-scanner-scanner-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `contains` (syntactic)
