# Cleaning an Existing Topic

How to bring a topic that predates the decision model in line: a written scope, classification facets on every document, broken captures recaptured or quarantined, off-topic sources reviewed, and typed links. Every step is incremental and resumable through the receipts cache, and nothing leaves `raw/` without an explicit `kb review accept`.

**Before you start**

- `OPENROUTER_API_KEY` is set (`--import-claude` and the label import are the only steps that call no model).
- Budget: about US$ 0.30 per 1,000 documents to classify, US$ 0.30 to link, US$ 0.15 in generation, plus about US$ 0.15 for the contract preview. Each run stops at `[decisions].budget_usd` (US$ 1 by default) or `--budget`; re-run to continue from the cache.
- Commit or back up the vault. kb writes only its own frontmatter keys, `topic.yaml` and `<topic>/.decisions/`, but you will review its changes.

## 1. Write and accept the selection contract

```bash
kb topic contract <topic> --import-claude   # when CLAUDE.md already has a "## Selection contract" section
kb topic contract <topic> --draft           # otherwise
```

Edit `contract_draft` in `topic.yaml` following `selection-contract.md`, then:

```bash
kb topic contract <topic> --accept
```

Read the impact preview: counts per band and per `raw/` folder, and the 20 documents most likely off-topic. If material the topic collected on purpose shows up there, widen `adjacent`, `collected_on_purpose` or `collected_on_purpose_paths` and accept again. Type the topic slug only when the list looks like junk.

Skip this step for topics defined by construction (one note per entity of a directory): set `decisions.relevance: off` in `topic.yaml` instead.

## 2. Import the topic's own labels (when it has them)

When the topic has screening or curation files, import them as relevance labels so the preview, the queues and calibration can report precision:

```bash
kb review import-labels <topic> --from outputs/datasets/science-screening.jsonl \
  --id-field pmcid --match pmcid --decision-field decision --keep-values include,true
```

`--match` is `url`, `path`, `doi` or `pmcid`; files may be JSONL, JSON or CSV. Unmatched rows are counted and listed. Imported labels keep their origin, are reported separately, and never count toward the 30 labels that switch relevance gates to `apply`. Treat them as imperfect: a title-only screen keeps things a full read would drop.

## 3. Classify

```bash
kb classify <topic>
```

This writes `genre`, `depth`, `relevance`, `quality`, `concepts`, `summary`, `entities` and `questions` on sources, and `criterion`, `aliases`, `summary`, `genre`, `concepts` on articles. It never quarantines: broken captures go to the `recapture` queue and off-topic sources to the `remove` queue.

**Topic without articles.** Classification writes no `concepts` (there is nothing to choose from) and lint reports `vocabulary-missing`. Bootstrap a vocabulary once the summaries exist:

```bash
kb topic vocabulary <topic> --draft    # 12–40 concepts with criteria, and how many summaries mention each
# edit vocabulary_draft in topic.yaml: drop, merge or rename concepts
kb topic vocabulary <topic> --accept   # one stub article per concept under wiki/concepts/
kb classify <topic> --only-missing     # fills concepts
```

Stub articles (`stage: stub`) carry a `criterion`; compile them later with `compilation-guide.md`.

## 4. Work the two cleaning queues

Broken captures and off-topic sources need different actions, so they are separate queues. Handle `recapture` first: a page that comes back with full text is no longer a removal candidate.

```bash
kb review <topic> --queue recapture
kb review accept --topic <topic> <id>...        # refetch fresh and re-judge each item
kb review reject --topic <topic> <id>...        # keep the source as it is
```

Accepting a `recapture` item refetches the URL with Firecrawl's cache bypassed; a page that is now good gets its body replaced in place (`recaptured: <date>`), a page still broken is quarantined with its quality reason. Sources without `source_url` are listed as `no_source` and can only be quarantined.

```bash
kb review <topic> --queue remove
kb review accept --topic <topic> --queue remove --above 0.9   # bulk, only after reading the list
```

Read titles and paths before bulk-accepting. Several wrong items in the list mean the contract is too narrow: fix it (step 1) instead of rejecting items one by one. Quarantine removes the index lines and `sources:`/relation entries that pointed at each file and records them in the quarantine ledger; restoring a source later with `kb review accept` puts back the file and every removed line byte for byte. Body links in other documents stay and lint reports them as `link-to-quarantined`.

Prose counts written by agents ("66 scraped articles" in `CLAUDE.md` or the Dashboard) are not updated by kb; the run summary prints the new counts per folder so you can update them.

## 5. Link

```bash
kb link <topic> --dry-run   # optional
kb link <topic>
```

Confident links become frontmatter relations (`related`, `extends`, `prerequisite`, `example_of`) and `affects` on sources. Body links are inserted only when the topic's decision mode is `apply`; in `shadow` they wait in `kb review <topic> --queue link`. `contradicts` always goes to review.

## Afterwards

- `kb lint <topic>`: `unclassified`, `criterion-missing`, `off-topic-kept`, `recapture-pending`, `remove-pending` and `pending-review` should be empty or explained.
- `kb review import-links <topic>` turns existing human links into positive link labels. With ≥ 30 labels per purpose, `kb review calibrate <topic>` reports dev and holdout precision; `--write` stores topic-specific thresholds and lets relevance gates leave shadow.
- New ingests now run the same gates, classification and linking automatically.

Measured on a 34-topic research vault: US$ 2.03 for 8,588 sources (relevance and quality), 320 sources marked at P ≥ 0.8 and split into 101 broken captures and 219 off-topic, 13 of the 101 recovered with full text by a fresh refetch (about 160 Firecrawl credits).
