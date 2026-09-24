// Package actions performs the human verdicts of the review queue (spec
// §12.1): `kb review accept|reject` dispatches each item to the action of
// its queue (keep or quarantine a gated source, rescue a skipped item,
// recapture a broken page, write a relation or a body link, create a stub
// concept, record an OKF type verdict), appends the verdict as a label
// (spec §12.2) and resolves the item. Nothing leaves raw/ without an accept.
package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/classify"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/firecrawl"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/gate"
	"github.com/compozy/kb/internal/link"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/quality"
	"github.com/compozy/kb/internal/refs"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

// Verdicts.
const (
	Accept = "accept"
	Reject = "reject"
)

// What an action did (Result.Action).
const (
	DidRestored    = "restored"
	DidKept        = "kept"
	DidQuarantined = "quarantined"
	DidRescued     = "rescued"
	DidDismissed   = "dismissed"
	DidRecaptured  = "recaptured"
	DidWritten     = "written"
	DidDemoted     = "demoted"
	DidCreated     = "created"
	DidRecorded    = "recorded"
	DidNothing     = "nothing"
)

// ErrNotPending is returned for an item that is already resolved.
var ErrNotPending = errors.New("actions: item is not pending")

// Deps are the fetchers an action may need; both are optional and an action
// that needs a missing one fails without changing anything.
type Deps struct {
	// Scrape fetches a URL with explicit Firecrawl options (recapture).
	Scrape func(ctx context.Context, sourceURL string, opts firecrawl.ScrapeOptions) (*firecrawl.ScrapeResult, error)
	// Rescue fetches and ingests a skipped item, bypassing only the
	// pre-fetch gate, and returns the written path (skip accept).
	Rescue func(ctx context.Context, row gate.SkippedRow) (string, error)
}

// Result reports one applied verdict.
type Result struct {
	ID      string `json:"id"`
	Queue   string `json:"queue"`
	Verdict string `json:"verdict"`
	Subject string `json:"subject"`
	// Action is what was done (restored, kept, quarantined, rescued, ...).
	Action string `json:"action"`
	// Path is the document path after the action (the quarantine path for a
	// quarantined source).
	Path string `json:"path,omitempty"`
	// Message is one line for the terminal.
	Message string       `json:"message,omitempty"`
	Label   review.Label `json:"label"`
	// Closed lists sibling items (same source, other gate/remove/recapture
	// queue) resolved by the same action.
	Closed []string `json:"closed,omitempty"`
}

// LabelVerdict maps a verdict on a queue item to its label (spec §12.2
// convention: positive = keep / link, negative = remove / broken): gate,
// skip, link, contradiction, concept-proposal and okf-type accepts are
// positive; remove and recapture accepts (the source is off topic or
// broken) are negative, and so is accepting a link demotion. Rejects are
// the opposite.
func LabelVerdict(item review.Item, verdict string) string {
	accepted := verdict == Accept
	negativeOnAccept := false
	switch item.Queue {
	case review.QueueRemove, review.QueueRecapture:
		negativeOnAccept = true
	case review.QueueLink:
		negativeOnAccept = actionBool(item, "demote")
	}
	if accepted != negativeOnAccept {
		return review.VerdictPositive
	}
	return review.VerdictNegative
}

