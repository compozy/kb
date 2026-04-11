---
created: "2026-04-11"
domain: "kodebase-go"
generator: "kodebase"
sources:
  - "[[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/models/models.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client_test.go]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/.agents/skills/systematic-debugging]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/cmd/kodebase]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/adapter]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/cli]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/config]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/generate]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/generate/testdata/fixture-go-repo]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/root]]"
  - "[[kodebase-go/raw/codebase/indexes/languages/go]]"
  - "[[kodebase-go/raw/codebase/indexes/languages/ts]]"
stage: "compiled"
tags:
  - "kodebase-go"
  - "wiki"
  - "codebase"
  - "starter"
title: "Codebase Overview"
type: "wiki"
updated: "2026-04-11"
---

# Codebase Overview

Kodebase Go currently compiles into a Karpathy-style topic where the codebase itself is staged in `raw/codebase/` and this starter wiki provides the first synthesized navigation layer. The corpus contains 80 parsed source files, 901 symbols, and 2637 extracted relations.

Start with [[kodebase-go/wiki/concepts/Module Health|Module Health]] for coupling, [[kodebase-go/wiki/concepts/Complexity Hotspots|Complexity Hotspots]] for function-level complexity, and [[kodebase-go/wiki/concepts/Dead Code Report|Dead Code Report]] for likely cleanup candidates.

## Language Coverage

| Language | Files | Symbols |
| -------- | ----- | ------- |
| [[kodebase-go/raw/codebase/indexes/languages/go|go]] | 79 | 898 |
| [[kodebase-go/raw/codebase/indexes/languages/ts|ts]] | 1 | 3 |

## Directory Coverage
- [[kodebase-go/raw/codebase/indexes/directories/root|.]] · 1 files · instability=0
- [[kodebase-go/raw/codebase/indexes/directories/.agents/skills/systematic-debugging|.agents/skills/systematic-debugging]] · 1 files · instability=0
- [[kodebase-go/raw/codebase/indexes/directories/cmd/kodebase|cmd/kodebase]] · 1 files · instability=0
- [[kodebase-go/raw/codebase/indexes/directories/internal/adapter|internal/adapter]] · 8 files · instability=0
- [[kodebase-go/raw/codebase/indexes/directories/internal/cli|internal/cli]] · 24 files · instability=0
- [[kodebase-go/raw/codebase/indexes/directories/internal/config|internal/config]] · 3 files · instability=0
- [[kodebase-go/raw/codebase/indexes/directories/internal/generate|internal/generate]] · 4 files · instability=0
- [[kodebase-go/raw/codebase/indexes/directories/internal/generate/testdata/fixture-go-repo|internal/generate/testdata/fixture-go-repo]] · 1 files · instability=0

## Relation Hotspots
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]]
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Sources and Further Reading
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go]]
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go]]
- [[kodebase-go/raw/codebase/files/internal/models/models.go]]
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go]]
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go]]
- [[kodebase-go/raw/codebase/indexes/directories/.agents/skills/systematic-debugging]]
- [[kodebase-go/raw/codebase/indexes/directories/cmd/kodebase]]
- [[kodebase-go/raw/codebase/indexes/directories/internal/adapter]]
- [[kodebase-go/raw/codebase/indexes/directories/internal/cli]]
- [[kodebase-go/raw/codebase/indexes/directories/internal/config]]
- [[kodebase-go/raw/codebase/indexes/directories/internal/generate]]
- [[kodebase-go/raw/codebase/indexes/directories/internal/generate/testdata/fixture-go-repo]]
- [[kodebase-go/raw/codebase/indexes/directories/root]]
- [[kodebase-go/raw/codebase/indexes/languages/go]]
- [[kodebase-go/raw/codebase/indexes/languages/ts]]
