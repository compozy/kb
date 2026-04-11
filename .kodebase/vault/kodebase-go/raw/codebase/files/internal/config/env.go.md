---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 13
smells:
source_kind: "codebase-file"
source_path: "internal/config/env.go"
stage: "raw"
symbol_count: 5
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/config/env.go"
type: "source"
---

# Codebase File: internal/config/env.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/config--internal-config-env-go-l1|config (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/secrets--internal-config-env-go-l26|Secrets (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/loadsecretsfromenv--internal-config-env-go-l33|LoadSecretsFromEnv (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/validate--internal-config-env-go-l41|Validate (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/loaddotenvifpresent--internal-config-env-go-l47|LoadDotEnvIfPresent (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/config--internal-config-env-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/loaddotenvifpresent--internal-config-env-go-l47]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/loadsecretsfromenv--internal-config-env-go-l33]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/secrets--internal-config-env-go-l26]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-env-go-l41]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/loaddotenvifpresent--internal-config-env-go-l47]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/loadsecretsfromenv--internal-config-env-go-l33]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/secrets--internal-config-env-go-l26]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validate--internal-config-env-go-l41]]
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `github.com/joho/godotenv`
- `imports` (syntactic) -> `os`

## Backlinks
None
