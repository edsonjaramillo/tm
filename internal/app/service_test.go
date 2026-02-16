package app

import (
	"errors"
	"reflect"
	"testing"
)

type fakeSystem struct {
	basename    string
	basenameErr error
	inTmux      bool
	gitRepo     bool
}

func (f *fakeSystem) BasenamePWD() (string, error) {
	if f.basenameErr != nil {
		return "", f.basenameErr
	}
	return f.basename, nil
}

func (f *fakeSystem) InTmuxSession() bool {
	return f.inTmux
}

func (f *fakeSystem) IsGitRepository() bool {
	return f.gitRepo
}

type fakeTmux struct {
	startSessionName string
	startAuxTarget   string
	startAuxName     string
	killSessionName  string
	newWindowName    string
	renameWindowName string

	currentSessionName string
	currentWindowIndex string
	panes              []int
	sessionExists      map[string]bool
	listSessions       []string

	splitDirections []string
	selectedPanes   []int
	sentKeys        [][]string

	startSessionErr   error
	startAuxErr       error
	detachErr         error
	newWindowErr      error
	renameWindowErr   error
	splitWindowErr    error
	currentSessionErr error
	currentWindowErr  error
	panesErr          error
	selectPaneErr     error
	sendKeysErr       error
	listSessionsErr   error
	sessionExistsErr  error
	killSessionErr    error
	killServerErr     error

	detachCalled     bool
	killServerCalled bool
}

func (f *fakeTmux) StartSession(name string) error {
	f.startSessionName = name
	return f.startSessionErr
}

func (f *fakeTmux) StartAuxSession(target string, name string) error {
	f.startAuxTarget = target
	f.startAuxName = name
	return f.startAuxErr
}

func (f *fakeTmux) DetachClient() error {
	f.detachCalled = true
	return f.detachErr
}

func (f *fakeTmux) NewWindow(name string) error {
	f.newWindowName = name
	return f.newWindowErr
}

func (f *fakeTmux) RenameWindow(name string) error {
	f.renameWindowName = name
	return f.renameWindowErr
}

func (f *fakeTmux) SplitWindow(direction string) error {
	f.splitDirections = append(f.splitDirections, direction)
	return f.splitWindowErr
}

func (f *fakeTmux) CurrentSessionName() (string, error) {
	if f.currentSessionErr != nil {
		return "", f.currentSessionErr
	}
	return f.currentSessionName, nil
}

func (f *fakeTmux) CurrentWindowIndex() (string, error) {
	if f.currentWindowErr != nil {
		return "", f.currentWindowErr
	}
	return f.currentWindowIndex, nil
}

func (f *fakeTmux) Panes(session string, window string) ([]int, error) {
	if f.panesErr != nil {
		return nil, f.panesErr
	}
	return f.panes, nil
}

func (f *fakeTmux) SelectPane(pane int) error {
	f.selectedPanes = append(f.selectedPanes, pane)
	return f.selectPaneErr
}

func (f *fakeTmux) SendKeys(keys ...string) error {
	copied := append([]string{}, keys...)
	f.sentKeys = append(f.sentKeys, copied)
	return f.sendKeysErr
}

func (f *fakeTmux) ListSessions() ([]string, error) {
	if f.listSessionsErr != nil {
		return nil, f.listSessionsErr
	}
	return f.listSessions, nil
}

func (f *fakeTmux) SessionExists(name string) (bool, error) {
	if f.sessionExistsErr != nil {
		return false, f.sessionExistsErr
	}
	exists, found := f.sessionExists[name]
	if !found {
		return false, nil
	}
	return exists, nil
}

func (f *fakeTmux) KillSession(name string) error {
	f.killSessionName = name
	return f.killSessionErr
}

func (f *fakeTmux) KillServer() error {
	f.killServerCalled = true
	return f.killServerErr
}

