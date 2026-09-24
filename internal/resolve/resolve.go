// Package resolve implements Obsidian-style link resolution for kb vaults and
// the link extraction shared by lint, refs and link: a target resolves by
// topic-relative path, topic-slug-prefixed path, vault-relative path, or file
// stem, case-insensitively. Titles and frontmatter aliases never resolve a
// link; aliases are display text only (`[[target|alias]]`).
package resolve

import (
	"path"
	"sort"
	"strings"
)

// QuarantineDir is the topic-relative directory that holds quarantined
// sources: `raw/articles/x.md` is quarantined to `raw/_quarantine/articles/x.md`.
const QuarantineDir = "raw/_quarantine"

// File is one markdown document known to the index.
type File struct {
	// Path is the vault-relative slash path including the `.md` extension.
	Path string
	// TopicRel is the topic-relative path when the file is in the active
	// topic, else "".
	TopicRel string
	// InTopic reports whether the file belongs to the active topic.
	InTopic bool
	// Quarantined reports whether the file lives under raw/_quarantine/.
	Quarantined bool
	// Title is the frontmatter title, when present.
	Title string
	// Aliases are the frontmatter aliases (display text only).
	Aliases []string
}

// Index resolves link targets against a set of vault files.
type Index struct {
	topicSlug string
	files     []*File
	topicPath map[string][]*File
	vaultPath map[string][]*File
	stem      map[string][]*File
}

// NewIndex builds an index over files for the topic named topicSlug.
// Quarantined files are detected from their paths (in addition to the
// Quarantined flag) and indexed under both their quarantine path and their
// original path, so a link to the original path resolves to the quarantined
// file. When several files match one key, non-quarantined files win, then
// files in the topic, then the lexicographically first path.
func NewIndex(topicSlug string, files []File) *Index {
	index := &Index{
		topicSlug: strings.Trim(topicSlug, "/"),
		files:     make([]*File, 0, len(files)),
		topicPath: make(map[string][]*File, len(files)),
		vaultPath: make(map[string][]*File, len(files)),
		stem:      make(map[string][]*File, len(files)),
	}

	for position := range files {
		file := files[position]
		file.Path = cleanSlashPath(file.Path)
		file.TopicRel = cleanSlashPath(file.TopicRel)
		if !file.Quarantined {
			_, quarantinedTopic := QuarantineOriginal(file.TopicRel)
			_, quarantinedVault := QuarantineOriginal(file.Path)
			file.Quarantined = quarantinedTopic || quarantinedVault
		}
		index.files = append(index.files, &file)
	}

	sort.SliceStable(index.files, func(i, j int) bool {
		return fileLess(index.files[i], index.files[j])
	})

	for _, file := range index.files {
		vaultKey := linkKey(file.Path)
		index.vaultPath[vaultKey] = append(index.vaultPath[vaultKey], file)
		if original, ok := QuarantineOriginal(file.Path); ok {
			originalKey := linkKey(original)
			index.vaultPath[originalKey] = append(index.vaultPath[originalKey], file)
		}

		if file.InTopic && file.TopicRel != "" {
			topicKey := linkKey(file.TopicRel)
			index.topicPath[topicKey] = append(index.topicPath[topicKey], file)
			if original, ok := QuarantineOriginal(file.TopicRel); ok {
				originalKey := linkKey(original)
				index.topicPath[originalKey] = append(index.topicPath[originalKey], file)
			}
		}

		stemKey := path.Base(linkKey(file.Path))
		if stemKey != "" && stemKey != "." {
			index.stem[stemKey] = append(index.stem[stemKey], file)
		}
	}

	return index
}

// Files returns the indexed files in preference order.
func (ix *Index) Files() []*File {
	return append([]*File(nil), ix.files...)
}

// TopicSlug returns the slug of the active topic.
func (ix *Index) TopicSlug() string {
	return ix.topicSlug
}

// Resolve returns the file a link target points to, or nil. See ResolveWhere.
func (ix *Index) Resolve(target string) *File {
	return ix.ResolveWhere(target, nil)
}

// ResolveWhere resolves target like Obsidian and returns the first matching
// file accepted by accept (nil accepts all). The target is normalized with
// Normalize and matched case-insensitively, in order, as a topic-relative
// path, a topic-slug-prefixed path, a vault-relative path and a file stem.
func (ix *Index) ResolveWhere(target string, accept func(*File) bool) *File {
	if ix == nil {
		return nil
	}
	key := strings.ToLower(Normalize(target))
	if key == "" {
		return nil
	}

	if file := pick(ix.topicPath[key], accept); file != nil {
		return file
	}
	if ix.topicSlug != "" {
		if after, ok := strings.CutPrefix(key, strings.ToLower(ix.topicSlug)+"/"); ok {
			if file := pick(ix.topicPath[after], accept); file != nil {
				return file
			}
		}
	}
	if file := pick(ix.vaultPath[key], accept); file != nil {
		return file
	}
	if !strings.Contains(key, "/") {
		if file := pick(ix.stem[key], accept); file != nil {
			return file
		}
	}

	return nil
}

