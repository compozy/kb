package ingest

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/firecrawl"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/gate"
	"github.com/compozy/kb/internal/link"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/refs"
	"github.com/compozy/kb/internal/resolve"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/scope"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/topic"
	"github.com/compozy/kb/internal/vault"
)

// GateOptions are the per-source gate inputs of a decision-backed ingest.
type GateOptions struct {
	// Force skips gates 1–6 (the source is written as kept); classify and
	// link still run.
	Force bool
	// Refetch performs the one fresh stage-4 scrape of a URL source (nil for
	// other sources).
	Refetch func(ctx context.Context) (*firecrawl.ScrapeResult, error)
	// RequestedURL is the URL the fetch was asked for; FinalURL and
	// StatusCode are what the fetcher reported.
	RequestedURL string
	FinalURL     string
	StatusCode   int
	// SiteName is the site name the fetcher reported (Firecrawl
	// ogSiteName), for the "title equals the site name" quality rule.
	SiteName string
	// PlatformID is the `youtube:<id>` / `instagram:<shortcode>` identity.
	PlatformID string
	// Prefetch is the stage-2 verdict of a bulk item (nil for single
	// sources and rescued items).
	Prefetch *gate.PrefetchDecision
}

// Result is the outcome of one gated source: the classic ingest result plus
// what the gates decided. FilePath is empty for skipped sources and points
// into raw/_quarantine/ for quarantined ones.
type Result struct {
	models.IngestResult
	Triage       string   `json:"triage,omitempty"`
	TriageReason string   `json:"triage_reason,omitempty"`
	Quality      string   `json:"quality,omitempty"`
	DuplicateOf  string   `json:"duplicate_of,omitempty"`
	Shadow       []string `json:"shadow,omitempty"`
	Supersedes   []string `json:"supersedes,omitempty"`
	Refetched    bool     `json:"refetched,omitempty"`
	ReviewItem   string   `json:"review_item,omitempty"`
	Undecided    []string `json:"undecided,omitempty"`
	// Excluded reports that the source matches decisions.exclude: it was
	// written without any decision call and is left out of classify and
	// link; Note says so.
	Excluded bool   `json:"excluded,omitempty"`
	Note     string `json:"note,omitempty"`
}

// RunOptions configures a decision-backed ingest run.
type RunOptions struct {
	// Command names the ingest subcommand ("url", "channel", ...).
	Command string
	// Batch and Query are the run provenance (ingest_batch, ingest_query);
	// per-source Options.Batch/Query override them.
	Batch string
	Query string
	// Force skips gates 1–6 for every source of the run.
	Force bool
}

// Run is one decision-backed ingest run over a topic (spec §7): every
// source goes through the gates, is written with its triage and
// provenance through the owned-key writer, quarantined through
// internal/refs when a gate says so, and classified and linked at Finish.
// Sources of one run are processed one at a time.
type Run struct {
	s    *session.Session
	opts RunOptions

	mu      sync.Mutex
	checker *gate.Checker
	queue   *review.Store
	tally   gate.Tally
	entries []logEntry
	written []string
}

// logEntry is one source line of the run's log.md entry.
type logEntry struct {
	triage, reason, title, path, detail string
}

// NewRun starts a run on a session.
func NewRun(s *session.Session, opts RunOptions) *Run {
	return &Run{s: s, opts: opts, queue: review.Open(s.Root(), s.Now)}
}

// Session returns the run's session.
func (r *Run) Session() *session.Session { return r.s }

func (r *Run) gates() (*gate.Checker, error) {
	if r.checker != nil {
		return r.checker, nil
	}
	checker, err := gate.NewChecker(r.s)
	if err != nil {
		return nil, fmt.Errorf("ingest: %w", err)
	}
	r.checker = checker
	return checker, nil
}

