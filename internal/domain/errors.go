package domain

import (
	"errors"
	"fmt"
)

// Code identifies a domain-level error category.
type Code string

const (
	CodeNotInTmuxSession Code = "not_in_tmux_session"
	CodeNoAuxTarget      Code = "no_aux_target"
	CodeNotGitRepository Code = "not_git_repository"
	CodeNoSessionsToKill Code = "no_sessions_to_kill"
	CodeMissingSession   Code = "missing_session"
	CodeSessionNotFound  Code = "session_not_found"
	CodeInvalidSplit     Code = "invalid_split"
	CodeNoPanesFound     Code = "no_panes_found"
	CodeGetWorkingDir    Code = "get_working_dir"
)

// Error is a typed domain error with a stable message.
type Error struct {
	Code    Code
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// Is reports whether err is a domain error with the provided code.
func Is(err error, code Code) bool {
	var domainErr *Error
	if !errors.As(err, &domainErr) {
		return false
	}
	return domainErr.Code == code
}

func errf(code Code, format string, args ...any) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func ErrNotInTmuxSession() error {
	return errf(CodeNotInTmuxSession, "You are not in a tmux session")
}

func ErrNoAuxTarget(name string) error {
	return errf(CodeNoAuxTarget, "No target session named '%s' found for auxiliary session.", name)
}

func ErrNotGitRepository() error {
	return errf(CodeNotGitRepository, "Git repository not found in PWD")
}

func ErrNoSessionsToKill() error {
	return errf(CodeNoSessionsToKill, "No tmux sessions to kill")
}

func ErrMissingSession() error {
	return errf(CodeMissingSession, "Please provide a session name or use the --all flag to kill all sessions")
}

func ErrSessionNotFound(name string) error {
	return errf(CodeSessionNotFound, "Session %s does not exist", name)
}

func ErrInvalidSplitDirection(direction string) error {
	_ = direction
	return errf(CodeInvalidSplit, "Invalid direction for split window. Use -h or -v.")
}

func ErrNoPanesFound() error {
	return errf(CodeNoPanesFound, "Error getting panes output: no panes found")
}

func ErrGetWorkingDir() error {
	return errf(CodeGetWorkingDir, "Could not get current working directory")
}
