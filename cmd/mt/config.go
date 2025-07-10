package mt

import (
	"os"

	"github.com/Smartling/smartling-cli/services/helpers/config"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// FileConfig defines MT file config
type FileConfig struct {
	MT struct {
		DefaultSourceLocale  *string           `yaml:"default_source_locale,omitempty"`
		DefaultTargetLocales []string          `yaml:"default_target_locales,omitempty"`
		InputDirectory       *string           `yaml:"input_directory,omitempty"`
		OutputDirectory      *string           `yaml:"output_directory,omitempty"`
		FileFormat           *string           `yaml:"file_format,omitempty"`
		Directives           map[string]string `yaml:"directives,omitempty"`
		PollInterval         *int              `yaml:"poll_interval,omitempty"`
		Timeout              *int              `yaml:"timeout,omitempty"`
	} `yaml:"mt,omitempty"`
	Files map[string]FileConfigMT `yaml:"files"`
}

// FileConfigMT defines file config
type FileConfigMT struct {
	MT struct {
		Type       string            `yaml:"type,omitempty"`
		Directives map[string]string `yaml:"directives,omitempty,flow"`
	} `yaml:"mt,omitempty"`
}

// BindFileConfig binds file config
func BindFileConfig(cmd *cobra.Command) (FileConfig, error) {
	dir := resolveConfigDirectory(cmd)
	filename := resolveConfigFile(cmd)
	path, err := config.GetPath(dir, filename, false)
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
