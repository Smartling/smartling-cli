package resolve

import (
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/env"

	"github.com/spf13/cobra"
)

// ConfigDirectory resolves the operation directory.
func ConfigDirectory(cmd *cobra.Command) string {
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

// ConfigFile resolves the config file name.
func ConfigFile(cmd *cobra.Command) string {
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
