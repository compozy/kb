package mediadl

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
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
			`"channel_id":"UC123456789",`,
			`"uploader_id":"@ExampleChannel",`,
			`"channel_follower_count":20000,`,
			`"duration":6441,`,
			`"duration_string":"1:47:21",`,
			`"upload_date":"20240307",`,
			`"view_count":3271,`,
			`"like_count":77,`,
			`"comment_count":11,`,
			`"categories":["Science & Technology"],`,
			`"tags":["go"," distributed systems ",""],`,
			`"language":"en",`,
			`"live_status":"not_live",`,
			`"was_live":false,`,
			`"chapters":[{"title":"Intro"},{"title":"Deep dive"}],`,
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
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{
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
	if result.Metadata.ChannelID != "UC123456789" {
		t.Fatalf("channel id = %q, want UC123456789", result.Metadata.ChannelID)
	}
	if result.Metadata.UploaderID != "@ExampleChannel" {
		t.Fatalf("uploader id = %q, want @ExampleChannel", result.Metadata.UploaderID)
	}
	if result.Metadata.Duration != 6441*time.Second {
		t.Fatalf("duration = %v, want 6441s", result.Metadata.Duration)
	}
	if result.Metadata.DurationString != "1:47:21" {
		t.Fatalf("duration string = %q, want 1:47:21", result.Metadata.DurationString)
	}
	if result.Metadata.PublishDate != time.Date(2024, time.March, 7, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("publish date = %v", result.Metadata.PublishDate)
	}
	assertInt64Ptr(t, result.Metadata.ViewCount, 3271)
	assertInt64Ptr(t, result.Metadata.LikeCount, 77)
	assertInt64Ptr(t, result.Metadata.CommentCount, 11)
	assertInt64Ptr(t, result.Metadata.ChannelFollowerCount, 20000)
	if !reflect.DeepEqual(result.Metadata.Categories, []string{"Science & Technology"}) {
		t.Fatalf("categories = %#v", result.Metadata.Categories)
	}
	if !reflect.DeepEqual(result.Metadata.VideoTags, []string{"go", "distributed systems"}) {
		t.Fatalf("video tags = %#v", result.Metadata.VideoTags)
	}
	if result.Metadata.Language != "en" {
		t.Fatalf("language = %q, want en", result.Metadata.Language)
	}
	if result.Metadata.LiveStatus != "not_live" {
		t.Fatalf("live status = %q, want not_live", result.Metadata.LiveStatus)
	}
	assertBoolPtr(t, result.Metadata.WasLive, false)
	if result.Metadata.ChapterCount != 2 {
		t.Fatalf("chapter count = %d, want 2", result.Metadata.ChapterCount)
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

func TestMetadataFromYTDLPInfoHandlesMissingOptionalMetrics(t *testing.T) {
	t.Parallel()

	var info ytDLPInfo
	if err := json.Unmarshal([]byte(strings.Join([]string{
		`{"id":"dQw4w9WgXcQ",`,
		`"title":"Sparse Video",`,
		`"duration":90,`,
		`"like_count":null,`,
		`"comment_count":null,`,
		`"channel_follower_count":null,`,
		`"was_live":null,`,
		`"categories":["", "Education"],`,
		`"tags":null}`,
	}, "")), &info); err != nil {
		t.Fatalf("unmarshal sparse yt-dlp info: %v", err)
	}

	metadata := metadataFromYTDLPInfo(parsedRickRollURL(), info)
	if metadata.DurationString != "1:30" {
		t.Fatalf("duration string = %q, want fallback 1:30", metadata.DurationString)
	}
	if metadata.LikeCount != nil {
		t.Fatalf("like count = %#v, want nil", metadata.LikeCount)
	}
	if metadata.CommentCount != nil {
		t.Fatalf("comment count = %#v, want nil", metadata.CommentCount)
	}
	if metadata.ChannelFollowerCount != nil {
		t.Fatalf("channel follower count = %#v, want nil", metadata.ChannelFollowerCount)
	}
	if metadata.WasLive != nil {
		t.Fatalf("was live = %#v, want nil", metadata.WasLive)
	}
	if !reflect.DeepEqual(metadata.Categories, []string{"Education"}) {
		t.Fatalf("categories = %#v, want trimmed non-empty values", metadata.Categories)
	}
	if !reflect.DeepEqual(metadata.VideoTags, []string{}) {
		t.Fatalf("video tags = %#v, want empty list", metadata.VideoTags)
	}
}

func TestMetadataFromYTDLPInfoCapturesDescription(t *testing.T) {
	t.Parallel()

	var info ytDLPInfo
	if err := json.Unmarshal([]byte(`{"id":"abc","title":"Reel","description":"  Caption body  ","duration":12}`), &info); err != nil {
		t.Fatalf("unmarshal yt-dlp info: %v", err)
	}
	metadata := metadataFromYTDLPInfo(ParsedURL{CanonicalURL: "https://www.instagram.com/reel/abc/", VideoID: "abc"}, info)
	if metadata.Description != "Caption body" {
		t.Fatalf("description = %q, want trimmed caption body", metadata.Description)
	}
}

func TestYTDLPBackendListsPlaylistEntries(t *testing.T) {
	t.Parallel()

	playlist := `{"title":"Asimov Academy Videos","channel":"Asimov Academy","uploader":"Asimov","entries":[
  {"id":"aaaaaaaaaaa","title":"First","url":"https://www.youtube.com/watch?v=aaaaaaaaaaa"},
  {"id":"bbbbbbbbbbb","title":"Second"}
]}`
	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{metadataJSON: playlist})
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{}, retryPolicy{Attempts: 1})

	listing, err := backend.ListPlaylist(context.Background(), "https://www.youtube.com/@chan/videos", 5)
	if err != nil {
		t.Fatalf("ListPlaylist returned error: %v", err)
	}
	if listing.Title != "Asimov Academy Videos" || listing.Channel != "Asimov Academy" || listing.Uploader != "Asimov" {
		t.Fatalf("listing metadata = %+v, want title/channel/uploader from yt-dlp JSON", listing)
	}
	entries := listing.Entries
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d: %+v", len(entries), entries)
	}
	if entries[0].ID != "aaaaaaaaaaa" || entries[0].Title != "First" {
		t.Fatalf("unexpected first entry: %+v", entries[0])
	}
	if entries[0].URL != "https://www.youtube.com/watch?v=aaaaaaaaaaa" {
		t.Fatalf("first entry url = %q", entries[0].URL)
	}

	args := readYTDLPInvocationLog(t, logPath)[0]
	assertArgsContain(t, args, "--flat-playlist")
	assertArgsContain(t, args, "--playlist-end", "5")
	if containsArg(args, "--no-playlist") {
		t.Fatalf("--no-playlist must be absent for channel enumeration: %#v", args)
	}
}

