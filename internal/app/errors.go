package app

import (
	"errors"
	"fmt"
)

// Kind identifies application-level error classes for consistent handling.
type Kind string

const (
	KindValidation Kind = "validation"
	KindExecution  Kind = "execution"
)

// Error wraps domain and adapter failures with app-level intent.
type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string {
	switch {
	case e.Message == "" && e.Err != nil:
		return e.Err.Error()
	case e.Message != "" && e.Err == nil:
		return e.Message
	case e.Message == "" && e.Err == nil:
		return ""
	default:
		if e.Message == e.Err.Error() {
			return e.Message
		}
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
}

func (e *Error) Unwrap() error {
	return e.Err
}

// IsKind reports whether err is an app.Error with the provided kind.
func IsKind(err error, kind Kind) bool {
	var appErr *Error
	if !errors.As(err, &appErr) {
		return false
	}
	return appErr.Kind == kind
}

func newValidationError(err error) error {
	if err == nil {
		return nil
	}
	return &Error{
		Kind:    KindValidation,
		Message: err.Error(),
		Err:     err,
	}
}

func newExecutionError(message string, err error) error {
	if err == nil {
		return nil
	}
	return &Error{
		Kind:    KindExecution,
		Message: message,
		Err:     err,
	}
}