// Apply performs verdict (Accept or Reject) on one pending item, appends its
// label and resolves it. When the action fails nothing is labelled and the
// item stays pending.
func Apply(ctx context.Context, s *session.Session, deps Deps, item review.Item, verdict string) (Result, error) {
	if verdict != Accept && verdict != Reject {
		return Result{}, fmt.Errorf("actions: verdict must be accept or reject: %q", verdict)
	}
	if item.Status != "" && item.Status != review.StatusPending {
		return Result{}, fmt.Errorf("%w: %s is %s", ErrNotPending, item.ID, item.Status)
	}
	a := &applier{s: s, deps: deps, item: item, verdict: verdict}
	result := Result{ID: item.ID, Queue: item.Queue, Verdict: verdict, Subject: item.Subject, Action: DidNothing}
	var err error
	switch item.Queue {
	case review.QueueGate:
		err = a.gate(ctx, &result)
	case review.QueueSkip:
		err = a.skip(ctx, &result)
	case review.QueueRecapture:
		err = a.recapture(ctx, &result)
	case review.QueueRemove:
		err = a.remove(ctx, &result)
	case review.QueueLink, review.QueueContradiction:
		err = a.link(&result)
	case review.QueueConceptProposal:
		err = a.concept(ctx, &result)
	case review.QueueOKFType:
		a.okfType(&result)
	default:
		err = fmt.Errorf("actions: unknown queue %q", item.Queue)
	}
	if err != nil {
		return result, fmt.Errorf("actions: %s %s: %w", verdict, item.ID, err)
	}

	store := review.Open(s.Root(), s.Now)
	result.Label = review.Label{
		Subject: item.Subject, Purpose: item.Purpose, Question: item.Question, Target: item.Target,
		Verdict: LabelVerdict(item, verdict), ReceiptKey: item.ReceiptKey, Origin: review.OriginReview, ItemID: item.ID,
	}
	if err := store.AddLabel(result.Label); err != nil {
		return result, fmt.Errorf("actions: label %s: %w", item.ID, err)
	}
	status := review.StatusRejected
	if verdict == Accept {
		status = review.StatusAccepted
	}
	if _, err := store.Resolve(item.ID, status); err != nil {
		return result, fmt.Errorf("actions: resolve %s: %w", item.ID, err)
	}
	closed, err := closeSiblings(store, item, result.Action)
	result.Closed = closed
	if err != nil {
		return result, fmt.Errorf("actions: close siblings of %s: %w", item.ID, err)
	}
	return result, nil
}

// closeSiblings resolves the pending gate/remove/recapture items of the
// same source once an action settled it: a quarantine settles all three
// (gate rejected, remove and recapture accepted); keeping a source from
// the gate queue rejects its remove and recapture items. No labels are
// added for siblings.
func closeSiblings(store *review.Store, item review.Item, action string) ([]string, error) {
	sourceQueues := []string{review.QueueGate, review.QueueRemove, review.QueueRecapture}
	if !slices.Contains(sourceQueues, item.Queue) {
		return nil, nil
	}
	quarantined := action == DidQuarantined
	kept := item.Queue == review.QueueGate && (action == DidKept || action == DidRestored)
	if !quarantined && !kept {
		return nil, nil
	}
	pending, err := store.Items("", review.StatusPending)
	if err != nil {
		return nil, err
	}
	closed := make([]string, 0)
	for _, sibling := range pending {
		if sibling.ID == item.ID || sibling.Subject != item.Subject || !slices.Contains(sourceQueues, sibling.Queue) {
			continue
		}
		accepted := sibling.Queue != review.QueueGate
		if kept {
			accepted = !accepted
		}
		status := review.StatusRejected
		if accepted {
			status = review.StatusAccepted
		}
		if _, err := store.Resolve(sibling.ID, status); err != nil {
			return closed, err
		}
		closed = append(closed, sibling.ID)
	}
	return closed, nil
}

// applier holds one verdict being applied.
type applier struct {
	s       *session.Session
	deps    Deps
	item    review.Item
	verdict string
}

func (a *applier) accepted() bool { return a.verdict == Accept }

// gate: accept keeps the source (restoring it from quarantine), reject
// quarantines it with the item's reason.
func (a *applier) gate(ctx context.Context, result *Result) error {
	if a.accepted() {
		return a.keep(ctx, result)
	}
	return a.quarantine(ctx, result, a.reason(gate.ReasonOffTopic))
}

// skip: accept rescues (fetch and ingest the skipped item), reject
// dismisses it.
func (a *applier) skip(ctx context.Context, result *Result) error {
	if !a.accepted() {
		result.Action, result.Message = DidDismissed, "skipped item dismissed"
		return nil
	}
	id := actionString(a.item, "skipped_id")
	if id == "" {
		id = a.item.Subject
	}
	row, err := gate.FindSkipped(a.s.Root(), id)
	if err != nil {
		return err
	}
	if a.deps.Rescue == nil {
		return errors.New("rescue is not available")
	}
	written, err := a.deps.Rescue(ctx, row)
	if err != nil {
		return fmt.Errorf("rescue %s: %w", id, err)
	}
	result.Action, result.Path = DidRescued, written
	result.Message = "rescued " + row.URL
	if written != "" {
		result.Message += " → " + written
	}
	return nil
}

