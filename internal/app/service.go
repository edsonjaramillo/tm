package app

import (
	"strings"

	"github.com/edsonjaramillo/tm/internal/domain"
	"github.com/edsonjaramillo/tm/internal/ports"
)

// UseCases defines the application operations exposed to the CLI transport.
type UseCases interface {
	Start(req StartRequest) error
	Detach() error
	Editor(req WindowRequest) error
	Opencode(req WindowRequest) error
	Claude(req WindowRequest) error
	Codex(req WindowRequest) error
	Yazi(req WindowRequest) error
	Quads(req WindowRequest) error
	Dual(req WindowRequest) error
	Git(req WindowRequest) error
	Kill(req KillRequest) error
}

// Service orchestrates tmux workflows using injected ports.
type Service struct {
	tmux   ports.Tmux
	system ports.System
}

// NewService builds a Service with concrete ports.
func NewService(tmux ports.Tmux, system ports.System) *Service {
	return &Service{
		tmux:   tmux,
		system: system,
	}
}

func (s *Service) Start(req StartRequest) error {
	basename, err := s.system.BasenamePWD()
	if err != nil {
		return newValidationError(domain.ErrGetWorkingDir())
	}

	session := strings.TrimSpace(req.Session)
	if session == "" {
		session = basename
	} else {
		exists, existsErr := s.tmux.SessionExists(session)
		if existsErr != nil {
			return newExecutionError("check tmux session", existsErr)
		}
		if !exists {
			return newValidationError(domain.ErrSessionNotFound(session))
		}
	}

	if req.Aux {
		exists, existsErr := s.tmux.SessionExists(session)
		if existsErr != nil {
			return newExecutionError("check tmux session", existsErr)
		}
		if !exists {
			return newValidationError(domain.ErrNoAuxTarget(session))
		}

		if startErr := s.tmux.StartAuxSession(session, session+"_aux"); startErr != nil {
			return newExecutionError("start tmux auxiliary session", startErr)
		}
		return nil
	}

	if startErr := s.tmux.StartSession(session); startErr != nil {
		return newExecutionError("start tmux session", startErr)
	}
	return nil
}

func (s *Service) Detach() error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.tmux.DetachClient(); err != nil {
		return newExecutionError("detach tmux client", err)
	}
	return nil
}

func (s *Service) Editor(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.prepareWindow("editor", req.New); err != nil {
		return err
	}
	if err := s.tmux.SendKeys("nvim", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}
	return nil
}

func (s *Service) Opencode(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.prepareWindow("opencode", req.New); err != nil {
		return err
	}
	if err := s.tmux.SendKeys("opencode", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}
	return nil
}

func (s *Service) Claude(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.prepareWindow("claude", req.New); err != nil {
		return err
	}
	if err := s.tmux.SendKeys("claude", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}
	return nil
}

func (s *Service) Codex(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.prepareWindow("codex", req.New); err != nil {
		return err
	}
	if err := s.tmux.SendKeys("codex", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}
	return nil
}

func (s *Service) Yazi(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.prepareWindow("yazi", req.New); err != nil {
		return err
	}
	if err := s.tmux.SendKeys("yazi", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}
	return nil
}

func (s *Service) Quads(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.prepareWindow("shells", req.New); err != nil {
		return err
	}

	if err := s.tmux.SplitWindow(domain.SplitHorizontal); err != nil {
		return newExecutionError("split tmux window", err)
	}
	if err := s.tmux.SplitWindow(domain.SplitVertical); err != nil {
		return newExecutionError("split tmux window", err)
	}

	firstPane, err := s.firstPaneInCurrentWindow()
	if err != nil {
		return err
	}

	if err := s.tmux.SelectPane(firstPane); err != nil {
		return newExecutionError("select tmux pane", err)
	}
	if err := s.tmux.SplitWindow(domain.SplitVertical); err != nil {
		return newExecutionError("split tmux window", err)
	}
	if err := s.tmux.SelectPane(firstPane); err != nil {
		return newExecutionError("select tmux pane", err)
	}
	if err := s.tmux.SendKeys("clear", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}

	return nil
}

func (s *Service) Dual(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if err := s.prepareWindow("shells", req.New); err != nil {
		return err
	}

	if err := s.tmux.SplitWindow(domain.SplitHorizontal); err != nil {
		return newExecutionError("split tmux window", err)
	}

	firstPane, err := s.firstPaneInCurrentWindow()
	if err != nil {
		return err
	}

	if err := s.tmux.SelectPane(firstPane); err != nil {
		return newExecutionError("select tmux pane", err)
	}
	if err := s.tmux.SendKeys("clear", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}

	return nil
}

func (s *Service) Git(req WindowRequest) error {
	if err := s.ensureTmuxSession(); err != nil {
		return err
	}
	if !s.system.IsGitRepository() {
		return newValidationError(domain.ErrNotGitRepository())
	}
	if err := s.prepareWindow("git", req.New); err != nil {
		return err
	}
	if err := s.tmux.SendKeys("lazygit", "C-m"); err != nil {
		return newExecutionError("send keys to tmux pane", err)
	}
	return nil
}

func (s *Service) Kill(req KillRequest) error {
	if req.All {
		sessions, err := s.tmux.ListSessions()
		if err != nil {
			sessions = []string{}
		}
		if len(sessions) == 0 {
			return newValidationError(domain.ErrNoSessionsToKill())
		}
		if err := s.tmux.KillServer(); err != nil {
			return newExecutionError("kill tmux server", err)
		}
		return nil
	}

	session := strings.TrimSpace(req.Session)
	if session == "" {
		currentSession, err := s.tmux.CurrentSessionName()
		if err != nil {
			return newValidationError(domain.ErrNotInTmuxSession())
		}
		session = strings.TrimSpace(currentSession)
		if session == "" {
			return newValidationError(domain.ErrNotInTmuxSession())
		}
	}

	exists, err := s.tmux.SessionExists(session)
	if err != nil {
		exists = false
	}
	if !exists {
		return newValidationError(domain.ErrSessionNotFound(session))
	}

	if err := s.tmux.KillSession(session); err != nil {
		return newExecutionError("kill tmux session", err)
	}
	return nil
}

func (s *Service) ensureTmuxSession() error {
	if s.system.InTmuxSession() {
		return nil
	}
	return newValidationError(domain.ErrNotInTmuxSession())
}

func (s *Service) prepareWindow(name string, newWindow bool) error {
	if newWindow {
		if err := s.tmux.NewWindow(name); err != nil {
			return newExecutionError("create tmux window", err)
		}
		return nil
	}
	if err := s.tmux.RenameWindow(name); err != nil {
		return newExecutionError("rename tmux window", err)
	}
	return nil
}

func (s *Service) firstPaneInCurrentWindow() (int, error) {
	sessionName, err := s.tmux.CurrentSessionName()
	if err != nil {
		return 0, newExecutionError("get tmux session name", err)
	}
	windowIndex, err := s.tmux.CurrentWindowIndex()
	if err != nil {
		return 0, newExecutionError("get tmux window index", err)
	}
	panes, err := s.tmux.Panes(sessionName, windowIndex)
	if err != nil {
		return 0, newExecutionError("get tmux panes", err)
	}
	if len(panes) == 0 {
		return 0, newValidationError(domain.ErrNoPanesFound())
	}
	return panes[0], nil
}
