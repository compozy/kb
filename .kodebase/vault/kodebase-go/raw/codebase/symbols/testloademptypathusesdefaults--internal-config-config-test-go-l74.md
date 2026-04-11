---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 84
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/config/config_test.go"
stage: "raw"
start_line: 74
symbol_kind: "function"
symbol_name: "TestLoadEmptyPathUsesDefaults"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestLoadEmptyPathUsesDefaults"
type: "source"
---

# Codebase Symbol: TestLoadEmptyPathUsesDefaults

Source file: [[kodebase-go/raw/codebase/files/internal/config/config_test.go|internal/config/config_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 11
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestLoadEmptyPathUsesDefaults(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/load--internal-config-config-go-l59]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/config/config_test.go|internal/config/config_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config_test.go|internal/config/config_test.go]] via `exports` (syntactic)
