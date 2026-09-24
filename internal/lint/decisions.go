package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/contract/topicyaml"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/resolve"
	"github.com/compozy/kb/internal/review"
)

// maxListedPaths caps the paths named in one aggregated issue (unclassified)
// so a large vault gets one readable line instead of thousands of rows.
const maxListedPaths = 10

// Bank ids recorded in state rows per document kind (plan W3/W4 contracts):
// `unclassified` compares their recorded versions with the embedded banks.
var (
	sourceBankIDs  = []string{"relevance", "quality", "classify", "concept"}
	articleBankIDs = []string{"classify", "concept"}
)

// classifyBankID is the bank every classification pass records; a state row
// without it (for example one written only by link) is not a classification.
const classifyBankID = "classify"

// Keys kb writes once, when ingest creates the file, and never rewrites;
// value-hash ownership does not apply to them. `aliases` is merged, never
// replaced, so a user alias list is not a conflict either (spec §5.2).
var creationOnlyKeys = []string{"ingest_batch", "ingest_query", "aliases"}

// decisionFiles are the topic-relative paths lint reads besides markdown.
const (
	statePath  = resolve.DecisionsDir + "/state.jsonl"
	reviewPath = resolve.DecisionsDir + "/" + review.QueueFile
)

// findDecisionIssues reports the decision-workflow kinds of spec §11 from
// frontmatter, topic.yaml and the records under `.decisions/`. It never calls
// a model and never writes. Unreadable records become format errors on the
// record file instead of failing the whole lint run.
//
// The kinds that describe progress through the decision workflow
// (unclassified, contract-missing, vocabulary-missing, criterion-missing) only
// apply to a topic that has adopted it: it has records under `.decisions/`, an
// accepted contract, a contract draft or a vocabulary draft. A topic that was
// never touched by a decision-backed command keeps its structural lint.
func findDecisionIssues(state vaultState) []models.LintIssue {
	root := state.topicPath
	issues := make([]models.LintIssue, 0)

	settings, err := contract.LoadSettings(root)
	if err != nil {
		issues = append(issues, recordIssue(topicyaml.FileName, err))
		settings = contract.Settings{}
	}
	store, err := corpus.OpenState(root)
	if err != nil {
		issues = append(issues, recordIssue(statePath, err))
	}
	loaded, err := corpus.Load(root, corpus.LoadOptions{Exclude: settings.Decisions.Exclude})
	if err != nil {
		issues = append(issues, recordIssue(".", err))
		return issues
	}

	userKeys := keyConflicts(loaded, store)
	issues = append(issues, keyConflictIssues(userKeys)...)
	for _, file := range state.files {
		if file.parseErr != nil || !isDecisionDocument(file) {
			continue
		}
		issues = append(issues, ownedKeyIssues(file, userKeys[file.relativePath])...)
	}

	issues = append(issues, offTopicKeptIssues(loaded)...)
	issues = append(issues, reviewQueueIssues(root, state.topicSlug)...)
	if settings.ContractDraft != nil && !settings.ContractDraft.Empty() {
		issues = append(issues, newIssue(
			models.LintIssueKindContractDraftPending,
			models.SeverityInfo,
			topicyaml.FileName,
			fmt.Sprintf("contract_draft is waiting for review: run `kb topic contract %s --accept`", state.topicSlug),
			"contract_draft",
		))
	}

	if !decisionsAdopted(root, settings) {
		return issues
	}

	sources, articles := loaded.Sources(), loaded.Articles()
	if len(sources) > 0 && !settings.Accepted() && relevanceOn(settings) {
		issues = append(issues, newIssue(
			models.LintIssueKindContractMissing,
			models.SeverityWarning,
			topicyaml.FileName,
			fmt.Sprintf("topic has %d source(s) and no accepted selection contract; relevance gates stay in shadow: run `kb topic contract %s --draft` (or --import-claude), then --accept", len(sources), state.topicSlug),
			"contract",
		))
	}
	if len(sources) > 0 && len(articles) == 0 {
		next := "--draft"
		if len(settings.VocabularyDraft) > 0 {
			next = "--accept"
		}
		issues = append(issues, newIssue(
			models.LintIssueKindVocabularyMissing,
			models.SeverityWarning,
			topicyaml.FileName,
			fmt.Sprintf("topic has %d source(s) and no concept articles, so classification writes no concepts: run `kb topic vocabulary %s %s`", len(sources), state.topicSlug, next),
			"vocabulary",
		))
	}
	for _, article := range articles {
		if article.Criterion() != "" {
			continue
		}
		issues = append(issues, newIssue(
			models.LintIssueKindCriterionMissing,
			models.SeverityWarning,
			article.Path,
			fmt.Sprintf("article has no criterion, which concept questions judge against: run `kb classify %s --only-missing`", state.topicSlug),
			"criterion",
		))
	}
	if store != nil {
		activeContract := ""
		if settings.Accepted() {
			activeContract = settings.Contract.Hash()
		}
		issues = append(issues, unclassifiedIssue(loaded, store, activeContract, state.topicSlug)...)
	}

	return issues
}

