---
blast_radius: 2
centrality: 0.0661
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 333
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 308
symbol_kind: "function"
symbol_name: "buildUserRules"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildUserRules"
type: "source"
---

# Codebase Symbol: buildUserRules

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0661
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func buildUserRules(options ScanOptions) []ignoreRule {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildrules--internal-scanner-scanner-go-l344]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizepattern--internal-scanner-scanner-go-l335]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/collectignorerules--internal-scanner-scanner-go-l243]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
