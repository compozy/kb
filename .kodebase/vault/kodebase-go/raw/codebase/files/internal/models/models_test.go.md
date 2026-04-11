---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 6
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/models/models_test.go"
stage: "raw"
symbol_count: 3
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/models/models_test.go"
type: "source"
---

# Codebase File: internal/models/models_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/models--internal-models-models-test-go-l1|models (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testsupportedlanguages--internal-models-models-test-go-l5|TestSupportedLanguages (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testdocumentandbaseconstants--internal-models-models-test-go-l26|TestDocumentAndBaseConstants (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/models--internal-models-models-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testdocumentandbaseconstants--internal-models-models-test-go-l26]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsupportedlanguages--internal-models-models-test-go-l5]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testdocumentandbaseconstants--internal-models-models-test-go-l26]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testsupportedlanguages--internal-models-models-test-go-l5]]
- `imports` (syntactic) -> `testing`

## Backlinks
None
