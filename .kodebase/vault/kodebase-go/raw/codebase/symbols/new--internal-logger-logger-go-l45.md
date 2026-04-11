---
blast_radius: 4
centrality: 0.1299
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 76
exported: true
external_reference_count: 4
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: false
language: "go"
loc: 32
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/logger/logger.go"
stage: "raw"
start_line: 45
symbol_kind: "function"
symbol_name: "New"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: New"
type: "source"
---

# Codebase Symbol: New

Source file: [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 4
- External references: 4
- Centrality: 0.1299
- LOC: 32
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func New(level string, opts ...Option) (*slog.Logger, error) {
```

## Documentation
New constructs the shared runtime logger with JSON output and optional
record observation hooks.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parselevel--internal-logger-logger-go-l108]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnewlevelfiltering--internal-logger-logger-test-go-l122]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnewuppercaselevelkey--internal-logger-logger-test-go-l76]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnewwithobserver--internal-logger-logger-test-go-l100]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnewwritesjson--internal-logger-logger-test-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] via `exports` (syntactic)
