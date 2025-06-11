package clierror

import (
	"fmt"
)

// InvalidConfigValueError is an error type that indicates a configuration value
type InvalidConfigValueError struct {
	ValueName   string
	Description string
}

// Error returns string representation.
func (err InvalidConfigValueError) Error() string {
	return NewError(
		fmt.Errorf(`"%s" is specified but invalid`, err.ValueName),
		`"%s" %s.`,
		err.ValueName,
		err.Description,
	).Error()
}
