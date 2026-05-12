# Java Ingest Adapter — Task List

## Tasks

| # | Title | Status | Complexity | Dependencies |
|---|-------|--------|------------|--------------|
| 01 | Add Java language support to models and scanner | completed | medium | — |
| 02 | Integrate Tree-sitter Java language binding | completed | medium | — |
| 03 | Implement Java adapter MVP parsing pipeline | completed | high | task_01, task_02 |
| 04 | Register Java adapter in generate runner | completed | medium | task_03 |
| 05 | Add deep Java relation resolution with fallback | completed | high | task_03 |
| 06 | Validate Java ingest end-to-end with CLI and benchmark | completed | high | task_04, task_05 |
| 07 | Improve nested and inner Java type resolution | completed | high | — |
| 08 | Add wildcard import deep-resolution support | completed | high | — |
| 09 | Add deterministic policy for ambiguous import targets | completed | medium | task_08 |
| 10 | Add best-effort enterprise module metadata hints | completed | medium | — |
| 11 | Validate Phase 2 regression suite for Java fidelity | completed | high | task_07, task_08, task_09, task_10 |
| 12 | Add Java operational observability telemetry | completed | high | task_11 |
| 13 | Expand rollout benchmark corpus and reproducible gate | completed | high | task_11 |
| 14 | Stabilize JSON contract for automation consumers | completed | medium | task_11 |
| 15 | Create Java portfolio adoption playbook | completed | medium | task_12, task_13, task_14 |
| 16 | Add diagnostics governance checks in lint workflow | completed | high | task_12, task_11 |
| 17 | Harden large-scale Java ingest operational behavior | completed | high | task_12, task_11 |
