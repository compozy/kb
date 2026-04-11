---
blast_radius: 2
centrality: 0.137
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 53
exported: true
external_reference_count: 2
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 51
symbol_kind: "function"
symbol_name: "Generate"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: Generate"
type: "source"
---

# Codebase Symbol: Generate

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.137
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func Generate(ctx context.Context, opts models.GenerateOptions) (models.GenerationSummary, error) {
```

## Documentation
Generate runs the full repository-to-vault pipeline and returns a structured
summary of the generated topic.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l57]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testgeneraterequiresrootpath--internal-generate-generate-test-go-l325]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgeneraterespectscanceledcontext--internal-generate-generate-test-go-l337]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `exports` (syntactic)
