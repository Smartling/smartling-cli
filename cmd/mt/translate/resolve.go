package translate

import (
	"fmt"
	"os"
	"strings"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	output "github.com/Smartling/smartling-cli/output/mt"
	"github.com/Smartling/smartling-cli/services/helpers"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/env"
	srv "github.com/Smartling/smartling-cli/services/mt"

	api "github.com/Smartling/api-sdk-go/api/mt"
	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, fileConfig mtcmd.FileConfig) (srv.TranslateParams, error) {
	cnf, err := rootcmd.Config()
	if err != nil {
		output.RenderAndExitIfErr(clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		})
	}

	sourceLocaleParam := resolve.FallbackString(cmd.Flags().Lookup(sourceLocaleFlag), resolve.StringParam{
		FlagName: sourceLocaleFlag,
		Config:   fileConfig.MT.DefaultSourceLocale,
	})
	detectLanguageParam := resolve.FallbackBool(cmd.Flags().Lookup(detectLanguageFlag), resolve.BoolParam{
		FlagName: detectLanguageFlag,
	})
	inputDirectoryParam := resolve.FallbackString(cmd.Flags().Lookup(inputDirectoryFlag), resolve.StringParam{
		FlagName: inputDirectoryFlag,
		Config:   fileConfig.MT.InputDirectory,
	})
	outputDirectoryParam := resolve.FallbackString(cmd.Flags().Lookup(outputDirectoryFlag), resolve.StringParam{
		FlagName: outputDirectoryFlag,
		Config:   fileConfig.MT.OutputDirectory,
	})
	progressParam := resolve.FallbackBool(cmd.Flags().Lookup(progressFlag), resolve.BoolParam{
		FlagName: progressFlag,
	})
	overrideFileTypeParam := resolve.FallbackString(cmd.Flags().Lookup(overrideFileTypeFlag), resolve.StringParam{
		FlagName: overrideFileTypeFlag,
	})

	var accountIDConfig *string
	if cnf.AccountID != "" {
		accountIDConfig = &cnf.AccountID
	}
	accountUIDParam := resolve.FallbackString(cmd.Root().PersistentFlags().Lookup("account"), resolve.StringParam{
		FlagName: "account",
		Config:   accountIDConfig,
	})
	params := srv.TranslateParams{
		SourceLocale:     sourceLocaleParam,
		DetectLanguage:   detectLanguageParam,
		TargetLocales:    resolveTargetLocale(cmd, fileConfig),
		InputDirectory:   inputDirectoryParam,
		OutputDirectory:  outputDirectoryParam,
		Progress:         progressParam,
		OverrideFileType: overrideFileTypeParam,
		AccountUID:       api.AccountUID(accountUIDParam),
	}
	params.Directives, err = resolveDirectives(cmd, fileConfig)
	if err != nil {
		return srv.TranslateParams{}, fmt.Errorf("unable to resolve directives: %w", err)
	}

	return params, nil
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
