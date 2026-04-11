---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 165
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 55
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_test.go"
stage: "raw"
start_line: 111
symbol_kind: "function"
symbol_name: "TestGenerateCommandWritesJSONEventsWhenRequested"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestGenerateCommandWritesJSONEventsWhenRequested"
type: "source"
---

# Codebase Symbol: TestGenerateCommandWritesJSONEventsWhenRequested

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_test.go|internal/cli/generate_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 55
- Dead export: true
- Smells: `dead-export`, `feature-envy`, `long-function`

## Signature
```text
func TestGenerateCommandWritesJSONEventsWhenRequested(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/generate_test.go|internal/cli/generate_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/generate_test.go|internal/cli/generate_test.go]] via `exports` (syntactic)
