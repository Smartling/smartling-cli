package resolve

import (
	"github.com/Smartling/smartling-cli/output"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/spf13/cobra"
)

// OutputParams resolve OutputParams for subcommands
func OutputParams(cmd *cobra.Command, fileConfigMTFileFormat *string) (output.Params, error) {
	const outputTemplateFlag = "format"
	format, err := cmd.Parent().PersistentFlags().GetString("output")
	if err != nil {
		return output.Params{}, clierror.UIError{
			Operation:   "get output",
			Err:         err,
			Description: "unable to get output param",
		}
	}
	template := FallbackString(cmd.Flags().Lookup(outputTemplateFlag), StringParam{
		FlagName: outputTemplateFlag,
		Config:   fileConfigMTFileFormat,
	})

	mode, err := cmd.Parent().PersistentFlags().GetString("output-mode")
	if err != nil {
		return output.Params{}, clierror.UIError{
			Operation:   "get output mode",
			Err:         err,
			Description: "unable to get output mode param",
		}
	}
	return output.Params{
		Mode:     mode,
		Format:   format,
		Template: template,
	}, nil
}
