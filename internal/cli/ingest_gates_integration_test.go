//go:build integration

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/fakes"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/gate"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/refs"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/topic"
)

// gateEnv is a vault with one topic ("demo") wired to fake OpenRouter and
// Firecrawl servers through the environment, so the real commands and
// session opening run against them.
type gateEnv struct {
	vault string
	root  string
	or    *fakes.OpenRouter
	fc    *fakes.Firecrawl
}

// useFakeDecisionModel points the decision model at a fake OpenRouter that
// answers with its defaults (on-topic enough, no quality problem) and valid
// generated literals; ingest commands require the decision model.
func useFakeDecisionModel(t *testing.T) *fakes.OpenRouter {
	t.Helper()
	fake := fakes.NewOpenRouter(nil, fakeLiterals)
	t.Cleanup(fake.Close)
	t.Setenv("OPENROUTER_API_KEY", "test-key")
	t.Setenv("OPENROUTER_API_URL", fake.URL)
	return fake
}

// fakeLiterals answers the classify generation call with valid literals.
func fakeLiterals(schemaName, _, _ string) any {
	if schemaName == "document_literals" {
		return map[string]any{
			"summary":   "Notes collected for the topic and judged at ingest.",
			"entities":  []string{},
			"questions": []string{"What does the document cover?", "Why was it collected?", "How is it used?"},
		}
	}
	return map[string]any{}
}

func newGateEnv(t *testing.T, decide fakes.DecideFunc, scrape fakes.ScrapeFunc) *gateEnv {
	t.Helper()
	or := fakes.NewOpenRouter(decide, fakeLiterals)
	t.Cleanup(or.Close)
	fc := fakes.NewFirecrawl(scrape)
	t.Cleanup(fc.Close)

	configPath := filepath.Join(t.TempDir(), "kb.toml")
	writeFile(t, configPath, "")
	t.Setenv("APP_CONFIG", configPath)
	t.Setenv("OPENROUTER_API_KEY", "test-key")
	t.Setenv("OPENROUTER_API_URL", or.URL)
	t.Setenv("FIRECRAWL_API_KEY", "fc-key")
	t.Setenv("FIRECRAWL_API_URL", fc.URL)

	vault := t.TempDir()
	info, err := topic.New(vault, "demo", "Demo", "demo")
	if err != nil {
		t.Fatalf("topic.New: %v", err)
	}
	return &gateEnv{vault: vault, root: info.RootPath, or: or, fc: fc}
}

// acceptContract activates a contract; gates adds `decisions.gates: <mode>`.
func (e *gateEnv) acceptContract(t *testing.T, gates string) {
	t.Helper()
	c := &contract.Contract{
		Purpose:    "Research on decision models that judge text with typed questions.",
		Core:       []string{"decision models and typed judgments"},
		OutOfScope: []string{"cooking recipes with no decision method"},
	}
	if err := contract.SetContract(e.root, c); err != nil {
		t.Fatal(err)
	}
	if gates != "" {
		path := filepath.Join(e.root, "topic.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		writeFile(t, path, string(data)+"decisions:\n  gates: "+gates+"\n")
	}
}

// run executes a command and returns stdout, stderr and the error.
func (e *gateEnv) run(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	command := newRootCommand()
	var stdout, stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs(append(args, "--vault", e.vault))
	err := command.ExecuteContext(context.Background())
	return stdout.String(), stderr.String(), err
}

func (e *gateEnv) mustRun(t *testing.T, args ...string) (string, string) {
	t.Helper()
	stdout, stderr, err := e.run(t, args...)
	if err != nil {
		t.Fatalf("%v: %v\nstderr:\n%s", args, err, stderr)
	}
	return stdout, stderr
}

func (e *gateEnv) ingestURL(t *testing.T, args ...string) kingest.Result {
	t.Helper()
	stdout, _ := e.mustRun(t, append([]string{"ingest", "url", "--topic", "demo"}, args...)...)
	var result kingest.Result
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("ingest url stdout is not a result: %v\n%s", err, stdout)
	}
	return result
}

func (e *gateEnv) read(t *testing.T, rel string) (map[string]any, string) {
	t.Helper()
	return readMarkdownDocument(t, filepath.Join(e.root, filepath.FromSlash(rel)))
}

