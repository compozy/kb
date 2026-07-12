package mediadl

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
)

type stubWhisperCPPStarter struct {
	calls int
	err   error
}

func TestNewTranscriberCreatesWhisperCPPProvider(t *testing.T) {
	t.Parallel()

	transcriber := NewTranscriber(config.STTConfig{
		Provider: "whispercpp",
		Model:    "whisper.cpp-small",
	}, config.OpenRouterConfig{}, config.WhisperCPPConfig{
		Host: "127.0.0.1",
		Port: 8188,
	})
	client, ok := transcriber.(*WhisperCPPTranscriber)
	if !ok {
		t.Fatalf("transcriber = %T, want *WhisperCPPTranscriber", transcriber)
	}
	if client.Provider() != "whispercpp" {
		t.Fatalf("provider = %q", client.Provider())
	}
}

func (starter *stubWhisperCPPStarter) Ensure(context.Context) error {
	starter.calls++
	return starter.err
}

func TestWhisperCPPTranscriberSendsUnauthenticatedMultipartRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != whisperCPPTranscriptionsPath {
			t.Fatalf("path = %q, want %q", r.URL.Path, whisperCPPTranscriptionsPath)
		}
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("authorization = %q, want no header", got)
		}
		if err := r.ParseMultipartForm(1024 * 1024); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("model"); got != "whisper.cpp-small" {
			t.Fatalf("model = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"text":"local transcript"}`))
	}))
	defer server.Close()

	starter := &stubWhisperCPPStarter{}
	client := &WhisperCPPTranscriber{
		endpoint:   server.URL + whisperCPPTranscriptionsPath,
		model:      "whisper.cpp-small",
		language:   "auto",
		starter:    starter,
		httpClient: server.Client(),
	}
	transcript, err := client.Transcribe(context.Background(), []byte("audio"), "mp3")
	if err != nil {
		t.Fatalf("Transcribe returned error: %v", err)
	}
	if transcript != "local transcript" {
		t.Fatalf("transcript = %q", transcript)
	}
	if starter.calls != 1 {
		t.Fatalf("server ensure calls = %d, want 1", starter.calls)
	}
}

func TestWhisperCPPTranscriberReportsServerStartupFailure(t *testing.T) {
	t.Parallel()

	starter := &stubWhisperCPPStarter{err: errors.New("model missing")}
	client := &WhisperCPPTranscriber{starter: starter}
	_, err := client.Transcribe(context.Background(), []byte("audio"), "mp3")
	if err == nil || !strings.Contains(err.Error(), "ensure local server: model missing") {
		t.Fatalf("error = %v", err)
	}
}

func TestWhisperCPPServerStartsOnlyWhenEndpointIsUnavailable(t *testing.T) {
	t.Parallel()

	reachable := false
	launched := false
	server := &whisperCPPServer{
		config: config.WhisperCPPConfig{
			Host:           "127.0.0.1",
			Port:           8188,
			ModelPath:      "model.bin",
			StartupTimeout: "1s",
		},
		language:  "auto",
		reachable: func(context.Context, string) bool { return reachable },
		validate:  func() (string, error) { return "whisper-server", nil },
		launch: func(binary string, arguments []string) error {
			if binary != "whisper-server" {
				t.Fatalf("binary = %q", binary)
			}
			if got := strings.Join(arguments, " "); !strings.Contains(got, "--inference-path /v1/audio/transcriptions") {
				t.Fatalf("arguments = %q", got)
			}
			launched = true
			reachable = true
			return nil
		},
		interval: time.Millisecond,
	}
	if err := server.Ensure(context.Background()); err != nil {
		t.Fatalf("Ensure returned error: %v", err)
	}
	if !launched {
		t.Fatal("expected local server to be launched")
	}
}

func TestWhisperCPPServerReusesReachableEndpoint(t *testing.T) {
	t.Parallel()

	server := &whisperCPPServer{
		config:    config.WhisperCPPConfig{Host: "127.0.0.1", Port: 8188},
		reachable: func(context.Context, string) bool { return true },
		validate: func() (string, error) {
			t.Fatal("validate should not run when the server is reachable")
			return "", nil
		},
		launch: func(string, []string) error {
			t.Fatal("launch should not run when the server is reachable")
			return nil
		},
	}
	if err := server.Ensure(context.Background()); err != nil {
		t.Fatalf("Ensure returned error: %v", err)
	}
}