func TestStartAuxMissingTargetSession(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{sessionExists: map[string]bool{"repo": false}}
	system := &fakeSystem{basename: "repo"}
	svc := NewService(tmux, system)

	err := svc.Start(StartRequest{Aux: true})
	if err == nil {
		t.Fatal("Start(StartRequest{Aux:true}) error = nil, want non-nil")
	}
	if got, want := err.Error(), "No target session named 'repo' found for auxiliary session."; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if !IsKind(err, KindValidation) {
		t.Fatalf("IsKind(err, KindValidation) = false, want true (err=%v)", err)
	}
}

func TestStartUsesProvidedSessionWhenItExists(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{sessionExists: map[string]bool{"provided": true}}
	system := &fakeSystem{basename: "repo"}
	svc := NewService(tmux, system)

	err := svc.Start(StartRequest{Session: "provided"})
	if err != nil {
		t.Fatalf("Start(StartRequest{Session:%q}) error = %v, want nil", "provided", err)
	}
	if got, want := tmux.startSessionName, "provided"; got != want {
		t.Fatalf("startSessionName = %q, want %q", got, want)
	}
}

func TestStartWithProvidedMissingSessionReturnsValidationError(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{sessionExists: map[string]bool{"other": true}}
	system := &fakeSystem{basename: "repo"}
	svc := NewService(tmux, system)

	err := svc.Start(StartRequest{Session: "provided"})
	if err == nil {
		t.Fatal("Start(StartRequest{Session:\"provided\"}) error = nil, want non-nil")
	}
	if got, want := err.Error(), "Session provided does not exist"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if !IsKind(err, KindValidation) {
		t.Fatalf("IsKind(err, KindValidation) = false, want true (err=%v)", err)
	}
}

func TestStartAuxUsesProvidedSessionAsTarget(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{sessionExists: map[string]bool{"provided": true}}
	system := &fakeSystem{basename: "repo"}
	svc := NewService(tmux, system)

	err := svc.Start(StartRequest{Aux: true, Session: "provided"})
	if err != nil {
		t.Fatalf("Start(StartRequest{Aux:true, Session:%q}) error = %v, want nil", "provided", err)
	}
	if got, want := tmux.startAuxTarget, "provided"; got != want {
		t.Fatalf("startAuxTarget = %q, want %q", got, want)
	}
	if got, want := tmux.startAuxName, "provided_aux"; got != want {
		t.Fatalf("startAuxName = %q, want %q", got, want)
	}
}

func TestStartWithoutSessionUsesBasename(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{}
	system := &fakeSystem{basename: "repo"}
	svc := NewService(tmux, system)

	err := svc.Start(StartRequest{})
	if err != nil {
		t.Fatalf("Start(StartRequest{}) error = %v, want nil", err)
	}
	if got, want := tmux.startSessionName, "repo"; got != want {
		t.Fatalf("startSessionName = %q, want %q", got, want)
	}
}

func TestDetachOutsideTmuxSession(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{}
	system := &fakeSystem{inTmux: false}
	svc := NewService(tmux, system)

	err := svc.Detach()
	if err == nil {
		t.Fatal("Detach() error = nil, want non-nil")
	}
	if got, want := err.Error(), "You are not in a tmux session"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if !IsKind(err, KindValidation) {
		t.Fatalf("IsKind(err, KindValidation) = false, want true (err=%v)", err)
	}
}

func TestGitOutsideRepository(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{}
	system := &fakeSystem{inTmux: true, gitRepo: false}
	svc := NewService(tmux, system)

	err := svc.Git(WindowRequest{})
	if err == nil {
		t.Fatal("Git(WindowRequest{}) error = nil, want non-nil")
	}
	if got, want := err.Error(), "Git repository not found in PWD"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if !IsKind(err, KindValidation) {
		t.Fatalf("IsKind(err, KindValidation) = false, want true (err=%v)", err)
	}
}

func TestKillUsesCurrentSessionWhenSessionArgMissing(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{
		currentSessionName: "active-session",
		sessionExists:      map[string]bool{"active-session": true},
	}
	system := &fakeSystem{}
	svc := NewService(tmux, system)

	err := svc.Kill(KillRequest{})
	if err != nil {
		t.Fatalf("Kill(KillRequest{}) error = %v, want nil", err)
	}
	if got, want := tmux.killSessionName, "active-session"; got != want {
		t.Fatalf("killSessionName = %q, want %q", got, want)
	}
}

