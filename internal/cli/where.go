package cli

import (
	"fmt"

	"github.com/edihasaj/shipyard/internal/config"
	"github.com/spf13/cobra"
)

func newWhereCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "where",
		Short: "Print the resolved config home and repos dir",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Println("config home:", config.Home())
			fmt.Println("repos dir:  ", config.ReposDir())
		},
	}
}