// decisionsAdopted reports whether the topic entered the decision workflow.
func decisionsAdopted(root string, settings contract.Settings) bool {
	if settings.Accepted() || (settings.ContractDraft != nil && !settings.ContractDraft.Empty()) || len(settings.VocabularyDraft) > 0 {
		return true
	}
	entries, err := os.ReadDir(filepath.Join(root, resolve.DecisionsDir))
	return err == nil && len(entries) > 0
}

func relevanceOn(settings contract.Settings) bool {
	return !strings.EqualFold(strings.TrimSpace(settings.Decisions.Relevance), contract.RelevanceOff)
}

// isDecisionDocument reports whether file is a source (quarantined included)
// or a wiki article: the documents whose kb-owned keys lint validates.
func isDecisionDocument(file *vaultFile) bool {
	if strings.HasPrefix(file.relativePath, "wiki/concepts/") {
		return true
	}
	if isQuarantinedFile(file) {
		return true
	}
	return isSourceFile(file)
}

func recordIssue(filePath string, err error) models.LintIssue {
	return newIssue(models.LintIssueKindFormat, models.SeverityError, filePath, fmt.Sprintf("cannot read decision records: %v", err), "")
}

// keyConflicts returns, per document path, the kb-owned keys that belong to
// the user (spec §6 conflicts): the document has a state row and the key's
// current value hash differs from written[key] or was never written by kb.
func keyConflicts(loaded *corpus.Corpus, store *corpus.StateStore) map[string]map[string]struct{} {
	conflicts := make(map[string]map[string]struct{})
	if store == nil {
		return conflicts
	}
	for _, doc := range loaded.Documents() {
		row, ok := store.Get(doc.Path)
		if !ok {
			continue
		}
		for _, key := range decisions.OwnedKeys {
			value, present := doc.Frontmatter[key]
			if !present || slices.Contains(creationOnlyKeys, key) {
				continue
			}
			if written, ok := row.Written[key]; ok && written == corpus.ValueHash(value) {
				continue
			}
			if conflicts[doc.Path] == nil {
				conflicts[doc.Path] = make(map[string]struct{})
			}
			conflicts[doc.Path][key] = struct{}{}
		}
	}
	return conflicts
}

func keyConflictIssues(conflicts map[string]map[string]struct{}) []models.LintIssue {
	issues := make([]models.LintIssue, 0)
	for _, path := range sortedKeys(conflicts) {
		for _, key := range sortedKeys(conflicts[path]) {
			issues = append(issues, newIssue(
				models.LintIssueKindKeyConflict,
				models.SeverityWarning,
				path,
				fmt.Sprintf("kb-owned key %q holds a value kb did not write (or that was edited since); kb leaves it untouched", key),
				key,
			))
		}
	}
	return issues
}

// offTopicKeptIssues reports sources judged off-topic that were not
// quarantined (spec §11, §15).
func offTopicKeptIssues(loaded *corpus.Corpus) []models.LintIssue {
	issues := make([]models.LintIssue, 0)
	for _, doc := range loaded.Sources() {
		if doc.Relevance() != "off_topic" || doc.Triage() == "quarantined" {
			continue
		}
		issues = append(issues, newIssue(
			models.LintIssueKindOffTopicKept,
			models.SeverityWarning,
			doc.Path,
			"source is judged off_topic but still kept; review it in the remove queue or quarantine it",
			"relevance",
		))
	}
	return issues
}

