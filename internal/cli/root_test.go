package cli

import (
	"bytes"
	"sort"
	"strings"
	"testing"
)

func TestRootCommandIncludesExpectedSubcommands(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)

	got := make([]string, 0, len(cmd.Commands()))
	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			continue
		}
		got = append(got, sub.Name())
	}
	sort.Strings(got)

	want := []string{"claude", "codex", "detach", "dual", "editor", "git", "kill", "opencode", "quads", "start", "yazi"}
	if len(got) != len(want) {
		t.Fatalf("subcommand count = %d, want %d (got=%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("subcommand[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}

func TestRootCommandDefinesGlobalFlags(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)

	if cmd.PersistentFlags().Lookup("color") == nil {
		t.Fatal("persistent flag --color not found")
	}
	if cmd.PersistentFlags().Lookup("no-level") == nil {
		t.Fatal("persistent flag --no-level not found")
	}

	versionFlag := cmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Fatal("flag --version not found")
	}
	if versionFlag.Shorthand != "v" {
		t.Fatalf("version shorthand = %q, want %q", versionFlag.Shorthand, "v")
	}
}

func TestCompletionCommandGeneratesScript(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"completion", "zsh"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if out.Len() == 0 {
		t.Fatal("completion output is empty")
	}
	if !strings.Contains(out.String(), "compdef") {
		t.Fatalf("completion output missing expected zsh marker: %q", out.String())
	}
}

func TestKillCommandMapsSessionArgument(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	cmd.SetArgs([]string{"kill", "my-session"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if !useCases.killCalled {
		t.Fatal("Kill use-case was not called")
	}
	if got, want := useCases.killReq.Session, "my-session"; got != want {
		t.Fatalf("killReq.Session = %q, want %q", got, want)
	}
	if useCases.killReq.All {
		t.Fatal("killReq.All = true, want false")
	}
}

func TestKillCommandMapsAllFlag(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	cmd.SetArgs([]string{"kill", "--all"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if !useCases.killCalled {
		t.Fatal("Kill use-case was not called")
	}
	if !useCases.killReq.All {
		t.Fatal("killReq.All = false, want true")
	}
	if got, want := useCases.killReq.Session, ""; got != want {
		t.Fatalf("killReq.Session = %q, want empty", got)
	}
}

func TestStartCommandMapsSessionArgument(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	cmd.SetArgs([]string{"start", "my-session"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if !useCases.startCalled {
		t.Fatal("Start use-case was not called")
	}
	if got, want := useCases.startReq.Session, "my-session"; got != want {
		t.Fatalf("startReq.Session = %q, want %q", got, want)
	}
	if useCases.startReq.Aux {
		t.Fatal("startReq.Aux = true, want false")
	}
}

func TestStartCommandMapsNoSessionArgument(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	cmd.SetArgs([]string{"start"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if !useCases.startCalled {
		t.Fatal("Start use-case was not called")
	}
	if got, want := useCases.startReq.Session, ""; got != want {
		t.Fatalf("startReq.Session = %q, want empty", got)
	}
	if useCases.startReq.Aux {
		t.Fatal("startReq.Aux = true, want false")
	}
}

func TestStartCommandMapsAuxFlagWithSessionArgument(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	cmd.SetArgs([]string{"start", "my-session", "--aux"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if !useCases.startCalled {
		t.Fatal("Start use-case was not called")
	}
	if got, want := useCases.startReq.Session, "my-session"; got != want {
		t.Fatalf("startReq.Session = %q, want %q", got, want)
	}
	if !useCases.startReq.Aux {
		t.Fatal("startReq.Aux = false, want true")
	}
}

func TestYaziCommandMapsDefaultRequest(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	cmd.SetArgs([]string{"yazi"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if !useCases.yaziCalled {
		t.Fatal("Yazi use-case was not called")
	}
	if useCases.yaziReq.New {
		t.Fatal("yaziReq.New = true, want false")
	}
}

func TestYaziCommandMapsNewFlag(t *testing.T) {
	useCases := &fakeUseCases{}
	cmd := newRootCmdWithUseCases(useCases)
	cmd.SetArgs([]string{"yazi", "--new"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() error = %v, want nil", err)
	}
	if !useCases.yaziCalled {
		t.Fatal("Yazi use-case was not called")
	}
	if !useCases.yaziReq.New {
		t.Fatal("yaziReq.New = false, want true")
	}
}
