package mt

// InitRender inits render for mt subcommands
func InitRender(outputParams OutputParams, dataProvider TableDataProvider, files []string) Renderer {
	var render Renderer = &Static{}
	if outputParams.Mode == "dynamic" {
		render = &Dynamic{}
	}
	render.Init(dataProvider, files, outputParams.Format, outputParams.Template)
	return render
}
