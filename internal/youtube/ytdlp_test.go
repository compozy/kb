package youtube

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
)

func TestYTDLPBackendExtractsJSON3Captions(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ",`,
			`"title":"Example Video",`,
			`"channel":"Example Channel",`,
			`"duration":95,`,
			`"upload_date":"20240307",`,
			`"webpage_url":"https://www.youtube.com/watch?v=dQw4w9WgXcQ",`,
			`"subtitles":{"en":[{"ext":"json3"}]},`,
			`"automatic_captions":{"pt-BR":[{"ext":"json3"}]}}`,
		}, ""),
		captionExt: "json3",
		captionBody: strings.Join([]string{
			`{"events":[`,
			`{"tStartMs":0,"segs":[{"utf8":" Hello   world "}]},`,
			`{"tStartMs":5000,"segs":[{"utf8":"Second"},{"utf8":" line"}]}`,
			`]}`,
		}, ""),
	})
	backend := newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{
		YTDLPPath:   "custom-yt-dlp",
		Proxy:       "http://proxy.internal:8080",
		CookiesFile: "/tmp/youtube-cookies.txt",
		UserAgent:   "kb-test-agent",
	}, retryPolicy{Attempts: 5})

	result, err := backend.Extract(context.Background(), parsedRickRollURL(), nil)
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}

	if result.Metadata.Title != "Example Video" {
		t.Fatalf("title = %q, want Example Video", result.Metadata.Title)
	}
	if result.Metadata.Channel != "Example Channel" {
		t.Fatalf("channel = %q, want Example Channel", result.Metadata.Channel)
	}
	if result.Metadata.Duration != 95*time.Second {
		t.Fatalf("duration = %v, want 95s", result.Metadata.Duration)
	}
	if result.Metadata.PublishDate != time.Date(2024, time.March, 7, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("publish date = %v", result.Metadata.PublishDate)
	}
	if result.Source != TranscriptSourceCaptions {
		t.Fatalf("source = %q, want captions", result.Source)
	}
	if result.Language != "en" {
		t.Fatalf("language = %q, want en", result.Language)
	}
	wantMarkdown := "## 00:00\nHello world\n\n## 00:05\nSecond line"
	if result.Markdown != wantMarkdown {
		t.Fatalf("markdown = %q, want %q", result.Markdown, wantMarkdown)
	}

	invocations := readYTDLPInvocationLog(t, logPath)
	if len(invocations) != 2 {
		t.Fatalf("invocations = %#v, want metadata and captions", invocations)
	}
	assertArgsContain(t, invocations[0], "--dump-single-json")
	assertArgsContain(t, invocations[0], "--ignore-config")
	assertArgsContain(t, invocations[0], "--no-playlist")
	assertArgsContain(t, invocations[0], "--proxy", "http://proxy.internal:8080")
	assertArgsContain(t, invocations[0], "--cookies", "/tmp/youtube-cookies.txt")
	assertArgsContain(t, invocations[0], "--user-agent", "kb-test-agent")
	assertArgsContain(t, invocations[0], "--retries", "5")
	assertArgsContain(t, invocations[1], "--ignore-config")
	assertArgsContain(t, invocations[1], "--no-playlist")
	assertArgsContain(t, invocations[1], "--retries", "5")
	assertArgsContain(t, invocations[1], "--fragment-retries", "5")
	assertArgsContain(t, invocations[1], "--write-subs")
	assertArgsContain(t, invocations[1], "--sub-langs", "en")
	assertArgsContain(t, invocations[1], "--sub-format", "json3/vtt/best")
	if containsArg(invocations[1], "--write-auto-subs") {
		t.Fatalf("manual captions should not use --write-auto-subs: %#v", invocations[1])
	}
}

