package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Cmd creates a new command without interactive I/O
func Cmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)

	return cmd
}

// CmdInteractive creates a new command with stdin/stdout connected to the terminal
func CmdInteractive(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return cmd
	}

	return cmd
}

// GetBasenamePWD returns the basename of the current working directory
func GetBasenamePWD() string {
	pwd, err := os.Getwd()
	if err != nil {
		Exit("Could not get current working directory")
	}

	pathParts := strings.Split(pwd, string(os.PathSeparator))

	return pathParts[len(pathParts)-1]
}

// Exit prints a red error message and exits with status 1
func Exit(message string) {
	fmt.Println("\033[31m" + message + "\033[0m")
	os.Exit(1)
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := Cmd("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}
