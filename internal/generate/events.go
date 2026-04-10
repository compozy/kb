package generate

import "context"

// EventKind identifies one generation lifecycle event.
type EventKind string

const (
	EventStageStarted   EventKind = "stage_started"
	EventStageProgress  EventKind = "stage_progress"
	EventStageCompleted EventKind = "stage_completed"
	EventStageFailed    EventKind = "stage_failed"
)

// Event reports structured progress from the generate pipeline.
type Event struct {
	Kind           EventKind      `json:"kind"`
	Stage          string         `json:"stage"`
	Completed      int            `json:"completed,omitempty"`
	DurationMillis int64          `json:"durationMillis,omitempty"`
	Error          string         `json:"error,omitempty"`
	Fields         map[string]any `json:"fields,omitempty"`
	Total          int            `json:"total,omitempty"`
}

// Observer receives structured generation events.
type Observer interface {
	ObserveGenerateEvent(context.Context, Event)
}

// ObserverFunc adapts a function into an Observer.
type ObserverFunc func(context.Context, Event)

// ObserveGenerateEvent implements Observer.
func (fn ObserverFunc) ObserveGenerateEvent(ctx context.Context, event Event) {
	if fn != nil {
		fn(ctx, event)
	}
}

type noopObserver struct{}

func (noopObserver) ObserveGenerateEvent(context.Context, Event) {}
