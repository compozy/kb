package classify

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

// Proposal limits (spec §5.2 "New concepts").
const (
	// MinNoneForProposals is the number of sources with primary concept
	// `none` that triggers concept proposals.
	MinNoneForProposals = 5
	// MaxProposals caps the concepts proposed in one run.
	MaxProposals = 5
	// maxProposalSources caps the source lines given to the proposal call.
	maxProposalSources = 60
	// enoughTextThreshold is the P(enough_text_to_judge) under which depth
	// is not written (spec §8).
	enoughTextThreshold = 0.5
)

// State-row fact recording the primary concept outcome of a source, so the
// concept proposal trigger counts across the topic, not only this run.
const (
	// FactPrimaryConcept is the fact name.
	FactPrimaryConcept = "primary_concept"
	// PrimaryNone records a decided primary concept `none`.
	PrimaryNone = "none"
	// PrimaryMatched records a decided primary concept other than `none`.
	PrimaryMatched = "concept"
)

// Options configures Run.
type Options struct {
	// OnlyMissing judges only documents without a state row or with a
	// missing facet key (the bootstrap follow-up after vocabulary accept).
	OnlyMissing bool
	// All re-judges every document regardless of the state record.
	All bool
	// Paths restricts the run to these topic-relative paths (ingest).
	Paths []string
}

// Run classifies the documents of the session's topic (spec §8). Per
// document it asks two requests: the shared gate judgment (JudgeGate: role +
// quality nouls; sources only) and one classification request (composite
// classify + concept bank: kind, depth, enough_text_to_judge,
// primary_concept, mentions_concept_{id}), plus one generation call for
// summary/entities/questions when they are missing or kb wrote them and the
// body changed. Articles first get their criterion and aliases generated, so
// the concept options of the same run use them. Code maps answers to
// frontmatter and writes through the session's writer; undecided answers
// never write and leave the document unclassified for the next run. Sources
// above the review band go to the recapture or remove review queues (never
// quarantined), and ≥5 sources with primary concept `none` produce concept
// proposals. Incremental by body hash, contract and bank versions unless
// opts.All. The error is fatal only (auth, cancellation, I/O); the report
// is returned alongside it with what was done.
func Run(ctx context.Context, s *session.Session, opts Options) (Report, error) {
	c, err := s.Corpus()
	if err != nil {
		return newReport(), fmt.Errorf("classify: load corpus: %w", err)
	}
	if missing := missingPaths(c, opts.Paths); len(missing) > 0 {
		s.ReloadCorpus()
		if c, err = s.Corpus(); err != nil {
			return newReport(), fmt.Errorf("classify: load corpus: %w", err)
		}
	}

	r := newRun(s, c, opts)
	docs := r.selected()
	r.report.Documents = len(docs)

	articles := make([]*corpus.Document, 0)
	for _, doc := range docs {
		if doc.Kind == corpus.KindArticle && r.needsArticleLiterals(doc) {
			articles = append(articles, doc)
		}
	}
	if err := forEach(ctx, s.Engine.Concurrency(), articles, r.articleLiterals); err != nil {
		return r.finish(), err
	}

	r.vocab = newVocabulary(c.Articles(), r.linkTarget)
	r.report.NoVocabulary = r.vocab.empty()
	r.report.CriterionMissing = r.vocab.criterionMissing()

	work := make([]*corpus.Document, 0, len(docs))
	for _, doc := range docs {
		if reason := r.skipReason(doc); reason != "" {
			r.report.Skipped[reason]++
			continue
		}
		work = append(work, doc)
	}
	if err := forEach(ctx, s.Engine.Concurrency(), work, r.classifyDoc); err != nil {
		return r.finish(), err
	}

	if err := r.proposeConcepts(ctx); err != nil {
		return r.finish(), err
	}
	return r.finish(), nil
}

