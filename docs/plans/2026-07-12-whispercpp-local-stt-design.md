# Native whisper.cpp transcription

## Goal

Let `kb ingest youtube` use a local `whisper.cpp` server without an OpenAI API
key or an external PowerShell wrapper.

## Design

Add `whispercpp` as an STT provider and a TOML-backed `[whispercpp]` section
with the server executable path, model path, host, and port. When that provider
is selected, the media extractor probes the configured local endpoint. If it is
not listening, it starts `whisper-server` with the configured model and the
OpenAI-compatible transcription route, then waits for the endpoint to become
reachable. An already-running server is reused and no server is stopped by
`kb`.

The `whispercpp` transcriber sends the existing multipart request shape to
`/v1/audio/transcriptions`, but deliberately sends no `Authorization` header.
The OpenAI and OpenRouter providers keep their current behavior.

## Configuration

```toml
[stt]
provider = "whispercpp"
model = "whisper.cpp-small"

[whispercpp]
server_path = "C:/Users/me/.local/whisper.cpp/cuda/Release/whisper-server.exe"
model_path = "C:/Users/me/.local/whisper.cpp/models/ggml-small.bin"
host = "127.0.0.1"
port = 8188
```

## Validation

Tests cover configuration validation, provider selection, an unauthenticated
multipart request, and server startup/reuse behavior through injected process
and readiness dependencies. The repository completion gate is `make verify`.