func TestYTDLPBackendExtractsVTTFallback(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"VTT Video",`,
			`"subtitles":{"en":[{"ext":"vtt"}]},"automatic_captions":{}}`,
		}, ""),
		captionExt: "vtt",
		captionBody: strings.Join([]string{
			"WEBVTT",
			"",
			"00:00:01.250 --> 00:00:02.500",
			"First caption",
			"",
			"00:02.000 --> 00:03.000",
			"Second caption",
			"",
		}, "\n"),
	})
	backend := newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

	result, err := backend.Extract(context.Background(), parsedRickRollURL(), nil)
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}

	wantMarkdown := "## 00:01\nFirst caption\n\n## 00:02\nSecond caption"
	if result.Markdown != wantMarkdown {
		t.Fatalf("markdown = %q, want %q", result.Markdown, wantMarkdown)
	}
}

func TestYTDLPBackendSelectsPreferredAutomaticCaption(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"Preferred",`,
			`"subtitles":{"en":[{"ext":"json3"}]},`,
			`"automatic_captions":{"pt-BR":[{"ext":"json3"}]}}`,
		}, ""),
		captionExt:  "json3",
		captionBody: `{"events":[{"tStartMs":0,"segs":[{"utf8":"Preferido"}]}]}`,
	})
	backend := newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

	result, err := backend.Extract(context.Background(), parsedRickRollURL(), []string{"pt"})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Language != "pt-BR" {
		t.Fatalf("language = %q, want pt-BR", result.Language)
	}
	invocations := readYTDLPInvocationLog(t, logPath)
	assertArgsContain(t, invocations[1], "--write-auto-subs")
	assertArgsContain(t, invocations[1], "--sub-langs", "pt-BR")
}

func TestExtractorFailsWhenYTDLPIsUnavailable(t *testing.T) {
	t.Parallel()

	extractor := &Extractor{
		ytDLP: &ytDLPBackend{
			binaryPath: "missing-yt-dlp",
			lookPath: func(string) (string, error) {
				return "", exec.ErrNotFound
			},
			commandContext: exec.CommandContext,
		},
	}

	_, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
	if err == nil {
		t.Fatal("expected Extract to fail")
	}
	if !errors.Is(err, errYTDLPUnavailable) {
		t.Fatalf("error = %v, want yt-dlp unavailable", err)
	}
}

func TestExtractorUsesYTDLPWhenAvailable(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Primary","subtitles":{"en":[{"ext":"json3"}]}}`,
		captionExt:   "json3",
		captionBody:  `{"events":[{"tStartMs":0,"segs":[{"utf8":"Primary transcript"}]}]}`,
	})
	extractor := &Extractor{
		ytDLP: newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
	}

	result, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if !strings.Contains(result.Markdown, "Primary transcript") {
		t.Fatalf("markdown = %q, want primary transcript", result.Markdown)
	}
}

func TestExtractorUsesSTTWhenYTDLPProvesCaptionsUnavailable(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"No Captions","subtitles":{},"automatic_captions":{}}`,
		audioExt:     "mp3",
		audioBody:    "audio-from-ytdlp",
	})
	stt := &stubSTTClient{transcript: "Fallback transcript"}
	extractor := &Extractor{
		ytDLP:     newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:       stt,
		sttConfig: config.Default().STT,
	}

	result, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{
		TranscriptionPolicy: TranscriptionPolicyAuto,
	})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Source != TranscriptSourceSTT {
		t.Fatalf("source = %q, want stt", result.Source)
	}
	if result.Markdown != "## 00:00\nFallback transcript" {
		t.Fatalf("markdown = %q", result.Markdown)
	}
	if !reflect.DeepEqual(stt.audio, []byte("audio-from-ytdlp\n")) {
		t.Fatalf("stt audio = %q, want audio-from-ytdlp", string(stt.audio))
	}
}

func TestExtractorAutoPolicyUsesSTTWhenOnlyAutomaticCaptionsExist(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Automatic Only","subtitles":{},"automatic_captions":{"en":[{"ext":"json3"}]}}`,
		audioExt:     "mp3",
		audioBody:    "audio-from-ytdlp",
		captionExit:  1,
		captionErr:   "caption download should not run for auto policy with only automatic captions",
	})
	stt := &stubSTTClient{transcript: "STT transcript"}
	extractor := &Extractor{
		ytDLP:     newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:       stt,
		sttConfig: config.Default().STT,
	}

	result, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{
		TranscriptionPolicy: TranscriptionPolicyAuto,
	})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Source != TranscriptSourceSTT {
		t.Fatalf("source = %q, want stt", result.Source)
	}
	if result.Markdown != "## 00:00\nSTT transcript" {
		t.Fatalf("markdown = %q", result.Markdown)
	}
	invocations := readYTDLPInvocationLog(t, logPath)
	if len(invocations) != 2 {
		t.Fatalf("invocations = %#v, want metadata and audio", invocations)
	}
	if containsArg(invocations[1], "--write-subs") || containsArg(invocations[1], "--write-auto-subs") {
		t.Fatalf("caption download should not run for auto policy without manual captions: %#v", invocations[1])
	}
}

