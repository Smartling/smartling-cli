package mt

import (
	"os"

	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	"github.com/Smartling/smartling-cli/services/helpers/config"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// FileConfig defines MT file config
type FileConfig struct {
	MT struct {
		DefaultSourceLocale  *string           `yaml:"default_source_locale,omitzero"`
		DefaultTargetLocales []string          `yaml:"default_target_locales,omitzero"`
		InputDirectory       *string           `yaml:"input_directory,omitzero"`
		OutputDirectory      *string           `yaml:"output_directory,omitzero"`
		FileFormat           *string           `yaml:"file_format,omitzero"`
		Directives           map[string]string `yaml:"directives,omitzero"`
		PollInterval         *int              `yaml:"poll_interval,omitzero"`
		Timeout              *int              `yaml:"timeout,omitzero"`
	} `yaml:"mt,omitzero"`
	Files map[string]FileConfigMT `yaml:"files"`
}

// FileConfigMT defines file config
type FileConfigMT struct {
	MT struct {
		Type       string            `yaml:"type,omitzero"`
		Directives map[string]string `yaml:"directives,omitzero,flow"`
	} `yaml:"mt,omitzero"`
}

// BindFileConfig binds file config
func BindFileConfig(cmd *cobra.Command) (FileConfig, error) {
	dir := resolve.ConfigDirectory(cmd)
	filename := resolve.ConfigFile(cmd)
	path, err := config.GetPath(dir, filename, false)
	if err != nil {
		return FileConfig{}, err
	}
	var config FileConfig
	data, err := os.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return config, err
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, nil
}
