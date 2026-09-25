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

## Decision Model Errors

These occur in commands that use the decision or generation model (`ingest` except `codebase`, `classify`, `link`, `find`, `review accept|reject`, `topic contract --draft|--accept`, `topic vocabulary`, `promote` with `[okf].types`, `okf check --decide`).

| Error Message / Status | Cause | Recovery |
|------------------------|-------|----------|
| `<command>: ... OPENROUTER_API_KEY is not set ...` | The decision model is required and not configured | Set `OPENROUTER_API_KEY` or `[openrouter].api_key` |
| `--decisions must be shadow or apply` | Invalid mode override | Pass `--decisions shadow` or `--decisions apply` |
| HTTP 401 / 402 / 403 from OpenRouter | Bad key, no credit, or no access to the model | Fix the key or account; these stop the run instead of retrying |
| `undecided:budget` in the run summary | The run reached `--budget` / `[decisions].budget_usd` | Re-run; cached answers are reused and the run continues |
| `undecided:timeout`, `undecided:retries`, `undecided:invalid_receipt` | Transport or receipt failure after retries | Re-run later; the document keeps its previous state and is never treated as a "no" |
| `not_checked:state_too_large` | The document state exceeds `[decisions].max_state_bytes` | Nothing to do per run; the document is listed and skipped |
| `skipped:user-key` / `key-conflict` | A kb-owned key holds a value kb did not write | Keep it (kb leaves it) or delete it to let kb write it |
| `skipped:changed` | The file changed while kb was writing | Re-run |
| `skipped:locked` | The file has `locked: true` | Remove `locked` to let kb write it |
| `topic contract: pass exactly one of --draft, --import-claude or --accept` | Zero or several actions | Run one action per call |
| Self-check conflicts on `--accept` | A kept line is excluded by an `out_of_scope` line | Qualify the `out_of_scope` line (`selection-contract.md`); `--force` only for a false alarm |
| `--budget must be a finite amount in US$ >= 0` / `invalid amount` | Negative or non-numeric `--budget` | Pass a non-negative amount; `--budget 0` runs on cached answers only |
| `not_checked:excluded` / `decisions.exclude: kept without a judgment` | The file matches `topic.yaml` `decisions.exclude` | Intended: excluded files never reach a model. Remove the glob to have them judged |
| `triage: review` with `triage_reason: undecided` after an ingest | The gate's judgment failed (timeout, invalid receipt, budget) | Re-run the ingest, or decide the `gate` review item yourself; the source was never treated as kept or rejected |
| `<command>: invalid topic question bank under <topic>/.decisions/banks` | A topic extra question bank fails the schema check | Fix the file (`decision-workflow.md`, Topic extra questions) or move it out of `banks/` |
| `--decisions=apply ignored for relevance: ...` in the run summary | No accepted contract, or `decisions.relevance: off` | Accept a contract first; relevance cannot be applied without one |
| `calibration is stale (<contract, model or bank> changed)` in the run summary | The stored calibration was made under another contract, model or question bank | Re-judge with `kb classify`, then `kb review calibrate <topic> --write` with the user's approval; relevance gates stay in shadow until then |
| `skipped <id>: already resolved by <id>` during a bulk `kb review accept\|reject` | An earlier item in the same run closed this one (same source) | Nothing to do |
| `actions: required judgment is undecided` on a `recapture` accept | The fresh page could not be judged | Re-run the accept later; the body, metadata and item stay unchanged |
| Restore prints manual repair entries (file, line or key, original text) | A file touched by the quarantine changed since | Re-insert the printed text by hand where it belongs; kb never forces it |

## OKF Errors

These occur in `kb promote` and `kb okf check`.

| Error Message | Cause | Recovery |
|---------------|-------|----------|
| `promote: target topic must use mode okf` | `--to` points at a wiki topic | Pass an existing `mode: okf` topic, or create one with `kb topic new <slug> <title> <domain> --mode okf` |
| `required flag(s) "to" not set` | `--to` was omitted | Pass `--to <okf-topic>` |
| `--type is required` | No `[okf].types` vocabulary, or the type suggestion was not confident (the error lists the top candidates) | Pass `--type <Type>`, ideally from `[okf].types` |
| `promote: source document not found: <path>` | The source path does not resolve | Use a vault-relative path to an existing wiki document, e.g. `<topic>/wiki/concepts/<Article>.md` |
| `okf check` exits non-zero with diagnostics | `severity=error` concepts (missing/empty `type`, bad frontmatter); under `--strict` also warnings (missing producer fields, off-vocabulary types) | Fix the listed concepts; see `okf-mode.md` |

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
| `warning: qmd collection "<name>" uses pattern "...", not kb's "..."` on `kb index` | The collection was created before kb added its mask, so quarantined files stay indexed (kb filters them from results) | `qmd collection remove <name>`, then `kb index --topic <topic>` |

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