func TestExtractorSTTPolicyIgnoresYTDLPCaptions(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Captioned","subtitles":{"en":[{"ext":"json3"}]},"automatic_captions":{}}`,
		audioExt:     "mp3",
		audioBody:    "audio-from-ytdlp",
		captionExit:  1,
		captionErr:   "caption download should not run",
	})
	stt := &stubSTTClient{transcript: "Forced STT transcript"}
	extractor := &Extractor{
		ytDLP:     newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:       stt,
		sttConfig: config.Default().STT,
	}

	result, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{
		TranscriptionPolicy: TranscriptionPolicySTT,
	})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Source != TranscriptSourceSTT {
		t.Fatalf("source = %q, want stt", result.Source)
	}
	if result.Markdown != "## 00:00\nForced STT transcript" {
		t.Fatalf("markdown = %q", result.Markdown)
	}
	invocations := readYTDLPInvocationLog(t, logPath)
	if len(invocations) != 2 {
		t.Fatalf("invocations = %#v, want metadata and audio", invocations)
	}
	if containsArg(invocations[1], "--write-subs") || containsArg(invocations[1], "--write-auto-subs") {
		t.Fatalf("caption download should not run for STT policy: %#v", invocations[1])
	}
}

func TestExtractorDoesNotUseSTTWhenYTDLPSeesCaptionsButFetchFails(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Captioned","subtitles":{"en":[{"ext":"json3"}]},"automatic_captions":{}}`,
		captionExit:  1,
		captionErr:   "HTTP Error 429: Too Many Requests",
	})
	stt := &stubSTTClient{transcript: "should not run"}
	extractor := &Extractor{
		ytDLP: newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:   stt,
	}

	_, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
	if err == nil {
		t.Fatal("expected Extract to fail")
	}
	if !errors.Is(err, errYTDLPCaptionFetchFailed) {
		t.Fatalf("error = %v, want yt-dlp caption fetch failure", err)
	}
	var youtubeErr *Error
	if !errors.As(err, &youtubeErr) || youtubeErr.Kind != ErrorKindNetworkBlocked {
		t.Fatalf("error = %v, want network_blocked detail", err)
	}
	if stt.called {
		t.Fatal("STT should not run when yt-dlp already observed captions")
	}
}

func TestYTDLPBackendDownloadsAudioWithConfiguredNetworkArgs(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		audioExt:  "mp3",
		audioBody: "audio-bytes",
	})
	backend := newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{
		YTDLPPath:   "custom-yt-dlp",
		Proxy:       "http://proxy.internal:8080",
		CookiesFile: "/tmp/youtube-cookies.txt",
		UserAgent:   "kb-test-agent",
	}, retryPolicy{Attempts: 4, Backoff: 2 * time.Second})

	audio, err := backend.downloadAudio(context.Background(), parsedRickRollURL().CanonicalURL, "mp3")
	if err != nil {
		t.Fatalf("downloadAudio returned error: %v", err)
	}
	defer audio.Cleanup()
	data, err := os.ReadFile(audio.Path)
	if err != nil {
		t.Fatalf("read audio: %v", err)
	}
	if string(data) != "audio-bytes\n" {
		t.Fatalf("audio data = %q", string(data))
	}
	if audio.Format != "mp3" {
		t.Fatalf("audio format = %q, want mp3", audio.Format)
	}

	invocations := readYTDLPInvocationLog(t, logPath)
	if len(invocations) != 1 {
		t.Fatalf("invocations = %#v, want one audio invocation", invocations)
	}
	assertArgsContain(t, invocations[0], "--extract-audio")
	assertArgsContain(t, invocations[0], "--audio-format", "mp3")
	assertArgsContain(t, invocations[0], "--proxy", "http://proxy.internal:8080")
	assertArgsContain(t, invocations[0], "--cookies", "/tmp/youtube-cookies.txt")
	assertArgsContain(t, invocations[0], "--user-agent", "kb-test-agent")
	assertArgsContain(t, invocations[0], "--retries", "4")
	assertArgsContain(t, invocations[0], "--fragment-retries", "4")
	assertArgsContain(t, invocations[0], "--retry-sleep", "2")
}

