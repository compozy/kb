---
blast_radius: 10
centrality: 0.1888
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 270
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 11
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner_test.go"
stage: "raw"
start_line: 259
symbol_kind: "function"
symbol_name: "writeTestFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeTestFile"
type: "source"
---

# Codebase Symbol: writeTestFile

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 10
- External references: 1
- Centrality: 0.1888
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func writeTestFile(t *testing.T, rootPath string, relativePath string, contents string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceintegrationnestedproject--internal-scanner-scanner-integration-test-go-l11]] via `calls` (syntactic)
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
