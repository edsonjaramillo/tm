package cli

import "errors"

var errAlreadyReportedFailure = errors.New("command failed with details already reported")

// IsAlreadyReportedFailure reports whether an error is a command failure already shown to the user.
func IsAlreadyReportedFailure(err error) bool {
	return errors.Is(err, errAlreadyReportedFailure)
}
