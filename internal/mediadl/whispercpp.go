package mediadl

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/config"
)

const whisperCPPTranscriptionsPath = "/v1/audio/transcriptions"

// WhisperCPPTranscriber calls a locally managed whisper.cpp server.
type WhisperCPPTranscriber struct {
	endpoint string
	model    string
	language string
	prompt   string
	starter  whisperCPPStarter

	httpClient *http.Client
}

type whisperCPPStarter interface {
	Ensure(context.Context) error
}

// NewWhisperCPPTranscriber constructs the local whisper.cpp STT provider.
func NewWhisperCPPTranscriber(sttConfig config.STTConfig, whisperConfig config.WhisperCPPConfig) *WhisperCPPTranscriber {
	sttConfig = normalizeSTTConfig(sttConfig)
	return &WhisperCPPTranscriber{
		endpoint:   whisperCPPEndpoint(whisperConfig),
		model:      sttConfig.Model,
		language:   sttConfig.Language,
		prompt:     sttConfig.Prompt,
		starter:    newWhisperCPPServer(whisperConfig, sttConfig.Language),
		httpClient: http.DefaultClient,
	}
}

func whisperCPPEndpoint(cfg config.WhisperCPPConfig) string {
	host := strings.TrimSpace(cfg.Host)
	if host == "" {
		host = "127.0.0.1"
	}
	port := cfg.Port
	if port == 0 {
		port = 8188
	}
	return "http://" + net.JoinHostPort(host, strconv.Itoa(port)) + whisperCPPTranscriptionsPath
}

func (client *WhisperCPPTranscriber) Provider() string {
	return "whispercpp"
}

func (client *WhisperCPPTranscriber) Model() string {
	if client == nil {
		return ""
	}
	return client.model
}

func (client *WhisperCPPTranscriber) Transcribe(ctx context.Context, audio []byte, format string) (string, error) {
	if client == nil {
		return "", errors.New("whispercpp transcribe: client is nil")
	}
	if client.starter == nil {
		return "", errors.New("whispercpp transcribe: server starter is nil")
	}
	if err := client.starter.Ensure(ctx); err != nil {
		return "", fmt.Errorf("whispercpp transcribe: ensure local server: %w", err)
	}
	return transcribeMultipart(
		ctx,
		"whispercpp",
		"",
		false,
		client.endpoint,
		client.model,
		client.language,
		client.prompt,
		client.httpClient,
		audio,
		format,
	)
}

type whisperCPPServer struct {
	config   config.WhisperCPPConfig
	language string

	mu        sync.Mutex
	reachable func(context.Context, string) bool
	validate  func() (string, error)
	launch    func(string, []string) error
	interval  time.Duration
}

func newWhisperCPPServer(cfg config.WhisperCPPConfig, language string) *whisperCPPServer {
	return &whisperCPPServer{
		config:    cfg,
		language:  language,
		reachable: whisperCPPReachable,
		validate: func() (string, error) {
			return validateWhisperCPPPaths(cfg)
		},
		launch:   launchWhisperCPPServer,
		interval: 100 * time.Millisecond,
	}
}

func (server *whisperCPPServer) Ensure(ctx context.Context) error {
	if server == nil {
		return errors.New("whispercpp server is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	server.mu.Lock()
	defer server.mu.Unlock()

	address := net.JoinHostPort(strings.TrimSpace(server.config.Host), strconv.Itoa(server.config.Port))
	if server.reachable(ctx, address) {
		return nil
	}
	binary, err := server.validate()
	if err != nil {
		return err
	}
	if err := server.launch(binary, server.arguments()); err != nil {
		return fmt.Errorf("start whispercpp server: %w", err)
	}

	timeout, err := server.config.StartupTimeoutValue()
	if err != nil {
		return err
	}
	startupCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	ticker := time.NewTicker(server.interval)
	defer ticker.Stop()
	for {
		if server.reachable(startupCtx, address) {
			return nil
		}
		select {
		case <-startupCtx.Done():
			return fmt.Errorf("whispercpp server did not become reachable at %s within %s: %w", address, timeout, startupCtx.Err())
		case <-ticker.C:
		}
	}
}

func (server *whisperCPPServer) arguments() []string {
	language := strings.TrimSpace(server.language)
	if language == "" {
		language = "auto"
	}
	return []string{
		"--model", server.config.ModelPath,
		"--host", server.config.Host,
		"--port", strconv.Itoa(server.config.Port),
		"--inference-path", whisperCPPTranscriptionsPath,
		"--convert",
		"--language", language,
	}
}

func whisperCPPReachable(ctx context.Context, address string) bool {
	connection, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		return false
	}
	_ = connection.Close()
	return true
}

func validateWhisperCPPPaths(cfg config.WhisperCPPConfig) (string, error) {
	binary, err := exec.LookPath(strings.TrimSpace(cfg.ServerPath))
	if err != nil {
		return "", fmt.Errorf("whispercpp server not found; install whisper.cpp or configure whispercpp.server_path: %w", err)
	}
	info, err := os.Stat(strings.TrimSpace(cfg.ModelPath))
	if err != nil {
		return "", fmt.Errorf("whispercpp model not found at %q: %w", cfg.ModelPath, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("whispercpp model path %q is a directory", cfg.ModelPath)
	}
	return binary, nil
}

func launchWhisperCPPServer(binary string, arguments []string) error {
	command := exec.Command(binary, arguments...)
	if err := command.Start(); err != nil {
		return err
	}
	go func() {
		_ = command.Wait()
	}()
	return nil
}
