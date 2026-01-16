package kill

import (
	"context"

	"edsonjaramillo/tm/internal/common/arguments"
	"edsonjaramillo/tm/internal/common/shell"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

// Command defines the kill command configuration
var Command = &cli.Command{
	Name:      "kill",
	Usage:     "kill a tmux session",
	UsageText: "tm kill",
	Arguments: []cli.Argument{
		arguments.SessionArg,
	},
	Action: action,

	Flags: []cli.Flag{
		allFlag,
	},
}

// action handles the kill command execution
// Supports three modes:
// 1. --all flag: kills all tmux sessions
// 2. No session name: kills session named after current directory
// 3. Session name provided: kills the specified session
func action(_ context.Context, command *cli.Command) error {
	all := command.Bool("all")

	if all {
		_, numberOfSessions := tmux.ListSessions()
		if numberOfSessions == 0 {
			shell.Exit("No tmux sessions to kill")
		}
		tmux.KillAllSessions()
		return nil
	}

	session := command.StringArg("session")
	if session == "" {
		tmux.AllowIfInSession()

		basename := shell.GetBasenamePWD()

		if !tmux.CheckIfSessionExists(basename) {
			shell.Exit(basename + " session does not exist")
		}

		tmux.KillSession(basename)
		return nil
	}

	if !tmux.CheckIfSessionExists(session) {
		shell.Exit("Session " + session + " does not exist")
	}

	tmux.KillSession(session)

	return nil
}

// allFlag defines the --all flag to kill all sessions
var allFlag = &cli.BoolFlag{
	Name:     "all",
	Usage:    "kill all tmux sessions",
	OnlyOnce: true,
	Required: false,
}
