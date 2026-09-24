package link

import (
	"path"
	"strings"

	"github.com/compozy/kb/internal/resolve"
)

// Targets renders relation targets and resolves existing link targets for
// one topic against the vault index (spec §6, plan "Relation target form").
type Targets struct {
	index *resolve.Index
	// topicPrefix is the vault-relative directory of the topic ("demo").
	topicPrefix string
	stemCount   map[string]int
}

// NewTargets builds the target helper over the vault index files. topicPrefix
// is the vault-relative slash path of the topic root.
func NewTargets(index *resolve.Index, files []resolve.File, topicPrefix string) *Targets {
	targets := &Targets{
		index:       index,
		topicPrefix: strings.Trim(topicPrefix, "/"),
		stemCount:   make(map[string]int, len(files)),
	}
	for _, file := range files {
		targets.stemCount[stemKey(file.Path)]++
	}
	return targets
}

// Form returns the relation target of the topic document at topicRel: its
// file stem, or its vault-relative path without `.md` when the stem is
// ambiguous across the vault. Callers wrap it as `[[<form>]]`.
func (t *Targets) Form(topicRel string) string {
	stem := strings.TrimSuffix(path.Base(topicRel), path.Ext(topicRel))
	if t == nil || t.stemCount[strings.ToLower(stem)] <= 1 {
		return stem
	}
	vaultRel := strings.TrimSuffix(topicRel, path.Ext(topicRel))
	if t.topicPrefix != "" {
		vaultRel = t.topicPrefix + "/" + vaultRel
	}
	return vaultRel
}

// Wikilink returns the quoted-list entry for a relation target: `[[<form>]]`.
func (t *Targets) Wikilink(topicRel string) string {
	return "[[" + t.Form(topicRel) + "]]"
}

// Resolve returns the topic-relative path a link target written in the
// document at fromTopicRel points to, or "" when it resolves to nothing in
// the topic.
func (t *Targets) Resolve(fromTopicRel, target string) string {
	if t == nil || t.index == nil {
		return ""
	}
	from := fromTopicRel
	if t.topicPrefix != "" {
		from = t.topicPrefix + "/" + fromTopicRel
	}
	file := t.index.ResolveFrom(from, target)
	if file == nil || !file.InTopic {
		return ""
	}
	return file.TopicRel
}

func stemKey(filePath string) string {
	base := path.Base(strings.ReplaceAll(filePath, "\\", "/"))
	return strings.ToLower(strings.TrimSuffix(base, path.Ext(base)))
}
