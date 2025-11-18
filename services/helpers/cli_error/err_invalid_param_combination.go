package clierror

import (
	"fmt"
	"strings"
)

// ErrConflictingParams creates a new conflicting param combination error
func ErrConflictingParams(params []string) error {
	cp := append([]string(nil), params...)
	return &errConflictingParams{
		params: cp,
	}
}

type errConflictingParams struct {
	params []string
}

// Error returns the string representation of the error
func (e *errConflictingParams) Error() string {
	return fmt.Sprintf("conflicting parameters: %s", strings.Join(e.params, ", "))
}
