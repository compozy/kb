package contract

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/compozy/kb/internal/contract/topicyaml"
)

// SectionHeading is the CLAUDE.md heading that carries the contract.
const SectionHeading = "## Selection contract"

// sectionNote tells readers the section is generated.
const sectionNote = "<!-- Rendered by kb from topic.yaml `contract:` (the only source of truth). Edit it through `kb topic contract`; manual edits here are overwritten on the next accept. -->"

// notSetMarker prefixes every placeholder value of an empty contract; the
// parser treats such values as empty.
const notSetMarker = "_not set"

const claudeFileName = "CLAUDE.md"

// ErrNoSection is returned when a CLAUDE.md has no `## Selection contract`.
var ErrNoSection = errors.New("contract: CLAUDE.md has no ## Selection contract section")

var (
	bulletPattern      = regexp.MustCompile(`^[-*+]\s+\*\*([^*]+?)\*\*\s*:?\s*(.*)$`)
	subBulletPattern   = regexp.MustCompile(`^\s+[-*+]\s+(.*)$`)
	headingLinePattern = regexp.MustCompile(`^#{1,2}\s`)
	fencePattern       = regexp.MustCompile("^\\s*(```|~~~)")
)

type fieldID int

const (
	fieldNone fieldID = iota
	fieldPurpose
	fieldCore
	fieldAdjacent
	fieldCollected
	fieldCollectedPaths
	fieldOutOfScope
)

var fieldLabels = map[string]fieldID{
	"purpose":                      fieldPurpose,
	"core":                         fieldCore,
	"central":                      fieldCore,
	"adjacent (keep)":              fieldAdjacent,
	"adjacent":                     fieldAdjacent,
	"collected on purpose":         fieldCollected,
	"collected on purpose (paths)": fieldCollectedPaths,
	"out of scope":                 fieldOutOfScope,
}

// RenderSection renders the `## Selection contract` markdown section for c,
// ending with a single newline. List items are joined with "; " inside one
// bullet (the format used by existing vault topics); a list whose items
// contain a top-level ";" is rendered as nested sub-bullets so it re-imports
// losslessly. A nil or empty contract renders placeholder bullets.
func RenderSection(c *Contract) string {
	var builder strings.Builder
	builder.WriteString(SectionHeading)
	builder.WriteString("\n\n")
	builder.WriteString(sectionNote)
	builder.WriteString("\n\n")

	if c.Empty() {
		builder.WriteString("- **Purpose:** " + notSetMarker + ": run `kb topic contract <topic> --draft` (or `--import-claude`), review the draft in topic.yaml, then `kb topic contract <topic> --accept`._\n")
		for _, label := range []string{"Core", "Adjacent (keep)", "Collected on purpose", "Out of scope"} {
			builder.WriteString("- **" + label + ":** " + notSetMarker + "_\n")
		}
		return builder.String()
	}

	normalized := c.Normalized()
	builder.WriteString("- **Purpose:** " + normalized.Purpose + "\n")
	writeListBullet(&builder, "Core", normalized.Core, false)
	writeListBullet(&builder, "Adjacent (keep)", normalized.Adjacent, false)
	writeListBullet(&builder, "Collected on purpose", normalized.CollectedOnPurpose, false)
	if len(normalized.CollectedOnPurposePaths) > 0 {
		writeListBullet(&builder, "Collected on purpose (paths)", normalized.CollectedOnPurposePaths, true)
	}
	writeListBullet(&builder, "Out of scope", normalized.OutOfScope, false)
	return builder.String()
}

func writeListBullet(builder *strings.Builder, label string, items []string, code bool) {
	rendered := make([]string, len(items))
	nested := false
	for index, item := range items {
		if code {
			rendered[index] = "`" + item + "`"
		} else {
			rendered[index] = item
		}
		if len(splitTopLevel(item)) > 1 {
			nested = true
		}
	}
	builder.WriteString("- **" + label + ":**")
	switch {
	case len(rendered) == 0:
		builder.WriteString(" none\n")
	case nested:
		builder.WriteString("\n")
		for _, item := range rendered {
			builder.WriteString("  - " + item + "\n")
		}
	default:
		builder.WriteString(" " + strings.Join(rendered, "; ") + "\n")
	}
}

