package cli

import (
	"github.com/edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

func newQuadsCmd(useCases app.UseCases) *cobra.Command {
	var newWindow bool

	cmd := &cobra.Command{
		Use:   "quads",
		Short: "start a 4x4 pane setup",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCases.Quads(app.WindowRequest{New: newWindow})
		},
	}

	cmd.Flags().BoolVar(&newWindow, "new", false, "create a new tmux window if none exists")
	return cmd
}
