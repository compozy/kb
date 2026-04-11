---
blast_radius: 5
centrality: 0.1457
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 363
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: true
language: "go"
loc: 70
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
start_line: 294
symbol_kind: "function"
symbol_name: "parsedFileFixture"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parsedFileFixture"
type: "source"
---

# Codebase Symbol: parsedFileFixture

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: true
- Blast radius: 5
- External references: 0
- Centrality: 0.1457
- LOC: 70
- Dead export: false
- Smells: `bottleneck`, `long-function`

## Signature
```text
func parsedFileFixture(relativePath string) models.ParsedFile {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/capitalize--internal-graph-normalize-test-go-l365]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphattachessymbolidstoparentfiles--internal-graph-normalize-test-go-l105]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphdeduplicatesfilessymbolsexternalnodesandrelations--internal-graph-normalize-test-go-l58]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphordersdiagnosticsbystagefilepathandmessage--internal-graph-normalize-test-go-l173]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphpassesthroughsingleparsedfile--internal-graph-normalize-test-go-l22]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphsortscollectionsdeterministically--internal-graph-normalize-test-go-l132]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `contains` (syntactic)
