package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/edihasaj/shipyard/internal/config"
	"github.com/edihasaj/shipyard/internal/assets"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init <repo>",
		Short: "Scaffold a per-repo config from the schema template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			dir := config.ReposDir()
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
			// Drop the documented schema next to configs the first time.
			schema := filepath.Join(dir, "_schema.yml")
			if _, err := os.Stat(schema); os.IsNotExist(err) {
				_ = os.WriteFile(schema, assets.SchemaYML, 0o644)
			}

			dest := filepath.Join(dir, key+".yml")
			if _, err := os.Stat(dest); err == nil && !force {
				return fmt.Errorf("config already exists: %s (use --force to overwrite)", dest)
			}
			if err := os.WriteFile(dest, assets.SchemaYML, 0o644); err != nil {
				return err
			}
			fmt.Printf("created %s\n", dest)
			fmt.Printf("edit it (set path, task_source, gates, push), then: shipyard %s \"<task>\"\n", key)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing config")
	return cmd
}
