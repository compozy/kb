package generate

import (
	"errors"
	"testing"
	"time"
)

func TestMedianDurationFromSamples(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		samples []time.Duration
		want    time.Duration
		wantErr error
	}{
		{
			name: "odd sample count returns middle",
			samples: []time.Duration{
				4 * time.Millisecond,
				1 * time.Millisecond,
				3 * time.Millisecond,
			},
			want: 3 * time.Millisecond,
		},
		{
			name: "even sample count returns midpoint average",
			samples: []time.Duration{
				8 * time.Millisecond,
				2 * time.Millisecond,
				4 * time.Millisecond,
				10 * time.Millisecond,
			},
			want: 6 * time.Millisecond,
		},
		{
			name:    "empty samples fail",
			samples: nil,
			wantErr: errEmptyBenchmarkSamples,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := medianDurationFromSamples(testCase.samples)
			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("medianDurationFromSamples() error = %v, want %v", err, testCase.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("medianDurationFromSamples() error = %v, want nil", err)
			}
			if got != testCase.want {
				t.Fatalf("medianDurationFromSamples() = %s, want %s", got, testCase.want)
			}
		})
	}
}

func TestCanonicalJavaBenchmarkFixtures(t *testing.T) {
	t.Parallel()

	fixtures := canonicalJavaBenchmarkFixtures()
	if len(fixtures) != 3 {
		t.Fatalf("canonicalJavaBenchmarkFixtures() len = %d, want 3", len(fixtures))
	}

	expectedProfiles := []javaBenchmarkProfile{
		javaBenchmarkProfileSingleModuleLibrary,
		javaBenchmarkProfileSpringService,
		javaBenchmarkProfileMultiModuleEnterprise,
	}
	for idx, expectedProfile := range expectedProfiles {
		if fixtures[idx].Profile != expectedProfile {
			t.Fatalf("canonicalJavaBenchmarkFixtures()[%d].Profile = %q, want %q", idx, fixtures[idx].Profile, expectedProfile)
		}
		if fixtures[idx].Label == "" {
			t.Fatalf("canonicalJavaBenchmarkFixtures()[%d].Label is empty", idx)
		}
	}
}

func TestBenchmarkGenerateOptions(t *testing.T) {
	t.Parallel()

	options := benchmarkGenerateOptions("/tmp/canonical-repo")
	if options.RootPath != "/tmp/canonical-repo" {
		t.Fatalf("RootPath = %q, want /tmp/canonical-repo", options.RootPath)
	}
	if !options.DryRun {
		t.Fatalf("DryRun = %t, want true", options.DryRun)
	}
	if options.Semantic {
		t.Fatalf("Semantic = %t, want false", options.Semantic)
	}
}

func TestCanonicalJavaBenchmarkPolicy(t *testing.T) {
	t.Parallel()

	policy := canonicalJavaBenchmarkPolicy()
	if policy.RepeatCount != 3 {
		t.Fatalf("RepeatCount = %d, want 3", policy.RepeatCount)
	}
	if policy.OverheadBudget != 1.20 {
		t.Fatalf("OverheadBudget = %.2f, want 1.20", policy.OverheadBudget)
	}
}
