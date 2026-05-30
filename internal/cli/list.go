package cli

import (
	"fmt"

	"github.com/edihasaj/shipyard/internal/config"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured repos",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			keys, err := config.List()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Printf("no repos configured in %s\n", config.ReposDir())
				fmt.Println("add one with: shipyard init <repo>")
				return nil
			}
			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		},
	}
}