// remove: accept quarantines the source as off_topic, reject keeps it.
func (a *applier) remove(ctx context.Context, result *Result) error {
	if !a.accepted() {
		result.Action, result.Message = DidKept, "kept"
		return nil
	}
	return a.quarantine(ctx, result, a.reason(gate.ReasonOffTopic))
}

// reason returns the item's Action reason or fallback.
func (a *applier) reason(fallback string) string {
	if reason := actionString(a.item, "reason"); reason != "" {
		return reason
	}
	return fallback
}

// keep writes `triage: kept` (dropping kb's triage_reason), restoring the
// source from quarantine first when it is there.
func (a *applier) keep(ctx context.Context, result *Result) error {
	s := a.s
	if active, ok, err := activeQuarantine(s.Root(), a.item.Subject); err != nil {
		return err
	} else if ok {
		restored, err := refs.Restore(ctx, refs.Options{
			VaultPath: s.VaultPath, TopicRoot: s.Root(), ID: active.ID, Writer: s.Writer, State: s.State, Now: s.Now,
		})
		if err != nil {
			return err
		}
		s.ReloadCorpus()
		result.Action, result.Path = DidRestored, restored.Path
		result.Message = fmt.Sprintf("restored %s from %s (%d references replayed)", restored.Path, active.Path, max(len(restored.Restored)-1, 0))
		if len(restored.Manual) > 0 {
			result.Message += fmt.Sprintf("; %d references changed since and need manual repair", len(restored.Manual))
		}
		return nil
	}
	doc, err := corpus.ReadDocument(s.Root(), a.item.Subject, corpus.KindSource)
	if err != nil {
		return err
	}
	if _, err := s.Writer.Apply(doc, map[string]any{"triage": gate.TriageKept, "triage_reason": nil}, s.StateMeta()); err != nil {
		return err
	}
	result.Action, result.Path, result.Message = DidKept, a.item.Subject, "kept "+a.item.Subject
	return nil
}

// quarantine moves the source to raw/_quarantine/ through the refs ledger
// (a source already there is left alone).
func (a *applier) quarantine(ctx context.Context, result *Result, reason string) error {
	s := a.s
	if active, ok, err := activeQuarantine(s.Root(), a.item.Subject); err != nil {
		return err
	} else if ok {
		result.Action, result.Path = DidQuarantined, active.Path
		result.Message = "already quarantined at " + active.Path
		return nil
	}
	moved, err := refs.Quarantine(ctx, refs.Options{
		VaultPath: s.VaultPath, TopicRoot: s.Root(), Path: a.item.Subject, Reason: reason,
		Writer: s.Writer, State: s.State, Now: s.Now,
	})
	if err != nil {
		return err
	}
	s.ReloadCorpus()
	result.Action, result.Path = DidQuarantined, moved.NewPath
	result.Message = fmt.Sprintf("quarantined %s (%s) → %s", a.item.Subject, reason, moved.NewPath)
	return nil
}

// activeQuarantine finds the active quarantine of a topic-relative path
// (original or quarantine path).
func activeQuarantine(topicRoot, subject string) (refs.Quarantined, bool, error) {
	active, err := refs.List(topicRoot)
	if err != nil {
		return refs.Quarantined{}, false, err
	}
	for _, entry := range slices.Backward(active) {
		if entry.Original == subject || entry.Path == subject {
			return entry, true, nil
		}
	}
	return refs.Quarantined{}, false, nil
}