func (e *gateEnv) file(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(e.root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func (e *gateEnv) pending(t *testing.T, queue string) []review.Item {
	t.Helper()
	items, err := review.Open(e.root, nil).Pending(queue)
	if err != nil {
		t.Fatal(err)
	}
	return items
}

func (e *gateEnv) labels(t *testing.T) []review.Label {
	t.Helper()
	labels, err := review.LoadLabels(e.root)
	if err != nil {
		t.Fatal(err)
	}
	return labels
}

// topicRel strips the "demo/" prefix of an ingest result path.
func topicRel(filePath string) string { return strings.TrimPrefix(filePath, "demo/") }

// article returns a markdown body of about n words on a subject.
func article(subject string, n int) string {
	var body strings.Builder
	body.WriteString("# " + subject + "\n\n")
	for i := range n {
		fmt.Fprintf(&body, "%s%d ", strings.ToLower(strings.Fields(subject)[0]), i)
		if i%15 == 14 {
			body.WriteString(".\n\n")
		}
	}
	return body.String() + "\n"
}

// pages answers scrapes from a URL → response table (404 for unknown URLs).
func pages(table map[string]fakes.ScrapeResponse) fakes.ScrapeFunc {
	return func(req fakes.ScrapeRequest) fakes.ScrapeResponse {
		if response, ok := table[req.URL]; ok {
			return response
		}
		return fakes.ScrapeResponse{Markdown: "# Not Found\n", Title: "Not Found", StatusCode: 404}
	}
}

func scrapedURLs(fc *fakes.Firecrawl) []string {
	urls := make([]string, 0)
	for _, req := range fc.Requests() {
		urls = append(urls, req.URL)
	}
	return urls
}

// TestIngestGatesQualityQuarantine: a 404 capture survives the fresh
// refetch and is quarantined with its reason, recorded in the ledger and the
// gate review queue, without any decision call; so is a site root.
func TestIngestGatesQualityQuarantine(t *testing.T) {
	env := newGateEnv(t, nil, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/": {Markdown: article("Example company home page", 500), Title: "Welcome to our platform"},
	}))

	missing := env.ingestURL(t, "https://example.com/posts/missing")
	if missing.Triage != gate.TriageQuarantined || missing.TriageReason != "error_page" {
		t.Fatalf("404 result = %+v", missing)
	}
	if !strings.HasPrefix(missing.FilePath, "demo/raw/_quarantine/articles/") {
		t.Fatalf("quarantined path = %q", missing.FilePath)
	}
	values, _ := env.read(t, topicRel(missing.FilePath))
	if values["triage"] != "quarantined" || values["triage_reason"] != "error_page" || values["ingest_batch"] == nil {
		t.Fatalf("quarantined frontmatter = %v", values)
	}
	active, err := refs.List(env.root)
	if err != nil || len(active) != 1 || active[0].Reason != "error_page" {
		t.Fatalf("ledger = %+v, %v", active, err)
	}
	requests := env.fc.Requests()
	if len(requests) != 2 || requests[1].MaxAge == nil || *requests[1].MaxAge != 0 {
		t.Fatalf("expected one scrape plus one fresh refetch, got %+v", requests)
	}
	if gated := env.pending(t, review.QueueGate); len(gated) != 1 || gated[0].Action["quarantined"] != true {
		t.Fatalf("gate queue = %+v", gated)
	}

	root := env.ingestURL(t, "https://example.com/")
	if root.Triage != gate.TriageQuarantined || root.TriageReason != "not_an_article" {
		t.Fatalf("root URL result = %+v", root)
	}
	if calls := len(env.or.Calls()); calls != 0 {
		t.Fatalf("code flags must quarantine without decision calls, got %d", calls)
	}
	logText := env.file(t, "log.md")
	if !strings.Contains(logText, "kept 0, review 0, quarantined 1, skipped 0, duplicates 0") || !strings.Contains(logText, "quarantined (not_an_article)") {
		t.Fatalf("log.md entry:\n%s", logText)
	}
}

