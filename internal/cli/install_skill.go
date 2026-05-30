package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/edihasaj/shipyard/internal/assets"
	"github.com/spf13/cobra"
)

func newInstallSkillCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "install-skill",
		Short: "Install the ship-task pipeline skill for the agent",
		Long: `Writes the bundled ship-task SKILL.md so the agent can run /ship-task.
Defaults to ~/.claude/skills/ship-task; override with --dir or $SHIPYARD_SKILLS_DIR.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if dir == "" {
				dir = os.Getenv("SHIPYARD_SKILLS_DIR")
			}
			if dir == "" {
				home, _ := os.UserHomeDir()
				dir = filepath.Join(home, ".claude", "skills")
			}
			target := filepath.Join(dir, "ship-task")
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			out := filepath.Join(target, "SKILL.md")
			if err := os.WriteFile(out, assets.SkillMD, 0o644); err != nil {
				return err
			}
			fmt.Printf("installed skill: %s\n", out)
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "", "skills directory (default: ~/.claude/skills)")
	return cmd
}
