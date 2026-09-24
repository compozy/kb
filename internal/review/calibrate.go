package review

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/decisions"
)

// Calibration floors and targets (spec §12.3).
const (
	// MinDevLabels is the dev floor per purpose.
	MinDevLabels = 30
	// MinHoldoutLabels is the holdout floor per purpose.
	MinHoldoutLabels = 10
	// TargetPrecision is the dev precision the sweep looks for.
	TargetPrecision = 0.90
	// devBuckets is the number of sha256(subject) mod 10 buckets (0..6) that
	// form the dev split; 7..9 are holdout.
	devBuckets = 7
)

// CalibrationFile is <topic>/.decisions/calibration.json. The JSON shape is
// the one internal/session reads (session.Calibration); it is mirrored here
// because session imports review.
const CalibrationFile = "calibration.json"

// ErrNotEnoughLabels is returned by WriteCalibration when no purpose has a
// recommendation.
var ErrNotEnoughLabels = errors.New("review: not enough labels; nothing written")

// qualityQuestions are the quality nouls whose maximum is the quality gate
// quantity.
var qualityQuestions = []string{"paywall_or_login", "error_or_placeholder_page", "thin_or_boilerplate", "no_speech_content"}

// purposeSpec describes how one purpose is calibrated.
type purposeSpec struct {
	purpose  string
	quantity string
	apply    string
	review   string
	// target is the label verdict the gate quantity predicts.
	target string
}

// calibratedPurposes lists the calibrated purposes and their gate quantity:
//   - relevance: P(off_topic) of the `role` choice predicts a negative
//     (off-topic) label; threshold relevance_quarantine;
//   - link: P(yes) of the `should_link_<id>` noul predicts a positive label;
//     threshold link_apply;
//   - quality: the maximum of the quality nouls predicts a negative (broken
//     capture) label; threshold quality_apply.
var calibratedPurposes = []purposeSpec{
	{purpose: PurposeRelevance, quantity: "P(off_topic) of role predicts a negative label", apply: "relevance_quarantine", review: "relevance_review", target: VerdictNegative},
	{purpose: PurposeLink, quantity: "P(yes) of should_link_<id> predicts a positive label", apply: "link_apply", review: "link_review", target: VerdictPositive},
	{purpose: PurposeQuality, quantity: "max quality noul predicts a negative label", apply: "quality_apply", review: "quality_review", target: VerdictNegative},
}

// CalibrateOptions configures Calibrate.
type CalibrateOptions struct {
	// ContractHash is the active contract hash; labels shown in a preview
	// made under another contract are forced to dev.
	ContractHash string
	// Thresholds are the current thresholds; nil means the defaults merged
	// with topic.yaml `decisions.thresholds`.
	Thresholds decisions.Thresholds
	// Now defaults to time.Now.
	Now func() time.Time
}

// Metrics are the scores of one split at one threshold. Labels without a
// judged gate quantity count in Total only.
type Metrics struct {
	// Total is the number of labels in the split.
	Total int `json:"total"`
	// N is the number of labels joined to a judged gate quantity.
	N int `json:"n"`
	// Predicted counts joined labels at or above the apply threshold.
	Predicted int `json:"predicted"`
	// TruePositive counts predicted labels whose verdict is the target.
	TruePositive int `json:"true_positive"`
	// Precision = TruePositive / Predicted.
	Precision float64 `json:"precision"`
	// Recall = TruePositive / joined labels whose verdict is the target.
	Recall float64 `json:"recall"`
	// Coverage = N / Total: the share of labels the model judged.
	Coverage float64 `json:"coverage"`
	// ReviewRate = joined labels in [review, apply) / N.
	ReviewRate float64 `json:"review_rate"`
}

// SweepPoint is the dev score at one candidate threshold.
type SweepPoint struct {
	Threshold float64 `json:"threshold"`
	Dev       Metrics `json:"dev"`
}

