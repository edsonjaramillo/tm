package opencode

import (
	"context"

	"edsonjaramillo/tm/internal/common/flags"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

// Command defines the opencode command configuration
var Command = &cli.Command{
	Name:      "opencode",
	Usage:     "start a opencode window",
	UsageText: "tm opencode",
	Action:    action,
	Flags: []cli.Flag{
		flags.NewFlag,
	},
}

// action handles the opencode command execution
// Renames the current window to "opencode" and starts the opencode CLI
func action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()

	needsNewWindow := command.Bool("new")
	if needsNewWindow {
		tmux.NewWindow("opencode")
	} else {
		tmux.RenameWindow("opencode")
	}

	tmux.SendKeys("opencode", "C-m")

	return nil
}