// ResolveFrom resolves target as written in the document at fromPath
// (vault-relative). Targets starting with `./` or `../` are joined with the
// directory of fromPath and matched as vault-relative paths; every other
// target resolves like Resolve.
func (ix *Index) ResolveFrom(fromPath, target string) *File {
	trimmed := strings.ReplaceAll(strings.TrimSpace(target), "\\", "/")
	if strings.HasPrefix(trimmed, "./") || strings.HasPrefix(trimmed, "../") {
		joined := path.Join(path.Dir(cleanSlashPath(fromPath)), stripLinkDecorations(trimmed))
		if !strings.HasPrefix(joined, "../") && joined != ".." {
			if file := pick(ix.vaultPath[linkKey(joined)], nil); file != nil {
				return file
			}
		}
	}

	return ix.Resolve(target)
}

// Normalize reduces a link target to its path part: it strips `[[`/`]]` (and
// a leading `!` embed marker), `|display`, `#heading`, surrounding spaces,
// backslashes (to `/`), a leading `./`, surrounding slashes and a trailing
// `.md`. The result keeps its case.
func Normalize(target string) string {
	trimmed := strings.TrimSpace(target)
	trimmed = strings.TrimPrefix(trimmed, "!")
	trimmed = strings.TrimPrefix(trimmed, "[[")
	trimmed = strings.TrimSuffix(trimmed, "]]")
	trimmed = stripLinkDecorations(trimmed)
	trimmed = strings.ReplaceAll(trimmed, "\\", "/")
	for strings.HasPrefix(trimmed, "./") {
		trimmed = strings.TrimPrefix(trimmed, "./")
	}
	trimmed = strings.Trim(trimmed, "/")
	if strings.HasSuffix(strings.ToLower(trimmed), ".md") {
		trimmed = trimmed[:len(trimmed)-len(".md")]
	}

	return strings.TrimSpace(trimmed)
}

// QuarantineOriginal maps a quarantined path (topic-relative
// `raw/_quarantine/<rest>` or vault-relative `<slug>/raw/_quarantine/<rest>`)
// to the path it had before quarantine (`raw/<rest>` or `<slug>/raw/<rest>`).
func QuarantineOriginal(filePath string) (string, bool) {
	cleaned := cleanSlashPath(filePath)
	prefix := QuarantineDir + "/"
	if rest, ok := strings.CutPrefix(cleaned, prefix); ok && rest != "" {
		return "raw/" + rest, true
	}

	slug, rest, found := strings.Cut(cleaned, "/")
	if !found || slug == "" {
		return "", false
	}
	if inner, ok := strings.CutPrefix(rest, prefix); ok && inner != "" {
		return slug + "/raw/" + inner, true
	}

	return "", false
}

// QuarantinePath maps a topic-relative `raw/<rest>` path to its quarantine
// location `raw/_quarantine/<rest>`.
func QuarantinePath(topicRel string) (string, bool) {
	cleaned := cleanSlashPath(topicRel)
	rest, ok := strings.CutPrefix(cleaned, "raw/")
	if !ok || rest == "" || strings.HasPrefix(cleaned, QuarantineDir+"/") {
		return "", false
	}

	return QuarantineDir + "/" + rest, true
}

func stripLinkDecorations(target string) string {
	if index := strings.IndexAny(target, "|#"); index >= 0 {
		target = target[:index]
	}

	return strings.TrimSpace(target)
}

func linkKey(filePath string) string {
	return strings.ToLower(Normalize(filePath))
}

func cleanSlashPath(filePath string) string {
	cleaned := strings.Trim(strings.ReplaceAll(strings.TrimSpace(filePath), "\\", "/"), "/")
	if cleaned == "" {
		return ""
	}

	return path.Clean(cleaned)
}

func fileLess(left, right *File) bool {
	if left.Quarantined != right.Quarantined {
		return !left.Quarantined
	}
	if left.InTopic != right.InTopic {
		return left.InTopic
	}

	return left.Path < right.Path
}

func pick(candidates []*File, accept func(*File) bool) *File {
	for _, candidate := range candidates {
		if accept == nil || accept(candidate) {
			return candidate
		}
	}

	return nil
}
