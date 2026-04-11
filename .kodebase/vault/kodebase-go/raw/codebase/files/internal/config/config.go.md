---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0.5
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 28
smells:
source_kind: "codebase-file"
source_path: "internal/config/config.go"
stage: "raw"
symbol_count: 12
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/config/config.go"
type: "source"
---

# Codebase File: internal/config/config.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 1
- Instability: 0.5
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/config--internal-config-config-go-l1|config (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/config--internal-config-config-go-l15|Config (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/appconfig--internal-config-config-go-l23|AppConfig (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/serverconfig--internal-config-config-go-l29|ServerConfig (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/logconfig--internal-config-config-go-l35|LogConfig (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/default--internal-config-config-go-l40|Default (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/load--internal-config-config-go-l59|Load (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/decodefile--internal-config-config-go-l72|decodeFile (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l95|Validate (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l112|Validate (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l125|Validate (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l133|Validate (method)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/appconfig--internal-config-config-go-l23]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/config--internal-config-config-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/config--internal-config-config-go-l15]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/decodefile--internal-config-config-go-l72]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/default--internal-config-config-go-l40]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/load--internal-config-config-go-l59]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/logconfig--internal-config-config-go-l35]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/serverconfig--internal-config-config-go-l29]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l112]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l125]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l133]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l95]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/appconfig--internal-config-config-go-l23]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/config--internal-config-config-go-l15]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/default--internal-config-config-go-l40]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/load--internal-config-config-go-l59]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/logconfig--internal-config-config-go-l35]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/serverconfig--internal-config-config-go-l29]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l112]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l125]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l133]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-config-go-l95]]
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/BurntSushi/toml`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
