//go:build integration

package adapter

import (
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestRustAdapterIntegrationResolvesWorkspaceCrateImports(t *testing.T) {
	t.Parallel()

	parsedFiles := parseRustSources(t, map[string]string{
		"Cargo.toml": `
[workspace]
members = ["crates/core", "crates/app"]
`,
		"crates/core/Cargo.toml": `
[package]
name = "openfang-core"
version = "0.1.0"
edition = "2021"
`,
		"crates/core/src/lib.rs": `
pub mod util;
`,
		"crates/core/src/util.rs": `
pub fn helper() {}
`,
		"crates/app/Cargo.toml": `
[package]
name = "openfang-app"
version = "0.1.0"
edition = "2021"
`,
		"crates/app/src/lib.rs": `
use openfang_core::util::helper;

pub fn run() {
    helper();
}
`,
	})

	appFile := mustFindParsedFile(t, parsedFiles, "crates/app/src/lib.rs")
	utilFile := mustFindParsedFile(t, parsedFiles, "crates/core/src/util.rs")
	helper := mustFindSymbol(t, utilFile.Symbols, "helper")
	run := mustFindSymbol(t, appFile.Symbols, "run")

	if !hasRelation(appFile.Relations, appFile.File.ID, utilFile.File.ID, models.RelImports) {
		t.Fatal("expected workspace crate import to resolve to util.rs")
	}
	if !hasRelation(appFile.Relations, appFile.File.ID, helper.ID, models.RelReferences) {
		t.Fatal("expected workspace crate import to reference helper()")
	}
	if !hasRelation(appFile.Relations, run.ID, helper.ID, models.RelCalls) {
		t.Fatal("expected run() to call helper() across workspace crates")
	}
}
