package clierror

import (
	"fmt"
	"strings"
)

// ErrIncompatibleParams creates a new incompatible params error
func ErrIncompatibleParams(param string, incompatibleWith []string) error {
	cp := append([]string(nil), incompatibleWith...)
	return &errIncompatibleParams{
		param:            param,
		incompatibleWith: cp,
	}
}

type errIncompatibleParams struct {
	param            string
	incompatibleWith []string
}

// Error returns the string representation of the error
func (e *errIncompatibleParams) Error() string {
	return fmt.Sprintf("parameter %s is incompatible with: %s", e.param, strings.Join(e.incompatibleWith, ", "))
}
