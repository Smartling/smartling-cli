package resolve

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Smartling/api-sdk-go/helpers/uid"
	"github.com/Smartling/smartling-cli/services/helpers/env"
	"github.com/spf13/cobra"

	"github.com/spf13/pflag"
)

// StringParam defines resolve string param
type StringParam struct {
	FlagName string
	Config   *string
}

// IntParam defines resolve int param
type IntParam struct {
	FlagName string
	Config   *int
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

// FallbackAccount resolves the account UID from flag or config.
func FallbackAccount(flag *pflag.Flag, accountIDConfig string) (uid.AccountUID, error) {
	var config *string
	if accountIDConfig != "" {
		config = &accountIDConfig
	}
	accountUIDParam := FallbackString(flag, StringParam{
		FlagName: "account",
		Config:   config,
	})
	accountUID := uid.AccountUID(accountUIDParam)
	if err := accountUID.Validate(); err != nil {
		return "", err
	}
	return accountUID, nil
}

// FallbackStringArray resolves a []string from flag → env (comma-separated) → config.
func FallbackStringArray(cmd *cobra.Command, flagName string, configVal []string) []string {
	flag := cmd.Flags().Lookup(flagName)
	if flag != nil && flag.Changed {
		vals, _ := cmd.Flags().GetStringArray(flagName)
		return vals
	}
	envVarName := env.VarNameFromCLIFlagName(flagName)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return strings.Split(val, ",")
	}
	if len(configVal) > 0 {
		return configVal
	}
	return nil
}

// FallbackInt resolve int value from hierarchy of fallbacks
func FallbackInt(flag *pflag.Flag, param IntParam) (int, error) {
	// return flag value if it was changed
	if flag != nil && flag.Changed {
		v, err := strconv.Atoi(flag.Value.String())
		if err != nil {
			return 0, fmt.Errorf("invalid --%s (int): %w", param.FlagName, err)
		}
		return v, nil
	}
	// return env value if it is available
	envVarName := env.VarNameFromCLIFlagName(param.FlagName)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		v, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("invalid %s env var (int): %w", envVarName, err)
		}
		return v, nil
	}
	// return config value if it is available
	if param.Config != nil {
		return *param.Config, nil
	}
	// return zero value if no flag
	if flag == nil {
		return 0, nil
	}
	// return default flag value
	v, err := strconv.Atoi(flag.DefValue)
	if err != nil {
		return 0, fmt.Errorf("invalid default value for --%s (int): %w", param.FlagName, err)
	}
	return v, nil
}

// FallbackDate resolves a time.Time from flag → env → config (all RFC3339 strings).
func FallbackDate(cmd *cobra.Command, flagName string, configVal string) (time.Time, error) {
	raw := ""
	flag := cmd.Flags().Lookup(flagName)
	if flag != nil && flag.Changed {
		raw = flag.Value.String()
	} else {
		envVarName := env.VarNameFromCLIFlagName(flagName)
		if val, isSet := os.LookupEnv(envVarName); isSet {
			raw = val
		} else {
			raw = configVal
		}
	}
	if raw == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid --%s (RFC3339): %w", flagName, err)
	}
	return t, nil
}
