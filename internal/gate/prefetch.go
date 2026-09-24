package gate

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

// MaxItemsPerRequest caps the items of one pre-fetch relevance request.
const MaxItemsPerRequest = 20

// Pre-fetch actions.
const (
	ActionFetch  = "fetch"
	ActionReview = "review"
	ActionSkip   = "skip"
)

// Pre-fetch reasons for a fetch that was not decided by the role answer.
const (
	PrefetchCollectedPath = "collected_on_purpose_path"
	PrefetchRelevanceOff  = "relevance_off"
	// PrefetchExcluded marks an item whose would-be path matches
	// decisions.exclude: it is fetched without a pre-fetch call.
	PrefetchExcluded = "excluded"
)

// relevanceBank holds the `role_item_{id}` question.
var relevanceBank = questions.MustLoad("relevance")

// Item is one bulk item before fetch: metadata only.
type Item struct {
	// ID is the stable skipped id; derived from URL with SkippedID when
	// empty.
	ID          string
	URL         string
	Title       string
	Description string
	// Channel is the channel, uploader or site name.
	Channel string
	// Date is the publication date as text.
	Date string
	// Path is the would-be topic-relative path (for collected_on_purpose_paths).
	Path string
	// SourceKind is the would-be source kind (provenance).
	SourceKind string
}

// PrefetchOptions carries the run provenance recorded with every item.
type PrefetchOptions struct {
	Batch string
	Query string
}

// PrefetchDecision is the pre-fetch verdict of one item.
type PrefetchDecision struct {
	Item Item `json:"-"`
	// Action is fetch, review (fetch and write with triage: review) or skip.
	Action string `json:"action"`
	// Reason explains a fetch without a role answer (collected path,
	// relevance off, undecided:<reason>).
	Reason     string  `json:"reason,omitempty"`
	Role       string  `json:"role,omitempty"`
	POffTopic  float64 `json:"p_off_topic"`
	PKept      float64 `json:"p_kept"`
	ReceiptKey string  `json:"receipt_key,omitempty"`
	// Shadow is the would-be action a shadow gate did not apply ("skip").
	Shadow string `json:"shadow,omitempty"`
}

// PrefetchAction bands one item: skip when P(off_topic) ≥ quarantine, fetch
// when P(core)+P(adjacent)+P(collected_on_purpose) ≥ fetch, else review.
func PrefetchAction(pOffTopic, pKept, quarantine, fetch float64) string {
	switch {
	case pOffTopic >= quarantine:
		return ActionSkip
	case pKept >= fetch:
		return ActionFetch
	default:
		return ActionReview
	}
}

// Prefetch runs stage 2 over bulk items: items whose would-be path matches
// decisions.exclude or collected_on_purpose_paths, and every item of a topic
// with relevance off, are fetched without a call; the rest are judged in requests of at most
// MaxItemsPerRequest items with the state
// `{contract, items:[{id,title,description,channel,date,provenance}]}` and
// one `role_item_{id}` question per item. In apply mode, skipped items are
// appended to skipped.jsonl and enqueued as `skip` review items; in shadow
// mode nothing is skipped and the would-be skip is returned in Shadow (the
// item is fetched and written with triage: review). Undecided items are
// fetched (a failed judgment is never a "no"). The error is fatal only.
func Prefetch(ctx context.Context, s *session.Session, items []Item, opts PrefetchOptions) ([]PrefetchDecision, error) {
	out := make([]PrefetchDecision, len(items))
	pending := make([]int, 0, len(items))
	for index, item := range items {
		if item.ID == "" {
			item.ID = SkippedID(item.URL)
		}
		out[index] = PrefetchDecision{Item: item, Action: ActionFetch}
		switch {
		case s.Excluded(item.Path):
			out[index].Reason = PrefetchExcluded
		case !s.RelevanceEnabled():
			out[index].Reason = PrefetchRelevanceOff
		case s.Contract.MatchesCollectedPath(item.Path):
			out[index].Reason = PrefetchCollectedPath
			out[index].Role = "collected_on_purpose"
			out[index].PKept = 1
		default:
			pending = append(pending, index)
		}
	}

	mode := s.RelevanceGateMode()
	quarantine, fetch := s.Threshold("relevance_quarantine"), s.Threshold("relevance_fetch")
	queue := review.Open(s.Root(), s.Now)
	for start := 0; start < len(pending); start += MaxItemsPerRequest {
		chunk := pending[start:min(start+MaxItemsPerRequest, len(pending))]
		request, err := prefetchRequest(s, out, chunk, opts)
		if err != nil {
			return out, fmt.Errorf("gate: pre-fetch relevance: %w", err)
		}
		result, err := s.Engine.Decide(ctx, request)
		if err != nil {
			return out, fmt.Errorf("gate: pre-fetch relevance: %w", err)
		}
		for position, index := range chunk {
			decision := &out[index]
			answer := result.Answers["role_item_"+itemQuestionID(position)]
			decision.ReceiptKey = answer.Receipt
			if !answer.Decided() {
				reason := answer.Reason
				if reason == "" {
					reason = decisions.ReasonMissingAnswer
				}
				decision.Reason = "undecided:" + reason
				continue
			}
			decision.Role = answer.Choice
			decision.POffTopic = answer.Probs["off_topic"]
			decision.PKept = answer.Probs["core"] + answer.Probs["adjacent"] + answer.Probs["collected_on_purpose"]
			decision.Action = PrefetchAction(decision.POffTopic, decision.PKept, quarantine, fetch)
			if decision.Action == ActionSkip && mode != session.ModeApply {
				decision.Action, decision.Shadow = ActionReview, ActionSkip
			}
			if decision.Action == ActionSkip {
				if err := recordSkip(s, queue, *decision, opts); err != nil {
					return out, err
				}
			}
		}
	}
	return out, nil
}