// TestIngestGatesDecisionsExclude: a source whose would-be path matches
// topic.yaml decisions.exclude is written kept with a clear note, and no
// request about it reaches the fake decisions or generation endpoints
// (spec §14: exclude keeps files out of every call).
func TestIngestGatesDecisionsExclude(t *testing.T) {
	env := newGateEnv(t, nil, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/posts/private-memo": {Markdown: article("Private memo about decision engines", 400), Title: "Private memo"},
	}))
	env.acceptContract(t, "apply")
	path := filepath.Join(env.root, "topic.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, path, string(data)+"  exclude:\n    - raw/articles/private-*.md\n")

	result := env.ingestURL(t, "https://example.com/posts/private-memo")
	if result.Triage != gate.TriageKept || !result.Excluded || result.Note != gate.ExcludedNote || len(result.Undecided) != 0 {
		t.Fatalf("result = %+v", result)
	}
	if !strings.HasPrefix(result.FilePath, "demo/raw/articles/private-") {
		t.Fatalf("path = %q", result.FilePath)
	}
	if calls, gen := len(env.or.Calls()), len(env.or.GenCalls()); calls != 0 || gen != 0 {
		t.Fatalf("an excluded source must not reach any model: %d decision and %d generation calls", calls, gen)
	}
	values, _ := env.read(t, topicRel(result.FilePath))
	if values["triage"] != "kept" || values["triage_reason"] != nil || values["summary"] != nil {
		t.Fatalf("frontmatter = %v", values)
	}
	if gated := env.pending(t, review.QueueGate); len(gated) != 0 {
		t.Fatalf("gate queue = %+v", gated)
	}
	logText := env.file(t, "log.md")
	if !strings.Contains(logText, "excluded 1 (decisions.exclude, not judged)") || !strings.Contains(logText, gate.ExcludedNote) {
		t.Fatalf("log.md entry:\n%s", logText)
	}
}

// TestIngestGatesSiteNameTitle: Firecrawl's ogSiteName reaches the code
// quality check, so a capture titled with the site name alone is flagged
// not_an_article without any decision call.
func TestIngestGatesSiteNameTitle(t *testing.T) {
	env := newGateEnv(t, nil, pages(map[string]fakes.ScrapeResponse{
		"https://www.catapultsports.com/products/vector": {Markdown: article("Catapult vector wearable", 500), Title: "Catapult", SiteName: "Catapult"},
	}))
	result := env.ingestURL(t, "https://www.catapultsports.com/products/vector")
	if result.Triage != gate.TriageQuarantined || result.TriageReason != "not_an_article" {
		t.Fatalf("result = %+v", result)
	}
	if calls := len(env.or.Calls()); calls != 0 {
		t.Fatalf("decision calls = %d, want 0", calls)
	}
}

// TestIngestGatesUndecidedIsReview: a gate judgment that fails (an invalid
// receipt for `role`) writes triage: review with triage_reason undecided,
// queues a gate item and counts undecided in the summary and log.md; it is
// never written as kept (spec §2 principle 3).
func TestIngestGatesUndecidedIsReview(t *testing.T) {
	env := newGateEnv(t, func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "role" {
			return map[string]any{"type": "choice", "choice": "core"}
		}
		return nil
	}, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/posts/typed-judges": {Markdown: article("Typed judges in practice", 400), Title: "Typed judges in practice"},
	}))

	stdout, stderr := env.mustRun(t, "ingest", "url", "https://example.com/posts/typed-judges", "--topic", "demo")
	var result kingest.Result
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatal(err)
	}
	if result.Triage != gate.TriageReview || result.TriageReason != gate.ReasonUndecided || result.ReviewItem == "" {
		t.Fatalf("result = %+v", result)
	}
	values, _ := env.read(t, topicRel(result.FilePath))
	if values["triage"] != "review" || values["triage_reason"] != "undecided" {
		t.Fatalf("frontmatter = %v", values)
	}
	gated := env.pending(t, review.QueueGate)
	if len(gated) != 1 || gated[0].Question != "role" || gated[0].Purpose != "relevance" {
		t.Fatalf("gate items = %+v", gated)
	}
	if !strings.Contains(stderr, "undecided 1") || !strings.Contains(env.file(t, "log.md"), "review 1, quarantined 0, skipped 0, duplicates 0, undecided 1") {
		t.Fatalf("summary must count undecided:\n%s\nlog.md:\n%s", stderr, env.file(t, "log.md"))
	}
}

