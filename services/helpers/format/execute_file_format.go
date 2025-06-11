package format

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	sdk "github.com/Smartling/api-sdk-go"
)

// ExecuteFileFormat executes a file format using the provided configuration and data.
// It returns the rendered string, and an error if any.
func ExecuteFileFormat(
	config config.Config,
	file sdk.File,
	fallback string,
	getter func(config config.FileConfig) string,
	data interface{},
) (string, error) {
	local, err := config.GetFileConfig(file.FileURI)
	if err != nil {
		return "", err
	}

	template := getter(local)

	if template == "" {
		template = fallback
	}

	format, err := Compile(template)
	if err != nil {
		return "", err
	}

	result, err := format.Execute(data)
	if err != nil {
		return "", err
	}

	return result, nil
}
