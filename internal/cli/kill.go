package cli

import (
	"strings"

	"edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

type sessionLister interface {
	ListSessions() ([]string, error)
}

func newKillCmd(useCases app.UseCases, sessions sessionLister) *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "kill [session]",
		Short: "kill a tmux session",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := ""
			if len(args) == 1 {
				session = args[0]
			}
			return useCases.Kill(app.KillRequest{
				All:     all,
				Session: session,
			})
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "kill all tmux sessions")
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 || all || sessions == nil {
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