// run is the state of one classify run.
type run struct {
	s      *session.Session
	corpus *corpus.Corpus
	opts   Options
	// rows is the state record as it was when the run started: writes during
	// the run must not hide a body change from later decisions.
	rows       map[string]*corpus.StateRow
	vocab      *vocabulary
	hostLines  map[string]map[string]int
	queue      *review.Store
	thresholds decisions.Thresholds
	stems      map[string]int
	vaultRel   string

	mu     sync.Mutex
	report Report
	none   []noneSource
	// primaryAsked lists the sources whose primary concept was asked in
	// this run: their fact in r.rows is superseded.
	primaryAsked map[string]bool
	// undecidedDocs collects the documents left with an undecided answer
	// or literal.
	undecidedDocs map[string]bool
	// taken maps every lower-cased article title and alias to its article
	// path, for alias collision checks.
	taken map[string]string
}

// noneSource is a source whose primary concept was `none`.
type noneSource struct {
	path, title, summary string
}

func newRun(s *session.Session, c *corpus.Corpus, opts Options) *run {
	r := &run{
		s:          s,
		corpus:     c,
		opts:       opts,
		rows:       map[string]*corpus.StateRow{},
		hostLines:  quality.HostLineCounts(c.Sources()),
		queue:      review.Open(s.Root(), s.Now),
		thresholds: s.Ref.Thresholds,
		stems:      map[string]int{},
		report:     newReport(),
		taken:      map[string]string{},

		primaryAsked:  map[string]bool{},
		undecidedDocs: map[string]bool{},
	}
	r.report.RelevanceOff = !s.RelevanceEnabled()
	exists := func(p string) bool { return c.ByPath(p) != nil }
	for _, doc := range c.Documents() {
		if row, ok := s.State.Lookup(doc.Path, doc.BodyHash, exists); ok {
			r.rows[doc.Path] = row
		}
		r.stems[strings.ToLower(stemOf(doc.Path))]++
	}
	for _, article := range c.Articles() {
		for _, name := range append([]string{article.Title}, article.Aliases...) {
			if key := strings.ToLower(strings.TrimSpace(name)); key != "" {
				if _, set := r.taken[key]; !set {
					r.taken[key] = article.Path
				}
			}
		}
	}
	if s.VaultPath != "" {
		if rel, err := filepath.Rel(s.VaultPath, s.Root()); err == nil {
			r.vaultRel = filepath.ToSlash(rel)
		}
	}
	return r
}

func (r *run) finish() Report {
	r.mu.Lock()
	defer r.mu.Unlock()
	sort.Strings(r.report.ArticlesChanged)
	r.report.UndecidedDocuments = append([]string{}, slices.Sorted(maps.Keys(r.undecidedDocs))...)
	return r.report
}

func (r *run) tally(fn func(report *Report)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fn(&r.report)
}

// selected returns the documents of the run (opts.Paths or every document).
func (r *run) selected() []*corpus.Document {
	if len(r.opts.Paths) == 0 {
		return r.corpus.Documents()
	}
	docs := make([]*corpus.Document, 0, len(r.opts.Paths))
	seen := map[string]bool{}
	for _, p := range r.opts.Paths {
		doc := r.corpus.ByPath(p)
		if doc == nil {
			r.report.Skipped["missing"]++
			continue
		}
		if !seen[doc.Path] {
			seen[doc.Path] = true
			docs = append(docs, doc)
		}
	}
	return docs
}

func missingPaths(c *corpus.Corpus, paths []string) []string {
	var missing []string
	for _, p := range paths {
		if c.ByPath(p) == nil {
			missing = append(missing, p)
		}
	}
	return missing
}

// banksFor returns the banks recorded in the state row of doc.
func banksFor(doc *corpus.Document) []*questions.Bank {
	if doc.Kind == corpus.KindArticle {
		return ArticleBanks()
	}
	return SourceBanks()
}

func bankVersions(banks []*questions.Bank) map[string]string {
	versions := make(map[string]string, len(banks))
	for _, bank := range banks {
		versions[bank.ID] = bank.Version
	}
	return versions
}

