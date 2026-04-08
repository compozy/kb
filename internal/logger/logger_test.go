package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		level     string
		wantErr   bool
		wantLevel slog.Level
	}{
		{name: "debug level", level: "debug", wantLevel: slog.LevelDebug},
		{name: "info level", level: "info", wantLevel: slog.LevelInfo},
		{name: "warn level", level: "warn", wantLevel: slog.LevelWarn},
		{name: "error level", level: "error", wantLevel: slog.LevelError},
		{name: "empty defaults to info", level: "", wantLevel: slog.LevelInfo},
		{name: "case insensitive", level: "DEBUG", wantLevel: slog.LevelDebug},
		{name: "whitespace trimmed", level: "  info  ", wantLevel: slog.LevelInfo},
		{name: "invalid level", level: "trace", wantErr: true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger, err := New(tc.level, WithWriter(&bytes.Buffer{}))
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if logger == nil {
				t.Fatal("expected non-nil logger")
			}
		})
	}
}

func TestNewWritesJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger, err := New("info", WithWriter(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logger.Info("test message", "key", "value")

	var record map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("output is not valid JSON: %v\nraw: %s", err, buf.String())
	}
	if record["msg"] != "test message" {
		t.Errorf("expected msg 'test message', got %q", record["msg"])
	}
	if record["key"] != "value" {
		t.Errorf("expected key 'value', got %q", record["key"])
	}
}

func TestNewUppercaseLevelKey(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger, err := New("info", WithWriter(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logger.Info("test")

	var record map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	level, ok := record["level"].(string)
	if !ok {
		t.Fatal("level key not found or not a string")
	}
	if level != "INFO" {
		t.Errorf("expected uppercase level 'INFO', got %q", level)
	}
}

func TestNewWithObserver(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	obs := &testObserver{}
	logger, err := New("info", WithWriter(&buf), WithObserver(obs))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logger.Info("observed message", "status", "ok")

	obs.mu.Lock()
	defer obs.mu.Unlock()
	if len(obs.records) != 1 {
		t.Fatalf("expected 1 observed record, got %d", len(obs.records))
	}
	if obs.records[0].Message != "observed message" {
		t.Errorf("expected message 'observed message', got %q", obs.records[0].Message)
	}
}

func TestNewLevelFiltering(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger, err := New("warn", WithWriter(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logger.Info("should be filtered")

	if buf.Len() != 0 {
		t.Errorf("expected no output for info at warn level, got: %s", buf.String())
	}
}

type testObserver struct {
	mu      sync.Mutex
	records []slog.Record
}

func (o *testObserver) Observe(_ context.Context, record slog.Record) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.records = append(o.records, record)
}
