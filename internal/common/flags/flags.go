package flags

import "github.com/urfave/cli/v3"

var NewFlag = &cli.BoolFlag{
	Name:     "new",
	Usage:    "create a new tmux window if none exists",
	OnlyOnce: true,
}
