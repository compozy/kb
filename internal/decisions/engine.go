package decisions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/questions"
)

const (
	// DefaultAPIURL is OpenRouter's API base; the engine posts to
	// <api_url>/alpha/decisions.
	DefaultAPIURL = "https://openrouter.ai/api"

	// MaxQuestionsPerRequest is the batch size: larger question sets are split
	// into several requests over the same state.
	MaxQuestionsPerRequest = 48

	// MaxStateTokens caps the estimated tokens of one state. The estimate is
	// bytes/4 of the canonical JSON (about four bytes per token for English
	// text and JSON punctuation); it only matters when max_state_bytes is
	// raised above 120 KB, and it keeps the state plus the longest question
	// under the model's 32k-token limit.
	MaxStateTokens = 30000

	// listPricePerInputToken is used to charge the budget when a response
	// carries tokens but no cost (US$ 0.042 per million input tokens,
	// checked 2026-09-22). The receipt still records cost_unknown.
	listPricePerInputToken = 0.042 / 1_000_000

	backoffBase     = 500 * time.Millisecond
	backoffCap      = 5 * time.Second
	maxRetryAfter   = 30 * time.Second
	maxResponseSize = 32 << 20
)

var (
	// ErrNotConfigured is returned by New when no OpenRouter API key is set.
	ErrNotConfigured = errors.New("decisions: OPENROUTER_API_KEY is not set; the decision model is required ([decisions] + [openrouter])")

	// ErrAuth marks a 401/402/403 from the decision endpoint. After the first
	// one, every later Decide returns it immediately.
	ErrAuth = errors.New("decisions: the decision endpoint refused the key or credits")

	// ErrInvalidRequest marks a request the engine cannot send: no or
	// duplicate question ids, malformed criteria, or a state that is not
	// JSON-serializable. It is a caller bug, not a model outcome.
	ErrInvalidRequest = errors.New("decisions: invalid request")
)

// Options configure an Engine. Only Config and APIKey are required.
type Options struct {
	Config     config.DecisionsConfig
	APIKey     string
	APIURL     string
	HTTPClient *http.Client
	// Budget is the run budget; nil creates one from Config.BudgetUSD.
	Budget *Budget
	// Receipts is the receipts store; nil creates one. Share it with the
	// generation client through Engine.Receipts.
	Receipts *Receipts
	Now      func() time.Time
	// Sleep waits between retries; it must return ctx.Err() when ctx ends.
	Sleep  func(ctx context.Context, d time.Duration) error
	Logger *slog.Logger
}

// Engine asks the decision model typed questions and returns banded answers.
// It is safe for concurrent use.
type Engine struct {
	model    string
	endpoint string
	apiKey   string
	deadline time.Duration
	retries  int
	maxState int
	client   *http.Client
	budget   *Budget
	receipts *Receipts
	now      func() time.Time
	sleep    func(ctx context.Context, d time.Duration) error
	logger   *slog.Logger
	sem      chan struct{}
	stats    stats

	pauseMu    sync.Mutex
	pauseUntil time.Time

	fatalMu sync.Mutex
	fatal   error
}

// New builds an Engine. It returns ErrNotConfigured when opts.APIKey is empty
// and a validation error when opts.Config is unusable.
func New(opts Options) (*Engine, error) {
	if strings.TrimSpace(opts.APIKey) == "" {
		return nil, ErrNotConfigured
	}
	if err := opts.Config.Validate(); err != nil {
		return nil, fmt.Errorf("decisions: %w", err)
	}
	deadline, err := opts.Config.DeadlineDuration()
	if err != nil {
		return nil, fmt.Errorf("decisions: %w", err)
	}
	apiURL := strings.TrimRight(strings.TrimSpace(opts.APIURL), "/")
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}
	engine := &Engine{
		model:    strings.TrimSpace(opts.Config.Model),
		endpoint: apiURL + "/alpha/decisions",
		apiKey:   strings.TrimSpace(opts.APIKey),
		deadline: deadline,
		retries:  opts.Config.Retries,
		maxState: opts.Config.MaxStateBytes,
		client:   opts.HTTPClient,
		budget:   opts.Budget,
		receipts: opts.Receipts,
		now:      opts.Now,
		sleep:    opts.Sleep,
		logger:   opts.Logger,
		sem:      make(chan struct{}, opts.Config.Concurrency),
	}
	if engine.client == nil {
		engine.client = &http.Client{}
	}
	if engine.budget == nil {
		engine.budget = NewBudget(opts.Config.BudgetUSD)
	}
	if engine.receipts == nil {
		engine.receipts = NewReceipts()
	}
	if engine.now == nil {
		engine.now = time.Now
	}
	if engine.sleep == nil {
		engine.sleep = SleepContext
	}
	if engine.logger == nil {
		engine.logger = slog.New(slog.DiscardHandler)
	}
	return engine, nil
}

