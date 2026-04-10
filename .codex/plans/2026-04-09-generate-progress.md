# Generate Progress UX

## Summary

- Separate pipeline execution events from presentation so `generate` can serve both humans and machines cleanly.
- Keep the final command result on `stdout` as JSON.
- Render human-friendly progress on `stderr` only when interactive output is enabled.
- Preserve structured stage events on `stderr` when `--log-format=json` is requested.

## Key Changes

- Introduce structured generation progress events in `internal/generate`.
- Add CLI flags `--progress auto|always|never` and `--log-format text|json`.
- Use a lightweight progress renderer for TTY mode, with incremental progress for `parse` and `write`.
- Fall back to stable line-based status messages when `stderr` is not a TTY or progress is disabled.

## Test Plan

- Unit test event emission order and progress totals for `parse` and `write`.
- CLI tests for default flag values and output mode selection.
- Integration tests to verify JSON stays on `stdout` and progress/logging stays on `stderr`.
- Run focused package tests and `make verify`.

## Assumptions

- Default behavior is `--progress=auto` and `--log-format=text`.
- JSON logs on `stderr` disable animated progress rendering to keep the stream machine-readable.
- `github.com/schollz/progressbar/v3` is the initial dependency for terminal progress rendering.
