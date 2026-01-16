package detach

import (
	"context"

	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

// Command defines the detach command configuration
var Command = &cli.Command{
	Name:      "detach",
	Usage:     "detach a tmux session",
	UsageText: "tm detach",
	Action:    action,
}

// action handles the detach command execution
func action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()
	tmux.DetachFromSession()
	return nil
}
