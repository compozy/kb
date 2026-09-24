package contract

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const userTopicYAML = `# Topic metadata, hand-edited.
slug: agent-swarm
title: Agent Swarm # inline comment
domain: agent-swarm
mode: wiki
qmd_collection: agent-swarm
owner:
  name: Pedro # who curates it
  tags: [a, b]
decisions:
  # keep links in shadow for now
  mode: shadow
  exclude:
    - raw/private/**
`

func writeTopicYAML(t *testing.T, content string) string {
	t.Helper()
	topicRoot := t.TempDir()
	writeTestFile(t, filepath.Join(topicRoot, "topic.yaml"), content)
	return topicRoot
}

func assertPreservesUserYAML(t *testing.T, topicRoot string) string {
	t.Helper()
	got := readTestFile(t, filepath.Join(topicRoot, "topic.yaml"))
	for _, fragment := range []string{
		"# Topic metadata, hand-edited.\nslug: agent-swarm\ntitle: Agent Swarm # inline comment\ndomain: agent-swarm\nmode: wiki\nqmd_collection: agent-swarm\n",
		"owner:\n  name: Pedro # who curates it\n  tags: [a, b]\n",
		"decisions:\n  # keep links in shadow for now\n  mode: shadow\n",
		"  exclude:\n    - raw/private/**\n",
	} {
		if !strings.Contains(got, fragment) {
			t.Fatalf("topic.yaml lost user content %q:\n%s", fragment, got)
		}
	}
	return got
}

