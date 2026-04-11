---
blast_radius: 4
centrality: 0.1299
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 34
exported: true
external_reference_count: 4
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/logger/logger.go"
stage: "raw"
start_line: 28
symbol_kind: "function"
symbol_name: "WithWriter"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WithWriter"
type: "source"
---

# Codebase Symbol: WithWriter

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 4
- External references: 4
- Centrality: 0.1299
- LOC: 7
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func WithWriter(writer io.Writer) Option {
```

## Documentation
WithWriter overrides the log destination. Tests use this to capture the
JSON stream deterministically.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnewlevelfiltering--internal-logger-logger-test-go-l122]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnewuppercaselevelkey--internal-logger-logger-test-go-l76]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnewwithobserver--internal-logger-logger-test-go-l100]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnewwritesjson--internal-logger-logger-test-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `exports` (syntactic)