// PurposeReport is the calibration of one purpose.
type PurposeReport struct {
	Purpose  string `json:"purpose"`
	Quantity string `json:"quantity"`
	// ApplyName and ReviewName are the threshold names calibrated.
	ApplyName  string  `json:"apply_name"`
	ReviewName string  `json:"review_name"`
	Current    float64 `json:"current"`
	Review     float64 `json:"review"`
	// DevLabels and HoldoutLabels count joined labels per split (imported
	// ones included); ImportedLabels counts the imported ones among them.
	DevLabels      int `json:"dev_labels"`
	HoldoutLabels  int `json:"holdout_labels"`
	ImportedLabels int `json:"imported_labels"`
	// ReviewDevLabels counts the joined dev labels given in `kb review`
	// (not imported); the dev floor must be met by these alone (spec §12.2:
	// imported labels are never the only evidence behind an apply).
	ReviewDevLabels int `json:"review_dev_labels"`
	// ImportedOnly flags a purpose whose dev floor is met only with imported
	// labels: metrics are reported, no threshold is recommended.
	ImportedOnly bool `json:"imported_only,omitempty"`
	// Demoted counts labels forced to dev by a preview under another
	// contract.
	Demoted int `json:"demoted"`
	// Unjoined counts labels with no judged gate quantity.
	Unjoined int `json:"unjoined"`
	// Recommended reports whether Chosen is a recommendation.
	Recommended bool    `json:"recommended"`
	Chosen      float64 `json:"chosen,omitempty"`
	// Note explains a missing recommendation ("not enough labels", ...).
	Note           string       `json:"note,omitempty"`
	DevCurrent     Metrics      `json:"dev_current"`
	HoldoutCurrent Metrics      `json:"holdout_current"`
	DevChosen      Metrics      `json:"dev_chosen"`
	HoldoutChosen  Metrics      `json:"holdout_chosen"`
	Imported       ImportedPart `json:"imported"`
	Sweep          []SweepPoint `json:"sweep"`
}

// ImportedPart scores imported labels alone, at the chosen threshold (or the
// current one when there is no recommendation).
type ImportedPart struct {
	Dev     Metrics `json:"dev"`
	Holdout Metrics `json:"holdout"`
}

// Report is a calibration report over every calibrated purpose.
type Report struct {
	Time     string          `json:"time"`
	Contract string          `json:"contract"`
	Purposes []PurposeReport `json:"purposes"`
}

// scored is one label with its split and gate quantity.
type scored struct {
	label    Label
	dev      bool
	demoted  bool
	imported bool
	value    float64
	joined   bool
}

// Calibrate measures, per purpose, precision, recall, coverage and review
// rate at the current thresholds and over a sweep, from labels only (spec
// §12.3). Each label is joined to a receipt (Label.ReceiptKey when it carries
// the needed answer, else the latest decided receipt for the same subject;
// link labels use the question id of the label or of the matching review
// item). Labels are split by the first byte of sha256(subject) mod 10:
// buckets 0–6 are dev, 7–9 holdout; subjects that appeared in a preview
// under another contract are forced to dev. Thresholds 0.50..0.95 (step
// 0.05) are swept on dev only: the lowest threshold with dev precision ≥
// 0.90 wins, else the best F0.5. Below 30 dev or 10 holdout joined labels a
// purpose gets "not enough labels" and no recommendation; when the dev floor
// is met only with imported labels (fewer than 30 dev labels given in `kb
// review`) the purpose is flagged "imported labels only" and gets no
// recommendation either (spec §12.2, §20).
func Calibrate(topicRoot string, opts CalibrateOptions) (Report, error) {
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	thresholds := opts.Thresholds
	if thresholds == nil {
		settings, err := contract.LoadSettings(topicRoot)
		if err != nil {
			return Report{}, fmt.Errorf("review: %w", err)
		}
		thresholds = decisions.Thresholds(nil).Merge(settings.Decisions.Thresholds)
	}
	labels, err := LoadLabels(topicRoot)
	if err != nil {
		return Report{}, err
	}
	receipts, err := decisions.LoadReceipts(topicRoot)
	if err != nil {
		return Report{}, fmt.Errorf("review: %w", err)
	}
	previews, err := LoadPreviews(topicRoot)
	if err != nil {
		return Report{}, err
	}
	items, err := Open(topicRoot, nil).Items("", "")
	if err != nil {
		return Report{}, err
	}

	demoted := map[string]bool{}
	for _, preview := range previews {
		if preview.Contract == opts.ContractHash {
			continue
		}
		for _, subject := range preview.Subjects {
			demoted[subject] = true
		}
	}
	joiner := newReceiptJoiner(receipts, items)

	report := Report{Time: now().UTC().Format(time.RFC3339), Contract: opts.ContractHash}
	latest := latestLabels(labels)
	for _, spec := range calibratedPurposes {
		rows := make([]scored, 0)
		for _, label := range latest {
			if label.Purpose != spec.purpose {
				continue
			}
			value, ok := joiner.quantity(spec.purpose, label)
			row := scored{label: label, value: value, joined: ok, imported: isImported(label.Origin)}
			row.dev = DevBucket(label.Subject)
			if !row.dev && demoted[label.Subject] {
				row.dev, row.demoted = true, true
			}
			rows = append(rows, row)
		}
		report.Purposes = append(report.Purposes, calibratePurpose(spec, rows, thresholds))
	}
	return report, nil
}

