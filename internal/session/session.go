// Package session bundles everything one decision-backed kb run needs for a
// topic: resolved topic, topic.yaml settings and active contract, the
// decision engine and generation client sharing one budget and receipts log,
// the state store, the owned-key writer, a lazily loaded corpus, and the
// decision modes (body links, relevance gates, quality gates).
package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/contract"
	"github.com/compozy/kb/internal/corpus"
	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/generation"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/questions"
	"github.com/compozy/kb/internal/topic"
)

// Decision mode values.
const (
	ModeShadow = "shadow"
	ModeApply  = "apply"
)

// CalibrationFile is the per-topic calibration record written by
// `kb review calibrate --write` (spec §12.3).
const CalibrationFile = "calibration.json"

// MinCalibrationLabels is the label floor behind an automatic relevance
// `apply` (spec §7 default mode).
const MinCalibrationLabels = 30

// ErrDecisionsRequired is returned when the decision model is not configured.
var ErrDecisionsRequired = errors.New("decision model is required")

// Flags are the per-run overrides shared by decision-backed commands.
type Flags struct {
	// BudgetUSD overrides [decisions].budget_usd when > 0.
	BudgetUSD float64
	// Decisions overrides the decision mode for this run: "", shadow or apply.
	Decisions string
}

// Options configures Open.
type Options struct {
	Config    config.Config
	VaultPath string
	Topic     string
	// Command names the kb command for error messages ("kb classify").
	Command    string
	Flags      Flags
	HTTPClient *http.Client
	Now        func() time.Time
	Logger     *slog.Logger
}

// Session is one decision-backed run over one topic.
type Session struct {
	Config    config.Config
	VaultPath string
	Topic     models.TopicInfo
	Settings  contract.Settings
	// Contract is the active contract, nil when none is accepted.
	Contract *contract.Contract
	Ref      decisions.TopicRef
	Engine   *decisions.Engine
	Gen      *generation.Client
	State    *corpus.StateStore
	Writer   *corpus.Writer
	Now      func() time.Time
	Logger   *slog.Logger
	Flags    Flags

	corpusMu sync.Mutex
	corpus   *corpus.Corpus
	// extras are the topic's extra question banks (all purposes).
	extras []*questions.Bank
}

// Open resolves the topic and builds a session. It fails fast, naming what is
// missing, when the decision model is not configured (spec §2 principle 2).
func Open(opts Options) (*Session, error) {
	command := strings.TrimSpace(opts.Command)
	if command == "" {
		command = "kb"
	}
	if strings.TrimSpace(opts.Config.OpenRouter.APIKey) == "" {
		return nil, fmt.Errorf(
			"%s: %w: OPENROUTER_API_KEY is not set (configure [openrouter].api_key or the OPENROUTER_API_KEY env var; [decisions] model %q)",
			command, ErrDecisionsRequired, opts.Config.Decisions.Model,
		)
	}
	switch opts.Flags.Decisions {
	case "", ModeShadow, ModeApply:
	default:
		return nil, fmt.Errorf("%s: --decisions must be shadow or apply: %q", command, opts.Flags.Decisions)
	}

	info, err := topic.Resolve(opts.VaultPath, opts.Topic)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", command, err)
	}
	settings, err := contract.LoadSettings(info.RootPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", command, err)
	}

	now := opts.Now
	if now == nil {
		now = time.Now
	}
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	decisionsConfig := opts.Config.Decisions
	var budget *decisions.Budget
	if opts.Flags.BudgetUSD > 0 {
		budget = decisions.NewBudget(opts.Flags.BudgetUSD)
	}
	engine, err := decisions.New(decisions.Options{
		Config:     decisionsConfig,
		APIKey:     opts.Config.OpenRouter.APIKey,
		APIURL:     opts.Config.OpenRouter.APIURL,
		HTTPClient: opts.HTTPClient,
		Budget:     budget,
		Now:        now,
		Logger:     logger,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", command, err)
	}
	gen, err := generation.New(generation.Options{
		Config:     opts.Config.Generation,
		APIKey:     opts.Config.OpenRouter.APIKey,
		APIURL:     opts.Config.OpenRouter.APIURL,
		HTTPClient: opts.HTTPClient,
		Budget:     engine.Budget(),
		Receipts:   engine.Receipts(),
		Now:        now,
		Logger:     logger,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", command, err)
	}

	store, err := corpus.OpenState(info.RootPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", command, err)
	}
	extras, err := questions.LoadTopicExtras(info.RootPath)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid topic question bank under %s: %w", command,
			filepath.Join(info.RootPath, decisions.ReceiptsDir, "banks"), err)
	}

	var active *contract.Contract
	if settings.Accepted() {
		normalized := settings.Contract.Normalized()
		active = &normalized
	}
	thresholds := decisions.Thresholds(nil).Merge(decisionsConfig.Thresholds).Merge(settings.Decisions.Thresholds)

	return &Session{
		Config:    opts.Config,
		VaultPath: opts.VaultPath,
		Topic:     info,
		Settings:  settings,
		Contract:  active,
		Ref: decisions.TopicRef{
			Slug:       info.Slug,
			Root:       info.RootPath,
			Contract:   active.Hash(),
			Thresholds: thresholds,
			Exclude:    settings.Decisions.Exclude,
		},
		Engine: engine,
		Gen:    gen,
		State:  store,
		Writer: corpus.NewWriter(info.RootPath, store, now),
		Now:    now,
		Logger: logger,
		Flags:  opts.Flags,
		extras: extras,
	}, nil
}

