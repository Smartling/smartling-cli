package format

import (
	"bytes"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"text/template"
)

// Format is format for rendering templates.
type Format struct {
	*template.Template

	Source string
}

// Execute executes the format template with the provided data.
// It returns the rendered string, and an error if any.
func (format *Format) Execute(data any) (string, error) {
	buffer := &bytes.Buffer{}
	if err := format.Template.Execute(buffer, data); err != nil {
		return "", ExecutionError{
			Cause:  err,
			Format: format.Source,
			Data:   data,
		}
	}

	return buffer.String(), nil
}

var (
	// UsePullFormat returns the format for pull files.
	UsePullFormat = func(config config.FileConfig) string {
		return config.Pull.Format
	}
)
