package youtube

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/config"
)

type audioChunk struct {
	Path   string
	Format string
	Offset time.Duration
}

type transcribedChunk struct {
	Text   string
	Offset time.Duration
}

// NewTranscriber constructs the configured STT backend.
func NewTranscriber(sttConfig config.STTConfig, openRouterConfig config.OpenRouterConfig) sttClient {
	sttConfig = normalizeSTTConfig(sttConfig)
	switch strings.ToLower(strings.TrimSpace(sttConfig.Provider)) {
	case "openrouter":
		return NewOpenRouterClientWithPrompt(openRouterConfigForSTT(sttConfig, openRouterConfig), openRouterPrompt(sttConfig))
	default:
		return NewOpenAITranscriber(sttConfig)
	}
}

func openRouterPrompt(sttConfig config.STTConfig) string {
	parts := make([]string, 0, 2)
	if prompt := strings.TrimSpace(sttConfig.Prompt); prompt != "" {
		parts = append(parts, prompt)
	}
	if language := strings.TrimSpace(sttConfig.Language); language != "" && !strings.EqualFold(language, "auto") {
		parts = append(parts, "The spoken language is "+language+".")
	}
	return strings.Join(parts, "\n")
}

func normalizeSTTConfig(sttConfig config.STTConfig) config.STTConfig {
	defaults := config.Default().STT
	if strings.TrimSpace(sttConfig.Provider) == "" {
		sttConfig.Provider = defaults.Provider
	}
	if strings.TrimSpace(sttConfig.APIURL) == "" {
		sttConfig.APIURL = defaults.APIURL
	}
	if strings.TrimSpace(sttConfig.Model) == "" {
		sttConfig.Model = defaults.Model
	}
	if strings.TrimSpace(sttConfig.Language) == "" {
		sttConfig.Language = defaults.Language
	}
	if strings.TrimSpace(sttConfig.AudioFormat) == "" {
		sttConfig.AudioFormat = defaults.AudioFormat
	}
	if strings.TrimSpace(sttConfig.ChunkDuration) == "" {
		sttConfig.ChunkDuration = defaults.ChunkDuration
	}
	if sttConfig.MaxChunkBytes == 0 {
		sttConfig.MaxChunkBytes = defaults.MaxChunkBytes
	}
	if sttConfig.Concurrency == 0 {
		sttConfig.Concurrency = defaults.Concurrency
	}
	if strings.TrimSpace(sttConfig.FFmpegPath) == "" {
		sttConfig.FFmpegPath = defaults.FFmpegPath
	}
	return sttConfig
}

func openRouterConfigForSTT(sttConfig config.STTConfig, openRouterConfig config.OpenRouterConfig) config.OpenRouterConfig {
	defaultSTT := config.Default().STT
	defaultOpenRouter := config.Default().OpenRouter
	apiURL := strings.TrimSpace(openRouterConfig.APIURL)
	if apiURL == "" {
		apiURL = defaultOpenRouter.APIURL
	}
	if candidate := strings.TrimSpace(sttConfig.APIURL); candidate != "" && candidate != defaultSTT.APIURL {
		apiURL = candidate
	}
	model := strings.TrimSpace(openRouterConfig.STTModel)
	if model == "" {
		model = defaultOpenRouter.STTModel
	}
	if candidate := strings.TrimSpace(sttConfig.Model); candidate != "" && candidate != defaultSTT.Model {
		model = candidate
	}
	apiKey := strings.TrimSpace(openRouterConfig.APIKey)
	if candidate := strings.TrimSpace(sttConfig.APIKey); candidate != "" {
		apiKey = candidate
	}
	return config.OpenRouterConfig{
		APIKey:   apiKey,
		APIURL:   apiURL,
		STTModel: model,
	}
}

