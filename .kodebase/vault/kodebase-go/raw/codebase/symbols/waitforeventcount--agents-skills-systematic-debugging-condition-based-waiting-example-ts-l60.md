---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 89
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "ts"
loc: 30
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: ".agents/skills/systematic-debugging/condition-based-waiting-example.ts"
stage: "raw"
start_line: 60
symbol_kind: "function"
symbol_name: "waitForEventCount"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "ts"
  - "function"
title: "Codebase Symbol: waitForEventCount"
type: "source"
---

# Codebase Symbol: waitForEventCount

Source file: [[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts|.agents/skills/systematic-debugging/condition-based-waiting-example.ts]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 30
- Dead export: true
- Smells: `dead-export`

## Signature
```text
function waitForEventCount(
  threadManager: ThreadManager,
  threadId: string,
  eventType: LaceEventType,
  count: number,
  timeoutMs = 5000
): Promise<LaceEvent[]>
```

## Documentation
Wait for a specific number of events of a given type

@param threadManager - The thread manager to query
@param threadId - Thread to check for events
@param eventType - Type of event to wait for
@param count - Number of events to wait for
@param timeoutMs - Maximum time to wait (default 5000ms)
@returns Promise resolving to all matching events once count is reached

Example:
// Wait for 2 AGENT_MESSAGE events (initial response + continuation)
await waitForEventCount(threadManager, agentThreadId, 'AGENT_MESSAGE', 2);

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts|.agents/skills/systematic-debugging/condition-based-waiting-example.ts]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts|.agents/skills/systematic-debugging/condition-based-waiting-example.ts]] via `exports` (syntactic)
