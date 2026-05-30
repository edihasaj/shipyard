// Command shipyard points an agent at a repo + task and runs an end-to-end
// task pipeline: resolve task -> branch -> implement -> gate -> review ->
// PR-ready. The pipeline logic lives in an installable skill (SKILL.md); this
// binary is the launcher + config layer around it.
package main

import (
	"fmt"
	"os"

	"github.com/edihasaj/shipyard/internal/cli"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := cli.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, "shipyard:", err)
		os.Exit(1)
	}
}
