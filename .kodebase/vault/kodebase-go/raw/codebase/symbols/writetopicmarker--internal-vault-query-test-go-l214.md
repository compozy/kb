---
blast_radius: 5
centrality: 0.1586
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 221
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/query_test.go"
stage: "raw"
start_line: 214
symbol_kind: "function"
symbol_name: "writeTopicMarker"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeTopicMarker"
type: "source"
---

# Codebase Symbol: writeTopicMarker

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.1586
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func writeTopicMarker(t *testing.T, topicPath string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testlistavailabletopicsreturnssortedtopics--internal-vault-query-test-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryautoresolvessingletopic--internal-vault-query-test-go-l72]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenmultipletopicsexist--internal-vault-query-test-go-l95]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryfindsvaultbywalkingup--internal-vault-query-test-go-l12]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryprefersexplicitvault--internal-vault-query-test-go-l41]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]] via `contains` (syntactic)
