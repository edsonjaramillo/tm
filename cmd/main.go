package main

import (
	"os"

	"edsonjaramillo/tm/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		if !cli.IsAlreadyReportedFailure(err) {
			_ = cli.PrintRootError(os.Stderr, err)
		}
		os.Exit(1)
	}
}
