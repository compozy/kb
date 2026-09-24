# kb-jev implementation plan (spec revision 3)

Source of truth: `.wayfinder/kb-jev/spec.md` (rev 3). This file fixes the cross-package
contracts so parallel slices fit together. When this file and the spec disagree, the spec
wins; fix this file.

**Never reintroduce an `ai_` prefix anywhere** (keys, code, docs, tests).

## Waves

- W0 (root): this plan, shared skeleton types (`internal/decisions/types.go`, `internal/questions/q.go`).
- W1 (parallel worktrees):
  - A `decisions + generation + config + questions (loader + all banks)`
  - B `frontmatter editor + resolve (extracted from lint) + corpus (loader/BM25/mentions/excerpt/provenance/state store/writer)`
  - C `contract + topic.yaml settings + topic templates + ingest provenance (ingest_batch/ingest_query, owned-key guard)`
- W2 (root): `internal/session`, integrate W1, `make verify`.
- W3 (parallel worktrees):
  - D `review (queue/labels/import-labels/import-links/calibrate) + refs (quarantine/restore/ledger) + firecrawl freshness`
  - E `classify + criterion/aliases/concept proposals + topic vocabulary + topic contract CLI (draft/import/accept+preview)`
  - F `link + find`
  - G `lint new kinds + okf/promote changes`
- W4 (parallel): H `gate + ingest wiring + review actions + review CLI + ingest flags`, I `docs/skills/templates/AGENTS/CLAUDE/README/config.example`.
- W5 (root): integration tests pass, real run on a topic copy, Codex review, PR, CI, release.

## Package map (new)

| Package | Owns |
|---|---|
| `internal/decisions` | `Engine.Decide`, transport, receipt validation, receipts+cache, banding, budget, concurrency, redaction, run summary, `OwnedKeys` |
| `internal/generation` | `Client.Generate` over chat completions (json_schema strict, reasoning disabled), fallback model, shared budget+receipts |
| `internal/questions` | embedded banks `banks/*.json`, loader (guard, schema check, hash), templating `{id}`, topic extra banks |
| `internal/resolve` | vault link index + Obsidian-style resolution + link extraction (wikilinks, markdown links, frontmatter links); shared by lint/refs/link |
| `internal/corpus` | topic document loader, body hash, BM25, alias dictionary + mention scan, section-aware excerpt, provenance state, `.decisions/state.jsonl`, byte-preserving owned-key writer |
| `internal/contract` | selection contract parse/validate/hash/render, CLAUDE.md import, topic.yaml settings (contract, drafts, decisions block) read + node-preserving write, path globs (`**`) |
| `internal/session` | per-run bundle: config, topic, settings, contract, engine, generation, corpus (lazy), state store, writer, review queue, modes |
| `internal/review` | queue store, labels, import-labels, import-links, calibrate (holdout), previews ledger |
| `internal/refs` | inbound-reference scan, quarantine/restore, `quarantine-ledger.jsonl` |
| `internal/classify` | facets pass, criterion, aliases, entities/questions, concept proposals, vocabulary draft/accept |
| `internal/link` | candidates, judgment, write policy (d), reverse pass |
| `internal/find` | retrieval + facets |
| `internal/gate` | dedupe, prefetch relevance, code quality, refetch, post-fetch quality+relevance, near-duplicate |
| `internal/actions` | review accept/reject dispatcher (imports refs/link/gate/classify) |

Import direction (no cycles): `questions, resolve` ← `decisions` ← `generation`; `frontmatter, resolve, decisions` ← `corpus`; `contract` (imports only frontmatter/yaml); `session` ← {decisions, generation, contract, corpus, review}; `review` imports corpus/decisions/contract only; `refs` imports resolve/corpus; `classify, link, find, gate` import session; `actions` imports all of them; `ingest` imports session/gate/classify/link; `cli` imports everything.

## Topic files under `<topic>/.decisions/`

All JSONL is append-only, one JSON object per line, UTF-8, `\n`. Readers take the last row per key.