// skipReason decides whether doc goes through the facet pass: "" means yes.
func (r *run) skipReason(doc *corpus.Document) string {
	switch {
	case doc.Locked():
		return "locked"
	case doc.Kind == corpus.KindArticle && isStub(doc):
		return "stub"
	case r.opts.All:
		return ""
	}
	row := r.rows[doc.Path]
	if r.opts.OnlyMissing {
		if row == nil || r.missingFacet(doc) {
			return ""
		}
		return "complete"
	}
	contractHash := r.s.ContractHash()
	if doc.Kind == corpus.KindArticle {
		// Article states carry no contract, so a contract change does not
		// re-judge them.
		contractHash = ""
	}
	if corpus.Unclassified(doc, row, contractHash, bankVersions(banksFor(doc))) {
		return ""
	}
	return "unchanged"
}

// missingFacet reports a document lacking one of the keys classify writes.
func (r *run) missingFacet(doc *corpus.Document) bool {
	keys := []string{"genre", "summary", "entities"}
	if doc.Kind == corpus.KindSource {
		keys = append(keys, "questions")
		if r.s.RelevanceEnabled() {
			keys = append(keys, "relevance")
		}
		if !r.vocab.empty() {
			keys = append(keys, "concepts")
		}
	} else {
		keys = append(keys, "criterion")
		if len(r.vocab.concepts) > 1 {
			keys = append(keys, "concepts")
		}
	}
	return slices.ContainsFunc(keys, func(key string) bool { return !hasValue(doc, key) })
}

// hasValue reports a present, non-empty frontmatter key (an empty list
// counts as present: it records a judgment with no result).
func hasValue(doc *corpus.Document, key string) bool {
	value, ok := doc.Frontmatter[key]
	if !ok || value == nil {
		return false
	}
	if text, isString := value.(string); isString {
		return strings.TrimSpace(text) != ""
	}
	return true
}

// needsLiteral reports whether kb should (re)generate the literal key of
// doc: it is missing, or kb wrote it (value hash matches the state record)
// from another body than the current one (the per-key body hash of the
// state row, which other writes never advance). A value the user owns is
// never regenerated.
func needsLiteral(doc *corpus.Document, row *corpus.StateRow, key string) bool {
	if !hasValue(doc, key) {
		return true
	}
	if row == nil || row.Written[key] != corpus.ValueHash(doc.Frontmatter[key]) {
		return false
	}
	return row.KeyBody(key) != doc.BodyHash
}

// ownsKey reports a present key whose value kb wrote and nobody edited
// since (its value hash equals the state record's written hash).
func ownsKey(doc *corpus.Document, row *corpus.StateRow, key string) bool {
	value, ok := doc.Frontmatter[key]
	return ok && value != nil && row != nil && row.Written[key] != "" && row.Written[key] == corpus.ValueHash(value)
}

func hasWords(doc *corpus.Document) bool {
	return strings.ContainsFunc(doc.Body, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) })
}

// needsArticleLiterals reports an article that lacks a criterion (or has a
// kb-written one and a changed body) or has no aliases.
func (r *run) needsArticleLiterals(article *corpus.Document) bool {
	if article.Locked() {
		return false
	}
	if r.needsCriterion(article) {
		return true
	}
	return len(article.Aliases) == 0
}

func (r *run) needsCriterion(article *corpus.Document) bool {
	return !isStub(article) && hasWords(article) && needsLiteral(article, r.rows[article.Path], "criterion")
}

// articleLiterals generates the criterion and aliases of one article and
// writes them (phase 1, before the vocabulary is built).
func (r *run) articleLiterals(ctx context.Context, article *corpus.Document) error {
	updates := map[string]any{}
	if r.needsCriterion(article) {
		criterion, err := generateCriterion(ctx, r.s, article)
		switch reason, soft := softGenerationError(err); {
		case soft:
			r.tally(func(report *Report) {
				report.Undecided["criterion:"+reason]++
				if errors.Is(err, errInvalidLiteral) {
					report.InvalidLiterals["criterion"]++
				}
			})
			r.markUndecided(article.Path)
		case err != nil:
			return err
		default:
			updates["criterion"] = criterion
		}
	}
	if len(article.Aliases) == 0 {
		proposed, err := generateAliases(ctx, r.s, article)
		switch reason, soft := softGenerationError(err); {
		case soft:
			r.tally(func(report *Report) { report.Undecided["aliases:"+reason]++ })
			r.markUndecided(article.Path)
		case err != nil:
			return err
		default:
			if kept := r.claimAliases(article, proposed); len(kept) > 0 {
				updates["aliases"] = kept
			}
		}
	}
	if len(updates) == 0 {
		return nil
	}
	result, err := r.s.Writer.Apply(article, updates, r.s.StateMeta())
	if err != nil {
		return fmt.Errorf("classify: write %s: %w", article.Path, err)
	}
	r.tally(func(report *Report) {
		recordWrite(report, result)
		if slices.Contains(result.Written, "aliases") {
			report.ArticlesChanged = append(report.ArticlesChanged, article.Path)
		}
	})
	return nil
}

