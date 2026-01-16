package arguments

import "github.com/urfave/cli/v3"

var SessionArg = &cli.StringArg{
	Name:      "session",
	UsageText: "name of the tmux session",
}
