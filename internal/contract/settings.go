package contract

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/compozy/kb/internal/contract/topicyaml"
)

// topic.yaml keys owned by this package.
const (
	keyContract        = "contract"
	keyContractDraft   = "contract_draft"
	keyVocabularyDraft = "vocabulary_draft"
	keyDecisions       = "decisions"
	keyThresholds      = "thresholds"
)

// Decision mode values accepted in topic.yaml `decisions:`.
const (
	ModeShadow   = "shadow"
	ModeApply    = "apply"
	RelevanceOn  = "on"
	RelevanceOff = "off"
)

// ErrNoDraft is returned by ActivateDraft when topic.yaml has no contract_draft.
var ErrNoDraft = errors.New("contract: topic.yaml has no contract_draft; run `kb topic contract <topic> --draft` or `--import-claude` first")

// VocabularyItem is one proposed concept of topic.yaml `vocabulary_draft:`.
type VocabularyItem struct {
	Title     string `yaml:"title" json:"title"`
	Criterion string `yaml:"criterion" json:"criterion"`
	Mentions  int    `yaml:"mentions,omitempty" json:"mentions,omitempty"`
}

// DecisionSettings is topic.yaml `decisions:`.
type DecisionSettings struct {
	// Mode is the body-link insertion mode: "" (inherit), shadow or apply.
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// Gates is the relevance gate mode: "" (inherit), shadow or apply.
	Gates string `yaml:"gates,omitempty" json:"gates,omitempty"`
	// Relevance is "" (on), on or off.
	Relevance string `yaml:"relevance,omitempty" json:"relevance,omitempty"`
	// Thresholds override [decisions.thresholds] by name.
	Thresholds map[string]float64 `yaml:"thresholds,omitempty" json:"thresholds,omitempty"`
	// Exclude globs keep matching topic-relative files out of every call.
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

// Settings is the part of topic.yaml read by the decision pipeline.
type Settings struct {
	Contract        *Contract        `yaml:"contract,omitempty" json:"contract,omitempty"`
	ContractDraft   *Contract        `yaml:"contract_draft,omitempty" json:"contract_draft,omitempty"`
	VocabularyDraft []VocabularyItem `yaml:"vocabulary_draft,omitempty" json:"vocabulary_draft,omitempty"`
	Decisions       DecisionSettings `yaml:"decisions,omitempty" json:"decisions"`
}

// Accepted reports whether the topic has an active, valid contract.
func (s Settings) Accepted() bool {
	return !s.Contract.Empty() && s.Contract.Validate() == nil
}

// SettingsPath returns <topicRoot>/topic.yaml.
func SettingsPath(topicRoot string) string {
	return filepath.Join(topicRoot, topicyaml.FileName)
}

// LoadSettings reads the decision settings from <topicRoot>/topic.yaml. A
// missing file or missing keys yield zero values; enum values, thresholds and
// exclude globs are validated.
func LoadSettings(topicRoot string) (Settings, error) {
	path := SettingsPath(topicRoot)
	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Settings{}, nil
	}
	if err != nil {
		return Settings{}, fmt.Errorf("contract: read %q: %w", path, err)
	}
	var settings Settings
	if err := yaml.Unmarshal(content, &settings); err != nil {
		return Settings{}, fmt.Errorf("contract: parse %q: %w", path, err)
	}
	if settings.Contract != nil {
		normalized := settings.Contract.Normalized()
		settings.Contract = &normalized
	}
	if settings.ContractDraft != nil {
		normalized := settings.ContractDraft.Normalized()
		settings.ContractDraft = &normalized
	}
	settings.VocabularyDraft = cleanVocabulary(settings.VocabularyDraft)
	if err := settings.Decisions.normalize(); err != nil {
		return Settings{}, fmt.Errorf("contract: %q: %w", path, err)
	}
	return settings, nil
}

func (d *DecisionSettings) normalize() error {
	d.Mode = strings.ToLower(strings.TrimSpace(d.Mode))
	d.Gates = strings.ToLower(strings.TrimSpace(d.Gates))
	d.Relevance = strings.ToLower(strings.TrimSpace(d.Relevance))
	if err := checkEnum("decisions.mode", d.Mode, ModeShadow, ModeApply); err != nil {
		return err
	}
	if err := checkEnum("decisions.gates", d.Gates, ModeShadow, ModeApply); err != nil {
		return err
	}
	if err := checkEnum("decisions.relevance", d.Relevance, RelevanceOn, RelevanceOff); err != nil {
		return err
	}
	if err := validateThresholds(d.Thresholds); err != nil {
		return err
	}
	d.Exclude = cleanList(d.Exclude)
	for _, pattern := range d.Exclude {
		if err := ValidatePattern(pattern); err != nil {
			return fmt.Errorf("decisions.exclude: %w", err)
		}
	}
	return nil
}

