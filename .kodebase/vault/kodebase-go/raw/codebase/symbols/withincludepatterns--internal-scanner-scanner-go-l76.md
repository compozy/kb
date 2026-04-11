---
blast_radius: 1
centrality: 0.0615
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 80
exported: true
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 76
symbol_kind: "function"
symbol_name: "WithIncludePatterns"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WithIncludePatterns"
type: "source"
---

# Codebase Symbol: WithIncludePatterns

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0615
- LOC: 5
- Dead export: false
- Smells: None

## Signature
```text
func WithIncludePatterns(patterns ...string) Option {
```

## Documentation
WithIncludePatterns configures user include patterns that re-include paths.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceincludepatternrestrictsresults--internal-scanner-scanner-test-go-l123]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `exports` (syntactic)
