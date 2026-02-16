package domain

import "testing"

func TestValidateSplitDirection(t *testing.T) {
	t.Parallel()

	cases := []string{SplitHorizontal, SplitVertical}
	for _, direction := range cases {
		direction := direction
		t.Run(direction, func(t *testing.T) {
			t.Parallel()
			if err := ValidateSplitDirection(direction); err != nil {
				t.Fatalf("ValidateSplitDirection(%q) error = %v, want nil", direction, err)
			}
		})
	}
}

func TestValidateSplitDirectionInvalid(t *testing.T) {
	t.Parallel()

	err := ValidateSplitDirection("-x")
	if err == nil {
		t.Fatal("ValidateSplitDirection(-x) error = nil, want non-nil")
	}
	if !Is(err, CodeInvalidSplit) {
		t.Fatalf("Is(err, CodeInvalidSplit) = false, want true (err=%v)", err)
	}
	if got, want := err.Error(), "Invalid direction for split window. Use -h or -v."; got != want {
		t.Fatalf("error message = %q, want %q", got, want)
	}
}
