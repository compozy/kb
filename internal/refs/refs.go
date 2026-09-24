// Package refs finds the inbound references to a topic document and moves a
// source into and out of quarantine without leaving dead links (spec §7.1).
//
// References are resolved with internal/resolve exactly like lint does:
// topic-relative path, topic-slug-prefixed path, vault-relative path or file
// stem, with Obsidian semantics (titles and aliases never resolve a link).
package refs

import (
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/resolve"
)

// Reference kinds.
const (
	// KindIndexLine is a list item, numbered item or table row of a
	// `wiki/index/` document that links to the target.
	KindIndexLine = "index-line"
	// KindFrontmatter is a frontmatter value that links to the target.
	KindFrontmatter = "frontmatter"
	// KindBody is a body wikilink or markdown link to the target.
	KindBody = "body"
)

// IndexDir is the topic-relative directory whose list items and table rows
// are removed when every link they carry points to quarantined files.
const IndexDir = "wiki/index/"

// RemovableKeys are the frontmatter list keys whose items pointing to a
// quarantined file are removed (spec §7.1).
var RemovableKeys = []string{
	"sources", "related", "extends", "prerequisite", "example_of", "contradicts",
	"affects", "supersedes", "concepts", "informed_by",
}

// Ref is one inbound reference to a topic document.
type Ref struct {
	// File is the topic-relative path of the referring document.
	File string `json:"file"`
	// Kind is KindIndexLine, KindFrontmatter or KindBody.
	Kind string `json:"kind"`
	// Key is the frontmatter key for KindFrontmatter references.
	Key string `json:"key,omitempty"`
	// Line is the 1-based line number in the referring file.
	Line int `json:"line"`
	// Text is the referring line (index lines, frontmatter items) or the
	// link markup (body links).
	Text string `json:"text"`
}

// span is a reference plus what quarantine may do with it.
type span struct {
	ref Ref
	// removable reports whether lines [start, end) may be deleted.
	removable bool
	start     int
	end       int
}

// fileScan is the analysis of one referring file.
type fileScan struct {
	rel     string
	abs     string
	content string
	spans   []span
}

// scanner resolves links from topic documents to one target file.
type scanner struct {
	index      *resolve.Index
	targetPath string
}

// Scan returns the inbound references to target (a topic-relative path such
// as `raw/articles/x.md`) from every other non-quarantined document of the
// topic at topicRoot, in file then line order.
func Scan(vaultPath, topicRoot, target string) ([]Ref, error) {
	scans, err := scanTopic(vaultPath, topicRoot, target)
	if err != nil {
		return nil, err
	}
	refs := make([]Ref, 0)
	for _, scan := range scans {
		for _, sp := range scan.spans {
			refs = append(refs, sp.ref)
		}
	}
	return refs, nil
}

func scanTopic(vaultPath, topicRoot, target string) ([]fileScan, error) {
	index, files, err := resolve.LoadVaultIndex(vaultPath, topicRoot)
	if err != nil {
		return nil, fmt.Errorf("refs: %w", err)
	}
	cleanTarget := cleanRel(target)
	var targetFile *resolve.File
	for position := range files {
		if files[position].InTopic && strings.EqualFold(files[position].TopicRel, cleanTarget) {
			targetFile = &files[position]
			break
		}
	}
	if targetFile == nil {
		return nil, fmt.Errorf("refs: %s is not a document of topic %s", cleanTarget, index.TopicSlug())
	}

	sc := scanner{index: index, targetPath: targetFile.Path}
	scans := make([]fileScan, 0)
	for _, file := range files {
		if !file.InTopic || file.Quarantined || file.Path == targetFile.Path {
			continue
		}
		abs := joinTopic(topicRoot, file.TopicRel)
		content, err := readFile(abs)
		if err != nil {
			return nil, err
		}
		spans := sc.analyze(file, content)
		if len(spans) == 0 {
			continue
		}
		scans = append(scans, fileScan{rel: file.TopicRel, abs: abs, content: content, spans: spans})
	}
	sort.Slice(scans, func(i, j int) bool { return scans[i].rel < scans[j].rel })
	return scans, nil
}

// pointsToTarget reports whether link (written in the document at fromPath)
// resolves to the target.
func (sc scanner) pointsToTarget(fromPath, link string) bool {
	file := sc.index.ResolveFrom(fromPath, link)
	return file != nil && file.Path == sc.targetPath
}

// pointsAway reports whether link resolves to the target or to a file that
// is already quarantined.
func (sc scanner) pointsAway(fromPath, link string) bool {
	file := sc.index.ResolveFrom(fromPath, link)
	return file != nil && (file.Path == sc.targetPath || file.Quarantined)
}

func (sc scanner) analyze(file resolve.File, content string) []span {
	lines := splitLines(content)
	values, _, parseErr := frontmatter.Parse(content)
	locked := parseErr == nil && frontmatter.GetBool(values, "locked")

	bodyStart := 0
	spans := make([]span, 0)
	if closing, ok := frontmatterClose(lines); ok {
		bodyStart = closing + 1
		spans = append(spans, sc.frontmatterSpans(file, lines, closing, locked)...)
	}
	spans = append(spans, sc.bodySpans(file, lines, bodyStart, locked)...)
	return spans
}