func (extractor *Extractor) extractSTTFromYTDLP(
	ctx context.Context,
	parsed parsedVideoURL,
	info ytDLPInfo,
	result *Result,
	transcriptErr error,
) (*Result, error) {
	if result == nil {
		result = &Result{Metadata: metadataFromYTDLPInfo(parsed, info)}
	}
	if !extractor.shouldAttemptSTT(result.TranscriptionPolicy) {
		return result, transcriptErr
	}
	if extractor.ytDLP == nil {
		return result, errors.Join(transcriptErr, errors.New("youtube stt: yt-dlp backend is nil"))
	}
	if extractor.stt == nil {
		return result, errors.Join(transcriptErr, errors.New("youtube stt: transcriber is nil"))
	}

	sttConfig := normalizeSTTConfig(extractor.sttConfig)
	audio, audioErr := extractor.ytDLP.downloadAudio(ctx, parsed.CanonicalURL, sttConfig.AudioFormat)
	if audioErr != nil {
		return result, errors.Join(transcriptErr, fmt.Errorf("youtube stt: %w", audioErr))
	}
	defer audio.Cleanup()

	markdown, sttErr := extractor.transcribeAudioPath(ctx, audio.Path, audio.Format, result.Metadata.Duration)
	if sttErr != nil {
		return result, errors.Join(transcriptErr, fmt.Errorf("youtube stt: %w", sttErr))
	}
	result.Markdown = markdown
	result.Source = TranscriptSourceSTT
	result.Language = ""
	result.CaptionKind = ""
	result.STTProvider = extractor.stt.Provider()
	result.STTModel = extractor.stt.Model()
	return result, nil
}

func (extractor *Extractor) transcribeAudioPath(ctx context.Context, audioPath string, format string, duration time.Duration) (string, error) {
	chunks, cleanup, err := extractor.prepareAudioChunks(ctx, audioPath, format, duration)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		return "", err
	}
	return extractor.transcribeChunks(ctx, chunks)
}

func (extractor *Extractor) prepareAudioChunks(
	ctx context.Context,
	audioPath string,
	format string,
	duration time.Duration,
) ([]audioChunk, func(), error) {
	sttConfig := normalizeSTTConfig(extractor.sttConfig)
	chunkDuration, err := sttConfig.ChunkDurationValue()
	if err != nil {
		return nil, nil, err
	}
	info, err := os.Stat(audioPath)
	if err != nil {
		return nil, nil, fmt.Errorf("stat audio file %q: %w", audioPath, err)
	}
	if info.Size() == 0 {
		return nil, nil, errors.New("audio file is empty")
	}
	if info.Size() <= sttConfig.MaxChunkBytes && (duration <= 0 || duration <= chunkDuration) {
		return []audioChunk{{Path: audioPath, Format: format, Offset: 0}}, nil, nil
	}
	chunks, cleanup, err := extractor.segmentAudio(ctx, audioPath, format, chunkDuration, sttConfig.MaxChunkBytes)
	if err != nil {
		return nil, cleanup, err
	}
	return chunks, cleanup, nil
}

func (extractor *Extractor) segmentAudio(
	ctx context.Context,
	audioPath string,
	format string,
	chunkDuration time.Duration,
	maxChunkBytes int64,
) ([]audioChunk, func(), error) {
	sttConfig := normalizeSTTConfig(extractor.sttConfig)
	ffmpegPath, err := exec.LookPath(strings.TrimSpace(sttConfig.FFmpegPath))
	if err != nil {
		return nil, nil, fmt.Errorf("ffmpeg not found; install ffmpeg or configure stt.ffmpeg_path: %w", err)
	}
	tempDir, err := os.MkdirTemp("", "kb-stt-chunks-*")
	if err != nil {
		return nil, nil, fmt.Errorf("create audio chunk directory: %w", err)
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		format = strings.TrimSpace(strings.ToLower(sttConfig.AudioFormat))
	}
	pattern := filepath.Join(tempDir, "chunk-%03d."+format)
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-y",
		"-i", audioPath,
		"-map", "0:a:0",
		"-f", "segment",
		"-segment_time", strconv.FormatInt(int64(chunkDuration.Seconds()), 10),
		"-reset_timestamps", "1",
		pattern,
	}
	command := exec.CommandContext(ctx, ffmpegPath, args...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		cleanup()
		diagnostics := strings.TrimSpace(stderr.String())
		if diagnostics != "" {
			return nil, nil, fmt.Errorf("ffmpeg segment audio: %w: %s", err, diagnostics)
		}
		return nil, nil, fmt.Errorf("ffmpeg segment audio: %w", err)
	}

	chunks, err := collectAudioChunks(tempDir, format, chunkDuration, maxChunkBytes)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return chunks, cleanup, nil
}

