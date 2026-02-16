package cli

import (
	"github.com/edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

func newCodexCmd(useCases app.UseCases) *cobra.Command {
	var newWindow bool

	cmd := &cobra.Command{
		Use:   "codex",
		Short: "start a codex window",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCases.Codex(app.WindowRequest{New: newWindow})
		},
	}

	cmd.Flags().BoolVar(&newWindow, "new", false, "create a new tmux window if none exists")
	return cmd
}
