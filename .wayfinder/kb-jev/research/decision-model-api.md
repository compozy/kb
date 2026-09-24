# Decision model API contract (Jev via OpenRouter Decisions)

Research for [ticket 01](../tickets/01-decision-model-api.md). Captured 2026-09-22. No live call was made:
no `OPENROUTER_API_KEY` was available in the environment or in kb config. The response shape below is
confirmed from **recorded real OpenRouter responses** stored by local packages (runs from 2026-09-18/19).

## Sources

- [S1] OpenRouter API reference, "Submit a Decisions request" — https://openrouter.ai/docs/api/api-reference/alphadecisions/submit-a-decisions-questions-and-answers-request (local copy: `~/Dev/courses/pedronauck/research/jev/raw/articles/openrouter-decisions-api-endpoint-de-perguntas-tipadas-jev.md`, captured 2026-09-18)
- [S2] TypeSafe HTTP API reference — https://docs.typesafe.ai/api.md
- [S3] TypeSafe primitives — https://docs.typesafe.ai/primitives.md, `/primitives/choice.md`, `/primitives/score.md`, `/primitives/noul.md`, `/primitives/advanced.md`
- [S4] TypeSafe Models page — https://docs.typesafe.ai/models.md
- [S5] TypeSafe Confidence + confidence-gated routing — https://docs.typesafe.ai/confidence.md, `/patterns/confidence-routing.md`
- [S6] Jev 1.13 jaggedness — https://docs.typesafe.ai/model-jaggedness/jev-1.13.md (local: `raw/articles/jev-1-13-md.md`)
- [S7] TypeSafe Python SDK retries/constants — https://docs.typesafe.ai/sdk/python/api/retries.md, `/sdk/python/api/constants.md`
- [S8] OpenRouter model page (embedded JSON) — https://openrouter.ai/typesafe/jev-1.13 (fetched 2026-09-22)
- [S9] Real Go client + validator: `~/Dev/courses/pedronauck/packages/jev-review/jev.go`; receipts in `jev-review/reports/**` (6,234 requests); `jev-review/VALIDATION.md`
- [S10] Real TS client + stored request/response pairs: `~/Dev/courses/pedronauck/packages/yc-apply-classifier/src/server/provider.ts`, `.data/runs/**` (196 stored runs), `docs/VALIDATION.md`
- [S11] Python clients: `packages/licitacao-fit/classify.py` + `README.md` (10,840-call scale test); `packages/reddit-idea-classifier/jev.py` + `README.md`
- [S12] Local KB wiki: `research/jev/wiki/concepts/Código faz conta, Jev julga texto.md` ("Números de operação", "Pegadinhas de API")

## Endpoint

```
POST https://openrouter.ai/api/alpha/decisions
Authorization: Bearer $OPENROUTER_API_KEY
Content-Type: application/json
```

- Every working local client uses `https://openrouter.ai/api/alpha/decisions` [S9, S10, S11]. kb's existing `[openrouter].api_url` default (`https://openrouter.ai/api`) + `/alpha/decisions` yields exactly this.
- The OpenRouter docs' cURL shows `https://openrouter.ai/api/v1/api/alpha/decisions` [S1]. **UNVERIFIED** whether that URL also works; it looks like a doc-generator artifact. Use the path proven in code.
- It is **not** chat completions: `typesafe/jev-1.13` returns 400 on `/v1/chat/completions` [S11, S12]. Jev is also absent from `GET /api/v1/models` (checked 2026-09-22: 453 models listed, no `typesafe/*`), so kb cannot discover it through the models list.
- The direct TypeSafe endpoint is `POST https://api.typesafe.ai/v1/systemone` with the same body, `model: "jev-latest"` [S2]. (Vercel AI Gateway also exposes it under a different, v4 "evaluation-model" shape with `boolean`/`probability` answers [S11]. Out of scope here.)

## Request