// RenderIntoClaude returns claudeMarkdown with its `## Selection contract`
// section replaced by the rendering of c (up to the next H1/H2 heading), or
// with the section inserted before the first `## ` heading (appended at the
// end when there is none). Every other byte is preserved.
func RenderIntoClaude(claudeMarkdown string, c *Contract) string {
	section := RenderSection(c)
	lines := splitLinesKeepEnds(claudeMarkdown)
	start, end := findSection(lines)
	if start >= 0 {
		replacement := section
		if end < len(lines) {
			replacement += "\n"
		}
		return strings.Join(lines[:start], "") + replacement + strings.Join(lines[end:], "")
	}

	insertAt := firstH2(lines)
	if insertAt < 0 {
		prefix := claudeMarkdown
		switch {
		case prefix == "":
			return section
		case strings.HasSuffix(prefix, "\n\n"):
		case strings.HasSuffix(prefix, "\n"):
			prefix += "\n"
		default:
			prefix += "\n\n"
		}
		return prefix + section
	}
	before := strings.Join(lines[:insertAt], "")
	if before != "" && !strings.HasSuffix(before, "\n\n") {
		if strings.HasSuffix(before, "\n") {
			before += "\n"
		} else {
			before += "\n\n"
		}
	}
	return before + section + "\n" + strings.Join(lines[insertAt:], "")
}

// ParseClaudeSection parses the `## Selection contract` section of a topic
// CLAUDE.md into a contract. Bullets `Purpose`, `Core` (or `Central`),
// `Adjacent (keep)` (or `Adjacent`), `Collected on purpose`, `Collected on
// purpose (paths)` and `Out of scope` are read; inline list values are split
// on top-level ";" (never inside parentheses or brackets) and nested
// sub-bullets are taken as items. A missing section returns ErrNoSection.
func ParseClaudeSection(claudeMarkdown string) (*Contract, error) {
	lines := splitLinesKeepEnds(claudeMarkdown)
	start, end := findSection(lines)
	if start < 0 {
		return nil, ErrNoSection
	}

	type fieldText struct {
		inline []string
		items  []string
	}
	fields := map[fieldID]*fieldText{}
	current := fieldNone
	inComment := false
	afterBlank := false
	for _, raw := range lines[start+1 : end] {
		line := strings.TrimRight(raw, "\r\n")
		trimmed := strings.TrimSpace(line)
		if inComment {
			if strings.Contains(trimmed, "-->") {
				inComment = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "<!--") {
			inComment = !strings.Contains(trimmed, "-->")
			continue
		}
		if trimmed == "" {
			afterBlank = true
			continue
		}
		wasBlank := afterBlank
		afterBlank = false
		if match := bulletPattern.FindStringSubmatch(line); match != nil {
			label := strings.ToLower(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(match[1]), ":")))
			current = fieldLabels[label]
			if current == fieldNone {
				continue
			}
			entry := &fieldText{}
			fields[current] = entry
			if value := strings.TrimSpace(match[2]); value != "" {
				entry.inline = append(entry.inline, value)
			}
			continue
		}
		if current == fieldNone {
			continue
		}
		entry := fields[current]
		if match := subBulletPattern.FindStringSubmatch(line); match != nil {
			entry.items = append(entry.items, strings.TrimSpace(match[1]))
			continue
		}
		if line != trimmed || (!wasBlank && !strings.HasPrefix(trimmed, "-") && !strings.HasPrefix(trimmed, "*")) {
			// Wrapped continuation of the bullet value.
			entry.inline = append(entry.inline, trimmed)
			continue
		}
		current = fieldNone
	}

	parsed := &Contract{}
	for id, entry := range fields {
		inline := strings.Join(entry.inline, " ")
		if isPlaceholder(inline) {
			inline = ""
		}
		if id == fieldPurpose {
			parsed.Purpose = strings.TrimSpace(strings.Join(append([]string{inline}, entry.items...), " "))
			continue
		}
		var items []string
		if inline != "" && !isNone(inline) {
			items = append(items, splitTopLevel(inline)...)
		}
		items = append(items, entry.items...)
		if id == fieldCollectedPaths {
			for index, item := range items {
				items[index] = strings.Trim(strings.TrimSpace(item), "`")
			}
		}
		switch id {
		case fieldCore:
			parsed.Core = items
		case fieldAdjacent:
			parsed.Adjacent = items
		case fieldCollected:
			parsed.CollectedOnPurpose = items
		case fieldCollectedPaths:
			parsed.CollectedOnPurposePaths = items
		case fieldOutOfScope:
			parsed.OutOfScope = items
		}
	}
	normalized := parsed.Normalized()
	return &normalized, nil
}

