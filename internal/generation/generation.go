// Package generation is kb's client for the generation model: short literals
// (summaries, aliases, criteria, drafts) over OpenRouter chat completions
// with a strict JSON schema and reasoning disabled. It shares the run budget
// and the receipts log with internal/decisions; every literal it returns is
// still validated by the caller's code before use.
package generation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/decisions"
)

// Failure reasons recorded in receipts and wrapped in ErrFailed.
const (
	ReasonInvalidOutput = "invalid_output"
	ReasonTimeout       = decisions.ReasonTimeout
	ReasonRetries       = decisions.ReasonRetries
	ReasonHTTPStatus    = decisions.ReasonHTTPStatus
	ReasonBudget        = decisions.ReasonBudget
)

const (
	backoffBase     = 500 * time.Millisecond
	backoffCap      = 5 * time.Second
	maxRetryAfter   = 30 * time.Second
	maxResponseSize = 32 << 20
)

var (
	// ErrNotConfigured is returned by New when no OpenRouter API key is set.
	ErrNotConfigured = errors.New("generation: OPENROUTER_API_KEY is not set; the generation model is required ([generation] + [openrouter])")

	// ErrBudget is returned when the shared run budget is exhausted before a
	// call; nothing is sent.
	ErrBudget = errors.New("generation: run budget exhausted")

	// ErrFailed is returned when every attempt on every model failed; the
	// wrapped message names the last reason.
	ErrFailed = errors.New("generation: no valid output")

	// ErrInvalidRequest marks a request that cannot be sent (missing kind,
	// prompt, schema name, or a schema that is not a JSON object schema).
	ErrInvalidRequest = errors.New("generation: invalid request")
)

// Options configure a Client. Only Config and APIKey are required. Share
// Budget and Receipts with the decision engine (Engine.Budget,
// Engine.Receipts) so a run has one ceiling and one log per topic.
type Options struct {
	Config     config.GenerationConfig
	APIKey     string
	APIURL     string
	HTTPClient *http.Client
	Budget     *decisions.Budget
	Receipts   *decisions.Receipts
	Now        func() time.Time
	Sleep      func(ctx context.Context, d time.Duration) error
	Logger     *slog.Logger
}

// Request is one generation request.
type Request struct {
	// Topic locates the receipts log (Topic.Root) and records the contract.
	Topic decisions.TopicRef
	// Kind names the literal (summary, aliases, criterion, contract_draft,
	// ...); receipts record purpose "generate:<kind>".
	Kind string
	// Subject is the document path or id, for receipts.
	Subject string
	System  string
	Prompt  string
	// SchemaName and Schema form the strict json_schema response format. The
	// schema must describe a JSON object.
	SchemaName string
	Schema     json.RawMessage
}

// Client generates schema-checked JSON. It is safe for concurrent use.
type Client struct {
	model    string
	fallback string
	endpoint string
	apiKey   string
	deadline time.Duration
	retries  int
	client   *http.Client
	budget   *decisions.Budget
	receipts *decisions.Receipts
	now      func() time.Time
	sleep    func(ctx context.Context, d time.Duration) error
	logger   *slog.Logger

	mu      sync.Mutex
	summary decisions.GenerationSummary
	fatal   error
}

