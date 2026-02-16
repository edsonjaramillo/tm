package domain

import "testing"

func TestIs(t *testing.T) {
	t.Parallel()

	err := ErrNotInTmuxSession()
	if !Is(err, CodeNotInTmuxSession) {
		t.Fatal("Is(ErrNotInTmuxSession, CodeNotInTmuxSession) = false, want true")
	}
	if Is(err, CodeNotGitRepository) {
		t.Fatal("Is(ErrNotInTmuxSession, CodeNotGitRepository) = true, want false")
	}
}
