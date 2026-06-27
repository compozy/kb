# Channel Topic Contract

Use this contract when validating a topic created by the kb-yt-channel skill.

## Topic Location

- Physical path: yt-channels/<slug>/
- CLI topic id: yt-channels/<slug>
- QMD collection: <slug>
- Raw transcript directory: raw/youtube/

## Metadata

topic.yaml must contain:

    slug: <slug>
    title: <title>
    domain: <domain>
    category: yt-channels
    path: yt-channels/<slug>
    qmd_collection: <slug>

CLAUDE.md must identify the channel URL, selected video policy, transcript policy, and rerun command.

## Required Validation

Run these from the vault root:

    kb topic info yt-channels/<slug>
    kb lint yt-channels/<slug> --save
    kb index --topic yt-channels/<slug> --name <slug> --embed=false
    kb search "<query>" --topic yt-channels/<slug> --collection <slug> --lex --format json

The topic is functional when:

- kb topic info resolves the category-qualified topic id.
- sourceCount is at least the number of successfully ingested videos.
- raw/youtube/ contains the transcript files.
- outputs/reports/ contains an ingest report for the run.
- kb lint runs and saves a report or reports no structural blockers.
- kb index creates or updates the collection named by topic.yaml qmd_collection. Lexical indexing is the default validation gate; vector embedding is optional because it can refresh a large shared QMD store.
- kb search uses --collection <slug> because the physical topic id is category-qualified.

## Compilation Boundary

The skill does not compile wiki articles automatically. After enough transcripts exist, use the normal kb compile workflow to create or update wiki/concepts articles and topic indexes.