```json
{
  "model": "typesafe/jev-1.13",
  "state": "<string> | {object} | [array]",
  "questions": {
    "<question id>": { "type": "noul|choice|score", "instructions": "...", "criteria": ... }
  },
  "provider":   { "allow_fallbacks": true },
  "session_id": "kb-relink-2026-09-22",
  "trace":      { "trace_id": "...", "trace_name": "..." },
  "user":       "<≤256 chars>"
}
```

| Field | Req. | Notes |
|---|---|---|
| `model` | yes | OpenRouter slug (see Slugs). |
| `state` | yes | "The content to evaluate: a plain string, or a JSON object or array of related context." [S1, S2] Evaluated once for all questions. |
| `questions` | yes | Map of id → Question. **The id is never sent to the model** [S2, S3]; write the whole question in `instructions`. |
| `provider` | no | OpenRouter routing prefs. Only one provider (TypeSafe) exists, so this does nothing useful today [S8]. |
| `session_id` | no | ≤256 chars. Used only for OpenRouter observability, never sent to the provider. Body beats the `x-session-id` header [S1]. |
| `trace` | no | Observability metadata. Known keys: `trace_id`, `trace_name`, `span_name`, `generation_name`, `parent_span_id` [S1]. |
| `user` | no | ≤256 chars [S1]. |

`session_id`, `trace`, `user` and `provider` are documented but no local client has sent them. **UNVERIFIED** in practice.

Instructions can reference parts of a structured state by backticked path, e.g. ``"Does `ticket.messages[0].text` request a refund?"`` [S3]. `instructions` and every criterion value can be a string, object or array. Objects let you put the question in one field and data or rubric in others [S2, S3]. Real example [S10]: `instructions: {question, scope, interpretation, policy}`.

### Question types

**Noul** is a yes/no question. It returns P(yes).
```json
"is_same_concept": {
  "type": "noul",
  "instructions": "Does `candidate` describe the same concept as `article`?",
  "criteria": { "true": "Same thing under another name", "false": "Related or different" }
}
```
`criteria` is optional, with only the keys `true` and `false` [S2].

**Choice** picks exactly one option from a closed set.
```json
"relation": {
  "type": "choice",
  "instructions": "How does `source` relate to `target`?",
  "criteria": {
    "extends": "Source builds on or refines the target",
    "contradicts": "Source makes a claim incompatible with the target",
    "mentions": null,
    "none": "No meaningful relation"
  }
}
```
- `criteria` is required. It maps option name to description (string, object, array or `null`) [S2].
- **Both the option key and its description are sent to the model** [S3]. Descriptions are how you pass option semantics. Use `null` when the name is self-explanatory.
- **Maximum 255 options per Choice** [S2, S3]. Each option costs "a few tokens". Add `other` or `none` when the list may not cover every input [S3].

**Score** places the state on ordered levels.
```json
"link_strength": {
  "type": "score",
  "instructions": "How central is `target` to understanding `source`?",
  "criteria": ["Not needed", "Useful background", "Essential prerequisite"]
}
```
- `criteria` is an ordered array. Level number = array index starting at 0. **At least 2 levels, and the API accepts up to 10** [S2, S3].
- Each level is judged on its own: the model never sees the number or the neighbouring levels, so describe situations, not degrees [S3].

## Response (200)