func TestLoadSettings(t *testing.T) {
	t.Parallel()

	t.Run("Should return zero settings when topic.yaml is missing", func(t *testing.T) {
		t.Parallel()
		settings, err := LoadSettings(t.TempDir())
		if err != nil {
			t.Fatalf("LoadSettings returned error: %v", err)
		}
		if !reflect.DeepEqual(settings, Settings{}) || settings.Accepted() {
			t.Fatalf("settings = %#v, want zero and not accepted", settings)
		}
	})

	t.Run("Should read every key and ignore unrelated ones", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, `slug: x
contract:
  purpose: " Purpose "
  core: [Core line, " "]
  adjacent: []
  collected_on_purpose: []
  collected_on_purpose_paths: ["raw/audit/**"]
  out_of_scope: [Junk]
contract_draft:
  purpose: Draft
  core: [Draft core]
vocabulary_draft:
  - title: Swarm runtimes
    criterion: Discusses a swarm runtime.
    mentions: 4
  - title: " "
decisions:
  mode: Apply
  gates: shadow
  relevance: off
  thresholds:
    link_apply: 0.9
  exclude: ["raw/private/**"]
unrelated: {a: 1}
`)
		settings, err := LoadSettings(topicRoot)
		if err != nil {
			t.Fatalf("LoadSettings returned error: %v", err)
		}
		wantContract := &Contract{
			Purpose:                 "Purpose",
			Core:                    []string{"Core line"},
			Adjacent:                []string{},
			CollectedOnPurpose:      []string{},
			CollectedOnPurposePaths: []string{"raw/audit/**"},
			OutOfScope:              []string{"Junk"},
		}
		if !reflect.DeepEqual(settings.Contract, wantContract) {
			t.Fatalf("contract = %#v, want %#v", settings.Contract, wantContract)
		}
		if settings.ContractDraft == nil || settings.ContractDraft.Purpose != "Draft" {
			t.Fatalf("contract draft = %#v", settings.ContractDraft)
		}
		wantVocabulary := []VocabularyItem{{Title: "Swarm runtimes", Criterion: "Discusses a swarm runtime.", Mentions: 4}}
		if !reflect.DeepEqual(settings.VocabularyDraft, wantVocabulary) {
			t.Fatalf("vocabulary = %#v, want %#v", settings.VocabularyDraft, wantVocabulary)
		}
		wantDecisions := DecisionSettings{
			Mode:       "apply",
			Gates:      "shadow",
			Relevance:  "off",
			Thresholds: map[string]float64{"link_apply": 0.9},
			Exclude:    []string{"raw/private/**"},
		}
		if !reflect.DeepEqual(settings.Decisions, wantDecisions) {
			t.Fatalf("decisions = %#v, want %#v", settings.Decisions, wantDecisions)
		}
		if !settings.Accepted() {
			t.Fatal("Accepted() = false for a valid contract")
		}
	})

	t.Run("Should not accept the empty scaffold contract", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, "slug: x\ncontract:\n  purpose: \"\"\n  core: []\n  adjacent: []\n  collected_on_purpose: []\n  out_of_scope: []\n")
		settings, err := LoadSettings(topicRoot)
		if err != nil {
			t.Fatalf("LoadSettings returned error: %v", err)
		}
		if settings.Accepted() {
			t.Fatal("Accepted() = true for an empty contract")
		}
	})

	invalid := []struct {
		name    string
		yaml    string
		wantErr string
	}{
		{"Should reject an unknown decisions.mode", "decisions:\n  mode: live\n", "decisions.mode must be one of shadow|apply"},
		{"Should reject an unknown decisions.gates", "decisions:\n  gates: on\n", "decisions.gates must be one of shadow|apply"},
		{"Should reject an unknown decisions.relevance", "decisions:\n  relevance: maybe\n", "decisions.relevance must be one of on|off"},
		{"Should reject a threshold outside [0, 1]", "decisions:\n  thresholds:\n    link_apply: 1.5\n", "decisions.thresholds.link_apply must be within [0, 1]"},
		{"Should reject an exclude glob that does not compile", "decisions:\n  exclude: [\"raw/[x\"]\n", "decisions.exclude"},
		{"Should reject malformed YAML", "contract: [\n", "parse"},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := LoadSettings(writeTopicYAML(t, tt.yaml))
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("LoadSettings() = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestSettingsWritersPreserveTopicYAML(t *testing.T) {
	t.Parallel()

	vocabulary := []VocabularyItem{{Title: "Swarm runtimes", Criterion: "Discusses a swarm runtime.", Mentions: 3}}
	tests := []struct {
		name  string
		write func(topicRoot string) error
		check func(t *testing.T, settings Settings)
	}{
		{
			name:  "Should save a contract draft",
			write: func(topicRoot string) error { return SaveContractDraft(topicRoot, validContract()) },
			check: func(t *testing.T, settings Settings) {
				if settings.ContractDraft.Hash() != validContract().Hash() || settings.Contract != nil {
					t.Fatalf("settings = %#v", settings)
				}
			},
		},
		{
			name:  "Should set the active contract",
			write: func(topicRoot string) error { return SetContract(topicRoot, validContract()) },
			check: func(t *testing.T, settings Settings) {
				if settings.Contract.Hash() != validContract().Hash() || !settings.Accepted() {
					t.Fatalf("settings = %#v", settings)
				}
			},
		},
		{
			name:  "Should save a vocabulary draft",
			write: func(topicRoot string) error { return SaveVocabularyDraft(topicRoot, vocabulary) },
			check: func(t *testing.T, settings Settings) {
				if !reflect.DeepEqual(settings.VocabularyDraft, vocabulary) {
					t.Fatalf("vocabulary = %#v", settings.VocabularyDraft)
				}
			},
		},
		{
			name: "Should clear a vocabulary draft",
			write: func(topicRoot string) error {
				if err := SaveVocabularyDraft(topicRoot, vocabulary); err != nil {
					return err
				}
				return ClearVocabularyDraft(topicRoot)
			},
			check: func(t *testing.T, settings Settings) {
				if settings.VocabularyDraft != nil {
					t.Fatalf("vocabulary = %#v, want cleared", settings.VocabularyDraft)
				}
			},
		},
		{
			name: "Should set thresholds inside the existing decisions block",
			write: func(topicRoot string) error {
				return SetThresholds(topicRoot, map[string]float64{"link_review": 0.55, "link_apply": 0.9})
			},
			check: func(t *testing.T, settings Settings) {
				want := map[string]float64{"link_apply": 0.9, "link_review": 0.55}
				if !reflect.DeepEqual(settings.Decisions.Thresholds, want) || settings.Decisions.Mode != "shadow" {
					t.Fatalf("decisions = %#v", settings.Decisions)
				}
			},
		},
		{
			name: "Should activate a draft and drop the draft key",
			write: func(topicRoot string) error {
				if err := SaveContractDraft(topicRoot, validContract()); err != nil {
					return err
				}
				activated, err := ActivateDraft(topicRoot)
				if err != nil {
					return err
				}
				if activated.Hash() != validContract().Hash() {
					return errors.New("activated contract differs from the draft")
				}
				return nil
			},
			check: func(t *testing.T, settings Settings) {
				if settings.ContractDraft != nil || !settings.Accepted() {
					t.Fatalf("settings = %#v", settings)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			topicRoot := writeTopicYAML(t, userTopicYAML)
			if err := tt.write(topicRoot); err != nil {
				t.Fatalf("write returned error: %v", err)
			}
			assertPreservesUserYAML(t, topicRoot)
			settings, err := LoadSettings(topicRoot)
			if err != nil {
				t.Fatalf("LoadSettings returned error: %v", err)
			}
			tt.check(t, settings)
		})
	}
}

func TestSettingsWritersEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("Should replace an existing contract in place keeping key order", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, "slug: x\ncontract:\n  purpose: old\n  core: [old]\ntitle: X\n")
		if err := SetContract(topicRoot, validContract()); err != nil {
			t.Fatalf("SetContract returned error: %v", err)
		}
		got := readTestFile(t, filepath.Join(topicRoot, "topic.yaml"))
		if !strings.HasPrefix(got, "slug: x\ncontract:\n  purpose: How agents coordinate as a swarm.\n") || !strings.HasSuffix(got, "title: X\n") {
			t.Fatalf("topic.yaml order changed:\n%s", got)
		}
	})

	t.Run("Should refuse to set an invalid contract and leave the file untouched", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, userTopicYAML)
		err := SetContract(topicRoot, &Contract{Purpose: "p"})
		if err == nil || !strings.Contains(err.Error(), "core needs at least one line") {
			t.Fatalf("SetContract() = %v, want validation error", err)
		}
		if got := readTestFile(t, filepath.Join(topicRoot, "topic.yaml")); got != userTopicYAML {
			t.Fatalf("topic.yaml changed:\n%s", got)
		}
	})

	t.Run("Should fail activation without a draft", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, userTopicYAML)
		if _, err := ActivateDraft(topicRoot); !errors.Is(err, ErrNoDraft) {
			t.Fatalf("ActivateDraft() = %v, want ErrNoDraft", err)
		}
	})

	t.Run("Should keep an invalid draft and not activate it", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, userTopicYAML)
		draft := validContract()
		draft.OutOfScope = append(draft.OutOfScope, draft.Core[0])
		if err := SaveContractDraft(topicRoot, draft); err != nil {
			t.Fatalf("SaveContractDraft returned error: %v", err)
		}
		before := readTestFile(t, filepath.Join(topicRoot, "topic.yaml"))
		if _, err := ActivateDraft(topicRoot); err == nil || !strings.Contains(err.Error(), "appears in both core and out_of_scope") {
			t.Fatalf("ActivateDraft() = %v, want conflict error", err)
		}
		if after := readTestFile(t, filepath.Join(topicRoot, "topic.yaml")); after != before {
			t.Fatalf("topic.yaml changed after failed activation:\n%s", after)
		}
	})

	t.Run("Should create topic.yaml when missing and remove thresholds on an empty map", func(t *testing.T) {
		t.Parallel()
		topicRoot := t.TempDir()
		if err := SetThresholds(topicRoot, map[string]float64{"link_apply": 0.8}); err != nil {
			t.Fatalf("SetThresholds returned error: %v", err)
		}
		if got := readTestFile(t, filepath.Join(topicRoot, "topic.yaml")); got != "decisions:\n  thresholds:\n    link_apply: 0.8\n" {
			t.Fatalf("topic.yaml =\n%q", got)
		}
		if err := SetThresholds(topicRoot, nil); err != nil {
			t.Fatalf("SetThresholds(nil) returned error: %v", err)
		}
		if got := readTestFile(t, filepath.Join(topicRoot, "topic.yaml")); got != "decisions: {}\n" {
			t.Fatalf("topic.yaml =\n%q", got)
		}
	})

	t.Run("Should reject invalid writer inputs", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, userTopicYAML)
		if err := SetThresholds(topicRoot, map[string]float64{"link_apply": -0.1}); err == nil {
			t.Fatal("SetThresholds accepted a negative threshold")
		}
		if err := SaveVocabularyDraft(topicRoot, []VocabularyItem{{Title: " "}}); err == nil {
			t.Fatal("SaveVocabularyDraft accepted an empty draft")
		}
		if err := SaveContractDraft(topicRoot, nil); err == nil {
			t.Fatal("SaveContractDraft accepted a nil draft")
		}
		if got := readTestFile(t, filepath.Join(topicRoot, "topic.yaml")); got != userTopicYAML {
			t.Fatalf("topic.yaml changed:\n%s", got)
		}
	})

	t.Run("Should refuse a topic.yaml whose root is not a mapping", func(t *testing.T) {
		t.Parallel()
		topicRoot := writeTopicYAML(t, "- a\n- b\n")
		if err := SetContract(topicRoot, validContract()); err == nil || !strings.Contains(err.Error(), "root must be a mapping") {
			t.Fatalf("SetContract() = %v, want mapping error", err)
		}
	})
}
