package mt

import "github.com/Smartling/smartling-cli/services/mt"

const DefaultTranslateTemplate = "{{name .File}}_{{.Locale}}{{ext .File}}"

// RenderTranslate renders MT translate
func RenderTranslate(output []mt.TranslateOutput, outputFormat, outputTemplate string) error {
	return nil
}
