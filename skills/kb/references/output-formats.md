# Output Format Reference

All `inspect` and `search` commands support three output formats via `--format`.

## Format Selection

| Format | Flag | Use Case |
|--------|------|----------|
| table | `--format table` | Human-readable display (default) |
| json | `--format json` | Programmatic parsing by agents |
| tsv | `--format tsv` | Piping to Unix tools |

Always use `--format json` when parsing output programmatically.

## Inspect Output (Tabular Commands)

Tabular inspect commands (`smells`, `dead-code`, `complexity`, `blast-radius`, `coupling`, `circular-deps`) return rows with typed columns.

### JSON Example (`inspect complexity --top 2 --format json`)

```json
[
  {
    "symbol_name": "parseConfig",
    "symbol_kind": "function",
    "source_path": "src/config.ts",
    "cyclomatic_complexity": 12,
    "loc": 45,
    "blast_radius": 8,
    "smells": ["high-complexity"]
  },
  {
    "symbol_name": "resolveImports",
    "symbol_kind": "function",
    "source_path": "src/resolver.ts",
    "cyclomatic_complexity": 9,
    "loc": 32,
    "blast_radius": 5,
    "smells": []
  }
]
```

### TSV Example

```
symbol_name	symbol_kind	source_path	cyclomatic_complexity	loc	blast_radius	smells
parseConfig	function	src/config.ts	12	45	8	high-complexity
resolveImports	function	src/resolver.ts	9	32	5	
```

## Inspect Output (Detail Commands)

Detail commands (`symbol`, `file`) return field-value pairs when a single entity matches.

### JSON Example (`inspect symbol parseConfig --format json`)

```json
[
  {"field": "symbol_name", "value": "parseConfig"},
  {"field": "symbol_kind", "value": "function"},
  {"field": "source_path", "value": "src/config.ts"},
  {"field": "loc", "value": 45},
  {"field": "blast_radius", "value": 8},
  {"field": "smells", "value": ["high-complexity"]},
  {"field": "outgoing_relations", "value": [
    {"target_path": "src/utils.ts", "type": "imports", "confidence": "syntactic"}
  ]},
  {"field": "backlinks", "value": [
    {"target_path": "src/main.ts", "type": "calls", "confidence": "semantic"}
  ]}
]
```

## Ingest Codebase Output

`kb ingest codebase` always outputs JSON to stdout (no `--format` flag).

```json
{
  "command": "generate",
  "rootPath": "/path/to/repo",
  "vaultPath": "/path/to/repo/.kb/vault",
  "topicPath": "/path/to/repo/.kb/vault/my-project",
  "topicSlug": "my-project",
  "filesScanned": 120,
  "filesParsed": 95,
  "filesSkipped": 25,
  "symbolsExtracted": 430,
  "relationsEmitted": 1200,
  "rawDocumentsWritten": 95,
  "wikiDocumentsWritten": 12,
  "indexDocumentsWritten": 5,
  "timings": {
    "scanMillis": 45,
    "selectAdaptersMillis": 2,
    "parseMillis": 1200,
    "normalizeMillis": 80,
    "metricsMillis": 150,
    "renderMillis": 300,
    "writeMillis": 200,
    "totalMillis": 1977
  },
  "diagnostics": []
}
```

## Search Output

### JSON Example (`search "auth middleware" --format json`)

```json
[
  {
    "path": "raw-codebase/src/auth/middleware.md",
    "score": 0.89,
    "preview": "Authentication middleware that validates JWT tokens..."
  }
]
```

## Index Output

`kb index` always outputs JSON to stdout (no `--format` flag).

```json
{
  "collectionName": "my-project",
  "embedRequested": true,
  "embedResult": {
    "docsProcessed": 95,
    "chunksEmbedded": 320,
    "errors": 0,
    "durationMs": 4500
  },
  "forceEmbed": false,
  "status": {
    "collection": {
      "name": "my-project",
      "path": "qmd://collections/my-project",
      "pattern": "",
      "documents": 95,
      "lastUpdated": "2026-04-10T12:00:00Z"
    },
    "hasVectorIndex": true,
    "needsEmbedding": 0,
    "totalDocuments": 95
  },
  "topicPath": "/path/to/vault/my-project",
  "topicSlug": "my-project",
  "updateResult": {
    "collections": 1,
    "indexed": 95,
    "updated": 0,
    "unchanged": 0,
    "removed": 0,
    "needsEmbedding": 95
  },
  "vaultPath": "/path/to/vault"
}
```

## Empty Results

| Format | Empty Output |
|--------|-------------|
| json | `[]` |
| table | `No results.` followed by newline |
| tsv | Header row only (no data rows) |
