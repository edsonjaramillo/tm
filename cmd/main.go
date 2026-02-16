package main

import (
	"os"

	"github.com/edsonjaramillo/tm/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		if !cli.IsAlreadyReportedFailure(err) {
			_ = cli.PrintRootError(os.Stderr, err)
		}
		os.Exit(1)
	}
}
