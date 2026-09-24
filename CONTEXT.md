# kb

A topic-based knowledge base in the Karpathy KB pattern: raw sources are ingested into a topic and compiled into interlinked wiki articles.

## Knowledge base

**Topic**:
A self-contained knowledge base under the vault root, with its own raw sources, wiki, and scope.
_Avoid_: Project, collection

**Source**:
A raw document ingested into a topic's `raw/` tree, carrying frontmatter with its provenance.
_Avoid_: Raw file, input

**Article**:
A compiled wiki document that synthesizes one concept from many sources.
_Avoid_: Page, note, concept file

**Alias**:
An alternative name of an article (acronym, synonym, other language) that should resolve to it.
_Avoid_: Synonym, keyword

## Graph

**Link**:
A wikilink `[[X]]` written in a document body.
_Avoid_: Reference, backlink (a backlink is a link seen from its target)

**Relation**:
A typed edge between two documents recorded in frontmatter, such as `related`, `extends`, `prerequisite`, or `contradicts`.
_Avoid_: Link, edge, association

**Affects**:
The relation from a source to the articles whose content it should change on the next compile.
_Avoid_: Impacts, touches

## Decisions

**Candidate**:
A pair of documents, or a mention inside a document, proposed by code as a possible link or relation.
_Avoid_: Suggestion, match

**Decision**:
A typed judgment with a calibrated probability, returned by the decision model for one question about one state.
_Avoid_: Classification, prediction, label

**Decision model**:
The provider-neutral model that answers typed questions without generating text; Jev is the default.
_Avoid_: Classifier, Jev (in code and config)

**Confidence band**:
The outcome bucket of a decision: apply, review, or ignore.
_Avoid_: Threshold, tier

**Receipt**:
The stored record of one decision-model call: the raw answers plus what produced them (question, scope, model version, cost).
_Avoid_: Log entry, cache entry

**Decision mode**:
Per-topic setting of whether decisions may change document bodies and gate outcomes: shadow (record only) or apply.
_Avoid_: Dry run, preview

**Selection contract**:
The written scope of a topic (purpose, core, adjacent, collected on purpose, out of scope) against which relevance is judged.
_Avoid_: Scope line, topic description

**Question bank**:
A versioned set of typed questions for one purpose.
_Avoid_: Prompt, rubric file

**Facet**:
A classification dimension stored on a document, such as kind, depth, relevance role, or concepts.
_Avoid_: Tag, label

**Quarantine**:
Where gated sources are kept out of every pipeline without being deleted.
_Avoid_: Trash, rejected

**Review queue**:
The pending review-band decisions of a topic, awaiting a human verdict.
_Avoid_: Inbox, backlog

**Label**:
A human verdict on a decision, kept as evidence for calibration. An **imported label** comes from a topic's own screening or curation file and records that origin.
_Avoid_: Annotation, feedback

**Criterion**:
The definitional sentence of an article that concept questions judge against, distinct from its summary.
_Avoid_: Definition, description

**Provenance**:
Where and how a source was collected (path, host, ingest batch, ingest query), carried in every decision state.
_Avoid_: Origin, metadata

**Impact preview**:
The shadow run a contract draft passes through before it becomes active: what it would keep, review and quarantine.
_Avoid_: Dry run, simulation

**Quarantine ledger**:
The record of every index line and frontmatter entry removed when a source is quarantined, replayed on restore.
_Avoid_: Undo log, backup

**Holdout**:
The share of labels never used to choose thresholds, used only to report their precision.
_Avoid_: Test set, validation set
