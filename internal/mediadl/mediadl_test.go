package mediadl

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
)

func TestExtractorRequiresYTDLPBackend(t *testing.T) {
	t.Parallel()

	t.Run("Should fail without yt-dlp backend", func(t *testing.T) {
		t.Parallel()

		_, err := (&Extractor{}).Extract(context.Background(), parsedRickRollURL(), ExtractOptions{})
		if err == nil {
			t.Fatal("expected Extract to fail without yt-dlp backend")
		}
		if !strings.Contains(err.Error(), "yt-dlp backend is required") {
			t.Fatalf("error = %q, want yt-dlp requirement", err.Error())
		}
	})
}

func TestFormatTranscriptMarkdownUsesTimestampHeaders(t *testing.T) {
	t.Parallel()

	t.Run("Should use timestamp headers", func(t *testing.T) {
		t.Parallel()

		transcript := []transcriptSegment{
			{StartMs: 0, Text: " first   segment "},
			{StartMs: 3_723_000, Text: "later segment"},
			{StartMs: 4_000_000, Text: "   "},
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
	})
}

func TestNewExtractorConstructsYTDLPAndSTT(t *testing.T) {
	t.Parallel()

	t.Run("Should construct yt-dlp and stt backends", func(t *testing.T) {
		t.Parallel()

		backendConfig := BackendConfig{
			Platform:      "youtube",
			YTDLPPath:     "/opt/bin/yt-dlp",
			Proxy:         "http://proxy.internal:8080",
			CookiesFile:   "/tmp/youtube-cookies.txt",
			UserAgent:     "kb-test-agent",
			RetryAttempts: 4,
			RetryBackoff:  "250ms",
		}
		extractor := NewExtractor(backendConfig, config.STTConfig{
			Provider:    "openai",
			APIKey:      "openai-key",
			APIURL:      "https://api.openai.test",
			Model:       "gpt-4o-transcribe",
			Language:    "auto",
			AudioFormat: "mp3",
		}, config.OpenRouterConfig{})
		if extractor == nil {
			t.Fatal("expected extractor")
		}
		if extractor.stt == nil {
			t.Fatal("expected STT client")
		}
		if extractor.ytDLP == nil {
			t.Fatal("expected yt-dlp backend")
		}
		if extractor.ytDLP.binaryPath != "/opt/bin/yt-dlp" {
			t.Fatalf("yt-dlp path = %q, want /opt/bin/yt-dlp", extractor.ytDLP.binaryPath)
		}
		if extractor.ytDLP.cfg != backendConfig {
			t.Fatalf("yt-dlp config = %#v, want %#v", extractor.ytDLP.cfg, backendConfig)
		}
		if extractor.ytDLP.retry.Attempts != 4 {
			t.Fatalf("yt-dlp retry attempts = %d, want 4", extractor.ytDLP.retry.Attempts)
		}
		if extractor.ytDLP.retry.Backoff != 250*time.Millisecond {
			t.Fatalf("yt-dlp retry backoff = %v, want 250ms", extractor.ytDLP.retry.Backoff)
		}
		if extractor.setupErr != nil {
			t.Fatalf("setupErr = %v, want nil", extractor.setupErr)
		}
	})
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

	if !strings.Contains(err.Error(), "media private") {
		t.Fatalf("error = %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Fatal("expected wrapped error")
	}
}

func TestShouldAttemptSTT(t *testing.T) {
	t.Parallel()

	if (&Extractor{}).shouldAttemptSTT(TranscriptionPolicyAuto) {
		t.Fatal("expected false when STT client is missing")
	}

	extractor := &Extractor{stt: &stubSTTClient{}}
	if !extractor.shouldAttemptSTT(TranscriptionPolicyAuto) {
		t.Fatal("expected auto policy to request STT when a transcriber is present")
	}
	if !extractor.shouldAttemptSTT(TranscriptionPolicySTT) {
		t.Fatal("expected stt policy to request STT when a transcriber is present")
	}
	if extractor.shouldAttemptSTT(TranscriptionPolicyCaptions) {
		t.Fatal("expected captions policy to avoid STT")
	}
}

func TestExtractorSTTErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("Should fail forced STT without transcriber", func(t *testing.T) {
		t.Parallel()

		scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
			metadataJSON: `{"id":"dQw4w9WgXcQ","title":"No Captions","subtitles":{},"automatic_captions":{}}`,
		})
		extractor := &Extractor{
			ytDLP: newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
		}

		_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
			TranscriptionPolicy: TranscriptionPolicySTT,
		})
		if err == nil {
			t.Fatal("expected forced STT to fail without a transcriber")
		}
		if !strings.Contains(err.Error(), "media stt: transcriber is nil") {
			t.Fatalf("error = %q, want missing transcriber", err.Error())
		}
		if IsTranscriptUnavailable(err) {
			t.Fatalf("error = %v, should not be classified as transcript unavailable", err)
		}
	})

	t.Run("Should not mark audio download failures as transcript unavailable", func(t *testing.T) {
		t.Parallel()

		scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
			metadataJSON: `{"id":"dQw4w9WgXcQ","title":"No Captions","subtitles":{},"automatic_captions":{}}`,
			audioExit:    1,
			audioErr:     "network down",
		})
		extractor := &Extractor{
			ytDLP:     newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
			stt:       &stubSTTClient{transcript: "should not run"},
			sttConfig: config.Default().STT,
		}

		_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
			TranscriptionPolicy: TranscriptionPolicyAuto,
		})
		if err == nil {
			t.Fatal("expected audio download failure")
		}
		if IsTranscriptUnavailable(err) {
			t.Fatalf("error = %v, should not be classified as transcript unavailable", err)
		}
	})

	t.Run("Should not mark generic STT failures as transcript unavailable", func(t *testing.T) {
		t.Parallel()

		scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
			metadataJSON: `{"id":"dQw4w9WgXcQ","title":"No Captions","subtitles":{},"automatic_captions":{}}`,
			audioExt:     "mp3",
			audioBody:    "audio-bytes",
		})
		extractor := &Extractor{
			ytDLP:     newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{}),
			stt:       &stubSTTClient{err: errors.New("stt backend failed")},
			sttConfig: config.Default().STT,
		}

		_, err := extractor.Extract(context.Background(), parsedRickRollURL(), ExtractOptions{
			TranscriptionPolicy: TranscriptionPolicyAuto,
		})
		if err == nil {
			t.Fatal("expected STT backend failure")
		}
		if IsTranscriptUnavailable(err) {
			t.Fatalf("error = %v, should not be classified as transcript unavailable", err)
		}
	})
}