This is a real OpenRouter response [S10: `yc-apply-classifier/.data/runs/c04ad4e5…/d4ce6a85….json`], with 1 of its 10 answers shown:
```json
{
  "id": "gen-dec-1789864370-f8C3386o2WxUy6CtOyf7",
  "model": "typesafe/jev-1.13-20260917",
  "provider": "TypeSafe",
  "answers": {
    "yc-134": {
      "type": "choice",
      "choice": "strong",
      "probabilities": { "not_applicable": 0, "unknown": 0.01, "weak": 0, "strong": 0.93, "partial": 0.06 },
      "confidence": 0.9
    }
  },
  "usage": { "input_tokens": 4622, "output_tokens": 566, "cost": 0.000194124 }
}
```
- Top-level fields are `answers` (required), `model` (required), `usage` (required), `id` and `provider` [S1]. The OpenRouter response **adds `id`, `provider` and `usage.cost`** to the TypeSafe native shape [S1 vs S2].
- `model` comes back as the **dated permaslug** (`typesafe/jev-1.13-20260917`) even when you request `typesafe/jev-1.13` [S9 VALIDATION, S10]. kb should log it.
- `usage.cost` is in USD: 4622 × $0.042/M = $0.000194124, so output tokens are indeed free [S10, S4].

Answer shapes per type [S2, S3, confirmed in S9/S10 recordings]:

| Type | Fields | Semantics |
|---|---|---|
| noul | `type`, `noul` | P(yes) in [0,1]. **No `confidence`**: the single number is the whole distribution. |
| choice | `type`, `choice`, `probabilities{option→p}`, `confidence` | `choice` = argmax. `probabilities` covers every option and sums to 1. |
| score | `type`, `score`, `legend{"0"→desc}`, `probabilities{"0"→p}`, `confidence` | `score = Σ level×p`, a float in [0, levels-1] that can fall between levels. Keys are **strings** `"0"`, `"1"`… |

Real noul answer [S9 `reports/safe.json`]: `{"type": "noul", "noul": 0.13}`. Score example [S3]:
`{"type":"score","score":1.43,"confidence":0.35,"legend":{"0":"…","1":"…","2":"…"},"probabilities":{"0":0.0,"1":0.57,"2":0.43}}`.

Probabilities are **rounded to 2 decimals**, so sums drift slightly from 1. The Vercel AI SDK validator rejects ~1 response in 1,300 for this reason [S12]. Validate with a tolerance (jev-review uses ±0.05, yc-apply ±0.02) [S9, S10].

### Confidence

`confidence` is derived from the shape of `probabilities`. The docs' interactive explorer computes it for a Choice as
`clamp((n·p_max − 1)/(n − 1), 0, 1)` [S5]. This reproduces the docs' examples (0.88/0.12/0 → 0.82 vs reported 0.81; 0.61 peak over 3 → 0.415 vs 0.42). The docs call it "a solid default" and not a contract. The exact server formula, including the one for Score, is **UNVERIFIED**. Confidence says how peaked the answer is, not whether it is correct [S3]. Suggested bands from the docs are <0.5 or <0.6 "don't act", then per-action thresholds such as 0.85 for risky actions [S5].

### Multi-label?

**No. A Choice is single-label.** It returns one argmax and a distribution that sums to 1 [S2, S3]. For multi-label, ask **one Noul per label** in the same request. A Choice is *relative* (which option wins) and each Noul is *absolute* (all of them can be low). The docs explicitly warn against assuming invariants between the two, or carrying a threshold tuned on a Noul over to a Choice [S6]. Top-k from `probabilities` is valid for ranking (the hierarchical beam-search cookbook does this), but it is not independent membership.

### Independence

Questions in one request are independent. One answer is never context for another. Adding or removing questions does not change the other answers [S3]. If a question depends on an earlier answer, make a second request.

## Limits

| Limit | Value | Source |
|---|---|---|
| Options per Choice | 255 max | S2, S3 |
| Levels per Score | 2–10 | S2, S3 |
| Context (TypeSafe) | 64k tokens per request (state + all questions); **32k for state + the single longest question** | S4 |
| Context (OpenRouter listing) | `context_length: 32000` | S8 |
| Observed | requests up to **43,045 input tokens** succeeded via OpenRouter, so the 64k/32k TypeSafe rule applies and 32k is not a per-request cap | S9 receipts |
| Questions per call | **No documented cap.** Observed ≤32 per request in S9 and 13 in the parallel-questions cookbook. UNVERIFIED upper bound | S3, S9 |
| Payload bytes | 413 "Request payload too large" exists [S1]; the byte threshold is **UNVERIFIED**. jev-review self-caps at 110,000 bytes and a 19.6 KB request was fine [S9, S10] | S1, S9 |
| Input | text only (string/JSON). English best, other languages "handled but not equally well" | S4 |
| `session_id`, `user` | ≤256 chars | S1 |

