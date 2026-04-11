---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 221
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 25
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils_test.go"
stage: "raw"
start_line: 197
symbol_kind: "function"
symbol_name: "TestPathHelpersHandleEmptyInputs"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestPathHelpersHandleEmptyInputs"
type: "source"
---

# Codebase Symbol: TestPathHelpersHandleEmptyInputs

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils_test.go|internal/vault/pathutils_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 25
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestPathHelpersHandleEmptyInputs(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils_test.go|internal/vault/pathutils_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils_test.go|internal/vault/pathutils_test.go]] via `exports` (syntactic)
