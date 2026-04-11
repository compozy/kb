---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 15
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/generate/events.go"
stage: "raw"
symbol_count: 8
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/generate/events.go"
type: "source"
---

# Codebase File: internal/generate/events.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/generate--internal-generate-events-go-l1|generate (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/eventkind--internal-generate-events-go-l6|EventKind (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/event--internal-generate-events-go-l16|Event (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/observer--internal-generate-events-go-l27|Observer (interface)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/observerfunc--internal-generate-events-go-l32|ObserverFunc (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-generate-events-go-l35|ObserveGenerateEvent (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/noopobserver--internal-generate-events-go-l41|noopObserver (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-generate-events-go-l43|ObserveGenerateEvent (method)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/event--internal-generate-events-go-l16]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/eventkind--internal-generate-events-go-l6]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generate--internal-generate-events-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/noopobserver--internal-generate-events-go-l41]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-generate-events-go-l35]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-generate-events-go-l43]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observer--internal-generate-events-go-l27]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observerfunc--internal-generate-events-go-l32]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/event--internal-generate-events-go-l16]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/eventkind--internal-generate-events-go-l6]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-generate-events-go-l35]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observegenerateevent--internal-generate-events-go-l43]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observer--internal-generate-events-go-l27]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observerfunc--internal-generate-events-go-l32]]
- `imports` (syntactic) -> `context`

## Backlinks
None
