package ports

// System abstracts host environment checks used by app workflows.
type System interface {
	BasenamePWD() (string, error)
	InTmuxSession() bool
	IsGitRepository() bool
}
