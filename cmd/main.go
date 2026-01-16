// Package main is the entry point for the CLI application
package main

import (
	"context"
	"log"
	"os"

	"edsonjaramillo/tm/internal/commands/detach"
	"edsonjaramillo/tm/internal/commands/editor"
	"edsonjaramillo/tm/internal/commands/git"
	"edsonjaramillo/tm/internal/commands/kill"
	"edsonjaramillo/tm/internal/commands/opencode"
	"edsonjaramillo/tm/internal/commands/quads"
	"edsonjaramillo/tm/internal/commands/start"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "tm",
		Version: "0.1.0",
		Usage:   `tm is a tmux helper`,
		Commands: []*cli.Command{
			start.Command,
			detach.Command,
			editor.Command,
			opencode.Command,
			quads.Command,
			git.Command,
			kill.Command,
		},
		HideHelpCommand: true,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
