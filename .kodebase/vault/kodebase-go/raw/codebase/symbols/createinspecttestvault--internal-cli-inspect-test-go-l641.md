---
blast_radius: 3
centrality: 0.137
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 790
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: true
language: "go"
loc: 150
outgoing_relation_count: 2
smells:
  - "bottleneck"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 641
symbol_kind: "function"
symbol_name: "createInspectTestVault"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createInspectTestVault"
type: "source"
---

# Codebase Symbol: createInspectTestVault

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: true
- Blast radius: 3
- External references: 1
- Centrality: 0.137
- LOC: 150
- Dead export: false
- Smells: `bottleneck`, `long-function`

## Signature
```text
func createInspectTestVault(t *testing.T) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mkdirall--internal-cli-inspect-test-go-l802]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writeinspectmarkdown--internal-cli-inspect-test-go-l792]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testinspectcommandvalidationerrors--internal-cli-inspect-helpers-test-go-l214]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectcommandjsonformatproducesvalidjson--internal-cli-inspect-test-go-l550]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectcommandtsvformatproducesheaderandrows--internal-cli-inspect-test-go-l572]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
