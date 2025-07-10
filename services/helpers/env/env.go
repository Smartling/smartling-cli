package env

import "strings"

const prefix = "SMARTLING_CLI_"

// VarNameFromCLIFlagName returns environment variable name for given CLI flag name
func VarNameFromCLIFlagName(cliFlagName string) string {
	return prefix + strings.ToUpper(cliFlagName)
}
