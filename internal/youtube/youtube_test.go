package youtube

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	ytdl "github.com/kkdai/youtube/v2"

	"github.com/compozy/kb/internal/config"
)

func TestExtractorExtractReturnsMetadataAndCaptionMarkdown(t *testing.T) {
	t.Parallel()

	video := &ytdl.Video{
		ID:          "dQw4w9WgXcQ",
		Title:       "Example Video",
		Author:      "Example Channel",
		Duration:    95 * time.Second,
		PublishDate: time.Date(2024, time.March, 7, 14, 30, 0, 0, time.UTC),
		CaptionTracks: []ytdl.CaptionTrack{
			{LanguageCode: "en", Kind: ""},
		},
	}

	client := &stubYouTubeClient{
		video: video,
		transcripts: map[string]ytdl.VideoTranscript{
			"en": {
				{StartMs: 0, Text: " Hello   world "},
				{StartMs: 5_000, Text: "Second line"},
			},
		},
	}

	extractor := &Extractor{youtube: client}
	result, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}

	if client.videoURL != "https://www.youtube.com/watch?v=dQw4w9WgXcQ" {
		t.Fatalf("video URL = %q", client.videoURL)
	}

	if result.Metadata.VideoID != "dQw4w9WgXcQ" {
		t.Fatalf("video ID = %q, want %q", result.Metadata.VideoID, "dQw4w9WgXcQ")
	}
	if result.Metadata.Title != "Example Video" {
		t.Fatalf("title = %q, want %q", result.Metadata.Title, "Example Video")
	}
	if result.Metadata.Channel != "Example Channel" {
		t.Fatalf("channel = %q, want %q", result.Metadata.Channel, "Example Channel")
	}
	if result.Metadata.Duration != 95*time.Second {
		t.Fatalf("duration = %v, want %v", result.Metadata.Duration, 95*time.Second)
	}
	if !result.Metadata.PublishDate.Equal(video.PublishDate) {
		t.Fatalf("publish date = %v, want %v", result.Metadata.PublishDate, video.PublishDate)
	}
	if result.Source != TranscriptSourceCaptions {
		t.Fatalf("source = %q, want %q", result.Source, TranscriptSourceCaptions)
	}
	if result.Language != "en" {
		t.Fatalf("language = %q, want %q", result.Language, "en")
	}

	wantMarkdown := strings.Join([]string{
		"## 00:00",
		"Hello world",
		"",
		"## 00:05",
		"Second line",
	}, "\n")
	if result.Markdown != wantMarkdown {
		t.Fatalf("markdown = %q, want %q", result.Markdown, wantMarkdown)
	}
}

func TestExtractorFallsBackToSTTWhenTranscriptUnavailableAndAllowed(t *testing.T) {
	t.Parallel()

	video := &ytdl.Video{
		ID:     "dQw4w9WgXcQ",
		Title:  "Fallback Video",
		Author: "Example Channel",
		CaptionTracks: []ytdl.CaptionTrack{
			{LanguageCode: "en", Kind: "asr"},
		},
		Formats: ytdl.FormatList{
			{MimeType: `audio/mp4; codecs="mp4a.40.2"`, AudioChannels: 2},
		},
	}

	client := &stubYouTubeClient{
		video:          video,
		transcriptErrs: map[string]error{"en": ytdl.ErrTranscriptDisabled},
		audioData:      []byte("audio-bytes"),
	}
	stt := &stubSTTClient{
		configured: true,
		transcript: "Hello from fallback",
	}

	extractor := &Extractor{
		youtube: client,
		stt:     stt,
	}

	result, err := extractor.Extract(context.Background(), "https://youtu.be/dQw4w9WgXcQ", ExtractOptions{})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}

	if !stt.called {
		t.Fatal("expected STT fallback to be invoked")
	}
	if stt.format != "m4a" {
		t.Fatalf("format = %q, want %q", stt.format, "m4a")
	}
	if !reflect.DeepEqual(stt.audio, []byte("audio-bytes")) {
		t.Fatalf("audio = %q, want audio-bytes", string(stt.audio))
	}
	if result.Source != TranscriptSourceSTT {
		t.Fatalf("source = %q, want %q", result.Source, TranscriptSourceSTT)
	}
	if result.Markdown != "## 00:00\nHello from fallback" {
		t.Fatalf("markdown = %q", result.Markdown)
	}
}

