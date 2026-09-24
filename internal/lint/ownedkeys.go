package lint

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/compozy/kb/internal/decisions"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/questions"
)

// Stages accepted for wiki concept articles: compiled articles written by the
// agent and stub articles created by `kb topic vocabulary --accept` or an
// accepted concept proposal (spec §5.2).
const (
	stageCompiled = "compiled"
	stageStub     = "stub"
)

// Owned-key limits (spec §6).
const (
	maxSummaryChars   = 400
	maxCriterionChars = 300
	maxEntities       = 15
	minQuestions      = 3
	maxQuestions      = 6
	maxDepth          = 3.0
)

var quotedWikilinkPattern = regexp.MustCompile(`^\[\[[^\[\]\n]+\]\]$`)

var (
	triageValues       = []string{"kept", "review", "quarantined"}
	triageReasonValues = []string{"off_topic", "paywall", "error_page", "thin", "not_an_article", "no_speech", "duplicate"}
	qualityValues      = []string{"thin", "error_page", "paywall", "not_an_article", "no_speech"}
)

// wikilinkListKeys hold lists of quoted wikilinks (spec §6).
var wikilinkListKeys = append([]string{"concepts", "affects", "supersedes"}, decisions.RelationKeys...)

// stageVariant adapts the wiki concept schema to the article's stage: a stub
// article is valid with `stage: stub` and an empty `sources` list; any other
// stage must be `compiled`.
func stageVariant(spec schemaSpec, file *vaultFile) schemaSpec {
	if !isWikiConceptPath(file.relativePath) || file.parseErr != nil {
		return spec
	}
	if strings.TrimSpace(stringValue(file.frontmatter["stage"])) != stageStub {
		return spec
	}

	variant := spec
	variant.expected = maps.Clone(spec.expected)
	variant.expected["stage"] = stageStub
	variant.required = slices.DeleteFunc(slices.Clone(spec.required), func(field string) bool {
		return field == "sources"
	})
	return variant
}

// ownedKeyIssues validates the shapes of kb-owned keys present in a source or
// article (spec §6). Keys that belong to the user (key-conflict) are skipped:
// kb's schema does not apply to them. Type and enum violations are errors;
// size limits (summary, criterion, entities, questions) are warnings.
func ownedKeyIssues(file *vaultFile, userKeys map[string]struct{}) []models.LintIssue {
	issues := make([]models.LintIssue, 0)
	report := func(severity models.DiagnosticSeverity, key, message string) {
		issues = append(issues, newIssue(models.LintIssueKindFormat, severity, file.relativePath, message, key))
	}

	for _, key := range sortedKeys(file.frontmatter) {
		if _, isUser := userKeys[key]; isUser {
			continue
		}
		value := file.frontmatter[key]
		if value == nil {
			continue
		}
		switch {
		case key == "locked":
			if _, ok := value.(bool); !ok {
				report(models.SeverityError, key, `frontmatter field "locked" must be true or false`)
			}
		case key == "triage":
			checkEnum(report, key, value, triageValues)
		case key == "triage_reason":
			checkEnum(report, key, value, triageReasonValues)
		case key == "quality":
			checkEnum(report, key, value, qualityValues)
		case key == "genre":
			checkEnum(report, key, value, bankOptions("classify", "kind"))
		case key == "relevance":
			checkEnum(report, key, value, bankOptions("relevance", "role"))
		case key == "depth":
			depth, ok := numberValue(value)
			if !ok || depth < 0 || depth > maxDepth {
				report(models.SeverityError, key, `frontmatter field "depth" must be a number from 0 to 3`)
			}
		case key == "summary":
			checkText(report, key, value, maxSummaryChars)
		case key == "criterion":
			checkText(report, key, value, maxCriterionChars)
		case key == "ingest_batch", key == "ingest_query":
			if _, ok := value.(string); !ok {
				report(models.SeverityError, key, fmt.Sprintf("frontmatter field %q must be a string", key))
			}
		case key == "recaptured":
			if _, isTime := value.(time.Time); !isTime && !isDateString(value) {
				report(models.SeverityError, key, `frontmatter field "recaptured" must be a valid ISO date`)
			}
		case key == "aliases":
			if _, isString := value.(string); !isString && !isStringList(value) {
				report(models.SeverityError, key, `frontmatter field "aliases" must be a list of strings`)
			}
		case key == "entities":
			checkStringList(report, key, value, 0, maxEntities)
		case key == "questions":
			checkStringList(report, key, value, minQuestions, maxQuestions)
		case slices.Contains(wikilinkListKeys, key):
			if !isWikilinkList(value) {
				report(models.SeverityError, key, fmt.Sprintf(`frontmatter field %q must be a list of quoted wikilinks ("[[File name]]")`, key))
			}
		}
	}

	return issues
}

func checkEnum(report func(models.DiagnosticSeverity, string, string), key string, value any, allowed []string) {
	text, ok := value.(string)
	if ok && slices.Contains(allowed, strings.TrimSpace(text)) {
		return
	}
	report(models.SeverityError, key, fmt.Sprintf("frontmatter field %q must be one of %s", key, strings.Join(allowed, ", ")))
}

func checkText(report func(models.DiagnosticSeverity, string, string), key string, value any, limit int) {
	text, ok := value.(string)
	if !ok {
		report(models.SeverityError, key, fmt.Sprintf("frontmatter field %q must be a string", key))
		return
	}
	if count := utf8.RuneCountInString(strings.TrimSpace(text)); count > limit {
		report(models.SeverityWarning, key, fmt.Sprintf("frontmatter field %q has %d characters, limit is %d", key, count, limit))
	}
}

func checkStringList(report func(models.DiagnosticSeverity, string, string), key string, value any, minItems, maxItems int) {
	if !isStringList(value) {
		report(models.SeverityError, key, fmt.Sprintf("frontmatter field %q must be a list of strings", key))
		return
	}
	count := listLen(value)
	if count < minItems || count > maxItems {
		report(models.SeverityWarning, key, fmt.Sprintf("frontmatter field %q has %d items, expected %d to %d", key, count, minItems, maxItems))
	}
}

func isStringList(value any) bool {
	switch typed := value.(type) {
	case []string:
		return true
	case []any:
		for _, item := range typed {
			if _, ok := item.(string); !ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func isWikilinkList(value any) bool {
	if !isStringList(value) {
		return false
	}
	for _, item := range stringItems(value) {
		if !quotedWikilinkPattern.MatchString(strings.TrimSpace(item)) {
			return false
		}
	}
	return true
}

func stringItems(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				items = append(items, text)
			}
		}
		return items
	default:
		return nil
	}
}

func listLen(value any) int {
	switch typed := value.(type) {
	case []string:
		return len(typed)
	case []any:
		return len(typed)
	default:
		return 0
	}
}

func numberValue(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float64:
		return typed, true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func isDateString(value any) bool {
	text, ok := value.(string)
	if !ok {
		return false
	}
	text = strings.TrimSpace(text)
	for _, layout := range []string{"2006-01-02", time.RFC3339} {
		if _, err := time.Parse(layout, text); err == nil {
			return true
		}
	}
	return false
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}

// bankOptions returns the option names of a choice question of an embedded
// bank, so enum checks follow the bank the classifier asks.
func bankOptions(bankID, questionID string) []string {
	bank, err := questions.Load(bankID)
	if err != nil {
		return nil
	}
	question, err := bank.Question(questionID, nil)
	if err != nil {
		return nil
	}
	criteria, ok := question.Criteria.(map[string]any)
	if !ok {
		return nil
	}
	return slices.Sorted(maps.Keys(criteria))
}