// claimAliases validates proposed aliases against every other article's
// title and aliases (including aliases claimed earlier in this run) and
// reserves the kept ones for article.
func (r *run) claimAliases(article *corpus.Document, proposed []string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	kept := filterAliases(proposed, article.Title, func(key string) bool {
		owner, ok := r.taken[key]
		return ok && owner != article.Path
	})
	for _, alias := range kept {
		r.taken[strings.ToLower(alias)] = article.Path
	}
	return kept
}

// classifyDoc runs the facet pass for one document (phase 2).
func (r *run) classifyDoc(ctx context.Context, doc *corpus.Document) error {
	row := r.rows[doc.Path]
	source := doc.Kind == corpus.KindSource

	var gate *GateJudgment
	if source {
		judgment, err := JudgeGate(ctx, r.s, doc, GateOptions{
			IsTranscript: IsTranscriptKind(doc.SourceKind()),
			Flags:        r.codeFlags(doc),
		})
		if err != nil {
			return err
		}
		gate = &judgment
	}

	asked, err := r.askFacets(ctx, doc)
	if err != nil {
		return err
	}
	out := mapFacets(facetInput{
		source:       source,
		relevanceOn:  r.s.RelevanceEnabled(),
		hasRelevance: hasValue(doc, "relevance"),
		ownsDepth:    ownsKey(doc, row, "depth"),
		gate:         gate,
		answers:      asked.answers,
		options:      asked.options,
		candidates:   asked.candidates,
		thresholds:   r.thresholds,
	})

	dropped, invalid := 0, []string(nil)
	if r.needsDocLiterals(doc, row) {
		lit, err := generateLiterals(ctx, r.s, doc)
		switch reason, soft := softGenerationError(err); {
		case soft:
			out.undecided = append(out.undecided, "literals:"+reason)
			out.complete = false
		case err != nil:
			return err
		default:
			if lit.summary != "" {
				out.updates["summary"] = lit.summary
			}
			if len(lit.entities) > 0 {
				out.updates["entities"] = lit.entities
			}
			if source && len(lit.questions) > 0 {
				out.updates["questions"] = lit.questions
			}
			dropped, invalid = lit.droppedEntities, lit.invalid
		}
	}

	meta := r.s.StateMeta(banksFor(doc)...)
	if !out.complete {
		// An undecided answer leaves the document unclassified: blank bank
		// versions make the next run judge it again (through the cache for
		// the answers that were decided).
		for id := range meta.Banks {
			meta.Banks[id] = ""
		}
		r.markUndecided(doc.Path)
	}
	if source && out.primaryAsked {
		fact := ""
		switch {
		case out.primaryNone:
			fact = PrimaryNone
		case out.primaryDecided:
			fact = PrimaryMatched
		}
		meta.Facts = map[string]string{FactPrimaryConcept: fact}
		r.mu.Lock()
		r.primaryAsked[doc.Path] = true
		r.mu.Unlock()
	}
	result, err := r.s.Writer.Apply(doc, out.updates, meta)
	if err != nil {
		return fmt.Errorf("classify: write %s: %w", doc.Path, err)
	}

	r.tally(func(report *Report) {
		report.Judged++
		if out.complete {
			report.Decided++
		}
		recordWrite(report, result)
		for _, reason := range out.undecided {
			report.Undecided[reason]++
		}
		report.DroppedEntities += dropped
		for _, key := range invalid {
			report.InvalidLiterals[key]++
		}
		if source && out.primaryNone {
			report.PrimaryNone++
		}
	})
	if source && out.primaryNone {
		summary := doc.Summary()
		if generated, ok := out.updates["summary"].(string); ok && summary == "" {
			summary = generated
		}
		r.mu.Lock()
		r.none = append(r.none, noneSource{path: doc.Path, title: doc.Title, summary: summary})
		r.mu.Unlock()
	}
	if source {
		return r.enqueue(doc, gate, out)
	}
	return nil
}

