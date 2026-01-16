package editor

import (
	"context"

	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

// Command defines the editor command configuration
var Command = &cli.Command{
	Name:      "editor",
	Usage:     "start a editor window",
	UsageText: "tm editor",
	Action:    action,
}

// action handles the editor command execution
// Renames the current window to "editor" and starts nvim
func action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()

	tmux.RenameWindow("editor")

	tmux.SendKeys("nvim", "C-m")

	return nil
}
