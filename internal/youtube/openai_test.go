package youtube

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/config"
)

func TestOpenAITranscriberSendsMultipartTranscriptionRequest(t *testing.T) {
	t.Parallel()

	audio := []byte("audio-bytes")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v1/audio/transcriptions" {
			t.Fatalf("path = %q, want /api/v1/audio/transcriptions", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer openai-key" {
			t.Fatalf("authorization = %q", got)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
			t.Fatalf("content-type = %q", r.Header.Get("Content-Type"))
		}
		if err := r.ParseMultipartForm(1024 * 1024); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("model"); got != "gpt-4o-transcribe" {
			t.Fatalf("model = %q", got)
		}
		if got := r.FormValue("response_format"); got != "json" {
			t.Fatalf("response_format = %q", got)
		}
		if got := r.FormValue("language"); got != "pt" {
			t.Fatalf("language = %q", got)
		}
		if got := r.FormValue("prompt"); got != "Domain context" {
			t.Fatalf("prompt = %q", got)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("form file: %v", err)
		}
		defer func() {
			_ = file.Close()
		}()
		data, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		if string(data) != string(audio) {
			t.Fatalf("audio payload = %q", string(data))
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"text":"transcribed words"}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewOpenAITranscriber(config.STTConfig{
		APIKey:   "openai-key",
		APIURL:   server.URL + "/api",
		Model:    "gpt-4o-transcribe",
		Language: "pt",
		Prompt:   "Domain context",
	})
	transcript, err := client.Transcribe(context.Background(), audio, "mp3")
	if err != nil {
		t.Fatalf("Transcribe returned error: %v", err)
	}
	if transcript != "transcribed words" {
		t.Fatalf("transcript = %q", transcript)
	}
}

func TestOpenAITranscriberReturnsHelpfulErrorWhenAPIKeyMissing(t *testing.T) {
	t.Parallel()

	client := NewOpenAITranscriber(config.STTConfig{})
	_, err := client.Transcribe(context.Background(), []byte("audio"), "mp3")
	if err == nil {
		t.Fatal("expected Transcribe to fail")
	}
	if !strings.Contains(err.Error(), "missing API key") || !strings.Contains(err.Error(), "OPENAI_API_KEY") {
		t.Fatalf("error = %q, want API key guidance", err.Error())
	}
}

func TestOpenAITranscriberReturnsAPIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error":{"message":"bad audio"}}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewOpenAITranscriber(config.STTConfig{
		APIKey: "openai-key",
		APIURL: server.URL,
		Model:  "gpt-4o-transcribe",
	})
	_, err := client.Transcribe(context.Background(), []byte("audio"), "mp3")
	if err == nil {
		t.Fatal("expected Transcribe to fail")
	}
	if !strings.Contains(err.Error(), "bad audio") {
		t.Fatalf("error = %q, want API detail", err.Error())
	}
}
