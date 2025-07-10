package resolve

import (
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/env"

	"github.com/spf13/pflag"
)

// StringParam defines resolve string param
type StringParam struct {
	FlagName string
	Config   *string
}

// BoolParam defines resolve bool param
type BoolParam struct {
	FlagName string
	Config   *bool
}

// FallbackString resolve string value from hierarchy of fallbacks
func FallbackString(flag *pflag.Flag, param StringParam) string {
	// return flag value if it was changed
	if flag != nil && flag.Changed {
		return flag.Value.String()
	}
	// return env value if it is available
	envVarName := env.VarNameFromCLIFlagName(param.FlagName)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	// return config value if it is available
	if param.Config != nil {
		return *param.Config
	}
	// return zero value if no flag
	if flag == nil {
		return ""
	}
	// return default flag value
	return flag.DefValue
}

// FallbackBool resolve bool value from hierarchy of fallbacks
func FallbackBool(flag *pflag.Flag, param BoolParam) bool {
	// return flag value if it was changed
	if flag != nil && flag.Changed {
		return true
	}
	// return env value if it is available
	envVarName := env.VarNameFromCLIFlagName(param.FlagName)
	if _, isSet := os.LookupEnv(envVarName); isSet {
		return true
	}
	// return config value if it is available
	if param.Config != nil {
		return *param.Config
	}
	// return zero value if no flag
	if flag == nil {
		return false
	}
	// return default flag value
	return true
}
