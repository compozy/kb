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
outgoing_relation_count: 4
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/version/version.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/version/version.go"
type: "source"
---

# Codebase File: internal/version/version.go

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
- [[kodebase-go/raw/codebase/symbols/version--internal-version-version-go-l1|version (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/string--internal-version-version-go-l11|String (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/string--internal-version-version-go-l11]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/version--internal-version-version-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/string--internal-version-version-go-l11]]
- `imports` (syntactic) -> `fmt`

## Backlinks
None
