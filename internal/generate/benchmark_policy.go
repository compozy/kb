package generate

import (
	"errors"
	"sort"
	"time"

	"github.com/compozy/kb/internal/models"
)

const (
	javaIngestOverheadBudget = 1.20
	javaBenchmarkRepeatCount = 3
)

type javaBenchmarkProfile string

const (
	javaBenchmarkProfileSingleModuleLibrary   javaBenchmarkProfile = "single-module-library"
	javaBenchmarkProfileSpringService         javaBenchmarkProfile = "spring-service"
	javaBenchmarkProfileMultiModuleEnterprise javaBenchmarkProfile = "multi-module-enterprise"
)

var errEmptyBenchmarkSamples = errors.New("benchmark samples cannot be empty")

type javaBenchmarkFixture struct {
	Profile javaBenchmarkProfile
	Label   string
}

type javaBenchmarkPolicy struct {
	RepeatCount    int
	OverheadBudget float64
}

func canonicalJavaBenchmarkPolicy() javaBenchmarkPolicy {
	return javaBenchmarkPolicy{
		RepeatCount:    javaBenchmarkRepeatCount,
		OverheadBudget: javaIngestOverheadBudget,
	}
}

func canonicalJavaBenchmarkFixtures() []javaBenchmarkFixture {
	return []javaBenchmarkFixture{
		{
			Profile: javaBenchmarkProfileSingleModuleLibrary,
			Label:   "single module library",
		},
		{
			Profile: javaBenchmarkProfileSpringService,
			Label:   "spring-style service",
		},
		{
			Profile: javaBenchmarkProfileMultiModuleEnterprise,
			Label:   "multi-module enterprise",
		},
	}
}

func benchmarkGenerateOptions(rootPath string) models.GenerateOptions {
	return models.GenerateOptions{
		RootPath: rootPath,
		DryRun:   true,
		Semantic: false,
	}
}

func medianDurationFromSamples(samples []time.Duration) (time.Duration, error) {
	if len(samples) == 0 {
		return 0, errEmptyBenchmarkSamples
	}

	ordered := append([]time.Duration(nil), samples...)
	sort.Slice(ordered, func(left int, right int) bool {
		return ordered[left] < ordered[right]
	})

	middle := len(ordered) / 2
	if len(ordered)%2 == 1 {
		return ordered[middle], nil
	}

	return (ordered[middle-1] + ordered[middle]) / 2, nil
}