func TestExtractorDoesNotInvokeSTTWhenNotAllowed(t *testing.T) {
	t.Parallel()

	video := &ytdl.Video{
		ID:            "dQw4w9WgXcQ",
		Title:         "Captions Missing",
		CaptionTracks: []ytdl.CaptionTrack{{LanguageCode: "en"}},
		Formats: ytdl.FormatList{
			{MimeType: `audio/mp4; codecs="mp4a.40.2"`, AudioChannels: 2},
		},
	}

	client := &stubYouTubeClient{
		video:          video,
		transcriptErrs: map[string]error{"en": ytdl.ErrTranscriptDisabled},
		audioData:      []byte("audio-bytes"),
	}
	stt := &stubSTTClient{configured: false}

	extractor := &Extractor{
		youtube: client,
		stt:     stt,
	}

	result, err := extractor.Extract(context.Background(), "https://www.youtube.com/shorts/dQw4w9WgXcQ", ExtractOptions{})
	if err == nil {
		t.Fatal("expected Extract to fail")
	}
	if result == nil {
		t.Fatal("expected partial result with metadata")
	}
	if stt.called {
		t.Fatal("did not expect STT fallback to run")
	}

	var youtubeErr *Error
	if !errors.As(err, &youtubeErr) {
		t.Fatalf("expected structured error, got %T", err)
	}
	if youtubeErr.Kind != ErrorKindTranscriptUnavailable {
		t.Fatalf("kind = %q, want %q", youtubeErr.Kind, ErrorKindTranscriptUnavailable)
	}
}

func TestExtractorReturnsStructuredErrorsForVideoAccessFailures(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		err      error
		wantKind ErrorKind
	}{
		{name: "private", err: ytdl.ErrVideoPrivate, wantKind: ErrorKindPrivate},
		{name: "age restricted", err: ytdl.ErrLoginRequired, wantKind: ErrorKindAgeRestricted},
		{name: "unavailable", err: &ytdl.ErrPlayabiltyStatus{Status: "ERROR", Reason: "Video unavailable"}, wantKind: ErrorKindUnavailable},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			extractor := &Extractor{
				youtube: &stubYouTubeClient{videoErr: tc.err},
			}

			_, err := extractor.Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
			if err == nil {
				t.Fatal("expected Extract to fail")
			}

			var youtubeErr *Error
			if !errors.As(err, &youtubeErr) {
				t.Fatalf("expected structured error, got %T", err)
			}
			if youtubeErr.Kind != tc.wantKind {
				t.Fatalf("kind = %q, want %q", youtubeErr.Kind, tc.wantKind)
			}
		})
	}
}

func TestExtractorRejectsInvalidYouTubeURL(t *testing.T) {
	t.Parallel()

	extractor := &Extractor{youtube: &stubYouTubeClient{}}

	_, err := extractor.Extract(context.Background(), "https://example.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
	if err == nil {
		t.Fatal("expected Extract to fail")
	}

	var youtubeErr *Error
	if !errors.As(err, &youtubeErr) {
		t.Fatalf("expected structured error, got %T", err)
	}
	if youtubeErr.Kind != ErrorKindInvalidURL {
		t.Fatalf("kind = %q, want %q", youtubeErr.Kind, ErrorKindInvalidURL)
	}
}

func TestParseVideoURLHandlesCommonYouTubeFormats(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
	}{
		{name: "watch", raw: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
		{name: "short", raw: "https://youtu.be/dQw4w9WgXcQ"},
		{name: "shorts", raw: "https://www.youtube.com/shorts/dQw4w9WgXcQ"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := parseVideoURL(tc.raw)
			if err != nil {
				t.Fatalf("parseVideoURL returned error: %v", err)
			}
			if parsed.VideoID != "dQw4w9WgXcQ" {
				t.Fatalf("video ID = %q, want %q", parsed.VideoID, "dQw4w9WgXcQ")
			}
			if parsed.CanonicalURL != "https://www.youtube.com/watch?v=dQw4w9WgXcQ" {
				t.Fatalf("canonical url = %q", parsed.CanonicalURL)
			}
		})
	}
}

