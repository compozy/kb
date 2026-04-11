---
blast_radius: 1
centrality: 0.0615
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 87
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
start_line: 83
symbol_kind: "function"
symbol_name: "WithExcludePatterns"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WithExcludePatterns"
type: "source"
---

# Codebase Symbol: WithExcludePatterns

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
func WithExcludePatterns(patterns ...string) Option {
```

## Documentation
WithExcludePatterns configures user exclude patterns.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludepatternremovesmatches--internal-scanner-scanner-test-go-l144]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `exports` (syntactic)