// Precheck is stage 1 before fetching: a source whose URL or platform id
// already exists in the topic (quarantined sources included) is skipped and
// reported, so no transcript, STT or scrape is paid for it. It returns the
// skip result and true; false means fetch. Force and shadow never skip.
func (r *Run) Precheck(sourceURL, title string, kind models.SourceKind) (Result, bool, error) {
	if r.opts.Force {
		return Result{}, false, nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	checker, err := r.gates()
	if err != nil {
		return Result{}, false, err
	}
	out, skip := checker.Skip(gate.Ref{URL: sourceURL})
	if !skip {
		return Result{}, false, nil
	}
	if strings.TrimSpace(title) == "" {
		title = sourceURL
	}
	return r.recordSkip(out, title, kind), true, nil
}

// Prefetch is stage 2 over bulk items (see gate.Prefetch). Skipped items
// are counted and logged; the caller fetches the rest.
func (r *Run) Prefetch(ctx context.Context, items []gate.Item) ([]gate.PrefetchDecision, error) {
	if r.opts.Force || len(items) == 0 {
		out := make([]gate.PrefetchDecision, len(items))
		for index, item := range items {
			if item.ID == "" {
				item.ID = gate.SkippedID(item.URL)
			}
			out[index] = gate.PrefetchDecision{Item: item, Action: gate.ActionFetch}
		}
		return out, nil
	}
	decided, err := gate.Prefetch(ctx, r.s, items, gate.PrefetchOptions{Batch: r.opts.Batch, Query: r.opts.Query})
	if err != nil {
		return decided, fmt.Errorf("ingest: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, decision := range decided {
		if decision.Action != gate.ActionSkip {
			continue
		}
		out := gate.Outcome{Triage: gate.TriageSkipped, Reason: gate.ReasonOffTopic, Stage: gate.StagePrefetch}
		r.tally.Add(out)
		title := decision.Item.Title
		if title == "" {
			title = decision.Item.URL
		}
		r.entries = append(r.entries, logEntry{
			triage: gate.TriageSkipped, reason: gate.ReasonOffTopic, title: title,
			detail: fmt.Sprintf("%s, P(off_topic) %.2f, rescue with --rescue %s", decision.Item.URL, decision.POffTopic, decision.Item.ID),
		})
	}
	return decided, nil
}

// Skipped returns the skipped.jsonl row of a skipped id (for --rescue).
func (r *Run) Skipped(id string) (gate.SkippedRow, error) {
	return gate.FindSkipped(r.s.Root(), id)
}

// MarkRescued records that a skipped item was ingested.
func (r *Run) MarkRescued(row gate.SkippedRow) error {
	return gate.MarkRescued(r.s.Root(), row, r.s.Now().UTC().Format(time.RFC3339))
}

// Ingest runs one fetched source through gates 1 and 3–6 (stage 2 already
// ran for bulk items; options.Gate.Prefetch carries its verdict), writes it
// with provenance and triage through the owned-key writer, quarantines it
// through internal/refs when the outcome says so, and queues a `gate`
// review item for the review band and gate quarantines. Classification and
// linking of the written sources happen in Finish.
func (r *Run) Ingest(ctx context.Context, options Options) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	s := r.s
	topicInfo := s.Topic
	if err := topic.EnsureCurrentSkeleton(topicInfo.RootPath); err != nil {
		return Result{}, fmt.Errorf("ingest: ensure topic skeleton: %w", err)
	}
	sourceDirectory, err := rawDirectoryForSourceKind(options.SourceKind)
	if err != nil {
		return Result{}, fmt.Errorf("ingest: %w", err)
	}
	title, markdown, err := resolveMarkdown(ctx, options)
	if err != nil {
		return Result{}, fmt.Errorf("ingest: %w", err)
	}
	now := options.ScrapedAt.UTC()
	if options.ScrapedAt.IsZero() {
		now = s.Now().UTC()
	}
	options.Batch = firstNonEmpty(strings.TrimSpace(options.Batch), r.opts.Batch)
	options.Query = firstNonEmpty(strings.TrimSpace(options.Query), r.opts.Query)
	force := r.opts.Force || options.Gate.Force

	checker, err := r.gates()
	if err != nil {
		return Result{}, err
	}
	if !force {
		if out, skip := checker.Skip(gate.Ref{URL: options.SourceURL, PlatformID: options.Gate.PlatformID}); skip {
			return r.recordSkip(out, title, options.SourceKind), nil
		}
	}

	slug, err := uniqueGatedSlug(topicInfo.RootPath, sourceDirectory, vault.SlugifySegment(title))
	if err != nil {
		return Result{}, fmt.Errorf("ingest: allocate slug: %w", err)
	}
	rel := path.Join("raw", sourceDirectory, slug+".md")
	built := &builder{topicInfo: topicInfo, options: options, rel: rel, now: now}

	var out gate.Outcome
	var doc *corpus.Document
	if force {
		out = gate.Outcome{Triage: gate.TriageKept, Stage: "force"}
		doc, err = built.build(title, markdown)
	} else {
		out, doc, err = checker.Evaluate(ctx, fetched(options, built, title, markdown))
	}
	if err != nil {
		return Result{}, err
	}
	if out.Triage == gate.TriageDuplicateSkipped {
		return r.recordSkip(out, doc.Title, options.SourceKind), nil
	}
	return r.write(ctx, built, doc, out)
}

// fetched assembles the gate input of one source.
func fetched(options Options, built *builder, title, markdown string) gate.Fetched {
	in := gate.Fetched{
		Title: title, Markdown: markdown, Build: built.build,
		RequestedURL: options.Gate.RequestedURL, FinalURL: options.Gate.FinalURL, StatusCode: options.Gate.StatusCode,
		SiteName: options.Gate.SiteName,
		Refetch:  options.Gate.Refetch, PlatformID: options.Gate.PlatformID,
		Transcript: classify.IsTranscriptKind(string(options.SourceKind)),
	}
	if decision := options.Gate.Prefetch; decision != nil {
		in.PrefetchReview = decision.Action == gate.ActionReview
		in.PrefetchShadow = decision.Shadow
		in.PrefetchPOffTopic = decision.POffTopic
		in.PrefetchReceipt = decision.ReceiptKey
	}
	return in
}

// write stores a kept, review or quarantined source.
func (r *Run) write(ctx context.Context, built *builder, doc *corpus.Document, out gate.Outcome) (Result, error) {
	s := r.s
	root := s.Root()
	base, err := built.values(doc.Title, built.markdown, false)
	if err != nil {
		return Result{}, fmt.Errorf("ingest: %w", err)
	}
	content, err := frontmatter.Generate(base, built.markdown)
	if err != nil {
		return Result{}, fmt.Errorf("ingest: generate frontmatter: %w", err)
	}
	absolute := filepath.Join(root, filepath.FromSlash(built.rel))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		return Result{}, fmt.Errorf("ingest: create target directory: %w", err)
	}
	if err := writeExclusive(absolute, content); err != nil {
		return Result{}, fmt.Errorf("ingest: write %q: %w", absolute, err)
	}

	stored, err := corpus.ReadDocument(root, built.rel, corpus.KindSource)
	if err != nil {
		return Result{}, fmt.Errorf("ingest: %w", err)
	}
	updates := map[string]any{}
	if batch := strings.TrimSpace(built.options.Batch); batch != "" {
		updates["ingest_batch"] = batch
	}
	if query := strings.TrimSpace(built.options.Query); query != "" {
		updates["ingest_query"] = query
	}
	if out.Triage != gate.TriageQuarantined {
		updates["triage"] = out.Triage
	}
	if out.Triage == gate.TriageReview && out.Reason != "" {
		updates["triage_reason"] = out.Reason
	}
	if out.Quality != "" {
		updates["quality"] = out.Quality
	}
	if len(out.Supersedes) > 0 {
		updates["supersedes"] = out.Supersedes
	}
	if _, err := s.Writer.Apply(stored, updates, s.StateMeta()); err != nil {
		return Result{}, fmt.Errorf("ingest: write owned keys: %w", err)
	}

	finalRel := built.rel
	quarantined := out.Triage == gate.TriageQuarantined
	if quarantined {
		moved, err := refs.Quarantine(ctx, refs.Options{
			VaultPath: s.VaultPath, TopicRoot: root, Path: built.rel, Reason: out.Reason,
			Writer: s.Writer, State: s.State, Now: s.Now,
		})
		if err != nil {
			return Result{}, fmt.Errorf("ingest: quarantine %s: %w", built.rel, err)
		}
		finalRel = moved.NewPath
	}

	result := Result{
		IngestResult: models.IngestResult{
			Topic: s.Topic.Slug, SourceType: built.options.SourceKind,
			FilePath: path.Join(s.Topic.Slug, finalRel), Title: doc.Title,
		},
		Triage: out.Triage, TriageReason: out.Reason, Quality: out.Quality, DuplicateOf: out.DuplicateOf,
		Shadow: out.Shadow, Supersedes: out.Supersedes, Refetched: out.Refetched, Undecided: out.Undecided,
		Excluded: out.Excluded,
	}
	if out.Excluded {
		result.Note = gate.ExcludedNote
	}
	if item, ok := gate.ReviewItem(out, built.rel, doc.Title); ok {
		item.ID = review.ItemID(item.Queue, item.Subject, item.Target, item.Question)
		if _, err := r.queue.Add(item); err != nil {
			return result, fmt.Errorf("ingest: review queue: %w", err)
		}
		result.ReviewItem = item.ID
	}

	if refreshed, err := corpus.ReadDocument(root, finalRel, corpus.KindSource); err == nil {
		r.checker.Register(refreshed, quarantined)
	}
	if !quarantined && !out.Excluded {
		// Excluded sources stay out of every call, so they are not
		// classified or linked either.
		r.written = append(r.written, built.rel)
	}
	r.tally.Add(out)
	details := make([]string, 0, len(out.Shadow)+1)
	if out.Evidence != "" {
		details = append(details, out.Evidence)
	}
	details = append(details, out.Shadow...)
	r.entries = append(r.entries, logEntry{
		triage: out.Triage, reason: out.Reason, title: doc.Title, path: result.FilePath, detail: strings.Join(details, "; "),
	})
	return result, nil
}

