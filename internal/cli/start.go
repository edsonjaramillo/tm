package cli

import (
	"strings"

	"github.com/edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

func newStartCmd(useCases app.UseCases, sessions sessionLister) *cobra.Command {
	var aux bool

	cmd := &cobra.Command{
		Use:   "start [session]",
		Short: "start a tmux session",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := ""
			if len(args) == 1 {
				session = args[0]
			}
			return useCases.Start(app.StartRequest{
				Aux:     aux,
				Session: session,
			})
		},
	}

	cmd.Flags().BoolVar(&aux, "aux", false, "Auxiliary flag for secondary session")
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 || sessions == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		list, err := sessions.ListSessions()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		matches := make([]string, 0, len(list))
		for _, session := range list {
			if strings.HasPrefix(session, toComplete) {
				matches = append(matches, session)
			}
		}

		return matches, cobra.ShellCompDirectiveNoFileComp
	}
	return cmd
}
