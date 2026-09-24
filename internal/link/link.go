// Package link relates a topic's documents to its wiki articles (spec §9):
// code proposes candidates (mentions, BM25, shared structure, optional
// semantic neighbours), the decision model judges one request per document,
// and code writes frontmatter relations and body links under write policy
// (d). kb only adds links; it never removes one.
package link

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/resolve"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

// State and bank bookkeeping (plan "Bank ids recorded in state rows").
const (
	// BankID is the question bank asked by link.
	BankID = "link"
	// CandidatesBank is the state-row bank key holding the candidate hash.
	CandidatesBank = "link_candidates"
	// InsertedLinksFile is the log of every body link kb inserted, under
	// <topic>/.decisions/.
	InsertedLinksFile = "inserted-links.jsonl"
	// excerptTokens bounds the document excerpt in the judgment state.
	excerptTokens = 6000
	// maxMentionsPerCandidate and maxMentions bound mention_sense questions.
	maxMentionsPerCandidate = 3
	maxMentions             = 40
	// maxListedDocuments caps the undecided documents named in Lines.
	maxListedDocuments = 10
)

// Skip reasons reported in Report.Skipped.
const (
	SkipUnchanged    = "unchanged"
	SkipNoCandidates = "no-candidates"
	SkipLocked       = "locked"
	SkipChanged      = "changed"
	SkipUserKey      = "user-key"
	SkipQuarantined  = "quarantined"
	SkipNotFound     = "not-found"
)

// Options configures Run.
type Options struct {
	// Paths limits the run to these topic-relative documents (ingest); empty
	// means every document of the topic.
	Paths []string
	// All relinks documents whose body and candidate set are unchanged.
	All bool
	// DryRun asks (cached) decisions but writes nothing: no files, state
	// rows, queue items or inserted-links rows.
	DryRun bool
	// Neighbours is the optional semantic-neighbour source (nil = off).
	Neighbours Neighbours
}

// Report summarizes a link run.
type Report struct {
	DryRun    bool
	BodyMode  string
	Documents int
	Judged    int
	// Applied counts new targets per relation key (and `affects`).
	Applied        map[string]int
	Inserted       int
	Proposals      int
	Reviews        int
	Demotions      int
	Contradictions int
	// Undecided counts undecided answers.
	Undecided int
	// Decided counts judged documents whose every answer was decided.
	Decided int
	// UndecidedDocuments lists, sorted, the judged documents with an
	// undecided answer: they keep no link bookkeeping, so the next run
	// judges them again (through the cache for the decided answers).
	UndecidedDocuments []string
	Skipped            map[string]int
	// Changes lists what was (or, in a dry run, would be) changed.
	Changes []string
}

// Lines renders the report for the terminal.
func (r Report) Lines() []string {
	prefix := "link"
	if r.DryRun {
		prefix = "link (dry run)"
	}
	lines := []string{
		fmt.Sprintf("%s: %d documents, %d judged, body links %s", prefix, r.Documents, r.Judged, r.BodyMode),
		fmt.Sprintf("coverage: %d/%d judged documents fully decided", r.Decided, r.Judged),
	}
	if r.DryRun {
		lines = append(lines, "dry run: decisions are still asked (cached answers are free, new ones spend the run budget); nothing was written")
	}
	applied := make([]string, 0, len(r.Applied))
	for _, key := range slices.Sorted(maps.Keys(r.Applied)) {
		if key != KeyAffects {
			applied = append(applied, fmt.Sprintf("%s %d", key, r.Applied[key]))
		}
	}
	if len(applied) == 0 {
		applied = append(applied, "none")
	}
	lines = append(lines, "relations: "+strings.Join(applied, ", "))
	if r.BodyMode == session.ModeApply {
		lines = append(lines, fmt.Sprintf("body links: %d inserted", r.Inserted))
	} else {
		lines = append(lines, fmt.Sprintf("body links: %d proposed (shadow; set decisions.mode: apply to insert)", r.Proposals))
	}
	lines = append(lines,
		fmt.Sprintf("affects: %d", r.Applied[KeyAffects]),
		fmt.Sprintf("review: %d link, %d demotions, %d contradictions", r.Reviews, r.Demotions, r.Contradictions),
		fmt.Sprintf("undecided: %d", r.Undecided),
	)
	if len(r.UndecidedDocuments) > 0 {
		lines = append(lines, fmt.Sprintf("undecided documents (%d, judged again next run): %s", len(r.UndecidedDocuments), corpus.ListPaths(r.UndecidedDocuments, maxListedDocuments)))
	}
	if len(r.Skipped) > 0 {
		skipped := make([]string, 0, len(r.Skipped))
		for _, reason := range slices.Sorted(maps.Keys(r.Skipped)) {
			skipped = append(skipped, fmt.Sprintf("skipped:%s %d", reason, r.Skipped[reason]))
		}
		lines = append(lines, strings.Join(skipped, ", "))
	}
	if r.DryRun {
		for _, change := range r.Changes {
			lines = append(lines, "  would "+change)
		}
	}
	return lines
}

