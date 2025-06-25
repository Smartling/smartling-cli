package detect

import (
	"os"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"

	sdk "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

func resolveFileType(cmd *cobra.Command) string {
	if cmd.Flags().Changed(fileTypeFlag) {
		return fileType
	}
	if val, isSet := os.LookupEnv(fileTypeFlag); isSet {
		return val
	}
	return fileType
}

func resolveOutputTemplate(cmd *cobra.Command, fileConfig mtcmd.FileConfig) string {
	if cmd.Flags().Changed(outputTemplateFlag) {
		return outputTemplate
	}
	if val, isSet := os.LookupEnv(outputTemplateFlag); isSet {
		return val
	}
	if fileConfig.MT.FileFormat != nil {
		return *fileConfig.MT.FileFormat
	}
	return outputTemplate
}

func resolveAccountUID(cmd *cobra.Command, cnfAccountID string) (sdk.AccountUID, error) {
	if cmd.Root().Flags().Changed("account") {
		val, err := cmd.Root().PersistentFlags().GetString("account")
		if err != nil {
			return "", err
		}
		return sdk.AccountUID(val), nil
	}
	if val, isSet := os.LookupEnv("account"); isSet {
		return sdk.AccountUID(val), nil
	}
	if cnfAccountID != "" {
		return sdk.AccountUID(cnfAccountID), nil
	}
	val, err := cmd.Root().PersistentFlags().GetString("account")
	if err != nil {
		return "", err
	}
	return sdk.AccountUID(val), nil
}

func resolveProjectID(cmd *cobra.Command, cnfProjectID string) (string, error) {
	if cmd.Root().Flags().Changed("project") {
		val, err := cmd.Root().PersistentFlags().GetString("project")
		if err != nil {
			return "", err
		}
		return val, nil
	}
	if val, isSet := os.LookupEnv("project"); isSet {
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
