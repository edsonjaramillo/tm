package cli

import (
	"github.com/edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

func newOpencodeCmd(useCases app.UseCases) *cobra.Command {
	var newWindow bool

	cmd := &cobra.Command{
		Use:   "opencode",
		Short: "start a opencode window",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCases.Opencode(app.WindowRequest{New: newWindow})
		},
	}

	cmd.Flags().BoolVar(&newWindow, "new", false, "create a new tmux window if none exists")
	return cmd
}