// reviewQueueIssues reports open contradictions one by one and pending
// review items as one count per queue (spec §11, §12.1).
func reviewQueueIssues(root, topicSlug string) []models.LintIssue {
	store := review.Open(root, nil)
	pending, err := store.Items("", review.StatusPending)
	if err != nil {
		return []models.LintIssue{recordIssue(reviewPath, err)}
	}

	issues := make([]models.LintIssue, 0)
	counts := make(map[string]int)
	for _, item := range pending {
		counts[item.Queue]++
		if item.Queue != review.QueueContradiction {
			continue
		}
		issues = append(issues, newIssue(
			models.LintIssueKindContradiction,
			models.SeverityWarning,
			item.Subject,
			fmt.Sprintf("open contradiction (P=%.2f, review item %s): run `kb review list %s --queue %s`", item.Probability, item.ID, topicSlug, review.QueueContradiction),
			item.Target,
		))
	}

	for _, queue := range review.Queues {
		count := counts[queue]
		if count == 0 || queue == review.QueueContradiction {
			continue
		}
		kind := models.LintIssueKindPendingReview
		switch queue {
		case review.QueueRecapture:
			kind = models.LintIssueKindRecapturePending
		case review.QueueRemove:
			kind = models.LintIssueKindRemovePending
		}
		issues = append(issues, newIssue(
			kind,
			models.SeverityInfo,
			reviewPath,
			fmt.Sprintf("%d pending item(s) in queue %q: run `kb review list %s --queue %s`", count, queue, topicSlug, queue),
			queue,
		))
	}
	return issues
}

// unclassifiedIssue aggregates the documents whose classification is
// missing or out of date into one issue per topic: no state row, a changed
// body, a row written under another contract, a row without a classification
// bank, or a recorded bank version that differs from the embedded bank. At
// most maxListedPaths paths are named; the count covers all of them.
func unclassifiedIssue(loaded *corpus.Corpus, store *corpus.StateStore, activeContract, topicSlug string) []models.LintIssue {
	current := currentBankVersions()
	reasons := make(map[string]int)
	paths := make([]string, 0)
	exists := func(path string) bool { return loaded.ByPath(path) != nil }

	for _, doc := range loaded.Documents() {
		bankIDs := sourceBankIDs
		if doc.Kind == corpus.KindArticle {
			bankIDs = articleBankIDs
		}
		row, ok := store.Lookup(doc.Path, doc.BodyHash, exists)
		reason := unclassifiedReason(doc, row, ok, activeContract, bankIDs, current)
		if reason == "" {
			continue
		}
		reasons[reason]++
		paths = append(paths, doc.Path)
	}
	if len(paths) == 0 {
		return nil
	}

	parts := make([]string, 0, len(reasons))
	for _, reason := range sortedKeys(reasons) {
		parts = append(parts, fmt.Sprintf("%s: %d", reason, reasons[reason]))
	}
	listed := paths[:min(len(paths), maxListedPaths)]
	message := fmt.Sprintf(
		"%d document(s) need classification (%s); run `kb classify %s`: %s",
		len(paths), strings.Join(parts, ", "), topicSlug, strings.Join(listed, ", "),
	)
	if extra := len(paths) - len(listed); extra > 0 {
		message += fmt.Sprintf(" (+%d more)", extra)
	}
	return []models.LintIssue{newIssue(
		models.LintIssueKindUnclassified,
		models.SeverityInfo,
		statePath,
		message,
		fmt.Sprintf("%d", len(paths)),
	)}
}

func unclassifiedReason(doc *corpus.Document, row *corpus.StateRow, ok bool, activeContract string, bankIDs []string, current map[string]string) string {
	switch {
	case !ok || row == nil:
		return "no state record"
	case row.BodyHash != doc.BodyHash:
		return "body changed"
	case activeContract != "" && row.Contract != activeContract:
		return "contract changed"
	case row.Banks[classifyBankID] == "":
		return "not classified"
	}
	for _, id := range bankIDs {
		recorded, has := row.Banks[id]
		if has && current[id] != "" && recorded != current[id] {
			return "question bank changed"
		}
	}
	return ""
}

func currentBankVersions() map[string]string {
	versions := make(map[string]string)
	for _, id := range sourceBankIDs {
		if bank, err := questions.Load(id); err == nil {
			versions[id] = bank.Version
		}
	}
	return versions
}
