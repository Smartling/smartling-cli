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
		DefaultSourceLocale  *string           `yaml:"default_source_locale"`
		DefaultTargetLocales []string          `yaml:"default_target_locales"`
		OutputDirectory      *string           `yaml:"output_directory"`
		FileFormat           *string           `yaml:"file_format"`
		Directives           map[string]string `yaml:"directives"`
		PollInterval         *int              `yaml:"poll_interval"`
		Timeout              *int              `yaml:"timeout"`
	} `yaml:"mt"`
	Files map[string]MTFileConfig `yaml:"files"`
}

type MTFileConfig struct {
	MT struct {
		Type       string            `yaml:"type,omitempty"`
		Directives map[string]string `yaml:"directives,omitempty,flow"`
	} `yaml:"mt,omitempty"`
}

func BindFileConfig(cmd *cobra.Command) (FileConfig, error) {
	dir, err := resolveConfigDirectory(cmd)
	if err != nil {
		return FileConfig{}, err
	}
	filename, err := resolveConfigFile(cmd)
	if err != nil {
		return FileConfig{}, err
	}
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

func resolveConfigDirectory(cmd *cobra.Command) (string, error) {
	if cmd.Root().Flags().Changed("directory") {
		val, err := cmd.Root().PersistentFlags().GetString("directory")
		if err != nil {
			return "", err
		}
		return val, nil
	}
	if val, isSet := os.LookupEnv("directory"); isSet {
		return val, nil
	}
	val, err := cmd.Root().PersistentFlags().GetString("directory")
	if err != nil {
		return "", err
	}
	return val, nil
}

func resolveConfigFile(cmd *cobra.Command) (string, error) {
	if cmd.Root().Flags().Changed("config") {
		val, err := cmd.Root().PersistentFlags().GetString("config")
		if err != nil {
			return "", err
		}
		return val, nil
	}
	if val, isSet := os.LookupEnv("config"); isSet {
		return val, nil
	}
	val, err := cmd.Root().PersistentFlags().GetString("config")
	if err != nil {
		return "", err
	}
	return val, nil
}
