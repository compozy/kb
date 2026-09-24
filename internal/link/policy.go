package link

import (
	"fmt"
	"slices"
	"sort"
	"strconv"

	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
)

// Relation names written by link (spec §6, §9.3).
const (
	RelationRelated     = "related"
	RelationContradicts = "contradicts"
	relationNone        = "none"
	// KeyAffects is the source → article relation written from affects_{id}.
	KeyAffects = "affects"
	// minRelationConfidence is the relation-choice confidence below which a
	// link falls back to `related`.
	minRelationConfidence = 0.5
)

// Thresholds are the named gates of the write policy.
type Thresholds struct {
	Apply   float64 // link_apply
	Review  float64 // link_review
	Mention float64 // mention_sense
	Affects float64 // affects
}

// ThresholdsFrom reads the link thresholds of a session.
func ThresholdsFrom(s *session.Session) Thresholds {
	return Thresholds{
		Apply:   s.Threshold("link_apply"),
		Review:  s.Threshold("link_review"),
		Mention: s.Threshold("mention_sense"),
		Affects: s.Threshold("affects"),
	}
}

// MentionRef is one mention asked about in the judgment: `mention_sense_{N}`.
type MentionRef struct {
	N int
	// Candidate is the candidate ID the mention points to.
	Candidate string
	Start     int
	End       int
	Text      string
	Sentence  string
}

// Plan is everything the write policy needs about one document, computed
// by code before the judgment.
type Plan struct {
	Doc        *corpus.Document
	Candidates []Candidate
	Mentions   []MentionRef
	// Forms maps candidate path → relation target form (stem or vault path).
	Forms map[string]string
	// Relations holds the document's current relation lists (raw entries)
	// for every relation key and `affects`.
	Relations map[string][]string
	// Present maps key → candidate paths already listed under that key.
	Present map[string]map[string]bool
	// KBOwned reports keys whose current value kb wrote (value hash equals
	// the state row's written hash).
	KBOwned map[string]bool
	// BodyLinked lists candidate paths the body already links to.
	BodyLinked map[string]bool
}

// Insertion is one body link kb inserts (apply) or proposes (shadow).
type Insertion struct {
	Target   string // candidate path
	Form     string
	Relation string
	Text     string
	Start    int
	End      int
	P        float64
	Receipt  string
}

// Outcome is the write policy's verdict for one document.
type Outcome struct {
	// Updates are the owned-key writes (relation lists, affects).
	Updates map[string]any
	// Added lists, per key, the targets newly added.
	Added map[string][]string
	// Insertions are the body links to insert (apply mode only).
	Insertions []Insertion
	// Items are review-queue items to add.
	Items          []review.Item
	Proposals      int
	Reviews        int
	Contradictions int
	Demotions      int
	Undecided      int
}

