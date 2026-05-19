package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/compozy/kb/internal/config"
)

const openAITranscriptionsPath = "/v1/audio/transcriptions"

type openAITranscriptionResponse struct {
	Text  string               `json:"text"`
	Error *openAIResponseError `json:"error"`
}

type openAIResponseError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// OpenAITranscriber calls OpenAI's audio transcriptions endpoint.
type OpenAITranscriber struct {
	apiKey   string
	apiURL   string
	model    string
	language string
	prompt   string

	httpClient *http.Client
}

// NewOpenAITranscriber constructs an OpenAI STT provider from runtime config.
func NewOpenAITranscriber(cfg config.STTConfig) *OpenAITranscriber {
	defaults := config.Default().STT
	apiURL := strings.TrimSpace(cfg.APIURL)
	if apiURL == "" {
		apiURL = defaults.APIURL
	}
	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = defaults.Model
	}
	language := strings.TrimSpace(cfg.Language)
	if language == "" {
		language = defaults.Language
	}
	return &OpenAITranscriber{
		apiKey:     strings.TrimSpace(cfg.APIKey),
		apiURL:     strings.TrimRight(apiURL, "/"),
		model:      model,
		language:   language,
		prompt:     strings.TrimSpace(cfg.Prompt),
		httpClient: http.DefaultClient,
	}
}

func (client *OpenAITranscriber) Configured() bool {
	return client != nil && strings.TrimSpace(client.apiKey) != ""
}

func (client *OpenAITranscriber) Provider() string {
	return "openai"
}

func (client *OpenAITranscriber) Model() string {
	if client == nil {
		return ""
	}
	return client.model
}

func (client *OpenAITranscriber) Transcribe(ctx context.Context, audio []byte, format string) (string, error) {
	if client == nil {
		return "", errors.New("openai transcribe: client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(client.apiKey) == "" {
		return "", errors.New("openai transcribe: missing API key; set stt.api_key or OPENAI_API_KEY")
	}
	if len(audio) == 0 {
		return "", errors.New("openai transcribe: audio is required")
	}
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		return "", errors.New("openai transcribe: audio format is required")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", "audio."+format)
	if err != nil {
		return "", fmt.Errorf("openai transcribe: create file part: %w", err)
	}
	if _, err := fileWriter.Write(audio); err != nil {
		return "", fmt.Errorf("openai transcribe: write file part: %w", err)
	}
	if err := writer.WriteField("model", client.model); err != nil {
		return "", fmt.Errorf("openai transcribe: write model field: %w", err)
	}
	if err := writer.WriteField("response_format", "json"); err != nil {
		return "", fmt.Errorf("openai transcribe: write response format field: %w", err)
	}
	if language := strings.TrimSpace(client.language); language != "" && !strings.EqualFold(language, "auto") {
		if err := writer.WriteField("language", language); err != nil {
			return "", fmt.Errorf("openai transcribe: write language field: %w", err)
		}
	}
	if client.prompt != "" {
		if err := writer.WriteField("prompt", client.prompt); err != nil {
			return "", fmt.Errorf("openai transcribe: write prompt field: %w", err)
		}
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("openai transcribe: close multipart body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpointURL(), &body)
	if err != nil {
		return "", fmt.Errorf("openai transcribe: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	httpClient := client.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", fmt.Errorf("openai transcribe: request canceled: %w", ctxErr)
		}
		return "", fmt.Errorf("openai transcribe: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openai transcribe: read response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("openai transcribe: request failed with status %d: %s", resp.StatusCode, parseOpenAIError(responseBody))
	}

	var payload openAITranscriptionResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return "", fmt.Errorf("openai transcribe: parse response: %w", err)
	}
	if payload.Error != nil && strings.TrimSpace(payload.Error.Message) != "" {
		return "", fmt.Errorf("openai transcribe: api error: %s", payload.Error.Message)
	}
	text := strings.TrimSpace(payload.Text)
	if text == "" {
		return "", errors.New("openai transcribe: empty transcription response")
	}
	return text, nil
}

func (client *OpenAITranscriber) endpointURL() string {
	apiURL := strings.TrimRight(client.apiURL, "/")
	if strings.HasSuffix(apiURL, "/v1") {
		return apiURL + "/audio/transcriptions"
	}
	return apiURL + openAITranscriptionsPath
}

func parseOpenAIError(body []byte) string {
	var payload openAITranscriptionResponse
	if err := json.Unmarshal(body, &payload); err == nil {
		if payload.Error != nil && strings.TrimSpace(payload.Error.Message) != "" {
			return strings.TrimSpace(payload.Error.Message)
		}
	}
	return strings.TrimSpace(string(body))
}
