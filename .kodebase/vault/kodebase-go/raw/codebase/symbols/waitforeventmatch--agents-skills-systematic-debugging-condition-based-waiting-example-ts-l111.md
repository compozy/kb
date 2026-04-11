---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 136
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "ts"
loc: 26
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: ".agents/skills/systematic-debugging/condition-based-waiting-example.ts"
stage: "raw"
start_line: 111
symbol_kind: "function"
symbol_name: "waitForEventMatch"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "ts"
  - "function"
title: "Codebase Symbol: waitForEventMatch"
type: "source"
---

# Codebase Symbol: waitForEventMatch

Source file: [[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts|.agents/skills/systematic-debugging/condition-based-waiting-example.ts]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 26
- Dead export: true
- Smells: `dead-export`

## Signature
```text
function waitForEventMatch(
  threadManager: ThreadManager,
  threadId: string,
  predicate: (event: LaceEvent) => boolean,
  description: string,
  timeoutMs = 5000
): Promise<LaceEvent>
```

## Documentation
Wait for an event matching a custom predicate
Useful when you need to check event data, not just type

@param threadManager - The thread manager to query
@param threadId - Thread to check for events
@param predicate - Function that returns true when event matches
@param description - Human-readable description for error messages
@param timeoutMs - Maximum time to wait (default 5000ms)
@returns Promise resolving to the first matching event

Example:
// Wait for TOOL_RESULT with specific ID
await waitForEventMatch(
threadManager,
agentThreadId,
(e) => e.type === 'TOOL_RESULT' && e.data.id === 'call_123',
'TOOL_RESULT with id=call_123'
);

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts|.agents/skills/systematic-debugging/condition-based-waiting-example.ts]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts|.agents/skills/systematic-debugging/condition-based-waiting-example.ts]] via `exports` (syntactic)
