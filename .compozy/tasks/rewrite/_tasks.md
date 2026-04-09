# Kodebase Go Port — Task List

**Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use it as the behavioral and structural source for this Go rewrite on every task.

## Tasks

| # | Title | Status | Complexity | Dependencies |
|---|-------|--------|------------|--------------|
| 01 | Rename module and bootstrap cobra CLI skeleton | pending | medium | — |
| 02 | Define domain models | pending | medium | task_01 |
| 03 | Implement workspace scanner | completed | medium | task_02 |
| 04 | Set up tree-sitter infrastructure | completed | low | task_01 |
| 05 | Implement Go language adapter | completed | high | task_02, task_04 |
| 06 | Implement TypeScript/JavaScript adapter | completed | high | task_02, task_04 |
| 07 | Implement graph normalizer | completed | medium | task_02 |
| 08 | Implement metrics engine | completed | high | task_02, task_07 |
| 09 | Implement vault path and text utilities | completed | low | task_01 |
| 10 | Implement document renderer | pending | critical | task_02, task_09 |
| 11 | Implement vault writer | completed | medium | task_10 |
| 12 | Wire generate command end-to-end | completed | high | task_03, task_05, task_06, task_08, task_11 |
| 13 | Implement vault reader and query resolver | completed | medium | task_02, task_09 |
| 14 | Implement output formatter | completed | low | task_01 |
| 15 | Wire inspect analysis subcommands | completed | high | task_13, task_14 |
| 16 | Wire inspect lookup subcommands | completed | medium | task_15 |
| 17 | Implement QMD shell client | completed | medium | task_01 |
| 18 | Wire search and index-vault commands | completed | medium | task_13, task_17 |
| 19 | End-to-end integration test | pending | medium | task_12, task_16, task_18 |
| 20 | Update documentation and final verification | completed | low | task_19 |
