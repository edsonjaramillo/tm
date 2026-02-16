package ports

// Tmux abstracts tmux interactions for app orchestration.
type Tmux interface {
	StartSession(name string) error
	StartAuxSession(target string, name string) error
	DetachClient() error
	NewWindow(name string) error
	RenameWindow(name string) error
	SplitWindow(direction string) error
	CurrentSessionName() (string, error)
	CurrentWindowIndex() (string, error)
	Panes(session string, window string) ([]int, error)
	SelectPane(pane int) error
	SendKeys(keys ...string) error
	ListSessions() ([]string, error)
	SessionExists(name string) (bool, error)
	KillSession(name string) error
	KillServer() error
}
