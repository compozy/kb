---
blast_radius: 2
centrality: 0.0661
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 306
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/scanner/scanner.go"
stage: "raw"
start_line: 298
symbol_kind: "function"
symbol_name: "readGitIgnoreRules"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: readGitIgnoreRules"
type: "source"
---

# Codebase Symbol: readGitIgnoreRules

Source file: [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0661
- LOC: 9
- Dead export: false
- Smells: None

## Signature
```text
func readGitIgnoreRules(rootPath string, relativePath string) ([]ignoreRule, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildrules--internal-scanner-scanner-go-l344]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/collectignorerules--internal-scanner-scanner-go-l243]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] via `contains` (syntactic)
