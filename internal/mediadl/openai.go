package mediadl

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
	return transcribeMultipart(
		ctx,
		"openai",
		client.apiKey,
		true,
		client.endpointURL(),
		client.model,
		client.language,
		client.prompt,
		client.httpClient,
		audio,
		format,
	)
}

func transcribeMultipart(
	ctx context.Context,
	provider string,
	apiKey string,
	requireAPIKey bool,
	endpoint string,
	model string,
	language string,
	prompt string,
	httpClient *http.Client,
	audio []byte,
	format string,
) (string, error) {
	prefix := provider + " transcribe"
	if ctx == nil {
		ctx = context.Background()
	}
	if requireAPIKey && strings.TrimSpace(apiKey) == "" {
		return "", fmt.Errorf("%s: missing API key; set stt.api_key or OPENAI_API_KEY", prefix)
	}
	if len(audio) == 0 {
		return "", fmt.Errorf("%s: audio is required", prefix)
	}
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		return "", fmt.Errorf("%s: audio format is required", prefix)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", "audio."+format)
	if err != nil {
		return "", fmt.Errorf("%s: create file part: %w", prefix, err)
	}
	if _, err := fileWriter.Write(audio); err != nil {
		return "", fmt.Errorf("%s: write file part: %w", prefix, err)
	}
	if err := writer.WriteField("model", model); err != nil {
		return "", fmt.Errorf("%s: write model field: %w", prefix, err)
	}
	if err := writer.WriteField("response_format", "json"); err != nil {
		return "", fmt.Errorf("%s: write response format field: %w", prefix, err)
	}
	if language := strings.TrimSpace(language); language != "" && !strings.EqualFold(language, "auto") {
		if err := writer.WriteField("language", language); err != nil {
			return "", fmt.Errorf("%s: write language field: %w", prefix, err)
		}
	}
	if prompt != "" {
		if err := writer.WriteField("prompt", prompt); err != nil {
			return "", fmt.Errorf("%s: write prompt field: %w", prefix, err)
		}
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("%s: close multipart body: %w", prefix, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &body)
	if err != nil {
		return "", fmt.Errorf("%s: build request: %w", prefix, err)
	}
	if requireAPIKey {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", fmt.Errorf("%s: request canceled: %w", prefix, ctxErr)
		}
		return "", fmt.Errorf("%s: request failed: %w", prefix, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: read response: %w", prefix, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("%s: request failed with status %d: %s", prefix, resp.StatusCode, parseOpenAIError(responseBody))
	}

	var payload openAITranscriptionResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return "", fmt.Errorf("%s: parse response: %w", prefix, err)
	}
	if payload.Error != nil && strings.TrimSpace(payload.Error.Message) != "" {
		return "", fmt.Errorf("%s: api error: %s", prefix, payload.Error.Message)
	}
	text := strings.TrimSpace(payload.Text)
	if text == "" {
		return "", fmt.Errorf("%s: empty transcription response", prefix)
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