// needsDocLiterals reports a document whose summary/entities/questions must
// be generated.
func (r *run) needsDocLiterals(doc *corpus.Document, row *corpus.StateRow) bool {
	if !hasWords(doc) {
		return false
	}
	keys := []string{"summary", "entities"}
	if doc.Kind == corpus.KindSource {
		keys = append(keys, "questions")
	}
	return slices.ContainsFunc(keys, func(key string) bool { return needsLiteral(doc, row, key) })
}

func (r *run) codeFlags(doc *corpus.Document) []quality.Flag {
	return CodeFlags(doc, r.hostLines)
}

// CodeFlags runs the code quality checks (spec §7 stage 3) over a stored
// source; hostLines comes from quality.HostLineCounts over the topic's
// sources. The thin rule only applies to web captures (an http(s) source_url,
// not a transcript or a bookmark cluster): the length of a local note or a
// transcript is not a capture failure. `kb classify` and the impact preview
// share it so their bands agree.
func CodeFlags(doc *corpus.Document, hostLines map[string]map[string]int) []quality.Flag {
	host, _ := doc.Provenance()["source_host"].(string)
	return quality.Check(quality.Input{
		Title:     doc.Title,
		Body:      doc.Body,
		SourceURL: doc.SourceURL(),
		HostLines: hostLines[host],
		SkipThin:  !webCapture(doc),
	})
}

func webCapture(doc *corpus.Document) bool {
	url := strings.ToLower(doc.SourceURL())
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	kind := doc.SourceKind()
	return !IsTranscriptKind(kind) && kind != string(models.SourceKindBookmarkCluster)
}

// askedFacets is the classification request and its answers.
type askedFacets struct {
	answers    map[string]decisions.Answer
	options    []*concept
	candidates []*concept
}

// askFacets asks the classification request of doc: kind (all), depth and
// enough_text_to_judge (sources), and, when the topic has other articles,
// primary_concept over the concept options plus mentions_concept_{id} for
// the code-built candidates.
func (r *run) askFacets(ctx context.Context, doc *corpus.Document) (askedFacets, error) {
	asked := askedFacets{}
	qs := []questions.Q{facetBank.MustQuestion("kind", nil)}
	var terms []string
	state := map[string]any{}
	if doc.Kind == corpus.KindSource {
		qs = append(qs, facetBank.MustQuestion("depth", nil), facetBank.MustQuestion("enough_text_to_judge", nil))
		state["contract"] = contractState(r.s.Contract)
	}
	if !r.vocab.empty() {
		asked.options = r.vocab.options(doc)
	}
	if len(asked.options) > 0 {
		criteria := make(map[string]string, len(asked.options))
		for _, option := range asked.options {
			criteria[option.id] = option.optionText()
		}
		qs = append(qs, facetBank.WithCriteria(facetBank.MustQuestion("primary_concept", nil), criteria))
		asked.candidates = r.vocab.candidates(doc)
		entries := make([]map[string]any, 0, len(asked.candidates))
		for _, candidate := range asked.candidates {
			qs = append(qs, facetBank.MustQuestion("mentions_concept_{id}", map[string]string{"id": candidate.id}))
			entries = append(entries, candidate.stateEntry())
			terms = append(terms, candidate.title)
			terms = append(terms, candidate.aliases...)
		}
		state["concepts"] = entries
	}
	state["document"] = corpus.DocumentState(doc, ExcerptTokens, terms)

	candidateIDs := make([]string, 0, len(asked.candidates))
	for _, candidate := range asked.candidates {
		candidateIDs = append(candidateIDs, candidate.id)
	}
	bank, qs, err := r.s.WithExtras(decisions.PurposeClassify, facetBank, qs, candidateIDs...)
	if err != nil {
		return asked, fmt.Errorf("classify: facets for %s: %w", doc.Path, err)
	}
	result, err := r.s.Engine.Decide(ctx, r.s.Request(decisions.PurposeClassify, doc.Path, bank, state, qs))
	if err != nil {
		return asked, fmt.Errorf("classify: facets for %s: %w", doc.Path, err)
	}
	asked.answers = result.Answers
	return asked, nil
}

