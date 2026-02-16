package cli

import (
	"edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

func newGitCmd(useCases app.UseCases) *cobra.Command {
	var newWindow bool

	cmd := &cobra.Command{
		Use:   "git",
		Short: "start a git window",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCases.Git(app.WindowRequest{New: newWindow})
		},
	}

	cmd.Flags().BoolVar(&newWindow, "new", false, "create a new tmux window if none exists")
	return cmd
}
