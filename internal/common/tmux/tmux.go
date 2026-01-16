package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"edsonjaramillo/tm/internal/common/shell"
)

// Session Functions
func StartSession(name string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "new-session", "-A", "-s", name)
}

func KillSession(name string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "kill-session", "-t", name)
}

func KillAllSessions() *exec.Cmd {
	return shell.CmdInteractive("tmux", "kill-server")
}

func SwitchSession(session string) *exec.Cmd {
	name := session + ":1"
	return shell.CmdInteractive("tmux", "switch-client", "-t", name)
}

func ListSessions() ([]string, int) {
	cmd := shell.Cmd("tmux", "list-sessions", "-F", "#S")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{}, 0
	}

	sessions := outputCleaner(output, "\n")

	return sessions, len(sessions)
}

func DetachFromSession() *exec.Cmd {
	return shell.CmdInteractive("tmux", "detach-client")
}

func RenameSession(newName string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "rename-session", newName)
}

func GetSessionName() string {
	cmd := shell.Cmd("tmux", "display-message", "-p", "#S")
	output, err := cmd.CombinedOutput()
	if err != nil {
		shell.Exit("Error getting session name: " + err.Error())
	}

	return outputCleaner(output, "\n")[0]
}

func GetWindowIndex() string {
	cmd := shell.Cmd("tmux", "display-message", "-p", "#I")
	output, err := cmd.CombinedOutput()
	if err != nil {
		shell.Exit("Error getting window index: " + err.Error())
	}
	return outputCleaner(output, "\n")[0]
}

func ExitIfNotInSession() {
	_, found := os.LookupEnv("TMUX")
	if found {
		shell.Exit("You are already in a tmux session")
	}
}

func AllowIfInSession() {
	_, found := os.LookupEnv("TMUX")
	if !found {
		shell.Exit("You are not in a tmux session")
	}
}

func CheckIfInSession() bool {
	_, found := os.LookupEnv("TMUX")
	return found
}

func CheckIfSessionExists(name string) bool {
	cmd := shell.Cmd("tmux", "has-session", "-t", name)
	err := cmd.Run()
	return err == nil
}

// Window Functions
func NewWindow(name string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "new-window", "-n", name)
}

func SplitWindow(direction string) *exec.Cmd {
	if direction != "-h" && direction != "-v" {
		shell.Exit("Invalid direction for split window. Use -h or -v.")
	}

	return shell.CmdInteractive("tmux", "split-window", direction)
}

func RenameWindow(newName string) *exec.Cmd {
	return shell.CmdInteractive("tmux", "rename-window", newName)
}

// Pane Functions
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

	for _, pane := range numberedPanes {
		fmt.Println("Found pane:", pane)
	}

	return numberedPanes
}

func SelectPane(pane int) *exec.Cmd {
	return shell.CmdInteractive("tmux", "select-pane", "-t", strconv.Itoa(pane))
}

// Misc Functions
func SendKeys(keys ...string) *exec.Cmd {
	args := append([]string{"send-keys"}, keys...)
	return shell.CmdInteractive("tmux", args...)
}

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