func (r *Report) skip(reason string) {
	if r.Skipped == nil {
		r.Skipped = map[string]int{}
	}
	r.Skipped[reason]++
}

// runner holds what every document of one run shares.
type runner struct {
	s         *session.Session
	bank      *questions.Bank
	targets   *Targets
	generator *Generator
	queue     *review.Store
	th        Thresholds
	bodyMode  string
	dryRun    bool
	// exists reports whether a topic-relative path is a document of the
	// corpus (rename detection in state lookups).
	exists func(path string) bool

	logMu sync.Mutex
}

// job is one document to judge: the plan plus the candidate hash.
type job struct {
	plan Plan
	hash string
	// meta is the state-row metadata recorded with the write.
	meta corpus.StateMeta
}

// docResult is the per-document outcome, merged into the report in order.
type docResult struct {
	path    string
	judged  bool
	skip    string
	outcome Outcome
	changes []string
	applied map[string]int
	userKey bool
}

// Run links the topic's documents (spec §9.4): new or changed documents by
// default, every document with All, only Paths when set.
func Run(ctx context.Context, s *session.Session, opts Options) (Report, error) {
	r, docs, err := newRunner(s, opts.Neighbours, opts.DryRun)
	if err != nil {
		return Report{}, err
	}
	report := r.newReport()

	selected := docs
	if len(opts.Paths) > 0 {
		c, err := s.Corpus()
		if err != nil {
			return Report{}, fmt.Errorf("link: %w", err)
		}
		selected = make([]*corpus.Document, 0, len(opts.Paths))
		for _, p := range opts.Paths {
			doc := c.ByPath(p)
			if doc == nil {
				report.skip(SkipNotFound)
				continue
			}
			selected = append(selected, doc)
		}
	}

	jobs := make([]*job, 0, len(selected))
	for _, doc := range selected {
		if reason := skipReason(doc); reason != "" {
			report.Documents++
			report.skip(reason)
			continue
		}
		report.Documents++
		candidates := r.generator.Candidates(ctx, doc)
		if len(candidates) == 0 {
			report.skip(SkipNoCandidates)
			continue
		}
		hash := CandidateHash(candidates)
		if !opts.All && r.unchanged(doc, hash) {
			report.skip(SkipUnchanged)
			continue
		}
		meta := s.StateMeta(r.bank)
		meta.Banks[CandidatesBank] = hash
		jobs = append(jobs, &job{plan: r.plan(doc, candidates, r.generator.Mentions(doc, candidates)), hash: hash, meta: meta})
	}

	return r.execute(ctx, jobs, report)
}