// TestIngestGatesRefetchKeepsFullBody: a thin first capture is refetched
// with maxAge 0; the full body is kept and `quality` is not written.
func TestIngestGatesRefetchKeepsFullBody(t *testing.T) {
	full := article("Decision engines explained", 400)
	env := newGateEnv(t, nil, func(req fakes.ScrapeRequest) fakes.ScrapeResponse {
		if req.MaxAge != nil && *req.MaxAge == 0 {
			return fakes.ScrapeResponse{Markdown: full, Title: "Decision engines explained"}
		}
		return fakes.ScrapeResponse{Markdown: "# Decision engines explained\n\nLoading...\n", Title: "Decision engines explained"}
	})

	result := env.ingestURL(t, "https://blog.example.com/posts/decision-engines")
	if result.Triage != gate.TriageKept || !result.Refetched {
		t.Fatalf("result = %+v", result)
	}
	values, body := env.read(t, topicRel(result.FilePath))
	if body != full {
		t.Fatalf("kept body is not the refetched one:\n%s", body)
	}
	if _, ok := values["quality"]; ok {
		t.Fatalf("quality must not be written: %v", values)
	}
	if values["triage"] != "kept" {
		t.Fatalf("triage = %v", values["triage"])
	}
	if requests := env.fc.Requests(); len(requests) != 2 || requests[0].MaxAge != nil || requests[1].WaitFor == nil {
		t.Fatalf("scrape requests = %+v", requests)
	}
}

// TestIngestGatesShadowDefault: with an accepted contract and no
// calibration, relevance never quarantines (review + shadow note instead),
// while a quality failure is quarantined.
func TestIngestGatesShadowDefault(t *testing.T) {
	env := newGateEnv(t, func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "role" {
			return fakes.Pick("off_topic", 0.95)
		}
		return nil
	}, pages(map[string]fakes.ScrapeResponse{
		"https://food.example.com/recipes/pasta": {Markdown: article("Pasta recipes for the weekend", 400), Title: "Pasta recipes for the weekend"},
	}))
	env.acceptContract(t, "")

	offTopic := env.ingestURL(t, "https://food.example.com/recipes/pasta")
	if offTopic.Triage != gate.TriageReview || offTopic.TriageReason != "off_topic" {
		t.Fatalf("off-topic result = %+v", offTopic)
	}
	if !strings.HasPrefix(offTopic.FilePath, "demo/raw/articles/") || !slices.ContainsFunc(offTopic.Shadow, func(s string) bool {
		return strings.Contains(s, "quarantined:off_topic (relevance shadow)")
	}) {
		t.Fatalf("shadow relevance must not move the source: %+v", offTopic)
	}
	_, stderr := env.mustRun(t, "ingest", "url", "https://food.example.com/recipes/missing", "--topic", "demo")
	if !strings.Contains(stderr, "gate relevance: shadow (not calibrated") || !strings.Contains(stderr, "gate quality: apply") {
		t.Fatalf("run summary must state the gate modes:\n%s", stderr)
	}
	if active, _ := refs.List(env.root); len(active) != 1 || active[0].Reason != "error_page" {
		t.Fatalf("quality failure must be quarantined: %+v", active)
	}
}

// TestIngestGatesReviewBand: a quality noul in the review band writes
// triage: review, triage_reason and quality, and queues a gate item.
func TestIngestGatesReviewBand(t *testing.T) {
	env := newGateEnv(t, func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "thin_or_boilerplate" {
			return fakes.Noul(0.6)
		}
		return nil
	}, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/docs/pricing": {Markdown: article("Pricing tiers overview", 400), Title: "Pricing tiers overview"},
	}))

	result := env.ingestURL(t, "https://example.com/docs/pricing")
	if result.Triage != gate.TriageReview || result.TriageReason != "thin" || result.ReviewItem == "" {
		t.Fatalf("result = %+v", result)
	}
	values, _ := env.read(t, topicRel(result.FilePath))
	if values["triage"] != "review" || values["triage_reason"] != "thin" || values["quality"] != "thin" {
		t.Fatalf("frontmatter = %v", values)
	}
	gated := env.pending(t, review.QueueGate)
	if len(gated) != 1 || gated[0].Question != "thin_or_boilerplate" || gated[0].Purpose != "quality" {
		t.Fatalf("gate items = %+v", gated)
	}
}

