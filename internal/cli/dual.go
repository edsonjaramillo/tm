package cli

import (
	"edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

func newDualCmd(useCases app.UseCases) *cobra.Command {
	var newWindow bool

	cmd := &cobra.Command{
		Use:   "dual",
		Short: "start a 2-pane setup",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCases.Dual(app.WindowRequest{New: newWindow})
		},
	}

	cmd.Flags().BoolVar(&newWindow, "new", false, "create a new tmux window if none exists")
	return cmd
}