func TestYTDLPBackendExtractsVTTFallback(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"VTT Video",`,
			`"language":"en",`,
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
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

	result, err := backend.Extract(context.Background(), parsedRickRollURL(), nil)
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}

	wantMarkdown := "## 00:01\nFirst caption\n\n## 00:02\nSecond caption"
	if result.Markdown != wantMarkdown {
		t.Fatalf("markdown = %q, want %q", result.Markdown, wantMarkdown)
	}
}

func TestYTDLPBackendDefaultsToOriginalAutomaticCaption(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"Portuguese",`,
			`"language":"pt",`,
			`"subtitles":{},`,
			`"automatic_captions":{"pt-orig":[{"ext":"json3"}],"pt":[{"ext":"json3"}],"en":[{"ext":"json3"}],"es":[{"ext":"json3"}],"fr":[{"ext":"json3"}]}}`,
		}, ""),
		captionExt:  "json3",
		captionBody: `{"events":[{"tStartMs":0,"segs":[{"utf8":"Ola mundo"}]}]}`,
	})
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

	result, err := backend.Extract(context.Background(), parsedRickRollURL(), nil)
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Language != "pt-orig" {
		t.Fatalf("language = %q, want pt-orig", result.Language)
	}
	invocations := readYTDLPInvocationLog(t, logPath)
	assertArgsContain(t, invocations[1], "--write-auto-subs")
	assertArgsContain(t, invocations[1], "--sub-langs", "pt-orig")
	if containsArgSequence(invocations[1], "--sub-langs", "en") {
		t.Fatalf("default selection must not request translated English: %#v", invocations[1])
	}
}

