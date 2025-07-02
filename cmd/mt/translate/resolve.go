package translate

import (
	"os"
	"strings"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	"github.com/Smartling/smartling-cli/services/helpers"
	"github.com/Smartling/smartling-cli/services/helpers/env"

	"github.com/spf13/cobra"
)

func resolveTargetLocale(cmd *cobra.Command, fileConfig mtcmd.FileConfig) []string {
	if cmd.Flags().Changed(targetLocaleFlag) {
		return targetLocales
	}
	envVarName := env.VarNameFromCLIFlagName(targetLocaleFlag)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return strings.Split(val, ",")
	}
	if len(fileConfig.MT.DefaultTargetLocales) > 0 {
		return fileConfig.MT.DefaultTargetLocales
	}
	return targetLocales
}

func resolveDirectives(cmd *cobra.Command, fileConfig mtcmd.FileConfig) (map[string]string, error) {
	if cmd.Flags().Changed(directiveFlag) {
		return helpers.MKeyValueToMap(directive)
	}
	envVarName := env.VarNameFromCLIFlagName(directiveFlag)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return helpers.MKeyValueToMap(strings.Split(val, ","))
	}
	if len(fileConfig.MT.Directives) > 0 {
		return fileConfig.MT.Directives, nil
	}
	return helpers.MKeyValueToMap(directive)
}