func TestYTDLPDefaultCaptionSelection(t *testing.T) {
	t.Parallel()

	t.Run("Should accept manual captions when original language metadata is missing", func(t *testing.T) {
		t.Parallel()

		scriptPath, _ := writeFakeYTDLP(t, fakeYTDLPOptions{
			metadataJSON: `{"id":"dQw4w9WgXcQ","title":"Manual Caption","subtitles":{"en":[{"ext":"json3"}]},"automatic_captions":{}}`,
			captionExt:   "json3",
			captionBody:  `{"events":[{"tStartMs":0,"segs":[{"utf8":"Manual caption"}]}]}`,
		})
		backend := newFakeYTDLPBackend(scriptPath, BackendConfig{YTDLPPath: "yt-dlp"}, retryPolicy{})

		result, err := backend.Extract(context.Background(), parsedRickRollURL(), nil)
		if err != nil {
			t.Fatalf("Extract returned error: %v", err)
		}
		if result.Language != "en" {
			t.Fatalf("language = %q, want en", result.Language)
		}
		if result.CaptionKind != CaptionKindManual {
			t.Fatalf("caption kind = %q, want manual", result.CaptionKind)
		}
		if result.Markdown != "## 00:00\nManual caption" {
			t.Fatalf("markdown = %q", result.Markdown)
		}
	})
}

