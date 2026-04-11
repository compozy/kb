---
blast_radius: 8
centrality: 0.2248
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 77
exported: true
external_reference_count: 8
has_smells: true
incoming_relation_count: 10
is_dead_export: false
is_long_function: true
language: "go"
loc: 61
outgoing_relation_count: 3
smells:
  - "bottleneck"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize.go"
stage: "raw"
start_line: 17
symbol_kind: "function"
symbol_name: "NormalizeGraph"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: NormalizeGraph"
type: "source"
---

# Codebase Symbol: NormalizeGraph

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: true
- Blast radius: 8
- External references: 8
- Centrality: 0.2248
- LOC: 61
- Dead export: false
- Smells: `bottleneck`, `long-function`

## Signature
```text
func NormalizeGraph(rootPath string, parsedFiles []models.ParsedFile) models.GraphSnapshot {
```

## Documentation
NormalizeGraph merges parsed files into a single deterministically ordered graph snapshot.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/attachsymbolids--internal-graph-normalize-go-l87]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isrenderablefile--internal-graph-normalize-go-l79]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquebykey--internal-graph-normalize-go-l145]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphmergesoverlappingimportsacrossparsedfiles--internal-graph-normalize-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphattachessymbolidstoparentfiles--internal-graph-normalize-test-go-l105]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphdeduplicatesfilessymbolsexternalnodesandrelations--internal-graph-normalize-test-go-l58]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphomitsdiagnosticonlyfiles--internal-graph-normalize-test-go-l203]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphordersdiagnosticsbystagefilepathandmessage--internal-graph-normalize-test-go-l173]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphpassesthroughsingleparsedfile--internal-graph-normalize-test-go-l22]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphreturnsemptysnapshotfornoparsedfiles--internal-graph-normalize-test-go-l10]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphsortscollectionsdeterministically--internal-graph-normalize-test-go-l132]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]] via `exports` (syntactic)
