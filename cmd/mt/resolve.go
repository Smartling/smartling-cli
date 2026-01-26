package mt

import (
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/env"

	"github.com/spf13/cobra"
)

func resolveConfigDirectory(cmd *cobra.Command) string {
	flag := cmd.Root().PersistentFlags().Lookup("operation-directory")
	if flag == nil {
		return ""
	}
	if flag.Changed {
		return flag.Value.String()
	}
	envVarName := env.VarNameFromCLIFlagName("operation-directory")
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	return flag.DefValue
}

func resolveConfigFile(cmd *cobra.Command) string {
	flag := cmd.Root().PersistentFlags().Lookup("config")
	if flag == nil {
		return ""
	}
	if flag.Changed {
		return flag.Value.String()
	}
	envVarName := env.VarNameFromCLIFlagName("config")
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	return flag.DefValue
}
