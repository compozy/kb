# Search and Index Command Reference

Both commands require the QMD binary on PATH. Install with `npm install -g @tobilu/qmd`.

---

## Search Command

### Usage

```
kodebase search <query> [flags]
```

The `<query>` argument is the search text (required, non-empty).

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--lex` | bool | `false` | Use BM25 keyword search only |
| `--vec` | bool | `false` | Use vector similarity search only |
| `--limit` | int | `10` | Maximum number of results to return |
| `--min-score` | float | `0` | Minimum score threshold for returned matches |
| `--full` | bool | `false` | Show the full matched document content instead of snippets |
| `--all` | bool | `false` | Return all matches above the minimum score threshold |
| `--collection` | string | `""` | Use an explicit QMD collection name instead of deriving from the topic |
| `--format` | string | `table` | Output format: `table`, `json`, or `tsv` |
| `--vault` | string | `""` | Vault root path (used when deriving the collection name) |
| `--topic` | string | `""` | Topic slug (used when deriving the collection name) |

### Search Modes

| Mode | Flag | QMD Command | Description |
|------|------|-------------|-------------|
| Hybrid | (default) | `query` | Combines lexical and vector search |
| Lexical | `--lex` | `search` | BM25 keyword search only |
| Vector | `--vec` | `vsearch` | Embedding-based semantic search |

The `--lex` and `--vec` flags are mutually exclusive. Omit both for hybrid mode.

### Output Columns

| Column | Type | Description |
|--------|------|-------------|
| `path` | string | Vault-relative path of the matched document |
| `score` | float | Relevance score |
| `preview` | string | Snippet of matched content (or full content if `--full` is set) |

### Collection Name Derivation

When `--collection` is omitted, the collection name is derived from the topic slug:
1. Resolve the vault and topic (same logic as inspect commands)
2. Use the `topicSlug` as the collection name

### Example Invocations

```bash
# Hybrid search (default)
kodebase search "authentication middleware" --format json

# Lexical search with higher result limit
kodebase search "parseConfig" --lex --limit 20 --format json

# Vector search with score threshold
kodebase search "error handling patterns" --vec --min-score 0.5 --format json

# Full document content
kodebase search "auth" --full --format json

# Explicit collection name
kodebase search "auth" --collection my-project --format json
```

---

## Index Command

### Usage

```
kodebase index [flags]
```

Alias: `kodebase index-vault`

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--vault` | string | `""` | Vault root path |
| `--topic` | string | `""` | Topic slug inside the vault |
| `--name` | string | `""` | Override the derived QMD collection name |
| `--embed` | bool | `true` | Run embedding after syncing files |
| `--force-embed` | bool | `false` | Force re-embedding all documents |
| `--context` | string | `""` | Attach human-written collection context to improve search relevance |

### Idempotent Behavior

The index command is idempotent. It checks `qmd status` first and selects the operation:
- If the collection already exists: performs an **update** (syncs changes)
- If the collection does not exist: performs an **add** (creates and populates)

Run `kodebase index` repeatedly without side effects.

### Output Schema (indexResultPayload)

```
{
  "collectionName": string,       // QMD collection name (= topic slug or --name override)
  "embedRequested": bool,         // whether --embed was true
  "embedResult": {                // present only if embedding was performed
    "docsProcessed": int,
    "chunksEmbedded": int,
    "errors": int,
    "durationMs": int
  },
  "forceEmbed": bool,             // whether --force-embed was set
  "status": {
    "collection": {               // null if collection was just created
      "name": string,
      "path": string,
      "pattern": string,
      "documents": int,
      "lastUpdated": string
    },
    "hasVectorIndex": bool,
    "needsEmbedding": int,
    "totalDocuments": int
  },
  "topicPath": string,            // absolute path to the topic directory
  "topicSlug": string,            // topic identifier
  "updateResult": {
    "collections": int,
    "indexed": int,
    "updated": int,
    "unchanged": int,
    "removed": int,
    "needsEmbedding": int
  },
  "vaultPath": string             // absolute path to the vault root
}
```

### Example Invocations

```bash
# Index with default settings (embed enabled)
kodebase index

# Index with custom context for search relevance
kodebase index --context "React application with Redux state management"

# Force re-embedding all documents
kodebase index --force-embed

# Index without embedding (sync files only)
kodebase index --embed=false

# Index with explicit vault and topic
kodebase index --vault /path/to/vault --topic my-project

# Index with custom collection name
kodebase index --name custom-collection
```
