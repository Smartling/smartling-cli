package mt

import (
	output "github.com/Smartling/smartling-cli/output/mt"
)

// InitRender inits render for mt subcommands
func InitRender(outputParams output.OutputParams, dataProvider output.TableDataProvider, files []string) (output.Renderer, error) {
	var render output.Renderer = &output.Static{}
	if outputParams.Mode == "dynamic" {
		render = &output.Dynamic{}
	}
	render.Init(dataProvider, files, outputParams.Format, outputParams.Template)
	return render, nil
}
