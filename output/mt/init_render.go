package mt

import "github.com/Smartling/smartling-cli/output"

// InitRender inits render for mt subcommands
func InitRender(outputParams output.Params, dataProvider TableDataProvider, files []string, targetLocalesQnt uint8) Renderer {
	var render Renderer = &Static{}
	if outputParams.Mode == "dynamic" {
		render = &Dynamic{}
	}
	render.Init(dataProvider, files, targetLocalesQnt, outputParams.Format, outputParams.Template)
	return render
}