// TestIngestGatesExactDedupeAndForce: re-ingesting the same URL (tracking
// parameters and trailing slash aside) is skipped before any fetch; --force
// ingests it again.
func TestIngestGatesExactDedupeAndForce(t *testing.T) {
	env := newGateEnv(t, nil, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/posts/typed-judges": {Markdown: article("Typed judges in practice", 400), Title: "Typed judges in practice"},
	}))

	first := env.ingestURL(t, "https://example.com/posts/typed-judges")
	if first.Triage != gate.TriageKept {
		t.Fatalf("first = %+v", first)
	}
	scrapes := len(env.fc.Requests())
	again := env.ingestURL(t, "https://www.example.com/posts/typed-judges/?utm_source=feed")
	if again.Triage != gate.TriageDuplicateSkipped || again.DuplicateOf != topicRel(first.FilePath) || again.FilePath != "" {
		t.Fatalf("duplicate = %+v", again)
	}
	if len(env.fc.Requests()) != scrapes {
		t.Fatal("a duplicate URL must not be fetched")
	}
	forced := env.ingestURL(t, "https://example.com/posts/typed-judges", "--force")
	if forced.Triage != gate.TriageKept || forced.FilePath == first.FilePath {
		t.Fatalf("forced = %+v", forced)
	}
}

// TestIngestGatesYouTubeDedupeBeforeExtraction: a single YouTube video
// already in the topic (quarantined included) is skipped before yt-dlp.
func TestIngestGatesYouTubeDedupeBeforeExtraction(t *testing.T) {
	env := newGateEnv(t, nil, nil)
	writeMarkdownDocument(t, env.root, "raw/_quarantine/youtube/old-talk.md", map[string]any{
		"title": "Old talk", "type": "source", "source_kind": "youtube-transcript",
		"source_url": "https://www.youtube.com/watch?v=abcdefghijk", "video_id": "abcdefghijk", "triage": "quarantined",
	}, "transcript\n")

	stdout, _ := env.mustRun(t, "ingest", "youtube", "https://youtu.be/abcdefghijk", "--topic", "demo")
	var result kingest.Result
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatal(err)
	}
	if result.Triage != gate.TriageDuplicateSkipped || result.DuplicateOf != "raw/_quarantine/youtube/old-talk.md" {
		t.Fatalf("result = %+v", result)
	}
	if !strings.Contains(env.file(t, "log.md"), "duplicate of quarantined raw/_quarantine/youtube/old-talk.md") {
		t.Fatalf("log.md:\n%s", env.file(t, "log.md"))
	}
}

// TestIngestGatesNearDuplicateSupersedes: a new version of an existing
// same-site source is kept and records `supersedes`.
func TestIngestGatesNearDuplicateSupersedes(t *testing.T) {
	env := newGateEnv(t, func(call fakes.Call, q fakes.Question) any {
		if strings.HasPrefix(q.ID, "relation_c") && strings.Contains(call.StateString(), "Typed judges handbook") {
			return fakes.Pick("new_version", 0.9)
		}
		return nil
	}, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/handbook/v2": {Markdown: article("Typed judges handbook second edition", 420), Title: "Typed judges handbook (2nd edition)"},
	}))
	writeMarkdownDocument(t, env.root, "raw/articles/typed-judges-handbook.md", map[string]any{
		"title": "Typed judges handbook", "type": "source", "source_kind": "article",
		"source_url": "https://example.com/handbook/v1",
	}, article("Typed judges handbook first edition", 400))

	result := env.ingestURL(t, "https://example.com/handbook/v2")
	if result.Triage != gate.TriageKept || !slices.Equal(result.Supersedes, []string{"[[typed-judges-handbook]]"}) {
		t.Fatalf("result = %+v", result)
	}
	values, _ := env.read(t, topicRel(result.FilePath))
	if got := frontmatter.GetStringSlice(values, "supersedes"); !slices.Equal(got, []string{"[[typed-judges-handbook]]"}) {
		t.Fatalf("supersedes = %v", got)
	}
}