func TestFormatTranscriptMarkdownUsesTimestampHeaders(t *testing.T) {
	t.Parallel()

	transcript := ytdl.VideoTranscript{
		{StartMs: 0, Text: "first segment"},
		{StartMs: 3_723_000, Text: "later segment"},
	}

	got := formatTranscriptMarkdown(transcript)
	want := strings.Join([]string{
		"## 00:00",
		"first segment",
		"",
		"## 01:02:03",
		"later segment",
	}, "\n")
	if got != want {
		t.Fatalf("markdown = %q, want %q", got, want)
	}
}

func TestNewExtractorConstructsDefaultClients(t *testing.T) {
	t.Parallel()

	extractor := NewExtractor(config.OpenRouterConfig{})
	if extractor == nil {
		t.Fatal("expected extractor")
	}
	if extractor.youtube == nil {
		t.Fatal("expected YouTube client")
	}
	if extractor.stt == nil {
		t.Fatal("expected STT client")
	}
}

func TestErrorFormattingAndUnwrap(t *testing.T) {
	t.Parallel()

	cause := errors.New("boom")
	err := &Error{
		Kind:    ErrorKindPrivate,
		URL:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		Message: "video is private",
		Err:     cause,
	}

	if !strings.Contains(err.Error(), "youtube private") {
		t.Fatalf("error = %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Fatal("expected wrapped error")
	}
}

func TestShouldAttemptSTT(t *testing.T) {
	t.Parallel()

	if (&Extractor{}).shouldAttemptSTT(ExtractOptions{}) {
		t.Fatal("expected false when STT client is missing")
	}

	extractor := &Extractor{stt: &stubSTTClient{configured: false}}
	if !extractor.shouldAttemptSTT(ExtractOptions{EnableSTTFallback: true}) {
		t.Fatal("expected explicit fallback flag to enable STT")
	}

	extractor = &Extractor{stt: &stubSTTClient{configured: true}}
	if !extractor.shouldAttemptSTT(ExtractOptions{}) {
		t.Fatal("expected configured client to enable STT")
	}
}

func TestOrderedCaptionTracksPrefersPreferredLanguageAndManualTracks(t *testing.T) {
	t.Parallel()

	tracks := []ytdl.CaptionTrack{
		{LanguageCode: "fr", Kind: "asr"},
		{LanguageCode: "en", Kind: "asr"},
		{LanguageCode: "en-US", Kind: ""},
		{LanguageCode: "de", Kind: ""},
	}

	got := orderedCaptionTracks(tracks, []string{" en ", "en"})
	if len(got) != 4 {
		t.Fatalf("ordered length = %d, want 4", len(got))
	}
	if got[0].LanguageCode != "en-US" {
		t.Fatalf("first language = %q, want %q", got[0].LanguageCode, "en-US")
	}
	if got[1].LanguageCode != "en" {
		t.Fatalf("second language = %q, want %q", got[1].LanguageCode, "en")
	}
	if !languageMatches("en-US", "en") {
		t.Fatal("expected language prefix match")
	}
}

func TestPickAudioFormatPrefersSupportedAudioMimeTypes(t *testing.T) {
	t.Parallel()

	format, normalized, err := pickAudioFormat(ytdl.FormatList{
		{MimeType: `audio/webm; codecs="opus"`, AudioChannels: 2, URL: "https://cdn.example.com/audio.webm"},
		{MimeType: `audio/mp4; codecs="mp4a.40.2"`, AudioChannels: 2, URL: "https://cdn.example.com/audio.m4a"},
	})
	if err != nil {
		t.Fatalf("pickAudioFormat returned error: %v", err)
	}
	if normalized != "m4a" {
		t.Fatalf("normalized format = %q, want %q", normalized, "m4a")
	}
	if !strings.Contains(format.MimeType, "audio/mp4") {
		t.Fatalf("mime type = %q, want audio/mp4", format.MimeType)
	}

	if got := normalizeAudioFormat("", "https://cdn.example.com/audio.mp3"); got != "mp3" {
		t.Fatalf("normalizeAudioFormat by extension = %q, want %q", got, "mp3")
	}
	if got := normalizeAudioFormat(`audio/ogg; codecs="opus"`, ""); got != "ogg" {
		t.Fatalf("normalizeAudioFormat by mime = %q, want %q", got, "ogg")
	}
}

func TestPickAudioFormatReturnsStructuredErrorWhenAudioMissing(t *testing.T) {
	t.Parallel()

	_, _, err := pickAudioFormat(ytdl.FormatList{
		{MimeType: `video/mp4; codecs="avc1.640028"`},
	})
	if err == nil {
		t.Fatal("expected pickAudioFormat to fail")
	}

	var youtubeErr *Error
	if !errors.As(err, &youtubeErr) {
		t.Fatalf("expected structured error, got %T", err)
	}
	if youtubeErr.Kind != ErrorKindAudioUnavailable {
		t.Fatalf("kind = %q, want %q", youtubeErr.Kind, ErrorKindAudioUnavailable)
	}
}

func TestDownloadAudioReturnsStructuredErrorForEmptyStream(t *testing.T) {
	t.Parallel()

	video := &ytdl.Video{
		ID: "dQw4w9WgXcQ",
		Formats: ytdl.FormatList{
			{MimeType: `audio/mp4; codecs="mp4a.40.2"`, AudioChannels: 2},
		},
	}

	extractor := &Extractor{
		youtube: &stubYouTubeClient{
			audioData: []byte{},
		},
	}

	_, _, err := extractor.downloadAudio(context.Background(), video)
	if err == nil {
		t.Fatal("expected downloadAudio to fail")
	}

	var youtubeErr *Error
	if !errors.As(err, &youtubeErr) {
		t.Fatalf("expected structured error, got %T", err)
	}
	if youtubeErr.Kind != ErrorKindAudioUnavailable {
		t.Fatalf("kind = %q, want %q", youtubeErr.Kind, ErrorKindAudioUnavailable)
	}
}

type stubYouTubeClient struct {
	video          *ytdl.Video
	videoErr       error
	videoURL       string
	transcripts    map[string]ytdl.VideoTranscript
	transcriptErrs map[string]error
	audioData      []byte
}

func (client *stubYouTubeClient) GetVideoContext(_ context.Context, rawURL string) (*ytdl.Video, error) {
	client.videoURL = rawURL
	if client.videoErr != nil {
		return nil, client.videoErr
	}
	if client.video == nil {
		return nil, errors.New("video not configured")
	}
	return client.video, nil
}

func (client *stubYouTubeClient) GetTranscriptCtx(_ context.Context, _ *ytdl.Video, lang string) (ytdl.VideoTranscript, error) {
	if err := client.transcriptErrs[lang]; err != nil {
		return nil, err
	}
	if transcript, ok := client.transcripts[lang]; ok {
		return transcript, nil
	}
	return nil, ytdl.ErrTranscriptDisabled
}

func (client *stubYouTubeClient) GetStreamContext(_ context.Context, _ *ytdl.Video, _ *ytdl.Format) (io.ReadCloser, int64, error) {
	return io.NopCloser(strings.NewReader(string(client.audioData))), int64(len(client.audioData)), nil
}

type stubSTTClient struct {
	configured bool
	called     bool
	audio      []byte
	format     string
	transcript string
	err        error
}

func (client *stubSTTClient) Configured() bool {
	return client.configured
}

func (client *stubSTTClient) Transcribe(_ context.Context, audio []byte, format string) (string, error) {
	client.called = true
	client.audio = append([]byte(nil), audio...)
	client.format = format
	if client.err != nil {
		return "", client.err
	}
	return client.transcript, nil
}