// facetInput is everything mapFacets reads.
type facetInput struct {
	source       bool
	relevanceOn  bool
	hasRelevance bool
	// ownsDepth reports a `depth` kb wrote and nobody edited since.
	ownsDepth  bool
	gate       *GateJudgment
	answers    map[string]decisions.Answer
	options    []*concept
	candidates []*concept
	thresholds decisions.Thresholds
}

// queueDecision is a review item code decided to add.
type queueDecision struct {
	reason, question string
	p                float64
}

// facetOutcome is what code does with the answers of one document.
type facetOutcome struct {
	updates     map[string]any
	undecided   []string
	complete    bool
	primaryNone bool
	// primaryAsked and primaryDecided report the primary_concept question.
	primaryAsked   bool
	primaryDecided bool
	recapture      *queueDecision
	remove         *queueDecision
}

// mapFacets maps answers to frontmatter updates and queue decisions (spec
// §8, §12.1). Undecided answers never write; the only exception is
// `relevance: unknown` on a source that has no relevance yet, which records
// that the role could not be judged (never `off_topic`).
func mapFacets(in facetInput) facetOutcome {
	out := facetOutcome{updates: map[string]any{}, complete: true}
	undecided := func(id string, answer decisions.Answer) {
		out.undecided = append(out.undecided, id+":"+answerReason(answer))
		out.complete = false
	}
	th := in.thresholds

	if in.source && in.gate != nil {
		gate := in.gate
		if len(gate.Undecided) > 0 {
			out.undecided = append(out.undecided, gate.Undecided...)
			out.complete = false
		}
		if in.relevanceOn {
			switch {
			case gate.RoleDecided:
				out.updates["relevance"] = gate.Role
			case !in.hasRelevance:
				out.updates["relevance"] = RoleUnknown
			}
		}
		if reason, _ := gate.QualityReason(th.Get("quality_apply")); reason != "" {
			out.updates["quality"] = reason
		}
		if reason, p := gate.QualityReason(th.Get("quality_review")); reason != "" {
			out.recapture = &queueDecision{reason: reason, question: gate.QualityQuestion(reason), p: p}
		} else if gate.RoleAsked && gate.RoleDecided && gate.QualityDecided && gate.POffTopic >= th.Get("relevance_review") {
			out.remove = &queueDecision{reason: RoleOffTopic, question: "role", p: gate.POffTopic}
		}
	}

	if kind := in.answers["kind"]; kind.Decided() {
		out.updates["genre"] = kind.Choice
	} else {
		undecided("kind", kind)
	}

	if in.source {
		enough := in.answers["enough_text_to_judge"]
		depth := in.answers["depth"]
		p, decided := enough.P("")
		switch {
		case !decided:
			undecided("enough_text_to_judge", enough)
		case p < enoughTextThreshold:
			// Too little text to judge depth: a depth kb wrote from an
			// earlier body is stale and goes (a user's depth stays).
			if in.ownsDepth {
				out.updates["depth"] = nil
			}
		case !depth.Decided() || depth.Score == nil:
			undecided("depth", depth)
		default:
			out.updates["depth"] = depthValue(*depth.Score)
		}
	}

	if len(in.options) > 0 {
		out.mapConcepts(in, undecided)
	}
	return out
}

// depthValue rounds the expected depth level to one decimal; whole levels
// are written as integers so the frontmatter reads `depth: 2`.
func depthValue(score float64) any {
	rounded := math.Round(score*10) / 10
	if rounded == math.Trunc(rounded) {
		return int(rounded)
	}
	return rounded
}

