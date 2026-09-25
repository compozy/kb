package repohealth

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindPackageBoundaryViolations(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		path    string
		source  string
		want    bool
		wantErr bool
	}{
		{name: "previously unchecked package", path: "internal/ingest/ingest.go", source: `package ingest; import _ "github.com/compozy/kb/internal/cli"`, want: true},
		{name: "nested CLI package", path: "internal/nested/worker/worker.go", source: `package worker; import "github.com/compozy/kb/internal/cli/commands"`, want: true},
		{name: "raw import string", path: "internal/topic/topic.go", source: "package topic; import `github.com/compozy/kb/internal/cli`", want: true},
		{name: "test import", path: "internal/topic/topic_test.go", source: `package topic_test; import "github.com/compozy/kb/internal/cli"`, want: true},
		{name: "other platform", path: "internal/topic/topic_windows.go", source: `package topic; import "github.com/compozy/kb/internal/cli"`, want: true},
		{name: "CLI can use CLI subpackages", path: "internal/cli/root.go", source: `package cli; import "github.com/compozy/kb/internal/cli/commands"`},
		{name: "comment is not an import", path: "internal/topic/topic.go", source: "package topic\n// github.com/compozy/kb/internal/cli\n"},
		{name: "similar package name", path: "internal/topic/topic.go", source: `package topic; import "github.com/compozy/kb/internal/client"`},
		{name: "fixture is not compiled", path: "internal/topic/testdata/example.go", source: `package example; import "github.com/compozy/kb/internal/cli"`},
		{name: "vendored source", path: "internal/topic/vendor/example.go", source: `package example; import "github.com/compozy/kb/internal/cli"`},
		{name: "invalid imports fail closed", path: "internal/topic/topic.go", source: `package topic; import "unterminated`, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(tc.path))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(tc.source), 0o644); err != nil {
				t.Fatal(err)
			}
			violations, err := FindPackageBoundaryViolations(root)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, want error %v", err, tc.wantErr)
			}
			if !tc.want {
				if len(violations) != 0 {
					t.Fatalf("unexpected violations: %v", violations)
				}
				return
			}
			if len(violations) != 1 || violations[0].Path != tc.path || !strings.Contains(violations[0].Reason, "imports github.com/compozy/kb/internal/cli") {
				t.Fatalf("violations = %v, want forbidden import in %s", violations, tc.path)
			}
		})
	}
}

func TestFindPackageBoundaryViolationsRequiresReadableSourceRoot(t *testing.T) {
	t.Parallel()
	if _, err := FindPackageBoundaryViolations(t.TempDir()); err == nil {
		t.Fatal("missing source directory must not pass the boundary check")
	}
}
