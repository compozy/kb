package youtube

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
)

func TestNewTranscriberUsesOpenRouterCredentialsWhenOpenAIEnvIsAlsoSet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kb.toml")
	if err := os.WriteFile(path, []byte(strings.Join([]string{
		"[app]",
		"name = \"kb\"",
		"env = \"development\"",
		"",
		"[openrouter]",
		"stt_model = \"openrouter/stt\"",
	}, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv(config.EnvSTTProvider, "openrouter")
	t.Setenv(config.EnvOpenRouterAPIKey, "openrouter-key")
	t.Setenv(config.EnvOpenRouterAPIURL, "https://openrouter.internal/api")
	t.Setenv(config.EnvOpenAIAPIKey, "openai-key")
	t.Setenv(config.EnvOpenAIAPIURL, "https://openai.internal")

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.STT.APIKey == "openai-key" {
		t.Fatal("openai key leaked into stt config while provider is openrouter")
	}
	if cfg.STT.APIURL == "https://openai.internal" {
		t.Fatal("openai api url leaked into stt config while provider is openrouter")
	}

	transcriber := NewTranscriber(cfg.STT, cfg.OpenRouter)
	client, ok := transcriber.(*OpenRouterClient)
	if !ok {
		t.Fatalf("transcriber = %T, want *OpenRouterClient", transcriber)
	}
	if client.apiKey != "openrouter-key" {
		t.Fatalf("api key = %q, want openrouter-key", client.apiKey)
	}
	if client.apiURL != "https://openrouter.internal/api" {
		t.Fatalf("api url = %q, want openrouter api url", client.apiURL)
	}
	if client.model != "openrouter/stt" {
		t.Fatalf("model = %q, want openrouter/stt", client.model)
	}
}

func TestNewTranscriberCarriesOpenRouterPromptAndLanguage(t *testing.T) {
	t.Parallel()

	transcriber := NewTranscriber(config.STTConfig{
		Provider: "openrouter",
		Language: "pt",
		Prompt:   "Domain context.",
	}, config.OpenRouterConfig{
		APIKey:   "openrouter-key",
		APIURL:   "https://openrouter.internal/api",
		STTModel: "openrouter/stt",
	})

	client, ok := transcriber.(*OpenRouterClient)
	if !ok {
		t.Fatalf("transcriber = %T, want *OpenRouterClient", transcriber)
	}
	prompt := client.transcriptionPrompt()
	for _, want := range []string{
		defaultTranscriptionPrompt,
		"Domain context.",
		"The spoken language is pt.",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt %q does not contain %q", prompt, want)
		}
	}
}

func TestTranscribeAudioPathSegmentsLongAudioAndPreservesOffsets(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	audioPath := filepath.Join(dir, "input.mp3")
	if err := os.WriteFile(audioPath, []byte("original audio"), 0o644); err != nil {
		t.Fatalf("write audio: %v", err)
	}
	ffmpegPath := filepath.Join(dir, "ffmpeg")
	if err := os.WriteFile(ffmpegPath, []byte(strings.Join([]string{
		"#!/bin/sh",
		"set -eu",
		"last=\"\"",
		"for arg in \"$@\"; do last=\"$arg\"; done",
		"out_dir=$(dirname \"$last\")",
		"mkdir -p \"$out_dir\"",
		"printf '%s' 'first audio' > \"$out_dir/chunk-000.mp3\"",
		"printf '%s' 'second audio' > \"$out_dir/chunk-001.mp3\"",
	}, "\n")+"\n"), 0o755); err != nil {
		t.Fatalf("write fake ffmpeg: %v", err)
	}
	stt := &stubSTTClient{
		transcripts: []string{"First transcript", "Second transcript"},
	}
	extractor := &Extractor{
		stt: stt,
		sttConfig: config.STTConfig{
			Provider:      "openai",
			APIURL:        "https://api.openai.com",
			Model:         "gpt-4o-transcribe",
			Language:      "auto",
			AudioFormat:   "mp3",
			ChunkDuration: "10s",
			MaxChunkBytes: 100,
			Concurrency:   1,
			FFmpegPath:    ffmpegPath,
		},
	}

	markdown, err := extractor.transcribeAudioPath(context.Background(), audioPath, "mp3", 25*time.Second)
	if err != nil {
		t.Fatalf("transcribeAudioPath returned error: %v", err)
	}
	want := strings.Join([]string{
		"## 00:00",
		"First transcript",
		"",
		"## 00:10",
		"Second transcript",
	}, "\n")
	if markdown != want {
		t.Fatalf("markdown = %q, want %q", markdown, want)
	}
	if stt.format != "mp3" {
		t.Fatalf("format = %q, want mp3", stt.format)
	}
}
