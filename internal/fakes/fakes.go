// Package fakes provides httptest servers standing in for the three network
// boundaries of decision-backed kb flows: the OpenRouter decisions endpoint
// (Jev), OpenRouter chat completions (generation) and Firecrawl scrape. Tests
// use them to exercise real vault flows without the network (spec §17).
package fakes

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// Question is one decoded question from a decisions request.
type Question struct {
	ID           string
	Type         string
	Instructions any
	Criteria     any
}

// Options returns the option names of a choice question in sorted order.
func (q Question) Options() []string {
	criteria, ok := q.Criteria.(map[string]any)
	if !ok {
		return nil
	}
	return slices.Sorted(maps.Keys(criteria))
}

// Levels returns the number of score levels.
func (q Question) Levels() int {
	levels, ok := q.Criteria.([]any)
	if !ok {
		return 0
	}
	return len(levels)
}

// Call is one decoded decisions request.
type Call struct {
	Model     string
	State     json.RawMessage
	Questions map[string]Question
}

// StateString returns the raw JSON state as text (for substring rules).
func (c Call) StateString() string { return string(c.State) }

// DecideFunc answers one question of a call. Returning nil falls back to the
// server default (noul 0.05; choice picks the exit option `none`, `unknown`,
// `other` or `different` when present, else the first option; score 0).
type DecideFunc func(call Call, q Question) any

// GenerateFunc answers one generation request: schema name and the user
// prompt; it returns the JSON object to place in the message content.
type GenerateFunc func(schemaName, system, prompt string) any

// OpenRouter is a fake OpenRouter serving /alpha/decisions and
// /v1/chat/completions. Point [openrouter].api_url at URL.
type OpenRouter struct {
	*httptest.Server

	mu          sync.Mutex
	decide      DecideFunc
	generate    GenerateFunc
	calls       []Call
	genCalls    []GenRequest
	CostPerCall float64
}

// GenRequest is one recorded generation request.
type GenRequest struct {
	Model      string
	SchemaName string
	System     string
	Prompt     string
}

// NewOpenRouter starts a fake OpenRouter. Close it with t.Cleanup(srv.Close).
func NewOpenRouter(decide DecideFunc, generate GenerateFunc) *OpenRouter {
	fake := &OpenRouter{decide: decide, generate: generate, CostPerCall: 0.0001}
	fake.Server = httptest.NewServer(http.HandlerFunc(fake.serve))
	return fake
}

// SetDecide replaces the decide rule.
func (f *OpenRouter) SetDecide(decide DecideFunc) {
	f.mu.Lock()
	f.decide = decide
	f.mu.Unlock()
}

// SetGenerate replaces the generate rule.
func (f *OpenRouter) SetGenerate(generate GenerateFunc) {
	f.mu.Lock()
	f.generate = generate
	f.mu.Unlock()
}

// Calls returns a copy of every decisions request received.
func (f *OpenRouter) Calls() []Call {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.calls)
}

// GenCalls returns a copy of every generation request received.
func (f *OpenRouter) GenCalls() []GenRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.genCalls)
}

// QuestionCount counts received questions whose id starts with prefix.
func (f *OpenRouter) QuestionCount(prefix string) int {
	count := 0
	for _, call := range f.Calls() {
		for id := range call.Questions {
			if strings.HasPrefix(id, prefix) {
				count++
			}
		}
	}
	return count
}

func (f *OpenRouter) serve(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch {
	case strings.HasSuffix(r.URL.Path, "/alpha/decisions"):
		f.serveDecisions(w, body)
	case strings.HasSuffix(r.URL.Path, "/chat/completions"):
		f.serveGeneration(w, body)
	default:
		http.NotFound(w, r)
	}
}

func (f *OpenRouter) serveDecisions(w http.ResponseWriter, body []byte) {
	var payload struct {
		Model     string                     `json:"model"`
		State     json.RawMessage            `json:"state"`
		Questions map[string]json.RawMessage `json:"questions"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	call := Call{Model: payload.Model, State: payload.State, Questions: map[string]Question{}}
	for id, raw := range payload.Questions {
		var q struct {
			Type         string `json:"type"`
			Instructions any    `json:"instructions"`
			Criteria     any    `json:"criteria"`
		}
		if err := json.Unmarshal(raw, &q); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		call.Questions[id] = Question{ID: id, Type: q.Type, Instructions: q.Instructions, Criteria: q.Criteria}
	}

	f.mu.Lock()
	f.calls = append(f.calls, call)
	decide := f.decide
	cost := f.CostPerCall
	f.mu.Unlock()

	answers := map[string]any{}
	for id, q := range call.Questions {
		var answer any
		if decide != nil {
			answer = decide(call, q)
		}
		if answer == nil {
			answer = defaultAnswer(q)
		}
		answers[id] = normalizeAnswer(q, answer)
	}
	writeJSON(w, map[string]any{
		"id":       fmt.Sprintf("gen-dec-fake-%d", len(f.Calls())),
		"model":    "typesafe/jev-1.13-20260917",
		"provider": "TypeSafe",
		"answers":  answers,
		"usage":    map[string]any{"input_tokens": len(body) / 4, "output_tokens": 0, "cost": cost},
	})
}

func (f *OpenRouter) serveGeneration(w http.ResponseWriter, body []byte) {
	var payload struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		ResponseFormat struct {
			JSONSchema struct {
				Name string `json:"name"`
			} `json:"json_schema"`
		} `json:"response_format"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := GenRequest{Model: payload.Model, SchemaName: payload.ResponseFormat.JSONSchema.Name}
	for _, message := range payload.Messages {
		switch message.Role {
		case "system":
			request.System += message.Content
		case "user":
			request.Prompt += message.Content
		}
	}
	f.mu.Lock()
	f.genCalls = append(f.genCalls, request)
	generate := f.generate
	cost := f.CostPerCall
	f.mu.Unlock()

	var output any = map[string]any{}
	if generate != nil {
		output = generate(request.SchemaName, request.System, request.Prompt)
	}
	content, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"id":    fmt.Sprintf("gen-fake-%d", len(f.GenCalls())),
		"model": payload.Model,
		"choices": []any{map[string]any{
			"message": map[string]any{"role": "assistant", "content": string(content)},
		}},
		"usage": map[string]any{
			"prompt_tokens":             len(body) / 4,
			"completion_tokens":         len(content) / 4,
			"cost":                      cost,
			"completion_tokens_details": map[string]any{"reasoning_tokens": 0},
		},
	})
}

