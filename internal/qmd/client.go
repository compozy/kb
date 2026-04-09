// Package qmd provides a shell-backed client for the QMD CLI.
package qmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	// DefaultBinaryPath is the default executable name used to locate QMD.
	DefaultBinaryPath = "qmd"
	// InstallCommand is the recommended install command for the QMD CLI.
	InstallCommand = "npm install -g @tobilu/qmd"
)

var (
	// ErrQMDUnavailable indicates the QMD CLI could not be found on PATH.
	ErrQMDUnavailable = errors.New("qmd is not available")

	updateCountsPattern     = regexp.MustCompile(`Indexed:\s+(\d+)\s+new,\s+(\d+)\s+updated,\s+(\d+)\s+unchanged,\s+(\d+)\s+removed`)
	collectionsPattern      = regexp.MustCompile(`Updating\s+(\d+)\s+collection\(s\)`)
	embeddingHintPattern    = regexp.MustCompile(`\((\d+)\s+unique hashes need vectors\)`)
	totalPattern            = regexp.MustCompile(`Total:\s+(\d+)`)
	vectorsPattern          = regexp.MustCompile(`Vectors:\s+(\d+)`)
	pendingPattern          = regexp.MustCompile(`Pending:\s+(\d+)`)
	collectionHeaderPattern = regexp.MustCompile(`^(.+?)\s+\((qmd://.+)\)$`)
	collectionFilesPattern  = regexp.MustCompile(`Files:\s+(\d+)\s+\(updated\s+(.+)\)`)
	embedDonePattern        = regexp.MustCompile(`Embedded\s+(\d+)\s+chunks\s+from\s+(\d+)\s+documents\s+in\s+([^\n]+)`)
	embedErrorPattern       = regexp.MustCompile(`(\d+)\s+chunks failed`)
	durationTokenPattern    = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(ms|s|m|h)`)
	ansiPattern             = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)
)

type lookPathFunc func(string) (string, error)
type commandContextFunc func(context.Context, string, ...string) *exec.Cmd

// SearchMode selects the QMD retrieval mode for Search.
type SearchMode string

const (
	// SearchModeHybrid runs QMD's hybrid `query` command.
	SearchModeHybrid SearchMode = "hybrid"
	// SearchModeLexical runs QMD's lexical `search` command.
	SearchModeLexical SearchMode = "lexical"
	// SearchModeVector runs QMD's semantic `vsearch` command.
	SearchModeVector SearchMode = "vector"
)

// IndexOperation selects whether Index creates a collection or updates one.
type IndexOperation string

const (
	// IndexOperationAdd creates a collection and performs the initial sync.
	IndexOperationAdd IndexOperation = "add"
	// IndexOperationUpdate refreshes an existing collection.
	IndexOperationUpdate IndexOperation = "update"
)

// SearchOptions configures a QMD search invocation.
type SearchOptions struct {
	Query      string
	Mode       SearchMode
	Limit      int
	All        bool
	MinScore   *float64
	Full       bool
	Collection string
}

// SearchResult is the normalized Go representation of a QMD search hit.
type SearchResult struct {
	DocID   string  `json:"docid,omitempty"`
	Path    string  `json:"path"`
	Title   string  `json:"title"`
	Snippet string  `json:"snippet"`
	Score   float64 `json:"score"`
}

// IndexOptions configures a collection sync invocation.
type IndexOptions struct {
	Operation      IndexOperation
	VaultPath      string
	CollectionName string
	Context        string
	Embed          bool
	ForceEmbed     bool
}

// UpdateResult mirrors QMD's store update summary.
type UpdateResult struct {
	Collections    int `json:"collections"`
	Indexed        int `json:"indexed"`
	Updated        int `json:"updated"`
	Unchanged      int `json:"unchanged"`
	Removed        int `json:"removed"`
	NeedsEmbedding int `json:"needsEmbedding"`
}

// EmbedResult mirrors QMD's embedding summary.
type EmbedResult struct {
	DocsProcessed  int `json:"docsProcessed"`
	ChunksEmbedded int `json:"chunksEmbedded"`
	Errors         int `json:"errors"`
	DurationMs     int `json:"durationMs"`
}

// CollectionInfo summarizes one indexed QMD collection.
type CollectionInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Pattern     string `json:"pattern"`
	Documents   int    `json:"documents"`
	LastUpdated string `json:"lastUpdated"`
}

// IndexStatus mirrors QMD's status view for collection health.
type IndexStatus struct {
	TotalDocuments int              `json:"totalDocuments"`
	NeedsEmbedding int              `json:"needsEmbedding"`
	HasVectorIndex bool             `json:"hasVectorIndex"`
	Collections    []CollectionInfo `json:"collections"`
}

// IndexResult groups the collection sync summary, optional embedding summary,
// and the final index status snapshot.
type IndexResult struct {
	CollectionName string       `json:"collectionName"`
	UpdateResult   UpdateResult `json:"updateResult"`
	EmbedResult    *EmbedResult `json:"embedResult,omitempty"`
	Status         IndexStatus  `json:"status"`
}

// ClientOption mutates a QMDClient configuration.
type ClientOption func(*QMDClient)

// QMDClient executes QMD shell commands with context-aware process management.
type QMDClient struct {
	binaryPath string
	indexName  string

	lookPath       lookPathFunc
	commandContext commandContextFunc
}

// NewClient constructs a QMDClient with sensible defaults.
func NewClient(options ...ClientOption) *QMDClient {
	client := &QMDClient{
		binaryPath:     DefaultBinaryPath,
		lookPath:       exec.LookPath,
		commandContext: exec.CommandContext,
	}
	for _, option := range options {
		if option != nil {
			option(client)
		}
	}
	return client
}

// WithBinaryPath overrides the executable used for QMD invocations.
func WithBinaryPath(path string) ClientOption {
	return func(client *QMDClient) {
		client.binaryPath = strings.TrimSpace(path)
	}
}

// WithIndexName routes QMD commands to a named SQLite index.
func WithIndexName(name string) ClientOption {
	return func(client *QMDClient) {
		client.indexName = strings.TrimSpace(name)
	}
}

// Status returns the current QMD index status summary.
func (client *QMDClient) Status(ctx context.Context) (IndexStatus, error) {
	stdout, _, err := client.run(ctx, client.statusCommand())
	if err != nil {
		return IndexStatus{}, err
	}

	status, err := parseIndexStatus(stdout)
	if err != nil {
		return IndexStatus{}, fmt.Errorf("qmd status: parse output: %w", err)
	}

	return status, nil
}

// Search executes a QMD search command and returns normalized results.
func (client *QMDClient) Search(ctx context.Context, options SearchOptions) ([]SearchResult, error) {
	query := strings.TrimSpace(options.Query)
	if query == "" {
		return nil, fmt.Errorf("qmd search: query is required")
	}

	command, err := client.searchCommand(options)
	if err != nil {
		return nil, err
	}

	stdout, _, err := client.run(ctx, command)
	if err != nil {
		return nil, err
	}

	var rawResults []searchResultPayload
	if err := json.Unmarshal([]byte(stdout), &rawResults); err != nil {
		return nil, fmt.Errorf("qmd search: parse JSON output: %w", err)
	}

	results := make([]SearchResult, 0, len(rawResults))
	for _, rawResult := range rawResults {
		normalized := rawResult.normalize(options.Full)
		if options.MinScore != nil && normalized.Score < *options.MinScore {
			continue
		}
		results = append(results, normalized)
	}

	return results, nil
}

// Index creates or updates a QMD collection and returns the structured summary.
func (client *QMDClient) Index(ctx context.Context, options IndexOptions) (IndexResult, error) {
	operation, err := normalizeIndexOperation(options.Operation)
	if err != nil {
		return IndexResult{}, err
	}

	collectionName := strings.TrimSpace(options.CollectionName)
	if collectionName == "" {
		return IndexResult{}, fmt.Errorf("qmd index: collection name is required")
	}
	if options.ForceEmbed && !options.Embed {
		return IndexResult{}, fmt.Errorf("qmd index: force embed requires embed=true")
	}

	syncCommand, err := client.indexCommand(operation, options)
	if err != nil {
		return IndexResult{}, err
	}

	stdout, _, err := client.run(ctx, syncCommand)
	if err != nil {
		return IndexResult{}, err
	}

	updateResult, err := parseUpdateResult(stdout, operation)
	if err != nil {
		return IndexResult{}, fmt.Errorf("qmd index: parse update output: %w", err)
	}

	if contextText := strings.TrimSpace(options.Context); contextText != "" {
		if _, _, err := client.run(ctx, commandSpec{
			label: "context add",
			args: client.baseArgs(
				"context",
				"add",
				fmt.Sprintf("qmd://%s/", collectionName),
				contextText,
			),
		}); err != nil {
			return IndexResult{}, err
		}
	}

	var embedResult *EmbedResult
	if options.Embed {
		embedOutput, _, err := client.run(ctx, client.embedCommand(options.ForceEmbed))
		if err != nil {
			return IndexResult{}, err
		}

		parsedEmbedResult, err := parseEmbedResult(embedOutput)
		if err != nil {
			return IndexResult{}, fmt.Errorf("qmd index: parse embed output: %w", err)
		}
		embedResult = &parsedEmbedResult
	}

	statusOutput, _, err := client.run(ctx, client.statusCommand())
	if err != nil {
		return IndexResult{}, err
	}

	status, err := parseIndexStatus(statusOutput)
	if err != nil {
		return IndexResult{}, fmt.Errorf("qmd index: parse status output: %w", err)
	}

	return IndexResult{
		CollectionName: collectionName,
		UpdateResult:   updateResult,
		EmbedResult:    embedResult,
		Status:         status,
	}, nil
}

func (client *QMDClient) searchCommand(options SearchOptions) (commandSpec, error) {
	mode, commandName, err := normalizeSearchMode(options.Mode)
	if err != nil {
		return commandSpec{}, err
	}

	args := client.baseArgs(commandName, "--json")
	if options.All {
		args = append(args, "--all")
	} else if options.Limit > 0 {
		args = append(args, "-n", strconv.Itoa(options.Limit))
	}
	if options.MinScore != nil {
		args = append(args, "--min-score", strconv.FormatFloat(*options.MinScore, 'f', -1, 64))
	}
	if options.Full {
		args = append(args, "--full")
	}
	if collection := strings.TrimSpace(options.Collection); collection != "" {
		args = append(args, "-c", collection)
	}
	args = append(args, strings.TrimSpace(options.Query))

	return commandSpec{
		label: fmt.Sprintf("search (%s)", mode),
		args:  args,
	}, nil
}

func (client *QMDClient) indexCommand(operation IndexOperation, options IndexOptions) (commandSpec, error) {
	switch operation {
	case IndexOperationAdd:
		vaultPath := strings.TrimSpace(options.VaultPath)
		if vaultPath == "" {
			return commandSpec{}, fmt.Errorf("qmd index: vault path is required for add")
		}
		return commandSpec{
			label: "collection add",
			args: client.baseArgs(
				"collection",
				"add",
				vaultPath,
				"--name",
				strings.TrimSpace(options.CollectionName),
			),
		}, nil
	case IndexOperationUpdate:
		return commandSpec{
			label: "update",
			args:  client.baseArgs("update", strings.TrimSpace(options.CollectionName)),
		}, nil
	default:
		return commandSpec{}, fmt.Errorf("qmd index: unsupported operation %q", operation)
	}
}

func (client *QMDClient) embedCommand(force bool) commandSpec {
	args := client.baseArgs("embed")
	if force {
		args = append(args, "-f")
	}

	return commandSpec{
		label: "embed",
		args:  args,
	}
}

func (client *QMDClient) statusCommand() commandSpec {
	return commandSpec{
		label: "status",
		args:  client.baseArgs("status"),
	}
}

func (client *QMDClient) baseArgs(args ...string) []string {
	result := make([]string, 0, len(args)+2)
	if indexName := strings.TrimSpace(client.indexName); indexName != "" {
		result = append(result, "--index", indexName)
	}
	result = append(result, args...)
	return result
}

func (client *QMDClient) resolveBinary() (string, error) {
	binaryPath := strings.TrimSpace(client.binaryPath)
	if binaryPath == "" {
		binaryPath = DefaultBinaryPath
	}

	resolvedPath, err := client.lookPath(binaryPath)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%w: install it with `%s`", ErrQMDUnavailable, InstallCommand)
		}
		return "", fmt.Errorf("resolve qmd binary %q: %w", binaryPath, err)
	}

	return resolvedPath, nil
}

func (client *QMDClient) run(ctx context.Context, spec commandSpec) (string, string, error) {
	binaryPath, err := client.resolveBinary()
	if err != nil {
		return "", "", err
	}

	command := client.commandContext(ctx, binaryPath, spec.args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err = command.Run()
	stdoutText := strings.TrimSpace(stdout.String())
	stderrText := strings.TrimSpace(stderr.String())
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return stdoutText, stderrText, fmt.Errorf("qmd %s: %w", spec.label, ctxErr)
		}

		diagnostics := cleanDiagnostics(stderrText)
		if diagnostics == "" {
			diagnostics = cleanDiagnostics(stdoutText)
		}
		if diagnostics != "" {
			return stdoutText, stderrText, fmt.Errorf("qmd %s: %w: %s", spec.label, err, diagnostics)
		}

		return stdoutText, stderrText, fmt.Errorf("qmd %s: %w", spec.label, err)
	}

	return stdoutText, stderrText, nil
}

type commandSpec struct {
	label string
	args  []string
}

type searchResultPayload struct {
	DocID       string  `json:"docid"`
	Score       float64 `json:"score"`
	File        string  `json:"file"`
	FilePath    string  `json:"filepath"`
	DisplayPath string  `json:"displayPath"`
	Title       string  `json:"title"`
	Snippet     string  `json:"snippet"`
	Body        string  `json:"body"`
	BestChunk   string  `json:"bestChunk"`
}

func (payload searchResultPayload) normalize(full bool) SearchResult {
	return SearchResult{
		DocID:   payload.DocID,
		Path:    firstNonEmpty(payload.DisplayPath, payload.FilePath, payload.File),
		Title:   payload.Title,
		Snippet: payload.resolveSnippet(full),
		Score:   payload.Score,
	}
}

func (payload searchResultPayload) resolveSnippet(full bool) string {
	if full {
		return firstNonEmpty(payload.Body, payload.BestChunk, payload.Snippet)
	}
	return firstNonEmpty(payload.BestChunk, payload.Snippet, payload.Body)
}

func normalizeSearchMode(mode SearchMode) (SearchMode, string, error) {
	switch normalized := SearchMode(strings.ToLower(strings.TrimSpace(string(mode)))); normalized {
	case "", SearchModeHybrid:
		return SearchModeHybrid, "query", nil
	case SearchModeLexical, "lex":
		return SearchModeLexical, "search", nil
	case SearchModeVector, "vec":
		return SearchModeVector, "vsearch", nil
	default:
		return "", "", fmt.Errorf("qmd search: unsupported mode %q", mode)
	}
}

func normalizeIndexOperation(operation IndexOperation) (IndexOperation, error) {
	switch normalized := IndexOperation(strings.ToLower(strings.TrimSpace(string(operation)))); normalized {
	case "":
		return IndexOperationAdd, nil
	case IndexOperationAdd:
		return IndexOperationAdd, nil
	case IndexOperationUpdate:
		return normalized, nil
	default:
		return "", fmt.Errorf("qmd index: unsupported operation %q", operation)
	}
}

func parseUpdateResult(output string, operation IndexOperation) (UpdateResult, error) {
	text := cleanOutput(output)
	matches := updateCountsPattern.FindStringSubmatch(text)
	if matches == nil {
		return UpdateResult{}, fmt.Errorf("missing indexed summary in %q", text)
	}

	indexed, err := strconv.Atoi(matches[1])
	if err != nil {
		return UpdateResult{}, err
	}
	updated, err := strconv.Atoi(matches[2])
	if err != nil {
		return UpdateResult{}, err
	}
	unchanged, err := strconv.Atoi(matches[3])
	if err != nil {
		return UpdateResult{}, err
	}
	removed, err := strconv.Atoi(matches[4])
	if err != nil {
		return UpdateResult{}, err
	}

	collections := 1
	if collectionMatches := collectionsPattern.FindStringSubmatch(text); collectionMatches != nil {
		parsedCollections, err := strconv.Atoi(collectionMatches[1])
		if err != nil {
			return UpdateResult{}, err
		}
		collections = parsedCollections
	} else if operation == IndexOperationUpdate {
		collections = 1
	}

	needsEmbedding := 0
	if embeddingMatches := embeddingHintPattern.FindStringSubmatch(text); embeddingMatches != nil {
		needsEmbedding, err = strconv.Atoi(embeddingMatches[1])
		if err != nil {
			return UpdateResult{}, err
		}
	}

	return UpdateResult{
		Collections:    collections,
		Indexed:        indexed,
		Updated:        updated,
		Unchanged:      unchanged,
		Removed:        removed,
		NeedsEmbedding: needsEmbedding,
	}, nil
}

func parseEmbedResult(output string) (EmbedResult, error) {
	text := cleanOutput(output)
	if text == "" {
		return EmbedResult{}, fmt.Errorf("missing embed output")
	}
	if strings.Contains(text, "All content hashes already have embeddings") || strings.Contains(text, "No non-empty documents to embed") {
		return EmbedResult{}, nil
	}

	matches := embedDonePattern.FindStringSubmatch(text)
	if matches == nil {
		return EmbedResult{}, fmt.Errorf("missing embed summary in %q", text)
	}

	chunksEmbedded, err := strconv.Atoi(matches[1])
	if err != nil {
		return EmbedResult{}, err
	}
	docsProcessed, err := strconv.Atoi(matches[2])
	if err != nil {
		return EmbedResult{}, err
	}
	durationMs, err := parseHumanDurationMilliseconds(matches[3])
	if err != nil {
		return EmbedResult{}, err
	}

	errorsCount := 0
	if errorMatches := embedErrorPattern.FindStringSubmatch(text); errorMatches != nil {
		errorsCount, err = strconv.Atoi(errorMatches[1])
		if err != nil {
			return EmbedResult{}, err
		}
	}

	return EmbedResult{
		DocsProcessed:  docsProcessed,
		ChunksEmbedded: chunksEmbedded,
		Errors:         errorsCount,
		DurationMs:     durationMs,
	}, nil
}

func parseIndexStatus(output string) (IndexStatus, error) {
	text := cleanOutput(output)
	lines := strings.Split(text, "\n")

	status := IndexStatus{}
	embeddedDocuments := 0
	var currentCollection *CollectionInfo
	inCollections := false

	flushCollection := func() {
		if currentCollection != nil {
			status.Collections = append(status.Collections, *currentCollection)
			currentCollection = nil
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "Total:"):
			total, err := parseSingleInteger(totalPattern, trimmed)
			if err != nil {
				return IndexStatus{}, fmt.Errorf("parse total documents: %w", err)
			}
			status.TotalDocuments = total
		case strings.HasPrefix(trimmed, "Vectors:"):
			vectors, err := parseSingleInteger(vectorsPattern, trimmed)
			if err != nil {
				return IndexStatus{}, fmt.Errorf("parse vector count: %w", err)
			}
			embeddedDocuments = vectors
		case strings.HasPrefix(trimmed, "Pending:"):
			pending, err := parseSingleInteger(pendingPattern, trimmed)
			if err != nil {
				return IndexStatus{}, fmt.Errorf("parse pending embeddings: %w", err)
			}
			status.NeedsEmbedding = pending
		case trimmed == "Collections":
			inCollections = true
		case inCollections && strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    "):
			flushCollection()
			headerMatches := collectionHeaderPattern.FindStringSubmatch(trimmed)
			if headerMatches == nil {
				inCollections = false
				continue
			}
			currentCollection = &CollectionInfo{
				Name: headerMatches[1],
				Path: headerMatches[2],
			}
		case inCollections && strings.HasPrefix(line, "    ") && currentCollection != nil:
			switch {
			case strings.HasPrefix(trimmed, "Pattern:"):
				currentCollection.Pattern = strings.TrimSpace(strings.TrimPrefix(trimmed, "Pattern:"))
			case strings.HasPrefix(trimmed, "Files:"):
				filesMatches := collectionFilesPattern.FindStringSubmatch(trimmed)
				if filesMatches == nil {
					return IndexStatus{}, fmt.Errorf("parse collection file summary: %q", trimmed)
				}

				documents, err := strconv.Atoi(filesMatches[1])
				if err != nil {
					return IndexStatus{}, err
				}
				currentCollection.Documents = documents
				currentCollection.LastUpdated = filesMatches[2]
			}
		case inCollections && !strings.HasPrefix(line, " "):
			flushCollection()
			inCollections = false
		}
	}

	flushCollection()
	status.HasVectorIndex = embeddedDocuments > 0 || (status.TotalDocuments > 0 && status.NeedsEmbedding == 0)

	if status.TotalDocuments == 0 && len(status.Collections) == 0 {
		if strings.Contains(text, "No collections.") {
			return status, nil
		}
		return IndexStatus{}, fmt.Errorf("missing status summary in %q", text)
	}

	return status, nil
}

func parseSingleInteger(pattern *regexp.Regexp, input string) (int, error) {
	matches := pattern.FindStringSubmatch(input)
	if matches == nil {
		return 0, fmt.Errorf("no integer match in %q", input)
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	return value, nil
}

func parseHumanDurationMilliseconds(input string) (int, error) {
	text := strings.TrimSpace(input)
	if text == "" {
		return 0, nil
	}

	totalMilliseconds := 0.0
	matches := durationTokenPattern.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("unsupported duration %q", input)
	}

	for _, match := range matches {
		value, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, err
		}

		switch match[2] {
		case "ms":
			totalMilliseconds += value
		case "s":
			totalMilliseconds += value * 1000
		case "m":
			totalMilliseconds += value * 60 * 1000
		case "h":
			totalMilliseconds += value * 60 * 60 * 1000
		default:
			return 0, fmt.Errorf("unsupported duration unit %q", match[2])
		}
	}

	return int(totalMilliseconds), nil
}

func cleanOutput(output string) string {
	return strings.TrimSpace(cleanDiagnostics(output))
}

func cleanDiagnostics(output string) string {
	cleaned := ansiPattern.ReplaceAllString(output, "")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")
	return strings.TrimSpace(cleaned)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
