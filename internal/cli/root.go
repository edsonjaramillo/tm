package cli

import (
	"fmt"
	"os"

	execadapter "edsonjaramillo/tm/internal/adapters/exec"
	"edsonjaramillo/tm/internal/app"

	"github.com/spf13/cobra"
)

// Version is the CLI version string set at build time via ldflags.
var Version = "0.1.0"

// Execute constructs and runs the root CLI command tree.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	system := execadapter.NewSystem()
	tmux := execadapter.NewTmux(os.Stdin, os.Stdout, os.Stderr)
	service := app.NewService(tmux, system)
	return newRootCmdWithDependencies(service, tmux)
}

func newRootCmdWithUseCases(useCases app.UseCases) *cobra.Command {
	return newRootCmdWithDependencies(useCases, nil)
}

func newRootCmdWithDependencies(useCases app.UseCases, sessions sessionLister) *cobra.Command {
	var colorFlag string
	var noLevelFlag bool

	cmd := &cobra.Command{
		Use:           "tm",
		Short:         "tm is a tmux helper",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			showVersion, err := cmd.Flags().GetBool("version")
			if err != nil {
				return err
			}
			if showVersion {
				_, printErr := fmt.Fprintln(cmd.OutOrStdout(), Version)
				return printErr
			}
			return cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			mode, err := parseColorMode(colorFlag)
			if err != nil {
				return err
			}

			showLevel := true
			noLevelFromEnv, envSet, err := parseNoLevelEnv()
			if err != nil {
				return err
			}
			if envSet {
				showLevel = !noLevelFromEnv
			}
			if cmd.Flags().Changed("no-level") {
				showLevel = !noLevelFlag
			}

			setOutputColorMode(mode)
			setOutputShowLevel(showLevel)
			return nil
		},
	}

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	cmd.InitDefaultCompletionCmd()

	cmd.PersistentFlags().StringVar(&colorFlag, "color", colorModeAuto, "Colorize output: auto, always, never")
	cmd.PersistentFlags().BoolVar(&noLevelFlag, "no-level", false, "Hide output level labels (INFO, OK, WARN, ERROR)")
	cmd.Flags().BoolP("version", "v", false, "print the version")

	cmd.AddCommand(newStartCmd(useCases, sessions))
	cmd.AddCommand(newDetachCmd(useCases))
	cmd.AddCommand(newEditorCmd(useCases))
	cmd.AddCommand(newOpencodeCmd(useCases))
	cmd.AddCommand(newClaudeCmd(useCases))
	cmd.AddCommand(newCodexCmd(useCases))
	cmd.AddCommand(newQuadsCmd(useCases))
	cmd.AddCommand(newDualCmd(useCases))
	cmd.AddCommand(newGitCmd(useCases))
	cmd.AddCommand(newKillCmd(useCases, sessions))

	return cmd
}
