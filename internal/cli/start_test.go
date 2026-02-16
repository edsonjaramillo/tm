package cli

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestStartCompletionReturnsSessions(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{sessions: []string{"alpha", "beta"}}
	cmd := newStartCmd(useCases, sessions)

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

func TestStartCompletionFiltersByPrefix(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{sessions: []string{"alpha", "beta", "alphabeta"}}
	cmd := newStartCmd(useCases, sessions)

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

func TestStartCompletionReturnsEmptyWhenListerFails(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{err: errors.New("list failed")}
	cmd := newStartCmd(useCases, sessions)

	got, directive := cmd.ValidArgsFunction(cmd, []string{}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(got) != 0 {
		t.Fatalf("completion = %v, want empty", got)
	}
}

func TestStartCompletionReturnsEmptyWhenArgAlreadyProvided(t *testing.T) {
	useCases := &fakeUseCases{}
	sessions := &fakeSessionLister{sessions: []string{"alpha", "beta"}}
	cmd := newStartCmd(useCases, sessions)

	got, directive := cmd.ValidArgsFunction(cmd, []string{"alpha"}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(got) != 0 {
		t.Fatalf("completion = %v, want empty", got)
	}
}
