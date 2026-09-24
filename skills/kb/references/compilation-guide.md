# Wiki Article Compilation Guide

Writing standards for articles in `<topic>/wiki/concepts/`. These are the primary output of the knowledge base and the interface the LLM answers queries against.

**Division of labour.** The agent writes and updates article prose. `kb` never edits article prose: `kb classify` writes the article's `criterion`, `aliases`, `summary`, `genre` and `concepts` frontmatter, and `kb link` writes typed relations (`related`, `extends`, ...) plus, in `apply` mode, body links at the first qualifying mention. Do not hand-maintain those keys; editing one makes it yours and kb stops updating it (see `frontmatter-schemas.md`).

## Target characteristics

- **Length:** 3000-4000 words. Split into sub-articles when exceeded.
- **Audience:** senior practitioners in the topic's field. Assume foundational literacy; explain domain-specific terms on first use.
- **Standalone:** a reader should be able to learn the topic from one article alone.
- **Dense wikilinks:** target 10-30 `[[wikilinks]]` per article, including cross-topic links where relevant.
- **Cited:** every factual claim traces back to a file in `raw/` listed under `sources:` frontmatter.

## Voice and style

- **Domain knowledge, not personal wiki.** Write as a reference anyone in the field could use. No "what this means for [person]" sections. No builder profiles. No first-person narration.
- **Declarative, technical, neutral.** Avoid hype. Avoid hedging. State what is true, with sources backing it.
- **Concrete examples.** Prefer code blocks, tables, and diagrams over prose when describing structures, comparisons, or flows.

## Required sections

Every wiki article has:

1. **H1 title** — matches the filename.
2. **Lead paragraph** — 2-4 sentences establishing what the topic is, why it matters, and scoping the article.
3. **Core sections** (H2) — the substantive body. Exact structure depends on the topic but should follow a consistent hierarchy.
4. **Sources and Further Reading** — bulleted list of every cited source plus related wikilinks.

Optional sections depending on topic:

- **Comparison tables** — when the article surveys alternatives
- **Code examples** — runnable snippets demonstrating the concepts
- **Architecture diagrams** — ASCII or Mermaid
- **Trade-offs** — explicit pros/cons when the topic has design tensions
- **Future direction** — where the field is heading, if well-established

## Wikilink density

Wikilinks and frontmatter relations are the knowledge graph. Every mention of a related concept should be a wikilink on first occurrence, and ideally a second time in a later section. Examples of good density:

- Mention of another concept article → `[[Concept Name]]`
- Mention of a protocol, tool, or framework that has its own article → `[[Tool Name]]`
- Cross-topic reference → `[[other-topic/wiki/concepts/Article Name|Display Name]]`

Do not wikilink every occurrence of common words. Do not wikilink authors or organizations unless they have their own article.

Write the links you mean while drafting. `kb link` adds the ones you missed afterwards (see *Backlinks and relations* below); it only adds links and never removes one you wrote.

## Sourcing rules

- **Every article cites real sources.** Do not write from general knowledge alone. If the corpus does not contain the claim, either ingest a new source (`kb ingest url/file`) or omit the claim.
- **Frontmatter `sources:`** lists every raw file that informed the article, as wikilinks.
- **Inline attributions** are allowed but not required. A Sources section at the bottom is mandatory.
- **Direct quotes** require quotation marks and a source reference.

## Anti-patterns

- Articles that summarize a single source — instead, synthesize across multiple sources, or cite the single source as reference and link to the raw file.
- Articles with no incoming wikilinks (orphans) — every article should be reachable via the link graph.
- Articles with no outgoing wikilinks — every article should participate in the graph.
- Prose that could apply to any topic — be specific to this topic's vocabulary, patterns, and tensions.
- Restating prerequisites at length — link to the prerequisite article and move on.

## When updating an existing article

1. Load the current article fully.
2. **Find the sources that affect it.** `kb link` marks each source with `affects: ["[[Article]]"]` when the source should change what the article says, and `kb lint <topic-id>` reports a `needs-compile` issue for every article older than a source that affects it (the source's `scraped` is newer than the article's `updated`):

   ```bash
   kb lint <topic-id> --format json    # kind "needs-compile": filePath = the article, target = the affecting source
   ```

   Load those sources fully. When the topic has not been linked yet, fall back to sources added since the article's `updated:` date.
