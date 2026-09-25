// Package jsonl provides append-only record writes shared by decision logs.
package jsonl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Append writes one compact JSON record, as returned by json.Marshal, and
// syncs it before returning. The record must not include its trailing newline.
// An unterminated previous record is separated, never truncated, so concurrent
// appenders cannot remove acknowledged rows. Readers may ignore torn fragments.
func Append(path string, record []byte) (err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close %s: %w", path, closeErr))
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	row := make([]byte, 0, len(record)+2)
	if info.Size() > 0 {
		var last [1]byte
		if _, err := file.ReadAt(last[:], info.Size()-1); err != nil {
			return fmt.Errorf("read log tail %s: %w", path, err)
		}
		if last[0] != '\n' {
			// Include the separator in the same append as the new row. A
			// concurrent appender can add only a harmless extra blank line.
			row = append(row, '\n')
		}
	}
	row = append(append(row, record...), '\n')
	if _, err := file.Write(row); err != nil {
		return fmt.Errorf("append %s: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync %s: %w", path, err)
	}
	return nil
}
