//go:build integration

package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kconfig "github.com/compozy/kb/internal/config"
)

// TestIngestInstagramCommandEndToEndWithFakeYTDLP exercises the real
// `kb ingest instagram` path (CLI -> instagram.Extractor -> mediadl yt-dlp
// subprocess -> ingest -> vault file) using a fake yt-dlp binary so the test is
// deterministic and never touches the network. The caption policy keeps STT out
// of the path.
func TestIngestInstagramCommandEndToEndWithFakeYTDLP(t *testing.T) {
	vault := t.TempDir()
	scriptPath := writeFakeInstagramYTDLP(t)

	configPath := filepath.Join(t.TempDir(), "kb.toml")
	configBody := strings.Join([]string{
		"[app]",
		`name = "kb"`,
		`env = "development"`,
		"",
		"[instagram]",
		`yt_dlp_path = "` + scriptPath + `"`,
		`transcription = "captions"`,
		"retry_attempts = 1",
		`retry_backoff = "1s"`,
	}, "\n") + "\n"
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv(kconfig.EnvConfigPath, configPath)

	runRootArgs(t, "topic", "new", "reels-test", "Reels Test", "reels", "--vault", vault)
	runRootArgs(t, "ingest", "instagram", "https://www.instagram.com/reel/SHORTCODE123/",
		"--topic", "reels-test", "--vault", vault, "--transcribe", "captions")

	rawDir := filepath.Join(vault, "reels-test", "raw", "instagram")
	entries, err := os.ReadDir(rawDir)
	if err != nil {
		t.Fatalf("read instagram raw dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected exactly one ingested file, got %d: %v", len(entries), entries)
	}

	data, err := os.ReadFile(filepath.Join(rawDir, entries[0].Name()))
	if err != nil {
		t.Fatalf("read ingested file: %v", err)
	}
	doc := string(data)
	for _, want := range []string{
		"source_kind: instagram-video",
		"shortcode: SHORTCODE123",
		"transcript_source: captions",
		"## Caption\nExploring the cosmos.",
		"## Transcript\n## 00:00\nExploring the cosmos",
	} {
		if !strings.Contains(doc, want) {
			t.Fatalf("ingested document missing %q\n---\n%s", want, doc)
		}
	}
}

func runRootArgs(t *testing.T, args ...string) {
	t.Helper()
	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs(args)
	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("command %v returned error: %v", args, err)
	}
}

func writeFakeInstagramYTDLP(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "yt-dlp")
	script := strings.Join([]string{
		"#!/bin/sh",
		"set -eu",
		`case " $* " in`,
		`  *" --dump-single-json "*)`,
		"    cat <<'JSON'",
		`{"id":"SHORTCODE123","title":"Reel by NASA","description":"Exploring the cosmos.","uploader":"nasa","like_count":4200,"comment_count":88,"duration":32,"webpage_url":"https://www.instagram.com/reel/SHORTCODE123/","subtitles":{"en":[{"ext":"json3"}]},"automatic_captions":{}}`,
		"JSON",
		"    ;;",
		"  *)",
		`    out_dir=""`,
		`    previous=""`,
		`    for arg in "$@"; do`,
		`      if [ "$previous" = "--paths" ]; then out_dir=${arg#home:}; fi`,
		`      previous="$arg"`,
		"    done",
		`    [ -n "$out_dir" ] || exit 7`,
		`    mkdir -p "$out_dir"`,
		`    cat > "$out_dir/SHORTCODE123.en.json3" <<'CAPTION'`,
		`{"events":[{"tStartMs":0,"segs":[{"utf8":"Exploring the cosmos"}]}]}`,
		"CAPTION",
		"    ;;",
		"esac",
	}, "\n") + "\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake yt-dlp: %v", err)
	}
	return scriptPath
}