3. Identify what changed in the sources (new techniques, corrections, new terminology).
4. **Propose each change with a structured diff before writing.** Present to the user:

   > **Current:** `<quote the existing text>`
   >
   > **Proposed:** `<replacement text>`
   >
   > **Reason:** `<why this change is warranted>`
   >
   > **Source:** `<raw/ file path or URL backing the new claim>`

   Always include **Source**. An edit without a source citation creates untraceability — future compile passes won't know why the change was made. Ask for confirmation per page. Do not batch-apply changes.

5. **Run a contradiction sweep.** If the new information contradicts something in the wiki, the contradicted claim may appear in more than one article. Before rewriting, grep every article for the contradicted claim:

   ```bash
   grep -rln "<contradicted claim or key term>" <topic>/wiki/concepts/
   ```

   Update all occurrences, not just the most obvious one. Silent contradictions across articles are the worst failure mode of a multi-article wiki.

6. **Check downstream effects.** List the articles that link to the one being updated, in the body or in a frontmatter relation list (`related`, `extends`, `prerequisite`, `example_of`, `contradicts`), and ask for each: *does the update change anything that page asserts?* If yes, flag it explicitly and offer to update it with the same Current/Proposed/Reason/Source flow. Open `contradicts` items are listed by `kb review <topic-id> --queue contradiction` and reported by lint as `contradiction`.

   ```bash
   grep -rln "\[\[<Article Title>" <topic>/wiki/concepts/
   ```

7. Update the article in place, preserving structure where possible.
8. Bump `updated:` in frontmatter (this clears `needs-compile` for the sources you absorbed).
9. Add any new `sources:` entries.
10. Check that existing wikilinks still resolve; add new ones for newly-introduced concepts.
11. Run `kb classify <topic-id> --only-missing` (refreshes `criterion` and `summary` when the body changed) and `kb link <topic-id>` (see below).

## Backlinks and relations (`kb link`)

A compounding wiki needs bidirectional links: every article needs incoming links from the documents that discuss its concept. `kb link` replaces the manual grep-based backlink audit:

```bash
kb classify <topic-id> --only-missing   # a new article gets its criterion and aliases, which linking judges against
kb link <topic-id> --dry-run            # optional: see what would change
kb link <topic-id>
```

- Code proposes candidates (title and alias mentions, BM25, shared `concepts` or `sources`, optional qmd neighbours); the decision model judges each one. When an article is created or gains aliases, kb scans every document for mentions of it and judges those documents against it (the automated backlink audit).
- Confident links become frontmatter relations in every mode. Body links `[[File name|matched text]]` at the first qualifying mention are inserted only when the topic's decision mode is `apply`; in `shadow` they are queued in `kb review <topic-id> --queue link`.
- Mentions inside code blocks, headings, URLs and existing wikilinks are skipped, and a mention in a different sense is rejected by the decision model.
- `kb link` only adds. Remove a wrong link by hand; a relation kb wrote that you edit becomes yours.

Review what landed in the review band with `kb review <topic-id> --queue link` and accept or reject it; each verdict becomes a label used by `kb review calibrate`.

## When to split an article

Split when any of these hold:

- Word count exceeds 4000
- A single H2 section exceeds 1000 words
- The article covers two distinct sub-topics that warrant their own entries
- Multiple other articles would benefit from linking to a sub-section (that sub-section deserves its own article)

After splitting, update the parent article to wikilink to the new sub-article(s) and update the topic's indexes.

## When to write a new article (vs extend an existing one)

Write new when:

- Three or more existing articles wikilink to the concept as a dead link
- A query answer keeps synthesizing the same cross-article content (that synthesis deserves its own article)
- The topic is a distinct concept with its own sources, patterns, and terminology

Extend existing when:

- The new material is a refinement, example, or sub-aspect of an existing concept
- The sub-topic would be under 500 words on its own
