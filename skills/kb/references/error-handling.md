# Error Handling Reference

Categorized error messages from the `kb` CLI with causes and recovery steps.

## Vault Resolution Errors

These occur when `inspect`, `search`, or `index` cannot locate a vault or topic.

| Error Message | Cause | Recovery |
|---------------|-------|----------|
| `unable to find a vault from <path>. walked up looking for kb.toml or .kb/vault/` | No `kb.toml` or legacy `.kb/vault/` exists above the working directory | Add `kb.toml` at the vault root, run `kb ingest codebase <path> --topic <topic-id>` to bootstrap a legacy vault, or pass `--vault <path>` |
| `Vault path was not found or is not a directory: <path>` | The `--vault` flag points to a nonexistent path | Verify the vault path exists and is a directory |
| `no topics were found in <path>. expected child directories containing CLAUDE.md` | The vault exists but no configured topic glob matches a directory with `CLAUDE.md` | Run `kb topic new`, ingest a codebase, or update `[vault].topic_globs` in `kb.toml` |
| `multiple topics were found in <path>: <slug1>, <slug2>` | The vault contains more than one topic and no `--topic` flag was provided | Re-run the command with `--topic <topic-id>`; nested topic ids are relative paths like `harness/goclaw` |
| `topic name is required when topic is specified` | The `--topic` flag was provided but with an empty or whitespace-only value | Provide a non-empty topic id |
| `topic "<topic>" is missing CLAUDE.md` | The topic directory is missing the marker/schema file | Create or restore `CLAUDE.md`, or choose a topic listed by `kb topic list` |

## Inspect Lookup Errors

These occur when `inspect symbol`, `inspect file`, `inspect backlinks`, or `inspect deps` cannot resolve the target entity.

| Error Message | Cause | Recovery |
|---------------|-------|----------|
| `no symbols matched "<query>"` | No symbol name contains the query as a case-insensitive substring | Use `kb inspect smells` or `kb inspect complexity` to discover valid symbol names |
| `multiple symbols matched "<query>": <name1>, <name2>` | More than one symbol matched the query | Re-run with a more specific query string |
| `no file matched "<path>"` | No file in the vault has the given `source_path` value | Use the exact source-relative path as stored in vault frontmatter (e.g., `src/config.ts` not `./src/config.ts`) |
| `no symbol or file matched "<query>"` | The query matched neither a file source path nor a symbol name | Re-run with a specific symbol name or an exact source path |

## QMD Errors

These occur when `search` or `index` cannot communicate with the QMD binary.

| Error Message | Cause | Recovery |
|---------------|-------|----------|
| `<command>: QMD is not available to kb. Install it with 'npm install -g @tobilu/qmd' and ensure 'qmd' is on PATH` | The `qmd` binary was not found on the system PATH | Run `npm install -g @tobilu/qmd` and verify with `qmd --version` |
| `<command>: <qmd error details>` | QMD returned an error during execution | Read the stderr diagnostics from QMD for details; common causes include missing collections or corrupted index files |

## Flag Validation Errors

These occur before any command execution when flag combinations are invalid.

| Error Message | Cause | Recovery |
|---------------|-------|----------|
| `ingest codebase: --title and --domain are bootstrap-only and cannot be used when topic "<topic-id>" already exists` | Bootstrap-only metadata flags were used while re-ingesting an existing topic | Remove `--title` / `--domain`, or create a new topic id if you intend a distinct topic |
| `choose at most one search mode flag: --lex or --vec` | Both `--lex` and `--vec` were provided to `search` | Use only one mode selector, or omit both for hybrid mode |
| `--force-embed cannot be used together with --embed=false` | Contradictory embedding flags on `index` | Remove `--force-embed` or set `--embed=true` |
| `--limit must be >= 1. received <N>` | The `--limit` flag on `search` was set to zero or negative | Provide a positive integer for `--limit` |
| `--min-score must be >= 0. received <N>` | The `--min-score` flag on `search` was set to a negative value | Provide a non-negative value for `--min-score` |
| `--top must be >= 1. received <N>` | The `--top` flag on `inspect complexity` was set to zero or negative | Provide a positive integer for `--top` |
| `--min must be >= 0. received <N>` | The `--min` flag on `inspect blast-radius` was set to negative | Provide a non-negative integer for `--min` |
| `invalid --format "<value>": expected one of "table", "json", "tsv"` | An unsupported format string was provided | Use `table`, `json`, or `tsv` |

## KB Workflow Errors

These occur during knowledge base maintenance operations.

| Error | Cause | Recovery |
|-------|-------|----------|
| `kb` not found on PATH | The `kb` binary is not installed or not on PATH | Install the `kb` binary and verify with `kb version` |
| Topic not found | The specified topic id does not exist in the vault or is not matched by `topic_globs` | Run `kb topic list` to see available topics, update `kb.toml`, scaffold a direct topic with `kb topic new <slug> <title> <domain>`, or bootstrap a nested topic that matches `topic_globs` |
| YouTube caption extraction fails before transcript ingest | The required `yt-dlp` backend is missing, outdated, blocked, or failed | Install or update `yt-dlp`, or set `[youtube].yt_dlp_path` / `YOUTUBE_YT_DLP_PATH` to the intended executable |
| YouTube `network_blocked` / blocked captions or audio | YouTube returned 400/403/429/5xx through `yt-dlp`, often due to bot detection, authentication, rate limiting, or datacenter IP blocking | After confirming `yt-dlp` is installed and current, configure `[youtube].proxy`, `[youtube].cookies_file`, `[youtube].user_agent`, or run from a trusted/residential network |
| Instagram caption/transcript extraction fails | The required `yt-dlp` backend is missing, outdated, or `[instagram].yt_dlp_path` points at the wrong executable | Install or update `yt-dlp`, or set `[instagram].yt_dlp_path` to the intended executable (`kb ingest instagram`/`channel` share the same engine as `youtube`) |
| Instagram `network_blocked` / blocked reel | Instagram returned 400/403/429/5xx through `yt-dlp`, often due to rate limiting or a missing session | After confirming `yt-dlp` is current, set `[instagram].cookies_file` (separate from YouTube — Instagram needs its own session) and optionally `[instagram].proxy`, or run from a trusted network |
| Article exceeds 4000 words | A wiki article has grown beyond the recommended length | Extract a sub-topic into its own article and wikilink to it, rather than padding |
| Cross-topic wikilink ambiguity | Two topics contain articles with the same title | Disambiguate with the full path: `[[other-topic/wiki/concepts/Article Name\|Display Name]]` |
| `log.md` missing in existing topic | The topic was created before `log.md` was standard, or it was accidentally deleted | Let the next write operation autocure the skeleton, or create manually and backfill from git: `git log --format='## [%ad] <op> \| %s' --date=short <topic>/` |
| Log entry conflicts with git | Apparent duplication between `log.md` and git history | The log is a human/LLM-readable audit trail, not a replacement for git. Let them coexist: git records *what changed*, `log.md` records *what the knowledge base did* |

## General Errors

| Error Message | Cause | Recovery |
|---------------|-------|----------|
| `a search query is required` | Empty or whitespace-only query passed to `search` | Provide a non-empty search query string |
| `a symbol name is required` | Empty query passed to `inspect symbol` | Provide a non-empty symbol name |
| `a file path is required` | Empty path passed to `inspect file` | Provide a non-empty source path |
| `a symbol name or file path is required` | Empty query passed to `inspect backlinks` or `inspect deps` | Provide a non-empty symbol name or file path |
