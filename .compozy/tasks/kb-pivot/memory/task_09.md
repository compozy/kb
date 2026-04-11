# Task Memory: task_09.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Add EPUB and image converters under `internal/convert/` with spine-ordered EPUB chapter rendering and build-tagged OCR behavior aligned to `task_09.md` and `_techspec.md`.

## Important Decisions

- Treat `_techspec.md` and the current `internal/models.Converter` signature as authoritative over the stale `ADR-003` method signature.
- Keep image converter registration unconditional in `registry.go` by compiling the same `ImageConverter` type from `ocr.go` or `ocr_noop.go` depending on build tags.
- Implement image metadata extraction without a second new dependency so task scope stays tight around the required optional `gosseract` addition.
- Generate OCR test images in Go test helpers instead of adding extra binary fixtures; only the required `sample.epub` fixture is stored in `internal/convert/testdata/`.

## Learnings

- The current environment is missing both the `tesseract` executable and the `tesseract`/`lept` `pkg-config` entries, so native OCR-tag verification is likely limited to non-linked paths in this run.
- `go test -tags ocr ./internal/convert` currently fails in this environment while compiling `github.com/otiai10/gosseract/v2` because `leptonica/allheaders.h` is absent.
- Default-path validation is clean: `go test ./internal/convert`, `go test -cover ./internal/convert` (81.2%), and `make verify` all pass after the converter changes.

## Files / Surfaces

- `internal/convert/epub.go`
- `internal/convert/image_common.go`
- `internal/convert/ocr.go`
- `internal/convert/ocr_noop.go`
- `internal/convert/epub_test.go`
- `internal/convert/image_test.go`
- `internal/convert/image_fixture_test.go`
- `internal/convert/image_fixture_ocr_test.go`
- `internal/convert/ocr_noop_test.go`
- `internal/convert/ocr_test.go`
- `internal/convert/testdata/sample.epub`
- `internal/convert/registry.go`
- `internal/convert/registry_test.go`
- `go.mod`
- `go.sum`

## Errors / Corrections

- Initial `cy-workflow-memory` skill path from the session index did not exist under `~/.agents`; the installed skill is available at `/home/pedronauck/Projects/compozy/skills/cy-workflow-memory/SKILL.md`.
- `make verify` initially failed because OCR-only test helpers lived in an always-built test file; moving them into an `//go:build ocr` helper file restored lint cleanliness on the default build.

## Ready for Next Run

- If this task needs to be fully closed in tracking, install the native Tesseract/Leptonica development packages first, then rerun `go test -tags ocr ./internal/convert` before changing task status or committing.