// Noul builds a noul answer.
func Noul(p float64) any { return map[string]any{"type": "noul", "noul": p} }

// Pick builds a choice answer giving option p and spreading the rest evenly
// over the question's other options (the fake fills the distribution).
func Pick(option string, p float64) any {
	return pickAnswer{option: option, p: p}
}

// Dist builds a choice answer from an explicit distribution; missing options
// get 0 and the argmax becomes the choice.
func Dist(probs map[string]float64) any { return distAnswer(probs) }

// Level builds a score answer concentrated on one level.
func Level(level int) any { return levelAnswer(level) }

type pickAnswer struct {
	option string
	p      float64
}

type distAnswer map[string]float64

type levelAnswer int

func defaultAnswer(q Question) any {
	switch q.Type {
	case "noul":
		return Noul(0.05)
	case "choice":
		options := q.Options()
		for _, exit := range []string{"none", "unknown", "other", "different"} {
			if slices.Contains(options, exit) {
				return Pick(exit, 0.9)
			}
		}
		if len(options) > 0 {
			return Pick(options[0], 0.9)
		}
	case "score":
		return Level(0)
	}
	return Noul(0.05)
}

func normalizeAnswer(q Question, answer any) any {
	switch typed := answer.(type) {
	case pickAnswer:
		options := q.Options()
		probs := map[string]float64{}
		rest := 0.0
		if len(options) > 1 {
			rest = (1 - typed.p) / float64(len(options)-1)
		}
		for _, option := range options {
			probs[option] = rest
		}
		probs[typed.option] = typed.p
		return choiceAnswer(probs)
	case distAnswer:
		probs := map[string]float64{}
		for _, option := range q.Options() {
			probs[option] = typed[option]
		}
		return choiceAnswer(probs)
	case levelAnswer:
		levels := max(q.Levels(), 2)
		probs := map[string]float64{}
		legend := map[string]any{}
		for level := range levels {
			key := strconv.Itoa(level)
			probs[key] = 0
			legend[key] = "level " + key
		}
		probs[strconv.Itoa(int(typed))] = 1
		return map[string]any{"type": "score", "score": float64(typed), "probabilities": probs, "legend": legend, "confidence": 1.0}
	default:
		return answer
	}
}

func choiceAnswer(probs map[string]float64) any {
	best, bestP := "", -1.0
	for _, option := range slices.Sorted(maps.Keys(probs)) {
		if probs[option] > bestP {
			best, bestP = option, probs[option]
		}
	}
	n := float64(len(probs))
	confidence := 1.0
	if n > 1 {
		confidence = max(0, min(1, (n*bestP-1)/(n-1)))
	}
	return map[string]any{"type": "choice", "choice": best, "probabilities": probs, "confidence": confidence}
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

// ScrapeRequest is one decoded Firecrawl /v2/scrape request.
type ScrapeRequest struct {
	URL             string   `json:"url"`
	Formats         []string `json:"formats"`
	MaxAge          *int64   `json:"maxAge,omitempty"`
	WaitFor         *int64   `json:"waitFor,omitempty"`
	OnlyMainContent *bool    `json:"onlyMainContent,omitempty"`
}

// ScrapeResponse is what the fake returns for one scrape.
type ScrapeResponse struct {
	Markdown   string
	Title      string
	SourceURL  string
	StatusCode int
}

// ScrapeFunc answers one scrape request.
type ScrapeFunc func(req ScrapeRequest) ScrapeResponse

// Firecrawl is a fake Firecrawl server. Point [firecrawl].api_url at URL.
type Firecrawl struct {
	*httptest.Server

	mu       sync.Mutex
	scrape   ScrapeFunc
	requests []ScrapeRequest
}

// NewFirecrawl starts a fake Firecrawl server.
func NewFirecrawl(scrape ScrapeFunc) *Firecrawl {
	fake := &Firecrawl{scrape: scrape}
	fake.Server = httptest.NewServer(http.HandlerFunc(fake.serve))
	return fake
}

// Requests returns a copy of every scrape request received.
func (f *Firecrawl) Requests() []ScrapeRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.requests)
}

func (f *Firecrawl) serve(w http.ResponseWriter, r *http.Request) {
	var request ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.mu.Lock()
	f.requests = append(f.requests, request)
	scrape := f.scrape
	f.mu.Unlock()

	response := ScrapeResponse{Markdown: "# Page\n\nBody.", Title: "Page", SourceURL: request.URL, StatusCode: 200}
	if scrape != nil {
		response = scrape(request)
	}
	if response.SourceURL == "" {
		response.SourceURL = request.URL
	}
	if response.StatusCode == 0 {
		response.StatusCode = 200
	}
	writeJSON(w, map[string]any{
		"success": true,
		"data": map[string]any{
			"markdown": response.Markdown,
			"metadata": map[string]any{
				"title":      response.Title,
				"sourceURL":  request.URL,
				"url":        response.SourceURL,
				"statusCode": response.StatusCode,
			},
		},
	})
}