// ExtraBanks returns the topic's extra question banks for purpose, loaded
// from <topic>/.decisions/banks/*.json at Open (spec §4.3). Their questions
// are added to built-in requests of that purpose; their answers are only
// recorded in receipts and never replace a built-in answer.
func (s *Session) ExtraBanks(purpose decisions.Purpose) []*questions.Bank {
	out := make([]*questions.Bank, 0)
	for _, bank := range s.extras {
		if bank.Purpose == string(purpose) {
			out = append(out, bank)
		}
	}
	return out
}

// WithExtras adds the topic's extra questions of purpose to a built-in
// request: it returns the bank to send (the built-in bank composed with the
// extras, or the built-in bank unchanged without extras, so the cache key
// only changes when a topic adds questions) and qs followed by the extra
// questions (templates instantiated once per element id, see
// questions.Instantiate).
func (s *Session) WithExtras(purpose decisions.Purpose, bank *questions.Bank, qs []questions.Q, elements ...string) (*questions.Bank, []questions.Q, error) {
	extraBank, extra, err := questions.Instantiate(s.ExtraBanks(purpose), elements)
	if err != nil {
		return bank, qs, fmt.Errorf("topic %s extra %s questions: %w", s.Topic.Slug, purpose, err)
	}
	if extraBank == nil {
		return bank, qs, nil
	}
	return questions.Compose(bank, extraBank), append(slices.Clone(qs), extra...), nil
}

// Root returns the topic root path.
func (s *Session) Root() string { return s.Topic.RootPath }

// ContractHash returns the active contract hash ("" without a contract).
func (s *Session) ContractHash() string { return s.Ref.Contract }

// Threshold returns a named threshold for this topic.
func (s *Session) Threshold(name string) float64 { return s.Ref.Thresholds.Get(name) }

// Corpus loads the topic corpus once per session (decisions.exclude applied).
func (s *Session) Corpus() (*corpus.Corpus, error) {
	s.corpusMu.Lock()
	defer s.corpusMu.Unlock()
	if s.corpus != nil {
		return s.corpus, nil
	}
	loaded, err := corpus.Load(s.Root(), corpus.LoadOptions{Exclude: s.Settings.Decisions.Exclude})
	if err != nil {
		return nil, err
	}
	s.corpus = loaded
	return loaded, nil
}

// ReloadCorpus drops the cached corpus so the next Corpus call re-reads disk.
func (s *Session) ReloadCorpus() {
	s.corpusMu.Lock()
	s.corpus = nil
	s.corpusMu.Unlock()
}

