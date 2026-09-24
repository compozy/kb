package resolve_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/compozy/kb/internal/resolve"
)

func testIndex() *resolve.Index {
	return resolve.NewIndex("topic", []resolve.File{
		{Path: "topic/raw/articles/Source Note.md", TopicRel: "raw/articles/Source Note.md", InTopic: true, Title: "A Very Different Title"},
		{Path: "topic/wiki/concepts/Agents.md", TopicRel: "wiki/concepts/Agents.md", InTopic: true, Title: "Agents", Aliases: []string{"AI agents"}},
		{Path: "topic/raw/_quarantine/articles/junk.md", TopicRel: "raw/_quarantine/articles/junk.md", InTopic: true},
		{Path: "topic/raw/articles/dup.md", TopicRel: "raw/articles/dup.md", InTopic: true},
		{Path: "topic/raw/youtube/dup.md", TopicRel: "raw/youtube/dup.md", InTopic: true},
		{Path: "aaa-other/wiki/concepts/dup.md"},
		{Path: "aaa-other/wiki/concepts/Elsewhere.md"},
		{Path: "topic/raw/_quarantine/articles/shadowed.md", TopicRel: "raw/_quarantine/articles/shadowed.md", InTopic: true},
		{Path: "aaa-other/raw/articles/shadowed.md"},
	})
}

func TestIndexResolve(t *testing.T) {
	t.Parallel()

	index := testIndex()
	tests := []struct {
		name            string
		target          string
		wantPath        string
		wantQuarantined bool
	}{
		{name: "topic relative path", target: "raw/articles/Source Note", wantPath: "topic/raw/articles/Source Note.md"},
		{name: "topic relative path with extension", target: "raw/articles/Source Note.md", wantPath: "topic/raw/articles/Source Note.md"},
		{name: "slug prefixed path", target: "topic/wiki/concepts/Agents", wantPath: "topic/wiki/concepts/Agents.md"},
		{name: "vault relative path in other topic", target: "aaa-other/wiki/concepts/Elsewhere", wantPath: "aaa-other/wiki/concepts/Elsewhere.md"},
		{name: "stem case insensitive", target: "source note", wantPath: "topic/raw/articles/Source Note.md"},
		{name: "wikilink markup display and heading", target: "[[Agents#Overview|the agents]]", wantPath: "topic/wiki/concepts/Agents.md"},
		{name: "embed markup", target: "![[Agents]]", wantPath: "topic/wiki/concepts/Agents.md"},
		{name: "backslashes and dot slash", target: `./wiki\concepts\Agents.md`, wantPath: "topic/wiki/concepts/Agents.md"},
		{name: "trailing slash", target: "Agents/", wantPath: "topic/wiki/concepts/Agents.md"},
		{name: "stem in other topic", target: "Elsewhere", wantPath: "aaa-other/wiki/concepts/Elsewhere.md"},
		{name: "ambiguous stem prefers topic then path order", target: "dup", wantPath: "topic/raw/articles/dup.md"},
		{name: "quarantine original path", target: "raw/articles/junk", wantPath: "topic/raw/_quarantine/articles/junk.md", wantQuarantined: true},
		{name: "quarantine original slug path", target: "topic/raw/articles/junk", wantPath: "topic/raw/_quarantine/articles/junk.md", wantQuarantined: true},
		{name: "quarantine path itself", target: "raw/_quarantine/articles/junk", wantPath: "topic/raw/_quarantine/articles/junk.md", wantQuarantined: true},
		{name: "quarantined stem", target: "junk", wantPath: "topic/raw/_quarantine/articles/junk.md", wantQuarantined: true},
		{name: "non-quarantined stem wins over quarantined", target: "shadowed", wantPath: "aaa-other/raw/articles/shadowed.md"},
		{name: "title does not resolve", target: "A Very Different Title"},
		{name: "alias does not resolve", target: "AI agents"},
		{name: "partial path does not resolve", target: "concepts/Agents"},
		{name: "empty", target: "[[ ]]"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := index.Resolve(tc.target)
			if tc.wantPath == "" {
				if got != nil {
					t.Fatalf("Resolve(%q) = %q, want nil", tc.target, got.Path)
				}
				return
			}
			if got == nil {
				t.Fatalf("Resolve(%q) = nil, want %q", tc.target, tc.wantPath)
			}
			if got.Path != tc.wantPath || got.Quarantined != tc.wantQuarantined {
				t.Fatalf("Resolve(%q) = %q quarantined=%v, want %q quarantined=%v", tc.target, got.Path, got.Quarantined, tc.wantPath, tc.wantQuarantined)
			}
		})
	}
}

func TestIndexResolveWhereAndFrom(t *testing.T) {
	t.Parallel()

	index := testIndex()
	onlyOther := func(file *resolve.File) bool { return !file.InTopic }
	if got := index.ResolveWhere("dup", onlyOther); got == nil || got.Path != "aaa-other/wiki/concepts/dup.md" {
		t.Fatalf("ResolveWhere(dup, other) = %#v", got)
	}
	if got := index.ResolveWhere("Agents", onlyOther); got != nil {
		t.Fatalf("ResolveWhere(Agents, other) = %#v, want nil", got)
	}

	if got := index.ResolveFrom("topic/wiki/concepts/Agents.md", "../../raw/articles/dup.md"); got == nil || got.Path != "topic/raw/articles/dup.md" {
		t.Fatalf("ResolveFrom relative = %#v", got)
	}
	if got := index.ResolveFrom("topic/wiki/concepts/Agents.md", "./Agents.md"); got == nil || got.Path != "topic/wiki/concepts/Agents.md" {
		t.Fatalf("ResolveFrom dot = %#v", got)
	}
	if got := index.ResolveFrom("topic/wiki/concepts/Agents.md", "../../../../outside.md"); got != nil {
		t.Fatalf("ResolveFrom outside = %#v, want nil", got)
	}

	var nilIndex *resolve.Index
	if got := nilIndex.Resolve("x"); got != nil {
		t.Fatalf("nil index Resolve = %#v", got)
	}
	if index.TopicSlug() != "topic" || len(index.Files()) != 9 {
		t.Fatalf("TopicSlug/Files = %q/%d", index.TopicSlug(), len(index.Files()))
	}
}