| File | Row | Owner |
|---|---|---|
| `receipts.jsonl` | `{key, time, purpose, subject, bank, bank_version, bank_hash, contract, model, reported_model, route, attempts, latency_ms, input_tokens, cost, cost_unknown, status, reason, questions:[ids], answers:{id: raw answer}}` | decisions / generation (`purpose: "generate:<kind>"`, `answers: {"output": <json>}`) |
| `state.jsonl` | `{path, body_hash, contract, banks:{bank_id: version}, written:{key: value_hash}, updated}` | corpus |
| `skipped.jsonl` | `{id, time, batch, query, url, title, description, channel, date, p_off_topic, receipt_key, rescued:false}` | gate |
| `quarantine-ledger.jsonl` | `{quarantine_id, time, op:"move"|"edit"|"restore", file, original, line_or_key, before, after, hash_before, hash_after, reason}` | refs |
| `review.jsonl` | review `Item` (below); status changes are new rows with the same id | review |
| `labels.jsonl` | `{time, subject, purpose, question, verdict:"positive"|"negative", receipt_key, origin:"review"|"import:<file>@<sha256>"|"import-links", decided_by, item_id}` | review |
| `previews.jsonl` | `{time, contract, subjects:[paths shown or used in precision]}` | contract CLI (E) |
| `calibration.json` | `{time, purposes:{<purpose>:{dev_labels, holdout_labels, thresholds:{...}, dev:{precision,recall,coverage,review_rate}, holdout:{...}}}}` | review |
| `banks/*.json` | user extra question banks (same schema as built-in, `purpose` required) | user |

Quarantine location: `<topic>/raw/_quarantine/<original topic-relative path without leading raw/>` (so `raw/articles/x.md` → `raw/_quarantine/articles/x.md`). Everything under `raw/_quarantine/` and `.decisions/` is excluded by `corpus.Load`.

## Owned frontmatter keys (spec §6)

`decisions.OwnedKeys` = `triage, triage_reason, genre, depth, relevance, concepts, summary, criterion, entities, questions, quality, ingest_batch, ingest_query, related, extends, prerequisite, example_of, contradicts, affects, supersedes, aliases, recaptured`. `locked` is user-owned (never written). Relation keys: `related, extends, prerequisite, example_of, contradicts`. `aliases` is special: kb merges (never removes), so the value-hash conflict rule does not block merging new aliases into user aliases.

Value hash = sha256 over canonical JSON of the normalized Go value (strings trimmed, lists in written order). Body hash = sha256 of body bytes after frontmatter with `\r\n` → `\n`.

## Core types (fixed; see `internal/decisions/types.go`, `internal/questions/q.go`)

```go
// questions
type Q struct { ID string; Type string /*noul|choice|score*/; Instructions any; Criteria any }
type Bank struct { ID, Version, Hash, Purpose, Guard string; ... }
func Load(id string) (*Bank, error)            // embedded
func (b *Bank) Question(id string, vars map[string]string) (Q, error) // {placeholders} substituted, guard applied
func (b *Bank) WithCriteria(q Q, criteria any) Q

// decisions
type Purpose string // relevance quality duplicate classify concept link mention find okf_type contract_check
type Status string  // decided undecided not_checked
type Band string    // apply review ignore
type TopicRef struct { Slug, Root string; Contract string /*hash*/; Thresholds Thresholds }
type Request struct { Topic TopicRef; Purpose Purpose; Subject string; Bank *questions.Bank /* id+version+hash for key & receipt */; State any; Questions []questions.Q }
type Answer struct { Status Status; Reason string; Band Band; Noul *float64; Choice string; Score *float64; Probs map[string]float64; Confidence *float64; Type string }
type Result struct { Answers map[string]Answer; ReceiptID string; CostUSD float64; CacheHit bool }
func (e *Engine) Decide(ctx, Request) (Result, error) // error only for fatal (auth 401/402/403, ctx cancel); everything else is per-answer status
```

Band per answer is computed by the engine from `Thresholds.For(purpose)` (apply, review pair): noul → P(yes); choice → confidence; score → confidence. Callers that need another gate quantity (sum of option probs, argmax prob, P(off_topic)) use `decisions.BandFor(value, apply, review)` with a named threshold. Thresholds are recomputed from receipts on every read (no call).

