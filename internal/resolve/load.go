package resolve

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/compozy/kb/internal/frontmatter"
)

// DecisionsDir is the per-topic directory of decision bookkeeping; it holds
// no vault documents and is never indexed.
const DecisionsDir = ".decisions"

// DiscoverTopicRoots returns topicPath plus every sibling directory of
// vaultPath that carries a topic marker (`CLAUDE.md`), sorted.
func DiscoverTopicRoots(vaultPath, topicPath string) ([]string, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("read vault path %q: %w", vaultPath, err)
	}

	roots := make([]string, 0, len(entries)+1)
	seen := map[string]struct{}{topicPath: {}}
	roots = append(roots, topicPath)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		candidate := filepath.Join(vaultPath, entry.Name())
		if _, exists := seen[candidate]; exists {
			continue
		}

		markerPath := filepath.Join(candidate, "CLAUDE.md")
		info, err := os.Stat(markerPath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("stat topic marker %q: %w", markerPath, err)
		}
		if info.IsDir() {
			continue
		}

		seen[candidate] = struct{}{}
		roots = append(roots, candidate)
	}

	sort.Strings(roots)
	return roots, nil
}

// LoadVaultIndex walks the topic at topicPath and its sibling topics under
// vaultPath (see DiscoverTopicRoots), reads each markdown file's frontmatter
// title and aliases, and returns the index plus the files in preference
// order. `.decisions/` directories and the topic-root `AGENTS.md` are
// skipped; unparsable frontmatter leaves Title and Aliases empty.
func LoadVaultIndex(vaultPath, topicPath string) (*Index, []File, error) {
	cleanTopic := filepath.Clean(topicPath)
	cleanVault := filepath.Clean(vaultPath)
	roots, err := DiscoverTopicRoots(cleanVault, cleanTopic)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve: %w", err)
	}

	files := make([]File, 0)
	for _, root := range roots {
		rootSlug := filepath.Base(root)
		inTopic := root == cleanTopic
		walkErr := filepath.WalkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("walk %q: %w", current, walkErr)
			}
			if entry.IsDir() {
				if entry.Name() == DecisionsDir && current != root {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
				return nil
			}

			relative, err := filepath.Rel(root, current)
			if err != nil {
				return fmt.Errorf("resolve relative path for %q: %w", current, err)
			}
			relative = filepath.ToSlash(relative)
			if relative == "AGENTS.md" {
				return nil
			}

			content, err := os.ReadFile(current)
			if err != nil {
				return fmt.Errorf("read %q: %w", current, err)
			}

			file := File{Path: path.Join(rootSlug, relative), InTopic: inTopic}
			if inTopic {
				file.TopicRel = relative
			}
			if values, _, parseErr := frontmatter.Parse(string(content)); parseErr == nil {
				file.Title = strings.TrimSpace(frontmatter.GetString(values, "title"))
				file.Aliases = Aliases(values)
			}
			files = append(files, file)
			return nil
		})
		if walkErr != nil {
			return nil, nil, fmt.Errorf("resolve: %w", walkErr)
		}
	}

	index := NewIndex(filepath.Base(cleanTopic), files)
	ordered := make([]File, 0, len(files))
	for _, file := range index.files {
		ordered = append(ordered, *file)
	}

	return index, ordered, nil
}

// Aliases reads the Obsidian `aliases` property: a list of strings or a
// single string. Empty entries are dropped.
func Aliases(values map[string]any) []string {
	raw := frontmatter.GetStringSlice(values, "aliases")
	if single, isString := values["aliases"].(string); isString {
		raw = []string{single}
	}

	aliases := make([]string, 0, len(raw))
	for _, alias := range raw {
		if trimmed := strings.TrimSpace(alias); trimmed != "" {
			aliases = append(aliases, trimmed)
		}
	}
	if len(aliases) == 0 {
		return nil
	}

	return aliases
}