func TestFindYTDLPAudioFileRejectsUnsupportedFormats(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audio.flac"), []byte("audio"), 0o644); err != nil {
		t.Fatalf("write audio: %v", err)
	}
	_, _, err := findYTDLPAudioFile(dir)
	if err == nil {
		t.Fatal("expected unsupported audio format to fail")
	}
	if !strings.Contains(err.Error(), "no audio file") {
		t.Fatalf("error = %v, want no audio file produced", err)
	}
}

func TestExtractorReportsYTDLPMetadataFailure(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataExit: 1,
		metadataErr:  "yt-dlp protocol failed",
	})
	extractor := &Extractor{
		ytDLP: newFakeYTDLPBackend(scriptPath, config.YouTubeConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
	}

	_, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
	if err == nil {
		t.Fatal("expected Extract to fail")
	}
	message := err.Error()
	if !strings.Contains(message, "yt-dlp backend") || strings.Contains(message, "legacy") {
		t.Fatalf("error = %q, want only yt-dlp diagnostics", message)
	}
}

type fakeYTDLPOptions struct {
	metadataJSON string
	metadataExit int
	metadataErr  string
	captionExt   string
	captionBody  string
	captionExit  int
	captionErr   string
	audioExt     string
	audioBody    string
	audioExit    int
	audioErr     string
}

