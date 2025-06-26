package translate

import (
	"os"
	"strings"

	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	"github.com/Smartling/smartling-cli/services/helpers"
	"github.com/Smartling/smartling-cli/services/helpers/env"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

func resolveSourceLocale(cmd *cobra.Command, fileConfig mtcmd.FileConfig) string {
	if cmd.Flags().Changed(sourceLocaleFlag) {
		return sourceLocale
	}
	envVarName := env.VarNameFromCLIFlagName(sourceLocaleFlag)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	if fileConfig.MT.DefaultSourceLocale != nil {
		return *fileConfig.MT.DefaultSourceLocale
	}
	return sourceLocale
}

func resolveDetectLanguage(cmd *cobra.Command) bool {
	if cmd.Flags().Changed(detectLanguageFlag) {
		return detectLanguage
	}
	envVarName := env.VarNameFromCLIFlagName(detectLanguageFlag)
	if _, isSet := os.LookupEnv(envVarName); isSet {
		return true
	}
	return detectLanguage
}

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

func resolveOutputDirectory(cmd *cobra.Command, fileConfig mtcmd.FileConfig) string {
	if cmd.Flags().Changed(outputDirectoryFlag) {
		return outputDirectory
	}
	envVarName := env.VarNameFromCLIFlagName(outputDirectoryFlag)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	if fileConfig.MT.OutputDirectory != nil {
		return *fileConfig.MT.OutputDirectory
	}
	return outputDirectory
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

func resolveProgress(cmd *cobra.Command) bool {
	if cmd.Flags().Changed(progressFlag) {
		return progress
	}
	envVarName := env.VarNameFromCLIFlagName(progressFlag)
	if _, isSet := os.LookupEnv(envVarName); isSet {
		return true
	}
	return progress
}

func resolveOverrideFileType(cmd *cobra.Command) string {
	if cmd.Flags().Changed(overrideFileTypeFlag) {
		return overrideFileType
	}
	envVarName := env.VarNameFromCLIFlagName(overrideFileTypeFlag)
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	return overrideFileType
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