Default thresholds (`[decisions.thresholds]`, overridable by `topic.yaml decisions.thresholds`):
`link_apply 0.85, link_review 0.60, mention_sense 0.85, affects 0.80, relevance_quarantine 0.80, relevance_review 0.50, relevance_fetch 0.50, quality_apply 0.80, quality_review 0.50, duplicate 0.80, primary_concept 0.30, concept_noul 0.70, find_keep 0.50, find_window 0.35, okf_type 0.80, contract_conflict 0.70`.

## Modes

- Body-link insertion mode: `[decisions].mode` (default shadow) → `topic.yaml decisions.mode` → `--decisions` flag.
- Relevance gate mode: shadow unless (contract accepted AND (`calibration.json` has ≥30 relevance labels OR `decisions.gates: apply`)); `--decisions` overrides for the run.
- Quality gates + exact dedupe: apply by default; `--decisions=shadow` puts them in shadow for that run.
- `decisions.relevance: off` disables relevance judgments (quality/classification still run).
- Shadow never moves, skips or quarantines; would-be outcomes go to the run summary and receipts.

## Review queue item

```go
type Item struct {
  ID string        // "rv-" + first 10 hex of sha256(queue|subject|target|question)
  Queue string     // gate | skip | link | contradiction | concept-proposal | okf-type | recapture | remove
  Purpose string   // decisions purpose
  Subject string   // topic-relative path, skipped id, or query id
  Target string    // candidate path / concept title, optional
  Question string  // question id
  Probability float64
  ReceiptKey string
  Evidence string  // one line for humans
  Action map[string]any // payload the actions dispatcher needs (relation, mention, reason, ...)
  Status string    // pending | accepted | rejected
  Created, Resolved string
}
```

## CLI surface (new/changed)

- `kb classify <topic> [--only-missing] [--budget] [--decisions]`
- `kb link <topic> [--all] [--dry-run] [--budget] [--decisions]`
- `kb find <topic> "<question>" [--limit 10] [--kind] [--concept] [--min-depth] [--relevance] [--json] [--explain]`; `kb find --facets <topic>`
- `kb review <topic>` = `kb review list <topic> [--queue]`; `kb review accept|reject --topic <t> <id>... [--queue q] [--all-purpose p] [--above p] [--topic-wide]`; `kb review import-links <topic>`; `kb review import-labels <topic> --from --id-field --match url|path|doi|pmcid --decision-field --keep-values`; `kb review calibrate <topic> [--write]`
- `kb topic contract <topic> --draft|--import-claude|--accept [--force] [--yes]`
- `kb topic vocabulary <topic> --draft|--accept`
- every `kb ingest` (except `codebase`): `--budget --force --decisions=shadow|apply --batch <name>`; bulk (`channel`, `bookmarks`): `--rescue <id>`
- `kb promote --type` optional (suggested); `kb okf check` advisory findings; `kb lint` new kinds (never calls a model)

Decision-model requirement: `session.Open` fails fast with `decisions: OPENROUTER_API_KEY is not set; kb <cmd> requires the decision model ([decisions] + [openrouter])`. `version`, `topic list|info|new`, `help`, `inspect`, `lint`, `index`, `search`, `ingest codebase`, `generate`, `migrate` never need it.

## Conventions for every slice

- Table-driven tests, `t.TempDir()`, fakes only at HTTP boundaries (`httptest`). Integration tests behind `//go:build integration` next to the package.
- Use `internal/frontmatter` helpers; never hand-assemble YAML except inside the byte-preserving editor.
- Keep Cobra thin; logic in packages.
- `make lint` zero findings (golangci v2 + modernize). Run `go vet ./...` and `go test ./internal/<pkg>/...` for your packages before finishing.
- Commit in your worktree with a conventional message; do not touch files outside your claim.

## W3/W4 contracts (added after W1/W2 merged at f606a3d)