func TestYTDLPBackendPrefersManualOriginalCaption(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"Manual Portuguese",`,
			`"language":"pt",`,
			`"subtitles":{"pt":[{"ext":"json3"}]},`,
			`"automatic_captions":{"pt-orig":[{"ext":"json3"}],"en":[{"ext":"json3"}]}}`,
		}, ""),
		captionExt:  "json3",
		captionBody: `{"events":[{"tStartMs":0,"segs":[{"utf8":"Manual"}]}]}`,
	})
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

	result, err := backend.Extract(context.Background(), parsedRickRollURL(), nil)
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Language != "pt" {
		t.Fatalf("language = %q, want pt", result.Language)
	}
	invocations := readYTDLPInvocationLog(t, logPath)
	assertArgsContain(t, invocations[1], "--write-subs")
	assertArgsContain(t, invocations[1], "--sub-langs", "pt")
}

func TestYTDLPBackendDefaultsToOriginalEnglishCaption(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"English",`,
			`"language":"en",`,
			`"subtitles":{},`,
			`"automatic_captions":{"en-orig":[{"ext":"json3"}],"en":[{"ext":"json3"}],"pt":[{"ext":"json3"}]}}`,
		}, ""),
		captionExt:  "json3",
		captionBody: `{"events":[{"tStartMs":0,"segs":[{"utf8":"Hello"}]}]}`,
	})
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

	result, err := backend.Extract(context.Background(), parsedRickRollURL(), nil)
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Language != "en-orig" {
		t.Fatalf("language = %q, want en-orig", result.Language)
	}
	invocations := readYTDLPInvocationLog(t, logPath)
	assertArgsContain(t, invocations[1], "--sub-langs", "en-orig")
}

func TestYTDLPBackendSelectsPreferredAutomaticCaption(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"Preferred",`,
			`"language":"pt-BR",`,
			`"subtitles":{"en":[{"ext":"json3"}]},`,
			`"automatic_captions":{"pt-BR":[{"ext":"json3"}]}}`,
		}, ""),
		captionExt:  "json3",
		captionBody: `{"events":[{"tStartMs":0,"segs":[{"utf8":"Preferido"}]}]}`,
	})
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

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

func TestExtractorRejectsTranslatedCaptionWhenDisabled(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"Portuguese",`,
			`"language":"pt",`,
			`"subtitles":{},`,
			`"automatic_captions":{"pt-orig":[{"ext":"json3"}],"es":[{"ext":"json3"}]}}`,
		}, ""),
		captionExit: 1,
		captionErr:  "translated caption should not be downloaded",
	})
	extractor := &Extractor{ytDLP: newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})}

	_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
		PreferredLanguages: []string{"es"},
	})
	if err == nil {
		t.Fatal("expected translated caption selection to fail")
	}
	var mediaErr *Error
	if !errors.As(err, &mediaErr) || mediaErr.Kind != ErrorKindTranscriptUnavailable {
		t.Fatalf("error = %v, want transcript_unavailable", err)
	}
	if !strings.Contains(err.Error(), "translated captions are disabled") {
		t.Fatalf("error = %v, want translated captions diagnostic", err)
	}
}