// SleepContext waits for d or until ctx ends, returning ctx.Err() then.
func SleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// Model returns the pinned decision model slug.
func (e *Engine) Model() string { return e.model }

// Concurrency returns the number of requests the engine runs at once, so
// callers can size their fan-out.
func (e *Engine) Concurrency() int { return cap(e.sem) }

// Budget returns the run budget shared with generation.
func (e *Engine) Budget() *Budget { return e.budget }

// Receipts returns the receipts store shared with generation.
func (e *Engine) Receipts() *Receipts { return e.receipts }

// Summary returns a snapshot of the run summary (Generation left empty).
func (e *Engine) Summary() Summary { return e.stats.snapshot() }

// Decide asks every question of r against r.State and returns one banded
// answer per question id. It returns an error only for fatal conditions:
// ErrAuth (401/402/403, sticky for the engine), ctx cancellation, and
// ErrInvalidRequest. Every other failure is a per-answer status:
// undecided (timeout, retries, budget, invalid_receipt, http_status,
// missing_answer) or not_checked (state_too_large, or excluded: a Subject
// matching r.Topic.Exclude is answered without any call).
func (e *Engine) Decide(ctx context.Context, r Request) (Result, error) {
	if err := e.fatalErr(); err != nil {
		return Result{}, err
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	prepared, err := prepareQuestions(r.Questions)
	if err != nil {
		return Result{}, err
	}
	result := Result{Answers: make(map[string]Answer, len(prepared))}
	if r.Topic.Excluded(r.Subject) {
		// decisions.exclude keeps the subject out of every call (spec §14):
		// no state is serialized, nothing is sent, no receipt is written.
		for _, pq := range prepared {
			result.Answers[pq.q.ID] = notChecked(pq.q.Type, ReasonExcluded)
		}
		e.stats.answers(r.Purpose, result.Answers)
		return result, nil
	}
	// Redaction happens before hashing so a secret never enters a cache key.
	stateJSON, err := CanonicalJSON(Redact(r.State))
	if err != nil {
		return Result{}, fmt.Errorf("%w: state: %v", ErrInvalidRequest, err)
	}

	if len(stateJSON) > e.maxState || (len(stateJSON)+3)/4 > MaxStateTokens {
		for _, pq := range prepared {
			result.Answers[pq.q.ID] = notChecked(pq.q.Type, ReasonStateTooLarge)
		}
		e.stats.answers(r.Purpose, result.Answers)
		return result, nil
	}

	apply, review := r.Topic.Thresholds.For(r.Purpose)
	allHits := true
	for start := 0; start < len(prepared); start += MaxQuestionsPerRequest {
		batch := prepared[start:min(start+MaxQuestionsPerRequest, len(prepared))]
		outcome, err := e.decideBatch(ctx, r, stateJSON, batch)
		if err != nil {
			return Result{}, err
		}
		if start == 0 {
			result.ReceiptID = outcome.key
		}
		allHits = allHits && outcome.cacheHit
		result.CostUSD += outcome.cost
		for id, answer := range outcome.answers {
			answer.Receipt = outcome.key
			result.Answers[id] = band(answer, apply, review)
		}
	}
	result.CacheHit = allHits
	e.stats.answers(r.Purpose, result.Answers)
	return result, nil
}

type batchOutcome struct {
	key      string
	answers  map[string]Answer
	cost     float64
	cacheHit bool
}

// keyQuestion is the hashed form of a question (Q hides its id from JSON).
type keyQuestion struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Instructions any    `json:"instructions"`
	Criteria     any    `json:"criteria,omitempty"`
}

