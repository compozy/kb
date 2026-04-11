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
source_path: "internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go"
stage: "raw"
symbol_count: 2
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go"
type: "source"
---

# Codebase File: internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go

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
- [[kodebase-go/raw/codebase/symbols/greeter--internal-generate-testdata-fixture-go-repo-internal-greeter-greeter-go-l1|greeter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/hello--internal-generate-testdata-fixture-go-repo-internal-greeter-greeter-go-l6|Hello (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/greeter--internal-generate-testdata-fixture-go-repo-internal-greeter-greeter-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hello--internal-generate-testdata-fixture-go-repo-internal-greeter-greeter-go-l6]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hello--internal-generate-testdata-fixture-go-repo-internal-greeter-greeter-go-l6]]
- `imports` (syntactic) -> `fmt`

## Backlinks
None