func TestNormalize(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"[[Note]]":                 "Note",
		"  [[Note|Display]]  ":     "Note",
		"[[Note#Heading|Display]]": "Note",
		"![[Image Note]]":          "Image Note",
		`folder\Note.md`:           "folder/Note",
		"./folder/Note.MD":         "folder/Note",
		"/folder/Note/":            "folder/Note",
		"Note#^block":              "Note",
		"":                         "",
	}
	for input, want := range tests {
		if got := resolve.Normalize(input); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestQuarantinePaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{input: "raw/_quarantine/articles/x.md", want: "raw/articles/x.md", ok: true},
		{input: "topic/raw/_quarantine/articles/x.md", want: "topic/raw/articles/x.md", ok: true},
		{input: "raw/articles/x.md"},
		{input: "raw/_quarantine"},
		{input: "wiki/concepts/x.md"},
	}
	for _, tc := range tests {
		got, ok := resolve.QuarantineOriginal(tc.input)
		if got != tc.want || ok != tc.ok {
			t.Errorf("QuarantineOriginal(%q) = %q, %v; want %q, %v", tc.input, got, ok, tc.want, tc.ok)
		}
	}

	if got, ok := resolve.QuarantinePath("raw/articles/x.md"); !ok || got != "raw/_quarantine/articles/x.md" {
		t.Fatalf("QuarantinePath = %q, %v", got, ok)
	}
	for _, input := range []string{"raw/_quarantine/articles/x.md", "wiki/x.md", "raw/"} {
		if got, ok := resolve.QuarantinePath(input); ok {
			t.Fatalf("QuarantinePath(%q) = %q, want not ok", input, got)
		}
	}
}

func TestLoadVaultIndex(t *testing.T) {
	t.Parallel()

	vault := t.TempDir()
	writeFile(t, vault, "topic/CLAUDE.md", "# Topic\n")
	writeFile(t, vault, "topic/AGENTS.md", "# Agents file\n")
	writeFile(t, vault, "topic/wiki/concepts/Agents.md", "---\ntitle: Agents\naliases:\n    - AI agents\n---\nbody\n")
	writeFile(t, vault, "topic/raw/articles/single.md", "---\ntitle: Single\naliases: One\n---\nbody\n")
	writeFile(t, vault, "topic/raw/articles/broken.md", "---\ntitle: [broken\n---\nbody\n")
	writeFile(t, vault, "topic/.decisions/notes.md", "ignored\n")
	writeFile(t, vault, "other/CLAUDE.md", "# Other\n")
	writeFile(t, vault, "other/wiki/concepts/Else.md", "---\ntitle: Else\n---\n")
	writeFile(t, vault, "not-a-topic/wiki/concepts/Hidden.md", "---\ntitle: Hidden\n---\n")

	index, files, err := resolve.LoadVaultIndex(vault, filepath.Join(vault, "topic"))
	if err != nil {
		t.Fatalf("LoadVaultIndex returned error: %v", err)
	}

	gotPaths := make([]string, 0, len(files))
	for _, file := range files {
		gotPaths = append(gotPaths, file.Path)
	}
	wantPaths := []string{
		"topic/CLAUDE.md",
		"topic/raw/articles/broken.md",
		"topic/raw/articles/single.md",
		"topic/wiki/concepts/Agents.md",
		"other/CLAUDE.md",
		"other/wiki/concepts/Else.md",
	}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Fatalf("paths = %#v, want %#v", gotPaths, wantPaths)
	}

	agents := index.Resolve("agents")
	if agents == nil || agents.Title != "Agents" || !reflect.DeepEqual(agents.Aliases, []string{"AI agents"}) || agents.TopicRel != "wiki/concepts/Agents.md" {
		t.Fatalf("agents = %#v", agents)
	}
	if single := index.Resolve("single"); single == nil || !reflect.DeepEqual(single.Aliases, []string{"One"}) {
		t.Fatalf("single = %#v", single)
	}
	if broken := index.Resolve("broken"); broken == nil || broken.Title != "" {
		t.Fatalf("broken = %#v", broken)
	}
	if other := index.Resolve("Else"); other == nil || other.InTopic || other.TopicRel != "" {
		t.Fatalf("other = %#v", other)
	}
	if hidden := index.Resolve("Hidden"); hidden != nil {
		t.Fatalf("hidden = %#v, want nil (not a topic)", hidden)
	}

	if _, _, err := resolve.LoadVaultIndex(filepath.Join(vault, "missing"), filepath.Join(vault, "missing", "topic")); err == nil {
		t.Fatal("LoadVaultIndex(missing vault) returned nil error")
	}
}

func writeFile(t *testing.T, root, relative, content string) {
	t.Helper()

	absolute := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", relative, err)
	}
	if err := os.WriteFile(absolute, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relative, err)
	}
}