func TestExtractorAllowsExplicitTranslatedCaptionWhenEnabled(t *testing.T) {
	t.Parallel()

	scriptPath, logPath := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"Portuguese",`,
			`"language":"pt",`,
			`"subtitles":{},`,
			`"automatic_captions":{"pt-orig":[{"ext":"json3"}],"es":[{"ext":"json3"}]}}`,
		}, ""),
		captionExt:  "json3",
		captionBody: `{"events":[{"tStartMs":0,"segs":[{"utf8":"Hola"}]}]}`,
	})
	extractor := &Extractor{ytDLP: newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})}

	result, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
		PreferredLanguages:      []string{"es"},
		AllowTranslatedCaptions: true,
	})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Language != "es" {
		t.Fatalf("language = %q, want es", result.Language)
	}
	invocations := readYTDLPInvocationLog(t, logPath)
	assertArgsContain(t, invocations[1], "--sub-langs", "es")
}

func TestExtractorFailsWhenNoOriginalCaptionAvailable(t *testing.T) {
	t.Parallel()

	scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
		metadataJSON: strings.Join([]string{
			`{"id":"dQw4w9WgXcQ","title":"Portuguese",`,
			`"language":"pt",`,
			`"subtitles":{},`,
			`"automatic_captions":{"en":[{"ext":"json3"}],"es":[{"ext":"json3"}]}}`,
		}, ""),
		captionExit: 1,
		captionErr:  "caption download should not run",
	})
	extractor := &Extractor{ytDLP: newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})}

	_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{})
	if err == nil {
		t.Fatal("expected missing original caption to fail")
	}
	var mediaErr *Error
	if !errors.As(err, &mediaErr) || mediaErr.Kind != ErrorKindTranscriptUnavailable {
		t.Fatalf("error = %v, want transcript_unavailable", err)
	}
	if !strings.Contains(err.Error(), "no original-language caption available") {
		t.Fatalf("error = %v, want original-language diagnostic", err)
	}
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

	_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{})
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
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Primary","language":"en","subtitles":{"en":[{"ext":"json3"}]}}`,
		captionExt:   "json3",
		captionBody:  `{"events":[{"tStartMs":0,"segs":[{"utf8":"Primary transcript"}]}]}`,
	})
	extractor := &Extractor{
		ytDLP: newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
	}

	result, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{})
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
		ytDLP:     newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:       stt,
		sttConfig: config.Default().STT,
	}

	result, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
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
		ytDLP:     newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:       stt,
		sttConfig: config.Default().STT,
	}

	result, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
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
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Captioned","language":"en","subtitles":{"en":[{"ext":"json3"}]},"automatic_captions":{}}`,
		audioExt:     "mp3",
		audioBody:    "audio-from-ytdlp",
		captionExit:  1,
		captionErr:   "caption download should not run",
	})
	stt := &stubSTTClient{transcript: "Forced STT transcript"}
	extractor := &Extractor{
		ytDLP:     newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:       stt,
		sttConfig: config.Default().STT,
	}

	result, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
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
		metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Captioned","language":"en","subtitles":{"en":[{"ext":"json3"}]},"automatic_captions":{}}`,
		captionExit:  1,
		captionErr:   "HTTP Error 429: Too Many Requests",
	})
	stt := &stubSTTClient{transcript: "should not run"}
	extractor := &Extractor{
		ytDLP: newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		stt:   stt,
	}

	_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{})
	if err == nil {
		t.Fatal("expected Extract to fail")
	}
	if !errors.Is(err, errYTDLPCaptionFetchFailed) {
		t.Fatalf("error = %v, want yt-dlp caption fetch failure", err)
	}
	var mediaErr *Error
	if !errors.As(err, &mediaErr) || mediaErr.Kind != ErrorKindRateLimited {
		t.Fatalf("error = %v, want rate_limited detail", err)
	}
	if !strings.Contains(err.Error(), `caption language "en"`) {
		t.Fatalf("error = %v, want requested caption language", err)
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
	backend := newFakeYTDLPBackend(scriptPath, BackendConfig{
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
		ytDLP: newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
	}

	_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{})
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

func newFakeYTDLPBackend(scriptPath string, cfg BackendConfig, retry retryPolicy) *ytDLPBackend {
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

func parsedRickRollURL() ParsedURL {
	return ParsedURL{
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
	return slices.Contains(args, want)
}

func containsArgSequence(args []string, want ...string) bool {
	for index := 0; index <= len(args)-len(want); index++ {
		if reflect.DeepEqual(args[index:index+len(want)], want) {
			return true
		}
	}
	return false
}

func assertInt64Ptr(t *testing.T, got *int64, want int64) {
	t.Helper()

	if got == nil {
		t.Fatalf("value = nil, want %d", want)
	}
	if *got != want {
		t.Fatalf("value = %d, want %d", *got, want)
	}
}

func assertBoolPtr(t *testing.T, got *bool, want bool) {
	t.Helper()

	if got == nil {
		t.Fatalf("value = nil, want %t", want)
	}
	if *got != want {
		t.Fatalf("value = %t, want %t", *got, want)
	}
}

func shellQuoteYTDLPTest(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
