package internal

import "strings"

func BuildMultiError(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	if len(errors) == 1 {
		return errors[0]
	}

	return &multiError{errors: errors}
}

type multiError struct {
	errors []error
}

func (m multiError) Error() string {
	var sb strings.Builder
	for i, e := range m.errors {
		sb.WriteString(e.Error())
		if i < len(m.errors)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (m multiError) Unwrap() []error {
	return m.errors
}
