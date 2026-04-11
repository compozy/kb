package convert

import (
	"os"
	"path/filepath"
	"testing"
)

func readConvertFixtureBytes(t *testing.T, fixtureName string) []byte {
	t.Helper()

	path := filepath.Join("testdata", fixtureName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}

	return data
}

func warningList(t *testing.T, metadata map[string]any) []string {
	t.Helper()

	warnings, ok := metadata["warnings"].([]string)
	if !ok {
		t.Fatalf("warnings = %#v, want []string", metadata["warnings"])
	}

	return warnings
}
