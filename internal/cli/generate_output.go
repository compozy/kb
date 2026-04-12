package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"

	kgenerate "github.com/user/kb/internal/generate"
)

type generateProgressMode string

const (
	generateProgressAuto   generateProgressMode = "auto"
	generateProgressAlways generateProgressMode = "always"
	generateProgressNever  generateProgressMode = "never"
)

type generateLogFormat string

const (
	generateLogFormatJSON generateLogFormat = "json"
	generateLogFormatText generateLogFormat = "text"
)

func parseGenerateProgressMode(value string) (generateProgressMode, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", string(generateProgressAuto):
		return generateProgressAuto, nil
	case string(generateProgressAlways):
		return generateProgressAlways, nil
	case string(generateProgressNever):
		return generateProgressNever, nil
	default:
		return "", fmt.Errorf("unsupported progress mode %q (want auto, always, or never)", value)
	}
}

func parseGenerateLogFormat(value string) (generateLogFormat, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", string(generateLogFormatText):
		return generateLogFormatText, nil
	case string(generateLogFormatJSON):
		return generateLogFormatJSON, nil
	default:
		return "", fmt.Errorf("unsupported log format %q (want text or json)", value)
	}
}

func newGenerateObserver(
	writer io.Writer,
	progressMode generateProgressMode,
	logFormat generateLogFormat,
) kgenerate.Observer {
	if logFormat == generateLogFormatJSON {
		return &generateJSONObserver{encoder: json.NewEncoder(writer)}
	}

	liveProgress := false
	switch progressMode {
	case generateProgressAlways:
		liveProgress = true
	case generateProgressAuto:
		liveProgress = isTerminalWriter(writer)
	}

	return &generateTextObserver{
		liveProgress: liveProgress,
		writer:       writer,
	}
}

func isTerminalWriter(writer io.Writer) bool {
	fdWriter, ok := writer.(interface{ Fd() uintptr })
	if !ok {
		return false
	}

	return term.IsTerminal(int(fdWriter.Fd()))
}

type generateJSONObserver struct {
	encoder *json.Encoder
	mu      sync.Mutex
}

func (o *generateJSONObserver) ObserveGenerateEvent(_ context.Context, event kgenerate.Event) {
	if o == nil || o.encoder == nil {
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	_ = o.encoder.Encode(event)
}

type generateTextObserver struct {
	activeBar    *progressbar.ProgressBar
	activeStage  string
	liveProgress bool
	mu           sync.Mutex
	writer       io.Writer
}

func (o *generateTextObserver) ObserveGenerateEvent(_ context.Context, event kgenerate.Event) {
	if o == nil {
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	switch event.Kind {
	case kgenerate.EventStageStarted:
		o.handleStageStarted(event)
	case kgenerate.EventStageProgress:
		o.handleStageProgress(event)
	case kgenerate.EventStageCompleted:
		o.handleStageCompleted(event)
	case kgenerate.EventStageFailed:
		o.handleStageFailed(event)
	}
}

func (o *generateTextObserver) handleStageStarted(event kgenerate.Event) {
	label := humanStageLabel(event.Stage)
	if o.liveProgress && isDeterminateStage(event) {
		o.activeStage = event.Stage
		o.activeBar = progressbar.NewOptions(
			event.Total,
			progressbar.OptionClearOnFinish(),
			progressbar.OptionSetDescription("generate: "+label),
			progressbar.OptionSetElapsedTime(false),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionSetWidth(24),
			progressbar.OptionSetWriter(o.writer),
			progressbar.OptionShowCount(),
			progressbar.OptionThrottle(65*time.Millisecond),
		)
		_ = o.activeBar.Set64(0)
		return
	}

	_, _ = fmt.Fprintf(o.writer, "generate: %s started\n", label)
}

func (o *generateTextObserver) handleStageProgress(event kgenerate.Event) {
	if !o.liveProgress || o.activeBar == nil || o.activeStage != event.Stage {
		return
	}
	_ = o.activeBar.Set64(int64(event.Completed))
}

func (o *generateTextObserver) handleStageCompleted(event kgenerate.Event) {
	if o.activeBar != nil && o.activeStage == event.Stage {
		_ = o.activeBar.Finish()
		o.activeBar = nil
		o.activeStage = ""
	}

	_, _ = fmt.Fprintf(
		o.writer,
		"generate: %s completed in %dms%s\n",
		humanStageLabel(event.Stage),
		event.DurationMillis,
		completedSuffix(event),
	)
}

func (o *generateTextObserver) handleStageFailed(event kgenerate.Event) {
	if o.activeBar != nil && o.activeStage == event.Stage {
		_, _ = fmt.Fprintln(o.writer)
		o.activeBar = nil
		o.activeStage = ""
	}

	_, _ = fmt.Fprintf(
		o.writer,
		"generate: %s failed after %dms: %s\n",
		humanStageLabel(event.Stage),
		event.DurationMillis,
		event.Error,
	)
}

func isDeterminateStage(event kgenerate.Event) bool {
	switch event.Stage {
	case "parse", "write":
		return event.Total > 0
	default:
		return false
	}
}

func humanStageLabel(stage string) string {
	return strings.ReplaceAll(stage, "_", " ")
}

func completedSuffix(event kgenerate.Event) string {
	if event.Total <= 0 {
		return ""
	}

	return fmt.Sprintf(" (%d/%d)", event.Completed, event.Total)
}
