package arguments

import "github.com/urfave/cli/v3"

// SessionArg defines the session argument for commands that accept a session name
var SessionArg = &cli.StringArg{
	Name:      "session",
	UsageText: "name of the tmux session",
}
