# Decision Workflow Reference

How kb's decision model gates ingests, classifies, links, finds and queues work for review, and what each command reads and writes. Read the section you need: **Commands** for flags, **Modes and thresholds** before changing `topic.yaml` `decisions`, **Review queues** before accepting or rejecting items, **Calibration** before `kb review calibrate --write`, **Files** when inspecting `<topic>/.decisions/`.

The decision model (Jev, `typesafe/jev-1.13` over OpenRouter) answers typed questions (yes/no, choice, score) with probabilities and writes no text. The generation model writes only short, code-validated literals (summaries, criteria, aliases, entities, questions, contract and vocabulary drafts). Code proposes candidates, applies thresholds and writes files. You write article prose and give the human verdicts.

## Commands

```bash
kb topic contract <topic> --import-claude            # CLAUDE.md "## Selection contract" → contract_draft (no model)
kb topic contract <topic> --draft [--budget <usd>]   # generated draft → contract_draft
kb topic contract <topic> --accept [--force] [--yes] [--budget <usd>]   # self-check + impact preview, then activate
kb topic vocabulary <topic> --draft | --accept       # 12–40 concepts for a topic without articles; accept creates stubs
kb classify <topic> [--only-missing] [--all] [--budget <usd>] [--decisions shadow|apply] [--format table|json|tsv]
kb link <topic> [--all] [--dry-run] [--no-qmd] [--budget <usd>] [--decisions shadow|apply]
kb find <topic> "<question>" [--limit N] [--kind <genre>] [--concept <title>] [--min-depth N] [--relevance <role>] [--explain] [--json] [--no-qmd] [--budget <usd>]
kb find --facets <topic> [--json]                    # counts per genre, relevance, depth, concepts (no model)
kb review <topic> | kb review list <topic> [--queue <queue>] [--format table|json]          # no model
kb review accept|reject --topic <topic> <id>... [--queue q] [--all-purpose p] [--above p] [--topic-wide] [--budget <usd>]
kb review import-links <topic>                       # no model
kb review import-labels <topic> --from <file> --id-field <f> --match url|path|doi|pmcid --decision-field <f> --keep-values include,true [--decided-by-field f]   # no model
kb review calibrate <topic> [--write] [--format]     # no model
```

- **Requirement.** Every command above that is not marked "no model", plus every `kb ingest` except `codebase`, `kb promote` with a non-empty `[okf].types` and `kb okf check --decide`, needs `OPENROUTER_API_KEY` (or `[openrouter].api_key`) and refuses to start without it.
- **Budget.** `--budget <usd>` caps the run (default `[decisions].budget_usd`, US$ 1). `--budget 0` is cache-only: cached answers are used and no call is made. A negative or non-numeric value is rejected before anything runs. Past the budget, remaining items are `undecided:budget`; re-running resumes through the cache.
- **Failures are never a "no".** `undecided:<reason>` (timeout, retries, invalid_receipt, budget) and `not_checked:<reason>` (state_too_large, excluded) keep the document's previous state and are listed in the run summary with coverage. A 401/402/403 stops the run.
- **Incremental.** `kb classify` and `kb link` skip documents whose body, contract and question-bank versions are unchanged since they were judged. `--all` re-judges everything; `kb classify --only-missing` judges only documents without a record or with a missing facet key.
- **`kb link --dry-run`** writes nothing but still asks decisions: cached answers are free, uncached ones spend the budget.
- **Reverse pass.** `kb classify` (for articles that gained aliases) and `kb topic vocabulary --accept` (for new stubs) scan every document for mentions of those articles and judge them with the article as the only link candidate.
- **qmd candidates.** When `qmd` is on PATH and the topic's collection is fresh, `kb link` and `kb find` add vector neighbours as candidates; `--no-qmd` turns that off. qmd scores never decide anything.

## Selection contract

Relevance is judged against `topic.yaml` `contract:` (purpose, core, adjacent, collected_on_purpose, collected_on_purpose_paths, out_of_scope), rendered into the topic `CLAUDE.md`. Drafting rules and the accept workflow: `selection-contract.md`.

- `--accept` blocks on self-check conflicts (a kept line that an `out_of_scope` line would exclude, P ≥ `contract_conflict`) unless `--force`, runs the impact preview over every source (a stratified sample of 500 when larger), and activates only after the topic slug is typed (or `--yes`). Show the preview to the user; only they decide whether the contract describes the collection.
- The preview's receipts are reused by the next `kb classify` when the draft is accepted unchanged.
- Documents under `collected_on_purpose_paths` get `relevance: collected_on_purpose` from code and are never quarantined for relevance.

## Ingest gates

