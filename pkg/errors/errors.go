package errors

import (
	"strings"
)

// MultiError is a collection of errors that implements the error interface.
//
// Example:
//
//	errs := &errors.MultiError{}
//	errs.Append(err1, err2)
//	if err := errs.ErrorOrNil(); err != nil { ... }
type MultiError struct {
	Errors []error
}

// Append adds errors to the collection. Nil errors are ignored.
func (e *MultiError) Append(errs ...error) {
	for _, err := range errs {
		if err != nil {
			e.Errors = append(e.Errors, err)
		}
	}
}

// Error implements the error interface.
//
// Returns:
//
// A concatenated string of all error messages, separated by semicolons.
func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	var sb strings.Builder
	sb.WriteString("multiple errors occurred: ")
	for i, err := range e.Errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Unwrap returns the errors as a slice for errors.Is/As support (Go 1.20+).
func (e *MultiError) Unwrap() []error {
	return e.Errors
}

// HasErrors returns true if there are any errors in the collection.
func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ErrorOrNil returns nil if there are no errors, otherwise returns itself.
func (e *MultiError) ErrorOrNil() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return e
}

// Join is a helper to join multiple errors into a single error.
//
// Example:
//
//	err := errors.Join(err1, err2)
func Join(errs ...error) error {
	e := &MultiError{}
	e.Append(errs...)
	return e.ErrorOrNil()
}