func TestKillMissingSessionOutsideTmux(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{currentSessionErr: errors.New("not in tmux")}
	system := &fakeSystem{}
	svc := NewService(tmux, system)

	err := svc.Kill(KillRequest{})
	if err == nil {
		t.Fatal("Kill(KillRequest{}) error = nil, want non-nil")
	}
	if got, want := err.Error(), "You are not in a tmux session"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if !IsKind(err, KindValidation) {
		t.Fatalf("IsKind(err, KindValidation) = false, want true (err=%v)", err)
	}
}

func TestKillAllFailsWhenNoSessions(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{listSessions: []string{}}
	system := &fakeSystem{}
	svc := NewService(tmux, system)

	err := svc.Kill(KillRequest{All: true})
	if err == nil {
		t.Fatal("Kill(KillRequest{All:true}) error = nil, want non-nil")
	}
	if got, want := err.Error(), "No tmux sessions to kill"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
	if !IsKind(err, KindValidation) {
		t.Fatalf("IsKind(err, KindValidation) = false, want true (err=%v)", err)
	}
}

func TestDualCreatesHorizontalSplitAndClearsFirstPane(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{
		currentSessionName: "s",
		currentWindowIndex: "1",
		panes:              []int{1, 2},
	}
	system := &fakeSystem{inTmux: true}
	svc := NewService(tmux, system)

	err := svc.Dual(WindowRequest{})
	if err != nil {
		t.Fatalf("Dual(WindowRequest{}) error = %v, want nil", err)
	}
	if got, want := tmux.renameWindowName, "shells"; got != want {
		t.Fatalf("renameWindowName = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(tmux.splitDirections, []string{"-h"}) {
		t.Fatalf("splitDirections = %#v, want %#v", tmux.splitDirections, []string{"-h"})
	}
	if !reflect.DeepEqual(tmux.selectedPanes, []int{1}) {
		t.Fatalf("selectedPanes = %#v, want %#v", tmux.selectedPanes, []int{1})
	}
	if len(tmux.sentKeys) != 1 {
		t.Fatalf("sentKeys calls = %d, want 1", len(tmux.sentKeys))
	}
	if !reflect.DeepEqual(tmux.sentKeys[0], []string{"clear", "C-m"}) {
		t.Fatalf("sentKeys[0] = %#v, want %#v", tmux.sentKeys[0], []string{"clear", "C-m"})
	}
}

func TestClaudeSendsClaudeBinary(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{}
	system := &fakeSystem{inTmux: true}
	svc := NewService(tmux, system)

	err := svc.Claude(WindowRequest{})
	if err != nil {
		t.Fatalf("Claude(WindowRequest{}) error = %v, want nil", err)
	}
	if len(tmux.sentKeys) != 1 {
		t.Fatalf("sentKeys calls = %d, want 1", len(tmux.sentKeys))
	}
	if !reflect.DeepEqual(tmux.sentKeys[0], []string{"claude", "C-m"}) {
		t.Fatalf("sentKeys[0] = %#v, want %#v", tmux.sentKeys[0], []string{"claude", "C-m"})
	}
}

func TestCodexSendsCodexBinary(t *testing.T) {
	t.Parallel()

	tmux := &fakeTmux{}
	system := &fakeSystem{inTmux: true}
	svc := NewService(tmux, system)

	err := svc.Codex(WindowRequest{})
	if err != nil {
		t.Fatalf("Codex(WindowRequest{}) error = %v, want nil", err)
	}
	if len(tmux.sentKeys) != 1 {
		t.Fatalf("sentKeys calls = %d, want 1", len(tmux.sentKeys))
	}
	if !reflect.DeepEqual(tmux.sentKeys[0], []string{"codex", "C-m"}) {
		t.Fatalf("sentKeys[0] = %#v, want %#v", tmux.sentKeys[0], []string{"codex", "C-m"})
	}
}