// mapConcepts writes `concepts` = the primary concept when P ≥
// primary_concept and ≠ none, plus every candidate whose mentions noul ≥
// concept_noul; only when every concept answer was decided.
func (out *facetOutcome) mapConcepts(in facetInput, undecided func(string, decisions.Answer)) {
	byID := make(map[string]*concept, len(in.options))
	for _, option := range in.options {
		byID[option.id] = option
	}
	links := make([]string, 0, 4)
	add := func(c *concept) {
		if c != nil && !slices.Contains(links, c.link) {
			links = append(links, c.link)
		}
	}
	complete := true
	primary := in.answers["primary_concept"]
	out.primaryAsked = true
	out.primaryDecided = primary.Decided()
	switch {
	case !primary.Decided():
		undecided("primary_concept", primary)
		complete = false
	case primary.Choice == "none":
		out.primaryNone = true
	case primary.Probs[primary.Choice] >= in.thresholds.Get("primary_concept"):
		add(byID[primary.Choice])
	}
	for _, candidate := range in.candidates {
		id := "mentions_concept_" + candidate.id
		answer := in.answers[id]
		p, ok := answer.P("")
		if !ok {
			undecided(id, answer)
			complete = false
			continue
		}
		if p >= in.thresholds.Get("concept_noul") {
			add(candidate)
		}
	}
	if complete {
		out.updates["concepts"] = links
	}
}

// enqueue adds the backfill review items of a source (spec §12.1).
func (r *run) enqueue(doc *corpus.Document, gate *GateJudgment, out facetOutcome) error {
	var items []review.Item
	if d := out.recapture; d != nil {
		action := map[string]any{"reason": d.reason}
		if strings.TrimSpace(doc.SourceURL()) == "" {
			action["no_source"] = true
		}
		evidence := doc.Title + " — " + d.reason
		if index := slices.IndexFunc(gate.Flags, func(flag quality.Flag) bool { return flag.Code == d.reason }); index >= 0 {
			evidence += ": " + gate.Flags[index].Detail
		}
		items = append(items, review.Item{
			Queue: review.QueueRecapture, Purpose: string(decisions.PurposeQuality), Subject: doc.Path,
			Question: d.question, Probability: d.p, ReceiptKey: gate.ReceiptKey, Evidence: evidence, Action: action,
		})
	}
	if d := out.remove; d != nil {
		items = append(items, review.Item{
			Queue: review.QueueRemove, Purpose: string(decisions.PurposeRelevance), Subject: doc.Path,
			Question: d.question, Probability: d.p, ReceiptKey: gate.ReceiptKey,
			Evidence: fmt.Sprintf("%s — P(off_topic) %.2f", doc.Title, d.p), Action: map[string]any{"reason": d.reason},
		})
	}
	return r.addItems(items)
}

func (r *run) addItems(items []review.Item) error {
	for _, item := range items {
		added, err := r.queue.Add(item)
		if err != nil {
			return fmt.Errorf("classify: review queue: %w", err)
		}
		if added {
			queue := item.Queue
			r.tally(func(report *Report) { report.Queued[queue]++ })
		}
	}
	return nil
}

// markUndecided records a document left with an undecided answer or
// literal, for the run summary.
func (r *run) markUndecided(path string) {
	r.mu.Lock()
	r.undecidedDocs[path] = true
	r.mu.Unlock()
}

// noneSources returns every source of the topic whose primary concept is
// `none`: the ones judged in this run plus, for the sources this run did not
// ask, the primary-concept fact recorded in their state row by earlier runs.
// Sorted by path.
func (r *run) noneSources() []noneSource {
	r.mu.Lock()
	none := slices.Clone(r.none)
	asked := maps.Clone(r.primaryAsked)
	r.mu.Unlock()
	for _, doc := range r.corpus.Sources() {
		if asked[doc.Path] || doc.Locked() || strings.EqualFold(doc.Triage(), "quarantined") {
			continue
		}
		if row := r.rows[doc.Path]; row != nil && row.Facts[FactPrimaryConcept] == PrimaryNone {
			none = append(none, noneSource{path: doc.Path, title: doc.Title, summary: doc.Summary()})
		}
	}
	sort.Slice(none, func(i, j int) bool { return none[i].path < none[j].path })
	return none
}