// Reverse is the backlink audit of spec §9.4: every document mentioning the
// article's title or aliases is judged with that article as its only
// candidate. It does not touch the incremental link bookkeeping.
func Reverse(ctx context.Context, s *session.Session, articlePath string) (Report, error) {
	r, docs, err := newRunner(s, nil, false)
	if err != nil {
		return Report{}, err
	}
	c, err := s.Corpus()
	if err != nil {
		return Report{}, fmt.Errorf("link: %w", err)
	}
	article := c.ByPath(articlePath)
	if article == nil || article.Kind != corpus.KindArticle {
		return Report{}, fmt.Errorf("link: reverse: %q is not an article of topic %q", articlePath, s.Topic.Slug)
	}
	single := corpus.NewDictionary([]*corpus.Document{article})
	candidate := Candidate{
		ID: candidateID(0), Path: article.Path, Title: article.Title, Aliases: slices.Clone(article.Aliases),
		Criterion: article.Criterion(), Summary: article.Summary(), Origins: []string{OriginMention},
	}

	report := r.newReport()
	jobs := make([]*job, 0)
	for _, doc := range docs {
		if doc.Path == article.Path {
			continue
		}
		mentions := single.Mentions(doc.Body)
		if len(mentions) == 0 {
			continue
		}
		report.Documents++
		if reason := skipReason(doc); reason != "" {
			report.skip(reason)
			continue
		}
		meta := corpus.StateMeta{Contract: s.ContractHash(), Banks: map[string]string{}}
		jobs = append(jobs, &job{plan: r.plan(doc, []Candidate{candidate}, mentions), meta: meta})
	}
	return r.execute(ctx, jobs, report)
}

// ReverseAll runs the reverse pass for each article in turn (spec §9.4: an
// article was created or gained aliases) and merges the reports. The corpus
// is reloaded between articles so a document written for one article is read
// fresh for the next.
func ReverseAll(ctx context.Context, s *session.Session, articlePaths []string) (Report, error) {
	total := Report{BodyMode: s.BodyMode(), Applied: map[string]int{}, Skipped: map[string]int{}}
	for _, articlePath := range articlePaths {
		s.ReloadCorpus()
		report, err := Reverse(ctx, s, articlePath)
		total.add(report)
		if err != nil {
			return total, err
		}
	}
	s.ReloadCorpus()
	return total, nil
}

// add folds another report's counts into r.
func (r *Report) add(other Report) {
	r.Documents += other.Documents
	r.Judged += other.Judged
	r.Inserted += other.Inserted
	r.Proposals += other.Proposals
	r.Reviews += other.Reviews
	r.Demotions += other.Demotions
	r.Contradictions += other.Contradictions
	r.Undecided += other.Undecided
	r.Decided += other.Decided
	r.UndecidedDocuments = mergeSorted(r.UndecidedDocuments, other.UndecidedDocuments)
	for key, count := range other.Applied {
		r.Applied[key] += count
	}
	for key, count := range other.Skipped {
		r.Skipped[key] += count
	}
	r.Changes = append(r.Changes, other.Changes...)
}

func newRunner(s *session.Session, neighbours Neighbours, dryRun bool) (*runner, []*corpus.Document, error) {
	bank, err := questions.Load(BankID)
	if err != nil {
		return nil, nil, fmt.Errorf("link: %w", err)
	}
	c, err := s.Corpus()
	if err != nil {
		return nil, nil, fmt.Errorf("link: %w", err)
	}
	index, files, err := resolve.LoadVaultIndex(s.VaultPath, s.Root())
	if err != nil {
		return nil, nil, fmt.Errorf("link: %w", err)
	}
	prefix, err := filepath.Rel(filepath.Clean(s.VaultPath), filepath.Clean(s.Root()))
	if err != nil || strings.HasPrefix(prefix, "..") {
		prefix = filepath.Base(s.Root())
	}
	r := &runner{
		s:         s,
		bank:      bank,
		targets:   NewTargets(index, files, filepath.ToSlash(prefix)),
		generator: NewGenerator(c.Articles(), neighbours),
		queue:     review.Open(s.Root(), s.Now),
		th:        ThresholdsFrom(s),
		bodyMode:  s.BodyMode(),
		dryRun:    dryRun,
		exists:    func(p string) bool { return c.ByPath(p) != nil },
	}
	docs := make([]*corpus.Document, 0)
	for _, doc := range c.Documents() {
		if !IsHub(doc.Path) {
			docs = append(docs, doc)
		}
	}
	return r, docs, nil
}