// recapture: accept refetches with fresh options and re-judges the new
// body; a page that is now good gets its body replaced in place
// (`recaptured: <date>`, kb's `quality` and `triage_reason` dropped); a
// page still failing, or a source without source_url, is quarantined with
// its quality reason. Reject keeps the source as it is.
func (a *applier) recapture(ctx context.Context, result *Result) error {
	if !a.accepted() {
		result.Action, result.Message = DidKept, "kept as captured"
		return nil
	}
	s := a.s
	doc, err := corpus.ReadDocument(s.Root(), a.item.Subject, corpus.KindSource)
	if err != nil {
		return err
	}
	reason := a.reason(quality.Thin)
	sourceURL := strings.TrimSpace(doc.SourceURL())
	if actionBool(a.item, "no_source") || !strings.HasPrefix(strings.ToLower(sourceURL), "http") {
		return a.quarantine(ctx, result, reason)
	}
	if a.deps.Scrape == nil {
		return errors.New("recapture needs Firecrawl")
	}
	fresh, err := a.deps.Scrape(ctx, sourceURL, firecrawl.RefetchOptions(s.Config.Firecrawl))
	if err != nil {
		return fmt.Errorf("refetch %s: %w", sourceURL, err)
	}
	if reason, failing, err := a.stillFailing(ctx, doc, fresh); err != nil {
		return err
	} else if failing {
		return a.quarantine(ctx, result, reason)
	}
	updates := map[string]any{
		"recaptured": s.Now().UTC().Format(frontmatter.DateLayout), "quality": nil, "triage_reason": nil,
	}
	written, err := s.Writer.ApplyBody(doc, fresh.Markdown, updates, s.StateMeta())
	if err != nil {
		return err
	}
	if !written.BodyWritten {
		return fmt.Errorf("%s was not rewritten (%s)", a.item.Subject, written.Status)
	}
	s.ReloadCorpus()
	result.Action, result.Path = DidRecaptured, a.item.Subject
	result.Message = fmt.Sprintf("recaptured %s (%d → %d characters)", a.item.Subject, len(doc.Body), len(fresh.Markdown))
	return nil
}

// stillFailing re-judges a refetched body: code flags first (no call), then
// the quality nouls of the shared gate judgment at the review threshold.
func (a *applier) stillFailing(ctx context.Context, doc *corpus.Document, fresh *firecrawl.ScrapeResult) (string, bool, error) {
	s := a.s
	candidate := *doc
	candidate.Body = fresh.Markdown
	candidate.BodyHash = corpus.BodyHash(fresh.Markdown)
	c, err := s.Corpus()
	if err != nil {
		return "", false, err
	}
	host, _ := doc.Provenance()["source_host"].(string)
	transcript := classify.IsTranscriptKind(doc.SourceKind())
	flags := quality.Check(quality.Input{
		Title: doc.Title, Body: fresh.Markdown, SourceURL: doc.SourceURL(), FinalURL: fresh.FinalURL,
		RequestedURL: doc.SourceURL(), StatusCode: fresh.StatusCode,
		HostLines: quality.HostLineCounts(c.Sources())[host],
		SkipThin:  transcript || doc.SourceKind() == string(models.SourceKindBookmarkCluster),
	})
	if len(flags) > 0 {
		return flags[0].Code, true, nil
	}
	judgment, err := classify.JudgeGate(ctx, s, &candidate, classify.GateOptions{IsTranscript: transcript})
	if err != nil {
		return "", false, err
	}
	if reason, _ := judgment.QualityReason(s.Threshold("quality_review")); reason != "" {
		return reason, true, nil
	}
	return "", false, nil
}

