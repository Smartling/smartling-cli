package format

import (
	"github.com/Smartling/smartling-cli/services/helpers/config"

	"github.com/Smartling/api-sdk-go"
)

func ExecuteFileFormat(
	config config.Config,
	file smartling.File,
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