// recordSkip counts and logs a source that is not written.
func (r *Run) recordSkip(out gate.Outcome, title string, kind models.SourceKind) Result {
	r.tally.Add(out)
	r.entries = append(r.entries, logEntry{triage: out.Triage, reason: out.Reason, title: title, detail: out.Evidence})
	return Result{
		IngestResult: models.IngestResult{Topic: r.s.Topic.Slug, SourceType: kind, Title: title},
		Triage:       out.Triage, TriageReason: out.Reason, DuplicateOf: out.DuplicateOf, Shadow: out.Shadow,
		Refetched: out.Refetched,
	}
}

// Tally returns the run's counts so far.
func (r *Run) Tally() gate.Tally {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tally
}

// Summary is what a finished run reports.
type Summary struct {
	Tally gate.Tally `json:"tally"`
	// Modes names the gate mode per gate.
	Modes []string `json:"modes"`
	// Classify and Link are the reports of the in-run classification and
	// linking of the written sources (nil when nothing was written).
	Classify *classify.Report `json:"classify,omitempty"`
	Link     *link.Report     `json:"link,omitempty"`
	CostUSD  float64          `json:"cost_usd"`
	// Folders are the source counts per raw/ folder after the run, set when
	// the run quarantined a source (spec §7.1: the agent updates prose
	// counts from them).
	Folders []scope.FolderCount `json:"folders,omitempty"`
}