// TestIngestGatesBulkSkipAndRescue: in apply mode an off-topic URL of a list
// is never fetched, is logged in skipped.jsonl and queued as a skip item;
// --rescue ingests it and records a positive label.
func TestIngestGatesBulkSkipAndRescue(t *testing.T) {
	env := newGateEnv(t, func(call fakes.Call, q fakes.Question) any {
		if !strings.HasPrefix(q.ID, "role_item_") {
			return nil
		}
		id := strings.TrimPrefix(q.ID, "role_item_")
		var state struct {
			Items []struct{ ID, Title string } `json:"items"`
		}
		_ = json.Unmarshal(call.State, &state)
		for _, item := range state.Items {
			if item.ID == id && strings.Contains(item.Title, "Pasta") {
				return fakes.Pick("off_topic", 0.95)
			}
		}
		return fakes.Pick("core", 0.9)
	}, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/posts/judges": {Markdown: article("Judges with typed questions", 400), Title: "Judges with typed questions"},
		"https://food.example.com/pasta":   {Markdown: article("Pasta for the weekend", 400), Title: "Pasta for the weekend"},
	}))
	env.acceptContract(t, "apply")
	list := filepath.Join(t.TempDir(), "reading-list.txt")
	writeFile(t, list, "# reading list\nhttps://example.com/posts/judges Judges with typed questions\nhttps://food.example.com/pasta Pasta for the weekend\n")

	stdout, stderr := env.mustRun(t, "ingest", "url", "--from", list, "--topic", "demo")
	var summary urlListSummary
	if err := json.Unmarshal([]byte(stdout), &summary); err != nil {
		t.Fatalf("summary: %v\n%s", err, stdout)
	}
	if len(summary.Results) != 1 || len(summary.Skipped) != 1 || summary.Skipped[0].Reason != "off_topic" {
		t.Fatalf("summary = %+v\n%s", summary, stderr)
	}
	if slices.Contains(scrapedURLs(env.fc), "https://food.example.com/pasta") {
		t.Fatal("a skipped item must never be fetched")
	}
	rows, err := gate.LoadSkipped(env.root)
	if err != nil || len(rows) != 1 || rows[0].URL != "https://food.example.com/pasta" || rows[0].Rescued || rows[0].Query != "reading-list.txt" {
		t.Fatalf("skipped.jsonl = %+v, %v", rows, err)
	}
	skips := env.pending(t, review.QueueSkip)
	if len(skips) != 1 || skips[0].Action["skipped_id"] != rows[0].ID {
		t.Fatalf("skip items = %+v", skips)
	}
	kept, _ := env.read(t, topicRel(summary.Results[0].FilePath))
	if kept["ingest_query"] != "Judges with typed questions" || kept["triage"] != "kept" {
		t.Fatalf("kept frontmatter = %v", kept)
	}

	rescued := env.ingestURL(t, "--rescue", rows[0].ID)
	if rescued.Triage == "" || !strings.HasPrefix(rescued.FilePath, "demo/raw/articles/") {
		t.Fatalf("rescued = %+v", rescued)
	}
	if !slices.Contains(scrapedURLs(env.fc), "https://food.example.com/pasta") {
		t.Fatal("rescue must fetch the skipped item")
	}
	rows, _ = gate.LoadSkipped(env.root)
	if !rows[0].Rescued {
		t.Fatalf("row not marked rescued: %+v", rows[0])
	}
	if len(env.pending(t, review.QueueSkip)) != 0 {
		t.Fatal("the skip item must be resolved by the rescue")
	}
	labels := env.labels(t)
	if len(labels) != 1 || labels[0].Verdict != review.VerdictPositive || labels[0].Purpose != "relevance" {
		t.Fatalf("labels = %+v", labels)
	}
}

