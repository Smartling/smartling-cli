package detect

import (
	"errors"
	"os"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	"github.com/Smartling/smartling-cli/services/helpers/env"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

func resolveFileType(cmd *cobra.Command) string {
	if cmd.Flags().Changed(fileTypeFlag) {
		return fileType
	}
	envVarName := env.VarNameFromCLIFlagName(fileTypeFlag)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	return fileType
}

func resolveOutputTemplate(cmd *cobra.Command, fileConfig mtcmd.FileConfig) string {
	if cmd.Flags().Changed(outputTemplateFlag) {
		return outputTemplate
	}
	envVarName := env.VarNameFromCLIFlagName(outputTemplateFlag)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	if fileConfig.MT.FileFormat != nil {
		return *fileConfig.MT.FileFormat
	}
	return outputTemplate
}

func resolveAccountUID(cmd *cobra.Command, cnfAccountID string) (api.AccountUID, error) {
	flagName := "account"
	flag := cmd.Root().PersistentFlags().Lookup(flagName)
	if flag == nil {
		return "", errors.New(flagName + " flag is not defined")
	}
	if flag.Changed {
		return api.AccountUID(flag.Value.String()), nil
	}
	envVarName := env.VarNameFromCLIFlagName(flagName)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return api.AccountUID(val), nil
	}
	if cnfAccountID != "" {
		return api.AccountUID(cnfAccountID), nil
	}
	return api.AccountUID(flag.DefValue), nil
}

func resolveProjectID(cmd *cobra.Command, cnfProjectID string) (string, error) {
	flagName := "project"
	flag := cmd.Root().PersistentFlags().Lookup(flagName)
	if flag == nil {
		return "", errors.New(flagName + " flag is not defined")
	}
	if flag.Changed {
		return flag.Value.String(), nil
	}
	envVarName := env.VarNameFromCLIFlagName(flagName)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val, nil
	}
	if cnfProjectID != "" {
		return cnfProjectID, nil
	}
	return flag.DefValue, nil
}
