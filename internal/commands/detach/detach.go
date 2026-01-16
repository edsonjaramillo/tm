package detach

import (
	"context"

	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "detach",
	Usage:     "detach a tmux session",
	UsageText: "tm detach",
	Action:    Action,
}

func Action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()
	tmux.DetachFromSession()
	return nil
}
