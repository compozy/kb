---
blast_radius: 3
centrality: 0.1673
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 59
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: true
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 57
symbol_kind: "function"
symbol_name: "GenerateWithObserver"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GenerateWithObserver"
type: "source"
---

# Codebase Symbol: GenerateWithObserver

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1673
- LOC: 3
- Dead export: true
- Smells: `bottleneck`, `dead-export`

## Signature
```text
func GenerateWithObserver(ctx context.Context, opts models.GenerateOptions, observer Observer) (models.GenerationSummary, error) {
```

## Documentation
GenerateWithObserver runs the full pipeline and emits structured events to
the provided observer while returning the final summary.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrunner--internal-generate-generate-go-l61]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-go-l51]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `exports` (syntactic)
