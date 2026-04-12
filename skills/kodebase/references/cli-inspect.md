# Inspect Command Reference

## Usage

```
kodebase inspect <subcommand> [flags]
```

## Shared Flags (All Subcommands)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `table` | Output format: `table`, `json`, or `tsv` |
| `--vault` | string | `""` | Vault root path (auto-discovered from cwd if omitted) |
| `--topic` | string | `""` | Topic slug inside the vault (auto-detected if only one topic exists) |

## Vault Auto-Discovery

When `--vault` is omitted, the CLI walks up from the current working directory looking for `.kb/vault/`. If `--topic` is omitted and only one topic exists, it is selected automatically. If multiple topics exist, the command fails with an error listing available slugs.

---

## Subcommands

### 1. smells

List symbols and files with detected code smells.

```
kodebase inspect smells [--type <smell-type>] [--format json]
```

**Flags:** `--type` (string) -- filter to a specific smell type (e.g., `long-function`, `high-complexity`, `dead-export`, `orphan-file`, `god-file`)

**Output Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `kind` | string | `"symbol"` or `"file"` |
| `name` | string | Symbol name or file source path |
| `source_path` | string | Source-relative file path |
| `symbol_kind` | string | Symbol kind (empty for files) |
| `smells` | string[] | List of detected smell types |

---

### 2. dead-code

List dead exports and orphan files.

```
kodebase inspect dead-code [--format json]
```

**Output Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `kind` | string | `"symbol"` or `"file"` |
| `name` | string | Symbol name or file source path |
| `source_path` | string | Source-relative file path |
| `symbol_kind` | string | Symbol kind (empty for files) |
| `reason` | string | `"dead-export"` or `"orphan-file"` |
| `smells` | string[] | List of detected smell types |

---

### 3. complexity

Rank functions by cyclomatic complexity (descending).

```
kodebase inspect complexity [--top N] [--format json]
```

**Flags:** `--top` (int, default 20) -- maximum number of rows to return

**Output Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `symbol_name` | string | Function or method name |
| `symbol_kind` | string | `"function"` or `"method"` |
| `source_path` | string | Source-relative file path |
| `cyclomatic_complexity` | int | Cyclomatic complexity score |
| `loc` | int | Lines of code |
| `blast_radius` | int | Transitive dependents count |
| `smells` | string[] | Detected smell types |

---

### 4. blast-radius

Rank symbols by blast radius (how many symbols transitively depend on a given symbol).

```
kodebase inspect blast-radius [--min N] [--top N] [--format json]
```

**Flags:**
- `--min` (int, default 0) -- minimum blast radius threshold
- `--top` (int, default 0) -- maximum rows to return (0 = all)

**Output Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `symbol_name` | string | Symbol name |
| `source_path` | string | Source-relative file path |
| `blast_radius` | int | Count of unique transitive dependents |
| `centrality` | float | Betweenness centrality score (0-1) |
| `external_reference_count` | int | References from outside the symbol's module |
| `smells` | string[] | Detected smell types |

---

### 5. coupling

Rank files by instability (Martin coupling metric).

```
kodebase inspect coupling [--unstable] [--format json]
```

**Flags:** `--unstable` (bool) -- only show files with instability > 0.5

**Output Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `source_path` | string | Source-relative file path |
| `afferent_coupling` | int | Files that import this file (Ca) |
| `efferent_coupling` | int | Files this file imports (Ce) |
| `instability` | float | Ce / (Ca + Ce); 1.0 = completely unstable |
| `has_circular_dependency` | bool | Participates in a circular import chain |
| `smells` | string[] | Detected smell types |

---

### 6. symbol \<name\>

Lookup symbols by case-insensitive substring match.

```
kodebase inspect symbol <name> [--format json]
```

**Behavior:**
- **No matches:** Returns error with suggestion to use `inspect smells` or `inspect complexity`
- **Single match:** Returns detailed field-value pairs (see detail output below)
- **Multiple matches:** Returns summary table

**Summary Table Columns** (multiple matches):

| Column | Type | Description |
|--------|------|-------------|
| `symbol_name` | string | Symbol name |
| `symbol_kind` | string | Symbol kind |
| `source_path` | string | Source-relative file path |
| `start_line` | int | Start line in source |
| `language` | string | Source language |
| `smells` | string[] | Detected smell types |

**Detail Fields** (single match):

| Field | Type |
|-------|------|
| `relative_path` | string |
| `symbol_name` | string |
| `symbol_kind` | string |
| `source_path` | string |
| `language` | string |
| `exported` | bool |
| `start_line` | int |
| `end_line` | int |
| `signature` | string |
| `loc` | int |
| `blast_radius` | int |
| `centrality` | float |
| `cyclomatic_complexity` | int |
| `external_reference_count` | int |
| `is_dead_export` | bool |
| `is_long_function` | bool |
| `smells` | string[] |
| `outgoing_relations` | relation[] |
| `backlinks` | relation[] |

Each relation entry has: `target_path` (string), `type` (string: imports|calls|references), `confidence` (string: semantic|syntactic).

---

### 7. file \<path\>

Lookup a file by its exact source path.

```
kodebase inspect file <path> [--format json]
```

**Detail Fields:**

| Field | Type |
|-------|------|
| `relative_path` | string |
| `source_path` | string |
| `language` | string |
| `symbol_count` | int |
| `symbols` | string[] (name + kind pairs) |
| `afferent_coupling` | int |
| `efferent_coupling` | int |
| `instability` | float |
| `is_orphan_file` | bool |
| `is_god_file` | bool |
| `has_circular_dependency` | bool |
| `smells` | string[] |
| `outgoing_relations` | relation[] |
| `backlinks` | relation[] |

---

### 8. backlinks \<name-or-path\>

Show incoming references for a symbol or file.

```
kodebase inspect backlinks <name-or-path> [--format json]
```

**Entity Resolution:** Tries exact file path match first, falls back to single symbol name match.

**Output Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `target_path` | string | Path of the referencing entity |
| `type` | string | Relation type: `imports`, `calls`, `references` |
| `confidence` | string | `semantic` or `syntactic` |

---

### 9. deps \<name-or-path\>

Show outgoing dependencies for a symbol or file.

```
kodebase inspect deps <name-or-path> [--format json]
```

**Entity Resolution:** Same as backlinks (file path first, then symbol name).

**Output Columns:**

| Column | Type | Description |
|--------|------|-------------|
| `target_path` | string | Path of the dependency |
| `type` | string | Relation type: `imports`, `calls`, `references` |
| `confidence` | string | `semantic` or `syntactic` |

---

### 10. circular-deps

List files that participate in circular dependencies.

```
kodebase inspect circular-deps [--format json]
```

**Behavior:**
- If cycles exist: returns a table of participating files
- If no cycles: returns `{"message": "no circular dependencies found"}`

**Output Columns** (when cycles exist):

| Column | Type | Description |
|--------|------|-------------|
| `source_path` | string | Source-relative file path |
| `afferent_coupling` | int | Files that import this file |
| `efferent_coupling` | int | Files this file imports |
| `instability` | float | Coupling instability metric |
| `smells` | string[] | Detected smell types |
