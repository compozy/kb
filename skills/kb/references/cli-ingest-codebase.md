# Ingest Codebase Command Reference

## Usage

```
kb ingest codebase <path> [flags]
```

The `<path>` argument is the root directory of the source repository to analyze (required).

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--topic` | string | `""` | Topic slug for the ingested codebase (required) |
| `--vault` | string | `""` | Vault root where the generated topic will be written. Defaults to `<path>/.kb/vault` |
| `--output` | string | `""` | Deprecated alias for `--vault` |
| `--title` | string | `""` | Bootstrap-only topic title override for a missing topic |
| `--domain` | string | `""` | Bootstrap-only topic domain override for a missing topic |
| `--include` | string[] | `nil` | Re-include a path pattern that would otherwise be ignored; repeatable |
| `--exclude` | string[] | `nil` | Exclude an additional path pattern from scanning; repeatable |
| `--semantic` | bool | `false` | Enable semantic analysis when the underlying adapters support it |
| `--progress` | string | `auto` | Progress rendering mode: `auto`, `always`, or `never` |
| `--log-format` | string | `text` | Stderr event format: `text` or `json` |

## Non-Interactive Usage

When invoking from an agent context, always set `--progress never` to prevent TTY progress bars from corrupting stdout output.

```
kb ingest codebase /path/to/repo --topic my-project --progress never
```

On first run, this command bootstraps the topic automatically. If the topic already exists, `--title` and `--domain` are rejected so re-ingest cannot silently mutate topic identity.

## Pipeline Stages

The codebase ingestion pipeline executes these stages in order:

1. **scan** -- Discover source files by language
2. **select_adapters** -- Choose language parsers (tree-sitter for TS/JS, Go parser)
3. **parse** -- Extract AST nodes, symbols, and relations
4. **normalize** -- Merge per-file graphs into a unified snapshot, resolve imports
5. **metrics** -- Compute complexity, coupling, blast radius, dead code, smells
6. **render** -- Generate markdown documents and Base definitions
7. **write** -- Persist vault files to disk

## Supported Languages

| Language | Extensions | Adapter |
|----------|-----------|---------|
| TypeScript | `.ts` | tree-sitter |
| TSX | `.tsx` | tree-sitter |
| JavaScript | `.js` | tree-sitter |
| JSX | `.jsx` | tree-sitter |
| Go | `.go` | tree-sitter |

## Output Schema (GenerationSummary)

The command writes JSON to stdout. Parse the following fields:

```
{
  "command": string,           // always "generate"
  "rootPath": string,          // absolute path to the analyzed repository
  "vaultPath": string,         // absolute path to the vault root
  "topicPath": string,         // absolute path to the topic directory
  "topicSlug": string,         // topic identifier (use for --topic in later commands)
  "filesScanned": int,         // total files discovered
  "filesParsed": int,          // files successfully parsed
  "filesSkipped": int,         // files skipped (unsupported or excluded)
  "symbolsExtracted": int,     // total symbols extracted
  "relationsEmitted": int,     // total relation edges
  "rawDocumentsWritten": int,  // per-file markdown documents
  "wikiDocumentsWritten": int, // concept wiki articles
  "indexDocumentsWritten": int, // index pages
  "timings": {
    "scanMillis": int,
    "selectAdaptersMillis": int,
    "parseMillis": int,
    "normalizeMillis": int,
    "metricsMillis": int,
    "renderMillis": int,
    "writeMillis": int,
    "totalMillis": int
  },
  "diagnostics": [             // structured warnings/errors
    {
      "code": string,
      "severity": "warning" | "error",
      "stage": "scan" | "parse" | "render" | "write" | "validate",
      "message": string,
      "filePath": string?,
      "language": string?,
      "detail": string?
    }
  ]
}
```

## Vault Structure

After ingestion, the topic directory contains:

```
<vaultPath>/<topicSlug>/
  raw/codebase/     # One markdown file per source file and symbol snapshot
  wiki/codebase/    # Generated codebase wiki articles and indexes
  wiki/index/       # Topic landing pages with bridges into wiki/codebase/
  bases/            # Obsidian Base view definitions
  CLAUDE.md         # Topic marker and schema document
  AGENTS.md         # Symlink to CLAUDE.md
  log.md            # Append-only audit log
```

## Default Path Derivation

- If `--vault` is omitted: vault path defaults to `<rootPath>/.kb/vault`
- If `--output` is provided: it behaves like `--vault` but is deprecated
- Full topic path: `<vaultPath>/<topicSlug>/`