Every `kb ingest` except `codebase`: (1) exact dedupe by normalized URL, platform id and body hash, including quarantined sources; (2) pre-fetch relevance for `channel` and URL lists (`kb ingest url <url>...` or `--from <file>`), skipping items at P(off_topic) ≥ `relevance_quarantine` into `.decisions/skipped.jsonl` (rescue with `--rescue <id>`); (3) code quality flags `not_an_article`, `thin`, `error_page`; (4) one fresh Firecrawl refetch of a flagged or short URL capture, keeping the longer body; (5) quality and relevance judgment; (6) near-duplicate (`same_content` quarantines as `duplicate`, `new_version` writes `supersedes`); (7) write, classify and link in the same run. `kb ingest bookmarks` ingests a cluster file as one source and gets stages 1 and 3–7.

- **Outcomes.** `triage: kept`, `review` (review band, with a `gate` review item) or `quarantined` (moved to `raw/_quarantine/`). When the judgment itself fails (undecided), the source is written `triage: review` with `triage_reason: undecided`, never kept silently.
- **`--force`** skips stages 1–6; classify and link still run.
- **`log.md`** gets one entry per run with kept, review, quarantined, skipped, duplicate, undecided counts and the cost.

## Modes and thresholds

A **shadow** outcome is recorded (receipts, run summary, a review item) and changes nothing; **apply** acts.

| What | Default | Becomes `apply` when | One-run override |
|---|---|---|---|
| Exact dedupe and quality gates | apply | always | `--decisions shadow` |
| Relevance gates (skip, quarantine) | shadow | an accepted contract **and** either a current calibration with ≥ 30 relevance labels given in `kb review`, or `decisions.gates: apply` | `--decisions shadow\|apply` (apply is ignored for relevance without an accepted contract or with `relevance: off`) |
| Body-link insertion by `kb link` | shadow | `[decisions].mode`, then `topic.yaml` `decisions.mode`, is `apply` | `--decisions shadow\|apply` |

Frontmatter relations are written in both modes; `contradicts` never auto-applies. The run summary states each gate's mode and why.

`topic.yaml`:

```yaml
decisions:
  mode: shadow                 # body-link insertion
  gates: shadow                # relevance gates (shadow until calibrated)
  relevance: on                # off for topics defined by construction (one note per entity)
  thresholds: {}               # named overrides; `kb review calibrate --write` fills them
  exclude: ["raw/private/**"]  # topic-relative globs never sent to any model
```

- **`exclude`** keeps matching files out of every decision and generation call (`not_checked:excluded`). An ingested file on an excluded path is written `triage: kept` with a note, and is not classified, linked or used as a near-duplicate candidate.
- **Thresholds** (defaults, `[decisions.thresholds]`, per-topic overrides above): `relevance_quarantine` 0.80, `relevance_review` 0.50, `relevance_fetch` 0.50, `quality_apply` 0.80, `quality_review` 0.50, `duplicate` 0.80, `link_apply` 0.85, `link_review` 0.60, `mention_sense` 0.85, `affects` 0.80, `primary_concept` 0.30, `concept_noul` 0.70, `find_keep` 0.50, `find_window` 0.35, `okf_type` 0.80, `contract_conflict` 0.70.

## Classify, link and find outputs

- **`kb classify`** writes `genre`, `depth` (only when there is enough text), `relevance`, `quality` (when a check fires), `concepts`, `summary`, `entities`, `questions` on sources; `genre`, `concepts`, `summary`, `criterion`, `aliases` on articles. Stub articles (`stage: stub`) get aliases only. A generated literal that fails validation is not written and the document stays pending for the next run. Classification never quarantines: broken captures go to `recapture`, off-topic sources (P(off_topic) ≥ `relevance_review`, no quality flag) to `remove`. When ≥ 5 sources across the topic match no article, it queues up to 5 `concept-proposal` items.
- **`kb link`** writes relation lists (`related`, `extends`, `prerequisite`, `example_of`) and `affects` in both modes, inserts `[[target|matched text]]` at the first confident mention only in apply mode (logged in `inserted-links.jsonl`), queues review-band links, contradictions and shadow insertions, and never removes a link. A kb-written relation that drops below the apply band becomes a `link` review item, not a deletion.
- **`kb find`** ranks candidates by P(answers) × specificity × genre match × depth, keeps P ≥ `find_keep` within `find_window` of the best, and names every drop (`low-probability`, `below-window`, `limit`, `filtered-facet`, `quarantined`, `undecided:<reason>`) under `--explain`.

## Review queues

