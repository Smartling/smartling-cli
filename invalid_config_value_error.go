package main

import (
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
)

type InvalidConfigValueError struct {
	ValueName   string
	Description string
}

func (err InvalidConfigValueError) Error() string {
	return clierror.NewError(
		fmt.Errorf(`"%s" is specified but invalid`, err.ValueName),
		`"%s" %s.`,
		err.ValueName,
		err.Description,
	).Error()
}
