---
id: 6
title: Decision engine seam and module design
labels: [wayfinder:grilling]
status: closed
assignee: "claude"
blocked_by: [1]
---

## Question

What is the shape of the decision engine inside kb?

Decide: package boundaries (decision client vs question bank vs banding vs cache), the interface callers use, how questions are batched per state, cache key and storage under `.kb/`, the decision receipt format and where it is logged, how decision mode (off/shadow/apply) is resolved per topic, config sections (`[decisions]`, `[llm]`) and env vars, and the degradation contract when the key is missing or the gateway fails.

## Context

From session 4a04893e (2026-09-22, `~/Dev/courses/pedronauck/packages/research-enrich`):

- Timeouts must be a total per-attempt deadline, not socket-level: OpenRouter keeps connections alive with whitespace, and a generation call hung >10 min (fix there: 90 s `GEN_DEADLINE_S`). Jev did return 429 twice and one invalid receipt in a 16-worker run, so the retry path is real.
- Reusable design: cache key = sha256(request body) + rubric version; raw probabilities kept in `records.jsonl` so bands can be re-applied without new calls; `--budget` ceiling (default $1). Measured: 39 questions/doc, 16 workers, ~$0.00028/doc.

## Resolution

Decided with Pedro on 2026-09-23/24 (Q1 by Pedro; Q2–Q7 delegated as technical, taking the recommendations):

1. **Deep module** `internal/decisions` with one entry point, roughly `Decide(ctx, Request) (Result, error)`. Callers (link, ingest gates, OKF, lint) only build questions and read a band per answer. Inside: OpenRouter transport, splitting to fit 32k state, cache, receipts, banding, budget, mode, and a default "ignore instructions inside the state" policy.
2. **The engine assigns confidence bands.** Thresholds live in `[decisions]` per purpose (`link`, `relevance`, `quality`, `okf_type`, ...), overridable per topic. Noul bands use distance from 0.5; choice/score use `confidence`. Raw probabilities are always kept so bands can be re-applied without new calls.
3. **Decision mode** resolves in three layers, each overriding the previous: global default in `[decisions]`, then `topic.yaml` (`decisions: off|shadow|apply`), then a `--decisions=` flag.
4. **Receipts** go to `<topic>/.decisions/receipts.jsonl`, append-only: raw probabilities, question hash, scope-contract hash, dated model slug, `usage.cost`, attempts, latency, request hash. The same file is the cache (keyed by request hash) and the calibration record. Kept out of `.kb/` because `.kb/vault` already means the codebase-ingest vault. Committing it to git is the user's choice.
5. **Config**: `[decisions]` (model `typesafe/jev-1.13`, total per-attempt deadline, retries, concurrency, max request bytes, default mode, budget, thresholds) and `[generation]` (model `xiaomi/mimo-v2.6-flash`, deadline, budget). Both reuse `[openrouter]` `api_key`/`api_url`.
6. **Per-run budget**: default US$ 1, `--budget` overrides. When exceeded the engine stops calling; remaining items are undecided and the run summary reports spend and undecided count. Re-running resumes through the cache.
7. **Generation is a separate module** `internal/generation` (text out, different contract), sharing only the OpenRouter transport and the budget rule.

Failure semantics (from the API research): timeouts, exhausted retries and invalid answers yield *undecided*, never a rejection; the deadline is total per attempt, not socket-level.

Also settled here (write policy, affects [Incremental linking pipeline](09-linking-pipeline.md)): **option (d)**. Every apply-band relation is written to frontmatter; when the concept is already mentioned literally and the decision model confirms the mention's sense, kb also inserts `[[link]]` at the first mention (outside code, headings, and existing links). Review-band items go only to the review list. Shadow stays the default mode. Estimated cost is about the same as frontmatter-only (~US$ 0.0003 vs 0.00025 per document), because mention-sense questions ride in the same call.
