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
		auxFlag,
	},
}

// action handles the start command execution
// Supports three modes:
// 2. --aux flag: starts session named after current directory with "_aux" suffix
// 3. Default: starts session named after current directory
func action(_ context.Context, command *cli.Command) error {
	basename := shell.GetBasenamePWD()

	isAux := command.Bool("aux")
	if isAux {
		basename = basename + "_aux"
	}
	tmux.StartSession(basename)

	return nil
}

// auxFlag defines the --aux flag for creating an auxiliary session
var auxFlag = &cli.BoolFlag{
	Name:     "aux",
	Usage:    "Auxiliary flag for secondary session",
	OnlyOnce: true,
}