func (sc scanner) frontmatterSpans(file resolve.File, lines []string, closing int, locked bool) []span {
	spans := make([]span, 0)
	key := ""
	itemIndent := -1
	for i := 1; i < closing; i++ {
		line := trimEOL(lines[i])
		indent := leadingSpaces(line)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if indent == 0 && !isListItem(trimmed) {
			key, itemIndent = topKey(line), -1
			_, rest, _ := strings.Cut(line, ":")
			if links := stringLinks(rest); sc.anyTarget(file.Path, links) {
				spans = append(spans, span{ref: Ref{File: file.TopicRel, Kind: KindFrontmatter, Key: key, Line: i + 1, Text: trimmed}})
			}
			continue
		}
		if key == "" || !isListItem(trimmed) {
			continue
		}
		end := i + 1
		for end < closing {
			next := trimEOL(lines[end])
			if strings.TrimSpace(next) == "" || leadingSpaces(next) <= indent {
				break
			}
			end++
		}
		if itemIndent < 0 {
			itemIndent = indent
		}
		links := stringLinks(strings.Join(lines[i:end], ""))
		if sc.anyTarget(file.Path, links) {
			removable := !locked && indent == itemIndent && slices.Contains(RemovableKeys, key) && sc.allAway(file.Path, links)
			spans = append(spans, span{
				ref:       Ref{File: file.TopicRel, Kind: KindFrontmatter, Key: key, Line: i + 1, Text: trimmed},
				removable: removable,
				start:     i,
				end:       end,
			})
		}
		i = end - 1
	}
	return spans
}

func (sc scanner) bodySpans(file resolve.File, lines []string, bodyStart int, locked bool) []span {
	body := strings.Join(lines[bodyStart:], "")
	links := resolve.BodyLinks(body)
	if len(links) == 0 {
		return nil
	}
	lineOffsets := make([]int, 0, len(lines)-bodyStart)
	offset := 0
	for _, line := range lines[bodyStart:] {
		lineOffsets = append(lineOffsets, offset)
		offset += len(line)
	}
	lineOf := func(position int) int {
		return sort.Search(len(lineOffsets), func(i int) bool { return lineOffsets[i] > position }) - 1
	}

	byLine := map[int][]resolve.Link{}
	order := make([]int, 0)
	for _, link := range links {
		line := lineOf(link.Start)
		if _, seen := byLine[line]; !seen {
			order = append(order, line)
		}
		byLine[line] = append(byLine[line], link)
	}

	inIndex := strings.HasPrefix(file.TopicRel, IndexDir)
	spans := make([]span, 0)
	for _, line := range order {
		lineLinks := byLine[line]
		targets := make([]string, 0, len(lineLinks))
		for _, link := range lineLinks {
			targets = append(targets, link.Target)
		}
		if !sc.anyTarget(file.Path, targets) {
			continue
		}
		absolute := bodyStart + line
		text := trimEOL(lines[absolute])
		if inIndex && isIndexEntry(text) {
			spans = append(spans, span{
				ref:       Ref{File: file.TopicRel, Kind: KindIndexLine, Line: absolute + 1, Text: strings.TrimSpace(text)},
				removable: !locked && sc.allAway(file.Path, targets),
				start:     absolute,
				end:       absolute + 1,
			})
			continue
		}
		for _, link := range lineLinks {
			if sc.pointsToTarget(file.Path, link.Target) {
				spans = append(spans, span{ref: Ref{File: file.TopicRel, Kind: KindBody, Line: absolute + 1, Text: body[link.Start:link.End]}})
			}
		}
	}
	return spans
}

func (sc scanner) anyTarget(fromPath string, links []string) bool {
	return slices.ContainsFunc(links, func(link string) bool { return sc.pointsToTarget(fromPath, link) })
}

func (sc scanner) allAway(fromPath string, links []string) bool {
	if len(links) == 0 {
		return false
	}
	for _, link := range links {
		if !sc.pointsAway(fromPath, link) {
			return false
		}
	}
	return true
}

// stringLinks returns the normalized wikilink targets in text.
func stringLinks(text string) []string {
	return resolve.FrontmatterLinks(map[string]any{"v": text})["v"]
}

// isIndexEntry reports whether an index line is a list item, a numbered item
// or a table row (not a table separator).
func isIndexEntry(line string) bool {
	trimmed := strings.TrimSpace(line)
	if isListItem(trimmed) || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "+ ") {
		return true
	}
	digits := len(trimmed) - len(strings.TrimLeft(trimmed, "0123456789"))
	if digits > 0 && (strings.HasPrefix(trimmed[digits:], ". ") || strings.HasPrefix(trimmed[digits:], ") ")) {
		return true
	}
	if strings.HasPrefix(trimmed, "|") {
		return strings.Trim(trimmed, "|-: ") != ""
	}
	return false
}

func isListItem(trimmed string) bool {
	return trimmed == "-" || strings.HasPrefix(trimmed, "- ")
}

func topKey(line string) string {
	key, _, _ := strings.Cut(line, ":")
	return strings.Trim(strings.TrimSpace(key), `"'`)
}

func leadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " \t"))
}

// frontmatterClose returns the index of the closing `---` line when lines
// start with a frontmatter block.
func frontmatterClose(lines []string) (int, bool) {
	if len(lines) == 0 || trimEOL(lines[0]) != "---" {
		return 0, false
	}
	for i := 1; i < len(lines); i++ {
		if trimEOL(lines[i]) == "---" {
			return i, true
		}
	}
	return 0, false
}

// splitLines splits text into lines that keep their line endings, so
// strings.Join(splitLines(text), "") == text.
func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	lines := make([]string, 0, strings.Count(text, "\n")+1)
	for text != "" {
		index := strings.IndexByte(text, '\n')
		if index < 0 {
			lines = append(lines, text)
			break
		}
		lines = append(lines, text[:index+1])
		text = text[index+1:]
	}
	return lines
}

func trimEOL(line string) string {
	return strings.TrimRight(line, "\r\n")
}

func cleanRel(rel string) string {
	cleaned := strings.Trim(strings.ReplaceAll(strings.TrimSpace(rel), "\\", "/"), "/")
	if cleaned == "" {
		return ""
	}
	return path.Clean(cleaned)
}
