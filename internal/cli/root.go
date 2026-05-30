// Package cli wires the shipyard command tree.
package cli

import (
	"github.com/spf13/cobra"
)

// Execute runs the root command.
func Execute(version string) error {
	root := &cobra.Command{
		Use:   "shipyard <repo> <task> [notes...]",
		Short: "Point an agent at a repo + task and run an end-to-end ship pipeline",
		Long: `shipyard runs a per-task pipeline against a managed repo:
resolve the task (Jira key, GitHub issue, URL, or free text) -> branch in the
repo's convention -> implement -> run gates (lint/typecheck/test) -> security +
code review -> write a PR description -> optionally smoke-test -> stop PR-ready
or open a PR, per the repo's config.

The pipeline itself is a skill (run 'shipyard install-skill'); this binary is
the launcher + config layer. Per-repo configs live under the config home
(see 'shipyard where').`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		RunE:          runShip, // default: treat args as <repo> <task...>
	}
	bindShipFlags(root)

	root.AddCommand(
		newListCmd(),
		newInitCmd(),
		newInstallSkillCmd(),
		newWhereCmd(),
	)
	return root.Execute()
}
