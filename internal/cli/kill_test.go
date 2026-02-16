package cli

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

type fakeSessionLister struct {
	sessions []string
	err      error
}

func (f *fakeSessionLister) ListSessions() ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.sessions, nil
}

func TestKillCompletionReturnsSessions(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{sessions: []string{"alpha", "beta"}}
	cmd := newKillCmd(useCases, sessions)

	got, directive := cmd.ValidArgsFunction(cmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(got) != 2 {
		t.Fatalf("completion count = %d, want 2 (got=%v)", len(got), got)
	}
	if got[0] != "alpha" || got[1] != "beta" {
		t.Fatalf("completion = %v, want [alpha beta]", got)
	}
}

func TestKillCompletionFiltersByPrefix(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{sessions: []string{"alpha", "beta", "alphabeta"}}
	cmd := newKillCmd(useCases, sessions)

	got, directive := cmd.ValidArgsFunction(cmd, []string{}, "alp")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(got) != 2 {
		t.Fatalf("completion count = %d, want 2 (got=%v)", len(got), got)
	}
	if got[0] != "alpha" || got[1] != "alphabeta" {
		t.Fatalf("completion = %v, want [alpha alphabeta]", got)
	}
}

func TestKillCompletionReturnsEmptyWhenAllFlagSet(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{sessions: []string{"alpha", "beta"}}
	cmd := newKillCmd(useCases, sessions)
	if err := cmd.Flags().Set("all", "true"); err != nil {
		t.Fatalf("Flags().Set(all) error = %v, want nil", err)
	}

	got, directive := cmd.ValidArgsFunction(cmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(got) != 0 {
		t.Fatalf("completion = %v, want empty", got)
	}
}

func TestKillCompletionReturnsEmptyWhenListerFails(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{err: errors.New("list failed")}
	cmd := newKillCmd(useCases, sessions)

	got, directive := cmd.ValidArgsFunction(cmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(got) != 0 {
		t.Fatalf("completion = %v, want empty", got)
	}
}
