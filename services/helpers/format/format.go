package format

import (
	"bytes"
	"github.com/Smartling/smartling-cli/services/helpers/config"
	"text/template"
)

type Format struct {
	*template.Template

	Source string
}

func (format *Format) Execute(data interface{}) (string, error) {
	buffer := &bytes.Buffer{}

	err := format.Template.Execute(buffer, data)
	if err != nil {
		return "", ExecutionError{
			Cause:  err,
			Format: format.Source,
			Data:   data,
		}
	}

	return buffer.String(), nil
}

var (
	UsePullFormat = func(config config.FileConfig) string {
		return config.Pull.Format
	}
)