// Excluded reports whether a topic-relative path matches decisions.exclude
// (the same match the engine and the generation client apply to every
// request Subject).
func (s *Session) Excluded(topicRel string) bool {
	return s.Ref.Excluded(topicRel)
}

// BodyMode is the body-link insertion mode: [decisions].mode, then
// topic.yaml decisions.mode, then --decisions.
func (s *Session) BodyMode() string {
	mode := normalizeMode(s.Config.Decisions.Mode, ModeShadow)
	mode = normalizeMode(s.Settings.Decisions.Mode, mode)
	return normalizeMode(s.Flags.Decisions, mode)
}

// RelevanceEnabled is false when topic.yaml sets decisions.relevance: off.
func (s *Session) RelevanceEnabled() bool {
	return !strings.EqualFold(strings.TrimSpace(s.Settings.Decisions.Relevance), contract.RelevanceOff)
}

// RelevanceGateMode is shadow until the topic has an accepted contract and
// either a stored calibration with ≥30 relevance labels or decisions.gates:
// apply (spec §7). --decisions overrides it for the run, except that
// --decisions=apply never enables relevance quarantine or skips without an
// accepted contract or with relevance off (spec §5.1: relevance then runs in
// shadow).
func (s *Session) RelevanceGateMode() string {
	if s.Flags.Decisions == ModeShadow {
		return ModeShadow
	}
	if !s.RelevanceEnabled() || s.Contract == nil {
		return ModeShadow
	}
	if s.Flags.Decisions == ModeApply {
		return ModeApply
	}
	if strings.EqualFold(s.Settings.Decisions.Gates, ModeApply) {
		return ModeApply
	}
	if strings.EqualFold(s.Settings.Decisions.Gates, ModeShadow) {
		return ModeShadow
	}
	if s.RelevanceCalibrated() {
		return ModeApply
	}
	return ModeShadow
}

// QualityGateMode is apply unless --decisions=shadow.
func (s *Session) QualityGateMode() string {
	if s.Flags.Decisions == ModeShadow {
		return ModeShadow
	}
	return ModeApply
}

// GateModeReason explains the relevance gate mode for run summaries.
func (s *Session) GateModeReason() string {
	noContract := "no accepted contract: run `kb topic contract " + s.Topic.Slug + " --draft|--import-claude` then `--accept`"
	switch {
	case s.Flags.Decisions == ModeShadow:
		return "set by --decisions"
	case !s.RelevanceEnabled() && s.Flags.Decisions == ModeApply:
		return "--decisions=apply ignored: relevance off in topic.yaml"
	case !s.RelevanceEnabled():
		return "relevance off in topic.yaml"
	case s.Contract == nil && s.Flags.Decisions == ModeApply:
		return "--decisions=apply ignored for relevance: " + noContract
	case s.Contract == nil:
		return noContract
	case s.Flags.Decisions == ModeApply:
		return "set by --decisions"
	case strings.EqualFold(s.Settings.Decisions.Gates, ModeApply):
		return "decisions.gates: apply in topic.yaml"
	case strings.EqualFold(s.Settings.Decisions.Gates, ModeShadow):
		return "decisions.gates: shadow in topic.yaml"
	case s.RelevanceCalibrated():
		return "calibrated"
	default:
		return fmt.Sprintf("not calibrated: needs ≥%d relevance labels and `kb review calibrate --write`", MinCalibrationLabels)
	}
}

// RelevanceCalibrated reports whether calibration.json records ≥30
// relevance labels given in review (dev + holdout minus imported): imported
// labels never are the only evidence behind an automatic apply (spec §20).
func (s *Session) RelevanceCalibrated() bool {
	record, err := ReadCalibration(s.Root())
	if err != nil || record == nil {
		return false
	}
	entry, ok := record.Purposes[string(decisions.PurposeRelevance)]
	if !ok {
		return false
	}
	return entry.DevLabels+entry.HoldoutLabels-entry.ImportedLabels >= MinCalibrationLabels
}

// StateMeta builds the state-row metadata for a write made from answers of
// the given banks under the active contract.
func (s *Session) StateMeta(banks ...*questions.Bank) corpus.StateMeta {
	meta := corpus.StateMeta{Contract: s.ContractHash(), Banks: map[string]string{}}
	for _, bank := range banks {
		if bank != nil {
			meta.Banks[bank.ID] = bank.Version
		}
	}
	return meta
}

