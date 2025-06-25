package detect

import (
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
	if cmd.Root().Flags().Changed("account") {
		val, err := cmd.Root().PersistentFlags().GetString("account")
		if err != nil {
			return "", err
		}
		return api.AccountUID(val), nil
	}
	envVarName := env.VarNameFromCLIFlagName("account")
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return api.AccountUID(val), nil
	}
	if cnfAccountID != "" {
		return api.AccountUID(cnfAccountID), nil
	}
	val, err := cmd.Root().PersistentFlags().GetString("account")
	if err != nil {
		return "", err
	}
	return api.AccountUID(val), nil
}

func resolveProjectID(cmd *cobra.Command, cnfProjectID string) (string, error) {
	if cmd.Root().Flags().Changed("project") {
		val, err := cmd.Root().PersistentFlags().GetString("project")
		if err != nil {
			return "", err
		}
		return val, nil
	}
	envVarName := env.VarNameFromCLIFlagName("project")
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val, nil
	}
	if cnfProjectID != "" {
		return cnfProjectID, nil
	}
	val, err := cmd.Root().PersistentFlags().GetString("project")
	if err != nil {
		return "", err
	}
	return val, nil
}
