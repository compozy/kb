//go:build integration

package okf

import (
	"context"
	"path/filepath"
	"testing"
)

func TestOfficialBundlesPassLenientConformance(t *testing.T) {
	t.Parallel()

	for _, bundle := range []string{"ga4", "stackoverflow", "crypto_bitcoin"} {
		bundle := bundle
		t.Run(bundle, func(t *testing.T) {
			t.Parallel()
			issues, err := Check(context.Background(), filepath.Join("testdata", "official", bundle), CheckOptions{})
			if err != nil {
				t.Fatalf("Check returned error: %v", err)
			}
			if HasErrors(issues) {
				t.Fatalf("official bundle %s has error issues: %#v", bundle, issues)
			}
		})
	}
}
