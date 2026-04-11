---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 8
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
source_path: "internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go"
stage: "raw"
start_line: 6
symbol_kind: "function"
symbol_name: "Hello"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: Hello"
type: "source"
---

# Codebase Symbol: Hello

Source file: [[kodebase-go/raw/codebase/files/internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go|internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go]]

## Kind
`function`

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
func Hello(name string) string {
```

## Documentation
Hello formats a deterministic greeting for the fixture program.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go|internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go|internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go]] via `exports` (syntactic)
