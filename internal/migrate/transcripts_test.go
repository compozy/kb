package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/topic"
)

func TestTranscriptsMovesLegacyMarkdownIntoYouTubeDirectory(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()
	if _, err := topic.New(vaultPath, "systems-design", "Systems Design", "systems"); err != nil {
		t.Fatalf("create topic: %v", err)
	}
	topicPath := filepath.Join(vaultPath, "systems-design")
	writeFile(t, filepath.Join(topicPath, "raw", "transcripts", "talk.md"), "# Talk\n")
	writeFile(t, filepath.Join(topicPath, "raw", "transcripts", ".hidden.md"), "# Hidden\n")
	writeFile(t, filepath.Join(topicPath, "raw", "transcripts", "notes.txt"), "notes\n")

	result, err := Transcripts(vaultPath, "systems-design", time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Transcripts returned error: %v", err)
	}

	if result.Moved != 1 {
		t.Fatalf("moved = %d, want 1", result.Moved)
	}
	assertFileExists(t, filepath.Join(topicPath, "raw", "youtube", "talk.md"))
	assertMissing(t, filepath.Join(topicPath, "raw", "transcripts", "talk.md"))
	assertFileExists(t, filepath.Join(topicPath, "raw", "transcripts", ".hidden.md"))
	assertFileExists(t, filepath.Join(topicPath, "raw", "transcripts", "notes.txt"))

	logContent := readFile(t, filepath.Join(topicPath, "log.md"))
	if !strings.Contains(logContent, "## [2026-05-18] migrate | transcripts") {
		t.Fatalf("log missing migration entry:\n%s", logContent)
	}
}

func TestTranscriptsFailsBeforeMovingWhenTargetExists(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()
	if _, err := topic.New(vaultPath, "systems-design", "Systems Design", "systems"); err != nil {
		t.Fatalf("create topic: %v", err)
	}
	topicPath := filepath.Join(vaultPath, "systems-design")
	writeFile(t, filepath.Join(topicPath, "raw", "transcripts", "talk.md"), "# Legacy\n")
	writeFile(t, filepath.Join(topicPath, "raw", "youtube", "talk.md"), "# Current\n")

	_, err := Transcripts(vaultPath, "systems-design", time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC))
	if err == nil || !strings.Contains(err.Error(), "target already exists") {
		t.Fatalf("Transcripts error = %v, want target conflict", err)
	}

	if got := readFile(t, filepath.Join(topicPath, "raw", "transcripts", "talk.md")); got != "# Legacy\n" {
		t.Fatalf("legacy file changed: %q", got)
	}
	if got := readFile(t, filepath.Join(topicPath, "raw", "youtube", "talk.md")); got != "# Current\n" {
		t.Fatalf("target file changed: %q", got)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("%s is a directory", path)
	}
}

func assertMissing(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be missing, err=%v", path, err)
	}
}