// itemQuestionID is the element id of the item at position in a request.
func itemQuestionID(position int) string { return "i" + strconv.Itoa(position+1) }

// prefetchRequest builds the request of one chunk: the built-in
// `role_item_{id}` per item plus the topic's extra relevance questions
// (templates once per item id; their answers are only recorded in
// receipts).
func prefetchRequest(s *session.Session, out []PrefetchDecision, chunk []int, opts PrefetchOptions) (decisions.Request, error) {
	stateItems := make([]map[string]any, 0, len(chunk))
	qs := make([]questions.Q, 0, len(chunk))
	ids := make([]string, 0, len(chunk))
	for position, index := range chunk {
		item := out[index].Item
		id := itemQuestionID(position)
		entry := map[string]any{"id": id, "provenance": itemProvenance(item, opts)}
		setText(entry, "title", item.Title)
		setText(entry, "description", item.Description)
		setText(entry, "channel", item.Channel)
		setText(entry, "date", item.Date)
		stateItems = append(stateItems, entry)
		qs = append(qs, relevanceBank.MustQuestion("role_item_{id}", map[string]string{"id": id}))
		ids = append(ids, id)
	}
	bank, qs, err := s.WithExtras(decisions.PurposeRelevance, relevanceBank, qs, ids...)
	if err != nil {
		return decisions.Request{}, err
	}
	contractState := map[string]any{}
	if !s.Contract.Empty() {
		contractState = s.Contract.QuestionState()
	}
	state := map[string]any{"contract": contractState, "items": stateItems}
	subject := "prefetch:" + opts.Batch
	return s.Request(decisions.PurposeRelevance, subject, bank, state, qs), nil
}

func itemProvenance(item Item, opts PrefetchOptions) map[string]any {
	provenance := map[string]any{}
	setText(provenance, "path", item.Path)
	setText(provenance, "source_kind", item.SourceKind)
	if key := NormalizeURL(item.URL); key != "" {
		host, _, _ := strings.Cut(key, "/")
		host, _, _ = strings.Cut(host, "?")
		setText(provenance, "source_host", host)
	}
	setText(provenance, "ingest_batch", opts.Batch)
	setText(provenance, "ingest_query", opts.Query)
	return provenance
}

func setText(values map[string]any, key, value string) {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		values[key] = trimmed
	}
}

// recordSkip appends the skipped row and the `skip` review item.
func recordSkip(s *session.Session, queue *review.Store, decision PrefetchDecision, opts PrefetchOptions) error {
	item := decision.Item
	row := SkippedRow{
		ID: item.ID, Time: s.Now().UTC().Format(time.RFC3339), Batch: opts.Batch, Query: opts.Query,
		URL: item.URL, Title: item.Title, Description: item.Description, Channel: item.Channel, Date: item.Date,
		POffTopic: decision.POffTopic, ReceiptKey: decision.ReceiptKey,
	}
	if err := AppendSkipped(s.Root(), row); err != nil {
		return err
	}
	label := item.Title
	if label == "" {
		label = item.URL
	}
	_, err := queue.Add(review.Item{
		Queue: review.QueueSkip, Purpose: string(decisions.PurposeRelevance), Subject: item.ID,
		Question: "role_item", Probability: decision.POffTopic, ReceiptKey: decision.ReceiptKey,
		Evidence: fmt.Sprintf("skipped before fetch: %s (%s) — P(off_topic) %.2f", label, item.URL, decision.POffTopic),
		Action:   map[string]any{"skipped_id": item.ID, "url": item.URL},
	})
	if err != nil {
		return fmt.Errorf("gate: skip review item: %w", err)
	}
	return nil
}