// TestReviewAcceptRestoresGateQuarantine: rejecting a review-band source
// quarantines it (index line and sources entry removed); accepting its gate
// quarantine item restores it with every touched file byte-identical.
func TestReviewAcceptRestoresGateQuarantine(t *testing.T) {
	env := newGateEnv(t, func(_ fakes.Call, q fakes.Question) any {
		if q.ID == "paywall_or_login" {
			return fakes.Noul(0.6)
		}
		return nil
	}, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/posts/judging": {Markdown: article("Judging documents cheaply", 400), Title: "Judging documents cheaply"},
	}))
	result := env.ingestURL(t, "https://example.com/posts/judging")
	rel := topicRel(result.FilePath)
	if result.Triage != gate.TriageReview {
		t.Fatalf("result = %+v", result)
	}
	stem := strings.TrimSuffix(filepath.Base(rel), ".md")
	writeFile(t, filepath.Join(env.root, "wiki/index/Source Index.md"), "# Source Index\n\n- [["+stem+"]] — judging\n- other line\n")
	writeMarkdownDocument(t, env.root, "wiki/concepts/Judges.md", map[string]any{
		"title": "Judges", "type": "wiki", "stage": "compiled", "sources": []string{"[[" + stem + "]]", "[[elsewhere]]"},
	}, "# Judges\n\nSee [["+stem+"]].\n")
	indexBefore, articleBefore := env.file(t, "wiki/index/Source Index.md"), env.file(t, "wiki/concepts/Judges.md")
	sourceBefore := env.file(t, rel)

	stdout, _ := env.mustRun(t, "review", "reject", "--topic", "demo", result.ReviewItem)
	if !strings.Contains(stdout, "rejected "+result.ReviewItem) {
		t.Fatalf("reject output:\n%s", stdout)
	}
	if strings.Contains(env.file(t, "wiki/index/Source Index.md"), stem) || strings.Contains(frontmatterSources(t, env), stem) {
		t.Fatal("quarantine must remove the index line and the sources entry")
	}
	active, _ := refs.List(env.root)
	if len(active) != 1 || active[0].Reason != "paywall" {
		t.Fatalf("active quarantines = %+v", active)
	}

	// The gate queue item a gate quarantine produces.
	item, _ := gate.ReviewItem(gate.Outcome{Triage: gate.TriageQuarantined, Reason: "paywall", Stage: gate.StageQuality, Purpose: "quality", Question: "paywall_or_login", Probability: 0.85}, rel, "Judging documents cheaply")
	item.Question = "paywall_or_login:requarantine"
	if _, err := review.Open(env.root, nil).Add(item); err != nil {
		t.Fatal(err)
	}
	gated := env.pending(t, review.QueueGate)
	if len(gated) != 1 {
		t.Fatalf("gate items = %+v", gated)
	}
	stdout, _ = env.mustRun(t, "review", "accept", "--topic", "demo", gated[0].ID)
	if !strings.Contains(stdout, "restored "+rel) {
		t.Fatalf("accept output:\n%s", stdout)
	}
	if env.file(t, "wiki/index/Source Index.md") != indexBefore || env.file(t, "wiki/concepts/Judges.md") != articleBefore {
		t.Fatal("restore must put the index line and the sources entry back byte for byte")
	}
	restored := env.file(t, rel)
	if want := strings.Replace(strings.Replace(sourceBefore, "triage: review\n", "triage: kept\n", 1), "triage_reason: paywall\n", "", 1); restored != want {
		t.Fatalf("restored source differs beyond its triage keys:\n--- want\n%s\n--- got\n%s", want, restored)
	}
	labels := env.labels(t)
	if len(labels) != 2 || labels[0].Verdict != review.VerdictNegative || labels[1].Verdict != review.VerdictPositive || labels[1].Purpose != "quality" {
		t.Fatalf("labels = %+v", labels)
	}
}

func frontmatterSources(t *testing.T, env *gateEnv) string {
	t.Helper()
	values, _ := env.read(t, "wiki/concepts/Judges.md")
	return strings.Join(frontmatter.GetStringSlice(values, "sources"), ",")
}