func (r *runner) newReport() Report {
	return Report{DryRun: r.dryRun, BodyMode: r.bodyMode, Applied: map[string]int{}, Skipped: map[string]int{}}
}

func skipReason(doc *corpus.Document) string {
	switch {
	case doc.Locked():
		return SkipLocked
	case strings.EqualFold(doc.Triage(), "quarantined"):
		return SkipQuarantined
	default:
		return ""
	}
}

// unchanged reports a document whose link judgment is current: the link
// bank was judged on its current body (recorded per bank, so writes by other
// commands never make it look current) with the same bank version and
// candidate set. An undecided judgment records blank versions, so it is
// never unchanged.
func (r *runner) unchanged(doc *corpus.Document, hash string) bool {
	row, ok := r.lookup(doc)
	if !ok {
		return false
	}
	return row.JudgedBody(BankID) == doc.BodyHash && row.Banks[BankID] == r.bank.Version && row.Banks[CandidatesBank] == hash
}

// lookup returns the state row of doc, following a rename by body hash.
func (r *runner) lookup(doc *corpus.Document) (*corpus.StateRow, bool) {
	return r.s.State.Lookup(doc.Path, doc.BodyHash, r.exists)
}

// plan gathers, by code, what the write policy needs about doc.
func (r *runner) plan(doc *corpus.Document, candidates []Candidate, mentions []corpus.Mention) Plan {
	plan := Plan{
		Doc:        doc,
		Candidates: candidates,
		Forms:      make(map[string]string, len(candidates)),
		Relations:  map[string][]string{},
		Present:    map[string]map[string]bool{},
		KBOwned:    map[string]bool{},
		BodyLinked: map[string]bool{},
	}
	idByPath := make(map[string]string, len(candidates))
	for _, candidate := range candidates {
		plan.Forms[candidate.Path] = r.targets.Form(candidate.Path)
		idByPath[candidate.Path] = candidate.ID
	}

	row, _ := r.lookup(doc)
	for _, key := range append(slices.Clone(decisions.RelationKeys), KeyAffects) {
		entries := stringList(doc.Frontmatter, key)
		if len(entries) == 0 {
			continue
		}
		plan.Relations[key] = entries
		plan.Present[key] = map[string]bool{}
		for _, entry := range entries {
			if target := r.targets.Resolve(doc.Path, resolve.Normalize(entry)); target != "" {
				plan.Present[key][target] = true
			}
		}
		if row != nil && row.Written[key] != "" && row.Written[key] == corpus.ValueHash(doc.Frontmatter[key]) {
			plan.KBOwned[key] = true
		}
	}
	for _, link := range resolve.BodyLinks(doc.Body) {
		if target := r.targets.Resolve(doc.Path, link.Target); target != "" {
			plan.BodyLinked[target] = true
		}
	}

	perCandidate := map[string]int{}
	for _, mention := range mentions {
		id, ok := idByPath[mention.Article.Path]
		if !ok || plan.BodyLinked[mention.Article.Path] || perCandidate[id] >= maxMentionsPerCandidate || len(plan.Mentions) >= maxMentions {
			continue
		}
		perCandidate[id]++
		plan.Mentions = append(plan.Mentions, MentionRef{
			N: len(plan.Mentions) + 1, Candidate: id, Start: mention.Start, End: mention.End,
			Text: doc.Body[mention.Start:mention.End], Sentence: mention.Sentence,
		})
	}
	return plan
}

