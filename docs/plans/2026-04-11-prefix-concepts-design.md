# Prefix Generated Wiki Concepts

## Goal

Prevent generated starter articles in `wiki/concepts/` from colliding with manual or unrelated notes by giving every managed concept file a stable `Kodebase - ` filename prefix.

## Approved Design

- Prefix only the generated concept filename and Obsidian target path.
- Keep the visible article title, frontmatter `title`, and rendered link labels unchanged.
- Centralize the behavior in `internal/vault.GetWikiConceptPath(...)` so renderers, manifests, and tests stay consistent.

## Expected Result

- Generated files are written as `wiki/concepts/Kodebase - <name>.md`.
- Internal links resolve to the prefixed target path while still displaying labels like `Codebase Overview`.
- `CLAUDE.md` continues to list human-readable article names without leaking the filename prefix.

## Verification

- `go test ./internal/vault ./internal/generate`
- `go test -tags integration ./internal/vault ./internal/generate`
- `make verify`
