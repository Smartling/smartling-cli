package jobs

import (
	"fmt"

	"github.com/Smartling/smartling-cli/services/jobs"
)

// OutputFormat defines behaviour to format and render data
type OutputFormat interface {
	FormatAndRender(data jobs.ProgressOutput)
}

// GetOutputFormat returns OutputFormat for given string
func GetOutputFormat(outputFormat string) OutputFormat {
	switch outputFormat {
	case "json":
		return JsonOutputFormat{}
	case "simple":
		return SimpleOutputFormat{}
	}
	return SimpleOutputFormat{}
}

// JsonOutputFormat is json output format for rendering data as json
type JsonOutputFormat struct{}

// FormatAndRender marshals table data into JSON and prints it.
func (j JsonOutputFormat) FormatAndRender(data jobs.ProgressOutput) {
	fmt.Println(string(data.Json))
}

// SimpleOutputFormat is simple output format for rendering data using a text template
type SimpleOutputFormat struct{}

// FormatAndRender formats the data rows using the stored template string
// and outputs the rendered result
func (s SimpleOutputFormat) FormatAndRender(data jobs.ProgressOutput) {
	fmt.Println("Total word count: ", data.TotalWordCount)
	fmt.Println("Percent complete: ", data.PercentComplete)
}
