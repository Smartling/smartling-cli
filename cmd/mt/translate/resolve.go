package translate

import (
	"os"
	"strings"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	"github.com/Smartling/smartling-cli/services/helpers"

	sdk "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

func resolveSourceLocale(cmd *cobra.Command, fileConfig mtcmd.FileConfig) string {
	if cmd.Flags().Changed(sourceLocaleFlag) {
		return sourceLocale
	}
	if val, isSet := os.LookupEnv(sourceLocaleFlag); isSet {
		return val
	}
	if fileConfig.MT.DefaultSourceLocale != nil {
		return *fileConfig.MT.DefaultSourceLocale
	}
	return sourceLocale
}

func resolveDetectLanguage(cmd *cobra.Command) string {
	if cmd.Flags().Changed(detectLanguageFlag) {
		return detectLanguage
	}
	if val, isSet := os.LookupEnv(detectLanguageFlag); isSet {
		return val
	}
	return detectLanguage
}

func resolveTargetLocale(cmd *cobra.Command, fileConfig mtcmd.FileConfig) []string {
	if cmd.Flags().Changed(targetLocaleFlag) {
		return targetLocale
	}
	if val, isSet := os.LookupEnv(targetLocaleFlag); isSet {
		return strings.Split(val, ",")
	}
	if len(fileConfig.MT.DefaultTargetLocales) > 0 {
		return fileConfig.MT.DefaultTargetLocales
	}
	return targetLocale
}

func resolveDirectory(cmd *cobra.Command, fileConfig mtcmd.FileConfig) string {
	if cmd.Flags().Changed(directoryFlag) {
		return directory
	}
	if val, isSet := os.LookupEnv(directoryFlag); isSet {
		return val
	}
	if fileConfig.MT.OutputDirectory != nil {
		return *fileConfig.MT.OutputDirectory
	}
	return directory
}

func resolveDirectives(cmd *cobra.Command, fileConfig mtcmd.FileConfig) (map[string]string, error) {
	if cmd.Flags().Changed(directiveFlag) {
		return helpers.MKeyValueToMap(directive)
	}
	if val, isSet := os.LookupEnv(directiveFlag); isSet {
		return helpers.MKeyValueToMap(strings.Split(val, ","))
	}
	if len(fileConfig.MT.Directives) > 0 {
		return fileConfig.MT.Directives, nil
	}
	return helpers.MKeyValueToMap(directive)
}

func resolveProgress(cmd *cobra.Command) bool {
	if cmd.Flags().Changed(progressFlag) {
		return progress
	}
	if _, isSet := os.LookupEnv(progressFlag); isSet {
		return true
	}
	return progress
}

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