// proposeConcepts asks, when ≥ MinNoneForProposals sources of the topic
// have primary concept `none` (this run's judgments plus the facts earlier
// runs recorded), one generation call for up to MaxProposals new concepts
// from those sources' summaries, and queues them as concept-proposal items.
// It never creates articles.
func (r *run) proposeConcepts(ctx context.Context) error {
	if r.vocab.empty() {
		return nil
	}
	none := r.noneSources()
	r.tally(func(report *Report) { report.PrimaryNoneTotal = len(none) })
	if len(none) < MinNoneForProposals {
		return nil
	}
	lines := make([]string, 0, maxProposalSources)
	for _, source := range none[:min(len(none), maxProposalSources)] {
		line := "- " + source.title
		if source.summary != "" {
			line += ": " + source.summary
		}
		lines = append(lines, line)
	}
	existing := make([]string, 0, len(r.vocab.concepts))
	for _, c := range r.vocab.concepts {
		existing = append(existing, c.title)
	}
	prompt := strings.Join([]string{
		fmt.Sprintf("The sources below match none of the existing concepts of the knowledge-base topic. Propose at most %d new concepts that would cover them; each concept becomes a wiki article.", MaxProposals),
		"- title: the concept's name as a reader would look it up, 1 to 8 words, no dates; never one of the existing concepts.",
		"- criterion: " + criterionRules,
		"",
		"Topic: " + r.s.Topic.Title,
		contractLines(r.s.Contract),
		"",
		"Existing concepts: " + strings.Join(existing, "; "),
		"",
		"Sources (title: summary):",
		strings.Join(lines, "\n"),
	}, "\n")
	proposed, err := generateConcepts(ctx, r.s, KindProposals, r.s.Topic.Slug, prompt)
	if reason, soft := softGenerationError(err); soft {
		r.tally(func(report *Report) { report.Undecided["concept_proposals:"+reason]++ })
		return nil
	} else if err != nil {
		return err
	}

	taken := map[string]bool{}
	r.mu.Lock()
	for key := range r.taken {
		taken[key] = true
	}
	r.mu.Unlock()
	items := make([]review.Item, 0, MaxProposals)
	for _, concept := range validConcepts(proposed, taken, MaxProposals) {
		items = append(items, review.Item{
			Queue:    review.QueueConceptProposal,
			Purpose:  string(decisions.PurposeConcept),
			Subject:  path.Join(ConceptsDir, sanitizeFileName(concept.Title)+".md"),
			Target:   concept.Title,
			Question: "primary_concept",
			Evidence: fmt.Sprintf("proposed from %d sources whose primary concept is none", len(none)),
			Action:   map[string]any{"title": concept.Title, "criterion": concept.Criterion},
		})
	}
	return r.addItems(items)
}

// linkTarget renders the `concepts` value of an article.
func (r *run) linkTarget(article *corpus.Document) string {
	return linkStem(article, func(stem string) bool { return r.stems[strings.ToLower(stem)] > 1 }, r.vaultRel)
}

func stemOf(p string) string {
	return strings.TrimSuffix(path.Base(p), path.Ext(p))
}

// recordWrite adds one writer result to the report.
func recordWrite(report *Report, result corpus.WriteResult) {
	report.Writes[result.Status]++
	for _, key := range result.Written {
		report.Facets[key]++
	}
	for _, key := range result.SkippedUserKeys {
		report.SkippedUserKeys[key]++
	}
}

// forEach runs fn over items on up to workers goroutines. The first error
// cancels the remaining work and is returned; a canceled parent context is
// returned as its error.
func forEach(ctx context.Context, workers int, items []*corpus.Document, fn func(context.Context, *corpus.Document) error) error {
	if len(items) == 0 {
		return ctx.Err()
	}
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	workers = max(1, min(workers, len(items)))

	jobs := make(chan *corpus.Document)
	var (
		wg       sync.WaitGroup
		once     sync.Once
		firstErr error
	)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				if ctx.Err() != nil {
					continue
				}
				if err := fn(ctx, item); err != nil {
					once.Do(func() {
						firstErr = err
						cancel(err)
					})
				}
			}
		}()
	}
feed:
	for _, item := range items {
		select {
		case jobs <- item:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}
	return context.Cause(ctx)
}
