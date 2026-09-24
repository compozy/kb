package lint

import (
	"fmt"
	"sort"
	"strings"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/resolve"
)

// addFrontmatterLinks resolves every wikilink found in frontmatter values
// (spec §6: frontmatter links are part of the graph). A link that resolves
// nowhere is a frontmatter-dead-link, a link to a quarantined source is
// link-to-quarantined, and a link that resolves inside the topic counts as an
// incoming edge for orphan detection. `sources` on wiki articles is already
// checked by missing-source, so only its quarantined targets are reported
// here.
func addFrontmatterLinks(state vaultState, incoming map[string]map[string]struct{}) []models.LintIssue {
	issues := make([]models.LintIssue, 0)
	seen := make(map[string]struct{})

	for _, file := range state.allFiles {
		if file.parseErr != nil {
			continue
		}
		links := resolve.FrontmatterLinks(file.frontmatter)
		for _, key := range sortedKeys(links) {
			for _, target := range links[key] {
				resolved := state.resolveTarget(target, false)
				switch {
				case resolved == nil:
					if !file.inTopic || (key == "sources" && isWikiConceptPath(file.relativePath)) {
						continue
					}
					dedupe := "fm\x00" + file.relativePath + "\x00" + key + "\x00" + target
					if _, exists := seen[dedupe]; exists {
						continue
					}
					seen[dedupe] = struct{}{}
					issues = append(issues, newIssue(
						models.LintIssueKindFrontmatterDeadLink,
						models.SeverityError,
						file.relativePath,
						fmt.Sprintf("frontmatter %q link target does not exist", key),
						target,
					))
				case isQuarantinedFile(resolved):
					if file.inTopic {
						issues = appendQuarantinedLink(issues, seen, file, target, key)
					}
				case resolved.inTopic:
					if incoming[resolved.relativePath] == nil {
						incoming[resolved.relativePath] = make(map[string]struct{})
					}
					incoming[resolved.relativePath][file.relativePath] = struct{}{}
				}
			}
		}
	}

	return issues
}

// appendQuarantinedLink reports a link whose target lives under
// raw/_quarantine/ (spec §7.1) once per file, key and target. key is "" for
// body links.
func appendQuarantinedLink(issues []models.LintIssue, seen map[string]struct{}, file *vaultFile, target, key string) []models.LintIssue {
	dedupe := "q\x00" + file.relativePath + "\x00" + key + "\x00" + target
	if _, exists := seen[dedupe]; exists {
		return issues
	}
	seen[dedupe] = struct{}{}

	message := "link target is a quarantined source"
	if key != "" {
		message = fmt.Sprintf("frontmatter %q link target is a quarantined source", key)
	}
	return append(issues, newIssue(
		models.LintIssueKindLinkToQuarantined,
		models.SeverityWarning,
		file.relativePath,
		message+"; restore it with `kb review accept` or edit the link",
		target,
	))
}

// isQuarantinedFile reports whether file lives under raw/_quarantine/.
func isQuarantinedFile(file *vaultFile) bool {
	if file == nil {
		return false
	}
	_, quarantined := resolve.QuarantineOriginal(file.vaultRelativePath)
	return quarantined
}

// isSourceFile reports whether file is a non-quarantined topic source that
// decisions judge: under raw/, outside raw/codebase/ and not a codebase
// snapshot document.
func isSourceFile(file *vaultFile) bool {
	if file == nil || !file.inTopic || !strings.HasPrefix(file.relativePath, "raw/") || isQuarantinedFile(file) {
		return false
	}
	if strings.HasPrefix(file.relativePath, "raw/codebase/") {
		return false
	}
	return !strings.HasPrefix(strings.TrimSpace(frontmatter.GetString(file.frontmatter, "source_kind")), "codebase-")
}

// affectsPair is an article and a source whose `affects` list links it.
type affectsPair struct {
	article *vaultFile
	source  *vaultFile
}

// affectsPairs returns the article/source pairs linked by a source's
// `affects` list, keyed by pairKey.
func affectsPairs(state vaultState) map[string]affectsPair {
	pairs := make(map[string]affectsPair)
	for _, file := range state.files {
		if file.parseErr != nil || !isSourceFile(file) {
			continue
		}
		for _, target := range resolve.FrontmatterLinks(file.frontmatter)["affects"] {
			article := state.resolveTarget(target, false)
			if article == nil || !article.inTopic || !isWikiConceptPath(article.relativePath) {
				continue
			}
			pairs[pairKey(article, file)] = affectsPair{article: article, source: file}
		}
	}
	return pairs
}

func pairKey(article, source *vaultFile) string {
	return article.relativePath + "\x00" + source.relativePath
}

// findNeedsCompile reports every article older than a source that affects
// it: the source's `scraped` is newer than the article's `updated` (spec
// §11). It replaces the date-only stale check for those pairs.
func findNeedsCompile(state vaultState) []models.LintIssue {
	issues := make([]models.LintIssue, 0)
	pairs := affectsPairs(state)
	for _, key := range sortedKeys(pairs) {
		article, source := pairs[key].article, pairs[key].source
		if article.parseErr != nil {
			continue
		}
		updated := frontmatter.GetTime(article.frontmatter, "updated")
		scraped := frontmatter.GetTime(source.frontmatter, "scraped")
		if updated.IsZero() || scraped.IsZero() || !updated.Before(scraped) {
			continue
		}
		issues = append(issues, newIssue(
			models.LintIssueKindNeedsCompile,
			models.SeverityWarning,
			article.relativePath,
			fmt.Sprintf(
				"source scraped %s affects this article, last updated %s; recompile it",
				scraped.Format(frontmatter.DateLayout),
				updated.Format(frontmatter.DateLayout),
			),
			strings.TrimSuffix(source.relativePath, ".md"),
		))
	}
	return issues
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
