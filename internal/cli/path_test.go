package cli

import (
	"path/filepath"
	"testing"
)

func absoluteTestPath(t testing.TB, path string) string {
	t.Helper()
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("absolute path %q: %v", path, err)
	}
	return absolute
}