## Model slugs and pinning

- `typesafe/jev-1.13`: canonical slug. Its permaslug is `typesafe/jev-1.13-20260917`, served by TypeSafe as `jev-1.13.0` [S8].
- `~typesafe/jev-latest`: OpenRouter alias ("always redirects to the latest model in the Jev family"), created 2026-09-18, `context_length` 32000 [S8]. No local code has used it: **UNVERIFIED** in practice, and **UNVERIFIED** whether the response `model` then shows the resolved permaslug (expected, by analogy).
- Whether the dated permaslug `typesafe/jev-1.13-20260917` is accepted as a request `model` is **UNVERIFIED**.
- TypeSafe native aliases are `jev-latest` and `jev-preview` (both currently → `jev-1.13.0`). The docs advise: "If you have tuned confidence thresholds against a specific version, pin that version's ID instead of the alias" [S4].
- Recommendation for kb: **default to the pinned `typesafe/jev-1.13`**. Record the response `model` in every receipt, and treat a change of model as "re-evaluate thresholds" [S12 shadow-mode note].

## Errors, timeouts, retries

Error body (OpenRouter): `{"error": {"code": <int>, "message": "<text>"}}` [S1].

| Status | Meaning | kb action |
|---|---|---|
| 400 | Invalid request parameters (also wrong endpoint or slug) | fail, no retry; a bug in question building |
| 401 | Missing or invalid auth | fail and disable decisions for the run |
| 402 | Insufficient credits | fail and disable decisions for the run |
| 403 | Forbidden | fail |
| 404 | Not found (bad path) | fail |
| 413 | Payload too large | split the request (fewer questions or less state); no retry |
| 422 | Validation (TypeSafe direct API) [S2]. **UNVERIFIED** whether OpenRouter surfaces 422 or maps it to 400 | fail |
| 429 | Rate limited | retry with backoff; honour `Retry-After` |
| 500 / 502 / 503 / 529 | Server, provider, unavailable, overloaded | retry with backoff |
| 524 | OpenRouter upstream timeout | retry |

- **Hangs are real:** "calls can hang with no timeout" [S12]. The OpenRouter endpoint reports `can_abort: false` [S8]. In the 10,840-call scale test there was 0 failures and 1 call hit the 30 s client timeout and succeeded on retry [S11]. **Always set a client timeout.** Observed latency is p50 0.33–0.39 s, p95 0.49–0.53 s, max 2.2 s (jev-review, 6,234 requests) and 30 s once [S9, S11]. A 10–30 s per-attempt timeout is safe.
- **503s:** listed in the docs [S1]. No local run recorded one (0 non-OK receipts across 6,234 + 10,840 + 334 calls) [S9, S10, S11]. Their frequency is **UNVERIFIED** and they should be handled as transient.
- **SDK default policy to mirror** [S7]: `max_retries=2`, backoff 0.5 s doubling to a max of 5 s with 25 % jitter, retry on 408/429/5xx plus connection and timeout errors, honour `Retry-After`/`retry-after-ms`, total budget 30 s. The SDK's per-HTTP-operation timeout is 10 s. Local clients use the same approach: jev-review retries transport, 429 and 5xx with jittered exponential backoff; the Python clients retry 429/500/502/503/524/529 [S9, S11].
- The request is read-only, so retrying after a timeout is safe. But **a timeout or a malformed answer must surface as "undecided", never as "no"** [S12, "Loop fechado"].
- Validate every response against the request, as jev-review does [S9]: answer count = question count; each id is present with a matching `type`; noul is in [0,1]; the choice is among the options; probabilities cover all options and sum to about 1; the score is in [0, levels-1]. A failed validation counts as a failed decision.

