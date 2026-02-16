package domain

const (
	SplitHorizontal = "-h"
	SplitVertical   = "-v"
)

// ValidateSplitDirection ensures only tmux-supported split directions are used.
func ValidateSplitDirection(direction string) error {
	switch direction {
	case SplitHorizontal, SplitVertical:
		return nil
	default:
		return ErrInvalidSplitDirection(direction)
	}
}