// TestReviewAcceptRecaptureRemoveAndLink covers the backfill and link
// queues: recapture replaces the body in place, remove quarantines as
// off_topic, a link insert writes the relation and the body link; every
// verdict appends a label.
func TestReviewAcceptRecaptureRemoveAndLink(t *testing.T) {
	full := article("Calibrating decision thresholds", 400)
	env := newGateEnv(t, nil, pages(map[string]fakes.ScrapeResponse{
		"https://example.com/posts/calibration": {Markdown: full, Title: "Calibrating decision thresholds"},
	}))
	writeMarkdownDocument(t, env.root, "raw/articles/calibration.md", map[string]any{
		"title": "Calibrating decision thresholds", "type": "source", "source_kind": "article",
		"source_url": "https://example.com/posts/calibration", "custom": "kept as is",
	}, "# Calibrating\n\nLoading...\n")
	writeMarkdownDocument(t, env.root, "raw/articles/recipes.md", map[string]any{
		"title": "Weekend recipes", "type": "source", "source_kind": "article",
	}, article("Weekend recipes", 300))
	writeMarkdownDocument(t, env.root, "wiki/concepts/Decision Models.md", map[string]any{
		"title": "Decision Models", "type": "wiki", "stage": "compiled",
	}, "# Decision Models\n")
	body := "# Notes\n\nWe compare decision models across tasks.\n"
	writeMarkdownDocument(t, env.root, "raw/articles/notes.md", map[string]any{
		"title": "Notes", "type": "source", "source_kind": "document",
	}, body)

	store := review.Open(env.root, nil)
	mentionStart := strings.Index("We compare decision models across tasks.\n", "decision models") + len("# Notes\n\n")
	items := []review.Item{
		{Queue: review.QueueRecapture, Purpose: "quality", Subject: "raw/articles/calibration.md", Question: "code:thin", Probability: 1, Action: map[string]any{"reason": "thin"}},
		{Queue: review.QueueRemove, Purpose: "relevance", Subject: "raw/articles/recipes.md", Question: "role", Probability: 0.9, Action: map[string]any{"reason": "off_topic"}},
		{Queue: review.QueueLink, Purpose: "mention", Subject: "raw/articles/notes.md", Target: "wiki/concepts/Decision Models.md", Question: "mention_sense", Probability: 0.9,
			Action: map[string]any{"insert": true, "relation": "related", "target": "Decision Models", "mention_text": "decision models", "mention_start": mentionStart}},
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		item.ID = review.ItemID(item.Queue, item.Subject, item.Target, item.Question)
		if _, err := store.Add(item); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, item.ID)
	}

	stdout, _ := env.mustRun(t, "review", "accept", "--topic", "demo", ids[0])
	if !strings.Contains(stdout, "recaptured raw/articles/calibration.md") {
		t.Fatalf("recapture output:\n%s", stdout)
	}
	values, gotBody := env.read(t, "raw/articles/calibration.md")
	if gotBody != full || values["custom"] != "kept as is" || values["recaptured"] == nil {
		t.Fatalf("recaptured document: %v\n%s", values, gotBody)
	}

	env.mustRun(t, "review", "accept", "--topic", "demo", "--queue", "remove", "--above", "0.8")
	if _, err := os.Stat(filepath.Join(env.root, "raw/_quarantine/articles/recipes.md")); err != nil {
		t.Fatalf("remove accept must quarantine: %v", err)
	}
	if active, _ := refs.List(env.root); len(active) != 1 || active[0].Reason != "off_topic" {
		t.Fatalf("quarantines = %+v", active)
	}

	env.mustRun(t, "review", "accept", "--topic", "demo", ids[2])
	values, gotBody = env.read(t, "raw/articles/notes.md")
	if got := frontmatter.GetStringSlice(values, "related"); !slices.Equal(got, []string{"[[Decision Models]]"}) {
		t.Fatalf("related = %v", got)
	}
	if !strings.Contains(gotBody, "We compare [[Decision Models|decision models]] across tasks.") {
		t.Fatalf("body link not inserted:\n%s", gotBody)
	}
	if !strings.Contains(env.file(t, ".decisions/inserted-links.jsonl"), `"subject":"raw/articles/notes.md"`) {
		t.Fatal("inserted link must be logged")
	}

	labels := env.labels(t)
	verdicts := make([]string, 0, len(labels))
	for _, label := range labels {
		verdicts = append(verdicts, label.Purpose+":"+label.Verdict)
	}
	if want := []string{"quality:negative", "relevance:negative", "mention:positive"}; !slices.Equal(verdicts, want) {
		t.Fatalf("labels = %v, want %v", verdicts, want)
	}
	if left := env.pending(t, ""); len(left) != 0 {
		t.Fatalf("pending after accepts = %+v", left)
	}
}