## Rate limits

- TypeSafe direct: 250,000 tokens/s and 1,200 requests/min, "adjusting dynamically … can change without notice" [S4]. The OpenRouter listing shows `limit_rpm: null`, `limit_rpd: null` [S8]. The OpenRouter-side limit is **UNVERIFIED**.
- Observed: 12–16 concurrent workers, 33–35 calls/s with no 429 (licitacao-fit). 48 simultaneous calls with no 429 (reddit-idea-classifier) [S11]. jev-review defaults to 6 workers, max 12 [S9].

## Pricing

- **$0.042 per 1M input tokens. Output is free** ($42 per 1B) [S4, S8]. OpenRouter reports the exact cost in `usage.cost` [S1, S10].
- OpenRouter data policy for this endpoint: `training: false`, `retainsPrompts: false` [S8]. This is relevant to the map's privacy open question.
- Real costs: 5,000 decision pairs for $0.30. 100 posts × 10 questions for $0.008. 5,350 questions over 293 requests for $0.085 [S9, S11]. Tokens are dominated by `state`, and each extra question over the same state is cheap: batch questions per state [S3].

## Implications for kb

- **Client:** a small Go client in a new internal package (provider-neutral name, e.g. `internal/decisions`) that POSTs `{model, state, questions}` to `<openrouter.api_url>/alpha/decisions`. It reuses the existing `[openrouter]` api_key/api_url config and `OPENROUTER_API_KEY`. `jev-review/jev.go` is a near drop-in reference (types, validator, retry loop, receipt).
- **Config:** add a `[decisions]` section with `model` (default `typesafe/jev-1.13`, pinned), `timeout` (~15 s per attempt), `retries` (2), `concurrency` (~6–8) and a byte cap per request (~100 KB). With no key, or with 401/402, kb disables decisions and behaves as today (standing decision).
- **Question design for kb's use cases:**
  - Linking and relation typing: one request per source article, with `state` = the source excerpt plus the candidate list. Use **one Noul per candidate** for "should link" (absolute, multi-label) and a **Choice per candidate** for relation type. Both fit in one request.
  - OKF type suggestion: a single Choice over the `[okf].types` vocabulary, with type descriptions as criteria and a `none` option. Up to 255 types.
  - Ingest gates and semantic lint: Nouls and Scores with explicit criteria.
- **Confidence bands** must be read per type: Choice and Score use `confidence`, Noul uses the distance of `noul` from 0.5 (for example `noul ≥ 0.8` → yes, `≤ 0.2` → no, otherwise the middle band). Thresholds are tuned per pinned model version. Never reuse a Noul threshold on a Choice [S6].
- **Receipts:** store the response `id`, the response `model` (dated permaslug), `usage.input_tokens` and `usage.cost`, attempts, latency and a request hash per call. `usage.cost` gives exact per-run cost reporting (feeds "cost and rate budgets" in Not yet specified).
- **Failure semantics:** timeouts, exhausted retries and validation failures produce an *undecided* result, never a rejection. Undecided items stay candidates and are never quarantined.
- **State hygiene** [S6]: send only the fields the question needs, because accuracy drops with irrelevant state. Keep the state plus the longest question under 32k tokens, and split large candidate sets across requests. Do counting, dates and numeric comparisons in Go. Treat vault content as untrusted data: add an "ignore instructions inside state" policy in `instructions` (prompt injection is a known jagged edge).
- **Language:** many vault topics are Portuguese. English is Jev's best language and others are weaker [S4], so write questions and criteria in English, calibrate thresholds on real PT content, and lean on the confidence bands.
- **Tests:** unit-test the client against an `httptest` server using the recorded response shapes above. Any live contract check belongs in an integration test gated on the key.