// Request builds a decisions request for this topic.
func (s *Session) Request(purpose decisions.Purpose, subject string, bank *questions.Bank, state any, qs []questions.Q) decisions.Request {
	return decisions.Request{Topic: s.Ref, Purpose: purpose, Subject: subject, Bank: bank, State: state, Questions: qs}
}

// Summary returns the combined decisions + generation run summary.
func (s *Session) Summary() decisions.Summary {
	summary := s.Engine.Summary()
	summary.Generation = s.Gen.Summary()
	return summary
}

// SummaryLines renders the run summary, prefixed with budget state.
func (s *Session) SummaryLines() []string {
	lines := s.Summary().Lines()
	if budget := s.Engine.Budget(); budget != nil && budget.Exceeded() {
		lines = append(lines, fmt.Sprintf("budget: US$ %.4f limit reached; re-run to resume through the cache", budget.Limit()))
	}
	return lines
}

// WriteSummary prints the run summary to w.
func (s *Session) WriteSummary(w io.Writer) {
	for _, line := range s.SummaryLines() {
		_, _ = fmt.Fprintln(w, line)
	}
}

// PrivacyNotice prints, once per topic, where judged content is sent (spec
// §14). The marker lives under .decisions/.
func (s *Session) PrivacyNotice(w io.Writer) {
	marker := filepath.Join(s.Root(), decisions.ReceiptsDir, "privacy-notice")
	if _, err := os.Stat(marker); err == nil {
		return
	}
	_, _ = fmt.Fprintf(w,
		"privacy: judged content from %q is sent to OpenRouter. Decisions (%s) use a no-retention route; generation (%s, fallback %s) follows its provider's data policy — check it at https://openrouter.ai/%s. Use topic.yaml decisions.exclude to keep files out of every call.\n",
		s.Topic.Slug, s.Engine.Model(), s.Gen.Model(), s.Config.Generation.FallbackModel, s.Gen.Model(),
	)
	if err := os.MkdirAll(filepath.Dir(marker), 0o755); err == nil {
		_ = os.WriteFile(marker, []byte(s.Now().UTC().Format(time.RFC3339)+"\n"), 0o644)
	}
}

func normalizeMode(value, fallback string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case ModeShadow:
		return ModeShadow
	case ModeApply:
		return ModeApply
	default:
		return fallback
	}
}

// Calibration is the calibration.json record (spec §12.3).
type Calibration struct {
	Time     string                        `json:"time"`
	Purposes map[string]CalibrationPurpose `json:"purposes"`
}

// CalibrationPurpose is the stored calibration of one purpose.
type CalibrationPurpose struct {
	DevLabels      int                `json:"dev_labels"`
	HoldoutLabels  int                `json:"holdout_labels"`
	ImportedLabels int                `json:"imported_labels,omitempty"`
	Thresholds     map[string]float64 `json:"thresholds"`
	Dev            CalibrationMetrics `json:"dev"`
	Holdout        CalibrationMetrics `json:"holdout"`
}

// CalibrationMetrics are precision/recall/coverage/review rate at the chosen
// thresholds.
type CalibrationMetrics struct {
	Precision  float64 `json:"precision"`
	Recall     float64 `json:"recall"`
	Coverage   float64 `json:"coverage"`
	ReviewRate float64 `json:"review_rate"`
	N          int     `json:"n"`
}

// CalibrationPath returns <topic>/.decisions/calibration.json.
func CalibrationPath(topicRoot string) string {
	return filepath.Join(topicRoot, decisions.ReceiptsDir, CalibrationFile)
}

// ReadCalibration reads calibration.json; a missing file returns nil, nil.
func ReadCalibration(topicRoot string) (*Calibration, error) {
	data, err := os.ReadFile(CalibrationPath(topicRoot))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read calibration: %w", err)
	}
	var record Calibration
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("parse calibration: %w", err)
	}
	return &record, nil
}