// SyncClaude re-renders the `## Selection contract` section of
// <topicRoot>/CLAUDE.md from the active topic.yaml contract. It is called after
// a contract is accepted. The file is only rewritten when it changes.
func SyncClaude(topicRoot string) error {
	settings, err := LoadSettings(topicRoot)
	if err != nil {
		return err
	}
	claudePath := filepath.Join(topicRoot, claudeFileName)
	content, err := os.ReadFile(claudePath)
	if err != nil {
		return fmt.Errorf("contract: sync CLAUDE.md: %w", err)
	}
	rendered := RenderIntoClaude(string(content), settings.Contract)
	if rendered == string(content) {
		return nil
	}
	if err := topicyaml.WriteFileAtomic(claudePath, []byte(rendered)); err != nil {
		return fmt.Errorf("contract: sync CLAUDE.md: %w", err)
	}
	return nil
}

func isPlaceholder(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), notSetMarker)
}

func isNone(value string) bool {
	switch strings.ToLower(strings.TrimSuffix(strings.TrimSpace(value), ".")) {
	case "none", "n/a", "-":
		return true
	default:
		return false
	}
}

// splitTopLevel splits value on ";" outside parentheses, brackets, braces and
// backticks, trimming items and dropping empty ones.
func splitTopLevel(value string) []string {
	var items []string
	depth := 0
	inCode := false
	start := 0
	for index, r := range value {
		switch {
		case r == '`':
			inCode = !inCode
		case inCode:
		case r == '(' || r == '[' || r == '{':
			depth++
		case (r == ')' || r == ']' || r == '}') && depth > 0:
			depth--
		case r == ';' && depth == 0:
			if item := strings.TrimSpace(value[start:index]); item != "" {
				items = append(items, item)
			}
			start = index + 1
		}
	}
	if item := strings.TrimSpace(value[start:]); item != "" {
		items = append(items, item)
	}
	return items
}

// splitLinesKeepEnds splits text into lines that keep their terminators, so
// joining them reproduces the input byte for byte.
func splitLinesKeepEnds(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.SplitAfter(text, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// findSection returns the line range [start, end) of the contract section,
// or (-1, -1). Headings inside fenced code blocks are ignored.
func findSection(lines []string) (int, int) {
	start := -1
	inFence := false
	for index, raw := range lines {
		line := strings.TrimRight(raw, "\r\n")
		if fencePattern.MatchString(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if start < 0 {
			if strings.EqualFold(strings.TrimSpace(line), SectionHeading) {
				start = index
			}
			continue
		}
		if headingLinePattern.MatchString(line) {
			return start, index
		}
	}
	if start < 0 {
		return -1, -1
	}
	return start, len(lines)
}

func firstH2(lines []string) int {
	inFence := false
	for index, raw := range lines {
		line := strings.TrimRight(raw, "\r\n")
		if fencePattern.MatchString(line) {
			inFence = !inFence
			continue
		}
		if !inFence && strings.HasPrefix(line, "## ") {
			return index
		}
	}
	return -1
}
