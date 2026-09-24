// Package contract owns a topic's selection contract (spec §5.1): the
// five-part scope that relevance and classification questions are judged
// against. It parses, validates, hashes and renders the contract, imports it
// from an existing CLAUDE.md section, and reads and writes the topic.yaml
// settings that carry it (contract, drafts, vocabulary draft, decisions block)
// without disturbing any other key.
package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Contract is the selection contract stored under topic.yaml `contract:`.
type Contract struct {
	Purpose                 string   `yaml:"purpose" json:"purpose"`
	Core                    []string `yaml:"core" json:"core"`
	Adjacent                []string `yaml:"adjacent" json:"adjacent"`
	CollectedOnPurpose      []string `yaml:"collected_on_purpose" json:"collected_on_purpose"`
	CollectedOnPurposePaths []string `yaml:"collected_on_purpose_paths,omitempty" json:"collected_on_purpose_paths,omitempty"`
	OutOfScope              []string `yaml:"out_of_scope" json:"out_of_scope"`
}

// Normalized returns a copy with every string trimmed, empty list items
// dropped and every list non-nil (so a YAML render always shows all fields).
// A nil contract normalizes to the empty contract.
func (c *Contract) Normalized() Contract {
	if c == nil {
		return Contract{Core: []string{}, Adjacent: []string{}, CollectedOnPurpose: []string{}, OutOfScope: []string{}}
	}
	normalized := Contract{
		Purpose:            strings.TrimSpace(c.Purpose),
		Core:               cleanList(c.Core),
		Adjacent:           cleanList(c.Adjacent),
		CollectedOnPurpose: cleanList(c.CollectedOnPurpose),
		OutOfScope:         cleanList(c.OutOfScope),
	}
	if paths := cleanList(c.CollectedOnPurposePaths); len(paths) > 0 {
		normalized.CollectedOnPurposePaths = paths
	}
	return normalized
}

// Empty reports whether the contract carries no text at all.
func (c *Contract) Empty() bool {
	if c == nil {
		return true
	}
	normalized := c.Normalized()
	return normalized.Purpose == "" &&
		len(normalized.Core) == 0 &&
		len(normalized.Adjacent) == 0 &&
		len(normalized.CollectedOnPurpose) == 0 &&
		len(normalized.CollectedOnPurposePaths) == 0 &&
		len(normalized.OutOfScope) == 0
}

// Validate checks the rules of spec §5.1: a non-empty purpose and core, no
// line shared (trimmed, case-insensitive) between core/adjacent and
// out_of_scope, and collected_on_purpose_paths globs that compile.
func (c *Contract) Validate() error {
	if c == nil {
		return errors.New("contract: contract is empty")
	}
	normalized := c.Normalized()
	var problems []string
	if normalized.Purpose == "" {
		problems = append(problems, "purpose is required")
	}
	if len(normalized.Core) == 0 {
		problems = append(problems, "core needs at least one line")
	}
	excluded := make(map[string]struct{}, len(normalized.OutOfScope))
	for _, line := range normalized.OutOfScope {
		excluded[strings.ToLower(line)] = struct{}{}
	}
	for _, group := range []struct {
		name  string
		lines []string
	}{{"core", normalized.Core}, {"adjacent", normalized.Adjacent}} {
		for _, line := range group.lines {
			if _, clash := excluded[strings.ToLower(line)]; clash {
				problems = append(problems, fmt.Sprintf("%q appears in both %s and out_of_scope", line, group.name))
			}
		}
	}
	for _, pattern := range normalized.CollectedOnPurposePaths {
		if err := ValidatePattern(pattern); err != nil {
			problems = append(problems, fmt.Sprintf("collected_on_purpose_paths: %v", err))
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("contract: %s", strings.Join(problems, "; "))
	}
	return nil
}

// Hash is the sha256 hex digest of the canonical JSON of the normalized
// contract. It is "" for a nil or empty contract.
func (c *Contract) Hash() string {
	if c.Empty() {
		return ""
	}
	normalized := c.Normalized()
	// A struct of strings and string slices always marshals; field order is
	// fixed by the struct, which makes the encoding canonical.
	encoded, _ := json.Marshal(normalized)
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

// MatchesCollectedPath reports whether a topic-relative path matches one of
// the collected_on_purpose_paths globs.
func (c *Contract) MatchesCollectedPath(topicRelPath string) bool {
	if c == nil {
		return false
	}
	for _, pattern := range c.CollectedOnPurposePaths {
		if MatchPath(strings.TrimSpace(pattern), topicRelPath) {
			return true
		}
	}
	return false
}

// QuestionState is the contract fragment every relevance question carries.
// Path globs are resolved by code and never sent.
func (c *Contract) QuestionState() map[string]any {
	normalized := c.Normalized()
	return map[string]any{
		"purpose":              normalized.Purpose,
		"core":                 normalized.Core,
		"adjacent":             normalized.Adjacent,
		"collected_on_purpose": normalized.CollectedOnPurpose,
		"out_of_scope":         normalized.OutOfScope,
	}
}

func cleanList(values []string) []string {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}
