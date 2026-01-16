package opencode

import (
	"context"

	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "opencode",
	Usage:     "start a opencode window",
	UsageText: "tm opencode",
	Action:    Action,
}

func Action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()

	tmux.RenameWindow("opencode")

	tmux.SendKeys("opencode", "C-m")

	return nil
}
