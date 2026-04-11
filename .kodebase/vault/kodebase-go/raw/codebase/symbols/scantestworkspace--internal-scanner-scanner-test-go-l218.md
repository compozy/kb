---
blast_radius: 10
centrality: 0.2176
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 227
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 11
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner_test.go"
stage: "raw"
start_line: 218
symbol_kind: "function"
symbol_name: "scanTestWorkspace"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: scanTestWorkspace"
type: "source"
---

# Codebase Symbol: scanTestWorkspace

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 10
- External references: 0
- Centrality: 0.2176
- LOC: 10
- Dead export: false
- Smells: `bottleneck`, `feature-envy`

## Signature
```text
func scanTestWorkspace(t *testing.T, rootPath string, opts ...Option) *models.ScannedWorkspace {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l90]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceemptydirectoryreturnsemptyworkspace--internal-scanner-scanner-test-go-l176]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludepatternremovesmatches--internal-scanner-scanner-test-go-l144]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesgitdirectorybydefault--internal-scanner-scanner-test-go-l71]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesnodemodulesbydefault--internal-scanner-scanner-test-go-l55]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspacegroupsfilesbylanguage--internal-scanner-scanner-test-go-l196]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceignoresunsupportedextensions--internal-scanner-scanner-test-go-l160]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceincludepatternrestrictsresults--internal-scanner-scanner-test-go-l123]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsgitignorepatterns--internal-scanner-scanner-test-go-l87]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsnestedgitignorepatterns--internal-scanner-scanner-test-go-l104]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceroutessupportedfilesbylanguage--internal-scanner-scanner-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] via `contains` (syntactic)
