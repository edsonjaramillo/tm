package git

import (
	"context"

	"edsonjaramillo/tm/internal/common/shell"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "git",
	Usage:     "start a git window",
	UsageText: "tm git",
	Action:    Action,
	Flags:     []cli.Flag{},
}

func Action(_ context.Context, command *cli.Command) error {
	tmux.AllowIfInSession()

	if !shell.IsGitRepository() {
		shell.Exit("Git repository not found in PWD")
	}

	tmux.SplitWindow("-h")

	firstPane := tmux.GetPanesInSession(tmux.GetSessionName(), tmux.GetWindowIndex())[0]
	tmux.SelectPane(firstPane)

	tmux.SendKeys("lazygit", "C-m")

	return nil
}
