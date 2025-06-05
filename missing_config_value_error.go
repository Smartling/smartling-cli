package main

import (
	"fmt"
	"github.com/Smartling/smartling-cli/services/helpers/cli_error"
)

type MissingConfigValueError struct {
	ConfigPath string
	EnvVarName string
	ValueName  string
	OptionName string
	KeyName    string
}

func (err MissingConfigValueError) Error() string {
	return clierror.NewError(
		fmt.Errorf(
			"cannot find mandatory configuration parameter %q",
			err.ValueName,
		),

		"Please, specify either:\n"+
			"- Environment variable $%s;\n"+
			"- Command line option --%s=<%s>;\n"+
			"- Or set %q option in the configuration file:\n\n\t%s\n\t\t%s",
		err.EnvVarName,
		err.OptionName,
		err.KeyName,
		err.KeyName,
		err.ConfigPath,
		fmt.Sprintf(`%s: "PUT_VALUE_HERE"`, err.KeyName),
	).Error()
}