// execute judges and writes every job over the engine's concurrency.
func (r *runner) execute(ctx context.Context, jobs []*job, report Report) (Report, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]docResult, len(jobs))
	workers := max(1, r.s.Engine.Concurrency())
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	var errMu sync.Mutex
	var firstErr error

	for index, current := range jobs {
		if ctx.Err() != nil {
			break
		}
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			result, err := r.link(ctx, current)
			if err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				errMu.Unlock()
				cancel()
				return
			}
			results[index] = result
		}()
	}
	wg.Wait()
	if firstErr != nil {
		return report, firstErr
	}

	for _, result := range results {
		if result.skip != "" {
			report.skip(result.skip)
		}
		if result.userKey {
			report.skip(SkipUserKey)
		}
		if !result.judged {
			continue
		}
		report.Judged++
		report.Undecided += result.outcome.Undecided
		if result.outcome.Undecided > 0 {
			report.UndecidedDocuments = append(report.UndecidedDocuments, result.path)
		} else {
			report.Decided++
		}
		if result.skip != "" {
			continue
		}
		for key, count := range result.applied {
			report.Applied[key] += count
		}
		report.Inserted += len(result.outcome.Insertions)
		report.Proposals += result.outcome.Proposals
		report.Reviews += result.outcome.Reviews
		report.Demotions += result.outcome.Demotions
		report.Contradictions += result.outcome.Contradictions
		report.Changes = append(report.Changes, result.changes...)
	}
	sort.Strings(report.UndecidedDocuments)
	return report, nil
}

// link judges one document and applies the write policy.
func (r *runner) link(ctx context.Context, current *job) (docResult, error) {
	plan := current.plan
	doc := plan.Doc
	request, err := r.request(plan)
	if err != nil {
		return docResult{}, err
	}
	result, err := r.s.Engine.Decide(ctx, request)
	if err != nil {
		return docResult{}, fmt.Errorf("link %s: %w", doc.Path, err)
	}
	outcome := Decide(plan, result.Answers, r.th, r.bodyMode)
	res := docResult{path: doc.Path, judged: true, outcome: outcome, changes: describe(plan, outcome), applied: map[string]int{}}
	meta := current.meta
	if outcome.Undecided > 0 {
		// An undecided answer (timeout, budget, invalid receipt) is never a
		// "no": the document keeps no link bookkeeping, so the next
		// incremental run judges it again, through the cache for the answers
		// that were decided (spec §2.3, §4.1).
		meta.Banks = maps.Clone(meta.Banks)
		for _, id := range []string{BankID, CandidatesBank} {
			if _, ok := meta.Banks[id]; ok {
				meta.Banks[id] = ""
			}
		}
	}

	if r.dryRun {
		for key, paths := range outcome.Added {
			res.applied[key] = len(paths)
		}
		return res, nil
	}

	var write corpus.WriteResult
	if len(outcome.Insertions) > 0 {
		write, err = r.s.Writer.ApplyBody(doc, InsertLinks(doc.Body, outcome.Insertions), outcome.Updates, meta)
	} else {
		write, err = r.s.Writer.Apply(doc, outcome.Updates, meta)
	}
	if err != nil {
		return docResult{}, fmt.Errorf("link %s: %w", doc.Path, err)
	}
	switch write.Status {
	case corpus.StatusSkippedLocked:
		res.skip = SkipLocked
		return res, nil
	case corpus.StatusSkippedChanged:
		res.skip = SkipChanged
		return res, nil
	}
	res.userKey = len(write.SkippedUserKeys) > 0
	for key, paths := range outcome.Added {
		if slices.Contains(write.Written, key) {
			res.applied[key] = len(paths)
		}
	}
	if !write.BodyWritten {
		res.outcome.Insertions = nil
	} else if err := r.logInsertions(doc.Path, outcome.Insertions); err != nil {
		return docResult{}, err
	}
	for _, item := range outcome.Items {
		if _, err := r.queue.Add(item); err != nil {
			return docResult{}, fmt.Errorf("link %s: queue: %w", doc.Path, err)
		}
	}
	return res, nil
}