func calibratePurpose(spec purposeSpec, rows []scored, thresholds decisions.Thresholds) PurposeReport {
	current, review := thresholds.Get(spec.apply), thresholds.Get(spec.review)
	result := PurposeReport{
		Purpose: spec.purpose, Quantity: spec.quantity, ApplyName: spec.apply, ReviewName: spec.review,
		Current: current, Review: review,
	}
	var dev, holdout, importedDev, importedHoldout []scored
	for _, row := range rows {
		if row.demoted {
			result.Demoted++
		}
		if !row.joined {
			result.Unjoined++
		}
		if row.dev {
			dev = append(dev, row)
			if row.joined && !row.imported {
				result.ReviewDevLabels++
			}
		} else {
			holdout = append(holdout, row)
		}
		if row.imported {
			if row.joined {
				result.ImportedLabels++
			}
			if row.dev {
				importedDev = append(importedDev, row)
			} else {
				importedHoldout = append(importedHoldout, row)
			}
		}
	}
	result.DevCurrent = score(dev, current, review, spec.target)
	result.HoldoutCurrent = score(holdout, current, review, spec.target)
	result.DevLabels, result.HoldoutLabels = result.DevCurrent.N, result.HoldoutCurrent.N

	for step := 10; step <= 19; step++ {
		threshold := float64(step) * 0.05
		threshold = math.Round(threshold*100) / 100
		result.Sweep = append(result.Sweep, SweepPoint{Threshold: threshold, Dev: score(dev, threshold, review, spec.target)})
	}

	chosen := current
	switch {
	case result.DevLabels < MinDevLabels || result.HoldoutLabels < MinHoldoutLabels:
		result.Note = fmt.Sprintf("not enough labels (dev %d/%d, holdout %d/%d)", result.DevLabels, MinDevLabels, result.HoldoutLabels, MinHoldoutLabels)
	case result.ReviewDevLabels < MinDevLabels:
		result.ImportedOnly = true
		result.Note = fmt.Sprintf("imported labels only: %d/%d dev labels given in kb review; imported labels never set a threshold alone",
			result.ReviewDevLabels, MinDevLabels)
	default:
		if threshold, ok := chooseThreshold(result.Sweep); ok {
			result.Recommended, result.Chosen, chosen = true, threshold, threshold
		} else {
			result.Note = "no swept threshold predicts a correct label on dev"
		}
	}
	result.DevChosen = score(dev, chosen, review, spec.target)
	result.HoldoutChosen = score(holdout, chosen, review, spec.target)
	result.Imported = ImportedPart{
		Dev:     score(importedDev, chosen, review, spec.target),
		Holdout: score(importedHoldout, chosen, review, spec.target),
	}
	return result
}

// chooseThreshold returns the lowest threshold with dev precision ≥ 0.90,
// else the threshold with the best F0.5 (lowest on ties); false when no
// threshold has a true positive.
func chooseThreshold(sweep []SweepPoint) (float64, bool) {
	for _, point := range sweep {
		if point.Dev.Predicted > 0 && point.Dev.Precision >= TargetPrecision {
			return point.Threshold, true
		}
	}
	best, bestF := 0.0, 0.0
	for _, point := range sweep {
		if f := fHalf(point.Dev); f > bestF {
			best, bestF = point.Threshold, f
		}
	}
	return best, bestF > 0
}

func fHalf(m Metrics) float64 {
	if m.Precision == 0 && m.Recall == 0 {
		return 0
	}
	return 1.25 * m.Precision * m.Recall / (0.25*m.Precision + m.Recall)
}