// link writes a relation (link and contradiction queues). Accept adds
// `[[target]]` to the relation list (plus, for `insert` items, the body
// link at the proposed mention, logged in inserted-links.jsonl); accepting
// a demotion removes the kb-written value. Reject does nothing.
func (a *applier) link(result *Result) error {
	if !a.accepted() {
		result.Message = "dismissed"
		return nil
	}
	s := a.s
	relation := firstNonEmpty(actionString(a.item, "relation"), link.RelationRelated)
	if a.item.Queue == review.QueueContradiction {
		relation = link.RelationContradicts
	}
	if !slices.Contains(decisions.RelationKeys, relation) && relation != link.KeyAffects {
		return fmt.Errorf("unknown relation %q", relation)
	}
	target := actionString(a.item, "target")
	if target == "" {
		return errors.New("item has no target")
	}
	doc, err := corpus.ReadDocument(s.Root(), a.item.Subject, documentKind(a.item.Subject))
	if err != nil {
		return err
	}
	entries := frontmatterList(doc.Frontmatter, relation)

	if actionBool(a.item, "demote") {
		return a.demote(doc, relation, target, entries, result)
	}

	updates := map[string]any{}
	if !slices.ContainsFunc(entries, func(entry string) bool { return sameTarget(entry, target) }) {
		updates[relation] = append(slices.Clone(entries), "[["+target+"]]")
	}
	body, inserted := doc.Body, ""
	if actionBool(a.item, "insert") {
		body, inserted = insertLink(doc.Body, target, actionString(a.item, "mention_text"), actionInt(a.item, "mention_start"))
	}
	var write corpus.WriteResult
	if inserted != "" {
		write, err = s.Writer.ApplyBody(doc, body, updates, s.StateMeta())
	} else {
		write, err = s.Writer.Apply(doc, updates, s.StateMeta())
	}
	if err != nil {
		return err
	}
	if write.Status == corpus.StatusSkippedLocked || write.Status == corpus.StatusSkippedChanged {
		return fmt.Errorf("%s not written (%s)", a.item.Subject, write.Status)
	}
	result.Action, result.Path = DidWritten, a.item.Subject
	result.Message = fmt.Sprintf("%s [[%s]] on %s", relation, target, a.item.Subject)
	if slices.Contains(write.SkippedUserKeys, relation) {
		result.Message += fmt.Sprintf(" (the %s key belongs to the user; left unchanged)", relation)
	}
	if inserted != "" && write.BodyWritten {
		if err := logInsertion(s, a.item.Subject, a.item.Target, inserted); err != nil {
			return err
		}
		result.Message += fmt.Sprintf("; body link [[%s|%s]] inserted", target, inserted)
	} else if actionBool(a.item, "insert") {
		result.Message += "; the mention is no longer in the body, no body link inserted"
	}
	return nil
}

// demote removes a kb-written relation value; a list the user edited is
// left alone.
func (a *applier) demote(doc *corpus.Document, relation, target string, entries []string, result *Result) error {
	s := a.s
	row, ok := s.State.Get(doc.Path)
	if !ok || row.Written[relation] == "" || row.Written[relation] != corpus.ValueHash(doc.Frontmatter[relation]) {
		result.Message = fmt.Sprintf("%s on %s was edited by the user; left unchanged", relation, doc.Path)
		return nil
	}
	kept := slices.DeleteFunc(slices.Clone(entries), func(entry string) bool { return sameTarget(entry, target) })
	var value any = kept
	if len(kept) == 0 {
		value = nil
	}
	if _, err := s.Writer.Apply(doc, map[string]any{relation: value}, s.StateMeta()); err != nil {
		return err
	}
	result.Action, result.Path = DidDemoted, doc.Path
	result.Message = fmt.Sprintf("removed %s [[%s]] from %s", relation, target, doc.Path)
	return nil
}

// concept: accept creates the stub article (criterion recorded as
// kb-written) and runs the reverse link pass for it.
func (a *applier) concept(ctx context.Context, result *Result) error {
	if !a.accepted() {
		result.Message = "proposal dismissed"
		return nil
	}
	s := a.s
	title := firstNonEmpty(actionString(a.item, "title"), strings.TrimSuffix(filepath.Base(a.item.Subject), ".md"))
	rel, err := classify.CreateStub(s, title, actionString(a.item, "criterion"), s.Now())
	if err != nil && !errors.Is(err, classify.ErrArticleExists) {
		return err
	}
	s.ReloadCorpus()
	report, err := link.Reverse(ctx, s, rel)
	if err != nil {
		return err
	}
	result.Action, result.Path = DidCreated, rel
	result.Message = fmt.Sprintf("created %s; reverse link pass judged %d documents", rel, report.Judged)
	return nil
}

