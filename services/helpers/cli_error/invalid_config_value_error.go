package clierror

import (
	"fmt"
)

type InvalidConfigValueError struct {
	ValueName   string
	Description string
}

func (err InvalidConfigValueError) Error() string {
	return NewError(
		fmt.Errorf(`"%s" is specified but invalid`, err.ValueName),
		`"%s" %s.`,
		err.ValueName,
		err.Description,
	).Error()
}