// New builds a Client. It returns ErrNotConfigured when opts.APIKey is empty.
func New(opts Options) (*Client, error) {
	if strings.TrimSpace(opts.APIKey) == "" {
		return nil, ErrNotConfigured
	}
	if err := opts.Config.Validate(); err != nil {
		return nil, fmt.Errorf("generation: %w", err)
	}
	deadline, err := opts.Config.DeadlineDuration()
	if err != nil {
		return nil, fmt.Errorf("generation: %w", err)
	}
	apiURL := strings.TrimRight(strings.TrimSpace(opts.APIURL), "/")
	if apiURL == "" {
		apiURL = decisions.DefaultAPIURL
	}
	client := &Client{
		model:    strings.TrimSpace(opts.Config.Model),
		fallback: strings.TrimSpace(opts.Config.FallbackModel),
		endpoint: apiURL + "/v1/chat/completions",
		apiKey:   strings.TrimSpace(opts.APIKey),
		deadline: deadline,
		retries:  opts.Config.Retries,
		client:   opts.HTTPClient,
		budget:   opts.Budget,
		receipts: opts.Receipts,
		now:      opts.Now,
		sleep:    opts.Sleep,
		logger:   opts.Logger,
	}
	if client.fallback == client.model {
		client.fallback = ""
	}
	if client.client == nil {
		client.client = &http.Client{}
	}
	if client.receipts == nil {
		client.receipts = decisions.NewReceipts()
	}
	if client.now == nil {
		client.now = time.Now
	}
	if client.sleep == nil {
		client.sleep = decisions.SleepContext
	}
	if client.logger == nil {
		client.logger = slog.New(slog.DiscardHandler)
	}
	return client, nil
}

// Model returns the primary generation model.
func (c *Client) Model() string { return c.model }

// Summary returns the client's part of the run summary; set it as
// decisions.Summary.Generation before printing Summary.Lines.
func (c *Client) Summary() decisions.GenerationSummary {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.summary
}

// Generate returns the model's JSON output for req, validated against the
// top level of req.Schema. Identical requests (models, system, prompt,
// schema) are served from the topic's receipts. Each model gets 1+retries
// attempts under the per-attempt deadline, then the fallback model gets the
// same. Errors: ErrBudget (nothing sent), ErrFailed (every attempt failed),
// decisions.ErrAuth (401/402/403, sticky), ErrInvalidRequest, ctx errors.
func (c *Client) Generate(ctx context.Context, req Request) (json.RawMessage, error) {
	if err := c.fatalErr(); err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	schema, err := parseRequest(req)
	if err != nil {
		return nil, err
	}
	system := decisions.RedactString(req.System)
	prompt := decisions.RedactString(req.Prompt)
	schemaTree, err := decisions.CanonicalJSON(req.Schema)
	if err != nil {
		return nil, fmt.Errorf("%w: schema: %v", ErrInvalidRequest, err)
	}
	key, err := decisions.HashJSON(map[string]any{
		"route":       decisions.RouteChatCompletions,
		"model":       c.model,
		"fallback":    c.fallback,
		"system":      system,
		"prompt":      prompt,
		"schema_name": req.SchemaName,
		"schema":      json.RawMessage(schemaTree),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}
	if output, ok := c.fromCache(req.Topic.Root, key, schema); ok {
		c.record(func(s *decisions.GenerationSummary) { s.CacheHits++ })
		return output, nil
	}

	run := attemptLog{started: c.now(), lastReason: ReasonRetries}
	models := []string{c.model}
	if c.fallback != "" {
		models = append(models, c.fallback)
	}
	for index, model := range models {
		if index > 0 {
			c.record(func(s *decisions.GenerationSummary) { s.Fallbacks++ })
		}
		output, err := c.tryModel(ctx, model, system, prompt, req, schemaTree, schema, &run)
		if err != nil {
			if errors.Is(err, ErrBudget) {
				c.finish(req, key, &run, nil)
			}
			return nil, err
		}
		if output != nil {
			c.finish(req, key, &run, output)
			return output, nil
		}
	}
	c.finish(req, key, &run, nil)
	return nil, fmt.Errorf("%w: %s", ErrFailed, run.lastReason)
}

// attemptLog aggregates the attempts of one Generate call for its receipt.
type attemptLog struct {
	started     time.Time
	attempts    int
	reported    string
	responseID  string
	inputTokens *int64
	cost        *float64
	costUnknown bool
	lastReason  string
}

func (c *Client) tryModel(ctx context.Context, model, system, prompt string, req Request, schemaJSON []byte, schema *objectSchema, run *attemptLog) (json.RawMessage, error) {
	payload, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":   req.SchemaName,
				"strict": true,
				"schema": json.RawMessage(schemaJSON),
			},
		},
		"reasoning":   map[string]any{"enabled": false},
		"temperature": 0,
		"usage":       map[string]any{"include": true},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: encode request: %v", ErrInvalidRequest, err)
	}
	var retryWait time.Duration
	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			wait := retryWait
			if wait <= 0 {
				wait = backoff(attempt)
			}
			if err := c.sleep(ctx, wait); err != nil {
				return nil, err
			}
		}
		retryWait = 0
		if c.budget.Exceeded() {
			run.lastReason = ReasonBudget
			return nil, ErrBudget
		}
		status, header, body, err := c.post(ctx, payload)
		run.attempts++
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			run.lastReason = ReasonRetries
			if errors.Is(err, context.DeadlineExceeded) {
				run.lastReason = ReasonTimeout
			}
			continue
		}
		switch {
		case status == http.StatusOK:
			output, reason := c.readResponse(body, schema, run)
			if reason == "" {
				return output, nil
			}
			run.lastReason = reason
			c.logger.Debug("generation: output rejected", "model", model, "reason", reason)
		case status == http.StatusUnauthorized || status == http.StatusPaymentRequired || status == http.StatusForbidden:
			return nil, c.setFatal(fmt.Errorf("%w: generation HTTP %d: %s", decisions.ErrAuth, status, errorMessage(body)))
		case retryable(status):
			run.lastReason = ReasonRetries
			if wait, ok := retryAfter(header, c.now()); ok {
				retryWait = wait
			}
		default:
			// A model-specific rejection (for example an unsupported schema)
			// moves on to the fallback model.
			run.lastReason = ReasonHTTPStatus
			c.logger.Warn("generation: request rejected", "model", model, "status", status, "message", errorMessage(body))
			return nil, nil
		}
	}
	return nil, nil
}

