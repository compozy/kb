package youtube

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/kb/internal/config"
)

func TestOpenRouterClientTranscribeSendsExpectedRequest(t *testing.T) {
	t.Parallel()

	audio := []byte("audio-bytes")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/api/v1/chat/completions" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/api/v1/chat/completions")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer openrouter-key" {
			t.Fatalf("authorization = %q, want %q", got, "Bearer openrouter-key")
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("content-type = %q, want %q", got, "application/json")
		}

		var body chatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if body.Model != "acme/stt" {
			t.Fatalf("model = %q, want %q", body.Model, "acme/stt")
		}
		if body.Stream {
			t.Fatal("expected stream to be false")
		}
		if len(body.Messages) != 1 {
			t.Fatalf("messages length = %d, want 1", len(body.Messages))
		}
		if len(body.Messages[0].Content) != 2 {
			t.Fatalf("content length = %d, want 2", len(body.Messages[0].Content))
		}
		if body.Messages[0].Content[0].Type != "text" {
			t.Fatalf("first part type = %q, want %q", body.Messages[0].Content[0].Type, "text")
		}
		if body.Messages[0].Content[0].Text != defaultTranscriptionPrompt {
			t.Fatalf("prompt = %q, want %q", body.Messages[0].Content[0].Text, defaultTranscriptionPrompt)
		}
		if body.Messages[0].Content[1].Type != "input_audio" {
			t.Fatalf("audio part type = %q, want %q", body.Messages[0].Content[1].Type, "input_audio")
		}
		if body.Messages[0].Content[1].InputAudio == nil {
			t.Fatal("expected input_audio payload")
		}
		if body.Messages[0].Content[1].InputAudio.Format != "m4a" {
			t.Fatalf("audio format = %q, want %q", body.Messages[0].Content[1].InputAudio.Format, "m4a")
		}
		if body.Messages[0].Content[1].InputAudio.Data != base64.StdEncoding.EncodeToString(audio) {
			t.Fatalf("audio data = %q", body.Messages[0].Content[1].InputAudio.Data)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"choices": [
				{
					"message": {
						"role": "assistant",
						"content": "transcribed words"
					}
				}
			]
		}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewOpenRouterClient(config.OpenRouterConfig{
		APIKey:   "openrouter-key",
		APIURL:   server.URL + "/api",
		STTModel: "acme/stt",
	})

	transcript, err := client.Transcribe(context.Background(), audio, "m4a")
	if err != nil {
		t.Fatalf("Transcribe returned error: %v", err)
	}
	if transcript != "transcribed words" {
		t.Fatalf("transcript = %q, want %q", transcript, "transcribed words")
	}
}

func TestOpenRouterClientTranscribeParsesContentPartsResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"choices": [
				{
					"message": {
						"role": "assistant",
						"content": [
							{"type": "text", "text": "hello"},
							{"type": "text", "text": "world"}
						]
					}
				}
			]
		}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewOpenRouterClient(config.OpenRouterConfig{
		APIKey: "openrouter-key",
		APIURL: server.URL + "/api",
	})

	transcript, err := client.Transcribe(context.Background(), []byte("audio"), "mp3")
	if err != nil {
		t.Fatalf("Transcribe returned error: %v", err)
	}
	if transcript != "hello world" {
		t.Fatalf("transcript = %q, want %q", transcript, "hello world")
	}
}

func TestOpenRouterClientTranscribeReturnsHelpfulErrorWhenAPIKeyMissing(t *testing.T) {
	t.Parallel()

	client := NewOpenRouterClient(config.OpenRouterConfig{})

	_, err := client.Transcribe(context.Background(), []byte("audio"), "mp3")
	if err == nil {
		t.Fatal("expected Transcribe to fail")
	}
	if !strings.Contains(err.Error(), "missing API key") {
		t.Fatalf("error = %q, want missing API key message", err.Error())
	}
	if !strings.Contains(err.Error(), "OPENROUTER_API_KEY") {
		t.Fatalf("error = %q, want env guidance", err.Error())
	}
}

func TestOpenRouterClientConfigured(t *testing.T) {
	t.Parallel()

	if NewOpenRouterClient(config.OpenRouterConfig{}).Configured() {
		t.Fatal("expected client without key to be unconfigured")
	}
	if !NewOpenRouterClient(config.OpenRouterConfig{APIKey: "openrouter-key"}).Configured() {
		t.Fatal("expected client with key to be configured")
	}
}

func TestOpenRouterClientTranscribeReturnsAPIErrorOnNonSuccessStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error":{"message":"bad audio payload"}}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewOpenRouterClient(config.OpenRouterConfig{
		APIKey: "openrouter-key",
		APIURL: server.URL + "/api",
	})

	_, err := client.Transcribe(context.Background(), []byte("audio"), "mp3")
	if err == nil {
		t.Fatal("expected Transcribe to fail")
	}
	if !strings.Contains(err.Error(), "bad audio payload") {
		t.Fatalf("error = %q, want API error message", err.Error())
	}
}

func TestParseOpenRouterErrorAndTranscriptionContentHelpers(t *testing.T) {
	t.Parallel()

	if got := parseOpenRouterError([]byte(`{"choices":[{"error":{"message":"choice failure"}}]}`)); got != "choice failure" {
		t.Fatalf("parseOpenRouterError = %q, want %q", got, "choice failure")
	}
	if got := parseOpenRouterError([]byte(`plain text error`)); got != "plain text error" {
		t.Fatalf("parseOpenRouterError plain = %q, want %q", got, "plain text error")
	}

	text, err := parseTranscriptionContent(json.RawMessage(`"hello world"`))
	if err != nil {
		t.Fatalf("parseTranscriptionContent string returned error: %v", err)
	}
	if text != "hello world" {
		t.Fatalf("text = %q, want %q", text, "hello world")
	}

	text, err = parseTranscriptionContent(json.RawMessage(`{"type":"text","text":"single part"}`))
	if err != nil {
		t.Fatalf("parseTranscriptionContent object returned error: %v", err)
	}
	if text != "single part" {
		t.Fatalf("text = %q, want %q", text, "single part")
	}

	if _, err := parseTranscriptionContent(json.RawMessage(`""`)); err == nil {
		t.Fatal("expected empty transcription content to fail")
	}
}
