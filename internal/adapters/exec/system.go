package execadapter

import (
	"os"
	osExec "os/exec"
	"path/filepath"
)

// System implements ports.System with direct OS/process calls.
type System struct{}

func NewSystem() *System {
	return &System{}
}

func (s *System) BasenamePWD() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(pwd), nil
}

func (s *System) InTmuxSession() bool {
	_, found := os.LookupEnv("TMUX")
	return found
}

func (s *System) IsGitRepository() bool {
	cmd := osExec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}