// Lines renders the summary for the terminal (the caller appends the
// session's decisions summary).
func (sum Summary) Lines() []string {
	lines := append([]string{}, sum.Modes...)
	lines = append(lines, "ingest: "+sum.Tally.Line())
	if sum.Folders != nil {
		lines = append(lines, scope.FolderCountsLine(sum.Folders))
	}
	if sum.Classify != nil {
		lines = append(lines, sum.Classify.Lines()...)
	}
	if sum.Link != nil {
		lines = append(lines, sum.Link.Lines()...)
	}
	return lines
}

// ModeLines states the mode of every gate for the run summary.
func ModeLines(s *session.Session, force bool) []string {
	if force {
		return []string{"gates: skipped (--force); classify and link still run"}
	}
	relevance := s.RelevanceGateMode()
	if !s.RelevanceEnabled() {
		relevance = "off"
	}
	return []string{
		fmt.Sprintf("gate dedupe: %s", s.QualityGateMode()),
		fmt.Sprintf("gate quality: %s", s.QualityGateMode()),
		fmt.Sprintf("gate relevance: %s (%s)", relevance, s.GateModeReason()),
	}
}

// Finish classifies and links the sources written by the run (spec §7
// stage 7), then appends one log.md entry for the run with the
// kept/review/quarantined/skipped counts and the cost. The error is fatal
// only; the summary is returned with what was done.
func (r *Run) Finish(ctx context.Context) (Summary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	summary := Summary{Tally: r.tally, Modes: ModeLines(r.s, r.opts.Force)}
	var runErr error
	if len(r.written) > 0 {
		r.s.ReloadCorpus()
		classified, err := classify.Run(ctx, r.s, classify.Options{Paths: r.written})
		summary.Classify = &classified
		if err != nil {
			runErr = fmt.Errorf("ingest: classify: %w", err)
		} else {
			r.s.ReloadCorpus()
			linked, err := link.Run(ctx, r.s, link.Options{Paths: r.written})
			summary.Link = &linked
			if err != nil {
				runErr = fmt.Errorf("ingest: link: %w", err)
			}
		}
	}
	if r.tally.Quarantined > 0 {
		folders, err := scope.SourceFolderCounts(r.s.Root())
		if err != nil {
			runErr = errors.Join(runErr, fmt.Errorf("ingest: %w", err))
		} else {
			summary.Folders = folders
		}
	}
	decided := r.s.Summary()
	summary.CostUSD = decided.CostUSD + decided.Generation.CostUSD
	if len(r.entries) > 0 {
		if err := appendRunLogEntry(filepath.Join(r.s.Root(), "log.md"), r.s.Now().UTC(), r.opts, summary, decided.Calls, r.entries); err != nil {
			runErr = errors.Join(runErr, fmt.Errorf("ingest: append log entry: %w", err))
		}
		r.entries = nil
	}
	r.written = nil
	return summary, runErr
}

