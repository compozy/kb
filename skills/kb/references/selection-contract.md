# Selection Contract

Every relevance judgment in a wiki topic (ingest gates, `kb classify`, the `remove` queue) is made by a literal decision model against the topic's **selection contract**. The model does exactly what the contract says, so the contract decides precision more than the model does: rewrites of one contract moved precision by up to 30 points on the same documents. Use this reference when drafting or editing a contract by hand; `kb topic contract --draft` gives the generation model the same rules.

## Where it lives

- `topic.yaml` `contract:` is the only source of truth. `kb topic new` writes an empty one.
- The `## Selection contract` section of the topic `CLAUDE.md` is rendered from it on `kb topic new` and on every `--accept`. Do not edit that section: it is overwritten.
- Drafts live in `topic.yaml` `contract_draft:` and are never active until accepted.

## Shape

```yaml
contract:
  purpose: "What this topic is for, in one or two sentences."
  core:
    - "Subjects that are the reason the topic exists"
  adjacent:
    - "Neighbouring subjects that are kept, each with the reason it is kept"
  collected_on_purpose:
    - "Kinds of material deliberately kept even if they look off-topic (SDK docs, READMEs, dataset dumps, audit pages, ...)"
  collected_on_purpose_paths:          # optional
    - "raw/brand-audit-*/**"
  out_of_scope:
    - "Clear junk only, each item qualified so it cannot match a kept line"
```

| Field | Meaning |
| --- | --- |
| `purpose` | What the topic is for, in one or two sentences. Required. |
| `core` | The subjects that are the reason the topic exists. At least one line. |
| `adjacent` | Neighbouring subjects that are kept, each line with the reason it is kept. |
| `collected_on_purpose` | Kinds of material kept deliberately even though they look off-topic. |
| `collected_on_purpose_paths` | Topic-relative globs (`**` allowed) for `raw/` folders that a deliberate collection step produced (an audit, a dataset dump, an acquisition log). Documents there get `relevance: collected_on_purpose` from code, with no relevance call, and are never quarantined for relevance. Quality gates still run. |
| `out_of_scope` | Clear junk only. |

Write every line in English, even for a Portuguese corpus: the lines are question material. Code validation rejects an empty `purpose` or `core`, a line that appears in both `core`/`adjacent` and `out_of_scope`, and a glob that does not compile.

## Drafting rules

1. **Describe what the topic collected on purpose**, as shown by its `raw/` folders, source titles and any screening decisions, not a narrower ideal of what it should be about. A contract narrower than the real collection is coherent, passes every consistency check, and quarantines the topic's own material.
2. **`out_of_scope` lists only material that clearly serves none of the topic's purposes**, and each item is qualified so it cannot match a kept line. Write "equity studies with no microstructure, execution or forecasting method", not "equity studies".
3. **When unsure whether a subject belongs, put it in `adjacent`** with the reason it is kept, never in `out_of_scope`.
4. **No dates, market claims or narrative.** Each line says what a document must discuss.

Measured failures these rules prevent:

- Adding "any sport other than basketball or tennis … offering no method to transfer" to a sports-tech contract flagged the rugby and badminton studies the topic kept as comparison; precision fell from 0.88 to 0.64.
- Restricting a wearable-hardware contract to "motion and impact on a body" flagged the ECG, sweat and glucose sensors the topic had kept; precision fell to 0.38 and recovered to 0.75 once the contract described the real collection.
- A one-line scope flagged 49 of 160 documents off-topic, of which only 30–40% were junk. The five-part contract fixed it.
- A brand audit collected on purpose was judged junk page by page until its folder was a `collected_on_purpose_paths` glob.

## Inputs to read before drafting by hand

The same inputs `--draft` sends to the generation model:

- the topic title, domain and the scope text of the topic `CLAUDE.md`;
- the list of `raw/` subfolders with file counts;
- a sample of about 60 source titles, spread across folders;
- the article titles in `wiki/concepts/`;
- the topic's own curation or screening files, when present (for example `outputs/datasets/*screening*.jsonl`, `outputs/collection/curation-decisions.json`, `exclusions.jsonl`), with their decision reasons.

## Workflow

```bash
kb topic contract <topic> --import-claude   # an existing CLAUDE.md "## Selection contract" section → contract_draft (no model call)
kb topic contract <topic> --draft           # or: the generation model drafts contract_draft from the collection
# edit contract_draft in topic.yaml by hand, following the rules above
kb topic contract <topic> --accept          # self-check, impact preview, then type the slug to activate
```

`--import-claude` reads bullets `Purpose`, `Core` (or `Central`), `Adjacent (keep)`, `Collected on purpose`, `Collected on purpose (paths)` and `Out of scope`.

`--accept` runs two checks before activation:

1. **Self-check.** One question per (kept line, `out_of_scope` line) pair: does the out-of-scope line exclude what the kept line describes? Pairs at P ≥ 0.7 (`contract_conflict`) are printed and block acceptance unless `--force`. Fix the out-of-scope line by qualifying it (rule 2), not by forcing.
2. **Impact preview.** The relevance and quality questions over every source (a stratified sample of 500 by `raw/` folder when larger, about US$ 0.15), printing counts per band (kept / review / would-be quarantine), counts per `raw/` folder, and the 20 documents most likely off-topic with title and path. With labels (`kb review import-labels`), it also prints precision and recall of the would-be quarantine.

Read the top-20 list before typing the slug. If it contains material the topic collected on purpose, the contract is too narrow: add the subject to `adjacent` or `collected_on_purpose`, or the folder to `collected_on_purpose_paths`, and accept again. The preview receipts are cached, so the following `kb classify` reuses them when the draft is accepted unchanged. Labels shown in a preview are demoted to calibration dev-only once the contract changes.

## After acceptance

- Relevance gates still run in **shadow** until the topic has ≥ 30 relevance labels given in `kb review` with a stored calibration (`kb review calibrate --write`), or `topic.yaml` sets `decisions.gates: apply`.
- Topics whose membership is defined by construction (a directory mirror, one note per company) should set `decisions.relevance: off` instead of writing a contract.
- A changed contract changes the cache key: the next `kb classify` re-judges relevance for every source.
