package files

import (
	"fmt"
	"strings"
)

// ErrInvalidParamCombination creates a new invalid param combination error
func ErrInvalidParamCombination(params []string) error {
	cp := append([]string(nil), params...)
	return &errInvalidParamCombination{
		params: cp,
	}
}

type errInvalidParamCombination struct {
	params []string
}

// Error returns the string representation of the error
func (e *errInvalidParamCombination) Error() string {
	return fmt.Sprintf("invalid params combination: %s", strings.Join(e.params, ", "))
}