// appendRunLogEntry writes the run's log.md entry:
//
//	## [2026-09-24] ingest | url-2026-09-24-a1b2c3 (url)
//
//	kept 1, review 1, quarantined 1, skipped 0, duplicates 0 · decisions US$ 0.0004 (3 calls)
//
//	- kept: `Title` → `topic/raw/articles/x.md`
//	- quarantined (error_page): `Title` → `topic/raw/_quarantine/articles/y.md` — HTTP status 404
func appendRunLogEntry(logPath string, when time.Time, opts RunOptions, summary Summary, calls int, entries []logEntry) error {
	command := firstNonEmpty(opts.Command, "ingest")
	label := firstNonEmpty(opts.Batch, command)
	lines := []string{
		"",
		fmt.Sprintf("## [%s] ingest | %s (%s)", when.Format(frontmatter.DateLayout), label, command),
		"",
		fmt.Sprintf("%s · decisions US$ %.4f (%d calls)", summary.Tally.Line(), summary.CostUSD, calls),
		"",
	}
	for _, entry := range entries {
		status := entry.triage
		if entry.reason != "" && entry.triage != "kept" {
			status += " (" + entry.reason + ")"
		}
		line := fmt.Sprintf("- %s: `%s`", status, entry.title)
		if entry.path != "" {
			line += fmt.Sprintf(" → `%s`", entry.path)
		}
		if entry.detail != "" {
			line += " — " + entry.detail
		}
		lines = append(lines, line)
	}
	lines = append(lines, "")
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %q: %w", logPath, err)
	}
	if _, err := file.WriteString(strings.Join(lines, "\n")); err != nil {
		_ = file.Close()
		return fmt.Errorf("write log entry: %w", err)
	}
	return file.Close()
}

