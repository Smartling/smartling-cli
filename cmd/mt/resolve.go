package mt

import (
	"os"

	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/env"

	"github.com/spf13/cobra"
)

// ResolveOutputParams resolve OutputParams for subcommands
func ResolveOutputParams(cmd *cobra.Command, fileConfigMTFileFormat *string) (output.OutputParams, error) {
	const outputTemplateFlag = "format"
	format, err := cmd.Parent().PersistentFlags().GetString("output")
	if err != nil {
		return output.OutputParams{}, clierror.UIError{
			Operation:   "get output",
			Err:         err,
			Description: "unable to get output param",
		}
	}
	template := resolve.FallbackString(cmd.Flags().Lookup(outputTemplateFlag), resolve.StringParam{
		FlagName: outputTemplateFlag,
		Config:   fileConfigMTFileFormat,
	})

	mode, err := cmd.Parent().PersistentFlags().GetString("output-mode")
	if err != nil {
		return output.OutputParams{}, clierror.UIError{
			Operation:   "get output mode",
			Err:         err,
			Description: "unable to get output mode param",
		}
	}
	return output.OutputParams{
		Mode:     mode,
		Format:   format,
		Template: template,
	}, nil
}

func resolveConfigDirectory(cmd *cobra.Command) string {
	flag := cmd.Root().PersistentFlags().Lookup("operation-directory")
	if flag == nil {
		return ""
	}
	if flag.Changed {
		return flag.Value.String()
	}
	envVarName := env.VarNameFromCLIFlagName("operation-directory")
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	return flag.DefValue
}

func resolveConfigFile(cmd *cobra.Command) string {
	flag := cmd.Root().PersistentFlags().Lookup("config")
	if flag == nil {
		return ""
	}
	if flag.Changed {
		return flag.Value.String()
	}
	envVarName := env.VarNameFromCLIFlagName("config")
	if val, isSet := os.LookupEnv(envVarName); isSet {
		return val
	}
	return flag.DefValue
}
