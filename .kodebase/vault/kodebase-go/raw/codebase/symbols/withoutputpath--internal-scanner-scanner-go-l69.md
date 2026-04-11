---
blast_radius: 2
centrality: 0.0738
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 73
exported: true
external_reference_count: 2
has_smells: false
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 69
symbol_kind: "function"
symbol_name: "WithOutputPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WithOutputPath"
type: "source"
---

# Codebase Symbol: WithOutputPath

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.0738
- LOC: 5
- Dead export: false
- Smells: None

## Signature
```text
func WithOutputPath(path string) Option {
```

## Documentation
WithOutputPath excludes the generated output directory from scan results.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceintegrationnestedproject--internal-scanner-scanner-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceroutessupportedfilesbylanguage--internal-scanner-scanner-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `exports` (syntactic)
