---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 262
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 4
smells:
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 243
symbol_kind: "function"
symbol_name: "collectIgnoreRules"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: collectIgnoreRules"
type: "source"
---

# Codebase Symbol: collectIgnoreRules

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0723
- LOC: 20
- Dead export: false
- Smells: None

## Signature
```text
func collectIgnoreRules(rootPath string, options ScanOptions) ([]ignoreRule, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildrules--internal-scanner-scanner-go-l344]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/builduserrules--internal-scanner-scanner-go-l308]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectgitignorepaths--internal-scanner-scanner-go-l264]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readgitignorerules--internal-scanner-scanner-go-l298]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l95]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
