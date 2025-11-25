package translate

import (
	"fmt"
	"os"
	"strings"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	mtcmd "github.com/Smartling/smartling-cli/cmd/mt"
	"github.com/Smartling/smartling-cli/services/helpers"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/env"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"

	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, fileConfig mtcmd.FileConfig) (srv.TranslateParams, error) {
	rlog.Debugf("resolving params")
	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.TranslateParams{}, clierror.UIError{
			Operation:   "config",
			Err:         err,
			Description: "failed to read config",
		}
	}

	sourceLocaleParam := resolve.FallbackString(cmd.Flags().Lookup(sourceLocaleFlag), resolve.StringParam{
		FlagName: sourceLocaleFlag,
		Config:   fileConfig.MT.DefaultSourceLocale,
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

	accountUID, err := resolve.FallbackAccount(cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
	if err != nil {
		return srv.TranslateParams{}, err
	}
	params := srv.TranslateParams{
		SourceLocale:     sourceLocaleParam,
		TargetLocales:    resolveTargetLocale(cmd, fileConfig),
		InputDirectory:   inputDirectoryParam,
		OutputDirectory:  outputDirectoryParam,
		Progress:         progressParam,
		OverrideFileType: overrideFileTypeParam,
		AccountUID:       accountUID,
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
