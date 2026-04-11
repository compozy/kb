---
blast_radius: 4
centrality: 0.1504
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 372
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 29
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 344
symbol_kind: "function"
symbol_name: "buildRules"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildRules"
type: "source"
---

# Codebase Symbol: buildRules

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.1504
- LOC: 29
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func buildRules(relativeDirectory string, patterns []string) []ignoreRule {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isnegatedpattern--internal-scanner-scanner-go-l383]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/shouldkeeppattern--internal-scanner-scanner-go-l374]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/builduserrules--internal-scanner-scanner-go-l308]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/collectignorerules--internal-scanner-scanner-go-l243]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/readgitignorerules--internal-scanner-scanner-go-l298]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
