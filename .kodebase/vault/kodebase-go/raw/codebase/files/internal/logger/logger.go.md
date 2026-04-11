---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 28
smells:
source_kind: "codebase-file"
source_path: "internal/logger/logger.go"
stage: "raw"
symbol_count: 13
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/logger/logger.go"
type: "source"
---

# Codebase File: internal/logger/logger.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/logger--internal-logger-logger-go-l1|logger (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/logobserver--internal-logger-logger-go-l14|LogObserver (interface)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/option--internal-logger-logger-go-l19|Option (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/config--internal-logger-logger-go-l21|config (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/withwriter--internal-logger-logger-go-l28|WithWriter (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withobserver--internal-logger-logger-go-l37|WithObserver (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/new--internal-logger-logger-go-l45|New (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/observedhandler--internal-logger-logger-go-l78|observedHandler (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/enabled--internal-logger-logger-go-l83|Enabled (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/handle--internal-logger-logger-go-l87|Handle (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withattrs--internal-logger-logger-go-l94|WithAttrs (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withgroup--internal-logger-logger-go-l101|WithGroup (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parselevel--internal-logger-logger-go-l108|parseLevel (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/config--internal-logger-logger-go-l21]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/enabled--internal-logger-logger-go-l83]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/handle--internal-logger-logger-go-l87]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/logger--internal-logger-logger-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/logobserver--internal-logger-logger-go-l14]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/new--internal-logger-logger-go-l45]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/observedhandler--internal-logger-logger-go-l78]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/option--internal-logger-logger-go-l19]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parselevel--internal-logger-logger-go-l108]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withattrs--internal-logger-logger-go-l94]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withgroup--internal-logger-logger-go-l101]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withobserver--internal-logger-logger-go-l37]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withwriter--internal-logger-logger-go-l28]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/enabled--internal-logger-logger-go-l83]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/handle--internal-logger-logger-go-l87]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/logobserver--internal-logger-logger-go-l14]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/new--internal-logger-logger-go-l45]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/option--internal-logger-logger-go-l19]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withattrs--internal-logger-logger-go-l94]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withgroup--internal-logger-logger-go-l101]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withobserver--internal-logger-logger-go-l37]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withwriter--internal-logger-logger-go-l28]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `io`
- `imports` (syntactic) -> `log/slog`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `strings`

## Backlinks
None
