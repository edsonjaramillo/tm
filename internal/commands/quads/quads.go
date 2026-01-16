package quads

import (
	"context"

	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "quads",
	Usage:     "start a 4x4 pane setup",
	UsageText: "tm quads",
	Action:    Action,
	Flags:     []cli.Flag{},
}

func Action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()

	tmux.RenameWindow("shells")

	tmux.SplitWindow("-h")
	tmux.SplitWindow("-v")

	firstPane := tmux.GetPanesInSession(tmux.GetSessionName(), tmux.GetWindowIndex())[0]
	tmux.SelectPane(firstPane)

	tmux.SplitWindow("-v")
	tmux.SelectPane(firstPane)
	tmux.SendKeys("clear", "C-m")

	return nil
}
