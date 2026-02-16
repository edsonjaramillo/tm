package cli

import (
	"github.com/edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

func newDetachCmd(useCases app.UseCases) *cobra.Command {
	return &cobra.Command{
		Use:   "detach",
		Short: "detach a tmux session",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCases.Detach()
		},
	}
}
