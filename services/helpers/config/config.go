package config

import (
	"os"

	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"

	"dario.cat/mergo"
	"github.com/gobwas/glob"
	"github.com/goccy/go-yaml"
	
)

type FileConfig struct {
	Pull struct {
		Format string `yaml:"format,omitempty"`
	} `yaml:"pull,omitempty"`

	Push struct {
		Type       string            `yaml:"type,omitempty"`
		Directives map[string]string `yaml:"directives,omitempty,flow"`
	} `yaml:"push,omitempty"`
}

type Config struct {
	UserID    string `yaml:"user_id"`
	Secret    string `yaml:"secret"`
	AccountID string `yaml:"account_id"`
	ProjectID string `yaml:"project_id,omitempty"`
	Threads   uint32 `yaml:"threads"`

	Files map[string]FileConfig `yaml:"files"`

	Proxy string `yaml:"proxy,omitempty"`

	Path string `yaml:"-"`
}

func (config *Config) GetFileConfig(path string) (FileConfig, error) {
	var (
		match FileConfig
		found bool
	)

	for key, candidate := range config.Files {
		pattern, err := glob.Compile(key, '/')
		if err != nil {
			return FileConfig{}, clierror.NewError(
				hierr.Errorf(
					err,
					`unable to compile pattern from config file (key "%s")`,
					key,
				),

				`File match pattern is malformed. Check out help for more `+
					`information on globbing patterns.`,
			)
		}

		if pattern.Match(path) {
			match = candidate
			found = true
		}
	}

	defaults := config.Files["default"]

	if !found {
		return defaults, nil
	}

	err := mergo.Merge(&match, defaults)
	if err != nil {
		return FileConfig{}, clierror.NewError(
			hierr.Errorf(err, "unable to merge file config options"),
			`It's internal error. Consider reporting bug.`,
		)
	}

	return match, nil
}

func LoadConfigFromFile(filename string) (Config, error) {
	config := Config{
		Path: filename,
	}

	data, err := os.ReadFile(filename)
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