func (e *Engine) requestKey(r Request, state json.RawMessage, batch []preparedQuestion) (string, error) {
	list := make([]keyQuestion, len(batch))
	for index, pq := range batch {
		list[index] = keyQuestion{ID: pq.q.ID, Type: pq.q.Type, Instructions: pq.q.Instructions, Criteria: pq.q.Criteria}
	}
	bankID, bankVersion, bankHash := bankFields(r.Bank)
	key, err := HashJSON(map[string]any{
		"model":        e.model,
		"bank":         bankID,
		"bank_version": bankVersion,
		"bank_hash":    bankHash,
		"contract":     r.Topic.Contract,
		"state":        state,
		"questions":    list,
	})
	if err != nil {
		return "", fmt.Errorf("%w: questions: %v", ErrInvalidRequest, err)
	}
	return key, nil
}

func bankFields(bank *questions.Bank) (id, version, hash string) {
	if bank == nil {
		return "", "", ""
	}
	return bank.ID, bank.Version, bank.Hash
}

func (e *Engine) decideBatch(ctx context.Context, r Request, state json.RawMessage, batch []preparedQuestion) (batchOutcome, error) {
	key, err := e.requestKey(r, state, batch)
	if err != nil {
		return batchOutcome{}, err
	}
	if cached, ok := e.fromCache(r.Topic.Root, key, batch); ok {
		e.stats.cacheHit()
		return batchOutcome{key: key, answers: cached, cacheHit: true}, nil
	}

	questionMap := make(map[string]questions.Q, len(batch))
	ids := make([]string, len(batch))
	for index, pq := range batch {
		questionMap[pq.q.ID] = pq.q
		ids[index] = pq.q.ID
	}
	payload, err := json.Marshal(map[string]any{"model": e.model, "state": state, "questions": questionMap})
	if err != nil {
		return batchOutcome{}, fmt.Errorf("%w: encode request: %v", ErrInvalidRequest, err)
	}
	// A slot covers dispatch through accounting, so queued calls cannot start
	// before the previous call's cost or fatal outcome has been recorded.
	select {
	case e.sem <- struct{}{}:
	case <-ctx.Done():
		return batchOutcome{}, ctx.Err()
	}
	defer func() { <-e.sem }()
	if err := ctx.Err(); err != nil {
		return batchOutcome{}, err
	}
	// An identical request may have populated the cache while this one waited.
	if cached, ok := e.fromCache(r.Topic.Root, key, batch); ok {
		e.stats.cacheHit()
		return batchOutcome{key: key, answers: cached, cacheHit: true}, nil
	}

	started := e.now()
	call, callErr := e.call(ctx, payload)
	if callErr != nil {
		switch {
		case errors.Is(callErr, context.Canceled):
			call.reason = ReasonContextCanceled
		case errors.Is(callErr, context.DeadlineExceeded):
			call.reason = ReasonTimeout
		default:
			call.reason = ReasonHTTPStatus
		}
	}
	outcome := batchOutcome{key: key, answers: make(map[string]Answer, len(batch))}
	if call.attempts == 0 {
		// Nothing reached the network: no receipt row.
		for _, pq := range batch {
			outcome.answers[pq.q.ID] = undecided(pq.q.Type, call.reason)
		}
		return outcome, callErr
	}

	bankID, bankVersion, bankHash := bankFields(r.Bank)
	row := Receipt{
		Key:         key,
		Time:        e.now().UTC().Format(time.RFC3339Nano),
		Purpose:     string(r.Purpose),
		Subject:     r.Subject,
		Bank:        bankID,
		BankVersion: bankVersion,
		BankHash:    bankHash,
		Contract:    r.Topic.Contract,
		Model:       e.model,
		Route:       RouteDecisions,
		Attempts:    call.attempts,
		LatencyMS:   e.now().Sub(started).Milliseconds(),
		Questions:   ids,
		Status:      StatusUndecided,
		Reason:      call.reason,
	}
	if call.response != nil {
		response := call.response
		row.ReportedModel = response.Model
		row.ResponseID = response.ID
		row.InputTokens = response.inputTokens()
		row.Cost = response.cost()
		row.CostUnknown = row.Cost == nil
		e.charge(row.Cost, row.InputTokens)
		e.stats.call(row.Cost)
		if row.Cost != nil {
			outcome.cost = *row.Cost
		}
		// Decode answers after usage: rejected output can still be billed.
		if err := json.Unmarshal(response.Answers, &row.Answers); err != nil || row.Answers == nil {
			row.Reason = ReasonInvalidReceipt
			for _, pq := range batch {
				outcome.answers[pq.q.ID] = undecided(pq.q.Type, ReasonInvalidReceipt)
			}
		} else {
			answers, allDecided := validateAnswers(batch, row.Answers)
			outcome.answers = answers
			if allDecided {
				row.Status, row.Reason = StatusDecided, ""
			} else {
				row.Reason = ReasonInvalidReceipt
			}
		}
	} else {
		row.CostUnknown = true
		e.stats.call(nil)
		for _, pq := range batch {
			outcome.answers[pq.q.ID] = undecided(pq.q.Type, call.reason)
		}
	}
	if err := e.receipts.Append(r.Topic.Root, row); err != nil {
		e.logger.Warn("decisions: receipt not written", "topic", r.Topic.Slug, "key", key, "error", err)
	}
	return outcome, callErr
}

