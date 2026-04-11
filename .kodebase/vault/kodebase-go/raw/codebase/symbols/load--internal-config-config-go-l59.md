---
blast_radius: 3
centrality: 0.1802
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 70
exported: true
external_reference_count: 3
has_smells: true
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 3
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/config/config.go"
stage: "raw"
start_line: 59
symbol_kind: "function"
symbol_name: "Load"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: Load"
type: "source"
---

# Codebase Symbol: Load

Source file: [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 3
- External references: 3
- Centrality: 0.1802
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func Load(path string) (Config, error) {
```

## Documentation
Load reads and validates the TOML config file, then overlays runtime secrets
from the environment.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/decodefile--internal-config-config-go-l72]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/default--internal-config-config-go-l40]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/loadsecretsfromenv--internal-config-env-go-l33]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testloadconfigroundtrip--internal-config-config-test-go-l33]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testloademptypathusesdefaults--internal-config-config-test-go-l74]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testloadrejectsunknownkeys--internal-config-config-test-go-l86]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `exports` (syntactic)