// builder assembles the would-be document of one source for the gates and
// the final write.
type builder struct {
	topicInfo models.TopicInfo
	options   Options
	rel       string
	now       time.Time
	// markdown is the body of the last built document (the refetched body
	// when stage 4 kept it).
	markdown string
}

// values returns the frontmatter of the source; withOwned adds the
// provenance keys (ingest_batch, ingest_query), which the final write adds
// through the owned-key writer instead.
func (b *builder) values(title, markdown string, withOwned bool) (map[string]any, error) {
	options := b.options
	if !withOwned {
		options.Batch, options.Query = "", ""
	}
	return buildFrontmatter(b.topicInfo, options, title, markdown, b.now)
}

// build returns the would-be document at the allocated path.
func (b *builder) build(title, markdown string) (*corpus.Document, error) {
	values, err := b.values(title, markdown, true)
	if err != nil {
		return nil, fmt.Errorf("ingest: %w", err)
	}
	content, err := frontmatter.Generate(values, markdown)
	if err != nil {
		return nil, fmt.Errorf("ingest: generate frontmatter: %w", err)
	}
	parsed, body, err := frontmatter.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("ingest: parse generated frontmatter: %w", err)
	}
	b.markdown = markdown
	return &corpus.Document{
		Path:        b.rel,
		AbsPath:     filepath.Join(b.topicInfo.RootPath, filepath.FromSlash(b.rel)),
		Title:       strings.TrimSpace(title),
		Kind:        corpus.KindSource,
		Frontmatter: maps.Clone(parsed),
		Raw:         content,
		Body:        body,
		BodyHash:    corpus.BodyHash(body),
		Aliases:     resolve.Aliases(parsed),
		ModTime:     b.now,
	}, nil
}

// uniqueGatedSlug allocates a slug free both in raw/<dir>/ and in
// raw/_quarantine/<dir>/, so a later quarantine never collides.
func uniqueGatedSlug(topicRoot, sourceDirectory, baseSlug string) (string, error) {
	live := filepath.Join(topicRoot, "raw", filepath.FromSlash(sourceDirectory))
	quarantine := filepath.Join(topicRoot, filepath.FromSlash(resolve.QuarantineDir), filepath.FromSlash(sourceDirectory))
	candidate := baseSlug
	for index := 2; ; index++ {
		free := true
		for _, dir := range []string{live, quarantine} {
			_, err := os.Stat(filepath.Join(dir, candidate+".md"))
			if err == nil {
				free = false
				break
			}
			if !errors.Is(err, os.ErrNotExist) {
				return "", fmt.Errorf("stat %q: %w", filepath.Join(dir, candidate+".md"), err)
			}
		}
		if free {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", baseSlug, index)
	}
}

// WouldBePath returns the topic-relative path a source of kind titled title
// would get before uniqueness suffixes (for collected_on_purpose_paths
// matching before fetch).
func WouldBePath(kind models.SourceKind, title string) string {
	directory, err := rawDirectoryForSourceKind(kind)
	if err != nil {
		return ""
	}
	slug := vault.SlugifySegment(strings.TrimSpace(title))
	if slug == "" {
		slug = "untitled-source"
	}
	return path.Join("raw", directory, slug+".md")
}

// writeExclusive creates path with content, failing when it exists.
func writeExclusive(path, content string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