// fromCache returns re-validated answers of a decided receipt for key.
func (e *Engine) fromCache(topicRoot, key string, batch []preparedQuestion) (map[string]Answer, bool) {
	row, ok, err := e.receipts.Lookup(topicRoot, key)
	if err != nil {
		e.logger.Warn("decisions: receipts unreadable, cache disabled for this lookup", "root", topicRoot, "error", err)
		return nil, false
	}
	if !ok {
		return nil, false
	}
	answers, allDecided := validateAnswers(batch, row.Answers)
	if !allDecided {
		return nil, false
	}
	return answers, true
}

func (e *Engine) charge(cost *float64, inputTokens *int64) {
	switch {
	case cost != nil:
		e.budget.Charge(*cost)
	case inputTokens != nil:
		e.budget.Charge(float64(*inputTokens) * listPricePerInputToken)
	}
}

func (e *Engine) fatalErr() error {
	e.fatalMu.Lock()
	defer e.fatalMu.Unlock()
	return e.fatal
}

func (e *Engine) setFatal(err error) error {
	e.fatalMu.Lock()
	defer e.fatalMu.Unlock()
	if e.fatal == nil {
		e.fatal = err
	}
	return e.fatal
}

// apiResponse is the decisions endpoint's 200 body.
type apiResponse struct {
	ID      string          `json:"id"`
	Model   string          `json:"model"`
	Answers json.RawMessage `json:"answers"`
	Usage   *struct {
		InputTokens json.RawMessage `json:"input_tokens"`
		Cost        json.RawMessage `json:"cost"`
	} `json:"usage"`
}

func (r *apiResponse) inputTokens() *int64 {
	if r.Usage == nil {
		return nil
	}
	var tokens *int64
	if err := json.Unmarshal(r.Usage.InputTokens, &tokens); err != nil || tokens == nil || *tokens < 0 {
		return nil
	}
	return tokens
}

func (r *apiResponse) cost() *float64 {
	if r.Usage == nil {
		return nil
	}
	return ParseCost(r.Usage.Cost)
}

// ParseCost reads a reported cost (a JSON number or numeric string). A
// missing, null, negative or non-numeric cost is unknown (nil), never zero.
func ParseCost(raw json.RawMessage) *float64 {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil
	}
	if raw[0] == '"' {
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			return nil
		}
		raw = []byte(strings.TrimSpace(text))
	}
	value, err := strconv.ParseFloat(string(raw), 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return nil
	}
	return &value
}

// callResult is the transport outcome of one decision request.
type callResult struct {
	response *apiResponse
	attempts int
	reason   string // set when no usable response
}

