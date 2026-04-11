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
language: "ts"
outgoing_relation_count: 8
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: ".agents/skills/systematic-debugging/condition-based-waiting-example.ts"
stage: "raw"
symbol_count: 3
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "ts"
title: "Codebase File: .agents/skills/systematic-debugging/condition-based-waiting-example.ts"
type: "source"
---

# Codebase File: .agents/skills/systematic-debugging/condition-based-waiting-example.ts

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`ts`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
Complete implementation of condition-based waiting utilities
From: Lace test infrastructure improvements (2025-10-03)
Context: Fixed 15 flaky tests by replacing arbitrary timeouts

## Symbols
- [[kodebase-go/raw/codebase/symbols/waitforevent--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l20|waitForEvent (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/waitforeventcount--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l60|waitForEventCount (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/waitforeventmatch--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l111|waitForEventMatch (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/waitforevent--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l20]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/waitforeventcount--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l60]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/waitforeventmatch--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l111]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/waitforevent--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l20]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/waitforeventcount--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l60]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/waitforeventmatch--agents-skills-systematic-debugging-condition-based-waiting-example-ts-l111]]
- `imports` (syntactic) -> `~/threads/thread-manager`
- `imports` (syntactic) -> `~/threads/types`

## Backlinks
None
