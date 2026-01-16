package git

import (
	"context"

	"edsonjaramillo/tm/internal/common/shell"
	"edsonjaramillo/tm/internal/common/tmux"

	"github.com/urfave/cli/v3"
)

// Command defines the git command configuration
var Command = &cli.Command{
	Name:      "git",
	Usage:     "start a git window",
	UsageText: "tm git",
	Action:    action,
	Flags:     []cli.Flag{},
}

// action handles the git command execution
// Checks if in a git repository, creates a horizontal split pane,
// and launches lazygit in the first pane
func action(_ context.Context, command *cli.Command) error {
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