func score(rows []scored, apply, review float64, target string) Metrics {
	m := Metrics{Total: len(rows)}
	targets, inReview := 0, 0
	for _, row := range rows {
		if !row.joined {
			continue
		}
		m.N++
		isTarget := row.label.Verdict == target
		if isTarget {
			targets++
		}
		switch {
		case row.value >= apply:
			m.Predicted++
			if isTarget {
				m.TruePositive++
			}
		case row.value >= review:
			inReview++
		}
	}
	m.Precision = ratio(m.TruePositive, m.Predicted)
	m.Recall = ratio(m.TruePositive, targets)
	m.Coverage = ratio(m.N, m.Total)
	m.ReviewRate = ratio(inReview, m.N)
	return m
}

func ratio(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

// DevBucket reports whether subject falls in the dev split: the first byte
// of sha256(subject) mod 10 is 0–6.
func DevBucket(subject string) bool {
	sum := sha256.Sum256([]byte(subject))
	return int(sum[0])%10 < devBuckets
}

func isImported(origin string) bool {
	return strings.HasPrefix(origin, OriginImportPrefix) || origin == OriginImportLinks
}

// latestLabels keeps the last label per purpose/subject/target, in first-
// appearance order.
func latestLabels(labels []Label) []Label {
	index := map[string]int{}
	out := make([]Label, 0, len(labels))
	for _, label := range labels {
		key := label.Purpose + "\x00" + label.Subject + "\x00" + label.Target
		if position, ok := index[key]; ok {
			out[position] = label
			continue
		}
		index[key] = len(out)
		out = append(out, label)
	}
	return out
}

// rawAnswer is the part of a receipt answer calibration reads.
type rawAnswer struct {
	Noul          *float64           `json:"noul"`
	Probabilities map[string]float64 `json:"probabilities"`
}

// receiptJoiner finds the gate quantity of a label in the receipts.
type receiptJoiner struct {
	byKey     map[string]decisions.Receipt
	bySubject map[string][]decisions.Receipt
	linkItems map[string]Item
}

func newReceiptJoiner(receipts []decisions.Receipt, items []Item) receiptJoiner {
	j := receiptJoiner{byKey: map[string]decisions.Receipt{}, bySubject: map[string][]decisions.Receipt{}, linkItems: map[string]Item{}}
	for _, receipt := range receipts {
		if receipt.Status != decisions.StatusDecided {
			continue
		}
		j.byKey[receipt.Key] = receipt
		j.bySubject[receipt.Subject] = append(j.bySubject[receipt.Subject], receipt)
	}
	for _, item := range items {
		if item.Purpose == PurposeLink || item.Queue == QueueLink {
			j.linkItems[item.Subject+"\x00"+item.Target] = item
		}
	}
	return j
}

func (j receiptJoiner) quantity(purpose string, label Label) (float64, bool) {
	receiptKey := label.ReceiptKey
	extract := func(r decisions.Receipt) (float64, bool) { return quantityOf(purpose, "", r) }
	if purpose == PurposeLink {
		question := label.Question
		if !strings.HasPrefix(question, QuestionShouldLink+"_") {
			item, ok := j.linkItems[label.Subject+"\x00"+label.Target]
			if !ok || !strings.HasPrefix(item.Question, QuestionShouldLink+"_") {
				return 0, false
			}
			question = item.Question
			if receiptKey == "" {
				receiptKey = item.ReceiptKey
			}
		}
		extract = func(r decisions.Receipt) (float64, bool) { return quantityOf(purpose, question, r) }
	}
	if receipt, ok := j.byKey[receiptKey]; ok && receiptKey != "" {
		if value, ok := extract(receipt); ok {
			return value, true
		}
	}
	candidates := j.bySubject[label.Subject]
	for _, candidate := range slices.Backward(candidates) {
		if value, ok := extract(candidate); ok {
			return value, true
		}
	}
	return 0, false
}

func quantityOf(purpose, question string, receipt decisions.Receipt) (float64, bool) {
	decode := func(id string) (rawAnswer, bool) {
		raw, ok := receipt.Answers[id]
		if !ok {
			return rawAnswer{}, false
		}
		var answer rawAnswer
		if json.Unmarshal(raw, &answer) != nil {
			return rawAnswer{}, false
		}
		return answer, true
	}
	switch purpose {
	case PurposeRelevance:
		answer, ok := decode(QuestionRole)
		if !ok {
			return 0, false
		}
		p, ok := answer.Probabilities["off_topic"]
		return p, ok
	case PurposeLink:
		answer, ok := decode(question)
		if !ok || answer.Noul == nil {
			return 0, false
		}
		return *answer.Noul, true
	case PurposeQuality:
		best, found := 0.0, false
		for _, id := range qualityQuestions {
			if answer, ok := decode(id); ok && answer.Noul != nil {
				best, found = math.Max(best, *answer.Noul), true
			}
		}
		return best, found
	}
	return 0, false
}

// calibrationRecord mirrors session.Calibration.
type calibrationRecord struct {
	Time     string                        `json:"time"`
	Purposes map[string]calibrationPurpose `json:"purposes"`
}

type calibrationPurpose struct {
	DevLabels      int                `json:"dev_labels"`
	HoldoutLabels  int                `json:"holdout_labels"`
	ImportedLabels int                `json:"imported_labels,omitempty"`
	Thresholds     map[string]float64 `json:"thresholds"`
	Dev            calibrationMetrics `json:"dev"`
	Holdout        calibrationMetrics `json:"holdout"`
}

type calibrationMetrics struct {
	Precision  float64 `json:"precision"`
	Recall     float64 `json:"recall"`
	Coverage   float64 `json:"coverage"`
	ReviewRate float64 `json:"review_rate"`
	N          int     `json:"n"`
}

// WriteCalibration stores every recommended threshold in topic.yaml
// `decisions.thresholds` (other thresholds are kept) and records the
// calibration of those purposes in <topic>/.decisions/calibration.json
// (entries of other purposes are kept). Without any recommendation nothing
// is written and ErrNotEnoughLabels is returned.
func WriteCalibration(topicRoot string, report Report) error {
	record := calibrationRecord{Time: report.Time, Purposes: map[string]calibrationPurpose{}}
	path := filepath.Join(topicRoot, decisions.ReceiptsDir, CalibrationFile)
	if data, err := os.ReadFile(path); err == nil {
		var existing calibrationRecord
		if json.Unmarshal(data, &existing) == nil && existing.Purposes != nil {
			record.Purposes = existing.Purposes
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("review: read calibration: %w", err)
	}

	chosen := map[string]float64{}
	for _, purpose := range report.Purposes {
		if !purpose.Recommended {
			continue
		}
		chosen[purpose.ApplyName] = purpose.Chosen
		record.Purposes[purpose.Purpose] = calibrationPurpose{
			DevLabels:      purpose.DevLabels,
			HoldoutLabels:  purpose.HoldoutLabels,
			ImportedLabels: purpose.ImportedLabels,
			Thresholds:     map[string]float64{purpose.ApplyName: purpose.Chosen},
			Dev:            storedMetrics(purpose.DevChosen),
			Holdout:        storedMetrics(purpose.HoldoutChosen),
		}
	}
	if len(chosen) == 0 {
		return ErrNotEnoughLabels
	}

	settings, err := contract.LoadSettings(topicRoot)
	if err != nil {
		return fmt.Errorf("review: %w", err)
	}
	merged := maps.Clone(settings.Decisions.Thresholds)
	if merged == nil {
		merged = map[string]float64{}
	}
	maps.Copy(merged, chosen)
	if err := contract.SetThresholds(topicRoot, merged); err != nil {
		return fmt.Errorf("review: %w", err)
	}

	if record.Time == "" {
		record.Time = time.Now().UTC().Format(time.RFC3339)
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("review: encode calibration: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("review: create %s: %w", filepath.Dir(path), err)
	}
	temp := path + ".tmp"
	if err := os.WriteFile(temp, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("review: write calibration: %w", err)
	}
	if err := os.Rename(temp, path); err != nil {
		return fmt.Errorf("review: write calibration: %w", err)
	}
	return nil
}

func storedMetrics(m Metrics) calibrationMetrics {
	return calibrationMetrics{Precision: m.Precision, Recall: m.Recall, Coverage: m.Coverage, ReviewRate: m.ReviewRate, N: m.N}
}
