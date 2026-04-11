---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 19
domain: "kodebase-go"
end_line: 304
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 217
outgoing_relation_count: 7
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 88
symbol_kind: "method"
symbol_name: "GenerateWithObserver"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: GenerateWithObserver"
type: "source"
---

# Codebase Symbol: GenerateWithObserver

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 19
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 217
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func (r runner) GenerateWithObserver(ctx context.Context, opts models.GenerateOptions, observer Observer) (models.GenerationSummary, error) {
```

## Documentation
GenerateWithObserver runs the configured pipeline and reports structured
events to the provided observer.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapternames--internal-generate-generate-go-l445]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/elapsedmillis--internal-generate-generate-go-l454]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filesforadapter--internal-generate-generate-go-l423]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/languagenames--internal-generate-generate-go-l436]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvepaths--internal-generate-generate-go-l340]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/selectadapters--internal-generate-generate-go-l407]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/workspacelanguages--internal-generate-generate-go-l390]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `exports` (syntactic)
