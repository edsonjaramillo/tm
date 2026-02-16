package cli

import "github.com/edsonjaramillo/tm/internal/app"

type fakeUseCases struct {
	startReq app.StartRequest
	killReq  app.KillRequest

	detachCalled bool
	startCalled  bool
	killCalled   bool
	gitCalled    bool
	quadsCalled  bool
	dualCalled   bool
	editorCalled bool
	openCalled   bool
	claudeCalled bool
	codexCalled  bool

	err error
}

func (f *fakeUseCases) Start(req app.StartRequest) error {
	f.startCalled = true
	f.startReq = req
	return f.err
}

func (f *fakeUseCases) Detach() error {
	f.detachCalled = true
	return f.err
}

func (f *fakeUseCases) Editor(req app.WindowRequest) error {
	f.editorCalled = true
	_ = req
	return f.err
}

func (f *fakeUseCases) Opencode(req app.WindowRequest) error {
	f.openCalled = true
	_ = req
	return f.err
}

func (f *fakeUseCases) Claude(req app.WindowRequest) error {
	f.claudeCalled = true
	_ = req
	return f.err
}

func (f *fakeUseCases) Codex(req app.WindowRequest) error {
	f.codexCalled = true
	_ = req
	return f.err
}

func (f *fakeUseCases) Quads(req app.WindowRequest) error {
	f.quadsCalled = true
	_ = req
	return f.err
}

func (f *fakeUseCases) Dual(req app.WindowRequest) error {
	f.dualCalled = true
	_ = req
	return f.err
}

func (f *fakeUseCases) Git(req app.WindowRequest) error {
	f.gitCalled = true
	_ = req
	return f.err
}

func (f *fakeUseCases) Kill(req app.KillRequest) error {
	f.killCalled = true
	f.killReq = req
	return f.err
}
