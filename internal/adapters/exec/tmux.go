package execadapter

import (
	"errors"
	"io"
	"os"
	osExec "os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/edsonjaramillo/tm/internal/domain"
)

// Tmux implements ports.Tmux by shelling out to tmux.
type Tmux struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewTmux(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Tmux {
	if stdin == nil {
		stdin = os.Stdin
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	return &Tmux{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (t *Tmux) StartSession(name string) error {
	return t.runInteractive("new-session", "-A", "-s", name)
}

func (t *Tmux) StartAuxSession(target string, name string) error {
	return t.runInteractive("new-session", "-A", "-s", name, "-t", target)
}

func (t *Tmux) DetachClient() error {
	return t.runInteractive("detach-client")
}

func (t *Tmux) NewWindow(name string) error {
	return t.runInteractive("new-window", "-n", name)
}

func (t *Tmux) RenameWindow(name string) error {
	return t.runInteractive("rename-window", name)
}

func (t *Tmux) SplitWindow(direction string) error {
	if err := domain.ValidateSplitDirection(direction); err != nil {
		return err
	}
	return t.runInteractive("split-window", direction)
}

func (t *Tmux) CurrentSessionName() (string, error) {
	output, err := t.combinedOutput("display-message", "-p", "#S")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (t *Tmux) CurrentWindowIndex() (string, error) {
	output, err := t.combinedOutput("display-message", "-p", "#I")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (t *Tmux) Panes(session string, window string) ([]int, error) {
	target := session + ":" + window
	output, err := t.combinedOutput("list-panes", "-t", target, "-F", "#P")
	if err != nil {
		return nil, err
	}

	lines := splitLines(output)
	panes := make([]int, 0, len(lines))
	for _, line := range lines {
		paneStr := strings.TrimPrefix(line, "%")
		pane, convErr := strconv.Atoi(paneStr)
		if convErr != nil {
			return nil, convErr
		}
		panes = append(panes, pane)
	}
	sort.Ints(panes)
	return panes, nil
}

func (t *Tmux) SelectPane(pane int) error {
	return t.runInteractive("select-pane", "-t", strconv.Itoa(pane))
}

func (t *Tmux) SendKeys(keys ...string) error {
	args := append([]string{"send-keys"}, keys...)
	return t.runInteractive(args...)
}

func (t *Tmux) ListSessions() ([]string, error) {
	output, err := t.combinedOutput("list-sessions", "-F", "#S")
	if err != nil {
		if isExitError(err) {
			return []string{}, nil
		}
		return nil, err
	}
	return splitLines(output), nil
}

func (t *Tmux) SessionExists(name string) (bool, error) {
	cmd := osExec.Command("tmux", "has-session", "-t", name)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	if isExitError(err) {
		return false, nil
	}
	return false, err
}

func (t *Tmux) KillSession(name string) error {
	return t.runInteractive("kill-session", "-t", name)
}

func (t *Tmux) KillServer() error {
	return t.runInteractive("kill-server")
}

func (t *Tmux) runInteractive(args ...string) error {
	cmd := osExec.Command("tmux", args...)
	cmd.Stdin = t.stdin
	cmd.Stdout = t.stdout
	cmd.Stderr = t.stderr
	return cmd.Run()
}

func (t *Tmux) combinedOutput(args ...string) ([]byte, error) {
	cmd := osExec.Command("tmux", args...)
	return cmd.CombinedOutput()
}

func splitLines(output []byte) []string {
	splits := strings.Split(string(output), "\n")
	lines := make([]string, 0, len(splits))
	for _, line := range splits {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func isExitError(err error) bool {
	var exitErr *osExec.ExitError
	return errors.As(err, &exitErr)
}
