package kill

import (
	"context"

	"edsonjaramillo/tm/internal/common/arguments"
	"edsonjaramillo/tm/internal/common/shell"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "kill",
	Usage:     "kill a tmux session",
	UsageText: "tm kill",
	Arguments: []cli.Argument{
		arguments.SessionArg,
	},
	Action: Action,

	Flags: []cli.Flag{
		allFlag,
	},
}

func Action(_ context.Context, command *cli.Command) error {
	all := command.Bool("all")

	// If --all flag is provided delete all tmux sessions
	if all {
		// Check if there are any tmux sessions
		// If there are none, exit with error
		_, numberOfSessions := tmux.ListSessions()
		if numberOfSessions == 0 {
			shell.Exit("No tmux sessions to kill")
		}
		tmux.KillAllSessions()
		return nil
	}

	// If no session is provided
	session := command.StringArg("session")
	if session == "" {
		tmux.AllowIfInSession()

		// Get basename of current directory and check if a tmux session with that name exists
		basename := shell.GetBasenamePWD()

		// If it does not exist, exit with error
		if !tmux.CheckIfSessionExists(basename) {
			shell.Exit(basename + " session does not exist")
		}

		tmux.KillSession(basename)
		return nil
	}

	// Check if the session provided exists
	if !tmux.CheckIfSessionExists(session) {
		shell.Exit("Session " + session + " does not exist")
	}

	tmux.KillSession(session)

	return nil
}

var allFlag = &cli.BoolFlag{
	Name:     "all",
	Usage:    "kill all tmux sessions",
	OnlyOnce: true,
	Required: false,
}
