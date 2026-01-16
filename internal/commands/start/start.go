package start

import (
	"context"

	"edsonjaramillo/tm/internal/common/shell"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

// Command defines the start command configuration
var Command = &cli.Command{
	Name:      "start",
	Usage:     "start a tmux session",
	UsageText: "tm start",
	Action:    action,
	Flags: []cli.Flag{
		switchFlag,
		auxFlag,
	},
}

// action handles the start command execution
// Supports three modes:
// 1. --switch flag: switches to existing session, or starts it if not in tmux
// 2. --aux flag: starts session named after current directory with "_aux" suffix
// 3. Default: starts session named after current directory
func action(_ context.Context, command *cli.Command) error {
	sessionSwitch := command.String("switch")

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

	basename := shell.GetBasenamePWD()

	isAux := command.Bool("aux")
	if isAux {
		basename = basename + "_aux"
	}
	tmux.StartSession(basename)

	return nil
}

// switchFlag defines the --switch flag to switch to a session
var switchFlag = &cli.StringFlag{
	Name:     "switch",
	Usage:    "Switch to the session if already inside a tmux session",
	OnlyOnce: true,
}

// auxFlag defines the --aux flag for creating an auxiliary session
var auxFlag = &cli.BoolFlag{
	Name:     "aux",
	Usage:    "Auxiliary flag for secondary session",
	OnlyOnce: true,
}
