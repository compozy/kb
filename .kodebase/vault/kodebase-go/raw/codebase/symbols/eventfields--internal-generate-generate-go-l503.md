---
blast_radius: 3
centrality: 0.1802
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 522
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 503
symbol_kind: "function"
symbol_name: "eventFields"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: eventFields"
type: "source"
---

# Codebase Symbol: eventFields

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.1802
- LOC: 20
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func eventFields(attrs ...any) map[string]any {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/emitstagecompleted--internal-generate-generate-go-l481]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/emitstageprogress--internal-generate-generate-go-l471]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/emitstagestarted--internal-generate-generate-go-l462]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