func writeFakeYTDLP(t *testing.T, options fakeYTDLPOptions) (string, string) {
	t.Helper()

	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "yt-dlp")
	logPath := filepath.Join(dir, "args.log")

	var builder strings.Builder
	builder.WriteString("#!/bin/sh\n")
	builder.WriteString("set -eu\n")
	builder.WriteString("LOG_PATH=" + shellQuoteYTDLPTest(logPath) + "\n")
	builder.WriteString("for arg in \"$@\"; do printf '%s\\n' \"$arg\" >> \"$LOG_PATH\"; done\n")
	builder.WriteString("printf '%s\\n' '---' >> \"$LOG_PATH\"\n")
	builder.WriteString("case \" $* \" in\n")
	builder.WriteString("  *\" --dump-single-json \"*)\n")
	if options.metadataErr != "" {
		builder.WriteString("    printf '%s\\n' " + shellQuoteYTDLPTest(options.metadataErr) + " >&2\n")
	}
	if options.metadataExit != 0 {
		builder.WriteString("    exit " + strconv.Itoa(options.metadataExit) + "\n")
	} else {
		builder.WriteString("    cat <<'JSON'\n")
		builder.WriteString(options.metadataJSON)
		if !strings.HasSuffix(options.metadataJSON, "\n") {
			builder.WriteString("\n")
		}
		builder.WriteString("JSON\n")
	}
	builder.WriteString("    ;;\n")
	builder.WriteString("  *\" --extract-audio \"*)\n")
	if options.audioErr != "" {
		builder.WriteString("    printf '%s\\n' " + shellQuoteYTDLPTest(options.audioErr) + " >&2\n")
	}
	if options.audioExit != 0 {
		builder.WriteString("    exit " + strconv.Itoa(options.audioExit) + "\n")
	} else if options.audioExt != "" {
		builder.WriteString("    out_dir=\"\"\n")
		builder.WriteString("    previous=\"\"\n")
		builder.WriteString("    for arg in \"$@\"; do\n")
		builder.WriteString("      if [ \"$previous\" = \"--paths\" ]; then out_dir=${arg#home:}; fi\n")
		builder.WriteString("      previous=\"$arg\"\n")
		builder.WriteString("    done\n")
		builder.WriteString("    [ -n \"$out_dir\" ] || exit 7\n")
		builder.WriteString("    mkdir -p \"$out_dir\"\n")
		builder.WriteString("    cat > \"$out_dir/dQw4w9WgXcQ." + options.audioExt + "\" <<'AUDIO'\n")
		builder.WriteString(options.audioBody)
		if !strings.HasSuffix(options.audioBody, "\n") {
			builder.WriteString("\n")
		}
		builder.WriteString("AUDIO\n")
	}
	builder.WriteString("    ;;\n")
	builder.WriteString("  *)\n")
	if options.captionErr != "" {
		builder.WriteString("    printf '%s\\n' " + shellQuoteYTDLPTest(options.captionErr) + " >&2\n")
	}
	if options.captionExit != 0 {
		builder.WriteString("    exit " + strconv.Itoa(options.captionExit) + "\n")
	} else if options.captionExt != "" {
		builder.WriteString("    out_dir=\"\"\n")
		builder.WriteString("    previous=\"\"\n")
		builder.WriteString("    for arg in \"$@\"; do\n")
		builder.WriteString("      if [ \"$previous\" = \"--paths\" ]; then out_dir=${arg#home:}; fi\n")
		builder.WriteString("      previous=\"$arg\"\n")
		builder.WriteString("    done\n")
		builder.WriteString("    [ -n \"$out_dir\" ] || exit 7\n")
		builder.WriteString("    mkdir -p \"$out_dir\"\n")
		builder.WriteString("    cat > \"$out_dir/dQw4w9WgXcQ.en." + options.captionExt + "\" <<'CAPTION'\n")
		builder.WriteString(options.captionBody)
		if !strings.HasSuffix(options.captionBody, "\n") {
			builder.WriteString("\n")
		}
		builder.WriteString("CAPTION\n")
	}
	builder.WriteString("    ;;\n")
	builder.WriteString("esac\n")

	if err := os.WriteFile(scriptPath, []byte(builder.String()), 0o755); err != nil {
		t.Fatalf("write fake yt-dlp: %v", err)
	}
	return scriptPath, logPath
}

func newFakeYTDLPBackend(scriptPath string, cfg config.YouTubeConfig, retry retryPolicy) *ytDLPBackend {
	backend := newYTDLPBackend(cfg, retry)
	backend.lookPath = func(string) (string, error) {
		return scriptPath, nil
	}
	return backend
}

func readYTDLPInvocationLog(t *testing.T, path string) [][]string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read invocation log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	invocations := make([][]string, 0, 2)
	current := make([]string, 0, 16)
	for _, line := range lines {
		if line == "---" {
			invocations = append(invocations, append([]string(nil), current...))
			current = current[:0]
			continue
		}
		current = append(current, line)
	}
	if len(current) > 0 {
		invocations = append(invocations, current)
	}
	return invocations
}

func parsedRickRollURL() parsedVideoURL {
	return parsedVideoURL{
		CanonicalURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		VideoID:      "dQw4w9WgXcQ",
	}
}

func assertArgsContain(t *testing.T, args []string, want ...string) {
	t.Helper()

	if len(want) == 1 {
		if !containsArg(args, want[0]) {
			t.Fatalf("args %#v do not contain %q", args, want[0])
		}
		return
	}
	for index := 0; index <= len(args)-len(want); index++ {
		if reflect.DeepEqual(args[index:index+len(want)], want) {
			return
		}
	}
	t.Fatalf("args %#v do not contain sequence %#v", args, want)
}

func containsArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}

func shellQuoteYTDLPTest(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
