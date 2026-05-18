// Package migrate contains explicit vault layout migrations.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/topic"
)

// TranscriptResult describes a raw/transcripts to raw/youtube migration.
type TranscriptResult struct {
	Topic      string   `json:"topic"`
	Moved      int      `json:"moved"`
	Files      []string `json:"files"`
	LegacyPath string   `json:"legacyPath"`
	TargetPath string   `json:"targetPath"`
}

// Transcripts moves legacy raw/transcripts markdown files into raw/youtube.
func Transcripts(vaultPath, topicRef string, when time.Time) (TranscriptResult, error) {
	topicInfo, err := topic.Resolve(vaultPath, topicRef)
	if err != nil {
		return TranscriptResult{}, fmt.Errorf("migrate transcripts: resolve topic: %w", err)
	}
	if err := topic.EnsureCurrentSkeleton(topicInfo.RootPath); err != nil {
		return TranscriptResult{}, fmt.Errorf("migrate transcripts: ensure topic skeleton: %w", err)
	}

	legacyDir := filepath.Join(topicInfo.RootPath, "raw", "transcripts")
	targetDir := filepath.Join(topicInfo.RootPath, "raw", "youtube")
	result := TranscriptResult{
		Topic:      topicInfo.Slug,
		LegacyPath: legacyDir,
		TargetPath: targetDir,
	}

	entries, err := os.ReadDir(legacyDir)
	if errors.Is(err, os.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return TranscriptResult{}, fmt.Errorf("migrate transcripts: read %q: %w", legacyDir, err)
	}

	moves := make([]fileMove, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		sourcePath := filepath.Join(legacyDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())
		if _, err := os.Stat(targetPath); err == nil {
			return TranscriptResult{}, fmt.Errorf("migrate transcripts: target already exists: %s", targetPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return TranscriptResult{}, fmt.Errorf("migrate transcripts: stat target %q: %w", targetPath, err)
		}
		moves = append(moves, fileMove{source: sourcePath, target: targetPath, name: entry.Name()})
	}
	if len(moves) == 0 {
		return result, nil
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return TranscriptResult{}, fmt.Errorf("migrate transcripts: create %q: %w", targetDir, err)
	}
	for _, move := range moves {
		if err := os.Rename(move.source, move.target); err != nil {
			return TranscriptResult{}, fmt.Errorf("migrate transcripts: move %q to %q: %w", move.source, move.target, err)
		}
		result.Files = append(result.Files, filepath.ToSlash(filepath.Join(topicInfo.Slug, "raw", "youtube", move.name)))
	}
	result.Moved = len(result.Files)

	if err := appendLogEntry(filepath.Join(topicInfo.RootPath, "log.md"), when, result.Files); err != nil {
		return TranscriptResult{}, err
	}
	return result, nil
}

type fileMove struct {
	source string
	target string
	name   string
}

func appendLogEntry(logPath string, when time.Time, files []string) error {
	if when.IsZero() {
		when = time.Now().UTC()
	}
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("migrate transcripts: open log %q: %w", logPath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	entry := strings.Join([]string{
		"",
		fmt.Sprintf("## [%s] migrate | transcripts", when.UTC().Format(frontmatter.DateLayout)),
		"",
		fmt.Sprintf("Moved %d legacy transcript files from `raw/transcripts` to `raw/youtube`.", len(files)),
		"",
	}, "\n")
	if _, err := file.WriteString(entry); err != nil {
		return fmt.Errorf("migrate transcripts: write log entry: %w", err)
	}
	return nil
}