func TestParseTranscriptionPolicy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		value           string
		want            TranscriptionPolicy
		wantErrContains string
	}{
		{name: "Should default empty policy to captions", value: "", want: TranscriptionPolicyCaptions},
		{name: "Should parse captions policy", value: "captions", want: TranscriptionPolicyCaptions},
		{name: "Should parse auto policy", value: "auto", want: TranscriptionPolicyAuto},
		{name: "Should parse stt policy", value: "stt", want: TranscriptionPolicySTT},
		{name: "Should trim and normalize policy", value: " STT ", want: TranscriptionPolicySTT},
		{name: "Should reject invalid policy", value: "maybe", wantErrContains: "transcription policy must be captions, auto, or stt"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseTranscriptionPolicy(tc.value)
			if tc.wantErrContains != "" {
				if err == nil {
					t.Fatal("expected invalid transcription policy to fail")
				}
				if !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Fatalf("error = %q, want %q", err.Error(), tc.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseTranscriptionPolicy(%q) returned error: %v", tc.value, err)
			}
			if got != tc.want {
				t.Fatalf("policy = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestLanguageHelpers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		assert func(t *testing.T)
	}{
		{
			name: "Should normalize and deduplicate language tags",
			assert: func(t *testing.T) {
				got := normalizeLanguages([]string{" en ", "EN", "", "pt-BR"})
				if len(got) != 2 || got[0] != "en" || got[1] != "pt-br" {
					t.Fatalf("normalized languages = %#v", got)
				}
			},
		},
		{
			name: "Should match language prefixes",
			assert: func(t *testing.T) {
				if !languageMatches("en-US", "en") {
					t.Fatal("expected language prefix match")
				}
			},
		},
		{
			name: "Should reject unrelated languages",
			assert: func(t *testing.T) {
				if languageMatches("de", "en") {
					t.Fatal("did not expect unrelated language match")
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.assert(t)
		})
	}
}

func TestIsStatusCodeNetworkBlocked(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		err                error
		wantNetworkBlocked bool
		wantRateLimited    bool
	}{
		{name: "Should classify forbidden status as network blocked", err: errors.New("unexpected status code: 403"), wantNetworkBlocked: true},
		{name: "Should classify too many requests as rate limited", err: errors.New("unexpected status code: 429"), wantRateLimited: true},
		{name: "Should not classify rate limits as generic network blocked", err: errors.New("unexpected status code: 429"), wantRateLimited: true},
		{name: "Should ignore unrelated status codes", err: errors.New("unexpected status code: 404")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := isStatusCodeNetworkBlocked(tc.err); got != tc.wantNetworkBlocked {
				t.Fatalf("network blocked = %v, want %v", got, tc.wantNetworkBlocked)
			}
			if got := isStatusCodeRateLimited(tc.err); got != tc.wantRateLimited {
				t.Fatalf("rate limited = %v, want %v", got, tc.wantRateLimited)
			}
		})
	}
}

type stubSTTClient struct {
	called      bool
	audio       []byte
	format      string
	transcript  string
	transcripts []string
	err         error
}

func (client *stubSTTClient) Provider() string {
	return "stub"
}

func (client *stubSTTClient) Model() string {
	return "stub-model"
}

func (client *stubSTTClient) Transcribe(_ context.Context, audio []byte, format string) (string, error) {
	client.called = true
	client.audio = append([]byte(nil), audio...)
	client.format = format
	if client.err != nil {
		return "", client.err
	}
	if len(client.transcripts) > 0 {
		transcript := client.transcripts[0]
		client.transcripts = client.transcripts[1:]
		return transcript, nil
	}
	return client.transcript, nil
}