Revised waves: W3 = D (review import/calibrate/list CLI + refs + firecrawl freshness), E (classify + quality code checks + criterion/aliases/entities/questions + concept proposals + vocabulary + `kb classify` + `kb topic vocabulary`), F (link + find + CLIs), G (lint kinds + okf/promote). W4 = H (gate + ingest wiring + actions dispatcher + `kb review accept|reject`), J (`kb topic contract --draft/--import-claude/--accept` with self-check + impact preview), I (docs/skills/templates).

Already in tree: `internal/session` (Open, modes, Corpus, StateMeta, Request, WriteSummary, PrivacyNotice, calibration.json schema + ReadCalibration), `internal/fakes` (fake OpenRouter decisions+generation, fake Firecrawl), `internal/review/store.go` (Item, Label, Store: Add/Resolve/Get/Items/Pending/PendingCounts/AddLabel/Labels, LoadLabels), `internal/cli/session.go` (`openSession`, `bindDecisionFlags`).

- **Gate judgment (one shared request).** `classify.JudgeGate(ctx, s, doc, opts)` asks ONE request = composite of banks `relevance` + `quality` over state `{contract: contract.QuestionState(), document: corpus.DocumentState(doc, 12000, nil)}` with questions `role` (omitted when relevance is off, the path matches `collected_on_purpose_paths`, or there is no contract → then relevance is shadow-only... still ask `role` without contract? NO: without a contract, ask `role` with `contract: {}` only when relevance is on; record in shadow) + the four quality nouls (`no_speech_content` only for `youtube-transcript`/`instagram-video` source kinds). Gate stage 5 (ingest), the impact preview (`kb topic contract --accept`) and `kb classify` all call it so receipts are shared by cache key. Composite banks: `questions.Compose(banks ...*Bank) *Bank` (id = ids joined by "+", version = versions joined by "+", hash = sha256 of hashes).
- **Bank ids recorded in state rows** (`StateMeta.Banks`): sources `{relevance, quality, classify, concept}`, articles `{classify, concept}`; link adds `link` and `link_candidates` (hash of sorted candidate ids). Lint's `unclassified` compares these with `questions.Load(id).Version`.
- **Code quality checks** live in `internal/quality` (E): `Check(Input{Title, Body, SourceURL, FinalURL, RequestedURL, StatusCode, SiteName, HostLines map[string]int}) []Flag` with flags `not_an_article`, `thin`, `error_page` (spec §7 stage 3). Gate and classify both use it.
- **Inserted links log.** `link` appends `{time, subject, target, text, mode}` to `<topic>/.decisions/inserted-links.jsonl` for every body link it inserts; `kb review import-links` treats every other body/frontmatter link between topic documents as human-written.
- **Stub articles.** Vocabulary accept and concept-proposal accept create `wiki/concepts/<Title>.md` with frontmatter `title, type: wiki, stage: stub, domain, tags: [<domain>, wiki, stub], created, updated, sources: [], criterion`. Lint accepts `stage: stub|compiled` for wiki concepts.
- **Relation target form.** Quoted wikilinks `"[[<stem>]]"`; when the stem is ambiguous in the vault index use the vault-relative path without `.md`. Body insertion writes `[[<stem>|<matched text>]]`.
- **Queue items by producer.** gate: `gate` (review band), `skip` (pre-fetch skip, Action `{skipped_id}`); classify: `recapture` (code flag or quality noul ≥ quality_review; Action `{reason}`), `remove` (P(off_topic) ≥ relevance_review and no quality reason), `concept-proposal` (Action `{title, criterion}`); link: `link` (review band or shadow body insertion; Action `{relation, target, insert, mention_text, mention_start}`), `contradiction`; okf: `okf-type`.
- **contract_check state** (changed after the live run): `{pairs:[{n, kept, excluded}]}` — `kept` = core/adjacent line, `excluded` = out_of_scope line. The old `a`/`b` wording was inverted by Jev on 8/8 live cases.
- Live regression (`-tags jevlive`): 83 cases over 9 banks pass on typesafe/jev-1.13 (≈US$0.003 total).