// chatResponse is the subset of a chat-completions body kb reads.
type chatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content *string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens            *int64          `json:"prompt_tokens"`
		Cost                    json.RawMessage `json:"cost"`
		CompletionTokensDetails *struct {
			ReasoningTokens int64 `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
}

// readResponse charges the budget for a 200 body and returns the validated
// output, or the rejection reason.
func (c *Client) readResponse(body []byte, schema *objectSchema, run *attemptLog) (json.RawMessage, string) {
	var response chatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		run.costUnknown = true
		return nil, ReasonInvalidOutput
	}
	run.reported = response.Model
	run.responseID = response.ID
	var cost *float64
	if response.Usage != nil {
		cost = decisions.ParseCost(response.Usage.Cost)
		if response.Usage.PromptTokens != nil {
			tokens := *response.Usage.PromptTokens
			if run.inputTokens != nil {
				tokens += *run.inputTokens
			}
			run.inputTokens = &tokens
		}
	}
	if cost == nil {
		run.costUnknown = true
	} else {
		total := *cost
		if run.cost != nil {
			total += *run.cost
		}
		run.cost = &total
		c.budget.Charge(*cost)
	}
	if response.Usage != nil && response.Usage.CompletionTokensDetails != nil && response.Usage.CompletionTokensDetails.ReasoningTokens > 0 {
		return nil, ReasonInvalidOutput
	}
	if len(response.Choices) == 0 || response.Choices[0].Message.Content == nil {
		return nil, ReasonInvalidOutput
	}
	output, err := validateOutput(*response.Choices[0].Message.Content, schema)
	if err != nil {
		return nil, ReasonInvalidOutput
	}
	return output, ""
}

// finish writes the receipt row of a Generate call that reached the network
// and updates the summary.
func (c *Client) finish(req Request, key string, run *attemptLog, output json.RawMessage) {
	if run.attempts == 0 {
		return
	}
	row := decisions.Receipt{
		Key:           key,
		Time:          c.now().UTC().Format(time.RFC3339Nano),
		Purpose:       "generate:" + req.Kind,
		Subject:       req.Subject,
		Contract:      req.Topic.Contract,
		Model:         c.model,
		ReportedModel: run.reported,
		Route:         decisions.RouteChatCompletions,
		ResponseID:    run.responseID,
		Attempts:      run.attempts,
		LatencyMS:     c.now().Sub(run.started).Milliseconds(),
		InputTokens:   run.inputTokens,
		Cost:          run.cost,
		CostUnknown:   run.costUnknown || run.cost == nil,
		Status:        decisions.StatusUndecided,
		Reason:        run.lastReason,
		Questions:     []string{"output"},
	}
	if output != nil {
		row.Status, row.Reason = decisions.StatusDecided, ""
		row.Answers = map[string]json.RawMessage{"output": output}
	}
	c.record(func(s *decisions.GenerationSummary) {
		s.Calls++
		if run.cost != nil {
			s.CostUSD += *run.cost
		}
		if row.CostUnknown {
			s.CostUnknown++
		}
		if output == nil {
			s.Failures++
		}
	})
	if err := c.receipts.Append(req.Topic.Root, row); err != nil {
		c.logger.Warn("generation: receipt not written", "topic", req.Topic.Slug, "key", key, "error", err)
	}
}

func (c *Client) fromCache(topicRoot, key string, schema *objectSchema) (json.RawMessage, bool) {
	row, ok, err := c.receipts.Lookup(topicRoot, key)
	if err != nil {
		c.logger.Warn("generation: receipts unreadable, cache disabled for this lookup", "root", topicRoot, "error", err)
		return nil, false
	}
	if !ok {
		return nil, false
	}
	raw, ok := row.Answers["output"]
	if !ok {
		return nil, false
	}
	output, err := validateOutput(string(raw), schema)
	if err != nil {
		return nil, false
	}
	return output, true
}

func (c *Client) record(update func(*decisions.GenerationSummary)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	update(&c.summary)
}

func (c *Client) fatalErr() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.fatal
}

func (c *Client) setFatal(err error) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.fatal == nil {
		c.fatal = err
	}
	return c.fatal
}

func (c *Client) post(ctx context.Context, payload []byte) (int, http.Header, []byte, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, c.deadline)
	defer cancel()
	request, err := http.NewRequestWithContext(attemptCtx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return 0, nil, nil, attemptError(attemptCtx, err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize))
	if err != nil {
		return 0, nil, nil, attemptError(attemptCtx, err)
	}
	return response.StatusCode, response.Header, body, nil
}

func attemptError(attemptCtx context.Context, err error) error {
	if errors.Is(attemptCtx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("%w: %v", context.DeadlineExceeded, err)
	}
	return err
}

func retryable(status int) bool {
	switch status {
	case http.StatusRequestTimeout, http.StatusTooManyRequests:
		return true
	default:
		return status >= 500 && status <= 599
	}
}

func backoff(retry int) time.Duration {
	delay := backoffBase << max(retry-1, 0)
	if delay > backoffCap || delay <= 0 {
		delay = backoffCap
	}
	return time.Duration(float64(delay) * (0.75 + rand.Float64()*0.5))
}

func retryAfter(header http.Header, now time.Time) (time.Duration, bool) {
	if header == nil {
		return 0, false
	}
	if ms, err := strconv.ParseFloat(strings.TrimSpace(header.Get("retry-after-ms")), 64); err == nil && ms >= 0 {
		return min(time.Duration(ms*float64(time.Millisecond)), maxRetryAfter), true
	}
	value := strings.TrimSpace(header.Get("Retry-After"))
	if seconds, err := strconv.ParseFloat(value, 64); err == nil && seconds >= 0 {
		return min(time.Duration(seconds*float64(time.Second)), maxRetryAfter), true
	}
	if date, err := http.ParseTime(value); err == nil {
		return min(max(date.Sub(now), 0), maxRetryAfter), true
	}
	return 0, false
}

func errorMessage(body []byte) string {
	var envelope struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Error.Message != "" {
		return envelope.Error.Message
	}
	text := strings.TrimSpace(string(body))
	if len(text) > 200 {
		text = text[:200]
	}
	return text
}
