---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 84
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 82
symbol_kind: "method"
symbol_name: "Generate"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: Generate"
type: "source"
---

# Codebase Symbol: Generate

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 3
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func (r runner) Generate(ctx context.Context, opts models.GenerateOptions) (models.GenerationSummary, error) {
```

## Documentation
Generate runs the full repository-to-vault pipeline using the configured
runner dependencies. Tests use this to substitute stage implementations.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `exports` (syntactic)
