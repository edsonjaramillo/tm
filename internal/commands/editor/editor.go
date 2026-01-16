package editor

import (
	"context"

	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "editor",
	Usage:     "start a editor window",
	UsageText: "tm editor",
	Action:    Action,
}

func Action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()

	tmux.RenameWindow("editor")

	tmux.SendKeys("nvim", "C-m")

	return nil
}
