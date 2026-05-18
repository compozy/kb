package repohealth

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindWindowsPathViolations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		paths      []string
		wantReason string
	}{
		{
			name:  "portable path",
			paths: []string{".codex/plans/2026-05-18-vault-topic-youtube.md"},
		},
		{
			name:       "asterisk",
			paths:      []string{".codex/plans/20260518-143218*vault-topic-youtube.md"},
			wantReason: "Windows-invalid character '*'",
		},
		{
			name:       "colon",
			paths:      []string{"docs/release:v1.md"},
			wantReason: "Windows-invalid character ':'",
		},
		{
			name:       "control character",
			paths:      []string{"docs/bad\x1fname.md"},
			wantReason: "ASCII control character",
		},
		{
			name:       "trailing period",
			paths:      []string{"docs/archive./note.md"},
			wantReason: "ends with a period",
		},
		{
			name:       "reserved device name with extension",
			paths:      []string{"docs/con.md"},
			wantReason: "Windows-reserved device name",
		},
		{
			name:       "empty component",
			paths:      []string{"docs//note.md"},
			wantReason: "empty component",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			violations := FindWindowsPathViolations(tt.paths)
			if tt.wantReason == "" {
				if len(violations) != 0 {
					t.Fatalf("FindWindowsPathViolations() = %#v, want no violations", violations)
				}
				return
			}

			if len(violations) != 1 {
				t.Fatalf("FindWindowsPathViolations() returned %d violations, want 1: %#v", len(violations), violations)
			}
			if !strings.Contains(violations[0].Reason, tt.wantReason) {
				t.Fatalf("violation reason = %q, want substring %q", violations[0].Reason, tt.wantReason)
			}
		})
	}
}

func TestTrackedRepositoryPathsAreWindowsPortable(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is unavailable")
	}

	repoRoot, err := gitOutput("rev-parse", "--show-toplevel")
	if err != nil {
		t.Skipf("not running inside a git worktree: %v", err)
	}

	pathsOutput, err := gitOutput("-C", repoRoot, "ls-files", "--cached", "--others", "--exclude-standard", "-z")
	if err != nil {
		t.Fatalf("list repository files: %v", err)
	}

	paths, err := existingRepositoryPaths(repoRoot, splitNullTerminated(pathsOutput))
	if err != nil {
		t.Fatal(err)
	}
	violations := FindWindowsPathViolations(paths)
	if len(violations) > 0 {
		t.Fatalf("repository paths must be portable to Windows checkouts:\n%s", FormatPathViolations(violations))
	}
}

func gitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func splitNullTerminated(output string) []string {
	if output == "" {
		return nil
	}

	parts := strings.Split(output, "\x00")
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func existingRepositoryPaths(repoRoot string, paths []string) ([]string, error) {
	existing := make([]string, 0, len(paths))
	for _, candidate := range paths {
		if candidate == "" {
			continue
		}
		_, err := os.Lstat(filepath.Join(repoRoot, filepath.FromSlash(candidate)))
		if err == nil {
			existing = append(existing, candidate)
			continue
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		return nil, err
	}
	return existing, nil
}
