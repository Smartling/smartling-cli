package mt

import (
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/spf13/cobra"
)

// InitRender inits render for mt subcommands
func InitRender(cmd *cobra.Command, fileConfigMTFileFormat *string, files []string) (output.Renderer, error) {
	const outputTemplateFlag = "format"
	outFormat, err := cmd.Parent().PersistentFlags().GetString("output")
	if err != nil {
		return nil, clierror.UIError{
			Operation:   "get output",
			Err:         err,
			Description: "unable to get output param",
		}
	}
	outTemplate := resolve.FallbackString(cmd.Flags().Lookup(outputTemplateFlag), resolve.StringParam{
		FlagName: outputTemplateFlag,
		Config:   fileConfigMTFileFormat,
	})

	outMode, err := cmd.Parent().PersistentFlags().GetString("output-mode")
	if err != nil {
		return nil, clierror.UIError{
			Operation:   "get output mode",
			Err:         err,
			Description: "unable to get output mode param",
		}
	}
	var render output.Renderer = &output.Static{}
	if outMode == "dynamic" {
		render = &output.Dynamic{}
	}

	var dataProvider output.TranslateDataProvider
	render.Init(dataProvider, files, outFormat, outTemplate)
	return render, nil
}