// okfType records the verdict only.
func (a *applier) okfType(result *Result) {
	result.Action = DidRecorded
	suggested := firstNonEmpty(actionString(a.item, "type"), a.item.Target)
	result.Message = "verdict recorded; re-run `kb promote " + a.item.Subject + " --to <okf-topic> --type <type>`"
	if suggested != "" && a.accepted() {
		result.Message = "verdict recorded; re-run `kb promote " + a.item.Subject + " --to <okf-topic> --type " + suggested + "`"
	}
}

// insertLink inserts `[[target|text]]` at start when the body still reads
// text there (else at the first plain occurrence of text outside a link);
// it returns the body and the linked text ("" when nothing was inserted).
func insertLink(body, target, text string, start int) (string, string) {
	if text == "" {
		return body, ""
	}
	if start < 0 || start+len(text) > len(body) || body[start:start+len(text)] != text || insideLink(body, start) {
		start = -1
		for offset := 0; offset < len(body); {
			index := strings.Index(body[offset:], text)
			if index < 0 {
				break
			}
			if !insideLink(body, offset+index) {
				start = offset + index
				break
			}
			offset += index + len(text)
		}
	}
	if start < 0 {
		return body, ""
	}
	return body[:start] + "[[" + target + "|" + text + "]]" + body[start+len(text):], text
}

// insideLink reports whether offset sits inside `[[...]]`.
func insideLink(body string, offset int) bool {
	open := strings.LastIndex(body[:offset], "[[")
	if open < 0 {
		return false
	}
	return !strings.Contains(body[open:offset], "]]")
}

// insertedLink mirrors link's inserted-links.jsonl row.
type insertedLink struct {
	Time    string `json:"time"`
	Subject string `json:"subject"`
	Target  string `json:"target"`
	Text    string `json:"text"`
	Mode    string `json:"mode"`
}

var insertedMu sync.Mutex

// logInsertion appends a row to inserted-links.jsonl (mode "review": the
// link was inserted by an accepted review item).
func logInsertion(s *session.Session, subject, target, text string) error {
	insertedMu.Lock()
	defer insertedMu.Unlock()
	dir := filepath.Join(s.Root(), decisions.ReceiptsDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("inserted links: %w", err)
	}
	line, err := json.Marshal(insertedLink{Time: s.Now().UTC().Format(time.RFC3339), Subject: subject, Target: target, Text: text, Mode: "review"})
	if err != nil {
		return fmt.Errorf("inserted links: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(dir, link.InsertedLinksFile), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("inserted links: %w", err)
	}
	if _, err := file.Write(append(line, '\n')); err != nil {
		_ = file.Close()
		return fmt.Errorf("inserted links: %w", err)
	}
	return file.Close()
}

// sameTarget reports whether a relation entry (`[[x]]`, `[[x|alias]]` or a
// bare name) points to target.
func sameTarget(entry, target string) bool {
	name := strings.TrimSpace(entry)
	name = strings.TrimSuffix(strings.TrimPrefix(name, "[["), "]]")
	name, _, _ = strings.Cut(name, "|")
	name, _, _ = strings.Cut(name, "#")
	return strings.EqualFold(strings.TrimSuffix(strings.TrimSpace(name), ".md"), strings.TrimSuffix(strings.TrimSpace(target), ".md"))
}

func documentKind(rel string) corpus.Kind {
	if strings.HasPrefix(rel, "wiki/") {
		return corpus.KindArticle
	}
	return corpus.KindSource
}

func frontmatterList(values map[string]any, key string) []string {
	if single, ok := values[key].(string); ok {
		return []string{single}
	}
	return frontmatter.GetStringSlice(values, key)
}

func actionString(item review.Item, key string) string {
	value, _ := item.Action[key].(string)
	return strings.TrimSpace(value)
}

func actionBool(item review.Item, key string) bool {
	value, _ := item.Action[key].(bool)
	return value
}

func actionInt(item review.Item, key string) int {
	switch value := item.Action[key].(type) {
	case float64:
		return int(value)
	case int:
		return value
	case json.Number:
		n, _ := value.Int64()
		return int(n)
	}
	return -1
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
