package tmux

import (
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"edsonjaramillo/tm/internal/common/shell"
)

// Session Functions

// StartSession creates or attaches to a tmux session with the given name
func StartSession(name string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "new-session", "-A", "-s", name)
}

func StartAuxSession(target string, name string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "new-session", "-A", "-s", name, "-t", target)
}

// KillSession terminates the specified tmux session
func KillSession(name string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "kill-session", "-t", name)
}

// KillAllSessions terminates the tmux server and all sessions
func KillAllSessions() *exec.Cmd {
	return shell.CmdInteractive("tmux", "kill-server")
}

// ListSessions returns a list of all tmux session names and their count
func ListSessions() ([]string, int) {
	cmd := shell.Cmd("tmux", "list-sessions", "-F", "#S")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{}, 0
	}

	sessions := outputCleaner(output, "\n")

	return sessions, len(sessions)
}

// DetachFromSession detaches the current client from the session
func DetachFromSession() *exec.Cmd {
	return shell.CmdInteractive("tmux", "detach-client")
}

// RenameSession renames the current session to the new name
func RenameSession(newName string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "rename-session", newName)
}

// GetSessionName returns the name of the current tmux session
func GetSessionName() string {
	cmd := shell.Cmd("tmux", "display-message", "-p", "#S")
	output, err := cmd.CombinedOutput()
	if err != nil {
		shell.Exit("Error getting session name: " + err.Error())
	}

	return outputCleaner(output, "\n")[0]
}

// GetWindowIndex returns the index of the current window
func GetWindowIndex() string {
	cmd := shell.Cmd("tmux", "display-message", "-p", "#I")
	output, err := cmd.CombinedOutput()
	if err != nil {
		shell.Exit("Error getting window index: " + err.Error())
	}
	return outputCleaner(output, "\n")[0]
}

// ExitIfNotInSession exits if the TMUX environment variable is set
func ExitIfNotInSession() {
	_, found := os.LookupEnv("TMUX")
	if found {
		shell.Exit("You are already in a tmux session")
	}
}

// AllowIfInSession exits if not currently in a tmux session
func AllowIfInSession() {
	_, found := os.LookupEnv("TMUX")
	if !found {
		shell.Exit("You are not in a tmux session")
	}
}

// CheckIfInSession returns true if currently in a tmux session
func CheckIfInSession() bool {
	_, found := os.LookupEnv("TMUX")
	return found
}

// CheckIfSessionExists returns true if a session with the given name exists
func CheckIfSessionExists(name string) bool {
	cmd := shell.Cmd("tmux", "has-session", "-t", name)
	err := cmd.Run()
	return err == nil
}

// Window Functions

// NewWindow creates a new tmux window with the specified name
func NewWindow(name string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "new-window", "-n", name)
}

// SplitWindow splits the current pane in the specified direction (-h for horizontal, -v for vertical)
func SplitWindow(direction string) *exec.Cmd {
	if direction != "-h" && direction != "-v" {
		shell.Exit("Invalid direction for split window. Use -h or -v.")
	}

	return shell.CmdInteractive("tmux", "split-window", direction)
}

// RenameWindow renames the current window to the new name
func RenameWindow(newName string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "rename-window", newName)
}

// Pane Functions

// GetPanesInSession returns a sorted list of pane IDs in the specified session and window
func GetPanesInSession(session string, window string) []int {
	target := session + ":" + window
	cmd := shell.Cmd("tmux", "list-panes", "-t", target, "-F", "#P")
	output, err := cmd.CombinedOutput()
	if err != nil {
		shell.Exit("Error getting panes output: " + err.Error())
	}

	panes := outputCleaner(output, "\n")
	numberedPanes := []int{}
	for _, pane := range panes {
		pane = strings.TrimPrefix(pane, "%")
		paneNum, err := strconv.Atoi(pane)
		if err != nil {
			shell.Exit("Error parsing pane number: " + err.Error())
		}
		numberedPanes = append(numberedPanes, paneNum)
	}

	sort.Ints(numberedPanes)

	return numberedPanes
}

// SelectPane selects the pane with the specified ID
func SelectPane(pane int) *exec.Cmd {
	return shell.CmdInteractive("tmux", "select-pane", "-t", strconv.Itoa(pane))
}

// Misc Functions

// SendKeys simulates typing the specified keys into the current pane
func SendKeys(keys ...string) *exec.Cmd {
	args := append([]string{"send-keys"}, keys...)
	return shell.CmdInteractive("tmux", args...)
}

// outputCleaner splits command output by delimiter and removes empty strings
func outputCleaner(output []byte, delimiter string) []string {
	splits := strings.Split(string(output), delimiter)
	lines := []string{}
	for _, line := range splits {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
