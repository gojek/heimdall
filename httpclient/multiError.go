package httpclient

import "strings"

// MultiError is a container for a list of errors.
type MultiError struct {
	errors []error
}

// HasError checks if MultiError has any error.
func (m MultiError) HasError() error {
	if len(m.errors) > 0 {
		return m
	}
	return nil
}

// MultiError implements error interface.
func (m MultiError) Error() string {
	formattedError := make([]string, len(m.errors))

	for i, err := range m.errors {
		formattedError[i] = err.Error()
	}

	return strings.Join(formattedError, ", ")
}

// ErrorList returns a list of all errors in the MultiError with the first one at the front.
func (m MultiError) ErrorList() []error {
	return m.errors
}

// Cause returns the first error in the MultiError or nil if there are none
func (m MultiError) Cause() error {
	if len(m.errors) == 0 {
		return nil
	}
	return m.errors[0]
}

// Unwrap returns the first error in the MultiError or nil if there are none
func (m MultiError) Unwrap() error {
	return m.Cause()
}