| Queue | Accept | Reject |
|---|---|---|
| `gate` | keep the source (restores it when quarantined) | quarantine it with the item's reason |
| `skip` | rescue: fetch and ingest the skipped item | dismiss |
| `recapture` | refetch fresh and re-judge: a longer, clean page replaces the body in place (`recaptured: <date>`); still broken or no `source_url` → quarantine; an undecided re-judgment leaves the item pending | keep as is |
| `remove` | quarantine as `off_topic` | keep |
| `link` | write the relation (and insert the body link for insertion items; remove the kb-written value for demotion items) | nothing |
| `contradiction` | write `contradicts` | nothing |
| `concept-proposal` | create a stub article with its criterion and run the reverse pass | nothing |
| `okf-type` | record the verdict; re-run `kb promote --type` | record the verdict |

- Every accept and reject becomes a **label** in `labels.jsonl` (positive = keep/link, negative = remove/broken). A failed action writes no label and leaves the item pending.
- Bulk selection (`--queue`, `--all-purpose`, `--above`, `--topic-wide`) skips items an earlier action in the same run already closed, and stops at the first authentication error.
- Read the items before bulk-accepting. Several wrong `remove` items mean the contract is too narrow: fix the contract instead of rejecting them one by one.
- Restoring a quarantined source replays the quarantine ledger; an entry whose file changed since is printed (file, line or key, original text) for manual repair instead of being forced.

## Calibration

- `kb review calibrate <topic>` reports precision, recall, coverage and review rate for `relevance`, `link` and `quality`, from labels only, at the current thresholds and over a sweep.
- Labels are split deterministically: `sha256(subject) mod 10` buckets 0–6 are **dev** (thresholds are chosen there), 7–9 **holdout** (only reported). Labels shown in an impact preview become dev-only once that contract changes.
- Labels join only scores from the **active decision context** (contract hash, decision model, question-bank version). After a contract, model or bank change, the stored calibration is **stale**: relevance gates return to shadow and the run summary says `calibration is stale (... changed)`. Re-judge (`kb classify`) and re-run `kb review calibrate --write`.
- A purpose needs ≥ 30 dev labels given in `kb review` and ≥ 10 holdout labels; imported labels are reported separately and never produce thresholds on their own.
- Show the dev and holdout numbers to the user and `--write` only with their approval: it stores thresholds in `topic.yaml` `decisions.thresholds` and the record in `calibration.json`.

## Topic extra questions

`<topic>/.decisions/banks/*.json` adds questions to the built-in requests of one purpose, without replacing any built-in question; their answers are recorded in receipts only. Same schema as the built-in banks: `id` (lowercase snake_case), `version`, `purpose` (`relevance`, `quality`, `classify`, `link`, `find`, `okf_type`, `duplicate`), `guard`, and `questions` with `id`, `type` (`noul`, `choice`, `score`), `instructions`, `criteria`, `source`. Use `{id}` in a question id and its text to ask it once per candidate. An invalid bank stops the command with `invalid topic question bank under ...`.

## Files under `<topic>/.decisions/`

Records and caches; the graph is the markdown. JSONL files are append-only.

| File | Holds |
|---|---|
| `receipts.jsonl` | every decision and generation call: purpose, subject, bank version, contract hash, model, cost, latency, status, raw answers (also the cache) |
| `state.jsonl` | per document: body hash, contract, bank versions, value hashes of keys kb wrote, and the body each bank and key was judged on |
| `review.jsonl`, `labels.jsonl` | review items with their status changes; human and imported labels with their origin |
| `calibration.json` | last calibration: thresholds, label counts, dev/holdout metrics, decision context per purpose |
| `previews.jsonl` | documents shown in each impact preview |
| `skipped.jsonl` | bulk items skipped before fetch (rescuable) |
| `quarantine-ledger.jsonl` | every quarantine move and removed index/frontmatter line, replayed on restore |
| `inserted-links.jsonl` | body links inserted by `kb link` |
| `banks/*.json` | optional topic extra questions |

Do not edit these files by hand; delete `receipts.jsonl` only to force every decision to be asked again (it costs the full run).

## Configuration

`[decisions]` (`model`, `deadline`, `retries`, `concurrency`, `max_state_bytes`, `budget_usd`, `mode`, `[decisions.thresholds]`), `[generation]` (`model`, `fallback_model`, `deadline`, `retries`, `summary_language`), `[firecrawl]` refetch keys (`refetch_wait_ms`, `refetch_only_main_content`) and `[okf].type_descriptions` in `kb.toml`; env `OPENROUTER_API_KEY` (required), `KB_DECISIONS_MODEL`, `KB_GENERATION_MODEL`. `config.example.toml` in the kb repo documents every key. The first call in a topic prints a privacy notice: judged excerpts are sent to OpenRouter.
