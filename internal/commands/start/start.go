package start

import (
	"context"

	"edsonjaramillo/tm/internal/common/shell"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "start",
	Usage:     "start a tmux session",
	UsageText: "tm start",
	Action:    Action,
	Flags: []cli.Flag{
		switchFlag,
		auxFlag,
	},
}

func Action(_ context.Context, command *cli.Command) error {
	sessionSwitch := command.String("switch")

	// If switch flag is provided, switch to that session
	if sessionSwitch != "" {
		if !tmux.CheckIfSessionExists(sessionSwitch) {
			shell.Exit("Session " + sessionSwitch + " does not exist")
		}
		if tmux.CheckIfInSession() {
			tmux.SwitchSession(sessionSwitch)
			return nil
		} else {
			tmux.StartSession(sessionSwitch)
			return nil
		}
	}

	// Get basename of current directory and check if a tmux session with that name exists
	basename := shell.GetBasenamePWD()

	// If aux flag is provided, append _aux to the basename
	// aux session is used for having a secondary session for the same project
	isAux := command.Bool("aux")
	if isAux {
		basename = basename + "_aux"
	}
	// Start a new session with the basename
	tmux.StartSession(basename)

	return nil
}

var switchFlag = &cli.StringFlag{
	Name:     "switch",
	Usage:    "Switch to the session if already inside a tmux session",
	OnlyOnce: true,
}

var auxFlag = &cli.BoolFlag{
	Name:     "aux",
	Usage:    "Auxiliary flag for secondary session",
	OnlyOnce: true,
}
