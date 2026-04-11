---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 8
smells:
source_kind: "codebase-file"
source_path: "cmd/kodebase/main.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: cmd/kodebase/main.go"
type: "source"
---

# Codebase File: cmd/kodebase/main.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: true
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/main--cmd-kodebase-main-go-l1|main (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/main--cmd-kodebase-main-go-l13|main (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/main--cmd-kodebase-main-go-l13]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/main--cmd-kodebase-main-go-l1]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/cli`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `os/signal`
- `imports` (syntactic) -> `syscall`

## Backlinks
None
