package youtube

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/compozy/kb/internal/config"
)

const (
	openRouterChatCompletionPath = "/v1/chat/completions"
	defaultTranscriptionPrompt   = "Transcribe this audio verbatim. Return only the transcript text."
)

type chatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []chatCompletionMessage `json:"messages"`
	Stream   bool                    `json:"stream"`
}

type chatCompletionMessage struct {
	Role    string                      `json:"role"`
	Content []chatCompletionContentPart `json:"content"`
}

type chatCompletionContentPart struct {
	Type       string                  `json:"type"`
	Text       string                  `json:"text,omitempty"`
	InputAudio *chatCompletionAudioRef `json:"input_audio,omitempty"`
}

type chatCompletionAudioRef struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

type chatCompletionResponse struct {
	Error   *chatCompletionResponseError `json:"error"`
	Choices []struct {
		Error   *chatCompletionResponseError `json:"error"`
		Message struct {
			Content json.RawMessage `json:"content"`
			Role    string          `json:"role"`
		} `json:"message"`
	} `json:"choices"`
}

type chatCompletionResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type chatCompletionTextPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// OpenRouterClient calls the OpenRouter chat completions API for STT fallback.
type OpenRouterClient struct {
	apiKey string
	apiURL string
	model  string

	httpClient *http.Client
}

// NewOpenRouterClient constructs an STT client from runtime configuration.
func NewOpenRouterClient(cfg config.OpenRouterConfig) *OpenRouterClient {
	defaults := config.Default().OpenRouter

	apiURL := strings.TrimSpace(cfg.APIURL)
	if apiURL == "" {
		apiURL = defaults.APIURL
	}

	model := strings.TrimSpace(cfg.STTModel)
	if model == "" {
		model = defaults.STTModel
	}

	return &OpenRouterClient{
		apiKey:     strings.TrimSpace(cfg.APIKey),
		apiURL:     strings.TrimRight(apiURL, "/"),
		model:      model,
		httpClient: http.DefaultClient,
	}
}

// Configured reports whether the client has the credentials needed to call the
// OpenRouter API.
func (client *OpenRouterClient) Configured() bool {
	return client != nil && strings.TrimSpace(client.apiKey) != ""
}

// Transcribe sends audio bytes to OpenRouter and returns the transcript text.
func (client *OpenRouterClient) Transcribe(ctx context.Context, audio []byte, format string) (string, error) {
	if client == nil {
		return "", errors.New("openrouter transcribe: client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(client.apiKey) == "" {
		return "", errors.New("openrouter transcribe: missing API key; set openrouter.api_key or OPENROUTER_API_KEY")
	}
	if len(audio) == 0 {
		return "", errors.New("openrouter transcribe: audio is required")
	}

	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		return "", errors.New("openrouter transcribe: audio format is required")
	}

	body, err := json.Marshal(chatCompletionRequest{
		Model: client.model,
		Messages: []chatCompletionMessage{
			{
				Role: "user",
				Content: []chatCompletionContentPart{
					{
						Type: "text",
						Text: defaultTranscriptionPrompt,
					},
					{
						Type: "input_audio",
						InputAudio: &chatCompletionAudioRef{
							Data:   base64.StdEncoding.EncodeToString(audio),
							Format: format,
						},
					},
				},
			},
		},
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("openrouter transcribe: encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpointURL(), bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("openrouter transcribe: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := client.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", fmt.Errorf("openrouter transcribe: request canceled: %w", ctxErr)
		}
		return "", fmt.Errorf("openrouter transcribe: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openrouter transcribe: read response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("openrouter transcribe: request failed with status %d: %s", resp.StatusCode, parseOpenRouterError(responseBody))
	}

	var payload chatCompletionResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return "", fmt.Errorf("openrouter transcribe: parse response: %w", err)
	}

	if payload.Error != nil && strings.TrimSpace(payload.Error.Message) != "" {
		return "", fmt.Errorf("openrouter transcribe: api error: %s", payload.Error.Message)
	}

	if len(payload.Choices) == 0 {
		return "", errors.New("openrouter transcribe: empty response choices")
	}
	if payload.Choices[0].Error != nil && strings.TrimSpace(payload.Choices[0].Error.Message) != "" {
		return "", fmt.Errorf("openrouter transcribe: api error: %s", payload.Choices[0].Error.Message)
	}

	transcript, err := parseTranscriptionContent(payload.Choices[0].Message.Content)
	if err != nil {
		return "", fmt.Errorf("openrouter transcribe: %w", err)
	}

	return transcript, nil
}

func (client *OpenRouterClient) endpointURL() string {
	return client.apiURL + openRouterChatCompletionPath
}

func parseOpenRouterError(body []byte) string {
	var payload chatCompletionResponse
	if err := json.Unmarshal(body, &payload); err == nil {
		if payload.Error != nil && strings.TrimSpace(payload.Error.Message) != "" {
			return strings.TrimSpace(payload.Error.Message)
		}
		if len(payload.Choices) > 0 && payload.Choices[0].Error != nil && strings.TrimSpace(payload.Choices[0].Error.Message) != "" {
			return strings.TrimSpace(payload.Choices[0].Error.Message)
		}
	}

	return strings.TrimSpace(string(body))
}

func parseTranscriptionContent(content json.RawMessage) (string, error) {
	if len(content) == 0 || string(content) == "null" {
		return "", errors.New("empty transcription response")
	}

	var rawText string
	if err := json.Unmarshal(content, &rawText); err == nil {
		rawText = strings.TrimSpace(rawText)
		if rawText == "" {
			return "", errors.New("empty transcription response")
		}
		return rawText, nil
	}

	var parts []chatCompletionTextPart
	if err := json.Unmarshal(content, &parts); err == nil {
		textParts := make([]string, 0, len(parts))
		for _, part := range parts {
			if strings.TrimSpace(part.Text) == "" {
				continue
			}
			textParts = append(textParts, strings.TrimSpace(part.Text))
		}

		if len(textParts) == 0 {
			return "", errors.New("empty transcription response")
		}

		return strings.Join(textParts, " "), nil
	}

	var part chatCompletionTextPart
	if err := json.Unmarshal(content, &part); err == nil && strings.TrimSpace(part.Text) != "" {
		return strings.TrimSpace(part.Text), nil
	}

	return "", errors.New("empty transcription response")
}
