package editor

import (
	"context"

	"edsonjaramillo/tm/internal/common/flags"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

// Command defines the editor command configuration
var Command = &cli.Command{
	Name:      "editor",
	Usage:     "start a editor window",
	UsageText: "tm editor",
	Action:    action,
	Flags: []cli.Flag{
		flags.NewFlag,
	},
}

// action handles the editor command execution
// Renames the current window to "editor" and starts nvim
func action(_ context.Context, command *cli.Command) error {
	needsNewWindow := command.Bool("new")
	if needsNewWindow {
		tmux.NewWindow("editor")
	} else {
		tmux.RenameWindow("editor")
	}

	tmux.AllowIfInSession()

	tmux.RenameWindow("editor")

	tmux.SendKeys("nvim", "C-m")

	return nil
}
