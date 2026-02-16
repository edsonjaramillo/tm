package app

// StartRequest captures inputs for tm start.
type StartRequest struct {
	Aux     bool
	Session string
}

// WindowRequest captures inputs for commands that manage tmux windows.
type WindowRequest struct {
	New bool
}

// KillRequest captures inputs for tm kill.
type KillRequest struct {
	All     bool
	Session string
}