func checkEnum(name, value string, allowed ...string) error {
	if value == "" {
		return nil
	}
	if slices.Contains(allowed, value) {
		return nil
	}
	return fmt.Errorf("%s must be one of %s, got %q", name, strings.Join(allowed, "|"), value)
}

func validateThresholds(thresholds map[string]float64) error {
	for _, name := range slices.Sorted(maps.Keys(thresholds)) {
		if strings.TrimSpace(name) == "" {
			return errors.New("decisions.thresholds: threshold name cannot be empty")
		}
		if value := thresholds[name]; value < 0 || value > 1 {
			return fmt.Errorf("decisions.thresholds.%s must be within [0, 1], got %v", name, value)
		}
	}
	return nil
}

func cleanVocabulary(items []VocabularyItem) []VocabularyItem {
	if len(items) == 0 {
		return nil
	}
	cleaned := make([]VocabularyItem, 0, len(items))
	for _, item := range items {
		item.Title = strings.TrimSpace(item.Title)
		item.Criterion = strings.TrimSpace(item.Criterion)
		if item.Title == "" {
			continue
		}
		cleaned = append(cleaned, item)
	}
	return cleaned
}

// SaveContractDraft writes c under topic.yaml `contract_draft:` (drafts are
// not validated; ActivateDraft validates).
func SaveContractDraft(topicRoot string, c *Contract) error {
	if c == nil {
		return errors.New("contract: draft is nil")
	}
	normalized := c.Normalized()
	return updateSettings(topicRoot, func(root *yaml.Node) error {
		return topicyaml.Set(root, keyContractDraft, normalized)
	})
}

// SetContract validates c and writes it as the active topic.yaml `contract:`.
func SetContract(topicRoot string, c *Contract) error {
	if err := c.Validate(); err != nil {
		return err
	}
	normalized := c.Normalized()
	return updateSettings(topicRoot, func(root *yaml.Node) error {
		return topicyaml.Set(root, keyContract, normalized)
	})
}

// ActivateDraft validates topic.yaml `contract_draft:`, moves it to
// `contract:` and removes the draft key. Nothing is written when the draft is
// missing or invalid. Callers re-render CLAUDE.md with SyncClaude afterwards.
func ActivateDraft(topicRoot string) (*Contract, error) {
	var activated Contract
	err := updateSettings(topicRoot, func(root *yaml.Node) error {
		draftNode := topicyaml.Get(root, keyContractDraft)
		if draftNode == nil || (draftNode.Kind == yaml.ScalarNode && draftNode.Tag == "!!null") {
			return ErrNoDraft
		}
		var draft Contract
		if err := draftNode.Decode(&draft); err != nil {
			return fmt.Errorf("contract: decode contract_draft: %w", err)
		}
		if err := draft.Validate(); err != nil {
			return err
		}
		activated = draft.Normalized()
		if err := topicyaml.Set(root, keyContract, activated); err != nil {
			return err
		}
		topicyaml.Delete(root, keyContractDraft)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &activated, nil
}

// SaveVocabularyDraft writes items under topic.yaml `vocabulary_draft:`.
func SaveVocabularyDraft(topicRoot string, items []VocabularyItem) error {
	cleaned := cleanVocabulary(items)
	if len(cleaned) == 0 {
		return errors.New("contract: vocabulary draft is empty")
	}
	return updateSettings(topicRoot, func(root *yaml.Node) error {
		return topicyaml.Set(root, keyVocabularyDraft, cleaned)
	})
}

// ClearVocabularyDraft removes topic.yaml `vocabulary_draft:` (no-op when
// absent).
func ClearVocabularyDraft(topicRoot string) error {
	return updateSettings(topicRoot, func(root *yaml.Node) error {
		topicyaml.Delete(root, keyVocabularyDraft)
		return nil
	})
}

// SetThresholds replaces topic.yaml `decisions.thresholds:`, keeping every
// other `decisions` key. An empty map removes the thresholds key.
func SetThresholds(topicRoot string, thresholds map[string]float64) error {
	if err := validateThresholds(thresholds); err != nil {
		return fmt.Errorf("contract: %w", err)
	}
	return updateSettings(topicRoot, func(root *yaml.Node) error {
		decisions, err := topicyaml.Mapping(root, keyDecisions)
		if err != nil {
			return err
		}
		if len(thresholds) == 0 {
			topicyaml.Delete(decisions, keyThresholds)
			return nil
		}
		return topicyaml.Set(decisions, keyThresholds, thresholds)
	})
}

func updateSettings(topicRoot string, edit func(root *yaml.Node) error) error {
	if strings.TrimSpace(topicRoot) == "" {
		return errors.New("contract: topic root is required")
	}
	path := SettingsPath(topicRoot)
	document, err := topicyaml.Load(path)
	if err != nil {
		return fmt.Errorf("contract: %w", err)
	}
	if err := edit(document.Root()); err != nil {
		return err
	}
	if err := document.Save(path); err != nil {
		return fmt.Errorf("contract: %w", err)
	}
	return nil
}