func collectAudioChunks(directory string, format string, chunkDuration time.Duration, maxChunkBytes int64) ([]audioChunk, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read audio chunks: %w", err)
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(strings.TrimPrefix(filepath.Ext(entry.Name()), "."), format) {
			paths = append(paths, filepath.Join(directory, entry.Name()))
		}
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		return nil, errors.New("ffmpeg did not produce audio chunks")
	}
	chunks := make([]audioChunk, 0, len(paths))
	for index, chunkPath := range paths {
		info, err := os.Stat(chunkPath)
		if err != nil {
			return nil, fmt.Errorf("stat audio chunk %q: %w", chunkPath, err)
		}
		if info.Size() == 0 {
			return nil, fmt.Errorf("audio chunk %q is empty", chunkPath)
		}
		if info.Size() > maxChunkBytes {
			return nil, fmt.Errorf("audio chunk %q is %d bytes, exceeds stt.max_chunk_bytes %d", chunkPath, info.Size(), maxChunkBytes)
		}
		chunks = append(chunks, audioChunk{
			Path:   chunkPath,
			Format: format,
			Offset: time.Duration(index) * chunkDuration,
		})
	}
	return chunks, nil
}

func (extractor *Extractor) transcribeChunks(ctx context.Context, chunks []audioChunk) (string, error) {
	if len(chunks) == 0 {
		return "", errors.New("no audio chunks to transcribe")
	}
	if extractor.stt == nil {
		return "", errors.New("transcriber is nil")
	}
	sttConfig := normalizeSTTConfig(extractor.sttConfig)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]transcribedChunk, len(chunks))
	sem := make(chan struct{}, sttConfig.Concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	setErr := func(err error) {
		if err == nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
			cancel()
		}
	}

	for index, chunk := range chunks {
		index, chunk := index, chunk
		select {
		case <-ctx.Done():
			setErr(ctx.Err())
			continue
		case sem <- struct{}{}:
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			data, err := os.ReadFile(chunk.Path)
			if err != nil {
				setErr(fmt.Errorf("read audio chunk %q: %w", chunk.Path, err))
				return
			}
			text, err := extractor.stt.Transcribe(ctx, data, chunk.Format)
			if err != nil {
				setErr(err)
				return
			}
			text = normalizeTranscriptText(text)
			if text == "" {
				setErr(errors.New("empty transcription response"))
				return
			}
			results[index] = transcribedChunk{Text: text, Offset: chunk.Offset}
		}()
	}
	wg.Wait()
	if firstErr != nil {
		return "", firstErr
	}
	return formatSTTChunksMarkdown(results), nil
}

func formatSTTChunksMarkdown(chunks []transcribedChunk) string {
	var builder strings.Builder
	wrote := false
	for _, chunk := range chunks {
		text := normalizeTranscriptText(chunk.Text)
		if text == "" {
			continue
		}
		if wrote {
			builder.WriteString("\n\n")
		}
		builder.WriteString("## ")
		builder.WriteString(formatTimestamp(int(chunk.Offset / time.Millisecond)))
		builder.WriteString("\n")
		builder.WriteString(text)
		wrote = true
	}
	return builder.String()
}