// Decide applies write policy (d) (spec §9.3) to the answers of one
// judgment. It is pure: nothing is written.
//
//   - should_link ≥ Apply: the target is added to the list of the argmax
//     relation (`none` or relation confidence < 0.5 → `related`) unless the
//     document already lists it under any relation, in both modes;
//     `contradicts` never applies and becomes a contradiction item.
//   - an applied target with a mention whose mention_sense ≥ Mention gets one
//     body link at its first qualifying mention: inserted when bodyMode is
//     apply, else proposed as a `link` review item.
//   - Review ≤ should_link < Apply: a `link` review item.
//   - below Apply, a target kb itself wrote to a relation list becomes a
//     `link` demotion item; nothing is ever removed.
//   - affects_{id} ≥ Affects on a source: the target is added to `affects`.
func Decide(plan Plan, answers map[string]decisions.Answer, th Thresholds, bodyMode string) Outcome {
	out := Outcome{Updates: map[string]any{}, Added: map[string][]string{}}
	subject := plan.Doc.Path
	additions := map[string][]string{}

	for _, candidate := range plan.Candidates {
		form := plan.Forms[candidate.Path]
		should, ok := answers["should_link_"+candidate.ID]
		p, decided := should.P("")
		if !ok || !decided {
			out.Undecided++
			continue
		}
		relation, relationP := relationOf(answers["relation_"+candidate.ID])

		switch {
		case p >= th.Review && relation == RelationContradicts:
			out.Contradictions++
			out.Items = append(out.Items, review.Item{
				Queue: review.QueueContradiction, Purpose: string(decisions.PurposeLink), Subject: subject,
				Target: candidate.Path, Question: "relation", Probability: relationP, ReceiptKey: should.Receipt,
				Evidence: fmt.Sprintf("%s contradicts %s (relation %.2f, should_link %.2f)", subject, candidate.Title, relationP, p),
				Action:   map[string]any{"relation": RelationContradicts, "target": form},
			})
		case p >= th.Apply:
			if presentUnder(plan, candidate.Path) == "" {
				additions[relation] = append(additions[relation], candidate.Path)
			}
			if insertion, ok := firstQualifyingMention(plan, candidate, answers, th.Mention); ok {
				insertion.Form, insertion.Relation = form, relation
				if bodyMode == session.ModeApply {
					out.Insertions = append(out.Insertions, insertion)
				} else {
					out.Proposals++
					out.Items = append(out.Items, review.Item{
						Queue: review.QueueLink, Purpose: string(decisions.PurposeMention), Subject: subject,
						Target: candidate.Path, Question: "mention_sense", Probability: insertion.P, ReceiptKey: insertion.Receipt,
						Evidence: fmt.Sprintf("insert [[%s|%s]] in %s (shadow)", form, insertion.Text, subject),
						Action: map[string]any{
							"insert": true, "relation": relation, "target": form,
							"mention_text": insertion.Text, "mention_start": insertion.Start,
						},
					})
				}
			}
		default:
			if key := presentUnder(plan, candidate.Path); key != "" && plan.KBOwned[key] {
				out.Demotions++
				out.Items = append(out.Items, review.Item{
					Queue: review.QueueLink, Purpose: string(decisions.PurposeLink), Subject: subject,
					Target: candidate.Path, Question: "demote:" + key, Probability: p, ReceiptKey: should.Receipt,
					Evidence: fmt.Sprintf("kb-written %s %s dropped to should_link %.2f", key, candidate.Title, p),
					Action:   map[string]any{"relation": key, "target": form, "demote": true},
				})
				continue
			}
			if p >= th.Review {
				out.Reviews++
				out.Items = append(out.Items, review.Item{
					Queue: review.QueueLink, Purpose: string(decisions.PurposeLink), Subject: subject,
					Target: candidate.Path, Question: "should_link", Probability: p, ReceiptKey: should.Receipt,
					Evidence: fmt.Sprintf("should_link %.2f: %s → %s (%s)", p, subject, candidate.Title, relation),
					Action:   map[string]any{"relation": relation, "target": form},
				})
			}
		}
	}

	if plan.Doc.Kind == corpus.KindSource {
		for _, candidate := range plan.Candidates {
			answer, ok := answers["affects_"+candidate.ID]
			p, decided := answer.P("")
			if !ok || !decided {
				out.Undecided++
				continue
			}
			if p >= th.Affects && !plan.Present[KeyAffects][candidate.Path] {
				additions[KeyAffects] = append(additions[KeyAffects], candidate.Path)
			}
		}
	}

	for key, paths := range additions {
		entries := slices.Clone(plan.Relations[key])
		for _, p := range paths {
			entries = append(entries, "[["+plan.Forms[p]+"]]")
		}
		out.Updates[key] = entries
		out.Added[key] = paths
	}
	sort.Slice(out.Insertions, func(i, j int) bool { return out.Insertions[i].Start < out.Insertions[j].Start })
	return out
}

// relationOf returns the relation key for a relation choice: the argmax
// option, with `none`, an undecided answer or a confidence below 0.5
// falling back to `related`. The second value is the argmax probability.
func relationOf(answer decisions.Answer) (string, float64) {
	if !answer.Decided() || answer.Choice == "" {
		return RelationRelated, 0
	}
	p := answer.Probs[answer.Choice]
	if answer.Choice == RelationContradicts {
		return RelationContradicts, p
	}
	if answer.Choice == relationNone || !slices.Contains(decisions.RelationKeys, answer.Choice) {
		return RelationRelated, p
	}
	if answer.Confidence != nil && *answer.Confidence < minRelationConfidence {
		return RelationRelated, p
	}
	return answer.Choice, p
}

// presentUnder returns the relation key that already lists path, or "".
func presentUnder(plan Plan, candidatePath string) string {
	for _, key := range decisions.RelationKeys {
		if plan.Present[key][candidatePath] {
			return key
		}
	}
	return ""
}

func firstQualifyingMention(plan Plan, candidate Candidate, answers map[string]decisions.Answer, threshold float64) (Insertion, bool) {
	if plan.BodyLinked[candidate.Path] {
		return Insertion{}, false
	}
	for _, mention := range plan.Mentions {
		if mention.Candidate != candidate.ID {
			continue
		}
		answer := answers["mention_sense_"+strconv.Itoa(mention.N)]
		p, ok := answer.P("")
		if !ok || p < threshold {
			continue
		}
		return Insertion{
			Target: candidate.Path, Text: mention.Text, Start: mention.Start, End: mention.End,
			P: p, Receipt: answer.Receipt,
		}, true
	}
	return Insertion{}, false
}

// InsertLinks returns body with `[[<form>|<text>]]` at every insertion.
// Insertions must not overlap; they are applied from the end so offsets stay
// valid.
func InsertLinks(body string, insertions []Insertion) string {
	ordered := slices.Clone(insertions)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Start > ordered[j].Start })
	for _, insertion := range ordered {
		if insertion.Start < 0 || insertion.End > len(body) || insertion.Start >= insertion.End {
			continue
		}
		text := body[insertion.Start:insertion.End]
		body = body[:insertion.Start] + "[[" + insertion.Form + "|" + text + "]]" + body[insertion.End:]
	}
	return body
}
