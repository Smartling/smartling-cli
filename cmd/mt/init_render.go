package mt

import (
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	output "github.com/Smartling/smartling-cli/output/mt"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"github.com/spf13/cobra"
)

// InitRender for mt subcommands
func InitRender(cmd *cobra.Command, fileConfigMTFileFormat *string, files []string) output.Renderer {
	const outputTemplateFlag = "format"
	outFormat, err := cmd.Parent().PersistentFlags().GetString("output")
	if err != nil {
		output.RenderAndExitIfErr(clierror.UIError{
			Operation:   "get output",
			Err:         err,
			Description: "unable to get output param",
		})
	}
	outTemplate := resolve.FallbackString(cmd.Flags().Lookup(outputTemplateFlag), resolve.StringParam{
		FlagName: outputTemplateFlag,
		Config:   fileConfigMTFileFormat,
	})

	var render output.Renderer = &output.Static{}
	outMode, err := cmd.Parent().PersistentFlags().GetString("output-mode")
	if err != nil {
		output.RenderAndExitIfErr(clierror.UIError{
			Operation:   "get output mode",
			Err:         err,
			Description: "unable to get output mode param",
		})
	}
	if outMode == "dynamic" {
		render = &output.Dynamic{}
	}

	var dataProvider output.TranslateDataProvider
	render.Init(dataProvider, files, outFormat, outTemplate)
	return render
}