// request builds the one judgment request of a document (spec §9.2).
func (r *runner) request(plan Plan) (decisions.Request, error) {
	doc := plan.Doc
	terms := make([]string, 0, len(plan.Candidates))
	candidates := make([]map[string]any, 0, len(plan.Candidates))
	qs := make([]questions.Q, 0, len(plan.Candidates)*3+len(plan.Mentions))
	titleByID := make(map[string]string, len(plan.Candidates))
	for _, candidate := range plan.Candidates {
		criterion := candidate.Criterion
		if criterion == "" {
			criterion = candidate.Summary
		}
		entry := map[string]any{"id": candidate.ID, "title": candidate.Title, "criterion": criterion}
		if len(candidate.Aliases) > 0 {
			entry["aliases"] = candidate.Aliases
		}
		candidates = append(candidates, entry)
		terms = append(terms, candidate.Title)
		titleByID[candidate.ID] = candidate.Title

		ids := []string{"should_link_{id}", "relation_{id}"}
		if doc.Kind == corpus.KindSource {
			ids = append(ids, "affects_{id}")
		}
		for _, id := range ids {
			q, err := r.bank.Question(id, map[string]string{"id": candidate.ID})
			if err != nil {
				return decisions.Request{}, fmt.Errorf("link: %w", err)
			}
			qs = append(qs, q)
		}
	}
	mentions := make([]map[string]any, 0, len(plan.Mentions))
	for _, mention := range plan.Mentions {
		mentions = append(mentions, map[string]any{
			"n": mention.N, "id": mention.Candidate, "candidate": titleByID[mention.Candidate],
			"text": mention.Text, "sentence": mention.Sentence,
		})
		q, err := r.bank.Question("mention_sense_{n}", map[string]string{"n": strconv.Itoa(mention.N)})
		if err != nil {
			return decisions.Request{}, fmt.Errorf("link: %w", err)
		}
		qs = append(qs, q)
	}

	document := corpus.DocumentState(doc, excerptTokens, terms)
	if summary := doc.Summary(); summary != "" {
		document["summary"] = summary
	}
	state := map[string]any{"document": document, "candidates": candidates}
	if len(mentions) > 0 {
		state["mentions"] = mentions
	}
	return r.s.Request(decisions.PurposeLink, doc.Path, r.bank, state, qs), nil
}

// insertedLink is one row of .decisions/inserted-links.jsonl.
type insertedLink struct {
	Time    string `json:"time"`
	Subject string `json:"subject"`
	Target  string `json:"target"`
	Text    string `json:"text"`
	Mode    string `json:"mode"`
}

func (r *runner) logInsertions(subject string, insertions []Insertion) error {
	if len(insertions) == 0 {
		return nil
	}
	r.logMu.Lock()
	defer r.logMu.Unlock()

	dir := filepath.Join(r.s.Root(), decisions.ReceiptsDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("link: inserted links: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(dir, InsertedLinksFile), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("link: inserted links: %w", err)
	}
	now := r.s.Now().UTC().Format(time.RFC3339)
	var writeErr error
	for _, insertion := range insertions {
		line, err := json.Marshal(insertedLink{Time: now, Subject: subject, Target: insertion.Target, Text: insertion.Text, Mode: r.bodyMode})
		if err != nil {
			writeErr = err
			break
		}
		if _, err := file.Write(append(line, '\n')); err != nil {
			writeErr = err
			break
		}
	}
	return errors.Join(writeErr, file.Close())
}

// describe renders the changes an outcome makes, for dry runs.
func describe(plan Plan, outcome Outcome) []string {
	changes := make([]string, 0)
	for _, key := range slices.Sorted(maps.Keys(outcome.Added)) {
		for _, p := range outcome.Added[key] {
			changes = append(changes, fmt.Sprintf("add %s [[%s]] to %s", key, plan.Forms[p], plan.Doc.Path))
		}
	}
	for _, insertion := range outcome.Insertions {
		changes = append(changes, fmt.Sprintf("insert [[%s|%s]] in %s at byte %d", insertion.Form, insertion.Text, plan.Doc.Path, insertion.Start))
	}
	items := slices.Clone(outcome.Items)
	sort.SliceStable(items, func(i, j int) bool { return items[i].Target < items[j].Target })
	for _, item := range items {
		changes = append(changes, fmt.Sprintf("queue %s item: %s", item.Queue, item.Evidence))
	}
	return changes
}

// mergeSorted returns the sorted union of two path lists.
func mergeSorted(a, b []string) []string {
	merged := append(slices.Clone(a), b...)
	sort.Strings(merged)
	return slices.Compact(merged)
}