// call posts payload with retries. It returns an error only when the run
// must stop (ErrAuth, ctx cancellation).
func (e *Engine) call(ctx context.Context, payload []byte) (callResult, error) {
	result := callResult{}
	lastReason := ReasonRetries
	pausedBy429 := false
	for attempt := 0; attempt <= e.retries; attempt++ {
		if attempt > 0 && !pausedBy429 {
			if err := e.sleep(ctx, backoff(attempt)); err != nil {
				return result, err
			}
		}
		if err := e.waitPause(ctx); err != nil {
			return result, err
		}
		if err := e.fatalErr(); err != nil {
			return result, err
		}
		pausedBy429 = false
		if e.budget.Exceeded() {
			result.reason = ReasonBudget
			return result, nil
		}

		status, header, body, err := e.post(ctx, payload)
		result.attempts++
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return result, ctxErr
			}
			lastReason = ReasonRetries
			if errors.Is(err, context.DeadlineExceeded) {
				lastReason = ReasonTimeout
			}
			e.logger.Debug("decisions: attempt failed", "attempt", result.attempts, "error", err)
			continue
		}
		switch {
		case status == http.StatusOK:
			var response apiResponse
			if err := json.Unmarshal(body, &response); err != nil {
				result.reason = ReasonInvalidReceipt
				return result, nil
			}
			result.response = &response
			return result, nil
		case status == http.StatusUnauthorized || status == http.StatusPaymentRequired || status == http.StatusForbidden:
			return result, e.setFatal(fmt.Errorf("%w: HTTP %d: %s", ErrAuth, status, errorMessage(body)))
		case retryableStatus(status):
			lastReason = ReasonRetries
			if status == http.StatusTooManyRequests {
				wait, ok := retryAfter(header, e.now())
				if !ok {
					wait = backoff(attempt + 1)
				}
				e.pause(wait)
				pausedBy429 = true
			}
			e.logger.Debug("decisions: retryable status", "status", status, "attempt", result.attempts)
		default:
			e.logger.Warn("decisions: request rejected", "status", status, "message", errorMessage(body))
			result.reason = ReasonHTTPStatus
			return result, nil
		}
	}
	result.reason = lastReason
	return result, nil
}

// post runs one attempt under the per-attempt deadline, body read included.
func (e *Engine) post(ctx context.Context, payload []byte) (int, http.Header, []byte, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, e.deadline)
	defer cancel()
	request, err := http.NewRequestWithContext(attemptCtx, http.MethodPost, e.endpoint, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+e.apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	response, err := e.client.Do(request)
	if err != nil {
		return 0, nil, nil, attemptError(attemptCtx, err)
	}
	defer func() { _ = response.Body.Close() }()
	// OpenRouter may stream keep-alive whitespace before the JSON; reading
	// the whole body under the attempt deadline and letting the JSON decoder
	// skip leading whitespace handles it.
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

func (e *Engine) pause(wait time.Duration) {
	e.pauseMu.Lock()
	defer e.pauseMu.Unlock()
	until := e.now().Add(wait)
	if until.After(e.pauseUntil) {
		e.pauseUntil = until
	}
}

// waitPause blocks while a 429 has paused the whole pool. Another worker
// may extend the pause while this one sleeps, so the shared deadline is
// re-read after every wake until it has passed.
func (e *Engine) waitPause(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		e.pauseMu.Lock()
		wait := e.pauseUntil.Sub(e.now())
		e.pauseMu.Unlock()
		if wait <= 0 {
			return nil
		}
		if err := e.sleep(ctx, wait); err != nil {
			return err
		}
	}
}

func retryableStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout, http.StatusTooManyRequests, 524, 529:
		return true
	default:
		return status >= 500 && status <= 599
	}
}

// backoff is the jittered exponential delay before retry n (1-based):
// 0.5 s doubling to a 5 s cap, ±25 %.
func backoff(retry int) time.Duration {
	delay := backoffBase << max(retry-1, 0)
	if delay > backoffCap || delay <= 0 {
		delay = backoffCap
	}
	jitter := 0.75 + rand.Float64()*0.5
	return time.Duration(float64(delay) * jitter)
}

// retryAfter reads retry-after-ms or Retry-After (seconds or HTTP date),
// capped at 30 s.
func retryAfter(header http.Header, now time.Time) (time.Duration, bool) {
	if header == nil {
		return 0, false
	}
	var wait time.Duration
	if ms, err := strconv.ParseFloat(strings.TrimSpace(header.Get("retry-after-ms")), 64); err == nil && ms >= 0 {
		wait = time.Duration(min(ms, float64(maxRetryAfter/time.Millisecond)) * float64(time.Millisecond))
	} else if value := strings.TrimSpace(header.Get("Retry-After")); value != "" {
		if seconds, err := strconv.ParseFloat(value, 64); err == nil && seconds >= 0 {
			wait = time.Duration(min(seconds, maxRetryAfter.Seconds()) * float64(time.Second))
		} else if date, err := http.ParseTime(value); err == nil {
			wait = max(date.Sub(now), 0)
		} else {
			return 0, false
		}
	} else {
		return 0, false
	}
	return min(wait, maxRetryAfter), true
}

// errorMessage extracts OpenRouter's {"error":{"message"}} or a short body.
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
