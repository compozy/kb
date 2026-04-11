---
created: "2026-04-11"
domain: "kodebase-go"
generator: "kodebase"
sources:
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/models/models.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/pathutils.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/generate.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/scanner/scanner.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/writer.go]]"
stage: "compiled"
tags:
  - "kodebase-go"
  - "wiki"
  - "codebase"
  - "starter"
title: "Dependency Hotspots"
type: "wiki"
updated: "2026-04-11"
---

# Dependency Hotspots

Hotspots are file-level raw snapshots with the highest observed relation density in the normalized graph. They are the best starting points when deciding what to compile into deeper conceptual articles next.

## Ranked Hotspots

| File | Relation count |
| ---- | -------------- |
| [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] | 70 |
| [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] | 69 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] | 57 |
| [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] | 55 |
| [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] | 49 |
| [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] | 48 |
| [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] | 45 |
| [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] | 43 |
| [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] | 43 |
| [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] | 43 |

## Interpretation

A high relation count usually indicates a coordination layer, a shared utility, or an entry point. Cross-check these files against [[kodebase-go/wiki/concepts/Module Health|Module Health]] to distinguish stable modules from unstable ones.

## Sources and Further Reading
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go]]
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go]]
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go]]
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go]]
- [[kodebase-go/raw/codebase/files/internal/models/models.go]]
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go]]
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go]]
- [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go]]
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go]]
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go]]
